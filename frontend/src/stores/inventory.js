import { defineStore } from 'pinia'
import apiService from '../services/api'

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    items: [],
    currentItem: null,
    stockTransactions: [],
    alerts: [],
    analytics: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      limit: 20
    },
    stockPagination: {
      total: 0,
      page: 1,
      limit: 50
    },
    stats: {
      total_value: 0,
      low_stock_count: 0,
      expiring_count: 0
    }
  }),

  getters: {
    isLoading: (state) => state.loading,
    getTotalItems: (state) => state.pagination.total,
    getLowStockItems: (state) => state.stats.low_stock_count,
    getExpiringItems: (state) => state.stats.expiring_count,
    getTotalValue: (state) => state.stats.total_value,
    getItemsByCategory: (state) => {
      const categories = {}
      state.items.forEach(item => {
        if (!categories[item.category]) {
          categories[item.category] = []
        }
        categories[item.category].push(item)
      })
      return categories
    }
  },

  actions: {
    async fetchItems(params = {}) {
      this.loading = true
      this.error = null

      const queryParams = new URLSearchParams({
        page: params.page || this.pagination.page,
        limit: params.limit || this.pagination.limit,
        ...params
      })

      try {
        const result = await apiService.get(`/api/inventory/items?${queryParams}`)
        if (result.success) {
          this.items = result.data.items
          this.pagination = {
            total: result.data.total,
            page: result.data.page,
            limit: result.data.limit
          }
          // Store the aggregated stats from the server
          if (result.data.stats) {
            this.stats = result.data.stats
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch inventory items'
        console.error('Error fetching inventory items:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchItem(itemId) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.get(`/api/inventory/items/${itemId}`)
        if (result.success) {
          this.currentItem = result.data
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch inventory item'
        console.error('Error fetching inventory item:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async createItem(itemData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/inventory/items', itemData)
        if (result.success) {
          // Add to items list if we're on the current page
          if (this.items.length < this.pagination.limit) {
            this.items.unshift(result.data)
            this.pagination.total++
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create inventory item'
        console.error('Error creating inventory item:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async updateItem(itemId, itemData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.put(`/api/inventory/items/${itemId}`, itemData)
        if (result.success) {
          // Update item in the list
          const index = this.items.findIndex(item => item.id === parseInt(itemId))
          if (index !== -1) {
            this.items[index] = result.data
          }
          // Update current item if it's the same
          if (this.currentItem && this.currentItem.id === parseInt(itemId)) {
            this.currentItem = result.data
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to update inventory item'
        console.error('Error updating inventory item:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async deleteItem(itemId) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.delete(`/api/inventory/items/${itemId}`)
        if (result.success) {
          // Remove from items list
          this.items = this.items.filter(item => item.id !== parseInt(itemId))
          this.pagination.total--
          return { success: true }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to delete inventory item'
        console.error('Error deleting inventory item:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async createStockTransaction(transactionData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/inventory/stock-transactions', transactionData)
        if (result.success) {
          // Update the item's stock in the list
          const itemIndex = this.items.findIndex(item => item.id === transactionData.item_id)
          if (itemIndex !== -1) {
            // Refresh the item data to get updated stock
            const itemResult = await this.fetchItem(transactionData.item_id)
            if (itemResult.success) {
              this.items[itemIndex] = itemResult.data
            }
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create stock transaction'
        console.error('Error creating stock transaction:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchStockTransactions(itemId, params = {}) {
      this.loading = true
      this.error = null

      const queryParams = new URLSearchParams({
        page: params.page || this.stockPagination.page,
        limit: params.limit || this.stockPagination.limit,
        ...params
      })

      try {
        const result = await apiService.get(`/api/inventory/items/${itemId}/stock-transactions?${queryParams}`)
        if (result.success) {
          this.stockTransactions = result.data.transactions
          this.stockPagination = {
            total: result.data.total,
            page: result.data.page,
            limit: result.data.limit
          }
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch stock transactions'
        console.error('Error fetching stock transactions:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchAlerts(params = {}) {
      this.loading = true
      this.error = null

      const queryParams = new URLSearchParams({
        page: params.page || 1,
        limit: params.limit || 50,
        ...params
      })

      try {
        const result = await apiService.get(`/api/inventory/alerts?${queryParams}`)
        if (result.success) {
          this.alerts = result.data.alerts
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch inventory alerts'
        console.error('Error fetching inventory alerts:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchAnalytics(params = {}) {
      this.loading = true
      this.error = null

      const queryParams = new URLSearchParams(params)

      try {
        const result = await apiService.get(`/api/inventory/analytics?${queryParams}`)
        if (result.success) {
          this.analytics = result.data.analytics
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch inventory analytics'
        console.error('Error fetching inventory analytics:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async createRestockOrder(orderData) {
      this.loading = true
      this.error = null

      try {
        const result = await apiService.post('/api/inventory/restock-orders', orderData)
        if (result.success) {
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to create restock order'
        console.error('Error creating restock order:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async fetchRestockOrders(params = {}) {
      this.loading = true
      this.error = null

      const queryParams = new URLSearchParams({
        page: params.page || 1,
        limit: params.limit || 50,
        ...params
      })

      try {
        const result = await apiService.get(`/api/inventory/restock-orders?${queryParams}`)
        if (result.success) {
          return { success: true, data: result.data }
        } else {
          this.error = result.error
          return { success: false, error: result.error }
        }
      } catch (error) {
        this.error = 'Failed to fetch restock orders'
        console.error('Error fetching restock orders:', error)
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    // Utility methods
    clearError() {
      this.error = null
    },

    resetState() {
      this.items = []
      this.currentItem = null
      this.stockTransactions = []
      this.alerts = []
      this.analytics = null
      this.loading = false
      this.error = null
      this.pagination = {
        total: 0,
        page: 1,
        limit: 20
      }
      this.stockPagination = {
        total: 0,
        page: 1,
        limit: 50
      }
      this.stats = {
        total_value: 0,
        low_stock_count: 0,
        expiring_count: 0
      }
    }
  }
})