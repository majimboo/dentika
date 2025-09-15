# Dentika Notification System

A comprehensive, real-time notification system with database persistence, NATS integration, and full-stack support.

## Overview

The notification system provides:

- **Database Persistence**: All notifications are stored in the database with read/dismiss status
- **Real-time Delivery**: NATS.js integration for instant notification delivery
- **Multi-scope Support**: User-specific, clinic-wide, and system-wide notifications
- **Rich Content**: Support for actions, icons, colors, and scheduling
- **API Integration**: RESTful API endpoints for notification management
- **Frontend Integration**: Vue.js store and components with real-time updates

## Architecture

### Backend Components

1. **Models** (`server/models/notification.go`)
   - `Notification`: Main notification model
   - `NotificationRecipient`: Tracks delivery and read status per user
   - Enums for types, scopes, and statuses

2. **Services** (`server/services/notification_service.go`)
   - `NotificationService`: Core business logic
   - NATS integration for real-time distribution
   - Convenience methods for common notification types

3. **Handlers** (`server/handlers/notification.go`)
   - REST API endpoints
   - User authentication and authorization
   - Request/response handling

### Frontend Components

1. **Store** (`frontend/src/stores/notification.js`)
   - Pinia store for state management
   - API integration methods
   - Real-time NATS integration

2. **Composables** (`frontend/src/composables/useNats.js`)
   - NATS connection management
   - Real-time message handling
   - Automatic reconnection

3. **Components**
   - `NotificationPanel.vue`: Dropdown and mobile panels
   - `NotificationToast.vue`: Toast notifications
   - `Notifications.vue`: Full notifications page

## NATS Subject Structure

All subjects are prefixed with `dentika.*` to prevent conflicts:

```
dentika.user.{userId}.notifications          # User-specific notifications
dentika.user.{userId}.appointments           # User appointment updates
dentika.user.{userId}.patients                # User patient updates
dentika.clinic.{clinicId}.notifications      # Clinic-wide notifications
dentika.clinic.{clinicId}.appointments       # Clinic appointment updates
dentika.clinic.{clinicId}.patients           # Clinic patient updates
dentika.system.notifications                 # System-wide notifications
dentika.system.alerts                        # System alerts
dentika.heartbeat                            # Client heartbeat
```

## Database Schema

### Notifications Table
```sql
CREATE TABLE notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    user_id BIGINT NULL,           -- NULL for clinic/system wide
    clinic_id BIGINT NULL,         -- NULL for system wide
    data JSON,                     -- Additional data
    actions JSON,                  -- Action buttons
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    is_dismissed BOOLEAN DEFAULT FALSE,
    dismissed_at TIMESTAMP NULL,
    scheduled_for TIMESTAMP NULL,  -- For delayed notifications
    expires_at TIMESTAMP NULL,     -- When notification expires
    created_by_id BIGINT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_clinic_id (clinic_id),
    INDEX idx_is_read (is_read),
    INDEX idx_is_dismissed (is_dismissed),
    INDEX idx_scheduled_for (scheduled_for),
    INDEX idx_expires_at (expires_at)
);
```

### Notification Recipients Table
```sql
CREATE TABLE notification_recipients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    notification_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    is_dismissed BOOLEAN DEFAULT FALSE,
    dismissed_at TIMESTAMP NULL,
    delivered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_notification_user (notification_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_is_dismissed (is_dismissed)
);
```

## API Endpoints

### GET /api/notifications
Get user's notifications with pagination and filtering.

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 10, max: 100)
- `filter` (all, unread, read)

**Response:**
```json
{
  "success": true,
  "data": {
    "notifications": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25
    }
  }
}
```

### GET /api/notifications/unread-count
Get count of unread notifications for the user.

**Response:**
```json
{
  "success": true,
  "data": {
    "count": 5
  }
}
```

### PUT /api/notifications/:id/read
Mark a specific notification as read.

### PUT /api/notifications/mark-all-read
Mark all user's notifications as read.

### PUT /api/notifications/:id/dismiss
Dismiss a specific notification.

### POST /api/notifications (Admin only)
Create a new notification.

**Request Body:**
```json
{
  "title": "System Maintenance",
  "message": "Scheduled maintenance tonight at 2 AM",
  "type": "system_alert",
  "user_id": null,           // null for clinic/system wide
  "clinic_id": null,         // null for system wide
  "scheduled_for": "2024-01-15T02:00:00Z",
  "expires_at": "2024-01-16T02:00:00Z",
  "actions": [
    {
      "label": "View Details",
      "action": "navigate",
      "url": "/maintenance"
    }
  ]
}
```

### POST /api/notifications/test
Create a test notification for the current user.

## Frontend Usage

### Using the Notification Store

```javascript
import { useNotificationStore } from '@/stores/notification'

const notificationStore = useNotificationStore()

// Initialize (fetch existing notifications)
await notificationStore.initialize()

// Fetch notifications with pagination
await notificationStore.fetchNotifications(page, limit, filter)

// Get unread count
await notificationStore.fetchUnreadCount()

// Mark as read
await notificationStore.markAsReadAPI(notificationId)

// Mark all as read
await notificationStore.markAllAsReadAPI()

// Dismiss notification
await notificationStore.dismissNotification(notificationId)

// Create test notification
await notificationStore.createTestNotification()

// Real-time: Add notification from NATS
notificationStore.addNotificationFromNATS(natsNotification)
```

### Using the NATS Composable

```javascript
import { useNats } from '@/composables/useNats'

const {
  connection,
  isConnected,
  connect,
  disconnect,
  publish
} = useNats()

// Connection is automatic on mount if authenticated
// Publishing custom notifications
await publish('dentika.user.123.notifications', {
  title: 'Custom Notification',
  message: 'This is a custom message'
})
```

### Server-Side Usage

```go
// Initialize notification service
notificationService := services.NewNotificationService(natsConn)

// Create appointment reminder
err := notificationService.CreateAppointmentReminder(appointment, 15)

// Create appointment update
err := notificationService.CreateAppointmentUpdate(appointment, "rescheduled")

// Create patient update
err := notificationService.CreatePatientUpdate(patient, "created")

// Create custom notification
notification, err := notificationService.CreateNotification(services.CreateNotificationRequest{
    Title:    "Custom Notification",
    Message:  "This is a custom notification",
    Type:     models.NotificationTypeInfo,
    UserID:   &userID,
    Actions: []models.NotificationAction{
        {
            Label:  "View Details",
            Action: "navigate",
            URL:    "/details",
        },
    },
})

// Mark as read
err := notificationService.MarkAsRead(notificationID, userID)

// Get user notifications
notifications, total, err := notificationService.GetUserNotifications(userID, 10, 0)
```

## Notification Types

- `success`: Success messages (green)
- `error`: Error messages (red)
- `warning`: Warning messages (yellow)
- `info`: Information messages (blue)
- `appointment_reminder`: Appointment reminders (blue)
- `appointment_update`: Appointment changes (purple)
- `patient_update`: Patient information changes (green)
- `system_alert`: System alerts (red)
- `clinic_announcement`: Clinic announcements (blue)
- `peer_review_update`: Peer review updates (purple)
- `inventory_alert`: Inventory alerts (yellow)

## Actions

Notifications can include action buttons:

```json
{
  "actions": [
    {
      "label": "View Appointment",
      "action": "view-appointment",
      "url": "/appointments/123"
    },
    {
      "label": "Snooze 5min",
      "action": "snooze-reminder",
      "minutes": 5
    }
  ]
}
```

## Configuration

### Environment Variables

**Server:**
- `NATS_URL`: NATS server URL (default: `nats://localhost:4222`)

**Frontend:**
The NATS WebSocket URLs are configured in `useNats.js`:
- Localhost: `wss://localhost:9222`
- Production: `wss://nats.majidarif.com:9222`

## Testing

Use the `NotificationTest.vue` component to test the notification system:

```html
<NotificationTest />
```

This component provides:
- Test notification creation
- Mark all as read functionality
- Refresh notifications
- Clear all notifications
- Real-time statistics display
- Recent notifications list

## Best Practices

1. **Always use the API methods** for persistent notifications
2. **Use the convenience methods** for common notification types
3. **Include meaningful actions** for better user engagement
4. **Set expiration dates** for time-sensitive notifications
5. **Use appropriate notification types** for proper styling and categorization
6. **Test real-time delivery** using multiple browser sessions
7. **Handle offline scenarios** gracefully with automatic reconnection

## Future Enhancements

- Email/SMS notification delivery
- Push notifications for mobile apps
- Notification templates and customization
- Advanced scheduling and recurring notifications
- Analytics and notification engagement metrics
- User notification preferences and filtering