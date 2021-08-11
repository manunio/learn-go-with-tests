package main

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func NewStubPlayerStore(scores map[string]int, winCalls []string) *StubPlayerStore {
	return &StubPlayerStore{scores, winCalls}
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
