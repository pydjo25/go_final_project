package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(startData, nowStr)
	if err != nil {
		http.Error(w, "date and repeat parameters are required", http.StatusBadRequest)
	}

	if date == "" || repeat == "" {
		http.Error(w, "date and repeat parameters are required", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

func contWeekDay(slWeekDays []int, wDay int) bool {
	if wDay == 0 {
		wDay = 7
	} else {
		wDay = wDay
	}

	return slices.Contains(slWeekDays, wDay)
}

func dCount(now, date time.Time, pieces []string) (string, error) {
	log.Println(pieces)
	if len(pieces) < 2 {
		return "", errors.New("repetition days not specified")
	}

	interval, err := strconv.Atoi(pieces[1])
	if err != nil {
		return "", errors.New("error converting repetition string to number")
	}

	if interval == 0 || interval > 400 {
		return "", errors.New("invalid number of recurrence days")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return formatDate(date), nil
}

func yCount(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return formatDate(date), nil
}

func wCount(now, date time.Time, pieces []string) (string, error) {
	if len(pieces) < 2 {
		return "", errors.New("days of the week not specified")
	}

	slDays := strings.Split(pieces[1], ",")
	slWeekDays := make([]int, 0, len(slDays))

	for _, days := range slDays {
		day, err := strconv.Atoi(days)
		if err != nil {
			return "", errors.New("incorrect day of the week")
		}
		if day < 1 || day > 7 {
			return "", errors.New("day of the week must be from 1 to 7")
		}
		slWeekDays = append(slWeekDays, day)
	}

	for !afterNow(date, now) || !contWeekDay(slWeekDays, int(date.Weekday())) {
		date = date.AddDate(0, 0, 1)

		if date.Sub(date).Hours()/24 > 7 {
			break
		}
	}

	return formatDate(date), nil

}

func contMonthDay(slMonthDays []int, date time.Time) bool {
	currencyDay := date.Day()

	for _, day := range slMonthDays {
		if day == currencyDay {
			return true
		}
		if day == -1 && currencyDay == daysInMonth(date.Year(), date.Month()) {
			return true
		}
		if day == -2 && currencyDay == daysInMonth(date.Year(), date.Month())-1 {
			return true
		}
	}
	return false
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func contain(sl []int, val int) bool {
	return slices.Contains(sl, val)
}

func mCount(now, date time.Time, pieces []string) (string, error) {
	if len(pieces) < 2 {
		return "", errors.New("days of the week not specified")
	}

	slDays := strings.Split(pieces[1], ",")
	slMonthDays := make([]int, 0, len(slDays))

	for _, days := range slDays {
		day, err := strconv.Atoi(days)
		if err != nil {
			return "", errors.New("incorrect day of the week")
		}
		if (day < -2 || day == 0 || day > 31) && day != -1 {
			return "", errors.New("day of month must be from 1 to 31 or -1, -2")
		}
		slMonthDays = append(slMonthDays, day)
	}

	var months []int

	if len(pieces) > 2 {
		slMonths := strings.Split(pieces[2], ",")
		months = make([]int, 0, len(slMonths))

		for _, slMonth := range slMonths {
			month, err := strconv.Atoi(slMonth)
			if err != nil {
				return "", errors.New("incorrect month")
			}
			if month < 1 || month > 12 {
				return "", errors.New("month must be in the range from 1 to 12")
			}
			months = append(months, month)
		}
	}
	for {
		if afterNow(date, now) &&
			contMonthDay(slMonthDays, date) &&
			(len(months) == 0 || contain(months, int(date.Month()))) {
			break
		}

		date = date.AddDate(0, 0, 1)

		if date.Year() > now.Year()+1 ||
			(date.Year() == now.Year()+5 && date.Month() > now.Month()) {
			return "", errors.New("could not find a suitable date")
		}
	}
	return formatDate(date), nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(startData, dstart)
	if err != nil {
		return "", errors.New("error parsing the resulting date")
	}

	if repeat == "" {
		return "", errors.New("repetition rule not specified")
	}

	pieces := strings.Split(repeat, " ") //
	if len(pieces) == 0 {
		return "", errors.New("repetition rule cannot be empty")
	}

	switch pieces[0] {
	case "d":
		return dCount(now, date, pieces)
	case "w":
		return wCount(now, date, pieces)
	case "m":
		return mCount(now, date, pieces)
	case "y":
		return yCount(now, date)
	default:
		return "", errors.New("error invalid character")
	}
}
