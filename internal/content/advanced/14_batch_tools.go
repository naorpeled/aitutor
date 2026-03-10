package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      14,
		Title:   "Batch Tool Calls",
		Tier:    types.Advanced,
		Summary: "Per-tool execution policies and parallel batching",
		VizBuilder: func(w, h int) interface{} { return viz.NewBatchToolModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Batch Tool Calls"},
			{Kind: types.Paragraph, Content: "AI assistants can call multiple tools simultaneously in a single response. But not every tool is safe to run in parallel — some have side effects that require exclusive execution. Modern AI systems let each tool declare its own batch policy."},
			{Kind: types.Heading, Content: "Per-Tool Batch Policies"},
			{Kind: types.Paragraph, Content: "Each tool declares whether it can be batched (run alongside other tools) or must run alone (sequential). Read-only tools like file readers and search tools are typically batchable. Tools with side effects like file editors and shell commands often must run sequentially."},
			{Kind: types.Code, Content: "  Example tool set:\n  Tool         Policy        Why\n  ────         ──────        ───\n  file read    ⚡ batchable   Read-only, no side effects\n  search       ⚡ batchable   Read-only search\n  grep         ⚡ batchable   Read-only search\n  file edit    🔒 sequential  Modifies files\n  file write   🔒 sequential  Creates/overwrites files\n  shell        🔒 sequential  Arbitrary side effects"},
			{Kind: types.Heading, Content: "How Batching Works"},
			{Kind: types.Paragraph, Content: "The AI groups consecutive batchable tool calls into a single round trip. When it hits a sequential-only tool, it flushes the batch and runs that tool alone. This creates an execution plan that maximizes parallelism while respecting safety constraints."},
			{Kind: types.Code, Content: "  Tool calls in order:         Execution plan:\n  ───────────────────         ───────────────\n  Read(go.mod)      ─┐\n  Read(main.go)     ─┤ batch   → Round 1 (3 parallel)\n  Grep(\"TODO\")      ─┘\n  Edit(main.go)     ─── alone  → Round 2 (1 alone)\n  Bash(go build)    ─── alone  → Round 3 (1 alone)\n  Bash(go test)     ─── alone  → Round 4 (1 alone)\n\n  Result: 4 round trips instead of 6"},
			{Kind: types.Heading, Content: "Why Not Batch Everything?"},
			{Kind: types.Bullet, Content: "Edit + Edit on the same file could conflict\nBash commands may depend on prior Edit results\nWrite could create a file that another tool reads\nSequential tools need the results of previous steps"},
			{Kind: types.Heading, Content: "Maximizing Parallelism"},
			{Kind: types.Paragraph, Content: "Understanding batch policies helps you structure requests for maximum speed. Front-load your reads — ask the AI to gather all information first, then make changes. This naturally groups batchable reads together."},
			{Kind: types.Code, Content: "  Slow (interleaved):           Fast (reads first):\n  ──────────────────           ──────────────────\n  Read A → Edit A               Read A ─┐\n  Read B → Edit B               Read B ─┤ 1 round trip\n  Read C → Edit C               Read C ─┘\n  = 6 round trips               Edit A → Edit B → Edit C\n                                = 4 round trips"},
			{Kind: types.Callout, Content: "Try it in the visualization: toggle each tool's batch policy and see how the execution plan changes. Notice how grouping reads together saves round trips."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Why are Read and Grep typically marked as batchable?",
				Choices:    []string{"They're faster than other tools", "They're read-only with no side effects", "They use less tokens", "They always return small results"},
				CorrectIdx: 1,
				Explanation: "Read-only tools have no side effects, so running them in parallel is safe — they can't interfere with each other.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What happens when the AI encounters a sequential-only tool in a batch?",
				Choices:    []string{"It skips the tool", "It converts the tool to batchable", "It flushes the current batch and runs the tool alone", "It waits for user confirmation"},
				CorrectIdx: 2,
				Explanation: "Sequential-only tools cause the AI to flush any pending batch, execute it, then run the sequential tool on its own before continuing.",
			},
		},
	})
}
