package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func parseDate(value string) (time.Time, bool) {
	if len(value) != 8 {
		return time.Time{}, false
	}
	parsed, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, false
	}
	if parsed.Year() < 1 || parsed.Year() > 9999 {
		return time.Time{}, false
	}
	return parsed, true
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, ok := parseDate(dstart)
	if !ok {
		return "", fmt.Errorf("invalid start date")
	}
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is required")
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "y":
		return yearDate(now, date), nil
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid repeat interval")
		}
		return dayDate(now, date, interval), nil
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid repeat format")
		}
		days, ok := parseMonthDays(parts[1])
		if !ok || len(days) == 0 {
			return "", fmt.Errorf("invalid month days")
		}
		months := allMonths()
		if len(parts) == 3 {
			var okm bool
			months, okm = parseMonths(parts[2])
			if !okm || len(months) == 0 {
				return "", fmt.Errorf("invalid months")
			}
		}
		return monthDate(now, date, days, months), nil
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}
		weekdays, ok := parseWeekDays(parts[1])
		if !ok || len(weekdays) == 0 {
			return "", fmt.Errorf("invalid week days")
		}
		return weekDate(now, date, weekdays), nil
	default:
		return "", fmt.Errorf("unsupported repeat type")
	}
}

func yearDate(now, date time.Time) string {
	boundary := now
	if date.After(now) {
		boundary = date
	}
	candidate := date
	for !candidate.After(boundary) {
		candidate = candidate.AddDate(1, 0, 0)
		if candidate.Year() > 9999 {
			return ""
		}
	}
	return candidate.Format(layout)
}

func dayDate(now, date time.Time, interval int) string {
	boundary := now
	if date.After(now) {
		boundary = date
	}
	candidate := date
	for !candidate.After(boundary) {
		candidate = candidate.AddDate(0, 0, interval)
		if candidate.Year() > 9999 {
			return ""
		}
	}
	return candidate.Format(layout)
}

func weekDate(now, date time.Time, weekdays []int) string {
	boundary := now
	if date.After(now) {
		boundary = date
	}
	candidate := boundary.AddDate(0, 0, 1)
	for i := 0; i < 14; i++ {
		week := int(candidate.Weekday())
		if week == 0 {
			week = 7
		}
		if containsInt(weekdays, week) {
			return candidate.Format(layout)
		}
		candidate = candidate.AddDate(0, 0, 1)
	}
	return ""
}

func monthDate(now, date time.Time, days []int, months []int) string {
	boundary := now
	if date.After(now) {
		boundary = date
	}
	search := make([]time.Time, 0, 48)
	yearLimit := boundary.Year() + 3
	for year := boundary.Year(); year <= yearLimit; year++ {
		for _, month := range months {
			for _, day := range days {
				if day == 0 || day < -31 || day > 31 {
					return ""
				}
				searchDay := day
				if searchDay < 0 {
					searchDay = lastDayOfMonth(year, time.Month(month)) + searchDay + 1
				}
				candidate := safeDate(year, time.Month(month), searchDay)
				if candidate.IsZero() {
					return ""
				}
				if candidate.After(boundary) {
					search = append(search, candidate)
				}
			}
		}
	}
	if len(search) == 0 {
		return ""
	}
	best := search[0]
	for _, candidate := range search {
		if candidate.Before(best) {
			best = candidate
		}
	}
	return best.Format(layout)
}

func parseMonthDays(value string) ([]int, bool) {
	parts := strings.Split(value, ",")
	days := make([]int, 0, len(parts))
	negCount := 0
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, false
		}
		day, err := strconv.Atoi(item)
		if err != nil {
			return nil, false
		}
		if day == 0 || day < -31 || day > 31 {
			return nil, false
		}
		if day < 0 {
			negCount++
		}
		days = append(days, day)
	}
	if negCount > 1 {
		return nil, false
	}
	return days, true
}

func parseMonths(value string) ([]int, bool) {
	parts := strings.Split(value, ",")
	months := make([]int, 0, len(parts))
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, false
		}
		month, err := strconv.Atoi(item)
		if err != nil || month < 1 || month > 12 {
			return nil, false
		}
		months = append(months, month)
	}
	return months, true
}

func parseWeekDays(value string) ([]int, bool) {
	parts := strings.Split(value, ",")
	weekdays := make([]int, 0, len(parts))
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, false
		}
		weekday, err := strconv.Atoi(item)
		if err != nil || weekday < 1 || weekday > 7 {
			return nil, false
		}
		weekdays = append(weekdays, weekday)
	}
	return weekdays, true
}

func allMonths() []int {
	months := make([]int, 0, 12)
	for m := 1; m <= 12; m++ {
		months = append(months, m)
	}
	return months
}

func containsInt(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func lastDayOfMonth(year int, month time.Month) int {
	next := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return next.AddDate(0, 0, -1).Day()
}

func safeDate(year int, month time.Month, day int) time.Time {
	if day < 1 {
		return time.Time{}
	}
	if month < time.January || month > time.December {
		return time.Time{}
	}
	if day > 31 {
		return time.Time{}
	}
	candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if candidate.Month() != month || candidate.Day() != day {
		return time.Time{}
	}
	return candidate
}
