package middleware

import (
	"strings"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		var authToken models.AuthToken
		if err := database.DB.Preload("User").Where("token = ?", tokenString).First(&authToken).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		if authToken.IsExpired() {
			database.DB.Delete(&authToken)
			return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
		}

		// Check if user is active
		if !authToken.User.IsActive {
			return c.Status(403).JSON(fiber.Map{"error": "Account is inactive"})
		}

		c.Locals("user_id", authToken.UserID)
		c.Locals("user", authToken.User)
		c.Locals("auth_token", &authToken)
		return c.Next()
	}
}

// RoleMiddleware creates middleware that checks if user has required roles
func RoleMiddleware(requiredRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)

		if !user.HasRole(requiredRoles...) {
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		return c.Next()
	}
}

// ClinicAccessMiddleware ensures user can only access their clinic's data
func ClinicAccessMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)

		// Super admin can access everything
		if user.IsSuperAdmin() {
			return c.Next()
		}

		// Other users must be assigned to a clinic
		if user.ClinicID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "User not assigned to any clinic"})
		}

		c.Locals("clinic_id", *user.ClinicID)
		return c.Next()
	}
}
