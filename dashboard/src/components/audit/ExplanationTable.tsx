'use client';

import { useState, useMemo } from 'react';
import { ArrowDown, ArrowUp, ArrowUpDown, ChevronRight, FileSearch } from 'lucide-react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Button } from '@/components/ui/button';
import type { ExplainResponse } from '@/lib/api/types';

interface ExplanationTableProps {
  items: ExplainResponse[];
  onRowClick: (id: string) => void;
}

type SortKey = 'target' | 'final_value' | 'confidence' | 'created_at';
type SortDir = 'asc' | 'desc';

function confidenceDotColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-emerald-500';
  if (confidence >= 0.5) return 'bg-amber-500';
  return 'bg-rose-500';
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

function formatFullDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZoneName: 'short',
  });
}

function SortIcon({
  column,
  sortKey,
  sortDir,
}: {
  column: SortKey;
  sortKey: SortKey;
  sortDir: SortDir;
}) {
  if (column !== sortKey)
    return <ArrowUpDown className="ml-1 inline size-3 opacity-40" />;
  return sortDir === 'asc' ? (
    <ArrowUp className="ml-1 inline size-3" />
  ) : (
    <ArrowDown className="ml-1 inline size-3" />
  );
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
      <div className="flex flex-col items-center justify-center py-16 space-y-4">
        <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
          <FileSearch className="h-6 w-6 text-muted-foreground" />
        </div>
        <div className="text-center space-y-1">
          <h3 className="text-base font-medium">No results found</h3>
          <p className="text-sm text-muted-foreground max-w-sm">
            Try adjusting your filters or search terms.
          </p>
        </div>
        <Button size="sm" variant="outline" onClick={() => onRowClick('')}>
          Clear filters
        </Button>
      </div>
    );
  }

  const columns: [SortKey, string][] = [
    ['target', 'Name'],
    ['final_value', 'Value'],
    ['confidence', 'Confidence'],
    ['created_at', 'Created'],
  ];

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            {columns.map(([key, label]) => (
              <TableHead
                key={key}
                className="cursor-pointer select-none text-xs font-medium text-muted-foreground uppercase tracking-wide"
                onClick={() => handleSort(key)}
              >
                {label}
                <SortIcon column={key} sortKey={sortKey} sortDir={sortDir} />
              </TableHead>
            ))}
            <TableHead className="w-20 text-right text-xs font-medium text-muted-foreground uppercase tracking-wide">
              Components
            </TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {sorted.map((item) => (
            <TableRow
              key={item.id}
              className="cursor-pointer hover:bg-accent/50 transition-colors focus-visible:bg-accent/50 focus-visible:outline-none"
              onClick={() => onRowClick(item.id)}
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === 'Enter') onRowClick(item.id);
              }}
            >
              <TableCell className="font-mono text-sm font-medium">
                {item.target}
              </TableCell>
              <TableCell className="tabular-nums text-sm text-right w-24">
                {item.final_value.toFixed(2)}
              </TableCell>
              <TableCell className="w-28">
                <span className="inline-flex items-center gap-2">
                  <span
                    className={`h-2 w-2 rounded-full ${confidenceDotColor(item.confidence)}`}
                  />
                  <span className="text-sm tabular-nums">
                    {(item.confidence * 100).toFixed(0)}%
                  </span>
                </span>
              </TableCell>
              <TableCell className="w-32 text-right">
                <Tooltip>
                  <TooltipTrigger className="text-sm text-muted-foreground cursor-default">
                    {relativeTime(item.metadata.created_at)}
                  </TooltipTrigger>
                  <TooltipContent>
                    {formatFullDate(item.metadata.created_at)}
                  </TooltipContent>
                </Tooltip>
              </TableCell>
              <TableCell className="w-20 text-right text-sm text-muted-foreground">
                {item.breakdown?.length ?? 0}
              </TableCell>
              <TableCell className="w-10 text-center">
                <ChevronRight className="size-4 text-muted-foreground" />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TooltipProvider>
  );
}
