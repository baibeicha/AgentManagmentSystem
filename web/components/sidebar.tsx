'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import {
  LayoutDashboard,
  Server,
  TerminalSquare,
  Settings2,
  AlertTriangle,
  Settings,
  X
} from 'lucide-react';

const navItems = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Inventory & Discovery', href: '/inventory', icon: Server },
  { name: 'Command & Control', href: '/command-and-control', icon: TerminalSquare },
  { name: 'Automation', href: '/automation', icon: Settings2 },
  { name: 'Incidents', href: '/incidents', icon: AlertTriangle },
];

interface SidebarProps {
  isMobileOpen?: boolean;
  setIsMobileOpen?: (open: boolean) => void;
}

export default function Sidebar({ isMobileOpen, setIsMobileOpen }: SidebarProps) {
  const pathname = usePathname();

  return (
    <>
      {/* Mobile backdrop */}
      {isMobileOpen && (
        <div 
          className="md:hidden fixed inset-0 z-40 bg-black/60 backdrop-blur-sm"
          onClick={() => setIsMobileOpen?.(false)}
        />
      )}
      <nav className={cn(
        "fixed left-0 top-0 h-screen w-64 bg-zinc-900/40 backdrop-blur-xl border-r border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),inset_-1px_0_10px_rgba(255,255,255,0.02)] flex flex-col p-6 z-50 transition-transform duration-300 md:translate-x-0",
        isMobileOpen ? "translate-x-0" : "-translate-x-full"
      )}>
        <div className="mb-8 flex items-start justify-between">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-bold tracking-tighter text-sky-400 drop-shadow-[0_0_15px_rgba(14,165,233,0.3)]">
              AEGIS_OS
            </h1>
            <span className="font-mono text-xs text-zinc-500">V2.0.4-STABLE</span>
          </div>
          {isMobileOpen && (
            <button 
              className="md:hidden p-1 text-zinc-400 hover:text-zinc-100 transition-colors"
              onClick={() => setIsMobileOpen?.(false)}
            >
              <X className="h-5 w-5" />
            </button>
          )}
        </div>

        <div className="flex-1 space-y-2 overflow-y-auto pr-2">
          {navItems.map((item) => {
            const isActive = pathname.startsWith(item.href);
            return (
              <Link
                key={item.name}
                href={item.href}
                onClick={() => setIsMobileOpen?.(false)}
                className={cn(
                  'group flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold tracking-wider uppercase transition-all duration-300 ease-out active:scale-95 hover:bg-white/[0.02] hover:scale-[1.01]',
                  isActive
                    ? 'border-r-2 border-sky-500 bg-sky-500/10 text-sky-400 shadow-[inset_0_0_20px_rgba(14,165,233,0.1)]'
                    : 'text-zinc-400 hover:text-zinc-200'
                )}
              >
                <item.icon
                  className={cn(
                    'h-5 w-5 transition-all duration-300',
                    isActive ? 'drop-shadow-[0_0_8px_rgba(14,165,233,0.5)]' : ''
                  )}
                />
                <span className="text-[11px]">{item.name}</span>
              </Link>
            );
          })}
        </div>

        <div className="mt-auto pt-4">
          <Link
            href="/settings"
            onClick={() => setIsMobileOpen?.(false)}
            className={cn(
              'group flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold tracking-wider uppercase transition-all duration-300 ease-out active:scale-95 hover:bg-white/[0.02] hover:scale-[1.01]',
              pathname.startsWith('/settings')
                ? 'border-r-2 border-sky-500 bg-sky-500/10 text-sky-400 shadow-[inset_0_0_20px_rgba(14,165,233,0.1)]'
                : 'text-zinc-400 hover:text-zinc-200'
            )}
          >
            <Settings className="h-5 w-5" />
            <span className="text-[11px]">Settings</span>
          </Link>
        </div>
      </nav>
    </>
  );
}
