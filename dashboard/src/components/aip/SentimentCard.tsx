'use client';

import Link from 'next/link';
import { Card, CardContent } from '@/components/ui/card';
import type { AIPSentimentData } from '@/lib/api/types';

function sentimentStyle(label: string) {
  const lower = label.toLowerCase();
  if (lower === 'bullish')
    return {
      text: 'text-emerald-600 dark:text-emerald-400',
      bg: 'bg-emerald-500/10',
      border: 'border-emerald-500/30',
      badge: 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300',
    };
  if (lower === 'bearish')
    return {
      text: 'text-rose-600 dark:text-rose-400',
      bg: 'bg-rose-500/10',
      border: 'border-rose-500/30',
      badge: 'bg-rose-500/15 text-rose-700 dark:text-rose-300',
    };
  return {
    text: 'text-muted-foreground',
    bg: 'bg-muted/50',
    border: 'border-border',
    badge: 'bg-muted text-muted-foreground',
  };
}

function trendIndicator(trend: number) {
  if (trend > 0.01) return { arrow: '\u2191', color: 'text-emerald-500' };
  if (trend < -0.01) return { arrow: '\u2193', color: 'text-rose-500' };
  return { arrow: '\u2192', color: 'text-muted-foreground' };
}

interface SentimentCardProps {
  data: AIPSentimentData;
}

export function SentimentCard({ data }: SentimentCardProps) {
  const style = sentimentStyle(data.sentiment_label);
  const trend = trendIndicator(data.trend);

  return (
    <Link href={`/aip/${data.ticker}`} className="group block">
      <Card
        className={`h-full border ${style.border} transition-all duration-150 group-hover:shadow-sm group-hover:ring-1 group-hover:ring-primary/20`}
      >
        <CardContent className="space-y-3 p-4">
          {/* Ticker + Label */}
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold tracking-tight">
              {data.ticker}
            </span>
            <span
              className={`rounded-full px-2 py-0.5 text-xs font-medium ${style.badge}`}
            >
              {data.sentiment_label}
            </span>
          </div>

          {/* Sentiment value + trend */}
          <div className="flex items-baseline gap-2">
            <span className={`text-2xl font-bold tabular-nums ${style.text}`}>
              {data.sentiment_7d >= 0 ? '+' : ''}
              {data.sentiment_7d.toFixed(3)}
            </span>
            <span className={`text-sm font-medium ${trend.color}`}>
              {trend.arrow}{' '}
              {data.trend >= 0 ? '+' : ''}
              {data.trend.toFixed(3)}
            </span>
          </div>

          {/* Article count */}
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>
              {data.article_count_7d} article{data.article_count_7d !== 1 ? 's' : ''} (7d)
            </span>
            <span>{(data.positive_ratio * 100).toFixed(0)}% positive</span>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
