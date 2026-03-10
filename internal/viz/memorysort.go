package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type memoryItem struct {
	Text       string
	ShouldSave bool
	UserChoice string // "save", "skip", or ""
}

// MemorySortModel lets users categorize items as save vs don't-save.
type MemorySortModel struct {
	items     []memoryItem
	current   int
	submitted bool
	score     int
}

func NewMemorySortModel(w, h int) Model {
	return &MemorySortModel{
		items: []memoryItem{
			{Text: "Always use bun instead of npm", ShouldSave: true},
			{Text: "The current task is to fix bug #423", ShouldSave: false},
			{Text: "Project uses PostgreSQL with pgx driver", ShouldSave: true},
			{Text: "I'm currently on the feature/auth branch", ShouldSave: false},
			{Text: "Run 'make lint' before committing", ShouldSave: true},
			{Text: "The user just asked about API design", ShouldSave: false},
			{Text: "Error messages should include request IDs", ShouldSave: true},
			{Text: "I found a potential fix in server.go line 42", ShouldSave: false},
			{Text: "The team prefers table-driven tests in Go", ShouldSave: true},
			{Text: "Working on implementing the login page now", ShouldSave: false},
		},
	}
}

func (m *MemorySortModel) Init() tea.Cmd { return nil }

func (m *MemorySortModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.submitted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
				m.current = 0
				m.submitted = false
				m.score = 0
				for i := range m.items {
					m.items[i].UserChoice = ""
				}
			}
		}
		return m, nil
	}

	if m.current >= len(m.items) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("s", "1"))):
			m.items[m.current].UserChoice = "save"
			if m.items[m.current].ShouldSave {
				m.score++
			}
			m.current++
			if m.current >= len(m.items) {
				m.submitted = true
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("d", "2"))):
			m.items[m.current].UserChoice = "skip"
			if !m.items[m.current].ShouldSave {
				m.score++
			}
			m.current++
			if m.current >= len(m.items) {
				m.submitted = true
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.current = 0
			m.submitted = false
			m.score = 0
			for i := range m.items {
				m.items[i].UserChoice = ""
			}
		}
	}
	return m, nil
}

func (m *MemorySortModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	good := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	bad := lipgloss.NewStyle().Foreground(ui.ColorIncorrect).Bold(true)
	save := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	skip := lipgloss.NewStyle().Foreground(ui.ColorIntermediate).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Memory Sort Exercise"))
	lines = append(lines, dim.Render("  Should the AI save this to memory? Press [s]ave or [d]on't save"))
	lines = append(lines, "")

	if m.submitted {
		lines = append(lines, good.Render(fmt.Sprintf("  Score: %d/%d", m.score, len(m.items))))
		lines = append(lines, "")

		for _, item := range m.items {
			correctAction := "SAVE"
			if !item.ShouldSave {
				correctAction = "SKIP"
			}

			wasCorrect := (item.UserChoice == "save" && item.ShouldSave) ||
				(item.UserChoice == "skip" && !item.ShouldSave)

			style := good
			icon := "✓"
			if !wasCorrect {
				style = bad
				icon = "✗"
			}

			lines = append(lines, fmt.Sprintf("  %s %s %s",
				style.Render(icon),
				text.Render(item.Text),
				dim.Render("["+correctAction+"]")))
		}

		lines = append(lines, "")
		lines = append(lines, dim.Render("  Save: stable patterns, conventions, user preferences"))
		lines = append(lines, dim.Render("  Skip: session-specific details, in-progress work"))
		lines = append(lines, "", dim.Render("  [r] Try again"))
	} else {
		// Show already sorted items
		for i := 0; i < m.current; i++ {
			item := m.items[i]
			if item.UserChoice == "save" {
				lines = append(lines, save.Render("  ✓ SAVE  ")+dim.Render(item.Text))
			} else {
				lines = append(lines, skip.Render("  ✗ SKIP  ")+dim.Render(item.Text))
			}
		}

		if m.current < len(m.items) {
			lines = append(lines, "")
			lines = append(lines, dim.Render(fmt.Sprintf("  Item %d of %d:", m.current+1, len(m.items))))
			lines = append(lines, "")
			lines = append(lines, highlight.Render("  \""+m.items[m.current].Text+"\""))
			lines = append(lines, "")
			lines = append(lines, save.Render("  [s/1] 💾 Save to memory")+"    "+skip.Render("[d/2] 🗑  Don't save"))
		}

		lines = append(lines, "", dim.Render("  [s/1] Save  [d/2] Don't save  [r] Restart"))
	}

	return strings.Join(lines, "\n")
}
