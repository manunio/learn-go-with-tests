package poker

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	dummyBlindAlerter = &SpyBlindAlerter{}
	dummyPlayerStore  = &StubPlayerStore{mu: &sync.RWMutex{}}
	dummyStdIn        = &bytes.Buffer{}
	dummyStdOut       = &bytes.Buffer{}
)

type SpyBlindAlerter struct {
	alerts []ScheduleAlert
}

type ScheduleAlert struct {
	at     time.Duration
	amount int
}

func (s ScheduleAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, ScheduleAlert{at, amount})
}

func TestCLI(t *testing.T) {
	var spyDummyAlerter = &SpyBlindAlerter{}
	t.Run("record Manu wins from user input", func(t *testing.T) {

		in := strings.NewReader("Manu wins\n")
		playerStore := NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil)

		cli := NewCLI(playerStore, in, dummyStdOut, spyDummyAlerter)
		cli.PlayPoker()

		AssertPlayerWin(t, playerStore, "Manu")

	})

	t.Run("record Maxx wins from user input", func(t *testing.T) {
		in := strings.NewReader("Maxx wins\n")
		playerStore := NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil)

		cli := NewCLI(playerStore, in, dummyStdOut, spyDummyAlerter)
		cli.PlayPoker()

		AssertPlayerWin(t, playerStore, "Maxx")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {

		in := strings.NewReader("Maxx wins\n")
		playerStore := NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil)
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStore, in, dummyStdOut, blindAlerter)
		cli.PlayPoker()

		cases := []ScheduleAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}

	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(dummyPlayerStore, in, stdout, blindAlerter)
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

		cases := []ScheduleAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}

	})

}

func assertScheduleAlert(t *testing.T, got, want ScheduleAlert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}

	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}
