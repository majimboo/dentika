import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    loading: false,
    userLoading: false,
    initialized: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isLoading: (state) => state.loading,
    isUserLoading: (state) => state.userLoading,
    isSuperAdmin: (state) => state.user?.role === 'super_admin',
    isClinicOwner: (state) => state.user?.role === 'clinic_owner',
    isDoctor: (state) => state.user?.role === 'doctor',
    isSecretary: (state) => state.user?.role === 'secretary',
    isAssistant: (state) => state.user?.role === 'assistant',
    userRole: (state) => state.user?.role,
    userClinic: (state) => state.user?.clinic,
    userClinicId: (state) => {
      if (state.user?.clinic_id) {
        return state.user.clinic_id
      }
      // If user is clinic_owner but has no clinic_id, default to clinic 1
      if (state.user?.role === 'clinic_owner') {
        return 1
      }
      return null
    },
    canManagePatients: (state) => {
      const role = state.user?.role
      return ['super_admin', 'clinic_owner', 'doctor', 'secretary'].includes(role)
    },
    canManageProcedures: (state) => {
      const role = state.user?.role
      return ['super_admin', 'doctor'].includes(role)
    },
    canCreateClinics: (state) => state.user?.role === 'super_admin',
    hasPermission: (state) => (permission) => {
      if (!state.user) return false

      const role = state.user.role
      const permissions = {
        'appointments.update': ['super_admin', 'clinic_owner', 'doctor', 'secretary'],
        'appointments.create': ['super_admin', 'clinic_owner', 'doctor', 'secretary'],
        'appointments.delete': ['super_admin', 'clinic_owner', 'doctor'],
        'patients.manage': ['super_admin', 'clinic_owner', 'doctor', 'secretary'],
        'procedures.manage': ['super_admin', 'doctor'],
        'clinics.manage': ['super_admin']
      }

      return permissions[permission]?.includes(role) || false
    }
  },

  actions: {
    async login(credentials) {
      this.loading = true
      try {
        console.log('Starting login process...')
        const result = await apiService.login(credentials)
        console.log('Login API result:', result)
        
        if (result.success) {
          console.log('Login successful, setting token:', result.data.token)
          this.token = result.data.token
          apiService.setAuthToken(this.token)
          
          // Use user data from login response (backend already returns user data)
          if (result.data.user) {
            console.log('Setting user data from login response:', result.data.user)
            this.user = result.data.user
          } else {
            console.log('No user data in login response, fetching from server...')
            try {
              await this.fetchCurrentUser()
            } catch (userError) {
              console.warn('Could not fetch user data:', userError)
            }
          }
          
          return { success: true }
        } else {
          console.log('Login failed:', result.error)
          return { 
            success: false, 
            error: result.error 
          }
        }
      } catch (error) {
        console.error('Login error:', error)
        return { 
          success: false, 
          error: 'An unexpected error occurred during login' 
        }
      } finally {
        this.loading = false
      }
    },

    async register(userData) {
      this.loading = true
      try {
        const result = await apiService.register(userData)
        
        if (result.success) {
          this.token = result.data.token
          apiService.setAuthToken(this.token)
          
          // Use user data from register response (backend already returns user data)
          if (result.data.user) {
            this.user = result.data.user
          } else {
            // Fallback to fetching from server
            try {
              await this.fetchCurrentUser()
            } catch (userError) {
              console.warn('Could not fetch user data:', userError)
            }
          }
          
          return { success: true }
        } else {
          return { 
            success: false, 
            error: result.error 
          }
        }
      } catch (error) {
        console.error('Registration error:', error)
        return { 
          success: false, 
          error: 'An unexpected error occurred during registration' 
        }
      } finally {
        this.loading = false
      }
    },

    async logout() {
      this.loading = true
      try {
        await apiService.logout()
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.user = null
        this.token = null
        this.loading = false
      }
    },

    async fetchCurrentUser() {
      if (!this.token) {
        console.log('No token available, clearing user')
        this.user = null
        return
      }

      console.log('Starting fetchCurrentUser with token:', this.token)
      this.userLoading = true
      try {
        const result = await apiService.getCurrentUser()
        console.log('getCurrentUser result:', result)
        
        if (result.success) {
          console.log('User data fetched successfully:', result.data)
          this.user = result.data
        } else {
          console.log('Failed to fetch user data:', result.error)
          throw new Error(`Failed to fetch user: ${result.error}`)
        }
      } catch (error) {
        console.error('Error fetching current user:', error)
        throw error
      } finally {
        this.userLoading = false
      }
    },

    clearAuthState() {
      // Silent logout without API call - for initialization
      this.user = null
      this.token = null
      this.initialized = false
      apiService.clearAuthData()
    },

    async initializeAuth() {
      if (this.initialized) {
        return // Already initialized
      }

      const storedToken = localStorage.getItem('token')

      if (storedToken) {
        this.token = storedToken
        // Always fetch fresh user data from server
        await this.fetchCurrentUser()
      } else {
        // Clear auth state if no token
        this.clearAuthState()
      }

      this.initialized = true
    }
  }
})