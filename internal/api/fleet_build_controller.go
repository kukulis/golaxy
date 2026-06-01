package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"log"
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	divisionId := c.Query("division_id")
	raceId := race.ID
	if c.Query("all") == "true" {
		if explicit := c.Query("race_id"); explicit != "" {
			raceId = explicit
		} else {
			raceId = ""
		}
	}
	c.JSON(http.StatusOK, controller.fleetBuildRepository.GetAll(divisionId, raceId))
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
// @Success 200 {array} galaxy.ShipModelAssignment
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/ship-model-assignments [get]
func (controller *FleetBuildController) GetAssignedShipModels(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	existing := controller.fleetBuildRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	assignedModels := controller.fleetBuildRepository.FindAssignedShipModels(id)
	for _, a := range assignedModels {
		a.ShipModel = controller.shipModelRepository.Get(a.ShipModelID)
		if a.ShipModel != nil {
			a.ShipModel.TotalMass = a.ShipModel.CalculateTotalMass()
			a.ResultMass = a.CalculateResultMass()
		}
	}
	c.JSON(http.StatusOK, assignedModels)
}

// GetShipModelAssignment godoc
// @Summary Get a ship model assignment by ID
// @Tags fleet-builds
// @Produce json
// @Param id path string true "ShipModelAssignmentInner ID"
// @Success 200 {object} galaxy.ShipModelAssignment
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-model-assignment/{id} [get]
func (controller *FleetBuildController) GetShipModelAssignment(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")

	log.Printf("The ship model assignment id %s", id)

	assignment := controller.fleetBuildRepository.FindShipModelAssignment(id)

	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ship Model Assignment not found"})
		return
	}

	fleetBuild := controller.fleetBuildRepository.Get(assignment.FleetBuildID)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Broken ship model assignment the fleet build id is wrong."})
		return
	}

	if race.Role != galaxy.RoleAdmin {
		if fleetBuild.RaceId != race.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "This ship model assignment does not belong to your race"})
			return
		}
	}

	assignment.ShipModel = controller.shipModelRepository.Get(assignment.ShipModelID)
	if assignment.ShipModel != nil {
		assignment.ShipModel.TotalMass = assignment.ShipModel.CalculateTotalMass()
		assignment.ResultMass = assignment.CalculateResultMass()
	}

	c.JSON(http.StatusOK, assignment)

}

// AddShipModelAssignment godoc
// @Summary Assign a ship model to a fleet build
// @Tags fleet-builds
// @Accept json
// @Produce json
// @Param assignment body galaxy.ShipModelAssignment true "Assignment data"
// @Success 201 {object} galaxy.ShipModelAssignment
// @Success 200 {object} galaxy.ShipModelAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-model-assignment [post]
func (controller *FleetBuildController) AddShipModelAssignment(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var assignment galaxy.ShipModelAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fleetBuild := controller.fleetBuildRepository.Get(assignment.FleetBuildID)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}
	if race.Role != galaxy.RoleAdmin && fleetBuild.RaceId != race.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "This fleet build does not belong to your race"})
		return
	}

	existingAssignment := controller.fleetBuildRepository.FindAssignedShipModel(assignment.FleetBuildID, assignment.ShipModelID)

	if existingAssignment != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("This ship model assignment already exists (%s, %s)", assignment.FleetBuildID, assignment.ShipModelID)})
		return
	}

	err := controller.fleetBuildRepository.AssignShipModel(&assignment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// UpdateShipModelAssignment godoc
// @Summary Update a ship model assignment
// @Tags fleet-builds
// @Accept json
// @Produce json
// @Param id path string true "ShipModelAssignment ID"
// @Param assignment body galaxy.ShipModelAssignment true "Assignment data"
// @Success 200 {object} galaxy.ShipModelAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-model-assignment/{id} [post]
func (controller *FleetBuildController) UpdateShipModelAssignment(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignmentId := c.Param("id")

	existingAssignment := controller.fleetBuildRepository.FindShipModelAssignment(assignmentId)
	if existingAssignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ship Model Assignment not found"})
		return
	}

	fleetBuild := controller.fleetBuildRepository.Get(existingAssignment.FleetBuildID)
	if fleetBuild != nil && race.Role != galaxy.RoleAdmin && fleetBuild.RaceId != race.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ship Model Assignment not found"})
		return
	}

	var updatedAssignment galaxy.ShipModelAssignment
	if err := c.ShouldBindJSON(&updatedAssignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingAssignment.ShipModelID = updatedAssignment.ShipModelID
	existingAssignment.Amount = updatedAssignment.Amount
	existingAssignment.ResultMass = updatedAssignment.ResultMass

	err := controller.fleetBuildRepository.UpdateShipModelAssignment(existingAssignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingAssignment)
}

// UnassignShipModel godoc
// @Summary Unassign a ship model from a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "ShipModelAssignment ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ship-models-assignment/{id} [delete]
func (controller *FleetBuildController) UnassignShipModel(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignment := controller.fleetBuildRepository.FindShipModelAssignment(c.Param("id"))
	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel assignment not found"})
		return
	}

	if race.Role != galaxy.RoleAdmin {
		fleetBuild := controller.fleetBuildRepository.Get(assignment.FleetBuildID)
		if fleetBuild != nil && fleetBuild.RaceId != race.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "ShipModel assignment not found"})
			return
		}
	}

	controller.fleetBuildRepository.UnassignShipModel(assignment.FleetBuildID, assignment.ShipModelID)
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	fleetBuild.AssignedShipModels = make([]galaxy.ShipModelAssignmentInner, 0, len(assignments))
	for _, a := range assignments {
		shipModel := controller.shipModelRepository.Get(a.ShipModelID)
		if shipModel != nil {
			fleetBuild.AssignedShipModels = append(fleetBuild.AssignedShipModels, galaxy.ShipModelAssignmentInner{
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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
				Owner: race.ID,
			}
			ships = append(ships, ship)
		}
	}

	fleet := galaxy.NewFleet(ships)
	fleet.ID = uuid.New().String()
	fleet.Owner = race.ID
	fmt.Printf("Built fleet %s for user %s with %d ships\n", fleet.ID, race.ID, len(ships))

	controller.fleetRepository.Upsert(fleet)
	controller.fleetRepository.UpsertDivisionFleet(&galaxy.DivisionFleet{
		DivisionId: fleetBuild.DivisionId,
		UserId:     race.ID,
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
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fleetBuildId := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	divisionFleet := controller.fleetRepository.GetDivisionFleet(fleetBuild.DivisionId, race.ID)
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

// GetTechnologies godoc
// @Summary Get researched technologies for a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {object} galaxy.Technologies
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/technologies [get]
func (controller *FleetBuildController) GetTechnologies(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fleetBuild := controller.fleetBuildRepository.Get(c.Param("id"))
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}
	c.JSON(http.StatusOK, fleetBuild.CalculateTechnologies())
}

// CalculateShipTech godoc
// @Summary Calculate ship tech for a ship model assigned to a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "ShipModelAssignment ID"
// @Success 200 {object} galaxy.ShipTech
// @Failure 404 {object} map[string]string
// @Router /ship-models-assignment/{id}/calculate-ship-tech [get]
func (controller *FleetBuildController) CalculateShipTech(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignmentId := c.Param("id")

	assignment := controller.fleetBuildRepository.FindShipModelAssignment(assignmentId)
	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild assignment not found"})
		return
	}

	fleetBuild := controller.fleetBuildRepository.Get(assignment.FleetBuildID)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	shipModel := controller.shipModelRepository.Get(assignment.ShipModelID)
	if shipModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("ShipModel by id %s not found", assignment.ShipModelID)})
		return
	}

	shipTech := fleetBuild.CalculateShipTech(shipModel)
	c.JSON(http.StatusOK, shipTech)
}

type AssignmentShipTech struct {
	AssignmentID string          `json:"assignment_id"`
	ShipTech     galaxy.ShipTech `json:"ship_tech"`
}

// CalculateAssignmentsShipTechs godoc
// @Summary Calculate ship tech for all assignments of a fleet build
// @Tags fleet-builds
// @Produce json
// @Param id path string true "FleetBuild ID"
// @Success 200 {array} AssignmentShipTech
// @Failure 404 {object} map[string]string
// @Router /fleet-builds/{id}/calculate-assignments-ship-techs [get]
func (controller *FleetBuildController) CalculateAssignmentsShipTechs(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fleetBuildId := c.Param("id")
	fleetBuild := controller.fleetBuildRepository.Get(fleetBuildId)
	if fleetBuild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FleetBuild not found"})
		return
	}

	assignments := controller.fleetBuildRepository.FindAssignedShipModels(fleetBuildId)
	result := make([]AssignmentShipTech, 0, len(assignments))
	for _, a := range assignments {
		shipModel := controller.shipModelRepository.Get(a.ShipModelID)
		if shipModel == nil {
			continue
		}
		result = append(result, AssignmentShipTech{
			AssignmentID: a.ID,
			ShipTech:     fleetBuild.CalculateShipTech(shipModel),
		})
	}
	c.JSON(http.StatusOK, result)
}

// ready fleet build

func (controller *FleetBuildController) DivisionGetReadyFleetBuilds(c *gin.Context) {
	// TODO return list of fleet builds that are from all races in this division
}
