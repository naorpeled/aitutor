package quiz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/pkg/types"
)

// MultipleChoiceModel handles a multiple choice question.
type MultipleChoiceModel struct {
	Question   types.QuizQuestion
	cursor     int
	answered   bool
	correct    bool
}

func NewMultipleChoice(q types.QuizQuestion) MultipleChoiceModel {
	return MultipleChoiceModel{Question: q}
}

func (m MultipleChoiceModel) Init() tea.Cmd { return nil }

func (m MultipleChoiceModel) Update(msg tea.Msg) (MultipleChoiceModel, tea.Cmd) {
	if m.answered {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.Question.Choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			if len(m.Question.Choices) > 0 {
				m.cursor = 0
				m.submit()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			if len(m.Question.Choices) > 1 {
				m.cursor = 1
				m.submit()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("3"))):
			if len(m.Question.Choices) > 2 {
				m.cursor = 2
				m.submit()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("4"))):
			if len(m.Question.Choices) > 3 {
				m.cursor = 3
				m.submit()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			m.submit()
		}
	}
	return m, nil
}

func (m *MultipleChoiceModel) submit() {
	m.answered = true
	m.correct = m.cursor == m.Question.CorrectIdx
}

func (m MultipleChoiceModel) IsAnswered() bool { return m.answered }
func (m MultipleChoiceModel) IsCorrect() bool  { return m.correct }

func (m MultipleChoiceModel) View() string {
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9fafb"))
	choiceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#38bdf8")).Bold(true)
	correctMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80")).Bold(true)
	wrongMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#f87171")).Bold(true)

	var lines []string
	lines = append(lines, promptStyle.Render("  "+m.Question.Prompt))
	lines = append(lines, "")

	for i, choice := range m.Question.Choices {
		prefix := fmt.Sprintf("  %d) ", i+1)
		style := choiceStyle

		if m.answered {
			if i == m.Question.CorrectIdx {
				prefix = correctMark.Render(fmt.Sprintf("  %d) ✓ ", i+1))
				style = correctMark
			} else if i == m.cursor && !m.correct {
				prefix = wrongMark.Render(fmt.Sprintf("  %d) ✗ ", i+1))
				style = wrongMark
			}
		} else if i == m.cursor {
			prefix = selectedStyle.Render(fmt.Sprintf("  %d) ▸ ", i+1))
			style = selectedStyle
		}

		lines = append(lines, prefix+style.Render(choice))
	}

	if m.answered {
		lines = append(lines, "")
		if m.correct {
			lines = append(lines, RenderCorrect(m.Question.Explanation))
		} else {
			lines = append(lines, RenderIncorrect(m.Question.Explanation))
		}
	}

	return strings.Join(lines, "\n")
}
