package mocking

import (
	"fmt"
	"io"
	"time"
)

const finalWorld = "Go!"
const countdownStart = 3
const sleep = "sleep"
const write = "write"

type Sleeper interface {
	Sleep()
}

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

type CountdownOperationsSpy struct {
	Calls []string
}

type ConfigurableSleep struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		_, _ = fmt.Fprintln(out, i)
	}

	sleeper.Sleep()
	_, _ = fmt.Fprint(out, finalWorld)
}
