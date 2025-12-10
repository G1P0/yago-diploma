package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateExample = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("missing repeat data")
	}

	date, err := time.Parse(dateExample, dstart)
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
		return date.Format(dateExample), nil
	case "w":
		return nextWeeklyDate(now, date, parts)
	case "m":
		return nextMonthlyDate(now, date, parts)
	default:
		return "", errors.New("unsupported repeat format")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateExample), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateExample, nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid now format")
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := w.Write([]byte(next)); err != nil {
		return
	}
}

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, now.Location())
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

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

func nextWeeklyDate(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid w rule format")
	}

	rawDays := strings.Split(parts[1], ",")
	if len(rawDays) == 0 {
		return "", errors.New("invalid w day list")
	}

	var days []int
	for _, d := range rawDays {
		d = strings.TrimSpace(d)
		if d == "" {
			return "", errors.New("invalid w day")
		}
		val, err := strconv.Atoi(d)
		if err != nil {
			return "", errors.New("invalid w day")
		}
		if val < 1 || val > 7 {
			return "", errors.New("w day out of range")
		}
		days = append(days, val)
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}

		if containsInt(days, weekdayIndex(date)) {
			break
		}
	}

	return date.Format(dateExample), nil
}

func weekdayIndex(t time.Time) int {
	wd := int(t.Weekday())
	return ((wd + 6) % 7) + 1
}

func nextMonthlyDate(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 && len(parts) != 3 {
		return "", errors.New("invalid m rule format")
	}

	rawDays := strings.Split(parts[1], ",")
	if len(rawDays) == 0 {
		return "", errors.New("invalid m day list")
	}

	var (
		days       []int
		useLast    bool
		usePreLast bool
	)

	for _, d := range rawDays {
		d = strings.TrimSpace(d)
		if d == "" {
			return "", errors.New("invalid m day")
		}
		val, err := strconv.Atoi(d)
		if err != nil {
			return "", errors.New("invalid m day")
		}

		switch {
		case val >= 1 && val <= 31:
			days = append(days, val)
		case val == -1:
			useLast = true
		case val == -2:
			usePreLast = true
		default:
			return "", errors.New("m day out of range")
		}
	}

	var (
		months        []int
		restrictMonth bool
	)

	if len(parts) == 3 {
		rawMonths := strings.Split(parts[2], ",")
		if len(rawMonths) == 0 {
			return "", errors.New("invalid m month list")
		}
		for _, m := range rawMonths {
			m = strings.TrimSpace(m)
			if m == "" {
				return "", errors.New("invalid m month value")
			}
			val, err := strconv.Atoi(m)
			if err != nil {
				return "", errors.New("invalid m month value")
			}
			if val < 1 || val > 12 {
				return "", errors.New("m month out of range")
			}
			months = append(months, val)
		}
		restrictMonth = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}

		m := int(date.Month())
		if restrictMonth && !containsInt(months, m) {
			continue
		}

		day := date.Day()
		match := false

		if containsInt(days, day) {
			match = true
		}

		if !match && (useLast || usePreLast) {
			last := lastDayOfMonth(date.Year(), date.Month())
			if useLast && day == last {
				match = true
			}
			if usePreLast && day == last-1 {
				match = true
			}
		}

		if match {
			break
		}
	}

	return date.Format(dateExample), nil
}

func lastDayOfMonth(year int, month time.Month) int {
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Day()
}

func containsInt(list []int, v int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
