import { useState } from 'react';

export function useDateSelection(initialDaysAgo = 2) {
  const [selectedDate, setSelectedDate] = useState(() => {
    const date = new Date();
    date.setDate(date.getDate() - initialDaysAgo);
    return date;
  });

  const dateString = selectedDate ? selectedDate.toISOString().split('T')[0] : '';

  return {
    selectedDate,
    setSelectedDate,
    dateString
  };
} 