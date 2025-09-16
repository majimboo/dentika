package handlers

import (
	"strconv"

	"dentika/server/database"
	"dentika/server/models"

	"github.com/gofiber/fiber/v2"
)

type CreateClinicRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Tagline string `json:"tagline"`
}

type CreateBranchRequest struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	IsMainBranch  bool   `json:"is_main_branch"`
	Schedule      string `json:"schedule"`
	IsClosedToday bool   `json:"is_closed_today"`
}

func GetClinics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var clinics []models.Clinic
	query := database.DB.Preload("Branches").Preload("Staff")

	// Super admin can see all clinics, others only their own
	if !user.IsSuperAdmin() {
		query = query.Where("id = ?", user.ClinicID)
	}

	if err := query.Find(&clinics).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch clinics"})
	}

	return c.JSON(clinics)
}

func GetClinic(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}

	// Check access
	if !user.CanAccessClinic(uint(clinicID)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var clinic models.Clinic
	if err := database.DB.Preload("Branches").Preload("Staff").First(&clinic, clinicID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
	}

	return c.JSON(clinic)
}

func CreateClinic(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Only super admin can create clinics
	if !user.IsSuperAdmin() {
		return c.Status(403).JSON(fiber.Map{"error": "Only super admin can create clinics"})
	}

	var req CreateClinicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Clinic name is required"})
	}

	clinic := models.Clinic{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		Email:    req.Email,
		Website:  req.Website,
		Tagline:  req.Tagline,
		IsActive: true,
	}

	if err := database.DB.Create(&clinic).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create clinic"})
	}

	// Create main branch
	mainBranch := models.Branch{
		Name:         req.Name + " - Main",
		Address:      req.Address,
		Phone:        req.Phone,
		IsMainBranch: true,
		IsActive:     true,
		ClinicID:     clinic.ID,
	}

	if err := database.DB.Create(&mainBranch).Error; err != nil {
		// If branch creation fails, we should probably rollback clinic creation
		database.DB.Delete(&clinic)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create main branch"})
	}

	// Reload with relationships
	database.DB.Preload("Branches").Preload("Staff").First(&clinic, clinic.ID)

	return c.Status(201).JSON(clinic)
}

func UpdateClinic(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}

	// Check access - super admin or admin
	if !user.IsSuperAdmin() && (user.ClinicID != uint(clinicID) || !user.HasRole(models.Admin)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var clinic models.Clinic
	if err := database.DB.First(&clinic, clinicID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Clinic not found"})
	}

	var req CreateClinicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update clinic fields
	if req.Name != "" {
		clinic.Name = req.Name
	}
	if req.Address != "" {
		clinic.Address = req.Address
	}
	if req.Phone != "" {
		clinic.Phone = req.Phone
	}
	if req.Email != "" {
		clinic.Email = req.Email
	}
	if req.Website != "" {
		clinic.Website = req.Website
	}
	if req.Tagline != "" {
		clinic.Tagline = req.Tagline
	}

	if err := database.DB.Save(&clinic).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update clinic"})
	}

	// Reload with relationships
	database.DB.Preload("Branches").Preload("Staff").First(&clinic, clinic.ID)

	return c.JSON(clinic)
}

func GetClinicBranches(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}

	// Check access
	if !user.CanAccessClinic(uint(clinicID)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var branches []models.Branch
	if err := database.DB.Where("clinic_id = ?", clinicID).Find(&branches).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch branches"})
	}

	return c.JSON(branches)
}

func CreateBranch(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}

	// Check access - super admin or admin
	if !user.IsSuperAdmin() && (user.ClinicID != uint(clinicID) || !user.HasRole(models.Admin)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var req CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Branch name is required"})
	}

	// If this is supposed to be main branch, unset other main branches
	if req.IsMainBranch {
		database.DB.Model(&models.Branch{}).Where("clinic_id = ?", clinicID).Update("is_main_branch", false)
	}

	branch := models.Branch{
		Name:          req.Name,
		Address:       req.Address,
		Phone:         req.Phone,
		IsMainBranch:  req.IsMainBranch,
		IsActive:      true,
		Schedule:      req.Schedule,
		IsClosedToday: req.IsClosedToday,
		ClinicID:      uint(clinicID),
	}

	if err := database.DB.Create(&branch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create branch"})
	}

	return c.Status(201).JSON(branch)
}

func DeleteClinic(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}

	// Only super admin can delete clinics
	if !user.IsSuperAdmin() {
		return c.Status(403).JSON(fiber.Map{"error": "Only super admin can delete clinics"})
	}

	// Use a transaction to ensure all or nothing
	tx := database.DB.Begin()

	// Deactivate users associated with the clinic
	if err := tx.Model(&models.User{}).Where("clinic_id = ?", clinicID).Update("is_active", false).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to deactivate users in clinic"})
	}

	// Soft delete branches of the clinic
	if err := tx.Where("clinic_id = ?", clinicID).Delete(&models.Branch{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete branches of clinic"})
	}

	// Soft delete the clinic itself
	if err := tx.Delete(&models.Clinic{}, clinicID).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete clinic"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.SendStatus(204) // No Content
}

func UpdateBranch(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}
	branchID, err := strconv.ParseUint(c.Params("branch_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid branch ID"})
	}

	// Check access - super admin or admin
	if !user.IsSuperAdmin() && (user.ClinicID != uint(clinicID) || !user.HasRole(models.Admin)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var branch models.Branch
	if err := database.DB.Where("id = ? AND clinic_id = ?", branchID, clinicID).First(&branch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Branch not found"})
	}

	var req CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// If this is supposed to be main branch, unset other main branches
	if req.IsMainBranch && !branch.IsMainBranch {
		database.DB.Model(&models.Branch{}).Where("clinic_id = ?", clinicID).Update("is_main_branch", false)
	}

	branch.Name = req.Name
	branch.Address = req.Address
	branch.Phone = req.Phone
	branch.IsMainBranch = req.IsMainBranch
	branch.Schedule = req.Schedule
	branch.IsClosedToday = req.IsClosedToday

	if err := database.DB.Save(&branch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update branch"})
	}

	return c.JSON(branch)
}

func DeleteBranch(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	clinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid clinic ID"})
	}
	branchID, err := strconv.ParseUint(c.Params("branch_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid branch ID"})
	}

	// Check access - super admin or admin
	if !user.IsSuperAdmin() && (user.ClinicID != uint(clinicID) || !user.HasRole(models.Admin)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	var branch models.Branch
	if err := database.DB.Where("id = ? AND clinic_id = ?", branchID, clinicID).First(&branch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Branch not found"})
	}

	// Prevent deleting the main branch
	if branch.IsMainBranch {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot delete the main branch. Please assign another branch as main first."})
	}

	if err := database.DB.Delete(&branch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete branch"})
	}

	return c.SendStatus(204)
}
