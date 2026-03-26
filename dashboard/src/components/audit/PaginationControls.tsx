'use client';

import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface PaginationControlsProps {
  total: number;
  limit: number;
  nextCursor?: string;
  onNext: () => void;
  onPrev: () => void;
  /** Current page offset (number of items already viewed). Defaults to 0. */
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

  return (
    <div className="flex items-center justify-between border-t pt-3">
      <p className="text-sm text-muted-foreground">
        {total === 0 ? 'No results' : `Showing ${start}\u2013${end} of ${total}`}
      </p>
      <div className="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={!hasPrev}
          onClick={onPrev}
        >
          <ChevronLeft className="size-4" />
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={!hasNext}
          onClick={onNext}
        >
          Next
          <ChevronRight className="size-4" />
        </Button>
      </div>
    </div>
  );
}
