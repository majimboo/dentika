import axios from 'axios'
import { useConnectionStore } from '../stores/connection'

axios.interceptors.response.use(
  (response) => {
    // If we get a successful response, mark server as reachable
    const connectionStore = useConnectionStore()
    connectionStore.handleNetworkRecovery()
    return response
  },
  (error) => {
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
      const connectionStore = useConnectionStore()
      connectionStore.handleNetworkError()
    }

    return Promise.reject(error)
  }
)

export default axios
