# Gastown Viewer Intent

> Mission Control dashboard for [Gastown](https://github.com/steveyegge/gastown) multi-agent workspaces.

## What It Does

**Gastown Viewer** provides real-time visibility into your Gas Town agent swarms:

- **Agent Dashboard**: See all agents (Mayor, Deacon, Witness, Refinery, Polecats, Crew) with live status
- **Rig Overview**: Monitor project rigs with agent health and activity
- **Convoy Tracking**: Track batch work progress across rigs
- **Beads Integration**: Kanban board view of issues managed by your agents
- **Web + TUI**: Browser dashboard or terminal interface

## Quickstart

### Prerequisites

- Go 1.22+
- Node.js 20+
- [Gastown](https://github.com/steveyegge/gastown) installed at `~/gt`
- [Beads](https://github.com/steveyegge/beads) (`bd` CLI in PATH)

### Run

```bash
# Start daemon + web UI
make dev

# Open http://localhost:5173
# Toggle between "Beads" and "Gas Town" tabs
```

### Verify

```bash
# Health check
curl http://localhost:7070/api/v1/health

# Gas Town status
curl http://localhost:7070/api/v1/town/status
# {"healthy":true,"active_agents":5,"total_agents":8,"active_rigs":2}

# List agents
curl http://localhost:7070/api/v1/town/agents
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Gastown Viewer Intent                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐      ┌──────────────┐      ┌──────────────┐  │
│   │   gvi-tui    │      │   Web UI     │      │  External    │  │
│   │  (Bubbletea) │      │ (React+Vite) │      │   Clients    │  │
│   └──────┬───────┘      └──────┬───────┘      └──────┬───────┘  │
│          │                     │                     │          │
│          └─────────────────────┼─────────────────────┘          │
│                                │                                 │
│                                ▼                                 │
│                    ┌───────────────────────┐                    │
│                    │       gvid Daemon     │                    │
│                    │     localhost:7070    │                    │
│                    └───────────┬───────────┘                    │
│                                │                                 │
│              ┌─────────────────┼─────────────────┐              │
│              ▼                                   ▼              │
│   ┌───────────────────────┐         ┌───────────────────────┐  │
│   │   Gastown Adapter     │         │    Beads Adapter      │  │
│   │   (reads ~/gt/)       │         │   (shells to `bd`)    │  │
│   └───────────┬───────────┘         └───────────┬───────────┘  │
│               │                                 │               │
│               ▼                                 ▼               │
│   ┌───────────────────────┐         ┌───────────────────────┐  │
│   │      Gas Town         │         │     .beads/ state     │  │
│   │  ~/gt (rigs, agents)  │         │   (issues, deps)      │  │
│   └───────────────────────┘         └───────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Gas Town Concepts

| Concept | Description |
|---------|-------------|
| **Town** | Workspace root (`~/gt`) containing all rigs and town-level agents |
| **Mayor** | Town coordinator - routes work across rigs |
| **Deacon** | Town patrol - monitors health and escalates issues |
| **Rig** | Project container with its own agent pool |
| **Witness** | Rig-level overseer - manages polecat lifecycle |
| **Refinery** | Merge queue processor for the rig |
| **Polecats** | Transient workers spawned for specific tasks |
| **Crew** | Persistent user-managed workers in a rig |
| **Convoy** | Batch work tracking across multiple rigs |

## API Endpoints

### Gas Town

| Endpoint | Description |
|----------|-------------|
| `GET /api/v1/town/status` | Town health, agent/rig counts |
| `GET /api/v1/town` | Full town structure |
| `GET /api/v1/town/rigs` | List all rigs |
| `GET /api/v1/town/rigs/:name` | Single rig details |
| `GET /api/v1/town/agents` | All agents with status |
| `GET /api/v1/town/convoys` | Active convoys |
| `GET /api/v1/town/mail/:address` | Agent mail inbox |

### Beads (Issues)

| Endpoint | Description |
|----------|-------------|
| `GET /api/v1/health` | Health check |
| `GET /api/v1/board` | Kanban board view |
| `GET /api/v1/issues` | List issues |
| `GET /api/v1/issues/:id` | Issue details |
| `GET /api/v1/graph` | Dependency graph |
| `GET /api/v1/events` | SSE event stream |

## Configuration

```bash
# Custom Gas Town location
go run ./cmd/gvid --town /path/to/gt

# Custom port
go run ./cmd/gvid --port 8080

# All options
go run ./cmd/gvid --help
```

## Project Structure

```
gastown-viewer-intent/
├── cmd/
│   ├── gvid/              # Daemon
│   └── gvi-tui/           # TUI client
├── internal/
│   ├── api/               # HTTP handlers
│   ├── gastown/           # Gas Town adapter (reads ~/gt)
│   ├── beads/             # Beads adapter (bd CLI)
│   └── model/             # Domain types
├── web/                   # React + Vite frontend
└── Makefile
```

## License

MIT

## Related Projects

- [Gastown](https://github.com/steveyegge/gastown) - Multi-agent workspace orchestrator
- [Beads](https://github.com/steveyegge/beads) - Local-first issue tracking with dependencies
