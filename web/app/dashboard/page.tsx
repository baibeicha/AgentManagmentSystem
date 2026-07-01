'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Activity, ShieldAlert, Cpu, Network, TerminalSquare } from 'lucide-react';
import { cn } from '@/lib/utils';
import { api } from '@/lib/api';
import { useLiveMetrics } from '@/hooks/useLiveMetrics';

interface KpiData {
  label: string;
  value: string;
  change: string;
  status: 'optimal' | 'warning' | 'critical';
  icon: any;
}

interface TopDevice {
  device_id: string;
  alias: string;
  value: number;
}

export default function DashboardPage() {
  const [selectedGroup, setSelectedGroup] = useState('all');
  const [chartData, setChartData] = useState<any[]>([]);
  const [topExhausted, setTopExhausted] = useState<TopDevice[]>([]);

  // Real-time metrics hook (Optional: could be used to update KPIs)
  const { metrics, isConnected } = useLiveMetrics();

  // We could use the live metrics to update these KPIs dynamically.
  // For now, we will just stub them out as we did before, but they can be updated via the `metrics` state.
  const [kpis, setKpis] = useState<KpiData[]>([
    { label: 'Active Agents', value: '1,248', change: '+12', status: 'optimal', icon: Activity },
    { label: 'Open Incidents', value: '3', change: '-2', status: 'warning', icon: ShieldAlert },
    { label: 'Avg CPU Load', value: '42%', change: '+5%', status: 'optimal', icon: Cpu },
    { label: 'Network I/O', value: '1.2 GB/s', change: '+0.1', status: 'optimal', icon: Network },
  ]);

  useEffect(() => {
    // If we receive live metrics, update KPIs
    if (metrics) {
       // Example logic to map metrics to KPIs if the backend stream provides it
       // setKpis([...])
    }
  }, [metrics]);

  useEffect(() => {
    const fetchAggregatedMetrics = async () => {
      try {
        const queryParams = new URLSearchParams({
          metric_type: 'cpu',
          from: new Date(Date.now() - 3600000).toISOString(), // Last hour
          to: new Date().toISOString(),
          step: '5m'
        });

        if (selectedGroup !== 'all') {
          queryParams.append('group_id', selectedGroup);
        }

        const { data } = await api.get(`/api/v1/metrics/aggregated?${queryParams.toString()}`);

        // Ensure data is mapped correctly for Recharts
        if (data && data.values) {
           const formattedData = data.values.map((v: any[]) => ({
             time: new Date(parseInt(v[0]) * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
             avg_cpu: parseFloat(v[1]),
             // If we had memory, we'd map it here. We'll simulate it for visual completeness if backend only returns CPU.
             avg_memory: parseFloat(v[1]) * 0.8
           }));
           setChartData(formattedData);
        }
      } catch (err) {
        console.error('Failed to fetch aggregated metrics', err);
      }
    };

    fetchAggregatedMetrics();
  }, [selectedGroup]);

  useEffect(() => {
    const fetchTopExhausted = async () => {
      try {
        const { data } = await api.get('/api/v1/metrics/top?metric_type=cpu&limit=5');
        setTopExhausted(data || []);
      } catch (err) {
        console.error('Failed to fetch top exhausted hosts', err);
      }
    };

    fetchTopExhausted();
  }, []);

  // For visual consistency, mock incidents for the Triage Feed until Phase 3 where we fetch them.
  const MOCK_INCIDENTS = [
    { id: 'INC-1042', severity: 'critical', title: 'DB Replication Lag Exceeded', time: '10m ago' },
    { id: 'INC-1043', severity: 'warning', title: 'High Memory Usage on worker pool', time: '2h ago' },
  ];

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header */}
      <div>
        <h1 className="text-3xl font-black tracking-tighter text-zinc-100 flex items-center gap-3">
          Global Telemetry
          {isConnected ? (
            <span className="flex items-center gap-1 text-[10px] uppercase font-bold text-emerald-500 bg-emerald-500/10 px-2 py-1 rounded border border-emerald-500/20">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /> Live
            </span>
          ) : (
            <span className="flex items-center gap-1 text-[10px] uppercase font-bold text-zinc-500 bg-zinc-800/50 px-2 py-1 rounded border border-zinc-700">
              Connecting...
            </span>
          )}
        </h1>
        <p className="text-sm text-zinc-500">Real-time infrastructure performance and security events.</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {kpis.map((kpi, idx) => {
          const Icon = kpi.icon;
          return (
            <div key={idx} className="group relative overflow-hidden rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-6 transition-all hover:bg-zinc-900 shadow-sm">
              <div className="flex items-center justify-between mb-4">
                <span className="text-sm font-semibold tracking-wider uppercase text-zinc-500">{kpi.label}</span>
                <Icon className={cn("h-5 w-5", 
                  kpi.status === 'optimal' ? 'text-emerald-500 animate-pulse' : 
                  kpi.status === 'warning' ? 'text-amber-500' : 'text-rose-500 animate-pulse'
                )} />
              </div>
              <div className="flex items-end justify-between">
                <div className="text-3xl font-black text-zinc-100 font-mono tracking-tight">{kpi.value}</div>
                <div className={cn("text-xs font-bold font-mono px-2 py-1 rounded bg-zinc-950 border",
                  kpi.change.startsWith('+') ? 'text-emerald-400 border-emerald-500/20' : 'text-rose-400 border-rose-500/20'
                )}>
                  {kpi.change}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* Chart Panel (Span 2) */}
        <div className="lg:col-span-2 rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-6 shadow-sm flex flex-col">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-bold text-zinc-100">Aggregated Cluster Metrics</h2>
              <p className="text-xs text-zinc-500 font-mono">/api/v1/metrics/aggregated{selectedGroup !== 'all' ? `?group_id=${selectedGroup}` : ''}</p>
            </div>
            
            <div className="flex items-center gap-6">
              <select 
                value={selectedGroup}
                onChange={(e) => setSelectedGroup(e.target.value)}
                className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-3 py-1.5 text-[11px] font-bold text-sky-400 focus:outline-none focus:border-sky-500 cursor-pointer shadow-inner uppercase tracking-wider transition-colors outline-none"
              >
                <option value="all">All Devices</option>
                {/* Realistically, groups would be fetched, but we use static options for now to maintain UI */}
                <option value="web-tier">Web Tier</option>
                <option value="databases">Databases</option>
                <option value="worker-pool">Worker Pool</option>
              </select>
              <div className="flex items-center gap-4 text-xs font-mono font-medium hidden sm:flex">
                <div className="flex items-center gap-2"><span className="w-2 h-2 rounded-full bg-emerald-500" /> CPU</div>
                <div className="flex items-center gap-2"><span className="w-2 h-2 rounded-full bg-sky-500" /> Memory</div>
              </div>
            </div>
          </div>
          
          <div className="h-[300px] w-full mt-auto">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData} margin={{ top: 10, right: 0, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorCpu" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorMem" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#0ea5e9" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#0ea5e9" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#27272a" vertical={false} />
                <XAxis dataKey="time" stroke="#52525b" fontSize={11} tickMargin={10} axisLine={false} tickLine={false} />
                <YAxis stroke="#52525b" fontSize={11} axisLine={false} tickLine={false} tickFormatter={(v) => `${v}%`} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#18181b', borderColor: '#27272a', borderRadius: '8px', fontSize: '12px' }}
                  itemStyle={{ color: '#d4d4d8' }}
                />
                <Area type="monotone" dataKey="avg_cpu" stroke="#10b981" strokeWidth={2} fillOpacity={1} fill="url(#colorCpu)" />
                <Area type="monotone" dataKey="avg_memory" stroke="#0ea5e9" strokeWidth={2} fillOpacity={1} fill="url(#colorMem)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Side panels (Span 1) */}
        <div className="space-y-6">
          
          {/* Top Exhausted Hosts */}
          <div className="rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-6 shadow-sm">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-sm font-bold text-zinc-100 uppercase tracking-wider">Top Exhausted Hosts</h2>
              <Activity className="h-4 w-4 text-rose-500 animate-pulse" />
            </div>
            <div className="space-y-4">
              {topExhausted.map((host, idx) => (
                <div key={idx} className="group flex items-center justify-between rounded p-2 hover:bg-zinc-800/50 transition-colors">
                  <div className="flex flex-col">
                    <span className="text-sm font-bold text-zinc-200">{host.alias || 'Unknown'}</span>
                    <span className="text-[10px] font-mono text-zinc-500">{host.device_id.substring(0, 8)}...</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="flex flex-col items-end">
                      <span className="text-[10px] uppercase font-bold text-zinc-500">CPU</span>
                      <span className="text-xs font-mono text-rose-400">{host.value}%</span>
                    </div>
                  </div>
                </div>
              ))}
              {topExhausted.length === 0 && (
                <div className="text-xs text-zinc-500 font-mono text-center py-4">No data available</div>
              )}
            </div>
          </div>

          {/* Quick Incidents Feed */}
          <div className="rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-6 shadow-sm">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-sm font-bold text-zinc-100 uppercase tracking-wider">Triage Feed</h2>
              <TerminalSquare className="h-4 w-4 text-zinc-500" />
            </div>
            <div className="space-y-3">
              {MOCK_INCIDENTS.map((inc, idx) => (
                <div key={idx} className="flex gap-3 border-l-2 p-2" 
                     style={{ borderLeftColor: inc.severity === 'critical' ? '#f43f5e' : inc.severity === 'warning' ? '#f59e0b' : '#3f3f46' }}>
                  <div className="flex flex-col w-full">
                    <div className="flex items-center justify-between">
                      <span className="text-[10px] font-mono font-bold text-zinc-500">{inc.id}</span>
                      <span className="text-[10px] font-medium text-zinc-600">{inc.time}</span>
                    </div>
                    <span className="text-sm text-zinc-200 mt-0.5">{inc.title}</span>
                  </div>
                </div>
              ))}
            </div>
            <Link href="/incidents" className="w-full mt-4 flex items-center justify-center rounded border border-white/5 bg-zinc-800/50 py-2 text-xs font-bold uppercase tracking-wider text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200 transition-colors">
              View All Incidents
            </Link>
          </div>

        </div>
      </div>
    </div>
  );
}
