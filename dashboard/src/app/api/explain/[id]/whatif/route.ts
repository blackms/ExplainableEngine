import type { NextRequest } from 'next/server';

const BACKEND_URL =
  process.env.BACKEND_URL ||
  'https://explainable-engine-516741092583.europe-west1.run.app';

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const body = await request.json();

  const res = await fetch(`${BACKEND_URL}/api/v1/explain/${id}/what-if`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  const data = await res.json();
  return Response.json(data, { status: res.status });
}
