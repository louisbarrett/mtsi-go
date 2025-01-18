import React from 'react';
import { Stack, NavLink } from '@mantine/core';

function NavBar({ activeTab, setActiveTab }) {
  const navItems = [
    { id: 'earth', label: 'Earth Status' },
    { id: 'sun', label: 'Solar Activity' },
    { id: 'orrery', label: 'Orrery View' },
  ];

  return (
    <Stack gap="sm">
      {navItems.map((item) => (
        <NavLink
          key={item.id}
          label={item.label}
          active={activeTab === item.id}
          onClick={() => setActiveTab(item.id)}
          variant="filled"
        />
      ))}
    </Stack>
  );
}

export default NavBar; 