import React from 'react';
import { AppShell, Title, Container } from '@mantine/core';
import NavBar from './NavBar';

function Layout({ children, activeTab, setActiveTab }) {
  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 300, breakpoint: 'sm' }}
      padding="md"
    >
      <AppShell.Header p="md">
        <Title order={1}>MTSI Dashboard</Title>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <NavBar activeTab={activeTab} setActiveTab={setActiveTab} />
      </AppShell.Navbar>

      <AppShell.Main>
        <Container size="xl">
          {children}
        </Container>
      </AppShell.Main>
    </AppShell>
  );
}

export default Layout; 