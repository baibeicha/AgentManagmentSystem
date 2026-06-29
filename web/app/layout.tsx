import type { Metadata } from 'next';
import { Inter, JetBrains_Mono } from 'next/font/google';
import './globals.css';
import Providers from '@/components/providers';
import RootLayoutWrapper from '@/components/root-layout-wrapper';

const inter = Inter({ subsets: ['latin'], variable: '--font-sans' });
const jetbrainsMono = JetBrains_Mono({ subsets: ['latin'], variable: '--font-mono' });

export const metadata: Metadata = {
  title: 'Agent Management System',
  description: 'Enterprise monitoring and security platform',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable} dark`}>
      <body className="bg-zinc-950 text-zinc-100 font-sans antialiased selection:bg-sky-500/30" suppressHydrationWarning>
        <Providers>
          <RootLayoutWrapper>
            {children}
          </RootLayoutWrapper>
        </Providers>
      </body>
    </html>
  );
}
