import type { NextRequest } from 'next/server';

const BACKEND_URL =
  process.env.BACKEND_URL ||
  'https://explainable-engine-516741092583.europe-west1.run.app';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const format = request.nextUrl.searchParams.get('format') || 'json';

  const res = await fetch(
    `${BACKEND_URL}/explain/${id}/graph?format=${format}`,
    { headers: { 'Content-Type': 'application/json' } }
  );

  const data = await res.json();
  return Response.json(data, { status: res.status });
}
