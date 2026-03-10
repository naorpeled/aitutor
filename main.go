package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/naorpeled/aitutor/internal/app"

	// Register all lessons via init()
	_ "github.com/naorpeled/aitutor/internal/content/beginner"
	_ "github.com/naorpeled/aitutor/internal/content/intermediate"
	_ "github.com/naorpeled/aitutor/internal/content/advanced"
)

func main() {
	p := tea.NewProgram(
		app.NewAppModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
