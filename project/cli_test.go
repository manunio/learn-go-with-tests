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
	// 	dummyStdIn        = &bytes.Buffer{}
	// 	dummyStdOut       = &bytes.Buffer{}
)

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

type ScheduledAlert struct {
	at     time.Duration
	amount int
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount})
}

func TestCLI(t *testing.T) {

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Star called with 7 but got %d", game.StartedWith)
		}
	})


}

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players ", func(t *testing.T) {

		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)
		game.Start(5)

		cases := []ScheduledAlert{
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
		checkSchedulingCase(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCase(t, cases, blindAlerter)

	})
}

func TestGame_Finish(t *testing.T) {
	store := dummyPlayerStore
	game := NewGame(dummyBlindAlerter, store)
	winner := "Manu"

	game.Finish(winner)
	AssertPlayerWin(t, store, winner)
}

func checkSchedulingCase(t *testing.T, cases []ScheduledAlert, blindAlerter *SpyBlindAlerter) {
	t.Helper()
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}
			got := blindAlerter.alerts[i]
			assertScheduleAlert(t, got, want)
		})
	}
}

func assertScheduleAlert(t *testing.T, got, want ScheduledAlert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}

	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}
