package main

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/wangzeping722/chinesecalendar"
	. "github.com/wangzeping722/chinesecalendar/internal"
)

const dateTypeWorkday = 1
const dateTypeHoliday = 2
const dateTypeInLieu = 3

type timeList []time.Time

func (l timeList) Len() int {
	return len(l)
}

func (l timeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l timeList) Less(i, j int) bool {
	return l[i].Before(l[j])
}

type arrangement struct {
	Holidays        map[time.Time]chinesecalendar.Holiday
	Workdays        map[time.Time]chinesecalendar.Holiday
	InLieuDays      map[time.Time]chinesecalendar.Holiday
	HolidayList     timeList
	WorkdayList     timeList
	InLieuDayList   timeList
	HolidayFieldMap map[chinesecalendar.Holiday]string
	MaxDay          time.Time
	MinDay          time.Time

	year    int
	month   int
	day     int
	holiday chinesecalendar.Holiday
	dayType int
}

func newArragement() *arrangement {
	return &arrangement{
		Holidays:   make(map[time.Time]chinesecalendar.Holiday),
		Workdays:   make(map[time.Time]chinesecalendar.Holiday),
		InLieuDays: make(map[time.Time]chinesecalendar.Holiday),
		HolidayFieldMap: map[chinesecalendar.Holiday]string{
			chinesecalendar.NewYearsDay:        "NewYearsDay",
			chinesecalendar.SpringFestival:     "SpringFestival",
			chinesecalendar.TombSweepingDay:    "TombSweepingDay",
			chinesecalendar.LabourDay:          "LabourDay",
			chinesecalendar.DragonBoatFestival: "DragonBoatFestival",
			chinesecalendar.NationalDay:        "NationalDay",
			chinesecalendar.MidAutumnFestival:  "MidAutumnFestival",
		},
		MaxDay: time.Time{},
		MinDay: Date(2099, 1, 1),
	}
}

func (ag *arrangement) generateHolidays() {
	v := reflect.ValueOf(ag)
	for i := 0; i < v.Type().NumMethod(); i++ {
		method := v.Type().Method(i)
		if strings.HasPrefix(method.Name, "Y") {
			v.MethodByName(method.Name).Call(nil)
		}
	}

	fn := func(m map[time.Time]chinesecalendar.Holiday, list *timeList) {
		for k := range m {
			if k.Before(ag.MinDay) {
				ag.MinDay = k
			}
			if k.After(ag.MaxDay) {
				ag.MaxDay = k
			}
			*list = append(*list, k)
		}
	}
	fn(ag.Holidays, &ag.HolidayList)
	fn(ag.Workdays, &ag.WorkdayList)
	fn(ag.InLieuDays, &ag.InLieuDayList)
	sort.Sort(ag.HolidayList)
	sort.Sort(ag.WorkdayList)
	sort.Sort(ag.InLieuDayList)
}

func (ag *arrangement) Y2022() {
	// http://www.gov.cn/zhengce/content/2021-10/25/content_5644835.htm
	// ???????????????2022???1???1??????3???????????????3??????
	// ???????????????1???31??????2???6?????????????????????7??????1???29?????????????????????1???30???????????????????????????
	// ??????????????????4???3??????5?????????????????????3??????4???2???????????????????????????
	// ??????????????????4???30??????5???4?????????????????????5??????4???24?????????????????????5???7???????????????????????????
	// ??????????????????6???3??????5???????????????3??????
	// ??????????????????9???10??????12???????????????3??????
	// ??????????????????10???1??????7?????????????????????7??????10???8?????????????????????10???9???????????????????????????
	ag.yearAt(2022).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(1, 31).to(2, 6).work(1, 29).work(1, 30).inLieu(2, 3).to(2, 4).
		tsd().rest(4, 3).to(4, 5).work(4, 2).inLieu(4, 4).
		ld().rest(4, 30).to(5, 4).work(4, 24).work(5, 7).inLieu(5, 3).to(5, 4).
		dbf().rest(6, 3).to(6, 5).
		maf().rest(9, 10).to(9, 12).
		nd().rest(10, 1).to(10, 7).work(10, 8).work(10, 9).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2021() {
	// http://www.gov.cn/zhengce/content/2020-11/25/content_5564127.htm
	// ???????????????2021???1???1??????3???????????????3??????
	// ???????????????2???11??????17?????????????????????7??????2???7?????????????????????2???20???????????????????????????
	// ??????????????????4???3??????5?????????????????????3??????
	// ??????????????????5???1??????5?????????????????????5??????4???25?????????????????????5???8???????????????????????????
	// ??????????????????6???12??????14???????????????3??????
	// ??????????????????9???19??????21?????????????????????3??????9???18???????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???26?????????????????????10???9???????????????????????????
	ag.yearAt(2021).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(2, 11).to(2, 17).work(2, 7).work(2, 20).inLieu(2, 16).to(2, 17).
		tsd().rest(4, 3).to(4, 5).
		ld().rest(5, 1).to(5, 5).work(4, 25).work(5, 8).inLieu(5, 4).to(5, 5).
		dbf().rest(6, 12).to(6, 14).
		maf().rest(9, 19).to(9, 21).work(9, 18).inLieu(9, 20).
		nd().rest(10, 1).to(10, 7).work(9, 26).work(10, 9).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2020() {
	// http://www.gov.cn/zhengce/content/2019-11/21/content_5454164.htm
	// ???????????????2020???1???1???????????????1??????
	// ???????????????1???24??????30?????????????????????7??????1???19?????????????????????2???1???????????????????????????
	// ??????????????????4???4??????6?????????????????????3??????
	// ??????????????????5???1??????5?????????????????????5??????4???26?????????????????????5???9???????????????????????????
	// ??????????????????6???25??????27?????????????????????3??????6???28???????????????????????????
	// ??????????????????????????????10???1??????8?????????????????????8??????9???27?????????????????????10???10???????????????????????????

	// http://www.gov.cn/zhengce/content/2020-01/27/content_5472352.htm
	// ??????2020??????????????????2???2??????????????????????????????????????????2???3????????????????????????????????????
	ag.yearAt(2020).
		nyd().rest(1, 1).
		sf().rest(1, 24).to(2, 2).work(1, 19).inLieu(1, 29).
		tsd().rest(4, 4).to(4, 6).
		ld().rest(5, 1).to(5, 5).work(4, 26).work(5, 9).inLieu(5, 4).to(5, 5).
		dbf().rest(6, 25).to(6, 27).work(6, 28).inLieu(6, 26).
		nd().rest(10, 1).to(10, 8).work(9, 27).work(10, 10).inLieu(10, 7).to(10, 8)
}

func (ag *arrangement) Y2019() {
	// http://www.gov.cn/xinwen/2018-12/06/content_5346287.htm
	// ???????????????12???30??????1???1???????????????3?????? 12???29???????????????????????????
	// ???????????????2???4??????10?????????????????????7??????2???2?????????????????????2???3???????????????????????????
	// ??????????????????4???5??????????????????????????????
	// ??????????????????5???1???????????????1??????
	// ??????????????????6???7??????????????????????????????
	// ??????????????????9???13??????????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???29?????????????????????10???12????????????????????????

	// http://www.gov.cn/zhengce/content/2019-03/22/content_5375877.htm
	// ?????????????????????????????????2019???????????????????????????????????????
	// 2019???5???1??????4?????????????????????4??????4???28?????????????????????5???5???????????????????????????
	ag.yearAt(2019).
		nyd().rest(1, 1).
		sf().rest(2, 4).to(2, 10).work(2, 2).to(2, 3).inLieu(2, 4).inLieu(2, 8).
		tsd().rest(4, 5).to(4, 7).
		ld().rest(5, 1).to(5, 4).work(4, 28).work(5, 5).inLieu(5, 2).inLieu(5, 3).
		dbf().rest(6, 7).to(6, 9).
		maf().rest(9, 13).to(9, 15).
		nd().rest(10, 1).to(10, 7).work(9, 29).work(10, 12).inLieu(10, 4).inLieu(10, 7)
}

func (ag *arrangement) Y2018() {
	// http://www.gov.cn/zhengce/content/2017-11/30/content_5243579.htm
	// ???????????????1???1??????????????????????????????
	// ???????????????2???15??????21?????????????????????7??????2???11?????????????????????2???24???????????????????????????
	// ??????????????????4???5??????7?????????????????????3??????4???8???????????????????????????
	// ??????????????????4???29??????5???1?????????????????????3??????4???28???????????????????????????
	// ??????????????????6???18??????????????????????????????
	// ??????????????????9???24??????????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???29?????????????????????9???30???????????????????????????
	ag.yearAt(2018).
		nyd().rest(1, 1).
		sf().rest(2, 15).to(2, 21).work(2, 11).work(2, 24).inLieu(2, 19).to(2, 21).
		tsd().rest(4, 5).to(4, 7).work(4, 8).inLieu(4, 6).
		ld().rest(4, 29).to(5, 1).work(4, 28).inLieu(4, 30).
		dbf().rest(6, 18).
		nd().rest(10, 1).to(10, 7).work(9, 29).to(9, 30).inLieu(10, 4).to(10, 5).
		maf().rest(9, 24).
		nyd().rest(12, 30).to(12, 31).work(12, 29).inLieu(12, 31)
}

func (ag *arrangement) Y2017() {
	// http://www.gov.cn/zhengce/content/2016-12/01/content_5141603.htm
	// ???????????????1???1????????????1???2???????????????????????????
	// ???????????????1???27??????2???2?????????????????????7??????1???22?????????????????????2???4???????????????????????????
	// ??????????????????4???2??????4?????????????????????3??????4???1???????????????????????????
	// ??????????????????5???1??????????????????????????????
	// ??????????????????5???28??????30?????????????????????3??????5???27???????????????????????????
	// ??????????????????????????????10???1??????8?????????????????????8??????9???30???????????????????????????
	ag.yearAt(2017).
		nyd().rest(1, 1).to(1, 2).
		sf().rest(1, 27).to(2, 2).work(1, 22).work(2, 4).inLieu(2, 1).to(2, 2).
		tsd().rest(4, 2).to(4, 4).work(4, 1).inLieu(4, 3).
		ld().rest(5, 1).
		dbf().rest(5, 28).to(5, 30).work(5, 27).inLieu(5, 29).
		nd().rest(10, 1).to(10, 8).work(9, 30).inLieu(10, 6).
		maf().rest(10, 4) // ??????????????????????????????10???4????????????
}

func (ag *arrangement) Y2016() {
	// http://www.gov.cn/zhengce/content/2015-12/10/content_10394.htm
	// ???????????????1???1??????????????????????????????
	// ???????????????2???7??????13?????????????????????7??????2???6?????????????????????2???14???????????????????????????
	// ??????????????????4???4??????????????????????????????
	// ??????????????????5???1????????????5???2???????????????????????????
	// ??????????????????6???9??????11?????????????????????3??????6???12???????????????????????????
	// ??????????????????9???15??????17?????????????????????3??????9???18???????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????10???8?????????????????????10???9???????????????????????????
	ag.yearAt(2017).
		nyd().rest(1, 1).
		sf().rest(2, 7).to(2, 13).work(2, 6).work(2, 14).inLieu(2, 11).to(2, 12).
		tsd().rest(4, 4).
		ld().rest(5, 1).to(5, 2).
		dbf().rest(6, 9).to(6, 11).work(6, 12).inLieu(6, 10).
		maf().rest(9, 15).to(9, 17).work(9, 18).inLieu(9, 16).
		nd().rest(10, 1).to(10, 7).work(10, 8).to(10, 9).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2015() {
	// http://www.gov.cn/zhengce/content/2014-12/16/content_9302.htm
	// ???????????????1???1??????3?????????????????????3??????1???4???????????????????????????
	// ???????????????2???18??????24?????????????????????7??????2???15?????????????????????2???28???????????????????????????
	// ??????????????????4???5????????????4???6???????????????????????????
	// ??????????????????5???1??????????????????????????????
	// ??????????????????6???20????????????6???22???????????????????????????
	// ??????????????????9???27????????????
	// ??????????????????10???1??????7?????????????????????7??????10???10???????????????????????????
	ag.yearAt(2015).
		nyd().rest(1, 1).to(1, 3).work(1, 4).inLieu(1, 2).
		sf().rest(2, 18).to(2, 24).work(2, 15).work(2, 28).inLieu(2, 23).to(2, 24).
		tsd().rest(4, 5).to(4, 6).
		ld().rest(5, 1).
		dbf().rest(6, 20).rest(6, 22).
		maf().rest(9, 27).
		nd().rest(10, 1).to(10, 7).work(10, 10).inLieu(10, 7).
		afd().rest(9, 3).to(9, 4).work(9, 6).inLieu(9, 4)
}

func (ag *arrangement) Y2014() {
	// http://www.gov.cn/zwgk/2013-12/11/content_2546204.htm
	// ???????????????1???1?????????1??????
	// ???????????????1???31??????2???6?????????????????????7??????1???26?????????????????????2???8???????????????????????????
	// ??????????????????4???5????????????4???7???????????????????????????
	// ??????????????????5???1??????3?????????????????????3??????5???4???????????????????????????
	// ??????????????????6???2??????????????????????????????
	// ??????????????????9???8??????????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???28?????????????????????10???11???????????????????????????
	ag.yearAt(2014).
		nyd().rest(1, 1).
		sf().rest(1, 31).to(2, 6).work(1, 26).work(2, 8).inLieu(2, 5).to(2, 6).
		tsd().rest(4, 5).to(4, 7).
		ld().rest(5, 1).to(5, 3).work(5, 4).inLieu(5, 2).
		dbf().rest(6, 2).
		maf().rest(9, 8).
		nd().rest(10, 1).to(10, 7).work(9, 28).work(10, 11).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2013() {
	// http://www.gov.cn/zwgk/2012-12/10/content_2286598.htm
	// ???????????????1???1??????3?????????????????????3??????1???5?????????????????????1???6???????????????????????????
	// ???????????????2???9??????15?????????????????????7??????2???16?????????????????????2???17???????????????????????????
	// ??????????????????4???4??????6?????????????????????3??????4???7???????????????????????????
	// ??????????????????4???29??????5???1?????????????????????3??????4???27?????????????????????4???28???????????????????????????
	// ??????????????????6???10??????12?????????????????????3??????6???8?????????????????????6???9???????????????????????????
	// ??????????????????9???19??????21?????????????????????3??????9???22???????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???29?????????????????????10???12???????????????????????????
	ag.yearAt(2013).
		nyd().rest(1, 1).to(1, 3).work(1, 5).to(1, 6).inLieu(1, 2).to(1, 3).
		sf().rest(2, 9).to(2, 15).work(2, 16).to(2, 17).inLieu(2, 14).to(2, 15).
		tsd().rest(4, 4).to(4, 6).work(4, 7).inLieu(4, 5).
		ld().rest(4, 29).to(5, 1).work(4, 27).to(4, 28).inLieu(4, 29).to(4, 30).
		dbf().rest(6, 10).to(6, 12).work(6, 8).to(6, 9).inLieu(6, 10).to(6, 11).
		maf().rest(9, 19).to(9, 21).work(9, 22).inLieu(9, 20).
		nd().rest(10, 1).to(10, 7).work(9, 29).work(10, 12).inLieu(10, 4).inLieu(10, 7)
}

func (ag *arrangement) Y2012() {
	// http://www.gov.cn/zwgk/2011-12/06/content_2012097.htm
	// ???????????????2012???1???1??????3?????????????????????3??????2011???12???31???????????????????????????
	// ???????????????1???22??????28?????????????????????7??????1???21?????????????????????1???29???????????????????????????
	// ??????????????????4???2??????4?????????????????????3??????3???31?????????????????????4???1???????????????????????????
	// ??????????????????4???29??????5???1?????????????????????3??????4???28???????????????????????????
	// ??????????????????6???22??????24?????????????????????3??????
	// ??????????????????????????????9???30??????10???7?????????????????????8??????9???29???????????????????????????
	// ???????????????????????????????????????????????? ??( ?? ??? ??|||)???
	ag.yearAt(2012).
		nyd().rest(1, 1).to(1, 3).inLieu(1, 3).
		sf().rest(1, 22).to(1, 28).work(1, 21).work(1, 29).inLieu(1, 26).to(1, 27).
		tsd().rest(4, 2).to(4, 4).work(3, 31).work(4, 1).inLieu(4, 2).to(4, 3).
		ld().rest(4, 29).to(5, 1).work(4, 28).inLieu(4, 30).
		dbf().rest(6, 22).rest(6, 24).
		maf().rest(9, 30).
		nd().rest(10, 1).to(10, 7).work(9, 29).inLieu(10, 5)
}

func (ag *arrangement) Y2011() {
	// http://www.gov.cn/zwgk/2010-12/10/content_1762643.htm
	// ???????????????1???1??????3?????????????????????3??????
	// ???????????????2???2????????????????????????8?????????????????????7??????1???30?????????????????????2???12???????????????????????????
	// ??????????????????4???3??????5?????????????????????3??????4???2???????????????????????????
	// ??????????????????4???30??????5???2?????????????????????3??????
	// ??????????????????6???4??????6?????????????????????3??????
	// ??????????????????9???10??????12?????????????????????3??????
	// ??????????????????10???1??????7?????????????????????7??????10???8?????????????????????10???9???????????????????????????
	// ????????????????????????????????????????????????????????? ??( ?? ??? ??|||)???
	ag.yearAt(2011).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(2, 2).to(2, 8).work(1, 30).work(2, 12).inLieu(2, 7).to(2, 8).
		tsd().rest(4, 3).to(4, 5).work(4, 2).inLieu(4, 4).
		ld().rest(4, 30).to(5, 2).
		dbf().rest(6, 4).rest(6, 6).
		maf().rest(9, 10).to(9, 12).
		nd().rest(10, 1).to(10, 7).work(10, 8).to(10, 9).inLieu(10, 6).to(10, 7).
		nyd().work(12, 31)
}

func (ag *arrangement) Y2010() {
	// http://www.gov.cn/zwgk/2009-12/08/content_1482691.htm
	// ???????????????1???1??????3?????????????????????3??????
	// ???????????????2???13??????19?????????????????????7??????2???20?????????????????????21???????????????????????????
	// ??????????????????4???3??????5?????????????????????3??????
	// ??????????????????5???1??????3?????????????????????3??????
	// ??????????????????6???14??????16?????????????????????3??????6???12?????????????????????13???????????????????????????
	// ??????????????????9???22??????24?????????????????????3??????9???19?????????????????????25???????????????????????????
	// ??????????????????10???1??????7?????????????????????7??????9???26?????????????????????10???9???????????????????????????
	ag.yearAt(2010).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(2, 13).to(2, 19).work(2, 20).to(2, 21).inLieu(2, 18).to(2, 19).
		tsd().rest(4, 3).to(4, 5).
		ld().rest(5, 1).to(5, 3).
		dbf().rest(6, 14).to(6, 16).work(6, 12).to(6, 13).inLieu(6, 14).to(6, 15).
		maf().rest(9, 22).to(9, 24).work(9, 19).work(9, 25).inLieu(9, 23).to(9, 24).
		nd().rest(10, 1).to(10, 7).work(9, 26).work(10, 9).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2009() {
	// http://www.gov.cn/zwgk/2008-12/10/content_1174014.htm
	// ???????????????1???1??????3???????????????3??????
	// ?????????1???1????????????????????????????????????????????????1???3?????????????????????????????????
	// 1???4?????????????????????????????????1???2?????????????????????
	// 1???4???????????????????????????
	// ???????????????1???25??????31???????????????7??????
	// ?????????1???25????????????????????????????????????1???26??????????????????????????????????????????1???27????????????????????????????????????????????????????????????1???31?????????????????????????????????1???25?????????????????????????????????1???28?????????????????????1???24?????????????????????2???1
	// ???????????????????????????????????????1???29?????????????????????1???30?????????????????????
	// 1???24?????????????????????2???1???????????????????????????
	// ??????????????????4???4??????6???????????????3??????
	// ?????????4???4????????????????????????????????????????????????????????????4???5?????????????????????????????????
	// 4???4?????????????????????????????????4???6?????????????????????
	// ??????????????????5???1??????3???????????????3??????
	// ?????????5???1?????????????????????????????????????????????????????????????????????5???2?????????????????????5???3?????????????????????????????????
	// ??????????????????5???28??????30???????????????3??????
	// ?????????5???28????????????????????????????????????????????????????????????5???30?????????????????????????????????5???31?????????????????????????????????5???29?????????????????????
	// 5???31???????????????????????????
	// ??????????????????????????????10???1??????8???????????????8??????
	// ?????????10???1?????????????????????10???2?????????????????????10???3????????????????????????????????????????????????10???4?????????????????????????????????10???3???????????????????????????????????????????????????10???5?????????????????????10???6?????????????????????9???27?????????????????????10
	// ???10?????????????????????????????????10???7?????????????????????10???8?????????????????????
	// 9???27?????????????????????10???10???????????????????????????
	ag.yearAt(2009).
		nyd().rest(1, 1).to(1, 3).work(1, 4).inLieu(1, 2).
		sf().rest(1, 25).to(1, 31).work(1, 24).work(2, 1).inLieu(1, 29).to(1, 30).
		tsd().rest(4, 4).to(4, 6).
		ld().rest(5, 1).to(5, 3).
		dbf().rest(5, 28).to(5, 30).work(5, 31).inLieu(5, 29).
		nd().rest(10, 1).to(10, 8).work(9, 27).work(10, 10).inLieu(10, 7).to(10, 8).
		maf().rest(10, 3) // ??????????????????????????????10???3????????????
}

func (ag *arrangement) Y2008() {
	// """ http://www.gov.cn/zwgk/2007-12/18/content_837184.htm
	// ???????????????2007???12???30??????2008???1???1???????????????3??????
	// ?????????1???1???????????????????????????????????????12???30?????????????????????????????????12???29?????????????????????????????????12???31?????????????????????12???29???????????????????????????
	// ???????????????2???6??????12????????????????????????????????????????????????7??????
	// ?????????2???6??????????????????2???7??????????????????2???8??????????????????????????????????????????2???9?????????????????????2???10?????????????????????????????????2???2?????????????????????2???3???????????????????????????????????????2???11?????????????????????2???12?????????????????????2???2
	// ?????????????????????2???3???????????????????????????
	// ??????????????????4???4??????6???????????????3??????
	// ?????????4???4???????????????????????????????????????4???5?????????????????????4???6?????????????????????????????????
	// ????????????????????????????????????5???1??????3???????????????3??????
	// ?????????5???1????????????????????????5???3?????????????????????????????????5???4?????????????????????????????????5???2?????????????????????5???4???????????????????????????
	// ??????????????????6???7??????9???????????????3??????
	// ?????????6???7?????????????????????????????????6???8????????????????????????????????????????????????????????????6???8?????????????????????????????????6???9?????????????????????
	// ??????????????????9???13??????15???????????????3??????
	// ?????????9???13?????????????????????????????????9???14????????????????????????????????????????????????????????????9???14?????????????????????????????????9???15?????????????????????
	// ??????????????????9???29??????10???5???????????????7??????
	// ?????????10???1??????2??????3????????????????????????9???27?????????????????????9???28???????????????????????????????????????9???29?????????????????????30?????????????????????10???4?????????????????????5?????????????????????????????????
	// ??????????????????????????????????????????????????????
	ag.yearAt(2008).
		nyd().rest(1, 1).
		sf().rest(2, 6).to(2, 12).work(2, 2).to(2, 3).inLieu(2, 11).to(2, 12).
		tsd().rest(4, 4).to(4, 6).
		ld().rest(5, 1).to(5, 3).work(5, 4).inLieu(5, 2).
		dbf().rest(6, 7).to(6, 9).
		maf().rest(9, 13).to(9, 15).
		nd().rest(9, 29).to(10, 5).work(9, 27).to(9, 28).inLieu(9, 29).to(9, 30)
}

func (ag *arrangement) Y2007() {
	//	http://www.gov.cn/fwxx/sh/2006-12/18/content_471877.htm
	//	??????????????? 1???1??????3????????????????????????
	//	??????1???1????????????????????????2006???12???30?????????????????????31?????????????????????????????????????????????2007???1???2??????3??????2006???12???30?????????????????????12???31???????????????????????????
	//	???????????????2???18??????24?????????????????????????????????????????????7??????
	//	??????18??????19??????20????????????????????????17?????????????????????18?????????????????????25?????????????????????????????????????????????21?????????????????????22?????????????????????23?????????????????????24?????????????????????????????????17??????25????????????
	//	?????????????????????5???1??????7???????????????7??????
	//	?????????1??????2??????3????????????????????????4???28?????????????????????29???????????????????????????????????????5???4?????????????????????7?????????????????????5???5?????????????????????6?????????????????????????????????4???28??????29????????????
	//	?????????????????????10???1??????7???????????????7??????
	//	?????????1??????2??????3????????????????????????9???29?????????????????????30???????????????????????????????????????10???4?????????????????????5?????????????????????10???6?????????????????????7?????????????????????????????????9???29??????30????????????
	//	???????????????????????????????????????????????????????????????????????????????????????????????????
	ag.yearAt(2007).
		nyd().rest(1, 1).to(1, 3).inLieu(1, 2).to(1, 3).
		sf().rest(2, 18).to(2, 24).work(2, 17).work(2, 25).inLieu(2, 22).to(2, 23).
		ld().rest(5, 1).to(5, 7).work(4, 28).to(4, 29).inLieu(5, 4).inLieu(5, 7).
		nd().rest(10, 1).to(10, 7).work(9, 29).to(9, 30).inLieu(10, 4).to(10, 5).
		nyd().rest(12, 30).to(12, 31).work(12, 29).inLieu(12, 31)
}

func (ag *arrangement) Y2006() {
	//	http://www.gov.cn/jrzg/2005-12/22/content_133837.htm
	//	???????????????1???1??????3???????????????3??????
	//	??????1???1????????????????????????12???31???(?????????)???1???1???(?????????)?????????????????????1???2???(?????????)???3???(?????????)???12???31???(?????????)?????????
	//	???????????????1???29??????2???4???(??????????????????????????????)????????????7??????
	//	?????????29??????30??????31????????????????????????1???28???(?????????)???29???(?????????)???2???5???(?????????)?????????????????????2???1???(?????????)???2???(?????????)???3???(?????????)???2???4???(?????????)???????????????1???28??????2???5????????????
	//	?????????????????????5???1??????7???????????????7??????
	//	?????????1??????2??????3????????????????????????4???29???(?????????)???30???(?????????)?????????????????????5???4???(?????????)???5???(?????????)???5???6???(?????????)???7???(?????????)???????????????4???29??????30????????????
	//	?????????????????????10???1??????7???????????????7??????
	//	?????????1??????2??????3????????????????????????9???30???(?????????)???10???1???(?????????)???8???(?????????)?????????????????????10???4???(?????????)???5???(?????????)???6???(?????????)???10???7???(?????????)???????????????9???30??????10???8????????????

	//	???????????????????????????????????????????????????????????????????????????????????????????????????
	ag.yearAt(2006).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(1, 29).to(2, 4).work(1, 28).work(2, 5).inLieu(2, 2).to(2, 3).
		ld().rest(5, 1).to(5, 7).work(4, 29).to(4, 30).inLieu(5, 4).to(5, 5).
		nd().rest(10, 1).to(10, 7).work(9, 30).work(10, 8).inLieu(10, 5).to(10, 6).
		nyd().work(12, 30).to(12, 31)
}

func (ag *arrangement) Y2005() {
	//	https://zhidao.baidu.com/question/2299098.html
	//	???????????????????????????????????????2005???????????????????????????????????????????????????????????????????????????????????????
	//	???????????????1???1??????3???????????????3????????????1???1????????????????????????1???1???(?????????)???????????????1???3???(?????????)???1???2???(?????????)???????????????
	//	???????????????2???9??????15???(???????????????????????????)????????????7???????????????9??????10??????11?????????????????????
	//	2???12???(?????????)???13???(?????????)??????????????????2???5???(?????????)???6???(?????????)?????????????????????2???14???(?????????)???15???(?????????)???
	//	2???5??????6????????????
	//	?????????????????????5???1??????7???????????????7???????????????1??????2??????3????????????????????????4???30???(?????????)???5???1???(?????????)???8???(?????????)???????????????
	//	??????5???4???(?????????)???5???(?????????)???6???(?????????)???5???7???(?????????)???????????????4???30??????5???8????????????
	//	?????????????????????10???1??????7???????????????7???????????????1??????2??????3????????????????????????10???1???(?????????)???2???(?????????)???????????????
	//	??????10???4???(?????????)???5???(?????????)???10???8???(?????????)???9???(?????????)?????????????????????10???6???(?????????)???7???(?????????)???10???8??????9????????????
	ag.yearAt(2005).
		nyd().rest(1, 1).to(1, 3).
		sf().rest(2, 9).to(2, 15).work(2, 5).to(2, 6).inLieu(2, 14).to(2, 15).
		ld().rest(5, 1).to(5, 7).work(4, 30).work(5, 8).inLieu(5, 5).to(5, 6).
		nd().rest(10, 1).to(10, 7).work(10, 8).to(10, 9).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) Y2004() {
	//	https://zh.wikisource.org/zh-hans/????????????????????????2004?????????????????????????????????
	//	????????????????????????????????????????????????????????????????????????????????????
	//		??????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????2004???
	//	??????????????????????????????????????????????????????????????????????????????????????????
	//	???????????????1???1????????????
	//	???????????????1???22????????????28???????????????????????????????????????????????????7??????
	//		?????????22??????23??????24?????????????????????1???25????????????????????????????????????1???17?????????????????????18?????????????????????24?????????????????????????????????
	//		??????1???26?????????????????????27?????????????????????28?????????????????????1???17??????18????????????
	//	?????????????????????5???1????????????7???????????????7??????
	//		?????????1??????2??????3????????????????????????5???1?????????????????????2???????????????????????????????????????5???4?????????????????????5?????????????????????
	//		5???8?????????????????????5???9???????????????????????????????????????5???6?????????????????????7?????????????????????5???8??????9????????????
	//	?????????????????????10???1????????????7???????????????7??????
	//		?????????1??????2??????3????????????????????????10???2?????????????????????3???????????????????????????????????????10???4?????????????????????5?????????????????????
	//		10???9?????????????????????10???????????????????????????????????????10???6?????????????????????7?????????????????????10???9??????10????????????
	ag.yearAt(2004).
		nyd().rest(1, 1).
		sf().rest(1, 22).to(1, 28).work(1, 17).to(1, 18).inLieu(1, 27).to(1, 28).
		ld().rest(5, 1).to(5, 7).work(5, 8).to(5, 9).inLieu(5, 6).to(5, 7).
		nd().rest(10, 1).to(10, 7).work(10, 9).to(10, 10).inLieu(10, 6).to(10, 7)
}

func (ag *arrangement) yearAt(year int) *arrangement {
	ag.year = year
	return ag
}

// ?????? New Year's Day
func (ag *arrangement) nyd() *arrangement {
	return ag.mark(chinesecalendar.NewYearsDay)
}

// ?????? Sprint Festival
func (ag *arrangement) sf() *arrangement {
	return ag.mark(chinesecalendar.SpringFestival)
}

// ????????? Tomb-Sweeping Day
func (ag *arrangement) tsd() *arrangement {
	return ag.mark(chinesecalendar.TombSweepingDay)
}

// ????????? Labour Day
func (ag *arrangement) ld() *arrangement {
	return ag.mark(chinesecalendar.LabourDay)
}

// ????????? Dragon Boat Festival
func (ag *arrangement) dbf() *arrangement {
	return ag.mark(chinesecalendar.DragonBoatFestival)
}

// ????????? National Day
func (ag *arrangement) nd() *arrangement {
	return ag.mark(chinesecalendar.NationalDay)
}

// ????????? Mid-autumn Festival
func (ag *arrangement) maf() *arrangement {
	return ag.mark(chinesecalendar.MidAutumnFestival)
}

// ?????????????????????????????????????????????????????????70??????????????? Anti-Fascist 70th Day
func (ag *arrangement) afd() *arrangement {
	return ag.mark(chinesecalendar.NationalDay)
}

func (ag *arrangement) mark(holiday chinesecalendar.Holiday) *arrangement {
	ag.holiday = holiday
	return ag
}

func (ag *arrangement) work(month, day int) *arrangement {
	return ag.save(month, day, dateTypeWorkday)
}

func (ag *arrangement) rest(month, day int) *arrangement {
	return ag.save(month, day, dateTypeHoliday)
}

func (ag *arrangement) inLieu(month, day int) *arrangement {
	return ag.save(month, day, dateTypeInLieu)
}

func (ag *arrangement) save(month, day, dayType int) *arrangement {
	if ag.year == 0 {
		panic("should set year before saving holiday")
	}
	if ag.holiday.Name() == "" {
		panic("should set holiday before saving holiday")
	}
	ag.month = month
	ag.day = day
	ag.dayType = dayType
	t := time.Date(ag.year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	switch dayType {
	case dateTypeHoliday:
		ag.Holidays[t] = ag.holiday
	case dateTypeWorkday:
		ag.Workdays[t] = ag.holiday
	case dateTypeInLieu:
		ag.InLieuDays[t] = ag.holiday
	}
	return ag
}

func (ag *arrangement) to(month, day int) *arrangement {
	if ag.year == 0 || ag.month == 0 || ag.day == 0 {
		panic("should set year, month, day before saving holiday")
	}
	startDate := Date(ag.year, ag.month, ag.day)
	endDate := Date(ag.year, month, day)
	if endDate.Before(startDate) {
		panic("end date should be after start date")
	}
	days := int(endDate.Sub(startDate)/(24*time.Hour)) + 1
	for i := 0; i < days; i++ {
		t := startDate.Add(24 * time.Hour * time.Duration(i))
		switch ag.dayType {
		case dateTypeHoliday:
			ag.Holidays[t] = ag.holiday
		case dateTypeWorkday:
			ag.Workdays[t] = ag.holiday
		case dateTypeInLieu:
			ag.InLieuDays[t] = ag.holiday
		}
	}
	return ag
}

var arrangementTemplate = `// Code generated by "scripts/generator"; DO NOT EDIT.
// Code generated by "scripts/generator"; DO NOT EDIT.
// Code generated by "scripts/generator"; DO NOT EDIT.
package chinesecalendar

import (
	"time"
	. "github.com/wangzeping722/chinesecalendar/internal"
)

var (
	// ???????????????
	minDay = Date({{.MinDay.Year}}, {{.MinDay.Month | printf "%d"}}, {{.MinDay.Day}})
	maxDay = Date({{.MaxDay.Year}}, {{.MaxDay.Month | printf "%d"}}, {{.MaxDay.Day}})
	// ?????????
	holidays = map[time.Time]Holiday{
		{{range $key := .HolidayList}}{{with index $.Holidays $key}}Date({{$key.Year}}, {{$key.Month | printf "%d"}}, {{$key.Day}}):{{index $.HolidayFieldMap .}},
		{{end}}{{end}}
	}

	// ?????????
	workdays = map[time.Time]Holiday{
		{{range $key := .WorkdayList}}{{with index $.Workdays $key}}Date({{$key.Year}}, {{$key.Month | printf "%d"}}, {{$key.Day}}):{{index $.HolidayFieldMap .}},
		{{end}}{{end}}
	}

	// ?????????
	inLieuDays = map[time.Time]Holiday{
		{{range $key := .InLieuDayList}}{{with index $.InLieuDays $key}}Date({{$key.Year}}, {{$key.Month | printf "%d"}}, {{$key.Day}}):{{index $.HolidayFieldMap .}},
		{{end}}{{end}}
	}
)
`

func generate() string {
	ag := newArragement()
	ag.generateHolidays()
	t := template.Must(template.New("").Parse(arrangementTemplate))

	buffer := &bytes.Buffer{}
	t.Execute(buffer, ag)
	return buffer.String()
}

func main() {
	str := generate()
	file, err := os.Create("constants.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(str))
	if err != nil {
		panic(err)
	}

	exec.Command("gofmt", "-w constants.go")
}
