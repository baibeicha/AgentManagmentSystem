'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { UserPlus } from 'lucide-react';

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorSec, setErrorSec] = useState<string | null>(null);

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorSec(null);
    try {
      // Explicitly construct a plain object using only primitive string values
      // This strictly avoids any circular structure to JSON serialization
      const payload = JSON.stringify({ 
        email: email, 
        password: password 
      });
      console.log('Registration payload generated safely:', payload);
      router.push('/dashboard');
    } catch (err) {
      setErrorSec('An operational exception occurred during registration.');
      console.error(err);
    }
  };

  return (
    <div id="register-viewport" className="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-zinc-950 animate-in fade-in duration-700">
      <div id="register-card" className="w-full max-w-sm rounded-2xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-8 shadow-2xl backdrop-blur-md relative overflow-hidden">
        <div id="glow-header" className="absolute top-0 left-0 right-0 h-1 bg-amber-500/50 shadow-[0_0_15px_rgba(245,158,11,0.5)]" />
        
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-amber-500/10 border border-amber-500/20 shadow-inner">
            <UserPlus className="h-8 w-8 text-amber-500 drop-shadow-[0_0_8px_rgba(245,158,11,0.5)]" />
          </div>
          <h1 id="register-title" className="text-2xl font-black tracking-tighter text-zinc-100">Node Registration</h1>
          <p id="register-subtitle" className="mt-2 text-xs font-mono text-zinc-500">Request via /api/v1/auth/register</p>
        </div>

        {errorSec && (
          <div id="register-error-banner" className="mb-4 rounded border border-rose-500/20 bg-rose-500/5 p-3 text-xs text-rose-400 font-mono">
            {errorSec}
          </div>
        )}

        <form onSubmit={handleRegister} className="space-y-4">
          <div>
            <label className="mb-2 block text-[10px] font-bold uppercase tracking-wider text-zinc-500">Identity Email</label>
            <input 
              id="input-register-email"
              type="email" 
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="recruit@aegis.local"
              className="w-full rounded border border-white/10 bg-zinc-950 px-4 py-3 text-sm font-medium text-zinc-200 outline-none transition-all focus:border-amber-500 focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <div>
            <label className="mb-2 block text-[10px] font-bold uppercase tracking-wider text-zinc-500">Secure Passphrase</label>
            <input 
              id="input-register-passphrase"
              type="password" 
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••••••"
              className="w-full rounded border border-white/10 bg-zinc-950 px-4 py-3 text-sm font-medium text-zinc-200 outline-none transition-all focus:border-amber-500 focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <button 
            id="btn-register-submit"
            type="submit"
            className="mt-6 flex w-full items-center justify-center gap-2 rounded bg-amber-500 py-3 text-sm font-bold text-zinc-950 shadow-[0_0_15px_rgba(245,158,11,0.4)] transition-all hover:bg-amber-600"
          >
            Register New Node
          </button>
        </form>
        
        <div className="mt-6 text-center text-xs font-mono text-zinc-500">
          Already cleared? <Link href="/login" className="text-amber-500 hover:text-amber-400 ml-1">Authenticate here</Link>
        </div>
      </div>
    </div>
  );
}
