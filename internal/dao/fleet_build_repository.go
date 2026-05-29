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
	fleetBuildToShipModels []*galaxy.ShipModelAssignment
}

func NewFleetBuildRepository() *FleetBuildRepository {
	return &FleetBuildRepository{
		fleetBuildMap: make(map[string]*galaxy.FleetBuild),
	}
}

func (r *FleetBuildRepository) Get(id string) *galaxy.FleetBuild {
	return r.fleetBuildMap[id]
}

func (r *FleetBuildRepository) GetAll(divisionId, raceId string) []*galaxy.FleetBuild {
	fleetBuilds := slices.Collect(maps.Values(r.fleetBuildMap))

	if divisionId != "" {
		fleetBuilds = util.ArrayFilter(fleetBuilds, func(b *galaxy.FleetBuild) bool { return b.DivisionId == divisionId })
	}
	if raceId != "" {
		fleetBuilds = util.ArrayFilter(fleetBuilds, func(b *galaxy.FleetBuild) bool { return b.RaceId == raceId })
	}

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

func (r *FleetBuildRepository) FindAssignedShipModels(fleetBuildId string) []*galaxy.ShipModelAssignment {
	return util.ArrayFilter(r.fleetBuildToShipModels, func(b2s *galaxy.ShipModelAssignment) bool { return b2s.FleetBuildID == fleetBuildId })
}

func (r *FleetBuildRepository) FindShipModelAssignment(shipModelAssignmentId string) *galaxy.ShipModelAssignment {
	for _, b2m := range r.fleetBuildToShipModels {
		if b2m.ID == shipModelAssignmentId {
			return b2m
		}
	}

	return nil
}

func (r *FleetBuildRepository) FindAssignedShipModel(fleetBuildId, shipModelId string) *galaxy.ShipModelAssignment {
	for _, b2m := range r.fleetBuildToShipModels {
		if b2m.FleetBuildID == fleetBuildId && b2m.ShipModelID == shipModelId {
			return b2m
		}
	}

	return nil
}

// AssignShipModel assigns a ship model to a fleet build (upsert operation).
// Returns true if a new assignment was created, false if an existing assignment was updated.
func (r *FleetBuildRepository) AssignShipModel(fleetBuild2ShipModel *galaxy.ShipModelAssignment) bool {

	// search for existing shipModels assignments
	for _, b2s := range r.fleetBuildToShipModels {
		if b2s.ShipModelID == fleetBuild2ShipModel.ShipModelID && b2s.FleetBuildID == fleetBuild2ShipModel.FleetBuildID {
			b2s.Amount = fleetBuild2ShipModel.Amount
			b2s.ResultMass = fleetBuild2ShipModel.ResultMass
			b2s.ShipModel = fleetBuild2ShipModel.ShipModel

			// Updated existing
			return false
		}
	}

	// Create new assignment
	r.fleetBuildToShipModels = append(r.fleetBuildToShipModels, fleetBuild2ShipModel)

	// Created new
	return true
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
