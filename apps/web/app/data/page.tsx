import { cookies } from 'next/headers'
import { getTasks } from '@/lib/api'
import { ACCESS_COOKIE } from '../auth/shared'
import { RefreshButton } from '../refresh-button'

type Task = {
  id: string
  task: string
  description: string
  is_completed: boolean
  created_at: string
}

export default async function DataPage() {
  const store = await cookies()
  const token = store.get(ACCESS_COOKIE)?.value

  if (!token) {
    return (
      <div className="max-w-lg space-y-4">
        <h1 className="text-lg font-mono font-semibold">Data</h1>
        <p className="text-sm font-mono text-neutral-400">
          Log in on the{' '}
          <a href="/auth" className="underline hover:text-white">
            auth page
          </a>{' '}
          first.
        </p>
      </div>
    )
  }

  let rows: Task[] = []
  let error = false
  try {
    const data = await getTasks(token)
    rows = (Array.isArray(data?.data) ? (data.data as Task[]) : []).slice(0, 20)
  } catch {
    error = true
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-mono font-semibold">Data — tasks</h1>
        <RefreshButton />
      </div>

      {error && (
        <p className="text-sm font-mono text-red-400">Failed to fetch tasks</p>
      )}

      {!error && rows.length === 0 && (
        <p className="text-sm font-mono text-neutral-500">No tasks found.</p>
      )}

      {rows.length > 0 && (
        <table className="w-full text-sm font-mono border-collapse">
          <thead>
            <tr className="text-left text-neutral-500 border-b border-neutral-800">
              <th className="pb-2 pr-6 font-normal">id</th>
              <th className="pb-2 pr-6 font-normal">title</th>
              <th className="pb-2 pr-6 font-normal">description</th>
              <th className="pb-2 pr-6 font-normal">done</th>
              <th className="pb-2 font-normal">created</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((t) => (
              <tr
                key={t.id}
                className="border-b border-neutral-900 hover:bg-neutral-900"
              >
                <td className="py-2 pr-6 text-neutral-500 font-mono text-xs max-w-24 truncate">
                  {t.id}
                </td>
                <td className="py-2 pr-6">{t.task}</td>
                <td className="py-2 pr-6 text-neutral-400">{t.description}</td>
                <td className="py-2 pr-6">{t.is_completed ? 'yes' : 'no'}</td>
                <td className="py-2 text-neutral-500">
                  {new Date(t.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <p className="text-xs text-neutral-600 font-mono">
        {rows.length} / 20 rows shown
      </p>
    </div>
  )
}
