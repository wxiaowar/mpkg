//	精确到0.1秒

package time

import (
	itime "time"
)

const (
	TIME_LAYOUT_1 = "2006-01-02 15:04:05"
	TIME_LAYOUT_2 = "2006-01-02"
)

var Now itime.Time
var Local *itime.Location

func init() {
	Now = itime.Now().Round(itime.Second)
	Local, _ = itime.LoadLocation("Local")
	go refresh()
}

func refresh() {
	for {
		Now = itime.Now().Round(itime.Second)
		itime.Sleep(100 * itime.Millisecond)
	}
}

// 获取时间界限，如：today  返回stm: 2015-05-01 00:00:00  etm: 2015-05-02: 00:00:00
func TmLime(tmflag string) (stm, etm string) {
	stm = "1970-01-01 00:00:00"
	etm = "2070-01-01 00:00:00"
	if "today" == tmflag {
		stm = Now.Format("2006-01-02") + " 00:00:00"
		etm_tm := Now.AddDate(0, 0, 1)
		etm = etm_tm.Format("2006-01-02") + " 00:00:00"
	} else if "yesterday" == tmflag {
		stm_tm := Now.AddDate(0, 0, -1)
		stm = stm_tm.Format("2006-01-02") + " 00:00:00"
		etm = Now.Format("2006-01-02") + " 00:00:00"
	}
	return
}


