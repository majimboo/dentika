<template>
  <!-- Notification Panel -->
  <div class="relative notification-panel">
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

    <!-- Notification Panel - Mobile: full-screen, Desktop: dropdown -->
    <div
      v-if="isOpen"
      class="fixed inset-0 lg:top-0 lg:right-0 lg:bottom-0 lg:w-96 bg-white lg:shadow-xl lg:border-l lg:border-gray-200 flex flex-col z-50 overflow-hidden"
      @click.stop
    >
      <!-- Panel Header - Mobile: matches nav height, Desktop: sidebar header -->
      <div class="h-14 lg:h-16 px-4 py-0 lg:py-4 border-b border-gray-200 flex items-center justify-between bg-white">
        <h3 class="text-lg lg:text-xl font-semibold text-gray-900">Notifications</h3>
        <div class="flex items-center space-x-2">
          <button
            v-if="hasUnread"
            @click="markAllAsRead"
            class="text-xs text-blue-600 hover:text-blue-800 transition-colors"
          >
            Mark all read
          </button>
          <button
            @click="closePanel"
            class="p-1 rounded-md hover:bg-gray-100 transition-colors"
          >
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Notifications List -->
      <div class="flex-1 overflow-y-auto bg-white">
        <div v-if="notifications.length === 0" class="h-screen lg:min-h-full flex flex-col items-center justify-center p-8 lg:p-4 text-center text-gray-500 text-sm bg-white">
          <svg class="w-16 lg:w-12 h-16 lg:h-12 mx-auto mb-4 lg:mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 8a6 6 0 0112 0c0 7 3 9 3 9H3s3-2 3-9zM13.73 21a2 2 0 01-3.46 0"></path>
          </svg>
          <p>No notifications yet</p>
          <p class="text-xs mt-1">You'll see updates about appointments and patients here</p>
        </div>

        <div v-else class="py-1">
          <div
            v-for="notification in displayedNotifications"
            :key="notification.id"
            class="notification-item px-4 py-3 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 transition-colors cursor-pointer"
            :class="{ 'bg-blue-50/50': !notification.read }"
            @click="handleNotificationClick(notification)"
          >
            <div class="flex items-start space-x-3">
              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between">
                  <p class="text-sm font-medium text-gray-900 truncate">
                    {{ notification.title }}
                  </p>
                  <div class="flex items-center space-x-2 ml-2">
                    <span class="text-xs text-gray-500">
                      {{ formatTime(notification.timestamp) }}
                    </span>
                    <div
                      v-if="!notification.read"
                      class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0"
                    ></div>
                  </div>
                </div>
                <p class="text-sm text-gray-600 mt-1">
                  {{ notification.message }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Panel Footer -->
      <div v-if="notifications.length > 0" class="px-4 py-3 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
          <button
            @click="clearAll"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            Clear all
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()

// Component state
const isOpen = ref(false)

// Computed properties
const notifications = computed(() => notificationStore.recentNotifications)
const unreadCount = computed(() => notificationStore.unreadCount)
const hasUnread = computed(() => notificationStore.hasUnread)

const displayedNotifications = computed(() => {
  return notifications.value.slice(0, 5)
})

// Methods
const togglePanel = () => {
  isOpen.value = !isOpen.value
  // Prevent background scrolling on mobile
  if (isOpen.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const closePanel = () => {
  isOpen.value = false
  // Re-enable background scrolling
  document.body.style.overflow = ''
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const clearAll = () => {
  notificationStore.clearAll()
  closePanel()
}

const handleNotificationClick = (notification) => {
  if (!notification.read) {
    notificationStore.markAsRead(notification.id)
  }
  closePanel()
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
})

// Close on outside click
const handleClickOutside = (event) => {
  if (isOpen.value && !event.target.closest('.notification-panel')) {
    closePanel()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  // Ensure body scrolling is restored if component unmounts while open
  document.body.style.overflow = ''
})
</script>

<style scoped>
@reference "tailwindcss";

.notification-item {
  transition: all 0.2s ease-in-out;
}
</style>