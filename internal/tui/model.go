package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/intent-solutions-io/gastown-viewer-intent/internal/model"
)

// View represents the current view mode.
type View int

const (
	ViewBoard View = iota
	ViewIssue
)

// Model is the main TUI model.
type Model struct {
	client   *Client
	board    *BoardResponse
	issue    *model.Issue
	view     View
	err      error
	loading  bool
	spinner  spinner.Model
	help     help.Model
	keys     keyMap
	cursor   int // selected column
	issueCur int // selected issue in column
	width    int
	height   int
}

// keyMap defines keybindings.
type keyMap struct {
	Left    key.Binding
	Right   key.Binding
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Back    key.Binding
	Refresh key.Binding
	Quit    key.Binding
	Help    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Enter, k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Up, k.Down},
		{k.Enter, k.Back, k.Refresh, k.Quit},
	}
}

var defaultKeys = keyMap{
	Left:    key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("<-/h", "left")),
	Right:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("->/l", "right")),
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("^/k", "up")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("v/j", "down")),
	Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
	Back:    key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "back")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}

// New creates a new TUI model.
func New(apiURL string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		client:  NewClient(apiURL),
		spinner: s,
		help:    help.New(),
		keys:    defaultKeys,
		loading: true,
		width:   80,
		height:  24,
	}
}

// Messages
type boardMsg *BoardResponse
type issueMsg *model.Issue
type errMsg error

func (m Model) fetchBoard() tea.Msg {
	board, err := m.client.Board()
	if err != nil {
		return errMsg(err)
	}
	return boardMsg(board)
}

func (m Model) fetchIssue(id string) tea.Cmd {
	return func() tea.Msg {
		issue, err := m.client.Issue(id)
		if err != nil {
			return errMsg(err)
		}
		return issueMsg(issue)
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchBoard)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchBoard)

		case key.Matches(msg, m.keys.Back):
			if m.view == ViewIssue {
				m.view = ViewBoard
				m.issue = nil
			}
			return m, nil

		case key.Matches(msg, m.keys.Left):
			if m.view == ViewBoard && m.board != nil && m.cursor > 0 {
				m.cursor--
				m.issueCur = 0
			}
			return m, nil

		case key.Matches(msg, m.keys.Right):
			if m.view == ViewBoard && m.board != nil && m.cursor < len(m.board.Columns)-1 {
				m.cursor++
				m.issueCur = 0
			}
			return m, nil

		case key.Matches(msg, m.keys.Up):
			if m.view == ViewBoard && m.issueCur > 0 {
				m.issueCur--
			}
			return m, nil

		case key.Matches(msg, m.keys.Down):
			if m.view == ViewBoard && m.board != nil && m.cursor < len(m.board.Columns) {
				col := m.board.Columns[m.cursor]
				if m.issueCur < len(col.Issues)-1 {
					m.issueCur++
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			if m.view == ViewBoard && m.board != nil {
				if m.cursor < len(m.board.Columns) {
					col := m.board.Columns[m.cursor]
					if m.issueCur < len(col.Issues) {
						issue := col.Issues[m.issueCur]
						m.loading = true
						return m, tea.Batch(m.spinner.Tick, m.fetchIssue(issue.ID))
					}
				}
			}
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case boardMsg:
		m.loading = false
		m.board = msg
		m.err = nil
		return m, nil

	case issueMsg:
		m.loading = false
		m.issue = msg
		m.view = ViewIssue
		m.err = nil
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		return m, nil
	}

	return m, nil
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(24)

	selectedColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(0, 1).
				Width(24)

	issueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	selectedIssueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	statusPending = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	statusInProgress = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	statusDone = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	statusBlocked = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
)

// View renders the model.
func (m Model) View() string {
	if m.err != nil {
		return m.viewError()
	}

	if m.loading && m.board == nil {
		return m.viewLoading()
	}

	var content string
	switch m.view {
	case ViewBoard:
		content = m.viewBoard()
	case ViewIssue:
		content = m.viewIssue()
	}

	helpView := m.help.View(m.keys)
	return content + "\n\n" + helpView
}

func (m Model) viewLoading() string {
	return fmt.Sprintf("\n  %s Loading...\n\n", m.spinner.View())
}

func (m Model) viewError() string {
	return fmt.Sprintf("\n  %s\n\n  %s\n\n  Press 'r' to retry or 'q' to quit.\n",
		errorStyle.Render("Error connecting to daemon:"),
		m.err.Error())
}

func (m Model) viewBoard() string {
	if m.board == nil {
		return m.viewLoading()
	}

	var b strings.Builder

	// Title
	title := titleStyle.Render("Gastown Viewer Intent")
	if m.loading {
		title += " " + m.spinner.View()
	}
	b.WriteString(title + "\n\n")

	// Columns
	var columns []string
	for i, col := range m.board.Columns {
		colStyle := columnStyle
		if i == m.cursor {
			colStyle = selectedColumnStyle
		}

		// Column header with status color
		var headerStyle lipgloss.Style
		switch col.Status {
		case model.StatusPending:
			headerStyle = statusPending
		case model.StatusInProgress:
			headerStyle = statusInProgress
		case model.StatusDone:
			headerStyle = statusDone
		case model.StatusBlocked:
			headerStyle = statusBlocked
		default:
			headerStyle = statusPending
		}

		header := headerStyle.Render(fmt.Sprintf("%s (%d)", col.Label, col.Count))

		// Issues
		var issues []string
		for j, issue := range col.Issues {
			style := issueStyle
			if i == m.cursor && j == m.issueCur {
				style = selectedIssueStyle
			}
			// Truncate title
			title := issue.Title
			if len(title) > 20 {
				title = title[:17] + "..."
			}
			issues = append(issues, style.Render(title))
		}

		content := header + "\n" + strings.Join(issues, "\n")
		columns = append(columns, colStyle.Render(content))
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, columns...))

	return b.String()
}

func (m Model) viewIssue() string {
	if m.issue == nil {
		return m.viewLoading()
	}

	var b strings.Builder

	// Back navigation hint
	b.WriteString(labelStyle.Render("< Press ESC to go back") + "\n\n")

	// Title
	b.WriteString(titleStyle.Render(m.issue.Title) + "\n")

	// Status and priority
	var statusStyle lipgloss.Style
	switch m.issue.Status {
	case model.StatusPending:
		statusStyle = statusPending
	case model.StatusInProgress:
		statusStyle = statusInProgress
	case model.StatusDone:
		statusStyle = statusDone
	case model.StatusBlocked:
		statusStyle = statusBlocked
	}
	b.WriteString(fmt.Sprintf("%s  %s\n\n",
		statusStyle.Render(string(m.issue.Status)),
		labelStyle.Render(fmt.Sprintf("[%s]", m.issue.Priority))))

	// ID
	b.WriteString(labelStyle.Render("ID: ") + m.issue.ID + "\n\n")

	// Description
	if m.issue.Description != "" {
		b.WriteString(labelStyle.Render("Description:\n"))
		// Wrap description
		desc := m.issue.Description
		if len(desc) > 500 {
			desc = desc[:497] + "..."
		}
		b.WriteString(desc + "\n\n")
	}

	// Done when
	if len(m.issue.DoneWhen) > 0 {
		b.WriteString(labelStyle.Render("Done when:\n"))
		for _, item := range m.issue.DoneWhen {
			b.WriteString("  - " + item + "\n")
		}
		b.WriteString("\n")
	}

	// Dependencies
	if len(m.issue.Blocks) > 0 {
		b.WriteString(labelStyle.Render("Blocks:\n"))
		for _, dep := range m.issue.Blocks {
			b.WriteString(fmt.Sprintf("  - %s (%s)\n", dep.Title, dep.ID))
		}
		b.WriteString("\n")
	}

	if len(m.issue.BlockedBy) > 0 {
		b.WriteString(labelStyle.Render("Blocked by:\n"))
		for _, dep := range m.issue.BlockedBy {
			b.WriteString(fmt.Sprintf("  - %s (%s)\n", dep.Title, dep.ID))
		}
		b.WriteString("\n")
	}

	return detailStyle.Width(m.width - 4).Render(b.String())
}
