package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
)

func TestGetFleetBuildsList(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token", Role: galaxy.RolePlayer})
	raceRepo.Upsert(&galaxy.Race{ID: "race-a", Name: "Race A", Token: "admin-token", Role: galaxy.RoleAdmin})
	auth := NewMemoryAuthenticationManager(raceRepo)
	fleetBuildRepository := dao.NewFleetBuildRepository()

	fleetBuildRepository.Upsert(&galaxy.FleetBuild{ID: "fb-1", DivisionId: "div-1", RaceId: "race-1"})
	fleetBuildRepository.Upsert(&galaxy.FleetBuild{ID: "fb-2", DivisionId: "div-1", RaceId: "race-1"})
	fleetBuildRepository.Upsert(&galaxy.FleetBuild{ID: "fb-3", DivisionId: "div-1", RaceId: "race-a"})

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
				body := w.Body.String()
				if err := json.Unmarshal([]byte(body), &gotFleetBuilds); err != nil {
					t.Fatalf("test case %s failed to unmarshal response: %v", testCase.name, err)
					return
				}

				if !reflect.DeepEqual(gotFleetBuilds, testCase.expectedFleetBuilds) {
					t.Errorf("test case [%s] expected results: %v got FleetBuilds %v", testCase.name, testCase.expectedFleetBuilds, gotFleetBuilds)

					diff := util.DiffStruct(testCase.expectedFleetBuilds, gotFleetBuilds)
					t.Errorf("test case [%s] diff:\n%v", testCase.name, diff)
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

		{name: "all",
			query:                "",
			authHeader:           "Bearer admin-token",
			expectedResponseCode: http.StatusOK,
			expectedFleetBuilds: []*galaxy.FleetBuild{
				{ID: "fb-1", DivisionId: "div-1", RaceId: "race-1"},
				{ID: "fb-2", DivisionId: "div-1", RaceId: "race-1"},
				{ID: "fb-3", DivisionId: "div-1", RaceId: "race-a"},
			},
		},
		// TODO more test cases
	}
}
