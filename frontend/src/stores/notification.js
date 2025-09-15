import { defineStore } from 'pinia'
import { ref, computed, nextTick } from 'vue'
import apiService from '../services/api.js'

export const useNotificationStore = defineStore('notification', () => {
  // State
  const notifications = ref([])
  const unreadCount = ref(0)
  const isLoading = ref(false)
  const isInitialized = ref(false)
  const error = ref(null)

  // Settings for different display areas
  const toastSettings = ref({
    autoRemoveDelay: 5000,
    maxVisible: 5,
    position: 'top-right'
  })

  // Event listeners for real-time updates
  const eventListeners = ref(new Set())

  // Getters
  const recentNotifications = computed(() => {
    return notifications.value.slice(0, 10) // Show last 10 notifications
  })

  const hasUnread = computed(() => unreadCount.value > 0)

  const unreadNotifications = computed(() => {
    return notifications.value.filter(n => !n.is_read)
  })

  const readNotifications = computed(() => {
    return notifications.value.filter(n => n.is_read)
  })

  const notificationsByType = computed(() => {
    const grouped = {}
    notifications.value.forEach(notification => {
      const type = notification.type || 'info'
      if (!grouped[type]) grouped[type] = []
      grouped[type].push(notification)
    })
    return grouped
  })

  const toastNotifications = computed(() => {
    return notifications.value
      .filter(n => n.showAsToast && !n.dismissedFromToast)
      .slice(0, toastSettings.value.maxVisible)
  })

  // Actions - Central notification management
  const addNotification = (notification, options = {}) => {
    const {
      source = 'local', // 'local', 'api', 'nats'
      showAsToast = true,
      persistInList = true,
      autoRemove = true
    } = options

    // Normalize notification structure
    const newNotification = {
      // Core properties
      id: notification.id || Date.now() + Math.random(),
      title: notification.title || 'Notification',
      message: notification.message || '',
      type: notification.type || 'info',

      // Timing
      timestamp: notification.created_at ? new Date(notification.created_at) : new Date(),
      created_at: notification.created_at || new Date().toISOString(),

      // Status - use is_read from backend (from notification_recipients join)
      is_read: notification.is_read || false,

      // Display settings
      showAsToast,
      persistInList,
      dismissedFromToast: false,
      autoRemove, // Store the auto-remove setting

      // Data
      icon: notification.icon,
      color: notification.color,
      data: notification.data || {},
      actions: notification.actions || [],

      // Source tracking
      source,

      // Original notification data
      ...notification
    }

    // Check for duplicates (especially important for NATS messages)
    const existingIndex = notifications.value.findIndex(n =>
      n.id === newNotification.id ||
      (n.title === newNotification.title && n.message === newNotification.message &&
       Math.abs(new Date(n.created_at) - new Date(newNotification.created_at)) < 1000)
    )

    if (existingIndex !== -1) {
      // Update existing notification instead of creating duplicate
      notifications.value[existingIndex] = { ...notifications.value[existingIndex], ...newNotification }
      emitNotificationEvent('updated', notifications.value[existingIndex])
      return notifications.value[existingIndex]
    }

    // Add to notifications list
    if (persistInList) {
      notifications.value.unshift(newNotification)

      // Limit total notifications in memory
      if (notifications.value.length > 100) {
        notifications.value = notifications.value.slice(0, 100)
      }
    }

    // Update unread count
    if (!newNotification.is_read) {
      unreadCount.value++
    }

    // Auto-remove logic
    if (autoRemove && newNotification.showAsToast) {
      setTimeout(() => {
        dismissFromToast(newNotification.id)
      }, toastSettings.value.autoRemoveDelay)
    }

    // Emit events for components to react
    emitNotificationEvent('added', newNotification)

    return newNotification
  }

  // Event system for components to react to changes
  const emitNotificationEvent = (eventType, notification) => {
    nextTick(() => {
      eventListeners.value.forEach(callback => {
        try {
          callback(eventType, notification)
        } catch (error) {
          console.error('Error in notification event listener:', error)
        }
      })
    })
  }

  const subscribe = (callback) => {
    eventListeners.value.add(callback)
    return () => eventListeners.value.delete(callback)
  }

  // Enhanced management methods
  const removeNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      const notification = notifications.value[index]
      if (!notification.is_read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      notifications.value.splice(index, 1)
      emitNotificationEvent('removed', notification)
    }
  }

  const dismissFromToast = (id) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.dismissedFromToast = true
      emitNotificationEvent('dismissedFromToast', notification)
    }
  }

  const showInToast = (id) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.dismissedFromToast = false
      notification.showAsToast = true
      emitNotificationEvent('showInToast', notification)
    }
  }

  const markAsRead = (id) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification && !notification.is_read) {
      notification.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
      emitNotificationEvent('read', notification)
    }
  }

  const markAllAsRead = () => {
    let changedCount = 0
    notifications.value.forEach(notification => {
      if (!notification.is_read) {
        notification.is_read = true
        changedCount++
        emitNotificationEvent('read', notification)
      }
    })
    unreadCount.value = 0
    if (changedCount > 0) {
      emitNotificationEvent('allRead', { count: changedCount })
    }
  }

  const clearAll = () => {
    notifications.value = []
    unreadCount.value = 0
  }

  // API integration methods
  const fetchNotifications = async (page = 1, limit = 10, filter = 'all') => {
    try {
      isLoading.value = true
      console.log('Fetching notifications from API:', { page, limit, filter })
      const result = await apiService.get(`/api/notifications?page=${page}&limit=${limit}&filter=${filter}`)
      console.log('API response:', result)

      if (result.success && result.data?.success) {
        const fetchedNotifications = result.data.data?.notifications || []
        notifications.value = fetchedNotifications
        console.log('Loaded notifications:', fetchedNotifications.length)
        // Update unread count based on fetched notifications
        updateUnreadCount()
      } else {
        console.error('API request failed:', result)
      }

      return result
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const fetchUnreadCount = async () => {
    try {
      const result = await apiService.get('/api/notifications/unread-count')

      if (result.success && result.data?.success) {
        unreadCount.value = result.data.data?.count || 0
      }

      return result
    } catch (error) {
      console.error('Failed to fetch unread count:', error)
      throw error
    }
  }

  const markAsReadAPI = async (notificationId) => {
    // Update frontend state immediately for snappy UI
    markAsRead(notificationId)

    // Send to API for persistence in background
    try {
      const result = await apiService.put(`/api/notifications/${notificationId}/read`)
      return result
    } catch (error) {
      console.error('Failed to persist mark as read to server:', error)
      // Optionally revert frontend state if API call fails
      // For now, we'll keep the frontend state as-is since it improves UX
      return { success: false, error: error.message }
    }
  }

  const markAllAsReadAPI = async () => {
    // Update frontend state immediately for snappy UI
    markAllAsRead()

    // Send to API for persistence in background
    try {
      const result = await apiService.put('/api/notifications/mark-all-read')
      return result
    } catch (error) {
      console.error('Failed to persist mark all as read to server:', error)
      // Keep frontend state as-is for better UX
      return { success: false, error: error.message }
    }
  }

  const dismissNotification = async (notificationId) => {
    // Remove from frontend state immediately for snappy UI
    removeNotification(notificationId)

    // Send to API for persistence in background
    try {
      const result = await apiService.put(`/api/notifications/${notificationId}/dismiss`)
      return result
    } catch (error) {
      console.error('Failed to persist notification dismissal to server:', error)
      // Keep frontend state as-is for better UX
      return { success: false, error: error.message }
    }
  }

  const createTestNotification = async () => {
    try {
      const result = await apiService.post('/api/notifications/test')

      if (result.success && result.data?.success) {
        // The notification will be delivered via NATS
        console.log('Test notification created')
      }

      return result
    } catch (error) {
      console.error('Failed to create test notification:', error)
      throw error
    }
  }

  // Helper method to update unread count from current notifications
  const updateUnreadCount = () => {
    unreadCount.value = notifications.value.filter(n => !n.is_read).length
  }

  // NATS Integration - Use centralized addNotification
  const addNotificationFromNATS = (notification) => {
    return addNotification(notification, {
      source: 'nats',
      showAsToast: true,
      persistInList: true,
      autoRemove: true
    })
  }

  // API Integration - Use centralized addNotification
  const addNotificationFromAPI = (notification) => {
    return addNotification(notification, {
      source: 'api',
      showAsToast: false, // API notifications usually don't show as toast
      persistInList: true,
      autoRemove: false
    })
  }

  // Initialize store (fetch initial data)
  const initialize = async () => {
    if (isInitialized.value) return

    console.log('Initializing notification store...')
    try {
      isLoading.value = true
      error.value = null

      console.log('Fetching notifications and unread count...')
      const [notificationsResult] = await Promise.all([
        fetchNotifications(1, 50), // Load more initially
        fetchUnreadCount()
      ])

      console.log('Notifications result:', notificationsResult)

      // Process fetched notifications through centralized manager
      if (notificationsResult?.success && notificationsResult.data?.success) {
        const fetchedNotifications = notificationsResult.data.data?.notifications || []
        console.log('Processing fetched notifications:', fetchedNotifications.length)

        // Clear existing and add from API
        notifications.value = []
        unreadCount.value = 0

        if (Array.isArray(fetchedNotifications)) {
          fetchedNotifications.forEach(notification => {
            addNotificationFromAPI(notification)
          })
        }
      } else {
        console.error('Failed to fetch notifications:', notificationsResult)
      }

      isInitialized.value = true
    } catch (error) {
      console.error('Failed to initialize notification store:', error)
      error.value = error.message || 'Failed to initialize notifications'
    } finally {
      isLoading.value = false
    }
  }

  // Notification types with default styling - All use centralized manager
  const showSuccess = (message, options = {}) => {
    return addNotification({
      type: 'success',
      title: 'Success',
      message,
      icon: 'check-circle',
      color: 'green',
      ...options
    }, {
      source: 'local',
      showAsToast: true,
      persistInList: true,
      autoRemove: true
    })
  }

  const showError = (message, options = {}) => {
    return addNotification({
      type: 'error',
      title: 'Error',
      message,
      icon: 'x-circle',
      color: 'red',
      ...options
    }, {
      source: 'local',
      showAsToast: true,
      persistInList: true,
      autoRemove: false // Keep error notifications until manually dismissed
    })
  }

  const showWarning = (message, options = {}) => {
    return addNotification({
      type: 'warning',
      title: 'Warning',
      message,
      icon: 'exclamation-triangle',
      color: 'yellow',
      ...options
    }, {
      source: 'local',
      showAsToast: true,
      persistInList: true,
      autoRemove: true
    })
  }

  const showInfo = (message, options = {}) => {
    return addNotification({
      type: 'info',
      title: 'Information',
      message,
      icon: 'information-circle',
      color: 'blue',
      ...options
    }, {
      source: 'local',
      showAsToast: true,
      persistInList: true,
      autoRemove: true
    })
  }

  const showAppointmentReminder = (appointment, minutesUntil) => {
    return addNotification({
      type: 'appointment_reminder',
      title: 'Appointment Reminder',
      message: `${appointment.patient_name} has an appointment in ${minutesUntil} minutes`,
      icon: 'clock',
      color: 'blue',
      appointmentId: appointment.id,
      temporary: false,
      actions: [
        {
          label: 'View',
          action: 'view-appointment',
          appointmentId: appointment.id
        },
        {
          label: 'Snooze 5min',
          action: 'snooze-reminder',
          minutes: 5
        }
      ]
    })
  }

  const showAppointmentUpdate = (appointment, updateType) => {
    const messages = {
      'scheduled': `New appointment scheduled for ${appointment.patient_name}`,
      'cancelled': `Appointment with ${appointment.patient_name} has been cancelled`,
      'confirmed': `Appointment with ${appointment.patient_name} has been confirmed`,
      'completed': `Appointment with ${appointment.patient_name} has been completed`,
      'rescheduled': `Appointment with ${appointment.patient_name} has been rescheduled`
    }

    return addNotification({
      type: 'appointment_update',
      title: 'Appointment Update',
      message: messages[updateType] || `Appointment with ${appointment.patient_name} has been updated`,
      icon: 'calendar',
      color: updateType === 'cancelled' ? 'red' : 'blue',
      appointmentId: appointment.id,
      actions: [
        {
          label: 'View',
          action: 'view-appointment',
          appointmentId: appointment.id
        }
      ]
    })
  }

  const showPatientUpdate = (patient, updateType) => {
    const messages = {
      'created': `New patient ${patient.first_name} ${patient.last_name} has been added`,
      'updated': `Patient ${patient.first_name} ${patient.last_name} information has been updated`
    }

    return addNotification({
      type: 'patient_update',
      title: 'Patient Update',
      message: messages[updateType] || `Patient ${patient.first_name} ${patient.last_name} has been updated`,
      icon: 'user',
      color: 'green',
      patientId: patient.id,
      actions: [
        {
          label: 'View',
          action: 'view-patient',
          patientId: patient.id
        }
      ]
    })
  }

  return {
    // State
    notifications,
    unreadCount,
    isLoading,
    isInitialized,
    error,
    toastSettings,

    // Getters (computed)
    recentNotifications,
    hasUnread,
    unreadNotifications,
    readNotifications,
    notificationsByType,
    toastNotifications,

    // Core Management (centralized)
    addNotification,
    removeNotification,
    markAsRead,
    markAllAsRead,
    clearAll,

    // Display Control
    dismissFromToast,
    showInToast,

    // Event System
    subscribe,
    emitNotificationEvent,

    // API Integration
    fetchNotifications,
    fetchUnreadCount,
    markAsReadAPI,
    markAllAsReadAPI,
    dismissNotification,
    createTestNotification,
    initialize,

    // Real-time Integration
    addNotificationFromNATS,
    addNotificationFromAPI,
    updateUnreadCount,

    // Convenience Methods (all use centralized addNotification)
    showSuccess,
    showError,
    showWarning,
    showInfo,
    showAppointmentReminder,
    showAppointmentUpdate,
    showPatientUpdate
  }
})