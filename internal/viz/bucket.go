package viz

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/naorpeled/aitutor/internal/ui"
)

type contextItem struct {
	Name     string
	Tokens   int
	Color    lipgloss.Color
	Char     rune
	Category string // "system", "tool", "conversation", "mcp"
}

type toolOption struct {
	Name     string
	Desc     string
	Tokens   int
	Category string
}

type mcpToolDef struct {
	Name    string
	Enabled bool
}

type mcpServerOption struct {
	Name      string
	ToolCount int // total tools (len of Tools)
	Tokens    int // tokens when all enabled
	Enabled   bool
	Expanded  bool
	Tools     []mcpToolDef
}

// compressionEntry records what happened during a compression step.
type compressionEntry struct {
	Before string
	After  string
	Saved  int
}

// mcpRow represents a row in the flattened MCP list view.
type mcpRow struct {
	isServer   bool
	serverIdx  int
	toolIdx    int // -1 for server rows
}

// BucketModel visualizes the context window as a bucket that fills up.
type BucketModel struct {
	width      int
	height     int
	items      []contextItem
	options    []toolOption
	mcpServers []mcpServerOption
	cursor     int
	capacity   int
	overflow   bool
	tab        int // 0=tools, 1=mcp servers

	// Compression demo state
	compressing    bool
	compressStep   int
	compressLog    []compressionEntry
	preCompressUse int
}

func NewBucketModel(w, h int) Model {
	m := &BucketModel{
		width:    w,
		height:   h,
		capacity: 200000,
		options: []toolOption{
			{Name: "Read small file", Desc: "50-line config", Tokens: 500, Category: "tool"},
			{Name: "Read large file", Desc: "2000-line source", Tokens: 8000, Category: "tool"},
			{Name: "Glob search", Desc: "Find *.go files", Tokens: 200, Category: "tool"},
			{Name: "Grep codebase", Desc: "Pattern search", Tokens: 3000, Category: "tool"},
			{Name: "Bash: run tests", Desc: "Test output", Tokens: 15000, Category: "tool"},
			{Name: "Bash: build", Desc: "Compile errors", Tokens: 5000, Category: "tool"},
			{Name: "Send a message", Desc: "Your reply", Tokens: 1000, Category: "conversation"},
			{Name: "AI response", Desc: "Detailed answer", Tokens: 3000, Category: "conversation"},
			{Name: "Paste error log", Desc: "Long stack trace", Tokens: 12000, Category: "conversation"},
			{Name: "Load huge file", Desc: "10k line file", Tokens: 45000, Category: "tool"},
		},
		mcpServers: makeMCPServers(),
	}
	m.items = []contextItem{
		{Name: "System Prompt", Tokens: 8000, Color: ui.ColorMuted, Char: '▒', Category: "system"},
		{Name: "CLAUDE.md", Tokens: 2000, Color: ui.ColorBeginner, Char: '░', Category: "system"},
	}
	return m
}

func makeMCPServers() []mcpServerOption {
	makeTools := func(names []string) []mcpToolDef {
		tools := make([]mcpToolDef, len(names))
		for i, n := range names {
			tools[i] = mcpToolDef{Name: n, Enabled: false}
		}
		return tools
	}
	servers := []mcpServerOption{
		{Name: "GitHub", Tools: makeTools([]string{
			"create_issue", "list_issues", "get_issue", "update_issue",
			"create_pr", "list_prs", "get_pr", "merge_pr",
			"create_review", "list_reviews", "list_commits", "get_commit",
			"create_branch", "list_branches", "delete_branch",
			"list_repos", "get_repo", "create_repo", "fork_repo",
			"list_releases", "create_release",
			"get_file_contents", "create_or_update_file",
			"search_code", "search_issues", "search_repos",
			"list_labels", "create_label", "add_labels",
			"list_milestones", "create_milestone",
			"list_workflows", "trigger_workflow", "list_workflow_runs",
		})},
		{Name: "Slack", Tools: makeTools([]string{
			"send_message", "read_channel", "list_channels",
			"create_channel", "archive_channel", "set_topic",
			"add_reaction", "remove_reaction", "upload_file",
			"list_users", "get_user", "send_dm",
			"search_messages", "list_reactions", "pin_message",
			"schedule_message", "update_message", "delete_message",
		})},
		{Name: "Jira", Tools: makeTools([]string{
			"create_issue", "get_issue", "update_issue", "delete_issue",
			"search_issues", "assign_issue", "transition_issue",
			"add_comment", "list_comments", "get_comment",
			"create_sprint", "get_sprint", "list_sprints",
			"list_projects", "get_project", "list_boards",
			"add_attachment", "list_attachments",
			"create_filter", "get_filter", "add_watcher", "log_work",
		})},
		{Name: "Linear", Tools: makeTools([]string{
			"create_issue", "get_issue", "update_issue", "list_issues",
			"search_issues", "create_project", "list_projects",
			"create_cycle", "list_cycles", "add_comment",
			"list_teams", "get_team", "create_label",
			"list_labels", "assign_issue", "list_views",
		})},
		{Name: "Filesystem", Tools: makeTools([]string{
			"read_file", "write_file", "list_directory",
			"create_directory", "move_file", "search_files",
			"get_file_info", "read_multiple_files",
			"delete_file", "copy_file", "watch_file",
		})},
		{Name: "Sentry", Tools: makeTools([]string{
			"list_issues", "get_issue", "resolve_issue",
			"list_events", "get_event", "list_projects",
			"get_project", "list_releases", "create_release",
		})},
		{Name: "PostgreSQL", Tools: makeTools([]string{
			"query", "list_tables", "describe_table",
			"list_schemas", "list_databases", "explain_query",
			"list_indexes",
		})},
		{Name: "Datadog", Tools: makeTools([]string{
			"list_monitors", "get_monitor", "create_monitor",
			"mute_monitor", "list_dashboards", "get_dashboard",
			"query_metrics", "list_logs", "search_logs",
			"list_incidents", "get_incident", "list_hosts",
			"get_host",
		})},
	}
	for i := range servers {
		servers[i].ToolCount = len(servers[i].Tools)
		servers[i].Tokens = servers[i].ToolCount * 200
	}
	return servers
}

func (m *BucketModel) Init() tea.Cmd { return nil }

// mcpRows builds a flat list of rows for MCP navigation.
func (m *BucketModel) mcpRows() []mcpRow {
	var rows []mcpRow
	for si, srv := range m.mcpServers {
		rows = append(rows, mcpRow{isServer: true, serverIdx: si, toolIdx: -1})
		if srv.Expanded {
			for ti := range srv.Tools {
				rows = append(rows, mcpRow{isServer: false, serverIdx: si, toolIdx: ti})
			}
		}
	}
	return rows
}

func (m *BucketModel) enabledMCPTokens() int {
	total := 0
	for _, s := range m.mcpServers {
		if !s.Enabled {
			continue
		}
		for _, t := range s.Tools {
			if t.Enabled {
				total += 200
			}
		}
	}
	return total
}

func (m *BucketModel) totalUsed() int {
	total := 0
	for _, item := range m.items {
		total += item.Tokens
	}
	total += m.enabledMCPTokens()
	return total
}

func (m *BucketModel) categoryTotals() (system, tool, conversation, mcp int) {
	for _, item := range m.items {
		switch item.Category {
		case "system":
			system += item.Tokens
		case "tool":
			tool += item.Tokens
		case "conversation":
			conversation += item.Tokens
		}
	}
	mcp = m.enabledMCPTokens()
	return
}

// runCompression simulates the multi-step compression process.
func (m *BucketModel) runCompression() {
	m.compressing = true
	m.compressStep = 0
	m.preCompressUse = m.totalUsed()
	m.compressLog = nil

	// Step 1: Drop old tool results (keep only last 3)
	toolItems := 0
	for i := len(m.items) - 1; i >= 0; i-- {
		if m.items[i].Category == "tool" {
			toolItems++
		}
	}
	if toolItems > 3 {
		dropped := 0
		savedTokens := 0
		var newItems []contextItem
		toolSeen := 0
		// Walk from end to keep most recent tools
		keep := make(map[int]bool)
		for i := len(m.items) - 1; i >= 0; i-- {
			if m.items[i].Category == "tool" {
				toolSeen++
				if toolSeen <= 3 {
					keep[i] = true
				}
			} else {
				keep[i] = true
			}
		}
		for i, item := range m.items {
			if keep[i] {
				newItems = append(newItems, item)
			} else {
				savedTokens += item.Tokens
				dropped++
			}
		}
		if dropped > 0 {
			m.compressLog = append(m.compressLog, compressionEntry{
				Before: fmt.Sprintf("%d tool results in context", toolItems),
				After:  fmt.Sprintf("Dropped %d oldest, kept 3 most recent", dropped),
				Saved:  savedTokens,
			})
			m.items = newItems
		}
	}

	// Step 2: Summarize old conversation messages (combine into summary)
	convItems := 0
	convTokens := 0
	for _, item := range m.items {
		if item.Category == "conversation" {
			convItems++
			convTokens += item.Tokens
		}
	}
	if convItems > 2 {
		// Keep only the last conversation message, summarize the rest
		var newItems []contextItem
		removedTokens := 0

		convCount := 0
		var keptItems []contextItem
		for _, item := range m.items {
			if item.Category == "conversation" {
				convCount++
				if convCount < convItems { // drop all but last
					removedTokens += item.Tokens
					continue
				}
			}
			keptItems = append(keptItems, item)
		}

		// Insert summary where old conversations were
		summaryTokens := removedTokens / 5 // summary is ~20% of original
		summary := contextItem{
			Name:     fmt.Sprintf("Summary (%d msgs)", convItems-1),
			Tokens:   summaryTokens,
			Color:    lipgloss.Color("#c084fc"),
			Char:     '◆',
			Category: "conversation",
		}

		// Insert summary before the last conversation item
		newItems = nil
		insertedSummary := false
		for _, item := range keptItems {
			if item.Category == "conversation" && !insertedSummary {
				newItems = append(newItems, summary)
				insertedSummary = true
			}
			newItems = append(newItems, item)
		}
		if !insertedSummary {
			newItems = append(newItems, summary)
		}

		saved := removedTokens - summaryTokens
		if saved > 0 {
			m.compressLog = append(m.compressLog, compressionEntry{
				Before: fmt.Sprintf("%d conversation messages (%s tokens)", convItems-1, formatTokens(removedTokens)),
				After:  fmt.Sprintf("Compressed to summary (%s tokens)", formatTokens(summaryTokens)),
				Saved:  saved,
			})
			m.items = newItems
		}
	}

	// Step 3: Truncate remaining large tool results
	for i := range m.items {
		if m.items[i].Category == "tool" && m.items[i].Tokens > 10000 {
			oldTokens := m.items[i].Tokens
			m.items[i].Tokens = oldTokens / 3
			saved := oldTokens - m.items[i].Tokens
			m.compressLog = append(m.compressLog, compressionEntry{
				Before: fmt.Sprintf("%s had %s tokens", m.items[i].Name, formatTokens(oldTokens)),
				After:  fmt.Sprintf("Truncated to %s tokens (kept key lines)", formatTokens(m.items[i].Tokens)),
				Saved:  saved,
			})
		}
	}

	m.overflow = m.totalUsed() > m.capacity
}

func (m *BucketModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// In compression view, Enter steps through, any other key exits
		if m.compressing {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
				if m.compressStep < len(m.compressLog) {
					m.compressStep++
				} else {
					m.compressing = false
					m.tab = 0
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
				m.compressing = false
				m.items = []contextItem{
					{Name: "System Prompt", Tokens: 8000, Color: ui.ColorMuted, Char: '▒', Category: "system"},
					{Name: "CLAUDE.md", Tokens: 2000, Color: ui.ColorBeginner, Char: '░', Category: "system"},
				}
				for i := range m.mcpServers {
					m.mcpServers[i].Enabled = false
					m.mcpServers[i].Expanded = false
					for j := range m.mcpServers[i].Tools {
						m.mcpServers[i].Tools[j].Enabled = false
					}
				}
				m.overflow = false
				m.tab = 0
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			m.tab = 0
			m.cursor = 0
		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			m.tab = 1
			m.cursor = 0
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.tab == 0 {
				if m.cursor < len(m.options)-1 {
					m.cursor++
				}
			} else {
				rows := m.mcpRows()
				if m.cursor < len(rows)-1 {
					m.cursor++
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("l"))):
			if m.tab == 1 {
				rows := m.mcpRows()
				if m.cursor < len(rows) && rows[m.cursor].isServer {
					m.mcpServers[rows[m.cursor].serverIdx].Expanded = true
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("h"))):
			if m.tab == 1 {
				rows := m.mcpRows()
				if m.cursor < len(rows) {
					row := rows[m.cursor]
					if row.isServer {
						m.mcpServers[row.serverIdx].Expanded = false
					} else {
						// Collapse parent server and move cursor to it
						m.mcpServers[row.serverIdx].Expanded = false
						// Find the server row index
						for ri, r := range m.mcpRows() {
							if r.isServer && r.serverIdx == row.serverIdx {
								m.cursor = ri
								break
							}
						}
					}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("c"))):
			// Trigger compression manually
			if m.totalUsed() > m.capacity/2 {
				m.runCompression()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if m.tab == 0 {
				opt := m.options[m.cursor]
				color := ui.ColorHighlight
				ch := '█'
				switch opt.Category {
				case "conversation":
					color = ui.ColorAccent
					ch = '▓'
				case "system":
					color = ui.ColorBeginner
					ch = '░'
				}
				m.items = append(m.items, contextItem{
					Name: opt.Name, Tokens: opt.Tokens, Color: color, Char: ch, Category: opt.Category,
				})
			} else {
				rows := m.mcpRows()
				if m.cursor < len(rows) {
					row := rows[m.cursor]
					if row.isServer {
						// Toggle entire server
						srv := &m.mcpServers[row.serverIdx]
						srv.Enabled = !srv.Enabled
						// When enabling, enable all tools; when disabling, disable all
						for i := range srv.Tools {
							srv.Tools[i].Enabled = srv.Enabled
						}
					} else {
						// Toggle individual tool
						srv := &m.mcpServers[row.serverIdx]
						srv.Tools[row.toolIdx].Enabled = !srv.Tools[row.toolIdx].Enabled
						// Update server enabled state: enabled if any tool is on
						anyOn := false
						for _, t := range srv.Tools {
							if t.Enabled {
								anyOn = true
								break
							}
						}
						srv.Enabled = anyOn
					}
				}
			}
			used := m.totalUsed()
			m.overflow = used > m.capacity
			// Auto-trigger compression on overflow
			if m.overflow {
				m.runCompression()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.items = []contextItem{
				{Name: "System Prompt", Tokens: 8000, Color: ui.ColorMuted, Char: '▒', Category: "system"},
				{Name: "CLAUDE.md", Tokens: 2000, Color: ui.ColorBeginner, Char: '░', Category: "system"},
			}
			for i := range m.mcpServers {
				m.mcpServers[i].Enabled = false
				m.mcpServers[i].Expanded = false
				for j := range m.mcpServers[i].Tools {
					m.mcpServers[i].Tools[j].Enabled = false
				}
			}
			m.overflow = false
			m.compressing = false
		}
	}
	return m, nil
}

func (m *BucketModel) View() string {
	if m.compressing {
		return m.viewCompression()
	}
	return m.viewNormal()
}

func (m *BucketModel) viewCompression() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	warn := lipgloss.NewStyle().Foreground(ui.ColorAdvanced).Bold(true)
	green := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)
	yellow := lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	strikethrough := lipgloss.NewStyle().Foreground(ui.ColorAdvanced).Strikethrough(true)
	arrow := lipgloss.NewStyle().Foreground(ui.ColorBeginner).Bold(true)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, warn.Render("  ⚠ CONTEXT OVERFLOW — COMPRESSION IN PROGRESS"))
	lines = append(lines, "")
	lines = append(lines, dim.Render(fmt.Sprintf("  Before: %s / %s tokens",
		formatTokens(m.preCompressUse), formatTokens(m.capacity))))
	lines = append(lines, "")

	lines = append(lines, accent.Render("  How compression works:"))
	lines = append(lines, dim.Render("  The system automatically compresses prior messages"))
	lines = append(lines, dim.Render("  to free space. Here's what happens step by step:"))
	lines = append(lines, "")

	tips := []string{
		"Keep sessions focused — avoid accumulating stale tool results.",
		"Use CLAUDE.md for persistent instructions so you don't repeat them each session.",
		"Break large tasks into smaller, focused sessions to stay within context.",
		"Reference files by path instead of pasting contents into the chat.",
		"Disable unused MCP servers — their tool definitions eat context silently.",
	}

	totalSaved := 0
	for i, entry := range m.compressLog {
		stepNum := i + 1
		if i < m.compressStep {
			// Revealed step
			lines = append(lines, highlight.Render(fmt.Sprintf("  Step %d:", stepNum)))
			lines = append(lines, "    "+strikethrough.Render(entry.Before))
			lines = append(lines, "    "+arrow.Render("→ ")+green.Render(entry.After))
			lines = append(lines, "    "+yellow.Render(fmt.Sprintf("Freed %s tokens", formatTokens(entry.Saved))))
			// Show a tip for this step
			tip := tips[i%len(tips)]
			lines = append(lines, "    "+dim.Render("💡 Tip: "+tip))
			lines = append(lines, "")
			totalSaved += entry.Saved
		} else {
			// Hidden step
			lines = append(lines, dim.Render(fmt.Sprintf("  Step %d: ...", stepNum)))
			lines = append(lines, "")
		}
	}

	// Summary
	if m.compressStep >= len(m.compressLog) {
		lines = append(lines, accent.Render("  ── Compression Complete ──"))
		lines = append(lines, "")

		afterUsed := m.totalUsed()
		lines = append(lines, dim.Render(fmt.Sprintf("  Before:  %s tokens", formatTokens(m.preCompressUse))))
		lines = append(lines, green.Render(fmt.Sprintf("  After:   %s tokens", formatTokens(afterUsed))))
		lines = append(lines, yellow.Render(fmt.Sprintf("  Saved:   %s tokens", formatTokens(totalSaved))))
		lines = append(lines, "")

		pct := float64(afterUsed) / float64(m.capacity) * 100
		barWidth := 30
		filled := int(pct / 100 * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}
		bar := green.Render(strings.Repeat("█", filled)) + dim.Render(strings.Repeat("░", barWidth-filled))
		lines = append(lines, fmt.Sprintf("  %s %.0f%%", bar, pct))
		lines = append(lines, "")

		lines = append(lines, accent.Render("  What was preserved:"))
		lines = append(lines, green.Render("    ✓ System prompt & CLAUDE.md (never compressed)"))
		lines = append(lines, green.Render("    ✓ Recent tool results (most relevant)"))
		lines = append(lines, green.Render("    ✓ Latest conversation (current context)"))
		lines = append(lines, "")
		lines = append(lines, warn.Render("  What was lost:"))
		lines = append(lines, warn.Render("    ✗ Old tool results (replaced by summaries)"))
		lines = append(lines, warn.Render("    ✗ Older conversation messages (condensed)"))
		lines = append(lines, warn.Render("    ✗ Verbose output (truncated to key lines)"))

		lines = append(lines, "")
		lines = append(lines, highlight.Render("  Press Enter to return  ")+dim.Render("[r] Reset"))
	} else {
		lines = append(lines, highlight.Render("  Press Enter/Space to reveal next step"))
		lines = append(lines, dim.Render(fmt.Sprintf("  (%d of %d steps)", m.compressStep, len(m.compressLog))))
	}

	return strings.Join(lines, "\n")
}

func (m *BucketModel) viewNormal() string {
	accent := lipgloss.NewStyle().Foreground(ui.ColorAccent).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	highlight := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true)
	warn := lipgloss.NewStyle().Foreground(ui.ColorAdvanced).Bold(true)
	green := lipgloss.NewStyle().Foreground(ui.ColorBeginner)
	mcpClr := lipgloss.NewStyle().Foreground(lipgloss.Color("#f97316")).Bold(true)
	tabActive := lipgloss.NewStyle().Foreground(ui.ColorHighlight).Bold(true).Underline(true)
	tabInactive := dim

	sysStyle := lipgloss.NewStyle().Foreground(ui.ColorMuted)
	toolStyle := lipgloss.NewStyle().Foreground(ui.ColorHighlight)
	convStyle := lipgloss.NewStyle().Foreground(ui.ColorAccent)

	used := m.totalUsed()
	sysTok, toolTok, convTok, mcpTok := m.categoryTotals()
	remaining := m.capacity - used
	if remaining < 0 {
		remaining = 0
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, accent.Render("  Context Window Bucket"))
	lines = append(lines, "")

	// ── Bucket ──
	bucketWidth := 34
	bucketHeight := 8
	innerWidth := bucketWidth - 2

	type rowInfo struct {
		char  rune
		color lipgloss.Color
	}
	rowData := make([]rowInfo, bucketHeight)
	currentRow := 0

	for _, item := range m.items {
		itemRows := int(float64(item.Tokens) / float64(m.capacity) * float64(bucketHeight))
		if itemRows < 1 && item.Tokens > 0 {
			itemRows = 1
		}
		for r := 0; r < itemRows && currentRow < bucketHeight; r++ {
			rowData[currentRow] = rowInfo{char: item.Char, color: item.Color}
			currentRow++
		}
	}
	if mcpTok > 0 {
		mcpRows := int(float64(mcpTok) / float64(m.capacity) * float64(bucketHeight))
		if mcpRows < 1 {
			mcpRows = 1
		}
		for r := 0; r < mcpRows && currentRow < bucketHeight; r++ {
			rowData[currentRow] = rowInfo{char: '▣', color: lipgloss.Color("#f97316")}
			currentRow++
		}
	}

	top := "  ┌" + strings.Repeat("─", innerWidth) + "┐"
	lines = append(lines, top)
	for i := bucketHeight - 1; i >= 0; i-- {
		if rowData[i].char != 0 {
			fill := strings.Repeat(string(rowData[i].char), innerWidth)
			colored := lipgloss.NewStyle().Foreground(rowData[i].color).Render(fill)
			lines = append(lines, "  │"+colored+"│")
		} else {
			lines = append(lines, "  │"+strings.Repeat(" ", innerWidth)+"│")
		}
	}
	bottom := "  └" + strings.Repeat("─", innerWidth) + "┘"
	lines = append(lines, bottom)

	pct := float64(used) / float64(m.capacity) * 100
	counterStyle := green
	if pct > 75 {
		counterStyle = lipgloss.NewStyle().Foreground(ui.ColorIntermediate)
	}
	if pct > 100 {
		counterStyle = warn
	}
	lines = append(lines, counterStyle.Render(fmt.Sprintf("  %s / %s tokens (%.0f%%)",
		formatTokens(used), formatTokens(m.capacity), pct)))
	lines = append(lines, "")

	// ── Breakdown ──
	barWidth := 22
	renderBar := func(label string, tokens int, style lipgloss.Style) string {
		ratio := float64(tokens) / float64(m.capacity)
		filled := int(ratio * float64(barWidth))
		if filled < 1 && tokens > 0 {
			filled = 1
		}
		if filled > barWidth {
			filled = barWidth
		}
		bar := style.Render(strings.Repeat("█", filled)) + dim.Render(strings.Repeat("░", barWidth-filled))
		pctStr := fmt.Sprintf("%.0f%%", ratio*100)
		// Pad the label to 13 chars BEFORE styling to avoid ANSI codes breaking alignment
		padded := label + strings.Repeat(" ", 13-len(label))
		if len(label) > 13 {
			padded = label
		}
		return fmt.Sprintf("  %s %s %5s %s", style.Render(padded), bar, dim.Render(pctStr), dim.Render(formatTokens(tokens)))
	}

	lines = append(lines, dim.Render("  ── Context Breakdown ──"))
	lines = append(lines, renderBar("System", sysTok, sysStyle))
	if toolTok > 0 {
		lines = append(lines, renderBar("Tool results", toolTok, toolStyle))
	}
	if convTok > 0 {
		lines = append(lines, renderBar("Conversation", convTok, convStyle))
	}
	if mcpTok > 0 {
		lines = append(lines, renderBar("MCP tool defs", mcpTok, mcpClr))
	}
	lines = append(lines, renderBar("Remaining", remaining, green))

	// Compression hint
	if pct > 50 && pct <= 100 {
		lines = append(lines, "")
		lines = append(lines, dim.Render("  [c] Simulate compression"))
	}

	lines = append(lines, "")

	// ── Tabs ──
	tab0 := tabInactive.Render("  Tools & Actions")
	tab1 := tabInactive.Render("  MCP Servers")
	if m.tab == 0 {
		tab0 = tabActive.Render("  Tools & Actions")
	} else {
		tab1 = tabActive.Render("  MCP Servers")
	}
	lines = append(lines, tab0+"    "+tab1+dim.Render("    [1/2] switch"))
	lines = append(lines, "")

	if m.tab == 0 {
		startOpt := m.cursor - 3
		if startOpt < 0 {
			startOpt = 0
		}
		endOpt := startOpt + 6
		if endOpt > len(m.options) {
			endOpt = len(m.options)
		}
		for i := startOpt; i < endOpt; i++ {
			opt := m.options[i]
			style := dim
			prefix := "  "
			if i == m.cursor {
				style = highlight
				prefix = "▸ "
			}
			tokLabel := fmt.Sprintf("+%s", formatTokens(opt.Tokens))
			lines = append(lines, fmt.Sprintf("  %s%-20s %6s  %s",
				prefix, style.Render(opt.Name), dim.Render(tokLabel), dim.Render(opt.Desc)))
		}
		lines = append(lines, "", dim.Render("  [↑/↓] Navigate  [Enter/Space] Add  [2] MCP  [r] Reset"))
	} else {
		totalMCPTools := 0
		for _, s := range m.mcpServers {
			totalMCPTools += s.ToolCount
		}
		lines = append(lines, warn.Render(fmt.Sprintf("  All servers = %s tokens just for definitions!",
			formatTokens(totalMCPTools*200))))
		lines = append(lines, "")

		rows := m.mcpRows()
		// Scrolling window
		startRow := m.cursor - 5
		if startRow < 0 {
			startRow = 0
		}
		endRow := startRow + 12
		if endRow > len(rows) {
			endRow = len(rows)
		}

		for ri := startRow; ri < endRow; ri++ {
			row := rows[ri]
			isCursor := ri == m.cursor

			if row.isServer {
				srv := m.mcpServers[row.serverIdx]
				style := dim
				prefix := "  "
				if isCursor {
					style = highlight
					prefix = "▸ "
				}
				checkbox := "[ ]"
				if srv.Enabled {
					checkbox = "[✓]"
					if !isCursor {
						style = mcpClr
					}
				}
				// Count enabled tools
				enabledCount := 0
				for _, t := range srv.Tools {
					if t.Enabled {
						enabledCount++
					}
				}
				expand := "▸"
				if srv.Expanded {
					expand = "▾"
				}
				toolLabel := fmt.Sprintf("%d/%d tools", enabledCount, srv.ToolCount)
				tokUsed := enabledCount * 200
				lines = append(lines, fmt.Sprintf("  %s%s %s %-12s %s  %s",
					prefix, style.Render(checkbox), dim.Render(expand),
					style.Render(srv.Name), dim.Render(toolLabel),
					dim.Render(formatTokens(tokUsed)+" tok")))
			} else {
				srv := m.mcpServers[row.serverIdx]
				tool := srv.Tools[row.toolIdx]
				style := dim
				prefix := "    "
				if isCursor {
					style = highlight
					prefix = "  ▸ "
				}
				checkbox := "○"
				if tool.Enabled {
					checkbox = "●"
					if !isCursor {
						style = mcpClr
					}
				}
				isLast := row.toolIdx == len(srv.Tools)-1
				branch := "├"
				if isLast {
					branch = "└"
				}
				lines = append(lines, fmt.Sprintf("  %s%s %s %s",
					prefix, dim.Render(branch), style.Render(checkbox),
					style.Render(tool.Name)))
			}
		}

		if len(rows) > 12 {
			lines = append(lines, dim.Render(fmt.Sprintf("  ... %d more (scroll with ↑/↓)", len(rows)-12)))
		}

		lines = append(lines, "")
		lines = append(lines, mcpClr.Render("  Tip: ")+dim.Render("Expand a server and disable tools you don't need."))
		lines = append(lines, dim.Render("  Each tool definition costs ~200 tokens of context."))
		lines = append(lines, "")
		lines = append(lines, dim.Render("  [↑/↓] Navigate  [Enter] Toggle  [l] Expand  [h] Collapse  [r] Reset"))
	}

	return strings.Join(lines, "\n")
}

func formatTokens(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%d", n)
}
