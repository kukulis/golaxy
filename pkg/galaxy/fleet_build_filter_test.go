package galaxy

import (
	"glaktika.eu/galaktika/pkg/util"
	"reflect"
	"testing"
)

func TestFleetBuildFilter_FromQuery(t *testing.T) {
	testCases := provideFleetBuildFilterFromQueryTestCases()

	for _, tc := range testCases {
		filter := NewFleetBuildFilter()
		filter.FromQuery(tc.query)

		if !reflect.DeepEqual(tc.expectedFleetBuildFilter, filter) {

			diff := util.DiffStruct(filter, tc.expectedFleetBuildFilter)
			t.Errorf("Failed test case %s:\n%s", tc.name, diff)
		}

	}
}

type fleetBuildFilterFromQueryTestCase struct {
	name                     string
	query                    map[string][]string
	expectedFleetBuildFilter *FleetBuildFilter
}

func provideFleetBuildFilterFromQueryTestCases() []fleetBuildFilterFromQueryTestCase {
	return []fleetBuildFilterFromQueryTestCase{
		{
			name:                     "empty",
			query:                    map[string][]string{},
			expectedFleetBuildFilter: NewFleetBuildFilter(),
		},
		{
			name:                     "race id given",
			query:                    map[string][]string{"race_id": {"10asdf"}},
			expectedFleetBuildFilter: NewFleetBuildFilter().SetRaceId("10asdf"),
		},
		{
			name:                     "division id given",
			query:                    map[string][]string{"division_id": {"grybas"}},
			expectedFleetBuildFilter: NewFleetBuildFilter().SetDivisionId("grybas"),
		},
		{
			name:  "both given",
			query: map[string][]string{"division_id": {"alpha"}, "race_id": {"rex"}},
			expectedFleetBuildFilter: NewFleetBuildFilter().
				SetDivisionId("alpha").
				SetRaceId("rex"),
		},

		{
			name:  "trash given",
			query: map[string][]string{"rikutygas": {"dyguzas"}},
			expectedFleetBuildFilter: NewFleetBuildFilter().
				SetDivisionId("").
				SetRaceId(""),
		},
	}
}
