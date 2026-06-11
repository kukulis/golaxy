package galaxy

import (
	"reflect"
	"testing"
)

func TestFleetBuildFilter_FromQuery(t *testing.T) {
	testCases := provideFleetBuildFilterFromQueryTestCases()

	for _, tc := range testCases {
		filter := NewFleetBuildFilter()
		filter.FromQuery(tc.query)

		if !reflect.DeepEqual(tc.expectedFleetBuildFilter, filter) {
			t.Errorf("expected %v, got %v", tc.expectedFleetBuildFilter, filter)
		}

	}
}

type fleetBuildFilterFromQueryTestCase struct {
	query                    map[string][]string
	expectedFleetBuildFilter *FleetBuildFilter
}

func provideFleetBuildFilterFromQueryTestCases() []fleetBuildFilterFromQueryTestCase {
	return []fleetBuildFilterFromQueryTestCase{
		{
			query:                    map[string][]string{},
			expectedFleetBuildFilter: NewFleetBuildFilter(),
		},
	}
}
