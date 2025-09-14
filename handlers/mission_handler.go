package handlers

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"net/http"
	"spy-cat-agency/database"
	"spy-cat-agency/models"
)

func CreateMission(c *gin.Context) {
	var req struct {
		CatID   string          `json:"catId" binding:"required,uuid"`
		Targets []models.Target `json:"targets" binding:"required,min=1,max=3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cat models.Cat
	if err := database.Session.First(&cat, "id = ?", req.CatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}
	if cat.Status != "Available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cat is not available for a mission"})
		return
	}

	tx := database.Session.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	mission := models.Mission{CatID: req.CatID}
	if err := tx.Create(&mission).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create mission"})
		return
	}

	for i := range req.Targets {
		req.Targets[i].MissionID = mission.ID
		if err := tx.Create(&req.Targets[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create targets"})
			return
		}
	}

	tx.Model(&cat).Update("status", "On Mission")

	tx.Commit()
	c.JSON(http.StatusCreated, mission)
}

func GetMissions(c *gin.Context) {
	var missions []models.Mission
	database.Session.Preload("Targets").Find(&missions)
	c.JSON(http.StatusOK, missions)
}

func GetMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := database.Session.Preload("Targets").First(&mission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}
	c.JSON(http.StatusOK, mission)
}

func CompleteMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := database.Session.Preload("Targets").First(&mission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	for _, target := range mission.Targets {
		if !target.IsCompleted {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All targets must be completed before marking the mission as complete"})
			return
		}
	}

	tx := database.Session.Begin()
	tx.Model(&mission).Update("is_completed", true)

	var cat models.Cat
	if err := tx.First(&cat, "id = ?", mission.CatID).Error; err == nil {
		tx.Model(&cat).Update("status", "Available")
	}

	tx.Commit()
	c.JSON(http.StatusOK, mission)
}

func AssignCat(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CatID string `json:"catId" binding:"required,uuid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cat models.Cat
	if err := database.Session.First(&cat, "id = ?", req.CatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}
	if cat.Status != "Available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cat is already on a mission"})
		return
	}

	var mission models.Mission
	if err := database.Session.First(&mission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}
	if mission.CatID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mission already has a cat assigned"})
		return
	}

	tx := database.Session.Begin()
	tx.Model(&mission).Update("cat_id", req.CatID)
	tx.Model(&cat).Update("status", "On Mission")
	tx.Commit()

	c.JSON(http.StatusOK, mission)
}

func DeleteMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := database.Session.First(&mission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.CatID != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete mission assigned to a cat"})
		return
	}

	database.Session.Delete(&mission)
	c.JSON(http.StatusOK, gin.H{"message": "Mission deleted successfully"})
}

func AddTarget(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := database.Session.Preload("Targets").First(&mission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.IsCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add a target to a completed mission"})
		return
	}

	if len(mission.Targets) >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mission already has the maximum number of targets (3)"})
		return
	}

	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target.MissionID = mission.ID
	if err := database.Session.Create(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add target"})
		return
	}

	c.JSON(http.StatusCreated, target)
}

func UpdateTargetNotes(c *gin.Context) {
	missionId := c.Param("missionId")
	targetId := c.Param("targetId")

	var target models.Target
	if err := database.Session.First(&target, "id = ? AND mission_id = ?", targetId, missionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if target.IsCompleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update notes on a completed target"})
		return
	}

	var mission models.Mission
	if err := database.Session.First(&mission, "id = ?", missionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.IsCompleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update notes on a target within a completed mission"})
		return
	}

	var req struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target.Notes = append(target.Notes, req.Note)
	if err := database.Session.Save(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update notes"})
		return
	}

	c.JSON(http.StatusOK, target)
}

func CompleteTarget(c *gin.Context) {
	missionId := c.Param("missionId")
	targetId := c.Param("targetId")

	var target models.Target
	if err := database.Session.First(&target, "id = ? AND mission_id = ?", targetId, missionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if target.IsCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target is already completed"})
		return
	}

	database.Session.Model(&target).Update("is_completed", true)
	c.JSON(http.StatusOK, target)
}

func DeleteTarget(c *gin.Context) {
	missionId := c.Param("missionId")
	targetId := c.Param("targetId")

	var target models.Target
	if err := database.Session.First(&target, "id = ? AND mission_id = ?", targetId, missionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	if target.IsCompleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete a completed target"})
		return
	}

	database.Session.Delete(&target)
	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}
