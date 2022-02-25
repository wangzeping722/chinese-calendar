package chinesecalendar

import "time"

type validateHolidayFunc func(t time.Time) bool

// IsWorkday 检查是否是工作日
// return false if the t is not in the range from 2004 to 2022
func IsWorkday(t time.Time) bool {
	var isValidate bool
	t, isValidate = validateDate(t)
	if !isValidate {
		return false
	}
	return isWorkday(t)
}

func isWorkday(t time.Time) bool {
	if _, inWorkDay := workdays[t]; inWorkDay {
		return true
	}
	weekday := t.Weekday()
	if _, inHoliday := holidays[t]; !inHoliday && weekday >= 1 && weekday <= 5 {
		return true
	}

	return false
}

// IsHoliday 检查是否节假日
// return false if the t is not in the range from 2004 to 2022
func IsHoliday(t time.Time) bool {
	var isValidate bool
	t, isValidate = validateDate(t)
	if !isValidate {
		return false
	}
	return isHoliday(t)
}

func isHoliday(t time.Time) bool {
	return !isWorkday(t)
}

// IsInLieu 检查是否调休日
// return false if the t is not in the range from 2004 to 2022
func IsInLieu(t time.Time) bool {
	var isValidate bool
	t, isValidate = validateDate(t)
	if !isValidate {
		return false
	}

	return isInLieu(t)
}

func isInLieu(t time.Time) bool {
	_, isInLieuDay := inLieuDays[t]
	return isInLieuDay
}

// GetHolidayDetail 获取节假日详细信息
func GetHolidayDetail(t time.Time) (Holiday, bool) {
	var isValidate bool
	t, isValidate = validateDate(t)
	if !isValidate {
		return Holiday{}, false
	}

	if _, ok := workdays[t]; ok {
		return Holiday{}, false
	}

	if hd, ok := holidays[t]; ok {
		return hd, true
	}
	weekday := t.Weekday()
	return Holiday{}, weekday == 0 || weekday == 6
}

func getDates(start, end time.Time, fn validateHolidayFunc) []time.Time {
	days := int(end.Sub(start)/oneDay + 1)
	list := make([]time.Time, 0)
	for i := 0; i < days; i++ {
		t := start.Add(time.Duration(i * int(oneDay)))
		if fn(t) {
			list = append(list, t)
		}
	}
	return list
}

// GetHolidays 获取时间区间内的节假日（包括起止时间），如果日期不符合，返回空切片
func GetHolidays(start, end time.Time, includeWeekends bool) ([]time.Time, error) {
	var isValidate bool
	start, isValidate = validateDate(start)
	if !isValidate {
		return []time.Time{}, ErrUnSupportDate
	}
	end, isValidate = validateDate(end)
	if !isValidate {
		return []time.Time{}, ErrUnSupportDate
	}

	list := make([]time.Time, 0)
	if includeWeekends {
		list = getDates(start, end, IsHoliday)
	} else {
		days := int(end.Sub(start)/oneDay + 1)
		for i := 0; i < days; i++ {
			t := start.Add(time.Duration(i * int(oneDay)))
			if _, ok := holidays[t]; ok {
				list = append(list, t)
			}
		}
	}
	return list, nil
}

// GetWorkdays 获取时间区间内的（包括起止时间），如果日期不符合，返回空切片
func GetWorkdays(start, end time.Time) ([]time.Time, error) {
	var isValidate bool
	start, isValidate = validateDate(start)
	if !isValidate {
		return []time.Time{}, ErrUnSupportDate
	}
	end, isValidate = validateDate(end)
	if !isValidate {
		return []time.Time{}, ErrUnSupportDate
	}

	return getDates(start, end, isWorkday), nil
}
