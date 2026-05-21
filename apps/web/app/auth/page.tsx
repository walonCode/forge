'use client'
import { useState } from 'react'
import { login, signup, checkAuth } from '@/lib/api'
import { useToken } from '@/lib/providers'

type Mode = 'login' | 'signup'

export default function AuthPage() {
  const { token, setToken } = useToken()
  const [mode, setMode] = useState<Mode>('login')
  const [name, setName] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [authOk, setAuthOk] = useState<boolean | null>(null)

  function switchMode(next: Mode) {
    setMode(next)
    setError(null)
    setName('')
    setUsername('')
    setPassword('')
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res =
        mode === 'login'
          ? await login({ username, password })
          : await signup({ name, username, password })

      const at = res.data?.access_token ?? null
      setToken(at)
      if (at) {
        const ok = await checkAuth(at)
        setAuthOk(ok)
      }
    } catch {
      setError(mode === 'login' ? 'Invalid username or password' : 'Signup failed — username may already exist')
      setToken(null)
      setAuthOk(false)
    } finally {
      setLoading(false)
    }
  }

  function handleClear() {
    setToken(null)
    setAuthOk(null)
    setError(null)
  }

  return (
    <div className="max-w-sm space-y-6">
      <h1 className="text-lg font-mono font-semibold">Auth</h1>

      <div className="flex gap-4 font-mono text-sm border-b border-neutral-800 pb-3">
        {(['login', 'signup'] as Mode[]).map((m) => (
          <button
            key={m}
            onClick={() => switchMode(m)}
            className={`transition-colors ${mode === m ? 'text-white' : 'text-neutral-500 hover:text-neutral-300'}`}
          >
            {m}
          </button>
        ))}
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {mode === 'signup' && (
          <div className="space-y-1">
            <label className="text-xs font-mono text-neutral-400">name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              minLength={2}
              className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
            />
          </div>
        )}

        <div className="space-y-1">
          <label className="text-xs font-mono text-neutral-400">username</label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            minLength={2}
            className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
          />
        </div>

        <div className="space-y-1">
          <label className="text-xs font-mono text-neutral-400">password</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            className="w-full bg-neutral-900 border border-neutral-700 rounded px-3 py-2 text-sm font-mono outline-none focus:border-neutral-400 transition-colors"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm font-mono bg-white text-black rounded hover:bg-neutral-200 disabled:opacity-50 transition-colors"
        >
          {loading ? `${mode === 'login' ? 'logging in' : 'signing up'}…` : mode}
        </button>
      </form>

      {error && <p className="text-red-400 text-sm font-mono">{error}</p>}

      <div className="space-y-2 font-mono text-sm">
        <div className="flex items-center gap-3">
          <span
            className={`w-2.5 h-2.5 rounded-full shrink-0 ${token ? 'bg-green-500' : 'bg-neutral-600'}`}
          />
          <span className="text-neutral-400">
            {token ? 'token stored in memory' : 'no token'}
          </span>
        </div>

        {authOk !== null && (
          <div className="flex items-center gap-3">
            <span
              className={`w-2.5 h-2.5 rounded-full shrink-0 ${authOk ? 'bg-green-500' : 'bg-red-500'}`}
            />
            <span className="text-neutral-400">
              GET /tasks →{' '}
              <span className={authOk ? 'text-green-400' : 'text-red-400'}>
                {authOk ? '200 ok' : 'failed'}
              </span>
            </span>
          </div>
        )}
      </div>

      {token && (
        <button
          onClick={handleClear}
          className="text-xs font-mono text-neutral-500 hover:text-white transition-colors"
        >
          clear token
        </button>
      )}
    </div>
  )
}
