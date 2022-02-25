package chinesecalendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsHoliday(t *testing.T) {
	dates := []time.Time{
		time.Date(2004, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(2017, 5, 30, 0, 0, 0, 0, time.Local),
		time.Date(2022, 10, 6, 0, 0, 0, 0, time.Local),
	}
	for _, date := range dates {
		assert.Equal(t, true, IsHoliday(date))
		assert.Equal(t, false, IsWorkday(date))
	}
}

func TestIsWorkDay(t *testing.T) {
	dates := []time.Time{
		time.Date(2004, 1, 5, 0, 0, 0, 0, time.Local),
		time.Date(2021, 2, 25, 0, 0, 0, 0, time.Local),
		time.Date(2022, 2, 25, 0, 0, 0, 0, time.Local),
	}
	for _, date := range dates {
		assert.Equal(t, true, IsWorkday(date))
		assert.Equal(t, false, IsHoliday(date))
	}
}

func TestGetHolidayDetail(t *testing.T) {
	args := []struct {
		date          time.Time
		expectHoliday Holiday
	}{
		{time.Date(2004, 1, 1, 0, 0, 0, 0, time.Local), NewYearsDay},
		{time.Date(2014, 4, 7, 0, 0, 0, 0, time.Local), TombSweepingDay},
		{time.Date(2022, 9, 10, 0, 0, 0, 0, time.Local), MidAutumnFestival},
	}

	for _, arg := range args {
		holiday, isHoliday := GetHolidayDetail(arg.date)
		assert.Equal(t, true, isHoliday)
		assert.Equal(t, arg.expectHoliday, holiday)
	}
}

func TestOverRangeDate(t *testing.T) {
	dates := []time.Time{
		time.Date(2001, 1, 5, 0, 0, 0, 0, time.Local),
		time.Date(2088, 2, 25, 0, 0, 0, 0, time.Local),
	}

	for _, date := range dates {
		assert.Equal(t, false, IsWorkday(date))
		assert.Equal(t, false, IsHoliday(date))
		assert.Equal(t, false, IsInLieu(date))
	}
}
