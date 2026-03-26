'use client';

import { useRouter } from 'next/navigation';
import Link from 'next/link';
import {
  Plus,
  Search,
  Activity,
  ChevronRight,
  Layers,
} from 'lucide-react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useStats } from '@/lib/api/hooks';
import { useExplanationList } from '@/lib/api/hooks';

function confidenceColor(value: number) {
  if (value >= 0.8) return { text: 'text-emerald-500', bg: 'bg-emerald-500', bgMuted: 'bg-emerald-50', label: 'High confidence' };
  if (value >= 0.5) return { text: 'text-amber-500', bg: 'bg-amber-500', bgMuted: 'bg-amber-50', label: 'Moderate confidence' };
  return { text: 'text-rose-500', bg: 'bg-rose-500', bgMuted: 'bg-rose-50', label: 'Low -- review recommended' };
}

function formatRelativeTime(dateString: string): string {
  const now = Date.now();
  const then = new Date(dateString).getTime();
  const diffMs = now - then;
  const diffSec = Math.floor(diffMs / 1000);

  if (diffSec < 60) return 'just now';
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHour = Math.floor(diffMin / 60);
  if (diffHour < 24) return `${diffHour}h ago`;
  const diffDay = Math.floor(diffHour / 24);
  if (diffDay < 7) return `${diffDay}d ago`;
  return new Date(dateString).toLocaleDateString();
}

function StatCard({
  label,
  value,
  subtitle,
  loading,
}: {
  label: string;
  value: string;
  subtitle?: string;
  loading?: boolean;
}) {
  return (
    <Card className="transition-shadow duration-150 hover:shadow-sm">
      <CardContent className="space-y-2 p-5">
        <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
          {label}
        </p>
        {loading ? (
          <div className="h-10 w-20 animate-pulse rounded bg-muted" />
        ) : (
          <p className="text-4xl font-bold tabular-nums text-foreground">
            {value}
          </p>
        )}
        {subtitle && (
          <p className="text-xs text-muted-foreground">{subtitle}</p>
        )}
      </CardContent>
    </Card>
  );
}

function ConfidenceDot({ value }: { value: number }) {
  const conf = confidenceColor(value);
  return (
    <span className="flex items-center gap-1.5">
      <span className={`inline-block size-2 rounded-full ${conf.bg}`} />
      <span className={`text-sm tabular-nums ${conf.text}`}>
        {(value * 100).toFixed(0)}%
      </span>
    </span>
  );
}

export default function HomePage() {
  const router = useRouter();
  const { data: stats, isLoading: statsLoading } = useStats();
  const { data: recentData, isLoading: recentLoading } = useExplanationList({ limit: 5 });

  const recentItems = recentData?.items ?? [];
  const totalExplanations = stats?.total_explanations ?? 0;

  // Compute average confidence from recent items
  const avgConfidence =
    recentItems.length > 0
      ? recentItems.reduce((sum, item) => sum + item.confidence, 0) / recentItems.length
      : 0;

  // Find the most recent item's timestamp
  const lastActivity =
    recentItems.length > 0 ? recentItems[0].metadata.created_at : null;

  return (
    <div className="space-y-6">
      {/* Welcome header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Welcome back</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Here&apos;s what&apos;s happening with your explanations.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Link
            href="/audit"
            className="inline-flex h-7 items-center gap-1 rounded-[min(var(--radius-md),12px)] border border-border bg-background px-2.5 text-[0.8rem] font-medium text-foreground transition-colors hover:bg-muted"
          >
            <Search className="size-3.5" />
            Search
          </Link>
          <Link
            href="/audit"
            className="inline-flex h-7 items-center gap-1 rounded-[min(var(--radius-md),12px)] bg-primary px-2.5 text-[0.8rem] font-medium text-primary-foreground transition-colors hover:bg-primary/80"
          >
            <Plus className="size-3.5" />
            New Explanation
          </Link>
        </div>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="Total Explanations"
          value={statsLoading ? '--' : String(totalExplanations)}
          loading={statsLoading}
        />
        <StatCard
          label="Avg Confidence"
          value={
            recentLoading
              ? '--'
              : recentItems.length > 0
                ? `${(avgConfidence * 100).toFixed(0)}%`
                : '--'
          }
          subtitle={
            recentItems.length > 0
              ? confidenceColor(avgConfidence).label
              : undefined
          }
          loading={recentLoading}
        />
        <StatCard
          label="Active Alerts"
          value="0"
          subtitle="No active alerts"
        />
        <StatCard
          label="Last Activity"
          value={
            lastActivity ? formatRelativeTime(lastActivity) : '--'
          }
          loading={recentLoading}
        />
      </div>

      {/* Recent explanations */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="text-xl font-semibold">
            Recent Explanations
          </CardTitle>
          <Link
            href="/audit"
            className="text-sm text-primary hover:underline"
          >
            View all
          </Link>
        </CardHeader>
        <CardContent>
          {recentLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="flex items-center justify-between py-3">
                  <div className="h-4 w-32 animate-pulse rounded bg-muted" />
                  <div className="h-4 w-20 animate-pulse rounded bg-muted" />
                </div>
              ))}
            </div>
          ) : recentItems.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <Layers className="mb-3 size-10 text-muted-foreground/50" />
              <p className="text-sm font-medium text-foreground">
                No explanations yet
              </p>
              <p className="mt-1 text-xs text-muted-foreground">
                Create your first explanation to see it here.
              </p>
              <Link
                href="/audit"
                className="mt-4 inline-flex h-7 items-center gap-1 rounded-[min(var(--radius-md),12px)] bg-primary px-2.5 text-[0.8rem] font-medium text-primary-foreground transition-colors hover:bg-primary/80"
              >
                <Plus className="size-3.5" />
                Create Explanation
              </Link>
            </div>
          ) : (
            <div className="divide-y divide-border">
              {recentItems.map((item) => (
                <button
                  key={item.id}
                  type="button"
                  onClick={() => router.push(`/explain/${item.id}`)}
                  className="flex w-full items-center justify-between rounded-md px-4 py-3 text-left transition-colors hover:bg-accent/50"
                >
                  <div className="flex items-center gap-3">
                    <span className="truncate text-sm font-medium font-mono">
                      {item.target}
                    </span>
                    <span className="text-sm tabular-nums text-muted-foreground">
                      {item.final_value.toFixed(2)}
                    </span>
                  </div>
                  <div className="flex items-center gap-4">
                    <ConfidenceDot value={item.confidence} />
                    <span className="w-16 text-right text-xs text-muted-foreground">
                      {formatRelativeTime(item.metadata.created_at)}
                    </span>
                    <ChevronRight className="size-4 text-muted-foreground" />
                  </div>
                </button>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Quick actions */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Link href="/audit" className="group">
          <Card className="h-full transition-all duration-150 group-hover:shadow-sm group-hover:ring-1 group-hover:ring-primary/20">
            <CardContent className="flex items-start gap-3 p-5">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Plus className="size-5" />
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  New Explanation
                </p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  Submit components to generate an explainability breakdown
                </p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/audit" className="group">
          <Card className="h-full transition-all duration-150 group-hover:shadow-sm group-hover:ring-1 group-hover:ring-primary/20">
            <CardContent className="flex items-start gap-3 p-5">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Search className="size-5" />
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  Browse Archive
                </p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  Search and filter through past explanations
                </p>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/monitor" className="group">
          <Card className="h-full transition-all duration-150 group-hover:shadow-sm group-hover:ring-1 group-hover:ring-primary/20">
            <CardContent className="flex items-start gap-3 p-5">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Activity className="size-5" />
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  Live Monitor
                </p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  Watch explanations and alerts in real-time
                </p>
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
}
