'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { ListOptions } from '@/lib/api/types';

interface FilterPanelProps {
  filters: ListOptions;
  onChange: (filters: ListOptions) => void;
  total: number;
}

const DEBOUNCE_MS = 300;

type ConfidencePreset = 'all' | 'high' | 'moderate' | 'low';
type DatePreset = 'all' | '24h' | '7d' | '30d';

function confidencePresetToRange(preset: ConfidencePreset): {
  min_confidence?: number;
  max_confidence?: number;
} {
  switch (preset) {
    case 'high':
      return { min_confidence: 0.8, max_confidence: undefined };
    case 'moderate':
      return { min_confidence: 0.5, max_confidence: 0.8 };
    case 'low':
      return { min_confidence: undefined, max_confidence: 0.5 };
    default:
      return { min_confidence: undefined, max_confidence: undefined };
  }
}

function rangeToConfidencePreset(
  min?: number,
  max?: number,
): ConfidencePreset {
  if (min === 0.8 && max === undefined) return 'high';
  if (min === 0.5 && max === 0.8) return 'moderate';
  if (min === undefined && max === 0.5) return 'low';
  return 'all';
}

function datePresetToRange(preset: DatePreset): {
  from?: string;
  to?: string;
} {
  if (preset === 'all') return { from: undefined, to: undefined };
  const now = new Date();
  const ms =
    preset === '24h'
      ? 24 * 60 * 60 * 1000
      : preset === '7d'
        ? 7 * 24 * 60 * 60 * 1000
        : 30 * 24 * 60 * 60 * 1000;
  const from = new Date(now.getTime() - ms).toISOString().split('T')[0];
  return { from, to: undefined };
}

function rangeToDatePreset(from?: string): DatePreset {
  if (!from) return 'all';
  const diff = Date.now() - new Date(from).getTime();
  const day = 24 * 60 * 60 * 1000;
  if (Math.abs(diff - day) < day * 0.1) return '24h';
  if (Math.abs(diff - 7 * day) < day * 0.5) return '7d';
  if (Math.abs(diff - 30 * day) < day * 2) return '30d';
  return 'all';
}

export function FilterPanel({ filters, onChange, total }: FilterPanelProps) {
  const [searchText, setSearchText] = useState(filters.target ?? '');
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleSearchChange = useCallback(
    (value: string) => {
      setSearchText(value);
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => {
        onChange({
          ...filters,
          target: value || undefined,
          cursor: undefined,
        });
      }, DEBOUNCE_MS);
    },
    [filters, onChange],
  );

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  const handleConfidenceChange = (value: string | null) => {
    if (!value) return;
    const preset = value as ConfidencePreset;
    const range = confidencePresetToRange(preset);
    onChange({
      ...filters,
      ...range,
      cursor: undefined,
    });
  };

  const handleDateChange = (value: string | null) => {
    if (!value) return;
    const preset = value as DatePreset;
    const range = datePresetToRange(preset);
    onChange({
      ...filters,
      ...range,
      cursor: undefined,
    });
  };

  const handleReset = () => {
    setSearchText('');
    onChange({ limit: filters.limit });
  };

  const hasActiveFilters =
    !!filters.target ||
    filters.min_confidence !== undefined ||
    filters.max_confidence !== undefined ||
    !!filters.from ||
    !!filters.to;

  const confidencePreset = rangeToConfidencePreset(
    filters.min_confidence,
    filters.max_confidence,
  );
  const datePreset = rangeToDatePreset(filters.from);

  return (
    <div className="flex items-center gap-3 flex-wrap mb-4">
      <div className="relative w-64">
        <Search className="pointer-events-none absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Search target..."
          className="pl-8"
          value={searchText}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleSearchChange(e.target.value)
          }
        />
      </div>

      <Select
        value={confidencePreset}
        onValueChange={handleConfidenceChange}
      >
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Confidence" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All</SelectItem>
          <SelectItem value="high">High (&ge;80%)</SelectItem>
          <SelectItem value="moderate">Moderate</SelectItem>
          <SelectItem value="low">Low (&lt;50%)</SelectItem>
        </SelectContent>
      </Select>

      <Select value={datePreset} onValueChange={handleDateChange}>
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Date range" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All time</SelectItem>
          <SelectItem value="24h">Last 24h</SelectItem>
          <SelectItem value="7d">Last 7 days</SelectItem>
          <SelectItem value="30d">Last 30 days</SelectItem>
        </SelectContent>
      </Select>

      {hasActiveFilters && (
        <Button variant="ghost" size="sm" onClick={handleReset}>
          <RotateCcw className="size-3.5" />
          Reset
        </Button>
      )}

      <span className="ml-auto text-sm text-muted-foreground">
        {total} result{total !== 1 ? 's' : ''}
      </span>
    </div>
  );
}
