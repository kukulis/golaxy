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

func (controller *BattleController) GetBattle(c *gin.Context) {
	// TODO id from request query
	id := "1"
	c.JSON(http.StatusOK, controller.battleRepository.GetBattle(id))
}
