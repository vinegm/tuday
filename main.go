package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"dday/internal/config"
	"dday/internal/tui"
)

const (
	logFile  = "/tmp/dday.log"
	logLevel = "debug"
)

var version = "dev"

func main() {
	cfg, cfgErr := config.Load(config.DefaultPath())

	debug := flag.Bool("debug", false, "enable debug mode")
	debugFile := flag.String("debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	showVersion := flag.Bool("version", false, "print the version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Println("dday version", version)
		return
	}

	log.SetOutput(io.Discard)
	if *debug {
		file, err := tea.LogToFile(*debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	if cfgErr != nil && !os.IsNotExist(cfgErr) {
		log.Printf("Could not load config: %v", cfgErr)
	}

	program := tea.NewProgram(tui.New(cfg))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
