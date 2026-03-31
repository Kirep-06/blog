package handler

import (
	"net/http"
	"strconv"

	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

func ListAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := service.ListAllPosts(service.AdminPostFilter{
		Search:    c.Query("q"),
		Published: c.Query("published"),
		Page:      page,
		PageSize:  pageSize,
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

func GetAnyPost(c *gin.Context) {
	post, err := service.GetAnyPost(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}
