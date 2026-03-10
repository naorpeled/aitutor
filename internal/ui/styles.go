package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorBeginner     = lipgloss.Color("#4ade80") // green
	ColorIntermediate = lipgloss.Color("#facc15") // yellow
	ColorAdvanced     = lipgloss.Color("#f87171") // red
	ColorAccent       = lipgloss.Color("#818cf8") // indigo
	ColorMuted        = lipgloss.Color("#6b7280") // gray
	ColorBright       = lipgloss.Color("#f9fafb") // near-white
	ColorDim          = lipgloss.Color("#374151") // dark gray
	ColorCorrect      = lipgloss.Color("#4ade80") // green
	ColorIncorrect    = lipgloss.Color("#f87171") // red
	ColorHighlight    = lipgloss.Color("#38bdf8") // sky blue

	// Header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright).
			Background(lipgloss.Color("#1e1b4b")).
			Padding(0, 1)

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Background(lipgloss.Color("#111827")).
			Padding(0, 1)

	FooterKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight)

	FooterDescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(0, 1)

	SidebarActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAccent)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	SidebarCompletedStyle = lipgloss.NewStyle().
				Foreground(ColorBeginner)

	// Content area
	ContentStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Theory blocks
	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright).
			MarginBottom(1)

	ParagraphStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d1d5db"))

	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a5f3fc")).
			Background(lipgloss.Color("#1e293b")).
			Padding(0, 1)

	CalloutStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Foreground(lipgloss.Color("#c4b5fd")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	BulletStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d1d5db")).
			PaddingLeft(2)

	// Tier badges
	TierBadge = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Foreground(ColorBeginner).Bold(true),
		1: lipgloss.NewStyle().Foreground(ColorIntermediate).Bold(true),
		2: lipgloss.NewStyle().Foreground(ColorAdvanced).Bold(true),
	}
)

// TierColor returns the color for a given tier.
func TierColor(tier int) lipgloss.Color {
	switch tier {
	case 0:
		return ColorBeginner
	case 1:
		return ColorIntermediate
	case 2:
		return ColorAdvanced
	default:
		return ColorMuted
	}
}
