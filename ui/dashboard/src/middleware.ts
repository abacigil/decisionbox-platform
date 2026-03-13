import { NextRequest, NextResponse } from 'next/server';

// Runtime API proxy — reads API_URL env var at request time (not build time).
// Proxies all /api/* requests to the backend API server-side.
// Browser never talks to the API directly — only to the dashboard.
//
// Configuration via env var:
//   API_URL=http://decisionbox-api:8080  (K8s cluster-internal)
//   API_URL=http://localhost:8080         (local dev, default)
export async function middleware(request: NextRequest) {
  const { pathname, search } = request.nextUrl;

  // Only proxy /api/* requests (not /health or other dashboard routes)
  if (!pathname.startsWith('/api/')) {
    return NextResponse.next();
  }

  const apiUrl = process.env.API_URL || 'http://localhost:8080';
  const targetUrl = `${apiUrl}${pathname}${search}`;

  // Forward the request to the backend API
  const headers = new Headers(request.headers);
  // Remove host header (will be set by fetch to the target)
  headers.delete('host');

  const response = await fetch(targetUrl, {
    method: request.method,
    headers,
    body: request.body,
    // @ts-expect-error duplex is needed for streaming request bodies
    duplex: 'half',
  });

  // Forward the response back to the client
  const responseHeaders = new Headers(response.headers);
  // Remove transfer-encoding to avoid issues with Next.js
  responseHeaders.delete('transfer-encoding');

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: responseHeaders,
  });
}

export const config = {
  // Only run middleware on /api/* paths
  matcher: '/api/:path*',
};
