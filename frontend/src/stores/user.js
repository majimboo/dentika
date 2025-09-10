import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    doctors: [],
    loading: false,
    error: null
  }),

  getters: {
    isLoading: (state) => state.loading,
    getActiveUsers: (state) => state.users.filter(user => user.is_active),
    getDoctors: (state) => state.doctors.filter(doctor => doctor.is_active),
    getUserById: (state) => (id) => state.users.find(user => user.id === id),
    getDoctorById: (state) => (id) => state.doctors.find(doctor => doctor.id === id),
    getUsersByRole: (state) => (role) =>
      state.users.filter(user => user.role === role && user.is_active)
  },

  actions: {
    async fetchUsers(params = {}) {
      this.loading = true
      this.error = null

      try {
        const queryParams = new URLSearchParams(params)
        const result = await apiService.get('/api/users?' + queryParams.toString())

        if (result.success) {
          this.users = result.data || []
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          console.warn('Failed to fetch users:', result.error)
          // Don't fail completely, just return empty array
          this.users = []
          return { success: true, data: [] }
        }
      } catch (error) {
        console.error('Error fetching users:', error)
        this.error = 'Failed to fetch users'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchDoctors() {
      this.loading = true
      this.error = null

      try {
        const result = await this.fetchUsers({ role: 'doctor' })

        if (result.success) {
          this.doctors = result.data || []
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          console.warn('Failed to fetch doctors:', result.error)
          // Don't fail completely, just return empty array
          this.doctors = []
          return { success: true, data: [] }
        }
      } catch (error) {
        console.error('Error fetching doctors:', error)
        this.error = 'Failed to fetch doctors'
        // Don't fail completely, just return empty array
        this.doctors = []
        return { success: true, data: [] }
      } finally {
        this.loading = false
      }
    },

    async createUser(userData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/users', userData)

        if (result.success) {
          // Add to local state
          this.users.push(result.data)
          // If it's a doctor, also add to doctors list
          if (result.data.role === 'doctor') {
            this.doctors.push(result.data)
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        console.error('Error creating user:', error)
        this.error = 'Failed to create user'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async updateUser(userId, userData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.put(`/api/users/${userId}`, userData)

        if (result.success) {
          // Update in local state
          const index = this.users.findIndex(user => user.id === userId)
          if (index !== -1) {
            this.users[index] = result.data
          }

          // Update in doctors list if it's a doctor
          if (result.data.role === 'doctor') {
            const doctorIndex = this.doctors.findIndex(doctor => doctor.id === userId)
            if (doctorIndex !== -1) {
              this.doctors[doctorIndex] = result.data
            } else {
              this.doctors.push(result.data)
            }
          } else {
            // Remove from doctors list if role changed
            const doctorIndex = this.doctors.findIndex(doctor => doctor.id === userId)
            if (doctorIndex !== -1) {
              this.doctors.splice(doctorIndex, 1)
            }
          }

          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        console.error('Error updating user:', error)
        this.error = 'Failed to update user'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    clearError() {
      this.error = null
    }
  }
})