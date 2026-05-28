package api

import (
	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"net/http"
)

type BattleController struct {
	authenticationManager AuthenticationManager
	battleRepository      *dao.BattleRepository
}

func NewBattleController(authenticationManager AuthenticationManager, repository *dao.BattleRepository) *BattleController {
	return &BattleController{authenticationManager: authenticationManager, battleRepository: repository}
}

// GetBattle godoc
// @Summary Get battle
// @Tags battles
// @Produce json
// @Success 200 {object} galaxy.Battle
// @Router /battle [get]
func (controller *BattleController) GetBattle(c *gin.Context) {
	race := controller.authenticationManager.AuthenticateFromContext(c)
	if race == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// TODO id from request query
	id := "1"
	c.JSON(http.StatusOK, controller.battleRepository.GetBattle(id))
}
