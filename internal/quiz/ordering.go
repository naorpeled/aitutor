package quiz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/pkg/types"
)

// OrderingModel handles an ordering/sequencing question.
type OrderingModel struct {
	Question types.QuizQuestion
	order    []int
	cursor   int
	answered bool
	correct  bool
}

func NewOrdering(q types.QuizQuestion) OrderingModel {
	order := make([]int, len(q.Choices))
	for i := range order {
		order[i] = i
	}
	return OrderingModel{Question: q, order: order}
}

func (m OrderingModel) Init() tea.Cmd { return nil }

func (m OrderingModel) Update(msg tea.Msg) (OrderingModel, tea.Cmd) {
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
			if m.cursor < len(m.order)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("K"))):
			// Move item up
			if m.cursor > 0 {
				m.order[m.cursor], m.order[m.cursor-1] = m.order[m.cursor-1], m.order[m.cursor]
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("J"))):
			// Move item down
			if m.cursor < len(m.order)-1 {
				m.order[m.cursor], m.order[m.cursor+1] = m.order[m.cursor+1], m.order[m.cursor]
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			m.answered = true
			m.correct = m.checkOrder()
		}
	}
	return m, nil
}

func (m OrderingModel) checkOrder() bool {
	for i, idx := range m.order {
		if idx != i {
			return false
		}
	}
	return true
}

func (m OrderingModel) IsAnswered() bool { return m.answered }
func (m OrderingModel) IsCorrect() bool  { return m.correct }

func (m OrderingModel) View() string {
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9fafb"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d5db"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#38bdf8")).Bold(true)
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))

	var lines []string
	lines = append(lines, promptStyle.Render("  "+m.Question.Prompt))
	lines = append(lines, dim.Render("  (↑/↓ navigate, Shift+J/K to reorder, Enter to submit)"))
	lines = append(lines, "")

	for i, idx := range m.order {
		style := itemStyle
		prefix := fmt.Sprintf("  %d. ", i+1)
		if i == m.cursor && !m.answered {
			style = selectedStyle
			prefix = fmt.Sprintf("  %d. ▸ ", i+1)
		}
		lines = append(lines, prefix+style.Render(m.Question.Choices[idx]))
	}

	if m.answered {
		lines = append(lines, "")
		if m.correct {
			lines = append(lines, RenderCorrect(m.Question.Explanation))
		} else {
			lines = append(lines, RenderIncorrect(m.Question.Explanation))
			lines = append(lines, "")
			correct := lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80"))
			lines = append(lines, correct.Render("  Correct order:"))
			for i, choice := range m.Question.Choices {
				lines = append(lines, correct.Render(fmt.Sprintf("    %d. %s", i+1, choice)))
			}
		}
	}

	return strings.Join(lines, "\n")
}
