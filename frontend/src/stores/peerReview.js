import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const usePeerReviewStore = defineStore('peerReview', () => {
  // State
  const cases = ref([])
  const currentCase = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const getCaseById = computed(() => {
    return (id) => cases.value.find(caseItem => caseItem.id === id)
  })

  const getCasesByStatus = computed(() => {
    return (status) => cases.value.filter(caseItem => caseItem.status === status)
  })

  const getCasesByVisibility = computed(() => {
    return (visibility) => cases.value.filter(caseItem => caseItem.visibility === visibility)
  })

  // Actions
  const fetchCases = async (filters = {}) => {
    loading.value = true
    error.value = null

    try {
      const params = new URLSearchParams(filters)
      const response = await fetch(`/api/peer-review/cases?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        cases.value = data.cases
        return data
      } else {
        const errorData = await response.json()
        error.value = errorData.error || 'Failed to fetch cases'
        throw new Error(error.value)
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchCase = async (id) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/peer-review/cases/${id}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const caseData = await response.json()
        currentCase.value = caseData

        // Update case in cases array if it exists
        const index = cases.value.findIndex(c => c.id === id)
        if (index !== -1) {
          cases.value[index] = caseData
        }

        return caseData
      } else {
        const errorData = await response.json()
        error.value = errorData.error || 'Failed to fetch case'
        throw new Error(error.value)
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const createCase = async (caseData) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch('/api/peer-review/cases', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(caseData)
      })

      if (response.ok) {
        const newCase = await response.json()
        cases.value.unshift(newCase.case) // Add to beginning of list
        return newCase
      } else {
        const errorData = await response.json()
        error.value = errorData.error || 'Failed to create case'
        throw new Error(error.value)
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addComment = async (caseId, commentData) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/peer-review/cases/${caseId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(commentData)
      })

      if (response.ok) {
        const comment = await response.json()

        // Update current case if it's the one being commented on
        if (currentCase.value && currentCase.value.id === caseId) {
          currentCase.value.comments = currentCase.value.comments || []
          currentCase.value.comments.push(comment)
        }

        // Update case in cases array
        const caseIndex = cases.value.findIndex(c => c.id === caseId)
        if (caseIndex !== -1) {
          cases.value[caseIndex].comments = cases.value[caseIndex].comments || []
          cases.value[caseIndex].comments.push(comment)
        }

        return comment
      } else {
        const errorData = await response.json()
        error.value = errorData.error || 'Failed to add comment'
        throw new Error(error.value)
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateCaseStatus = async (caseId, status) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/peer-review/cases/${caseId}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ status })
      })

      if (response.ok) {
        const result = await response.json()

        // Update current case
        if (currentCase.value && currentCase.value.id === caseId) {
          currentCase.value.status = status
        }

        // Update case in cases array
        const caseIndex = cases.value.findIndex(c => c.id === caseId)
        if (caseIndex !== -1) {
          cases.value[caseIndex].status = status
        }

        return result
      } else {
        const errorData = await response.json()
        error.value = errorData.error || 'Failed to update status'
        throw new Error(error.value)
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    cases.value = []
    currentCase.value = null
    loading.value = false
    error.value = null
  }

  return {
    // State
    cases,
    currentCase,
    loading,
    error,

    // Getters
    getCaseById,
    getCasesByStatus,
    getCasesByVisibility,

    // Actions
    fetchCases,
    fetchCase,
    createCase,
    addComment,
    updateCaseStatus,
    clearError,
    reset
  }
})