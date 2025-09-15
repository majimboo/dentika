<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-3 max-w-sm">
      <TransitionGroup name="toast">
        <div
          v-for="notification in visibleNotifications"
          :key="notification.id"
          class="toast-container relative bg-white rounded-xl shadow-xl border-0 p-0 max-w-sm cursor-grab active:cursor-grabbing overflow-hidden backdrop-blur-sm"
          :class="getToastClasses(notification)"
          :style="{ transform: getSwipeTransform(notification.id) }"
          @touchstart="handleTouchStart($event, notification.id)"
          @touchmove="handleTouchMove($event, notification.id)"
          @touchend="handleTouchEnd($event, notification.id)"
        >
          <!-- Color stripe on left -->
          <div class="absolute left-0 top-0 bottom-0 w-1" :class="getColorStripe(notification)"></div>

          <!-- Main content -->
          <div class="p-4 pl-6">
            <div class="flex items-start">
              <!-- Icon -->
              <div class="flex-shrink-0 mr-4">
                <div
                  class="w-10 h-10 rounded-full flex items-center justify-center shadow-sm"
                  :class="getIconClasses(notification)"
                >
                  <component :is="getIconComponent(notification)" class="w-5 h-5" />
                </div>
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0 pr-2">
                <div class="text-sm font-semibold text-gray-900 leading-tight">
                  {{ notification.title }}
                </div>
                <div class="text-sm text-gray-600 mt-1 leading-relaxed">
                  {{ notification.message }}
                </div>

                <!-- Actions -->
                <div v-if="notification.actions && notification.actions.length" class="mt-3 flex flex-wrap gap-2">
                  <button
                    v-for="action in notification.actions"
                    :key="action.label"
                    @click="handleAction(notification, action)"
                    class="text-xs font-medium px-3 py-1 rounded-full transition-all duration-200"
                    :class="getActionButtonClasses(notification)"
                  >
                    {{ action.label }}
                  </button>
                </div>

                <!-- Timestamp -->
                <div class="text-xs text-gray-400 mt-2">
                  {{ formatTime(notification.timestamp) }}
                </div>
              </div>

              <!-- Close Button -->
              <div class="flex-shrink-0">
                <button
                  @click="removeNotification(notification.id)"
                  class="w-8 h-8 rounded-full flex items-center justify-center text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-all duration-200"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Progress bar for auto-dismiss (shows time remaining) -->
          <div v-if="notification.showAsToast && !notification.dismissedFromToast && notification.autoRemove" class="absolute bottom-0 left-0 right-0 h-1 bg-gray-100 overflow-hidden">
            <div
              class="progress-bar h-full"
              :class="getProgressBarClass(notification)"
              :style="{ width: getProgressWidth(notification) + '%' }"
              :title="`${Math.ceil((getProgressWidth(notification) / 100) * (notificationStore.toastSettings.autoRemoveDelay / 1000))}s remaining`"
            ></div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

// Swipe handling
const swipeData = ref({})

// Timer for updating progress bars
const progressTimer = ref(null)
const progressUpdateTrigger = ref(0) // Used to trigger reactivity

// Show only recent notifications that haven't been dismissed
// Use the centralized toast notifications from the store
const visibleNotifications = computed(() => {
  return notificationStore.toastNotifications
})

// Start progress timer to update progress bars in real-time
const startProgressTimer = () => {
  if (progressTimer.value) return // Already running

  progressTimer.value = setInterval(() => {
    const autoRemoveNotifications = visibleNotifications.value.filter(n => n.autoRemove)
    if (autoRemoveNotifications.length > 0) {
      progressUpdateTrigger.value++ // Trigger reactivity
    } else {
      stopProgressTimer() // Stop timer when no auto-remove notifications
    }
  }, 100) // Update every 100ms for smooth animation
}

// Stop progress timer
const stopProgressTimer = () => {
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
    progressTimer.value = null
  }
}

// Get toast styling classes
const getToastClasses = (notification) => {
  const baseClasses = 'transform transition-all duration-300 ease-out hover:shadow-lg'

  const typeClasses = {
    success: 'border-green-50 bg-gradient-to-r from-green-50 to-white',
    error: 'border-red-50 bg-gradient-to-r from-red-50 to-white',
    warning: 'border-yellow-50 bg-gradient-to-r from-yellow-50 to-white',
    info: 'border-blue-50 bg-gradient-to-r from-blue-50 to-white',
    appointment_reminder: 'border-purple-50 bg-gradient-to-r from-purple-50 to-white'
  }

  const notificationType = notification.type || 'info'
  return `${baseClasses} ${typeClasses[notificationType] || typeClasses.info}`
}

// Get color stripe classes
const getColorStripe = (notification) => {
  const stripeClasses = {
    success: 'bg-gradient-to-b from-green-400 to-green-600',
    error: 'bg-gradient-to-b from-red-400 to-red-600',
    warning: 'bg-gradient-to-b from-yellow-400 to-yellow-600',
    info: 'bg-gradient-to-b from-blue-400 to-blue-600',
    appointment_reminder: 'bg-gradient-to-b from-purple-400 to-purple-600'
  }

  const notificationType = notification.type || 'info'
  return stripeClasses[notificationType] || stripeClasses.info
}

// Get action button classes
const getActionButtonClasses = (notification) => {
  const buttonClasses = {
    success: 'bg-green-100 text-green-700 hover:bg-green-200',
    error: 'bg-red-100 text-red-700 hover:bg-red-200',
    warning: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200',
    info: 'bg-blue-100 text-blue-700 hover:bg-blue-200',
    appointment_reminder: 'bg-purple-100 text-purple-700 hover:bg-purple-200'
  }

  const notificationType = notification.type || 'info'
  return buttonClasses[notificationType] || buttonClasses.info
}

// Get progress bar classes
const getProgressBarClass = (notification) => {
  const progressClasses = {
    success: 'bg-green-400',
    error: 'bg-red-400',
    warning: 'bg-yellow-400',
    info: 'bg-blue-400',
    appointment_reminder: 'bg-purple-400'
  }

  const notificationType = notification.type || 'info'
  return progressClasses[notificationType] || progressClasses.info
}

// Get icon styling classes
const getIconClasses = (notification) => {
  const iconClasses = {
    success: 'bg-green-100 text-green-600 border-2 border-green-200',
    error: 'bg-red-100 text-red-600 border-2 border-red-200',
    warning: 'bg-yellow-100 text-yellow-600 border-2 border-yellow-200',
    info: 'bg-blue-100 text-blue-600 border-2 border-blue-200',
    appointment_reminder: 'bg-purple-100 text-purple-600 border-2 border-purple-200'
  }

  const notificationType = notification.type || 'info'
  return iconClasses[notificationType] || iconClasses.info
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

// Get icon component (SVG icons to avoid FontAwesome dependency issues)
const getIconComponent = (notification) => {
  const iconType = notification.icon || notification.type
  const iconMap = {
    'check-circle': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'x-circle': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'exclamation-triangle': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.99-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z' })
    ]),
    'information-circle': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'clock': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'calendar': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z' })
    ]),
    'user': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' })
    ]),
    // Type-based icons
    'success': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'error': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z' })
    ]),
    'warning': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.99-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z' })
    ]),
    'info': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
    ])
  }

  return iconMap[iconType]?.() || iconMap['info']()
}

// Format timestamp for display
const formatTime = (timestamp) => {
  if (!timestamp) return ''

  const now = new Date()
  const time = new Date(timestamp)
  const diffMs = now - time
  const diffSeconds = Math.floor(diffMs / 1000)
  const diffMinutes = Math.floor(diffSeconds / 60)
  const diffHours = Math.floor(diffMinutes / 60)

  if (diffSeconds < 60) {
    return 'just now'
  } else if (diffMinutes < 60) {
    return `${diffMinutes}m ago`
  } else if (diffHours < 24) {
    return `${diffHours}h ago`
  } else {
    return time.toLocaleDateString()
  }
}

// Get progress bar width showing time remaining (reverse progress)
const getProgressWidth = (notification) => {
  // Access the trigger to make this function reactive
  progressUpdateTrigger.value

  // Use the timestamp from when the notification was created
  const createdTime = notification.timestamp || notification.created_at || notification.createdAt
  if (!createdTime) return 100

  const startTime = new Date(createdTime).getTime()
  const currentTime = Date.now()
  const elapsed = currentTime - startTime

  // Get auto-dismiss delay from store settings (default 5 seconds)
  const totalDuration = notificationStore.toastSettings.autoRemoveDelay || 5000

  // Calculate remaining time as percentage (100% = full time remaining, 0% = time up)
  const remainingPercentage = Math.max(0, 100 - (elapsed / totalDuration) * 100)

  return remainingPercentage
}

// Handle notification actions
const handleAction = (notification, action) => {
  console.log('=== ACTION DEBUG ===')
  console.log('Full action object:', JSON.stringify(action, null, 2))
  console.log('action.url:', action.url)
  console.log('action.action:', action.action)
  console.log('action.label:', action.label)
  console.log('Full notification:', JSON.stringify(notification, null, 2))
  console.log('notification.actions:', notification.actions)
  console.log('===================')

  switch (action.action) {
    case 'view-appointment':
      // Use URL if available, otherwise construct from appointmentId
      if (action.url) {
        router.push(action.url)
      } else if (action.appointmentId) {
        router.push(`/appointments/${action.appointmentId}`)
      }
      break
    case 'view-patient':
      // Use URL if available, otherwise construct from patientId
      if (action.url) {
        router.push(action.url)
      } else if (action.patientId) {
        router.push(`/patients/${action.patientId}`)
      }
      break
    case 'navigate':
      // Generic navigation action
      if (action.url) {
        router.push(action.url)
      }
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

// Remove notification from toast (but keep in notification list)
const removeNotification = (id) => {
  notificationStore.dismissFromToast(id)
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
// Watch for changes in visible notifications to manage timer
watch(visibleNotifications, (newNotifications) => {
  const autoRemoveNotifications = newNotifications.filter(n => n.autoRemove)
  if (autoRemoveNotifications.length > 0) {
    startProgressTimer()
  } else {
    stopProgressTimer()
  }
}, { immediate: true })

onMounted(() => {
  // Listen for WebSocket events
  window.addEventListener('appointment-reminder', handleAppointmentReminder)
  window.addEventListener('appointment-updated', handleAppointmentUpdate)
  window.addEventListener('patient-updated', handlePatientUpdate)
  window.addEventListener('app-notification', handleAppNotification)
  window.addEventListener('system-alert', handleSystemAlert)
})

onUnmounted(() => {
  // Clean up event listeners
  window.removeEventListener('appointment-reminder', handleAppointmentReminder)
  window.removeEventListener('appointment-updated', handleAppointmentUpdate)
  window.removeEventListener('patient-updated', handlePatientUpdate)
  window.removeEventListener('app-notification', handleAppNotification)
  window.removeEventListener('system-alert', handleSystemAlert)

  // Clean up timer
  stopProgressTimer()
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

.toast-enter-active {
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.55, 0.055, 0.675, 0.19);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-move {
  transition: transform 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.toast-container {
  max-width: 400px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.toast-container:hover {
  /* Removed translateY lift effect - keeping only shadow enhancement */
  transition: all 0.2s ease-out;
}

/* Close button hover effect */
.toast-container button:hover svg {
  transform: rotate(90deg);
  transition: transform 0.2s ease;
}

/* Progress bar animation - smooth real-time updates */
.toast-container .progress-bar {
  transition: width 0.1s linear;
}

/* Optional: pulse animation when almost done */
.toast-container .progress-bar[style*="width: 0%"],
.toast-container .progress-bar[style*="width: 1%"],
.toast-container .progress-bar[style*="width: 2%"],
.toast-container .progress-bar[style*="width: 3%"],
.toast-container .progress-bar[style*="width: 4%"],
.toast-container .progress-bar[style*="width: 5%"] {
  animation: progressPulse 0.5s ease-in-out infinite alternate;
}

@keyframes progressPulse {
  from { opacity: 0.8; }
  to { opacity: 1; }
}

/* Mobile adjustments */
@media (max-width: 640px) {
  .fixed.top-4.right-4 {
    @apply top-2 right-2 left-2 max-w-none;
  }

  .toast-container {
    @apply max-w-none;
    max-width: calc(100vw - 1rem);
  }

  .toast-container:hover {
    transform: none; /* Disable hover effects on mobile */
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