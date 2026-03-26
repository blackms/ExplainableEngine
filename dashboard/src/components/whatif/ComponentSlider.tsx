'use client';

import { useCallback } from 'react';
import { Slider } from '@/components/ui/slider';

interface ComponentSliderProps {
  name: string;
  originalValue: number;
  value: number;
  onChange: (value: number) => void;
  min: number;
  max: number;
  step: number;
}

function formatDelta(current: number, original: number): { text: string; positive: boolean } {
  const delta = current - original;
  const pct = original !== 0 ? (delta / Math.abs(original)) * 100 : 0;
  const sign = delta >= 0 ? '+' : '';
  return {
    text: `${sign}${delta.toFixed(4)} (${sign}${pct.toFixed(1)}%)`,
    positive: delta >= 0,
  };
}

export function ComponentSlider({
  name,
  originalValue,
  value,
  onChange,
  min,
  max,
  step,
}: ComponentSliderProps) {
  const delta = formatDelta(value, originalValue);
  const hasChanged = Math.abs(value - originalValue) > 0.001;

  const handleValueChange = useCallback(
    (newValue: number | readonly number[]) => {
      const v = typeof newValue === 'number' ? newValue : newValue[0];
      onChange(v);
    },
    [onChange],
  );

  const handleReset = useCallback(() => {
    onChange(originalValue);
  }, [onChange, originalValue]);

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium">{name}</label>
        <div className="flex items-center gap-2">
          <span className="text-sm tabular-nums">{value.toFixed(4)}</span>
          {hasChanged && (
            <>
              <span
                className={`text-xs tabular-nums font-medium ${
                  delta.positive ? 'text-green-600' : 'text-red-600'
                }`}
              >
                {delta.text}
              </span>
              <button
                type="button"
                onClick={handleReset}
                className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                title="Reset to original"
              >
                &#8634;
              </button>
            </>
          )}
        </div>
      </div>
      <Slider
        value={[value]}
        onValueChange={handleValueChange}
        min={min}
        max={max}
        step={step}
      />
      <div className="flex justify-between text-xs text-muted-foreground">
        <span>{min.toFixed(2)}</span>
        <span className="text-muted-foreground/60">original: {originalValue.toFixed(4)}</span>
        <span>{max.toFixed(2)}</span>
      </div>
    </div>
  );
}
