package ui

import (
	"fmt"
	"time"
)

func showData(t interface{}) string {
	return fmt.Sprintf("%+v", t)
}

func showTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02 15:04:05")
}
