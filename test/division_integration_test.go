package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/di"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func setupTestServer() *httptest.Server {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	apiRoute := router.Group("/api")
	di.CreateSingletons("test")
	di.RegisterRoutes(apiRoute)

	return httptest.NewServer(router)
}

func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	return client.Do(req)
}

func TestDivisionEndpoints(t *testing.T) {
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
			name:           "GET all divisions - empty",
			method:         "GET",
			path:           "/divisions",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				if string(body) != "null" {
					t.Errorf("Expected empty array or null, got: %s", string(body))
				}
			},
		},
		{
			name:   "POST create division",
			method: "POST",
			path:   "/divisions",
			body: map[string]interface{}{
				"id":               "div1",
				"resources_amount": 1000,
				"tech_attack":      5,
				"tech_defense":     3,
				"tech_engines":     4,
				"tech_cargo":       2,
			},
			expectedStatus: 201,
			validateBody: func(t *testing.T, body []byte) {
				var division galaxy.Division
				if err := json.Unmarshal(body, &division); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if division.ID != "div1" {
					t.Errorf("Expected ID 'div1', got: %s", division.ID)
				}
				if division.ResourcesAmount != 1000 {
					t.Errorf("Expected ResourcesAmount 1000, got: %d", division.ResourcesAmount)
				}
			},
		},
		{
			name:           "GET single division",
			method:         "GET",
			path:           "/divisions/div1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var division galaxy.Division
				if err := json.Unmarshal(body, &division); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if division.ID != "div1" {
					t.Errorf("Expected ID 'div1', got: %s", division.ID)
				}
			},
		},
		{
			name:           "GET all divisions - with data",
			method:         "GET",
			path:           "/divisions",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var divisions []galaxy.Division
				if err := json.Unmarshal(body, &divisions); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(divisions) != 1 {
					t.Errorf("Expected 1 division, got: %d", len(divisions))
				}
			},
		},
		{
			name:   "PUT update division",
			method: "PUT",
			path:   "/divisions/div1",
			body: map[string]interface{}{
				"resources_amount": 2000,
				"tech_attack":      10,
				"tech_defense":     8,
				"tech_engines":     6,
				"tech_cargo":       4,
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var division galaxy.Division
				if err := json.Unmarshal(body, &division); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if division.ResourcesAmount != 2000 {
					t.Errorf("Expected ResourcesAmount 2000, got: %d", division.ResourcesAmount)
				}
				if division.TechAttack != 10 {
					t.Errorf("Expected TechAttack 10, got: %d", division.TechAttack)
				}
			},
		},
		{
			name:           "DELETE division",
			method:         "DELETE",
			path:           "/divisions/div1",
			body:           nil,
			expectedStatus: 200,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["message"] != "Division deleted successfully" {
					t.Errorf("Unexpected message: %s", response["message"])
				}
			},
		},
		{
			name:           "GET non-existent division - 404",
			method:         "GET",
			path:           "/divisions/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "Division not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:   "PUT non-existent division - 404",
			method: "PUT",
			path:   "/divisions/nonexistent",
			body: map[string]interface{}{
				"resources_amount": 100,
			},
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "Division not found" {
					t.Errorf("Unexpected error message: %s", response["error"])
				}
			},
		},
		{
			name:           "DELETE non-existent division - 404",
			method:         "DELETE",
			path:           "/divisions/nonexistent",
			body:           nil,
			expectedStatus: 404,
			validateBody: func(t *testing.T, body []byte) {
				var response map[string]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if response["error"] != "Division not found" {
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

func TestDivisionEndpoints_CreateMultiple(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	baseURL := server.URL + "/api"

	// Create multiple divisions
	divisions := []map[string]interface{}{
		{"id": "div1", "resources_amount": 1000, "tech_attack": 5, "tech_defense": 3, "tech_engines": 4, "tech_cargo": 2},
		{"id": "div2", "resources_amount": 2000, "tech_attack": 7, "tech_defense": 5, "tech_engines": 6, "tech_cargo": 3},
		{"id": "div3", "resources_amount": 1500, "tech_attack": 6, "tech_defense": 4, "tech_engines": 5, "tech_cargo": 2},
	}

	for _, div := range divisions {
		resp, err := makeRequest("POST", baseURL+"/divisions", div)
		if err != nil {
			t.Fatalf("Failed to create division: %v", err)
		}
		_ = resp.Body.Close()

		if resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got: %d", resp.StatusCode)
		}
	}

	// Get all divisions
	resp, err := makeRequest("GET", baseURL+"/divisions", nil)
	if err != nil {
		t.Fatalf("Failed to get divisions: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	var result []galaxy.Division
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 divisions, got: %d", len(result))
	}

	// Verify they are sorted by ID
	if result[0].ID != "div1" || result[1].ID != "div2" || result[2].ID != "div3" {
		t.Errorf("Divisions not sorted correctly. Got: %v", result)
	}
}
