'use client';

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
import type { SensitivityRanking as SensitivityRankingType } from '@/lib/api/types';

interface SensitivityRankingProps {
  ranking: SensitivityRankingType[];
}

function CustomTooltipContent({
  active,
  payload,
}: {
  active?: boolean;
  payload?: Array<{ payload: SensitivityRankingType }>;
}) {
  if (!active || !payload || payload.length === 0) return null;
  const item = payload[0].payload;
  return (
    <div className="rounded-md bg-card px-3 py-2 text-sm ring-1 ring-foreground/10 shadow-md">
      <p className="font-medium">{item.name}</p>
      <p className="text-muted-foreground">Impact: {item.impact.toFixed(4)}</p>
      <p className="text-muted-foreground">Rank: #{item.rank}</p>
    </div>
  );
}

const BAR_COLORS = [
  'hsl(220, 70%, 50%)',
  'hsl(200, 70%, 50%)',
  'hsl(180, 70%, 50%)',
  'hsl(160, 70%, 50%)',
  'hsl(140, 70%, 50%)',
  'hsl(260, 70%, 50%)',
  'hsl(280, 70%, 50%)',
  'hsl(300, 70%, 50%)',
];

export function SensitivityRanking({ ranking }: SensitivityRankingProps) {
  const sorted = [...ranking].sort((a, b) => b.impact - a.impact);
  const chartHeight = Math.max(200, sorted.length * 40 + 40);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Which inputs matter most?</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={chartHeight}>
          <BarChart data={sorted} layout="vertical" margin={{ left: 20, right: 40 }}>
            <XAxis
              type="number"
              fontSize={12}
              tickFormatter={(v: number) => v.toFixed(2)}
            />
            <YAxis
              type="category"
              dataKey="name"
              width={120}
              tick={{ fontSize: 12 }}
            />
            <RechartsTooltip content={<CustomTooltipContent />} />
            <Bar
              dataKey="impact"
              radius={[0, 4, 4, 0]}
              label={{
                position: 'right',
                formatter: (v: unknown) => Number(v).toFixed(4),
                fontSize: 11,
              }}
            >
              {sorted.map((_, index) => (
                <Cell
                  key={index}
                  fill={BAR_COLORS[index % BAR_COLORS.length]}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
