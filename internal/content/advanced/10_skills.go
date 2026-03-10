package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      10,
		Title:   "Skills",
		Tier:    types.Advanced,
		Summary: "Reusable workflows and specialized knowledge",
		VizBuilder: func(w, h int) interface{} { return viz.NewSkillLoadModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Skills"},
			{Kind: types.Paragraph, Content: "Skills are reusable, composable packages of specialized knowledge and workflows that extend an AI assistant's capabilities. They provide domain-specific expertise, step-by-step processes, and guardrails for common tasks."},
			{Kind: types.Heading, Content: "How Skills Work"},
			{Kind: types.Paragraph, Content: "Skills are loaded on-demand — they aren't always in the context window. When a task matches a skill's domain, the skill is invoked and its instructions are loaded. This keeps the context window efficient while providing deep expertise when needed."},
			{Kind: types.Code, Content: "  User: \"Create a new MCP server\"\n         │\n         ▼\n  ┌──────────────────┐\n  │ Skill Detection  │\n  │ \"mcp-builder\"    │\n  └────────┬─────────┘\n           ▼\n  ┌──────────────────┐\n  │  Load Skill      │\n  │  Instructions    │\n  └────────┬─────────┘\n           ▼\n  ┌──────────────────┐\n  │ Follow Workflow  │\n  │ Step by Step     │\n  └──────────────────┘"},
			{Kind: types.Heading, Content: "Skill Types"},
			{Kind: types.Bullet, Content: "Rigid skills — must be followed exactly (TDD, debugging workflows)\nFlexible skills — adapt principles to context (design patterns)\nProcess skills — define HOW to approach tasks (brainstorming, planning)\nImplementation skills — guide execution (frontend-design, mcp-builder)"},
			{Kind: types.Heading, Content: "Skill Composition"},
			{Kind: types.Paragraph, Content: "Skills can invoke other skills. For example, a \"build feature\" workflow might first invoke brainstorming, then planning, then test-driven development — each skill providing expertise for its phase."},
			{Kind: types.Callout, Content: "Skills turn tribal knowledge into reproducible workflows. Instead of remembering complex processes, encode them as skills that execute consistently every time."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Why are skills loaded on-demand rather than always present?",
				Choices:    []string{"They're too slow to load", "To save context window space", "They require special permissions", "They're experimental features"},
				CorrectIdx: 1,
				Explanation: "Skills are loaded on-demand to keep the context window efficient — only loading deep expertise when needed.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which type of skill should be followed EXACTLY without adaptation?",
				Choices:    []string{"Flexible skills", "Pattern skills", "Rigid skills", "Implementation skills"},
				CorrectIdx: 2,
				Explanation: "Rigid skills (like TDD and debugging workflows) must be followed exactly — their discipline is the point.",
			},
		},
	})
}
