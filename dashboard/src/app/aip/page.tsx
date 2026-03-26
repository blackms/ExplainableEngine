'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  TrendingUp,
  Search,
  Newspaper,
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
import { useAIPMarketMood } from '@/lib/api/hooks';
import Link from 'next/link';

const POPULAR_TICKERS = ['AAPL', 'MSFT', 'NVDA', 'TSLA', 'AMZN', 'GOOGL', 'META', 'JPM'];

function trendIndicator(trend: number) {
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

function formatSentiment(value: number) {
  return (value >= 0 ? '+' : '') + value.toFixed(3);
}

export default function AIPPage() {
  const router = useRouter();
  const [tickerInput, setTickerInput] = useState('');
  const {
    data: moodData,
    isLoading: moodLoading,
    error: moodError,
    refetch: refetchMood,
  } = useAIPMarketMood();

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    const ticker = tickerInput.trim().toUpperCase();
    if (ticker) {
      router.push(`/aip/${ticker}`);
    }
  }

  const mood = moodData?.market_mood;
  const trend = mood ? trendIndicator(mood.overall_trend) : null;
  const TrendIcon = trend?.icon ?? ArrowRight;

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div>
        <div className="flex items-center gap-2">
          <TrendingUp className="size-6 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">
            AIP Sentiment Intelligence
          </h1>
        </div>
        <p className="mt-1 text-sm text-muted-foreground">
          AI-powered sentiment analysis for financial tickers, backed by
          explainable breakdowns.
        </p>
      </div>

      {/* Ticker search */}
      <form onSubmit={handleSearch} className="flex items-center gap-2">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <input
            type="text"
            value={tickerInput}
            onChange={(e) => setTickerInput(e.target.value)}
            placeholder="Search ticker (e.g. NVDA)"
            className="h-10 w-full rounded-md border border-input bg-background pl-9 pr-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
          />
        </div>
        <button
          type="submit"
          className="inline-flex h-10 items-center gap-1.5 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80"
        >
          Explain
        </button>
      </form>

      {/* Market Mood */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Newspaper className="size-5" />
            Market Mood
          </CardTitle>
        </CardHeader>
        <CardContent>
          {moodLoading ? (
            <div className="space-y-3">
              <div className="flex items-center gap-4">
                <Skeleton className="h-8 w-28" />
                <Skeleton className="h-6 w-36" />
              </div>
              <Skeleton className="h-5 w-40" />
            </div>
          ) : moodError ? (
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <AlertCircle className="size-5 text-rose-500 shrink-0" />
              <p>Could not fetch market mood data.</p>
              <button
                type="button"
                onClick={() => refetchMood()}
                className="inline-flex items-center gap-1 text-primary hover:underline"
              >
                <RefreshCw className="size-3.5" />
                Retry
              </button>
            </div>
          ) : mood ? (
            <div className="space-y-2">
              <div className="flex flex-wrap items-baseline gap-4">
                <span className="text-2xl font-bold tabular-nums">
                  Overall: {formatSentiment(mood.overall_sentiment)}
                </span>
                <span
                  className={`flex items-center gap-1 text-sm font-medium ${trend?.color}`}
                >
                  <TrendIcon className="size-4" />
                  Trend: {formatSentiment(mood.overall_trend)} ({trend?.label})
                </span>
              </div>
              <p className="text-sm text-muted-foreground">
                {mood.total_articles.toLocaleString()} articles analyzed
              </p>
              {moodData?.explanation?.id && (
                <Link
                  href={`/explain/${moodData.explanation.id}`}
                  className="inline-flex items-center gap-1 text-sm text-primary hover:underline"
                >
                  View Full Breakdown
                  <ArrowRight className="size-3.5" />
                </Link>
              )}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">
              AIP integration requires an API key. Contact your administrator.
            </p>
          )}
        </CardContent>
      </Card>

      {/* Quick Explains */}
      <div>
        <h2 className="mb-3 text-lg font-semibold">Quick Explains</h2>
        <div className="flex flex-wrap gap-2">
          {POPULAR_TICKERS.map((ticker) => (
            <Link
              key={ticker}
              href={`/aip/${ticker}`}
              className="inline-flex h-10 items-center justify-center rounded-md border border-input bg-background px-4 text-sm font-semibold transition-colors hover:bg-accent hover:text-accent-foreground"
            >
              {ticker}
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
