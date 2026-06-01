package web

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"log"
	"net/http"
)

type WebController struct {
	fleetBuildRepository *dao.FleetBuildRepository
}

func NewWebController(fleetBuildRepository *dao.FleetBuildRepository) *WebController {
	return &WebController{fleetBuildRepository: fleetBuildRepository}
}

func (wc *WebController) RenderDivision(c *gin.Context) {
	divisionId := c.Param("divisionId")
	c.HTML(http.StatusOK, "division_main", gin.H{"DivisionId": divisionId})
}

func (wc *WebController) RenderFleetBuilds(c *gin.Context) {
	divisionId := c.Param("divisionId")
	c.HTML(http.StatusOK, "fleet_builds", gin.H{"DivisionId": divisionId})
}

func (wc *WebController) RenderFleetBuild(c *gin.Context) {
	id := c.Param("id")
	fleetBuild := wc.fleetBuildRepository.Get(id)
	if fleetBuild == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "fleet_build_main", gin.H{"Id": id, "DivisionId": fleetBuild.DivisionId})
}

func (wc *WebController) RenderFleetBuildEdit(c *gin.Context) {
	id := c.Param("id")
	fleetBuild := wc.fleetBuildRepository.Get(id)
	if fleetBuild == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "fleet_build_edit", gin.H{"Id": id, "DivisionId": fleetBuild.DivisionId})
}

func (wc *WebController) RenderShipModelAssignment(c *gin.Context) {
	id := c.Param("id")

	shipModelAssignment := wc.fleetBuildRepository.FindShipModelAssignment(id)
	if shipModelAssignment == nil {
		c.Status(http.StatusNotFound)
		return
	}

	fleetBuild := wc.fleetBuildRepository.Get(shipModelAssignment.FleetBuildID)
	if fleetBuild == nil {
		log.Printf("Error: the ship model assignment contains wrong fleet build id %s", shipModelAssignment.FleetBuildID)

		c.Status(http.StatusNotFound)
		return
	}

	log.Printf("RenderShipModelAssignment:fleetBuild.DivisionId: %s", fleetBuild.DivisionId)

	c.HTML(http.StatusOK, "ship_model_assignment", gin.H{"DivisionId": fleetBuild.DivisionId, "FleetBuildId": fleetBuild.ID, "ShipModelAssignmentId": shipModelAssignment.ID})
}

func (wc *WebController) RenderShipModelList(c *gin.Context) {
	c.HTML(http.StatusOK, "ship_model_list", nil)
}

func (wc *WebController) RenderShipModelDetails(c *gin.Context) {
	c.HTML(http.StatusOK, "ship_model_details", wc.shipModelData(c))
}

func (wc *WebController) RenderShipModelEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "ship_model_edit", wc.shipModelData(c))
}

func (wc *WebController) shipModelData(c *gin.Context) gin.H {
	data := gin.H{"Id": c.Param("id")}
	fleetBuildId := c.Query("FleetBuildId")
	if fleetBuildId != "" {
		fleetBuild := wc.fleetBuildRepository.Get(fleetBuildId)
		if fleetBuild != nil {
			data["FleetBuildId"] = fleetBuildId
			data["DivisionId"] = fleetBuild.DivisionId
		}
	}
	return data
}

func (wc *WebController) RenderChallenges(c *gin.Context) {
	c.HTML(http.StatusOK, "challenges", nil)
}

func (wc *WebController) RenderChallengeEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "challenge_edit", gin.H{"Id": c.Param("id")})
}
