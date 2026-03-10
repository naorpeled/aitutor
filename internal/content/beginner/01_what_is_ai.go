package beginner

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      1,
		Title:   "What is an AI Coding Assistant?",
		Tier:    types.Beginner,
		Summary: "Understanding how AI assistants help you write code",
		VizBuilder: func(w, h int) interface{} { return viz.NewAgentLoopModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "What is an AI Coding Assistant?"},
			{Kind: types.Paragraph, Content: "An AI coding assistant is a tool powered by large language models (LLMs) that helps developers write, understand, and modify code. Unlike traditional autocomplete or linters, AI assistants understand natural language instructions and can reason about code at a high level."},
			{Kind: types.Heading, Content: "The Agent Loop"},
			{Kind: types.Paragraph, Content: "Modern AI coding assistants operate in an \"agent loop\" — a cycle of reading context, reasoning about what to do, taking actions (like editing files or running commands), and observing results. This loop continues until the task is complete."},
			{Kind: types.Code, Content: "  ┌──────────────────┐\n  │   User Request   │\n  └────────┬─────────┘\n           ▼\n  ┌─────────────────┐\n  │  Read Context   │◄──────┐\n  └────────┬────────┘       │\n           ▼                │\n  ┌─────────────────┐       │\n  │    Reason &     │       │\n  │     Plan        │       │\n  └────────┬────────┘       │\n           ▼                │\n  ┌─────────────────┐       │\n  │  Take Action    │       │\n  │  (tool call)    │       │\n  └────────┬────────┘       │\n           ▼                │\n  ┌─────────────────┐       │\n  │ Observe Result  │───────┘\n  └────────┬────────┘\n           ▼\n  ┌─────────────────┐\n  │    Response     │\n  └─────────────────┘"},
			{Kind: types.Heading, Content: "Key Capabilities"},
			{Kind: types.Bullet, Content: "Code generation — write new code from natural language descriptions\nCode editing — modify existing code with precise changes\nCode explanation — understand and explain complex codebases\nBug fixing — identify and fix issues in your code\nRefactoring — improve code structure while preserving behavior\nTest writing — generate tests for your code"},
			{Kind: types.Heading, Content: "Tools, Not Magic"},
			{Kind: types.Paragraph, Content: "AI assistants are tools that augment your abilities. They work best when you provide clear context, review their output, and guide them when they go astray. Understanding how they work — which is what this tutorial teaches — makes you far more effective at using them."},
			{Kind: types.Callout, Content: "The most productive developers don't just use AI assistants — they understand how they work under the hood. That's exactly what you'll learn in this tutorial."},
			{Kind: types.Callout, Content: "Learn more: Large Language Models — https://en.wikipedia.org/wiki/Large_language_model"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What is the 'agent loop' in AI coding assistants?",
				Choices:    []string{"A programming language feature", "A cycle of read context → reason → act → observe", "A type of infinite loop bug", "A user interface pattern"},
				CorrectIdx: 1,
				Explanation: "The agent loop is the core cycle: reading context, reasoning about what to do, taking action, and observing results.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which of these is NOT a key capability of AI coding assistants?",
				Choices:    []string{"Code generation", "Bug fixing", "Replacing developers entirely", "Refactoring"},
				CorrectIdx: 2,
				Explanation: "AI assistants augment developer abilities — they're tools, not replacements.",
			},
		},
	})
}
