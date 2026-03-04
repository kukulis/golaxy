package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"net/http"
)

type BattleController struct {
	battleRepository *dao.BattleRepository
}

func NewBattleController(repository *dao.BattleRepository) *BattleController {
	return &BattleController{battleRepository: repository}
}

// GetBattle godoc
// @Summary Get battle
// @Tags battles
// @Produce json
// @Success 200 {object} galaxy.Battle
// @Router /battle [get]
func (controller *BattleController) GetBattle(c *gin.Context) {
	// TODO id from request query
	id := "1"
	c.JSON(http.StatusOK, controller.battleRepository.GetBattle(id))
}
