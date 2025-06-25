package utils

import "time"

func ParseTime(t string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"15:04:05",
		"15:04",
		"2006-01-02 15:04",
	}
	for _, layout := range layouts {
		tm, err := time.Parse(layout, t)
		if err == nil {
			return tm, nil
		}
	}
	return time.Time{}, &time.ParseError{
		Value:   t,
		Message: "invalid time format",
	}
}
