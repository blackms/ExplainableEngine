import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

interface SensitivityQuickViewProps {
  drivers: DriverItem[];
  explanationId: string;
}

export function SensitivityQuickView({
  drivers,
  explanationId,
}: SensitivityQuickViewProps) {
  const topDrivers = [...drivers]
    .sort((a, b) => b.impact - a.impact)
    .slice(0, 3);

  if (topDrivers.length === 0) return null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Sensitivity</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <ul className="space-y-3">
          {topDrivers.map((driver) => (
            <li key={driver.name} className="space-y-1">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium truncate">
                  {driver.name}
                </span>
                <span className="text-xs text-muted-foreground ml-2 shrink-0">
                  {(driver.impact * 100).toFixed(1)}%
                </span>
              </div>
              <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                <div
                  className="h-full rounded-full bg-primary transition-all duration-300"
                  style={{ width: `${Math.min(driver.impact * 100, 100)}%` }}
                />
              </div>
            </li>
          ))}
        </ul>
        <Link
          href={`/explain/${explanationId}/whatif`}
          className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
        >
          Run What-if Analysis &rarr;
        </Link>
      </CardContent>
    </Card>
  );
}
