package handler

import (
	"net/http"

	"blog/internal/database"
	"blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ListTags(c *gin.Context) {
	var tags []model.Tag
	if err := database.DB.Order("name").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func CreateTag(c *gin.Context) {
	var in struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag := model.Tag{Name: in.Name, Slug: slug.Make(in.Name)}
	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "tag already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

func DeleteTag(c *gin.Context) {
	if err := database.DB.Delete(&model.Tag{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
