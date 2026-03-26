import type { NextRequest } from 'next/server';

const BACKEND_URL =
  process.env.BACKEND_URL ||
  'https://explainable-engine-516741092583.europe-west1.run.app';

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const params = searchParams.toString();
  const res = await fetch(
    `${BACKEND_URL}/api/v1/explain${params ? `?${params}` : ''}`,
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  const data = await res.json();
  return Response.json(data, { status: res.status });
}

export async function POST(request: Request) {
  const body = await request.json();

  const res = await fetch(`${BACKEND_URL}/api/v1/explain`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  const data = await res.json();
  return Response.json(data, { status: res.status });
}
