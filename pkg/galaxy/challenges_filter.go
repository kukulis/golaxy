package galaxy

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
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.DivisionId != "" && c.DivisionId != f.DivisionId {
		return false
	}

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
