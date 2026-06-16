'use client';

import React, { useState } from 'react';
import { 
  Activity, PlayCircle, GitMerge, Clock, Workflow, Plus, Trash2, 
  Settings, Code, ToggleLeft, ToggleRight, X, ChevronRight, CheckCircle2,
  AlertTriangle, Edit3
} from 'lucide-react';
import { cn } from '@/lib/utils';

// --- TYPES ---
type AlertRule = {
  rule_id: string;
  name: string;
  metric: 'cpu' | 'ram' | 'disk' | 'network';
  operator: '>' | '<' | '==' | '>=';
  threshold: number;
  duration: string;
  is_muted: boolean;
}

type PlaybookStep = {
  id: string;
  name: string;
  type: string;
  action: string;
}

type Playbook = {
  playbook_id: string;
  name: string;
  steps: PlaybookStep[];
}

type CronJob = {
  cron_id: string;
  expression: string;
  script_id: string;
  device_group_id: string;
}

// --- INITIAL STATE ---
const INITIAL_RULES: AlertRule[] = [
  { rule_id: 'RULE-001', name: 'High CPU Triage', metric: 'cpu', operator: '>=', threshold: 90, duration: '5m', is_muted: false },
  { rule_id: 'RULE-002', name: 'DB Memory Exhaustion', metric: 'ram', operator: '>=', threshold: 85, duration: '10m', is_muted: true }
];

const INITIAL_PLAYBOOKS: Playbook[] = [
  { 
    playbook_id: 'PB-100', 
    name: 'Auto-Scale Web Tier', 
    steps: [
      { id: 'S1', name: 'Detect Load', type: 'trigger', action: 'Listen on Port 80' },
      { id: 'S2', name: 'Spin Up Node', type: 'execution', action: 'Provision Instance' }
    ] 
  },
  { 
    playbook_id: 'PB-101', 
    name: 'Container Self-Healing', 
    steps: [
      { id: 'S1', name: 'OOM Killed', type: 'trigger', action: 'Filter Kubelet Logs' },
      { id: 'S2', name: 'Restart Pod', type: 'execution', action: 'kubectl cycle' },
      { id: 'S3', name: 'Alert Channel', type: 'notification', action: 'Slack #ops' }
    ] 
  }
];

const INITIAL_CRON: CronJob[] = [
  { cron_id: 'CRON-001', expression: '0 2 * * *', script_id: 'backup_db.sh', device_group_id: 'grp_prod_db' },
  { cron_id: 'CRON-002', expression: '0 0 * * 0', script_id: 'rotate_logs.py', device_group_id: 'grp_all_servers' }
];

export default function AutomationPage() {
  const [activeTab, setActiveTab] = useState<'rules' | 'playbooks' | 'cron'>('rules');
  
  const [rules, setRules] = useState<AlertRule[]>(INITIAL_RULES);
  const [playbooks, setPlaybooks] = useState<Playbook[]>(INITIAL_PLAYBOOKS);
  const [crons, setCrons] = useState<CronJob[]>(INITIAL_CRON);
  
  // Rule Modal Form
  const [isRuleModalOpen, setIsRuleModalOpen] = useState(false);
  const [ruleName, setRuleName] = useState('');
  const [ruleMetric, setRuleMetric] = useState<'cpu' | 'ram' | 'disk' | 'network'>('cpu');
  const [ruleOp, setRuleOp] = useState<'>' | '<' | '==' | '>='>('>=');
  const [ruleThreshold, setRuleThreshold] = useState('');
  const [ruleDuration, setRuleDuration] = useState('');

  // Cron Modal Form
  const [isCronModalOpen, setIsCronModalOpen] = useState(false);
  const [cronExp, setCronExp] = useState('');
  const [cronScript, setCronScript] = useState('');
  const [cronGroup, setCronGroup] = useState('');

  // Playbook Modal Form
  const [isPlaybookModalOpen, setIsPlaybookModalOpen] = useState(false);
  const [editingPlaybookId, setEditingPlaybookId] = useState<string | null>(null);
  const [pbName, setPbName] = useState('');
  const [pbSteps, setPbSteps] = useState<PlaybookStep[]>([]);

  // Handlers
  const toggleRuleMute = (id: string) => {
    setRules(prev => prev.map(r => r.rule_id === id ? { ...r, is_muted: !r.is_muted } : r));
  };
  const deleteRule = (id: string) => {
    setRules(prev => prev.filter(r => r.rule_id !== id));
  };

  const handleRuleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!ruleName || !ruleThreshold) return;
    const newRule: AlertRule = {
      rule_id: `RULE-${Math.floor(Math.random()*1000).toString().padStart(3, '0')}`,
      name: ruleName,
      metric: ruleMetric,
      operator: ruleOp,
      threshold: Number(ruleThreshold),
      duration: ruleDuration || '1m',
      is_muted: false
    };
    setRules(prev => [...prev, newRule]);
    setIsRuleModalOpen(false);
    setRuleName(''); setRuleThreshold(''); setRuleDuration('');
  };

  const deleteCron = (id: string) => {
    setCrons(prev => prev.filter(c => c.cron_id !== id));
  }

  const handleCronSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!cronExp || !cronScript) return;
    const newCron: CronJob = {
      cron_id: `CRON-${Math.floor(Math.random()*1000).toString().padStart(3, '0')}`,
      expression: cronExp,
      script_id: cronScript,
      device_group_id: cronGroup || 'default'
    };
    setCrons(prev => [...prev, newCron]);
    setIsCronModalOpen(false);
    setCronExp(''); setCronScript(''); setCronGroup('');
  }

  const openPlaybookModal = (pb?: Playbook) => {
    if (pb) {
      setEditingPlaybookId(pb.playbook_id);
      setPbName(pb.name);
      setPbSteps(pb.steps);
    } else {
      setEditingPlaybookId(null);
      setPbName('');
      setPbSteps([]);
    }
    setIsPlaybookModalOpen(true);
  };

  const handlePlaybookSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!pbName) return;
    
    if (editingPlaybookId) {
      setPlaybooks(prev => prev.map(p => p.playbook_id === editingPlaybookId ? { ...p, name: pbName, steps: pbSteps } : p));
    } else {
      const newPb: Playbook = {
        playbook_id: `PB-${Math.floor(Math.random()*1000).toString().padStart(3, '0')}`,
        name: pbName,
        steps: pbSteps
      };
      setPlaybooks(prev => [...prev, newPb]);
    }
    setIsPlaybookModalOpen(false);
  };

  const handleAddStep = () => {
    setPbSteps(prev => [...prev, {
      id: `S${Math.random()}`,
      name: '',
      type: 'execution',
      action: ''
    }]);
  };
  
  const handleUpdateStep = (id: string, field: keyof PlaybookStep, value: string) => {
    setPbSteps(prev => prev.map(s => s.id === id ? { ...s, [field]: value } : s));
  };

  const handleDeleteStep = (id: string) => {
    setPbSteps(prev => prev.filter(s => s.id !== id));
  };

  return (
    <div className="flex flex-col space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 pb-10">
      
      {/* Header */}
      <div className="border-b border-white/5 pb-6">
        <h1 className="text-3xl font-black tracking-tighter text-zinc-100 mb-2">Automation Engine</h1>
        <p className="text-sm text-zinc-500 cursor-default">Autonomous response controls, self-healing topological playbooks, and scheduled shell executions.</p>
      </div>

      {/* Tabs Layout */}
      <div className="flex items-center gap-2 border-b border-white/5 pb-px">
        <button 
          onClick={(e) => { e.preventDefault(); setActiveTab('rules'); }}
          className={cn("px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors border-b-2", activeTab === 'rules' ? "border-sky-500 text-sky-400" : "border-transparent text-zinc-500 hover:text-zinc-300")}
        >
          Alert Rules
        </button>
        <button 
          onClick={(e) => { e.preventDefault(); setActiveTab('playbooks'); }}
          className={cn("px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors border-b-2", activeTab === 'playbooks' ? "border-emerald-500 text-emerald-400" : "border-transparent text-zinc-500 hover:text-zinc-300")}
        >
          Playbooks
        </button>
        <button 
          onClick={(e) => { e.preventDefault(); setActiveTab('cron'); }}
          className={cn("px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors border-b-2", activeTab === 'cron' ? "border-amber-500 text-amber-400" : "border-transparent text-zinc-500 hover:text-zinc-300")}
        >
          Cron Jobs
        </button>
      </div>

      {/* TABS CONTENT */}
      <div className="mt-6">
        
        {/* --- TAB: ALERT RULES --- */}
        {activeTab === 'rules' && (
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-sm text-zinc-400 font-mono">Active environmental triggers computing in real-time.</h2>
              <button 
                onClick={(e) => { e.preventDefault(); setIsRuleModalOpen(true); }}
                className="flex items-center gap-2 bg-sky-500/10 text-sky-400 border border-sky-500/20 px-4 py-2 rounded text-xs font-bold uppercase tracking-wider hover:bg-sky-500/20 transition-colors"
              >
                <Plus className="w-4 h-4" /> Create Rule
              </button>
            </div>
            
            <div className="rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] overflow-hidden shadow-[0_8px_30px_rgba(0,0,0,0.12)]">
              <table className="w-full text-left text-sm whitespace-nowrap">
                <thead className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] border-b border-white/5 text-zinc-400 font-medium">
                  <tr>
                    <th className="px-6 py-4 font-mono font-normal">Rule ID</th>
                    <th className="px-6 py-4 font-mono font-normal">Name</th>
                    <th className="px-6 py-4 font-mono font-normal">Condition</th>
                    <th className="px-6 py-4 font-mono font-normal">Status</th>
                    <th className="px-6 py-4 text-right font-mono font-normal">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5 text-zinc-300">
                  {rules.map(rule => (
                    <tr key={rule.rule_id} className="transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95">
                      <td className="px-6 py-4 font-mono text-zinc-500 text-xs">{rule.rule_id}</td>
                      <td className="px-6 py-4 font-medium text-zinc-200">{rule.name}</td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2 font-mono text-xs">
                          <span className="text-rose-400 uppercase">{rule.metric}</span>
                          <span className="text-zinc-500">{rule.operator}</span>
                          <span className="text-zinc-200">{rule.threshold}</span>
                          <span className="text-zinc-600">for</span>
                          <span className="text-sky-400">{rule.duration}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <button 
                          onClick={(e) => { e.preventDefault(); toggleRuleMute(rule.rule_id); }}
                          className="flex items-center gap-2 focus:outline-none"
                        >
                          {rule.is_muted ? (
                            <ToggleLeft className="w-5 h-5 text-zinc-600" />
                          ) : (
                            <ToggleRight className="w-5 h-5 text-emerald-500 animate-pulse" />
                          )}
                          <span className={cn("text-xs font-bold uppercase", rule.is_muted ? "text-zinc-600" : "text-emerald-500 animate-pulse")}>
                            {rule.is_muted ? 'Muted' : 'Active'}
                          </span>
                        </button>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <button 
                          onClick={(e) => { e.preventDefault(); deleteRule(rule.rule_id); }}
                          className="text-zinc-600 hover:text-rose-500 animate-pulse transition-colors p-2"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                  {rules.length === 0 && (
                     <tr><td colSpan={5} className="px-6 py-12 text-center text-zinc-500 font-mono text-xs uppercase">No trigger rules active</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* --- TAB: PLAYBOOKS --- */}
        {activeTab === 'playbooks' && (
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-sm text-zinc-400 font-mono">Self-healing sequence routines mapped to system triggers.</h2>
              <button 
                onClick={(e) => { e.preventDefault(); openPlaybookModal(); }}
                className="flex items-center gap-2 bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-4 py-2 rounded text-xs font-bold uppercase tracking-wider hover:bg-emerald-500/20 transition-colors"
              >
                <GitMerge className="w-4 h-4" /> New Playbook Workflow
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {playbooks.map(pb => (
                <div key={pb.playbook_id} className="rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] p-$1 flex flex-col">
                  <div className="flex justify-between items-start mb-6">
                    <div>
                      <div className="text-xs font-mono text-emerald-500 animate-pulse mb-1">{pb.playbook_id}</div>
                      <h3 className="font-bold text-lg text-zinc-100">{pb.name}</h3>
                    </div>
                    <div>
                      <button className="text-zinc-600 hover:text-emerald-500 animate-pulse transition-all duration-300 ease-out active:scale-95 mr-3" onClick={(e) => {
                        e.preventDefault();
                        openPlaybookModal(pb);
                      }}>
                        <Edit3 className="w-4 h-4" />
                      </button>
                      <button className="text-zinc-600 hover:text-rose-500 animate-pulse transition-all duration-300 ease-out active:scale-95" onClick={(e) => {
                        e.preventDefault();
                        setPlaybooks(prev => prev.filter(p => p.playbook_id !== pb.playbook_id));
                      }}>
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  <div className="flex-1 space-y-3 relative before:absolute before:inset-y-0 before:left-[19px] before:w-0.5 before:bg-white/5 pl-2">
                    {pb.steps.map((step, idx) => (
                      <div key={step.id} className="relative pl-10">
                        <div className="absolute left-0 top-1.5 w-6 h-6 rounded-full bg-zinc-900 border-2 border-zinc-800 flex items-center justify-center -translate-x-[9px]">
                          {idx === 0 ? <Activity className="w-3 h-3 text-amber-500" /> : <PlayCircle className="w-3 h-3 text-sky-400" />}
                        </div>
                        <div className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] border border-white/5 rounded-lg p-3">
                          <div className="text-[10px] uppercase font-bold text-zinc-500 mb-1">{step.type}</div>
                          <div className="text-sm font-medium text-zinc-300">{step.name}</div>
                          <div className="text-xs font-mono text-zinc-500 mt-2 bg-black/30 p-2 rounded">{step.action}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
              {playbooks.length === 0 && (
                <div className="col-span-full py-12 text-center text-zinc-500 font-mono text-xs uppercase border border-dashed border-white/10 rounded-xl">No playbooks configured</div>
              )}
            </div>
          </div>
        )}

        {/* --- TAB: CRON JOBS --- */}
        {activeTab === 'cron' && (
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-sm text-zinc-400 font-mono">Unsupervised scheduling matrix for distributed targets.</h2>
              <button 
                onClick={(e) => { e.preventDefault(); setIsCronModalOpen(true); }}
                className="flex items-center gap-2 bg-amber-500/10 text-amber-400 border border-amber-500/20 px-4 py-2 rounded text-xs font-bold uppercase tracking-wider hover:bg-amber-500/20 transition-colors"
              >
                <Clock className="w-4 h-4" /> Schedule Task
              </button>
            </div>

            <div className="rounded-xl bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] overflow-hidden shadow-[0_8px_30px_rgba(0,0,0,0.12)]">
              <table className="w-full text-left text-sm whitespace-nowrap">
                <thead className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] border-b border-white/5 text-zinc-400 font-medium">
                  <tr>
                    <th className="px-6 py-4 font-mono font-normal">Job ID</th>
                    <th className="px-6 py-4 font-mono font-normal">Schedule</th>
                    <th className="px-6 py-4 font-mono font-normal">Target Script</th>
                    <th className="px-6 py-4 font-mono font-normal">Distribution Group</th>
                    <th className="px-6 py-4 text-right font-mono font-normal">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5 text-zinc-300">
                  {crons.map(cron => (
                    <tr key={cron.cron_id} className="transition-all duration-300 ease-out hover:bg-white/[0.02] hover:scale-[1.01] active:scale-95">
                      <td className="px-6 py-4 font-mono text-zinc-500 text-xs">{cron.cron_id}</td>
                      <td className="px-6 py-4">
                        <span className="font-mono text-sm bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] px-3 py-1 rounded text-amber-400">
                          {cron.expression}
                        </span>
                      </td>
                      <td className="px-6 py-4 font-mono text-zinc-300 text-xs flex items-center gap-2">
                        <Code className="w-4 h-4 text-zinc-600" /> {cron.script_id}
                      </td>
                      <td className="px-6 py-4 font-mono text-sky-400 text-xs">{cron.device_group_id}</td>
                      <td className="px-6 py-4 text-right">
                        <button 
                          onClick={(e) => { e.preventDefault(); deleteCron(cron.cron_id); }}
                          className="text-zinc-600 hover:text-rose-500 animate-pulse transition-colors p-2"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                  {crons.length === 0 && (
                     <tr><td colSpan={5} className="px-6 py-12 text-center text-zinc-500 font-mono text-xs uppercase">No cron jobs scheduled</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}

      </div>

      {/* --- RULE MODAL --- */}
      {isRuleModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-xl shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-700">
            <div className="flex justify-between items-center p-6 border-b border-white/5">
              <h3 className="font-black tracking-tighter text-zinc-200 text-xl">Create Alert Rule</h3>
              <button 
                onClick={(e) => { e.preventDefault(); setIsRuleModalOpen(false); }}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleRuleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Rule Name</label>
                <input 
                  type="text" 
                  value={ruleName}
                  onChange={e => setRuleName(e.target.value)}
                  className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 focus:outline-none focus:ring-2 focus:ring-sky-500/50 transition-shadow"
                  placeholder="e.g. Critical Heat Warning"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Target Metric</label>
                  <select 
                    value={ruleMetric}
                    onChange={e => setRuleMetric(e.target.value as any)}
                    className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 focus:outline-none focus:ring-2 focus:ring-sky-500/50 appearance-none"
                  >
                    <option value="cpu">CPU Usage (%)</option>
                    <option value="ram">RAM Usage (%)</option>
                    <option value="disk">Disk I/O</option>
                    <option value="network">Network Volume</option>
                  </select>
                </div>
                <div>
                  <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Operator</label>
                  <select 
                    value={ruleOp}
                    onChange={e => setRuleOp(e.target.value as any)}
                    className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 font-mono focus:outline-none focus:ring-2 focus:ring-sky-500/50 appearance-none"
                  >
                    <option value=">">&gt; (Greater than)</option>
                    <option value="<">&lt; (Less than)</option>
                    <option value=">=">&gt;= (GTE)</option>
                    <option value="==">== (Equals)</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Threshold Value</label>
                  <input 
                    type="number" 
                    value={ruleThreshold}
                    onChange={e => setRuleThreshold(e.target.value)}
                    className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 focus:outline-none focus:ring-2 focus:ring-sky-500/50 transition-shadow"
                    placeholder="90"
                    required
                  />
                </div>
                <div>
                  <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Sustain Duration</label>
                  <input 
                    type="text" 
                    value={ruleDuration}
                    onChange={e => setRuleDuration(e.target.value)}
                    className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 focus:outline-none focus:ring-2 focus:ring-sky-500/50 transition-shadow"
                    placeholder="e.g. 5m"
                    required
                  />
                </div>
              </div>

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button" 
                  onClick={(e) => { e.preventDefault(); setIsRuleModalOpen(false); }}
                  className="px-4 py-2 text-sm font-medium text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="px-6 py-2 text-sm font-bold bg-sky-500 hover:bg-sky-400 text-white rounded-lg transition-all duration-300 ease-out active:scale-95 shadow-lg shadow-sky-500/20"
                >
                  Install Rule
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- CRON MODAL --- */}
      {isCronModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-xl shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-700">
            <div className="flex justify-between items-center p-6 border-b border-white/5">
              <h3 className="font-black tracking-tighter text-zinc-200 text-xl">Schedule Cron Task</h3>
              <button 
                onClick={(e) => { e.preventDefault(); setIsCronModalOpen(false); }}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleCronSubmit} className="p-6 space-y-4">
              
              <div>
                <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Cron Expression</label>
                <input 
                  type="text" 
                  value={cronExp}
                  onChange={e => setCronExp(e.target.value)}
                  className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-amber-400 font-mono focus:outline-none focus:ring-2 focus:ring-amber-500/50 transition-shadow"
                  placeholder="0 5 * * *"
                  required
                />
              </div>

              <div>
                <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Target Script ID</label>
                <input 
                  type="text" 
                  value={cronScript}
                  onChange={e => setCronScript(e.target.value)}
                  className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 font-mono focus:outline-none focus:ring-2 focus:ring-amber-500/50 transition-shadow"
                  placeholder="e.g. node_maintenance.sh"
                  required
                />
              </div>

              <div>
                <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Distribution Group</label>
                <input 
                  type="text" 
                  value={cronGroup}
                  onChange={e => setCronGroup(e.target.value)}
                  className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-sky-400 font-mono focus:outline-none focus:ring-2 focus:ring-amber-500/50 transition-shadow"
                  placeholder="e.g. grp_frontend_nodes"
                />
              </div>

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button" 
                  onClick={(e) => { e.preventDefault(); setIsCronModalOpen(false); }}
                  className="px-4 py-2 text-sm font-medium text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="px-6 py-2 text-sm font-bold bg-amber-500 hover:bg-amber-400 text-zinc-950 rounded-lg transition-all duration-300 ease-out active:scale-95 shadow-lg shadow-amber-500/20"
                >
                  Confirm Schedule
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- PLAYBOOK MODAL --- */}
      {isPlaybookModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-xl shadow-2xl w-full max-w-2xl overflow-hidden animate-in fade-in slide-in-from-bottom-4 duration-700">
            <div className="flex justify-between items-center p-6 border-b border-white/5">
              <h3 className="font-black tracking-tighter text-zinc-200 text-xl">{editingPlaybookId ? 'Edit Playbook' : 'Create Playbook'}</h3>
              <button 
                onClick={(e) => { e.preventDefault(); setIsPlaybookModalOpen(false); }}
                className="text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handlePlaybookSubmit} className="p-6 space-y-6">
              
              <div>
                <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">Playbook Name</label>
                <input 
                  type="text" 
                  value={pbName}
                  onChange={e => setPbName(e.target.value)}
                  className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded-lg px-4 py-2.5 text-sm text-zinc-200 focus:outline-none focus:ring-2 focus:ring-emerald-500/50 transition-shadow"
                  placeholder="e.g. Service Restart Protocol"
                  required
                />
              </div>

              <div>
                <div className="flex items-center justify-between mb-4">
                  <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-500">Workflow Steps</label>
                  <button 
                    type="button"
                    onClick={handleAddStep}
                    className="flex items-center gap-1.5 bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-3 py-1.5 rounded text-[10px] font-bold uppercase tracking-wider hover:bg-emerald-500/20 transition-all duration-300 ease-out active:scale-95"
                  >
                    <Plus className="w-3 h-3" /> Add Step
                  </button>
                </div>
                
                <div className="space-y-3 max-h-[40vh] overflow-y-auto pr-2">
                  {pbSteps.map((step, idx) => (
                    <div key={step.id} className="bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] border border-white/5 rounded-lg p-4 flex gap-4 items-start">
                      <div className="w-8 h-8 rounded-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] flex items-center justify-center shrink-0 font-mono text-xs text-zinc-500">
                        {idx + 1}
                      </div>
                      <div className="flex-1 grid grid-cols-2 gap-3">
                        <div className="col-span-2 sm:col-span-1">
                          <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-600 mb-1">Step Name</label>
                          <input 
                            type="text"
                            value={step.name}
                            onChange={e => handleUpdateStep(step.id, 'name', e.target.value)}
                            className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-3 py-1.5 text-xs text-zinc-300 focus:outline-none focus:border-white/20"
                            placeholder="e.g. Stop Service"
                            required
                          />
                        </div>
                        <div className="col-span-2 sm:col-span-1">
                          <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-600 mb-1">Step Type</label>
                          <select 
                            value={step.type}
                            onChange={e => handleUpdateStep(step.id, 'type', e.target.value)}
                            className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-3 py-1.5 text-xs text-zinc-300 focus:outline-none focus:border-white/20 appearance-none"
                          >
                            <option value="trigger">Trigger</option>
                            <option value="execution">Execution</option>
                            <option value="notification">Notification</option>
                          </select>
                        </div>
                        <div className="col-span-2">
                          <label className="block text-[10px] font-bold uppercase tracking-wider text-zinc-600 mb-1">Action / Command</label>
                          <input 
                            type="text"
                            value={step.action}
                            onChange={e => handleUpdateStep(step.id, 'action', e.target.value)}
                            className="w-full bg-zinc-900/40 backdrop-blur-xl border border-white/[0.08] shadow-[inset_0_1px_0_rgba(255,255,255,0.1)] rounded px-3 py-1.5 text-xs font-mono text-emerald-400 focus:outline-none focus:border-white/20"
                            placeholder="e.g. systemctl stop nginx"
                            required
                          />
                        </div>
                      </div>
                      <button 
                        type="button"
                        onClick={() => handleDeleteStep(step.id)}
                        className="text-zinc-600 hover:text-rose-500 animate-pulse transition-colors mt-6 shrink-0"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                  {pbSteps.length === 0 && (
                    <div className="text-center py-8 text-zinc-600 text-xs font-mono border border-dashed border-white/5 rounded-lg bg-zinc-900/20">
                      No steps defined. Add a step to begin your workflow.
                    </div>
                  )}
                </div>
              </div>

              <div className="pt-4 flex justify-end gap-3 border-t border-white/5 mt-6">
                <button 
                  type="button" 
                  onClick={(e) => { e.preventDefault(); setIsPlaybookModalOpen(false); }}
                  className="px-4 py-2 text-sm font-medium text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="px-6 py-2 text-sm font-bold bg-emerald-500 hover:bg-emerald-400 text-zinc-950 rounded-lg transition-all duration-300 ease-out active:scale-95 shadow-lg shadow-emerald-500/20"
                >
                  {editingPlaybookId ? 'Save Playbook' : 'Create Playbook'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
