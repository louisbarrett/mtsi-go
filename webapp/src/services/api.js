import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

export async function fetchEarthData(date) {
  const response = await axios.get(`${API_BASE_URL}/earth`, {
    params: { date }
  });
  
  const imageUrl = response.data.imageData 
    ? `data:image/png;base64,${response.data.imageData}`
    : '';
    
  console.log('Earth API Response:', {
    hasImageData: !!response.data.imageData,
    imageUrlLength: imageUrl.length,
    lunarData: response.data.lunarData
  });

  return {
    lunarData: response.data.lunarData,
    imageUrl
  };
}

export async function fetchSunData(date) {
  const response = await axios.get(`${API_BASE_URL}/sun`, {
    params: { date }
  });
  
  const imageUrl = response.data.imageData 
    ? `data:image/png;base64,${response.data.imageData}`
    : '';
    
  console.log('Sun API Response:', {
    hasImageData: !!response.data.imageData,
    imageUrlLength: imageUrl.length
  });

  return {
    imageUrl
  };
} 