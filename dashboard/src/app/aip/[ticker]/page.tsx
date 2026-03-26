'use client';

import { use } from 'react';
import Link from 'next/link';
import {
  ArrowUpRight,
  ArrowDownRight,
  ArrowRight,
  AlertCircle,
  RefreshCw,
} from 'lucide-react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { BreakdownChart } from '@/components/explanation/BreakdownChart';
import { DriverRanking } from '@/components/explanation/DriverRanking';
import { ConfidencePanel } from '@/components/explanation/ConfidencePanel';
import { useAIPExplainTicker } from '@/lib/api/hooks';

function sentimentStyle(label: string) {
  const lower = label.toLowerCase();
  if (lower === 'bullish')
    return {
      text: 'text-emerald-600 dark:text-emerald-400',
      badge: 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300',
    };
  if (lower === 'bearish')
    return {
      text: 'text-rose-600 dark:text-rose-400',
      badge: 'bg-rose-500/15 text-rose-700 dark:text-rose-300',
    };
  return {
    text: 'text-muted-foreground',
    badge: 'bg-muted text-muted-foreground',
  };
}

function trendInfo(trend: number) {
  if (trend > 0.01)
    return {
      icon: ArrowUpRight,
      color: 'text-emerald-500',
      label: 'rising',
    };
  if (trend < -0.01)
    return {
      icon: ArrowDownRight,
      color: 'text-rose-500',
      label: 'declining',
    };
  return { icon: ArrowRight, color: 'text-muted-foreground', label: 'flat' };
}

function formatSigned(value: number) {
  return (value >= 0 ? '+' : '') + value.toFixed(3);
}

function TickerSkeleton() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Skeleton className="h-5 w-12" />
        <Skeleton className="h-8 w-56" />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardContent className="pt-4 space-y-2">
            <Skeleton className="h-10 w-28" />
            <Skeleton className="h-6 w-20" />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-28" />
          </CardHeader>
          <CardContent className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-4 w-full" />
            ))}
          </CardContent>
        </Card>
      </div>
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-48" />
        </CardHeader>
        <CardContent className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="flex items-center gap-3">
              <Skeleton className="h-4 w-[120px]" />
              <Skeleton
                className="h-8 rounded-full"
                style={{ width: `${80 - i * 20}%` }}
              />
            </div>
          ))}
        </CardContent>
      </Card>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-28" />
          </CardHeader>
          <CardContent className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="h-6 w-6 rounded-full" />
                <div className="flex-1 space-y-1">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-2 w-full rounded-full" />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-28" />
          </CardHeader>
          <CardContent className="space-y-3">
            <Skeleton className="h-12 w-full rounded-lg" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default function AIPTickerPage({
  params,
}: {
  params: Promise<{ ticker: string }>;
}) {
  const { ticker } = use(params);
  const tickerUpper = ticker.toUpperCase();
  const { data, isLoading, error, refetch } = useAIPExplainTicker(tickerUpper);

  if (isLoading) {
    return <TickerSkeleton />;
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <div className="text-center space-y-3">
          <AlertCircle className="mx-auto size-10 text-rose-500" />
          <h2 className="text-lg font-semibold">
            Could not fetch sentiment for {tickerUpper}
          </h2>
          <p className="text-sm text-muted-foreground">
            {error instanceof Error ? error.message : 'An unexpected error occurred.'}
          </p>
          <button
            type="button"
            onClick={() => refetch()}
            className="inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
          >
            <RefreshCw className="size-3.5" />
            Retry
          </button>
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <p className="text-sm text-muted-foreground">
          AIP integration requires an API key. Contact your administrator.
        </p>
      </div>
    );
  }

  const { aip_data, explanation } = data;
  const style = sentimentStyle(aip_data.sentiment_label);
  const trend = trendInfo(aip_data.trend);
  const TrendIcon = trend.icon;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <Link
            href="/aip"
            className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            &larr; Back
          </Link>
          <h1 className="text-2xl font-bold">
            {tickerUpper} Sentiment Analysis
          </h1>
        </div>
        {explanation?.id && (
          <Link
            href={`/explain/${explanation.id}`}
            className="inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm font-medium hover:bg-accent transition-colors"
          >
            Full Explanation
          </Link>
        )}
      </div>

      {/* Sentiment overview + AIP raw data */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Sentiment overview card */}
        <Card>
          <CardHeader>
            <CardTitle>Sentiment</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className={`text-4xl font-bold tabular-nums ${style.text}`}>
              {formatSigned(aip_data.sentiment_7d)}
            </div>
            <span
              className={`inline-block rounded-full px-2.5 py-0.5 text-sm font-medium ${style.badge}`}
            >
              {aip_data.sentiment_label}
            </span>
            <div
              className={`flex items-center gap-1 text-sm font-medium ${trend.color}`}
            >
              <TrendIcon className="size-4" />
              {formatSigned(aip_data.trend)} ({trend.label})
            </div>
            <p className="text-xs text-muted-foreground">
              Scaled score: {aip_data.sentiment_score_scaled}
            </p>
          </CardContent>
        </Card>

        {/* AIP raw data card */}
        <Card>
          <CardHeader>
            <CardTitle>AIP Raw Data</CardTitle>
          </CardHeader>
          <CardContent>
            <dl className="space-y-2 text-sm">
              <div className="flex items-center justify-between">
                <dt className="text-muted-foreground">7d Sentiment</dt>
                <dd className="font-medium tabular-nums">
                  {formatSigned(aip_data.sentiment_7d)}
                </dd>
              </div>
              <div className="flex items-center justify-between">
                <dt className="text-muted-foreground">30d Sentiment</dt>
                <dd className="font-medium tabular-nums">
                  {formatSigned(aip_data.sentiment_30d)}
                </dd>
              </div>
              <div className="flex items-center justify-between">
                <dt className="text-muted-foreground">Trend</dt>
                <dd
                  className={`flex items-center gap-1 font-medium tabular-nums ${trend.color}`}
                >
                  <TrendIcon className="size-3.5" />
                  {formatSigned(aip_data.trend)} ({trend.label})
                </dd>
              </div>
              <div className="flex items-center justify-between">
                <dt className="text-muted-foreground">Articles (7d)</dt>
                <dd className="font-medium tabular-nums">
                  {aip_data.article_count_7d}
                </dd>
              </div>
              <div className="flex items-center justify-between">
                <dt className="text-muted-foreground">Positive ratio</dt>
                <dd className="font-medium tabular-nums">
                  {(aip_data.positive_ratio * 100).toFixed(0)}%
                </dd>
              </div>
            </dl>
          </CardContent>
        </Card>
      </div>

      {/* Explanation Breakdown */}
      {explanation?.breakdown && explanation.breakdown.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Explanation Breakdown</CardTitle>
          </CardHeader>
          <CardContent>
            <BreakdownChart breakdown={explanation.breakdown} />
          </CardContent>
        </Card>
      )}

      {/* Drivers + Confidence */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {explanation?.top_drivers && explanation.top_drivers.length > 0 && (
          <DriverRanking drivers={explanation.top_drivers} />
        )}
        {explanation?.confidence_detail && (
          <ConfidencePanel
            overall={explanation.confidence}
            perNode={explanation.confidence_detail.per_node ?? {}}
            missingImpact={explanation.missing_impact}
          />
        )}
      </div>

      {/* AI Analysis link */}
      {explanation?.id && (
        <Card>
          <CardContent className="flex items-center justify-between p-4">
            <div>
              <p className="text-sm font-medium">AI Analysis</p>
              <p className="text-xs text-muted-foreground">
                Ask AI about this sentiment and get deeper insights.
              </p>
            </div>
            <Link
              href={`/explain/${explanation.id}`}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
            >
              Ask AI about this sentiment
            </Link>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
