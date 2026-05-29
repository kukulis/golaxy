package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func setupFleetBuildRouter(controller *FleetBuildController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/api")
	apiGroup.GET("/fleet-builds/:id", controller.GetFleetBuild)
	apiGroup.GET("/fleet-builds/:id/statistics", controller.GetStatistics)
	return router
}

type getFleetBuildTestCase struct {
	name               string
	id                 string
	storedFleetBuild   *galaxy.FleetBuild
	authHeader         string
	expectedStatus     int
	expectedFleetBuild *galaxy.FleetBuild
}

func getFleetBuildTestCases() []getFleetBuildTestCase {
	return []getFleetBuildTestCase{
		{
			name: "returns fleet build when found",
			id:   "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{
				ID:               "fb-1",
				DivisionId:       "div-1",
				RaceId:           "race-1",
				AttackResources:  10,
				DefenseResources: 20,
				EngineResources:  30,
				CargoResources:   40,
			},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
			expectedFleetBuild: &galaxy.FleetBuild{
				ID:               "fb-1",
				DivisionId:       "div-1",
				RaceId:           "race-1",
				AttackResources:  10,
				DefenseResources: 20,
				EngineResources:  30,
				CargoResources:   40,
				UsedResources:    100,
			},
		},
		{
			name:               "returns 404 when not found",
			id:                 "missing",
			storedFleetBuild:   nil,
			authHeader:         "Bearer test-token",
			expectedStatus:     http.StatusNotFound,
			expectedFleetBuild: nil,
		},
		{
			name:               "returns 401 without authorization",
			id:                 "fb-1",
			storedFleetBuild:   &galaxy.FleetBuild{ID: "fb-1"},
			authHeader:         "",
			expectedStatus:     http.StatusUnauthorized,
			expectedFleetBuild: nil,
		},
	}
}

func TestGetFleetBuild(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token"})
	auth := NewMemoryAuthenticationManager(raceRepo)

	for _, tc := range getFleetBuildTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			repo := dao.NewFleetBuildRepository()
			if tc.storedFleetBuild != nil {
				repo.Upsert(tc.storedFleetBuild)
			}

			controller := NewFleetBuildController(auth, repo, nil, nil, nil)
			router := setupFleetBuildRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/fleet-builds/"+tc.id, nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedFleetBuild != nil {
				var got galaxy.FleetBuild
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if got.ID != tc.expectedFleetBuild.ID ||
					got.DivisionId != tc.expectedFleetBuild.DivisionId ||
					got.RaceId != tc.expectedFleetBuild.RaceId ||
					got.AttackResources != tc.expectedFleetBuild.AttackResources ||
					got.DefenseResources != tc.expectedFleetBuild.DefenseResources ||
					got.EngineResources != tc.expectedFleetBuild.EngineResources ||
					got.CargoResources != tc.expectedFleetBuild.CargoResources {
					t.Errorf("expected %+v, got %+v", tc.expectedFleetBuild, got)
				}

				// TODO later
				//if got.UsedResources != tc.expectedFleetBuild.UsedResources {
				//	t.Errorf("expected fleet build used resources %+v, got %+v", tc.expectedFleetBuild.UsedResources, got.UsedResources)
				//}
			}
		})
	}
}

// --- GetStatistics ---

type getStatisticsTestCase struct {
	name               string
	fleetBuildId       string
	storedFleetBuild   *galaxy.FleetBuild
	storedDivision     *galaxy.Division
	authHeader         string
	expectedStatus     int
	expectedStatistics *galaxy.FleetBuildStatistics
}

func getStatisticsTestCases() []getStatisticsTestCase {
	return []getStatisticsTestCase{
		{
			name:         "returns statistics when fleet build and division found",
			fleetBuildId: "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{
				ID:               "fb-1",
				DivisionId:       "div-1",
				RaceId:           "race-1",
				AttackResources:  10,
				DefenseResources: 20,
				EngineResources:  30,
				CargoResources:   40,
			},
			storedDivision: &galaxy.Division{
				ID:              "div-1",
				ResourcesAmount: 200,
			},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
			expectedStatistics: &galaxy.FleetBuildStatistics{
				MaxResources:                 200,
				UsedResources:                100,
				UsedResourcesForShips:        0,
				UsedResourcesForTechnologies: 100,
				RemainingResources:           100,
				ExceedingResources:           0,
			},
		},
		{
			name:             "returns 404 when fleet build not found",
			fleetBuildId:     "missing",
			storedFleetBuild: nil,
			storedDivision:   &galaxy.Division{ID: "div-1", ResourcesAmount: 200},
			authHeader:       "Bearer test-token",
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:         "returns 404 when division not found",
			fleetBuildId: "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{
				ID:         "fb-1",
				DivisionId: "missing-div",
			},
			storedDivision: nil,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "returns 401 without authorization",
			fleetBuildId: "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{
				ID:         "fb-1",
				DivisionId: "div-1",
			},
			storedDivision: &galaxy.Division{ID: "div-1", ResourcesAmount: 200},
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
}

func TestGetStatistics(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token"})
	auth := NewMemoryAuthenticationManager(raceRepo)

	for _, tc := range getStatisticsTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			fleetBuildRepo := dao.NewFleetBuildRepository()
			if tc.storedFleetBuild != nil {
				fleetBuildRepo.Upsert(tc.storedFleetBuild)
			}

			divisionRepo := dao.NewDivisionRepository()
			if tc.storedDivision != nil {
				divisionRepo.Upsert(tc.storedDivision)
			}

			controller := NewFleetBuildController(auth, fleetBuildRepo, nil, dao.NewShipModelRepository(), divisionRepo)
			router := setupFleetBuildRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/fleet-builds/"+tc.fleetBuildId+"/statistics", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedStatistics != nil {
				var got galaxy.FleetBuildStatistics
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if got.MaxResources != tc.expectedStatistics.MaxResources ||
					got.UsedResources != tc.expectedStatistics.UsedResources ||
					got.UsedResourcesForShips != tc.expectedStatistics.UsedResourcesForShips ||
					got.UsedResourcesForTechnologies != tc.expectedStatistics.UsedResourcesForTechnologies ||
					got.RemainingResources != tc.expectedStatistics.RemainingResources ||
					got.ExceedingResources != tc.expectedStatistics.ExceedingResources {
					t.Errorf("expected %+v, got %+v", tc.expectedStatistics, got)
				}
			}
		})
	}
}
