package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
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

// CreateRace godoc
// @Summary Create a race
// @Tags races
// @Accept json
// @Produce json
// @Param race body galaxy.Race true "Race data"
// @Success 201 {object} galaxy.Race
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /races [post]
func (controller *RaceController) CreateRace(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if race.Role != galaxy.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var newRace galaxy.Race
	if err := c.ShouldBindJSON(&newRace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.raceRepository.Upsert(&newRace)
	c.JSON(http.StatusCreated, newRace)
}

// DeleteRace godoc
// @Summary Delete a race
// @Tags races
// @Produce json
// @Param id path string true "Race ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /races/{id} [delete]
func (controller *RaceController) DeleteRace(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if race.Role != galaxy.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	id := c.Param("id")
	if race.ID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete own race"})
		return
	}
	if controller.raceRepository.Get(id) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Race not found"})
		return
	}
	controller.raceRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "Race deleted successfully"})
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
