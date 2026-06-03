package galaxy

import "testing"

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

		// TODO more test cases
	}
}
