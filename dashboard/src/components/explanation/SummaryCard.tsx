'use client';

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ConfidenceGauge } from './ConfidenceGauge';

interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

interface ExplainMetadata {
  version: string;
  created_at: string;
  deterministic_hash: string;
  computation_type: string;
}

interface SummaryCardProps {
  target: string;
  finalValue: number;
  confidence: number;
  missingImpact: number;
  topDrivers: DriverItem[];
  metadata: ExplainMetadata;
}

function getMissingImpactSeverity(impact: number): 'warning' | 'critical' {
  return impact >= 0.5 ? 'critical' : 'warning';
}

export function SummaryCard({
  target,
  finalValue,
  confidence,
  missingImpact,
  topDrivers,
  metadata,
}: SummaryCardProps) {
  const topThree = [...topDrivers].sort((a, b) => a.rank - b.rank).slice(0, 3);
  const severity = getMissingImpactSeverity(missingImpact);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="text-lg">{target}</CardTitle>
            <p className="text-3xl font-bold tracking-tight mt-1">
              {finalValue.toFixed(4)}
            </p>
          </div>
          <ConfidenceGauge confidence={confidence} size={72} />
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {missingImpact > 0 && (
          <div
            className={`rounded-md px-3 py-2 text-sm font-medium ${
              severity === 'critical'
                ? 'bg-red-500/10 text-red-700 dark:text-red-400'
                : 'bg-yellow-500/10 text-yellow-700 dark:text-yellow-400'
            }`}
          >
            {(missingImpact * 100).toFixed(0)}% of input data is missing
          </div>
        )}

        <div>
          <h4 className="text-sm font-medium text-muted-foreground mb-2">Top Drivers</h4>
          <ol className="space-y-1.5">
            {topThree.map((driver) => (
              <li key={driver.name} className="flex items-center gap-2 text-sm">
                <Badge variant="outline" className="shrink-0 w-7 justify-center text-xs">
                  #{driver.rank}
                </Badge>
                <span className="font-medium truncate">{driver.name}</span>
                <span className="text-muted-foreground ml-auto shrink-0">
                  impact: {driver.impact.toFixed(2)}
                </span>
              </li>
            ))}
          </ol>
        </div>
      </CardContent>

      <CardFooter>
        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
          <span>Hash: {metadata.deterministic_hash.slice(0, 12)}...</span>
          <span>{new Date(metadata.created_at).toLocaleString()}</span>
          <span>{metadata.computation_type}</span>
          <span>v{metadata.version}</span>
        </div>
      </CardFooter>
    </Card>
  );
}
