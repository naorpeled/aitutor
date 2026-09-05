# AGENTS.md

## Project Overview

AITutor is an interactive terminal-based tutorial for AI coding concepts, built with Go and the Charm ecosystem (Bubbletea, Lipgloss, Bubbles). It teaches through theory, interactive visualizations, and quizzes across 17 lessons in 3 tiers (Beginner, Intermediate, Advanced).

## Architecture

- **Entry point:** `main.go` — blank-imports all content packages to trigger `init()` lesson registration
- **App model:** `internal/app/app.go` — root Bubbletea model orchestrating header, footer, sidebar, lesson model, progress, welcome screen, and help overlay
- **Lesson state machine:** `internal/lesson/model.go` — phases: Theory → Visualization → Quiz → Complete
- **Content:** `internal/content/{beginner,intermediate,advanced}/` — each file registers a `LessonDef` via `init()`
- **Visualizations:** `internal/viz/` — each viz implements the `viz.Model` interface (Init/Update/View)
- **Quiz:** `internal/quiz/` — supports MultipleChoice, FillBlank, and Ordering question types
- **Progress:** `internal/progress/` — persists to `~/.aitutor/progress.json`
- **Types:** `pkg/types/types.go` — shared types to avoid circular imports

## Key Patterns

### Adding a New Lesson

1. Create `internal/content/<tier>/NN_topic.go`
2. In `init()`, call `lesson.Register(types.LessonDef{...})` with ID, Title, Tier, Theory, VizBuilder, and Questions
3. Create the visualization in `internal/viz/` implementing `viz.Model` interface
4. The lesson auto-registers — no wiring needed
5. Update `internal/ui/header.go` default Total if needed (it gets overridden at runtime but keeps the default accurate)

### VizBuilder Uses `interface{}`

`VizBuilder` returns `interface{}` (not `viz.Model`) to avoid circular imports between `types` and `viz` packages. The `lesson` package does the type assertion.

### Key Handling

- App-level keys (quit, tab, n/p, →/← for phase navigation, ?) are handled in `app.go`
- `→` (right arrow) advances phases from any phase — this is the primary "next phase" key
- `←` (left arrow) and Backspace go back a phase
- `Enter` advances from Theory phase, but in Viz/Quiz phases it falls through to the viz/quiz for interaction
- Tab is reserved for sidebar toggle — visualizations must not use Tab
- Visualizations should use `h`/`l` for horizontal navigation, not left/right arrows (reserved for phase nav)

### Styling

All styles live in `internal/ui/styles.go`. Use the defined color constants (`ColorAccent`, `ColorBeginner`, etc.) rather than inline hex values.

### Type Naming

Avoid type name collisions across files in the same package. For example, `internal/viz/` has both `bucket.go` and `mcpcaller.go` — use distinct names like `mcpToolDef` vs `mcpTool` if both files need similar types.

## Code Conventions

- Go standard formatting (`gofmt`)
- No external dependencies beyond the Charm ecosystem
- Use keyed struct fields (e.g., `ui.KeyHint{Key: "q", Desc: "quit"}`) — `go vet` enforces this
- Lesson content is Go code (not YAML/Markdown) for type safety and compile-time checks
- Interactive visualizations in every lesson — no static-only diagrams
- Keep content vendor-neutral — avoid referencing specific AI products in lesson theory/examples

## Build & Run

```bash
make build    # builds ./aitutor binary
make run      # go run .
make vet      # go vet ./...
make install  # go install .
```

## Releasing

Releases are automated via GoReleaser and GitHub Actions (`.goreleaser.yml`, `.github/workflows/release.yml`).

To publish a new release, trigger the `Release` workflow via `workflow_dispatch` with the desired version (e.g., `0.1.5`). The workflow:
1. Verifies the version format, then runs `go vet` / `go test` / `go build` — nothing is published if the build is broken
2. Checks npm credentials up front, before any irreversible step
3. Bumps `package.json` version, commits, and pushes
4. Creates and pushes the git tag `vX.Y.Z`, then checks it out so the rest of the release publishes exactly the tagged tree
5. GoReleaser builds cross-platform binaries (macOS/Linux/Windows, amd64/arm64), creates a GitHub Release, and updates the Homebrew formula in `naorpeled/homebrew-tap`
6. Publishes the npm package `@aitutor/cli` with OIDC provenance

### Re-running a failed release

A release is four independent publishes (bump commit, git tag, GitHub Release + Homebrew, npm), so any one of them can fail on its own and leave the others already shipped. The workflow detects what already exists and skips it, so **re-dispatching the same version publishes only the missing pieces** — that is the way to repair a partial release. Do not hand-delete tags or releases to force a clean re-run.

Users install via Homebrew or npm:
```bash
brew tap naorpeled/tap
brew install aitutor
# or
npx @aitutor/cli
```

Secrets required:
- `HOMEBREW_TAP_TOKEN` — grants GoReleaser permission to push formula updates to the tap repo
- `NPM_TOKEN` — authenticates npm publish. Optional: if it is unset, the workflow publishes via npm [trusted publishing](https://docs.npmjs.com/trusted-publishers) using the `id-token` OIDC grant, which needs no secret and cannot expire. Trusted publishing is preferred — npm tokens expire and a silently expired one strands a release (npm answers an unauthorized `PUT` with a `404`, not a `401`).

## Generating the Demo GIF

```bash
brew install vhs
vhs demo.tape
```

This produces `demo.gif` used in the README.

## Testing Changes

After any change:
1. `go build ./...` — must compile without errors
2. `go vet ./...` — must pass with no warnings
3. Run the app and navigate through affected lessons to verify visualizations work
4. Test phase navigation: `→` to advance, `←` to go back through Theory → Viz → Quiz → Complete
