export interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  language: 'en' | 'it';
  narrativeLevel: 'basic' | 'advanced';
}

const STORAGE_KEY = 'ee_preferences';

const defaults: UserPreferences = {
  theme: 'system',
  language: 'en',
  narrativeLevel: 'basic',
};

export function getPreferences(): UserPreferences {
  if (typeof window === 'undefined') return defaults;
  const raw = localStorage.getItem(STORAGE_KEY);
  return raw ? { ...defaults, ...JSON.parse(raw) } : defaults;
}

export function savePreferences(
  prefs: Partial<UserPreferences>,
): UserPreferences {
  const current = getPreferences();
  const updated = { ...current, ...prefs };
  localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
  return updated;
}
