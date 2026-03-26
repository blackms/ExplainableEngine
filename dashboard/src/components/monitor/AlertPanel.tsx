'use client';

import { useState, useMemo, useCallback } from 'react';
import Link from 'next/link';
import { useExplanationList } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
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
  }>
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

function alertLabel(type: Alert['type']): string {
  return type === 'low_confidence' ? 'Low Confidence' : 'High Missing Data';
}

function alertVariant(type: Alert['type']) {
  return type === 'low_confidence' ? 'secondary' as const : 'destructive' as const;
}

export function AlertPanel() {
  const { data } = useExplanationList(
    { limit: 10 },
    { refetchInterval: 10000 }
  );

  const [dismissedIds, setDismissedIds] = useState<Set<string>>(new Set());

  const allAlerts = useMemo(
    () => deriveAlerts(data?.items ?? []),
    [data?.items]
  );

  const visibleAlerts = useMemo(
    () => allAlerts.filter((a) => !dismissedIds.has(a.id)),
    [allAlerts, dismissedIds]
  );

  const dismiss = useCallback((id: string) => {
    setDismissedIds((prev) => new Set(prev).add(id));
  }, []);

  const clearAll = useCallback(() => {
    setDismissedIds(new Set(allAlerts.map((a) => a.id)));
  }, [allAlerts]);

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div className="flex items-center gap-2">
          <CardTitle>Alerts</CardTitle>
          <AlertBadge count={visibleAlerts.length} />
        </div>
        {visibleAlerts.length > 0 && (
          <Button variant="ghost" size="sm" onClick={clearAll}>
            Clear all
          </Button>
        )}
      </CardHeader>
      <CardContent>
        {visibleAlerts.length === 0 ? (
          <p className="text-sm text-muted-foreground text-center py-6">
            No anomalies detected.
          </p>
        ) : (
          <div className="space-y-2">
            {visibleAlerts.map((alert) => (
              <div
                key={alert.id}
                className="flex items-center justify-between gap-2 rounded-lg border p-3 text-sm"
              >
                <div className="flex items-center gap-2 min-w-0 flex-1">
                  <Badge variant={alertVariant(alert.type)}>
                    {alertLabel(alert.type)}
                  </Badge>
                  <span className="truncate font-medium">{alert.target}</span>
                  <span className="text-muted-foreground shrink-0">
                    {(alert.value * 100).toFixed(0)}%
                  </span>
                </div>
                <div className="flex items-center gap-1 shrink-0">
                  <Link
                    href={`/explain/${alert.explanationId}`}
                    className="text-xs text-primary hover:underline whitespace-nowrap"
                  >
                    View &rarr;
                  </Link>
                  <Button
                    variant="ghost"
                    size="xs"
                    onClick={() => dismiss(alert.id)}
                    aria-label="Dismiss alert"
                  >
                    &times;
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
