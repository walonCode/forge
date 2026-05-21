import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import Link from 'next/link'
import { Providers } from '@/lib/providers'
import './globals.css'

const geistSans = Geist({ variable: '--font-geist-sans', subsets: ['latin'] })
const geistMono = Geist_Mono({ variable: '--font-geist-mono', subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Forge',
  description: 'Forge API dashboard',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col bg-background text-foreground">
        <Providers>
          <nav className="border-b border-neutral-800 px-6 py-3 flex gap-6 text-sm font-mono">
            <Link href="/" className="text-neutral-400 hover:text-white transition-colors">
              health
            </Link>
            <Link href="/metrics" className="text-neutral-400 hover:text-white transition-colors">
              metrics
            </Link>
            <Link href="/data" className="text-neutral-400 hover:text-white transition-colors">
              data
            </Link>
            <Link href="/auth" className="text-neutral-400 hover:text-white transition-colors">
              auth
            </Link>
          </nav>
          <main className="flex-1 p-6">{children}</main>
        </Providers>
      </body>
    </html>
  )
}
