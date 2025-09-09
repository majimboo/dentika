import { defineStore } from 'pinia'
import apiService from '../services/api'

export const usePatientStore = defineStore('patient', {
  state: () => ({
    patients: [],
    currentPatient: null,
    patientDentalRecords: [],
    loading: false,
    searchLoading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      limit: 50
    },
    searchQuery: ''
  }),

  getters: {
    isLoading: (state) => state.loading,
    isSearchLoading: (state) => state.searchLoading,
    getTotalPatients: (state) => state.pagination.total,
    getActivePatients: (state) => state.patients.filter(patient => patient.is_active),
    getCurrentPatientAge: (state) => {
      if (!state.currentPatient?.date_of_birth) return null
      const birthDate = new Date(state.currentPatient.date_of_birth)
      const today = new Date()
      let age = today.getFullYear() - birthDate.getFullYear()
      const monthDiff = today.getMonth() - birthDate.getMonth()
      if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
        age--
      }
      return age
    }
  },

  actions: {
    async fetchPatients(params = {}) {
      this.loading = true
      this.error = null
      
      const queryParams = new URLSearchParams({
        page: params.page || this.pagination.page,
        limit: params.limit || this.pagination.limit,
        ...params
      })
      
      if (this.searchQuery) {
        queryParams.set('search', this.searchQuery)
      }
      
      try {
        const result = await apiService.get(`/api/patients?${queryParams}`)
        if (result.success) {
          this.patients = result.data.patients
          this.pagination = {
            total: result.data.total,
            page: result.data.page,
            limit: result.data.limit
          }
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch patients'
        console.error('Error fetching patients:', error)
      } finally {
        this.loading = false
      }
    },

    async searchPatients(query) {
      this.searchQuery = query
      this.searchLoading = true
      
      try {
        await this.fetchPatients({ page: 1 })
      } finally {
        this.searchLoading = false
      }
    },

    async fetchPatient(patientId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/patients/${patientId}`)
        if (result.success) {
          this.currentPatient = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch patient'
        console.error('Error fetching patient:', error)
      } finally {
        this.loading = false
      }
    },

    async createPatient(patientData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.post('/api/patients', patientData)
        if (result.success) {
          this.patients.unshift(result.data)
          this.pagination.total++
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create patient'
        console.error('Error creating patient:', error)
        return { success: false, error: 'Failed to create patient' }
      } finally {
        this.loading = false
      }
    },

    async updatePatient(patientId, patientData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/patients/${patientId}`, patientData)
        if (result.success) {
          const index = this.patients.findIndex(p => p.id === patientId)
          if (index !== -1) {
            this.patients[index] = result.data
          }
          if (this.currentPatient && this.currentPatient.id === patientId) {
            this.currentPatient = result.data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update patient'
        console.error('Error updating patient:', error)
        return { success: false, error: 'Failed to update patient' }
      } finally {
        this.loading = false
      }
    },

    async deactivatePatient(patientId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.delete(`/api/patients/${patientId}`)
        if (result.success) {
          const index = this.patients.findIndex(p => p.id === patientId)
          if (index !== -1) {
            this.patients[index].is_active = false
          }
          return { success: true }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to deactivate patient'
        console.error('Error deactivating patient:', error)
        return { success: false, error: 'Failed to deactivate patient' }
      } finally {
        this.loading = false
      }
    },

    async fetchPatientDentalRecords(patientId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/patients/${patientId}/dental-records`)
        if (result.success) {
          this.patientDentalRecords = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch dental records'
        console.error('Error fetching dental records:', error)
      } finally {
        this.loading = false
      }
    },

    clearError() {
      this.error = null
    },

    clearSearch() {
      this.searchQuery = ''
      this.fetchPatients({ page: 1 })
    },

    setCurrentPatient(patient) {
      this.currentPatient = patient
    },

    clearCurrentPatient() {
      this.currentPatient = null
      this.patientDentalRecords = []
    }
  }
})