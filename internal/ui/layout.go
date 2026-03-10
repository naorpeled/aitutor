package ui

// Layout holds the computed dimensions for the split-pane layout.
type Layout struct {
	Width        int
	Height       int
	HeaderHeight int
	FooterHeight int
	SidebarWidth int
	ContentWidth int
	ContentHeight int
	SidebarOpen  bool
}

const (
	DefaultHeaderHeight = 1
	DefaultFooterHeight = 1
	DefaultSidebarWidth = 28
	MinContentWidth     = 40
)

// ComputeLayout calculates the layout dimensions from terminal size.
func ComputeLayout(width, height int, sidebarOpen bool) Layout {
	l := Layout{
		Width:        width,
		Height:       height,
		HeaderHeight: DefaultHeaderHeight,
		FooterHeight: DefaultFooterHeight,
		SidebarOpen:  sidebarOpen,
	}

	l.ContentHeight = height - l.HeaderHeight - l.FooterHeight
	if l.ContentHeight < 1 {
		l.ContentHeight = 1
	}

	if sidebarOpen {
		l.SidebarWidth = DefaultSidebarWidth
		if width-l.SidebarWidth < MinContentWidth {
			l.SidebarWidth = 0
			l.SidebarOpen = false
		}
	}

	l.ContentWidth = width - l.SidebarWidth
	if l.ContentWidth < 1 {
		l.ContentWidth = 1
	}

	return l
}
