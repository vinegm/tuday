package report

import (
	"slices"
	"strconv"
	"time"

	"dday/internal/doomsday"
)

// Bucket is one label of a breakdown,
// e.g. "Tuesday" or "1987".
type Bucket struct {
	Label   string
	Correct int
	Asked   int
	Median  float64
}

// Breakdown groups the attempts along
// one dimension of the date.
type Breakdown struct {
	Name    string
	Buckets []Bucket // in the order each label was first asked
}

// Breakdowns are the weekday, month and year profiles of a session.
func Breakdowns(attempts []doomsday.Attempt) []Breakdown {
	return []Breakdown{
		{"weekday", group(attempts, func(date time.Time) string {
			return date.Weekday().String()
		})},
		{"month", group(attempts, func(date time.Time) string {
			return date.Month().String()
		})},
		{"year", group(attempts, func(date time.Time) string {
			return strconv.Itoa(date.Year())
		})},
	}
}

// Worst is the label with the lowest accuracy, and
// false if every bucket was answered perfectly.
func (breakdown Breakdown) Worst() (string, bool) {
	if len(breakdown.Buckets) == 0 {
		return "", false
	}

	worst := breakdown.Buckets[0]

	for _, bucket := range breakdown.Buckets[1:] {
		if compareWorst(bucket, worst) < 0 {
			worst = bucket
		}
	}

	return worst.Label, worst.Correct < worst.Asked
}

// Ranked is the buckets sorted worst-first:
// lowest accuracy, then slowest.
func (breakdown Breakdown) Ranked() []Bucket {
	ranked := slices.Clone(breakdown.Buckets)
	slices.SortStableFunc(ranked, compareWorst)

	return ranked
}

// compareWorst orders A before B when A has the lower
// accuracy, or the same accuracy and a slower median.
func compareWorst(A, B Bucket) int {
	// WARN: Cross-multiplied
	left := A.Correct * B.Asked
	right := B.Correct * A.Asked
	if left != right {
		if left < right {
			return -1
		}

		return 1
	}

	switch {
	case A.Median > B.Median:
		return -1
	case A.Median < B.Median:
		return 1
	default:
		return 0
	}
}

// group buckets the attempts by the label returned by key.
func group(attempts []doomsday.Attempt, key func(time.Time) string) []Bucket {
	var labels []string
	rows := map[string][]doomsday.Attempt{}

	for _, attempt := range attempts {
		label := key(attempt.Date)

		if _, seen := rows[label]; !seen {
			labels = append(labels, label)
		}

		rows[label] = append(rows[label], attempt)
	}

	buckets := make([]Bucket, len(labels))

	for i, label := range labels {
		buckets[i] = Bucket{
			Label:   label,
			Correct: CountCorrect(rows[label]),
			Asked:   len(rows[label]),
			Median:  Median(rows[label]),
		}
	}

	return buckets
}
