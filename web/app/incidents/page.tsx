'use client';

import { useState, useEffect } from 'react';
import { ShieldAlert, AlertTriangle, Activity, Info, Check, CheckCircle2, Clock, TerminalSquare, GitMerge, X, Filter } from 'lucide-react';
import { cn } from '@/lib/utils';
import { api } from '@/lib/api';

interface Incident {
  incident_id: string;
  device_id: string;
  type: string;
  severity: 'info' | 'warning' | 'critical';
  status: 'OPEN' | 'ACKNOWLEDGED' | 'RESOLVED';
  created_at: string;
  description?: string;
}

export default function IncidentsPage() {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [filterStatus, setFilterStatus] = useState<'ALL' | 'OPEN' | 'ACKNOWLEDGED' | 'RESOLVED'>('ALL');
  const [runbookModalIncident, setRunbookModalIncident] = useState<Incident | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchIncidents = async () => {
      try {
        setIsLoading(true);
        const { data } = await api.get('/api/v1/incidents/active');
        // Add a mock description since the backend API doesn't provide it by default in the schema,
        // but the UI relies on it for context.
        const enhancedData = (data || []).map((inc: Incident) => ({
          ...inc,
          description: `Anomaly detected matching rule signature for ${inc.type}. Automatic triage initiated.`
        }));
        setIncidents(enhancedData);
      } catch (err) {
        console.error('Failed to fetch incidents', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchIncidents();
  }, []);

  const handleUpdateStatus = async (id: string, newStatus: 'ACKNOWLEDGED' | 'RESOLVED') => {
    try {
      const endpoint = newStatus === 'ACKNOWLEDGED' ? 'acknowledge' : 'resolve';
      await api.put(`/api/v1/incidents/${id}/${endpoint}`);

      // Mutate local state for instant UI feedback
      setIncidents(prev => prev.map(inc =>
        inc.incident_id === id ? { ...inc, status: newStatus } : inc
      ));
    } catch (err) {
      console.error(`Failed to update incident status to ${newStatus}`, err);
    }
  };

  const filteredIncidents = incidents.filter(inc => {
    if (filterStatus === 'ALL') return true;
    return inc.status === filterStatus;
  });

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700">
      
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-black tracking-tighter text-zinc-100 flex items-center gap-3">
            Active Incidents
            {incidents.some(i => i.status === 'OPEN' && i.severity === 'critical') && (
               <span className="flex items-center gap-1 text-[10px] uppercase font-bold text-rose-500 bg-rose-500/10 px-2 py-1 rounded border border-rose-500/20 animate-pulse">
                 Critical Outages
               </span>
            )}
          </h1>
          <p className="text-sm text-zinc-500">Respond to active infrastructural anomalies and security breaches.</p>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex items-center bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded p-1">
             {['ALL', 'OPEN', 'ACKNOWLEDGED', 'RESOLVED'].map(status => (
               <button
                 key={status}
                 onClick={() => setFilterStatus(status as any)}
                 className={cn(
                   "px-3 py-1 text-xs font-bold uppercase tracking-wider rounded transition-colors",
                   filterStatus === status ? "bg-zinc-800 text-zinc-100" : "text-zinc-500 hover:text-zinc-300"
                 )}
               >
                 {status}
               </button>
             ))}
          </div>
          <button className="flex items-center gap-2 rounded-lg bg-zinc-900/40 border border-white/[0.08] px-4 py-2 text-sm font-bold text-zinc-400 hover:text-zinc-200 transition-colors backdrop-blur-md shadow-sm">
            <Filter className="h-4 w-4" /> Filters
          </button>
        </div>
      </div>

      {/* Triage List */}
      <div className="flex flex-col gap-4">
        {isLoading ? (
           <div className="text-center py-10 text-zinc-500 text-sm font-mono">Loading incidents...</div>
        ) : filteredIncidents.length === 0 ? (
           <div className="text-center py-10">
             <ShieldAlert className="h-10 w-10 text-emerald-500/50 mx-auto mb-4" />
             <div className="text-zinc-400 font-bold">No incidents found</div>
             <p className="text-zinc-500 text-sm mt-1">Systems are operating within acceptable parameters.</p>
           </div>
        ) : (
          filteredIncidents.map(inc => {
            const isCritical = inc.severity === 'critical';
            const isWarning = inc.severity === 'warning';
            
            return (
              <div 
                key={inc.incident_id} 
                className={cn(
                  "relative flex flex-col gap-4 overflow-hidden rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-5 transition-all w-full",
                  isCritical ? "border-rose-500/30 shadow-[0_0_15px_rgba(244,63,94,0.1)] animate-in" :
                  isWarning ? "border-amber-500/20" : "border-white/5",
                  isCritical && inc.status === 'OPEN' ? "animate-pulse shadow-[0_0_12px_rgba(239,68,68,0.2)]" : ""
                )}
              >
                {/* Left Accent Strip */}
                <div className={cn(
                  "absolute left-0 top-0 bottom-0 w-1",
                  isCritical ? "bg-rose-500" : isWarning ? "bg-amber-500" : "bg-sky-500"
                )} />

                <div className="flex justify-between items-start">
                  <div className="flex items-start gap-4">
                    <div className={cn(
                      "flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border shadow-inner",
                      isCritical ? "bg-rose-500/10 border-rose-500/30 text-rose-500" :
                      isWarning ? "bg-amber-500/10 border-amber-500/30 text-amber-500" :
                      "bg-sky-500/10 border-sky-500/30 text-sky-400"
                    )}>
                      {isCritical ? <AlertTriangle className="h-5 w-5" /> : 
                       isWarning ? <Activity className="h-5 w-5" /> : 
                       <Info className="h-5 w-5" />}
                    </div>
                    <div>
                      <div className="flex items-center gap-3 mb-1">
                        <span className="font-mono text-sm font-bold text-zinc-200">{inc.incident_id.substring(0, 8)}...</span>
                        <span className={cn(
                          "rounded px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider border",
                          inc.status === 'OPEN' ? "bg-rose-500/10 text-rose-400 border-rose-500/20" :
                          inc.status === 'ACKNOWLEDGED' ? "bg-amber-500/10 text-amber-400 border-amber-500/20" :
                          "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                        )}>
                          {inc.status}
                        </span>
                        <span className="flex items-center gap-1 text-[11px] font-mono text-zinc-500">
                           <Clock className="h-3 w-3" /> {new Date(inc.created_at).toLocaleTimeString()}
                        </span>
                      </div>
                      <h3 className="text-lg font-bold text-zinc-100 tracking-tight">{inc.type}</h3>
                      <p className="text-sm text-zinc-400 mt-1 max-w-3xl leading-relaxed">{inc.description}</p>
                      
                      <div className="mt-3 flex items-center gap-2">
                        <span className="text-[10px] font-bold uppercase text-zinc-500">Target Resource:</span>
                        <span className="rounded bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-2 py-0.5 font-mono text-[11px] text-zinc-300 shadow-inner">
                          {inc.device_id}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Operational Actions */}
                  <div className="flex flex-col items-end gap-2 shrink-0">
                    {/* State Transmutations */}
                    {inc.status === 'OPEN' && (
                       <button 
                         onClick={() => handleUpdateStatus(inc.incident_id, 'ACKNOWLEDGED')}
                         className="flex items-center gap-2 rounded bg-amber-500/10 border border-amber-500/30 px-4 py-2 text-xs font-bold uppercase tracking-wider text-amber-400 hover:bg-amber-500/20 transition-all shadow-[0_0_10px_rgba(245,158,11,0.1)]"
                       >
                         <AlertTriangle className="h-4 w-4" /> Acknowledge
                       </button>
                    )}
                    {inc.status === 'ACKNOWLEDGED' && (
                       <button 
                         onClick={() => handleUpdateStatus(inc.incident_id, 'RESOLVED')}
                         className="flex items-center gap-2 rounded bg-emerald-500/10 border border-emerald-500/30 px-4 py-2 text-xs font-bold uppercase tracking-wider text-emerald-400 hover:bg-emerald-500/20 transition-all shadow-[0_0_10px_rgba(16,185,129,0.1)]"
                       >
                         <Check className="h-4 w-4" /> Resolve Incident
                       </button>
                    )}
                    {inc.status === 'RESOLVED' && (
                      <span className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-emerald-500 px-4 py-2 opacity-50">
                        <CheckCircle2 className="h-4 w-4" /> Closed
                      </span>
                    )}

                    {inc.status !== 'RESOLVED' && (
                      <button 
                        onClick={() => setRunbookModalIncident(inc)}
                        className="text-[11px] font-bold uppercase text-sky-400 hover:text-sky-300 mt-2 hover:underline"
                      >
                        View Runbooks
                      </button>
                    )}
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* Runbook Modal */}
      {runbookModalIncident && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="w-full max-w-3xl rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1),-10px_0_30px_rgba(0,0,0,0.5)] overflow-hidden flex flex-col max-h-[90vh]">
            <div className="flex items-center justify-between border-b border-white/[0.08] p-4 bg-zinc-900/40 backdrop-blur-xl shadow-[inset_0_1px_0_rgba(255,255,255,0.1)]">
              <div>
                <h3 className="font-bold text-zinc-100 flex items-center gap-2">
                  <GitMerge className="h-5 w-5 text-sky-400" />
                  Self-Healing Playbook
                </h3>
                <p className="text-[11px] text-zinc-500 font-mono mt-1">
                  Tracing execution for {runbookModalIncident.incident_id.substring(0,8)} on {runbookModalIncident.device_id.substring(0,8)}
                </p>
              </div>
              <button 
                onClick={() => setRunbookModalIncident(null)}
                className="text-zinc-500 hover:text-zinc-300 transition-colors bg-zinc-800 hover:bg-zinc-700 p-2 rounded-full"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            
            <div className="flex-1 overflow-y-auto p-8 bg-[#0a0a0c] relative">
               <div className="absolute inset-0 pattern-grid-lg text-white/[0.02] bg-[length:40px_40px]" />
               
               <div className="relative z-10 flex flex-col gap-6 max-w-lg mx-auto">
                 {/* Step 1 */}
                 <div className="relative pl-8">
                   <div className="absolute left-0 top-1.5 w-4 h-4 rounded-full bg-emerald-500 border-4 border-[#0a0a0c] shadow-[0_0_10px_rgba(16,185,129,0.5)] z-10" />
                   <div className="absolute left-[7px] top-4 bottom-[-32px] w-0.5 bg-emerald-500/30" />
                   <div className="bg-zinc-900/80 border border-white/5 rounded-lg p-4 shadow-sm backdrop-blur-sm">
                     <h4 className="text-sm font-bold text-emerald-400 mb-1 font-mono">1. DETECT OUTAGE</h4>
                     <p className="text-xs text-zinc-400">Telemetry engine registered {runbookModalIncident.type}. Metric thresholds engaged.</p>
                   </div>
                 </div>

                 {/* Step 2 */}
                 <div className="relative pl-8">
                   <div className="absolute left-0 top-1.5 w-4 h-4 rounded-full bg-amber-500 border-4 border-[#0a0a0c] shadow-[0_0_10px_rgba(245,158,11,0.5)] z-10 animate-pulse" />
                   <div className="absolute left-[7px] top-4 bottom-[-32px] w-0.5 bg-zinc-800" />
                   <div className="bg-zinc-900/80 border border-amber-500/20 rounded-lg p-4 shadow-[0_0_15px_rgba(245,158,11,0.05)] backdrop-blur-sm">
                     <h4 className="text-sm font-bold text-amber-400 mb-1 font-mono">2. TRIGGER HOT-SWAP REPLICA</h4>
                     <p className="text-xs text-zinc-400">Orchestrating container payload replacement. Spooling instance clone in secondary zone.</p>
                     <div className="mt-3 flex items-center justify-between border-t border-white/5 pt-3">
                       <span className="text-[10px] uppercase font-bold text-zinc-500">Executing node template...</span>
                       <TerminalSquare className="h-4 w-4 text-amber-500/50" />
                     </div>
                   </div>
                 </div>

                 {/* Step 3 */}
                 <div className="relative pl-8">
                   <div className="absolute left-0 top-1.5 w-4 h-4 rounded-full bg-zinc-800 border-4 border-[#0a0a0c] z-10" />
                   <div className="bg-zinc-900/40 border border-white/5 rounded-lg p-4 backdrop-blur-sm opacity-50">
                     <h4 className="text-sm font-bold text-zinc-500 mb-1 font-mono">3. NOTIFY TEAM</h4>
                     <p className="text-xs text-zinc-500">Dispatching structured summary block to Telegram Core Bot Chat channel.</p>
                   </div>
                 </div>
               </div>

            </div>
          </div>
        </div>
      )}

    </div>
  );
}
