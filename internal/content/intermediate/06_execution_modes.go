package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      6,
		Title:   "Execution Modes",
		Tier:    types.Intermediate,
		Summary: "Plan mode vs execution mode",
		VizBuilder: func(w, h int) interface{} { return viz.NewModePickerModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Execution Modes"},
			{Kind: types.Paragraph, Content: "AI coding assistants can operate in different modes depending on the task. The two primary modes are Plan Mode and Execution Mode, each optimized for different phases of work."},
			{Kind: types.Heading, Content: "Plan Mode"},
			{Kind: types.Paragraph, Content: "In plan mode, the AI focuses on analysis and planning without making changes. It reads files, explores the codebase, and produces a structured plan. This is ideal for complex tasks where you want to review the approach before implementation."},
			{Kind: types.Bullet, Content: "Read-only — no file modifications\nProduces structured plans with steps\nIdentifies files to change and potential risks\nGreat for architecture decisions and large refactors"},
			{Kind: types.Heading, Content: "Execution Mode"},
			{Kind: types.Paragraph, Content: "In execution mode, the AI actively makes changes — editing files, running commands, and iterating on results. This is the default mode for most tasks."},
			{Kind: types.Bullet, Content: "Full tool access — read, write, execute\nIterates based on results (test failures, build errors)\nBest for well-defined, focused tasks"},
			{Kind: types.Heading, Content: "When to Use Each"},
			{Kind: types.Code, Content: "  Plan Mode:                    Execution Mode:\n  ─────────────                 ────────────────\n  \"How should we restructure    \"Add input validation\n   the auth system?\"             to the signup handler\"\n\n  \"What's the best approach     \"Fix the NPE in\n   for adding caching?\"          UserService.java\"\n\n  \"Review this PR's              \"Write tests for\n   architecture\"                 the new endpoint\""},
			{Kind: types.Callout, Content: "A common workflow: start in plan mode to design the approach, review the plan, then switch to execution mode to implement it."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which mode should you use for 'How should we restructure the auth system?'",
				Choices:    []string{"Execution mode", "Plan mode", "Debug mode", "Test mode"},
				CorrectIdx: 1,
				Explanation: "Architecture questions benefit from plan mode — read-only exploration that produces a structured plan without making changes.",
			},
			{
				Kind:       types.FillBlank,
				Prompt:     "In plan mode, the AI operates in ___-only mode (no file modifications). What word fills the blank?",
				Answer:     "read",
				Explanation: "Plan mode is read-only — the AI analyzes and plans without making changes.",
			},
		},
	})
}
