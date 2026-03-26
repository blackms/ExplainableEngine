'use client';

import { Button } from '@/components/ui/button';

interface PaginationControlsProps {
  total: number;
  limit: number;
  nextCursor?: string;
  onNext: () => void;
  onPrev: () => void;
  currentOffset?: number;
}

export function PaginationControls({
  total,
  limit,
  nextCursor,
  onNext,
  onPrev,
  currentOffset = 0,
}: PaginationControlsProps) {
  const start = total === 0 ? 0 : currentOffset + 1;
  const end = Math.min(currentOffset + limit, total);
  const hasPrev = currentOffset > 0;
  const hasNext = !!nextCursor;

  if (total === 0) return null;

  return (
    <div className="flex items-center justify-between pt-4">
      <span className="text-sm text-muted-foreground">
        Showing {start}&ndash;{end} of {total} explanations
      </span>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={!hasPrev}
          onClick={onPrev}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={!hasNext}
          onClick={onNext}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
