package handlers

import (
	"strconv"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch users"})
	}
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Gender    string `json:"gender"`
		Avatar    string `json:"avatar"`
		Password  string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Username != "" {
		// Check if username already exists for another user
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Check if email already exists for another user
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
		}
		user.Email = req.Email
	}

	// Update profile fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.Password != "" {
		user.Password = req.Password
		if err := user.HashPassword(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := database.DB.Delete(&models.User{}, uint(userID)).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
