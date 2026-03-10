package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:      12,
		Title:   "Git Worktrees",
		Tier:    types.Advanced,
		Summary: "Isolated workspaces for parallel development",
		VizBuilder: func(w, h int) interface{} { return viz.NewWorktreeSimModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Git Worktrees"},
			{Kind: types.Paragraph, Content: "Git worktrees let you check out multiple branches of the same repository simultaneously, each in its own directory. For AI-assisted development, this enables true parallel work without conflicts."},
			{Kind: types.Heading, Content: "How Worktrees Work"},
			{Kind: types.Code, Content: "  ~/project/           ← main worktree (main branch)\n  ~/project-worktrees/\n    ├── feature-auth/   ← worktree (feature/auth branch)\n    ├── fix-bug-123/    ← worktree (fix/bug-123 branch)\n    └── refactor-api/   ← worktree (refactor/api branch)"},
			{Kind: types.Paragraph, Content: "Each worktree is a full working copy with its own branch, staging area, and working directory. They all share the same .git data, so they're lightweight and fast to create."},
			{Kind: types.Heading, Content: "Worktrees + AI Agents"},
			{Kind: types.Bullet, Content: "Subagents get isolated worktrees — no merge conflicts\nMain workspace stays clean while agents work in parallel\nChanges can be reviewed per-worktree before merging\nAutomatic cleanup when agent work is done"},
			{Kind: types.Heading, Content: "Common Commands"},
			{Kind: types.Code, Content: "  # Create a new worktree\n  git worktree add ../feature-x -b feature/x\n\n  # List all worktrees\n  git worktree list\n\n  # Remove a worktree\n  git worktree remove ../feature-x"},
			{Kind: types.Heading, Content: "Best Practices"},
			{Kind: types.Bullet, Content: "Use a sibling directory for worktrees (not inside the repo)\nName worktrees after their branch for clarity\nClean up worktrees when branches are merged\nAvoid having two worktrees on the same branch"},
			{Kind: types.Callout, Content: "Worktrees are the key enabler for safe parallel AI development — they give each agent its own sandbox to work in."},
			{Kind: types.Callout, Content: "Learn more: Git Worktrees — https://git-scm.com/docs/git-worktree"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:       types.MultipleChoice,
				Prompt:     "What do git worktrees share between copies?",
				Choices:    []string{"Working directory", "Staging area", ".git data", "Branch name"},
				CorrectIdx: 2,
				Explanation: "Worktrees share the same .git data (object store, refs), making them lightweight. Each has its own working directory, staging area, and branch.",
			},
			{
				Kind:       types.FillBlank,
				Prompt:     "What command creates a new git worktree? (start with 'git')",
				Answer:     "git worktree add",
				Explanation: "The command 'git worktree add <path> -b <branch>' creates a new worktree.",
			},
		},
	})
}
