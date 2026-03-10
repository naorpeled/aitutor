package viz

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

// LifecycleModel shows the hooks lifecycle timeline.
type LifecycleModel struct {
	width    int
	height   int
	step     int
	maxSteps int
}

func NewLifecycleModel(w, h int) Model {
	return &LifecycleModel{width: w, height: h, maxSteps: 5}
}

func (m *LifecycleModel) Init() tea.Cmd { return nil }

func (m *LifecycleModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.step < m.maxSteps {
				m.step++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.step = 0
		}
	}
	return m, nil
}

func (m *LifecycleModel) View() string {
	active := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	hook := lipgloss.NewStyle().Foreground(ui.ColorIntermediate).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight)
	blocked := lipgloss.NewStyle().Foreground(ui.ColorAdvanced).Bold(true)

	type event struct {
		label    string
		isHook   bool
		canBlock bool
	}

	events := []event{
		{label: "Session starts", isHook: false},
		{label: "SessionStart hook", isHook: true},
		{label: "User sends message", isHook: false},
		{label: "PromptSubmit hook", isHook: true, canBlock: true},
		{label: "AI decides to use Edit tool", isHook: false},
		{label: "PreToolUse hook", isHook: true, canBlock: true},
		{label: "Tool executes", isHook: false},
		{label: "PostToolUse hook", isHook: true},
		{label: "AI generates response", isHook: false},
		{label: "Notification hook", isHook: true},
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, highlight.Render("  Hooks Lifecycle Timeline"))
	lines = append(lines, "")

	shown := m.step * 2
	if shown > len(events) {
		shown = len(events)
	}

	for i := 0; i < len(events); i++ {
		e := events[i]
		style := dim
		connector := dim.Render("  │")

		if i < shown {
			if e.isHook {
				style = hook
				if e.canBlock {
					lines = append(lines, style.Render("  ◆ "+e.label)+blocked.Render(" [can block]"))
				} else {
					lines = append(lines, style.Render("  ◆ "+e.label))
				}
			} else {
				style = active
				lines = append(lines, style.Render("  ● "+e.label))
			}
		} else {
			lines = append(lines, dim.Render("  ○ "+e.label))
		}

		if i < len(events)-1 {
			lines = append(lines, connector)
		}
	}

	lines = append(lines, "")
	if m.step < m.maxSteps {
		lines = append(lines, highlight.Render("  Press Enter/Space to advance"))
	} else {
		lines = append(lines, active.Render("  ✓ Full lifecycle shown!"))
	}
	lines = append(lines, "", dim.Render("  [Enter/Space] Advance  [r] Reset"))

	return strings.Join(lines, "\n")
}
