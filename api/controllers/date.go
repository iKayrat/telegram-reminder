package controllers

import (
	"fmt"
	"time"
)

func Check(t time.Time, ch chan bool) {
	loc, err := time.LoadLocation("UTC+06:00")
	if err != nil {
		fmt.Println(err)
	}

	h, err := time.ParseInLocation("15:04", "0:30", loc)
	if err != nil {
		fmt.Println(err)
	}

	for {
		if t.Equal(h) {
			ch <- true
		}
	}
}
