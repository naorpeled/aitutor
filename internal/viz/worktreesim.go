package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type worktreeEntry struct {
	Branch string
	Path   string
	Agent  string
	Color  lipgloss.Color
}

// WorktreeSimModel simulates creating and managing git worktrees.
type WorktreeSimModel struct {
	width     int
	height    int
	worktrees []worktreeEntry
	cursor    int
}

var worktreeTemplates = []worktreeEntry{
	{Branch: "feature/auth", Path: "../wt/feature-auth", Agent: "Agent 1", Color: ui.ColorHighlight},
	{Branch: "fix/bug-42", Path: "../wt/fix-bug-42", Agent: "Agent 2", Color: ui.ColorIntermediate},
	{Branch: "refactor/api", Path: "../wt/refactor-api", Agent: "Agent 3", Color: ui.ColorAccent},
}

func NewWorktreeSimModel(w, h int) Model {
	return &WorktreeSimModel{
		width: w,
		height: h,
		worktrees: []worktreeEntry{
			{Branch: "main", Path: "~/project", Agent: "You", Color: ui.ColorBeginner},
		},
	}
}

func (m *WorktreeSimModel) Init() tea.Cmd { return nil }

func (m *WorktreeSimModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("a"))):
			// Add next worktree
			addIdx := len(m.worktrees) - 1
			if addIdx < len(worktreeTemplates) {
				m.worktrees = append(m.worktrees, worktreeTemplates[addIdx])
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("x"))):
			// Remove last worktree (can't remove main)
			if len(m.worktrees) > 1 {
				m.worktrees = m.worktrees[:len(m.worktrees)-1]
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.worktrees)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.worktrees = []worktreeEntry{
				{Branch: "main", Path: "~/project", Agent: "You", Color: ui.ColorBeginner},
			}
			m.cursor = 0
		}
	}
	return m, nil
}

func (m *WorktreeSimModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	active := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Worktree Manager"))
	lines = append(lines, dim.Render("  Create isolated workspaces for parallel development"))
	lines = append(lines, "")

	// Git worktree list
	lines = append(lines, dim.Render("  $ git worktree list"))
	lines = append(lines, "")

	for i, wt := range m.worktrees {
		style := lipgloss.NewStyle().Foreground(wt.Color).Bold(true)
		prefix := "  "
		if i == m.cursor {
			prefix = "▸ "
		}

		branchLabel := style.Render(fmt.Sprintf("%-20s", wt.Branch))
		pathLabel := dim.Render(fmt.Sprintf("%-25s", wt.Path))
		agentLabel := style.Render("[" + wt.Agent + "]")

		lines = append(lines, fmt.Sprintf("  %s%s  %s  %s", prefix, branchLabel, pathLabel, agentLabel))
	}

	lines = append(lines, "")

	// Simulated command
	if len(m.worktrees) < len(worktreeTemplates)+1 {
		next := worktreeTemplates[len(m.worktrees)-1]
		lines = append(lines, highlight.Render("  Next: ")+dim.Render(fmt.Sprintf("git worktree add %s -b %s", next.Path, next.Branch)))
	}

	// Visual tree
	lines = append(lines, "")
	lines = append(lines, dim.Render("  Repository Structure:"))
	lines = append(lines, "")
	for i, wt := range m.worktrees {
		style := lipgloss.NewStyle().Foreground(wt.Color)
		connector := "├──"
		if i == len(m.worktrees)-1 {
			connector = "└──"
		}
		if i == 0 {
			lines = append(lines, style.Render(fmt.Sprintf("  %s (.git shared)", wt.Path)))
		} else {
			lines = append(lines, style.Render(fmt.Sprintf("  %s %s [%s]", connector, wt.Path, wt.Branch)))
		}
	}

	lines = append(lines, "")
	lines = append(lines, dim.Render(fmt.Sprintf("  %d worktree(s) active", len(m.worktrees))))

	addCount := len(worktreeTemplates) + 1 - len(m.worktrees)
	if addCount > 0 {
		lines = append(lines, "")
		lines = append(lines, active.Render("  [a] Add worktree")+"  "+dim.Render(fmt.Sprintf("(%d available)", addCount)))
	}
	if len(m.worktrees) > 1 {
		lines = append(lines, dim.Render("  [x] Remove last worktree"))
	}

	lines = append(lines, "", dim.Render("  [a] Add  [x] Remove  [↑/↓] Navigate  [r] Reset"))

	return strings.Join(lines, "\n")
}
