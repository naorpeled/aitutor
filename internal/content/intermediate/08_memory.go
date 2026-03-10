package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      8,
		Title:   "Memory & Persistence",
		Tier:    types.Intermediate,
		Summary: "How AI assistants remember across sessions",
		VizBuilder: func(w, h int) interface{} { return viz.NewMemorySortModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Memory & Persistence"},
			{Kind: types.Paragraph, Content: "By default, each AI conversation starts fresh — the model has no memory of previous sessions. But several mechanisms exist to persist knowledge across conversations."},
			{Kind: types.Heading, Content: "Memory Layers"},
			{Kind: types.Code, Content: "  ┌─────────────────────────────────────┐\n  │         Session Memory              │\n  │  (conversation context, ephemeral)  │\n  ├─────────────────────────────────────┤\n  │        Auto Memory Files            │\n  │  (tool-managed persistent notes)    │\n  ├─────────────────────────────────────┤\n  │      Project Config Files           │\n  │  (AGENTS.md, CLAUDE.md, etc.)       │\n  ├─────────────────────────────────────┤\n  │        User Settings                │\n  │  (tool-specific user preferences)   │\n  └─────────────────────────────────────┘"},
			{Kind: types.Heading, Content: "Auto Memory"},
			{Kind: types.Paragraph, Content: "Some AI coding tools support auto memory — the AI can save notes to persistent files that are loaded in future sessions. It stores patterns, conventions, debugging insights, and user preferences — things confirmed across multiple interactions."},
			{Kind: types.Heading, Content: "What to Remember vs Not"},
			{Kind: types.Bullet, Content: "Remember: stable patterns, architecture decisions, user preferences, recurring solutions\nDon't remember: session-specific details, in-progress work, speculative conclusions"},
			{Kind: types.Callout, Content: "In tools that support persistent memory, you can explicitly tell the AI to remember something: \"Always use bun instead of npm\" — and it will persist this preference across sessions."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which memory layer is the MOST persistent?",
				Choices:    []string{"Session memory", "Auto memory files", "Project config files (AGENTS.md)", "Conversation history"},
				CorrectIdx: 2,
				Explanation: "Project config files like AGENTS.md are the most persistent — they're version-controlled and always loaded into every session.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What should you NOT save to auto memory?",
				Choices:    []string{"Stable coding patterns", "Architecture decisions", "Session-specific task details", "User preferences"},
				CorrectIdx: 2,
				Explanation: "Session-specific details (current task, in-progress work) shouldn't be saved — they're temporary context.",
			},
		},
	})
}
