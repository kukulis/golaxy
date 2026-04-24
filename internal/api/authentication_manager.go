package api

import "glaktika.eu/galaktika/pkg/galaxy"

type AuthenticationManager interface {
	Authenticate(token string) *galaxy.Race
	TokenValid(token string) bool
	AddToken(token string, race *galaxy.Race)
}

type MemoryAuthenticationManager struct {
	tokenToRace map[string]*galaxy.Race
}

func NewMemoryAuthenticationManager() *MemoryAuthenticationManager {
	return &MemoryAuthenticationManager{
		tokenToRace: make(map[string]*galaxy.Race),
	}
}

func (am *MemoryAuthenticationManager) Authenticate(token string) *galaxy.Race {
	return am.tokenToRace[token]
}

func (am *MemoryAuthenticationManager) TokenValid(token string) bool {
	_, yes := am.tokenToRace[token]
	return yes
}

func (am *MemoryAuthenticationManager) AddToken(token string, race *galaxy.Race) {
	am.tokenToRace[token] = race
}
