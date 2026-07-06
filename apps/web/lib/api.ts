import 'server-only'
import type { definitions } from '@/lib/openapi'

const BASE = process.env.API_URL ?? 'http://localhost:8080'

export async function getHealth(): Promise<
  definitions['health.HealthResponse']
> {
  const res = await fetch(`${BASE}/health`, { cache: 'no-store' })
  if (!res.ok) throw new Error(`health: ${res.status}`)
  return res.json()
}

export async function getMetrics(): Promise<string> {
  const res = await fetch(`${BASE}/metrics`, { cache: 'no-store' })
  if (!res.ok) throw new Error(`metrics: ${res.status}`)
  return res.text()
}

export async function getTasks(
  token: string,
): Promise<definitions['tasks.TaskResponse']> {
  const res = await fetch(`${BASE}/tasks`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  })
  if (!res.ok) throw new Error(`tasks: ${res.status}`)
  return res.json()
}

export async function login(
  body: definitions['auth.LoginRequest'],
): Promise<definitions['auth.AuthResponse']> {
  const res = await fetch(`${BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`login: ${res.status}`)
  return res.json()
}

export async function signup(
  body: definitions['auth.SignupRequest'],
): Promise<definitions['auth.AuthResponse']> {
  const res = await fetch(`${BASE}/auth/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`signup: ${res.status}`)
  return res.json()
}
