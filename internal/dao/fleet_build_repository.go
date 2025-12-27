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

func (r *FleetBuildRepository) FindAssignedShipModel(fleetBuildId, shipModelId string) *galaxy.FleetBuildToShipModel {
	for _, b2m := range r.fleetBuildToShipModels {
		if b2m.FleetBuildID == fleetBuildId && b2m.ShipModelID == shipModelId {
			return b2m
		}
	}

	return nil
}

// AssignShipModel assigns a ship model to a fleet build (upsert operation).
// Returns true if a new assignment was created, false if an existing assignment was updated.
func (r *FleetBuildRepository) AssignShipModel(fleetBuild2ShipModel *galaxy.FleetBuildToShipModel) bool {
	for _, b2s := range r.fleetBuildToShipModels {
		if b2s.ShipModelID == fleetBuild2ShipModel.ShipModelID && b2s.FleetBuildID == fleetBuild2ShipModel.FleetBuildID {
			// Update existing assignment
			b2s.Amount = fleetBuild2ShipModel.Amount
			b2s.ResultMass = fleetBuild2ShipModel.ResultMass
			b2s.ShipModel = fleetBuild2ShipModel.ShipModel
			return false // Updated existing
		}
	}

	// Create new assignment
	r.fleetBuildToShipModels = append(r.fleetBuildToShipModels, fleetBuild2ShipModel)
	return true // Created new
}

func (r *FleetBuildRepository) UnassignShipModel(fleetBuildId, shipModelId string) bool {
	foundIndex := -1
	for i, b2s := range r.fleetBuildToShipModels {
		if b2s.ShipModelID == shipModelId && b2s.FleetBuildID == fleetBuildId {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return false
	}

	r.fleetBuildToShipModels = slices.Delete(r.fleetBuildToShipModels, foundIndex, foundIndex+1)

	return true
}

func (r *FleetBuildRepository) ResetData() {
	r.fleetBuildMap = make(map[string]*galaxy.FleetBuild)
	r.fleetBuildToShipModels = nil
}
