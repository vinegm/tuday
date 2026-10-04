<div align="center">
<pre>
.   ▐▌            ▐▌              .
. ▗▄▟▙▄▖ ▄   ▄ ▗▞▀▜▌▗▞▀▜▌▄   ▄    .
.   ▐▙▄▖ █▄▄▄█ ▝▚▄▟▌▝▚▄▟▌█   █    .
.                         ▀▀▀█    .
.                        ▄   █    .
.                         ▀▀▀     .
<strong>Doomsday rule trainer TUI in Go</strong>
</pre>
</div>

tuday is a TUI trainer for the [Doomsday rule](https://en.wikipedia.org/wiki/Doomsday_rule), built with Go and the Bubble Tea TUI framework. It shows you a random date, you name its weekday against the clock, and when you stop it breaks down where you are weakest by weekday, month and year.

## Features

- Terminal UI using Bubble Tea and Lip Gloss
- Random dates from any range of years between 1 and 9999, picked from the menu or the config file
- A live timer per question, a running score and a strip of your recent answers
- Instant answers with `0`-`7`, or type a weekday name or any unambiguous prefix of one
- Optional hints: the year's doomsday next to each question, and a weekday number cheatsheet
- End-of-session results ranking weekdays, months and years worst-first (lowest accuracy, then slowest), with accuracy bars
- Configurable colors (see the Themes directory)

## Getting Started

Prerequisites: Go 1.25 or newer (see `go.mod`)

Clone the repository and run locally:

```bash
go run main.go
```

Or build a standalone binary:

```bash
go build -o tuday ./
./tuday
```

Command line flags:

- `--version` : print the version and exit
- `--debug` : enable debug logging
- `--debug-file <path>` : write debug logs to the specified file (defaults to `/tmp/tuday.log`)
- `--config <path>` : load the config from `<path>` instead of the default location
- `--no-config` : ignore the config file and run with the built-in defaults

## Configuration

Configuration is read at startup from `$XDG_CONFIG_HOME/tuday/config.yaml`, defaulting to `~/.config/tuday/config.yaml`. The file is optional: anything it leaves out keeps its built-in default (see `internal/config/config.go`), an unparsable file falls back to the defaults entirely, and an invalid year range falls back to the default range. An example config with every setting is provided at `Themes/example-config.yaml`.

Configurable items include the year range (or a single year), hints, and the colors of the UI. The year range and hints set the menu's starting values and can still be changed there.

## Controls

Menu:

- `↑`/`↓` or `k`/`j`: Move between items
- `←`/`→` or `h`/`l`: Change the highlighted setting (`shift` or `H`/`L` for steps of 10 years)
- Digits / `Backspace`: Type a year straight into the highlighted year
- `Enter` / `Space`: Select
- `Q`: Quit

Drill:

- `0`-`7`: Answer instantly (`1` = Monday through `6` = Saturday, Sunday is `0` or `7`)
- Letters then `Enter`: Answer by name or an unambiguous prefix (`wed`, `th`, `sa`); a word that isn't a weekday keeps the same date and the clock running
- `Ctrl-U`: Clear the answer
- `Esc`: Finish the session

Results:

- `←`/`→`, `h`/`l` or `Tab`: Switch between the weekday, month and year breakdowns
- `↑`/`↓` or `k`/`j`: Scroll a breakdown too long for the window
- `Enter`: Drill again with the same settings
- `Q` / `Esc`: Back to the menu

`Ctrl-C` quits from anywhere.

## Themes

Theme files are stored in the `Themes` directory; a few presets are included (catppuccin-mocha, gruvbox, tokyo-night). Copy one to `~/.config/tuday/config.yaml`, or merge its colors into your config.

## Development

```bash
make build   # build the binary
make run     # build and run
make test    # run the test suite
make ci      # format check, vet, lint and test
```

Linting and formatting use [golangci-lint](https://golangci-lint.run/welcome/install/).

## License

This project is licensed under the MIT License, see the [LICENSE](LICENSE) file for details.
