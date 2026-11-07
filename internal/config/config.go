package config

// Config is config.yaml.
// Keys are documented in Themes/example-config.yaml.
type Config struct {
	Drill Drill `yaml:"drill"`
	Theme Theme `yaml:"theme"`
}

// Drill is what a new drill starts from, changeable in the menu.
type Drill struct {
	StartYear    int  `yaml:"start_year"`
	EndYear      int  `yaml:"end_year"`
	SingleYear   bool `yaml:"single_year"`   // drill only start_year, ignoring end_year
	ShowDoomsday bool `yaml:"show_doomsday"` // show the year's doomsday alongside each date
	Cheatsheet   bool `yaml:"cheatsheet"`    // show the 0-7 weekday numbers under the prompt
}

// Default is the built-in configuration, used
// as the base the config file is layered over.
func Default() Config {
	return Config{
		Drill: Drill{
			StartYear:    1900,
			EndYear:      2099,
			SingleYear:   false,
			ShowDoomsday: false,
			Cheatsheet:   false,
		},
		Theme: defaultTheme(),
	}
}
