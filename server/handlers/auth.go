package handlers

import (
	"time"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !user.CheckPassword(req.Password) {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	tokenString, err := models.GenerateToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	authToken := models.AuthToken{
		Token:       tokenString,
		UserID:      user.ID,
		SessionData: "",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	if err := database.DB.Create(&authToken).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create session"})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user":  user,
	})
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username and password are required"})
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Check if email exists (if provided)
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
		}
	}

	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		Password:  req.Password,
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	tokenString, err := models.GenerateToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	authToken := models.AuthToken{
		Token:       tokenString,
		UserID:      user.ID,
		SessionData: "",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	if err := database.DB.Create(&authToken).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create session"})
	}

	return c.Status(201).JSON(fiber.Map{
		"token": tokenString,
		"user":  user,
	})
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	if err := database.DB.Where("token = ?", tokenString).Delete(&models.AuthToken{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not destroy session"})
	}

	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}

func GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	return c.JSON(user)
}
