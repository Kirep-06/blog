package handler

import (
	"net/http"
	"strconv"

	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

func ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, total, err := service.ListPosts(service.PostFilter{
		CategorySlug: c.Query("category"),
		TagSlug:      c.Query("tag"),
		Search:       c.Query("q"),
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPost(c *gin.Context) {
	post, err := service.GetPost(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

func CreatePost(c *gin.Context) {
	var in service.CreatePostInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("userID")
	post, err := service.CreatePost(userID, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func UpdatePost(c *gin.Context) {
	var in service.UpdatePostInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := service.UpdatePost(c.Param("slug"), in)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

func DeletePost(c *gin.Context) {
	if err := service.DeletePost(c.Param("slug")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
