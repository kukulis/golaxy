package dao

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"maps"
	"slices"
	"strings"
)

type ShipModelRepository struct {
	shipModelMap map[string]*galaxy.ShipModel
}

func NewShipModelRepository() *ShipModelRepository {
	return &ShipModelRepository{
		shipModelMap: make(map[string]*galaxy.ShipModel),
	}
}

func (r *ShipModelRepository) Get(id string) *galaxy.ShipModel {
	return r.shipModelMap[id]
}

func (r *ShipModelRepository) GetAll() []*galaxy.ShipModel {
	shipModels := slices.Collect(maps.Values(r.shipModelMap))

	slices.SortFunc(shipModels, func(a, b *galaxy.ShipModel) int {
		return strings.Compare(a.ID, b.ID)
	})

	return shipModels
}

func (r *ShipModelRepository) Upsert(shipModel *galaxy.ShipModel) {
	r.shipModelMap[shipModel.ID] = shipModel
}

func (r *ShipModelRepository) Delete(id string) {
	delete(r.shipModelMap, id)
}

func (r *ShipModelRepository) ResetData() {
	r.shipModelMap = make(map[string]*galaxy.ShipModel)
}
