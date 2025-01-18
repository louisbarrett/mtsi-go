import React, { useEffect, useState, useRef } from 'react';
import { fetchSunData } from '../services/api';
import DatePicker from './DatePicker';
import styles from '../styles/common.module.css';

function Sun() {
  const mounted = useRef(true);
  const [sunData, setSunData] = useState({ imageUrl: '' });
  const [loading, setLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState(() => {
    const date = new Date();
    date.setDate(date.getDate() - 2);
    return date;
  });

  useEffect(() => {
    mounted.current = true;
    
    const loadSunData = async () => {
      if (!mounted.current) return;
      setLoading(true);
      
      try {
        const dateString = selectedDate ? selectedDate.toISOString().split('T')[0] : '';
        const data = await fetchSunData(dateString);
        console.log('Sun Component Data:', {
          hasImageUrl: !!data.imageUrl,
          imageUrlLength: data.imageUrl?.length
        });
        
        if (mounted.current && data.imageUrl) {
          setSunData(data);
        }
      } catch (error) {
        console.error('Error loading Sun data:', error);
      } finally {
        if (mounted.current) {
          setLoading(false);
        }
      }
    };

    loadSunData();

    let interval;
    if (!selectedDate) {
      interval = setInterval(loadSunData, 300000);
    }
    
    return () => {
      mounted.current = false;
      if (interval) clearInterval(interval);
    };
  }, [selectedDate]);

  if (loading) {
    return <div className={styles.loading}>Loading Sun data...</div>;
  }

  return (
    <div>
      <DatePicker value={selectedDate} onChange={setSelectedDate} />
      <div className={styles.card}>
        <div className={styles.imageContainer}>
          {sunData.imageUrl ? (
            <img 
              src={sunData.imageUrl} 
              alt="Solar activity" 
              onError={(e) => {
                console.error('Image failed to load:', e);
                e.target.style.display = 'none';
              }}
            />
          ) : (
            <div className={styles.noImage}>No image available for selected date</div>
          )}
        </div>
        <h2>Solar Activity {selectedDate ? `for ${selectedDate.toLocaleDateString()}` : '(Current)'}</h2>
      </div>
    </div>
  );
}

export default Sun; 