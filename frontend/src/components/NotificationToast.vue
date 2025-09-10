<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
      <TransitionGroup name="toast">
        <div
          v-for="notification in visibleNotifications"
          :key="notification.id"
          class="toast-container bg-white rounded-lg shadow-lg border-l-4 p-4 max-w-sm cursor-grab active:cursor-grabbing"
          :class="getToastClasses(notification)"
          :style="{ transform: getSwipeTransform(notification.id) }"
          @touchstart="handleTouchStart($event, notification.id)"
          @touchmove="handleTouchMove($event, notification.id)"
          @touchend="handleTouchEnd($event, notification.id)"
        >
          <div class="flex items-start">
            <!-- Icon -->
            <div class="flex-shrink-0 mr-3">
              <div
                class="w-6 h-6 rounded-full flex items-center justify-center"
                :class="getIconClasses(notification)"
              >
                <font-awesome-icon
                  :icon="getIconName(notification.icon)"
                  class="w-4 h-4"
                />
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
              <font-awesome-icon icon="fa-solid fa-times" class="w-4 h-4" />
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
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

// Swipe handling
const swipeData = ref({})

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

// Get FontAwesome icon name
const getIconName = (iconType) => {
  const iconMap = {
    'check-circle': 'fa-solid fa-check-circle',
    'x-circle': 'fa-solid fa-times-circle',
    'exclamation-triangle': 'fa-solid fa-exclamation-triangle',
    'information-circle': 'fa-solid fa-info-circle',
    'clock': 'fa-solid fa-clock',
    'calendar': 'fa-solid fa-calendar',
    'user': 'fa-solid fa-user'
  }

  return iconMap[iconType] || 'fa-solid fa-info-circle'
}

// Get progress bar width for temporary notifications
const getProgressWidth = (notification) => {
  if (!notification.createdAt) return 100

  const elapsed = Date.now() - notification.createdAt
  const total = 5000 // 5 seconds
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

// Swipe handling methods
const handleTouchStart = (event, notificationId) => {
  const touch = event.touches[0]
  swipeData.value[notificationId] = {
    startX: touch.clientX,
    startY: touch.clientY,
    currentX: touch.clientX,
    isSwiping: false
  }
}

const handleTouchMove = (event, notificationId) => {
  if (!swipeData.value[notificationId]) return

  event.preventDefault()
  const touch = event.touches[0]
  const data = swipeData.value[notificationId]

  data.currentX = touch.clientX
  const deltaX = data.currentX - data.startX
  const deltaY = Math.abs(touch.clientY - data.startY)

  // Only consider it a swipe if horizontal movement is greater than vertical
  if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 10) {
    data.isSwiping = true
  }
}

const handleTouchEnd = (event, notificationId) => {
  if (!swipeData.value[notificationId]) return

  const data = swipeData.value[notificationId]
  const deltaX = data.currentX - data.startX

  // If swiped to the right with sufficient distance, dismiss the notification
  if (data.isSwiping && deltaX > 100) {
    removeNotification(notificationId)
  }

  // Reset swipe data
  delete swipeData.value[notificationId]
}

const getSwipeTransform = (notificationId) => {
  const data = swipeData.value[notificationId]
  if (!data || !data.isSwiping) return ''

  const deltaX = data.currentX - data.startX
  // Only allow rightward movement
  const translateX = Math.max(0, deltaX)
  return `translateX(${translateX}px)`
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
@import "../styles/main.css";

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

/* Swipe animations */
.toast-container {
  transition: transform 0.2s ease-out;
  will-change: transform;
}

.toast-container:active {
  transition: none;
}

/* Swipe feedback */
.toast-container.swipe-active {
  opacity: 0.8;
}
</style>