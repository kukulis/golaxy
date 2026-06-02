package galaxy

import "time"

const (
	ChallengeStatusPending  = "pending"
	ChallengeStatusRejected = "rejected"
	ChallengeStatusExecuted = "executed"
)

type Challenge struct {
	// Unique identifier. Ignored in update.
	ID string `json:"id"`
	// Division this challenge belongs to. Ignored in update.
	DivisionId string `json:"division_id"`
	// Race that issued the challenge. Ignored in update.
	ChallengerRaceId string `json:"challenger_race_id"`
	// Race that received the challenge. Ignored in update.
	ChallengeeRaceId string `json:"challengee_race_id"`
	// Fleet build selected by the challenger. Updatable by challenger and admin.
	FleetBuildAId string `json:"fleet_build_a_id"`
	// Fleet build selected by the challengee. Updatable by challengee and admin.
	FleetBuildBId string `json:"fleet_build_b_id"`
	// Challenger readiness flag. Updatable by challenger and admin.
	ReadyA bool `json:"ready_a"`
	// Challengee readiness flag. Updatable by challengee and admin.
	ReadyB bool `json:"ready_b"`
	// When the challenge was created. Ignored in update.
	CreatedAt time.Time `json:"created_at"`
	// When the challenger accepted. Updatable by admin only.
	AcceptedAAt *time.Time `json:"accepted_a_at"`
	// When the challengee accepted. Updatable by admin only.
	AcceptedBAt *time.Time `json:"accepted_b_at"`
	// Challenge status (pending/rejected/executed). Updatable by challengee (rejected) and admin.
	Status string `json:"status"`
	// When the challenge was executed. Updatable by admin only.
	ExecutedAt *time.Time `json:"executed_at"`
	// Battle report assigned after execution. Updatable by admin only.
	BattleReportId string `json:"battle_report_id"`

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

type ChallengesFilter struct {
	ChallengerId string
	ChallengeeId string
	AcceptedA    bool
	AcceptedB    bool

	// TODO filter by status
}

func (f *ChallengesFilter) Match(c *Challenge) bool {
	if f.ChallengerId != "" && c.ChallengerRaceId != f.ChallengerId {
		return false
	}
	if f.ChallengeeId != "" && c.ChallengeeRaceId != f.ChallengeeId {
		return false
	}
	if f.AcceptedA && c.AcceptedAAt == nil {
		return false
	}
	if f.AcceptedB && c.AcceptedBAt == nil {
		return false
	}
	return true
}
