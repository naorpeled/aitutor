package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:         18,
		Title:      "Building Reliable Agents",
		Tier:       types.Advanced,
		Summary:    "Compaction, long-term memory, and durable agent state",
		SourceFile: "internal/content/advanced/18_agents_building.go",
		VizBuilder: func(w, h int) interface{} { return viz.NewAgentBuildingModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Building Reliable Agents"},
			{Kind: types.Paragraph, Content: "An AI coding agent is more than a prompt wrapped around tools. Reliable agents manage state deliberately: what stays in the active context window, what gets compacted into a handoff summary, and what becomes long-term memory for future sessions."},
			{Kind: types.Heading, Content: "Compaction"},
			{Kind: types.Paragraph, Content: "Compaction is the process of summarizing a long working session into a smaller, high-signal state bundle. A good compact summary keeps the goal, root cause, files touched, decisions made, commands run, failing or passing tests, and the next concrete step."},
			{Kind: types.Code, Content: "  Good compact handoff:\n  - Goal: fix parser panic on empty input\n  - Evidence: TestParseEmpty panics in parser.go:42\n  - Changed: parser.go guard before token lookup\n  - Verified: go test ./parser passes\n  - Next: run full go test ./... and open PR"},
			{Kind: types.Paragraph, Content: "Compaction is not the same as forgetting. The agent should compact after a major phase, before context gets tight, or before handing work to another session. It should not compact away unresolved evidence during active debugging."},
			{Kind: types.Heading, Content: "Long-term Memory"},
			{Kind: types.Paragraph, Content: "Long-term memory stores durable facts that should survive beyond the current task. Good candidates include stable user preferences, repository conventions, repeated commands, and architectural decisions. Temporary stack traces, current branch names, or half-formed guesses should stay out of memory."},
			{Kind: types.Code, Content: "  Save to memory:\n  - This repo uses table-driven tests\n  - Run make test-unit before committing\n  - User prefers concise final summaries\n\n  Do not save:\n  - Current failing line number\n  - Today's temporary branch\n  - A hypothesis not yet verified"},
			{Kind: types.Heading, Content: "Agent State Design"},
			{Kind: types.Bullet, Content: "Active context - detailed working set for the current task\nCompacted state - concise handoff that preserves progress\nLong-term memory - durable facts reused across sessions\nExternal artifacts - commits, test logs, issues, and PRs that outlive chat history"},
			{Kind: types.Heading, Content: "Practical Rule"},
			{Kind: types.Paragraph, Content: "Before saving anything, ask: will this still be true and useful next week? Before compacting, ask: could another agent continue the task from this summary without re-reading everything?"},
			{Kind: types.Callout, Content: "Reliable agents separate short-lived working context from durable memory. That separation prevents stale assumptions while still preserving useful experience."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:        types.MultipleChoice,
				Prompt:      "What belongs in a compaction summary?",
				Choices:     []string{"Every line of every file read", "Goal, evidence, decisions, commands, and next steps", "Only the final answer", "Only user preferences"},
				CorrectIdx:  1,
				Explanation: "A compact summary should preserve the state needed to continue work without replaying the whole conversation.",
			},
			{
				Kind:        types.FillBlank,
				Prompt:      "Stable repo conventions should be stored in long-term ____.",
				Answer:      "memory",
				Explanation: "Long-term memory is for durable information that remains useful across sessions.",
			},
			{
				Kind:        types.MultipleChoice,
				Prompt:      "Which item should NOT be saved as long-term memory?",
				Choices:     []string{"The repo uses go vet before PRs", "The user prefers concise final answers", "The current stack trace from one debugging run", "The team uses table-driven tests"},
				CorrectIdx:  2,
				Explanation: "A stack trace from one debugging run is task-local evidence, not durable memory.",
			},
		},
	})
}
