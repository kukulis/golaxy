package api

import (
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

func newChallengeController() (*ChallengeController, *dao.RaceRepository, *dao.ChallengeRepository, *dao.DivisionRepository, *dao.FleetBuildRepository) {
	raceRepo := dao.NewRaceRepository()
	challengeRepo := dao.NewChallengeRepository()
	divisionRepo := dao.NewDivisionRepository()
	fleetBuildRepo := dao.NewFleetBuildRepository()
	authManager := NewMemoryAuthenticationManager(raceRepo)
	controller := NewChallengeController(authManager, challengeRepo, divisionRepo, raceRepo, fleetBuildRepo)
	return controller, raceRepo, challengeRepo, divisionRepo, fleetBuildRepo
}

func freshChallenge() *galaxy.Challenge {
	return &galaxy.Challenge{
		ID:               "c-1",
		DivisionId:       "div-1",
		ChallengerRaceId: "race-1",
		ChallengeeRaceId: "race-2",
		Status:           galaxy.ChallengeStatusPending,
	}
}
