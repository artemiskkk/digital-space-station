package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler { return &UploadHandler{} }

func (h *UploadHandler) UploadImage(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "图片上传暂未开放"})
}
