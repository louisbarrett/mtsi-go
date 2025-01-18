import React, { useEffect, useState, useRef } from 'react';
import { fetchEarthData } from '../services/api';
import { useDateSelection } from '../hooks/useDateSelection';
import DatePicker from './DatePicker';
import styles from '../styles/common.module.css';

function Earth() {
  const mounted = useRef(true);
  const { selectedDate, setSelectedDate, dateString } = useDateSelection(2);
  const [earthData, setEarthData] = useState({ imageUrl: '', lunarData: [] });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    mounted.current = true;
    
    const loadEarthData = async () => {
      if (!mounted.current) return;
      setLoading(true);
      setError(null);
      
      try {
        const data = await fetchEarthData(dateString);
        console.log('Earth Component Data:', {
          hasImageUrl: !!data.imageUrl,
          imageUrlLength: data.imageUrl?.length,
          lunarDataLength: data.lunarData?.length
        });
        
        if (mounted.current && data.imageUrl) {
          setEarthData(data);
        }
      } catch (error) {
        console.error('Error loading Earth data:', error);
        if (mounted.current) {
          setError(error.response?.data?.message || 'Failed to load Earth data');
        }
      } finally {
        if (mounted.current) {
          setLoading(false);
        }
      }
    };

    loadEarthData();
    
    return () => {
      mounted.current = false;
    };
  }, [dateString]);

  if (loading) {
    return <div className={styles.loading}>Loading Earth data...</div>;
  }

  if (error) {
    return (
      <div>
        <DatePicker value={selectedDate} onChange={setSelectedDate} />
        <div className={styles.error}>{error}</div>
      </div>
    );
  }

  return (
    <div>
      <DatePicker value={selectedDate} onChange={setSelectedDate} />
      <div className={styles.grid}>
        <div className={styles.card}>
          <div className={styles.imageContainer}>
            {earthData.imageUrl ? (
              <img 
                src={earthData.imageUrl} 
                alt="Earth from space" 
                onError={(e) => {
                  console.error('Image failed to load:', e);
                  e.target.style.display = 'none';
                }}
              />
            ) : (
              <div className={styles.noImage}>No image available for selected date</div>
            )}
          </div>
          <h2>Earth View {selectedDate ? `for ${selectedDate.toLocaleDateString()}` : '(Current)'}</h2>
        </div>
        <div className={styles.card}>
          <h3>Moon Phase</h3>
          {earthData.lunarData?.length > 0 ? (
            earthData.lunarData.map((data, index) => (
              <p key={index}>{data}</p>
            ))
          ) : (
            <p className={styles.noData}>No lunar data available</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default Earth; 