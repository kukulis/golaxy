package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type createChallengeTestCase struct {
	name           string
	queryParams    string
	authHeader     string
	setupRaces     []*galaxy.Race
	setupDivisions []*galaxy.Division
	setupChallenge *galaxy.Challenge
	expectedStatus int
}

func createChallengeTestCases() []createChallengeTestCase {
	return []createChallengeTestCase{
		{
			name:        "creates challenge successfully",
			queryParams: "id=11111111-1111-1111-1111-111111111111&challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
				{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "admin creates challenge with any challenger race",
			queryParams: "id=22222222-2222-2222-2222-222222222222&challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer admin-token",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
				{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
				{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "returns 401 without authorization",
			queryParams:    "id=33333333-3333-3333-3333-333333333333&challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:     "",
			setupRaces:     []*galaxy.Race{},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "returns 403 when challenger_race_id does not match logged-in race",
			queryParams: "id=44444444-4444-4444-4444-444444444444&challenger_race_id=race-other&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
				{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "returns 400 when id is missing",
			queryParams: "challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when challenger_race_id is missing",
			queryParams: "id=55555555-5555-5555-5555-555555555555&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when challengee_race_id is missing",
			queryParams: "id=66666666-6666-6666-6666-666666666666&challenger_race_id=race-1&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when division_id is missing",
			queryParams: "id=77777777-7777-7777-7777-777777777777&challenger_race_id=race-1&challengee_race_id=race-2",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when id is not a valid UUID",
			queryParams: "id=not-a-uuid&challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when division not found",
			queryParams: "id=88888888-8888-8888-8888-888888888888&challenger_race_id=race-1&challengee_race_id=race-2&division_id=missing-div",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 400 when challengee_race_id does not exist",
			queryParams: "id=aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa&challenger_race_id=race-1&challengee_race_id=race-missing&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "admin gets 400 when challenger_race_id does not exist",
			queryParams: "id=aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaaa&challenger_race_id=race-missing&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer admin-token",
			setupRaces: []*galaxy.Race{
				{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
				{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "returns 409 when challenge id already exists",
			queryParams: "id=99999999-9999-9999-9999-999999999999&challenger_race_id=race-1&challengee_race_id=race-2&division_id=div-1",
			authHeader:  "Bearer token-1",
			setupRaces: []*galaxy.Race{
				{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
				{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
			},
			setupDivisions: []*galaxy.Division{{ID: "div-1"}},
			setupChallenge: &galaxy.Challenge{ID: "99999999-9999-9999-9999-999999999999"},
			expectedStatus: http.StatusConflict,
		},
	}
}

func TestCreateChallenge(t *testing.T) {
	for _, tc := range createChallengeTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, divisionRepo, _ := newChallengeController()
			for _, r := range tc.setupRaces {
				raceRepo.Upsert(r)
			}
			for _, d := range tc.setupDivisions {
				divisionRepo.Upsert(d)
			}
			if tc.setupChallenge != nil {
				challengeRepo.Upsert(tc.setupChallenge)
			}

			router := setupChallengeRouter(controller)
			req, _ := http.NewRequest(http.MethodPost, "/api/challenges?"+tc.queryParams, nil)
			req.Header.Set("Authorization", tc.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.expectedStatus == http.StatusCreated {
				var challenge galaxy.Challenge
				if err := json.Unmarshal(w.Body.Bytes(), &challenge); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if challenge.Status != galaxy.ChallengeStatusPending {
					t.Errorf("expected status %q, got %q", galaxy.ChallengeStatusPending, challenge.Status)
				}
			}
		})
	}
}
