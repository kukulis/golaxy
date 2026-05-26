package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

type AuthenticationManager interface {
	Authenticate(token string) *galaxy.Race
	AuthenticateFromContext(c *gin.Context) *galaxy.Race
	TokenValid(token string) bool
}

type MemoryAuthenticationManager struct {
	raceRepository *dao.RaceRepository
}

func NewMemoryAuthenticationManager(raceRepository *dao.RaceRepository) *MemoryAuthenticationManager {
	return &MemoryAuthenticationManager{
		raceRepository: raceRepository,
	}
}

func (am *MemoryAuthenticationManager) Authenticate(token string) *galaxy.Race {
	return am.raceRepository.GetByToken(token)
}

func (am *MemoryAuthenticationManager) AuthenticateFromContext(c *gin.Context) *galaxy.Race {
	return am.Authenticate(bearerToken(c))
}

func (am *MemoryAuthenticationManager) TokenValid(token string) bool {
	return am.raceRepository.GetByToken(token) != nil
}
