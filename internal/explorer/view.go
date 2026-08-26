package explorer

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/spinner"
	"github.com/eduardofuncao/squix/internal/styles"
)

// chromeLines is the number of screen lines reserved for the title, separator
// and footer areas; the list viewport gets the rest.
func chromeLines() int { return 6 }

func (m Model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	if m.helpOverlay {
		return tea.NewView(m.renderHelpOverlay())
	}

	if m.relationsMode {
		return tea.NewView(m.renderRelationsPanel())
	}

	return tea.NewView(m.renderListView())
}

func (m Model) renderListView() string {
	var b strings.Builder

	title := styles.Title.Render("◆ explore")
	if m.conn != nil {
		title += styles.Faint.Render(
			fmt.Sprintf(" · %s (%s)", m.conn.GetName(), m.conn.GetDbType()),
		)
	}
	b.WriteString(title)
	b.WriteString("\n")

	b.WriteString(styles.Separator.Render(strings.Repeat("─", max(m.width-1, 1))))
	b.WriteString("\n")

	if m.numRows() == 0 {
		b.WriteString(styles.Faint.Render("No tables or views found"))
	}

	endRow := min(m.offsetY+m.visibleRows, m.numRows())
	for i := m.offsetY; i < endRow; i++ {
		b.WriteString(m.renderRow(i))
		b.WriteString("\n")
	}

	b.WriteString(m.renderListFooter())

	if m.statusMessage != "" {
		b.WriteString("\n")
		b.WriteString(m.statusMessage)
	}

	return b.String()
}

func (m Model) renderRow(idx int) string {
	r := m.rows[idx]
	selected := idx == m.cursor

	if r.kind == rowHeader {
		sec := m.sections[r.section]
		marker := "▾"
		if sec.folded {
			marker = "▸"
		}
		label := fmt.Sprintf("%s %s %s", marker, sec.title, styles.Faint.Render(fmt.Sprintf("(%d)", len(sec.items))))
		if selected {
			return styles.TableSelected.Render(label)
		}
		return styles.TableHeader.Render(label)
	}

	item := m.sections[r.section].items[r.item]
	prefix := "  "
	switch {
	case selected:
		return styles.TableSelected.Render(prefix + item)
	case isSearchMatch(idx, m.searchMatches):
		return styles.SearchMatch.Render(prefix + item)
	default:
		return styles.TableCell.Render(prefix + item)
	}
}

// isSearchMatch reports whether idx is in the sorted searchMatches slice.
func isSearchMatch(idx int, matches []int) bool {
	lo, hi := 0, len(matches)
	for lo < hi {
		mid := (lo + hi) / 2
		if matches[mid] < idx {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo < len(matches) && matches[lo] == idx
}

func (m Model) renderListFooter() string {
	var b strings.Builder

	if m.searchMode {
		cursorBefore := m.searchQuery[:m.searchCursor]
		cursorAfter := ""
		if m.searchCursor < len(m.searchQuery) {
			cursorAfter = m.searchQuery[m.searchCursor:]
		}

		input := styles.SearchMatch.Render("/") + " " + cursorBefore + "█" + cursorAfter
		b.WriteString("\n")
		b.WriteString(input)
		b.WriteString("\n")
		b.WriteString(styles.Faint.Render("Enter: search  Esc: cancel"))
		return b.String()
	}

	if m.uiVisibility.FooterKeymaps {
		km := m.keyMap
		var parts []string
		addEntry := func(key, desc string) {
			parts = append(parts, styles.TableHeader.Render(key)+":"+styles.Faint.Render(desc))
		}
		addEntry(km.FirstKey(config.ActionMoveUp)+km.FirstKey(config.ActionMoveDown), "navigate")
		addEntry(km.FirstKey(config.ActionExplorerSelect), "select *")
		addEntry(km.FirstKey(config.ActionExplorerColumns), "columns")
		addEntry(km.FirstKey(config.ActionExplorerRelations), "relations")
		addEntry(km.FirstKey(config.ActionExplorerToggleFold), "fold")
		addEntry(km.FirstKey(config.ActionSearch), "search")
		addEntry(km.FirstKey(config.ActionHelp), "help")
		addEntry(km.FirstKey(config.ActionQuit), "quit")
		b.WriteString("\n")
		b.WriteString("  " + strings.Join(parts, " "))
	} else {
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderRelationsPanel() string {
	var b strings.Builder

	loadingLine := ""
	if m.relationsLoading {
		stage := spinner.Stages[m.spinnerFrame%len(spinner.Stages)]
		elapsed := time.Since(m.relationsLoadStart).Seconds()
		loadingLine = styles.Success.Render(stage) + fmt.Sprintf(" %.2fs", elapsed)
	}

	lines := m.relationsLines()
	availableHeight := m.relationsViewportHeight()
	panelHeight := min(len(lines), availableHeight)
	if panelHeight < 1 {
		panelHeight = 1
	}

	start := m.relationsScroll
	end := min(start+panelHeight, len(lines))

	// Panel hugs its content: width fits the widest visible line.
	titleText := fmt.Sprintf("◆ relationships · %s (depth %d)", m.relationsTarget, relationsDepth)
	panelWidth := lipgloss.Width(titleText)
	maxLineWidth := max(m.width-2, 20)

	visible := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		line := truncateANSI(lines[i], maxLineWidth)
		visible = append(visible, line)
		if w := lipgloss.Width(line); w > panelWidth {
			panelWidth = w
		}
	}
	scrollInfo := ""
	if !m.relationsLoading && len(lines) > panelHeight {
		scrollInfo = fmt.Sprintf(" [%d-%d of %d]", start+1, end, len(lines))
	}
	km := m.keyMap
	closeKeys := ""
	if !m.relationsLoading {
		closeKeys = "  " + km.DisplayKeys(config.ActionQuit) + " close"
	}
	footerText := km.FirstKey(config.ActionMoveDown) + "/" +
		km.FirstKey(config.ActionMoveUp) + " scroll" + scrollInfo + closeKeys

	panelWidth = min(max(panelWidth+2, 24), max(m.width-1, 1))

	b.WriteString(styles.Title.Render(titleText))
	b.WriteString("\n")

	b.WriteString(styles.Separator.Render(strings.Repeat("─", panelWidth)))

	if m.relationsLoading {
		b.WriteString("\n")
		b.WriteString(loadingLine)
	} else {
		for _, line := range visible {
			b.WriteString("\n")
			b.WriteString(styles.TableCell.Render(line))
		}
	}

	// Footer hugs the content so the panel stays as tall as the tree.
	b.WriteString(footerText)

	return b.String()
}

// relationsViewportHeight is the number of tree lines the panel can show at
// once: everything except the title row, its separator and the footer.
func (m Model) relationsViewportHeight() int {
	h := m.height - 3
	if h < 3 {
		h = 3
	}
	return h
}

// relationsLines splits the tree into logical lines. Lines wider than the
// terminal are truncated at render time, so no width-dependent cache needed.
func (m Model) relationsLines() []string {
	return strings.Split(m.relationsText, "\n")
}

// truncateANSI shortens a styled line to width visible cells, preserving
// escape sequences and marking truncation with an ellipsis.
func truncateANSI(line string, width int) string {
	if lipgloss.Width(line) <= width {
		return line
	}
	return ansi.Truncate(line, max(width-1, 1), "…")
}

func (m Model) renderHelpOverlay() string {
	type keyBind struct {
		keys string
		desc string
	}
	type category struct {
		name  string
		binds []keyBind
	}

	km := m.keyMap

	categories := []category{
		{"Navigation", []keyBind{
			{km.DisplayKeys(config.ActionMoveUp), "Move up"},
			{km.DisplayKeys(config.ActionMoveDown), "Move down"},
			{km.DisplayKeys(config.ActionJumpFirstRow), "Jump to first row"},
			{km.DisplayKeys(config.ActionJumpLastRow), "Jump to last row"},
			{km.DisplayKeys(config.ActionPageUp), "Page up"},
			{km.DisplayKeys(config.ActionPageDown), "Page down"},
		}},
		{"Actions", []keyBind{
			{km.DisplayKeys(config.ActionExplorerSelect), "SELECT * from item"},
			{km.DisplayKeys(config.ActionExplorerColumns), "Show columns metadata"},
			{km.DisplayKeys(config.ActionExplorerRelations), "Show relationships tree"},
			{km.DisplayKeys(config.ActionExplorerToggleFold), "Fold/unfold section"},
		}},
		{"Search", []keyBind{
			{km.DisplayKeys(config.ActionSearch), "Search items"},
			{km.DisplayKeys(config.ActionNextMatch) + " / " + km.DisplayKeys(config.ActionPrevMatch), "Next/previous match"},
		}},
		{"Other", []keyBind{
			{km.DisplayKeys(config.ActionToggleFooter), "Toggle footer keybinds"},
			{km.DisplayKeys(config.ActionHelp), "Show/hide this help"},
			{km.DisplayKeys(config.ActionQuit), "Quit"},
		}},
	}

	separatorWidth := max(m.width-4, 0)

	var b strings.Builder
	b.WriteString(styles.Title.Render("Keyboard Shortcuts"))
	b.WriteString("\n")
	b.WriteString(styles.Separator.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	for _, cat := range categories {
		b.WriteString(styles.TableHeader.Render(cat.name))
		b.WriteString("\n")
		for _, bind := range cat.binds {
			line := fmt.Sprintf("  %-16s %s", bind.keys, bind.desc)
			b.WriteString(styles.Faint.Render(line))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.Separator.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n")
	b.WriteString(styles.Faint.Render(km.DisplayKeys(config.ActionHelpClose) + " to close"))

	return b.String()
}
