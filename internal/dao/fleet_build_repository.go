package dao

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
	"maps"
	"slices"
	"strings"
)

type FleetBuildRepository struct {
	fleetBuildMap map[string]*galaxy.FleetBuild
	// not very effective way, but this repository is for DEV purposes only
	fleetBuildToShipModels []*galaxy.FleetBuildToShipModel
}

func NewFleetBuildRepository() *FleetBuildRepository {
	return &FleetBuildRepository{
		fleetBuildMap: make(map[string]*galaxy.FleetBuild),
	}
}

func (r *FleetBuildRepository) Get(id string) *galaxy.FleetBuild {
	return r.fleetBuildMap[id]
}

func (r *FleetBuildRepository) GetAll() []*galaxy.FleetBuild {
	fleetBuilds := slices.Collect(maps.Values(r.fleetBuildMap))

	slices.SortFunc(fleetBuilds, func(a, b *galaxy.FleetBuild) int {
		return strings.Compare(a.ID, b.ID)
	})

	return fleetBuilds
}

func (r *FleetBuildRepository) Upsert(fleetBuild *galaxy.FleetBuild) {
	r.fleetBuildMap[fleetBuild.ID] = fleetBuild
}

func (r *FleetBuildRepository) Delete(id string) {
	delete(r.fleetBuildMap, id)
}

func (r *FleetBuildRepository) FindAssignedShipModels(fleetBuildId string) []*galaxy.FleetBuildToShipModel {
	return util.ArrayFilter(r.fleetBuildToShipModels, func(b2s *galaxy.FleetBuildToShipModel) bool { return b2s.FleetBuildID == fleetBuildId })
}

func (r *FleetBuildRepository) AssignShipModel(fleetBuild2ShipModel *galaxy.FleetBuildToShipModel) bool {
	for _, b2s := range r.fleetBuildToShipModels {
		if b2s.ShipModelID == fleetBuild2ShipModel.ShipModelID && b2s.FleetBuildID == fleetBuild2ShipModel.FleetBuildID {
			return false
		}
	}

	r.fleetBuildToShipModels = append(r.fleetBuildToShipModels, fleetBuild2ShipModel)

	return true
}

func (r *FleetBuildRepository) UnassignShipModel(fleetBuildId, shipModelId string) bool {
	foundIndex := -1
	for i, b2s := range r.fleetBuildToShipModels {
		if b2s.ShipModelID == fleetBuildId && b2s.FleetBuildID == shipModelId {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return false
	}

	slices.Delete(r.fleetBuildToShipModels, foundIndex, foundIndex)

	return true
}

func (r *FleetBuildRepository) ResetData() {
	r.fleetBuildMap = make(map[string]*galaxy.FleetBuild)
	r.fleetBuildToShipModels = nil
}
