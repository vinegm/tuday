package drill

import "charm.land/bubbles/v2/key"

// keyMap is the drill's keys.
type keyMap struct {
	Finish, Answer, Delete, Clear key.Binding

	// Weekday is the digits that answer
	// on their own, before anything is typed.
	Weekday key.Binding
}

var keys = keyMap{
	Finish: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Finish"),
	),
	Answer: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Answer"),
	),
	Delete: key.NewBinding(key.WithKeys("backspace")),
	Clear: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "Clear"),
	),
	Weekday: key.NewBinding(
		key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7"),
		key.WithHelp("0-7", "Answer"),
	),
}
