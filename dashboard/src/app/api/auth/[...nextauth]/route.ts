import { handlers } from '@/lib/auth';
import { NextResponse } from 'next/server';

const isAuthConfigured = !!(process.env.GOOGLE_CLIENT_ID && process.env.GOOGLE_CLIENT_SECRET);

export const GET = isAuthConfigured
  ? handlers.GET
  : () => NextResponse.json({ user: null, expires: '' });

export const POST = isAuthConfigured
  ? handlers.POST
  : () => NextResponse.json({ user: null, expires: '' });
