package handlers

import (
	"log"
	"net/http"
	"strconv"

	"dentika/server/models"
	"dentika/server/services"

	"github.com/gofiber/fiber/v2"
)

var notificationService *services.NotificationService

// SetNotificationService sets the notification service instance
func SetNotificationService(service *services.NotificationService) {
	notificationService = service
}

// GetUserNotifications retrieves notifications for the authenticated user
func GetUserNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter := c.Query("filter", "all") // all, unread, read

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get notifications
	notifications, total, err := notificationService.GetUserNotifications(user.ID, limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve notifications",
		})
	}

	// Apply filter if specified
	if filter != "all" {
		filteredNotifications := make([]models.Notification, 0)
		for _, notification := range notifications {
			switch filter {
			case "unread":
				if !notification.IsRead {
					filteredNotifications = append(filteredNotifications, notification)
				}
			case "read":
				if notification.IsRead {
					filteredNotifications = append(filteredNotifications, notification)
				}
			}
		}
		notifications = filteredNotifications
		total = int64(len(filteredNotifications))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"notifications": notifications,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

// GetUnreadNotificationCount returns the count of unread notifications
func GetUnreadNotificationCount(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	count, err := notificationService.GetUnreadCount(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"count": count,
		},
	})
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	notificationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	err = notificationService.MarkAsRead(uint(notificationID), user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification marked as read",
	})
}

// MarkAllNotificationsAsRead marks all notifications as read for the user
func MarkAllNotificationsAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Get all unread notifications for the user
	notifications, _, err := notificationService.GetUserNotifications(user.ID, 1000, 0)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve notifications",
		})
	}

	// Mark each as read
	for _, notification := range notifications {
		if !notification.IsRead {
			err := notificationService.MarkAsRead(notification.ID, user.ID)
			if err != nil {
				// Log error but continue with other notifications
				continue
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All notifications marked as read",
	})
}

// DismissNotification dismisses a notification for the user
func DismissNotification(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	notificationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	err = notificationService.MarkAsDismissed(uint(notificationID), user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to dismiss notification",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification dismissed",
	})
}

// CreateNotification creates a new notification (admin only)
func CreateNotification(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Check if user is admin or super admin
	if user.Role != models.Admin && user.Role != models.SuperAdmin {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	var req services.CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set created by
	req.CreatedByID = &user.ID

	// Validate required fields
	if req.Title == "" || req.Message == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and message are required",
		})
	}

	notification, err := notificationService.CreateNotification(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create notification",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"notification": notification,
		},
	})
}

// TestNotification creates a test notification for the current user
func TestNotification(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Add debug logging
	log.Printf("Creating test notification for user %d", user.ID)

	if notificationService == nil {
		log.Printf("ERROR: notificationService is nil")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Notification service not initialized",
		})
	}

	notification, err := notificationService.CreateNotification(services.CreateNotificationRequest{
		Title:       "Test Notification",
		Message:     "This is a test notification to verify the system is working correctly.",
		Type:        models.NotificationTypeInfo,
		UserID:      &user.ID,
		CreatedByID: &user.ID,
		Actions: []models.NotificationAction{
			{
				Label:  "View Dashboard",
				Action: "navigate",
				URL:    "/dashboard",
			},
		},
	})

	if err != nil {
		log.Printf("ERROR creating test notification: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create test notification: " + err.Error(),
		})
	}

	log.Printf("Successfully created test notification with ID %d", notification.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Test notification created",
		"data": fiber.Map{
			"notification": notification,
		},
	})
}

// SendAppointmentReminder sends appointment reminders (called by scheduler)
func SendAppointmentReminder(c *fiber.Ctx) error {
	// This would typically be called by a background job scheduler
	// For now, we'll allow admin users to trigger it manually

	user := c.Locals("user").(models.User)
	if user.Role != models.Admin && user.Role != models.SuperAdmin {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	// This is a placeholder - in a real system, you'd query for appointments
	// that need reminders and send them

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Appointment reminders processed",
	})
}

// NotificationStats returns notification statistics for admin users
func NotificationStats(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Check if user is admin or super admin
	if user.Role != models.Admin && user.Role != models.SuperAdmin {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	// Get basic stats - this would be more comprehensive in a real system
	stats := fiber.Map{
		"total_notifications": 0,
		"unread_count":       0,
		"dismissed_count":    0,
		"by_type":           fiber.Map{},
	}

	// You could implement actual statistics queries here

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}