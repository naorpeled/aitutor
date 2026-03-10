package beginner

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      4,
		Title:   "Prompt Engineering",
		Tier:    types.Beginner,
		Summary: "Writing effective instructions for AI assistants",
		VizBuilder: func(w, h int) interface{} { return viz.NewPromptImproveModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Prompt Engineering"},
			{Kind: types.Paragraph, Content: "The quality of what you get from an AI assistant depends heavily on how you ask. Prompt engineering is the practice of crafting effective instructions that lead to better results."},
			{Kind: types.Heading, Content: "Principles of Good Prompts"},
			{Kind: types.Bullet, Content: "Be specific — \"Fix the login bug where users get 401 on valid tokens\" vs \"fix the bug\"\nProvide context — mention relevant files, error messages, expected behavior\nState the goal — explain what you want to achieve, not just what to change\nSet constraints — \"don't modify the public API\" or \"use existing patterns\""},
			{Kind: types.Heading, Content: "Before vs After"},
			{Kind: types.Code, Content: "  ✗ Bad:  \"make it work\"\n  ✓ Good: \"The UserService.GetByEmail method returns nil\n           instead of an error when the database query\n           fails. Fix it to return a wrapped error.\""},
			{Kind: types.Code, Content: "  ✗ Bad:  \"add tests\"\n  ✓ Good: \"Add unit tests for the ParseConfig function\n           covering: valid YAML, missing required fields,\n           and invalid port numbers.\""},
			{Kind: types.Heading, Content: "Iterative Refinement"},
			{Kind: types.Paragraph, Content: "You don't need the perfect prompt on the first try. Start with a clear request, review the result, and refine. The AI remembers the conversation context, so you can course-correct naturally."},
			{Kind: types.Callout, Content: "The best prompts give the AI the same information you'd give a skilled colleague: what's the problem, what have you tried, and what does success look like?"},
			{Kind: types.Callout, Content: "Learn more: Prompt Engineering Guide — https://www.promptingguide.ai"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which is a better prompt for fixing a bug?",
				Choices:    []string{"\"fix the bug\"", "\"make it work\"", "\"Fix the 401 error in LoginHandler when valid JWT tokens are rejected\"", "\"debug the code\""},
				CorrectIdx: 2,
				Explanation: "Specific prompts with context (what, where, when) produce much better results than vague requests.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What should you NOT do with prompt engineering?",
				Choices:    []string{"Provide context about the error", "State your constraints", "Expect perfection on the first try", "Mention relevant file paths"},
				CorrectIdx: 2,
				Explanation: "Prompt engineering is iterative — start clear, review results, and refine.",
			},
		},
	})
}
