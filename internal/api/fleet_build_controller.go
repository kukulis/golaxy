package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type FleetBuildController struct {
	fleetBuildRepository *dao.FleetBuildRepository
	shipModelRepository  *dao.ShipModelRepository
}

func NewFleetBuildController(
	fleetBuildRepository *dao.FleetBuildRepository,
	shipModelRepository *dao.ShipModelRepository,
) *FleetBuildController {
	return &FleetBuildController{
		fleetBuildRepository: fleetBuildRepository,
		shipModelRepository:  shipModelRepository,
	}
}

// GetFleetBuild godoc
// @Summary Get a fleet build by ID
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {object} galaxy.FleetBuild
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id} [get]
func (controller *FleetBuildController) GetFleetBuild(c *gin.Context) {
	id := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(id)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}
	c.JSON(http.StatusOK, fleetBuild)
}

// GetAllFleetBuilds godoc
// @Summary List all fleet builds
// @Tags fleet-builds
// @Produce json
// @Success 200 {array} galaxy.FleetBuild
// @Router /fleet-builds [get]
func (controller *FleetBuildController) GetAllFleetBuilds(c *gin.Context) {
	c.JSON(http.StatusOK, controller.fleetBuildRepository.GetAll())
}

// CreateFleetBuild godoc
// @Summary Create a fleet build
// @Tags fleet-builds
// @Accept json
// @Produce json
// @Param fleetBuild body galaxy.FleetBuild true "FleetBuild data"
// @Success 201 {object} galaxy.FleetBuild
// @Failure 400 {object} map[string]string
// @Router /fleet-builds [post]
func (controller *FleetBuildController) CreateFleetBuild(c *gin.Context) {
	var fleetBuild galaxy.FleetBuild
	if err := c.ShouldBindJSON(&fleetBuild); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.fleetBuildRepository.Upsert(&fleetBuild)
	c.JSON(http.StatusCreated, fleetBuild)
}

// UpdateFleetBuild godoc
// @Summary Update a fleet build
// @Tags fleet-builds
// @Accept json
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Param fleetBuild body galaxy.FleetBuild true "FleetBuild data"
// @Success 200 {object} galaxy.FleetBuild
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id} [put]
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

// DeleteFleetBuild godoc
// @Summary Delete a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id} [delete]
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

// GetAssignedShipModels godoc
// @Summary List ship models assigned to a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {array} galaxy.FleetBuildToShipModel
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/ship-models [get]
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

// AssignShipModel godoc
// @Summary Assign a ship model to a fleet build
// @Tags fleet-builds
// @Accept json
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Param assignment body galaxy.FleetBuildToShipModel true "Assignment data"
// @Success 201 {object} galaxy.FleetBuildToShipModel
// @Success 200 {object} galaxy.FleetBuildToShipModel
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/ship-models [post]
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

// UnassignShipModel godoc
// @Summary Unassign a ship model from a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Param shipModelId path string true "ShipModel ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/ship-models/{shipModelId} [delete]
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

// TODO
func (controller *FleetBuildController) CalculateShipTech(c *gin.Context) {
	fleetBuildId := c.Param("id")
	shipModelId := c.Param("shipModelId")

	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	//build2shipModel := controller.fleetBuildRepository.FindAssignedShipModel(fleetBuildId, shipModelId)
	//build2shipModel.ShipModel = controller.shipModelRepository.Get(build2shipModel.ShipModelID)

	shipModel := controller.shipModelRepository.Get(shipModelId)

	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}

	shipTech := fleetBuild.CalculateShipTech(shipModel)

	c.JSON(http.StatusOK, shipTech)
}
