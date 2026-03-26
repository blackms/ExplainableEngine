'use client';

import { useEffect, useRef, useState } from 'react';
import {
  getPreferences,
  savePreferences,
  type UserPreferences,
} from '@/lib/preferences';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

export function PreferencesForm() {
  const [prefs, setPrefs] = useState<UserPreferences | null>(null);
  const [saved, setSaved] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    setPrefs(getPreferences());
  }, []);

  const handleChange = (partial: Partial<UserPreferences>) => {
    const updated = savePreferences(partial);
    setPrefs(updated);
    setSaved(true);
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => setSaved(false), 2000);

    // Apply theme change immediately
    if (partial.theme) {
      const root = document.documentElement;
      if (partial.theme === 'dark') {
        root.classList.add('dark');
      } else if (partial.theme === 'light') {
        root.classList.remove('dark');
      } else {
        const isDark = window.matchMedia(
          '(prefers-color-scheme: dark)',
        ).matches;
        root.classList.toggle('dark', isDark);
      }
    }
  };

  if (!prefs) return null;

  return (
    <div className="space-y-6">
      <div className="grid gap-2">
        <Label htmlFor="theme-select">Theme</Label>
        <Select
          value={prefs.theme}
          onValueChange={(val) =>
            handleChange({
              theme: val as UserPreferences['theme'],
            })
          }
        >
          <SelectTrigger id="theme-select">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="light">Light</SelectItem>
            <SelectItem value="dark">Dark</SelectItem>
            <SelectItem value="system">System</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="grid gap-2">
        <Label htmlFor="language-select">Default Language</Label>
        <Select
          value={prefs.language}
          onValueChange={(val) =>
            handleChange({
              language: val as UserPreferences['language'],
            })
          }
        >
          <SelectTrigger id="language-select">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="en">English</SelectItem>
            <SelectItem value="it">Italiano</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="grid gap-2">
        <Label htmlFor="narrative-select">Default Narrative Level</Label>
        <Select
          value={prefs.narrativeLevel}
          onValueChange={(val) =>
            handleChange({
              narrativeLevel: val as UserPreferences['narrativeLevel'],
            })
          }
        >
          <SelectTrigger id="narrative-select">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="basic">Basic</SelectItem>
            <SelectItem value="advanced">Advanced</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {saved && (
        <p className="text-sm text-muted-foreground animate-in fade-in">
          Saved
        </p>
      )}
    </div>
  );
}
