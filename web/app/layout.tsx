import type { Metadata } from 'next';
import './globals.css';
import Providers from '@/components/providers';

export const metadata: Metadata = {
  title: 'Aegis XDR / RMM',
  description: 'Enterprise monitoring and security platform',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased bg-zinc-950 text-zinc-200 selection:bg-sky-500/30">
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
