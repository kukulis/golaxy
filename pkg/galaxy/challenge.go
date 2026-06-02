package galaxy

import "time"

const (
	ChallengeStatusPending  = "pending"
	ChallengeStatusRejected = "rejected"
	ChallengeStatusExecuted = "executed"
)

type Challenge struct {
	// stored to db
	ID               string     `json:"id"`
	DivisionId       string     `json:"division_id"`
	ChallengerRaceId string     `json:"challenger_race_id"`
	ChallengeeRaceId string     `json:"challengee_race_id"`
	FleetBuildAId    string     `json:"fleet_build_a_id"`
	FleetBuildBId    string     `json:"fleet_build_b_id"`
	ReadyA           bool       `json:"ready_a"`
	ReadyB           bool       `json:"ready_b"`
	CreatedAt        time.Time  `json:"created_at"`
	AcceptedAAt      *time.Time `json:"accepted_a_at"`
	AcceptedBAt      *time.Time `json:"accepted_b_at"`
	Status           string     `json:"status"`
	ExecutedAt       *time.Time `json:"executed_at"`
	BattleReportId   string     `json:"battle_report_id"`

	// Not stored to db
	FleetBuildA *FleetBuild
	FleetBuildB *FleetBuild
}

func NewChallenge(
	id string,
	challengerRaceId string,
	challengeeRaceId string,
	divisionId string,
	createdAt time.Time,
) *Challenge {
	return &Challenge{
		ID:               id,
		DivisionId:       divisionId,
		ChallengerRaceId: challengerRaceId,
		ChallengeeRaceId: challengeeRaceId,
		CreatedAt:        createdAt,
		Status:           ChallengeStatusPending,
	}
}
