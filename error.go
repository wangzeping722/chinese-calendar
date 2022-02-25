package chinesecalendar

import "fmt"

var ErrUnSupportDate = fmt.Errorf("unsupported date, supported date range is %s - %s", minDay.Format(dateFormatYYYYMMDD), maxDay.Format(dateFormatYYYYMMDD))
