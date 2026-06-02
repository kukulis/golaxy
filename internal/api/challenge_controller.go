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

// GetChallenges godoc
// @Summary List challenges
// @Tags challenges
// @Produce json
// @Param challenger_id query string false "Filter by challenger race ID"
// @Param challengee_id query string false "Filter by challengee race ID"
// @Param accepted_a query string false "Filter by accepted A (1=true)"
// @Param accepted_b query string false "Filter by accepted B (1=true)"
// @Success 200 {array} galaxy.Challenge
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /challenges [get]
func (controller *ChallengeController) GetChallenges(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filter := galaxy.ChallengesFilter{
		ChallengerId: c.Query("challenger_id"),
		ChallengeeId: c.Query("challengee_id"),
		AcceptedA:    c.Query("accepted_a") == "1",
		AcceptedB:    c.Query("accepted_b") == "1",
	}

	if race.Role != galaxy.RoleAdmin {
		if filter.ChallengerId == "" && filter.ChallengeeId == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if filter.ChallengerId != "" && filter.ChallengerId != race.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if filter.ChallengeeId != "" && filter.ChallengeeId != race.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, controller.challengeRepository.GetAll(filter))
}

// GetChallenge godoc
// @Summary Get a challenge by ID
// @Tags challenges
// @Produce json
// @Param id path string true "Challenge ID"
// @Success 200 {object} galaxy.Challenge
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /challenges/{id} [get]
func (controller *ChallengeController) GetChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	challenge := controller.challengeRepository.Get(c.Param("id"))
	if challenge == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	if race.Role != galaxy.RoleAdmin && race.ID != challenge.ChallengerRaceId && race.ID != challenge.ChallengeeRaceId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func (controller *ChallengeController) UpdateChallenge(c *gin.Context) {
	panic("unimplemented")
}

func (controller *ChallengeController) DeleteChallenge(c *gin.Context) {
	panic("unimplemented")
}
