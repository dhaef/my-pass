package model

import "time"

func getNowTimeStamp() string {
	return time.Now().Format(time.RFC3339)
}
