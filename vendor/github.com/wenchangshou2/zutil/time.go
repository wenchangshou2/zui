package zutil

import "time"

func Now()int64{
	return time.Now().UnixNano() / 1000000
}
