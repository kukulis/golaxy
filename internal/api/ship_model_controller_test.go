package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func setupShipModelRouter(controller *ShipModelController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/api")
	apiGroup.GET("/ship-models/:id", controller.GetShipModel)
	apiGroup.PUT("/ship-models/:id", controller.UpdateShipModel)
	return router
}

// --- GetShipModel ---

type getShipModelTestCase struct {
	name              string
	id                string
	storedShipModel   *galaxy.ShipModel
	expectedStatus    int
	expectedShipModel *galaxy.ShipModel
}

func getShipModelTestCases() []getShipModelTestCase {
	return []getShipModelTestCase{
		{
			name: "returns ship model when found",
			id:   "sm-1",
			storedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Name:        "Fighter",
				Guns:        4,
				OneGunMass:  2,
				DefenseMass: 8,
				EngineMass:  16,
				CargoMass:   0,
				OwnerId:     "race-1",
			},
			expectedStatus: http.StatusOK,
			expectedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Name:        "Fighter",
				Guns:        4,
				OneGunMass:  2,
				DefenseMass: 8,
				EngineMass:  16,
				CargoMass:   0,
				OwnerId:     "race-1",
			},
		},
		{
			name:              "returns 404 when not found",
			id:                "missing",
			storedShipModel:   nil,
			expectedStatus:    http.StatusNotFound,
			expectedShipModel: nil,
		},
	}
}

func TestGetShipModel(t *testing.T) {
	for _, tc := range getShipModelTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			id, stored, expectedStatus, expectedModel := tc.id, tc.storedShipModel, tc.expectedStatus, tc.expectedShipModel

			repo := dao.NewShipModelRepository()
			if stored != nil {
				repo.Upsert(stored)
			}

			controller := NewShipModelController(NewMemoryAuthenticationManager(), repo)
			router := setupShipModelRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/ship-models/"+id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != expectedStatus {
				t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
			}

			if expectedModel != nil {
				var got galaxy.ShipModel
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if got != *expectedModel {
					t.Errorf("expected %+v, got %+v", *expectedModel, got)
				}
			}
		})
	}
}

// --- UpdateShipModel ---

type updateShipModelTestCase struct {
	name            string
	id              string
	storedShipModel *galaxy.ShipModel
	body            string
	headers         map[string]string
	expectedStatus    int
	expectedShipModel *galaxy.ShipModel
}

func buildUpdateShipModelTestCases() []updateShipModelTestCase {
	return []updateShipModelTestCase{
		{
			name: "updates ship model",
			id:   "sm-1",
			storedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Name:        "Fighter",
				Guns:        4,
				OneGunMass:  2,
				DefenseMass: 8,
				EngineMass:  16,
				CargoMass:   0,
				OwnerId:     "race-1",
			},
			body:           `{"name":"Heavy Fighter","guns":8,"one_gun_mass":2,"defense_mass":16,"engine_mass":16,"cargo_mass":0,"owner_id":"race-1"}`,
			headers:        map[string]string{"Authorization": "Bearer test-token"},
			expectedStatus: http.StatusOK,
			expectedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Name:        "Heavy Fighter",
				Guns:        8,
				OneGunMass:  2,
				DefenseMass: 16,
				EngineMass:  16,
				CargoMass:   0,
				OwnerId:     "race-1",
			},
		},
		{
			name: "returns 403 without authorization header",
			id:   "sm-1",
			storedShipModel: &galaxy.ShipModel{
				ID:          "sm-1",
				Name:        "Fighter",
				Guns:        4,
				OneGunMass:  2,
				DefenseMass: 8,
				EngineMass:  16,
				CargoMass:   0,
				OwnerId:     "race-1",
			},
			body:              `{"name":"Heavy Fighter","guns":8,"one_gun_mass":2,"defense_mass":16,"engine_mass":16,"cargo_mass":0}`,
			headers:           map[string]string{},
			expectedStatus:    http.StatusForbidden,
			expectedShipModel: nil,
		},
		{
			name:              "returns 404 when not found",
			id:                "missing",
			storedShipModel:   nil,
			body:              `{"name":"Fighter"}`,
			headers:           map[string]string{"Authorization": "Bearer test-token"},
			expectedStatus:    http.StatusNotFound,
			expectedShipModel: nil,
		},
		{
			name:              "returns 400 on invalid body",
			id:                "sm-1",
			storedShipModel:   &galaxy.ShipModel{ID: "sm-1"},
			body:              `not json`,
			headers:           map[string]string{"Authorization": "Bearer test-token"},
			expectedStatus:    http.StatusBadRequest,
			expectedShipModel: nil,
		},
	}
}

func TestUpdateShipModel(t *testing.T) {
	auth := NewMemoryAuthenticationManager()
	auth.AddToken("test-token", &galaxy.Race{ID: "race-1", Name: "Race One"})

	for _, tc := range buildUpdateShipModelTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			id, stored, body, headers, expectedStatus, expectedModel := tc.id, tc.storedShipModel, tc.body, tc.headers, tc.expectedStatus, tc.expectedShipModel

			repo := dao.NewShipModelRepository()
			if stored != nil {
				repo.Upsert(stored)
			}

			controller := NewShipModelController(auth, repo)
			router := setupShipModelRouter(controller)

			req := httptest.NewRequest(http.MethodPut, "/api/ship-models/"+id, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			for k, v := range headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != expectedStatus {
				t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
			}

			if expectedModel != nil {
				var got galaxy.ShipModel
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if got != *expectedModel {
					t.Errorf("expected %+v, got %+v", *expectedModel, got)
				}
			}
		})
	}
}
