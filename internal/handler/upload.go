package handler

import (
	"net/http"

	"blog/internal/service"
	"blog/internal/storage"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Provider storage.StorageProvider
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	fh, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image field required"})
		return
	}

	userID := c.GetUint("userID")
	img, err := service.UploadImage(c.Request.Context(), h.Provider, fh, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": img.URL})
}
