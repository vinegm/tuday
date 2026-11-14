package theme

import (
	"strings"

	"charm.land/bubbles/v2/key"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/tuday/internal/config"
)

// WindowWidth is the outer width, border included, every window renders at,
// so switching between them doesn't make the box jump.
const WindowWidth = 59

// Theme holds every style the windows render
// with, resolved from the configuration.
type Theme struct {
	config config.Config

	TextStyle             lg.Style
	TitleStyle            lg.Style
	SelectedMenuTextStyle lg.Style
	AccentStyle           lg.Style
	MutedStyle            lg.Style
	CorrectStyle          lg.Style
	WrongStyle            lg.Style

	WindowStyle lg.Style

	keyHintsStyle lg.Style
}

// NewTheme resolves a Theme from the config.
func NewTheme(cfg config.Config) *Theme {
	theme := &Theme{config: cfg}
	theme.setStyles()

	return theme
}

// InnerWidth is the room a window has for its content, inside border and padding.
func (theme *Theme) InnerWidth() int {
	return WindowWidth - theme.WindowStyle.GetHorizontalFrameSize()
}

// Sets the foreground, leaving it inherited when no color is configured.
func withFg(base lg.Style, color string) lg.Style {
	if color == "" {
		return base
	}

	return base.Foreground(lg.Color(color))
}

// Sets the background, leaving it inherited when no color is configured.
func withBg(base lg.Style, color string) lg.Style {
	if color == "" {
		return base
	}

	return base.Background(lg.Color(color))
}

// Updates styles based on the current configuration.
func (theme *Theme) setStyles() {
	colors := theme.config.Theme

	genericStyle := withFg(lg.NewStyle(), colors.Fg)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	theme.TextStyle = lg.NewStyle().Inherit(genericStyle)
	theme.TitleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	theme.SelectedMenuTextStyle = withFg(lg.NewStyle().Inherit(genericMenuTextStyle), colors.Selected)
	theme.AccentStyle = withFg(lg.NewStyle().Inherit(genericMenuTextStyle), colors.Accent)
	theme.MutedStyle = withFg(lg.NewStyle().Inherit(genericStyle), colors.Muted)
	theme.CorrectStyle = withFg(lg.NewStyle().Inherit(genericStyle), colors.Correct)
	theme.WrongStyle = withFg(lg.NewStyle().Inherit(genericStyle), colors.Wrong)

	theme.WindowStyle = withBg(lg.NewStyle(), colors.Bg).
		Width(WindowWidth).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(colors.Bg)).
		BorderForeground(lg.Color(colors.Border)).
		Align(lg.Center)

	keyHints := colors.KeyHints
	if keyHints == "" {
		keyHints = colors.Muted
	}

	theme.keyHintsStyle = withFg(lg.NewStyle(), keyHints)
}

// RenderHints draws the key hints bar from the help of each
// enabled binding, cut to fit in width.
func (theme *Theme) RenderHints(bindings []key.Binding, width int) string {
	parts := make([]string, 0, len(bindings))
	for _, binding := range bindings {
		help := binding.Help()
		if binding.Enabled() && help.Desc != "" {
			parts = append(parts, help.Desc+": "+help.Key)
		}
	}

	return theme.keyHintsStyle.MaxWidth(width).Render(strings.Join(parts, " | "))
}
