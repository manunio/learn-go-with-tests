package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
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

func TestGETPlayers(t *testing.T) {
	store := NewStubPlayerStore(&sync.RWMutex{}, map[string]int{"Pepper": 20, "Floyd": 10}, nil, nil)
	server := NewPlayerServer(store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

}

func TestStoreWins(t *testing.T) {
	store := NewStubPlayerStore(&sync.RWMutex{}, map[string]int{}, nil, nil)
	server := NewPlayerServer(store)

	t.Run("it record wins on POST", func(t *testing.T) {
		const player = "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})

}

func TestNewPlayerStoreConcurrently(t *testing.T) {
	t.Run("it runs safe concurrently", func(t *testing.T) {
		store := NewStubPlayerStore(&sync.RWMutex{},
			map[string]int{"Pepper": 20, "Floyd": 10},
			nil, nil)
		wantedStore := 1000

		var wg sync.WaitGroup
		wg.Add(wantedStore)

		for i := 0; i < wantedStore; i++ {
			i := i
			go func() {
				store.RecordWin(fmt.Sprint(i))
				wg.Done()
			}()
		}
		wg.Wait()
		gotStore := len(store.winCalls)
		if gotStore != wantedStore {
			t.Errorf("got %d want %d ", gotStore, wantedStore)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table in JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := NewStubPlayerStore(&sync.RWMutex{}, nil, nil, wantedLeague)
		server := NewPlayerServer(store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}
