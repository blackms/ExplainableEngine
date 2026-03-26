import { NextResponse } from 'next/server';
const BACKEND_URL = process.env.BACKEND_URL || 'https://explainable-engine-516741092583.europe-west1.run.app';

export async function GET() {
  const res = await fetch(`${BACKEND_URL}/api/v1/aip/explain/market-mood`, {
    headers: { 'Content-Type': 'application/json' },
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
