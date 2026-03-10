package quiz

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/pkg/types"
)

// QuizCompleteMsg is sent when all quiz questions are answered.
type QuizCompleteMsg struct {
	Score int
	Total int
}

// Model orchestrates a sequence of quiz questions.
type Model struct {
	questions []types.QuizQuestion
	current   int
	score     int
	done      bool

	// Current question models
	mcModel MultipleChoiceModel
	fbModel FillBlankModel
	orModel OrderingModel
}

func New(questions []types.QuizQuestion) Model {
	m := Model{questions: questions}
	if len(questions) > 0 {
		m.loadQuestion(0)
	}
	return m
}

func (m *Model) loadQuestion(idx int) {
	if idx >= len(m.questions) {
		m.done = true
		return
	}
	m.current = idx
	q := m.questions[idx]
	switch q.Kind {
	case types.MultipleChoice:
		m.mcModel = NewMultipleChoice(q)
	case types.FillBlank:
		m.fbModel = NewFillBlank(q)
	case types.Ordering:
		m.orModel = NewOrdering(q)
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.done || len(m.questions) == 0 {
		return m, nil
	}

	q := m.questions[m.current]

	switch q.Kind {
	case types.MultipleChoice:
		wasAnswered := m.mcModel.IsAnswered()
		var cmd tea.Cmd
		m.mcModel, cmd = m.mcModel.Update(msg)
		if !wasAnswered && m.mcModel.IsAnswered() {
			if m.mcModel.IsCorrect() {
				m.score++
			}
			return m, cmd
		}
		// If already answered, Enter advances to next question
		if wasAnswered {
			if kmsg, ok := msg.(tea.KeyMsg); ok && kmsg.String() == "enter" {
				m.advanceQuestion()
			}
		}
		return m, cmd

	case types.FillBlank:
		wasAnswered := m.fbModel.IsAnswered()
		var cmd tea.Cmd
		m.fbModel, cmd = m.fbModel.Update(msg)
		if !wasAnswered && m.fbModel.IsAnswered() {
			if m.fbModel.IsCorrect() {
				m.score++
			}
			return m, cmd
		}
		if wasAnswered {
			if kmsg, ok := msg.(tea.KeyMsg); ok && kmsg.String() == "enter" {
				m.advanceQuestion()
			}
		}
		return m, cmd

	case types.Ordering:
		wasAnswered := m.orModel.IsAnswered()
		var cmd tea.Cmd
		m.orModel, cmd = m.orModel.Update(msg)
		if !wasAnswered && m.orModel.IsAnswered() {
			if m.orModel.IsCorrect() {
				m.score++
			}
			return m, cmd
		}
		if wasAnswered {
			if kmsg, ok := msg.(tea.KeyMsg); ok && kmsg.String() == "enter" {
				m.advanceQuestion()
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m *Model) advanceQuestion() {
	m.loadQuestion(m.current + 1)
}

func (m Model) Done() bool { return m.done }
func (m Model) Score() int { return m.score }
func (m Model) Total() int { return len(m.questions) }

func (m Model) View() string {
	if len(m.questions) == 0 {
		return "  No quiz questions for this lesson."
	}

	if m.done {
		scoreColor := lipgloss.Color("#4ade80")
		if m.score < len(m.questions) {
			scoreColor = lipgloss.Color("#facc15")
		}
		if m.score == 0 {
			scoreColor = lipgloss.Color("#f87171")
		}

		scoreStyle := lipgloss.NewStyle().Foreground(scoreColor).Bold(true)
		dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))

		var lines []string
		lines = append(lines, "")
		lines = append(lines, scoreStyle.Render(fmt.Sprintf("  Quiz Complete! Score: %d/%d", m.score, len(m.questions))))
		lines = append(lines, "")
		if m.score == len(m.questions) {
			lines = append(lines, scoreStyle.Render("  Perfect score! 🎉"))
		} else {
			lines = append(lines, dim.Render("  Press Enter to mark lesson complete"))
		}
		return strings.Join(lines, "\n")
	}

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6b7280")).
		Render(fmt.Sprintf("  Question %d of %d", m.current+1, len(m.questions)))

	var questionView string
	q := m.questions[m.current]
	switch q.Kind {
	case types.MultipleChoice:
		questionView = m.mcModel.View()
	case types.FillBlank:
		questionView = m.fbModel.View()
	case types.Ordering:
		questionView = m.orModel.View()
	}

	return header + "\n\n" + questionView
}
