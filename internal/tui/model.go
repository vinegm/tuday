package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/tuday/internal/tui/drill"
	"github.com/vinegm/tuday/internal/tui/menu"
	"github.com/vinegm/tuday/internal/tui/results"

	"github.com/vinegm/tuday/internal/config"
	"github.com/vinegm/tuday/internal/doomsday"
	"github.com/vinegm/tuday/internal/report"
	"github.com/vinegm/tuday/internal/tui/theme"
)

// Smallest terminal the app renders correctly in.
//
//	WARN: Must be able to fit any window.
const (
	minWidth  = theme.WindowWidth + 2
	minHeight = 20
)

// Window is the interface that holds methods to handle
// updates and rendering for different windows.
type Window interface {
	Update(msg tea.Msg) tea.Cmd
	View(width, height int) string
}

// Model represents the state of the application.
type Model struct {
	width  int
	height int

	currentWindow Window
	menu          *menu.Window
	drill         *drill.Window
	results       *results.Window

	options drill.Options // of the drill last started, reused by "again"
}

// New creates a new model with initial state.
func New(cfg config.Config) Model {
	th := theme.NewTheme(cfg)
	options := drill.Options{
		Years: doomsday.Years{
			Start: cfg.StartYear,
			End:   cfg.EndYear,
		},
		ShowDoomsday: cfg.ShowDoomsday,
		Cheatsheet:   cfg.Cheatsheet,
	}

	menuWindow := menu.New(th, options, cfg.SingleYear)

	return Model{
		currentWindow: menuWindow,
		menu:          menuWindow,
		drill:         drill.New(th),
		results:       results.New(th),
		options:       options,
	}
}

// Init is the command run on start, none.
func (model Model) Init() tea.Cmd {
	return nil
}

// Update is the tea update loop.
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return model, tea.Quit
		}

	case menu.StartDrillMsg:
		model.options = msg.Options
		model.currentWindow = model.drill
		return model, model.drill.Start(msg.Options)

	case drill.FinishedMsg:
		if !model.finish(msg.Attempts) {
			model.currentWindow = model.menu
			return model, nil
		}

		model.currentWindow = model.results
		return model, nil

	case results.AgainMsg:
		model.currentWindow = model.drill
		return model, model.drill.Start(model.options)

	case results.ExitMsg:
		model.currentWindow = model.menu
		return model, nil
	}

	return model, model.currentWindow.Update(msg)
}

// finish loads the results window with a session. False when nothing was
// answered, so there is nothing to show.
func (model *Model) finish(attempts []doomsday.Attempt) bool {
	if len(attempts) == 0 {
		return false
	}

	model.results.Show(attempts, report.Breakdowns(attempts))

	return true
}

// View handles terminal size checks
// and delegates to the current window.
func (model Model) View() tea.View {
	content := model.viewContent()

	view := tea.NewView(
		lg.Place(
			model.width,
			model.height,
			lg.Center,
			lg.Center,
			content,
		),
	)
	view.AltScreen = true

	return view
}

func (model Model) viewContent() string {
	if model.width < minWidth || model.height < minHeight {
		return fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	}

	// Render centered base panel
	return model.currentWindow.View(model.width, model.height)
}
