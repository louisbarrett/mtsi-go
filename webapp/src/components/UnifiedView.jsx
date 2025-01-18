import React, { useState, useEffect } from 'react';
import { Card, Group, Button, Text, LoadingOverlay } from '@mantine/core';
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react';
import { fetchEarthData, fetchSunData } from '../services/api';
import styles from '../styles/UnifiedView.module.css';

function UnifiedView() {
  const [date, setDate] = useState(() => {
    const d = new Date();
    d.setDate(d.getDate() - 2);
    return d;
  });
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState({
    earth: { imageUrl: '', lunarData: [] },
    sun: { imageUrl: '' }
  });
  const [transition, setTransition] = useState(false);

  const loadData = async (targetDate) => {
    setTransition(true);
    const dateString = targetDate.toISOString().split('T')[0];
    
    try {
      const [earthResult, sunResult] = await Promise.all([
        fetchEarthData(dateString),
        fetchSunData(dateString)
      ]);

      setData({
        earth: earthResult,
        sun: sunResult
      });
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
      setTimeout(() => setTransition(false), 300);
    }
  };

  useEffect(() => {
    loadData(date);
  }, [date]);

  const changeDate = (delta) => {
    const newDate = new Date(date);
    newDate.setDate(date.getDate() + delta);
    
    const minDate = new Date(2015, 5, 13);
    const maxDate = new Date();
    maxDate.setDate(maxDate.getDate() - 2);

    if (newDate >= minDate && newDate <= maxDate) {
      setDate(newDate);
    }
  };

  if (loading) {
    return <div className={styles.loading}>Loading data...</div>;
  }

  return (
    <div className={styles.container}>
      <Group position="center" mb="xl">
        <Button onClick={() => changeDate(-1)}>
          <IconChevronLeft size={18} />
        </Button>
        <Text>{date.toLocaleDateString()}</Text>
        <Button onClick={() => changeDate(1)}>
          <IconChevronRight size={18} />
        </Button>
      </Group>

      <div className={styles.grid}>
        <Card className={`${styles.imageCard} ${transition ? styles.fadeOut : styles.fadeIn}`}>
          <div className={styles.imageWrapper}>
            {data.earth.imageUrl && (
              <img
                src={data.earth.imageUrl}
                alt="Earth view"
                className={styles.image}
              />
            )}
          </div>
          <Text size="lg" weight={500} mt="md">Earth View</Text>
        </Card>

        <Card className={`${styles.imageCard} ${transition ? styles.fadeOut : styles.fadeIn}`}>
          <div className={styles.imageWrapper}>
            {data.sun.imageUrl && (
              <img
                src={data.sun.imageUrl}
                alt="Solar activity"
                className={styles.image}
              />
            )}
          </div>
          <Text size="lg" weight={500} mt="md">Solar Activity</Text>
        </Card>

        <Card className={styles.lunarCard}>
          <Text size="lg" weight={500} mb="md">Lunar Data</Text>
          {data.earth.lunarData?.map((item, index) => (
            <Text key={index} mb="xs">{item}</Text>
          ))}
        </Card>
      </div>
    </div>
  );
}

export default UnifiedView; 