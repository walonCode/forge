'use client'
import { useRouter } from 'next/navigation'
import { useTransition } from 'react'

export function RefreshButton() {
  const router = useRouter()
  const [pending, startTransition] = useTransition()

  return (
    <button
      type="button"
      onClick={() => startTransition(() => router.refresh())}
      className="text-xs font-mono text-neutral-400 hover:text-white transition-colors"
    >
      {pending ? 'refreshing…' : 'refresh'}
    </button>
  )
}
