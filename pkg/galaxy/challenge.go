package galaxy

type Challenge struct {
	// stored to db
	ID               string `json:"id"`
	ChallengerRaceId string `json:"challenger_race_id"`
	FleetBuildAId    string `json:"fleet_build_a_id"`
	FleetBuildBId    string `json:"fleet_build_b_id"`
	BattleReportId   string `json:"battle_report_id"`

	// Not stored to db
	FleetBuildA *FleetBuild
	FleetBuildB *FleetBuild
}

func NewChallenge(id string, raceId string) *Challenge {
	return &Challenge{
		ID:               id,
		ChallengerRaceId: raceId,
		FleetBuildAId:    "",
		FleetBuildBId:    "",
		BattleReportId:   "",
		FleetBuildA:      nil,
		FleetBuildB:      nil,
	}
}
