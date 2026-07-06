// Plain (non-'use server') module: cookie names + shared types usable from
// both server and client code. A 'use server' file may only export async functions.
export const ACCESS_COOKIE = 'access_token'
export const REFRESH_COOKIE = 'refresh_token'

export type AuthState = { ok: boolean; error: string | null }
