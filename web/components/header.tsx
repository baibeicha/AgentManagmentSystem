'use client';

import { Search, Bell, User, Activity, Database, Zap, HardDrive, LogOut, Menu } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useState } from 'react';
import NotificationsDrawer from './notifications-drawer';
import useAuthStore from '@/hooks/useAuth';

function AdminSystemStatus() {
  return (
    <div className="hidden lg:flex items-center gap-4 rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-4 py-2">
      <div className="flex items-center gap-2">
        <ServerStatusIcon icon={Activity} status="online" label="Gateway" />
        <div className="h-4 w-px bg-white/10 mx-2" />
        <ServerStatusIcon icon={Database} status="online" label="PostgreSQL" />
        <ServerStatusIcon icon={Zap} status="online" label="Redis" />
        <ServerStatusIcon icon={HardDrive} status="degraded" label="Kafka" />
      </div>
    </div>
  );
}

function ServerStatusIcon({
  icon: Icon,
  status,
  label,
}: {
  icon: any;
  status: 'online' | 'degraded' | 'offline';
  label: string;
}) {
  return (
    <div className="group relative flex items-center justify-center">
      <Icon
        className={cn(
          'h-4 w-4',
          status === 'online' && 'text-emerald-500 drop-shadow-[0_0_5px_rgba(16,185,129,0.5)]',
          status === 'degraded' && 'text-amber-500 drop-shadow-[0_0_5px_rgba(245,158,11,0.5)]',
          status === 'offline' && 'text-rose-500 drop-shadow-[0_0_5px_rgba(226,29,72,0.5)]'
        )}
      />
      {/* Tooltip */}
      <div className="pointer-events-none absolute -bottom-8 left-1/2 -translate-x-1/2 whitespace-nowrap rounded bg-zinc-800 px-2 py-1 text-[10px] font-mono font-semibold text-zinc-200 opacity-0 transition-opacity group-hover:opacity-100 border border-white/10 z-50">
        {label}: {status.toUpperCase()}
      </div>
    </div>
  );
}

export default function Header({ onMenuClick }: { onMenuClick?: () => void }) {
  const [isNotificationsOpen, setIsNotificationsOpen] = useState(false);

  const logout = useAuthStore((state) => state.logout);
  const profile = useAuthStore((state) => state.profile);

  const handleLogout = async () => {
    await logout();
  };

  return (
    <>
      <header className="sticky top-0 z-40 w-full bg-zinc-900/40 backdrop-blur-xl border-b border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] flex h-16 items-center justify-between px-4 md:px-8">
        <div className="flex items-center gap-3">
          <button 
            className="md:hidden p-2 -ml-2 text-zinc-400 hover:text-zinc-100 transition-colors"
            onClick={onMenuClick}
          >
            <Menu className="h-6 w-6" />
          </button>
          <div className="flex items-center gap-4 text-lg md:text-xl font-black text-zinc-100 drop-shadow-md tracking-tighter">
            <span className="hidden sm:inline">Aegis XDR</span>
            <span className="sm:hidden">Aegis</span>
          </div>
        </div>

        <div className="relative mx-4 flex-1 max-w-xs md:max-w-md hidden sm:block">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-500" />
          <input
            type="text"
            placeholder="Search logs, hosts, or IP addresses..."
            className="w-full rounded-full border border-white/[0.08] bg-zinc-900/40 py-2 pl-10 pr-4 text-xs md:text-sm font-medium text-zinc-100 shadow-[inset_0_1px_0_rgba(255,255,255,0.05)] outline-none transition-all placeholder:text-zinc-500 focus:border-sky-500/50 focus:ring-1 focus:ring-sky-500/50 backdrop-blur-md"
          />
        </div>

        <div className="flex items-center gap-3 md:gap-6 text-zinc-400">
          {profile?.role === 'GLOBAL_ADMIN' && <AdminSystemStatus />}

          <button
            onClick={() => setIsNotificationsOpen(true)}
            className="group relative rounded-full p-2 outline-none transition-transform active:scale-95 hover:bg-zinc-800/50 hover:text-sky-400 focus:ring-1 focus:ring-sky-500/50"
          >
            <Bell className="h-4 w-4 md:h-5 md:w-5" />
            <span className="absolute right-2 top-2 h-1.5 w-1.5 md:h-2 md:w-2 rounded-full bg-rose-500 shadow-[0_0_10px_rgba(244,63,94,0.5)] animate-pulse" />
          </button>

          <button className="hidden sm:flex items-center justify-center h-8 w-8 rounded-full bg-zinc-900/40 border border-white/[0.08] hover:border-sky-500/50 transition-all active:scale-95 group overflow-hidden">
            <span className="text-xs font-bold text-zinc-300 group-hover:text-sky-400">
              {profile?.login ? profile.login.substring(0, 3).toUpperCase() : 'SYS'}
            </span>
          </button>
          
          <div className="hidden sm:block h-6 w-px bg-white/10" />

          <button 
            onClick={handleLogout}
            className="flex items-center justify-center rounded-full p-2 outline-none transition-colors hover:bg-rose-500/10 hover:text-rose-400 text-zinc-500"
            title="Secure Logout"
          >
            <LogOut className="h-4 w-4" />
          </button>
        </div>
      </header>
      <NotificationsDrawer
        isOpen={isNotificationsOpen}
        onClose={() => setIsNotificationsOpen(false)}
      />
    </>
  );
}
