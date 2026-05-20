package advanced

import (
	"strings"
	"testing"

	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/pkg/types"
)

func TestAgentsBuildingLessonCoversIssueTopics(t *testing.T) {
	var found *types.LessonDef
	for _, def := range lesson.All() {
		if def.ID == 18 {
			found = &def
			break
		}
	}

	if found == nil {
		t.Fatal("lesson 18 for agents building topics is not registered")
	}
	if found.Tier != types.Advanced {
		t.Fatalf("lesson 18 tier = %s, want Advanced", found.Tier)
	}
	if found.VizBuilder == nil {
		t.Fatal("lesson 18 must provide an interactive visualization")
	}
	if len(found.Questions) < 2 {
		t.Fatalf("lesson 18 has %d quiz questions, want at least 2", len(found.Questions))
	}

	var theory strings.Builder
	for _, block := range found.Theory {
		theory.WriteString(block.Content)
		theory.WriteString("\n")
	}
	content := strings.ToLower(theory.String())
	for _, topic := range []string{"compaction", "long-term memory"} {
		if !strings.Contains(content, topic) {
			t.Fatalf("lesson 18 theory does not cover %q", topic)
		}
	}
}
