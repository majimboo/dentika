<template>
  <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Notification System Test</h3>

    <div class="space-y-4">
      <!-- Test Buttons -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <button
          @click="createTestNotification"
          :disabled="loading"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'Creating...' : 'Test Notification' }}
        </button>

        <button
          @click="markAllAsRead"
          :disabled="loading"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors"
        >
          Mark All Read
        </button>

        <button
          @click="refreshNotifications"
          :disabled="loading"
          class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50 transition-colors"
        >
          Refresh
        </button>

        <button
          @click="clearAllNotifications"
          :disabled="loading"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors"
        >
          Clear All
        </button>
      </div>

      <!-- Stats -->
      <div class="bg-gray-50 rounded-lg p-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-blue-600">{{ notificationStore.notifications.length }}</div>
            <div class="text-sm text-gray-600">Total Notifications</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-red-600">{{ notificationStore.unreadCount }}</div>
            <div class="text-sm text-gray-600">Unread</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-green-600">{{ readCount }}</div>
            <div class="text-sm text-gray-600">Read</div>
          </div>
        </div>
      </div>

      <!-- Recent Notifications -->
      <div class="space-y-2">
        <h4 class="font-medium text-gray-900">Recent Notifications</h4>
        <div class="max-h-64 overflow-y-auto space-y-2">
          <div
            v-for="notification in notificationStore.recentNotifications"
            :key="notification.id"
            class="flex items-center justify-between p-3 border border-gray-200 rounded-lg"
            :class="{ 'bg-blue-50': !notification.is_read }"
          >
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">{{ notification.title }}</div>
              <div class="text-sm text-gray-600 truncate">{{ notification.message }}</div>
              <div class="text-xs text-gray-400">
                {{ formatTime(notification.created_at) }}
              </div>
            </div>
            <div class="flex items-center space-x-2 ml-3">
              <span
                v-if="!notification.is_read"
                class="w-2 h-2 bg-blue-500 rounded-full"
              ></span>
              <button
                @click="markAsRead(notification.id)"
                v-if="!notification.is_read"
                class="text-xs text-blue-600 hover:text-blue-800"
              >
                Mark Read
              </button>
              <button
                @click="dismissNotification(notification.id)"
                class="text-xs text-red-600 hover:text-red-800"
              >
                Dismiss
              </button>
            </div>
          </div>

          <div v-if="notificationStore.notifications.length === 0" class="text-center py-8 text-gray-500">
            No notifications yet. Click "Test Notification" to create one!
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useNotificationStore } from '../stores/notification'

const notificationStore = useNotificationStore()
const loading = ref(false)

const readCount = computed(() => {
  return notificationStore.notifications.filter(n => n.is_read).length
})

const createTestNotification = async () => {
  loading.value = true
  try {
    await notificationStore.createTestNotification()
    console.log('Test notification created successfully')
  } catch (error) {
    console.error('Failed to create test notification:', error)
    alert('Failed to create test notification')
  } finally {
    loading.value = false
  }
}

const markAllAsRead = async () => {
  loading.value = true
  try {
    await notificationStore.markAllAsReadAPI()
    console.log('All notifications marked as read')
  } catch (error) {
    console.error('Failed to mark all as read:', error)
    alert('Failed to mark all notifications as read')
  } finally {
    loading.value = false
  }
}

const refreshNotifications = async () => {
  loading.value = true
  try {
    await notificationStore.fetchNotifications()
    await notificationStore.fetchUnreadCount()
    console.log('Notifications refreshed')
  } catch (error) {
    console.error('Failed to refresh notifications:', error)
    alert('Failed to refresh notifications')
  } finally {
    loading.value = false
  }
}

const clearAllNotifications = () => {
  if (confirm('Are you sure you want to clear all notifications? This action cannot be undone.')) {
    notificationStore.clearAll()
  }
}

const markAsRead = async (notificationId) => {
  try {
    await notificationStore.markAsReadAPI(notificationId)
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

const dismissNotification = async (notificationId) => {
  try {
    await notificationStore.dismissNotification(notificationId)
  } catch (error) {
    console.error('Failed to dismiss notification:', error)
  }
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''

  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) return 'Just now'
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}m ago`
  }
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}h ago`
  }

  return date.toLocaleDateString()
}
</script>