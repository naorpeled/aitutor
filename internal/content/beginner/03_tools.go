package beginner

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      3,
		Title:   "Tools",
		Tier:    types.Beginner,
		Summary: "How AI assistants interact with your codebase",
		VizBuilder: func(w, h int) interface{} { return viz.NewToolFlowModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Tools: The AI's Hands"},
			{Kind: types.Paragraph, Content: "AI models can only think and generate text. To actually do things — read files, search code, run commands — they need tools. Tools are functions the AI can call to interact with the real world."},
			{Kind: types.Heading, Content: "Common Tool Categories"},
			{Kind: types.Bullet, Content: "File operations — Read, Write, Edit files\nSearch — Glob (find files by pattern), Grep (search content)\nExecution — Bash (run shell commands)\nNavigation — LSP (go to definition, find references)"},
			{Kind: types.Heading, Content: "The Tool Call Flow"},
			{Kind: types.Code, Content: "  1. AI decides it needs information\n  2. AI calls a tool (e.g., Read file.go)\n  3. Tool executes and returns results\n  4. AI processes the results\n  5. AI continues reasoning or calls another tool"},
			{Kind: types.Heading, Content: "Dedicated vs General Tools"},
			{Kind: types.Paragraph, Content: "Dedicated tools like Read, Edit, and Grep are preferred over general-purpose tools like Bash. They provide better safety, clearer intent, and easier review of what the AI is doing."},
			{Kind: types.Code, Content: "  ✓ Read(\"src/main.go\")        — clear, safe, reviewable\n  ✗ Bash(\"cat src/main.go\")    — opaque, harder to review"},
			{Kind: types.Heading, Content: "Permission Model"},
			{Kind: types.Paragraph, Content: "Tools operate under a permission system. Some tools run automatically (like reading files), while others require explicit approval (like running arbitrary shell commands or writing files). This keeps you in control."},
			{Kind: types.Callout, Content: "Think of tools as the AI's hands and eyes. Without them, it can only think. With them, it can explore, modify, and build."},
			{Kind: types.Callout, Content: "Learn more: Language Server Protocol — https://microsoft.github.io/language-server-protocol/ | Glob patterns — https://en.wikipedia.org/wiki/Glob_(programming)"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Why are dedicated tools (Read, Edit, Grep) preferred over Bash?",
				Choices:    []string{"They're faster", "Better safety, clearer intent, easier review", "They use less memory", "They're newer"},
				CorrectIdx: 1,
				Explanation: "Dedicated tools provide better safety, clearer intent, and make it easier to review what the AI is doing.",
			},
			{
				Kind:       types.Ordering,
				Prompt:     "Put the tool call flow in the correct order:",
				Choices:    []string{"AI decides it needs information", "AI calls a tool", "Tool executes and returns results", "AI processes results"},
				Explanation: "The flow is: decide → call → execute → process results.",
			},
		},
	})
}
