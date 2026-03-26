'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from './client';
import type {
  AskRequest,
  ExplainRequest,
  LLMNarrativeRequest,
  ListOptions,
  Modification,
  SummaryRequest,
} from './types';

export function useExplanation(id: string | undefined) {
  return useQuery({
    queryKey: ['explanation', id],
    queryFn: () => api.getExplanation(id!),
    enabled: !!id,
  });
}

export function useGraph(id: string | undefined, format: string = 'json') {
  return useQuery({
    queryKey: ['graph', id, format],
    queryFn: () => api.getGraph(id!, format),
    enabled: !!id,
  });
}

export function useNarrative(
  id: string | undefined,
  level: string = 'basic',
  lang: string = 'en'
) {
  return useQuery({
    queryKey: ['narrative', id, level, lang],
    queryFn: () => api.getNarrative(id!, level, lang),
    enabled: !!id,
  });
}

export function useCreateExplanation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (req: ExplainRequest) => api.createExplanation(req),
    onSuccess: (data) => {
      queryClient.setQueryData(['explanation', data.id], data);
    },
  });
}

export function useWhatIf(id: string | undefined) {
  return useMutation({
    mutationFn: (modifications: Modification[]) =>
      api.whatIf(id!, modifications),
  });
}

export function useExplanationList(
  opts?: ListOptions,
  queryOpts?: { refetchInterval?: number | false }
) {
  return useQuery({
    queryKey: ['explanations', opts],
    queryFn: () => api.listExplanations(opts),
    refetchInterval: queryOpts?.refetchInterval,
  });
}

export function useStats() {
  return useQuery({
    queryKey: ['stats'],
    queryFn: () => api.stats(),
    refetchInterval: 30000,
  });
}

export function useHealth() {
  return useQuery({
    queryKey: ['health'],
    queryFn: () => api.health(),
    refetchInterval: 30_000,
  });
}

export function useLLMNarrative(id: string | undefined) {
  return useMutation({
    mutationFn: (req: LLMNarrativeRequest) => api.generateLLMNarrative(id!, req),
  });
}

export function useAskQuestion(id: string | undefined) {
  return useMutation({
    mutationFn: (req: AskRequest) => api.askQuestion(id!, req),
  });
}

export function useGenerateSummary(id: string | undefined) {
  return useMutation({
    mutationFn: (req: SummaryRequest) => api.generateSummary(id!, req),
  });
}

// AIP Sentiment hooks

export function useAIPExplainTicker(ticker: string | undefined) {
  return useQuery({
    queryKey: ['aip-explain', ticker],
    queryFn: () => api.explainTicker(ticker!),
    enabled: !!ticker,
  });
}

export function useAIPMarketMood() {
  return useQuery({
    queryKey: ['aip-market-mood'],
    queryFn: () => api.explainMarketMood(),
    refetchInterval: 60_000,
  });
}
