package poker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const tenMS = 10 * time.Millisecond

var (
	dummyGame = &GameSpy{}
)

func TestGETPlayers(t *testing.T) {
	store := NewStubPlayerStore(&sync.RWMutex{}, map[string]int{"Pepper": 20, "Floyd": 10}, nil, nil)
	server, _ := NewPlayerServer(store, dummyGame)

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
	server, _ := NewPlayerServer(store, dummyGame)

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
		server, _ := NewPlayerServer(store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, _ := NewPlayerServer(NewStubPlayerStore(&sync.RWMutex{}, nil, nil, nil), dummyGame)

		request, _ := NewGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	// test for websocket
	t.Run("start a game with 3 players and declare Manu the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Manu"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		// time.Sleep(tenMS)
		assertGameStartedWith(t, game, 3)
		assertGameFinishCalledWith(t, game, winner)
		within(t, tenMS, func() {
			assertWebsocketGotMsg(t, ws, wantedBlindAlert)
		})
	})

}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}

func NewGameRequest() (*http.Request, error) {
	return http.NewRequest(http.MethodGet, "/game", nil)
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game *GameSpy) *PlayerServer {
	server, err := NewPlayerServer(store, dummyGame)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
