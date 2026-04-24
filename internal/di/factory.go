package di

import (
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func NewDivisionRepository() *dao.DivisionRepository {
	divisionRepository := dao.NewDivisionRepository()

	divisionRepository.Upsert(&galaxy.Division{ID: "alpha", ResourcesAmount: 500, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "beta", ResourcesAmount: 750, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})
	divisionRepository.Upsert(&galaxy.Division{ID: "gamma", ResourcesAmount: 1000, TechAttack: 1, TechDefense: 1, TechEngines: 1, TechCargo: 1})

	return divisionRepository
}
