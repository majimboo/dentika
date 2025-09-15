<template>
  <div class="notifications-page">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6 mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Notifications</h1>
          <p class="text-gray-600 mt-1">Manage your notifications and updates</p>
        </div>
        <div class="flex items-center space-x-2 md:space-x-3">
          <!-- Test Manager button hidden for production -->
          <!-- <button
            @click="showTestManager = !showTestManager"
            class="px-3 py-2 md:px-4 border border-gray-300 text-gray-700 text-sm rounded-lg hover:bg-gray-50 transition-colors"
          >
            {{ showTestManager ? 'Hide' : 'Show' }} Test Manager
          </button> -->
          <button
            v-if="hasUnread"
            @click="markAllAsRead"
            class="px-3 py-2 md:px-4 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors"
          >
            <span class="hidden sm:inline">Mark all as read</span>
            <span class="sm:hidden">Mark read</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Test Manager (hidden for production) -->
    <!-- <div v-if="showTestManager" class="mb-6">
      <NotificationTestManager />
    </div> -->

    <!-- Filter Tabs -->
    <div class="bg-white rounded-lg shadow-sm border border-neutral-200 mb-6">
      <div class="border-b border-gray-200">
        <nav class="flex space-x-4 md:space-x-6 lg:space-x-8 px-4 md:px-6 overflow-x-auto" aria-label="Tabs">
          <button
            v-for="tab in filterTabs"
            :key="tab.key"
            @click="activeFilter = tab.key"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors whitespace-nowrap flex-shrink-0',
              activeFilter === tab.key
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            {{ tab.label }}
            <span
              v-if="tab.count > 0"
              :class="[
                'ml-2 py-0.5 px-2 rounded-full text-xs',
                activeFilter === tab.key
                  ? 'bg-blue-100 text-blue-600'
                  : 'bg-gray-100 text-gray-900'
              ]"
            >
              {{ tab.count }}
            </span>
          </button>
        </nav>
      </div>
    </div>

    <!-- Notifications List -->
    <div class="bg-white rounded-lg shadow-sm border border-neutral-200">
      <!-- Empty State -->
      <div v-if="filteredNotifications.length === 0" class="px-6 py-12 text-center">
        <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 8a6 6 0 0112 0c0 7 3 9 3 9H3s3-2 3-9zM13.73 21a2 2 0 01-3.46 0"></path>
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">
          {{ activeFilter === 'all' ? 'No notifications yet' : `No ${activeFilter} notifications` }}
        </h3>
        <p class="text-gray-500">
          {{ activeFilter === 'all'
            ? "You'll see updates about appointments and patients here"
            : `No ${activeFilter} notifications to display`
          }}
        </p>
      </div>

      <!-- Notifications -->
      <div v-else class="divide-y divide-gray-100">
        <div
          v-for="notification in paginatedNotifications"
          :key="notification.id"
          class="px-6 py-4 hover:bg-gray-50 transition-colors cursor-pointer"
          :class="{ 'bg-blue-50': !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <div class="flex items-start space-x-4">
            <!-- Icon -->
            <div class="flex-shrink-0 mt-1">
              <div :class="getIconClass(notification)" class="w-10 h-10 rounded-full flex items-center justify-center">
                <component :is="getIconComponent(notification)" class="w-5 h-5" />
              </div>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center space-x-2">
                    <h3 class="text-base font-semibold text-gray-900">
                      {{ notification.title }}
                    </h3>
                    <span
                      :class="getTypeBadgeClass(notification.type)"
                      class="px-2 py-1 text-xs font-medium rounded-full"
                    >
                      {{ getTypeLabel(notification.type) }}
                    </span>
                  </div>
                  <p class="text-gray-600 mt-1 leading-5">
                    {{ notification.message }}
                  </p>
                </div>
                <div class="flex items-center space-x-3 ml-4 flex-shrink-0">
                  <span class="text-sm text-gray-500">
                    {{ formatDateTime(notification.timestamp) }}
                  </span>
                  <div
                    v-if="!notification.is_read"
                    class="w-3 h-3 bg-blue-500 rounded-full"
                  ></div>
                </div>
              </div>

              <!-- Action buttons -->
              <div v-if="notification.actions && notification.actions.length > 0" class="mt-3 flex flex-wrap gap-2">
                <button
                  v-for="action in notification.actions"
                  :key="action.label"
                  @click.stop="handleActionClick(action, notification)"
                  class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-md hover:bg-blue-200 transition-colors"
                >
                  {{ action.label }}
                </button>
              </div>

              <!-- Individual actions -->
              <div class="mt-3 flex items-center space-x-4 text-sm">
                <button
                  v-if="!notification.is_read"
                  @click.stop="markAsRead(notification.id)"
                  class="text-blue-600 hover:text-blue-800 transition-colors"
                >
                  Mark as read
                </button>
                <button
                  @click.stop="removeNotification(notification.id)"
                  class="text-red-600 hover:text-red-800 transition-colors"
                >
                  Remove
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Showing {{ startIndex + 1 }} to {{ Math.min(endIndex, filteredNotifications.length) }} of {{ filteredNotifications.length }} notifications
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="previousPage"
              :disabled="currentPage <= 1"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Previous
            </button>
            <span class="px-3 py-1 text-sm text-gray-700">
              Page {{ currentPage }} of {{ totalPages }}
            </span>
            <button
              @click="nextPage"
              :disabled="currentPage >= totalPages"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Clear All Button (Bottom for Safety) -->
    <div v-if="notifications.length > 0" class="mt-6 text-center">
      <button
        @click="clearAll"
        class="px-6 py-3 border border-red-300 text-red-700 bg-white text-sm rounded-lg hover:bg-red-50 hover:border-red-400 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
      >
        <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
        Clear All Notifications
      </button>
      <p class="text-xs text-gray-500 mt-2">This action cannot be undone</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'
import NotificationTestManager from '../components/NotificationTestManager.vue'

const router = useRouter()
const notificationStore = useNotificationStore()

// Test manager toggle
const showTestManager = ref(false)

// State
const activeFilter = ref('all')
const currentPage = ref(1)
const itemsPerPage = 10

// Computed properties
const notifications = computed(() => notificationStore.notifications)
const hasUnread = computed(() => notificationStore.hasUnread)

const filterTabs = computed(() => [
  {
    key: 'all',
    label: 'All',
    count: notifications.value.length
  },
  {
    key: 'unread',
    label: 'Unread',
    count: notifications.value.filter(n => !n.is_read).length
  },
  {
    key: 'appointment_reminder',
    label: 'Reminders',
    count: notifications.value.filter(n => n.type === 'appointment_reminder').length
  },
  {
    key: 'appointment_update',
    label: 'Appointments',
    count: notifications.value.filter(n => n.type === 'appointment_update').length
  },
  {
    key: 'patient_update',
    label: 'Patients',
    count: notifications.value.filter(n => n.type === 'patient_update').length
  }
])

const filteredNotifications = computed(() => {
  let filtered = notifications.value

  if (activeFilter.value === 'unread') {
    filtered = filtered.filter(n => !n.is_read)
  } else if (activeFilter.value !== 'all') {
    filtered = filtered.filter(n => n.type === activeFilter.value)
  }

  return filtered.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
})

const totalPages = computed(() => Math.ceil(filteredNotifications.value.length / itemsPerPage))
const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => startIndex.value + itemsPerPage)

const paginatedNotifications = computed(() => {
  return filteredNotifications.value.slice(startIndex.value, endIndex.value)
})

// Methods
const markAllAsRead = async () => {
  try {
    await notificationStore.markAllAsReadAPI()
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

const clearAll = async () => {
  if (confirm('Are you sure you want to clear all notifications? This action cannot be undone.')) {
    try {
      // Dismiss each notification individually
      const notificationsToRemove = [...notifications.value]
      for (const notification of notificationsToRemove) {
        await notificationStore.dismissNotification(notification.id)
      }
    } catch (error) {
      console.error('Failed to clear all notifications:', error)
    }
  }
}

const markAsRead = async (id) => {
  try {
    await notificationStore.markAsReadAPI(id)
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

const removeNotification = async (id) => {
  try {
    await notificationStore.dismissNotification(id)
  } catch (error) {
    console.error('Failed to remove notification:', error)
  }
}

const handleNotificationClick = async (notification) => {
  if (!notification.is_read) {
    try {
      await notificationStore.markAsReadAPI(notification.id)
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
    }
  }

  // Handle notification-specific routing
  if (notification.appointmentId) {
    router.push(`/appointments/${notification.appointmentId}`)
  } else if (notification.patientId) {
    router.push(`/patients/${notification.patientId}`)
  }
}

const handleActionClick = (action, notification) => {
  switch (action.action) {
    case 'view-appointment':
      if (action.url) {
        router.push(action.url)
      } else {
        console.error('No URL found in action')
      }
      break
    case 'view-patient':
      router.push(`/patients/${action.patientId}`)
      break
    case 'snooze-reminder':
      notificationStore.removeNotification(notification.id)
      break
  }
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// Watch filter changes to reset pagination
const resetPagination = () => {
  currentPage.value = 1
}

// Reset pagination when filter changes
const unwatchFilter = computed(() => activeFilter.value, resetPagination)

// Utility functions
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

const getTypeBadgeClass = (type) => {
  const classes = {
    success: 'bg-green-100 text-green-800',
    error: 'bg-red-100 text-red-800',
    warning: 'bg-yellow-100 text-yellow-800',
    info: 'bg-blue-100 text-blue-800',
    appointment_reminder: 'bg-blue-100 text-blue-800',
    appointment_update: 'bg-purple-100 text-purple-800',
    patient_update: 'bg-green-100 text-green-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getTypeLabel = (type) => {
  const labels = {
    success: 'Success',
    error: 'Error',
    warning: 'Warning',
    info: 'Info',
    appointment_reminder: 'Reminder',
    appointment_update: 'Appointment',
    patient_update: 'Patient'
  }
  return labels[type] || 'Notification'
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
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 002 2v12a2 2 0 002 2z' })
    ]),
    'user': () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' })
    ])
  }

  return iconMap[notification.icon] || iconMap['information-circle']
}

const formatDateTime = (timestamp) => {
  if (!timestamp) return ''

  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  // If less than a day, show relative time
  if (diff < 86400000) {
    if (diff < 60000) return 'Just now'
    if (diff < 3600000) {
      const minutes = Math.floor(diff / 60000)
      return `${minutes}m ago`
    }
    const hours = Math.floor(diff / 3600000)
    return `${hours}h ago`
  }

  // Otherwise show full date and time
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  })
}
</script>

<style scoped>
/* Mobile-friendly tab scrolling */
nav[aria-label="Tabs"] {
  scrollbar-width: thin;
  scrollbar-color: rgba(156, 163, 175, 0.3) transparent;
}

nav[aria-label="Tabs"]::-webkit-scrollbar {
  height: 4px;
}

nav[aria-label="Tabs"]::-webkit-scrollbar-track {
  background: transparent;
}

nav[aria-label="Tabs"]::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.3);
  border-radius: 2px;
}

nav[aria-label="Tabs"]::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.5);
}

/* Ensure tabs are properly visible on mobile */
@media (max-width: 768px) {
  nav[aria-label="Tabs"] {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  nav[aria-label="Tabs"] button {
    padding-left: 0.5rem;
    padding-right: 0.5rem;
    font-size: 0.875rem;
    min-width: fit-content;
  }
}

@media (max-width: 640px) {
  nav[aria-label="Tabs"] {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }

  nav[aria-label="Tabs"] button {
    padding-left: 0.25rem;
    padding-right: 0.25rem;
    font-size: 0.8125rem;
  }
}
</style>