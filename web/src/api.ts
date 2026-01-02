// API client for gvid daemon

const API_BASE = 'http://localhost:7070/api/v1';

export type Status = 'pending' | 'in_progress' | 'done' | 'blocked';
export type Priority = 'critical' | 'high' | 'medium' | 'low';

export interface IssueSummary {
  id: string;
  title: string;
  status: Status;
  priority: Priority;
}

export interface Issue {
  id: string;
  title: string;
  description: string;
  status: Status;
  priority: Priority;
  parent?: IssueSummary;
  children: IssueSummary[];
  blocks: IssueSummary[];
  blocked_by: IssueSummary[];
  done_when: string[];
  created_at: string;
  updated_at: string;
}

export interface Column {
  status: Status;
  label: string;
  count: number;
  issues: IssueSummary[];
}

export interface BoardResponse {
  columns: Column[];
  total: number;
}

export interface HealthResponse {
  status: string;
  beads_initialized: boolean;
  version: string;
  bd_version?: string;
  error?: string;
}

export async function fetchHealth(): Promise<HealthResponse> {
  const res = await fetch(`${API_BASE}/health`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchBoard(): Promise<BoardResponse> {
  const res = await fetch(`${API_BASE}/board`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchIssue(id: string): Promise<Issue> {
  const res = await fetch(`${API_BASE}/issues/${encodeURIComponent(id)}`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

// Gas Town Types

export type AgentRole = 'mayor' | 'deacon' | 'witness' | 'refinery' | 'crew' | 'polecat';
export type AgentStatus = 'active' | 'idle' | 'stuck' | 'offline' | 'unknown';

export interface Agent {
  role: AgentRole;
  name: string;
  rig?: string;
  status: AgentStatus;
  session?: string;
  molecule?: string;
}

export interface Rig {
  name: string;
  path: string;
  remote?: string;
  witness?: Agent;
  refinery?: Agent;
  polecats: Agent[];
  crew: Agent[];
}

export interface Convoy {
  id: string;
  title: string;
  status: string;
  issues: string[];
  progress: number;
  total: number;
}

export interface Town {
  root: string;
  name?: string;
  rigs: Rig[];
  mayor?: Agent;
  deacon?: Agent;
  convoys: Convoy[];
}

export interface TownStatus {
  healthy: boolean;
  town_root: string;
  active_agents: number;
  total_agents: number;
  active_rigs: number;
  open_convoys: number;
  error?: string;
}

export interface AgentsResponse {
  agents: Agent[];
  total: number;
  active: number;
  offline: number;
}

export interface RigsResponse {
  rigs: Rig[];
  total: number;
}

export interface ConvoysResponse {
  convoys: Convoy[];
  total: number;
}

// Gas Town API calls

export async function fetchTownStatus(): Promise<TownStatus> {
  const res = await fetch(`${API_BASE}/town/status`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchTown(): Promise<Town> {
  const res = await fetch(`${API_BASE}/town`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchRigs(): Promise<RigsResponse> {
  const res = await fetch(`${API_BASE}/town/rigs`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchAgents(): Promise<AgentsResponse> {
  const res = await fetch(`${API_BASE}/town/agents`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function fetchConvoys(): Promise<ConvoysResponse> {
  const res = await fetch(`${API_BASE}/town/convoys`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}
