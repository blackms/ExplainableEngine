import type {
  ExplainRequest,
  ExplainResponse,
  GraphResponse,
  Modification,
  NarrativeResult,
  SensitivityResult,
} from './types';

const API_BASE = '/api';

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  });
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(error.error || `API error: ${res.status}`);
  }
  return res.json();
}

export const api = {
  createExplanation: (req: ExplainRequest) =>
    fetchAPI<ExplainResponse>('/explain', {
      method: 'POST',
      body: JSON.stringify(req),
    }),

  getExplanation: (id: string) =>
    fetchAPI<ExplainResponse>(`/explain/${id}`),

  getGraph: (id: string, format: string = 'json') =>
    fetchAPI<GraphResponse>(`/explain/${id}/graph?format=${format}`),

  getNarrative: (id: string, level: string = 'basic', lang: string = 'en') =>
    fetchAPI<NarrativeResult>(
      `/explain/${id}/narrative?level=${level}&lang=${lang}`
    ),

  whatIf: (id: string, modifications: Modification[]) =>
    fetchAPI<SensitivityResult>(`/explain/${id}/whatif`, {
      method: 'POST',
      body: JSON.stringify({ modifications }),
    }),

  health: () =>
    fetchAPI<{ status: string; version: string; uptime: string }>('/health'),
};
