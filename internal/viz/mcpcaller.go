package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mcpTool struct {
	Name   string
	Desc   string
	Input  string
	Output string
}

type mcpSrv struct {
	Name  string
	Tools []mcpTool
}

// MCPCallerModel lets users browse servers and invoke tools.
type MCPCallerModel struct {
	width       int
	height      int
	servers     []mcpSrv
	serverIdx   int
	toolIdx     int
	called      bool
	callOutput  string
	inToolView  bool
}

func NewMCPCallerModel(w, h int) Model {
	return &MCPCallerModel{
		width:  w,
		height: h,
		servers: []mcpSrv{
			{
				Name: "GitHub",
				Tools: []mcpTool{
					{Name: "list_issues", Desc: "List open issues", Input: "repo: owner/myapp", Output: "#42 Fix login timeout\n#43 Add dark mode\n#44 Update dependencies"},
					{Name: "create_pr", Desc: "Create a pull request", Input: "title: \"Add validation\"\nbranch: feature/validate", Output: "PR #15 created successfully\nURL: github.com/owner/myapp/pull/15"},
					{Name: "read_file", Desc: "Read a file from the repo", Input: "path: src/main.go", Output: "package main\n\nfunc main() {\n    app.Run()\n}"},
				},
			},
			{
				Name: "Database",
				Tools: []mcpTool{
					{Name: "query", Desc: "Run a SQL query", Input: "SELECT count(*) FROM users", Output: "┌───────┐\n│ count │\n├───────┤\n│ 1,247 │\n└───────┘"},
					{Name: "list_tables", Desc: "List database tables", Input: "(no input needed)", Output: "users\norders\nproducts\nsessions"},
					{Name: "describe", Desc: "Describe a table schema", Input: "table: users", Output: "id       SERIAL PRIMARY KEY\nemail    VARCHAR(255) UNIQUE\nname     VARCHAR(100)\ncreated  TIMESTAMP DEFAULT NOW()"},
				},
			},
			{
				Name: "Slack",
				Tools: []mcpTool{
					{Name: "send_message", Desc: "Send a message to a channel", Input: "channel: #dev\nmsg: Deploy complete ✓", Output: "Message sent to #dev at 14:32"},
					{Name: "search", Desc: "Search messages", Input: "query: \"deployment error\"", Output: "3 results found:\n- @alice: \"deployment error on staging\" (2h ago)\n- @bob: \"fixed deployment error\" (1h ago)\n- @carol: \"no more deployment errors\" (30m ago)"},
				},
			},
		},
	}
}

func (m *MCPCallerModel) Init() tea.Cmd { return nil }

func (m *MCPCallerModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.called {
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter", " ", "backspace"))) {
				m.called = false
			}
			return m, nil
		}

		if m.inToolView {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
				if m.toolIdx > 0 {
					m.toolIdx--
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
				srv := m.servers[m.serverIdx]
				if m.toolIdx < len(srv.Tools)-1 {
					m.toolIdx++
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
				tool := m.servers[m.serverIdx].Tools[m.toolIdx]
				m.callOutput = tool.Output
				m.called = true
			case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
				m.inToolView = false
				m.toolIdx = 0
			}
		} else {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
				if m.serverIdx > 0 {
					m.serverIdx--
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
				if m.serverIdx < len(m.servers)-1 {
					m.serverIdx++
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
				m.inToolView = true
				m.toolIdx = 0
			}
		}
	}
	return m, nil
}

func (m *MCPCallerModel) View() string {
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#818cf8")).Bold(true)
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80")).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("#38bdf8")).Bold(true)
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#facc15"))
	output := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80"))

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  MCP Tool Explorer"))
	lines = append(lines, dim.Render("  Browse servers and call their tools"))
	lines = append(lines, "")

	if m.called {
		// Show call result
		tool := m.servers[m.serverIdx].Tools[m.toolIdx]
		srv := m.servers[m.serverIdx]
		lines = append(lines, highlight.Render(fmt.Sprintf("  Calling %s.%s...", srv.Name, tool.Name)))
		lines = append(lines, "")
		lines = append(lines, yellow.Render("  Input:"))
		for _, line := range strings.Split(tool.Input, "\n") {
			lines = append(lines, dim.Render("    "+line))
		}
		lines = append(lines, "")
		lines = append(lines, yellow.Render("  Output:"))
		for _, line := range strings.Split(m.callOutput, "\n") {
			lines = append(lines, output.Render("    "+line))
		}
		lines = append(lines, "")
		lines = append(lines, active.Render("  ✓ Tool call complete"))
		lines = append(lines, "", dim.Render("  [Enter] Back to tools"))
	} else if m.inToolView {
		// Show tools for current server
		srv := m.servers[m.serverIdx]
		lines = append(lines, highlight.Render("  "+srv.Name+" Server Tools"))
		lines = append(lines, "")
		for i, tool := range srv.Tools {
			style := text
			prefix := "  "
			if i == m.toolIdx {
				style = highlight
				prefix = "▸ "
			}
			lines = append(lines, fmt.Sprintf("  %s%s  %s", prefix, style.Render("ƒ "+tool.Name), dim.Render(tool.Desc)))
		}
		lines = append(lines, "")
		lines = append(lines, dim.Render("  [↑/↓] Navigate  [Enter] Call tool  [Bksp] Back to servers"))
	} else {
		// Server selection
		lines = append(lines, highlight.Render("  Select an MCP Server"))
		lines = append(lines, "")

		// Client diagram
		lines = append(lines, accent.Render("       ┌──────────┐"))
		lines = append(lines, accent.Render("       │  Claude   │"))
		lines = append(lines, accent.Render("       │ (client)  │"))
		lines = append(lines, accent.Render("       └────┬─────┘"))
		lines = append(lines, accent.Render("            │ MCP"))
		lines = append(lines, "")

		for i, srv := range m.servers {
			style := dim
			prefix := "    "
			if i == m.serverIdx {
				style = highlight
				prefix = "  ▸ "
			}
			toolCount := fmt.Sprintf("(%d tools)", len(srv.Tools))
			lines = append(lines, fmt.Sprintf("%s%s %s", prefix, style.Render("◆ "+srv.Name+" Server"), dim.Render(toolCount)))
		}

		lines = append(lines, "", dim.Render("  [↑/↓] Navigate  [Enter] Browse tools"))
	}

	return strings.Join(lines, "\n")
}
