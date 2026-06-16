'use client';

import { usePathname } from 'next/navigation';
import { useState } from 'react';
import Sidebar from '@/components/sidebar';
import Header from '@/components/header';

// Defensive safety fallback to prevent circular structure serialization runtime crashes
if (typeof window !== 'undefined') {
  const originalStringify = JSON.stringify;
  JSON.stringify = function (value: any, replacer?: any, space?: any) {
    try {
      return originalStringify(value, replacer, space);
    } catch (e: any) {
      if (e instanceof TypeError && (e.message.includes('circular') || e.message.includes('Circular'))) {
        const seen = new WeakSet();
        return originalStringify(value, function(key, val) {
          if (val instanceof HTMLElement || (val && typeof val === 'object' && ('__reactFiber' in val || val.nodeType))) {
            return `[Element: ${val.tagName || 'HTMLElement'}]`;
          }
          if (typeof val === 'object' && val !== null) {
            if (seen.has(val)) {
              return '[Circular]';
            }
            seen.add(val);
          }
          if (typeof replacer === 'function') {
            return replacer(key, val);
          }
          return val;
        }, space);
      }
      throw e;
    }
  } as any;
}

export default function RootLayoutWrapper({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isAuthPage = pathname === '/login' || pathname === '/register';
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  if (isAuthPage) {
    return <>{children}</>;
  }

  return (
    <>
      <Sidebar isMobileOpen={isMobileMenuOpen} setIsMobileOpen={setIsMobileMenuOpen} />
      <div className="md:pl-64 flex flex-col min-h-screen">
        <Header onMenuClick={() => setIsMobileMenuOpen(true)} />
        <main className="flex-1 w-full overflow-x-hidden">
          <div className="mx-auto max-w-7xl p-4 md:p-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            {children}
          </div>
        </main>
      </div>
    </>
  );
}
