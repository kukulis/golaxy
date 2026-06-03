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
	ReadyA       bool
	ReadyB       bool
	Status       string
	DivisionId   string
}

func NewChallengesFilter() *ChallengesFilter {
	return &ChallengesFilter{}
}

// Match returns true if the given challenge matches the filter parameters
func (f *ChallengesFilter) Match(c *Challenge) bool {
	if f.ChallengerId != "" && c.ChallengerRaceId != f.ChallengerId {
		return false
	}
	if f.ChallengeeId != "" && c.ChallengeeRaceId != f.ChallengeeId {
		return false
	}
	if f.ReadyA && !c.ReadyA {
		return false
	}
	if f.ReadyB && !c.ReadyB {
		return false
	}

	// TODO status
	// TODO divisionId

	return true
}

// FromQuery fluent query parser
func (f *ChallengesFilter) FromQuery(query map[string][]string) *ChallengesFilter {
	challengerIds, ok := query["challenger_id"]
	if ok && len(challengerIds) > 0 {
		f.ChallengerId = challengerIds[0]
	}

	challengeeIds, ok := query["challengee_id"]
	if ok && len(challengeeIds) > 0 {
		f.ChallengeeId = challengeeIds[0]
	}

	readyAs, ok := query["ready_a"]
	if ok && len(readyAs) > 0 {
		f.ReadyA = readyAs[0] == "1"
	}

	readyBs, ok := query["ready_b"]
	if ok && len(readyBs) > 0 {
		f.ReadyB = readyBs[0] == "1"
	}

	divisionIds, ok := query["division_id"]
	if ok && len(divisionIds) > 0 {
		f.DivisionId = divisionIds[0]
	}

	statuses, ok := query["status"]
	if ok && len(statuses) > 0 {
		f.Status = statuses[0]
	}

	return f
}
