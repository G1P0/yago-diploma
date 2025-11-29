package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("missing repeat data")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("cannot convert time to date")
	}

	var interval int

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		interval, err = intervalDays(parts)
		if err != nil {
			return "", err
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid y rule format")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil
	default:
		return "", errors.New("unsupported repeat format")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return d.After(n)
}

func intervalDays(parts []string) (int, error) {
	if len(parts) != 2 {
		return 0, errors.New("invalid d rule format")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid number")
	}

	if interval < 1 || interval > 400 {
		return 0, errors.New("interval out of range")
	}

	return interval, nil

}
