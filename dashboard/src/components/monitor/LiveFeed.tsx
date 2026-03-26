'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useExplanationList } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import type { ExplainResponse } from '@/lib/api/types';

function confidenceVariant(confidence: number) {
  if (confidence >= 0.8) return 'default';
  if (confidence >= 0.5) return 'secondary';
  return 'destructive';
}

function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

function FeedSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 5 }, (_, i) => (
        <Skeleton key={i} className="h-20 w-full" />
      ))}
    </div>
  );
}

interface FeedItemProps {
  item: ExplainResponse;
  isNew: boolean;
  onClick: () => void;
}

function FeedItem({ item, isNew, onClick }: FeedItemProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`w-full text-left transition-all ${
        isNew ? 'animate-pulse ring-2 ring-primary/50 rounded-xl' : ''
      }`}
    >
      <Card size="sm" className="hover:ring-2 hover:ring-primary/30 cursor-pointer transition-all">
        <CardContent className="flex items-center justify-between gap-4">
          <div className="min-w-0 flex-1">
            <p className="font-medium truncate">{item.target}</p>
            <p className="text-xs text-muted-foreground">
              Value: {item.final_value.toFixed(2)}
            </p>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <Badge variant={confidenceVariant(item.confidence)}>
              {(item.confidence * 100).toFixed(0)}%
            </Badge>
            <span className="text-xs text-muted-foreground whitespace-nowrap">
              {relativeTime(item.metadata.created_at)}
            </span>
          </div>
        </CardContent>
      </Card>
    </button>
  );
}

export function LiveFeed() {
  const router = useRouter();
  const [paused, setPaused] = useState(false);
  const previousIdsRef = useRef<Set<string>>(new Set());
  const [newIds, setNewIds] = useState<Set<string>>(new Set());

  const { data, isLoading, error } = useExplanationList(
    { limit: 10 },
    {
      refetchInterval: paused ? false : 10000,
    }
  );

  useEffect(() => {
    if (!data?.items) return;

    const currentIds = new Set(data.items.map((item) => item.id));
    const prev = previousIdsRef.current;

    // On first load, don't highlight anything
    if (prev.size > 0) {
      const freshIds = new Set<string>();
      for (const id of currentIds) {
        if (!prev.has(id)) {
          freshIds.add(id);
        }
      }
      if (freshIds.size > 0) {
        setNewIds(freshIds);
        // Clear highlight after 3 seconds
        const timer = setTimeout(() => setNewIds(new Set()), 3000);
        return () => clearTimeout(timer);
      }
    }

    previousIdsRef.current = currentIds;
  }, [data?.items]);

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Live Feed</CardTitle>
        <Button
          variant={paused ? 'default' : 'outline'}
          size="sm"
          onClick={() => setPaused((p) => !p)}
        >
          {paused ? 'Resume' : 'Pause'}
        </Button>
      </CardHeader>
      <CardContent>
        {error ? (
          <div className="rounded-lg border border-destructive/50 bg-destructive/5 p-4 text-sm text-destructive">
            Failed to load feed: {(error as Error).message}
          </div>
        ) : isLoading ? (
          <FeedSkeleton />
        ) : !data?.items.length ? (
          <p className="text-sm text-muted-foreground text-center py-8">
            No explanations yet.
          </p>
        ) : (
          <div className="space-y-3">
            {data.items.map((item) => (
              <FeedItem
                key={item.id}
                item={item}
                isNew={newIds.has(item.id)}
                onClick={() => router.push(`/explain/${item.id}`)}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
