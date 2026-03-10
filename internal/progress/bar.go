package progress

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	filledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80"))
	emptyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
)

// Bar renders a progress bar.
func Bar(completed, total, width int) string {
	if total == 0 || width < 5 {
		return ""
	}

	barWidth := width - 8 // room for "XX/XX "
	if barWidth < 5 {
		barWidth = 5
	}

	filled := 0
	if total > 0 {
		filled = completed * barWidth / total
	}
	if filled > barWidth {
		filled = barWidth
	}

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", barWidth-filled))

	return fmt.Sprintf("%s %d/%d", bar, completed, total)
}
