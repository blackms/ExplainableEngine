'use client';
import { use } from 'react';
import Link from 'next/link';
import { useAIPExplainTicker } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { BreakdownChart } from '@/components/explanation/BreakdownChart';
import { DriverRanking } from '@/components/explanation/DriverRanking';
import { ConfidenceGauge } from '@/components/explanation/ConfidenceGauge';

export default function AIPTickerPage({ params }: { params: Promise<{ ticker: string }> }) {
  const { ticker } = use(params);
  const { data, isLoading, error } = useAIPExplainTicker(ticker);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-64" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Skeleton className="h-40" />
          <Skeleton className="h-40" />
        </div>
        <Skeleton className="h-60" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="space-y-4">
        <Link href="/aip" className="text-sm text-muted-foreground hover:text-foreground">&larr; Back to AIP</Link>
        <Card>
          <CardContent className="p-8 text-center">
            <p className="text-lg font-medium">Could not load sentiment for {ticker}</p>
            <p className="text-sm text-muted-foreground mt-1">The ticker may not be available or AIP is not configured.</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const { explanation, aip_data } = data;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/aip" className="text-sm text-muted-foreground hover:text-foreground">&larr; Back</Link>
          <h1 className="text-2xl font-bold">{ticker} Sentiment Analysis</h1>
        </div>
        <Link href={`/explain/${explanation.id}`}>
          <Button variant="outline" size="sm">Full Explanation &rarr;</Button>
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle>Sentiment Score</CardTitle></CardHeader>
          <CardContent>
            <div className="text-4xl font-bold tabular-nums">
              {aip_data.sentiment_7d >= 0 ? '+' : ''}{aip_data.sentiment_7d.toFixed(3)}
            </div>
            <div className="mt-1 text-sm text-muted-foreground">{aip_data.sentiment_label || 'N/A'}</div>
            <div className="mt-3">
              <ConfidenceGauge confidence={explanation.confidence} size="sm" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>AIP Raw Data</CardTitle></CardHeader>
          <CardContent className="space-y-2 text-sm">
            <div className="flex justify-between"><span className="text-muted-foreground">7-day sentiment</span><span className="font-mono">{aip_data.sentiment_7d.toFixed(3)}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">30-day sentiment</span><span className="font-mono">{aip_data.sentiment_30d.toFixed(3)}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">Trend</span><span className={`font-mono ${aip_data.trend >= 0 ? 'text-emerald-500' : 'text-rose-500'}`}>{aip_data.trend >= 0 ? '+' : ''}{aip_data.trend.toFixed(3)}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">Articles (7d)</span><span>{aip_data.article_count_7d}</span></div>
            <div className="flex justify-between"><span className="text-muted-foreground">Positive ratio</span><span>{(aip_data.positive_ratio * 100).toFixed(0)}%</span></div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader><CardTitle>Explanation Breakdown</CardTitle></CardHeader>
        <CardContent>
          <BreakdownChart breakdown={explanation.breakdown} />
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <DriverRanking drivers={explanation.top_drivers} />
        <Card>
          <CardHeader><CardTitle>Data Quality</CardTitle></CardHeader>
          <CardContent>
            {explanation.missing_impact > 0 ? (
              <div className="rounded-lg border border-amber-300 bg-amber-50 p-3 text-sm dark:border-amber-700 dark:bg-amber-950">
                {(explanation.missing_impact * 100).toFixed(1)}% of input data is missing
              </div>
            ) : (
              <div className="rounded-lg border border-emerald-300 bg-emerald-50 p-3 text-sm dark:border-emerald-700 dark:bg-emerald-950">
                All data sources available
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
