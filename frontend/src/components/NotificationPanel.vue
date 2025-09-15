<template>
  <div class="relative">
    <!-- Notification Bell Button -->
    <button
      @click="togglePanel"
      class="p-2 rounded-lg hover:bg-gray-100 transition-colors relative"
      :class="{ 'bg-blue-50': isOpen }"
    >
      <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 8a6 6 0 0112 0c0 7 3 9 3 9H3s3-2 3-9zM13.73 21a2 2 0 01-3.46 0"></path>
      </svg>
      <!-- Notification Badge -->
      <span
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 min-w-[18px] h-[18px] bg-red-500 text-white text-xs rounded-full flex items-center justify-center font-medium"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>


    <!-- Desktop Dropdown Panel -->
    <div
      v-if="isOpen && !isMobile"
      class="absolute top-full right-0 mt-2 w-96 bg-white rounded-lg shadow-xl border border-gray-200 z-50 max-h-96 flex flex-col"
      @click.stop
    >
      <!-- Desktop Header -->
      <div class="px-4 py-3 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900">Notifications</h3>
          <button
            v-if="hasUnread"
            @click="markAllAsRead"
            class="text-sm text-blue-600 hover:text-blue-800 transition-colors"
          >
            Mark all read
          </button>
        </div>
      </div>

      <!-- Desktop Notifications List -->
      <div class="flex-1 overflow-y-auto">
        <div v-if="notifications.length === 0" class="px-4 py-8 text-center text-gray-500">
          <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 8a6 6 0 0112 0c0 7 3 9 3 9H3s3-2 3-9zM13.73 21a2 2 0 01-3.46 0"></path>
          </svg>
          <p class="text-sm">No notifications yet</p>
          <p class="text-xs mt-1 text-gray-400">You'll see updates here</p>
        </div>

        <div v-else>
          <div
            v-for="notification in displayedNotifications"
            :key="notification.id"
            class="px-4 py-3 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 transition-colors cursor-pointer"
            :class="{ 'bg-blue-50': !notification.read }"
            @click="handleNotificationClick(notification)"
          >
            <div class="flex items-start space-x-3">
              <!-- Icon -->
              <div class="flex-shrink-0 mt-0.5">
                <div :class="getIconClass(notification)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <component :is="getIconComponent(notification)" class="w-4 h-4" />
                </div>
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between">
                  <div class="flex-1">
                    <p class="text-sm font-medium text-gray-900">
                      {{ notification.title }}
                    </p>
                    <p class="text-sm text-gray-600 mt-1">
                      {{ notification.message }}
                    </p>
                  </div>
                  <div class="flex items-center space-x-2 ml-2 flex-shrink-0">
                    <span class="text-xs text-gray-500">
                      {{ formatTime(notification.timestamp) }}
                    </span>
                    <div
                      v-if="!notification.read"
                      class="w-2 h-2 bg-blue-500 rounded-full"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Desktop Footer -->
      <div v-if="notifications.length > 0" class="px-4 py-3 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
          <button
            @click="clearAll"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            Clear all
          </button>
          <button
            @click="viewAllNotifications"
            class="text-sm text-blue-600 hover:text-blue-800 transition-colors"
          >
            View all
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

// Component state
const isOpen = ref(false)
const isMobile = ref(window.innerWidth < 1024)

// Computed properties
const notifications = computed(() => notificationStore.recentNotifications)
const unreadCount = computed(() => notificationStore.unreadCount)
const hasUnread = computed(() => notificationStore.hasUnread)

const displayedNotifications = computed(() => {
  // Desktop dropdown shows only 5 most recent
  return notifications.value.slice(0, 5)
})

// Methods
const togglePanel = () => {
  if (isMobile.value) {
    // On mobile, navigate directly to notifications page
    router.push('/notifications')
    return
  }

  // On desktop, toggle the dropdown panel
  isOpen.value = !isOpen.value
}

const closePanel = () => {
  isOpen.value = false
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const clearAll = () => {
  notificationStore.clearAll()
  closePanel()
}

const viewAllNotifications = () => {
  closePanel()
  router.push('/notifications')
}

const handleNotificationClick = (notification) => {
  if (!notification.read) {
    notificationStore.markAsRead(notification.id)
  }

  // Handle notification-specific routing
  if (notification.appointmentId) {
    router.push(`/appointments/${notification.appointmentId}`)
  } else if (notification.patientId) {
    router.push(`/patients/${notification.patientId}`)
  }

  closePanel()
}

const handleActionClick = (action, notification) => {
  // Handle action-specific logic
  switch (action.action) {
    case 'view-appointment':
      router.push(`/appointments/${action.appointmentId}`)
      break
    case 'view-patient':
      router.push(`/patients/${action.patientId}`)
      break
    case 'snooze-reminder':
      // Implementation for snoozing reminders
      notificationStore.removeNotification(notification.id)
      break
  }
  closePanel()
}

const getIconClass = (notification) => {
  const colorClasses = {
    success: 'bg-green-100 text-green-600',
    error: 'bg-red-100 text-red-600',
    warning: 'bg-yellow-100 text-yellow-600',
    info: 'bg-blue-100 text-blue-600',
    appointment_reminder: 'bg-blue-100 text-blue-600',
    appointment_update: 'bg-purple-100 text-purple-600',
    patient_update: 'bg-green-100 text-green-600'
  }
  return colorClasses[notification.type] || 'bg-gray-100 text-gray-600'
}

const getIconComponent = (notification) => {
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
    ])
  }

  return iconMap[notification.icon] || iconMap['information-circle']
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''

  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time

  if (diff < 60000) return 'Just now'
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}m ago`
  }
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}h ago`
  }

  return time.toLocaleDateString()
}

// Handle resize to update mobile state
const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

// Close on outside click (desktop only)
const handleClickOutside = (event) => {
  if (isOpen.value && !isMobile.value && !event.target.closest('.relative')) {
    closePanel()
  }
}

// Add sample notifications
onMounted(() => {
  if (notifications.value.length === 0) {
    notificationStore.showAppointmentUpdate(
      { id: 1, patient_name: 'John Doe' },
      'scheduled'
    )

    notificationStore.showInfo('System maintenance scheduled for tonight', {
      title: 'System Update'
    })
  }

  window.addEventListener('resize', handleResize)
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('click', handleClickOutside)
})
</script>