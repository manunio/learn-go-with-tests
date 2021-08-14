package poker_test

import (
	poker "learn-go-with-tests/project"
	"strings"
	"sync"
	"testing"
)

func TestCLI(t *testing.T) {

	t.Run("record Manu wins from user input", func(t *testing.T) {

		in := strings.NewReader("Manu wins\n")
		playerStore := poker.NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil)

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Manu")

	})

	t.Run("record Maxx wins from user input", func(t *testing.T) {
		in := strings.NewReader("Maxx wins\n")
		playerStore := poker.NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil)

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Maxx")
	})

}
