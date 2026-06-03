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
				AcceptedA:    false,
				AcceptedB:    false,
			},
		},
	}
}
