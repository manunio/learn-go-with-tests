package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlinderAlertFunc func(duration time.Duration, amount int)

func (a BlinderAlertFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		_, err := fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
		check(err)
	})
}
