package dao

import (
	"errors"
	"maps"
	"slices"
	"strings"

	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
)

type FleetBuildRepository struct {
	fleetBuildMap map[string]*galaxy.FleetBuild
	// not very effective way, but this repository is for DEV purposes only
	fleetBuildToShipModels []*galaxy.ShipModelAssignment
}

func NewFleetBuildRepository() *FleetBuildRepository {
	return &FleetBuildRepository{
		fleetBuildMap: make(map[string]*galaxy.FleetBuild, 0),
	}
}

func (r *FleetBuildRepository) Get(id string) *galaxy.FleetBuild {
	return r.fleetBuildMap[id]
}

func (r *FleetBuildRepository) GetList(divisionId, raceId string) ([]*galaxy.FleetBuild, error) {
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

	return fleetBuilds, nil
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
func (r *FleetBuildRepository) AssignShipModel(fleetBuild2ShipModel *galaxy.ShipModelAssignment) error {
	r.fleetBuildToShipModels = append(r.fleetBuildToShipModels, fleetBuild2ShipModel)

	return nil
}

func (r *FleetBuildRepository) UpdateShipModelAssignment(updatedAssignment *galaxy.ShipModelAssignment) error {
	for _, existingAssignment := range r.fleetBuildToShipModels {
		if existingAssignment.ID == updatedAssignment.ID {
			existingAssignment.ShipModelID = updatedAssignment.ShipModelID
			existingAssignment.Amount = updatedAssignment.Amount
			existingAssignment.ResultMass = updatedAssignment.ResultMass

			return nil
		}
	}

	return errors.New("Could not updatedAssignment assignment with id " + updatedAssignment.ID)
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
