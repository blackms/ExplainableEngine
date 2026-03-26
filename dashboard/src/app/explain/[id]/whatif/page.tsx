import Link from 'next/link';
import { WhatIfPageClient } from './WhatIfPageClient';

export default async function WhatIfPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  const res = await fetch(
    `${process.env.BACKEND_URL || 'https://explainable-engine-516741092583.europe-west1.run.app'}/api/v1/explain/${id}`,
    { cache: 'no-store' },
  );

  if (!res.ok) {
    return (
      <div className="flex items-center justify-center py-20 text-muted-foreground">
        Explanation not found
      </div>
    );
  }

  const explanation = await res.json();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          What-if Analysis: {explanation.target}
        </h1>
        <Link
          href={`/explain/${id}`}
          className="text-sm text-primary hover:underline"
        >
          &larr; Back to detail
        </Link>
      </div>
      <WhatIfPageClient explanation={explanation} />
    </div>
  );
}
