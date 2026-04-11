package main

import (
	"strings"
	"testing"

	"github.com/naorpeled/aitutor/internal/app"

	// Register all lessons via init()
	_ "github.com/naorpeled/aitutor/internal/content/advanced"
	_ "github.com/naorpeled/aitutor/internal/content/beginner"
	_ "github.com/naorpeled/aitutor/internal/content/intermediate"
)

func TestSmokeUI(t *testing.T) {
	m := app.NewAppModel("test")
	m = m.SmokeInit()
	view := m.View()

	// Check for UI elements that are always present regardless of which lesson is active
	for _, want := range []string{"Theory", "quit", "next"} {
		if !strings.Contains(view, want) {
			t.Errorf("view missing %q\n\nview output:\n%s", want, view)
		}
	}

	// At least one tier name must appear in the header
	hasTier := strings.Contains(view, "Beginner") ||
		strings.Contains(view, "Intermediate") ||
		strings.Contains(view, "Advanced")
	if !hasTier {
		t.Errorf("view missing tier name (Beginner/Intermediate/Advanced)\n\nview output:\n%s", view)
	}
}
