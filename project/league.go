package main

import (
	"encoding/json"
	"errors"
	"io"
)

// NewLeague ..
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = errors.New("Problem parsing error")
	}
	return league, err
}
