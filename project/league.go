package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

// NewLeague ..
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("Problem parsing league, %v", err)
	}
	return league, err
}