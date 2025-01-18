package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
)

const (
	layoutISO = "2006-01-02"
)

var (
	flagDate = flag.String("date", time.Now().Add(time.Hour*-48).Format(layoutISO), "Date to pull data for in 2006-12-24 format")
	flagZone = flag.Int("zone", 6, `
	00 - 03 Australia
	05 - 07 Africa
	08 - 09 South America
	10 - 11 USA
	12      Pacific Ocean
	`)
)

func intToString(intValue int64) string {
	newvalue := strconv.FormatInt(intValue, 10)
	return newvalue

}

// GetCurrentLunarPhase fuck off golang
func GetCurrentLunarPhase(c chan string, date string) {
	UserTime, err := time.Parse(layoutISO, date)
	if err != nil {
		log.Fatal(err)
	}
	moonphase := map[string]string{
		"Dark Moon":       "🌑",
		"New Moon":        "🌑",
		"Waxing Crescent": "🌒",
		"1st Quarter":     "🌓",
		"Full Moon":       "🌕",
		"Waning Crescent": "🌔",
		"3rd Quarter":     "🌗",
		"Waning Gibbous":  "🌖",
		"Waxing Gibbous":  "🌔",
	}

	unixtime := strconv.FormatInt(UserTime.Unix(), 10)
	url := "http://api.farmsense.net/v1/moonphases/?d=" + string(unixtime)
	// Invoke initial web request
	Response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error retrieving content")
	}
	// Convert byte stream to string
	ResponseBody, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		fmt.Println("error decoding contents")
	}
	ResponseString, err := gabs.ParseJSON(ResponseBody)
	elements := ResponseString.Children()
	// fmt.Println(elements[0].String())

	ResponseString, _ = gabs.ParseJSON(elements[0].Bytes())
	if err != nil {
		log.Fatal("nah", err)
	}
	// TargetDate := stringCleaning(ResponseString.Path("TargetDate").String())
	// Moon := ResponseString.Path("Moon").Data()
	Age := ResponseString.Path("Age").String()
	Phase := ResponseString.Path("Phase").Data()
	EmojiPhase := moonphase[Phase.(string)]
	Distance := ResponseString.Path("Distance").String()
	Illumination := ResponseString.Path("Illumination").String()
	// AngularDiameter := stringCleaning(ResponseString.Path("AngularDiameter").String())
	DistanceToSun := ResponseString.Path("DistanceToSun").String()
	// c <- string(Moon.(string))
	c <- string("Days into cycle: " + string(Age))
	c <- string("Current Moon Phase: " + string(Phase.(string)+" "+EmojiPhase))
	c <- string("Distance From Earth: " + string(Distance) + " kilometers")
	c <- string("Lunar Illumination: " + string(Illumination))
	c <- string("Distance from the sun " + string(DistanceToSun) + " miles")
}

// GetCurrentAtenImage for later use
func GetCurrentAtenImage(c chan string, date string) {
	var currentTime time.Time
	var err error
	if date != "" {
		currentTime, err = time.Parse(layoutISO, date)
	} else {
		currentTime = time.Now()
	}
	year := intToString(int64(currentTime.Year()))
	month := intToString(int64(currentTime.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	day := intToString(int64(currentTime.Day()))
	if len(day) == 1 {
		day = "0" + day
	}

	TimeString := year + "-" + month + "-" + day

	solarURL := ("https://api.helioviewer.org/v1/takeScreenshot/?imageScale=2.4204409&layers=[SDO,AIA,AIA,304,1,100]&events=&eventLabels=true&scale=false&date=" + TimeString + "T00:00:00Z" + "&x1=-929.2475775696686&x2=486.70112763033143&y1=-970.7984919973343&y2=486.3069298026657&display=true&watermark=true&events=[CH,all,1]")
	// fmt.Println(solarURL)
	response, err := http.Get(solarURL)
	if err != nil {
		log.Fatal(err)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	atenImagePath := "aten.png"
	atenImage := ioutil.WriteFile(atenImagePath, responseBytes, os.FileMode(0777))
	if atenImage != nil {
	}
	c <- string("Wrote solar image data to " + atenImagePath)
}

// GetCurrentGaiaImage Stuff
func GetCurrentGaiaImage(c chan string, date string) {
	var currentTime time.Time
	currentTime, err := time.Parse(layoutISO, date)
	if err != nil {
		log.Fatal(err)
	}
	if currentTime.After(time.Now().AddDate(0, 0, -2)) {
		currentTime = time.Now().AddDate(0, 0, -2)
	}

	// Latest data is from two days ago
	// currentTime := time.Now()
	day := intToString(int64(currentTime.Day()))
	if len(day) == 1 {
		day = "0" + day
	}
	year := intToString(int64(currentTime.Year()))
	month := intToString(int64(currentTime.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	// image list URL
	URL := "https://epic.gsfc.nasa.gov/api/natural/date/" + string(year) + "-" + string(month) + "-" + string(day) + "?api_key=DEMO_KEY"
	// image data URl
	imageURL := "https://epic.gsfc.nasa.gov/archive/natural/" + string(year) + "/" + string(month) + "/" + string(day) + "/png/"
	// fmt.Println(URL)
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	parsedJSON, err := gabs.ParseJSON(responseBytes)
	if err != nil {
		log.Fatal(err, string(responseBytes))
	}
	if parsedJSON.Index(*flagZone).Search("image").Data() == nil {
		c <- "No image for " + currentTime.String()
		return
	}
	gaiaImagesList := parsedJSON.Index(*flagZone).Search("image").Data().(string)

	// 00 - 03 Australia
	// 05 - 07 Africa
	// 08 - 09 South America
	// 10 USA
	// 12 Pacific Ocean
	// fmt.Println(gaiaImagesList)

	gaiaURL := imageURL + string(gaiaImagesList) + ".png"
	response, err = http.Get(gaiaURL)
	if err != nil {
		log.Fatal(err)
	}
	responseBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	gaiaImagePath := "gaia.png"
	gaiaImage := ioutil.WriteFile(gaiaImagePath, responseBytes, os.FileMode(0777))
	if gaiaImage != nil {
		log.Fatal(err)
	}

	c <- "Wrote earth image data to " + gaiaImagePath

}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Enable CORS middleware
	router.HandleFunc("/api/earth", enableCORS(handleEarthData))
	router.HandleFunc("/api/sun", enableCORS(handleSunData))

	return router
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func handleEarthData(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received Earth data request for date: %s", r.URL.Query().Get("date"))

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	c := make(chan string)
	imageC := make(chan []byte)
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format(layoutISO)
	}

	go GetCurrentLunarPhase(c, date)
	go GetCurrentGaiaImageBytes(imageC, date)

	// Collect lunar data
	var lunarData []string
	for i := 0; i < 5; i++ {
		data := <-c
		lunarData = append(lunarData, data)
	}

	// Get image bytes
	imageBytes := <-imageC

	response := map[string]interface{}{
		"lunarData": lunarData,
		"imageData": imageBytes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSunData(w http.ResponseWriter, r *http.Request) {
	imageC := make(chan []byte)
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format(layoutISO)
	}

	go GetCurrentAtenImageBytes(imageC, date)
	imageBytes := <-imageC

	response := map[string]interface{}{
		"imageData": imageBytes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCurrentGaiaImageBytes(c chan []byte, date string) {
	var currentTime time.Time
	var err error
	if date != "" {
		currentTime, err = time.Parse(layoutISO, date)
	} else {
		currentTime = time.Now()
	}
	year := intToString(int64(currentTime.Year()))
	month := intToString(int64(currentTime.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	day := intToString(int64(currentTime.Day()))
	if len(day) == 1 {
		day = "0" + day
	}

	// image list URL
	URL := "https://epic.gsfc.nasa.gov/api/natural/date/" + string(year) + "-" + string(month) + "-" + string(day) + "?api_key=DEMO_KEY"
	// image data URl
	imageURL := "https://epic.gsfc.nasa.gov/archive/natural/" + string(year) + "/" + string(month) + "/" + string(day) + "/png/"

	response, err := http.Get(URL)
	if err != nil {
		c <- []byte("Error fetching image list")
		return
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c <- []byte("Error reading response")
		return
	}
	parsedJSON, err := gabs.ParseJSON(responseBytes)
	if err != nil {
		c <- []byte("Error parsing JSON")
		return
	}
	if parsedJSON.Index(*flagZone).Search("image").Data() == nil {
		c <- []byte(fmt.Sprintf("No image available for %s", currentTime.Format(layoutISO)))
		return
	}
	gaiaImagesList := parsedJSON.Index(*flagZone).Search("image").Data().(string)

	gaiaURL := imageURL + string(gaiaImagesList) + ".png"
	response, err = http.Get(gaiaURL)
	if err != nil {
		c <- []byte("Error fetching image")
		return
	}
	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c <- []byte("Error reading image data")
		return
	}

	c <- imageBytes
}

func GetCurrentAtenImageBytes(c chan []byte, date string) {
	var currentTime time.Time
	var err error
	if date != "" {
		currentTime, err = time.Parse(layoutISO, date)
	} else {
		currentTime = time.Now()
	}
	year := intToString(int64(currentTime.Year()))
	month := intToString(int64(currentTime.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	day := intToString(int64(currentTime.Day()))
	if len(day) == 1 {
		day = "0" + day
	}

	TimeString := year + "-" + month + "-" + day

	solarURL := ("https://api.helioviewer.org/v1/takeScreenshot/?imageScale=2.4204409&layers=[SDO,AIA,AIA,304,1,100]&events=&eventLabels=true&scale=false&date=" + TimeString + "T00:00:00Z" + "&x1=-929.2475775696686&x2=486.70112763033143&y1=-970.7984919973343&y2=486.3069298026657&display=true&watermark=true&events=[CH,all,1]")

	response, err := http.Get(solarURL)
	if err != nil {
		c <- []byte("Error fetching solar image")
		return
	}
	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c <- []byte("Error reading solar image data")
		return
	}

	c <- imageBytes
}

// Magic happens here
func main() {
	flag.Parse()

	// Setup router with our API endpoints
	router := setupRouter()

	// Create a file server for serving the images
	fs := http.FileServer(http.Dir("."))
	router.Handle("/aten.png", fs)
	router.Handle("/gaia.png", fs)

	// Start the server
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
