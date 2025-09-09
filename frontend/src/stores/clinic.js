import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useClinicStore = defineStore('clinic', {
  state: () => ({
    clinics: [],
    currentClinic: null,
    branches: [],
    doctors: [],
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
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch clinic'
        console.error('Error fetching clinic:', error)
        return { success: false, error: 'Failed to fetch clinic' }
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
        console.log('Fetching branches for clinic:', clinicId)
        const result = await apiService.get(`/api/clinics/${clinicId}/branches`)
        console.log('Branches API result:', result)
        if (result.success) {
          this.branches = result.data
          console.log('Loaded branches:', this.branches)
        } else {
          this.error = result.error
          console.error('Failed to fetch branches:', result.error)
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

    async fetchDoctors() {
      this.loading = true
      this.error = null

      try {
        console.log('Fetching doctors from API...')
        const result = await apiService.get('/api/users?role=doctor')
        console.log('Doctors API result:', result)
        if (result.success) {
          this.doctors = result.data.users || result.data
          console.log('Loaded doctors:', this.doctors)
        } else {
          this.error = result.error
          console.error('Failed to fetch doctors:', result.error)
        }
      } catch (error) {
        this.error = 'Failed to fetch doctors'
        console.error('Error fetching doctors:', error)
      } finally {
        this.loading = false
      }
    },

    async deleteClinic(clinicId) {
      this.loading = true
      this.error = null
      try {
        const result = await apiService.delete(`/api/clinics/${clinicId}`)
        if (result.success) {
          this.clinics = this.clinics.filter(c => c.id !== clinicId)
          return { success: true }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to delete clinic'
        console.error('Error deleting clinic:', error)
        return { success: false, error: 'Failed to delete clinic' }
      } finally {
        this.loading = false
      }
    },

    async updateBranch(clinicId, branchId, branchData) {
      this.loading = true
      this.error = null
      try {
        const result = await apiService.put(`/api/clinics/${clinicId}/branches/${branchId}`, branchData)
        if (result.success) {
          const clinicIndex = this.clinics.findIndex(c => c.id === clinicId)
          if (clinicIndex !== -1) {
            const branchIndex = this.clinics[clinicIndex].branches.findIndex(b => b.id === branchId)
            if (branchIndex !== -1) {
              this.clinics[clinicIndex].branches[branchIndex] = result.data
            }
          }
          if (this.currentClinic && this.currentClinic.id === clinicId) {
            const branchIndex = this.currentClinic.branches.findIndex(b => b.id === branchId)
            if (branchIndex !== -1) {
              this.currentClinic.branches[branchIndex] = result.data
            }
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update branch'
        console.error('Error updating branch:', error)
        return { success: false, error: 'Failed to update branch' }
      } finally {
        this.loading = false
      }
    },

    async deleteBranch(clinicId, branchId) {
      this.loading = true
      this.error = null
      try {
        const result = await apiService.delete(`/api/clinics/${clinicId}/branches/${branchId}`)
        if (result.success) {
          const clinicIndex = this.clinics.findIndex(c => c.id === clinicId)
          if (clinicIndex !== -1) {
            this.clinics[clinicIndex].branches = this.clinics[clinicIndex].branches.filter(b => b.id !== branchId)
          }
          if (this.currentClinic && this.currentClinic.id === clinicId) {
            this.currentClinic.branches = this.currentClinic.branches.filter(b => b.id !== branchId)
          }
          return { success: true }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to delete branch'
        console.error('Error deleting branch:', error)
        return { success: false, error: 'Failed to delete branch' }
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