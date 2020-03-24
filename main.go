package main

import (
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

// Magic happens here
func main() {
	flag.Parse()
	var DateTime string

	DateTime = *flagDate

	// Create data concurrency channel
	c := make(chan string)
	// Create state concorrency channel

	// Execute both functions in concurrency channel
	go GetCurrentLunarPhase(c, DateTime)
	go GetCurrentAtenImage(c, DateTime)
	go GetCurrentGaiaImage(c, DateTime)

	for i := 0; i < 7; i++ {
		fmt.Println(<-c)
	}

}
