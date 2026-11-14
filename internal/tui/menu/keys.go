package menu

import "charm.land/bubbles/v2/key"

// keyMap is the menu's keys.
type keyMap struct {
	Up, Down, Left, Right, JumpLeft, JumpRight key.Binding

	Digit, Delete, Select, Quit key.Binding

	// Move, Done, Change and Jump name the pairs
	// of keys above as one in the key hints bar.
	Move, Done, Change, Jump key.Binding
}

var keys = keyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k")),
	Down:      key.NewBinding(key.WithKeys("down", "j")),
	Left:      key.NewBinding(key.WithKeys("left", "h")),
	Right:     key.NewBinding(key.WithKeys("right", "l")),
	JumpLeft:  key.NewBinding(key.WithKeys("shift+left", "H")),
	JumpRight: key.NewBinding(key.WithKeys("shift+right", "L")),
	Digit: key.NewBinding(
		key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("0-9", "Type"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "Delete"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "Select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	Move: key.NewBinding(
		key.WithKeys("down", "j", "up", "k"),
		key.WithHelp("jk", "Move"),
	),
	Done: key.NewBinding(
		key.WithKeys("down", "j", "up", "k"),
		key.WithHelp("jk", "Done"),
	),
	Change: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("hl", "Change"),
	),
	Jump: key.NewBinding(
		key.WithKeys("shift+left", "H", "shift+right", "L"),
		key.WithHelp("HL", "±10"),
	),
}
