package util

import "time"

type LangString struct {
	De string `json:"de"`
	En string `json:"en"`
}

func SameDay(time1, time2 time.Time) bool {
	y1, m1, d1 := time1.Date()
	y2, m2, d2 := time2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
