<template>
  <div id="app">
    <router-view />
    
    <!-- Real-time notifications -->
    <NotificationToast />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useAuthStore } from './stores/auth'
import { useWebSocket } from './composables/useWebSocket'
import NotificationToast from './components/NotificationToast.vue'
import './api/axios'

const authStore = useAuthStore()

// Initialize WebSocket for real-time updates
const { connect, disconnect } = useWebSocket()

onMounted(() => {
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
})
</script>