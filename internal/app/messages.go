package app

import tea "github.com/charmbracelet/bubbletea"

// LessonChangeMsg is sent when the user navigates to a different lesson.
type LessonChangeMsg struct {
	LessonIdx int
}

// PhaseAdvanceMsg is sent when a lesson phase completes.
type PhaseAdvanceMsg struct{}

// PhaseBackMsg is sent when user goes back a phase.
type PhaseBackMsg struct{}

// ToggleSidebarMsg is sent when user presses Tab.
type ToggleSidebarMsg struct{}

// ShowHelpMsg toggles the help overlay.
type ShowHelpMsg struct{}

// LessonCompleteMsg signals a lesson has been fully completed.
type LessonCompleteMsg struct {
	LessonID int
}

// WindowSizeMsg wraps tea.WindowSizeMsg for convenience.
type WindowSizeMsg = tea.WindowSizeMsg
