package poker

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	dummyBlindAlerter = &SpyBlindAlerter{}
	dummyPlayerStore  = &StubPlayerStore{mu: &sync.RWMutex{}}
	// dummyStdIn        = &bytes.Buffer{}
	dummyStdOut = &bytes.Buffer{}
)

type GameSpy struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	FinishCalled     bool
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	fmt.Printf("Start called.. %d\n", numberOfPlayers)
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	_, err := out.Write(g.BlindAlert)
	check(err)
}

func (g *GameSpy) Finish(winner string) {
	fmt.Printf("Finished called.. %q\n", winner)
	g.FinishCalled = true
	g.FinishCalledWith = winner
}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Manu' as winner", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("3", "Manu wins")

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertGameFinishCalledWith(t, game, "Manu")
	})

	t.Run("start game with 8 players and record 'Maxx' as winner", func(t *testing.T) {
		game := &GameSpy{}
		in := userSends("8", "Maxx wins")

		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertGameFinishCalledWith(t, game, "Maxx")
	})

	t.Run("it prints error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("NaN")

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, PlayerPrompt, BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when the winner is declared incorrectly", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("8", "Lloyd is a killer")

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, PlayerPrompt, BadWinnerInputMessage)
	})

}

func assertGameNotFinished(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.FinishCalled {
		t.Errorf("game should not have finished")
	}
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartCalledWith == numberOfPlayersWanted
	})

	if !passed {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartCalledWith)
	}
}

func assertGameFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}

func userSends(message ...string) *strings.Reader {
	return strings.NewReader(strings.Join(message, "\n"))
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
