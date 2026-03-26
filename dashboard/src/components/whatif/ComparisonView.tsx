'use client';

import type { ExplainResponse, SensitivityResult } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DiffTable } from './DiffTable';

interface ComparisonViewProps {
  original: ExplainResponse;
  result: SensitivityResult;
}

export function ComparisonView({ original, result }: ComparisonViewProps) {
  const deltaSign = result.delta_value >= 0 ? '+' : '';
  const deltaColor = result.delta_value >= 0 ? 'text-green-600' : 'text-red-600';

  return (
    <Card>
      <CardHeader>
        <CardTitle>Comparison</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="flex items-center gap-3 text-lg">
          <span className="tabular-nums font-medium">
            {result.original_value.toFixed(4)}
          </span>
          <span className="text-muted-foreground">&rarr;</span>
          <span className="tabular-nums font-medium">
            {result.modified_value.toFixed(4)}
          </span>
          <span className={`text-sm font-medium ${deltaColor}`}>
            ({'\u0394'} {deltaSign}
            {result.delta_value.toFixed(4)}, {deltaSign}
            {result.delta_percentage.toFixed(1)}%)
          </span>
        </div>

        <DiffTable diffs={result.component_diffs} />
      </CardContent>
    </Card>
  );
}
