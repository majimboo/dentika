package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"dentika/server/database"
	"dentika/server/models"
)

const (
	MaxFileSize      = 5 << 20 // 5MB
	UploadBaseDir    = "uploads"
	AvatarDir        = "avatars"
	InventoryItemDir = "inventory-items"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func init() {
	// Create upload directories if they don't exist
	createUploadDirs()
}

func createUploadDirs() {
	dirs := []string{
		filepath.Join(UploadBaseDir, AvatarDir),
		filepath.Join(UploadBaseDir, InventoryItemDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
		}
	}
}

func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return uuid.New().String() + ext
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func saveFile(file *multipart.FileHeader, directory string) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	filename := generateFileName(file.Filename)
	filePath := filepath.Join(directory, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file contents
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func UploadAvatar(c *fiber.Ctx) error {
	// Get form parameters
	entityType := c.FormValue("type") // "user" or "patient"
	entityID := c.FormValue("id")

	if entityType == "" || entityID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Both 'type' (user/patient) and 'id' parameters are required",
		})
	}

	// Get the file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return c.Status(400).JSON(fiber.Map{
			"error": "File size too large. Maximum size is 5MB",
		})
	}

	// Validate file type
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed",
		})
	}

	// Create year/month directory structure
	now := time.Now()
	yearMonth := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	avatarDir := filepath.Join(UploadBaseDir, AvatarDir, yearMonth)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Save the file
	filename, err := saveFile(file, avatarDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Return the file path relative to uploads directory
	relativePath := filepath.Join(AvatarDir, yearMonth, filename)
	// Convert to forward slashes for URLs
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")

	// Update the database with the new avatar path
	if err := updateAvatarPath(entityType, entityID, relativePath); err != nil {
		// File was uploaded but database update failed - should we delete the file?
		// For now, just return an error
		return c.Status(500).JSON(fiber.Map{
			"error": "File uploaded but failed to update database: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"filename": filename,
		"path":     relativePath,
		"url":      "/uploads/" + relativePath,
	})
}

func DeleteAvatar(c *fiber.Ctx) error {
	avatarPath := c.Query("path")
	if avatarPath == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Avatar path is required",
		})
	}

	// Remove leading slash and "uploads/" if present
	avatarPath = strings.TrimPrefix(avatarPath, "/")
	avatarPath = strings.TrimPrefix(avatarPath, "uploads/")

	// Construct full file path
	fullPath := filepath.Join(UploadBaseDir, avatarPath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

func UploadInventoryItemImage(c *fiber.Ctx) error {
	// Get the file from form
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return c.Status(400).JSON(fiber.Map{
			"error": "File size too large. Maximum size is 5MB",
		})
	}

	// Validate file type
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed",
		})
	}

	// Create year/month directory structure
	now := time.Now()
	yearMonth := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	inventoryDir := filepath.Join(UploadBaseDir, InventoryItemDir, yearMonth)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(inventoryDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Save the file
	filename, err := saveFile(file, inventoryDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Return the file path relative to uploads directory
	relativePath := filepath.Join(InventoryItemDir, yearMonth, filename)
	// Convert to forward slashes for URLs
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")

	return c.JSON(fiber.Map{
		"success":  true,
		"filename": filename,
		"path":     relativePath,
		"url":      "/uploads/" + relativePath,
	})
}

func DeleteInventoryItemImage(c *fiber.Ctx) error {
	imagePath := c.Query("path")
	if imagePath == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Image path is required",
		})
	}

	// Remove leading slash and "uploads/" if present
	imagePath = strings.TrimPrefix(imagePath, "/")
	imagePath = strings.TrimPrefix(imagePath, "uploads/")

	// Construct full file path
	fullPath := filepath.Join(UploadBaseDir, imagePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// updateAvatarPath updates the avatar_path field for either a user or patient
func updateAvatarPath(entityType, entityIDStr, avatarPath string) error {
	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid entity ID: %v", err)
	}

	// Check if entity ID is valid (not 0)
	if entityID == 0 {
		return fmt.Errorf("invalid entity ID: cannot update avatar for unsaved %s", entityType)
	}

	fmt.Println(entityID, avatarPath)

	switch entityType {
	case "user":
		return database.DB.Model(&models.User{}).Where("id = ?", entityID).Update("avatar_path", avatarPath).Error
	case "patient":
		return database.DB.Model(&models.Patient{}).Where("id = ?", entityID).Update("avatar_path", avatarPath).Error
	default:
		return fmt.Errorf("invalid entity type: %s (must be 'user' or 'patient')", entityType)
	}
}
