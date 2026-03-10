package lesson

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/quiz"
	"github.com/naorpeled/aitutor/internal/ui"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

// Phase represents which phase of a lesson the user is in.
type Phase int

const (
	PhaseTheory Phase = iota
	PhaseViz
	PhaseQuiz
	PhaseComplete
)

func (p Phase) String() string {
	switch p {
	case PhaseTheory:
		return "Theory"
	case PhaseViz:
		return "Visualization"
	case PhaseQuiz:
		return "Quiz"
	case PhaseComplete:
		return "Complete"
	default:
		return ""
	}
}

// Model manages the state machine for a single lesson.
type Model struct {
	Def      types.LessonDef
	Phase    Phase
	IsLast   bool
	viewport viewport.Model
	width    int
	height   int
	ready    bool

	vizModel  viz.Model
	quizModel quiz.Model
}

func New(def types.LessonDef, width, height int) Model {
	vp := viewport.New(width, height-2)
	vp.SetContent(RenderTheory(def.Theory, width-2))

	m := Model{
		Def:      def,
		Phase:    PhaseTheory,
		viewport: vp,
		width:    width,
		height:   height,
		ready:    true,
	}

	// Build viz if available
	if def.VizBuilder != nil {
		if vm, ok := def.VizBuilder(width, height).(viz.Model); ok {
			m.vizModel = vm
		}
	}

	// Build quiz if questions available
	if len(def.Questions) > 0 {
		m.quizModel = quiz.New(def.Questions)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2
		m.viewport.SetContent(RenderTheory(m.Def.Theory, msg.Width-2))
		return m, nil
	}

	switch m.Phase {
	case PhaseTheory:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case PhaseViz:
		if m.vizModel != nil {
			var cmd tea.Cmd
			m.vizModel, cmd = m.vizModel.Update(msg)
			return m, cmd
		}
	case PhaseQuiz:
		var cmd tea.Cmd
		m.quizModel, cmd = m.quizModel.Update(msg)
		// Auto-advance when quiz is complete
		if m.quizModel.Done() {
			if kmsg, ok := msg.(tea.KeyMsg); ok && kmsg.String() == "enter" {
				m.Phase = PhaseComplete
			}
		}
		return m, cmd
	}

	return m, nil
}

// Advance moves to the next phase. Returns true if lesson is now complete.
func (m *Model) Advance() bool {
	switch m.Phase {
	case PhaseTheory:
		if m.vizModel != nil {
			m.Phase = PhaseViz
		} else if len(m.Def.Questions) > 0 {
			m.Phase = PhaseQuiz
		} else {
			m.Phase = PhaseComplete
			return true
		}
	case PhaseViz:
		if len(m.Def.Questions) > 0 {
			m.Phase = PhaseQuiz
		} else {
			m.Phase = PhaseComplete
			return true
		}
	case PhaseQuiz:
		m.Phase = PhaseComplete
		return true
	}
	return false
}

// GoBack moves to the previous phase.
func (m *Model) GoBack() {
	switch m.Phase {
	case PhaseViz:
		m.Phase = PhaseTheory
	case PhaseQuiz:
		if m.vizModel != nil {
			m.Phase = PhaseViz
		} else {
			m.Phase = PhaseTheory
		}
	case PhaseComplete:
		if len(m.Def.Questions) > 0 {
			m.Phase = PhaseQuiz
		} else if m.vizModel != nil {
			m.Phase = PhaseViz
		} else {
			m.Phase = PhaseTheory
		}
	}
}

func (m Model) View() string {
	phaseIndicator := lipgloss.NewStyle().
		Foreground(ui.ColorAccent).
		Bold(true).
		Render("─── " + m.Phase.String() + " ───")

	var content string
	switch m.Phase {
	case PhaseTheory:
		content = m.viewport.View()
	case PhaseViz:
		if m.vizModel != nil {
			content = m.vizModel.View()
		} else {
			content = "No visualization available"
		}
	case PhaseQuiz:
		content = m.quizModel.View()
	case PhaseComplete:
		msg := "✓ Lesson Complete!\n\nPress → or n for next lesson\nPress ← or p for previous lesson"
		if m.IsLast {
			msg = "✓ Lesson Complete!\n\nYou've finished the last lesson!\nPress ← or p to revisit previous lessons"
		}
		content = lipgloss.NewStyle().
			Foreground(ui.ColorBeginner).
			Bold(true).
			Align(lipgloss.Center).
			Width(m.width).
			Render(msg)
	}

	return lipgloss.JoinVertical(lipgloss.Left, phaseIndicator, "", content)
}
