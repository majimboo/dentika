import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useClinicStore = defineStore('clinic', {
  state: () => ({
    clinics: [],
    currentClinic: null,
    branches: [],
    loading: false,
    error: null
  }),

  getters: {
    isLoading: (state) => state.loading,
    getClinicsCount: (state) => state.clinics.length,
    getActiveBranches: (state) => state.branches.filter(branch => branch.is_active),
    getMainBranch: (state) => state.branches.find(branch => branch.is_main_branch)
  },

  actions: {
    async fetchClinics() {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get('/api/clinics')
        if (result.success) {
          this.clinics = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch clinics'
        console.error('Error fetching clinics:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchClinic(clinicId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/clinics/${clinicId}`)
        if (result.success) {
          this.currentClinic = result.data
          this.branches = result.data.branches || []
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch clinic'
        console.error('Error fetching clinic:', error)
      } finally {
        this.loading = false
      }
    },

    async createClinic(clinicData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.post('/api/clinics', clinicData)
        if (result.success) {
          this.clinics.push(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create clinic'
        console.error('Error creating clinic:', error)
        return { success: false, error: 'Failed to create clinic' }
      } finally {
        this.loading = false
      }
    },

    async updateClinic(clinicId, clinicData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/clinics/${clinicId}`, clinicData)
        if (result.success) {
          const index = this.clinics.findIndex(c => c.id === clinicId)
          if (index !== -1) {
            this.clinics[index] = result.data
          }
          if (this.currentClinic && this.currentClinic.id === clinicId) {
            this.currentClinic = result.data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update clinic'
        console.error('Error updating clinic:', error)
        return { success: false, error: 'Failed to update clinic' }
      } finally {
        this.loading = false
      }
    },

    async fetchBranches(clinicId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/clinics/${clinicId}/branches`)
        if (result.success) {
          this.branches = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch branches'
        console.error('Error fetching branches:', error)
      } finally {
        this.loading = false
      }
    },

    async createBranch(clinicId, branchData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.post(`/api/clinics/${clinicId}/branches`, branchData)
        if (result.success) {
          this.branches.push(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create branch'
        console.error('Error creating branch:', error)
        return { success: false, error: 'Failed to create branch' }
      } finally {
        this.loading = false
      }
    },

    clearError() {
      this.error = null
    },

    setCurrentClinic(clinic) {
      this.currentClinic = clinic
      this.branches = clinic.branches || []
    }
  }
})