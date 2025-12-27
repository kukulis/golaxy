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

func (controller *DivisionController) GetDivision(c *gin.Context) {
	id := c.Param("id")
	division := controller.divisionRepository.Get(id)
	if division == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}
	c.JSON(http.StatusOK, division)
}

func (controller *DivisionController) GetAllDivisions(c *gin.Context) {
	c.JSON(http.StatusOK, controller.divisionRepository.GetAll())
}

func (controller *DivisionController) CreateDivision(c *gin.Context) {
	var division galaxy.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.divisionRepository.Upsert(&division)
	c.JSON(http.StatusCreated, division)
}

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
