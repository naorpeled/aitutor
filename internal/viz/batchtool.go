package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type batchTool struct {
	Name      string
	Batchable bool // can run in parallel with others
	Args      string
	Result    string
	Done      bool
}

// BatchToolModel lets users toggle per-tool batch policies and see execution impact.
type BatchToolModel struct {
	width    int
	height   int
	tools    []batchTool
	cursor   int
	phase    int // 0=configure policies, 1=execute, 2=done
	execStep int // current execution step in phase 1
	plan     [][]int // computed execution plan: groups of tool indices
}

func NewBatchToolModel(w, h int) Model {
	m := &BatchToolModel{
		width:  w,
		height: h,
		tools: []batchTool{
			{Name: "Read", Batchable: true, Args: "go.mod", Result: "module github.com/app"},
			{Name: "Read", Batchable: true, Args: "main.go", Result: "package main..."},
			{Name: "Grep", Batchable: true, Args: "\"TODO\"", Result: "Found 3 matches"},
			{Name: "Glob", Batchable: true, Args: "**/*_test.go", Result: "Found 8 test files"},
			{Name: "Edit", Batchable: false, Args: "main.go", Result: "✓ Applied changes"},
			{Name: "Bash", Batchable: false, Args: "go build", Result: "✓ Build OK"},
			{Name: "Bash", Batchable: false, Args: "go test ./...", Result: "✓ All tests pass"},
			{Name: "Write", Batchable: false, Args: "config.yaml", Result: "✓ File created"},
		},
	}
	m.computePlan()
	return m
}

// computePlan groups tools into execution batches based on their policies.
func (m *BatchToolModel) computePlan() {
	m.plan = nil
	var currentBatch []int

	for i, t := range m.tools {
		if t.Batchable {
			// Batchable tools accumulate into the current batch
			currentBatch = append(currentBatch, i)
		} else {
			// Non-batchable: flush any pending batch, then run alone
			if len(currentBatch) > 0 {
				m.plan = append(m.plan, currentBatch)
				currentBatch = nil
			}
			m.plan = append(m.plan, []int{i})
		}
	}
	if len(currentBatch) > 0 {
		m.plan = append(m.plan, currentBatch)
	}
}

func (m *BatchToolModel) Init() tea.Cmd { return nil }

func (m *BatchToolModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.phase == 0 && m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.phase == 0 && m.cursor < len(m.tools)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			switch m.phase {
			case 0:
				// Toggle batchable policy for selected tool
				m.tools[m.cursor].Batchable = !m.tools[m.cursor].Batchable
				m.computePlan()
			case 1:
				// Execute next batch
				if m.execStep < len(m.plan) {
					for _, idx := range m.plan[m.execStep] {
						m.tools[idx].Done = true
					}
					m.execStep++
					if m.execStep >= len(m.plan) {
						m.phase = 2
					}
				}
			case 2:
				// do nothing, use r to reset
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
			if m.phase == 0 {
				// Switch to execution phase
				m.phase = 1
				m.execStep = 0
				for i := range m.tools {
					m.tools[i].Done = false
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.phase = 0
			m.execStep = 0
			m.cursor = 0
			for i := range m.tools {
				m.tools[i].Done = false
			}
			m.computePlan()
		}
	}
	return m, nil
}

func (m *BatchToolModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	done := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	pending := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	redStyle := lipgloss.NewStyle().Foreground(ui.ColorIncorrect).Bold(true)
	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("#f97316"))

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Batch Tool Calls — Per-Tool Policies"))
	lines = append(lines, "")

	// Show execution plan summary
	lines = append(lines, dim.Render(fmt.Sprintf("  Execution plan: %d round trips for %d tool calls",
		len(m.plan), len(m.tools))))

	batchable := 0
	for _, t := range m.tools {
		if t.Batchable {
			batchable++
		}
	}
	lines = append(lines, dim.Render(fmt.Sprintf("  Batchable: %d    Sequential-only: %d",
		batchable, len(m.tools)-batchable)))
	lines = append(lines, "")

	if m.phase == 0 {
		// ── Configure phase ──
		lines = append(lines, highlight.Render("  Configure each tool's batch policy:"))
		lines = append(lines, dim.Render("  Toggle whether each tool can run in parallel with others."))
		lines = append(lines, "")

		for i, t := range m.tools {
			style := dim
			prefix := "  "
			if i == m.cursor {
				prefix = "▸ "
				style = highlight
			}

			policy := done.Render("⚡ batchable")
			if !t.Batchable {
				policy = orange.Render("🔒 sequential")
			}

			toolCall := fmt.Sprintf("%s(%s)", t.Name, t.Args)
			lines = append(lines, fmt.Sprintf("  %s%-28s %s", prefix, style.Render(toolCall), policy))
		}

		lines = append(lines, "")

		// ── Show computed plan ──
		lines = append(lines, accent.Render("  Execution Plan:"))
		for i, group := range m.plan {
			roundLabel := fmt.Sprintf("  Round %d:", i+1)
			if len(group) > 1 {
				roundLabel = done.Render(roundLabel)
			} else {
				roundLabel = yellow.Render(roundLabel)
			}

			var names []string
			for _, idx := range group {
				t := m.tools[idx]
				names = append(names, fmt.Sprintf("%s(%s)", t.Name, t.Args))
			}

			parallel := ""
			if len(group) > 1 {
				parallel = done.Render(fmt.Sprintf(" [%d in parallel]", len(group)))
			} else {
				parallel = dim.Render(" [alone]")
			}

			lines = append(lines, fmt.Sprintf("  %s %s%s",
				roundLabel, dim.Render(strings.Join(names, " + ")), parallel))
		}

		lines = append(lines, "")
		lines = append(lines, dim.Render("  [↑/↓] Navigate  [Enter/Space] Toggle policy  [e] Execute  [r] Reset"))

	} else {
		// ── Execute / Done phase ──
		lines = append(lines, accent.Render("  Executing plan:"))
		lines = append(lines, "")

		for i, group := range m.plan {
			isActive := m.phase == 1 && i == m.execStep
			isDone := i < m.execStep || m.phase == 2

			roundStyle := pending
			icon := "○"
			if isDone {
				roundStyle = done
				icon = "✓"
			} else if isActive {
				roundStyle = yellow
				icon = "▸"
			}

			label := fmt.Sprintf("Round %d", i+1)
			if len(group) > 1 {
				label += fmt.Sprintf(" (%d parallel)", len(group))
			}
			lines = append(lines, fmt.Sprintf("  %s %s",
				roundStyle.Render(icon), roundStyle.Render(label)))

			for j, idx := range group {
				t := m.tools[idx]
				callStyle := pending
				callIcon := "○"
				if t.Done {
					callStyle = done
					callIcon = "✓"
				} else if isActive {
					callStyle = yellow
					callIcon = "◉"
				}

				bracket := "├"
				if j == len(group)-1 {
					bracket = "└"
				}

				toolCall := fmt.Sprintf("%s(%s)", t.Name, t.Args)
				result := ""
				if t.Done {
					result = " → " + done.Render(t.Result)
				}
				lines = append(lines, fmt.Sprintf("    %s %s %s%s",
					dim.Render(bracket), callStyle.Render(callIcon),
					callStyle.Render(toolCall), result))
			}
			lines = append(lines, "")
		}

		if m.phase == 2 {
			// All done — show savings
			seqTrips := len(m.tools) // if everything were sequential
			lines = append(lines, done.Render(fmt.Sprintf("  Done! %d round trips (vs %d if all sequential)",
				len(m.plan), seqTrips)))

			saved := seqTrips - len(m.plan)
			if saved > 0 {
				lines = append(lines, done.Render(fmt.Sprintf("  Saved %d round trips by batching!", saved)))
			} else {
				lines = append(lines, redStyle.Render("  No savings — try enabling batch on more tools!"))
			}
			lines = append(lines, "")
			lines = append(lines, dim.Render("  [r] Reset and reconfigure"))
		} else {
			lines = append(lines, highlight.Render("  [Enter/Space] Execute next round  ")+dim.Render("[r] Reset"))
		}
	}

	return strings.Join(lines, "\n")
}
