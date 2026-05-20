package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/naorpeled/aitutor/internal/ui"
)

type agentBuildingScenario struct {
	Prompt      string
	Answer      string
	Explanation string
}

// AgentBuildingModel teaches when to compact context and when to save memory.
type AgentBuildingModel struct {
	scenarios []agentBuildingScenario
	current   int
	choice    string
	answered  bool
	score     int
}

func NewAgentBuildingModel(w, h int) Model {
	return &AgentBuildingModel{
		scenarios: []agentBuildingScenario{
			{
				Prompt:      "The agent has read 20 files, tests are green, and only the final PR body remains.",
				Answer:      "compact",
				Explanation: "Compaction preserves the useful state and frees context before the next phase.",
			},
			{
				Prompt:      "The user says this repo always uses table-driven tests and make test-unit.",
				Answer:      "memory",
				Explanation: "Stable project conventions belong in long-term memory so future sessions benefit.",
			},
			{
				Prompt:      "The current stack trace points at parser.go:42 during this debugging session.",
				Answer:      "continue",
				Explanation: "This is short-lived task context. Keep it in the active window while debugging.",
			},
			{
				Prompt:      "The conversation is near the context limit, but root cause and next steps are clear.",
				Answer:      "compact",
				Explanation: "A compact summary should carry root cause, files touched, commands run, and next steps.",
			},
			{
				Prompt:      "The user prefers concise final answers across all projects.",
				Answer:      "memory",
				Explanation: "Durable user preferences are good long-term memory candidates.",
			},
		},
	}
}

func (m *AgentBuildingModel) Init() tea.Cmd { return nil }

func (m *AgentBuildingModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.current >= len(m.scenarios) {
			if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
				m.current = 0
				m.choice = ""
				m.answered = false
				m.score = 0
			}
			return m, nil
		}

		if m.answered {
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))) {
				m.current++
				m.choice = ""
				m.answered = false
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("1", "c"))):
			m.submit("compact")
		case key.Matches(msg, key.NewBinding(key.WithKeys("2", "m"))):
			m.submit("memory")
		case key.Matches(msg, key.NewBinding(key.WithKeys("3", "w"))):
			m.submit("continue")
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.current = 0
			m.choice = ""
			m.answered = false
			m.score = 0
		}
	}
	return m, nil
}

func (m *AgentBuildingModel) submit(choice string) {
	m.choice = choice
	m.answered = true
	if choice == m.scenarios[m.current].Answer {
		m.score++
	}
}

func (m *AgentBuildingModel) View() string {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, ui.AccentStyle.Render("  Agent Builder Decisions"))
	lines = append(lines, ui.MutedStyle.Render("  Choose what the agent should do with context and memory."))
	lines = append(lines, "")

	lines = append(lines, ui.MutedStyle.Render("  Active context"))
	lines = append(lines, ui.IntermediateStyle.Render("  [facts, files, trace, plan] -> compact when large but still useful"))
	lines = append(lines, "")
	lines = append(lines, ui.MutedStyle.Render("  Long-term memory"))
	lines = append(lines, ui.BeginnerStyle.Render("  [stable preferences, conventions, project facts] -> save for later"))
	lines = append(lines, "")

	if m.current >= len(m.scenarios) {
		lines = append(lines, ui.CorrectStyle.Render(fmt.Sprintf("  Exercise complete. Score: %d/%d", m.score, len(m.scenarios))))
		lines = append(lines, "", ui.MutedStyle.Render("  [r] Restart"))
		return strings.Join(lines, "\n")
	}

	scenario := m.scenarios[m.current]
	lines = append(lines, ui.MutedStyle.Render(fmt.Sprintf("  Scenario %d of %d", m.current+1, len(m.scenarios))))
	lines = append(lines, ui.BrightTextStyle.Render("  "+scenario.Prompt))
	lines = append(lines, "")

	if m.answered {
		if m.choice == scenario.Answer {
			lines = append(lines, ui.CorrectStyle.Render("  Correct."))
		} else {
			lines = append(lines, ui.IncorrectStyle.Render("  Not quite."))
			lines = append(lines, ui.HighlightStyle.Render("  Better answer: "+scenario.Answer))
		}
		lines = append(lines, ui.MutedStyle.Render("  "+scenario.Explanation))
		lines = append(lines, "", ui.HighlightStyle.Render("  Press Enter to continue"))
	} else {
		lines = append(lines, ui.IntermediateStyle.Render("  [1/c] Compact context"))
		lines = append(lines, ui.BeginnerStyle.Render("  [2/m] Save long-term memory"))
		lines = append(lines, ui.HighlightStyle.Render("  [3/w] Keep working in active context"))
	}

	lines = append(lines, "", ui.MutedStyle.Render("  [1/c] Compact  [2/m] Memory  [3/w] Work  [r] Restart"))
	return strings.Join(lines, "\n")
}
