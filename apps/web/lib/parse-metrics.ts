type MetricEntry = { labels: Record<string, string>; value: number }
export type MetricsMap = Map<string, MetricEntry[]>

export function parsePrometheus(text: string): MetricsMap {
  const result: MetricsMap = new Map()

  for (const raw of text.split('\n')) {
    const line = raw.trim()
    if (!line || line.startsWith('#')) continue

    const braceOpen = line.indexOf('{')
    let name: string
    const labels: Record<string, string> = {}
    let rest: string

    if (braceOpen !== -1) {
      name = line.slice(0, braceOpen)
      const braceClose = line.indexOf('}')
      for (const pair of line.slice(braceOpen + 1, braceClose).split(',')) {
        const eq = pair.indexOf('=')
        if (eq === -1) continue
        labels[pair.slice(0, eq).trim()] = pair.slice(eq + 1).replace(/^"|"$/g, '')
      }
      rest = line.slice(braceClose + 2)
    } else {
      const sp = line.indexOf(' ')
      name = line.slice(0, sp)
      rest = line.slice(sp + 1)
    }

    const value = parseFloat(rest.split(' ')[0])
    if (!isNaN(value)) {
      const entries = result.get(name) ?? []
      entries.push({ labels, value })
      result.set(name, entries)
    }
  }

  return result
}

export function metric(
  map: MetricsMap,
  name: string,
  labels?: Record<string, string>,
): number | null {
  const entries = map.get(name)
  if (!entries?.length) return null
  if (!labels) return entries[0].value
  return (
    entries.find((e) =>
      Object.entries(labels).every(([k, v]) => e.labels[k] === v),
    )?.value ?? null
  )
}
