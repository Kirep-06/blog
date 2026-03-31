package handler

import (
	"net/http"

	"blog/internal/database"
	"blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ListCategories(c *gin.Context) {
	var cats []model.Category
	if err := database.DB.Order("name").Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cats})
}

func CreateCategory(c *gin.Context) {
	var in struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := model.Category{Name: in.Name, Slug: slug.Make(in.Name)}
	if err := database.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

func DeleteCategory(c *gin.Context) {
	if err := database.DB.Delete(&model.Category{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
