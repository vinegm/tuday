package doomsday

import "time"

// Attempt is one answered question.
type Attempt struct {
	Date    time.Time
	Guess   time.Weekday
	Elapsed time.Duration
}

// Correct reports whether the guess was the date's weekday.
func (attempt Attempt) Correct() bool {
	return attempt.Guess == attempt.Date.Weekday()
}

// Seconds is how long the answer took.
func (attempt Attempt) Seconds() float64 {
	return attempt.Elapsed.Seconds()
}
