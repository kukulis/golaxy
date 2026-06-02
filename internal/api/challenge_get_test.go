package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type getChallengeTestCase struct {
	name           string
	id             string
	authHeader     string
	setupRaces     []*galaxy.Race
	setupChallenge *galaxy.Challenge
	expectedStatus int
}

func getChallengeTestCases() []getChallengeTestCase {
	challenge := &galaxy.Challenge{
		ID:               "challenge-1",
		ChallengerRaceId: "race-1",
		ChallengeeRaceId: "race-2",
		Status:           galaxy.ChallengeStatusPending,
	}
	return []getChallengeTestCase{
		{
			name:           "challenger can get challenge",
			id:             "challenge-1",
			authHeader:     "Bearer token-1",
			setupRaces:     []*galaxy.Race{{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer}},
			setupChallenge: challenge,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "challengee can get challenge",
			id:             "challenge-1",
			authHeader:     "Bearer token-2",
			setupRaces:     []*galaxy.Race{{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer}},
			setupChallenge: challenge,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "admin can get any challenge",
			id:             "challenge-1",
			authHeader:     "Bearer admin-token",
			setupRaces:     []*galaxy.Race{{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin}},
			setupChallenge: challenge,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 403 when race is neither challenger nor challengee",
			id:             "challenge-1",
			authHeader:     "Bearer token-3",
			setupRaces:     []*galaxy.Race{{ID: "race-3", Token: "token-3", Role: galaxy.RolePlayer}},
			setupChallenge: challenge,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns 404 when challenge not found",
			id:             "missing",
			authHeader:     "Bearer token-1",
			setupRaces:     []*galaxy.Race{{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer}},
			setupChallenge: nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 401 without authorization",
			id:             "challenge-1",
			authHeader:     "",
			setupRaces:     []*galaxy.Race{},
			setupChallenge: challenge,
			expectedStatus: http.StatusUnauthorized,
		},
	}
}

func TestGetChallenge(t *testing.T) {
	for _, tc := range getChallengeTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, _, _ := newChallengeController()
			for _, r := range tc.setupRaces {
				raceRepo.Upsert(r)
			}
			if tc.setupChallenge != nil {
				challengeRepo.Upsert(tc.setupChallenge)
			}

			router := setupChallengeRouter(controller)
			req, _ := http.NewRequest(http.MethodGet, "/api/challenges/"+tc.id, nil)
			req.Header.Set("Authorization", tc.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}
