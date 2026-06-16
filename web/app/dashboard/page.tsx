'use client';

import { 
  Activity, ShieldAlert, Cpu, Network, CheckCircle2, AlertTriangle, TerminalSquare
} from 'lucide-react';
import { 
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  BarChart, Bar
} from 'recharts';
import { cn } from '@/lib/utils';
import { useState } from 'react';
import Link from 'next/link';

// --- MOCKS ---
const MOCK_KPIS = [
  { label: 'Active Agents', value: '1,452', change: '+12', status: 'optimal', icon: Network },
  { label: 'Avg CPU Load', value: '45.2%', change: '+2.1%', status: 'warning', icon: Cpu },
  { label: 'Unresolved Incidents', value: '3', change: '-2', status: 'critical', icon: ShieldAlert },
  { label: 'Network Throughput', value: '1.2 TB/s', change: '+0.1 TB/s', status: 'optimal', icon: Activity },
];

const MOCK_HISTORICAL_METRICS_BASE = Array.from({ length: 24 }).map((_, i) => ({
  time: `${i}:00`,
  avg_cpu: Math.floor(Math.random() * 40 + 20),
  avg_memory: Math.floor(Math.random() * 30 + 40),
  network_io: Math.floor(Math.random() * 100 + 50),
}));

const MOCK_TOP_EXHAUSTED = [
  { host_id: 'db-master-01', hostname: 'eu-west-db1', cpu: 98.5, mem: 92.1 },
  { host_id: 'cache-redis-04', hostname: 'us-east-cache4', cpu: 94.2, mem: 88.0 },
  { host_id: 'worker-node-12', hostname: 'ap-south-worker12', cpu: 89.1, mem: 76.5 },
  { host_id: 'lb-ingress-02', hostname: 'eu-central-lb2', cpu: 85.0, mem: 60.2 },
];

const MOCK_INCIDENTS = [
  { id: 'INC-1042', severity: 'critical', title: 'DB Replication Lag Exceeded', time: '10m ago' },
  { id: 'INC-1043', severity: 'warning', title: 'High Memory Usage on worker pool', time: '2h ago' },
  { id: 'INC-1044', severity: 'info', title: 'Automated Snapshot Completed', time: '4h ago' },
];

export default function DashboardPage() {
  const [selectedGroup, setSelectedGroup] = useState('all');

  const multiplier = selectedGroup === 'all' ? 1 : selectedGroup === 'web-tier' ? 0.6 : selectedGroup === 'databases' ? 1.4 : 0.8;
  const currentChartData = MOCK_HISTORICAL_METRICS_BASE.map(d => ({
    time: d.time,
    avg_cpu: Math.min(d.avg_cpu * multiplier, 100),
    avg_memory: Math.min(d.avg_memory * multiplier, 100),
  }));

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header */}
      <div>
        <h1 className="text-3xl font-black tracking-tighter text-zinc-100">Global Telemetry</h1>
        <p className="text-sm text-zinc-500">Real-time infrastructure performance and security events.</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {MOCK_KPIS.map((kpi, idx) => {
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
              <AreaChart data={currentChartData} margin={{ top: 10, right: 0, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorCpu" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorMem" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#6366f1" stopOpacity={0}/>
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
                <Area type="monotone" dataKey="avg_memory" stroke="#6366f1" strokeWidth={2} fillOpacity={1} fill="url(#colorMem)" />
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
              {MOCK_TOP_EXHAUSTED.map((host, idx) => (
                <div key={idx} className="group flex items-center justify-between rounded p-2 hover:bg-zinc-800/50 transition-colors">
                  <div className="flex flex-col">
                    <span className="text-sm font-bold text-zinc-200">{host.hostname}</span>
                    <span className="text-[10px] font-mono text-zinc-500">{host.host_id}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="flex flex-col items-end">
                      <span className="text-[10px] uppercase font-bold text-zinc-500">CPU</span>
                      <span className="text-xs font-mono text-rose-400">{host.cpu}%</span>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="text-[10px] uppercase font-bold text-zinc-500">MEM</span>
                      <span className="text-xs font-mono text-amber-400">{host.mem}%</span>
                    </div>
                  </div>
                </div>
              ))}
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
