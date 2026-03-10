package quiz

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/pkg/types"
)

// FillBlankModel handles a fill-in-the-blank question.
type FillBlankModel struct {
	Question types.QuizQuestion
	input    textinput.Model
	answered bool
	correct  bool
}

func NewFillBlank(q types.QuizQuestion) FillBlankModel {
	ti := textinput.New()
	ti.Placeholder = "Type your answer..."
	ti.Focus()
	ti.Width = 40
	return FillBlankModel{Question: q, input: ti}
}

func (m FillBlankModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FillBlankModel) Update(msg tea.Msg) (FillBlankModel, tea.Cmd) {
	if m.answered {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("enter"))) {
			m.answered = true
			answer := strings.TrimSpace(strings.ToLower(m.input.Value()))
			expected := strings.TrimSpace(strings.ToLower(m.Question.Answer))
			m.correct = answer == expected
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m FillBlankModel) IsAnswered() bool { return m.answered }
func (m FillBlankModel) IsCorrect() bool  { return m.correct }

func (m FillBlankModel) View() string {
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9fafb"))

	var lines []string
	lines = append(lines, promptStyle.Render("  "+m.Question.Prompt))
	lines = append(lines, "")
	lines = append(lines, "  "+m.input.View())

	if m.answered {
		lines = append(lines, "")
		if m.correct {
			lines = append(lines, RenderCorrect(m.Question.Explanation))
		} else {
			expected := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80")).Bold(true).
				Render("  Answer: " + m.Question.Answer)
			lines = append(lines, RenderIncorrect(m.Question.Explanation))
			lines = append(lines, expected)
		}
	}

	return strings.Join(lines, "\n")
}
