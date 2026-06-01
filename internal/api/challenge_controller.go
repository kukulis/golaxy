package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

type ChallengeController struct {
	authenticationManager AuthenticationManager
	challengeRepository   *dao.ChallengeRepository
}

func NewChallengeController(authenticationManager AuthenticationManager, challengeRepository *dao.ChallengeRepository) *ChallengeController {
	return &ChallengeController{
		authenticationManager: authenticationManager,
		challengeRepository:   challengeRepository,
	}
}

// GetChallenge godoc
// @Summary Get a challenge by ID
// @Tags challenges
// @Produce json
// @Param id path string true "Challenge ID"
// @Success 200 {object} galaxy.Challenge
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /challenges/{id} [get]
func (controller *ChallengeController) GetChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	challenge := controller.challengeRepository.Get(id)
	if challenge == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}
	c.JSON(http.StatusOK, challenge)
}

// GetAllChallenges godoc
// @Summary List all challenges
// @Tags challenges
// @Produce json
// @Success 200 {array} galaxy.Challenge
// @Failure 401 {object} map[string]string
// @Router /challenges [get]
func (controller *ChallengeController) GetAllChallenges(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, controller.challengeRepository.GetAll())
}

// CreateChallenge godoc
// @Summary Create a challenge
// @Tags challenges
// @Accept json
// @Produce json
// @Param challenge body galaxy.Challenge true "Challenge data"
// @Success 201 {object} galaxy.Challenge
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /challenges [post]
func (controller *ChallengeController) CreateChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var challenge galaxy.Challenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if challenge.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	challenge.ChallengerRaceId = race.ID
	controller.challengeRepository.Upsert(&challenge)
	c.JSON(http.StatusCreated, challenge)
}

// UpdateChallenge godoc
// @Summary Update a challenge
// @Tags challenges
// @Accept json
// @Produce json
// @Param id path string true "Challenge ID"
// @Param challenge body galaxy.Challenge true "Challenge data"
// @Success 200 {object} galaxy.Challenge
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /challenges/{id} [put]
func (controller *ChallengeController) UpdateChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	existing := controller.challengeRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}
	if existing.ChallengerRaceId != race.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var challenge galaxy.Challenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	challenge.ID = id
	challenge.ChallengerRaceId = race.ID
	controller.challengeRepository.Upsert(&challenge)
	c.JSON(http.StatusOK, challenge)
}

// DeleteChallenge godoc
// @Summary Delete a challenge
// @Tags challenges
// @Produce json
// @Param id path string true "Challenge ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /challenges/{id} [delete]
func (controller *ChallengeController) DeleteChallenge(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	existing := controller.challengeRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}
	if existing.ChallengerRaceId != race.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	controller.challengeRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "Challenge deleted successfully"})
}
