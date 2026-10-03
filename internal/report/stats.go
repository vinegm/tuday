package report

import (
	"slices"

	"github.com/vinegm/tuday/internal/doomsday"
)

// CountCorrect is how many attempts were answered right.
func CountCorrect(attempts []doomsday.Attempt) int {
	correct := 0

	for _, attempt := range attempts {
		if attempt.Correct() {
			correct++
		}
	}

	return correct
}

// Mean is the average answer time in seconds.
func Mean(attempts []doomsday.Attempt) float64 {
	total := 0.0

	for _, attempt := range attempts {
		total += attempt.Seconds()
	}

	return total / float64(len(attempts))
}

// Median is the middle answer time in seconds.
func Median(attempts []doomsday.Attempt) float64 {
	times := answerTimes(attempts)
	slices.Sort(times)

	middle := len(times) / 2
	if len(times)%2 == 1 {
		return times[middle]
	}

	return (times[middle-1] + times[middle]) / 2
}

// Slowest is the longest answer time in seconds.
func Slowest(attempts []doomsday.Attempt) float64 {
	return slices.Max(answerTimes(attempts))
}

func answerTimes(attempts []doomsday.Attempt) []float64 {
	times := make([]float64, len(attempts))

	for i, attempt := range attempts {
		times[i] = attempt.Seconds()
	}

	return times
}
