package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type agentAction struct {
	Phase       string
	Description string
	Detail      string
}

// AgentLoopModel simulates walking through an agent loop interactively.
type AgentLoopModel struct {
	width   int
	height  int
	step    int
	actions []agentAction
}

func NewAgentLoopModel(w, h int) Model {
	return &AgentLoopModel{
		width:  w,
		height: h,
		actions: []agentAction{
			{Phase: "User Request", Description: "User asks: \"Add input validation to the signup handler\"", Detail: "The agent receives your task and begins working."},
			{Phase: "Read Context", Description: "Agent reads src/handlers/signup.go", Detail: "Gathering context — the agent needs to understand the current code."},
			{Phase: "Reason & Plan", Description: "Agent identifies: no email format check, no password length check", Detail: "Analyzing what's missing and planning the changes."},
			{Phase: "Take Action", Description: "Agent calls Edit to add email regex validation", Detail: "Tool call: Edit(signup.go, add validateEmail function)"},
			{Phase: "Observe Result", Description: "File updated successfully ✓", Detail: "The edit succeeded. But the task isn't done yet..."},
			// Second loop iteration
			{Phase: "Read Context", Description: "Agent re-reads the updated file", Detail: "Loop iteration 2 — checking current state after changes."},
			{Phase: "Reason & Plan", Description: "Agent sees: email validated, but password still unchecked", Detail: "More work needed — password validation is still missing."},
			{Phase: "Take Action", Description: "Agent calls Edit to add password length validation", Detail: "Tool call: Edit(signup.go, add validatePassword function)"},
			{Phase: "Observe Result", Description: "File updated successfully ✓", Detail: "Both validations added. Let's verify..."},
			// Third iteration - verify
			{Phase: "Read Context", Description: "Agent reads final file state", Detail: "Loop iteration 3 — final verification pass."},
			{Phase: "Reason & Plan", Description: "Both validations present. Task complete!", Detail: "All requirements met — ready to respond."},
			{Phase: "Response", Description: "\"Added email format and password length validation to signup handler\"", Detail: "Agent reports what it did back to the user."},
		},
	}
}

func (m *AgentLoopModel) Init() tea.Cmd { return nil }

func (m *AgentLoopModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.step < len(m.actions)-1 {
				m.step++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.step = 0
		}
	}
	return m, nil
}

func (m *AgentLoopModel) View() string {
	active := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	current := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	detail := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))

	phases := []string{"User Request", "Read Context", "Reason & Plan", "Take Action", "Observe Result", "Response"}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Drive the Agent Loop"))
	lines = append(lines, dim.Render("  Press Enter/Space to advance each step"))
	lines = append(lines, "")

	// Show loop diagram with current phase highlighted
	currentPhase := m.actions[m.step].Phase
	for _, phase := range phases {
		icon := "○"
		style := dim
		if phase == currentPhase {
			icon = "▸"
			style = current
		}
		// Check if this phase was already visited
		for i := 0; i < m.step; i++ {
			if m.actions[i].Phase == phase {
				if phase != currentPhase {
					icon = "✓"
					style = active
				}
				break
			}
		}
		lines = append(lines, fmt.Sprintf("  %s %s", style.Render(icon), style.Render(phase)))
		if phase != "Response" {
			lines = append(lines, dim.Render("  │"))
		}
	}

	lines = append(lines, "")

	// Current action detail
	a := m.actions[m.step]
	lines = append(lines, current.Render("  ── Current Step ──"))
	lines = append(lines, desc.Render("  "+a.Description))
	lines = append(lines, detail.Render("  "+a.Detail))
	lines = append(lines, "")

	// Progress
	lines = append(lines, dim.Render(fmt.Sprintf("  Step %d of %d", m.step+1, len(m.actions))))

	if m.step < len(m.actions)-1 {
		lines = append(lines, "")
		lines = append(lines, current.Render("  Press Enter/Space to continue"))
	} else {
		lines = append(lines, "")
		lines = append(lines, active.Render("  ✓ Agent loop complete! The task is done."))
	}

	lines = append(lines, "", dim.Render("  [Enter/Space] Next step  [r] Restart"))

	return strings.Join(lines, "\n")
}
