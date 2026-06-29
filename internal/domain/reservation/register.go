package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/user"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	zoneRepo := parkingzone.NewRepository(db)
	userRepo := user.NewRepository(db)
	service := NewService(repo, zoneRepo, userRepo)
	handler := NewHandler(service)

	jwtService := auth.NewJWTService(cfg.JwtSecret)
	authMiddleware := middlewares.AuthMiddleware(jwtService)
	adminMiddleware := middlewares.RoleMiddleware("admin")

	api := e.Group("/api/v1/reservations", authMiddleware)

	api.POST("", handler.CreateReservation)
	api.GET("/my-reservations", handler.GetMyReservations)
	api.DELETE("/:id", handler.CancelReservation)
	api.GET("", handler.GetAllReservations, adminMiddleware)
}
