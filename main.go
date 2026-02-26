package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"knowledge-capsule/app/handlers"
	"knowledge-capsule/app/middleware"
	_ "knowledge-capsule/docs"
	"knowledge-capsule/pkg/config"
	"knowledge-capsule/pkg/db"
	"knowledge-capsule/pkg/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Knowledge Capsule
// @version 1.0
// @description Knowledge Capsule - knowledge management backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token (or "Bearer &lt;token&gt;" for clarity). Get token from POST /api/auth/login.

// @BasePath /
func main() {
	// Load env variables first (for log level)
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load environment variables", "error", err)
		os.Exit(1)
	}

	// Set log level based on environment
	var level slog.Level
	if cfg.Env == "development" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	handlers.InitStores(database)

	if err := db.SeedSuperAdmin(database, cfg.SuperAdminEmail, cfg.SuperAdminPassword, cfg.SuperAdminName); err != nil {
		slog.Error("Failed to seed superadmin", "error", err)
		os.Exit(1)
	}

	handlers.InitChat(cfg.CORSOrigins)

	mux := http.NewServeMux()

	// Default routes
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/api", handlers.ApiRootHandler)

	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/test-ws", handlers.TestChatHandler)

	// Swagger
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	// Public routes
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)

	// User routes
	mux.Handle("/api/users", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(handlers.ListUsers))))
	mux.Handle("/api/users/", middleware.AuthMiddleware(http.HandlerFunc(handlers.UserHandler)))

	// Admin routes
	mux.Handle("/api/admin/search", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(handlers.GlobalSearch))))
	mux.Handle("/api/admin/admins", middleware.AuthMiddleware(middleware.RequireSuperAdmin(http.HandlerFunc(handlers.ListAdmins))))
	mux.Handle("/api/admin/users/", middleware.AuthMiddleware(middleware.RequireSuperAdmin(http.HandlerFunc(handlers.AdminUsersHandler))))

	// Protected routes
	mux.Handle("/api/topics", middleware.AuthMiddleware(http.HandlerFunc(handlers.TopicHandler)))
	mux.Handle("/api/topics/", middleware.AuthMiddleware(http.HandlerFunc(handlers.TopicByIDHandler)))
	mux.Handle("/api/capsules", middleware.AuthMiddleware(http.HandlerFunc(handlers.CapsuleHandler)))
	mux.Handle("/api/capsules/", middleware.AuthMiddleware(http.HandlerFunc(handlers.CapsuleByIDHandler)))

	// Chat & File Upload
	mux.Handle("/ws/chat", middleware.AuthMiddleware(http.HandlerFunc(handlers.ChatWebSocketHandler)))
	mux.Handle("/api/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadHandler)))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Wrap with CORS, logger, recover
	handler := middleware.CORS(cfg.CORSOrigins)(middleware.Recover(middleware.Logger(mux)))

	utils.InitJWTSecret(cfg.JWTSecret)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Server starting", "env", cfg.Env, "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
