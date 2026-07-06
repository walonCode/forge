import { cookies } from 'next/headers'
import { AuthForm } from './auth-form'
import { ACCESS_COOKIE } from './shared'

export default async function AuthPage() {
  const store = await cookies()
  const loggedIn = Boolean(store.get(ACCESS_COOKIE)?.value)

  return <AuthForm loggedIn={loggedIn} />
}
