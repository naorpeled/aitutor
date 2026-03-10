package viz

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the interface all visualizations implement.
type Model interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	View() string
}

// Box draws an ASCII box around content.
func Box(content string, width int) string {
	lines := strings.Split(content, "\n")
	innerWidth := width - 4
	if innerWidth < 1 {
		innerWidth = 1
	}

	top := "┌" + strings.Repeat("─", innerWidth+2) + "┐"
	bottom := "└" + strings.Repeat("─", innerWidth+2) + "┘"

	var result []string
	result = append(result, top)
	for _, line := range lines {
		padded := line
		if len(padded) < innerWidth {
			padded += strings.Repeat(" ", innerWidth-len(padded))
		} else if len(padded) > innerWidth {
			padded = padded[:innerWidth]
		}
		result = append(result, "│ "+padded+" │")
	}
	result = append(result, bottom)
	return strings.Join(result, "\n")
}

// Arrow returns a vertical arrow string.
func Arrow(length int) string {
	if length <= 0 {
		return "▼"
	}
	var lines []string
	for i := 0; i < length; i++ {
		lines = append(lines, "│")
	}
	lines = append(lines, "▼")
	return strings.Join(lines, "\n")
}

// HLine returns a horizontal line.
func HLine(width int) string {
	return strings.Repeat("─", width)
}

// CenterText centers text within a given width.
func CenterText(text string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(text)
}

// StaticModel is a simple viz that just displays static content.
type StaticModel struct {
	Content string
}

func NewStaticModel(content string) *StaticModel {
	return &StaticModel{Content: content}
}

func (m *StaticModel) Init() tea.Cmd         { return nil }
func (m *StaticModel) Update(tea.Msg) (Model, tea.Cmd) { return m, nil }
func (m *StaticModel) View() string           { return m.Content }
