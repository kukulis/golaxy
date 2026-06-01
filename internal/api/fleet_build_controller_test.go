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
	apiGroup.GET("/fleet-builds/:id/ship-models", controller.GetAssignedShipModels)
	apiGroup.GET("/ship-models-assignment/:id/calculate-ship-tech", controller.CalculateShipTech)
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

// --- GetAssignedShipModels ---

type getAssignedShipModelsTestCase struct {
	name             string
	fleetBuildId     string
	storedFleetBuild *galaxy.FleetBuild
	assignments      []*galaxy.ShipModelAssignment
	authHeader       string
	storedShipModels []*galaxy.ShipModel

	// expected part
	expectedStatus      int
	expectedAssignments []*galaxy.ShipModelAssignment
}

func getAssignedShipModelsTestCases() []getAssignedShipModelsTestCase {
	return []getAssignedShipModelsTestCase{
		{
			name:             "returns assigned ship models",
			fleetBuildId:     "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{ID: "fb-1"},
			assignments: []*galaxy.ShipModelAssignment{
				{FleetBuildID: "fb-1", ShipModelID: "sm-1", Amount: 3},
				{FleetBuildID: "fb-1", ShipModelID: "sm-2", Amount: 5},
			},
			authHeader: "Bearer test-token",

			storedShipModels: []*galaxy.ShipModel{
				{
					ID:          "sm-1",
					Name:        "First Ship",
					Guns:        1,
					OneGunMass:  1,
					DefenseMass: 0,
					EngineMass:  0,
					CargoMass:   0,
					OwnerId:     "alpha",
				},
				{
					ID:          "sm-2",
					Name:        "Second Ship",
					Guns:        1,
					OneGunMass:  1,
					DefenseMass: 1,
					EngineMass:  0,
					CargoMass:   0,
					OwnerId:     "alpha",
				},
			},

			expectedStatus: http.StatusOK,
			expectedAssignments: []*galaxy.ShipModelAssignment{
				{FleetBuildID: "fb-1", ShipModelID: "sm-1", Amount: 3,
					ShipModel: &galaxy.ShipModel{
						ID:   "sm-1",
						Name: "First Ship",
						//Guns:        0,
						//OneGunMass:  0,
						//DefenseMass: 0,
						//EngineMass:  0,
						//CargoMass:   0,
						//OwnerId:     "",
					}},
				{FleetBuildID: "fb-1", ShipModelID: "sm-2", Amount: 5,
					ShipModel: &galaxy.ShipModel{
						ID:   "sm-2",
						Name: "Second Ship",
						//Guns:        0,
						//OneGunMass:  0,
						//DefenseMass: 0,
						//EngineMass:  0,
						//CargoMass:   0,
						//OwnerId:     "",
					}},
			},
		},
		{
			name:                "returns empty list when no assignments",
			fleetBuildId:        "fb-1",
			storedFleetBuild:    &galaxy.FleetBuild{ID: "fb-1"},
			assignments:         nil,
			authHeader:          "Bearer test-token",
			expectedStatus:      http.StatusOK,
			expectedAssignments: []*galaxy.ShipModelAssignment{},
		},
		{
			name:             "returns 404 when fleet build not found",
			fleetBuildId:     "missing",
			storedFleetBuild: nil,
			authHeader:       "Bearer test-token",
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:             "returns 401 without authorization",
			fleetBuildId:     "fb-1",
			storedFleetBuild: &galaxy.FleetBuild{ID: "fb-1"},
			authHeader:       "",
			expectedStatus:   http.StatusUnauthorized,
		},
	}
}

func TestGetAssignedShipModels(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token"})
	auth := NewMemoryAuthenticationManager(raceRepo)

	for _, tc := range getAssignedShipModelsTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			shipModelRepository := dao.NewShipModelRepository()
			for _, shipModel := range tc.storedShipModels {
				shipModelRepository.Upsert(shipModel)
			}
			fleetBuildRepository := dao.NewFleetBuildRepository()
			if tc.storedFleetBuild != nil {
				fleetBuildRepository.Upsert(tc.storedFleetBuild)
			}
			for _, a := range tc.assignments {
				fleetBuildRepository.AssignShipModel(a)
			}

			controller := NewFleetBuildController(auth, fleetBuildRepository, nil, shipModelRepository, nil)
			router := setupFleetBuildRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/fleet-builds/"+tc.fleetBuildId+"/ship-models", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedAssignments != nil {
				var got []*galaxy.ShipModelAssignment
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(got) != len(tc.expectedAssignments) {
					t.Fatalf("expected %d assignments, got %d", len(tc.expectedAssignments), len(got))
				}
				for i, expected := range tc.expectedAssignments {
					if got[i].FleetBuildID != expected.FleetBuildID ||
						got[i].ShipModelID != expected.ShipModelID ||
						got[i].Amount != expected.Amount {
						t.Errorf("assignment %d: expected %+v, got %+v", i, expected, got[i])
					}

					// also asserting the ship model
					if got[i].ShipModel == nil && expected.ShipModel != nil {
						t.Fatalf("assignment %d: expected a ship model, but got nil", i)
					}

					if got[i].ShipModel.ID != expected.ShipModel.ID ||
						got[i].ShipModel.Name != expected.ShipModel.Name {
						t.Errorf("Assigned ship model %d: expected %+v, got %+v", i, expected.ShipModel, got[i].ShipModel)
					}
				}
			}
		})
	}
}

// --- CalculateShipTech ---

type calculateShipTechTestCase struct {
	name             string
	assignmentId     string
	assignment       *galaxy.ShipModelAssignment
	storedFleetBuild *galaxy.FleetBuild
	storedShipModel  *galaxy.ShipModel
	authHeader       string
	expectedStatus   int
	expectedShipTech *galaxy.ShipTech
}

func calculateShipTechTestCases() []calculateShipTechTestCase {
	return []calculateShipTechTestCase{
		{
			name:         "returns ship tech for assignment",
			assignmentId: "asgn-1",
			assignment: &galaxy.ShipModelAssignment{
				ID:           "asgn-1",
				FleetBuildID: "fb-1",
				ShipModelID:  "sm-1",
				Amount:       2,
			},
			storedFleetBuild: &galaxy.FleetBuild{
				ID:               "fb-1",
				RaceId:           "race-1",
				AttackResources:  0,
				DefenseResources: 0,
				EngineResources:  0,
				CargoResources:   0,
			},
			storedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Guns:        0,
				OneGunMass:  0,
				DefenseMass: 0,
				EngineMass:  1,
				CargoMass:   0,
			},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
			expectedShipTech: &galaxy.ShipTech{
				Guns:    0,
				Attack:  0,
				Defense: 0,
				Speed:   1,
				Mass:    1,
			},
		},
		{
			name:             "returns 404 when assignment not found",
			assignmentId:     "missing",
			assignment:       nil,
			storedFleetBuild: nil,
			storedShipModel:  nil,
			authHeader:       "Bearer test-token",
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:         "returns 404 when fleet build not found",
			assignmentId: "asgn-1",
			assignment: &galaxy.ShipModelAssignment{
				ID:           "asgn-1",
				FleetBuildID: "missing-fb",
				ShipModelID:  "sm-1",
			},
			storedFleetBuild: nil,
			storedShipModel:  nil,
			authHeader:       "Bearer test-token",
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:         "returns 404 when ship model not found",
			assignmentId: "asgn-1",
			assignment: &galaxy.ShipModelAssignment{
				ID:           "asgn-1",
				FleetBuildID: "fb-1",
				ShipModelID:  "missing-sm",
			},
			storedFleetBuild: &galaxy.FleetBuild{ID: "fb-1"},
			storedShipModel:  nil,
			authHeader:       "Bearer test-token",
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:         "returns 401 without authorization",
			assignmentId: "asgn-1",
			assignment: &galaxy.ShipModelAssignment{
				ID:           "asgn-1",
				FleetBuildID: "fb-1",
				ShipModelID:  "sm-1",
			},
			storedFleetBuild: &galaxy.FleetBuild{ID: "fb-1"},
			storedShipModel:  &galaxy.ShipModel{ID: "sm-1"},
			authHeader:       "",
			expectedStatus:   http.StatusUnauthorized,
		},
	}
}

func TestCalculateShipTech(t *testing.T) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token"})
	auth := NewMemoryAuthenticationManager(raceRepo)

	for _, tc := range calculateShipTechTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			fleetBuildRepo := dao.NewFleetBuildRepository()
			if tc.storedFleetBuild != nil {
				fleetBuildRepo.Upsert(tc.storedFleetBuild)
			}
			if tc.assignment != nil {
				fleetBuildRepo.AssignShipModel(tc.assignment)
			}

			shipModelRepo := dao.NewShipModelRepository()
			if tc.storedShipModel != nil {
				shipModelRepo.Upsert(tc.storedShipModel)
			}

			controller := NewFleetBuildController(auth, fleetBuildRepo, nil, shipModelRepo, nil)
			router := setupFleetBuildRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/ship-models-assignment/"+tc.assignmentId+"/calculate-ship-tech", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedShipTech != nil {
				var got galaxy.ShipTech
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if got != *tc.expectedShipTech {
					t.Errorf("expected %+v, got %+v", *tc.expectedShipTech, got)
				}
			}
		})
	}
}
