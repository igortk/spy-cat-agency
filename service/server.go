package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"spy-cat-agency/config"
	"spy-cat-agency/handlers"
	"spy-cat-agency/middleware"
)

func Run(httpCfg *config.HttpConfig) {
	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/cats", handlers.CreateCat)
		api.GET("/cats", handlers.GetCats)
		api.GET("/cats/:id", handlers.GetCat)
		api.PUT("/cats/:id/salary", handlers.UpdateCatSalary)
		api.DELETE("/cats/:id", handlers.DeleteCat)

		api.POST("/missions", handlers.CreateMission)
		api.GET("/missions", handlers.GetMissions)

		missionGroup := api.Group("/missions/:id")
		{
			missionGroup.GET("", handlers.GetMission)
			missionGroup.PUT("/complete", handlers.CompleteMission)
			missionGroup.POST("/assign-cat", handlers.AssignCat)
			missionGroup.DELETE("", handlers.DeleteMission)

			missionGroup.POST("/targets", handlers.AddTarget)
			missionGroup.PUT("/targets/:targetId/notes", handlers.UpdateTargetNotes)
			missionGroup.PUT("/targets/:targetId/complete", handlers.CompleteTarget)
			missionGroup.DELETE("/targets/:targetId", handlers.DeleteTarget)
		}
	}

	if err := router.Run(fmt.Sprintf(":%d", httpCfg.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
