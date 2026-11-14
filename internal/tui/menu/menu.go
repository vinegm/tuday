package menu

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/tuday/internal/doomsday"
	"github.com/vinegm/tuday/internal/tui/drill"
	"github.com/vinegm/tuday/internal/tui/theme"
)

type menuOption int

const (
	startDrill menuOption = iota
	singleYear
	startYear
	endYear
	showDoomsday
	cheatsheet
	quit
)

// Longest year a year row
// accepts typed in.
const maxYearDigits = 4

// menuItem is the label shown and
// the action taken when selected.
type menuItem struct {
	label  string
	action menuOption
}

// StartDrillMsg asks for control to
// start a drill with these options.
type StartDrillMsg struct {
	Options drill.Options
}

// Window is the main menu, where the user picks
// the year range and hints before drilling.
type Window struct {
	Choice  int
	Options drill.Options

	// Drill only Options.Years.Start.
	Single bool

	// Digits typed on a year row, applied as they are typed and
	// settled into a valid range once the row is left.
	typing bool
	typed  string

	theme *theme.Theme
}

// New creates a new Window with the provided theme and options.
func New(th *theme.Theme, options drill.Options, single bool) *Window {
	return &Window{
		Options: options,
		Single:  single,
		theme:   th,
	}
}

// Update handles the keys.
func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(msg)
	}

	return nil
}

// SetSize does nothing, the menu is drawn at its own size.
func (window *Window) SetSize(_, _ int) {}

// View renders the menu window as a string.
func (window *Window) View() string {
	var sb strings.Builder

	// WARN: Left-aligned so the art's lines are
	// padded to one width and center as a block.
	sb.WriteString(window.theme.TitleStyle.Align(lg.Left).Render(getMenuTitle()))
	sb.WriteString("\n")
	sb.WriteString(window.theme.MutedStyle.Render("Doomsday rule trainer"))
	sb.WriteString("\n")

	for i, item := range window.items() {
		sb.WriteString("\n")
		if window.Choice == i {
			sb.WriteString(window.theme.SelectedMenuTextStyle.Render(item.label))
			continue
		}

		sb.WriteString(window.theme.TextStyle.Render(item.label))
	}

	return window.theme.WindowStyle.Render(sb.String())
}

// ShortHelp is the keys for the key hints bar, for the selected row.
func (window *Window) ShortHelp() []key.Binding {
	if window.typing {
		return []key.Binding{keys.Digit, keys.Delete, keys.Done, keys.Quit}
	}

	bindings := []key.Binding{keys.Move}

	switch window.current() {
	case startDrill, quit:
		bindings = append(bindings, keys.Select)

	case singleYear, showDoomsday, cheatsheet:
		toggle := keys.Select
		toggle.SetHelp("enter", "Toggle")
		bindings = append(bindings, toggle)

	case startYear, endYear:
		bindings = append(bindings, keys.Change, keys.Jump, keys.Digit)
	}

	return append(bindings, keys.Quit)
}

func getMenuTitle() string {
	art := `
  ▐▌            ▐▌
▗▄▟▙▄▖ ▄   ▄ ▗▞▀▜▌▗▞▀▜▌▄   ▄
  ▐▙▄▖ █▄▄▄█ ▝▚▄▟▌▝▚▄▟▌█   █
                        ▀▀▀█
                       ▄   █
                        ▀▀▀  `

	return art[1:] // Remove leading newline
}

// rows lists the options shown, in order: one year row while drilling a
// single year, a start and an end row otherwise.
func (window *Window) rows() []menuOption {
	years := []menuOption{startYear, endYear}
	if window.Single {
		years = years[:1]
	}

	rows := []menuOption{startDrill, singleYear}
	rows = append(rows, years...)

	return append(rows, showDoomsday, cheatsheet, quit)
}

func (window *Window) items() []menuItem {
	options := window.Options
	rows := window.rows()
	items := make([]menuItem, len(rows))

	for i, row := range rows {
		var label string

		switch row {
		case startDrill:
			label = "Start Drill"

		case singleYear:
			label = "Single Year: " + onOff(window.Single)

		case startYear:
			label = "Start Year: "
			if window.Single {
				label = "Year: "
			}

			label += window.yearLabel(startYear, options.Years.Start)

		case endYear:
			label = "End Year: " + window.yearLabel(endYear, options.Years.End)

		case showDoomsday:
			label = "Doomsday Hint: " + onOff(options.ShowDoomsday)

		case cheatsheet:
			label = "Cheatsheet: " + onOff(options.Cheatsheet)

		case quit:
			label = "Quit"
		}

		items[i] = menuItem{label, row}
	}

	return items
}

// drillOptions are the options a drill started now would use.
func (window *Window) drillOptions() drill.Options {
	options := window.Options
	if window.Single {
		options.Years.End = options.Years.Start
	}

	return options
}

// yearLabel shows a year row's value,
// or what is being typed into it.
func (window *Window) yearLabel(row menuOption, year int) string {
	if window.typing && window.current() == row {
		return window.typed + "_"
	}

	return "◂ " + strconv.Itoa(year) + " ▸"
}

func onOff(on bool) string {
	if on {
		return "On"
	}

	return "Off"
}

// current is the highlighted row's option.
func (window *Window) current() menuOption {
	return window.rows()[window.Choice]
}

func (window *Window) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	rows := window.rows()

	switch {
	case key.Matches(msg, keys.Up):
		window.settle()
		if window.Choice > 0 {
			window.Choice--
		}

	case key.Matches(msg, keys.Down):
		window.settle()
		if window.Choice < len(rows)-1 {
			window.Choice++
		}

	case key.Matches(msg, keys.Left):
		window.adjust(-1)

	case key.Matches(msg, keys.Right):
		window.adjust(1)

	case key.Matches(msg, keys.JumpLeft):
		window.adjust(-10)

	case key.Matches(msg, keys.JumpRight):
		window.adjust(10)

	case key.Matches(msg, keys.Delete):
		window.erase()

	case key.Matches(msg, keys.Select):
		window.settle()

		switch window.current() {
		case startDrill:
			options := window.drillOptions()
			return func() tea.Msg {
				return StartDrillMsg{Options: options}
			}

		case singleYear, showDoomsday, cheatsheet:
			window.adjust(1)

		case quit:
			return tea.Quit
		}

	case key.Matches(msg, keys.Quit):
		return tea.Quit

	case key.Matches(msg, keys.Digit):
		window.typeDigit(msg.String())
	}

	return nil
}

// adjust steps the setting the highlighted row edits.
func (window *Window) adjust(step int) {
	options := &window.Options

	switch window.current() {
	case singleYear:
		window.Single = !window.Single

	case startYear:
		options.Years.Start += step

	case endYear:
		options.Years.End += step

	case showDoomsday:
		options.ShowDoomsday = !options.ShowDoomsday

	case cheatsheet:
		options.Cheatsheet = !options.Cheatsheet
	}

	window.settle()
}

// typeDigit enters a digit into the highlighted year row.
func (window *Window) typeDigit(key string) {
	row := window.current()

	if (row != startYear && row != endYear) || len(key) != 1 || key[0] < '0' || key[0] > '9' {
		return
	}

	if !window.typing {
		window.setTyped(key)
		return
	}

	if len(window.typed) < maxYearDigits {
		window.setTyped(window.typed + key)
	}
}

// erase drops the last digit of
// the highlighted year row.
func (window *Window) erase() {
	row := window.current()
	if row != startYear && row != endYear {
		return
	}

	typed := window.typed
	if !window.typing {
		typed = strconv.Itoa(window.year(row))
	}

	if typed != "" {
		window.setTyped(typed[:len(typed)-1])
	}
}

func (window *Window) setTyped(typed string) {
	window.typing = true
	window.typed = typed

	// Empty or not yet a valid year
	year, _ := strconv.Atoi(typed)

	if window.current() == startYear {
		window.Options.Years.Start = year
		return
	}

	window.Options.Years.End = year
}

func (window *Window) year(row menuOption) int {
	if row == startYear {
		return window.Options.Years.Start
	}

	return window.Options.Years.End
}

// settle ends any typing and pulls the
// years back into a valid range.
func (window *Window) settle() {
	window.typing = false
	window.typed = ""
	years := &window.Options.Years

	years.Start = clampYear(years.Start)
	years.End = clampYear(years.End)

	if years.Start > years.End && !window.Single {
		if window.current() == endYear {
			years.Start = years.End
			return
		}

		years.End = years.Start
	}
}

func clampYear(year int) int {
	return min(max(year, doomsday.MinYear), doomsday.MaxYear)
}
