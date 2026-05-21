'use client'
import { useQuery } from '@tanstack/react-query'
import { getHealth } from '@/lib/api'

export default function HealthPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['health'],
    queryFn: getHealth,
    refetchInterval: 5000,
    retry: false,
  })

  const ok = !isError && !isLoading && data?.status === 'ok'
  const dotColor = isLoading
    ? 'bg-yellow-500'
    : ok
      ? 'bg-green-500'
      : 'bg-red-500'

  return (
    <div className="max-w-md space-y-6">
      <h1 className="text-lg font-mono font-semibold">Health</h1>

      <div className="flex items-center gap-3">
        <span className={`w-3 h-3 rounded-full ${dotColor}`} />
        <span className="font-mono text-sm">
          {isLoading ? 'checking…' : ok ? 'ok' : 'degraded'}
        </span>
      </div>

      {data && (
        <dl className="font-mono text-sm space-y-2">
          <div className="flex gap-4">
            <dt className="text-neutral-500 w-24">uptime</dt>
            <dd>{data.uptime ?? '—'}</dd>
          </div>
          <div className="flex gap-4">
            <dt className="text-neutral-500 w-24">database</dt>
            <dd>{data.db ?? '—'}</dd>
          </div>
          <div className="flex gap-4">
            <dt className="text-neutral-500 w-24">version</dt>
            <dd>{data.version ?? '—'}</dd>
          </div>
        </dl>
      )}

      {isError && (
        <p className="text-red-400 text-sm font-mono">
          Could not reach API — is the server running on port 8080?
        </p>
      )}

      <p className="text-xs text-neutral-600 font-mono">polls every 5s</p>
    </div>
  )
}
