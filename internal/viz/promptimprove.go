package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type promptChallenge struct {
	BadPrompt   string
	Options     []string
	CorrectIdx  int
	Explanation string
}

// PromptImproveModel is an interactive exercise where users pick better prompts.
type PromptImproveModel struct {
	width      int
	height     int
	challenges []promptChallenge
	current    int
	cursor     int
	answered   bool
	correct    bool
	score      int
}

func NewPromptImproveModel(w, h int) Model {
	return &PromptImproveModel{
		width:  w,
		height: h,
		challenges: []promptChallenge{
			{
				BadPrompt: "fix the bug",
				Options: []string{
					"fix all the bugs in the project",
					"Fix the NullPointerException in UserService.getById() when the user ID doesn't exist in the database",
					"debug the code and make it work",
					"the code is broken, please help",
				},
				CorrectIdx:  1,
				Explanation: "Specific: names the error, the method, and the condition that triggers it.",
			},
			{
				BadPrompt: "add tests",
				Options: []string{
					"write some tests for the code",
					"test everything",
					"Add unit tests for ParseConfig covering: valid YAML input, missing required 'port' field, and port values outside 1-65535",
					"make sure the code works",
				},
				CorrectIdx:  2,
				Explanation: "Names the function, lists specific test cases, and defines edge cases.",
			},
			{
				BadPrompt: "make it faster",
				Options: []string{
					"optimize everything for speed",
					"The /api/users endpoint takes 3s to respond. Profile the SQL query in UserRepository.findAll() — it's doing N+1 queries. Use a JOIN instead.",
					"performance is bad, fix it",
					"use caching to make it fast",
				},
				CorrectIdx:  1,
				Explanation: "Identifies the endpoint, the metric, the root cause, and suggests a specific approach.",
			},
		},
	}
}

func (m *PromptImproveModel) Init() tea.Cmd { return nil }

func (m *PromptImproveModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.answered {
			if key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))) {
				m.current++
				m.cursor = 0
				m.answered = false
				if m.current >= len(m.challenges) {
					// Done — stay on last
					m.current = len(m.challenges)
				}
			}
			return m, nil
		}

		if m.current >= len(m.challenges) {
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			c := m.challenges[m.current]
			if m.cursor < len(c.Options)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			m.cursor = 0
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			m.cursor = 1
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("3"))):
			m.cursor = 2
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("4"))):
			if len(m.challenges[m.current].Options) > 3 {
				m.cursor = 3
				m.submit()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			m.submit()
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.current = 0
			m.cursor = 0
			m.answered = false
			m.score = 0
		}
	}
	return m, nil
}

func (m *PromptImproveModel) submit() {
	c := m.challenges[m.current]
	m.answered = true
	m.correct = m.cursor == c.CorrectIdx
	if m.correct {
		m.score++
	}
}

func (m *PromptImproveModel) View() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	bad := lipgloss.NewStyle().Foreground(ui.ColorIncorrect).Bold(true)
	good := lipgloss.NewStyle().Foreground(ui.ColorCorrect).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	explain := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Prompt Improvement Exercise"))
	lines = append(lines, dim.Render("  Pick the BEST replacement for each vague prompt"))
	lines = append(lines, "")

	if m.current >= len(m.challenges) {
		// Done
		lines = append(lines, good.Render(fmt.Sprintf("  Exercise Complete! Score: %d/%d", m.score, len(m.challenges))))
		lines = append(lines, "")
		if m.score == len(m.challenges) {
			lines = append(lines, good.Render("  Perfect! You know how to write effective prompts."))
		} else {
			lines = append(lines, text.Render("  Remember: be specific, provide context, name the problem."))
		}
		lines = append(lines, "", dim.Render("  [r] Try again"))
		return strings.Join(lines, "\n")
	}

	c := m.challenges[m.current]
	lines = append(lines, dim.Render(fmt.Sprintf("  Challenge %d of %d", m.current+1, len(m.challenges))))
	lines = append(lines, "")
	lines = append(lines, bad.Render("  Vague prompt: ")+text.Render("\""+c.BadPrompt+"\""))
	lines = append(lines, "")
	lines = append(lines, highlight.Render("  Which is the best replacement?"))
	lines = append(lines, "")

	for i, opt := range c.Options {
		prefix := fmt.Sprintf("  %d) ", i+1)
		style := text

		if m.answered {
			if i == c.CorrectIdx {
				prefix = good.Render(fmt.Sprintf("  %d) ✓ ", i+1))
				style = good
			} else if i == m.cursor && !m.correct {
				prefix = bad.Render(fmt.Sprintf("  %d) ✗ ", i+1))
				style = bad
			}
		} else if i == m.cursor {
			prefix = highlight.Render(fmt.Sprintf("  %d) ▸ ", i+1))
			style = highlight
		}

		// Wrap long options
		wrapped := opt
		if len(wrapped) > m.width-10 {
			wrapped = wrapped[:m.width-13] + "..."
		}
		lines = append(lines, prefix+style.Render(wrapped))
	}

	if m.answered {
		lines = append(lines, "")
		if m.correct {
			lines = append(lines, good.Render("  ✓ Correct!"))
		} else {
			lines = append(lines, bad.Render("  ✗ Not quite"))
		}
		lines = append(lines, explain.Render("  "+c.Explanation))
		lines = append(lines, "", highlight.Render("  Press Enter to continue"))
	}

	lines = append(lines, "", dim.Render("  [↑/↓] Navigate  [1-4] Select  [Enter] Confirm  [r] Restart"))

	return strings.Join(lines, "\n")
}
