<template>
  <div class="p-6 bg-white rounded-lg shadow-lg max-w-2xl">
    <h2 class="text-xl font-bold mb-4">Notification System Test</h2>

    <!-- System Status -->
    <div class="mb-6 p-4 bg-gray-50 rounded">
      <h3 class="font-semibold mb-2">System Status</h3>
      <div class="space-y-2 text-sm">
        <div>Total Notifications: {{ notifications.length }}</div>
        <div>Unread Count: {{ unreadCount }}</div>
        <div>Toast Notifications: {{ toastNotifications.length }}</div>
        <div>Store Initialized: {{ isInitialized ? 'Yes' : 'No' }}</div>
        <div>Loading: {{ isLoading ? 'Yes' : 'No' }}</div>
        <div v-if="error" class="text-red-600">Error: {{ error }}</div>
      </div>
    </div>

    <!-- Test Buttons -->
    <div class="space-y-3">
      <h3 class="font-semibold">Test Actions</h3>

      <div class="grid grid-cols-2 gap-3">
        <button
          @click="testSuccess"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          Test Success
        </button>

        <button
          @click="testError"
          class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        >
          Test Error
        </button>

        <button
          @click="testWarning"
          class="px-4 py-2 bg-yellow-600 text-white rounded hover:bg-yellow-700"
        >
          Test Warning
        </button>

        <button
          @click="testInfo"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Test Info
        </button>

        <button
          @click="testAppointmentReminder"
          class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700"
        >
          Test Appointment
        </button>

        <button
          @click="testNATSNotification"
          class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
        >
          Test NATS
        </button>

        <button
          @click="testAPINotification"
          class="px-4 py-2 bg-cyan-600 text-white rounded hover:bg-cyan-700"
        >
          Test API
        </button>

        <button
          @click="clearAll"
          class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
        >
          Clear All
        </button>
      </div>
    </div>

    <!-- Toast Debug Info -->
    <div class="mt-6 p-4 bg-yellow-50 rounded">
      <h3 class="font-semibold mb-2">Toast Debug Info</h3>
      <div class="text-sm space-y-1">
        <div>Toast notifications count: {{ toastNotifications.length }}</div>
        <div v-if="toastNotifications.length > 0">
          <div class="font-medium">Toast notifications:</div>
          <div v-for="toast in toastNotifications" :key="toast.id" class="ml-2 text-xs">
            - {{ toast.title }} (showAsToast: {{ toast.showAsToast }}, dismissedFromToast: {{ toast.dismissedFromToast }})
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Notifications Preview -->
    <div class="mt-6">
      <h3 class="font-semibold mb-2">Recent Notifications (Preview)</h3>
      <div class="space-y-2 max-h-60 overflow-y-auto">
        <div
          v-for="notification in recentNotifications"
          :key="notification.id"
          class="p-3 border rounded text-sm"
          :class="{
            'bg-green-50 border-green-200': notification.type === 'success',
            'bg-red-50 border-red-200': notification.type === 'error',
            'bg-yellow-50 border-yellow-200': notification.type === 'warning',
            'bg-blue-50 border-blue-200': notification.type === 'info',
            'bg-purple-50 border-purple-200': notification.type === 'appointment_reminder',
            'border-gray-300': !notification.is_read && !notification.read,
            'opacity-60': notification.is_read || notification.read
          }"
        >
          <div class="font-medium">{{ notification.title }}</div>
          <div class="text-gray-600">{{ notification.message }}</div>
          <div class="text-xs text-gray-500 mt-1">
            {{ formatDate(notification.timestamp || notification.created_at) }}
            | Source: {{ notification.source }}
            | Toast: {{ notification.showAsToast ? 'Yes' : 'No' }}
            | Read: {{ notification.is_read || notification.read ? 'Yes' : 'No' }}
          </div>
        </div>

        <div v-if="recentNotifications.length === 0" class="text-gray-500 text-center py-4">
          No notifications yet. Try the test buttons above!
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '../stores/notification.js'

const notificationStore = useNotificationStore()

// Get reactive references from store
const {
  notifications,
  unreadCount,
  isLoading,
  isInitialized,
  error,
  recentNotifications,
  toastNotifications
} = storeToRefs(notificationStore)

// Test functions
const testSuccess = () => {
  notificationStore.showSuccess('Operation completed successfully!', {
    title: 'Success Test'
  })
}

const testError = () => {
  notificationStore.showError('Something went wrong. Please try again.', {
    title: 'Error Test'
  })
}

const testWarning = () => {
  notificationStore.showWarning('This action cannot be undone.', {
    title: 'Warning Test'
  })
}

const testInfo = () => {
  notificationStore.showInfo('Here is some useful information.', {
    title: 'Info Test'
  })
}

const testAppointmentReminder = () => {
  notificationStore.showAppointmentReminder({
    id: Date.now(),
    patient_name: 'John Doe Test'
  }, 15)
}

const testNATSNotification = () => {
  // Simulate a NATS notification
  notificationStore.addNotificationFromNATS({
    id: Date.now(),
    title: 'NATS Test Notification',
    message: 'This notification came through NATS',
    type: 'info',
    created_at: new Date().toISOString()
  })
}

const testAPINotification = () => {
  // Simulate an API notification (but show as toast for testing)
  notificationStore.addNotification({
    id: Date.now(),
    title: 'API Test Notification',
    message: 'This notification came from the API',
    type: 'info',
    created_at: new Date().toISOString(),
    is_read: false
  }, {
    source: 'api',
    showAsToast: true, // Show as toast for testing purposes
    persistInList: true,
    autoRemove: true
  })
}

const clearAll = () => {
  notificationStore.clearAll()
}

// Helper function
const formatDate = (date) => {
  if (!date) return 'Unknown'
  return new Date(date).toLocaleTimeString()
}
</script>