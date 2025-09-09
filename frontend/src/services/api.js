import axios from 'axios'

class ApiService {
  constructor() {
    this.baseURL = '/api'
    this.setupInterceptors()
  }

  setupInterceptors() {
    // Request interceptor to add auth token
    axios.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor to handle authentication errors
    axios.interceptors.response.use(
      (response) => {
        return response
      },
      (error) => {
        const originalRequest = error.config
        
        // Handle 401 errors (except during login attempts)
        if (error.response?.status === 401) {
          // Don't auto-logout if this is a login attempt or logout attempt
          const isLoginRequest = originalRequest.url?.includes('/auth/login')
          const isLogoutRequest = originalRequest.url?.includes('/auth/logout')
          const isOnLoginPage = window.location.pathname === '/login'
          const isOnRegisterPage = window.location.pathname === '/register'
          
          if (!isLoginRequest && !isLogoutRequest && !isOnLoginPage && !isOnRegisterPage) {
            this.handleAuthError()
          }
        }
        
        return Promise.reject(error)
      }
    )
  }

  handleAuthError() {
    // When we get 401, we're already not authenticated on server
    // Just clear frontend auth state and redirect to login
    this.logoutOnUnauthorized()
    
    // Only redirect if we're not already on login/register pages
    const currentPath = window.location.pathname
    if (currentPath !== '/login' && currentPath !== '/register') {
      window.location.href = '/login'
    }
  }

  async request(method, url, data = null, config = {}) {
    try {
      const response = await axios({
        method,
        url: `${this.baseURL}${url}`,
        data,
        ...config
      })
      return {
        success: true,
        data: response.data,
        status: response.status
      }
    } catch (error) {
      console.error(`API ${method.toUpperCase()} ${url} error:`, error)
      
      return {
        success: false,
        error: this.getErrorMessage(error),
        status: error.response?.status,
        data: error.response?.data
      }
    }
  }

  get(url, config) {
    const correctedUrl = url.startsWith('/api') ? url.slice(4) : url
    return this.request('get', correctedUrl, null, config)
  }

  post(url, data, config) {
    const correctedUrl = url.startsWith('/api') ? url.slice(4) : url
    return this.request('post', correctedUrl, data, config)
  }

  put(url, data, config) {
    const correctedUrl = url.startsWith('/api') ? url.slice(4) : url
    return this.request('put', correctedUrl, data, config)
  }

  delete(url, config) {
    const correctedUrl = url.startsWith('/api') ? url.slice(4) : url
    return this.request('delete', correctedUrl, null, config)
  }

  getErrorMessage(error) {
    if (error.response?.data?.message) {
      return error.response.data.message
    }
    
    if (error.response?.data?.error) {
      return error.response.data.error
    }
    
    switch (error.response?.status) {
      case 400:
        return 'Invalid request. Please check your input.'
      case 401:
        return 'Invalid username or password.'
      case 403:
        return 'You do not have permission to perform this action.'
      case 404:
        return 'The requested resource was not found.'
      case 422:
        return 'Validation failed. Please check your input.'
      case 500:
        return 'Server error. Please try again later.'
      default:
        return error.message || 'An unexpected error occurred.'
    }
  }

  // Authentication methods
  async login(credentials) {
    return this.request('post', '/auth/login', credentials)
  }

  async register(userData) {
    return this.request('post', '/auth/register', userData)
  }

  async logout() {
    // Manual logout - call API to invalidate token on server
    let result = { success: true }
    
    const token = localStorage.getItem('token')
    if (token) {
      result = await this.request('post', '/auth/logout')
    }
    
    // Clear local storage regardless of API call result
    this.clearAuthData()
    return result
  }

  logoutOnUnauthorized() {
    // Automatic logout on 401 - don't call API (server already considers us logged out)
    this.clearAuthData()
  }

  // User management methods
  async getUsers() {
    return this.request('get', '/users')
  }

  async getUser(id) {
    return this.request('get', `/users/${id}`)
  }

  async getCurrentUser() {
    return this.request('get', '/auth/me')
  }

  async createUser(userData) {
    return this.request('post', '/users', userData)
  }

  async updateUser(id, userData) {
    return this.request('put', `/users/${id}`, userData)
  }

  async deleteUser(id) {
    return this.request('delete', `/users/${id}`)
  }

  // Utility methods
  isAuthenticated() {
    return !!localStorage.getItem('token')
  }

  setAuthToken(token) {
    localStorage.setItem('token', token)
  }

  clearAuthData() {
    localStorage.removeItem('token')
  }

  // Analytics methods
  async getDashboardMetrics(period = 'today') {
    return this.request('get', `/analytics/dashboard?period=${period}`)
  }

  // Procedure and diagnosis template methods
  async getProcedureTemplates(params = {}) {
    const queryParams = new URLSearchParams(params)
    return this.request('get', `/procedure-templates?${queryParams}`)
  }

  async createProcedureTemplate(templateData) {
    return this.request('post', '/procedure-templates', templateData)
  }

  async getDiagnosisTemplates(params = {}) {
    const queryParams = new URLSearchParams(params)
    return this.request('get', `/diagnosis-templates?${queryParams}`)
  }

  async createDiagnosisTemplate(templateData) {
    return this.request('post', '/diagnosis-templates', templateData)
  }

  // Avatar upload methods
  async uploadAvatar(formData) {
    return this.request('post', '/upload/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }

  async deleteAvatar(avatarPath) {
    return this.request('delete', '/upload/avatar', null, {
      params: { path: avatarPath }
    })
  }
}

// Create and export a single instance
const apiService = new ApiService()
export default apiService