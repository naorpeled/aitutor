package beginner

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      2,
		Title:   "Context Window",
		Tier:    types.Beginner,
		Summary: "How AI models see and process information",
		VizBuilder: func(w, h int) interface{} { return viz.NewBucketModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "The Context Window"},
			{Kind: types.Paragraph, Content: "Every AI model has a \"context window\" — a fixed amount of text (measured in tokens) it can process at once. Think of it as the model's working memory. Everything the model needs to understand must fit within this window."},
			{Kind: types.Heading, Content: "What Goes Into the Context Window?"},
			{Kind: types.Bullet, Content: "System prompt — instructions that define the model's behavior\nProject config files — tool-specific instructions loaded automatically (e.g., CLAUDE.md, .cursorrules, copilot-instructions.md)\nConversation history — all previous messages in the session\nTool results — output from file reads, searches, command execution\nFile contents — code files you're working with"},
			{Kind: types.Callout, Content: "Learn more: What are tokens? Try the interactive tokenizer — https://platform.openai.com/tokenizer"},
			{Kind: types.Heading, Content: "Token Limits"},
			{Kind: types.Paragraph, Content: "Context windows vary by model — for example, Claude offers 200K tokens, GPT-4o offers 128K, and Gemini offers up to 1M. While these are large, they're not infinite. Large codebases, long conversations, and verbose tool outputs can fill them up."},
			{Kind: types.Code, Content: "  Example: 200,000 tokens ≈\n  • 150,000 words\n  • 30,000 lines of code\n  • ~100 average source files"},
			{Kind: types.Heading, Content: "Context Window Management"},
			{Kind: types.Paragraph, Content: "Smart AI tools automatically manage the context window for you. They compress old messages, summarize tool results, and prioritize the most relevant information. But understanding this constraint helps you work more effectively."},
			{Kind: types.Callout, Content: "When you notice an AI assistant 'forgetting' something from earlier in a long conversation, it's likely because that information was compressed or evicted from the context window."},
			{Kind: types.Heading, Content: "The Hidden Cost: MCP Tool Definitions"},
			{Kind: types.Paragraph, Content: "Every MCP server you enable adds its tool definitions to the context window — even before you use any of those tools. Each tool definition costs ~200 tokens. Some servers are surprisingly bloated — GitHub alone exposes 34 tools."},
			{Kind: types.Code, Content: "  MCP servers enabled:        Context cost:\n  ─────────────────────       ──────────────\n  GitHub (34 tools)           6,800 tokens\n  Slack (18 tools)            3,600 tokens\n  Jira (22 tools)             4,400 tokens\n  Linear (16 tools)           3,200 tokens\n  ─────────────────────       ──────────────\n  Total                      18,000 tokens\n                              (wasted if unused!)"},
			{Kind: types.Callout, Content: "Only enable MCP servers you actively need — and disable individual tools within servers you do use. GitHub has 34 tools, but you may only need 5. Each unused tool definition wastes ~200 tokens of context. Try it in the visualization!"},
			{Kind: types.Heading, Content: "Strategies for Effective Context Use"},
			{Kind: types.Bullet, Content: "Be specific — targeted requests use less context than vague ones\nUse project config files (e.g., CLAUDE.md, .cursorrules) — persistent instructions don't need repeating\nBreak up large tasks — smaller, focused sessions are more effective\nReference files by path — let the AI read only what it needs\nDisable unused MCP servers — tool definitions eat context silently\nDisable individual tools within servers — keep only what you use"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What is the context window?",
				Choices:    []string{"A GUI window showing context menus", "The fixed amount of text an AI can process at once", "A debugging tool for viewing variables", "The terminal window where you type"},
				CorrectIdx: 1,
				Explanation: "The context window is the model's working memory — all input must fit within its token limit.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What happens when a conversation exceeds the context window?",
				Choices:    []string{"The AI crashes", "Older messages are compressed or evicted", "The context window automatically expands", "The AI asks you to start over"},
				CorrectIdx: 1,
				Explanation: "When the context window fills up, AI tools compress or evict older messages to make room for new content.",
			},
		},
	})
}
