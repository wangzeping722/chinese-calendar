package main

import (
	"fmt"
	"time"

	"github.com/wangzeping722/chinesecalendar"
)

func main() {
	// 判断是否节假日
	t := time.Date(2022, 2, 25, 0, 0, 0, 0, time.Local)
	fmt.Printf("IsHoliday: %v\n", chinesecalendar.IsHoliday(t))
	fmt.Printf("IsWorkday: %v\n", chinesecalendar.IsWorkday(t))
	// output:
	// IsHoliday: false
	// IsWorkday: true

	t1 := time.Date(2022, 2, 26, 0, 0, 0, 0, time.Local)
	fmt.Printf("IsHoliday: %v\n", chinesecalendar.IsHoliday(t1))
	fmt.Printf("IsWorkday: %v\n", chinesecalendar.IsWorkday(t1))
	// output:
	// IsHoliday: true
	// IsWorkday: false

	// 获取节日名
	t2 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	holiday, isHoliday := chinesecalendar.GetHolidayDetail(t2)
	if isHoliday {
		fmt.Printf("节日: %v\n", holiday.Name())
	}
	// output:
	// 节日: 元旦

	// 判断节日是否调休
	t3 := time.Date(2022, 2, 3, 0, 0, 0, 0, time.Local)
	fmt.Printf("IsInLieu: %v\n", chinesecalendar.IsInLieu(t3))
	// output:
	// IsInLieu: true
}
