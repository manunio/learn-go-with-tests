package poker

import (
	"fmt"
	"io"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlinderAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlinderAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

func Alerter(duration time.Duration, amount int, io io.Writer) {
	time.AfterFunc(duration, func() {
		_, err := fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
		check(err)
	})
}
