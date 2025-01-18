import React from 'react';
import { MantineProvider } from '@mantine/core';
import UnifiedView from './components/UnifiedView';
import styles from './styles/App.module.css';

function App() {
  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <div className={styles.container}>
        <header className={styles.header}>
          <h1>MTSI Dashboard</h1>
        </header>
        <main className={styles.main}>
          <UnifiedView />
        </main>
      </div>
    </MantineProvider>
  );
}

export default App; 