package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type scenario struct {
	Description string
	CorrectMode string // "plan" or "execution"
	Explanation string
}

// ModePickerModel presents scenarios and users pick plan vs execution mode.
type ModePickerModel struct {
	width     int
	height    int
	scenarios []scenario
	current   int
	choice    string
	answered  bool
	correct   bool
	score     int
}

func NewModePickerModel(w, h int) Model {
	return &ModePickerModel{
		width:  w,
		height: h,
		scenarios: []scenario{
			{
				Description: "\"How should we restructure the authentication system to support OAuth2?\"",
				CorrectMode: "plan",
				Explanation: "Architecture decisions need analysis first — plan mode explores without changing anything.",
			},
			{
				Description: "\"Fix the NPE in UserService.getById when user is null\"",
				CorrectMode: "execution",
				Explanation: "Clear, focused bug fix — execution mode can read, edit, and test the fix.",
			},
			{
				Description: "\"What's the best approach for adding real-time notifications?\"",
				CorrectMode: "plan",
				Explanation: "Open-ended design question — needs exploration and comparison of approaches.",
			},
			{
				Description: "\"Write unit tests for the new PaymentProcessor class\"",
				CorrectMode: "execution",
				Explanation: "Well-defined task — execution mode can write tests and run them immediately.",
			},
			{
				Description: "\"Review this PR's architecture and identify potential scaling issues\"",
				CorrectMode: "plan",
				Explanation: "Code review and analysis — plan mode reads code without modifying it.",
			},
			{
				Description: "\"Rename the 'data' variable to 'userRecords' across the codebase\"",
				CorrectMode: "execution",
				Explanation: "Mechanical refactoring — execution mode can find and replace efficiently.",
			},
		},
	}
}

func (m *ModePickerModel) Init() tea.Cmd { return nil }

func (m *ModePickerModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.answered {
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))) {
				m.current++
				m.answered = false
				m.choice = ""
				if m.current >= len(m.scenarios) {
					m.current = len(m.scenarios)
				}
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("p", "1"))):
			m.choice = "plan"
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("e", "2"))):
			m.choice = "execution"
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.current = 0
			m.answered = false
			m.choice = ""
			m.score = 0
		}
	}
	return m, nil
}

func (m *ModePickerModel) submit() {
	s := m.scenarios[m.current]
	m.answered = true
	m.correct = m.choice == s.CorrectMode
	if m.correct {
		m.score++
	}
}

func (m *ModePickerModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	good := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	bad := lipgloss.NewStyle().Foreground(ui.ColorIncorrect).Bold(true)
	plan := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	exec := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	explain := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Mode Picker Exercise"))
	lines = append(lines, dim.Render("  For each scenario, choose: Plan mode or Execution mode?"))
	lines = append(lines, "")

	if m.current >= len(m.scenarios) {
		lines = append(lines, good.Render(fmt.Sprintf("  Exercise Complete! Score: %d/%d", m.score, len(m.scenarios))))
		if m.score == len(m.scenarios) {
			lines = append(lines, good.Render("  Perfect! You can identify the right mode for any task."))
		}
		lines = append(lines, "", dim.Render("  [r] Try again"))
		return strings.Join(lines, "\n")
	}

	s := m.scenarios[m.current]
	lines = append(lines, dim.Render(fmt.Sprintf("  Scenario %d of %d", m.current+1, len(m.scenarios))))
	lines = append(lines, "")
	lines = append(lines, text.Render("  "+s.Description))
	lines = append(lines, "")

	if m.answered {
		if m.correct {
			lines = append(lines, good.Render("  ✓ Correct!"))
		} else {
			correctLabel := plan.Render("Plan Mode")
			if s.CorrectMode == "execution" {
				correctLabel = exec.Render("Execution Mode")
			}
			lines = append(lines, bad.Render("  ✗ Not quite — the answer is ")+correctLabel)
		}
		lines = append(lines, explain.Render("  "+s.Explanation))
		lines = append(lines, "", highlight.Render("  Press Enter to continue"))
	} else {
		lines = append(lines, plan.Render("  [p/1] 📋 Plan Mode")+"      "+exec.Render("[e/2] ⚡ Execution Mode"))
	}

	lines = append(lines, "", dim.Render("  [p/1] Plan  [e/2] Execution  [r] Restart"))

	return strings.Join(lines, "\n")
}
