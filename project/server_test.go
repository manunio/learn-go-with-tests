package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	mu       sync.RWMutex
	scores   map[string]int
	winCalls []string
}

func NewStubPlayerStore(scores map[string]int, winCalls []string) *StubPlayerStore {
	return &StubPlayerStore{sync.RWMutex{}, scores, winCalls}
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
	store := NewStubPlayerStore(map[string]int{"Pepper": 20, "Floyd": 10}, nil)
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
	store := NewStubPlayerStore(map[string]int{}, nil)
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

func assertResponseBody(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestNewPlayerStoreConcurrently(t *testing.T) {
	t.Run("it runs safe concurrently", func(t *testing.T) {
		store := NewStubPlayerStore(map[string]int{"Pepper": 20, "Floyd": 10}, nil)
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
	store := NewStubPlayerStore(map[string]int{}, nil)
	server := NewPlayerServer(store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
	})
}
