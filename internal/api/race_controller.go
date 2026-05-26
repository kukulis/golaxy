package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"net/http"
)

type RaceController struct {
	raceRepository *dao.RaceRepository
}

func NewRaceController(repository *dao.RaceRepository) *RaceController {
	return &RaceController{raceRepository: repository}
}

// GetAllRaces godoc
// @Summary List all races
// @Tags races
// @Produce json
// @Success 200 {array} galaxy.Race
// @Router /races [get]
func (controller *RaceController) GetAllRaces(c *gin.Context) {
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
	id := c.Param("id")
	race := controller.raceRepository.Get(id)
	if race == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Race not found"})
		return
	}
	c.JSON(http.StatusOK, race)
}
