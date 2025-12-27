package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type ShipModelController struct {
	shipModelRepository *dao.ShipModelRepository
}

func NewShipModelController(repository *dao.ShipModelRepository) *ShipModelController {
	return &ShipModelController{shipModelRepository: repository}
}

func (controller *ShipModelController) GetShipModel(c *gin.Context) {
	id := c.Param("id")
	shipModel := controller.shipModelRepository.Get(id)
	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}
	c.JSON(http.StatusOK, shipModel)
}

func (controller *ShipModelController) GetAllShipModels(c *gin.Context) {
	c.JSON(http.StatusOK, controller.shipModelRepository.GetAll())
}

func (controller *ShipModelController) CreateShipModel(c *gin.Context) {
	var shipModel galaxy.ShipModel
	if err := c.ShouldBindJSON(&shipModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.shipModelRepository.Upsert(&shipModel)
	c.JSON(http.StatusCreated, shipModel)
}

func (controller *ShipModelController) UpdateShipModel(c *gin.Context) {
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
	controller.shipModelRepository.Upsert(&shipModel)
	c.JSON(http.StatusOK, shipModel)
}

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
