package advanced

import (
	"strings"
	"testing"

	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/pkg/types"
)

func TestAgentsBuildingLessonCoversIssueTopics(t *testing.T) {
	defs := lesson.All()
	found := -1
	for i := range defs {
		if defs[i].ID == 18 {
			found = i
			break
		}
	}

	if found == -1 {
		t.Fatal("lesson 18 for agents building topics is not registered")
	}
	def := defs[found]
	if def.Tier != types.Advanced {
		t.Fatalf("lesson 18 tier = %v, want Advanced", def.Tier)
	}
	if def.VizBuilder == nil {
		t.Fatal("lesson 18 must provide an interactive visualization")
	}
	if len(def.Questions) < 2 {
		t.Fatalf("lesson 18 has %d quiz questions, want at least 2", len(def.Questions))
	}

	var theory strings.Builder
	for _, block := range def.Theory {
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
