'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

interface DriverRankingProps {
  drivers: DriverItem[];
}

export function DriverRanking({ drivers }: DriverRankingProps) {
  const sorted = [...drivers].sort((a, b) => a.rank - b.rank);
  const maxImpact = Math.max(...sorted.map((d) => d.impact), 0.01);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Top Drivers</CardTitle>
      </CardHeader>
      <CardContent>
        <ol className="space-y-3">
          {sorted.map((driver) => (
            <li key={driver.name} className="flex items-center gap-3">
              <span className="text-xs rounded-full bg-primary/10 text-primary w-6 h-6 flex items-center justify-center shrink-0 font-semibold">
                {driver.rank}
              </span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium truncate">
                    {driver.name}
                  </span>
                  <span className="text-xs text-muted-foreground ml-2 shrink-0 tabular-nums">
                    {driver.impact.toFixed(2)}
                  </span>
                </div>
                <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                  <div
                    className="h-full rounded-full bg-primary transition-all duration-300"
                    style={{
                      width: `${(driver.impact / maxImpact) * 100}%`,
                    }}
                  />
                </div>
              </div>
            </li>
          ))}
        </ol>
      </CardContent>
    </Card>
  );
}
