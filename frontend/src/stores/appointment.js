import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useAppointmentStore = defineStore('appointment', {
  state: () => ({
    appointments: [],
    upcomingAppointments: [],
    currentAppointment: null,
    loading: false,
    error: null,
    filters: {
      date: null,
      startDate: null,
      endDate: null,
      status: '',
      patientId: null,
      doctorId: null,
      branchId: null
    }
  }),

  getters: {
    isLoading: (state) => state.loading,
    getTodaysAppointments: (state) => {
      const today = new Date().toISOString().split('T')[0]
      return state.appointments.filter(apt => 
        apt.start_time.startsWith(today)
      )
    },
    getAppointmentsWithCountdown: (state) => {
      return state.upcomingAppointments.filter(apt => 
        apt.should_show_countdown
      )
    },
    getNextAppointment: (state) => {
      const now = new Date()
      return state.appointments
        .filter(apt => new Date(apt.start_time) > now && apt.status === 'scheduled')
        .sort((a, b) => new Date(a.start_time) - new Date(b.start_time))[0]
    },
    getAppointmentsByStatus: (state) => (status) => {
      return state.appointments.filter(apt => apt.status === status)
    },
    getAppointmentsByDate: (state) => (date) => {
      return state.appointments.filter(apt => 
        apt.start_time.startsWith(date)
      )
    }
  },

  actions: {
    async fetchAppointments(params = {}) {
      this.loading = true
      this.error = null
      
      const queryParams = new URLSearchParams()
      
      // Apply current filters
      if (this.filters.date) queryParams.set('date', this.filters.date)
      if (this.filters.startDate) queryParams.set('start_date', this.filters.startDate)
      if (this.filters.endDate) queryParams.set('end_date', this.filters.endDate)
      if (this.filters.status) queryParams.set('status', this.filters.status)
      if (this.filters.patientId) queryParams.set('patient_id', this.filters.patientId)
      if (this.filters.doctorId) queryParams.set('doctor_id', this.filters.doctorId)
      if (this.filters.branchId) queryParams.set('branch_id', this.filters.branchId)
      
      // Apply additional params
      Object.keys(params).forEach(key => {
        if (params[key]) queryParams.set(key, params[key])
      })
      
      try {
        const result = await apiService.get(`/api/appointments?${queryParams}`)
        if (result.success) {
          this.appointments = result.data.appointments || result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch appointments'
        console.error('Error fetching appointments:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchTodaysAppointments() {
      return this.fetchAppointments({ today: 'true' })
    },

    async fetchUpcomingAppointments() {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get('/api/appointments/upcoming')
        if (result.success) {
          this.upcomingAppointments = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch upcoming appointments'
        console.error('Error fetching upcoming appointments:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchAppointment(appointmentId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/appointments/${appointmentId}`)
        if (result.success) {
          this.currentAppointment = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch appointment'
        console.error('Error fetching appointment:', error)
      } finally {
        this.loading = false
      }
    },

    async createAppointment(appointmentData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.post('/api/appointments', appointmentData)
        if (result.success) {
          this.appointments.unshift(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create appointment'
        console.error('Error creating appointment:', error)
        return { success: false, error: 'Failed to create appointment' }
      } finally {
        this.loading = false
      }
    },

    async updateAppointment(appointmentId, appointmentData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/appointments/${appointmentId}`, appointmentData)
        if (result.success) {
          const index = this.appointments.findIndex(apt => apt.id === appointmentId)
          if (index !== -1) {
            this.appointments[index] = result.data
          }
          if (this.currentAppointment && this.currentAppointment.id === appointmentId) {
            this.currentAppointment = result.data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update appointment'
        console.error('Error updating appointment:', error)
        return { success: false, error: 'Failed to update appointment' }
      } finally {
        this.loading = false
      }
    },

    async updateAppointmentStatus(appointmentId, statusData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/appointments/${appointmentId}/status`, statusData)
        if (result.success) {
          const index = this.appointments.findIndex(apt => apt.id === appointmentId)
          if (index !== -1) {
            this.appointments[index] = { ...this.appointments[index], ...result.data }
          }
          if (this.currentAppointment && this.currentAppointment.id === appointmentId) {
            this.currentAppointment = { ...this.currentAppointment, ...result.data }
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update appointment status'
        console.error('Error updating appointment status:', error)
        return { success: false, error: 'Failed to update appointment status' }
      } finally {
        this.loading = false
      }
    },

    async markPatientArrived(appointmentId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.post(`/api/appointments/${appointmentId}/arrived`)
        if (result.success) {
          const index = this.appointments.findIndex(apt => apt.id === appointmentId)
          if (index !== -1) {
            this.appointments[index] = {
              ...this.appointments[index],
              patient_arrived: true,
              arrival_time: new Date().toISOString(),
              is_late: result.data.is_late
            }
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to mark patient as arrived'
        console.error('Error marking patient as arrived:', error)
        return { success: false, error: 'Failed to mark patient as arrived' }
      } finally {
        this.loading = false
      }
    },

    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
    },

    clearFilters() {
      this.filters = {
        date: null,
        startDate: null,
        endDate: null,
        status: '',
        patientId: null,
        doctorId: null,
        branchId: null
      }
    },

    setCurrentAppointment(appointment) {
      this.currentAppointment = appointment
    },

    clearCurrentAppointment() {
      this.currentAppointment = null
    },

    clearError() {
      this.error = null
    }
  }
})