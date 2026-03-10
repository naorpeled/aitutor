package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      11,
		Title:   "Subagents",
		Tier:    types.Advanced,
		Summary: "Parallel execution with specialized sub-processes",
		VizBuilder: func(w, h int) interface{} { return viz.NewFanoutModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Subagents"},
			{Kind: types.Paragraph, Content: "Subagents are specialized AI processes spawned by a main agent to handle specific tasks in parallel. They enable concurrent work on independent problems, dramatically speeding up complex tasks."},
			{Kind: types.Heading, Content: "How Subagents Work"},
			{Kind: types.Code, Content: "  ┌─────────────┐\n  │  Main Agent  │\n  │  (orchestr.) │\n  └──────┬──────┘\n    ┌────┼────┐\n    ▼    ▼    ▼\n  ┌───┐┌───┐┌───┐\n  │ A ││ B ││ C │  ← parallel execution\n  └─┬─┘└─┬─┘└─┬─┘\n    └────┼────┘\n         ▼\n  ┌─────────────┐\n  │  Main Agent  │\n  │  (combines)  │\n  └─────────────┘"},
			{Kind: types.Heading, Content: "Subagent Types"},
			{Kind: types.Bullet, Content: "Explore — fast codebase exploration, read-only\nGeneral-purpose — full tool access for complex tasks\nPlan — architectural planning, read-only\nSpecialized — custom agents for specific workflows (code-review, test-runner)"},
			{Kind: types.Heading, Content: "When to Use Subagents"},
			{Kind: types.Bullet, Content: "Multiple independent research questions\nParallel file searches across different areas\nImplementing unrelated changes simultaneously\nRunning tests while making other changes"},
			{Kind: types.Heading, Content: "Isolation"},
			{Kind: types.Paragraph, Content: "Subagents can run in isolated git worktrees — separate copies of the repository where they can make changes without affecting the main workspace. This prevents conflicts when multiple agents edit files simultaneously."},
			{Kind: types.Callout, Content: "Think of subagents as team members you can spin up instantly. Each gets a clear task, works independently, and reports back results."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What type of subagent is best for quick codebase exploration?",
				Choices:    []string{"General-purpose", "Plan", "Explore", "Specialized"},
				CorrectIdx: 2,
				Explanation: "Explore agents are fast, read-only agents specialized for searching and navigating codebases.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "How do subagents avoid merge conflicts when working in parallel?",
				Choices:    []string{"They lock files", "They use isolated git worktrees", "They take turns", "They can't — conflicts are inevitable"},
				CorrectIdx: 1,
				Explanation: "Subagents can run in isolated git worktrees — separate copies of the repository that prevent conflicts.",
			},
		},
	})
}
