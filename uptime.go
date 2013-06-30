package helpers

import (
	"time"
)

var start = time.Now().Round(time.Second)

func UptimeGet() time.Time {
	return start
}

func UptimeGetString() string {
	return time.Now().Round(time.Second).Sub(start).String()
}
