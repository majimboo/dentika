import axios from 'axios'
import { useConnectionStore } from '../stores/connection'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000'

axios.defaults.baseURL = API_BASE_URL

// Initialize connection store
let connectionStore = null

// Lazy load connection store to avoid circular dependencies
const getConnectionStore = async () => {
  if (!connectionStore) {
    // Import store dynamically to avoid circular dependency
    const module = await import('../stores/connection')
    connectionStore = module.useConnectionStore()
  }
  return connectionStore
}

axios.interceptors.response.use(
  async (response) => {
    // If we get a successful response, mark server as reachable
    if (connectionStore) {
      connectionStore.handleNetworkRecovery()
    } else {
      // Try to get store and mark as recovered
      try {
        const store = await getConnectionStore()
        if (store) store.handleNetworkRecovery()
      } catch (e) {
        // Ignore errors during store loading
      }
    }
    return response
  },
  async (error) => {
    // Handle authentication errors
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
      window.location.href = '/login'
      return Promise.reject(error)
    }

    // Handle network/connection errors
    const isNetworkError = !error.response || error.code === 'NETWORK_ERROR' || error.code === 'ECONNREFUSED' || error.code === 'ENOTFOUND'

    if (isNetworkError) {
      console.warn('Network error detected:', error.message)

      // Update connection status
      if (connectionStore) {
        connectionStore.handleNetworkError()
      } else {
        // Fallback if store not loaded yet
        try {
          const store = await getConnectionStore()
          if (store) store.handleNetworkError()
        } catch (e) {
          // Ignore errors during store loading
        }
      }
    }

    return Promise.reject(error)
  }
)

export default axios