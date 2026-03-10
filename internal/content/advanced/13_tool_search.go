package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      13,
		Title:   "Tool Search & Deferred Tools",
		Tier:    types.Advanced,
		Summary: "Lazy-loading tools to optimize context usage",
		VizBuilder: func(w, h int) interface{} { return viz.NewToolSearchModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "The Tool Search Problem"},
			{Kind: types.Paragraph, Content: "AI assistants can have hundreds of available tools — built-in tools, MCP server tools, and plugin tools. Loading all their definitions into the context window would waste thousands of tokens. The solution: deferred tools that load on demand."},
			{Kind: types.Heading, Content: "How Deferred Tools Work"},
			{Kind: types.Paragraph, Content: "Deferred tools are registered but NOT loaded into context. They appear in an <available-deferred-tools> list. When the AI needs one, it uses the ToolSearch tool to discover and load it — just in time."},
			{Kind: types.Code, Content: "  <available-deferred-tools>\n  mcp__slack__send_message\n  mcp__github__create_pr\n  mcp__database__query\n  NotebookEdit\n  WebSearch\n  </available-deferred-tools>"},
			{Kind: types.Heading, Content: "ToolSearch Query Modes"},
			{Kind: types.Bullet, Content: "Keyword search — \"slack message\" finds Slack messaging tools\nDirect selection — \"select:mcp__slack__send_message\" loads a specific tool\nRequired keyword — \"+slack send\" only searches Slack tools, ranked by \"send\"\nMulti-select — \"select:Read,Edit,Grep\" loads multiple tools at once"},
			{Kind: types.Heading, Content: "Why This Matters"},
			{Kind: types.Paragraph, Content: "Without deferred loading, 50 MCP tools at ~200 tokens each = 10,000 tokens wasted on tool definitions you may never use. With deferred loading, only the tools actually needed get loaded."},
			{Kind: types.Code, Content: "  Without deferred tools:\n  ┌──────────────────────────────┐\n  │ System Prompt (8k)           │\n  │ ALL tool definitions (10k)   │ ← wasted\n  │ Conversation (5k)            │\n  │ ...remaining: 177k           │\n  └──────────────────────────────┘\n\n  With deferred tools:\n  ┌──────────────────────────────┐\n  │ System Prompt (8k)           │\n  │ Core tools only (2k)         │ ← efficient\n  │ Conversation (5k)            │\n  │ ...remaining: 185k           │ ← 8k more!\n  └──────────────────────────────┘"},
			{Kind: types.Callout, Content: "Always use ToolSearch to load deferred tools BEFORE calling them. Calling a deferred tool without loading it first will fail."},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Why are some tools deferred instead of always loaded?",
				Choices:    []string{"They're not important", "To save context window tokens", "They're experimental", "To increase security"},
				CorrectIdx: 1,
				Explanation: "Deferred loading saves context window space by only loading tool definitions when they're actually needed.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which ToolSearch query would load a specific tool you know the name of?",
				Choices:    []string{"\"slack tools\"", "\"select:mcp__slack__send_message\"", "\"+find slack\"", "\"load slack\""},
				CorrectIdx: 1,
				Explanation: "Use \"select:tool_name\" for direct selection when you know the exact tool name.",
			},
		},
	})
}
