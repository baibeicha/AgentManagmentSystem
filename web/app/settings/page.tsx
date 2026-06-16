'use client';

import { useState } from 'react';
import { 
  User, ShieldAlert, Key, BellRing, History, Users, X, QrCode, Lock, Send, Edit
} from 'lucide-react';
import { cn } from '@/lib/utils';

type GlobalRole = 'GLOBAL_ADMIN' | 'TEAM_ADMIN' | 'USER';
type Permission = 'Metrics:View' | 'Terminal:Execute' | 'Incidents:Manage';
const ALL_PERMISSIONS: Permission[] = ['Metrics:View', 'Terminal:Execute', 'Incidents:Manage'];

type AppUser = {
  id: number;
  email: string;
  role: GlobalRole;
  permissions: Permission[];
  lastLogin: string;
};

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<'profile' | 'rbac' | 'notifications' | 'audit'>('profile');
  const [is2FAModalOpen, setIs2FAModalOpen] = useState(false);
  const [passwordForm, setPasswordForm] = useState({ current: '', new: '' });
  const [passwordStatus, setPasswordStatus] = useState<'idle' | 'success'>('idle');

  // RBAC State
  const [users, setUsers] = useState<AppUser[]>([
    { id: 1, email: 'admin@aegis-os.local', role: 'GLOBAL_ADMIN', permissions: ['Metrics:View', 'Terminal:Execute', 'Incidents:Manage'], lastLogin: 'Current Session' },
    { id: 2, email: 'l2-support@aegis-os.local', role: 'TEAM_ADMIN', permissions: ['Metrics:View', 'Incidents:Manage'], lastLogin: '2 hours ago' },
    { id: 3, email: 'read-only-dev@aegis-os.local', role: 'USER', permissions: ['Metrics:View'], lastLogin: '1 day ago' }
  ]);

  const [isUserModalOpen, setIsUserModalOpen] = useState(false);
  const [editingUserId, setEditingUserId] = useState<number | null>(null);
  const [userForm, setUserForm] = useState<{ email: string; password: string; role: GlobalRole; permissions: Permission[] }>({
    email: '',
    password: '',
    role: 'USER',
    permissions: ['Metrics:View']
  });

  const handlePasswordChange = (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordStatus('success');
    setTimeout(() => {
      setPasswordStatus('idle');
      setPasswordForm({ current: '', new: '' });
    }, 3000);
  };

  const openUserModal = (user?: AppUser) => {
    if (user) {
      setEditingUserId(user.id);
      setUserForm({ email: user.email, password: '', role: user.role, permissions: user.permissions });
    } else {
      setEditingUserId(null);
      setUserForm({ email: '', password: '', role: 'USER', permissions: ['Metrics:View'] });
    }
    setIsUserModalOpen(true);
  };

  const handleUserSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingUserId) {
      setUsers(users.map(u => u.id === editingUserId ? { ...u, role: userForm.role, permissions: userForm.permissions } : u));
    } else {
      const newUser: AppUser = {
        id: Date.now(),
        email: userForm.email,
        role: userForm.role,
        permissions: userForm.permissions,
        lastLogin: 'Never'
      };
      setUsers([...users, newUser]);
    }
    setIsUserModalOpen(false);
  };

  const togglePermission = (perm: Permission) => {
    setUserForm(prev => {
      const has = prev.permissions.includes(perm);
      return {
        ...prev,
        permissions: has ? prev.permissions.filter(p => p !== perm) : [...prev.permissions, perm]
      };
    });
  };

  const handleRevokeUser = (id: number) => {
    setUsers(users.filter(u => u.id !== id));
  };

  return (
    <div className="flex h-[calc(100vh-8rem)] flex-col space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header */}
      <div className="border-b border-white/5 pb-6">
        <h1 className="text-3xl font-black tracking-tighter text-zinc-100 mb-2">System Configuration</h1>
        <p className="text-sm text-zinc-500">Manage identity, access control, and platform integrations.</p>
      </div>

      <div className="flex flex-1 gap-8 overflow-hidden relative">
        
        {/* Sidebar Nav */}
        <div className="w-64 shrink-0 overflow-y-auto pr-4 space-y-1">
          <button
            onClick={() => setActiveTab('profile')}
            className={cn(
              "flex w-full items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold transition-all duration-200",
              activeTab === 'profile' 
                ? "bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-zinc-100 shadow-sm" 
                : "text-zinc-400 hover:bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] hover:text-zinc-200"
            )}
          >
            <User className={cn("h-4 w-4", activeTab === 'profile' ? "text-sky-400" : "")} />
            Profile & Security
          </button>
          
          <button
            onClick={() => setActiveTab('rbac')}
            className={cn(
              "flex w-full items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold transition-all duration-200",
              activeTab === 'rbac' 
                ? "bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-zinc-100 shadow-sm" 
                : "text-zinc-400 hover:bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] hover:text-zinc-200"
            )}
          >
            <Users className={cn("h-4 w-4", activeTab === 'rbac' ? "text-amber-400" : "")} />
            Team & RBAC
          </button>

          <button
            onClick={() => setActiveTab('notifications')}
            className={cn(
              "flex w-full items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold transition-all duration-200",
              activeTab === 'notifications' 
                ? "bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-zinc-100 shadow-sm" 
                : "text-zinc-400 hover:bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] hover:text-zinc-200"
            )}
          >
            <BellRing className={cn("h-4 w-4", activeTab === 'notifications' ? "text-emerald-400" : "")} />
            Dispatch Channels
          </button>

          <button
            onClick={() => setActiveTab('audit')}
            className={cn(
              "flex w-full items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold transition-all duration-200",
              activeTab === 'audit' 
                ? "bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-zinc-100 shadow-sm" 
                : "text-zinc-400 hover:bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] hover:text-zinc-200"
            )}
          >
            <History className={cn("h-4 w-4", activeTab === 'audit' ? "text-rose-400" : "")} />
            Audit & System Logs
          </button>
        </div>

        {/* Content Area */}
        <div className="flex-1 overflow-y-auto rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-8 shadow-sm relative">
          
          {activeTab === 'profile' && (
            <div className="max-w-2xl space-y-8 animate-in fade-in duration-300">
              <div>
                <h2 className="text-xl font-bold text-zinc-100 mb-6 flex items-center gap-2">
                  <User className="h-5 w-5 text-sky-400" /> Identity Settings
                </h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Email Address</label>
                    <input type="email" defaultValue="admin@aegis-os.local" className="w-full rounded border border-white/10 bg-zinc-950 px-4 py-2 text-sm text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all" />
                  </div>
                </div>
              </div>

              <div className="pt-6 border-t border-white/5">
                <h2 className="text-xl font-bold text-zinc-100 mb-6 flex items-center gap-2">
                  <Lock className="h-5 w-5 text-sky-400" /> Password Escalation
                </h2>
                <form onSubmit={handlePasswordChange} className="rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-$1 space-y-4 shadow-sm">
                  <div>
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Current Password</label>
                    <input 
                      type="password" 
                      required
                      value={passwordForm.current}
                      onChange={(e) => setPasswordForm(prev => ({ ...prev, current: e.target.value }))}
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2 text-sm text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all" 
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">New Password</label>
                    <input 
                      type="password" 
                      required
                      value={passwordForm.new}
                      onChange={(e) => setPasswordForm(prev => ({ ...prev, new: e.target.value }))}
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2 text-sm text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all" 
                    />
                  </div>
                  <div className="pt-2 flex items-center justify-between">
                    <button type="submit" className="rounded bg-sky-500 px-6 py-2 text-sm font-bold text-white hover:bg-sky-600 transition-all duration-300 ease-out active:scale-95 shadow-[0_0_15px_rgba(14,165,233,0.3)]">
                      Change Password
                    </button>
                    {passwordStatus === 'success' && (
                      <span className="text-emerald-500 animate-pulse text-sm font-bold animate-pulse">Password Updated!</span>
                    )}
                  </div>
                </form>
              </div>

              <div className="pt-6 border-t border-white/5">
                <h2 className="text-xl font-bold text-zinc-100 mb-6 flex items-center gap-2">
                  <Key className="h-5 w-5 text-sky-400" /> Authentication
                </h2>
                <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-4 flex items-center justify-between transition-colors hover:bg-emerald-500/10">
                  <div>
                    <div className="font-bold text-zinc-200 text-sm">Two-Factor Authentication (2FA)</div>
                    <div className="text-xs text-zinc-500 mt-1">Configured via TOTP (Google Authenticator)</div>
                  </div>
                  <button 
                    onClick={() => setIs2FAModalOpen(true)}
                    className="rounded bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-4 py-2 text-xs font-bold text-zinc-300 hover:bg-zinc-700 transition-colors"
                  >
                    Reconfigure
                  </button>
                </div>
              </div>

            </div>
          )}

          {activeTab === 'rbac' && (
             <div className="space-y-6 animate-in fade-in duration-300">
               <div className="flex items-center justify-between mb-6">
                <div>
                  <h2 className="text-xl font-bold text-zinc-100 flex items-center gap-2">
                    <ShieldAlert className="h-5 w-5 text-amber-400" /> Access Management
                  </h2>
                  <p className="text-sm text-zinc-500 mt-1">Control who has access to infrastructure commands.</p>
                </div>
                <button 
                  onClick={() => openUserModal()}
                  className="flex items-center gap-2 rounded bg-amber-500/10 px-4 py-2 text-sm font-bold text-amber-500 hover:bg-amber-500/20 transition-colors border border-amber-500/20"
                >
                  Add New User
                </button>
              </div>

              <div className="rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] overflow-hidden">
                <table className="w-full text-left text-sm text-zinc-400">
                  <thead className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-xs font-semibold uppercase tracking-wider text-zinc-500 border-b border-white/5">
                    <tr>
                      <th className="px-6 py-4">User</th>
                      <th className="px-6 py-4">Role</th>
                      <th className="px-6 py-4">Permissions</th>
                      <th className="px-6 py-4">Last Login</th>
                      <th className="px-6 py-4 text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-white/5">
                    {users.map(user => (
                      <tr key={user.id} className="transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95">
                        <td className="px-6 py-4 font-bold text-zinc-200">{user.email}</td>
                        <td className="px-6 py-4">
                          <span className={cn(
                            "rounded border px-2 py-1 font-mono text-[10px] uppercase",
                            user.role === 'GLOBAL_ADMIN' ? "bg-amber-500/10 border-amber-500/20 text-amber-400" :
                            user.role === 'TEAM_ADMIN' ? "bg-sky-500/10 border-sky-500/20 text-sky-400" :
                            "bg-zinc-800 border-white/10 text-zinc-300"
                          )}>
                            {user.role}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex flex-wrap gap-1">
                            {user.permissions.map(p => (
                              <span key={p} className="text-[9px] bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-1.5 py-0.5 rounded text-zinc-400 uppercase tracking-wider">{p}</span>
                            ))}
                          </div>
                        </td>
                        <td className="px-6 py-4 text-xs font-mono">{user.lastLogin}</td>
                        <td className="px-6 py-4 text-right">
                          <div className="flex justify-end gap-3">
                            <button 
                              onClick={() => openUserModal(user)}
                              className="flex items-center gap-1 text-xs text-sky-400 hover:text-sky-300 font-bold transition-colors"
                            >
                              <Edit className="w-3 h-3" /> Edit
                            </button>
                            <button 
                              onClick={() => handleRevokeUser(user.id)}
                              className="text-xs text-rose-400 hover:text-rose-300 font-bold transition-colors"
                            >
                              Revoke
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                    {users.length === 0 && (
                      <tr>
                        <td colSpan={5} className="px-6 py-4 text-center text-zinc-500 text-sm">
                          No users found.
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
             </div>
          )}

          {activeTab === 'notifications' && (
            <div className="max-w-2xl space-y-6 animate-in fade-in duration-300">
               <div>
                  <h2 className="text-xl font-bold text-zinc-100 flex items-center gap-2 mb-6">
                    <Send className="h-5 w-5 text-emerald-400" /> Platform Dispatch
                  </h2>
                </div>

                <div className="space-y-4">
                  <div className="rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-$1 group transition-colors hover:border-blue-500/30">
                    <div className="flex items-start justify-between">
                      <div>
                        <h3 className="font-bold text-zinc-200 flex items-center gap-2">Telegram Setup <span className="rounded bg-blue-500/10 px-1.5 py-0.5 text-[9px] font-bold uppercase text-blue-500 border border-blue-500/20">Active</span></h3>
                        <p className="text-xs text-zinc-500 mt-1">Routes alerts securely to configured Telegram ID</p>
                      </div>
                      <button className="text-xs font-bold text-sky-400 hover:text-sky-300 transition-all duration-300 ease-out active:scale-95">Configure</button>
                    </div>
                  </div>
                </div>
            </div>
          )}

          {activeTab === 'audit' && (
            <div className="flex flex-col h-full animate-in fade-in duration-300">
               <div>
                  <h2 className="text-xl font-bold text-zinc-100 flex items-center gap-2 mb-6">
                    <History className="h-5 w-5 text-rose-400" /> Immutable Audit Trail
                  </h2>
                  <p className="text-sm text-zinc-500 mb-6">Secure log of all administrative actions and system events. Retained for 365 days.</p>
                </div>
                
                <div className="flex-1 overflow-y-auto rounded border border-white/5 bg-[#0a0a0c] p-4 text-xs font-mono text-zinc-400 space-y-2 shadow-inner">
                  <div>[2024-03-10T10:15:00Z] <span className="text-amber-400">WARN</span> Rule &quot;Network Drop&quot; muted by test@aegis-os.local</div>
                  <div>[2024-03-10T09:30:12Z] <span className="text-emerald-400">INFO</span> Agent v2.0.4 deployed to host <span className="text-zinc-300">dev-104</span> by admin@aegis-os.local</div>
                  <div>[2024-03-10T08:45:00Z] <span className="text-sky-400">AUTH</span> Successful SSH proxy session established to <span className="text-zinc-300">db-master-01</span></div>
                  <div>[2024-03-09T22:15:00Z] <span className="text-emerald-400">INFO</span> Fleet backup snapshots completed across 12 managed nodes</div>
                  <div>[2024-03-09T18:02:11Z] <span className="text-rose-400">CRIT</span> Failed login attempt from IP 192.168.1.100 (reason: invalid 2FA)</div>
                  <div>[2024-03-09T15:20:00Z] <span className="text-emerald-400">INFO</span> Configuration &quot;notification.slack&quot; updated by admin@aegis-os.local</div>
                </div>
            </div>
          )}

        </div>
      </div>

      {/* 2FA Setup Modal Overlay */}
      {is2FAModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="w-full max-w-md rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),-10px_0_30px_rgba(0,0,0,0.5)] overflow-hidden flex flex-col">
            <div className="flex items-center justify-between border-b border-white/[0.08] p-4 bg-zinc-900/40 backdrop-blur-xl shadow-[inset_0_1px_0_rgba(255,255,255,0.1)]">
              <h3 className="font-bold text-zinc-100 flex items-center gap-2">
                <QrCode className="h-5 w-5 text-sky-400" />
                Setup 2FA (TOTP)
              </h3>
              <button 
                onClick={() => setIs2FAModalOpen(false)}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="p-6 space-y-6 flex-1 bg-transparent">
              <p className="text-sm text-zinc-400">
                Scan the QR code below with an authenticator app (e.g., Google Authenticator, Authy).
              </p>
              
              <div className="flex justify-center">
                {/* Mock QR Code Pattern */}
                <div className="h-48 w-48 bg-white rounded-xl flex items-center justify-center p-4 shadow-inner">
                  <div className="w-full h-full border-8 border-black relative grid grid-cols-5 grid-rows-5 gap-1 p-2 bg-white">
                     {Array.from({ length: 25 }).map((_, i) => (
                       <div key={i} className={cn("bg-black", (i % 2 === 0 || i % 7 === 0) ? "opacity-100" : "opacity-0")} />
                     ))}
                     <div className="absolute top-1 left-1 w-8 h-8 border-[6px] border-black bg-white" />
                     <div className="absolute top-1 right-1 w-8 h-8 border-[6px] border-black bg-white" />
                     <div className="absolute bottom-1 left-1 w-8 h-8 border-[6px] border-black bg-white" />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Confirmation Token</label>
                <input 
                  type="text" 
                  placeholder="000000"
                  maxLength={6}
                  className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-3 text-center text-lg font-mono tracking-[0.5em] text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all" 
                />
              </div>
            </div>
            <div className="border-t border-white/5 p-4 bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] flex justify-end gap-3">
              <button 
                onClick={() => setIs2FAModalOpen(false)}
                className="rounded px-4 py-2 text-sm font-bold text-zinc-400 hover:text-zinc-200 transition-colors"
              >
                Cancel
              </button>
              <button 
                onClick={() => setIs2FAModalOpen(false)}
                className="rounded bg-sky-500 px-6 py-2 text-sm font-bold text-white hover:bg-sky-600 transition-colors shadow-lg shadow-sky-500/20"
              >
                Verify & Save
              </button>
            </div>
          </div>
        </div>
      )}

      {/* User / RBAC Modal Overlay */}
      {isUserModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="w-full max-w-lg rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),-10px_0_30px_rgba(0,0,0,0.5)] overflow-hidden flex flex-col">
            <div className="flex items-center justify-between border-b border-white/[0.08] p-4 bg-zinc-900/40 backdrop-blur-xl shadow-[inset_0_1px_0_rgba(255,255,255,0.1)]">
              <h3 className="font-bold text-zinc-100 flex items-center gap-2">
                <Users className="h-5 w-5 text-amber-400" />
                {editingUserId ? 'Edit User Roles & Permissions' : 'Add New User'}
              </h3>
              <button 
                onClick={() => setIsUserModalOpen(false)}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
                type="button"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            
            <form onSubmit={handleUserSubmit}>
              <div className="p-6 space-y-6 flex-1 bg-transparent">
                <div className="grid grid-cols-2 gap-4">
                  <div className="col-span-2">
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Email Address</label>
                    <input 
                      type="email" 
                      required
                      disabled={!!editingUserId}
                      value={userForm.email}
                      onChange={e => setUserForm({ ...userForm, email: e.target.value })}
                      placeholder="user@aegis-os.local"
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2.5 text-sm text-zinc-200 focus:border-amber-500 focus:ring-1 focus:ring-amber-500 outline-none transition-all disabled:opacity-50" 
                    />
                  </div>
                  
                  {!editingUserId && (
                    <div className="col-span-2">
                      <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Initial Password</label>
                      <input 
                        type="password" 
                        required
                        value={userForm.password}
                        onChange={e => setUserForm({ ...userForm, password: e.target.value })}
                        placeholder="••••••••"
                        className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2.5 text-sm text-zinc-200 focus:border-amber-500 focus:ring-1 focus:ring-amber-500 outline-none transition-all" 
                      />
                    </div>
                  )}

                  <div className="col-span-2">
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Global Role</label>
                    <select
                      value={userForm.role}
                      onChange={e => setUserForm({ ...userForm, role: e.target.value as GlobalRole })}
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2.5 text-sm text-zinc-200 focus:border-amber-500 focus:ring-1 focus:ring-amber-500 outline-none transition-all appearance-none"
                    >
                      <option value="USER">User (Standard Access)</option>
                      <option value="TEAM_ADMIN">Team Admin</option>
                      <option value="GLOBAL_ADMIN">Global Admin</option>
                    </select>
                  </div>
                  
                  <div className="col-span-2 mt-2">
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-3">Custom Permissions</label>
                    <div className="space-y-3 p-4 bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg border border-white/5">
                      {ALL_PERMISSIONS.map(perm => (
                        <label key={perm} className="flex items-center gap-3 cursor-pointer group">
                          <div className="relative flex items-center justify-center w-4 h-4">
                            <input 
                              type="checkbox" 
                              checked={userForm.permissions.includes(perm)}
                              onChange={() => togglePermission(perm)}
                              className="peer appearance-none w-4 h-4 border border-white/20 rounded bg-zinc-950 checked:bg-amber-500 checked:border-amber-500 transition-colors"
                            />
                            <svg className="absolute w-3 h-3 text-zinc-950 opacity-0 peer-checked:opacity-100 pointer-events-none" viewBox="0 0 14 10" fill="none">
                              <path d="M1 5L4.5 8.5L13 1" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            </svg>
                          </div>
                          <span className="text-sm font-medium text-zinc-300 group-hover:text-zinc-100 transition-colors">
                            {perm}
                          </span>
                        </label>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t border-white/5 p-4 bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] flex justify-end gap-3 rounded-b-xl">
                <button 
                  type="button"
                  onClick={() => setIsUserModalOpen(false)}
                  className="rounded px-4 py-2 text-sm font-bold text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="rounded bg-amber-500 px-6 py-2 text-sm font-bold text-zinc-950 hover:bg-amber-400 transition-all duration-300 ease-out active:scale-95 shadow-lg shadow-amber-500/20"
                >
                  {editingUserId ? 'Save Changes' : 'Create User'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
