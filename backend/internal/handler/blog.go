package handler

import (
	"net/http"
	"strconv"

	"github.com/artemis/dss-backend/internal/model"
	"github.com/artemis/dss-backend/internal/store"
	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	store *store.Store
}

func NewBlogHandler(s *store.Store) *BlogHandler {
	return &BlogHandler{store: s}
}

func getUserID(c *gin.Context) int64 {
	if v, ok := c.Get("user_id"); ok {
		return v.(int64)
	}
	return 0
}

// Public: admin's published posts (for guest preview)
func (h *BlogHandler) ListPublicPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 { page = 1 }
	if size < 1 || size > 50 { size = 10 }

	list, err := h.store.ListPublicPosts(c, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Auth: user's own posts
func (h *BlogHandler) ListMyPosts(c *gin.Context) {
	uid := getUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))
	if page < 1 { page = 1 }

	list, err := h.store.ListPostsByUser(c, uid, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *BlogHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")
	uid := getUserID(c)
	post, err := h.store.GetPostBySlug(c, slug, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *BlogHandler) CreatePost(c *gin.Context) {
	uid := getUserID(c)
	var p model.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreatePost(c, &p, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *BlogHandler) UpdatePost(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p model.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := h.store.UpdatePost(c, &p, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *BlogHandler) DeletePost(c *gin.Context) {
	uid := getUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.store.DeletePost(c, id, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
