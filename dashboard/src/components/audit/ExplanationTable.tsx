'use client';

import { useState, useMemo } from 'react';
import { AlertTriangle, ArrowDown, ArrowUp, ArrowUpDown } from 'lucide-react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import type { ExplainResponse } from '@/lib/api/types';

interface ExplanationTableProps {
  items: ExplainResponse[];
  onRowClick: (id: string) => void;
}

type SortKey = 'target' | 'final_value' | 'confidence' | 'missing_impact' | 'created_at';
type SortDir = 'asc' | 'desc';

function confidenceBadge(confidence: number) {
  if (confidence >= 0.8) {
    return <Badge className="bg-green-600/15 text-green-700 dark:text-green-400">{(confidence * 100).toFixed(0)}%</Badge>;
  }
  if (confidence >= 0.5) {
    return <Badge className="bg-yellow-500/15 text-yellow-700 dark:text-yellow-400">{(confidence * 100).toFixed(0)}%</Badge>;
  }
  return <Badge className="bg-red-500/15 text-red-700 dark:text-red-400">{(confidence * 100).toFixed(0)}%</Badge>;
}

function relativeTime(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diff = now - then;
  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return 'just now';
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  const months = Math.floor(days / 30);
  return `${months}mo ago`;
}

function SortIcon({ column, sortKey, sortDir }: { column: SortKey; sortKey: SortKey; sortDir: SortDir }) {
  if (column !== sortKey) return <ArrowUpDown className="ml-1 inline size-3 opacity-40" />;
  return sortDir === 'asc'
    ? <ArrowUp className="ml-1 inline size-3" />
    : <ArrowDown className="ml-1 inline size-3" />;
}

export function ExplanationTable({ items, onRowClick }: ExplanationTableProps) {
  const [sortKey, setSortKey] = useState<SortKey>('created_at');
  const [sortDir, setSortDir] = useState<SortDir>('desc');

  const handleSort = (key: SortKey) => {
    if (key === sortKey) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortKey(key);
      setSortDir('desc');
    }
  };

  const sorted = useMemo(() => {
    const copy = [...items];
    copy.sort((a, b) => {
      let av: number | string;
      let bv: number | string;
      switch (sortKey) {
        case 'target':
          av = a.target;
          bv = b.target;
          break;
        case 'final_value':
          av = a.final_value;
          bv = b.final_value;
          break;
        case 'confidence':
          av = a.confidence;
          bv = b.confidence;
          break;
        case 'missing_impact':
          av = a.missing_impact;
          bv = b.missing_impact;
          break;
        case 'created_at':
          av = a.metadata.created_at;
          bv = b.metadata.created_at;
          break;
        default:
          return 0;
      }
      if (av < bv) return sortDir === 'asc' ? -1 : 1;
      if (av > bv) return sortDir === 'asc' ? 1 : -1;
      return 0;
    });
    return copy;
  }, [items, sortKey, sortDir]);

  if (items.length === 0) {
    return (
      <div className="flex min-h-40 items-center justify-center rounded-lg border border-dashed text-sm text-muted-foreground">
        No explanations found.
      </div>
    );
  }

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow>
            {(
              [
                ['target', 'Target'],
                ['final_value', 'Value'],
                ['confidence', 'Confidence'],
                ['missing_impact', 'Missing Impact'],
                ['created_at', 'Created At'],
              ] as [SortKey, string][]
            ).map(([key, label]) => (
              <TableHead
                key={key}
                className="cursor-pointer select-none"
                onClick={() => handleSort(key)}
              >
                {label}
                <SortIcon column={key} sortKey={sortKey} sortDir={sortDir} />
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {sorted.map((item) => (
            <TableRow
              key={item.id}
              className="cursor-pointer"
              onClick={() => onRowClick(item.id)}
            >
              <TableCell className="font-medium">{item.target}</TableCell>
              <TableCell>{item.final_value.toFixed(2)}</TableCell>
              <TableCell>{confidenceBadge(item.confidence)}</TableCell>
              <TableCell>
                {item.missing_impact > 0 ? (
                  <span className="inline-flex items-center gap-1 text-yellow-600 dark:text-yellow-400">
                    <AlertTriangle className="size-3.5" />
                    {(item.missing_impact * 100).toFixed(1)}%
                  </span>
                ) : (
                  <span className="text-muted-foreground">&mdash;</span>
                )}
              </TableCell>
              <TableCell>
                <Tooltip>
                  <TooltipTrigger className="text-left">
                    {relativeTime(item.metadata.created_at)}
                  </TooltipTrigger>
                  <TooltipContent>
                    {new Date(item.metadata.created_at).toLocaleString()}
                  </TooltipContent>
                </Tooltip>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TooltipProvider>
  );
}
