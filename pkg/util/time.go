package util

import (
	"time"
)

// GetNowUnixTimeInt64 return int64 unixTime
func GetNowUnixTimeInt64() int64 {
	now := time.Now()
	unix := now.Unix()
	return unix
}

// GetNow ไปใฎๆ้
func GetNow() *time.Time {
	t := time.Now().UTC()
	return &t
}
