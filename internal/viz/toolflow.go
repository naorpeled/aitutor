package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type toolStep struct {
	Name   string
	Desc   string
	Input  string
	Output string
}

// ToolFlowModel visualizes the Glob→Read→Edit tool chain.
type ToolFlowModel struct {
	width   int
	height  int
	step    int
	steps   []toolStep
}

func NewToolFlowModel(w, h int) Model {
	return &ToolFlowModel{
		width:  w,
		height: h,
		steps: []toolStep{
			{
				Name:   "Glob",
				Desc:   "Find files by pattern",
				Input:  "\"src/**/*.go\"",
				Output: "src/main.go\nsrc/app/app.go\nsrc/ui/styles.go",
			},
			{
				Name:   "Read",
				Desc:   "Read file contents",
				Input:  "src/app/app.go",
				Output: "func (m AppModel) Update...\n  // TODO: fix bug here\n  return m, nil",
			},
			{
				Name:   "Edit",
				Desc:   "Apply precise changes",
				Input:  "old: // TODO: fix bug\nnew: validated := check(m)",
				Output: "✓ File updated successfully",
			},
		},
	}
}

func (m *ToolFlowModel) Init() tea.Cmd { return nil }

func (m *ToolFlowModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.step < len(m.steps) {
				m.step++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.step = 0
		}
	}
	return m, nil
}

func (m *ToolFlowModel) View() string {
	active := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight)
	labelStyle := lipgloss.NewStyle().Foreground(ui.ColorIntermediate).Bold(true)

	var lines []string
	lines = append(lines, "")

	// Tool chain header
	var chain []string
	for i, s := range m.steps {
		if i < m.step {
			chain = append(chain, active.Render("✓ "+s.Name))
		} else if i == m.step {
			chain = append(chain, highlight.Render("▸ "+s.Name))
		} else {
			chain = append(chain, dim.Render("○ "+s.Name))
		}
		if i < len(m.steps)-1 {
			chain = append(chain, dim.Render(" → "))
		}
	}
	lines = append(lines, "  "+strings.Join(chain, ""))
	lines = append(lines, "")

	// Show details for current/completed steps
	for i := 0; i < m.step && i < len(m.steps); i++ {
		s := m.steps[i]
		lines = append(lines, fmt.Sprintf("  %s %s",
			active.Render(s.Name+":"),
			dim.Render(s.Desc)))
		lines = append(lines, fmt.Sprintf("    %s %s",
			labelStyle.Render("in:"),
			highlight.Render(s.Input)))
		for _, outLine := range strings.Split(s.Output, "\n") {
			lines = append(lines, fmt.Sprintf("    %s %s",
				labelStyle.Render("out:"),
				active.Render(outLine)))
		}
		lines = append(lines, "")
	}

	// Current step prompt
	if m.step < len(m.steps) {
		lines = append(lines, fmt.Sprintf("  %s",
			highlight.Render("Press Enter/Space to execute "+m.steps[m.step].Name)))
	} else {
		lines = append(lines, "  "+active.Render("✓ Tool chain complete!"))
	}

	lines = append(lines, "", "  [Enter/Space] Next step  [r] Reset")

	return strings.Join(lines, "\n")
}
