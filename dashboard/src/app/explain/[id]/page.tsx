import { SummaryCard } from '@/components/explanation/SummaryCard';
import { BreakdownChart } from '@/components/explanation/BreakdownChart';
import { DriverRanking } from '@/components/explanation/DriverRanking';

export default async function ExplanationPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  const res = await fetch(
    `${process.env.BACKEND_URL || 'https://explainable-engine-516741092583.europe-west1.run.app'}/api/v1/explain/${id}`,
    { cache: 'no-store' }
  );

  if (!res.ok) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <div className="text-center space-y-2">
          <h2 className="text-lg font-semibold">Explanation not found</h2>
          <p className="text-sm text-muted-foreground">
            Could not load explanation with ID: {id}
          </p>
        </div>
      </div>
    );
  }

  const explanation = await res.json();

  return (
    <div className="space-y-6 p-6">
      <SummaryCard
        target={explanation.target}
        finalValue={explanation.final_value}
        confidence={explanation.confidence}
        missingImpact={explanation.missing_impact}
        topDrivers={explanation.top_drivers}
        metadata={explanation.metadata}
      />
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <BreakdownChart breakdown={explanation.breakdown} />
        <DriverRanking drivers={explanation.top_drivers} />
      </div>
    </div>
  );
}
