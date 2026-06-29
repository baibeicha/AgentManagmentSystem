'use client';

import { useState, useEffect } from 'react';
import { X, CheckCircle2, AlertTriangle, Info, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

// Mock data strictly styled according to OpenAPI definitions
const MOCK_NOTIFICATIONS = [
  {
    id: 'notif-1',
    title: 'DEAD_MAN_SWITCH_TRIGGERED',
    message: 'Signal lost for SRV-CORE-01. Automated failover sequence initiated. Manual review required immediately.',
    severity: 'critical',
    is_read: false,
    created_at: new Date(Date.now() - 1000 * 60 * 2).toISOString(), // 2 mins ago
  },
  {
    id: 'notif-2',
    title: 'Host Offline',
    message: 'DB-REPLICA-03 unreachable. Ping timeout exceeded 3000ms. Retrying connection...',
    severity: 'warning',
    is_read: true,
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 3).toISOString(), // 3 hours ago
  },
  {
    id: 'notif-3',
    title: 'Agent Updated',
    message: 'v2.0.5 deployed to 12 nodes successfully. No anomalies detected during rollout.',
    severity: 'info',
    is_read: false,
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 12).toISOString(), // 12 hours ago
  },
];

const SEVERITY_CONFIG = {
  critical: {
    icon: AlertTriangle,
    colorClass: 'text-rose-500',
    bgClass: 'bg-rose-500/10',
    borderClass: 'border-rose-500/20',
  },
  warning: {
    icon: AlertCircle,
    colorClass: 'text-amber-500',
    bgClass: 'bg-amber-500/10',
    borderClass: 'border-amber-500/20',
  },
  info: {
    icon: Info,
    colorClass: 'text-zinc-400',
    bgClass: 'bg-zinc-800/50',
    borderClass: 'border-white/5',
  },
};

export default function NotificationsDrawer({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [filter, setFilter] = useState<'all' | 'unread' | 'critical'>('all');
  const [notifications, setNotifications] = useState(MOCK_NOTIFICATIONS);
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setIsMounted(true);
  }, []);

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  const filteredNotifications = notifications.filter((n) => {
    if (filter === 'unread') return !n.is_read;
    if (filter === 'critical') return n.severity === 'critical';
    return true;
  });

  const markAllAsRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })));
  };

  const formatRelativeTime = (isoString: string) => {
    // eslint-disable-next-line react-hooks/purity
    const diff = Date.now() - new Date(isoString).getTime();
    if (diff < 60000) return 'Just now';
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    return `${Math.floor(diff / 86400000)}d ago`;
  };

  return (
    <>
      {/* Backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm transition-opacity"
          onClick={onClose}
        />
      )}

      {/* Drawer */}
      <div
        className={cn(
          'fixed inset-y-0 right-0 z-50 flex w-96 flex-col bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] border-l border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),-10px_0_30px_rgba(0,0,0,0.5)] transition-transform duration-500 ease-out',
          isOpen ? 'translate-x-0' : 'translate-x-full'
        )}
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-white/5 p-6">
          <div className="flex items-center gap-3">
            <h2 className="text-xl font-bold text-zinc-100">System Notifications</h2>
            {unreadCount > 0 && (
              <span className="flex h-5 items-center justify-center rounded bg-sky-500/20 px-2 text-[11px] font-bold text-sky-400 border border-sky-500/30">
                {unreadCount} NEW
              </span>
            )}
          </div>
          <button
            onClick={onClose}
            className="rounded-full p-2 text-zinc-400 hover:bg-zinc-800 hover:text-zinc-100 transition-all active:scale-95"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Action Bar */}
        <div className="flex items-center justify-between border-b border-white/5 bg-zinc-900/30 px-6 py-3">
          <div className="flex gap-2">
            {(['all', 'unread', 'critical'] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setFilter(tab)}
                className={cn(
                  'rounded px-3 py-1 text-[11px] font-bold uppercase tracking-wider transition-all active:scale-95 border',
                  filter === tab
                    ? 'bg-zinc-800 border-zinc-700 text-zinc-100 shadow-sm'
                    : 'border-transparent text-zinc-500 hover:text-zinc-300'
                )}
              >
                {tab}
              </button>
            ))}
          </div>
          <button
            onClick={markAllAsRead}
            className="flex items-center gap-1.5 text-[11px] font-bold uppercase tracking-wider text-sky-400 hover:text-sky-300 transition-all active:scale-95"
          >
            <CheckCircle2 className="h-3.5 w-3.5" />
            Mark all Read
          </button>
        </div>

        {/* Scrollable List */}
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {filteredNotifications.length === 0 ? (
            <div className="flex h-full flex-col items-center justify-center text-zinc-500">
              <CheckCircle2 className="mb-2 h-8 w-8 opacity-20" />
              <p className="text-sm">No notifications found.</p>
            </div>
          ) : (
            filteredNotifications.map((notif) => {
              const config = SEVERITY_CONFIG[notif.severity as keyof typeof SEVERITY_CONFIG];
              const Icon = config.icon;

              return (
                <div
                  key={notif.id}
                  className={cn(
                    'group relative overflow-hidden rounded-lg border bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-4 transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01]',
                    !notif.is_read ? 'border-zinc-700' : 'border-white/[0.08]'
                  )}
                >
                  {/* Left accent strip */}
                  {!notif.is_read && (
                    <div
                      className={cn(
                        'absolute bottom-0 left-0 top-0 w-1',
                        notif.severity === 'critical' ? 'bg-rose-500' : 'bg-sky-500'
                      )}
                    />
                  )}

                  <div className="flex items-start gap-4">
                    <div
                      className={cn(
                        'flex h-8 w-8 shrink-0 items-center justify-center rounded border',
                        config.bgClass,
                        config.borderClass,
                        config.colorClass
                      )}
                    >
                      <Icon className="h-4 w-4" />
                    </div>
                    <div className="flex-1 space-y-1">
                      <div className="flex items-start justify-between gap-2">
                        <h3 className="font-mono text-sm font-semibold text-zinc-200">
                          {notif.title}
                        </h3>
                        <span className="shrink-0 text-[11px] font-medium text-zinc-500">
                          {isMounted ? formatRelativeTime(notif.created_at) : '...'}
                        </span>
                      </div>
                      <p className="text-sm text-zinc-400 leading-relaxed">
                        {notif.message}
                      </p>

                      {/* Action Triggers based on severity */}
                      {notif.severity === 'critical' && (
                        <div className="mt-3 flex gap-2">
                          <button className="rounded border border-rose-500/30 bg-rose-500/10 px-3 py-1.5 text-[11px] font-bold uppercase tracking-wider text-rose-400 transition-all duration-300 ease-out active:scale-95 hover:bg-rose-500/20">
                            Investigate Action
                          </button>
                        </div>
                      )}
                      {notif.severity === 'warning' && (
                        <div className="mt-3 flex gap-2">
                          <button className="rounded border border-zinc-700 bg-zinc-800 px-3 py-1.5 text-[11px] font-bold uppercase tracking-wider text-zinc-300 transition-all duration-300 ease-out active:scale-95 hover:bg-zinc-700">
                            View Logs
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>
    </>
  );
}
