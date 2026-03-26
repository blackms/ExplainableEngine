'use client';

import type { ComponentDiff } from '@/lib/api/types';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface DiffTableProps {
  diffs: ComponentDiff[];
}

export function DiffTable({ diffs }: DiffTableProps) {
  const sorted = [...diffs].sort(
    (a, b) => Math.abs(b.delta_value) - Math.abs(a.delta_value),
  );

  const maxAbsDelta = sorted.length > 0 ? Math.abs(sorted[0].delta_value) : 0;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Component</TableHead>
          <TableHead className="text-right">Original</TableHead>
          <TableHead className="text-right">Modified</TableHead>
          <TableHead className="text-right">Delta</TableHead>
          <TableHead className="text-right">Delta %</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {sorted.map((diff) => {
          const isLargest =
            maxAbsDelta > 0 &&
            Math.abs(diff.delta_value) >= maxAbsDelta * 0.8;
          const deltaColor =
            diff.delta_value > 0
              ? 'text-green-600'
              : diff.delta_value < 0
                ? 'text-red-600'
                : '';

          return (
            <TableRow
              key={diff.name}
              className={isLargest ? 'bg-muted/30' : ''}
            >
              <TableCell className="font-medium">{diff.name}</TableCell>
              <TableCell className="text-right tabular-nums">
                {diff.original_value.toFixed(4)}
              </TableCell>
              <TableCell className="text-right tabular-nums">
                {diff.modified_value.toFixed(4)}
              </TableCell>
              <TableCell className={`text-right tabular-nums font-medium ${deltaColor}`}>
                {diff.delta_value >= 0 ? '+' : ''}
                {diff.delta_value.toFixed(4)}
              </TableCell>
              <TableCell className={`text-right tabular-nums ${deltaColor}`}>
                {diff.delta_percentage >= 0 ? '+' : ''}
                {diff.delta_percentage.toFixed(1)}%
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
