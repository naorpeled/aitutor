package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type loopPhase int

const (
	loopRead loopPhase = iota
	loopThink
	loopAct
	loopObserve
)

func (p loopPhase) String() string {
	switch p {
	case loopRead:
		return "Read"
	case loopThink:
		return "Think"
	case loopAct:
		return "Act"
	case loopObserve:
		return "Observe"
	}
	return ""
}

type loopIteration struct {
	Read    string
	Think   string
	Act     string
	Observe string
	Success bool
}

// AgenticLoopModel walks through multiple iterations of an agentic loop.
type AgenticLoopModel struct {
	width      int
	height     int
	iterations []loopIteration
	iterIdx    int
	phase      loopPhase
	revealed   bool // whether current phase content is shown
	done       bool
}

func NewAgenticLoopModel(w, h int) Model {
	return &AgenticLoopModel{
		width:  w,
		height: h,
		iterations: []loopIteration{
			{
				Read:    "User: \"Fix the login bug — users get 500 error\"",
				Think:   "I need to find the login handler. Let me search for it.",
				Act:     "Grep(\"func.*login\", \"**/*.go\") → Found auth/handler.go:42",
				Observe: "Found the handler. Need to read it to understand the bug.",
				Success: false,
			},
			{
				Read:    "Read(auth/handler.go) → sees db.Query without error check",
				Think:   "The query result isn't checked for errors. A nil row causes the 500.",
				Act:     "Edit(auth/handler.go) → Added error check after db.Query",
				Observe: "Fix applied. Need to verify it compiles.",
				Success: false,
			},
			{
				Read:    "Bash(\"go build ./...\") → Build successful",
				Think:   "It compiles. Now I should run the tests to make sure nothing broke.",
				Act:     "Bash(\"go test ./auth/...\") → 1 test failed: TestLoginEmpty",
				Observe: "A test failed! The empty-input case hits a different code path. Need another iteration.",
				Success: false,
			},
			{
				Read:    "Read test output → TestLoginEmpty expects 400, got 500",
				Think:   "The empty email case also needs validation before the query.",
				Act:     "Edit(auth/handler.go) → Added input validation at top of handler",
				Observe: "Fix applied. Let me run tests again.",
				Success: false,
			},
			{
				Read:    "Bash(\"go test ./auth/...\") → All tests pass",
				Think:   "All tests pass. The fix handles both the nil-row bug and input validation.",
				Act:     "Bash(\"git add auth/handler.go && git diff --cached\")",
				Observe: "Changes look correct. Ready to report back to the user.",
				Success: true,
			},
		},
	}
}

func (m *AgenticLoopModel) Init() tea.Cmd { return nil }

func (m *AgenticLoopModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.done {
				return m, nil
			}
			if !m.revealed {
				m.revealed = true
				return m, nil
			}
			// Advance to next phase or next iteration
			m.revealed = false
			if m.phase < loopObserve {
				m.phase++
			} else {
				// Move to next iteration
				if m.iterIdx < len(m.iterations)-1 {
					m.iterIdx++
					m.phase = loopRead
				} else {
					m.done = true
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.iterIdx = 0
			m.phase = loopRead
			m.revealed = false
			m.done = false
		}
	}
	return m, nil
}

func (m *AgenticLoopModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	bright := lipgloss.NewStyle().Foreground(ui.ColorBright).Bold(true)
	green := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	red := lipgloss.NewStyle().Foreground(ui.ColorAdvanced).Bold(true)
	blue := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a5f3fc"))

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  The Agentic Loop"))
	lines = append(lines, "")

	// Loop diagram — highlight current phase
	phaseNames := []string{"Read", "Think", "Act", "Observe"}
	phaseStyle := func(i int) lipgloss.Style {
		if loopPhase(i) == m.phase && !m.done {
			return bright
		} else if loopPhase(i) < m.phase || m.iterIdx > 0 || m.done {
			return green
		}
		return dim
	}

	// Build the inner content: " Read ──> Think ──> Act ──> Observe "
	inner := fmt.Sprintf(" %s ──> %s ──> %s ──> %s ",
		phaseStyle(0).Render(phaseNames[0]),
		phaseStyle(1).Render(phaseNames[1]),
		phaseStyle(2).Render(phaseNames[2]),
		phaseStyle(3).Render(phaseNames[3]))

	// Measure the actual rendered content's visual width for matching borders
	innerWidth := lipgloss.Width(inner)
	border := strings.Repeat("─", innerWidth)

	lines = append(lines, dim.Render("  ┌"+border+"┐"))
	lines = append(lines, "  "+dim.Render("│")+inner+dim.Render("│"))
	lines = append(lines, dim.Render("  └"+border+"┘"))
	lines = append(lines, "")

	// Iteration counter
	iterLabel := fmt.Sprintf("  Iteration %d/%d", m.iterIdx+1, len(m.iterations))
	if m.done {
		lines = append(lines, green.Render(iterLabel+" — Done!"))
	} else {
		lines = append(lines, yellow.Render(iterLabel))
	}
	lines = append(lines, "")

	// Show completed iterations as compact summaries
	for i := 0; i < m.iterIdx; i++ {
		iter := m.iterations[i]
		icon := yellow.Render("↻")
		if iter.Success {
			icon = green.Render("✓")
		}
		// Compact summary of past iteration
		lines = append(lines, fmt.Sprintf("  %s Iter %d: %s",
			icon, i+1, dim.Render(truncate(iter.Act, 55))))
	}
	if m.iterIdx > 0 {
		lines = append(lines, "")
	}

	// Current iteration detail
	if !m.done {
		iter := m.iterations[m.iterIdx]
		lines = append(lines, accent.Render(fmt.Sprintf("  ── Iteration %d ──", m.iterIdx+1)))
		lines = append(lines, "")

		type phaseInfo struct {
			label   string
			content string
			style   lipgloss.Style
		}
		phaseData := []phaseInfo{
			{"Read", iter.Read, blue},
			{"Think", iter.Think, yellow},
			{"Act", iter.Act, codeStyle},
			{"Observe", iter.Observe, green},
		}

		for i, pd := range phaseData {
			lp := loopPhase(i)
			if lp < m.phase || (lp == m.phase && m.revealed) {
				// Show this phase
				label := pd.style.Render(fmt.Sprintf("  %-8s", pd.label))
				lines = append(lines, fmt.Sprintf("  %s %s", label, dim.Render(pd.content)))
			} else if lp == m.phase && !m.revealed {
				// Current phase, not yet revealed
				label := bright.Render(fmt.Sprintf("  ▸ %-6s", pd.label))
				lines = append(lines, fmt.Sprintf("  %s %s", label, dim.Render("Press Enter to reveal...")))
			} else {
				// Future phase
				lines = append(lines, fmt.Sprintf("  %s %s", dim.Render(fmt.Sprintf("  %-8s", pd.label)), dim.Render("...")))
			}
		}
	}

	lines = append(lines, "")

	// Summary when done
	if m.done {
		lines = append(lines, accent.Render("  ── Loop Complete ──"))
		lines = append(lines, "")
		lines = append(lines, green.Render(fmt.Sprintf("  Total iterations: %d", len(m.iterations))))
		lines = append(lines, dim.Render("  The AI kept looping until all tests passed."))
		lines = append(lines, dim.Render("  Each iteration: read context, reason about it,"))
		lines = append(lines, dim.Render("  take action, then observe the result."))
		lines = append(lines, "")
		lines = append(lines, dim.Render("  Key insight: the loop is self-correcting —"))
		lines = append(lines, dim.Render("  failures become input for the next iteration."))
		lines = append(lines, "")
		lines = append(lines, red.Render("  Without the loop: single-shot, hope for the best."))
		lines = append(lines, green.Render("  With the loop: iterative, self-correcting, reliable."))
		lines = append(lines, "")
		lines = append(lines, dim.Render("  [r] Reset"))
	} else {
		lines = append(lines, dim.Render("  [Enter/Space] Step through  [r] Reset"))
	}

	return strings.Join(lines, "\n")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
