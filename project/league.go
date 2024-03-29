package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

// League ..
type League []Player

// NewLeague ..
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("Problem parsing league, %v", err)
	}
	return league, err
}

// Find ..
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
