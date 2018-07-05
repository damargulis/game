package player

import (
	"time"
)

func sleep(t float64, ch chan int) {
	time.Sleep(time.Duration(t) * time.Millisecond)
	ch <- 0
}
