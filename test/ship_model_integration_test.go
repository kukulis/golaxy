package test

import (
	"encoding/json"
	"io"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

func TestShipModelEndpoints(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	baseURL := server.URL + "/api"

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "GET all ship-models - empty",
			method:         "GET",
			path:           "/ship-models",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				if string(body) != "null" {
					t.Errorf("Expected empty array or null, got: %s", string(body))
				}
			},
		},
		{
			name:   "POST create ship-model",
			method: "POST",
			path:   "/ship-models",
			body: map[string]interface{}{
				"id":           "sm1",
				"name":         "Cruiser",
				"guns":         10,
				"one_gun_mass": 5.5,
				"defense_mass": 100.0,
				"engine_mass":  200.0,
				"owner_id":     "player1",
			},
			expectedStatus: 201,
			validateBody: func(t *testing.T, body []byte) {
				var shipModel galaxy.ShipModel
				if err := json.Unmarshal(body, &shipModel); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if shipModel.ID != "sm1" {
					t.Errorf("Expected ID 'sm1', got: %s", shipModel.ID)
				}
				if shipModel.Name != "Cruiser" {
					t.Errorf("Expected Name 'Cruiser', got: %s", shipModel.Name)
				}
				if shipModel.Guns != 10 {
					t.Errorf("Expected Guns 10, got: %d", shipModel.Guns)
				}
			},
		},
		{
			name:           "GET single ship-model",
			method:         "GET",
			path:           "/ship-models/sm1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var shipModel galaxy.ShipModel
				if err := json.Unmarshal(body, &shipModel); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if shipModel.ID != "sm1" {
					t.Errorf("Expected ID 'sm1', got: %s", shipModel.ID)
				}
			},
		},
		{
			name:           "GET all ship-models - with data",
			method:         "GET",
			path:           "/ship-models",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var shipModels []galaxy.ShipModel
				if err := json.Unmarshal(body, &shipModels); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(shipModels) != 1 {
					t.Errorf("Expected 1 ship-model, got: %d", len(shipModels))
				}
			},
		},
		{
			name:   "PUT update ship-model",
			method: "PUT",
			path:   "/ship-models/sm1",
			body: map[string]interface{}{
				"name":         "Destroyer",
				"guns":         15,
				"one_gun_mass": 6.0,
				"defense_mass": 150.0,
				"engine_mass":  250.0,
				"owner_id":     "player2",
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var shipModel galaxy.ShipModel
				if err := json.Unmarshal(body, &shipModel); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if shipModel.Name != "Destroyer" {
					t.Errorf("Expected Name 'Destroyer', got: %s", shipModel.Name)
				}
				if shipModel.Guns != 15 {
					t.Errorf("Expected Guns 15, got: %d", shipModel.Guns)
				}
			},
		},
		{
			name:           "DELETE ship-model",
			method:         "DELETE",
			path:           "/ship-models/sm1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["message"] != "ShipModel deleted successfully" {
					t.Errorf("Unexpected message: %s", response["message"])
				}
			},
		},
		{
			name:           "GET non-existent ship-model - 404",
			method:         "GET",
			path:           "/ship-models/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "ShipModel not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:   "PUT non-existent ship-model - 404",
			method: "PUT",
			path:   "/ship-models/nonexistent",
			body: map[string]interface{}{
				"name": "Test",
			},
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "ShipModel not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:           "DELETE non-existent ship-model - 404",
			method:         "DELETE",
			path:           "/ship-models/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "ShipModel not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makeRequest(tt.method, baseURL+tt.path, tt.body)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got: %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if tt.validateBody != nil {
				tt.validateBody(t, body)
			}
		})
	}
}

func TestShipModelEndpoints_CreateMultiple(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	baseURL := server.URL + "/api"

	// Create multiple ship-models
	shipModels := []map[string]interface{}{
		{"id": "sm1", "name": "Cruiser", "guns": 10, "one_gun_mass": 5.0, "defense_mass": 100.0, "engine_mass": 200.0, "owner_id": "player1"},
		{"id": "sm2", "name": "Destroyer", "guns": 15, "one_gun_mass": 6.0, "defense_mass": 150.0, "engine_mass": 250.0, "owner_id": "player2"},
		{"id": "sm3", "name": "Frigate", "guns": 8, "one_gun_mass": 4.5, "defense_mass": 80.0, "engine_mass": 180.0, "owner_id": "player3"},
	}

	for _, sm := range shipModels {
		resp, err := makeRequest("POST", baseURL+"/ship-models", sm)
		if err != nil {
			t.Fatalf("Failed to create ship-model: %v", err)
		}
		_ = resp.Body.Close()

		if resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got: %d", resp.StatusCode)
		}
	}

	// Get all ship-models
	resp, err := makeRequest("GET", baseURL+"/ship-models", nil)
	if err != nil {
		t.Fatalf("Failed to get ship-models: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var result []galaxy.ShipModel
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 ship-models, got: %d", len(result))
	}

	// Verify they are sorted by ID
	if result[0].ID != "sm1" || result[1].ID != "sm2" || result[2].ID != "sm3" {
		t.Errorf("ShipModels not sorted correctly. Got: %v", result)
	}
}
