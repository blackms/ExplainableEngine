import type {
  AskRequest,
  AskResponse,
  ExplainRequest,
  ExplainResponse,
  GraphResponse,
  LLMNarrativeRequest,
  LLMNarrativeResponse,
  ListOptions,
  ListResult,
  Modification,
  NarrativeResult,
  SensitivityResult,
  SummaryRequest,
  SummaryResponse,
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

  listExplanations: (opts?: ListOptions) => {
    const params = new URLSearchParams();
    if (opts?.cursor) params.set('cursor', opts.cursor);
    if (opts?.limit) params.set('limit', String(opts.limit));
    if (opts?.target) params.set('target', opts.target);
    if (opts?.min_confidence !== undefined)
      params.set('min_confidence', String(opts.min_confidence));
    if (opts?.max_confidence !== undefined)
      params.set('max_confidence', String(opts.max_confidence));
    if (opts?.from) params.set('from', opts.from);
    if (opts?.to) params.set('to', opts.to);
    const qs = params.toString();
    return fetchAPI<ListResult>(`/explain${qs ? `?${qs}` : ''}`);
  },

  generateLLMNarrative: (id: string, req: LLMNarrativeRequest) =>
    fetchAPI<LLMNarrativeResponse>(`/explain/${id}/narrative-llm`, {
      method: 'POST',
      body: JSON.stringify(req),
    }),

  askQuestion: (id: string, req: AskRequest) =>
    fetchAPI<AskResponse>(`/explain/${id}/ask`, {
      method: 'POST',
      body: JSON.stringify(req),
    }),

  generateSummary: (id: string, req: SummaryRequest) =>
    fetchAPI<SummaryResponse>(`/explain/${id}/summary`, {
      method: 'POST',
      body: JSON.stringify(req),
    }),

  health: () =>
    fetchAPI<{ status: string; version: string; uptime: string }>('/health'),

  stats: () =>
    fetchAPI<{ total_explanations: number; status: string }>('/stats'),
};
