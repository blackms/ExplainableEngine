'use client';

import { useEffect, useState } from 'react';
import { getPreferences, savePreferences } from '@/lib/preferences';

export function ThemeToggle() {
  const [theme, setTheme] = useState<'light' | 'dark' | 'system'>('system');

  useEffect(() => {
    setTheme(getPreferences().theme);
  }, []);

  useEffect(() => {
    const root = document.documentElement;
    if (theme === 'dark') {
      root.classList.add('dark');
    } else if (theme === 'light') {
      root.classList.remove('dark');
    } else {
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      root.classList.toggle('dark', isDark);
    }
  }, [theme]);

  const toggle = () => {
    const next =
      theme === 'light' ? 'dark' : theme === 'dark' ? 'system' : 'light';
    setTheme(next);
    savePreferences({ theme: next });
  };

  return (
    <button
      onClick={toggle}
      className="text-xs px-2 py-1 rounded border border-border hover:bg-accent"
      title={`Theme: ${theme}`}
    >
      {theme === 'light' ? 'Light' : theme === 'dark' ? 'Dark' : 'Auto'}
    </button>
  );
}
