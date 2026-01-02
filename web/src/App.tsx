import { useEffect, useState } from 'react';
import type { BoardResponse, Issue, Column, IssueSummary } from './api';
import { fetchBoard, fetchIssue } from './api';
import './App.css';

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    pending: '#6b7280',
    in_progress: '#f59e0b',
    done: '#10b981',
    blocked: '#ef4444',
  };
  return (
    <span
      className="status-badge"
      style={{ backgroundColor: colors[status] || '#6b7280' }}
    >
      {status.replace('_', ' ')}
    </span>
  );
}

function IssueCard({
  issue,
  onClick,
}: {
  issue: IssueSummary;
  onClick: () => void;
}) {
  return (
    <div className="issue-card" onClick={onClick}>
      <div className="issue-title">{issue.title}</div>
      <div className="issue-meta">
        <StatusBadge status={issue.status} />
        <span className="issue-priority">{issue.priority}</span>
      </div>
    </div>
  );
}

function BoardColumn({
  column,
  onIssueClick,
}: {
  column: Column;
  onIssueClick: (id: string) => void;
}) {
  return (
    <div className="board-column">
      <div className="column-header">
        <span className="column-title">{column.label}</span>
        <span className="column-count">{column.count}</span>
      </div>
      <div className="column-issues">
        {column.issues.map((issue) => (
          <IssueCard
            key={issue.id}
            issue={issue}
            onClick={() => onIssueClick(issue.id)}
          />
        ))}
      </div>
    </div>
  );
}

function IssueDetail({
  issue,
  onClose,
}: {
  issue: Issue;
  onClose: () => void;
}) {
  return (
    <div className="issue-detail-overlay" onClick={onClose}>
      <div className="issue-detail" onClick={(e) => e.stopPropagation()}>
        <button className="close-btn" onClick={onClose}>
          &times;
        </button>
        <h2>{issue.title}</h2>
        <div className="issue-detail-meta">
          <StatusBadge status={issue.status} />
          <span className="issue-priority">[{issue.priority}]</span>
          <span className="issue-id">{issue.id}</span>
        </div>

        {issue.description && (
          <div className="issue-section">
            <h3>Description</h3>
            <p className="issue-description">{issue.description}</p>
          </div>
        )}

        {issue.done_when && issue.done_when.length > 0 && (
          <div className="issue-section">
            <h3>Done When</h3>
            <ul>
              {issue.done_when.map((item, i) => (
                <li key={i}>{item}</li>
              ))}
            </ul>
          </div>
        )}

        {issue.blocks && issue.blocks.length > 0 && (
          <div className="issue-section">
            <h3>Blocks</h3>
            <ul>
              {issue.blocks.map((dep) => (
                <li key={dep.id}>
                  {dep.title} <span className="dep-id">({dep.id})</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {issue.blocked_by && issue.blocked_by.length > 0 && (
          <div className="issue-section">
            <h3>Blocked By</h3>
            <ul>
              {issue.blocked_by.map((dep) => (
                <li key={dep.id}>
                  {dep.title} <span className="dep-id">({dep.id})</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {issue.children && issue.children.length > 0 && (
          <div className="issue-section">
            <h3>Children</h3>
            <ul>
              {issue.children.map((child) => (
                <li key={child.id}>
                  {child.title} <span className="dep-id">({child.id})</span>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

function App() {
  const [board, setBoard] = useState<BoardResponse | null>(null);
  const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadBoard();
    const interval = setInterval(loadBoard, 5000);
    return () => clearInterval(interval);
  }, []);

  async function loadBoard() {
    try {
      const data = await fetchBoard();
      setBoard(data);
      setError(null);
    } catch {
      setError(
        'Failed to connect to daemon. Is gvid running on localhost:7070?'
      );
    } finally {
      setLoading(false);
    }
  }

  async function handleIssueClick(id: string) {
    try {
      const issue = await fetchIssue(id);
      setSelectedIssue(issue);
    } catch (e) {
      console.error('Failed to fetch issue:', e);
    }
  }

  if (loading) {
    return (
      <div className="app">
        <div className="loading">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app">
        <div className="error">
          <h2>Connection Error</h2>
          <p>{error}</p>
          <button onClick={loadBoard}>Retry</button>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Gastown Viewer Intent</h1>
        <span className="total-count">{board?.total || 0} issues</span>
      </header>

      <div className="board">
        {board?.columns.map((column) => (
          <BoardColumn
            key={column.status}
            column={column}
            onIssueClick={handleIssueClick}
          />
        ))}
      </div>

      {selectedIssue && (
        <IssueDetail
          issue={selectedIssue}
          onClose={() => setSelectedIssue(null)}
        />
      )}
    </div>
  );
}

export default App;
