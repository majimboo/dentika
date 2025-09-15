<template>
  <div id="app">
    <router-view />

    <!-- Real-time notifications -->
    <NotificationToast />

    <!-- Connection status overlay -->
    <ConnectionOverlay />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useAuthStore } from './stores/auth'
import { useConnectionStore } from './stores/connection'
import { useNotificationStore } from './stores/notification'
import { useNats } from './composables/useNats'
import NotificationToast from './components/NotificationToast.vue'
import ConnectionOverlay from './components/ConnectionOverlay.vue'
import './api/axios'

const authStore = useAuthStore()
const connectionStore = useConnectionStore()
const notificationStore = useNotificationStore()

// Initialize NATS for real-time updates
const { connect, disconnect } = useNats()

onMounted(async () => {
  // Initialize connection monitoring
  connectionStore.initialize()

  // Auth is initialized by router guards before component mounting
  // No redundant initialization needed here

  // Initialize notifications after auth is ready
  if (authStore.token) {
    await notificationStore.initialize()
    connect()
  }

  // Watch for auth changes to connect/disconnect NATS and initialize notifications
  authStore.$subscribe(async (mutation, state) => {
    if (state.token && !state.isConnected) {
      await notificationStore.initialize()
      connect()
    } else if (!state.token) {
      disconnect()
      // Clear notifications when logged out
      notificationStore.clearAll()
    }
  })
})

onUnmounted(() => {
  disconnect()
  connectionStore.stopPeriodicChecks()
})
</script>