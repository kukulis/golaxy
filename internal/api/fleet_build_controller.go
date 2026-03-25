package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type FleetBuildController struct {
	authenticationManager AuthenticationManager
	fleetBuildRepository  *dao.FleetBuildRepository
	fleetRepository       *dao.FleetRepository
	shipModelRepository   *dao.ShipModelRepository
	divisionRepository    *dao.DivisionRepository
}

func NewFleetBuildController(
	authenticationManager AuthenticationManager,
	fleetBuildRepository *dao.FleetBuildRepository,
	fleetRepository *dao.FleetRepository,
	shipModelRepository *dao.ShipModelRepository,
	divisionRepository *dao.DivisionRepository,
) *FleetBuildController {
	return &FleetBuildController{
		authenticationManager: authenticationManager,
		fleetBuildRepository:  fleetBuildRepository,
		fleetRepository:       fleetRepository,
		shipModelRepository:   shipModelRepository,
		divisionRepository:    divisionRepository,
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
	if controller.divisionRepository.Get(fleetBuild.DivisionId) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Division not found"})
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

// GetStatistics godoc
// @Summary Get fleet build statistics
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {object} galaxy.FleetBuildStatistics
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/statistics [get]
func (controller *FleetBuildController) GetStatistics(c *gin.Context) {
	fleetBuildId := c.Param("id")

	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	division := controller.divisionRepository.Get(fleetBuild.DivisionId)
	if division == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}

	assignments := controller.fleetBuildRepository.FindAssignedShipModels(fleetBuildId)
	fleetBuild.AssignedShipModels = make([]galaxy.ShipModelAssignment, 0, len(assignments))
	for _, a := range assignments {
		shipModel := controller.shipModelRepository.Get(a.ShipModelID)
		if shipModel != nil {
			fleetBuild.AssignedShipModels = append(fleetBuild.AssignedShipModels, galaxy.ShipModelAssignment{
				ShipModel: *shipModel,
				Amount:    a.Amount,
			})
		}
	}

	statistics := fleetBuild.CalculateStatistics(division.ResourcesAmount)
	c.JSON(http.StatusOK, statistics)
}

// Build godoc
// @Summary Build a fleet from a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/build [post]
func (controller *FleetBuildController) Build(c *gin.Context) {
	token := c.GetHeader("token")
	if !controller.authenticationManager.TokenValid(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userId := controller.authenticationManager.Authenticate(token)

	fleetBuildId := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	assignments := controller.fleetBuildRepository.FindAssignedShipModels(fleetBuildId)
	var ships []*galaxy.Ship
	for _, a := range assignments {
		shipModel := controller.shipModelRepository.Get(a.ShipModelID)
		if shipModel == nil {
			continue
		}
		shipTech := fleetBuild.CalculateShipTech(shipModel)
		for i := 0; i < a.Amount; i++ {
			ship := &galaxy.Ship{
				ID:    uuid.New().String(),
				Name:  shipModel.Name,
				Tech:  shipTech,
				Owner: userId,
			}
			ships = append(ships, ship)
		}
	}

	fleet := galaxy.NewFleet(ships)
	fleet.ID = uuid.New().String()
	fleet.Owner = userId
	fmt.Printf("Built fleet %s for user %s with %d ships\n", fleet.ID, userId, len(ships))

	controller.fleetRepository.Upsert(fleet)
	controller.fleetRepository.UpsertDivisionFleet(&galaxy.DivisionFleet{
		DivisionId: fleetBuild.DivisionId,
		UserId:     userId,
		FleetId:    fleet.ID,
	})

	c.JSON(http.StatusOK, fleet)
}

// GetFleet godoc
// @Summary Get fleet for a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Param token header string true "Authentication token"
// @Success 200 {object} galaxy.Fleet
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/fleet [get]
func (controller *FleetBuildController) GetFleet(c *gin.Context) {
	token := c.GetHeader("token")
	if !controller.authenticationManager.TokenValid(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userId := controller.authenticationManager.Authenticate(token)

	fleetBuildId := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	divisionFleet := controller.fleetRepository.GetDivisionFleet(fleetBuild.DivisionId, userId)
	if divisionFleet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fleet not found"})
		return
	}

	fleet := controller.fleetRepository.Get(divisionFleet.FleetId)
	if fleet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fleet not found"})
		return
	}

	c.JSON(http.StatusOK, fleet)
}

// CalculateShipTech godoc
// @Summary Calculate ship tech for a ship model assigned to a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Param shipModelId path string true "ShipModel ID"
// @Success 200 {object} galaxy.ShipTech
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/ship-models/{shipModelId}/calculate-ship-tech [get]
func (controller *FleetBuildController) CalculateShipTech(c *gin.Context) {
	fleetBuildId := c.Param("id")
	shipModelId := c.Param("shipModelId")

	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	assignment := controller.fleetBuildRepository.FindAssignedShipModel(fleetBuildId, shipModelId)
	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel is not assigned to this FleetBuild"})
		return
	}

	shipModel := controller.shipModelRepository.Get(shipModelId)
	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel not found"})
		return
	}

	shipTech := fleetBuild.CalculateShipTech(shipModel)
	c.JSON(http.StatusOK, shipTech)
}
