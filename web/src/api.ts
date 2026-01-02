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
