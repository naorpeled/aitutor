package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the application.
type KeyMap struct {
	Quit      key.Binding
	Tab       key.Binding
	Next      key.Binding
	Prev      key.Binding
	Advance      key.Binding
	AdvancePhase key.Binding
	Back         key.Binding
	Help      key.Binding
	Up        key.Binding
	Down      key.Binding
	Select    key.Binding
	Space     key.Binding
	Reset     key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "toggle sidebar"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next lesson"),
	),
	Prev: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous lesson"),
	),
	Advance: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "advance phase"),
	),
	AdvancePhase: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next phase"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "left"),
		key.WithHelp("←/Bksp", "previous phase"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "select"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("Space", "interact"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
}
