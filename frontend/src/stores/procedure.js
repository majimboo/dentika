import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useProcedureStore = defineStore('procedure', {
  state: () => ({
    procedureTemplates: [],
    loading: false,
    error: null
  }),

  getters: {
    isLoading: (state) => state.loading,
    getActiveProcedureTemplates: (state) =>
      state.procedureTemplates.filter(template => template.is_active),
    getProcedureTemplatesByCategory: (state) => (category) =>
      state.procedureTemplates.filter(template =>
        template.category === category && template.is_active
      ),
    getProcedureTemplateById: (state) => (id) =>
      state.procedureTemplates.find(template => template.id === id)
  },

  actions: {
    async fetchProcedureTemplates(params = {}) {
      this.loading = true
      this.error = null

      try {
        const queryParams = new URLSearchParams(params)
        const result = await apiService.get('/api/procedure-templates?' + queryParams.toString())

        if (result.success) {
          this.procedureTemplates = result.data || []
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          console.warn('Failed to fetch procedure templates:', result.error)
          // Don't fail completely, just return empty array
          this.procedureTemplates = []
          return { success: true, data: [] }
        }
      } catch (error) {
        console.error('Error fetching procedure templates:', error)
        this.error = 'Failed to fetch procedure templates'
        // Don't fail completely, just return empty array
        this.procedureTemplates = []
        return { success: true, data: [] }
      } finally {
        this.loading = false
      }
    },

    async createProcedureTemplate(templateData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/procedure-templates', templateData)

        if (result.success) {
          // Add to local state
          this.procedureTemplates.push(result.data)
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        console.error('Error creating procedure template:', error)
        this.error = 'Failed to create procedure template'
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