'use client';

import { Card, CardContent } from '@/components/ui/card';
import { ConfidenceGauge } from './ConfidenceGauge';

interface ExplainMetadata {
  version: string;
  created_at: string;
  deterministic_hash: string;
  computation_type: string;
}

interface ValueCardProps {
  value: number;
  target: string;
  hash: string;
}

interface ConfidenceCardProps {
  confidence: number;
}

interface SummaryCardProps {
  target: string;
  finalValue: number;
  confidence: number;
  missingImpact: number;
  topDrivers: { name: string; impact: number; rank: number }[];
  metadata: ExplainMetadata;
}

function ValueCard({ value, target, hash }: ValueCardProps) {
  return (
    <Card>
      <CardContent className="pt-2">
        <span className="text-4xl font-bold tabular-nums">{value.toFixed(4)}</span>
        <p className="text-sm text-muted-foreground mt-1">{target}</p>
        <p className="text-xs text-muted-foreground/60 mt-0.5 font-mono">
          {hash.slice(0, 12)}
        </p>
      </CardContent>
    </Card>
  );
}

function ConfidenceCard({ confidence }: ConfidenceCardProps) {
  return (
    <Card>
      <CardContent className="pt-2">
        <ConfidenceGauge confidence={confidence} size="lg" />
      </CardContent>
    </Card>
  );
}

export { ValueCard, ConfidenceCard };

/**
 * SummaryCard is kept as a convenience wrapper that renders the two-column
 * metric layout. Individual cards are also exported for flexibility.
 */
export function SummaryCard({
  target,
  finalValue,
  confidence,
  metadata,
}: SummaryCardProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <ValueCard
        value={finalValue}
        target={target}
        hash={metadata.deterministic_hash}
      />
      <ConfidenceCard confidence={confidence} />
    </div>
  );
}
