package api

import (
	"encoding/json"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetFleetBuildsList(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token", Role: galaxy.RolePlayer})
	raceRepo.Upsert(&galaxy.Race{ID: "race-a", Name: "Race A", Token: "test2-token", Role: galaxy.RoleAdmin})
	auth := NewMemoryAuthenticationManager(raceRepo)
	fleetBuildRepository := dao.NewFleetBuildRepository()

	// TODO put more data to repository
	controller := NewFleetBuildController(auth, fleetBuildRepository, nil, nil, nil)
	router := setupFleetBuildRouter(controller)

	testCases := provideFleetBuildListTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/fleet-builds", nil)
			if testCase.authHeader != "" {
				req.Header.Set("Authorization", testCase.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// get results
			if w.Code != testCase.expectedResponseCode {
				t.Errorf("test case %s: expected status %d, gotFleetBuilds %d", testCase.name, testCase.expectedResponseCode, w.Code)
				return
			}

			if testCase.expectedFleetBuilds != nil {
				var gotFleetBuilds []*galaxy.FleetBuild
				if err := json.Unmarshal(w.Body.Bytes(), &gotFleetBuilds); err != nil {
					t.Fatalf("test case %s failed to unmarshal response: %v", testCase.name, err)
					return
				}

				if !reflect.DeepEqual(gotFleetBuilds, testCase.expectedFleetBuilds) {
					t.Fatalf("test case %s expected results: %v gotFleetBuilds %v", testCase.name, testCase.expectedFleetBuilds, gotFleetBuilds)
				}
			}
		})
	}
}

type fleetBuildListTestCase struct {
	name                 string
	query                string
	authHeader           string
	expectedResponseCode int
	expectedFleetBuilds  []*galaxy.FleetBuild
}

func provideFleetBuildListTestCases() []fleetBuildListTestCase {
	return []fleetBuildListTestCase{
		{name: "empty",
			query:                "",
			authHeader:           "",
			expectedResponseCode: http.StatusUnauthorized,
			expectedFleetBuilds:  nil,
		},

		// TODO more test cases
	}
}
