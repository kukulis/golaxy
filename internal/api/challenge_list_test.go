package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type getChallengesTestCase struct {
	name            string
	queryParams     string
	authHeader      string
	setupRaces      []*galaxy.Race
	setupChallenges []*galaxy.Challenge
	expectedStatus  int
	expectedCount   int
}

func getChallengesTestCases() []getChallengesTestCase {
	now := time.Now()
	challengeA := &galaxy.Challenge{ID: "c-1", ChallengerRaceId: "race-1", ChallengeeRaceId: "race-2", ReadyA: true}
	challengeB := &galaxy.Challenge{ID: "c-2", ChallengerRaceId: "race-1", ChallengeeRaceId: "race-3", ReadyB: true}
	challengeC := &galaxy.Challenge{ID: "c-3", ChallengerRaceId: "race-3", ChallengeeRaceId: "race-2", ReadyA: true, ReadyB: true, AcceptedAAt: &now, AcceptedBAt: &now}

	races := []*galaxy.Race{
		{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
		{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
		{ID: "race-3", Token: "token-3", Role: galaxy.RolePlayer},
		{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin},
	}
	allChallenges := []*galaxy.Challenge{challengeA, challengeB, challengeC}

	return []getChallengesTestCase{
		{
			name:            "challenger filters by challenger_id",
			queryParams:     "challenger_id=race-1",
			authHeader:      "Bearer token-1",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   2,
		},
		{
			name:            "challengee filters by challengee_id",
			queryParams:     "challengee_id=race-2",
			authHeader:      "Bearer token-2",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   2,
		},
		{
			name:            "admin gets all challenges without filter",
			queryParams:     "",
			authHeader:      "Bearer admin-token",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   3,
		},
		{
			name:            "filter by ready_a",
			queryParams:     "challenger_id=race-1&ready_a=1",
			authHeader:      "Bearer token-1",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   1,
		},
		{
			name:            "filter by ready_b",
			queryParams:     "challengee_id=race-2&ready_b=1",
			authHeader:      "Bearer token-2",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   1,
		},
		{
			name:            "returns 403 when non-admin provides no filter",
			queryParams:     "",
			authHeader:      "Bearer token-1",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "returns 403 when filtering by another race challenger_id",
			queryParams:     "challenger_id=race-3",
			authHeader:      "Bearer token-1",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "returns 403 when filtering by another race challengee_id",
			queryParams:     "challengee_id=race-3",
			authHeader:      "Bearer token-2",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "returns 401 without authorization",
			queryParams:     "challenger_id=race-1",
			authHeader:      "",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusUnauthorized,
		},
	}
}

func TestGetChallenges(t *testing.T) {
	for _, tc := range getChallengesTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, _, _ := newChallengeController()
			for _, r := range tc.setupRaces {
				raceRepo.Upsert(r)
			}
			for _, c := range tc.setupChallenges {
				challengeRepo.Upsert(c)
			}

			router := setupChallengeRouter(controller)
			url := "/api/challenges"
			if tc.queryParams != "" {
				url += "?" + tc.queryParams
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", tc.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.expectedStatus == http.StatusOK {
				var challenges []*galaxy.Challenge
				if err := json.Unmarshal(w.Body.Bytes(), &challenges); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(challenges) != tc.expectedCount {
					t.Errorf("expected %d challenges, got %d", tc.expectedCount, len(challenges))
				}
			}
		})
	}
}
