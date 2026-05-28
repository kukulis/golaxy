package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"net/http"
)

type RaceController struct {
	authenticationManager AuthenticationManager
	raceRepository        *dao.RaceRepository
}

func NewRaceController(authenticationManager AuthenticationManager, repository *dao.RaceRepository) *RaceController {
	return &RaceController{authenticationManager: authenticationManager, raceRepository: repository}
}

// GetCurrentRace godoc
// @Summary Get the current logged-in race
// @Tags races
// @Produce json
// @Success 200 {object} galaxy.Race
// @Failure 401 {object} map[string]string
// @Router /current-race [get]
func (controller *RaceController) GetCurrentRace(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, race)
}

// GetAllRaces godoc
// @Summary List all races
// @Tags races
// @Produce json
// @Success 200 {array} galaxy.Race
// @Router /races [get]
func (controller *RaceController) GetAllRaces(c *gin.Context) {
	// TODO: filter out Token field from response before returning
	c.JSON(http.StatusOK, controller.raceRepository.GetAll())
}

// GetRace godoc
// @Summary Get a race by ID
// @Tags races
// @Produce json
// @Param id path string true "Race ID"
// @Success 200 {object} galaxy.Race
// @Failure 404 {object} map[string]string
// @Router /races/{id} [get]
func (controller *RaceController) GetRace(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	found := controller.raceRepository.Get(id)
	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Race not found"})
		return
	}
	c.JSON(http.StatusOK, found)
}
