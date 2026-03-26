'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Activity, Pause, Play } from 'lucide-react';
import { useExplanationList } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import type { ExplainResponse } from '@/lib/api/types';

function confidenceDotColor(confidence: number): string {
  if (confidence >= 0.8) return 'bg-emerald-500';
  if (confidence >= 0.5) return 'bg-amber-500';
  return 'bg-rose-500';
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
        <Skeleton key={i} className="h-16 w-full rounded-lg" />
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
      className="w-full text-left"
    >
      <div
        className="flex items-center justify-between gap-4 rounded-lg border p-3 cursor-pointer transition-all hover:bg-accent/50"
        style={
          isNew
            ? {
                animation: 'feedHighlight 2s ease-out forwards',
              }
            : undefined
        }
      >
        <div className="min-w-0 flex-1">
          <p className="font-mono text-sm font-medium truncate">
            {item.target}
          </p>
          <p className="text-xs text-muted-foreground tabular-nums">
            {item.final_value.toFixed(2)}
          </p>
        </div>
        <div className="flex items-center gap-3 shrink-0">
          <span className="inline-flex items-center gap-1.5">
            <span
              className={`h-2 w-2 rounded-full ${confidenceDotColor(item.confidence)}`}
            />
            <span className="text-sm tabular-nums">
              {(item.confidence * 100).toFixed(0)}%
            </span>
          </span>
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {relativeTime(item.metadata.created_at)}
          </span>
        </div>
      </div>
    </button>
  );
}

const feedHighlightKeyframes = `
@keyframes feedHighlight {
  0% { background-color: color-mix(in srgb, var(--primary) 10%, transparent); }
  100% { background-color: transparent; }
}
`;

export function LiveFeed() {
  const router = useRouter();
  const [paused, setPaused] = useState(false);
  const previousIdsRef = useRef<Set<string>>(new Set());
  const [newIds, setNewIds] = useState<Set<string>>(new Set());

  const { data, isLoading, error } = useExplanationList(
    { limit: 10 },
    {
      refetchInterval: paused ? false : 10000,
    },
  );

  useEffect(() => {
    if (!data?.items) return;

    const currentIds = new Set(data.items.map((item) => item.id));
    const prev = previousIdsRef.current;

    if (prev.size > 0) {
      const freshIds = new Set<string>();
      for (const id of currentIds) {
        if (!prev.has(id)) {
          freshIds.add(id);
        }
      }
      if (freshIds.size > 0) {
        setNewIds(freshIds);
        const timer = setTimeout(() => setNewIds(new Set()), 2000);
        return () => clearTimeout(timer);
      }
    }

    previousIdsRef.current = currentIds;
  }, [data?.items]);

  return (
    <Card>
      <style dangerouslySetInnerHTML={{ __html: feedHighlightKeyframes }} />
      <CardHeader>
        <CardTitle>Live Feed</CardTitle>
        <CardAction>
          <div className="flex items-center gap-2">
            {!paused && (
              <span className="inline-flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400">
                <span className="relative flex h-2 w-2">
                  <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
                  <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
                </span>
                Live
              </span>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPaused((p) => !p)}
            >
              {paused ? (
                <>
                  <Play className="size-3.5" />
                  Resume
                </>
              ) : (
                <>
                  <Pause className="size-3.5" />
                  Pause
                </>
              )}
            </Button>
          </div>
        </CardAction>
      </CardHeader>
      <CardContent>
        {error ? (
          <div className="rounded-lg border border-destructive/50 bg-destructive/5 p-4 text-sm text-destructive">
            Failed to load feed: {(error as Error).message}
          </div>
        ) : isLoading ? (
          <FeedSkeleton />
        ) : !data?.items.length ? (
          <div className="flex flex-col items-center justify-center py-16 space-y-4">
            <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
              <Activity className="h-6 w-6 text-muted-foreground" />
            </div>
            <div className="text-center space-y-1">
              <h3 className="text-base font-medium">No recent activity</h3>
              <p className="text-sm text-muted-foreground max-w-sm">
                Explanations will appear here as they are created.
              </p>
            </div>
          </div>
        ) : (
          <div className="space-y-2">
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
