# Contributing to AITutor

Thanks for your interest in contributing to AITutor! This guide covers everything you need to get started.

## Getting Started

### Prerequisites

- Go 1.23+
- Make (optional, for convenience commands)

### Setup

```bash
git clone https://github.com/naorpeled/aitutor.git
cd aitutor
make build
```

### Build & Verify

```bash
make build    # builds ./aitutor binary
make run      # go run .
make vet      # go vet ./...
make install  # go install .
```

After any change, always run:

```bash
go build ./...
go vet ./...
```

Both must pass with no errors before submitting a PR.

## Architecture Overview

AITutor is built with Go and the [Charm](https://charm.sh) ecosystem (Bubbletea, Lipgloss, Bubbles).

```
main.go                        # Entry point — blank-imports content packages
internal/
  app/app.go                   # Root Bubbletea model (header, footer, sidebar, lessons)
  lesson/model.go              # Lesson state machine: Theory → Viz → Quiz → Complete
  content/
    beginner/                  # Lessons 1-4
    intermediate/              # Lessons 5-8, 15
    advanced/                  # Lessons 9-14
  viz/                         # Interactive visualizations (each implements viz.Model)
  quiz/                        # Quiz engine (MultipleChoice, FillBlank, Ordering)
  progress/                    # Progress persistence (~/.aitutor/progress.json)
  ui/                          # Styles, header, footer, sidebar
pkg/types/types.go             # Shared types (LessonDef, TheoryBlock, QuizQuestion)
```

### How Lessons Work

Each lesson file calls `lesson.Register()` in an `init()` function. The `main.go` blank-imports all content packages, which triggers registration automatically — no wiring needed.

Lessons have four phases: **Theory** (scrollable content) → **Visualization** (interactive ASCII viz) → **Quiz** (questions) → **Complete**.

## Adding a New Lesson

### 1. Create the Visualization

Create `internal/viz/<name>.go` implementing the `viz.Model` interface:

```go
type Model interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Model, tea.Cmd)
    View() string
}
```

Visualization rules:
- Use `ui.Color*` constants from `internal/ui/styles.go` — never inline hex values
- Use `Enter`/`Space` for primary interaction, `r` for reset
- **Do NOT use** `Tab` (reserved for sidebar toggle) or left/right arrows (reserved for phase navigation)
- Use `h`/`l` for horizontal navigation, `j`/`k` or up/down for vertical navigation
- Use helper functions: `viz.Box()`, `viz.Arrow()`, `viz.HLine()`, `viz.CenterText()`
- Use distinct type names to avoid collisions with other files in the `viz` package
- Always show key hints at the bottom of `View()`

### 2. Create the Lesson File

Create `internal/content/<tier>/NN_topic.go`:

```go
package beginner // or intermediate, advanced

import (
    "github.com/naorpeled/aitutor/internal/lesson"
    "github.com/naorpeled/aitutor/internal/viz"
    "github.com/naorpeled/aitutor/pkg/types"
)

func init() {
    lesson.Register(types.LessonDef{
        ID:      16,
        Title:   "Your Lesson Title",
        Tier:    types.Beginner,
        Summary: "Brief description for sidebar",
        VizBuilder: func(w, h int) interface{} {
            return viz.NewYourVizModel(w, h)
        },
        Theory: []types.TheoryBlock{
            {Kind: types.Heading, Content: "Main Heading"},
            {Kind: types.Paragraph, Content: "Explanation text."},
            {Kind: types.Code, Content: "  code example"},
            {Kind: types.Bullet, Content: "Point one\nPoint two\nPoint three"},
            {Kind: types.Callout, Content: "Key takeaway or link"},
        },
        Questions: []types.QuizQuestion{
            {
                Kind:        types.MultipleChoice,
                Prompt:      "Question?",
                Choices:     []string{"A", "B", "C", "D"},
                CorrectIdx:  1,
                Explanation: "Why B is correct.",
            },
        },
    })
}
```

### 3. Update the Header

In `internal/ui/header.go`, update the `Total` default in `NewHeaderModel()` to reflect the new lesson count.

### 4. Test It

```bash
go build ./...
go vet ./...
./aitutor  # navigate to your lesson, test all phases
```

## Content Guidelines

- **Vendor-neutral**: Do not reference specific AI products. Use generic terms like "AI coding assistants", "models", "LLMs".
- **Interactive visualizations required**: Every lesson must have an interactive viz — no static-only diagrams.
- **Mix quiz types**: Include 2-3 questions using different types (MultipleChoice, FillBlank, Ordering).
- **External references**: Add `Callout` blocks with links to Wikipedia, official docs, or other educational resources.
- **Keyed struct fields**: Always use `{Key: value}` syntax — `go vet` enforces this.

## Code Conventions

- Go standard formatting (`gofmt`)
- No external dependencies beyond the Charm ecosystem (Bubbletea, Lipgloss, Bubbles)
- Lesson content is Go code (not YAML/Markdown) for type safety and compile-time checks
- `VizBuilder` returns `interface{}` (not `viz.Model`) to avoid circular imports — the `lesson` package does the type assertion

## AI-Assisted Contributing with Claude Code

This project includes two Claude Code skills for AI-assisted lesson development:

### `/add-chapter` — Create a New Lesson

Invoke this skill when asking Claude Code to create a new lesson. It guides through the full workflow: choosing the lesson ID/tier, creating the visualization, writing the lesson file, updating the header, and verifying the build.

### `/edit-chapter` — Modify an Existing Lesson

Invoke this skill when asking Claude Code to edit theory content, fix quiz questions, improve visualizations, or restructure existing lessons.

These skills encode all the project's patterns and constraints, so Claude Code will follow the conventions automatically.

## Testing Changes

1. `go build ./...` — must compile without errors
2. `go vet ./...` — must pass with no warnings
3. Run the app and navigate through affected lessons to verify visualizations work
4. Test phase navigation: `→` to advance, `←` to go back through Theory → Viz → Quiz → Complete

## Generating the Demo GIF

```bash
brew install vhs
vhs demo.tape
```

This produces `demo.gif` used in the README.

## Submitting a PR

1. Fork the repo and create a branch
2. Make your changes following the guidelines above
3. Ensure `go build ./...` and `go vet ./...` pass
4. Test your changes manually in the terminal
5. Submit a PR with a clear description of what you added or changed
