"use client";

import { useQuery } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { fetchApiStatus } from "@/lib/api";

function statusLabel(query: {
  isPending: boolean;
  isError: boolean;
  data?: { status: string };
}) {
  if (query.isPending) {
    return "checking…";
  }
  if (query.isError) {
    return "down";
  }
  return query.data?.status ?? "unknown";
}

function AppStatus() {
  const health = useQuery({
    queryKey: ["status", "health"],
    queryFn: () => fetchApiStatus("/api/health"),
  });
  const ready = useQuery({
    queryKey: ["status", "ready"],
    queryFn: () => fetchApiStatus("/api/ready"),
  });

  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col gap-6 p-8">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">
          App connection status
        </h1>
        <p className="text-sm text-muted-foreground">
          Health is the Go process. Ready is Postgres.
        </p>
      </div>
      <dl className="space-y-3 text-sm">
        <div className="flex items-center justify-between gap-4">
          <dt className="text-muted-foreground">API /api/health</dt>
          <dd className="font-medium">{statusLabel(health)}</dd>
        </div>
        <div className="flex items-center justify-between gap-4">
          <dt className="text-muted-foreground">Database /api/ready</dt>
          <dd className="font-medium">{statusLabel(ready)}</dd>
        </div>
      </dl>
      <Button
        variant="outline"
        onClick={() => {
          void health.refetch();
          void ready.refetch();
        }}
      >
        Recheck
      </Button>
    </main>
  );
}

export default AppStatus;
