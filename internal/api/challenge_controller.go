package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

type ChallengeController struct {
	authenticationManager AuthenticationManager
	challengeRepository   *dao.ChallengeRepository
	divisionRepository    *dao.DivisionRepository
	raceRepository        *dao.RaceRepository
}

func NewChallengeController(
	authenticationManager AuthenticationManager,
	challengeRepository *dao.ChallengeRepository,
	divisionRepository *dao.DivisionRepository,
	raceRepository *dao.RaceRepository,
) *ChallengeController {
	return &ChallengeController{
		authenticationManager: authenticationManager,
		challengeRepository:   challengeRepository,
		divisionRepository:    divisionRepository,
		raceRepository:        raceRepository,
	}
}

// CreateChallenge godoc
// @Summary Create a challenge
// @Tags challenges
// @Produce json
// @Param id query string true "Challenge UUID"
// @Param challenger_race_id query string true "Challenger race ID"
// @Param challengee_race_id query string true "Challengee race ID"
// @Param division_id query string true "Division ID"
// @Success 201 {object} galaxy.Challenge
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /challenges [post]
func (controller *ChallengeController) CreateChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Query("id")
	challengerRaceId := c.Query("challenger_race_id")
	challengeeRaceId := c.Query("challengee_race_id")
	divisionId := c.Query("division_id")

	if id == "" || challengerRaceId == "" || challengeeRaceId == "" || divisionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, challenger_race_id, challengee_race_id and division_id are required"})
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a valid UUID"})
		return
	}

	if race.Role != galaxy.RoleAdmin && race.ID != challengerRaceId {
		c.JSON(http.StatusForbidden, gin.H{"error": "challenger_race_id must match your race"})
		return
	}

	if race.Role == galaxy.RoleAdmin && controller.raceRepository.Get(challengerRaceId) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Challenger race not found"})
		return
	}

	if controller.raceRepository.Get(challengeeRaceId) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Challengee race not found"})
		return
	}

	if controller.divisionRepository.Get(divisionId) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Division not found"})
		return
	}

	if controller.challengeRepository.Get(id) != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Challenge with this id already exists"})
		return
	}

	challenge := galaxy.NewChallenge(id, challengerRaceId, challengeeRaceId, divisionId, time.Now())
	controller.challengeRepository.Upsert(challenge)
	c.JSON(http.StatusCreated, challenge)
}

func (controller *ChallengeController) GetChallenges(c *gin.Context) {
	panic("unimplemented")
}

func (controller *ChallengeController) GetChallenge(c *gin.Context) {
	panic("unimplemented")
}

func (controller *ChallengeController) UpdateChallenge(c *gin.Context) {
	panic("unimplemented")
}

func (controller *ChallengeController) DeleteChallenge(c *gin.Context) {
	panic("unimplemented")
}
