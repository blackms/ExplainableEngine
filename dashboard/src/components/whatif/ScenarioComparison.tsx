'use client';

import type { SavedScenario } from '@/lib/scenarios';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface ScenarioComparisonProps {
  scenarios: SavedScenario[];
  originalValue: number;
}

export function ScenarioComparison({ scenarios, originalValue }: ScenarioComparisonProps) {
  if (scenarios.length < 2) return null;

  // Collect all unique component names across all scenarios
  const allComponents = new Set<string>();
  for (const scenario of scenarios) {
    for (const diff of scenario.result.component_diffs) {
      allComponents.add(diff.name);
    }
  }
  const componentNames = Array.from(allComponents).sort();

  // Determine which scenario has the best (smallest negative or largest positive) overall delta
  const bestIdx = scenarios.reduce((bestI, s, i, arr) =>
    s.result.delta_value > arr[bestI].result.delta_value ? i : bestI,
    0,
  );

  // For each component row, find the best scenario index
  function bestComponentIdx(componentName: string): number {
    let best = 0;
    let bestDelta = -Infinity;
    for (let i = 0; i < scenarios.length; i++) {
      const diff = scenarios[i].result.component_diffs.find((d) => d.name === componentName);
      const delta = diff?.delta_value ?? 0;
      if (delta > bestDelta) {
        bestDelta = delta;
        best = i;
      }
    }
    return best;
  }

  const formatDelta = (value: number) => {
    const sign = value >= 0 ? '+' : '';
    return `${sign}${value.toFixed(4)}`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Scenario Comparison</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[180px]">Metric</TableHead>
              {scenarios.map((s) => (
                <TableHead key={s.id} className="text-right">
                  {s.name}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {/* Overall value row */}
            <TableRow className="font-medium">
              <TableCell>Modified Value</TableCell>
              {scenarios.map((s, i) => {
                const deltaSign = s.result.delta_value >= 0 ? '+' : '';
                return (
                  <TableCell
                    key={s.id}
                    className={`text-right tabular-nums ${i === bestIdx ? 'text-green-600' : ''}`}
                  >
                    <div>{s.result.modified_value.toFixed(4)}</div>
                    <div className="text-xs text-muted-foreground">
                      {deltaSign}{s.result.delta_value.toFixed(4)} ({deltaSign}{s.result.delta_percentage.toFixed(1)}%)
                    </div>
                  </TableCell>
                );
              })}
            </TableRow>

            {/* Original value reference row */}
            <TableRow>
              <TableCell className="text-muted-foreground">Original Value</TableCell>
              {scenarios.map((s) => (
                <TableCell key={s.id} className="text-right tabular-nums text-muted-foreground">
                  {originalValue.toFixed(4)}
                </TableCell>
              ))}
            </TableRow>

            {/* Separator row */}
            <TableRow>
              <TableCell colSpan={scenarios.length + 1} className="py-1">
                <div className="border-t border-border" />
              </TableCell>
            </TableRow>

            {/* Component rows */}
            {componentNames.map((name) => {
              const bestI = bestComponentIdx(name);
              return (
                <TableRow key={name}>
                  <TableCell className="font-medium text-sm">{name}</TableCell>
                  {scenarios.map((s, i) => {
                    const diff = s.result.component_diffs.find((d) => d.name === name);
                    if (!diff) {
                      return (
                        <TableCell key={s.id} className="text-right text-muted-foreground">
                          --
                        </TableCell>
                      );
                    }
                    const hasDelta = Math.abs(diff.delta_value) > 0.0001;
                    return (
                      <TableCell
                        key={s.id}
                        className={`text-right tabular-nums ${i === bestI && hasDelta ? 'text-green-600' : ''}`}
                      >
                        <div>{diff.modified_value.toFixed(4)}</div>
                        {hasDelta && (
                          <div className="text-xs text-muted-foreground">
                            {formatDelta(diff.delta_value)}
                          </div>
                        )}
                      </TableCell>
                    );
                  })}
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
