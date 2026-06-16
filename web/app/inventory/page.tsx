'use client';

import { useState } from 'react';
import { 
  Server, Network, Search, Filter, MoreVertical, ShieldCheck, 
  AlertTriangle, Power, X, RefreshCw, Box, UploadCloud, Trash2, Cpu
} from 'lucide-react';
import { cn } from '@/lib/utils';

// --- MOCKS ---
const INITIAL_DEVICES = [
  { id: 'dev-101', hostname: 'prod-web-01', ip_address: '10.0.1.10', os_type: 'Linux', os_version: 'Ubuntu 22.04', status: 'online', agent_version: '2.0.4', last_seen: '2024-03-10T12:00:00Z', group_id: 'web-tier', tags: ['prod', 'web'] },
  { id: 'dev-102', hostname: 'prod-db-master', ip_address: '10.0.2.5', os_type: 'Linux', os_version: 'RHEL 9', status: 'online', agent_version: '2.0.4', last_seen: '2024-03-10T11:59:00Z', group_id: 'databases', tags: ['prod', 'db'] },
  { id: 'dev-103', hostname: 'corp-file-srv', ip_address: '192.168.1.50', os_type: 'Windows', os_version: 'Server 2022', status: 'offline', agent_version: '2.0.3', last_seen: '2024-03-09T08:15:00Z', group_id: 'corp-it', tags: ['windows', 'internal'] },
  { id: 'dev-104', hostname: 'staging-worker-1', ip_address: '10.0.3.15', os_type: 'Linux', os_version: 'Debian 12', status: 'degraded', agent_version: '2.0.4', last_seen: '2024-03-10T11:58:30Z', group_id: 'staging', tags: ['staging', 'worker'] },
];

const MOCK_DISCOVERY = [
  { mac_address: '00:1B:44:11:3A:B7', initial_ip: '192.168.1.105', hostname_hint: 'DESKTOP-J8K9P', os_guess: 'Windows 10/11', discovery_time: '2024-03-10T10:00:00Z' },
  { mac_address: '00:15:5D:8A:1C:12', initial_ip: '192.168.1.210', hostname_hint: 'unknown-device', os_guess: 'Linux 3.x/4.x', discovery_time: '2024-03-10T09:45:00Z' },
];

const MOCK_PROCESSES = [
  { pid: 1450, name: 'nginx', cpu_percent: 2.5, memory_mb: 45.2, user: 'www-data' },
  { pid: 3306, name: 'mysqld', cpu_percent: 15.0, memory_mb: 1024.5, user: 'mysql' },
  { pid: 1, name: 'systemd', cpu_percent: 0.1, memory_mb: 12.0, user: 'root' },
  { pid: 994, name: 'aegis-agent', cpu_percent: 1.2, memory_mb: 28.4, user: 'aegis' },
];

export default function InventoryPage() {
  const [activeTab, setActiveTab] = useState<'fleet' | 'discovery'>('fleet');
  const [devices, setDevices] = useState(INITIAL_DEVICES);
  const [selectedDevice, setSelectedDevice] = useState<typeof INITIAL_DEVICES[0] | null>(null);
  
  // Search & Filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  
  // Panel Action States
  const [isUpdateMode, setIsUpdateMode] = useState(false);
  const [selectedUpdateVersion, setSelectedUpdateVersion] = useState('v2.0.5-patch1');

  // Extract all unique tags
  const allTags = Array.from(new Set(INITIAL_DEVICES.flatMap(dev => dev.tags)));

  const toggleTag = (tag: string) => {
    setSelectedTags(prev => prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]);
  };

  const filteredDevices = devices.filter(d => {
    const matchesSearch = d.hostname.toLowerCase().includes(searchQuery.toLowerCase()) || d.ip_address.includes(searchQuery);
    const matchesStatus = statusFilter === 'all' || d.status === statusFilter;
    const matchesTags = selectedTags.length === 0 || selectedTags.every(tag => d.tags.includes(tag));
    return matchesSearch && matchesStatus && matchesTags;
  });

  const handleDeregister = () => {
    if (selectedDevice) {
      setDevices(prev => prev.filter(d => d.id !== selectedDevice.id));
      setSelectedDevice(null);
    }
  };

  const handleUpdateAgent = () => {
    if (selectedDevice) {
      // Mock update successful
      setDevices(prev => prev.map(d => 
        d.id === selectedDevice.id 
          ? { ...d, agent_version: selectedUpdateVersion.replace('v', '') }
          : d
      ));
      setSelectedDevice(prev => prev ? { ...prev, agent_version: selectedUpdateVersion.replace('v', '') } : null);
      setIsUpdateMode(false);
    }
  };

  return (
    <div className="relative flex h-[calc(100vh-8rem)] flex-col space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header & Tabs */}
      <div className="flex flex-col border-b border-white/5 pb-4">
        <h1 className="text-3xl font-black tracking-tighter text-zinc-100 mb-4">Inventory & Discovery</h1>
        
        <div className="flex gap-2">
          <button
            onClick={() => setActiveTab('fleet')}
            className={cn(
              "flex items-center gap-2 rounded-md px-4 py-2 text-sm font-bold uppercase tracking-wide transition-all",
              activeTab === 'fleet' ? "bg-sky-500/10 text-sky-400 border border-sky-500/20 shadow-[inset_0_0_10px_rgba(14,165,233,0.1)]" : "text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900"
            )}
          >
            <Server className="h-4 w-4" />
            Managed Fleet
          </button>
          <button
            onClick={() => setActiveTab('discovery')}
            className={cn(
              "flex items-center gap-2 rounded-md px-4 py-2 text-sm font-bold uppercase tracking-wide transition-all",
              activeTab === 'discovery' ? "bg-sky-500/10 text-sky-400 border border-sky-500/20 shadow-[inset_0_0_10px_rgba(14,165,233,0.1)]" : "text-zinc-400 hover:text-zinc-200 hover:bg-zinc-900"
            )}
          >
            <Network className="h-4 w-4" />
            LAN Discovery
          </button>
        </div>
      </div>
      
      {/* Control Toolbar Above Table */}
      {activeTab === 'fleet' && (
        <div className="flex items-center justify-between rounded-lg bg-zinc-900/40 border border-white/5 p-3 shrink-0">
          <div className="flex bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded overflow-hidden shadow-inner w-72">
             <div className="flex items-center pl-3 text-zinc-500 pointer-events-none">
               <Search className="h-4 w-4" />
             </div>
             <input 
               type="text"
               value={searchQuery}
               onChange={e => setSearchQuery(e.target.value)}
               placeholder="Search hostname or IP..."
               className="bg-transparent border-none text-sm px-3 py-1.5 w-full text-zinc-200 focus:outline-none placeholder:text-zinc-600 font-mono"
             />
          </div>
          
          <div className="flex items-center gap-4 text-xs">
            <div className="flex items-center gap-2 text-zinc-400 font-bold uppercase tracking-wider">
               <Filter className="h-3 w-3" /> Status:
            </div>
            <select
               value={statusFilter}
               onChange={e => setStatusFilter(e.target.value)}
               className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] text-zinc-200 font-bold px-3 py-1.5 rounded cursor-pointer outline-none focus:border-sky-500 text-sm shadow-sm transition-colors uppercase tracking-widest"
            >
               <option value="all">All Status</option>
               <option value="online">Online</option>
               <option value="offline">Offline</option>
               <option value="degraded">Degraded</option>
            </select>
          </div>
        </div>
      )}

      {/* Main Content Area */}
      <div className="flex-1 overflow-hidden relative">
        
        {/* Fleet Tab */}
        {activeTab === 'fleet' && (
          <div className="h-full flex flex-col">
            <div className="h-full overflow-y-auto rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)]">
               <table className="w-full text-left text-sm text-zinc-400">
                 <thead className="sticky top-0 bg-zinc-950/90 text-xs font-semibold uppercase tracking-wider text-zinc-500 border-b border-white/5 backdrop-blur-md z-10">
                   <tr>
                     <th className="px-6 py-4">Status</th>
                     <th className="px-6 py-4">Hostname / IP</th>
                     <th className="px-6 py-4">OS / Agent</th>
                     <th className="px-6 py-4">Group</th>
                     <th className="px-6 py-4 text-right">Actions</th>
                   </tr>
                 </thead>
                 <tbody className="divide-y divide-white/5">
                   {filteredDevices.length === 0 ? (
                     <tr>
                       <td colSpan={5} className="py-12 text-center text-zinc-500">
                         No managed devices found matching filters.
                       </td>
                     </tr>
                   ) : (
                     filteredDevices.map((dev) => (
                       <tr 
                        key={dev.id} 
                        onClick={() => setSelectedDevice(dev)}
                        className={cn(
                          "transition-colors transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95 cursor-pointer",
                          selectedDevice?.id === dev.id && "bg-white/[0.04]"
                        )}
                       >
                         <td className="px-6 py-4">
                           <div className="flex items-center gap-2">
                             <div className={cn(
                               "flex h-6 w-6 items-center justify-center rounded-full border shadow-sm",
                               dev.status === 'online' ? "bg-emerald-500/10 border-emerald-500/20 text-emerald-500 animate-pulse shadow-emerald-500/10" :
                               dev.status === 'degraded' ? "bg-amber-500/10 border-amber-500/20 text-amber-500 shadow-amber-500/10" :
                               "bg-rose-500/10 border-rose-500/20 text-rose-500 animate-pulse shadow-rose-500/10"
                             )}>
                               {dev.status === 'online' ? <ShieldCheck className="h-3 w-3" /> : 
                                dev.status === 'degraded' ? <AlertTriangle className="h-3 w-3" /> : 
                                <Power className="h-3 w-3" />}
                             </div>
                             <span className={cn("font-mono text-[10px] uppercase font-bold", 
                               dev.status === 'online' ? "text-emerald-400" :
                               dev.status === 'degraded' ? "text-amber-400" : "text-rose-400"
                             )}>{dev.status}</span>
                           </div>
                         </td>
                         <td className="px-6 py-4">
                           <div className="font-bold text-zinc-200">{dev.hostname}</div>
                           <div className="font-mono text-xs text-zinc-500">{dev.ip_address}</div>
                         </td>
                         <td className="px-6 py-4">
                           <div className="text-zinc-300">{dev.os_type} <span className="text-zinc-500">({dev.os_version})</span></div>
                           <div className="font-mono text-xs text-sky-400">v{dev.agent_version}</div>
                         </td>
                         <td className="px-6 py-4">
                           <span className="rounded bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-2 py-1 font-mono text-[10px] text-zinc-400">
                             {dev.group_id}
                           </span>
                         </td>
                         <td className="px-6 py-4 text-right">
                           <button className="p-2 text-zinc-500 hover:text-zinc-300 transition-all duration-300 ease-out active:scale-95 rounded-full hover:bg-zinc-800">
                             <MoreVertical className="h-4 w-4" />
                           </button>
                         </td>
                       </tr>
                     ))
                   )}
                 </tbody>
               </table>
            </div>
          </div>
        )}

        {/* Discovery Tab */}
        {activeTab === 'discovery' && (
          <div className="h-full overflow-y-auto rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-6 shadow-sm">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <h2 className="text-lg font-bold text-zinc-100">Unmanaged Devices Detected</h2>
                <p className="text-sm text-zinc-500">Devices found via active LAN sweep (ARP/ICMP). Deploy agents to manage.</p>
              </div>
              <button className="flex items-center gap-2 rounded bg-sky-500 px-4 py-2 text-sm font-bold text-white hover:bg-sky-600 transition-all duration-300 ease-out active:scale-95 shadow-[0_0_15px_rgba(14,165,233,0.3)]">
                <RefreshCw className="h-4 w-4" /> Rescan Network
              </button>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {MOCK_DISCOVERY.map((cand, idx) => (
                <div key={idx} className="flex flex-col rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-$1 relative overflow-hidden group shadow-sm">
                  <div className="absolute top-0 right-0 p-4 opacity-5 pointer-events-none group-hover:opacity-10 transition-opacity whitespace-nowrap text-right">
                    <Box className="w-24 h-24 text-sky-500" />
                  </div>
                  <div className="flex items-start justify-between mb-4 relative z-10">
                    <div>
                      <h3 className="font-bold text-zinc-200">{cand.hostname_hint}</h3>
                      <p className="font-mono text-xs text-zinc-500">{cand.initial_ip}</p>
                    </div>
                    <span className="rounded bg-zinc-800/50 border border-white/5 px-2 py-1 font-mono text-[10px] text-zinc-400">
                      {cand.mac_address}
                    </span>
                  </div>
                  <div className="mt-auto flex items-center justify-between border-t border-white/5 pt-4 relative z-10">
                    <div className="flex flex-col">
                      <span className="text-[10px] uppercase font-bold text-zinc-500">OS Guess</span>
                      <span className="text-sm font-medium text-zinc-300">{cand.os_guess}</span>
                    </div>
                    <button className="text-[11px] font-bold uppercase tracking-wider text-sky-400 hover:text-sky-300 transition-all duration-300 ease-out active:scale-95 border border-sky-500/20 bg-sky-500/10 px-3 py-1.5 rounded hover:bg-sky-500/20">
                      Push Agent Request
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

      </div>

      {/* Sliding Side Panel for Host Details */}
      <div 
        className={cn(
          "absolute right-0 top-0 bottom-0 w-[450px] border-l border-white/10 bg-zinc-950 shadow-2xl transition-transform duration-300 z-40 flex flex-col",
          selectedDevice ? "translate-x-0" : "translate-x-full"
        )}
      >
        {selectedDevice && (
          <>
            <div className="flex items-center justify-between border-b border-white/5 p-6 bg-zinc-900/30">
              <div>
                <h2 className="text-xl font-bold text-zinc-100 flex items-center gap-2">
                  <Server className="h-5 w-5 text-sky-400" />
                  {selectedDevice.hostname}
                </h2>
                <p className="font-mono text-sm text-zinc-500">{selectedDevice.ip_address}</p>
              </div>
              <button 
                onClick={() => setSelectedDevice(null)}
                className="rounded-full p-2 text-zinc-500 hover:bg-zinc-800 hover:text-zinc-200 transition-colors"
              >
                <X className="h-5 w-5" />
              </button>
            </div>

             <div className="flex-1 overflow-y-auto p-6 space-y-8">
              {/* Properties */}
              <div>
                <h3 className="text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-3">System Identity</h3>
                <div className="grid grid-cols-2 gap-4 rounded-lg border border-white/5 bg-zinc-900/30 p-4 shadow-inner">
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">OS Type</div>
                    <div className="text-sm font-medium text-zinc-200">{selectedDevice.os_type}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Version</div>
                    <div className="text-sm font-medium text-zinc-200">{selectedDevice.os_version}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Agent Version</div>
                    <div className="font-mono text-sm font-bold text-sky-400">v{selectedDevice.agent_version}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Group ID</div>
                    <div className="font-mono text-[11px] font-bold text-zinc-300 bg-zinc-800 inline-block px-1.5 py-0.5 rounded mt-1">{selectedDevice.group_id}</div>
                  </div>
                </div>
              </div>

               {/* Agent Management */}
               <div>
                  <h3 className="text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-3 flex items-center justify-between">
                    Agent Management
                  </h3>
                  <div className="rounded-lg border border-white/5 bg-zinc-900/30 p-4 space-y-4">
                    {isUpdateMode ? (
                      <div className="flex items-end gap-2 animate-in fade-in duration-200">
                        <div className="flex-1">
                           <label className="text-[10px] uppercase font-bold text-zinc-500 block mb-1">Target Version</label>
                           <select 
                            value={selectedUpdateVersion}
                            onChange={(e) => setSelectedUpdateVersion(e.target.value)}
                            className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-2 py-1.5 text-sm font-mono text-zinc-200 focus:outline-none focus:border-sky-500"
                           >
                              <option value="v2.0.4">v2.0.4 (Current)</option>
                              <option value="v2.0.5-patch1">v2.0.5-patch1 (Latest)</option>
                              <option value="v2.1.0-beta">v2.1.0-beta</option>
                           </select>
                        </div>
                        <button 
                          onClick={handleUpdateAgent}
                          className="bg-sky-500 hover:bg-sky-600 text-white px-4 py-1.5 rounded text-sm font-bold transition-all duration-300 ease-out active:scale-95"
                        >
                          Execute
                        </button>
                        <button 
                          onClick={() => setIsUpdateMode(false)}
                          className="bg-zinc-800 hover:bg-zinc-700 text-zinc-300 px-3 py-1.5 rounded text-sm font-bold transition-colors"
                        >
                          Cancel
                        </button>
                      </div>
                    ) : (
                      <div className="flex gap-2">
                        <button 
                          onClick={() => setIsUpdateMode(true)}
                          className="flex items-center gap-2 rounded bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-4 py-2 text-xs font-bold uppercase tracking-wider text-zinc-300 hover:bg-zinc-700 transition-colors"
                        >
                          <UploadCloud className="h-4 w-4" /> Trigger Remote Update
                        </button>
                      </div>
                    )}
                  </div>
               </div>

              {/* Processes */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-[11px] font-bold uppercase tracking-wider text-zinc-500 flex items-center gap-2">
                    <Cpu className="h-4 w-4" /> Active OS Processes
                  </h3>
                  <button className="text-[10px] font-medium text-sky-400 hover:underline">/api/v1/metrics/{selectedDevice.id}/processes</button>
                </div>
                <div className="rounded-lg border border-white/5 bg-zinc-900/30 overflow-hidden shadow-inner">
                  <table className="w-full text-left font-mono text-[11px]">
                    <thead className="bg-zinc-900/80 text-zinc-500 border-b border-white/5">
                      <tr>
                        <th className="p-2.5 font-normal">PID</th>
                        <th className="p-2.5 font-normal">NAME</th>
                        <th className="p-2.5 font-normal text-right">CPU</th>
                        <th className="p-2.5 font-normal text-right">RAM (MB)</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5 text-zinc-300">
                      {MOCK_PROCESSES.map(proc => (
                        <tr key={proc.pid} className="transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95">
                          <td className="p-2.5 text-zinc-500">{proc.pid}</td>
                          <td className="p-2.5 text-zinc-200 font-bold">{proc.name}</td>
                          <td className="p-2.5 text-right font-bold text-emerald-400">{proc.cpu_percent}%</td>
                          <td className="p-2.5 text-right font-bold text-amber-400">{proc.memory_mb}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

            </div>

            <div className="border-t border-white/5 p-4 bg-zinc-900/80 flex gap-2">
              <button 
                onClick={handleDeregister}
                className="flex flex-1 items-center justify-center gap-2 rounded bg-rose-500/10 border border-rose-500/20 py-2.5 text-xs font-bold uppercase tracking-wider text-rose-400 hover:bg-rose-500/20 transition-all duration-300 ease-out active:scale-95 shadow-sm shadow-rose-500/5 group"
              >
                <Trash2 className="h-4 w-4 group-hover:drop-shadow-[0_0_5px_rgba(244,63,94,0.5)]" /> De-register Host
              </button>
            </div>
          </>
        )}
      </div>

    </div>
  );
}
