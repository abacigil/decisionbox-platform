import type { NextConfig } from "next";

const apiUrl = process.env.API_URL || "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => [
    {
      // Proxy all /api/* requests to the backend API server-side.
      // Browser never talks to the API directly — only to the dashboard.
      // In K8s: API_URL=http://decisionbox-api:8080 (cluster-internal)
      // In dev: API_URL=http://localhost:8080 (default)
      source: "/api/:path*",
      destination: `${apiUrl}/api/:path*`,
    },
  ],
};

export default nextConfig;
