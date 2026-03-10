package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/pkg/types"
)

// HeaderModel renders the top bar with tier, lesson title, and progress.
type HeaderModel struct {
	Width       int
	Tier        types.Tier
	LessonTitle string
	Current     int
	Total       int
}

func NewHeaderModel() HeaderModel {
	return HeaderModel{Total: 15}
}

func (h HeaderModel) View() string {
	return h.ViewWithProgress("")
}

func (h HeaderModel) ViewWithProgress(progressBar string) string {
	tierColor := TierColor(int(h.Tier))
	tierBadge := lipgloss.NewStyle().
		Bold(true).
		Foreground(tierColor).
		Render(fmt.Sprintf(" %s ", h.Tier))

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBright).
		Render(h.LessonTitle)

	left := fmt.Sprintf("%s  %s", tierBadge, title)

	right := progressBar
	if right == "" && h.Total > 0 {
		right = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Render(fmt.Sprintf("%d/%d", h.Current, h.Total))
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := h.Width - leftWidth - rightWidth
	if gap < 0 {
		gap = 0
	}
	padding := lipgloss.NewStyle().Width(gap).Render("")

	row := left + padding + right

	return HeaderStyle.Width(h.Width).Render(row)
}
