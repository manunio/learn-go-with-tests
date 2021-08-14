package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

type Player struct {
	Name string
	Wins int
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	playerServer := new(PlayerServer)

	playerServer.store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(playerServer.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(playerServer.playerHandler))

	playerServer.Handler = router

	return playerServer
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	check(json.NewEncoder(w).Encode(p.store.GetLeague()))
}

func (p *PlayerServer) getLeagueTable() []Player {
	leagueTable := []Player{{"Chris", 20}}
	return leagueTable
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	_, err := fmt.Fprint(w, score)
	check(err)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
	return
}
