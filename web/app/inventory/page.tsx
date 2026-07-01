'use client';

import { useState, useEffect } from 'react';
import { Search, Server, Plus, Filter, Trash2, Cpu, Box, X, Activity, HardDrive, Network, UploadCloud } from 'lucide-react';
import { cn } from '@/lib/utils';
import { api } from '@/lib/api';

interface Device {
  device_id: string;
  alias: string;
  os: string;
  arch: string;
  ip_address: string;
  status: 'online' | 'offline';
  agent_version: string;
  group_id: string;
  tags: string[];
  last_heartbeat: string;
}

interface Process {
  pid: number;
  name: string;
  cpu_usage: number;
  ram_usage_bytes: number;
}

const MOCK_CANDIDATES = [
  { id: 'cand-1', ip_address: '10.0.5.12', mac_address: '00:1A:2B:3C:4D:5E' },
  { id: 'cand-2', ip_address: '10.0.5.15', mac_address: '00:1A:2B:3C:4D:5F' },
];

export default function InventoryPage() {
  const [activeTab, setActiveTab] = useState<'managed' | 'discovery'>('managed');
  const [selectedDevice, setSelectedDevice] = useState<Device | null>(null);
  const [isUpdateMode, setIsUpdateMode] = useState(false);
  const [selectedUpdateVersion, setSelectedUpdateVersion] = useState('v1.3.0');

  const [devices, setDevices] = useState<Device[]>([]);
  const [processes, setProcesses] = useState<Process[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchDevices = async () => {
      try {
        setIsLoading(true);
        const { data } = await api.get('/api/v1/devices');
        setDevices(data || []);
      } catch (err) {
        console.error('Failed to fetch devices', err);
      } finally {
        setIsLoading(false);
      }
    };

    if (activeTab === 'managed') {
      fetchDevices();
    }
  }, [activeTab]);

  useEffect(() => {
    const fetchProcesses = async () => {
      if (!selectedDevice) return;
      try {
        const { data } = await api.get(`/api/v1/metrics/${selectedDevice.device_id}/processes`);
        setProcesses(data || []);
      } catch (err) {
        console.error('Failed to fetch processes', err);
      }
    };

    if (selectedDevice) {
      fetchProcesses();
    }
  }, [selectedDevice]);

  const handleDeregister = async () => {
    if (!selectedDevice) return;
    try {
      await api.delete(`/api/v1/devices/${selectedDevice.device_id}`);
      setDevices(devices.filter(d => d.device_id !== selectedDevice.device_id));
      setSelectedDevice(null);
    } catch (err) {
      console.error('Failed to deregister device', err);
    }
  };

  const handleUpdateAgent = async () => {
    if (!selectedDevice) return;
    try {
      await api.post('/api/v1/devices/update', {
        device_id: selectedDevice.device_id,
        target_version: selectedUpdateVersion
      });
      setIsUpdateMode(false);
      // Optional: show a success toast here
    } catch (err) {
      console.error('Failed to update agent', err);
    }
  };

  return (
    <div className="relative h-[calc(100vh-6rem)] overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header Area */}
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-6 gap-4">
        <div>
          <h1 className="text-3xl font-black tracking-tighter text-zinc-100 flex items-center gap-3">
            Fleet Inventory
          </h1>
          <p className="text-sm text-zinc-500">Manage organizational hosts and deploy agents.</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg p-1 shadow-inner">
            <button
              onClick={() => setActiveTab('managed')}
              className={cn("px-4 py-1.5 rounded text-sm font-bold transition-all", activeTab === 'managed' ? "bg-zinc-800 text-zinc-100 shadow" : "text-zinc-500 hover:text-zinc-300")}
            >
              Managed ({devices.length})
            </button>
            <button
              onClick={() => setActiveTab('discovery')}
              className={cn("px-4 py-1.5 rounded text-sm font-bold transition-all flex items-center gap-2", activeTab === 'discovery' ? "bg-zinc-800 text-sky-400 shadow" : "text-zinc-500 hover:text-zinc-300")}
            >
              Discovery
              <span className="flex h-4 w-4 items-center justify-center rounded-full bg-sky-500/20 text-[9px] text-sky-400">2</span>
            </button>
          </div>
          <button className="flex items-center gap-2 rounded bg-sky-500 py-2 px-4 text-sm font-bold text-white shadow-[0_0_15px_rgba(14,165,233,0.4)] transition-all hover:bg-sky-600">
            <Plus className="h-4 w-4" /> Provision Key
          </button>
        </div>
      </div>

      <div className="h-full pb-20">
        
        {/* Managed Hosts Tab */}
        {activeTab === 'managed' && (
          <div className="h-full flex flex-col space-y-4 animate-in fade-in duration-500">
            {/* Search and Filter */}
            <div className="flex gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-500" />
                <input
                  type="text"
                  placeholder="Filter by hostname, IP, or tag..."
                  className="w-full rounded-lg border border-white/[0.08] bg-zinc-900/40 py-2 pl-10 pr-4 text-sm font-medium text-zinc-100 shadow-[inset_0_1px_0_rgba(255,255,255,0.05)] outline-none transition-all focus:border-sky-500 focus:ring-1 focus:ring-sky-500 backdrop-blur-md"
                />
              </div>
              <button className="flex items-center gap-2 rounded-lg bg-zinc-900/40 border border-white/[0.08] px-4 py-2 text-sm font-bold text-zinc-400 hover:text-zinc-200 transition-colors backdrop-blur-md shadow-sm">
                <Filter className="h-4 w-4" /> Filters
              </button>
            </div>

            {/* Table */}
            <div className="flex-1 rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] overflow-hidden flex flex-col shadow-sm">
              <div className="overflow-x-auto flex-1">
                <table className="w-full text-left text-sm whitespace-nowrap">
                  <thead className="bg-zinc-900/80 text-xs uppercase tracking-wider text-zinc-500 border-b border-white/5 sticky top-0 z-10 backdrop-blur-md">
                    <tr>
                      <th className="p-4 font-bold">Status</th>
                      <th className="p-4 font-bold">Hostname</th>
                      <th className="p-4 font-bold">IP Address</th>
                      <th className="p-4 font-bold">OS</th>
                      <th className="p-4 font-bold">Tags</th>
                      <th className="p-4 font-bold">Agent Version</th>
                      <th className="p-4 font-bold text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-white/5 text-zinc-300">
                    {isLoading ? (
                      <tr>
                        <td colSpan={7} className="p-4 text-center text-zinc-500 text-xs font-mono">Loading...</td>
                      </tr>
                    ) : devices.length === 0 ? (
                      <tr>
                        <td colSpan={7} className="p-4 text-center text-zinc-500 text-xs font-mono">No devices found</td>
                      </tr>
                    ) : (
                      devices.map(dev => (
                        <tr
                          key={dev.device_id}
                          className="transition-colors hover:bg-white/[0.02] group cursor-pointer"
                          onClick={() => setSelectedDevice(dev)}
                        >
                          <td className="p-4">
                            <div className="flex items-center gap-2">
                              <span className={cn("relative flex h-2.5 w-2.5")}>
                                {dev.status === 'online' && <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>}
                                <span className={cn("relative inline-flex rounded-full h-2.5 w-2.5", dev.status === 'online' ? 'bg-emerald-500' : 'bg-rose-500')}></span>
                              </span>
                              <span className="text-[10px] font-bold uppercase tracking-wider text-zinc-500">{dev.status}</span>
                            </div>
                          </td>
                          <td className="p-4 font-bold text-zinc-200">
                            <div className="flex items-center gap-2">
                              <Server className="h-4 w-4 text-sky-400 opacity-50 group-hover:opacity-100 transition-opacity" />
                              {dev.alias || dev.device_id.substring(0,8)}
                            </div>
                          </td>
                          <td className="p-4 font-mono text-xs text-zinc-400">{dev.ip_address}</td>
                          <td className="p-4">
                            <div className="flex items-center gap-1.5 text-xs font-bold text-zinc-400 uppercase">
                              {dev.os} {dev.arch}
                            </div>
                          </td>
                          <td className="p-4">
                            <div className="flex gap-1.5">
                              {(dev.tags || []).map(tag => (
                                <span key={tag} className="rounded bg-zinc-800 border border-white/5 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-zinc-400">
                                  {tag}
                                </span>
                              ))}
                            </div>
                          </td>
                          <td className="p-4 font-mono text-xs font-bold text-sky-400">
                            v{dev.agent_version}
                          </td>
                          <td className="p-4 text-right">
                            <button
                              onClick={(e) => { e.stopPropagation(); setSelectedDevice(dev); }}
                              className="text-[10px] font-bold uppercase tracking-wider text-zinc-500 hover:text-sky-400 transition-colors"
                            >
                              Inspect
                            </button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* Discovery Tab */}
        {activeTab === 'discovery' && (
          <div className="animate-in fade-in duration-500 space-y-6">
            <div className="rounded-xl bg-sky-500/10 border border-sky-500/20 p-4 flex items-start gap-4">
               <div className="p-2 bg-sky-500/20 rounded-full">
                  <Network className="h-5 w-5 text-sky-400" />
               </div>
               <div>
                 <h3 className="text-sm font-bold text-sky-400">Network Discovery Active</h3>
                 <p className="text-xs text-sky-400/70 mt-1">Managed nodes are passively sniffing ARP traffic to identify unmanaged devices on local subnets. Review candidates below to push agents.</p>
               </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {MOCK_CANDIDATES.map((cand, idx) => (
                <div key={idx} className="flex flex-col rounded-lg bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-5 relative overflow-hidden group shadow-sm">
                  <div className="absolute top-0 right-0 p-4 opacity-5 pointer-events-none group-hover:opacity-10 transition-opacity whitespace-nowrap text-right">
                    <Box className="w-24 h-24 text-sky-500" />
                  </div>
                  <div className="flex items-start justify-between mb-4 relative z-10">
                    <div>
                      <h3 className="font-bold text-zinc-200">Candidate #{idx+1}</h3>
                      <p className="font-mono text-xs text-zinc-500">{cand.ip_address}</p>
                    </div>
                    <span className="rounded bg-zinc-800/50 border border-white/5 px-2 py-1 font-mono text-[10px] text-zinc-400">
                      {cand.mac_address}
                    </span>
                  </div>
                  <div className="mt-auto flex items-center justify-between border-t border-white/5 pt-4 relative z-10">
                    <button className="text-[11px] font-bold uppercase tracking-wider text-sky-400 hover:text-sky-300 transition-all duration-300 ease-out active:scale-95 border border-sky-500/20 bg-sky-500/10 px-3 py-1.5 rounded hover:bg-sky-500/20">
                      Approve Candidate
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
                  {selectedDevice.alias || selectedDevice.device_id.substring(0,8)}
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
                    <div className="text-sm font-medium text-zinc-200 uppercase">{selectedDevice.os}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Arch</div>
                    <div className="text-sm font-medium text-zinc-200 uppercase">{selectedDevice.arch}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Agent Version</div>
                    <div className="font-mono text-sm font-bold text-sky-400">v{selectedDevice.agent_version}</div>
                  </div>
                  <div>
                    <div className="text-[10px] uppercase text-zinc-500">Group ID</div>
                    <div className="font-mono text-[11px] font-bold text-zinc-300 bg-zinc-800 inline-block px-1.5 py-0.5 rounded mt-1">{selectedDevice.group_id?.substring(0,8) || 'N/A'}</div>
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
                              <option value="v1.2.0">v1.2.0</option>
                              <option value="v1.3.0">v1.3.0</option>
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
                  <button className="text-[10px] font-medium text-sky-400 hover:underline">/api/v1/metrics/.../processes</button>
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
                      {processes.length === 0 ? (
                        <tr><td colSpan={4} className="p-2.5 text-center text-zinc-500">No data</td></tr>
                      ) : (
                        processes.map(proc => (
                          <tr key={proc.pid} className="transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95">
                            <td className="p-2.5 text-zinc-500">{proc.pid}</td>
                            <td className="p-2.5 text-zinc-200 font-bold">{proc.name}</td>
                            <td className="p-2.5 text-right font-bold text-emerald-400">{proc.cpu_usage}%</td>
                            <td className="p-2.5 text-right font-bold text-amber-400">{Math.round(proc.ram_usage_bytes / 1024 / 1024)}</td>
                          </tr>
                        ))
                      )}
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
