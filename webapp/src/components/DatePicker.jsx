import React from 'react';
import styles from '../styles/DatePicker.module.css';
import commonStyles from '../styles/common.module.css';

function DatePicker({ value, onChange }) {
  const maxDate = new Date();
  const minDate = new Date(2015, 5, 13);
  const twoDaysAgo = new Date(maxDate);
  twoDaysAgo.setDate(maxDate.getDate() - 2);

  const handleChange = (e) => {
    const date = new Date(e.target.value);
    if (date > twoDaysAgo) {
      alert('Images are only available up to 2 days ago');
      return;
    }
    onChange(date);
  };

  return (
    <div className={`${styles.datePickerContainer} ${commonStyles.card}`}>
      <label htmlFor="datePicker">Select Date</label>
      <input
        id="datePicker"
        type="date"
        value={value ? value.toISOString().split('T')[0] : ''}
        onChange={handleChange}
        min={minDate.toISOString().split('T')[0]}
        max={twoDaysAgo.toISOString().split('T')[0]}
      />
      <small>Images are available from 2015-06-13 up to 2 days ago</small>
    </div>
  );
}

export default DatePicker; 