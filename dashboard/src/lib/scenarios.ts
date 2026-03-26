import type { Modification, SensitivityResult } from '@/lib/api/types';

export interface SavedScenario {
  id: string;
  name: string;
  explanationId: string;
  modifications: Modification[];
  result: SensitivityResult;
  savedAt: string; // ISO timestamp
}

const STORAGE_KEY = 'ee_scenarios';

export function saveScenario(scenario: Omit<SavedScenario, 'id' | 'savedAt'>): SavedScenario {
  const saved: SavedScenario = {
    ...scenario,
    id: crypto.randomUUID(),
    savedAt: new Date().toISOString(),
  };
  const existing = getScenarios(scenario.explanationId);
  existing.push(saved);
  localStorage.setItem(`${STORAGE_KEY}_${scenario.explanationId}`, JSON.stringify(existing));
  return saved;
}

export function getScenarios(explanationId: string): SavedScenario[] {
  const raw = localStorage.getItem(`${STORAGE_KEY}_${explanationId}`);
  return raw ? JSON.parse(raw) : [];
}

export function deleteScenario(explanationId: string, scenarioId: string) {
  const existing = getScenarios(explanationId).filter(s => s.id !== scenarioId);
  localStorage.setItem(`${STORAGE_KEY}_${explanationId}`, JSON.stringify(existing));
}

export function clearScenarios(explanationId: string) {
  localStorage.removeItem(`${STORAGE_KEY}_${explanationId}`);
}
