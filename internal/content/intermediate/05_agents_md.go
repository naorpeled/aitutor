package intermediate

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      5,
		Title:   "AGENTS.md / CLAUDE.md",
		Tier:    types.Intermediate,
		Summary: "Project-specific AI configuration files",
		VizBuilder: func(w, h int) interface{} { return viz.NewClaudeMDBuilderModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "AGENTS.md & CLAUDE.md"},
			{Kind: types.Paragraph, Content: "CLAUDE.md (and AGENTS.md) are special configuration files that give AI assistants project-specific instructions. They're loaded automatically into the context window at the start of every conversation, acting as persistent memory for your project."},
			{Kind: types.Heading, Content: "File Hierarchy & Scope"},
			{Kind: types.Code, Content: "  ~/.claude/CLAUDE.md          ← user-level (all projects)\n  ~/project/.claude/CLAUDE.md  ← personal project config (git-ignored)\n  ~/project/CLAUDE.md          ← project root (committed, whole repo)\n  ~/project/src/CLAUDE.md      ← directory-level (scoped)\n  ~/project/AGENTS.md          ← cross-tool standard (committed)"},
			{Kind: types.Paragraph, Content: "Files cascade: user-level settings apply everywhere, project-level to the whole repo, and directory-level to specific areas. AGENTS.md serves the same purpose but is recognized by multiple AI coding tools — not just one — making it ideal for cross-tool teams."},
			{Kind: types.Heading, Content: "What to Put in CLAUDE.md"},
			{Kind: types.Bullet, Content: "Build & test commands — how to run the project\nCode conventions — naming, patterns, architecture decisions\nDo's and don'ts — project-specific rules and constraints\nFile structure — where things live and why\nCommon workflows — deployment, testing, review processes"},
			{Kind: types.Heading, Content: "Example"},
			{Kind: types.Code, Content: "  # CLAUDE.md\n\n  ## Build\n  - `make build` to compile\n  - `make test` to run all tests\n\n  ## Conventions\n  - Use snake_case for database columns\n  - All API handlers in internal/api/\n  - Never commit .env files"},
			{Kind: types.Callout, Content: "Think of CLAUDE.md as onboarding documentation for your AI assistant — the same things you'd tell a new team member on day one."},
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
				Prompt:     "Which CLAUDE.md file has the broadest scope?",
				Choices:    []string{"~/project/src/CLAUDE.md", "~/project/CLAUDE.md", "~/.claude/CLAUDE.md", "~/project/AGENTS.md"},
				CorrectIdx: 2,
				Explanation: "The user-level ~/.claude/CLAUDE.md applies to all projects, giving it the broadest scope.",
			},
		},
	})
}
