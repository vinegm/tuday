package results

import "charm.land/bubbles/v2/key"

// keyMap is the results' keys.
type keyMap struct {
	Next, Previous, Down, Up, Again, Menu key.Binding

	// Breakdown and Scroll name the pairs of keys
	// above as one in the key hints bar.
	Breakdown, Scroll key.Binding
}

var keys = keyMap{
	Next:     key.NewBinding(key.WithKeys("tab", "right", "l")),
	Previous: key.NewBinding(key.WithKeys("shift+tab", "left", "h")),
	Down:     key.NewBinding(key.WithKeys("down", "j")),
	Up:       key.NewBinding(key.WithKeys("up", "k")),
	Again: key.NewBinding(
		key.WithKeys("enter", "space", "r"),
		key.WithHelp("enter", "Again"),
	),
	Menu: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "Menu"),
	),
	Breakdown: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("hl", "Breakdown"),
	),
	Scroll: key.NewBinding(
		key.WithKeys("down", "j", "up", "k"),
		key.WithHelp("jk", "Scroll"),
	),
}
