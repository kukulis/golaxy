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

func setupChallengeRouter(controller *ChallengeController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/api")
	apiGroup.GET("/challenges", controller.GetAllChallenges)
	apiGroup.GET("/challenges/:id", controller.GetChallenge)
	apiGroup.POST("/challenges", controller.CreateChallenge)
	apiGroup.PUT("/challenges/:id", controller.UpdateChallenge)
	apiGroup.DELETE("/challenges/:id", controller.DeleteChallenge)
	return router
}

func newChallengeTestDeps() (*dao.RaceRepository, AuthenticationManager, *dao.ChallengeRepository) {
	raceRepo := dao.NewRaceRepository()
	raceRepo.Upsert(&galaxy.Race{ID: "race-1", Name: "Race One", Token: "test-token"})
	raceRepo.Upsert(&galaxy.Race{ID: "race-2", Name: "Race Two", Token: "other-token"})
	auth := NewMemoryAuthenticationManager(raceRepo)
	challengeRepo := dao.NewChallengeRepository()
	return raceRepo, auth, challengeRepo
}

// --- GetChallenge ---

func TestGetChallenge(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		stored         *galaxy.Challenge
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "returns challenge when found",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1", ChallengerRaceId: "race-1"},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 404 when not found",
			id:             "missing",
			stored:         nil,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 401 without authorization",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1"},
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, auth, repo := newChallengeTestDeps()
			if tc.stored != nil {
				repo.Upsert(tc.stored)
			}
			controller := NewChallengeController(auth, repo)
			router := setupChallengeRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/challenges/"+tc.id, nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("test case %q: expected status %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
		})
	}
}

// --- GetAllChallenges ---

func TestGetAllChallenges(t *testing.T) {
	tests := []struct {
		name           string
		stored         []*galaxy.Challenge
		authHeader     string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "returns all challenges",
			stored:         []*galaxy.Challenge{{ID: "ch-1"}, {ID: "ch-2"}},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "returns 401 without authorization",
			stored:         nil,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, auth, repo := newChallengeTestDeps()
			for _, c := range tc.stored {
				repo.Upsert(c)
			}
			controller := NewChallengeController(auth, repo)
			router := setupChallengeRouter(controller)

			req := httptest.NewRequest(http.MethodGet, "/api/challenges", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("test case %q: expected status %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
			if tc.expectedCount > 0 {
				var got []*galaxy.Challenge
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("test case %q: failed to unmarshal: %v", tc.name, err)
				}
				if len(got) != tc.expectedCount {
					t.Errorf("test case %q: expected %d challenges, got %d", tc.name, tc.expectedCount, len(got))
				}
			}
		})
	}
}

// --- CreateChallenge ---

func TestCreateChallenge(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		authHeader     string
		expectedStatus int
		expectedRaceId string
	}{
		{
			name:           "creates challenge with challenger_race_id from auth",
			body:           `{"id":"ch-1","fleet_build_a_id":"fb-1","fleet_build_b_id":"fb-2"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusCreated,
			expectedRaceId: "race-1",
		},
		{
			name:           "returns 400 when id is empty",
			body:           `{"id":"","fleet_build_a_id":"fb-1"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "returns 400 when id is missing",
			body:           `{"fleet_build_a_id":"fb-1"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "returns 401 without authorization",
			body:           `{"id":"ch-1"}`,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, auth, repo := newChallengeTestDeps()
			controller := NewChallengeController(auth, repo)
			router := setupChallengeRouter(controller)

			req := httptest.NewRequest(http.MethodPost, "/api/challenges", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("test case %q: expected status %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
			if tc.expectedRaceId != "" {
				var got galaxy.Challenge
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("test case %q: failed to unmarshal: %v", tc.name, err)
				}
				if got.ChallengerRaceId != tc.expectedRaceId {
					t.Errorf("test case %q: expected challenger_race_id %q, got %q", tc.name, tc.expectedRaceId, got.ChallengerRaceId)
				}
			}
		})
	}
}

// --- UpdateChallenge ---

func TestUpdateChallenge(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		stored         *galaxy.Challenge
		body           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "updates challenge",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1", ChallengerRaceId: "race-1", FleetBuildAId: "fb-1"},
			body:           `{"fleet_build_a_id":"fb-2","fleet_build_b_id":"fb-3"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 404 when not found",
			id:             "missing",
			stored:         nil,
			body:           `{"fleet_build_a_id":"fb-2"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 403 when not owner",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1", ChallengerRaceId: "race-2"},
			body:           `{"fleet_build_a_id":"fb-2"}`,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns 401 without authorization",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1"},
			body:           `{"fleet_build_a_id":"fb-2"}`,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, auth, repo := newChallengeTestDeps()
			if tc.stored != nil {
				repo.Upsert(tc.stored)
			}
			controller := NewChallengeController(auth, repo)
			router := setupChallengeRouter(controller)

			req := httptest.NewRequest(http.MethodPut, "/api/challenges/"+tc.id, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("test case %q: expected status %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
		})
	}
}

// --- DeleteChallenge ---

func TestDeleteChallenge(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		stored         *galaxy.Challenge
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "deletes challenge",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1", ChallengerRaceId: "race-1"},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 404 when not found",
			id:             "missing",
			stored:         nil,
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 403 when not owner",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1", ChallengerRaceId: "race-2"},
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns 401 without authorization",
			id:             "ch-1",
			stored:         &galaxy.Challenge{ID: "ch-1"},
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, auth, repo := newChallengeTestDeps()
			if tc.stored != nil {
				repo.Upsert(tc.stored)
			}
			controller := NewChallengeController(auth, repo)
			router := setupChallengeRouter(controller)

			req := httptest.NewRequest(http.MethodDelete, "/api/challenges/"+tc.id, nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("test case %q: expected status %d, got %d", tc.name, tc.expectedStatus, w.Code)
			}
		})
	}
}
