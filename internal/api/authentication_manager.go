package api

type AuthenticationManager interface {
	Authenticate(token string) string
	TokenValid(token string) bool
	AddToken(token string, userId string)
}
type MemoryAuthenticationManager struct {
	tokenToUserId map[string]string
}

func NewMemoryAuthenticationManager() *MemoryAuthenticationManager {
	return &MemoryAuthenticationManager{
		tokenToUserId: make(map[string]string),
	}
}

func (am *MemoryAuthenticationManager) Authenticate(token string) string {
	return am.tokenToUserId[token]
}

func (am *MemoryAuthenticationManager) TokenValid(token string) bool {
	_, yes := am.tokenToUserId[token]
	return yes
}

func (am *MemoryAuthenticationManager) AddToken(token string, userId string) {
	am.tokenToUserId[token] = userId
}
