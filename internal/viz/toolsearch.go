package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type deferredTool struct {
	Name   string
	Desc   string
	Loaded bool
}

// ToolSearchModel simulates the ToolSearch discovery flow.
type ToolSearchModel struct {
	width     int
	height    int
	tools     []deferredTool
	cursor    int
	phase     int // 0=browse deferred list, 1=search, 2=result
	searchResults []int
	contextUsed int
}

func NewToolSearchModel(w, h int) Model {
	return &ToolSearchModel{
		width:  w,
		height: h,
		contextUsed: 10000, // base context usage
		tools: []deferredTool{
			{Name: "mcp__slack__send_message", Desc: "Send a Slack message"},
			{Name: "mcp__slack__read_channel", Desc: "Read Slack channel messages"},
			{Name: "mcp__github__create_pr", Desc: "Create a GitHub pull request"},
			{Name: "mcp__github__list_issues", Desc: "List GitHub issues"},
			{Name: "mcp__database__query", Desc: "Run a database query"},
			{Name: "NotebookEdit", Desc: "Edit Jupyter notebook cells"},
			{Name: "WebSearch", Desc: "Search the web"},
			{Name: "WebFetch", Desc: "Fetch a web page"},
		},
	}
}

func (m *ToolSearchModel) Init() tea.Cmd { return nil }

func (m *ToolSearchModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.tools)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			// Load the selected tool
			if !m.tools[m.cursor].Loaded {
				m.tools[m.cursor].Loaded = true
				m.contextUsed += 200 // each tool def ~200 tokens
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
			// Simulate keyword search: find all slack tools
			m.searchResults = nil
			for i, t := range m.tools {
				if strings.Contains(strings.ToLower(t.Name), "slack") {
					m.searchResults = append(m.searchResults, i)
					if !m.tools[i].Loaded {
						m.tools[i].Loaded = true
						m.contextUsed += 200
					}
				}
			}
			m.phase = 1
		case key.Matches(msg, key.NewBinding(key.WithKeys("g"))):
			// Simulate keyword search: find all github tools
			m.searchResults = nil
			for i, t := range m.tools {
				if strings.Contains(strings.ToLower(t.Name), "github") {
					m.searchResults = append(m.searchResults, i)
					if !m.tools[i].Loaded {
						m.tools[i].Loaded = true
						m.contextUsed += 200
					}
				}
			}
			m.phase = 1
		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
			m.phase = 0
			m.searchResults = nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			for i := range m.tools {
				m.tools[i].Loaded = false
			}
			m.contextUsed = 10000
			m.phase = 0
			m.searchResults = nil
		}
	}
	return m, nil
}

func (m *ToolSearchModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	loaded := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	deferred := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  ToolSearch Simulator"))
	lines = append(lines, dim.Render("  Load deferred tools on demand"))
	lines = append(lines, "")

	// Context usage bar
	barWidth := 30
	filledRatio := float64(m.contextUsed) / 200000
	filled := int(filledRatio * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	bar := loaded.Render(strings.Repeat("█", filled)) + dim.Render(strings.Repeat("░", barWidth-filled))
	lines = append(lines, fmt.Sprintf("  Context: %s %dk/200k", bar, m.contextUsed/1000))

	loadedCount := 0
	for _, t := range m.tools {
		if t.Loaded {
			loadedCount++
		}
	}
	lines = append(lines, dim.Render(fmt.Sprintf("  %d/%d tools loaded (+%d tokens)", loadedCount, len(m.tools), loadedCount*200)))
	lines = append(lines, "")

	if m.phase == 1 && len(m.searchResults) > 0 {
		lines = append(lines, yellow.Render("  Search results (auto-loaded):"))
		for _, idx := range m.searchResults {
			t := m.tools[idx]
			lines = append(lines, loaded.Render(fmt.Sprintf("    ✓ %s — %s", t.Name, t.Desc)))
		}
		lines = append(lines, "")
		lines = append(lines, dim.Render("  [Bksp] Back to tool list"))
	}

	// Tool list
	lines = append(lines, highlight.Render("  Available Deferred Tools:"))
	for i, t := range m.tools {
		style := deferred
		icon := "○"
		label := "deferred"
		if t.Loaded {
			style = loaded
			icon = "✓"
			label = "loaded"
		}
		prefix := "  "
		if i == m.cursor {
			prefix = "▸ "
			if !t.Loaded {
				style = highlight
			}
		}
		lines = append(lines, fmt.Sprintf("  %s%s %-30s %s",
			prefix, style.Render(icon), style.Render(t.Name), dim.Render("["+label+"]")))
	}

	lines = append(lines, "")
	lines = append(lines, dim.Render("  [↑/↓] Navigate  [Enter] Load tool  [s] Search 'slack'  [g] Search 'github'  [r] Reset"))

	return strings.Join(lines, "\n")
}
