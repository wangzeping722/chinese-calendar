package chinesecalendar

import "time"

func validateDate(t time.Time) (time.Time, bool) {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if t.Year() < minDay.Year() || t.Year() > maxDay.Year() {
		return time.Time{}, false
	}

	return t, true
}
