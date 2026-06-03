package galaxy

import (
	"reflect"
	"testing"
)

func TestChallengeFilterFromQuery(t *testing.T) {
	testCases := prepareTestCasesForChallengeFilterFromQuery()

	for _, tc := range testCases {
		filter := &ChallengesFilter{}

		filter.FromQuery(tc.query)

		if !reflect.DeepEqual(tc.expectedFilter, filter) {
			t.Errorf("Test case %s : expected %v, got %v", tc.name, tc.expectedFilter, filter)
		}
	}
}

type filterFromQueryTestCase struct {
	name           string
	query          map[string][]string
	expectedFilter *ChallengesFilter
}

func prepareTestCasesForChallengeFilterFromQuery() []filterFromQueryTestCase {
	return []filterFromQueryTestCase{
		{
			name:  "wrong parameters",
			query: map[string][]string{"aaa": {"bbb"}},
			expectedFilter: &ChallengesFilter{
				ChallengerId: "",
				ChallengeeId: "",
				ReadyA:       false,
				ReadyB:       false,
			},
		},
		{
			name: "wrong parameters",
			query: map[string][]string{
				"challenger_id": {"123123"},
				"challengee_id": {"123124"},
				"ready_a":       {"1"},
				"ready_b":       {"0"},
				"division_id":   {"111"},
				"status":        {"executed"},
			},
			expectedFilter: &ChallengesFilter{
				ChallengerId: "123123",
				ChallengeeId: "123124",
				ReadyA:       true,
				ReadyB:       false,
				DivisionId:   "111",
				Status:       "executed",
			},
		},
	}
}
