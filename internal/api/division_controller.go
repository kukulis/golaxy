package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
)

type DivisionController struct {
	divisionRepository *dao.DivisionRepository
}

func NewDivisionController(repository *dao.DivisionRepository) *DivisionController {
	return &DivisionController{divisionRepository: repository}
}

// GetDivision godoc
// @Summary Get a division by ID
// @Tags divisions
// @Produce json
// @Param id path string true "Division ID"
// @Success 200 {object} galaxy.Division
// @Failure 404 {object} map[string]string
// @Router /divisions/{id} [get]
func (controller *DivisionController) GetDivision(c *gin.Context) {
	id := c.Param("id")
	division := controller.divisionRepository.Get(id)
	if division == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}
	c.JSON(http.StatusOK, division)
}

// GetAllDivisions godoc
// @Summary List all divisions
// @Tags divisions
// @Produce json
// @Success 200 {array} galaxy.Division
// @Router /divisions [get]
func (controller *DivisionController) GetAllDivisions(c *gin.Context) {
	c.JSON(http.StatusOK, controller.divisionRepository.GetAll())
}

// CreateDivision godoc
// @Summary Create a division
// @Tags divisions
// @Accept json
// @Produce json
// @Param division body galaxy.Division true "Division data"
// @Success 201 {object} galaxy.Division
// @Failure 400 {object} map[string]string
// @Router /divisions [post]
func (controller *DivisionController) CreateDivision(c *gin.Context) {
	var division galaxy.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.divisionRepository.Upsert(&division)
	c.JSON(http.StatusCreated, division)
}

// UpdateDivision godoc
// @Summary Update a division
// @Tags divisions
// @Accept json
// @Produce json
// @Param id path string true "Division ID"
// @Param division body galaxy.Division true "Division data"
// @Success 200 {object} galaxy.Division
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /divisions/{id} [put]
func (controller *DivisionController) UpdateDivision(c *gin.Context) {
	id := c.Param("id")
	existing := controller.divisionRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}

	var division galaxy.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	division.ID = id
	controller.divisionRepository.Upsert(&division)
	c.JSON(http.StatusOK, division)
}

// DeleteDivision godoc
// @Summary Delete a division
// @Tags divisions
// @Produce json
// @Param id path string true "Division ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /divisions/{id} [delete]
func (controller *DivisionController) DeleteDivision(c *gin.Context) {
	id := c.Param("id")
	existing := controller.divisionRepository.Get(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}

	controller.divisionRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "Division deleted successfully"})
}
