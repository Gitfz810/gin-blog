package util

import "time"

func GetTime() time.Time {
	//return time.Now().UnixNano() / int64(time.Millisecond)
	return time.Now()
}
