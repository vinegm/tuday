package config

// Theme is every colour the windows render with.
type Theme struct {
	Fg       string `yaml:"fg"`
	Bg       string `yaml:"bg"`
	Selected string `yaml:"selected"`
	Border   string `yaml:"border"`
	Accent   string `yaml:"accent"` // the date being asked
	Muted    string `yaml:"muted"`  // hints and timers
	Correct  string `yaml:"correct"`
	Wrong    string `yaml:"wrong"`
	KeyHints string `yaml:"key_hints"` // fg of the key hints bar at the bottom, defaults to muted
}

func defaultTheme() Theme {
	return Theme{
		Fg:       "",
		Bg:       "",
		Selected: "#5a5a5a",
		Border:   "",
		Accent:   "",
		Muted:    "#6c6c6c",
		Correct:  "#5faf5f",
		Wrong:    "#d75f5f",
		KeyHints: "",
	}
}
