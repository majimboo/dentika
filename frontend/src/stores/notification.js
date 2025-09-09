import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  // State
  const notifications = ref([])
  const unreadCount = ref(0)

  // Getters
  const recentNotifications = computed(() => {
    return notifications.value.slice(0, 10) // Show last 10 notifications
  })

  const hasUnread = computed(() => unreadCount.value > 0)

  // Actions
  const addNotification = (notification) => {
    const newNotification = {
      id: Date.now() + Math.random(),
      timestamp: new Date(),
      read: false,
      ...notification
    }

    notifications.value.unshift(newNotification)
    unreadCount.value++

    // Auto-remove after 30 seconds if it's a temporary notification
    if (notification.temporary !== false) {
      setTimeout(() => {
        removeNotification(newNotification.id)
      }, 30000)
    }

    return newNotification
  }

  const removeNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      const notification = notifications.value[index]
      if (!notification.read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      notifications.value.splice(index, 1)
    }
  }

  const markAsRead = (id) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification && !notification.read) {
      notification.read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  }

  const markAllAsRead = () => {
    notifications.value.forEach(notification => {
      notification.read = true
    })
    unreadCount.value = 0
  }

  const clearAll = () => {
    notifications.value = []
    unreadCount.value = 0
  }

  // Notification types with default styling
  const showSuccess = (message, options = {}) => {
    return addNotification({
      type: 'success',
      title: 'Success',
      message,
      icon: 'check-circle',
      color: 'green',
      ...options
    })
  }

  const showError = (message, options = {}) => {
    return addNotification({
      type: 'error',
      title: 'Error',
      message,
      icon: 'x-circle',
      color: 'red',
      temporary: false, // Keep error notifications until manually dismissed
      ...options
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
    
    // Getters
    recentNotifications,
    hasUnread,
    
    // Actions
    addNotification,
    removeNotification,
    markAsRead,
    markAllAsRead,
    clearAll,
    
    // Convenience methods
    showSuccess,
    showError,
    showWarning,
    showInfo,
    showAppointmentReminder,
    showAppointmentUpdate,
    showPatientUpdate
  }
})