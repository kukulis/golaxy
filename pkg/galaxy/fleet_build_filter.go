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
	// TODO
	return f
}
