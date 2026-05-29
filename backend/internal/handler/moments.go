package handler

import (
	"net/http"
	"strconv"

	"github.com/artemis/dss-backend/internal/model"
	"github.com/artemis/dss-backend/internal/store"
	"github.com/gin-gonic/gin"
)

type MomentsHandler struct {
	store *store.Store
}

func NewMomentsHandler(s *store.Store) *MomentsHandler {
	return &MomentsHandler{store: s}
}

func (h *MomentsHandler) ListMoments(c *gin.Context) {
	uid := getUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }

	list, err := h.store.ListMoments(c, uid, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *MomentsHandler) CreateMoment(c *gin.Context) {
	uid := getUserID(c)
	var m model.Moment
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateMoment(c, &m, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *MomentsHandler) DeleteMoment(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.store.DeleteMoment(c, id, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
