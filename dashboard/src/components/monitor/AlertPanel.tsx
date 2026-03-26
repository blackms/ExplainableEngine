'use client';

import { useState, useMemo, useCallback } from 'react';
import Link from 'next/link';
import { AlertTriangle, ShieldCheck, X, XCircle } from 'lucide-react';
import { useExplanationList } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AlertBadge } from './AlertBadge';

export interface Alert {
  id: string;
  type: 'low_confidence' | 'high_missing';
  explanationId: string;
  target: string;
  value: number;
  timestamp: string;
}

function deriveAlerts(
  items: Array<{
    id: string;
    target: string;
    confidence: number;
    missing_impact: number;
    metadata: { created_at: string };
  }>,
): Alert[] {
  const alerts: Alert[] = [];

  for (const item of items) {
    if (item.confidence < 0.5) {
      alerts.push({
        id: `low-conf-${item.id}`,
        type: 'low_confidence',
        explanationId: item.id,
        target: item.target,
        value: item.confidence,
        timestamp: item.metadata.created_at,
      });
    }
    if (item.missing_impact > 0.2) {
      alerts.push({
        id: `high-miss-${item.id}`,
        type: 'high_missing',
        explanationId: item.id,
        target: item.target,
        value: item.missing_impact,
        timestamp: item.metadata.created_at,
      });
    }
  }

  return alerts;
}

function AlertIcon({ type }: { type: Alert['type'] }) {
  if (type === 'low_confidence') {
    return <AlertTriangle className="size-4 text-amber-500 shrink-0" />;
  }
  return <XCircle className="size-4 text-rose-500 shrink-0" />;
}

function alertDescription(alert: Alert): string {
  if (alert.type === 'low_confidence') {
    return `Low confidence (${(alert.value * 100).toFixed(0)}%) on ${alert.target}`;
  }
  return `High missing data (${(alert.value * 100).toFixed(0)}%) on ${alert.target}`;
}

export function AlertPanel() {
  const { data } = useExplanationList(
    { limit: 10 },
    { refetchInterval: 10000 },
  );

  const [dismissedIds, setDismissedIds] = useState<Set<string>>(new Set());

  const allAlerts = useMemo(
    () => deriveAlerts(data?.items ?? []),
    [data?.items],
  );

  const visibleAlerts = useMemo(
    () => allAlerts.filter((a) => !dismissedIds.has(a.id)),
    [allAlerts, dismissedIds],
  );

  const dismiss = useCallback((id: string) => {
    setDismissedIds((prev) => new Set(prev).add(id));
  }, []);

  const clearAll = useCallback(() => {
    setDismissedIds(new Set(allAlerts.map((a) => a.id)));
  }, [allAlerts]);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <CardTitle>Alerts</CardTitle>
          <AlertBadge count={visibleAlerts.length} />
        </div>
        {visibleAlerts.length > 0 && (
          <CardAction>
            <Button variant="ghost" size="sm" onClick={clearAll}>
              Clear all
            </Button>
          </CardAction>
        )}
      </CardHeader>
      <CardContent>
        {visibleAlerts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 space-y-4">
            <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
              <ShieldCheck className="h-6 w-6 text-muted-foreground" />
            </div>
            <div className="text-center space-y-1">
              <h3 className="text-base font-medium">All systems normal</h3>
              <p className="text-sm text-muted-foreground max-w-sm">
                No active alerts. Everything looks good.
              </p>
            </div>
          </div>
        ) : (
          <div className="space-y-2">
            {visibleAlerts.map((alert) => (
              <div
                key={alert.id}
                className="flex items-start gap-3 rounded-lg border p-3 text-sm"
              >
                <AlertIcon type={alert.type} />
                <div className="min-w-0 flex-1">
                  <p className="text-sm leading-snug">
                    {alertDescription(alert)}
                  </p>
                  <Link
                    href={`/explain/${alert.explanationId}`}
                    className="text-xs text-primary hover:underline mt-1 inline-block"
                  >
                    View details
                  </Link>
                </div>
                <Button
                  variant="ghost"
                  size="icon-xs"
                  onClick={() => dismiss(alert.id)}
                  aria-label="Dismiss alert"
                >
                  <X className="size-3.5" />
                </Button>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
