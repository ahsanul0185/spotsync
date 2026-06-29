package middlewares

import (
	"net/http"
	"strings"

	"spotsync/internal/auth"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Forbidden: role not found",
				})
			}

			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Forbidden: insufficient permissions",
			})
		}
	}
}
