package drill

import (
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"strings"
	"time"
	"unicode"

	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/tuday/internal/doomsday"
	"github.com/vinegm/tuday/internal/tui/theme"
)

// How often the question timer on screen refreshes.
const tickInterval = 100 * time.Millisecond

// Longest answer the
// input accepts.
const maxInput = 12

// Options control what the
// drill asks and shows.
type Options struct {
	Years        doomsday.Years
	ShowDoomsday bool
	Cheatsheet   bool
	Seed         string // date order, random when empty
}

type (
	// FinishedMsg asks for control to end
	// the session, handing over what was answered.
	FinishedMsg struct{ Attempts []doomsday.Attempt }

	// Message carrying a timer tick, tagged
	// with the session that armed it.
	tickMsg struct{ tag int }
)

type outcome int

const (
	answeredRight outcome = iota
	answeredWrong
	notAWeekday
)

// feedback is what the previous
// submission turned out to be.
type feedback struct {
	outcome outcome
	attempt doomsday.Attempt
	input   string
}

// String form of a feedback line.
func (last feedback) String() string {
	switch last.outcome {
	case answeredRight:
		return fmt.Sprintf(
			"✓ %s was a %s",
			doomsday.FormatDate(last.attempt.Date),
			last.attempt.Date.Weekday(),
		)

	case answeredWrong:
		return fmt.Sprintf(
			"✗ %s was a %s",
			doomsday.FormatDate(last.attempt.Date),
			last.attempt.Date.Weekday(),
		)

	default:
		return fmt.Sprintf(
			"? %q is not a weekday",
			last.input,
		)
	}
}

// Window asks one date after another
// and times each answer.
type Window struct {
	options  Options
	attempts []doomsday.Attempt
	rng      *rand.Rand

	date    time.Time
	askedAt time.Time
	input   string
	last    *feedback

	tickTag int

	theme *theme.Theme
}

// New initializes a new Window.
//
//	NOTE: No date is asked until Start() is called.
func New(th *theme.Theme) *Window {
	return &Window{
		theme: th,
	}
}

// Start begins a fresh session with the given options.
func (window *Window) Start(options Options) tea.Cmd {
	window.options = options
	window.rng = rngFrom(options.Seed)
	window.attempts = nil
	window.last = nil
	window.tickTag++

	window.nextQuestion()

	return tick(window.tickTag)
}

// Update handles the keys and keeps the timer ticking.
func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(msg)

	case tickMsg:
		if msg.tag != window.tickTag {
			return nil
		}

		return tick(window.tickTag)
	}

	return nil
}

// Hints is the keys for the key hints bar.
func (window *Window) Hints() []theme.Hint {
	if window.input == "" {
		return []theme.Hint{
			{Key: "0-7", Desc: "Answer"},
			{Key: "esc", Desc: "Finish"},
		}
	}

	return []theme.Hint{
		{Key: "enter", Desc: "Answer"},
		{Key: "ctrl+u", Desc: "Clear"},
		{Key: "esc", Desc: "Finish"},
	}
}

func (window *Window) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		return window.finish()

	case "enter":
		return window.submit(window.input)

	case "backspace":
		runes := []rune(window.input)
		if len(runes) > 0 {
			window.input = string(runes[:len(runes)-1])
		}

		return nil

	case "ctrl+u":
		window.input = ""
		return nil
	}

	// NOTE: lone digit is a full
	// answer and doesnt wait
	isEmpty := window.input == ""
	isValidDigit := len(msg.Text) == 1 && msg.Text[0] >= '0' && msg.Text[0] <= '7'
	if isEmpty && isValidDigit {
		return window.submit(msg.Text)
	}

	if isLetters(msg.Text) && len([]rune(window.input+msg.Text)) <= maxInput {
		window.input += msg.Text
	}

	return nil
}

func (window *Window) submit(raw string) tea.Cmd {
	elapsed := time.Since(window.askedAt)
	text := strings.TrimSpace(raw)
	window.input = ""

	if text == "" {
		return nil
	}

	guess, ok := doomsday.ParseGuess(text)
	if !ok {
		// Same date again, the
		// clock keeps running
		window.last = &feedback{
			outcome: notAWeekday,
			input:   text,
		}
		return nil
	}

	attempt := doomsday.Attempt{
		Date:    window.date,
		Guess:   guess,
		Elapsed: elapsed,
	}

	window.attempts = append(
		window.attempts,
		attempt,
	)

	window.last = &feedback{
		outcome: answeredWrong,
		attempt: attempt,
	}

	if attempt.Correct() {
		window.last.outcome = answeredRight
	}

	window.nextQuestion()

	return nil
}

func (window *Window) nextQuestion() {
	window.date = doomsday.RandomDate(window.rng, window.options.Years)
	window.askedAt = time.Now()
	window.input = ""
}

// rngFrom is where a session's questions come from: the same seed
// asks the same ones, and an empty seed different ones every time.
func rngFrom(seed string) *rand.Rand {
	if seed == "" {
		return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}

	hash := fnv.New64a()
	hash.Write([]byte(seed))

	return rand.New(rand.NewPCG(hash.Sum64(), 0))
}

func (window *Window) finish() tea.Cmd {
	window.tickTag++
	attempts := window.attempts

	return func() tea.Msg {
		return FinishedMsg{Attempts: attempts}
	}
}

func tick(tag int) tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{tag: tag}
	})
}

func isLetters(text string) bool {
	if text == "" {
		return false
	}

	for _, char := range text {
		if !unicode.IsLetter(char) {
			return false
		}
	}

	return true
}
