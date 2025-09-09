package handlers

import (
	"strconv"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	// Get the current user from context
	user := c.Locals("user").(models.User)

	db := database.DB

	// If user is not a super admin, scope results to their clinic
	if !user.IsSuperAdmin() {
		if user.ClinicID == nil {
			// This case should ideally be handled by middleware, but as a safeguard:
			// A non-super-admin must belong to a clinic to see other users.
			return c.JSON([]models.User{})
		}
		db = db.Where("clinic_id = ?", user.ClinicID)
	}

	// Filter by role if the query parameter is provided
	role := c.Query("role")
	if role != "" {
		db = db.Where("role = ?", role)
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
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

type CreateUserRequest struct {
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Role      models.UserRole `json:"role"`
	ClinicID  *uint           `json:"clinic_id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
}

func CreateUser(c *fiber.Ctx) error {
	// Authorization check
	requestingUser := c.Locals("user").(models.User)
	if !requestingUser.IsSuperAdmin() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only super admins can create users"})
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validation
	if req.Username == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username, password, and role are required"})
	}

	// Check for existing username
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Check for existing email
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
	}
	
	// Check if clinic exists if provided
	if req.ClinicID != nil {
		var clinic models.Clinic
		if err := database.DB.First(&clinic, *req.ClinicID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid clinic ID"})
		}
	}

	// Create user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		ClinicID:  req.ClinicID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
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
