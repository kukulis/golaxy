package test

import (
	"encoding/json"
	"io"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

func TestFleetBuildEndpoints(t *testing.T) {
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
			name:           "GET all fleet-builds - empty",
			method:         "GET",
			path:           "/fleet-builds",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				if string(body) != "null" {
					t.Errorf("Expected empty array or null, got: %s", string(body))
				}
			},
		},
		{
			name:   "POST create fleet-build",
			method: "POST",
			path:   "/fleet-builds",
			body: map[string]interface{}{
				"id":                "fb1",
				"division_id":       "div1",
				"race_id":           "race1",
				"attack_resources":  100.5,
				"defense_resources": 200.0,
				"engine_resources":  150.0,
				"cargo_resources":   50.5,
			},
			expectedStatus: 201,
			validateBody: func(t *testing.T, body []byte) {
				var fleetBuild galaxy.FleetBuild
				if err := json.Unmarshal(body, &fleetBuild); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if fleetBuild.ID != "fb1" {
					t.Errorf("Expected ID 'fb1', got: %s", fleetBuild.ID)
				}
				if fleetBuild.DivisionId != "div1" {
					t.Errorf("Expected DivisionId 'div1', got: %s", fleetBuild.DivisionId)
				}
				if fleetBuild.AttackResources != 100.5 {
					t.Errorf("Expected AttackResources 100.5, got: %f", fleetBuild.AttackResources)
				}
			},
		},
		{
			name:           "GET single fleet-build",
			method:         "GET",
			path:           "/fleet-builds/fb1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var fleetBuild galaxy.FleetBuild
				if err := json.Unmarshal(body, &fleetBuild); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if fleetBuild.ID != "fb1" {
					t.Errorf("Expected ID 'fb1', got: %s", fleetBuild.ID)
				}
			},
		},
		{
			name:           "GET all fleet-builds - with data",
			method:         "GET",
			path:           "/fleet-builds",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var fleetBuilds []galaxy.FleetBuild
				if err := json.Unmarshal(body, &fleetBuilds); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(fleetBuilds) != 1 {
					t.Errorf("Expected 1 fleet-build, got: %d", len(fleetBuilds))
				}
			},
		},
		{
			name:   "PUT update fleet-build",
			method: "PUT",
			path:   "/fleet-builds/fb1",
			body: map[string]interface{}{
				"division_id":       "div2",
				"race_id":           "race2",
				"attack_resources":  200.0,
				"defense_resources": 300.0,
				"engine_resources":  250.0,
				"cargo_resources":   100.0,
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var fleetBuild galaxy.FleetBuild
				if err := json.Unmarshal(body, &fleetBuild); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if fleetBuild.AttackResources != 200.0 {
					t.Errorf("Expected AttackResources 200.0, got: %f", fleetBuild.AttackResources)
				}
				if fleetBuild.DivisionId != "div2" {
					t.Errorf("Expected DivisionId 'div2', got: %s", fleetBuild.DivisionId)
				}
			},
		},
		{
			name:           "DELETE fleet-build",
			method:         "DELETE",
			path:           "/fleet-builds/fb1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["message"] != "FleetBuild deleted successfully" {
					t.Errorf("Unexpected message: %s", response["message"])
				}
			},
		},
		{
			name:           "GET non-existent fleet-build - 404",
			method:         "GET",
			path:           "/fleet-builds/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "FleetBuild not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:   "PUT non-existent fleet-build - 404",
			method: "PUT",
			path:   "/fleet-builds/nonexistent",
			body: map[string]interface{}{
				"attack_resources": 100.0,
			},
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "FleetBuild not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:           "DELETE non-existent fleet-build - 404",
			method:         "DELETE",
			path:           "/fleet-builds/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "FleetBuild not found" {
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

func TestFleetBuildEndpoints_CreateMultiple(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	baseURL := server.URL + "/api"

	// Create multiple fleet-builds
	fleetBuilds := []map[string]interface{}{
		{"id": "fb1", "division_id": "div1", "race_id": "race1", "attack_resources": 100.0, "defense_resources": 200.0, "engine_resources": 150.0, "cargo_resources": 50.0},
		{"id": "fb2", "division_id": "div2", "race_id": "race2", "attack_resources": 150.0, "defense_resources": 250.0, "engine_resources": 200.0, "cargo_resources": 75.0},
		{"id": "fb3", "division_id": "div3", "race_id": "race3", "attack_resources": 120.0, "defense_resources": 220.0, "engine_resources": 180.0, "cargo_resources": 60.0},
	}

	for _, fb := range fleetBuilds {
		resp, err := makeRequest("POST", baseURL+"/fleet-builds", fb)
		if err != nil {
			t.Fatalf("Failed to create fleet-build: %v", err)
		}
		_ = resp.Body.Close()

		if resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got: %d", resp.StatusCode)
		}
	}

	// Get all fleet-builds
	resp, err := makeRequest("GET", baseURL+"/fleet-builds", nil)
	if err != nil {
		t.Fatalf("Failed to get fleet-builds: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var result []galaxy.FleetBuild
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 fleet-builds, got: %d", len(result))
	}

	// Verify they are sorted by ID
	if result[0].ID != "fb1" || result[1].ID != "fb2" || result[2].ID != "fb3" {
		t.Errorf("FleetBuilds not sorted correctly. Got: %v", result)
	}
}
func TestFleetBuildShipModelAssignment(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	baseURL := server.URL + "/api"

	// First test: verify 404 when fleet-build doesn't exist
	resp, _ := makeRequest("GET", baseURL+"/fleet-builds/fb1/ship-models", nil)
	_ = resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404 for non-existent fleet-build, got: %d", resp.StatusCode)
	}

	// Create a fleet-build for subsequent tests
	fleetBuild := map[string]interface{}{
		"id":                "fb1",
		"division_id":       "div1",
		"race_id":           "race1",
		"attack_resources":  100.0,
		"defense_resources": 200.0,
		"engine_resources":  150.0,
		"cargo_resources":   50.0,
	}
	resp, _ = makeRequest("POST", baseURL+"/fleet-builds", fleetBuild)
	_ = resp.Body.Close()

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "GET assigned ship-models - empty",
			method:         "GET",
			path:           "/fleet-builds/fb1/ship-models",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				if string(body) != "null" && string(body) != "[]" {
					t.Errorf("Expected empty array or null, got: %s", string(body))
				}
			},
		},
		{
			name:   "POST assign ship-model",
			method: "POST",
			path:   "/fleet-builds/fb1/ship-models",
			body: map[string]interface{}{
				"ship_model_id": "sm1",
				"amount":        5,
				"result_mass":   500.0,
			},
			expectedStatus: 201,
			validateBody: func(t *testing.T, body []byte) {
				var assignment galaxy.FleetBuildToShipModel
				if err := json.Unmarshal(body, &assignment); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if assignment.FleetBuildID != "fb1" {
					t.Errorf("Expected FleetBuildID 'fb1', got: %s", assignment.FleetBuildID)
				}
				if assignment.ShipModelID != "sm1" {
					t.Errorf("Expected ShipModelID 'sm1', got: %s", assignment.ShipModelID)
				}
				if assignment.Amount != 5 {
					t.Errorf("Expected Amount 5, got: %d", assignment.Amount)
				}
			},
		},
		{
			name:           "GET assigned ship-models - with data",
			method:         "GET",
			path:           "/fleet-builds/fb1/ship-models",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var assignments []galaxy.FleetBuildToShipModel
				if err := json.Unmarshal(body, &assignments); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(assignments) != 1 {
					t.Errorf("Expected 1 assignment, got: %d", len(assignments))
				}
				if len(assignments) > 0 && assignments[0].ShipModelID != "sm1" {
					t.Errorf("Expected ShipModelID 'sm1', got: %s", assignments[0].ShipModelID)
				}
			},
		},
		{
			name:   "POST assign same ship-model again - update (upsert)",
			method: "POST",
			path:   "/fleet-builds/fb1/ship-models",
			body: map[string]interface{}{
				"ship_model_id": "sm1",
				"amount":        10,
				"result_mass":   1000.0,
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var assignment galaxy.FleetBuildToShipModel
				if err := json.Unmarshal(body, &assignment); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if assignment.Amount != 10 {
					t.Errorf("Expected updated Amount 10, got: %d", assignment.Amount)
				}
				if assignment.ResultMass != 1000.0 {
					t.Errorf("Expected updated ResultMass 1000.0, got: %f", assignment.ResultMass)
				}
			},
		},
		{
			name:           "GET assigned ship-models - verify update persisted",
			method:         "GET",
			path:           "/fleet-builds/fb1/ship-models",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var assignments []galaxy.FleetBuildToShipModel
				if err := json.Unmarshal(body, &assignments); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(assignments) != 1 {
					t.Errorf("Expected 1 assignment, got: %d", len(assignments))
				}
				if len(assignments) > 0 {
					if assignments[0].Amount != 10 {
						t.Errorf("Expected persisted Amount 10, got: %d", assignments[0].Amount)
					}
					if assignments[0].ResultMass != 1000.0 {
						t.Errorf("Expected persisted ResultMass 1000.0, got: %f", assignments[0].ResultMass)
					}
				}
			},
		},
		{
			name:           "DELETE unassign ship-model",
			method:         "DELETE",
			path:           "/fleet-builds/fb1/ship-models/sm1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["message"] != "ShipModel unassigned successfully" {
					t.Errorf("Unexpected message: %s", response["message"])
				}
			},
		},
		{
			name:           "DELETE unassign non-existent assignment - 404",
			method:         "DELETE",
			path:           "/fleet-builds/fb1/ship-models/sm999",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "ShipModel assignment not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:   "POST assign to non-existent fleet-build - 404",
			method: "POST",
			path:   "/fleet-builds/nonexistent/ship-models",
			body: map[string]interface{}{
				"ship_model_id": "sm1",
				"amount":        5,
			},
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "FleetBuild not found" {
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
