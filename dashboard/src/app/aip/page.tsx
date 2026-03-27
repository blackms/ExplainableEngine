'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAIPMarketMood } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';

const POPULAR_TICKERS = ['AAPL', 'MSFT', 'NVDA', 'TSLA', 'AMZN', 'GOOGL', 'META', 'JPM'];

export default function AIPPage() {
  const [ticker, setTicker] = useState('');
  const router = useRouter();
  const { data: mood, isLoading: moodLoading, error: moodError } = useAIPMarketMood();

  const handleExplain = () => {
    if (ticker.trim()) router.push(`/aip/${ticker.trim().toUpperCase()}`);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">AIP Sentiment Intelligence</h1>
        <p className="text-sm text-muted-foreground">Explain sentiment scores powered by AIP</p>
      </div>

      <Card>
        <CardContent className="p-4">
          <form onSubmit={(e) => { e.preventDefault(); handleExplain(); }} className="flex gap-2">
            <Input
              placeholder="Enter ticker (e.g. AAPL)"
              value={ticker}
              onChange={(e) => setTicker(e.target.value)}
              className="max-w-xs uppercase"
            />
            <Button type="submit" disabled={!ticker.trim()}>Explain</Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Market Mood</CardTitle>
        </CardHeader>
        <CardContent>
          {moodLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-6 w-48" />
              <Skeleton className="h-4 w-32" />
            </div>
          ) : moodError ? (
            <p className="text-sm text-muted-foreground">Could not load market mood. AIP may not be configured.</p>
          ) : mood ? (
            <div className="space-y-2">
              <div className="flex items-baseline gap-3">
                <span className="text-3xl font-bold tabular-nums">
                  {mood.aip_data.overall_sentiment >= 0 ? '+' : ''}{mood.aip_data.overall_sentiment.toFixed(3)}
                </span>
                <span className={`text-sm ${mood.aip_data.overall_trend >= 0 ? 'text-emerald-500' : 'text-rose-500'}`}>
                  {mood.aip_data.overall_trend >= 0 ? '\u2191' : '\u2193'} {mood.aip_data.overall_trend.toFixed(3)}
                </span>
              </div>
              <p className="text-xs text-muted-foreground">{mood.aip_data.total_articles.toLocaleString()} articles analyzed</p>
              <Button variant="outline" size="sm" onClick={() => router.push(`/explain/${mood.explanation.id}`)}>
                View Full Breakdown
              </Button>
            </div>
          ) : null}
        </CardContent>
      </Card>

      <div>
        <h2 className="text-lg font-semibold mb-3">Quick Explain</h2>
        <div className="flex flex-wrap gap-2">
          {POPULAR_TICKERS.map(t => (
            <Button key={t} variant="outline" size="sm" onClick={() => router.push(`/aip/${t}`)}>
              {t}
            </Button>
          ))}
        </div>
      </div>
    </div>
  );
}
