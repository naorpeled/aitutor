package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type skill struct {
	Name   string
	Status string
}

// SkillLoadModel shows skill lazy-loading.
type SkillLoadModel struct {
	width  int
	height int
	skills []skill
	step   int
}

func NewSkillLoadModel(w, h int) Model {
	return &SkillLoadModel{
		width:  w,
		height: h,
		skills: []skill{
			{Name: "brainstorming", Status: "deferred"},
			{Name: "debugging", Status: "deferred"},
			{Name: "mcp-builder", Status: "deferred"},
			{Name: "frontend-design", Status: "deferred"},
			{Name: "test-driven-dev", Status: "deferred"},
		},
	}
}

func (m *SkillLoadModel) Init() tea.Cmd { return nil }

func (m *SkillLoadModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.step < len(m.skills) {
				m.skills[m.step].Status = "loaded"
				m.step++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.step = 0
			for i := range m.skills {
				m.skills[i].Status = "deferred"
			}
		}
	}
	return m, nil
}

func (m *SkillLoadModel) View() string {
	loaded := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	deferred := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Skill Registry"))
	lines = append(lines, "")

	// Context window representation
	loadedCount := 0
	for _, s := range m.skills {
		if s.Status == "loaded" {
			loadedCount++
		}
	}
	contextUsed := loadedCount * 3 // each skill ~3% of context
	contextBar := strings.Repeat("█", contextUsed) + strings.Repeat("░", 30-contextUsed)
	lines = append(lines, fmt.Sprintf("  Context: %s %d%%", highlight.Render(contextBar), contextUsed*100/30))
	lines = append(lines, "")

	for _, s := range m.skills {
		var icon, label string
		var style lipgloss.Style
		if s.Status == "loaded" {
			icon = "✓"
			label = "loaded"
			style = loaded
		} else {
			icon = "○"
			label = "deferred"
			style = deferred
		}
		lines = append(lines, fmt.Sprintf("  %s %-18s %s",
			style.Render(icon),
			highlight.Render(s.Name),
			style.Render("["+label+"]")))
	}

	lines = append(lines, "")
	if m.step < len(m.skills) {
		lines = append(lines, highlight.Render("  Press Enter/Space to load next skill"))
	} else {
		lines = append(lines, loaded.Render("  ✓ All skills loaded!"))
	}

	lines = append(lines, "")
	lines = append(lines, dim.Render("  Skills load on-demand to save context space"))
	lines = append(lines, "", dim.Render("  [Enter/Space] Load  [r] Reset"))

	return strings.Join(lines, "\n")
}
