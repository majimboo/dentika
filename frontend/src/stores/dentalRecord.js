import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useDentalRecordStore = defineStore('dentalRecord', {
  state: () => ({
    dentalRecords: [],
    currentRecord: null,
    toothHistory: [],
    recordHistory: [],
    dentalSnapshots: [],
    currentSnapshot: null,
    loading: false,
    error: null,
    selectedTooth: null
  }),

  getters: {
    isLoading: (state) => state.loading,
    getActiveDentalRecord: (state) => {
      return state.dentalRecords.find(record => record.is_active)
    },
    getPermanentTeethRecord: (state) => {
      return state.dentalRecords.find(record => 
        record.record_type === 'permanent' && record.is_active
      )
    },
    getPrimaryTeethRecord: (state) => {
      return state.dentalRecords.find(record => 
        record.record_type === 'primary' && record.is_active
      )
    },
    getToothByNumber: (state) => (toothNumber) => {
      if (!state.currentRecord?.teeth_data) return null
      return state.currentRecord.teeth_data.find(tooth => 
        tooth.tooth_number === toothNumber
      )
    },
    getTeethByCondition: (state) => (condition) => {
      if (!state.currentRecord?.teeth_data) return []
      return state.currentRecord.teeth_data.filter(tooth => 
        tooth.condition === condition
      )
    },
    getTeethNeedingTreatment: (state) => {
      if (!state.currentRecord?.teeth_data) return []
      return state.currentRecord.teeth_data.filter(tooth => 
        ['decay', 'damaged'].includes(tooth.condition)
      )
    }
  },

  actions: {
    async fetchPatientDentalRecords(patientId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/patients/${patientId}/dental-records`)
        if (result.success) {
          this.dentalRecords = result.data
          // Set current record to active record
          const activeRecord = this.dentalRecords.find(record => record.is_active)
          if (activeRecord) {
            this.currentRecord = activeRecord
          }
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

    async fetchDentalRecord(recordId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/dental-records/${recordId}`)
        if (result.success) {
          this.currentRecord = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch dental record'
        console.error('Error fetching dental record:', error)
      } finally {
        this.loading = false
      }
    },

    async updateToothCondition(recordId, toothData) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/dental-records/${recordId}/tooth`, toothData)
        if (result.success) {
          // Update current record with new teeth data
          if (this.currentRecord && this.currentRecord.id === recordId) {
            this.currentRecord.teeth_data = result.data.teeth_data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update tooth condition'
        console.error('Error updating tooth condition:', error)
        return { success: false, error: 'Failed to update tooth condition' }
      } finally {
        this.loading = false
      }
    },

    async bulkUpdateTeeth(recordId, updates, reason = '') {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/dental-records/${recordId}/bulk-update`, {
          updates,
          reason
        })
        if (result.success) {
          // Update current record with new teeth data
          if (this.currentRecord && this.currentRecord.id === recordId) {
            this.currentRecord.teeth_data = result.data.teeth_data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update teeth'
        console.error('Error updating teeth:', error)
        return { success: false, error: 'Failed to update teeth' }
      } finally {
        this.loading = false
      }
    },

    async activateDentalRecord(recordId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.put(`/api/dental-records/${recordId}/activate`)
        if (result.success) {
          // Update the records array
          this.dentalRecords.forEach(record => {
            record.is_active = record.id === recordId
          })
          // Set as current record
          const activatedRecord = this.dentalRecords.find(record => record.id === recordId)
          if (activatedRecord) {
            this.currentRecord = activatedRecord
          }
          return { success: true }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to activate dental record'
        console.error('Error activating dental record:', error)
        return { success: false, error: 'Failed to activate dental record' }
      } finally {
        this.loading = false
      }
    },

    async fetchRecordHistory(recordId) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/dental-records/${recordId}/history`)
        if (result.success) {
          this.recordHistory = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch record history'
        console.error('Error fetching record history:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchToothHistory(recordId, toothNumber) {
      this.loading = true
      this.error = null
      
      try {
        const result = await apiService.get(`/api/dental-records/${recordId}/tooth-history?tooth_number=${toothNumber}`)
        if (result.success) {
          this.toothHistory = result.data
        } else {
          this.error = result.error
        }
      } catch (error) {
        this.error = 'Failed to fetch tooth history'
        console.error('Error fetching tooth history:', error)
      } finally {
        this.loading = false
      }
    },

    setSelectedTooth(toothNumber) {
      this.selectedTooth = toothNumber
    },

    clearSelectedTooth() {
      this.selectedTooth = null
    },

    setCurrentRecord(record) {
      this.currentRecord = record
    },

    clearCurrentRecord() {
      this.currentRecord = null
      this.selectedTooth = null
      this.toothHistory = []
      this.recordHistory = []
    },

    clearError() {
      this.error = null
    },

    // Helper methods for tooth chart visualization
    getToothColor(tooth) {
      if (!tooth) return '#ffffff'
      
      const colorMap = {
        healthy: '#90EE90',    // Light green
        decay: '#FFB6C1',      // Light pink
        filled: '#87CEEB',     // Sky blue  
        crowned: '#DDA0DD',    // Plum
        root_canal: '#F0E68C', // Khaki
        extracted: '#FF6347',  // Tomato
        implant: '#20B2AA',    // Light sea green
        bridge: '#DEB887',     // Burlywood
        missing: '#D3D3D3'     // Light gray
      }
      
      return colorMap[tooth.condition] || '#ffffff'
    },

    getToothDisplayText(tooth) {
      if (!tooth) return ''

      const displayMap = {
        healthy: 'H',
        decay: 'D',
        filled: 'F',
        crowned: 'C',
        root_canal: 'RC',
        extracted: 'X',
        implant: 'I',
        bridge: 'B',
        missing: 'M'
      }

      return displayMap[tooth.condition] || ''
    },

    // Dental Chart Snapshot Methods
    async createDentalChartSnapshot(patientId, snapshotData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post(`/api/patients/${patientId}/dental-snapshots`, snapshotData)
        if (result.success) {
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create dental chart snapshot'
        console.error('Error creating dental chart snapshot:', error)
        return { success: false, error: 'Failed to create dental chart snapshot' }
      } finally {
        this.loading = false
      }
    },

    async fetchPatientDentalSnapshots(patientId, chartType = null) {
      this.loading = true
      this.error = null

      try {
        let url = `/api/patients/${patientId}/dental-snapshots`
        if (chartType) {
          url += `?chart_type=${chartType}`
        }

        const result = await apiService.get(url)
        if (result.success) {
          this.dentalSnapshots = result.data
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch dental snapshots'
        console.error('Error fetching dental snapshots:', error)
        return { success: false, error: 'Failed to fetch dental snapshots' }
      } finally {
        this.loading = false
      }
    },

    async fetchDentalSnapshot(snapshotId) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.get(`/api/dental-snapshots/${snapshotId}`)
        if (result.success) {
          this.currentSnapshot = result.data
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch dental snapshot'
        console.error('Error fetching dental snapshot:', error)
        return { success: false, error: 'Failed to fetch dental snapshot' }
      } finally {
        this.loading = false
      }
    },

    setCurrentSnapshot(snapshot) {
      this.currentSnapshot = snapshot
    },

    clearCurrentSnapshot() {
      this.currentSnapshot = null
    },

    // Enhanced tooth update with appointment linking
    async updateToothConditionWithAppointment(recordId, toothData, appointmentId = null) {
      this.loading = true
      this.error = null

      try {
        const payload = {
          ...toothData,
          appointment_id: appointmentId
        }

        const result = await apiService.put(`/api/dental-records/${recordId}/tooth`, payload)
        if (result.success) {
          // Update current record with new teeth data
          if (this.currentRecord && this.currentRecord.id === recordId) {
            this.currentRecord.teeth_data = result.data.teeth_data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update tooth condition'
        console.error('Error updating tooth condition:', error)
        return { success: false, error: 'Failed to update tooth condition' }
      } finally {
        this.loading = false
      }
    }
  }
})