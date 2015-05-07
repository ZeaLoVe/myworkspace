package util

import (
	"time"
)

type Clock struct {
	start time.Time
}

func (t *Clock) Start() {
	t.start = time.Now()
}

func (t *Clock) Seconds() float64 {
	d := time.Since(t.start)
	return d.Seconds()
}
