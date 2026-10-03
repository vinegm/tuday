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

type flags struct {
	debugFile   string
	debug       bool
	showVersion bool
}

func main() {
	cfg, cfgErr := config.Load(config.DefaultPath())

	opts := parseFlags()

	if opts.showVersion {
		fmt.Println("dday version", version)
		return
	}

	log.SetOutput(io.Discard)
	if opts.debug {
		file, err := tea.LogToFile(opts.debugFile, logLevel)
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

func parseFlags() flags {
	var opts flags

	flag.BoolVar(&opts.debug, "debug", false, "enable debug mode")
	flag.StringVar(&opts.debugFile, "debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	flag.BoolVar(&opts.showVersion, "version", false, "print the version and exit")
	flag.Parse()

	return opts
}
