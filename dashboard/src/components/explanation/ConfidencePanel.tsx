'use client';

import { useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfidenceGauge } from './ConfidenceGauge';

interface ConfidencePanelProps {
  overall: number;
  perNode: Record<string, number>;
  missingImpact: number;
}

function getBarColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-green-500';
  if (confidence >= 0.5) return 'bg-yellow-500';
  return 'bg-red-500';
}

function getBarBgColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-green-500/20';
  if (confidence >= 0.5) return 'bg-yellow-500/20';
  return 'bg-red-500/20';
}

export function ConfidencePanel({ overall, perNode, missingImpact }: ConfidencePanelProps) {
  const sortedNodes = useMemo(() => {
    return Object.entries(perNode)
      .map(([name, confidence]) => ({ name, confidence }))
      .sort((a, b) => a.confidence - b.confidence);
  }, [perNode]);

  const missingPct = (missingImpact * 100).toFixed(0);
  const isMissingCritical = missingImpact >= 0.2;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Confidence Analysis</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Overall confidence */}
        <div className="flex items-center gap-4">
          <ConfidenceGauge confidence={overall} size={96} />
          <div>
            <div className="text-sm font-medium">Overall Confidence</div>
            <div className="text-xs text-muted-foreground mt-0.5">
              Aggregated across all components
            </div>
          </div>
        </div>

        {/* Per-component breakdown */}
        {sortedNodes.length > 0 && (
          <div>
            <h4 className="text-sm font-medium text-muted-foreground mb-3">
              Per-Component Confidence
            </h4>
            <ul className="space-y-2">
              {sortedNodes.map(({ name, confidence }) => (
                <li key={name} className="space-y-1">
                  <div className="flex items-center justify-between text-xs">
                    <span className="font-medium truncate mr-2">{name}</span>
                    <span className="text-muted-foreground shrink-0">
                      {(confidence * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className={`h-2 w-full rounded-full overflow-hidden ${getBarBgColor(confidence)}`}>
                    <div
                      className={`h-full rounded-full transition-all duration-300 ${getBarColor(confidence)}`}
                      style={{ width: `${(confidence * 100).toFixed(1)}%` }}
                    />
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Missing data section */}
        {missingImpact > 0 && (
          <div>
            <h4 className="text-sm font-medium text-muted-foreground mb-2">
              Missing Data Impact
            </h4>
            <div className={`rounded-md px-3 py-2 text-sm ${
              isMissingCritical
                ? 'bg-red-500/10 text-red-700 dark:text-red-400'
                : 'bg-yellow-500/10 text-yellow-700 dark:text-yellow-400'
            }`}>
              <div className="flex items-center justify-between mb-1.5">
                <span className="font-medium">
                  {missingPct}% of input data is missing
                </span>
              </div>
              <div className={`h-2 w-full rounded-full overflow-hidden ${
                isMissingCritical ? 'bg-red-500/20' : 'bg-yellow-500/20'
              }`}>
                <div
                  className={`h-full rounded-full transition-all duration-300 ${
                    isMissingCritical ? 'bg-red-500' : 'bg-yellow-500'
                  }`}
                  style={{ width: `${missingPct}%` }}
                />
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
