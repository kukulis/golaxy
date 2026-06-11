package galaxy

type FleetBuildFilter struct {
	RaceId     string
	DivisionId string
}

func NewFleetBuildFilter() *FleetBuildFilter {
	return &FleetBuildFilter{
		RaceId:     "",
		DivisionId: "",
	}
}

func (f *FleetBuildFilter) SetRaceId(id string) *FleetBuildFilter {
	f.RaceId = id

	return f
}

func (f *FleetBuildFilter) SetDivisionId(id string) *FleetBuildFilter {
	f.DivisionId = id

	return f
}

func (f *FleetBuildFilter) FromQuery(query map[string][]string) *FleetBuildFilter {

	raceIds, ok := query["race_id"]

	if ok && len(raceIds) > 0 {
		f.RaceId = raceIds[0]
	}

	divisionIds, ok := query["division_id"]
	if ok && len(divisionIds) > 0 {
		f.DivisionId = divisionIds[0]
	}

	return f
}
