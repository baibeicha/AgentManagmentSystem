'use client';

import { useState, useRef, useEffect } from 'react';
import { 
  TerminalSquare, Play, Plus, Terminal, Code, ChevronDown, X
} from 'lucide-react';
import { cn } from '@/lib/utils';

// --- MOCKS ---
const INITIAL_SCRIPTS = [
  { 
    id: 'scr-1', 
    name: 'clean_docker_cache.sh', 
    lang: 'bash',
    code: `#!/bin/bash\n\necho "[INFO] Engaging Docker system prune..."\ndocker system prune -a --volumes -f\n\necho "[INFO] Validating disk space..."\ndf -h /\n\necho "[SUCCESS] Operation clean_docker_cache completed successfully."`,
  },
  { 
    id: 'scr-2', 
    name: 'restart_nginx_graceful.sh', 
    lang: 'bash',
    code: `#!/bin/bash\n\n# Test configuration first\nnginx -t\nif [ $? -eq 0 ]; then\n    echo "Configuration valid. Reloading..."\n    systemctl reload nginx\n    systemctl status nginx --no-pager\nelse\n    echo "[ERROR] Nginx configuration is invalid. Aborting reload."\n    exit 1\nfi`,
  },
];

const MOCK_DEVICES = [
  { id: 'nyc-edge-01', ip: '10.0.1.1' },
  { id: 'db-master-01', ip: '10.0.2.1' },
  { id: 'worker-node-12', ip: '10.0.3.12' },
];

export default function CommandAndControlPage() {
  const [scripts, setScripts] = useState(INITIAL_SCRIPTS);
  const [selectedDevice, setSelectedDevice] = useState(MOCK_DEVICES[0]);
  const [selectedScript, setSelectedScript] = useState<typeof INITIAL_SCRIPTS[0] | null>(null);
  const [terminalInput, setTerminalInput] = useState('');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [newScriptForm, setNewScriptForm] = useState({ name: '', lang: 'bash', code: '' });
  
  const [terminalHistory, setTerminalHistory] = useState<{ type: 'input' | 'output' | 'system', text: string, isError?: boolean }[]>([
    { type: 'system', text: 'Welcome to Aegis OS Shell v2.0.4. Connection established.' },
    { type: 'system', text: 'Authentication successful. Type "help" for available commands.' }
  ]);

  const terminalEndRef = useRef<HTMLDivElement>(null);

  const handleCommandSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!terminalInput.trim()) return;

    setTerminalHistory(prev => [...prev, { type: 'input', text: terminalInput }]);

    const cmd = terminalInput.trim().toLowerCase();
    setTerminalInput('');

    setTimeout(() => {
      if (cmd === 'help') {
        setTerminalHistory(prev => [...prev, { type: 'output', text: 'Available mock commands: help, clear, systemctl restart nginx, ping' }]);
      } else if (cmd === 'clear') {
        setTerminalHistory([{ type: 'system', text: 'Terminal cleared.' }]);
      } else if (cmd.includes('restart nginx')) {
        setTerminalHistory(prev => [...prev, { type: 'output', text: '[OK] Service restarted successfully.\nExit code: 0' }]);
      } else {
         setTerminalHistory(prev => [...prev, { type: 'output', text: `bash: ${cmd}: command not found`, isError: true }]);
      }
    }, 400);
  };

  const handleCreateScript = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newScriptForm.name || !newScriptForm.code) return;
    
    const newScript = {
      id: `scr-${Date.now()}`,
      name: newScriptForm.name,
      lang: newScriptForm.lang,
      code: newScriptForm.code
    };
    
    setScripts(prev => [...prev, newScript]);
    setIsCreateModalOpen(false);
    setSelectedScript(newScript);
    setNewScriptForm({ name: '', lang: 'bash', code: '' });
  };

  useEffect(() => {
    if (!selectedScript) {
       terminalEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [terminalHistory, selectedScript]);

  return (
    <div className="flex h-[calc(100vh-8rem)] flex-col space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header */}
      <div className="flex items-end justify-between border-b border-white/5 pb-6">
        <div>
          <h1 className="text-3xl font-black tracking-tighter text-zinc-100 mb-2 whitespace-nowrap">Command & Control</h1>
          <p className="text-sm text-zinc-500">Live bash shell and operational script execution library.</p>
        </div>
      </div>

      {/* Main Content Areas */}
      <div className="flex-1 overflow-hidden relative pb-4">
        
        <div className="flex h-full gap-6">
          
          {/* Left: Script Library */}
          <div className="w-1/3 flex flex-col rounded-xl border border-white/5 bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] overflow-hidden text-sm">
            <div className="flex items-center justify-between border-b border-white/5 bg-zinc-900 px-4 py-3 shrink-0">
              <span className="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center gap-2">
                <Code className="h-4 w-4" /> Template Library
              </span>
              <button onClick={() => setIsCreateModalOpen(true)} className="flex items-center gap-1 rounded bg-sky-500/10 border border-sky-500/20 px-2 py-1 text-[10px] font-bold text-sky-400 hover:bg-sky-500/20 transition-colors">
                <Plus className="h-3 w-3" /> New
              </button>
            </div>
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {scripts.map(script => (
                <div 
                  key={script.id} 
                  onClick={() => setSelectedScript(script)}
                  className={cn(
                    "group cursor-pointer rounded-lg border p-4 transition-all shadow-sm",
                    selectedScript?.id === script.id 
                      ? "bg-zinc-900 border-sky-500/50 shadow-[inset_0_0_15px_rgba(14,165,233,0.1)]" 
                      : "bg-zinc-950 border-white/5 hover:border-sky-500/30 hover:bg-zinc-900"
                  )}
                >
                  <div className="flex items-start justify-between mb-2">
                    <div className="font-semibold text-zinc-200 text-sm overflow-hidden text-ellipsis whitespace-nowrap">{script.name}</div>
                  </div>
                  <div className="flex items-center justify-between mt-4">
                    <span className="rounded bg-zinc-800 px-2 py-0.5 font-mono text-[10px] text-zinc-400 border border-white/5 uppercase">
                      {script.lang}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right: Interactive Terminal or Script Editor */}
          <div className="flex-1 flex flex-col rounded-xl border border-white/5 bg-[#0a0a0c] backdrop-blur-sm shadow-sm overflow-hidden font-mono shadow-inner relative">
            
            {/* Header: Device Selector */}
            <div className="flex items-center justify-between border-b border-white/5 bg-zinc-950 px-4 py-2 absolute top-0 left-0 right-0 z-10 shadow-sm shrink-0">
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-2 text-zinc-500 text-xs">
                  <Terminal className="h-4 w-4 text-emerald-400 drop-shadow-[0_0_5px_rgba(16,185,129,0.5)]" />
                  <span>P2P Link: <span className="text-emerald-500 font-bold">ESTABLISHED</span></span>
                </div>
                
                <div className="relative group flex items-center gap-2">
                  <span className="text-xs text-zinc-600">Target:</span>
                  <select 
                    value={selectedDevice.id}
                    onChange={(e) => {
                      const dev = MOCK_DEVICES.find(d => d.id === e.target.value);
                      if (dev) {
                        setSelectedDevice(dev);
                        setTerminalHistory([{ type: 'system', text: `Session switched to ${dev.id} (${dev.ip}).` }]);
                      }
                    }}
                    className="appearance-none bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-3 py-1 pr-8 text-xs font-bold text-sky-400 focus:outline-none focus:border-sky-500 cursor-pointer shadow-inner"
                  >
                    {MOCK_DEVICES.map(dev => (
                      <option key={dev.id} value={dev.id}>{dev.id}</option>
                    ))}
                  </select>
                  <ChevronDown className="absolute right-2 top-1/2 -translate-y-1/2 h-3 w-3 text-sky-400 pointer-events-none" />
                </div>
              </div>

              <div className="flex items-center gap-4">
                {selectedScript && (
                  <button 
                    onClick={() => setSelectedScript(null)}
                    className="text-[10px] font-bold text-zinc-400 hover:text-zinc-200 uppercase tracking-wider bg-zinc-800 hover:bg-zinc-700 px-2 py-1 rounded transition-colors"
                  >
                    Close Editor
                  </button>
                )}
                <div className="flex gap-2">
                  <span className="w-3 h-3 rounded-full bg-rose-500/80 border border-rose-500" />
                  <span className="w-3 h-3 rounded-full bg-amber-500/80 border border-amber-500" />
                  <span className="w-3 h-3 rounded-full bg-emerald-500/80 border border-emerald-500" />
                </div>
              </div>
            </div>

            {/* Viewport Content */}
            <div className="flex-1 mt-[45px] relative">
              
              {selectedScript ? (
                <div className="absolute inset-0 flex flex-col p-4 animate-in fade-in duration-300 bg-[#0a0a0c]">
                  <div className="flex items-center justify-between mb-4">
                      <span className="text-zinc-300 font-bold text-sm bg-zinc-900 px-3 py-1 rounded-md inline-block border border-white/10">{selectedScript.name}</span>
                      <button className="flex items-center gap-2 rounded bg-sky-500/10 border border-sky-500/30 px-4 py-2 text-xs font-bold uppercase tracking-wider text-sky-400 hover:bg-sky-500/20 transition-all shadow-[0_0_10px_rgba(14,165,233,0.1)]">
                        <Play className="h-4 w-4" /> Run on {selectedDevice.id}
                      </button>
                  </div>
                  <textarea 
                    readOnly
                    value={selectedScript.code}
                    className="flex-1 w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg p-4 text-emerald-400/90 text-sm shadow-inner outline-none resize-none leading-relaxed"
                  />
                </div>
              ) : (
                <div className="absolute inset-0 flex flex-col p-4 overflow-y-auto space-y-2 text-sm pb-16 scroll-smooth">
                  {terminalHistory.map((entry, idx) => (
                    <div key={idx} className="whitespace-pre-wrap">
                      {entry.type === 'system' && (
                        <span className="text-zinc-500 font-bold">{entry.text}</span>
                      )}
                      {entry.type === 'input' && (
                        <div><span className="text-sky-400">root@{selectedDevice.id}</span>:<span className="text-blue-400">~</span>$ {entry.text}</div>
                      )}
                      {entry.type === 'output' && (
                        <div className={cn("mt-1", entry.isError ? "text-rose-400" : "text-emerald-400")}>{entry.text}</div>
                      )}
                    </div>
                  ))}
                  <div ref={terminalEndRef} />
                </div>
              )}
              
              {!selectedScript && (
                <div className="absolute bottom-0 left-0 right-0 bg-[#0a0a0c] p-4 border-t border-white/5">
                  <form onSubmit={handleCommandSubmit} className="flex items-center gap-2 w-full">
                    <span className="text-sky-400 whitespace-nowrap">root@{selectedDevice.id}:<span className="text-blue-400">~</span>$</span>
                    <input 
                      type="text" 
                      value={terminalInput}
                      onChange={(e) => setTerminalInput(e.target.value)}
                      className="flex-1 bg-transparent border-none outline-none text-zinc-300 shadow-none"
                      autoFocus
                      autoComplete="off"
                      spellCheck="false"
                    />
                    <span className="animate-pulse w-2 h-4 bg-zinc-500/50 block" />
                  </form>
                </div>
              )}

            </div>
          </div>

        </div>
      </div>

      {/* Create Script Modal Overlay */}
      {isCreateModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="w-full max-w-2xl rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),-10px_0_30px_rgba(0,0,0,0.5)] overflow-hidden flex flex-col">
            <div className="flex items-center justify-between border-b border-white/5 p-4 bg-zinc-900/40 backdrop-blur-xl block">
              <h3 className="font-bold text-zinc-100 flex items-center gap-2">
                <Code className="h-5 w-5 text-sky-400" />
                Create Script Template
              </h3>
              <button 
                onClick={() => setIsCreateModalOpen(false)}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
                title="Close"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            
            <form onSubmit={handleCreateScript} className="flex flex-col">
              <div className="p-6 space-y-6 bg-transparent">
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Script Name</label>
                    <input 
                      type="text" 
                      required
                      placeholder="e.g. rotate_logs.sh"
                      value={newScriptForm.name}
                      onChange={(e) => setNewScriptForm(prev => ({...prev, name: e.target.value}))}
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2 text-sm font-medium text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all" 
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Interpreter Type</label>
                    <select 
                      value={newScriptForm.lang}
                      onChange={(e) => setNewScriptForm(prev => ({...prev, lang: e.target.value}))}
                      className="w-full rounded border border-white/10 bg-zinc-900 px-4 py-2 text-sm font-medium text-zinc-200 focus:border-sky-500 focus:ring-1 focus:ring-sky-500 outline-none transition-all"
                    >
                      <option value="bash">Bash Script (.sh)</option>
                      <option value="python">Python 3 (.py)</option>
                      <option value="powershell">PowerShell (.ps1)</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label className="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-2">Code Payload</label>
                  <textarea 
                    required
                    placeholder="#!/bin/bash"
                    value={newScriptForm.code}
                    onChange={(e) => setNewScriptForm(prev => ({...prev, code: e.target.value}))}
                    className="w-full rounded-lg border border-white/10 bg-[#0a0a0c] p-4 text-sm font-mono text-emerald-400/90 shadow-inner outline-none transition-all focus:border-sky-500 resize-none h-[250px]"
                  />
                </div>
              </div>
              <div className="border-t border-white/5 p-4 bg-zinc-900/40 backdrop-blur-xl flex justify-end gap-3">
                <button 
                  type="button"
                  onClick={() => setIsCreateModalOpen(false)}
                  className="rounded px-4 py-2 text-sm font-bold text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="rounded bg-sky-500 px-6 py-2 text-sm font-bold text-white hover:bg-sky-600 transition-all duration-300 ease-out active:scale-95 shadow-lg shadow-sky-500/20"
                >
                  Save Script
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
