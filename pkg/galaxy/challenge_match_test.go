package galaxy

import (
	"testing"
	"time"
)

func TestChallengesFilter_Match(t *testing.T) {
	for _, testCase := range challengesFilterMatchTestCases() {
		if testCase.result != testCase.filter.Match(testCase.challenge) {
			t.Errorf("%s: expected %v, got %v", testCase.name, testCase.result, testCase.filter.Match(testCase.challenge))
		}
	}
}

type challengesFilterMatchTestCase struct {
	name      string
	challenge *Challenge
	filter    *ChallengesFilter
	result    bool
}

func challengesFilterMatchTestCases() []challengesFilterMatchTestCase {
	return []challengesFilterMatchTestCase{
		{
			name:      "empty",
			challenge: &Challenge{},
			filter:    &ChallengesFilter{},
			result:    true,
		},
		{
			name: "by challenger",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           false,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "",
				ReadyA:       false,
				ReadyB:       false,
				Status:       "",
				DivisionId:   "",
			},
			result: true,
		},
		{
			name: "by challenger wrong",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           false,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch780",
				ChallengeeId: "",
				ReadyA:       false,
				ReadyB:       false,
				Status:       "",
				DivisionId:   "",
			},
			result: false,
		},
		{
			name: "by status wrong",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           false,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "",
				ReadyA:       false,
				ReadyB:       false,
				Status:       "executing",
				DivisionId:   "",
			},
			result: false,
		},
		{
			name: "by division wrong",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           false,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "",
				ReadyA:       false,
				ReadyB:       false,
				Status:       "",
				DivisionId:   "d111",
			},
			result: false,
		},

		{
			name: "by all ",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           true,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "ch799",
				ReadyA:       true,
				ReadyB:       true,
				Status:       "pending",
				DivisionId:   "d456",
			},
			result: true,
		},

		{
			name: "by ready b wrong",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           true,
				ReadyB:           false,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "ch799",
				ReadyA:       true,
				ReadyB:       true,
				Status:       "pending",
				DivisionId:   "d456",
			},
			result: false,
		},
		{
			name: "by ready a wrong",
			challenge: &Challenge{
				ID:               "123",
				DivisionId:       "d456",
				ChallengerRaceId: "ch789",
				ChallengeeRaceId: "ch799",
				FleetBuildAId:    "f456",
				FleetBuildBId:    "f457",
				ReadyA:           false,
				ReadyB:           true,
				CreatedAt:        time.Time{},
				AcceptedAAt:      nil,
				AcceptedBAt:      nil,
				Status:           "pending",
				ExecutedAt:       nil,
				BattleReportId:   "",
				FleetBuildA:      nil,
				FleetBuildB:      nil,
			},
			filter: &ChallengesFilter{
				ChallengerId: "ch789",
				ChallengeeId: "ch799",
				ReadyA:       true,
				ReadyB:       true,
				Status:       "pending",
				DivisionId:   "d456",
			},
			result: false,
		},
	}
}
