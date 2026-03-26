import Link from 'next/link';
import { ExplanationGraph } from '@/components/graph/ExplanationGraph';

export default async function GraphPage({
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
          <h2 className="text-lg font-semibold">Graph not available</h2>
          <p className="text-sm text-muted-foreground">
            Could not load explanation graph for ID: {id}
          </p>
          <Link
            href={`/explain/${id}`}
            className="text-sm text-primary hover:underline"
          >
            Back to detail
          </Link>
        </div>
      </div>
    );
  }

  const explanation = await res.json();

  if (!explanation.graph) {
    return (
      <div className="flex items-center justify-center min-h-[40vh]">
        <div className="text-center space-y-2">
          <h2 className="text-lg font-semibold">No graph data</h2>
          <p className="text-sm text-muted-foreground">
            This explanation does not include graph data.
          </p>
          <Link
            href={`/explain/${id}`}
            className="text-sm text-primary hover:underline"
          >
            Back to detail
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="relative">
      <div className="absolute top-4 left-4 z-10">
        <Link
          href={`/explain/${id}`}
          className="inline-flex items-center gap-1.5 rounded-md bg-background/80 px-3 py-1.5 text-sm font-medium ring-1 ring-foreground/10 backdrop-blur-sm hover:bg-background transition-colors"
        >
          &larr; Back to detail
        </Link>
      </div>
      <ExplanationGraph
        graph={explanation.graph}
        className="h-[calc(100vh-4rem)] w-full"
      />
    </div>
  );
}
