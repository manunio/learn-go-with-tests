package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type playerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *playerServer {
	playerServer := new(playerServer)

	playerServer.store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(playerServer.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(playerServer.playerHandler))

	playerServer.Handler = router

	return playerServer
}

func (p *playerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *playerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *playerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	_, _ = fmt.Fprint(w, score)
}

func (p *playerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
	return
}
