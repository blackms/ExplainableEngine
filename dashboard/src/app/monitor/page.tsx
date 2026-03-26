'use client';

import { StatsCards } from '@/components/monitor/StatsCards';
import { LiveFeed } from '@/components/monitor/LiveFeed';
import { AlertPanel } from '@/components/monitor/AlertPanel';

export default function MonitorPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Live Monitoring</h1>
        <p className="text-sm text-muted-foreground">
          Real-time overview of explanation activity and alerts
        </p>
      </div>
      <StatsCards />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <LiveFeed />
        </div>
        <div>
          <AlertPanel />
        </div>
      </div>
    </div>
  );
}
