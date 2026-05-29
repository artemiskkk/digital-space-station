package handler

import (
	"net/http"
	"time"

	"github.com/artemis/dss-backend/internal/config"
	"github.com/artemis/dss-backend/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store *store.Store
}

func NewAuthHandler(s *store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
		return
	}

	// Legacy single-password login (no username)
	if req.Username == "" {
		if req.Password == config.C.AdminPassword {
			// Find or assume admin user id=1
			user, _ := h.store.GetUserByUsername(c, "admin")
			uid := int64(1)
			role := "admin"
			if user != nil {
				uid = user.ID
				role = user.Role
			}
			token := signToken(uid, role)
			c.JSON(http.StatusOK, gin.H{"token": token, "user_id": uid, "role": role, "username": "admin"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	// Username + password login
	user, err := h.store.GetUserByUsername(c, req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token := signToken(user.ID, user.Role)
	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user.ID, "role": user.Role, "username": user.Username})
}

func signToken(userID int64, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     userID,
		"role":    role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	signed, _ := token.SignedString([]byte(config.C.JWTSecret))
	return signed
}

// ── User management (admin only) ──────────────────────────────────────────────

func (h *AuthHandler) ListUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}
	users, err := h.store.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": users})
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
		return
	}

	user, err := h.store.CreateUser(c, req.Username, string(hashed), "user")
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

type changePasswordRequest struct {
	UserID      int64  `json:"user_id"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	role, _ := c.Get("role")
	uid, _ := c.Get("user_id")

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Admin can change anyone's password; users can only change their own
	targetID := req.UserID
	if targetID == 0 {
		targetID = uid.(int64)
	}
	if role != "admin" && targetID != uid.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只能修改自己的密码"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := h.store.UpdatePassword(c, targetID, string(hashed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "密码已修改"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}
	id := c.Param("id")
	if err := h.store.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
