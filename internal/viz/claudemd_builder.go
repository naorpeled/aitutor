package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type claudemdItem struct {
	Text     string
	Category string // "good" or "bad"
	Selected bool
}

// ClaudeMDBuilderModel lets users build a CLAUDE.md by selecting good items.
type ClaudeMDBuilderModel struct {
	width    int
	height   int
	items    []claudemdItem
	cursor   int
	submitted bool
}

func NewClaudeMDBuilderModel(w, h int) Model {
	return &ClaudeMDBuilderModel{
		width:  w,
		height: h,
		items: []claudemdItem{
			{Text: "Build: `make build` to compile", Category: "good"},
			{Text: "Test: `npm test` to run all tests", Category: "good"},
			{Text: "The sky is blue", Category: "bad"},
			{Text: "Use snake_case for database columns", Category: "good"},
			{Text: "My favorite color is green", Category: "bad"},
			{Text: "Never commit .env files", Category: "good"},
			{Text: "API handlers live in internal/api/", Category: "good"},
			{Text: "Today is a good day", Category: "bad"},
			{Text: "Run `go vet` before committing", Category: "good"},
			{Text: "I like pizza", Category: "bad"},
		},
	}
}

func (m *ClaudeMDBuilderModel) Init() tea.Cmd { return nil }

func (m *ClaudeMDBuilderModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.submitted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
				m.submitted = false
				m.cursor = 0
				for i := range m.items {
					m.items[i].Selected = false
				}
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys(" "))):
			m.items[m.cursor].Selected = !m.items[m.cursor].Selected
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			m.submitted = true
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.cursor = 0
			for i := range m.items {
				m.items[i].Selected = false
			}
		}
	}
	return m, nil
}

func (m *ClaudeMDBuilderModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	good := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	bad := lipgloss.NewStyle().Foreground(ui.ColorIncorrect).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Build a CLAUDE.md"))
	lines = append(lines, dim.Render("  Select items that belong in a CLAUDE.md file"))
	lines = append(lines, dim.Render("  Space to toggle, Enter to submit"))
	lines = append(lines, "")

	if m.submitted {
		// Score
		correct := 0
		total := 0
		for _, item := range m.items {
			if item.Category == "good" {
				total++
				if item.Selected {
					correct++
				}
			}
			if item.Category == "bad" && !item.Selected {
				correct++
				total++
			} else if item.Category == "bad" {
				total++
			}
		}

		lines = append(lines, yellow.Render(fmt.Sprintf("  Score: %d/%d correct selections", correct, total)))
		lines = append(lines, "")

		// Show results
		lines = append(lines, good.Render("  Your CLAUDE.md:"))
		lines = append(lines, dim.Render("  ┌────────────────────────────────────────┐"))
		for _, item := range m.items {
			if item.Selected {
				style := good
				mark := "✓"
				if item.Category == "bad" {
					style = bad
					mark = "✗ (doesn't belong)"
				}
				lines = append(lines, "  │ "+style.Render(mark+" "+item.Text))
			}
		}

		// Show missed good items
		for _, item := range m.items {
			if item.Category == "good" && !item.Selected {
				lines = append(lines, "  │ "+dim.Render("✗ MISSED: "+item.Text))
			}
		}
		lines = append(lines, dim.Render("  └────────────────────────────────────────┘"))

		lines = append(lines, "", dim.Render("  [r] Try again"))
	} else {
		for i, item := range m.items {
			checkbox := "[ ]"
			style := text
			if item.Selected {
				checkbox = "[✓]"
				style = highlight
			}
			prefix := "  "
			if i == m.cursor {
				prefix = "▸ "
				style = highlight
			}
			lines = append(lines, fmt.Sprintf("  %s%s %s", prefix, style.Render(checkbox), style.Render(item.Text)))
		}

		selected := 0
		for _, item := range m.items {
			if item.Selected {
				selected++
			}
		}
		lines = append(lines, "")
		lines = append(lines, dim.Render(fmt.Sprintf("  %d items selected", selected)))
		lines = append(lines, "", dim.Render("  [↑/↓] Navigate  [Space] Toggle  [Enter] Submit  [r] Reset"))
	}

	return strings.Join(lines, "\n")
}
