package di

import (
	"glaktika.eu/galaktika/internal/api"
	"glaktika.eu/galaktika/internal/dao"
)

var BattleRepositoryInstance *dao.BattleRepository
var BattleControllerInstance *api.BattleController

func CreateSingletons() {
	BattleRepositoryInstance = dao.NewBattleRepository()
	BattleControllerInstance = api.NewBattleController(BattleRepositoryInstance)
}
