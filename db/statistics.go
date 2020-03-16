package db

import "time"

//统计
const (
	_static_key_pre      = "pxs_"
	STATISTICS_TYPE_TIME = 1 //时长统计
)

type StatisticsItem struct {
	Key       string
	Type      uint32
	EventName string
}

func GetStatisticsKey(s string) string {
	return _static_key_pre + s
}

//判断时间是否可加
func StatisticsAddTime(lastTime time.Time) (time.Time, time.Duration) {
	//5分钟内都可以加
	now := time.Now()
	if lastTime.Add(5 * time.Minute).After(now) {
		return now, now.Sub(lastTime)
	}
	return lastTime, 0
}
