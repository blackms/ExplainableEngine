'use client';

interface ConfidenceGaugeProps {
  confidence: number;
  size?: 'sm' | 'lg';
}

function confidenceColor(value: number): string {
  if (value >= 0.8) return '#10b981'; // emerald-500
  if (value >= 0.5) return '#f59e0b'; // amber-500
  return '#f43f5e'; // rose-500
}

function confidenceLabel(value: number): string {
  if (value >= 0.8) return 'High confidence';
  if (value >= 0.5) return 'Moderate — review recommended';
  return 'Low — action required';
}

export function ConfidenceGauge({ confidence, size = 'lg' }: ConfidenceGaugeProps) {
  const pct = (confidence * 100).toFixed(1);
  const color = confidenceColor(confidence);
  const label = confidenceLabel(confidence);

  if (size === 'sm') {
    return (
      <div className="flex items-center gap-2">
        <div className="h-2 flex-1 rounded-full bg-secondary">
          <div
            className="h-full rounded-full transition-all duration-500"
            style={{ width: `${pct}%`, backgroundColor: color }}
          />
        </div>
        <span className="text-xs text-muted-foreground shrink-0">{pct}%</span>
      </div>
    );
  }

  return (
    <div className="space-y-1">
      <span className="text-4xl font-bold tabular-nums">{pct}%</span>
      <div className="h-3 w-full rounded-full bg-secondary">
        <div
          className="h-full rounded-full transition-all duration-500"
          style={{ width: `${pct}%`, backgroundColor: color }}
        />
      </div>
      <span className="text-sm text-muted-foreground">{label}</span>
    </div>
  );
}
