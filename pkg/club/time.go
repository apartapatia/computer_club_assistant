package club

import "time"

const TimeFormat = "15:04"

type WorkingTime struct {
	Open  time.Time
	Close time.Time
}

func NewWorkingTime(open, close time.Time) *WorkingTime {
	return &WorkingTime{
		Open:  open,
		Close: close,
	}
}

func IsTimeWithinWorkingHours(workingTime WorkingTime, time time.Time) bool {
	if workingTime.Open.After(workingTime.Close) {
		return !time.Before(workingTime.Open) || !time.After(workingTime.Close)
	}

	return !time.Before(workingTime.Open) && !time.After(workingTime.Close)
}
