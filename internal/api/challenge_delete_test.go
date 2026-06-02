package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type deleteChallengeTestCase struct {
	name           string
	id             string
	authHeader     string
	setupRaces     []*galaxy.Race
	setupChallenge *galaxy.Challenge
	expectedStatus int
}

func deleteChallengeTestCases() []deleteChallengeTestCase {
	races := []*galaxy.Race{
		{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
		{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
		{ID: "race-3", Token: "token-3", Role: galaxy.RolePlayer},
		{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin},
	}
	return []deleteChallengeTestCase{
		{
			name:           "challenger deletes challenge",
			id:             "c-1",
			authHeader:     "Bearer token-1",
			setupRaces:     races,
			setupChallenge: freshChallenge(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "admin deletes challenge",
			id:             "c-1",
			authHeader:     "Bearer admin-token",
			setupRaces:     races,
			setupChallenge: freshChallenge(),
			expectedStatus: http.StatusOK,
		},
		{
			name:       "admin deletes executed challenge",
			id:         "c-1",
			authHeader: "Bearer admin-token",
			setupRaces: races,
			setupChallenge: &galaxy.Challenge{
				ID:               "c-1",
				DivisionId:       "div-1",
				ChallengerRaceId: "race-1",
				ChallengeeRaceId: "race-2",
				Status:           galaxy.ChallengeStatusExecuted,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 403 when challengee tries to delete",
			id:             "c-1",
			authHeader:     "Bearer token-2",
			setupRaces:     races,
			setupChallenge: freshChallenge(),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "returns 403 when unrelated race tries to delete",
			id:             "c-1",
			authHeader:     "Bearer token-3",
			setupRaces:     races,
			setupChallenge: freshChallenge(),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:       "returns 409 when challenger tries to delete executed challenge",
			id:         "c-1",
			authHeader: "Bearer token-1",
			setupRaces: races,
			setupChallenge: &galaxy.Challenge{
				ID:               "c-1",
				DivisionId:       "div-1",
				ChallengerRaceId: "race-1",
				ChallengeeRaceId: "race-2",
				Status:           galaxy.ChallengeStatusExecuted,
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "returns 404 when challenge not found",
			id:             "missing",
			authHeader:     "Bearer token-1",
			setupRaces:     races,
			setupChallenge: nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 401 without authorization",
			id:             "c-1",
			authHeader:     "",
			setupRaces:     races,
			setupChallenge: freshChallenge(),
			expectedStatus: http.StatusUnauthorized,
		},
	}
}

func TestDeleteChallenge(t *testing.T) {
	for _, tc := range deleteChallengeTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, _, _ := newChallengeController()
			for _, r := range tc.setupRaces {
				raceRepo.Upsert(r)
			}
			if tc.setupChallenge != nil {
				challengeRepo.Upsert(tc.setupChallenge)
			}

			router := setupChallengeRouter(controller)
			req, _ := http.NewRequest(http.MethodDelete, "/api/challenges/"+tc.id, nil)
			req.Header.Set("Authorization", tc.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.expectedStatus == http.StatusOK {
				if challengeRepo.Get(tc.id) != nil {
					t.Error("expected challenge to be deleted from repository")
				}
			}
		})
	}
}
