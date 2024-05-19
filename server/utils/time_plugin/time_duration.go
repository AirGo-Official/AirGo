package time_plugin

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}

// 获取月初第一天（00:00:00）到当天（23:59:59）的起止时间
func GetFirstToTodayForMonth() (start time.Time, end time.Time) {
	today := time.Now()
	return time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()), time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, -1, 0, today.Location())
}

// 距离第二天零点的时间间隔
func GetTimeIntervalBetweenNowAndMidnightTheNextDay() time.Duration {
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return zeroTime.Sub(now)
}
