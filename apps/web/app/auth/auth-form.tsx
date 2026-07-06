'use client'
import { useActionState, useState } from 'react'
import { authAction, logoutAction } from './actions'
import type { AuthState } from './shared'

type Mode = 'login' | 'signup'

const initialState: AuthState = { ok: false, error: null }

export function AuthForm({ loggedIn }: { loggedIn: boolean }) {
  const [mode, setMode] = useState<Mode>('login')
  const [state, formAction, pending] = useActionState(authAction, initialState)

  const isLoggedIn = loggedIn || state.ok

  return (
    <div className="max-w-sm space-y-6">
      <h1 className="text-lg font-mono font-semibold">Auth</h1>

      <div className="flex gap-4 font-mono text-sm border-b border-neutral-800 pb-3">
        {(['login', 'signup'] as Mode[]).map((m) => (
          <button
            key={m}
            type="button"
            onClick={() => setMode(m)}
            className={`transition-colors ${mode === m ? 'text-white' : 'text-neutral-500 hover:text-neutral-300'}`}
          >
            {m}
          </button>
        ))}
      </div>

      <form action={formAction} className="space-y-4">
        <input type="hidden" name="mode" value={mode} />

        {mode === 'signup' && (
          <div className="space-y-1">
            <label
              htmlFor="name"
              className="text-xs font-mono text-neutral-400"
            >
              name
            </label>
            <input
              id="name"
              type="text"
              name="name"
              required
              minLength={2}
              className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
            />
          </div>
        )}

        <div className="space-y-1">
          <label
            htmlFor="username"
            className="text-xs font-mono text-neutral-400"
          >
            username
          </label>
          <input
            id="username"
            type="text"
            name="username"
            required
            minLength={2}
            className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
          />
        </div>

        <div className="space-y-1">
          <label
            htmlFor="password"
            className="text-xs font-mono text-neutral-400"
          >
            password
          </label>
          <input
            id="password"
            type="password"
            name="password"
            required
            minLength={8}
            className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
          />
        </div>

        <button
          type="submit"
          disabled={pending}
          className="px-4 py-2 text-sm font-mono bg-white text-black rounded hover:bg-neutral-200 disabled:opacity-50 transition-colors"
        >
          {pending
            ? `${mode === 'login' ? 'logging in' : 'signing up'}…`
            : mode}
        </button>
      </form>

      {state.error && (
        <p className="text-red-400 text-sm font-mono">{state.error}</p>
      )}

      <div className="space-y-2 font-mono text-sm">
        <div className="flex items-center gap-3">
          <span
            className={`w-2.5 h-2.5 rounded-full shrink-0 ${isLoggedIn ? 'bg-green-500' : 'bg-neutral-600'}`}
          />
          <span className="text-neutral-400">
            {isLoggedIn
              ? 'logged in — token in httpOnly cookie'
              : 'not logged in'}
          </span>
        </div>
      </div>

      {isLoggedIn && (
        <form action={logoutAction}>
          <button
            type="submit"
            className="text-xs font-mono text-neutral-500 hover:text-white transition-colors"
          >
            log out
          </button>
        </form>
      )}
    </div>
  )
}
