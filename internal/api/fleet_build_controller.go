package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type FleetBuildController struct {
	fleetBuildRepository *dao.FleetBuildRepository
}

func NewFleetBuildController(repository *dao.FleetBuildRepository) *FleetBuildController {
	return &FleetBuildController{fleetBuildRepository: repository}
}

func (controller *FleetBuildController) GetFleetBuild(c *gin.Context) {
	id := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(id)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}
	c.JSON(http.StatusOK, fleetBuild)
}

func (controller *FleetBuildController) GetAllFleetBuilds(c *gin.Context) {
	c.JSON(http.StatusOK, controller.fleetBuildRepository.GetAll())
}

func (controller *FleetBuildController) CreateFleetBuild(c *gin.Context) {
	var fleetBuild galaxy.FleetBuild
	if err := c.ShouldBindJSON(&fleetBuild); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.fleetBuildRepository.Upsert(&fleetBuild)
	c.JSON(http.StatusCreated, fleetBuild)
}

func (controller *FleetBuildController) UpdateFleetBuild(c *gin.Context) {
	id := c.Param("id")
	existing := controller.fleetBuildRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	var fleetBuild galaxy.FleetBuild
	if err := c.ShouldBindJSON(&fleetBuild); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fleetBuild.ID = id
	controller.fleetBuildRepository.Upsert(&fleetBuild)
	c.JSON(http.StatusOK, fleetBuild)
}

func (controller *FleetBuildController) DeleteFleetBuild(c *gin.Context) {
	id := c.Param("id")
	existing := controller.fleetBuildRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	controller.fleetBuildRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "FleetBuild deleted successfully"})
}

func (controller *FleetBuildController) GetAssignedShipModels(c *gin.Context) {
	id := c.Param("id")
	existing := controller.fleetBuildRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	assignedModels := controller.fleetBuildRepository.FindAssignedShipModels(id)
	c.JSON(http.StatusOK, assignedModels)
}

func (controller *FleetBuildController) AssignShipModel(c *gin.Context) {
	fleetBuildId := c.Param("id")
	existing := controller.fleetBuildRepository.Get(fleetBuildId)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	var assignment galaxy.FleetBuildToShipModel
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment.FleetBuildID = fleetBuildId

	wasCreated := controller.fleetBuildRepository.AssignShipModel(&assignment)

	if wasCreated {
		c.JSON(http.StatusCreated, assignment)
	} else {
		c.JSON(http.StatusOK, assignment)
	}
}

func (controller *FleetBuildController) UnassignShipModel(c *gin.Context) {
	fleetBuildId := c.Param("id")
	shipModelId := c.Param("shipModelId")

	existing := controller.fleetBuildRepository.Get(fleetBuildId)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	success := controller.fleetBuildRepository.UnassignShipModel(fleetBuildId, shipModelId)
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel assignment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ShipModel unassigned successfully"})
}
