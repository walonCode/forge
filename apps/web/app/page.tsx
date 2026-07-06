import { getHealth } from '@/lib/api'
import { RefreshButton } from './refresh-button'

export default async function HealthPage() {
  let data: Awaited<ReturnType<typeof getHealth>> | null = null
  let error = false
  try {
    data = await getHealth()
  } catch {
    error = true
  }

  const ok = !error && data?.status === 'ok'
  const dotColor = ok ? 'bg-green-500' : 'bg-red-500'

  return (
    <div className="max-w-md space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-mono font-semibold">Health</h1>
        <RefreshButton />
      </div>

      <div className="flex items-center gap-3">
        <span className={`w-3 h-3 rounded-full ${dotColor}`} />
        <span className="font-mono text-sm">{ok ? 'ok' : 'degraded'}</span>
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

      {error && (
        <p className="text-red-400 text-sm font-mono">
          Could not reach API — is the server running on port 8080?
        </p>
      )}
    </div>
  )
}
