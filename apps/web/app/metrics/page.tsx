import { getMetrics } from '@/lib/api'
import { metric, parsePrometheus } from '@/lib/parse-metrics'
import { RefreshButton } from '../refresh-button'

function fmt(n: number | null, unit = ''): string {
  if (n === null) return '—'
  return `${n.toFixed(4)}${unit}`
}

function fmtBytes(n: number | null): string {
  if (n === null) return '—'
  if (n >= 1_048_576) return `${(n / 1_048_576).toFixed(1)} MB`
  return `${(n / 1024).toFixed(1)} KB`
}

function fmtInt(n: number | null): string {
  if (n === null) return '—'
  return Math.round(n).toString()
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-4">
      <dt className="text-neutral-500 w-36">{label}</dt>
      <dd className="tabular-nums">{value}</dd>
    </div>
  )
}

export default async function MetricsPage() {
  let raw: string | null = null
  let error = false
  try {
    raw = await getMetrics()
  } catch {
    error = true
  }

  const m = raw ? parsePrometheus(raw) : null

  return (
    <div className="max-w-lg space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-mono font-semibold">Metrics</h1>
        <RefreshButton />
      </div>

      {error && (
        <p className="text-sm font-mono text-red-400">
          Could not fetch /metrics
        </p>
      )}

      {m && (
        <div className="space-y-8">
          <section className="space-y-2">
            <h2 className="text-xs font-mono uppercase tracking-widest text-neutral-500">
              HTTP — not instrumented
            </h2>
            <dl className="font-mono text-sm space-y-2">
              <Row label="request rate" value="—" />
              <Row label="error rate" value="—" />
              <Row label="p50 latency" value="—" />
              <Row label="p95 latency" value="—" />
              <Row label="p99 latency" value="—" />
            </dl>
          </section>

          <section className="space-y-2">
            <h2 className="text-xs font-mono uppercase tracking-widest text-neutral-500">
              GC duration (go_gc_duration_seconds)
            </h2>
            <dl className="font-mono text-sm space-y-2">
              <Row
                label="p50"
                value={fmt(
                  metric(m, 'go_gc_duration_seconds', { quantile: '0.5' }),
                  's',
                )}
              />
              <Row
                label="p95"
                value={fmt(
                  metric(m, 'go_gc_duration_seconds', { quantile: '0.95' }),
                  's',
                )}
              />
              <Row
                label="p99"
                value={fmt(
                  metric(m, 'go_gc_duration_seconds', { quantile: '0.99' }),
                  's',
                )}
              />
            </dl>
          </section>

          <section className="space-y-2">
            <h2 className="text-xs font-mono uppercase tracking-widest text-neutral-500">
              Go runtime
            </h2>
            <dl className="font-mono text-sm space-y-2">
              <Row
                label="goroutines"
                value={fmtInt(metric(m, 'go_goroutines'))}
              />
              <Row label="threads" value={fmtInt(metric(m, 'go_threads'))} />
              <Row
                label="heap alloc"
                value={fmtBytes(metric(m, 'go_memstats_alloc_bytes'))}
              />
              <Row
                label="heap sys"
                value={fmtBytes(metric(m, 'go_memstats_sys_bytes'))}
              />
              <Row
                label="GC cycles"
                value={fmtInt(metric(m, 'go_gc_duration_seconds_count'))}
              />
            </dl>
          </section>
        </div>
      )}
    </div>
  )
}
