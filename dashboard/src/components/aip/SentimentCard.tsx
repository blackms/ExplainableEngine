'use client';
import Link from 'next/link';
import { Card, CardContent } from '@/components/ui/card';

interface SentimentCardProps {
  ticker: string;
  sentiment: number;
  label: string;
  trend: number;
  articleCount: number;
}

function sentimentColor(label: string): string {
  if (label.toLowerCase().includes('bullish')) return 'text-emerald-600';
  if (label.toLowerCase().includes('bearish')) return 'text-rose-600';
  return 'text-slate-500';
}

function trendArrow(trend: number): { icon: string; color: string } {
  if (trend > 0.01) return { icon: '\u2191', color: 'text-emerald-500' };
  if (trend < -0.01) return { icon: '\u2193', color: 'text-rose-500' };
  return { icon: '\u2192', color: 'text-slate-400' };
}

export function SentimentCard({ ticker, sentiment, label, trend, articleCount }: SentimentCardProps) {
  const t = trendArrow(trend);
  return (
    <Link href={`/aip/${ticker}`}>
      <Card className="cursor-pointer transition-shadow hover:shadow-md">
        <CardContent className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold">{ticker}</span>
            <span className={`text-xs font-medium ${sentimentColor(label)}`}>{label}</span>
          </div>
          <div className="mt-2 flex items-baseline gap-2">
            <span className="text-2xl font-bold tabular-nums">
              {sentiment >= 0 ? '+' : ''}{sentiment.toFixed(3)}
            </span>
            <span className={`text-sm ${t.color}`}>{t.icon} {trend >= 0 ? '+' : ''}{trend.toFixed(3)}</span>
          </div>
          <div className="mt-1 text-xs text-muted-foreground">{articleCount} articles</div>
        </CardContent>
      </Card>
    </Link>
  );
}
