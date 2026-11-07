package drill

import (
	"fmt"
	"strings"
	"time"

	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/tuday/internal/doomsday"
	"github.com/vinegm/tuday/internal/report"
)

const (
	// How many of the latest answers
	// the history strip shows.
	historyLength = 25

	help = "enter answer ~ 0-7 answer instantly ~ esc finish"
)

// SetSize does nothing, the drill is drawn at its own size.
func (window *Window) SetSize(_, _ int) {}

// View renders the drill.
func (window *Window) View() string {
	th := window.theme

	lines := []string{
		th.TitleStyle.Render("tuday") + th.MutedStyle.Render("  "+window.options.Years.Label()),
		th.TextStyle.Render(window.score()),
		window.history(),
		"",
		th.AccentStyle.Render(doomsday.FormatDate(window.date)),
		th.MutedStyle.Render(window.doomsdayHint()),
		"",
		window.inputLine(),
		th.MutedStyle.Render(fmt.Sprintf(
			"%.1fs",
			time.Since(window.askedAt).Seconds(),
		)),
		"",
	}

	lines = append(lines, window.feedbackLines()...)
	lines = append(
		lines,
		"",
		th.MutedStyle.Render(window.cheatsheet()),
		th.MutedStyle.Render(help),
	)

	return th.WindowStyle.Render(strings.Join(lines, "\n"))
}

func (window *Window) score() string {
	if len(window.attempts) == 0 {
		return "no answers yet"
	}

	correct := report.CountCorrect(window.attempts)

	percent := 100 * float64(correct) / float64(len(window.attempts))

	return fmt.Sprintf(
		"%d/%d correct (%.0f%%) · avg %.1fs",
		correct,
		len(window.attempts),
		percent,
		report.Mean(window.attempts),
	)
}

// history is a strip of dots, one per recent
// answer, colored by whether it was right.
func (window *Window) history() string {
	recent := window.attempts[max(0, len(window.attempts)-historyLength):]

	var sb strings.Builder
	for _, attempt := range recent {
		if attempt.Correct() {
			sb.WriteString(window.theme.CorrectStyle.Render("●"))
			continue
		}

		sb.WriteString(window.theme.WrongStyle.Render("●"))
	}

	return sb.String()
}

func (window *Window) doomsdayHint() string {
	if !window.options.ShowDoomsday {
		return ""
	}

	year := window.date.Year()

	return fmt.Sprintf("doomsday %d: %s", year, doomsday.Of(year))
}

func (window *Window) inputLine() string {
	prompt := window.theme.SelectedMenuTextStyle.Render("› ")

	if window.input == "" {
		return prompt + window.theme.MutedStyle.Render("weekday or 0-7")
	}

	return prompt + window.theme.TextStyle.Render(window.input) + window.theme.SelectedMenuTextStyle.Render("▏")
}

// feedbackLines are the verdict on the
// previous submission and its detail.
func (window *Window) feedbackLines() []string {
	last := window.last
	if last == nil {
		return []string{"", ""}
	}

	th := window.theme
	width := th.InnerWidth()

	switch last.outcome {
	case answeredRight:
		return []string{
			th.CorrectStyle.Render(lg.NewStyle().MaxWidth(width).Render(last.String())),
			th.MutedStyle.Render(fmt.Sprintf(
				"%.1fs",
				last.attempt.Seconds()),
			),
		}

	case answeredWrong:
		return []string{
			th.WrongStyle.Render(lg.NewStyle().MaxWidth(width).Render(last.String())),
			th.MutedStyle.Render(fmt.Sprintf(
				"you said %s in %.1fs",
				last.attempt.Guess,
				last.attempt.Seconds()),
			),
		}

	default:
		return []string{
			th.WrongStyle.Render(lg.NewStyle().MaxWidth(width).Render(last.String())),
			th.MutedStyle.Render("try again, the clock is still running"),
		}
	}
}

func (window *Window) cheatsheet() string {
	if !window.options.Cheatsheet {
		return ""
	}

	return doomsday.Cheatsheet()
}
