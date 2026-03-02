import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL =
  process.env.BACKEND_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  (process.env.NODE_ENV === "production"
    ? "https://satisfied-strength-production-fcdb.up.railway.app/api/v1"
    : "http://localhost:3000/api/v1");

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(request, await params);
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(request, await params);
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(request, await params);
}

export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(request, await params);
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxyRequest(request, await params);
}

async function proxyRequest(
  request: NextRequest,
  params: { path: string[] }
) {
  const path = params.path.join("/");
  const backendUrl = `${BACKEND_URL}/${path}`;

  const url = new URL(backendUrl);
  request.nextUrl.searchParams.forEach((value, key) => {
    url.searchParams.append(key, value);
  });

  const headers = new Headers();
  request.headers.forEach((value, key) => {
    if (!key.toLowerCase().startsWith("host") && !key.toLowerCase().startsWith("connection")) {
      headers.set(key, value);
    }
  });

  const body = request.method !== "GET" && request.method !== "HEAD"
    ? await request.text()
    : undefined;

  try {
    console.log(`[Proxy] ${request.method} ${url.toString()}`);

    const response = await fetch(url.toString(), {
      method: request.method,
      headers,
      body,
      credentials: "include",
    });

    console.log(`[Proxy] Response: ${response.status} ${response.statusText}`);

    const responseHeaders = new Headers(response.headers);

    const cookies = response.headers.getSetCookie?.() || [];
    console.log(`[Proxy] Set-Cookie headers: ${cookies.length}`);

    if (cookies.length > 0) {
      responseHeaders.delete("set-cookie");
      cookies.forEach((cookie) => {
        responseHeaders.append("set-cookie", cookie);
      });
    }

    const responseBody = await response.arrayBuffer();

    return new NextResponse(responseBody, {
      status: response.status,
      statusText: response.statusText,
      headers: responseHeaders,
    });
  } catch (error) {
    console.error("[Proxy] Error:", error);
    return NextResponse.json(
      { code: "PROXY_ERROR", message: "Failed to connect to backend" },
      { status: 502 }
    );
  }
}
