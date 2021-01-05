/*
	时间处理相关
*/
package mutil

import (
	"time"
)

// TimeMinAgo 返回前一分钟整秒
func TimeMinAgo(ts int64) int64 {
	Ts := time.Unix(ts-60, ts)

	//小时Ts
	hour := Ts.Hour()

	// 分钟Ts
	minute := Ts.Minute()
	mu := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), hour, minute, 0, 0, Ts.Location())
	return mu.Unix()
}

// TimeMinZero 返回时间戳的分钟整秒
func TimeMinZero(ts int64) int64 {
	timer := time.Unix(ts, 0)
	z := time.Date(timer.Year(), timer.Month(), timer.Day(), timer.Hour(), timer.Minute(), 0, 0, timer.Location())
	return z.Unix()
}

// TimeZeroToday 返回今日零点时间戳
func TimeZeroToday() (zeroTs int64) {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.Unix()
}

// TimeZeroDayAttach 返回今天起的天数整点时间戳(附加条件)
func TimeZeroDayAttach(n int) (ts time.Time) {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day()+n, 0, 0, 0, 0, t.Location())
	return newTime
}

// GetHmsByTs 获取数据suo在的时 分 秒对应创建索引
func TimeGetHmsByTs(ts int64) (int64, int64, int64) {
	Ts := time.Unix(ts, ts)

	//小时Ts
	hour := Ts.Hour()
	hu := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), hour, 0, 0, 0, Ts.Location())

	// 分钟Ts
	minute := Ts.Minute()
	mu := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), hour, minute, 0, 0, Ts.Location())

	// 秒Ts
	second := Ts.Second()
	su := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), hour, minute, second, 0, Ts.Location())
	return hu.Unix(), mu.Unix(), su.Unix()
}

func TimeDayStartEnd(ts int64) (int64, int64) {
	Ts := time.Unix(ts, ts)
	dayStart := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), 0, 0, 0, 0, Ts.Location())
	dayEnd := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), 23, 59, 0, 0, Ts.Location())

	return dayStart.Unix(), dayEnd.Unix()
}

func TimeFormatUtc() string {
	time.Now().Format("2006-01-02T15:04:05")
	return time.Now().UTC().Format("2006-01-02T15:04:05")
}

// TimeZeroForToday 返回今日零点时间戳
func TimeZeroForToday() (zeroTs int64) {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.Unix()
}
