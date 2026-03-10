package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type subagent struct {
	Name     string
	Task     string
	Status   string
	Progress int
}

// FanoutModel visualizes subagent parallel fan-out.
type FanoutModel struct {
	width  int
	height int
	agents []subagent
}

func NewFanoutModel(w, h int) Model {
	return &FanoutModel{
		width:  w,
		height: h,
		agents: []subagent{
			{Name: "Agent 1", Task: "Frontend components", Status: "pending", Progress: 0},
			{Name: "Agent 2", Task: "Backend API", Status: "pending", Progress: 0},
			{Name: "Agent 3", Task: "Test suite", Status: "pending", Progress: 0},
		},
	}
}

func (m *FanoutModel) Init() tea.Cmd { return nil }

func (m *FanoutModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			allDone := true
			for i := range m.agents {
				if m.agents[i].Progress < 100 {
					allDone = false
					m.agents[i].Status = "running"
					// Different agents progress at different rates
					increments := []int{25, 20, 30}
					m.agents[i].Progress += increments[i%3]
					if m.agents[i].Progress >= 100 {
						m.agents[i].Progress = 100
						m.agents[i].Status = "done"
					}
				}
			}
			if allDone {
				// Reset
				for i := range m.agents {
					m.agents[i].Progress = 0
					m.agents[i].Status = "pending"
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			for i := range m.agents {
				m.agents[i].Progress = 0
				m.agents[i].Status = "pending"
			}
		}
	}
	return m, nil
}

func (m *FanoutModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	green := lipgloss.NewStyle().Foreground(ui.ColorBeginner)
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	blue := lipgloss.NewStyle().Foreground(ui.ColorHighlight)

	var lines []string
	lines = append(lines, "")

	// Main agent
	lines = append(lines, accent.Render("            ┌───────────┐"))
	lines = append(lines, accent.Render("            │   Main    │"))
	lines = append(lines, accent.Render("            │   Agent   │"))
	lines = append(lines, accent.Render("            └─────┬─────┘"))
	lines = append(lines, accent.Render("           ┌──────┼──────┐"))
	lines = append(lines, accent.Render("           ▼      ▼      ▼"))
	lines = append(lines, "")

	// Subagents
	barWidth := 15
	for _, a := range m.agents {
		var statusIcon string
		var style lipgloss.Style
		switch a.Status {
		case "pending":
			statusIcon = "○"
			style = dim
		case "running":
			statusIcon = "◉"
			style = yellow
		case "done":
			statusIcon = "✓"
			style = green
		}

		filled := int(float64(a.Progress) / 100 * float64(barWidth))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

		lines = append(lines, fmt.Sprintf("  %s %s %s",
			style.Render(statusIcon),
			blue.Render(fmt.Sprintf("%-10s", a.Name)),
			style.Render(a.Task)))
		lines = append(lines, fmt.Sprintf("    %s %s",
			style.Render(bar),
			dim.Render(fmt.Sprintf("%3d%%", a.Progress))))
		lines = append(lines, "")
	}

	// Check if all done
	allDone := true
	for _, a := range m.agents {
		if a.Status != "done" {
			allDone = false
			break
		}
	}

	if allDone {
		lines = append(lines, green.Render("  ✓ All agents completed! Results combined."))
	} else {
		lines = append(lines, blue.Render("  Press Enter/Space to advance agents"))
	}

	lines = append(lines, "", dim.Render("  [Enter/Space] Advance  [r] Reset"))

	return strings.Join(lines, "\n")
}
