import { auth } from '@/lib/auth';
import { NextResponse } from 'next/server';

export default auth((req) => {
  // If auth is not configured (no Google credentials), allow all access
  if (!process.env.GOOGLE_CLIENT_ID) {
    return NextResponse.next();
  }

  if (!req.auth && !req.nextUrl.pathname.startsWith('/login')) {
    return NextResponse.redirect(new URL('/login', req.url));
  }

  return NextResponse.next();
});

export const config = {
  matcher: ['/((?!api/auth|_next/static|_next/image|favicon.ico).*)'],
};
