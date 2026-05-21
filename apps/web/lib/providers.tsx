'use client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createContext, useContext, useState } from 'react'

type TokenCtx = { token: string | null; setToken: (t: string | null) => void }

export const TokenContext = createContext<TokenCtx>({ token: null, setToken: () => {} })
export const useToken = () => useContext(TokenContext)

export function Providers({ children }: { children: React.ReactNode }) {
  const [client] = useState(() => new QueryClient())
  const [token, setToken] = useState<string | null>(null)

  return (
    <QueryClientProvider client={client}>
      <TokenContext.Provider value={{ token, setToken }}>
        {children}
      </TokenContext.Provider>
    </QueryClientProvider>
  )
}
