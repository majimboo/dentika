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
import { useWebSocket } from './composables/useWebSocket'
import NotificationToast from './components/NotificationToast.vue'
import ConnectionOverlay from './components/ConnectionOverlay.vue'
import './api/axios'

const authStore = useAuthStore()
const connectionStore = useConnectionStore()

// Initialize WebSocket for real-time updates
const { connect, disconnect } = useWebSocket()

onMounted(() => {
  // Initialize connection monitoring
  connectionStore.initialize()

  authStore.initializeAuth()

  // Connect WebSocket after auth is initialized
  if (authStore.token) {
    connect()
  }

  // Watch for auth changes to connect/disconnect WebSocket
  authStore.$subscribe((mutation, state) => {
    if (state.token && !state.isConnected) {
      connect()
    } else if (!state.token) {
      disconnect()
    }
  })
})

onUnmounted(() => {
  disconnect()
  connectionStore.stopPeriodicChecks()
})
</script>