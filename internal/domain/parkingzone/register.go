package parkingzone

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	jwtService := auth.NewJWTService(cfg.JwtSecret)
	authMiddleware := middlewares.AuthMiddleware(jwtService)
	adminMiddleware := middlewares.RoleMiddleware("admin")

	api := e.Group("/api/v1/zones")

	// Public routes
	api.GET("", handler.GetAllZones)
	api.GET("/:id", handler.GetZoneByID)

	// Admin only routes
	api.POST("", handler.CreateZone, authMiddleware, adminMiddleware)
	api.PUT("/:id", handler.UpdateZone, authMiddleware, adminMiddleware)
	api.DELETE("/:id", handler.DeleteZone, authMiddleware, adminMiddleware)
}
