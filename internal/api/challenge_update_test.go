package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type updateChallengeTestCase struct {
	name             string
	id               string
	body             string
	authHeader       string
	setupRaces       []*galaxy.Race
	setupChallenge   *galaxy.Challenge
	setupFleetBuilds []*galaxy.FleetBuild
	expectedStatus   int
	checkResponse    func(t *testing.T, challenge galaxy.Challenge)
}

func updateChallengeTestCases() []updateChallengeTestCase {
	races := []*galaxy.Race{
		{ID: "race-1", Token: "token-1", Role: galaxy.RolePlayer},
		{ID: "race-2", Token: "token-2", Role: galaxy.RolePlayer},
		{ID: "admin", Token: "admin-token", Role: galaxy.RoleAdmin},
	}
	fleetBuilds := []*galaxy.FleetBuild{
		{ID: "fb-a", DivisionId: "div-1", RaceId: "race-1"},
		{ID: "fb-b", DivisionId: "div-1", RaceId: "race-2"},
		{ID: "fb-wrong-div", DivisionId: "div-other", RaceId: "race-1"},
	}
	return []updateChallengeTestCase{
		{
			name:             "challenger updates ready_a and fleet_build_a_id",
			id:               "c-1",
			body:             `{"ready_a":true,"fleet_build_a_id":"fb-a"}`,
			authHeader:       "Bearer token-1",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusOK,
			checkResponse: func(t *testing.T, c galaxy.Challenge) {
				if !c.ReadyA {
					t.Error("expected ready_a to be true")
				}
				if c.FleetBuildAId != "fb-a" {
					t.Errorf("expected fleet_build_a_id fb-a, got %s", c.FleetBuildAId)
				}
				if c.ReadyB {
					t.Error("expected ready_b to remain false")
				}
			},
		},
		{
			name:             "challengee updates ready_b and fleet_build_b_id",
			id:               "c-1",
			body:             `{"ready_b":true,"fleet_build_b_id":"fb-b"}`,
			authHeader:       "Bearer token-2",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusOK,
			checkResponse: func(t *testing.T, c galaxy.Challenge) {
				if !c.ReadyB {
					t.Error("expected ready_b to be true")
				}
				if c.FleetBuildBId != "fb-b" {
					t.Errorf("expected fleet_build_b_id fb-b, got %s", c.FleetBuildBId)
				}
			},
		},
		{
			name:             "challengee rejects challenge",
			id:               "c-1",
			body:             `{"status":"rejected"}`,
			authHeader:       "Bearer token-2",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusOK,
			checkResponse: func(t *testing.T, c galaxy.Challenge) {
				if c.Status != galaxy.ChallengeStatusRejected {
					t.Errorf("expected status rejected, got %s", c.Status)
				}
			},
		},
		{
			name:             "challenger fields are ignored for challengee",
			id:               "c-1",
			body:             `{"ready_a":true,"fleet_build_a_id":"fb-a","ready_b":true}`,
			authHeader:       "Bearer token-2",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusOK,
			checkResponse: func(t *testing.T, c galaxy.Challenge) {
				if c.ReadyA {
					t.Error("expected ready_a to be silently ignored")
				}
				if c.FleetBuildAId != "" {
					t.Error("expected fleet_build_a_id to be silently ignored")
				}
			},
		},
		{
			name:             "admin updates accepted_a_at",
			id:               "c-1",
			body:             `{"accepted_a_at":"2026-01-01T00:00:00Z"}`,
			authHeader:       "Bearer admin-token",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusOK,
			checkResponse: func(t *testing.T, c galaxy.Challenge) {
				if c.AcceptedAAt == nil {
					t.Error("expected accepted_a_at to be set")
				}
			},
		},
		{
			name:             "returns 400 when fleet_build_a_id is from wrong division",
			id:               "c-1",
			body:             `{"fleet_build_a_id":"fb-wrong-div"}`,
			authHeader:       "Bearer token-1",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "returns 400 when fleet_build_b_id is from wrong division",
			id:               "c-1",
			body:             `{"fleet_build_b_id":"fb-wrong-div"}`,
			authHeader:       "Bearer token-2",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "returns 404 when challenge not found",
			id:               "missing",
			body:             `{"ready_a":true}`,
			authHeader:       "Bearer token-1",
			setupRaces:       races,
			setupChallenge:   nil,
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:             "returns 403 when race is unrelated",
			id:               "c-1",
			body:             `{"ready_a":true}`,
			authHeader:       "Bearer token-3",
			setupRaces:       append(races, &galaxy.Race{ID: "race-3", Token: "token-3", Role: galaxy.RolePlayer}),
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusForbidden,
		},
		{
			name:             "returns 401 without authorization",
			id:               "c-1",
			body:             `{"ready_a":true}`,
			authHeader:       "",
			setupRaces:       races,
			setupChallenge:   freshChallenge(),
			setupFleetBuilds: fleetBuilds,
			expectedStatus:   http.StatusUnauthorized,
		},
	}
}

func TestUpdateChallenge(t *testing.T) {
	for _, tc := range updateChallengeTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			controller, raceRepo, challengeRepo, _, fleetBuildRepo := newChallengeController()
			for _, r := range tc.setupRaces {
				raceRepo.Upsert(r)
			}
			if tc.setupChallenge != nil {
				challengeRepo.Upsert(tc.setupChallenge)
			}
			for _, fb := range tc.setupFleetBuilds {
				fleetBuildRepo.Upsert(fb)
			}

			router := setupChallengeRouter(controller)
			req, _ := http.NewRequest(http.MethodPut, "/api/challenges/"+tc.id, strings.NewReader(tc.body))
			req.Header.Set("Authorization", tc.authHeader)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.expectedStatus == http.StatusOK && tc.checkResponse != nil {
				var challenge galaxy.Challenge
				if err := json.Unmarshal(w.Body.Bytes(), &challenge); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				tc.checkResponse(t, challenge)
			}
		})
	}
}
