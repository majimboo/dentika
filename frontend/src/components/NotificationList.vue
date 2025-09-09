<template>
  <div class="space-y-3">
    <div v-if="notifications.length === 0" class="text-center py-8">
      <svg class="w-12 h-12 text-neutral-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
      </svg>
      <p class="text-neutral-500">No notifications yet</p>
    </div>

    <div v-else class="space-y-2">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="p-3 rounded-lg border border-neutral-200 hover:bg-neutral-50 transition-colors cursor-pointer"
        :class="{ 'bg-blue-50 border-blue-200': !notification.read }"
        @click="markAsRead(notification.id)"
      >
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0 mt-0.5">
            <div 
              class="w-2 h-2 rounded-full" 
              :class="notification.read ? 'bg-neutral-300' : 'bg-blue-500'"
            ></div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-2 mb-1">
              <h4 class="text-sm font-medium text-neutral-900 truncate">{{ notification.title }}</h4>
              <span class="text-xs text-neutral-500 flex-shrink-0">{{ formatTime(notification.createdAt) }}</span>
            </div>
            <p class="text-sm text-neutral-600 line-clamp-2">{{ notification.message }}</p>
            <div v-if="notification.type" class="mt-2">
              <span 
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium"
                :class="getTypeClasses(notification.type)"
              >
                {{ notification.type }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mark all as read button -->
    <div v-if="unreadCount > 0" class="pt-3 border-t border-neutral-200">
      <button 
        @click="markAllAsRead"
        class="w-full p-2 text-sm text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
      >
        Mark all as read ({{ unreadCount }})
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  name: 'NotificationList',
  setup() {
    // Mock notifications - replace with actual data
    const notifications = ref([
      {
        id: 1,
        title: 'System Update',
        message: 'The system has been successfully updated to version 2.1.0. New features include improved performance and bug fixes.',
        type: 'info',
        read: false,
        createdAt: new Date(Date.now() - 5 * 60 * 1000) // 5 minutes ago
      },
      {
        id: 2,
        title: 'New User Registered',
        message: 'John Doe has registered and is pending approval.',
        type: 'success',
        read: false,
        createdAt: new Date(Date.now() - 30 * 60 * 1000) // 30 minutes ago
      },
      {
        id: 3,
        title: 'Storage Warning',
        message: 'Your storage usage is at 85%. Consider cleaning up old files.',
        type: 'warning',
        read: true,
        createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000) // 2 hours ago
      },
      {
        id: 4,
        title: 'Failed Login Attempt',
        message: 'Multiple failed login attempts detected from IP 192.168.1.100.',
        type: 'error',
        read: true,
        createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000) // 24 hours ago
      }
    ])

    const unreadCount = computed(() => {
      return notifications.value.filter(n => !n.read).length
    })

    const markAsRead = (id) => {
      const notification = notifications.value.find(n => n.id === id)
      if (notification) {
        notification.read = true
      }
    }

    const markAllAsRead = () => {
      notifications.value.forEach(n => n.read = true)
    }

    const formatTime = (date) => {
      const now = new Date()
      const diff = now - date
      const minutes = Math.floor(diff / (1000 * 60))
      const hours = Math.floor(diff / (1000 * 60 * 60))
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))

      if (minutes < 1) return 'Just now'
      if (minutes < 60) return `${minutes}m`
      if (hours < 24) return `${hours}h`
      return `${days}d`
    }

    const getTypeClasses = (type) => {
      const classes = {
        info: 'bg-blue-100 text-blue-800',
        success: 'bg-green-100 text-green-800',
        warning: 'bg-yellow-100 text-yellow-800',
        error: 'bg-red-100 text-red-800'
      }
      return classes[type] || 'bg-neutral-100 text-neutral-800'
    }

    return {
      notifications,
      unreadCount,
      markAsRead,
      markAllAsRead,
      formatTime,
      getTypeClasses
    }
  }
}
</script>