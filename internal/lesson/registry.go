package lesson

import (
	"sort"

	"github.com/naorpeled/aitutor/pkg/types"
)

var registry []types.LessonDef

// Register adds a lesson to the global registry. Called from init() in content files.
func Register(def types.LessonDef) {
	registry = append(registry, def)
}

// All returns all registered lessons sorted by ID.
func All() []types.LessonDef {
	sorted := make([]types.LessonDef, len(registry))
	copy(sorted, registry)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

// Count returns how many lessons are registered.
func Count() int {
	return len(registry)
}
