package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      15,
		Title:   "The Agentic Loop",
		Tier:    types.Intermediate,
		Summary: "How AI agents iterate to solve problems",
		VizBuilder: func(w, h int) interface{} { return viz.NewAgenticLoopModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "The Agentic Loop"},
			{Kind: types.Paragraph, Content: "At the heart of every AI coding agent is a loop. The agent doesn't just generate code in one shot — it iterates. It reads context, reasons about it, takes an action, observes the result, and loops back. This is what makes agents fundamentally different from simple code generators."},
			{Kind: types.Heading, Content: "The Core Loop: Read → Think → Act → Observe"},
			{Kind: types.Paragraph, Content: "Every agentic system follows some variation of this pattern:"},
			{Kind: types.Code, Content: "  Read ──> Think ──> Act ──> Observe\n    ^                           │\n    └───────── loops back ──────┘\n\n  Read:     Gather context (files, errors, docs)\n  Think:    Reason about what to do next\n  Act:      Execute a tool (edit, run, search)\n  Observe:  Check the result, feed it back in"},
			{Kind: types.Paragraph, Content: "This loop runs until the task is complete or the agent determines it cannot proceed. The key insight is that the output of one iteration becomes the input of the next — making the process self-correcting."},
			{Kind: types.Heading, Content: "Why Loops Matter"},
			{Kind: types.Paragraph, Content: "Without loops, an AI would be limited to single-shot responses: generate code and hope it works. With loops, the agent can:"},
			{Kind: types.Bullet, Content: "Verify its own work by running tests\nSelf-correct when something fails\nGather more context when the first attempt doesn't have enough info\nBreak complex tasks into smaller steps, tackling each iteratively\nAdapt when the codebase doesn't match its expectations"},
			{Kind: types.Heading, Content: "Real Example: Fixing a Bug"},
			{Kind: types.Code, Content: "  Iteration 1: Search for the bug\n    Read:    Grep for the error message\n    Think:   Found the handler, need to read it\n    Act:     Read the file\n    Observe: See missing error check → need to fix\n\n  Iteration 2: Apply the fix\n    Read:    Understand the code around the bug\n    Think:   Add error handling after the query\n    Act:     Edit the file\n    Observe: Fix applied → need to test\n\n  Iteration 3: Verify\n    Read:    Run tests\n    Think:   One test failed! Different code path.\n    Act:     Edit to handle that case too\n    Observe: Tests pass → done!"},
			{Kind: types.Heading, Content: "Loop Variants"},
			{Kind: types.Paragraph, Content: "Different agentic frameworks name the steps differently, but the pattern is the same:"},
			{Kind: types.Code, Content: "  Pattern          Steps\n  ──────           ─────\n  OODA             Observe → Orient → Decide → Act\n  ReAct            Reason → Act → Observe\n  RALPH            Read → Act → Log → Plan → Hypothesize\n  General agent    Perceive → Decide → Execute → Evaluate"},
			{Kind: types.Callout, Content: "The names vary, but every agentic system follows the same principle: gather information, reason about it, take action, check the result, and repeat. The loop is what turns a language model into an agent."},
			{Kind: types.Heading, Content: "What Controls the Loop?"},
			{Kind: types.Bullet, Content: "Stop condition — the agent decides the task is complete (tests pass, user confirms)\nMax iterations — safety limit to prevent infinite loops\nError handling — the agent can break out if it's stuck\nUser intervention — the human can redirect or stop the agent\nToken budget — the context window limits how many iterations fit"},
			{Kind: types.Heading, Content: "Single-Shot vs Agentic"},
			{Kind: types.Code, Content: "  Single-shot:              Agentic:\n  ──────────               ────────\n  Prompt → Response         Prompt → Loop ──┐\n  (hope it's right)               ↑         │\n                                  └─────────┘\n                            (verify and correct)\n\n  Speed: Fast               Speed: Slower per task\n  Accuracy: Variable        Accuracy: High (self-correcting)\n  Complexity: Simple only   Complexity: Handles hard tasks"},
			{Kind: types.Callout, Content: "Try the visualization to step through a real debugging scenario. Watch how each iteration builds on the previous one — failures aren't dead ends, they're information that drives the next iteration."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:        types.MultipleChoice,
				Prompt:      "What makes an agentic loop different from a single-shot response?",
				Choices:     []string{"It uses more tokens", "It iterates — observing results and self-correcting", "It always writes better code", "It runs faster"},
				CorrectIdx:  1,
				Explanation: "The agentic loop iterates: it takes action, observes the result, and uses that observation to inform the next step. This self-correcting behavior is what separates agents from one-shot generation.",
			},
			{
				Kind:        types.MultipleChoice,
				Prompt:      "In a typical agentic loop, what happens when a test fails?",
				Choices:     []string{"The agent gives up and reports failure", "The failure becomes input for the next iteration", "The agent restarts from scratch", "The agent asks the user what to do"},
				CorrectIdx:  1,
				Explanation: "Failures are information. The test output becomes the 'Read' input of the next iteration — the agent analyzes what went wrong and tries a different approach.",
			},
		},
	})
}
