package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      5,
		Title:   "Project AI Config Files",
		Tier:    types.Intermediate,
		Summary: "Project-specific AI configuration files",
		VizBuilder: func(w, h int) interface{} { return viz.NewClaudeMDBuilderModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Project AI Config Files"},
			{Kind: types.Paragraph, Content: "AI coding tools use project-level configuration files to receive persistent instructions. These files are loaded automatically into the context window at the start of every conversation, acting as persistent memory for your project. Examples include AGENTS.md (cross-tool standard), CLAUDE.md (Claude Code), .cursorrules (Cursor), and copilot-instructions.md (GitHub Copilot)."},
			{Kind: types.Heading, Content: "Why They Matter: The Evidence"},
			{Kind: types.Paragraph, Content: "Research shows these files have a real, measurable impact. A study across 10 repositories and 124 pull requests found that AGENTS.md files correlated with a 28.64% median runtime reduction and 16.58% fewer output tokens, while maintaining the same task completion rates (Girgis et al., 2026)."},
			{Kind: types.Paragraph, Content: "Vercel ran an internal eval targeting 8 newer Next.js 16 APIs absent from model training data. Without config files, their agents hit a 53% pass rate. Adding a compressed 8KB AGENTS.md docs index raised that to 100% on their eval suite. The scope was narrow (their own hardened tests, single model), but it illustrates how much well-structured project context can help — especially for APIs the model hasn't seen before."},
			{Kind: types.Callout, Content: "Think of these config files as onboarding documentation for your AI assistant — the same things you'd tell a new team member on day one."},
			{Kind: types.Heading, Content: "File Hierarchy & Scope"},
			{Kind: types.Code, Content: "  Example hierarchy (Claude Code):\n  ~/.claude/CLAUDE.md          ← user-level (all projects)\n  ~/project/.claude/CLAUDE.md  ← personal project config (git-ignored)\n  ~/project/CLAUDE.md          ← project root (committed, whole repo)\n  ~/project/src/CLAUDE.md      ← directory-level (scoped)\n\n  Cross-tool standard:\n  ~/project/AGENTS.md          ← recognized by multiple AI tools"},
			{Kind: types.Paragraph, Content: "Most tools support a hierarchy where files cascade: user-level settings apply everywhere, project-level to the whole repo, and directory-level to specific areas. AGENTS.md is the cross-tool standard — recognized by multiple AI coding tools, making it ideal for teams using different editors."},
			{Kind: types.Heading, Content: "The Instruction Budget: Less Is More"},
			{Kind: types.Paragraph, Content: "Here's the catch: more instructions aren't always better. Research by Gloaguen et al. (2025) found that overly comprehensive config files can actually reduce task success rates while increasing inference cost by over 20%. Agents follow instructions faithfully — but unnecessary requirements trigger broader exploration that can derail the task."},
			{Kind: types.Paragraph, Content: "Frontier LLMs can follow roughly 150-200 instructions with reasonable consistency. Every token in your config file loads on every request, regardless of relevance. Large files waste context and confuse agents; small, focused files leave more capacity for the actual task."},
			{Kind: types.Heading, Content: "What to Put in Your Config File"},
			{Kind: types.Bullet, Content: "Project description — one sentence explaining what this project is\nBuild & test commands — how to compile, test, and lint\nPackage manager — only if it's not the default (e.g., bun instead of npm)\nCode conventions — naming patterns, architecture decisions\nDo's and don'ts — critical project-specific rules"},
			{Kind: types.Heading, Content: "What NOT to Put In"},
			{Kind: types.Bullet, Content: "Exhaustive file listings — paths drift quickly; describe capabilities instead\nEverything you know — keep it minimal, use progressive disclosure\nAuto-generated content — these prioritize comprehensiveness over restraint\nStale documentation — outdated info actively poisons the agent's context"},
			{Kind: types.Heading, Content: "Progressive Disclosure"},
			{Kind: types.Paragraph, Content: "For larger projects, keep AGENTS.md as a concise index and move detailed guidance into separate files. Vercel found that a compressed 8KB docs index (reduced from 40KB) paired with retrieval guidance outperformed having everything inline."},
			{Kind: types.Code, Content: "  # AGENTS.md (keep this lean ~150 lines)\n\n  ## Build\n  - `make build` to compile\n  - `make test` to run all tests\n\n  ## Conventions\n  - Use snake_case for database columns\n  - All API handlers in internal/api/\n  - Never commit .env files\n\n  ## Detailed Docs\n  - See docs/TYPESCRIPT.md for language conventions\n  - See docs/API.md for endpoint patterns"},
			{Kind: types.Callout, Content: "Treat your config file like production code: review additions critically, remove stale entries, and resist the urge to add rules reactively. An unmaintained file is worse than no file at all."},
			{Kind: types.Heading, Content: "References"},
			{Kind: types.Bullet, Content: "Girgis et al. (2026) — \"AGENTS.md Files and AI Coding Agent Efficiency\" — arxiv.org/abs/2601.20404\nGloaguen et al. (2025) — \"Evaluating AGENTS.md for Coding Agents\" — arxiv.org/abs/2602.11988\nVercel — \"AGENTS.md outperforms skills in our agent evals\" — vercel.com/blog/agents-md-outperforms-skills-in-our-agent-evals\nAI Hero — \"A Complete Guide to AGENTS.md\" — aihero.dev/a-complete-guide-to-agents-md"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What is the difference between CLAUDE.md and AGENTS.md?",
				Choices:    []string{"CLAUDE.md is for Claude, AGENTS.md is for other AIs", "CLAUDE.md is tool-specific, AGENTS.md is a cross-tool standard", "There is no difference", "AGENTS.md is deprecated"},
				CorrectIdx: 1,
				Explanation: "CLAUDE.md is specific to Claude Code and can exist at multiple scopes. AGENTS.md is a cross-tool standard recognized by multiple AI coding assistants, designed to be committed to version control.",
			},
			{
				Kind:       types.MultipleChoice,
				Prompt:     "Why can overly comprehensive config files hurt agent performance?",
				Choices:    []string{"They crash the AI", "Unnecessary instructions trigger broader exploration that derails tasks", "They exceed the file size limit", "They conflict with the system prompt"},
				CorrectIdx: 1,
				Explanation: "Research found that agents follow instructions faithfully — but unnecessary requirements cause broader exploration, reducing success rates while increasing inference cost by over 20%.",
			},
		},
	})
}
