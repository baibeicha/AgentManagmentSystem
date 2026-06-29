'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { ShieldCheck, LogIn } from 'lucide-react';
import { api } from '@/lib/api';
import useAuthStore from '@/hooks/useAuth';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorSec, setErrorSec] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const setTokens = useAuthStore((state) => state.setTokens);
  const getDeviceId = useAuthStore((state) => state.getDeviceId);

  useEffect(() => {
    // Ensure device ID is initialized on mount
    getDeviceId();
  }, [getDeviceId]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorSec(null);
    setIsSubmitting(true);
    try {
      const deviceId = getDeviceId();
      const payload = {
        login: email,
        password: password,
        device_id: deviceId,
      };

      const { data } = await api.post('/api/v1/auth/login', payload);

      if (data.access_token && data.refresh_token) {
        setTokens(data.access_token, data.refresh_token);
        router.push('/dashboard');
      } else {
         setErrorSec('Invalid credentials or account locked.');
      }
    } catch (err: any) {
      if (err.response?.status === 401) {
        setErrorSec('Invalid credentials or account locked due to brute-force protection.');
      } else {
        setErrorSec('An operational exception occurred during authentication.');
      }
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div id="login-viewport" className="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-zinc-950 animate-in fade-in duration-700">
      <div id="login-card" className="w-full max-w-sm rounded-2xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-8 shadow-2xl backdrop-blur-md relative overflow-hidden">
        <div id="glow-header" className="absolute top-0 left-0 right-0 h-1 bg-sky-500/50 shadow-[0_0_15px_rgba(14,165,233,0.5)]" />
        
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-sky-500/10 border border-sky-500/20 shadow-inner">
            <ShieldCheck className="h-8 w-8 text-sky-400 drop-shadow-[0_0_8px_rgba(14,165,233,0.5)]" />
          </div>
          <h1 id="login-title" className="text-2xl font-black tracking-tighter text-zinc-100">AMS Login</h1>
          <p id="login-subtitle" className="mt-2 text-xs font-mono text-zinc-500">Authenticate via /api/v1/auth/login</p>
        </div>

        {errorSec && (
          <div id="login-error-banner" className="mb-4 rounded border border-rose-500/20 bg-rose-500/5 p-3 text-xs text-rose-400 font-mono">
            {errorSec}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="mb-2 block text-[10px] font-bold uppercase tracking-wider text-zinc-500">Email / Identity</label>
            <input 
              id="input-identity-email"
              type="text"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="ops@aegis.local"
              className="w-full rounded border border-white/10 bg-zinc-950 px-4 py-3 text-sm font-medium text-zinc-200 outline-none transition-all focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
            />
          </div>
          <div>
            <label className="mb-2 block text-[10px] font-bold uppercase tracking-wider text-zinc-500">Passphrase</label>
            <input 
              id="input-passphrase"
              type="password" 
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••••••"
              className="w-full rounded border border-white/10 bg-zinc-950 px-4 py-3 text-sm font-medium text-zinc-200 outline-none transition-all focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
            />
          </div>
          <button 
            id="btn-submit-session"
            type="submit"
            disabled={isSubmitting}
            className={`mt-6 flex w-full items-center justify-center gap-2 rounded bg-sky-500 py-3 text-sm font-bold text-white shadow-[0_0_15px_rgba(14,165,233,0.4)] transition-all hover:bg-sky-600 ${isSubmitting ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            <LogIn className="h-4 w-4" /> {isSubmitting ? 'Authenticating...' : 'Authenticate Session'}
          </button>
        </form>
        
        <div className="mt-6 text-center text-xs font-mono text-zinc-500">
          Unregistered node? <Link href="/register" className="text-sky-400 hover:text-sky-300 ml-1">Create Access Token</Link>
        </div>
      </div>
    </div>
  );
}
