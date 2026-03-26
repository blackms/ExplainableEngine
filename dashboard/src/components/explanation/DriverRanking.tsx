'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

interface DriverRankingProps {
  drivers: DriverItem[];
}

function getRankColor(rank: number): 'default' | 'secondary' | 'outline' {
  if (rank === 1) return 'default';
  if (rank <= 3) return 'secondary';
  return 'outline';
}

export function DriverRanking({ drivers }: DriverRankingProps) {
  const sorted = [...drivers].sort((a, b) => a.rank - b.rank);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Top Drivers</CardTitle>
      </CardHeader>
      <CardContent>
        <ul className="space-y-3">
          {sorted.map((driver) => (
            <li key={driver.name} className="flex items-center gap-3">
              <Badge variant={getRankColor(driver.rank)} className="shrink-0 w-8 justify-center">
                #{driver.rank}
              </Badge>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium truncate">{driver.name}</span>
                  <span className="text-xs text-muted-foreground ml-2 shrink-0">
                    {driver.impact.toFixed(2)}
                  </span>
                </div>
                <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                  <div
                    className="h-full rounded-full bg-primary transition-all duration-300"
                    style={{ width: `${Math.min(driver.impact * 100, 100)}%` }}
                  />
                </div>
              </div>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
