'use server'
import { revalidatePath } from 'next/cache'
import { cookies } from 'next/headers'
import { login, signup } from '@/lib/api'
import type { definitions } from '@/lib/openapi'
import { ACCESS_COOKIE, type AuthState, REFRESH_COOKIE } from './shared'

const COOKIE_OPTS = {
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'lax' as const,
  path: '/',
}

async function storeTokens(
  data: definitions['auth.AuthResponseData'] | undefined,
): Promise<boolean> {
  if (!data?.access_token) return false
  const store = await cookies()
  store.set(ACCESS_COOKIE, data.access_token, {
    ...COOKIE_OPTS,
    maxAge: 60 * 60,
  })
  if (data.refresh_token) {
    store.set(REFRESH_COOKIE, data.refresh_token, {
      ...COOKIE_OPTS,
      maxAge: 60 * 60 * 24 * 7,
    })
  }
  return true
}

export async function authAction(
  _prev: AuthState,
  formData: FormData,
): Promise<AuthState> {
  const mode = String(formData.get('mode') ?? 'login')
  const username = String(formData.get('username') ?? '')
  const password = String(formData.get('password') ?? '')
  const name = String(formData.get('name') ?? '')

  try {
    const res =
      mode === 'signup'
        ? await signup({ name, username, password })
        : await login({ username, password })

    if (!(await storeTokens(res.data))) {
      return {
        ok: false,
        error:
          mode === 'signup' ? 'Signup failed' : 'Invalid username or password',
      }
    }
    revalidatePath('/', 'layout')
    return { ok: true, error: null }
  } catch {
    return {
      ok: false,
      error:
        mode === 'signup'
          ? 'Signup failed — username may already exist'
          : 'Invalid username or password',
    }
  }
}

export async function logoutAction(): Promise<void> {
  const store = await cookies()
  store.delete(ACCESS_COOKIE)
  store.delete(REFRESH_COOKIE)
  revalidatePath('/', 'layout')
}
