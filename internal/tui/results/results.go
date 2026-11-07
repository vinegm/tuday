package results

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/tuday/internal/doomsday"
	"github.com/vinegm/tuday/internal/report"
	"github.com/vinegm/tuday/internal/tui/theme"
)

// Lines the window spends around
// the breakdown table.
const chromeHeight = 12

const barWidth = 12

const help = "enter again · ←/→ breakdown · j/k scroll · q menu"

type (
	// AgainMsg asks for control to drill again
	// with the same options.
	AgainMsg struct{}

	// ExitMsg asks for control to
	// return to the menu.
	ExitMsg struct{}
)

// Window shows how a finished session
// went, one breakdown at a time.
type Window struct {
	attempts   []doomsday.Attempt
	breakdowns []report.Breakdown

	tab     int
	scroll  int
	visible int // rows of the breakdown shown at the window's size

	theme *theme.Theme
}

// New initializes a new Window.
//
//	NOTE: Nothing is shown until Show() is called.
func New(th *theme.Theme) *Window {
	return &Window{
		theme: th,
	}
}

// Show switches the window to a newly finished session.
func (window *Window) Show(attempts []doomsday.Attempt, breakdowns []report.Breakdown) {
	window.attempts = attempts
	window.breakdowns = breakdowns
	window.tab = 0
	window.scroll = 0
}

// Update handles the keys.
func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(msg.String())
	}

	return nil
}

func (window *Window) handleKey(key string) tea.Cmd {
	switch key {
	case "tab", "right", "l":
		window.tab = (window.tab + 1) % len(window.breakdowns)
		window.scroll = 0

	case "shift+tab", "left", "h":
		window.tab = (window.tab + len(window.breakdowns) - 1) % len(window.breakdowns)
		window.scroll = 0

	case "down", "j":
		window.scroll++

	case "up", "k":
		window.scroll--

	case "enter", "space", "r":
		return func() tea.Msg {
			return AgainMsg{}
		}

	case "q", "esc":
		return func() tea.Msg {
			return ExitMsg{}
		}
	}

	return nil
}

// SetSize fits the breakdown to height, with room for the longest
// one so that switching between them doesn't resize the window.
func (window *Window) SetSize(_, height int) {
	longest := 0
	for _, breakdown := range window.breakdowns {
		longest = max(longest, len(breakdown.Buckets))
	}

	window.visible = min(longest, max(1, height-chromeHeight))
	window.scroll = min(max(window.scroll, 0), max(0, window.rankedCount()-window.visible))
}

func (window *Window) rankedCount() int {
	return len(window.breakdowns[window.tab].Buckets)
}

// View renders the results window as a string.
func (window *Window) View() string {
	th := window.theme
	breakdown := window.breakdowns[window.tab]
	ranked := breakdown.Ranked()

	rows := make([]string, 0, window.visible)
	for _, bucket := range ranked[window.scroll:min(len(ranked), window.scroll+window.visible)] {
		rows = append(rows, window.row(bucket))
	}

	for len(rows) < window.visible {
		rows = append(rows, "")
	}

	lines := []string{
		th.TitleStyle.Render("Session results"),
		th.TextStyle.Render(window.summary()),
		"",
		window.tabs(),
		th.MutedStyle.Render(worstLine(breakdown)),
		"",
		lg.JoinVertical(lg.Left, rows...),
		th.MutedStyle.Render(window.scrollMarker(len(ranked), window.visible)),
		"",
		th.MutedStyle.Render(help),
	}

	return th.WindowStyle.Render(strings.Join(lines, "\n"))
}

func (window *Window) summary() string {
	correct := report.CountCorrect(window.attempts)
	percent := 100 * float64(correct) / float64(len(window.attempts))

	return fmt.Sprintf(
		"%d/%d correct (%.0f%%) ~ median %.1fs ~ worst %.1fs",
		correct,
		len(window.attempts),
		percent,
		report.Median(window.attempts),
		report.Slowest(window.attempts),
	)
}

func (window *Window) tabs() string {
	names := make([]string, len(window.breakdowns))

	for i, breakdown := range window.breakdowns {
		if i == window.tab {
			names[i] = window.theme.SelectedMenuTextStyle.Render("[" + breakdown.Name + "]")
			continue
		}

		names[i] = window.theme.MutedStyle.Render(" " + breakdown.Name + " ")
	}

	return strings.Join(names, "  ")
}

func worstLine(breakdown report.Breakdown) string {
	if worst, ok := breakdown.Worst(); ok {
		return "worst: " + worst
	}

	return "flawless"
}

// row is one bucket.
func (window *Window) row(bucket report.Bucket) string {
	th := window.theme
	filled := (barWidth*bucket.Correct + bucket.Asked/2) / bucket.Asked

	return fmt.Sprintf("%-9s  %s%s  %7s  %5.1fs",
		bucket.Label,
		th.CorrectStyle.Render(strings.Repeat("█", filled)),
		th.WrongStyle.Render(strings.Repeat("░", barWidth-filled)),
		fmt.Sprintf("%d/%d", bucket.Correct, bucket.Asked),
		bucket.Median,
	)
}

func (window *Window) scrollMarker(total, visible int) string {
	above := window.scroll
	below := total - visible - window.scroll

	switch {
	case above > 0 && below > 0:
		return fmt.Sprintf("↑ %d more · ↓ %d more", above, below)

	case above > 0:
		return fmt.Sprintf("↑ %d more", above)

	case below > 0:
		return fmt.Sprintf("↓ %d more", below)

	default:
		return ""
	}
}
