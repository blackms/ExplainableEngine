import { Skeleton } from '@/components/ui/skeleton';

export default function ExplanationLoading() {
  return (
    <div className="space-y-6 p-6">
      {/* Summary Card skeleton */}
      <div className="rounded-xl ring-1 ring-foreground/10 bg-card p-4 space-y-4">
        <div className="flex items-start justify-between">
          <div className="space-y-2">
            <Skeleton className="h-5 w-48" />
            <Skeleton className="h-9 w-32" />
          </div>
          <Skeleton className="h-[72px] w-[72px] rounded-full" />
        </div>
        <div className="space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-4 w-1/2" />
        </div>
        <div className="border-t pt-3">
          <Skeleton className="h-3 w-64" />
        </div>
      </div>

      {/* Chart + Ranking grid skeleton */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Breakdown Chart skeleton */}
        <div className="rounded-xl ring-1 ring-foreground/10 bg-card p-4 space-y-4">
          <Skeleton className="h-5 w-28" />
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-4 flex-1" />
              </div>
            ))}
          </div>
        </div>

        {/* Driver Ranking skeleton */}
        <div className="rounded-xl ring-1 ring-foreground/10 bg-card p-4 space-y-4">
          <Skeleton className="h-5 w-28" />
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="h-5 w-8 rounded-full" />
                <div className="flex-1 space-y-1">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-2 w-full rounded-full" />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
