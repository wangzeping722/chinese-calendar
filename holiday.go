package chinesecalendar

import "time"

var (
	dateFormatYYYYMMDD = "2006-01-02"
	oneDay             = 24 * time.Hour

	// 节假日定义
	NewYearsDay        = Holiday{"New Year's Day", "元旦", 1}
	SpringFestival     = Holiday{"Spring Festival", "春节", 3}
	TombSweepingDay    = Holiday{"Tomb-sweeping Day", "清明", 1}
	LabourDay          = Holiday{"Labour Day", "劳动节", 1}
	DragonBoatFestival = Holiday{"Dragon Boat Festival", "端午", 1}
	NationalDay        = Holiday{"National Day", "国庆节", 3}
	MidAutumnFestival  = Holiday{"Mid-autumn Festival", "中秋", 1}
	AntiFascist70thDay = Holiday{"Anti-Fascist 70th Day", "中国人民抗日战争暨世界反法西斯战争胜利70周年纪念日", 1}
)

type Holiday struct {
	engName string
	name    string
	days    int
}

func (h *Holiday) Name() string {
	return h.name
}

func (h *Holiday) EngName() string {
	return h.engName
}
