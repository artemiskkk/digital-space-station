package handler

import (
	"net/http"
	"strconv"

	"github.com/artemis/dss-backend/internal/model"
	"github.com/artemis/dss-backend/internal/store"
	"github.com/gin-gonic/gin"
)

type MilestonesHandler struct {
	store *store.Store
}

func NewMilestonesHandler(s *store.Store) *MilestonesHandler {
	return &MilestonesHandler{store: s}
}

func (h *MilestonesHandler) List(c *gin.Context) {
	uid := getUserID(c)
	items, err := h.store.ListMilestones(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []model.Milestone{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MilestonesHandler) Create(c *gin.Context) {
	uid := getUserID(c)
	var m model.Milestone
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateMilestone(c, &m, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (h *MilestonesHandler) Delete(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.store.DeleteMilestone(c, id, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
