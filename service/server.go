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
		api.POST("/cats", handlers.CreateCat)                 //Ability to create a spy cat in the system
		api.GET("/cats", handlers.GetCats)                    //Ability to list spy cats
		api.GET("/cats/:id", handlers.GetCat)                 //Ability to get a single spy cat
		api.PUT("/cats/:id/salary", handlers.UpdateCatSalary) //Ability to update spy cats’ information (Salary)
		api.DELETE("/cats/:id", handlers.DeleteCat)           //Ability to remove spy cats from the system

		api.POST("/missions", handlers.CreateMission) //Ability to create a mission in the system along with targets
		api.GET("/missions", handlers.GetMissions)    //Ability to list missions

		missionGroup := api.Group("/missions/:id")
		{
			missionGroup.GET("", handlers.GetMission)               //Ability to get a single mission
			missionGroup.PUT("/complete", handlers.CompleteMission) //Ability to update mission
			missionGroup.POST("/assign-cat", handlers.AssignCat)    //Ability to assign a cat to a mission
			missionGroup.DELETE("", handlers.DeleteMission)         //Ability to delete a mission

			missionGroup.POST("/targets", handlers.AddTarget)                //Ability to add targets to an existing mission
			missionGroup.DELETE("/targets/:targetId", handlers.DeleteTarget) //Ability to delete targets from an existing mission
			//Ability to update mission targets
			missionGroup.PUT("/targets/:targetId/notes", handlers.UpdateTargetNotes) //Ability to update Notes
			missionGroup.PUT("/targets/:targetId/complete", handlers.CompleteTarget) //Ability to mark it as completed
		}
	}

	if err := router.Run(fmt.Sprintf(":%d", httpCfg.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
