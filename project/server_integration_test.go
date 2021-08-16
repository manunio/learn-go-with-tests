package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "[]")
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(t, response.Code, http.StatusOK)
		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})

	// test for case when value is recorded by different clients(ex: cli)
	// but doesn not reflect it(ex: api) unless server is restarted.
	// TODO: check authenticity
	// t.Run("get score after post", func(t *testing.T) {
	// 	database, cleanDatabase := createTempFile(t, `[
	//         {"Name": "Cleo", "Wins": 10},
	//         {"Name": "Chris", "Wins": 33}]`)
	// 	defer cleanDatabase()
	//
	// 	var dummyBlindAlerter = &SpyBlindAlerter{}
	//
	// 	name := "Maxx"
	// 	store, err := NewFileSystemPlayerStore(database)
	// 	assertNoError(t, err)
	//
	// 	// cli client record requests.
	// 	in1 := strings.NewReader("Maxx wins\n")
	// 	cli := NewCLI(store, in1, dummyBlindAlerter)
	// 	cli.PlayPoker()
	//
	// 	in2 := strings.NewReader("Maxx wins\n")
	// 	cli = NewCLI(store, in2, dummyBlindAlerter)
	// 	cli.PlayPoker()
	//
	// 	// api get request
	// 	server := NewPlayerServer(store)
	//
	// 	request := newGetScoreRequest(name)
	// 	response := httptest.NewRecorder()
	//
	// 	server.ServeHTTP(response, request)
	//
	// 	assertStatus(t, response.Code, http.StatusOK)
	//
	// 	got, err := strconv.Atoi(response.Body.String())
	// 	assertNoError(t, err)
	// 	want := 2
	//
	// 	if got != want {
	// 		t.Errorf("got %d want %d", got, want)
	// 	}
	//
	// })

}
