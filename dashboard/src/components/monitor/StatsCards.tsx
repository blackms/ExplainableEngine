'use client';

import { useMemo } from 'react';
import { useExplanationList, useStats } from '@/lib/api/hooks';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

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

interface StatCardProps {
  label: string;
  value: string;
  loading?: boolean;
}

function StatCard({ label, value, loading }: StatCardProps) {
  return (
    <Card>
      <CardContent className="py-2">
        <p className="text-sm text-muted-foreground">{label}</p>
        {loading ? (
          <Skeleton className="h-8 w-20 mt-1" />
        ) : (
          <p className="text-2xl font-bold">{value}</p>
        )}
      </CardContent>
    </Card>
  );
}

export function StatsCards() {
  const { data: statsData, isLoading: statsLoading } = useStats();
  const { data: feedData, isLoading: feedLoading } = useExplanationList(
    { limit: 10 },
    { refetchInterval: 10000 }
  );

  const items = feedData?.items ?? [];

  const avgConfidence = useMemo(() => {
    if (items.length === 0) return 0;
    const sum = items.reduce((acc, item) => acc + item.confidence, 0);
    return sum / items.length;
  }, [items]);

  const anomalyCount = useMemo(() => {
    return items.filter(
      (item) => item.confidence < 0.5 || item.missing_impact > 0.2
    ).length;
  }, [items]);

  const latestTimestamp = useMemo(() => {
    if (items.length === 0) return null;
    return items.reduce((latest, item) => {
      const t = item.metadata.created_at;
      return t > latest ? t : latest;
    }, items[0].metadata.created_at);
  }, [items]);

  const isLoading = statsLoading || feedLoading;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard
        label="Total Explanations"
        value={statsData?.total_explanations?.toLocaleString() ?? '0'}
        loading={statsLoading}
      />
      <StatCard
        label="Avg. Confidence"
        value={`${(avgConfidence * 100).toFixed(1)}%`}
        loading={isLoading}
      />
      <StatCard
        label="Anomalies"
        value={String(anomalyCount)}
        loading={feedLoading}
      />
      <StatCard
        label="Latest"
        value={latestTimestamp ? relativeTime(latestTimestamp) : '--'}
        loading={feedLoading}
      />
    </div>
  );
}
