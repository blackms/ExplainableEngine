'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Slider } from '@/components/ui/slider';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import type { ListOptions } from '@/lib/api/types';

interface FilterPanelProps {
  filters: ListOptions;
  onChange: (filters: ListOptions) => void;
  total: number;
}

const DEBOUNCE_MS = 300;

const DEFAULT_CONFIDENCE: [number, number] = [0, 1];

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
    [filters, onChange]
  );

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  const confidenceValue: [number, number] = [
    filters.min_confidence ?? 0,
    filters.max_confidence ?? 1,
  ];

  const handleConfidenceChange = (value: number | readonly number[]) => {
    const values = value as readonly number[];
    onChange({
      ...filters,
      min_confidence: values[0] === 0 ? undefined : values[0],
      max_confidence: values[1] === 1 ? undefined : values[1],
      cursor: undefined,
    });
  };

  const handleDateChange = (field: 'from' | 'to', value: string) => {
    onChange({
      ...filters,
      [field]: value || undefined,
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

  return (
    <div className="rounded-lg border bg-card p-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        {/* Target search */}
        <div className="space-y-1.5">
          <Label htmlFor="audit-search" className="text-xs font-medium text-muted-foreground">
            Target
          </Label>
          <div className="relative">
            <Search className="pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
            <Input
              id="audit-search"
              placeholder="Search target..."
              className="pl-7"
              value={searchText}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleSearchChange(e.target.value)}
            />
          </div>
        </div>

        {/* Confidence range */}
        <div className="space-y-1.5">
          <Label className="text-xs font-medium text-muted-foreground">
            Confidence: {(confidenceValue[0] * 100).toFixed(0)}% &ndash; {(confidenceValue[1] * 100).toFixed(0)}%
          </Label>
          <div className="pt-2">
            <Slider
              min={0}
              max={1}
              step={0.05}
              value={confidenceValue}
              onValueChange={handleConfidenceChange}
            />
          </div>
        </div>

        {/* Date from */}
        <div className="space-y-1.5">
          <Label htmlFor="audit-from" className="text-xs font-medium text-muted-foreground">
            From
          </Label>
          <Input
            id="audit-from"
            type="date"
            value={filters.from ?? ''}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleDateChange('from', e.target.value)}
          />
        </div>

        {/* Date to */}
        <div className="space-y-1.5">
          <Label htmlFor="audit-to" className="text-xs font-medium text-muted-foreground">
            To
          </Label>
          <Input
            id="audit-to"
            type="date"
            value={filters.to ?? ''}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleDateChange('to', e.target.value)}
          />
        </div>
      </div>

      <div className="mt-3 flex items-center justify-between">
        <p className="text-xs text-muted-foreground">
          {total} result{total !== 1 ? 's' : ''}
        </p>
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={handleReset}>
            <RotateCcw className="size-3.5" />
            Reset filters
          </Button>
        )}
      </div>
    </div>
  );
}
