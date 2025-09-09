import { defineStore } from 'pinia'
import axios from 'axios'

export const useConnectionStore = defineStore('connection', {
  state: () => ({
    isOnline: true,
    isServerReachable: true,
    lastChecked: null,
    retryCount: 0,
    maxRetries: 3,
    checkInterval: null,
    isChecking: false
  }),

  getters: {
    isConnected: (state) => state.isOnline && state.isServerReachable,
    shouldShowOverlay: (state) => !state.isConnected && state.retryCount >= state.maxRetries,
    connectionStatus: (state) => {
      if (!state.isOnline) return 'offline'
      if (!state.isServerReachable) return 'server-down'
      return 'connected'
    }
  },

  actions: {
    async checkServerHealth() {
      if (this.isChecking) return

      this.isChecking = true
      const startTime = Date.now()

      try {
        // Use a simple health check endpoint or the auth/me endpoint
        const response = await axios.get('/api/auth/health', {
          timeout: 5000, // 5 second timeout
          validateStatus: (status) => status < 500 // Accept any status < 500 as success
        })

        this.isServerReachable = true
        this.retryCount = 0
        this.lastChecked = new Date()

        console.log(`Server health check successful (${Date.now() - startTime}ms)`)
      } catch (error) {
        console.warn('Server health check failed:', error.message)

        // Only count as unreachable if it's a network error, not auth errors
        if (!error.response || error.code === 'NETWORK_ERROR' || error.code === 'ECONNREFUSED') {
          this.isServerReachable = false
          this.retryCount = Math.min(this.retryCount + 1, this.maxRetries + 1)
        }

        this.lastChecked = new Date()
      } finally {
        this.isChecking = false
      }
    },

    handleNetworkError() {
      this.isServerReachable = false
      this.retryCount = Math.min(this.retryCount + 1, this.maxRetries + 1)
      this.lastChecked = new Date()

      // Start periodic health checks if not already running
      if (!this.checkInterval) {
        this.startPeriodicChecks()
      }
    },

    handleNetworkRecovery() {
      this.isServerReachable = true
      this.retryCount = 0
      this.lastChecked = new Date()

      // Stop periodic checks once we're back online
      if (this.checkInterval) {
        clearInterval(this.checkInterval)
        this.checkInterval = null
      }
    },

    startPeriodicChecks() {
      // Check every 30 seconds when server is unreachable
      this.checkInterval = setInterval(() => {
        this.checkServerHealth()
      }, 30000)
    },

    stopPeriodicChecks() {
      if (this.checkInterval) {
        clearInterval(this.checkInterval)
        this.checkInterval = null
      }
    },

    // Browser online/offline detection
    setupNetworkListeners() {
      window.addEventListener('online', () => {
        this.isOnline = true
        console.log('Browser is back online')
        // Check server immediately when browser comes back online
        setTimeout(() => this.checkServerHealth(), 1000)
      })

      window.addEventListener('offline', () => {
        this.isOnline = false
        console.log('Browser is offline')
      })

      // Initial check
      this.isOnline = navigator.onLine
    },

    initialize() {
      this.setupNetworkListeners()
      // Initial server health check
      this.checkServerHealth()
    }
  }
})