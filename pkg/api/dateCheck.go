package api

import (
	"fmt"
	"time"

	"main.go/pkg/db"
)

func checkAndSetDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = formatDate(now)
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format. Expected YYYYMMDD")
	}

	var nextDate string
	if task.Repeat != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("incorrect repetition rule: %v", err)
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = formatDate(now)
		} else {
			task.Date = nextDate
		}
	}

	return nil
}

func afterNow(now, date time.Time) bool {
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dateDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return nowDate.After(dateDate)
}

func formatDate(t time.Time) string {
	return t.Format("20060102")
}
