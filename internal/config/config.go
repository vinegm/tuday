package config

// Config is config.yaml.
// Keys are documented in Themes/example-config.yaml.
type Config struct {
	// Drill settings
	StartYear    int  `yaml:"start_year"`
	EndYear      int  `yaml:"end_year"`
	SingleYear   bool `yaml:"single_year"`   // drill only start_year, ignoring end_year
	ShowDoomsday bool `yaml:"show_doomsday"` // show the year's doomsday alongside each date
	Cheatsheet   bool `yaml:"cheatsheet"`    // show the 0-7 weekday numbers under the prompt

	// Style settings
	FgColor       string `yaml:"fg_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`
	AccentColor   string `yaml:"accent_color"` // the date being asked
	MutedColor    string `yaml:"muted_color"`  // hints, timers and help lines
	CorrectColor  string `yaml:"correct_color"`
	WrongColor    string `yaml:"wrong_color"`
}

// Default is the built-in configuration, used
// as the base the config file is layered over.
func Default() Config {
	return Config{
		StartYear:    1900,
		EndYear:      2099,
		SingleYear:   false,
		ShowDoomsday: false,
		Cheatsheet:   false,

		FgColor:       "",
		BgColor:       "",
		SelectedColor: "#5a5a5a",
		BorderColor:   "",
		AccentColor:   "",
		MutedColor:    "#6c6c6c",
		CorrectColor:  "#5faf5f",
		WrongColor:    "#d75f5f",
	}
}
