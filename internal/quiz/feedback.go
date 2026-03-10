package quiz

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	correctStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ade80")).
			Bold(true)

	incorrectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f87171")).
			Bold(true)

	explanationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#818cf8")).
				PaddingLeft(2)
)

func RenderCorrect(explanation string) string {
	result := correctStyle.Render("  ✓ Correct!")
	if explanation != "" {
		result += "\n" + explanationStyle.Render(explanation)
	}
	return result
}

func RenderIncorrect(explanation string) string {
	result := incorrectStyle.Render("  ✗ Incorrect")
	if explanation != "" {
		result += "\n" + explanationStyle.Render(explanation)
	}
	return result
}
