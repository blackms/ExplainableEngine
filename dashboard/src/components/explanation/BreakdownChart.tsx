'use client';

import { useState, useCallback, useMemo } from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  Cell,
} from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface BreakdownItem {
  node_id: string;
  label: string;
  value: number;
  weight: number;
  absolute_contribution: number;
  percentage: number;
  confidence: number;
  children?: BreakdownItem[];
}

interface BreakdownChartProps {
  breakdown: BreakdownItem[];
  onDrillDown?: (nodeId: string) => void;
}

interface BreadcrumbEntry {
  label: string;
  items: BreakdownItem[];
}

function getBarColor(confidence: number): string {
  const hue = confidence * 120; // 0 = red, 60 = yellow, 120 = green
  return `hsl(${hue}, 70%, 50%)`;
}

function CustomTooltipContent({
  active,
  payload,
}: {
  active?: boolean;
  payload?: Array<{ payload: BreakdownItem }>;
}) {
  if (!active || !payload || payload.length === 0) return null;
  const item = payload[0].payload;
  return (
    <div className="rounded-md bg-card px-3 py-2 text-sm ring-1 ring-foreground/10 shadow-md">
      <p className="font-medium">{item.label}</p>
      <p className="text-muted-foreground">Contribution: {item.absolute_contribution.toFixed(4)}</p>
      <p className="text-muted-foreground">Percentage: {item.percentage.toFixed(1)}%</p>
      <p className="text-muted-foreground">Weight: {item.weight.toFixed(4)}</p>
      <p className="text-muted-foreground">Confidence: {(item.confidence * 100).toFixed(1)}%</p>
      {item.children && item.children.length > 0 && (
        <p className="text-xs text-primary mt-1">Click to drill down</p>
      )}
    </div>
  );
}

export function BreakdownChart({ breakdown, onDrillDown }: BreakdownChartProps) {
  const [drillStack, setDrillStack] = useState<BreadcrumbEntry[]>([]);

  const currentItems = useMemo(() => {
    const items = drillStack.length > 0
      ? drillStack[drillStack.length - 1].items
      : breakdown;
    return [...items].sort((a, b) => b.percentage - a.percentage);
  }, [drillStack, breakdown]);

  const handleBarClick = useCallback(
    (item: BreakdownItem) => {
      if (item.children && item.children.length > 0) {
        setDrillStack((prev) => [...prev, { label: item.label, items: item.children! }]);
        onDrillDown?.(item.node_id);
      }
    },
    [onDrillDown]
  );

  const handleBreadcrumbClick = useCallback((index: number) => {
    if (index < 0) {
      setDrillStack([]);
    } else {
      setDrillStack((prev) => prev.slice(0, index + 1));
    }
  }, []);

  const chartHeight = Math.max(200, currentItems.length * 40 + 40);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Breakdown</CardTitle>
        {drillStack.length > 0 && (
          <nav className="flex items-center gap-1 text-sm text-muted-foreground mt-1">
            <button
              type="button"
              onClick={() => handleBreadcrumbClick(-1)}
              className="hover:text-foreground transition-colors"
            >
              Root
            </button>
            {drillStack.map((entry, i) => (
              <span key={i} className="flex items-center gap-1">
                <span>/</span>
                <button
                  type="button"
                  onClick={() => handleBreadcrumbClick(i)}
                  className={`hover:text-foreground transition-colors ${
                    i === drillStack.length - 1 ? 'text-foreground font-medium' : ''
                  }`}
                >
                  {entry.label}
                </button>
              </span>
            ))}
          </nav>
        )}
      </CardHeader>
      <CardContent>
        {drillStack.length > 0 && (
          <button
            type="button"
            onClick={() => setDrillStack((prev) => prev.slice(0, -1))}
            className="mb-3 text-sm text-primary hover:underline"
          >
            &larr; Back
          </button>
        )}
        <ResponsiveContainer width="100%" height={chartHeight}>
          <BarChart data={currentItems} layout="vertical" margin={{ left: 20, right: 40 }}>
            <XAxis type="number" domain={[0, 100]} tickFormatter={(v) => `${v}%`} fontSize={12} />
            <YAxis
              type="category"
              dataKey="label"
              width={120}
              tick={{ fontSize: 12 }}
            />
            <RechartsTooltip content={<CustomTooltipContent />} />
            <Bar
              dataKey="percentage"
              radius={[0, 4, 4, 0]}
              cursor="pointer"
              onClick={(_data: unknown, index: number) => {
                if (index >= 0 && index < currentItems.length) {
                  handleBarClick(currentItems[index]);
                }
              }}
              label={{ position: 'right', formatter: (v: unknown) => `${Number(v).toFixed(1)}%`, fontSize: 11 }}
            >
              {currentItems.map((item) => (
                <Cell key={item.node_id} fill={getBarColor(item.confidence)} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
