package dao

import "glaktika.eu/galaktika/pkg/galaxy"

type fleetKey struct {
	DivisionId string
	UserId     string
}

type FleetRepository struct {
	fleetMap       map[string]*galaxy.Fleet
	divisionFleets map[fleetKey]*galaxy.DivisionFleet
}

func NewFleetRepository() *FleetRepository {
	return &FleetRepository{
		fleetMap:       make(map[string]*galaxy.Fleet),
		divisionFleets: make(map[fleetKey]*galaxy.DivisionFleet),
	}
}

func (r *FleetRepository) Get(id string) *galaxy.Fleet {
	return r.fleetMap[id]
}

func (r *FleetRepository) Upsert(fleet *galaxy.Fleet) {
	r.fleetMap[fleet.ID] = fleet
}

func (r *FleetRepository) GetDivisionFleet(divisionId, userId string) *galaxy.DivisionFleet {
	return r.divisionFleets[fleetKey{DivisionId: divisionId, UserId: userId}]
}

func (r *FleetRepository) UpsertDivisionFleet(df *galaxy.DivisionFleet) {
	r.divisionFleets[fleetKey{DivisionId: df.DivisionId, UserId: df.UserId}] = df
}

func (r *FleetRepository) ResetData() {
	r.fleetMap = make(map[string]*galaxy.Fleet)
	r.divisionFleets = make(map[fleetKey]*galaxy.DivisionFleet)
}
