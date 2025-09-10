import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useConsentStore = defineStore('consent', {
  state: () => ({
    consentTemplates: [],
    consentForms: [],
    loading: false,
    error: null
  }),

  getters: {
    isLoading: (state) => state.loading,
    getActiveConsentTemplates: (state) =>
      state.consentTemplates.filter(template => template.is_active),
    getConsentTemplatesByCategory: (state) => (category) =>
      state.consentTemplates.filter(template =>
        template.category === category && template.is_active
      ),
    getConsentTemplateById: (state) => (id) =>
      state.consentTemplates.find(template => template.id === id),
    getDefaultConsentTemplate: (state) =>
      state.consentTemplates.find(template => template.is_default && template.is_active),
    getPatientConsentForms: (state) => (patientId) =>
      state.consentForms.filter(form => form.patient_id === parseInt(patientId))
  },

  actions: {
    async fetchConsentTemplates(params = {}) {
      this.loading = true
      this.error = null

      try {
        const queryParams = new URLSearchParams(params)
        const result = await apiService.get('/api/consent-templates?' + queryParams.toString())

        if (result.success) {
          this.consentTemplates = result.data || []
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          console.warn('Failed to fetch consent templates:', result.error)
          // Don't fail completely, just return empty array
          this.consentTemplates = []
          return { success: true, data: [] }
        }
      } catch (error) {
        console.error('Error fetching consent templates:', error)
        this.error = 'Failed to fetch consent templates'
        // Don't fail completely, just return empty array
        this.consentTemplates = []
        return { success: true, data: [] }
      } finally {
        this.loading = false
      }
    },

    async createConsentTemplate(templateData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/consent-templates', templateData)

        if (result.success) {
          // Add to local state
          this.consentTemplates.push(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        console.error('Error creating consent template:', error)
        this.error = 'Failed to create consent template'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async createConsentForm(formData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/consent-forms', formData)

        if (result.success) {
          // Add to local state
          this.consentForms.push(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        console.error('Error creating consent form:', error)
        this.error = 'Failed to create consent form'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchConsentForms(params = {}) {
      this.loading = true
      this.error = null

      try {
        const queryParams = new URLSearchParams(params)
        const result = await apiService.get('/api/consent-forms?' + queryParams.toString())

        if (result.success) {
          this.consentForms = result.data || []
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          console.warn('Failed to fetch consent forms:', result.error)
          this.consentForms = []
          return { success: true, data: [] }
        }
      } catch (error) {
        console.error('Error fetching consent forms:', error)
        this.error = 'Failed to fetch consent forms'
        this.consentForms = []
        return { success: true, data: [] }
      } finally {
        this.loading = false
      }
    },

    clearError() {
      this.error = null
    }
  }
})