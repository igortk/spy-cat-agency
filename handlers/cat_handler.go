package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spy-cat-agency/database"
	"spy-cat-agency/models"
	"spy-cat-agency/utils"
)

func CreateCat(c *gin.Context) {
	var cat models.Cat

	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidBreed(cat.Breed) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat breed"})
		return
	}

	if err := database.Session.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cat"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func GetCats(c *gin.Context) {
	var cats []models.Cat
	database.Session.Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func GetCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat

	if err := database.Session.First(&cat, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func UpdateCatSalary(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat

	if err := database.Session.First(&cat, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	var body struct {
		Salary float64 `json:"salary"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Salary < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Salary cannot be negative"})
		return
	}

	database.Session.Model(&cat).Update("salary", body.Salary)
	c.JSON(http.StatusOK, cat)
}

func DeleteCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat

	if err := database.Session.First(&cat, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	database.Session.Delete(&cat)
	c.JSON(http.StatusOK, gin.H{"message": "Cat deleted successfully"})
}
