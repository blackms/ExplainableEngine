'use client';

import { Badge } from '@/components/ui/badge';

interface AlertBadgeProps {
  count: number;
}

export function AlertBadge({ count }: AlertBadgeProps) {
  if (count === 0) return null;

  return (
    <Badge
      variant="destructive"
      className={count > 0 ? 'animate-pulse' : ''}
    >
      {count}
    </Badge>
  );
}
