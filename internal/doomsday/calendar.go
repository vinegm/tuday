package doomsday

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

// Bounds for the years that
// can be used in the game.
const (
	MinYear = 1
	MaxYear = 9999
)

const secondsPerDay = 24 * 60 * 60

// Years is the inclusive range
// dates are drawn from.
type Years struct {
	Start int
	End   int
}

// Label is the range as shown to the player.
func (years Years) Label() string {
	if years.Start == years.End {
		return strconv.Itoa(years.Start)
	}

	return fmt.Sprintf("%d-%d", years.Start, years.End)
}

// Validate reports an error for a range
// dates can't be drawn from.
func (years Years) Validate() error {
	if years.Start > years.End {
		return fmt.Errorf("start year %d is after end year %d", years.Start, years.End)
	}

	if years.Start < MinYear || years.End > MaxYear {
		return fmt.Errorf("years must be between %d and %d", MinYear, MaxYear)
	}

	return nil
}

// RandomDate picks a random day from the
// given range of years, inclusive.
func RandomDate(rng *rand.Rand, years Years) time.Time {
	first := time.Date(years.Start, time.January, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(years.End, time.December, 31, 0, 0, 0, 0, time.UTC)

	// WARN: time.Duration would overflow past ~292 years.
	days := (last.Unix()-first.Unix())/secondsPerDay + 1

	return first.AddDate(0, 0, int(rng.Int64N(days)))
}

// FormatDate renders a human readable
// date, e.g. "22 of September, 2026".
func FormatDate(date time.Time) string {
	return fmt.Sprintf(
		"%d of %s, %d",
		date.Day(),
		date.Month(),
		date.Year(),
	)
}

// Of returns the weekday that year's doomsday falls on.
func Of(year int) time.Weekday {
	century, rest := year/100, year%100
	anchor := (5*(century%4) + 2) % 7

	return time.Weekday((anchor + rest + rest/4) % 7)
}

// ParseGuess returns the weekday for "3", "wed" or
// "Wednesday", and false if there is no match.
//
//	NOTE: Sunday is accepted as both 0 and 7.
func ParseGuess(text string) (time.Weekday, bool) {
	text = strings.ToLower(strings.TrimSpace(text))

	if text == "" {
		return 0, false
	}

	if isDigits(text) {
		number, err := strconv.Atoi(text)
		if err != nil || number > 7 {
			return 0, false
		}

		return time.Weekday(number % 7), true
	}

	match, matches := time.Sunday, 0

	for day := time.Sunday; day <= time.Saturday; day++ {
		if strings.HasPrefix(strings.ToLower(day.String()), text) {
			match = day
			matches++
		}
	}

	return match, matches == 1
}

// Cheatsheet lists the numbers ParseGuess
// accepts for each weekday, on one line.
func Cheatsheet() string {
	entries := []string{"0/7 Sun"}

	for day := time.Monday; day <= time.Saturday; day++ {
		entries = append(
			entries,
			fmt.Sprintf(
				"%d %s",
				int(day),
				day.String()[:3],
			),
		)
	}

	return strings.Join(entries, "  ")
}

func isDigits(text string) bool {
	for _, char := range text {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
