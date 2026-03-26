import Link from 'next/link';
import { SummaryCard } from '@/components/explanation/SummaryCard';
import { BreakdownChart } from '@/components/explanation/BreakdownChart';
import { DriverRanking } from '@/components/explanation/DriverRanking';
import { ConfidencePanel } from '@/components/explanation/ConfidencePanel';
import { NarrativeViewer } from '@/components/explanation/NarrativeViewer';
import { SensitivityQuickView } from '@/components/explanation/SensitivityQuickView';
import { ExportPanel } from '@/components/explanation/ExportPanel';

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
      {/* Header with View Full Graph link */}
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-semibold sr-only">Explanation Detail</h1>
        <div className="ml-auto flex items-center gap-2">
          <ExportPanel explanation={explanation} />
          {explanation.graph && (
            <Link
              href={`/explain/${id}/graph`}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/80 transition-colors"
            >
              View Full Graph
            </Link>
          )}
        </div>
      </div>

      {/* Row 1: SummaryCard (full width) */}
      <SummaryCard
        target={explanation.target}
        finalValue={explanation.final_value}
        confidence={explanation.confidence}
        missingImpact={explanation.missing_impact}
        topDrivers={explanation.top_drivers}
        metadata={explanation.metadata}
      />

      {/* Row 2: BreakdownChart (left) + DriverRanking (right) */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <BreakdownChart breakdown={explanation.breakdown} />
        <DriverRanking drivers={explanation.top_drivers} />
      </div>

      {/* Row 3: ConfidencePanel (left) + NarrativeViewer (right) */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ConfidencePanel
          overall={explanation.confidence}
          perNode={explanation.confidence_detail?.per_node ?? {}}
          missingImpact={explanation.missing_impact}
        />
        <NarrativeViewer explanationId={id} />
      </div>

      {/* Row 4: Sensitivity Quick View */}
      {explanation.top_drivers && explanation.top_drivers.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <SensitivityQuickView
            drivers={explanation.top_drivers}
            explanationId={id}
          />
        </div>
      )}
    </div>
  );
}
