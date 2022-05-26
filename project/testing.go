package poker

import (
	"fmt"
	"io"
	"sync"
	"testing"
	"time"
)

type StubPlayerStore struct {
	mu       *sync.RWMutex
	scores   map[string]int
	winCalls []string
	league   []Player
}

func NewStubPlayerStore(
	mu *sync.RWMutex,
	scores map[string]int,
	winCalls []string,
	league []Player) *StubPlayerStore {
	return &StubPlayerStore{
		mu,
		scores,
		winCalls,
		league,
	}
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.winCalls = append(s.winCalls, name)
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

type ScheduledAlert struct {
	at     time.Duration
	amount int
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount})
}
