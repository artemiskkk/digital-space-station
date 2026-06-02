package main

import (
	"context"
	"log"
	"net/http"

	"github.com/artemis/dss-backend/internal/config"
	"github.com/artemis/dss-backend/internal/handler"
	"github.com/artemis/dss-backend/internal/middleware"
	"github.com/artemis/dss-backend/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.Load()

	poolCfg, err := pgxpool.ParseConfig(config.C.DatabaseURL)
	if err != nil {
		log.Fatalf("parse db config: %v", err)
	}
	poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("database connected")

	// Ensure admin user exists
	ensureAdminUser(db)

	s := store.New(db)

	authH := handler.NewAuthHandler(s)
	blogH := handler.NewBlogHandler(s)
	momentsH := handler.NewMomentsHandler(s)
	milestonesH := handler.NewMilestonesHandler(s)
	uploadH := handler.NewUploadHandler()

	r := gin.Default()
	r.Use(middleware.CORS())

	// ── Public routes (no auth) ───────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/api/auth/login", authH.Login)

	// Public preview: only admin's published posts (no full content)
	r.GET("/api/public/posts", blogH.ListPublicPosts)

	// ── Auth routes (logged-in users) ─────────────────────────────────────────
	auth := r.Group("/api", middleware.JWT())
	{
		// Posts (user's own)
		auth.GET("/posts", blogH.ListMyPosts)
		auth.GET("/posts/:slug", blogH.GetPost)
		auth.POST("/posts", blogH.CreatePost)
		auth.PUT("/posts/:id", blogH.UpdatePost)
		auth.DELETE("/posts/:id", blogH.DeletePost)

		// Moments (user's own)
		auth.GET("/moments", momentsH.ListMoments)
		auth.POST("/moments", momentsH.CreateMoment)
		auth.DELETE("/moments/:id", momentsH.DeleteMoment)

		// Milestones (user's own)
		auth.GET("/milestones", milestonesH.List)
		auth.POST("/milestones", milestonesH.Create)
		auth.DELETE("/milestones/:id", milestonesH.Delete)

		// Image upload
		auth.POST("/upload/image", uploadH.UploadImage)

		// User management (admin only)
		auth.GET("/users", authH.ListUsers)
		auth.POST("/users", authH.CreateUser)
		auth.PUT("/users/password", authH.ChangePassword)
		auth.DELETE("/users/:id", authH.DeleteUser)
	}

	log.Printf("server starting on :%s", config.C.Port)
	if err := r.Run(":" + config.C.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func ensureAdminUser(db *pgxpool.Pool) {
	var count int
	db.QueryRow(context.Background(), `SELECT COUNT(*) FROM users WHERE role='admin'`).Scan(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(config.C.AdminPassword), bcrypt.DefaultCost)
		db.Exec(context.Background(),
			`INSERT INTO users (id, username, password, role) VALUES (1, $1, $2, 'admin')`,
			config.C.AdminUsername, string(hashed))
		log.Println("admin user created")
	}
}
