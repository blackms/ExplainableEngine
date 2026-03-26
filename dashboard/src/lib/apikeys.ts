export interface ApiKey {
  id: string;
  name: string;
  key: string;
  prefix: string;
  createdAt: string;
}

const STORAGE_KEY = 'ee_api_keys';

export function generateApiKey(name: string): ApiKey {
  const key = `ee_${crypto.randomUUID().replace(/-/g, '')}`;
  const apiKey: ApiKey = {
    id: crypto.randomUUID(),
    name,
    key,
    prefix: key.substring(0, 11) + '...',
    createdAt: new Date().toISOString(),
  };
  const existing = getApiKeys();
  existing.push(apiKey);
  localStorage.setItem(STORAGE_KEY, JSON.stringify(existing));
  return apiKey;
}

export function getApiKeys(): ApiKey[] {
  if (typeof window === 'undefined') return [];
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return [];
  try {
    return JSON.parse(raw) as ApiKey[];
  } catch {
    return [];
  }
}

export function revokeApiKey(id: string): void {
  const keys = getApiKeys().filter((k) => k.id !== id);
  localStorage.setItem(STORAGE_KEY, JSON.stringify(keys));
}
