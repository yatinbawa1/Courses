import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
       
  // Grab the cookie from the incoming request
  const refreshToken = request.cookies.get('refreshToken')?.value;
  const { pathname } = request.nextUrl;

  // 1. If trying to access dashboard without a token -> Redirect to login
  if (pathname.startsWith('/dashboard') && !refreshToken) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // 2. If trying to access login page *with* a token -> Redirect to dashboard
  if (pathname === '/login' && refreshToken) {
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  return NextResponse.next();
}

// Only run middleware on specific paths to keep performance fast
export const config = {
  matcher: ['/dashboard/:path*', '/login'],
};