package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      7,
		Title:   "Hooks",
		Tier:    types.Intermediate,
		Summary: "Lifecycle hooks for AI assistant actions",
		VizBuilder: func(w, h int) interface{} { return viz.NewLifecycleModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Hooks"},
			{Kind: types.Paragraph, Content: "Hooks are user-defined shell commands that execute automatically in response to AI assistant events. They let you customize behavior, enforce policies, and integrate with your workflow."},
			{Kind: types.Heading, Content: "Hook Types"},
			{Kind: types.Bullet, Content: "PreToolUse — runs before a tool executes (can block it)\nPostToolUse — runs after a tool executes\nNotification — runs when the AI wants to notify you\nSessionStart — runs when a new conversation begins\nPromptSubmit — runs when you send a message"},
			{Kind: types.Heading, Content: "How Hooks Work"},
			{Kind: types.Code, Content: "  User sends message\n       │\n       ▼\n  ┌─PromptSubmit hook──┐\n  └────────────────────┘\n       │\n       ▼\n  AI decides to use tool\n       │\n       ▼\n  ┌─PreToolUse hook────┐  ← can BLOCK the tool\n  └────────────────────┘\n       │\n       ▼\n  Tool executes\n       │\n       ▼\n  ┌─PostToolUse hook───┐\n  └────────────────────┘"},
			{Kind: types.Heading, Content: "Example: Auto-format on Edit"},
			{Kind: types.Code, Content: "  // .claude/hooks.json\n  {\n    \"hooks\": {\n      \"PostToolUse\": [{\n        \"matcher\": \"Edit\",\n        \"command\": \"./scripts/format-on-save.sh\"\n      }]\n    }\n  }"},
			{Kind: types.Callout, Content: "Hooks are powerful for enforcing team standards — like running linters after every edit or blocking writes to protected files."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which hook type can BLOCK a tool from executing?",
				Choices:    []string{"PostToolUse", "Notification", "PreToolUse", "SessionStart"},
				CorrectIdx: 2,
				Explanation: "PreToolUse runs before a tool executes and can block it from running.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What is a practical use case for PostToolUse hooks?",
				Choices:    []string{"Blocking dangerous commands", "Auto-formatting files after edits", "Starting new sessions", "Sending notifications"},
				CorrectIdx: 1,
				Explanation: "PostToolUse hooks run after tool execution — perfect for auto-formatting code after an Edit.",
			},
		},
	})
}
