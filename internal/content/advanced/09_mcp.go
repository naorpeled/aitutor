package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      9,
		Title:   "MCP (Model Context Protocol)",
		Tier:    types.Advanced,
		Summary: "Extending AI with external tool servers",
		VizBuilder: func(w, h int) interface{} { return viz.NewMCPCallerModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Model Context Protocol (MCP)"},
			{Kind: types.Paragraph, Content: "MCP is an open protocol that lets AI assistants connect to external tool servers. Instead of building every capability into the AI itself, MCP allows a plug-in architecture where specialized servers provide tools, resources, and prompts."},
			{Kind: types.Heading, Content: "Architecture"},
			{Kind: types.Code, Content: "  ┌──────────┐     stdio/HTTP    ┌──────────────┐\n  │AI Client │◄════════════════►│  MCP Server  │\n  │          │                  │              │\n  └──────────┘                  │  ┌────────┐  │\n                                │  │ Tool 1 │  │\n                                │  ├────────┤  │\n                                │  │ Tool 2 │  │\n                                │  ├────────┤  │\n                                │  │Resource│  │\n                                │  └────────┘  │\n                                └──────────────┘"},
			{Kind: types.Heading, Content: "Key Concepts"},
			{Kind: types.Bullet, Content: "Tools — functions the AI can call (e.g., query database, send Slack message)\nResources — data the AI can read (e.g., documentation, API schemas)\nPrompts — reusable prompt templates for common tasks\nSampling — lets servers request the client to perform LLM completions\nTransports — communication channels (stdio for local, Streamable HTTP for remote)"},
			{Kind: types.Heading, Content: "Configuration"},
			{Kind: types.Code, Content: "  // MCP config (path varies: .claude/mcp.json, .cursor/mcp.json, etc.)\n  {\n    \"mcpServers\": {\n      \"github\": {\n        \"command\": \"gh-mcp-server\",\n        \"args\": [\"--repo\", \"owner/repo\"]\n      },\n      \"database\": {\n        \"command\": \"db-mcp-server\",\n        \"args\": [\"--connection\", \"postgres://...\"]\n      }\n    }\n  }"},
			{Kind: types.Callout, Content: "MCP turns AI assistants from closed systems into extensible platforms. Any developer can build an MCP server to give the AI new capabilities."},
			{Kind: types.Callout, Content: "Learn more: MCP Specification — https://spec.modelcontextprotocol.io/ | MCP Introduction — https://modelcontextprotocol.io/introduction"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What does MCP stand for?",
				Choices:    []string{"Model Control Protocol", "Model Context Protocol", "Machine Code Pipeline", "Multi-Channel Processor"},
				CorrectIdx: 1,
				Explanation: "MCP stands for Model Context Protocol — an open protocol for connecting AI to external tool servers.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Which transport does MCP use for local servers?",
				Choices:    []string{"HTTP", "WebSocket", "stdio", "gRPC"},
				CorrectIdx: 2,
				Explanation: "MCP uses stdio for local servers (process communication) and Streamable HTTP for remote servers.",
			},
		},
	})
}
