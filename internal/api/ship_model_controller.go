package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type ShipModelController struct {
	authenticationManager AuthenticationManager
	shipModelRepository   *dao.ShipModelRepository
}

func NewShipModelController(authenticationManager AuthenticationManager, repository *dao.ShipModelRepository) *ShipModelController {
	return &ShipModelController{authenticationManager: authenticationManager, shipModelRepository: repository}
}

// GetShipModel godoc
// @Summary Get a ship model by ID
// @Tags ship-models
// @Produce json
// @Param id path string true "ShipModel ID"
// @Success 200 {object} galaxy.ShipModel
// @Failure 404 {object} map[string]string
// @Router /ship-models/{id} [get]
func (controller *ShipModelController) GetShipModel(c *gin.Context) {
	id := c.Param("id")
	shipModel := controller.shipModelRepository.Get(id)
	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}
	c.JSON(http.StatusOK, shipModel)
}

// GetAllShipModels godoc
// @Summary List all ship models
// @Tags ship-models
// @Produce json
// @Success 200 {array} galaxy.ShipModel
// @Router /ship-models [get]
func (controller *ShipModelController) GetAllShipModels(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	ownerId := ""
	if race != nil {
		ownerId = race.ID
	}
	c.JSON(http.StatusOK, controller.shipModelRepository.GetAll(ownerId))
}

// CreateShipModel godoc
// @Summary Create a ship model
// @Tags ship-models
// @Accept json
// @Produce json
// @Param shipModel body galaxy.ShipModel true "ShipModel data"
// @Success 201 {object} galaxy.ShipModel
// @Failure 400 {object} map[string]string
// @Router /ship-models [post]
func (controller *ShipModelController) CreateShipModel(c *gin.Context) {
	var shipModel galaxy.ShipModel
	if err := c.ShouldBindJSON(&shipModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.shipModelRepository.Upsert(&shipModel)
	c.JSON(http.StatusCreated, shipModel)
}

// UpdateShipModel godoc
// @Summary Update a ship model
// @Tags ship-models
// @Accept json
// @Produce json
// @Param id path string true "ShipModel ID"
// @Param shipModel body galaxy.ShipModel true "ShipModel data"
// @Success 200 {object} galaxy.ShipModel
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-models/{id} [put]
func (controller *ShipModelController) UpdateShipModel(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	id := c.Param("id")
	existing := controller.shipModelRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}

	var shipModel galaxy.ShipModel
	if err := c.ShouldBindJSON(&shipModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipModel.ID = id
	shipModel.OwnerId = race.ID
	controller.shipModelRepository.Upsert(&shipModel)
	c.JSON(http.StatusOK, shipModel)
}

// CalculateShipTech godoc
// @Summary Calculate ship tech for a ship model
// @Tags ship-models
// @Accept json
// @Produce json
// @Param id path string true "ShipModel ID"
// @Param technologies body galaxy.Technologies true "Technologies data"
// @Success 200 {object} galaxy.ShipTech
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-models/{id}/calculate-ship-tech [post]
func (controller *ShipModelController) CalculateShipTech(c *gin.Context) {
	id := c.Param("id")
	shipModel := controller.shipModelRepository.Get(id)
	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}
	var tech galaxy.Technologies
	if err := c.ShouldBindJSON(&tech); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shipTech := shipModel.CalculateShipTech(&tech)
	c.JSON(http.StatusOK, shipTech)
}

// DeleteShipModel godoc
// @Summary Delete a ship model
// @Tags ship-models
// @Produce json
// @Param id path string true "ShipModel ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-models/{id} [delete]
func (controller *ShipModelController) DeleteShipModel(c *gin.Context) {
	id := c.Param("id")
	existing := controller.shipModelRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}

	controller.shipModelRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "ShipModel deleted successfully"})
}
