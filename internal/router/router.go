package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/internal/storage"

	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, storageProvider storage.StorageProvider) {
	engine.Use(middleware.CORS())

	uploadHandler := &handler.UploadHandler{Provider: storageProvider}

	api := engine.Group("/api")

	// Auth (public)
	api.POST("/auth/login", handler.Login)

	// Posts (public read)
	posts := api.Group("/posts")
	posts.GET("", handler.ListPosts)
	posts.GET("/:slug", handler.GetPost)

	// Posts (auth required)
	authed := api.Group("")
	authed.Use(middleware.Auth())
	authed.POST("/posts", handler.CreatePost)
	authed.PUT("/posts/:slug", handler.UpdatePost)
	authed.DELETE("/posts/:slug", handler.DeletePost)

	// Categories
	api.GET("/categories", handler.ListCategories)
	authed.POST("/categories", handler.CreateCategory)
	authed.DELETE("/categories/:id", handler.DeleteCategory)

	// Tags
	api.GET("/tags", handler.ListTags)
	authed.POST("/tags", handler.CreateTag)
	authed.DELETE("/tags/:id", handler.DeleteTag)

	// Upload
	authed.POST("/upload/image", uploadHandler.UploadImage)

	// Admin (auth required)
	authed.GET("/admin/posts", handler.ListAllPosts)
	authed.GET("/admin/posts/:slug", handler.GetAnyPost)
}
