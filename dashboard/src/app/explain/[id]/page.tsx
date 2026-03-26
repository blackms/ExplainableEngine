import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ValueCard, ConfidenceCard } from '@/components/explanation/SummaryCard';
import { BreakdownChart } from '@/components/explanation/BreakdownChart';
import { DriverRanking } from '@/components/explanation/DriverRanking';
import { ConfidencePanel } from '@/components/explanation/ConfidencePanel';
import { NarrativeViewer } from '@/components/explanation/NarrativeViewer';
import { SensitivityQuickView } from '@/components/explanation/SensitivityQuickView';
import { ExportPanel } from '@/components/explanation/ExportPanel';
import { AIAssistantSection } from '@/components/llm/AIAssistantSection';

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
      {/* Header: back + title + actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link
            href="/"
            className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            &larr; Back
          </Link>
          <h1 className="text-2xl font-bold">{explanation.target}</h1>
        </div>
        <div className="flex items-center gap-2">
          {explanation.graph && (
            <Link
              href={`/explain/${id}/graph`}
              className="inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm font-medium hover:bg-accent transition-colors"
            >
              View Graph
            </Link>
          )}
          <Link
            href={`/explain/${id}/whatif`}
            className="inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm font-medium hover:bg-accent transition-colors"
          >
            What-if
          </Link>
          <ExportPanel explanation={explanation} />
        </div>
      </div>

      {/* Metric cards: value + confidence */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <ValueCard
          value={explanation.final_value}
          target={explanation.target}
          hash={explanation.metadata.deterministic_hash}
        />
        <ConfidenceCard confidence={explanation.confidence} />
      </div>

      {/* Component Breakdown (full width) */}
      <Card>
        <CardHeader>
          <CardTitle>Component Breakdown</CardTitle>
        </CardHeader>
        <CardContent>
          <BreakdownChart breakdown={explanation.breakdown} />
        </CardContent>
      </Card>

      {/* Drivers + Data Quality */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <DriverRanking drivers={explanation.top_drivers} />
        <ConfidencePanel
          overall={explanation.confidence}
          perNode={explanation.confidence_detail?.per_node ?? {}}
          missingImpact={explanation.missing_impact}
        />
      </div>

      {/* Narrative + Sensitivity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <NarrativeViewer explanationId={id} />
        {explanation.top_drivers && explanation.top_drivers.length > 0 && (
          <SensitivityQuickView
            drivers={explanation.top_drivers}
            explanationId={id}
          />
        )}
      </div>

      {/* AI Assistant */}
      <AIAssistantSection explanationId={id} />
    </div>
  );
}
