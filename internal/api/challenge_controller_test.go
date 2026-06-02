package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/dao"
	"glaktika.eu/galaktika/pkg/galaxy"
)

func setupChallengeRouter(controller *ChallengeController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/api")
	apiGroup.POST("/challenges", controller.CreateChallenge)
	apiGroup.GET("/challenges", controller.GetChallenges)
	apiGroup.GET("/challenges/:id", controller.GetChallenge)
	apiGroup.PUT("/challenges/:id", controller.UpdateChallenge)
	apiGroup.DELETE("/challenges/:id", controller.DeleteChallenge)
	return router
}

func newChallengeController() (*ChallengeController, *dao.RaceRepository, *dao.ChallengeRepository, *dao.DivisionRepository) {
	raceRepo := dao.NewRaceRepository()
	challengeRepo := dao.NewChallengeRepository()
	divisionRepo := dao.NewDivisionRepository()
	authManager := NewMemoryAuthenticationManager(raceRepo)
	controller := NewChallengeController(authManager, challengeRepo, divisionRepo, raceRepo)
	return controller, raceRepo, challengeRepo, divisionRepo
}

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
			controller, raceRepo, challengeRepo, divisionRepo := newChallengeController()
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

type getChallengesTestCase struct {
	name             string
	queryParams      string
	authHeader       string
	setupRaces       []*galaxy.Race
	setupChallenges  []*galaxy.Challenge
	expectedStatus   int
	expectedCount    int
}

func getChallengesTestCases() []getChallengesTestCase {
	now := time.Now()
	challengeA := &galaxy.Challenge{ID: "c-1", ChallengerRaceId: "race-1", ChallengeeRaceId: "race-2", AcceptedAAt: &now}
	challengeB := &galaxy.Challenge{ID: "c-2", ChallengerRaceId: "race-1", ChallengeeRaceId: "race-3", AcceptedBAt: &now}
	challengeC := &galaxy.Challenge{ID: "c-3", ChallengerRaceId: "race-3", ChallengeeRaceId: "race-2", AcceptedAAt: &now, AcceptedBAt: &now}

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
			name:            "filter by accepted_a",
			queryParams:     "challenger_id=race-1&accepted_a=1",
			authHeader:      "Bearer token-1",
			setupRaces:      races,
			setupChallenges: allChallenges,
			expectedStatus:  http.StatusOK,
			expectedCount:   1,
		},
		{
			name:            "filter by accepted_b",
			queryParams:     "challengee_id=race-2&accepted_b=1",
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
			controller, raceRepo, challengeRepo, _ := newChallengeController()
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

func TestGetChallenge(t *testing.T) {
	for _, tc := range getChallengeTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, _ := newChallengeController()
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
