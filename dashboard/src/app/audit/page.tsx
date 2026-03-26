'use client';

import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useExplanationList } from '@/lib/api/hooks';
import { FilterPanel } from '@/components/audit/FilterPanel';
import { ExplanationTable } from '@/components/audit/ExplanationTable';
import { ExportButton } from '@/components/audit/ExportButton';
import { PaginationControls } from '@/components/audit/PaginationControls';
import { Skeleton } from '@/components/ui/skeleton';
import type { ListOptions } from '@/lib/api/types';

function AuditSkeleton() {
  return (
    <div className="space-y-4">
      <Skeleton className="h-28 w-full" />
      <Skeleton className="h-8 w-full" />
      {Array.from({ length: 5 }, (_, i) => (
        <Skeleton key={i} className="h-12 w-full" />
      ))}
      <Skeleton className="h-10 w-full" />
    </div>
  );
}

export default function AuditPage() {
  const router = useRouter();
  const [filters, setFilters] = useState<ListOptions>({ limit: 20 });
  const [pageOffsets, setPageOffsets] = useState<string[]>([]);
  const { data, isLoading, error } = useExplanationList(filters);

  const currentOffset = pageOffsets.length * (filters.limit ?? 20);

  const handleNext = useCallback(() => {
    if (!data?.next_cursor) return;
    setPageOffsets((prev) => [...prev, data.next_cursor!]);
    setFilters((f) => ({ ...f, cursor: data.next_cursor }));
  }, [data?.next_cursor]);

  const handlePrev = useCallback(() => {
    setPageOffsets((prev) => {
      const next = prev.slice(0, -1);
      const cursor = next.length > 0 ? next[next.length - 1] : undefined;
      setFilters((f) => ({ ...f, cursor }));
      return next;
    });
  }, []);

  const handleFiltersChange = useCallback((next: ListOptions) => {
    setPageOffsets([]);
    setFilters(next);
  }, []);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Audit Trail</h1>
        {data?.items && data.items.length > 0 && (
          <ExportButton items={data.items} />
        )}
      </div>

      <FilterPanel
        filters={filters}
        onChange={handleFiltersChange}
        total={data?.total ?? 0}
      />

      {error ? (
        <div className="rounded-lg border border-destructive/50 bg-destructive/5 p-4 text-sm text-destructive">
          Failed to load explanations: {(error as Error).message}
        </div>
      ) : isLoading ? (
        <AuditSkeleton />
      ) : (
        <>
          <ExplanationTable
            items={data?.items ?? []}
            onRowClick={(id) => router.push(`/explain/${id}`)}
          />
          <PaginationControls
            total={data?.total ?? 0}
            limit={filters.limit ?? 20}
            nextCursor={data?.next_cursor}
            currentOffset={currentOffset}
            onNext={handleNext}
            onPrev={handlePrev}
          />
        </>
      )}
    </div>
  );
}
