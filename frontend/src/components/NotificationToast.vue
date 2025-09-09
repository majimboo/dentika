<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
      <TransitionGroup name="toast" tag="div">
        <div
          v-for="notification in visibleNotifications"
          :key="notification.id"
          class="toast-container bg-white rounded-lg shadow-lg border-l-4 p-4 max-w-sm"
          :class="getToastClasses(notification)"
        >
          <div class="flex items-start">
            <!-- Icon -->
            <div class="flex-shrink-0 mr-3">
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center"
                :class="getIconClasses(notification)"
              >
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path v-if="notification.icon === 'check-circle'" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"></path>
                  <path v-else-if="notification.icon === 'x-circle'" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"></path>
                  <path v-else-if="notification.icon === 'exclamation-triangle'" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"></path>
                  <path v-else-if="notification.icon === 'information-circle'" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"></path>
                  <path v-else-if="notification.icon === 'clock'" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"></path>
                  <path v-else-if="notification.icon === 'calendar'" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zM4 9h12v8H4V9z"></path>
                  <path v-else-if="notification.icon === 'user'" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"></path>
                </svg>
              </div>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-gray-900">
                {{ notification.title }}
              </div>
              <div class="text-sm text-gray-600 mt-1">
                {{ notification.message }}
              </div>
              
              <!-- Actions -->
              <div v-if="notification.actions && notification.actions.length" class="mt-3 flex space-x-2">
                <button
                  v-for="action in notification.actions"
                  :key="action.label"
                  @click="handleAction(notification, action)"
                  class="text-xs font-medium text-blue-600 hover:text-blue-800 transition-colors"
                >
                  {{ action.label }}
                </button>
              </div>
            </div>

            <!-- Close button -->
            <button
              @click="removeNotification(notification.id)"
              class="flex-shrink-0 ml-2 text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"></path>
              </svg>
            </button>
          </div>

          <!-- Progress bar for temporary notifications -->
          <div
            v-if="notification.temporary !== false"
            class="mt-3 w-full bg-gray-200 rounded-full h-1"
          >
            <div
              class="h-1 rounded-full transition-all duration-300 ease-linear"
              :class="`bg-${notification.color}-500`"
              :style="{ width: getProgressWidth(notification) + '%' }"
            ></div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

// Show only recent notifications that haven't been dismissed
const visibleNotifications = computed(() => {
  return notificationStore.recentNotifications.filter(n => !n.dismissed)
})

// Get toast styling classes
const getToastClasses = (notification) => {
  const colorClasses = {
    green: 'border-l-green-500',
    red: 'border-l-red-500',
    yellow: 'border-l-yellow-500',
    blue: 'border-l-blue-500'
  }
  
  return colorClasses[notification.color] || 'border-l-gray-500'
}

// Get icon styling classes
const getIconClasses = (notification) => {
  const colorClasses = {
    green: 'bg-green-100 text-green-600',
    red: 'bg-red-100 text-red-600',
    yellow: 'bg-yellow-100 text-yellow-600',
    blue: 'bg-blue-100 text-blue-600'
  }
  
  return colorClasses[notification.color] || 'bg-gray-100 text-gray-600'
}

// Get progress bar width for temporary notifications
const getProgressWidth = (notification) => {
  if (!notification.createdAt) return 100
  
  const elapsed = Date.now() - notification.createdAt
  const total = 30000 // 30 seconds
  return Math.max(0, 100 - (elapsed / total) * 100)
}

// Handle notification actions
const handleAction = (notification, action) => {
  switch (action.action) {
    case 'view-appointment':
      router.push(`/appointments/${action.appointmentId}`)
      break
    case 'view-patient':
      router.push(`/patients/${action.patientId}`)
      break
    case 'snooze-reminder':
      // Implement snooze logic here
      console.log(`Snoozing reminder for ${action.minutes} minutes`)
      break
    default:
      console.log('Unknown action:', action)
  }
  
  // Mark notification as read and remove from view
  notificationStore.markAsRead(notification.id)
  removeNotification(notification.id)
}

// Remove notification
const removeNotification = (id) => {
  notificationStore.removeNotification(id)
}

// Set up event listeners for WebSocket events
onMounted(() => {
  // Listen for WebSocket events
  window.addEventListener('appointment-reminder', handleAppointmentReminder)
  window.addEventListener('appointment-updated', handleAppointmentUpdate)
  window.addEventListener('patient-updated', handlePatientUpdate)
  window.addEventListener('app-notification', handleAppNotification)
  window.addEventListener('system-alert', handleSystemAlert)
})

onUnmounted(() => {
  window.removeEventListener('appointment-reminder', handleAppointmentReminder)
  window.removeEventListener('appointment-updated', handleAppointmentUpdate)
  window.removeEventListener('patient-updated', handlePatientUpdate)
  window.removeEventListener('app-notification', handleAppNotification)
  window.removeEventListener('system-alert', handleSystemAlert)
})

// Event handlers
const handleAppointmentReminder = (event) => {
  const reminder = event.detail
  notificationStore.showAppointmentReminder(reminder.appointment, reminder.minutes_until)
}

const handleAppointmentUpdate = (event) => {
  const appointment = event.detail
  notificationStore.showAppointmentUpdate(appointment, appointment.status)
}

const handlePatientUpdate = (event) => {
  const patient = event.detail
  notificationStore.showPatientUpdate(patient, 'updated')
}

const handleAppNotification = (event) => {
  const notification = event.detail
  
  switch (notification.type) {
    case 'success':
      notificationStore.showSuccess(notification.message, notification.options)
      break
    case 'error':
      notificationStore.showError(notification.message, notification.options)
      break
    case 'warning':
      notificationStore.showWarning(notification.message, notification.options)
      break
    case 'info':
      notificationStore.showInfo(notification.message, notification.options)
      break
  }
}

const handleSystemAlert = (event) => {
  const alert = event.detail
  notificationStore.showError(alert.message, {
    title: 'System Alert',
    temporary: false
  })
}
</script>

<style scoped>
@reference "tailwindcss";

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}

.toast-container {
  max-width: 384px; /* w-96 */
}

/* Mobile adjustments */
@media (max-width: 640px) {
  .fixed.top-4.right-4 {
    @apply top-2 right-2 left-2 max-w-none;
  }
  
  .toast-container {
    @apply max-w-none;
  }
}
</style>