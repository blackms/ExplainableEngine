'use client';

import { useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface ConfidencePanelProps {
  overall: number;
  perNode: Record<string, number>;
  missingImpact: number;
}

function barColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-emerald-500';
  if (confidence >= 0.5) return 'bg-amber-500';
  return 'bg-rose-500';
}

function barBgColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-emerald-500/20';
  if (confidence >= 0.5) return 'bg-amber-500/20';
  return 'bg-rose-500/20';
}

export function ConfidencePanel({
  overall,
  perNode,
  missingImpact,
}: ConfidencePanelProps) {
  const sortedNodes = useMemo(() => {
    return Object.entries(perNode)
      .map(([name, confidence]) => ({ name, confidence }))
      .sort((a, b) => a.confidence - b.confidence);
  }, [perNode]);

  const missingPct = (missingImpact * 100).toFixed(0);
  const hasMissing = missingImpact > 0;
  const isMissingCritical = missingImpact >= 0.2;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Data Quality</CardTitle>
      </CardHeader>
      <CardContent className="space-y-5">
        {/* Missing data section */}
        {hasMissing ? (
          <div
            className={`rounded-lg border px-4 py-3 text-sm ${
              isMissingCritical
                ? 'border-rose-500/40 bg-rose-500/5 text-rose-700 dark:text-rose-400'
                : 'border-amber-500/40 bg-amber-500/5 text-amber-700 dark:text-amber-400'
            }`}
          >
            <p className="font-medium">
              {missingPct}% of input data is missing
            </p>
            <div
              className={`mt-2 h-2 w-full rounded-full overflow-hidden ${
                isMissingCritical ? 'bg-rose-500/20' : 'bg-amber-500/20'
              }`}
            >
              <div
                className={`h-full rounded-full transition-all duration-300 ${
                  isMissingCritical ? 'bg-rose-500' : 'bg-amber-500'
                }`}
                style={{ width: `${missingPct}%` }}
              />
            </div>
          </div>
        ) : (
          <div className="rounded-lg border border-emerald-500/40 bg-emerald-500/5 px-4 py-3 text-sm text-emerald-700 dark:text-emerald-400">
            <span className="font-medium">&#10003; No missing data</span>
          </div>
        )}

        {/* Per-node confidence sorted weakest first */}
        {sortedNodes.length > 0 && (
          <div>
            <h4 className="text-sm font-medium text-muted-foreground mb-3">
              Per-Node Confidence
            </h4>
            <ul className="space-y-2">
              {sortedNodes.map(({ name, confidence }) => (
                <li key={name} className="space-y-1">
                  <div className="flex items-center justify-between text-xs">
                    <span className="font-medium truncate mr-2">{name}</span>
                    <span className="text-muted-foreground shrink-0 tabular-nums">
                      {(confidence * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div
                    className={`h-2 w-full rounded-full overflow-hidden ${barBgColor(confidence)}`}
                  >
                    <div
                      className={`h-full rounded-full transition-all duration-300 ${barColor(confidence)}`}
                      style={{ width: `${(confidence * 100).toFixed(1)}%` }}
                    />
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
