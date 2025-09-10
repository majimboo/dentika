package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
