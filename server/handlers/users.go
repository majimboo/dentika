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
		// Non-super-admin users are scoped to their clinic (ClinicID is now non-nullable)
		db = db.Where("clinic_id = ?", user.ClinicID)
	}

	// Filter by role if the query parameter is provided
	role := c.Query("role")
	if role != "" {
		if role == "doctor" {
			// Include super admins as doctors for testing purposes
			db = db.Where("role IN (?, ?)", role, models.SuperAdmin)
		} else {
			db = db.Where("role = ?", role)
		}
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
	Username   string          `json:"username"`
	Email      string          `json:"email"`
	Password   string          `json:"password"`
	Role       models.UserRole `json:"role"`
	ClinicID   *uint           `json:"clinic_id"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	AvatarPath string          `json:"avatar_path"`
}

func CreateUser(c *fiber.Ctx) error {
	// Parse request first
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Authorization check
	requestingUser := c.Locals("user").(models.User)

	// Authorization logic
	if requestingUser.IsSuperAdmin() {
		// Super admin can create users for any clinic or without clinic
	} else if requestingUser.HasRole(models.Admin) {
		// Admin can only create users for their own clinic
		if req.ClinicID == nil || *req.ClinicID != requestingUser.ClinicID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admins can only create users for their own clinic"})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions to create users"})
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

	// Determine ClinicID - if not provided, use requestingUser's ClinicID
	var userClinicID uint
	if req.ClinicID != nil {
		userClinicID = *req.ClinicID
	} else {
		userClinicID = requestingUser.ClinicID
	}

	// Create user
	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		ClinicID:   userClinicID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		AvatarPath: req.AvatarPath,
		IsActive:   true,
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

	// Get the requesting user for authorization
	requestingUser := c.Locals("user").(models.User)

	var req struct {
		Username   string          `json:"username"`
		Email      string          `json:"email"`
		FirstName  string          `json:"first_name"`
		LastName   string          `json:"last_name"`
		Gender     string          `json:"gender"`
		AvatarPath string          `json:"avatar_path"`
		Password   string          `json:"password"`
		Role       models.UserRole `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Authorization logic for role updates
	if req.Role != "" {
		if requestingUser.IsSuperAdmin() {
			// Super admin can set any role
		} else if requestingUser.HasRole(models.Admin) {
			// Admin can only update roles for users in their clinic
			if user.ClinicID != requestingUser.ClinicID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admins can only update users in their own clinic"})
			}

			// Admin can only assign staff roles (not super_admin or admin)
			validRoles := []models.UserRole{models.Doctor, models.Secretary, models.Assistant}
			isValidRole := false
			for _, validRole := range validRoles {
				if req.Role == validRole {
					isValidRole = true
					break
				}
			}
			if !isValidRole {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admins can only assign Doctor, Secretary, or Assistant roles"})
			}
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions to update user roles"})
		}
	}

	// General authorization for non-role updates
	if req.Role == "" {
		if requestingUser.IsSuperAdmin() {
			// Super admin can update any user
		} else if requestingUser.HasRole(models.Admin) {
			// Admin can only update users in their clinic
			if user.ClinicID != requestingUser.ClinicID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admins can only update users in their own clinic"})
			}
		} else {
			// Users can only update their own profile
			if requestingUser.ID != uint(userID) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only update your own profile"})
			}
		}
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
	if req.AvatarPath != "" {
		user.AvatarPath = req.AvatarPath
	}

	if req.Password != "" {
		user.Password = req.Password
		if err := user.HashPassword(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
		}
	}

	if req.Role != "" {
		user.Role = req.Role
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
