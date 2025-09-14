<template>
  <div class="super-admin-platform-inventory-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Platform Inventory Management</h1>
        <p class="text-gray-600">Manage Dentika's platform inventory catalog</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/admin/shop/add"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          Add New Item
        </router-link>
        <button
          @click="exportInventory"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-download" class="w-4 h-4 mr-2" />
          Export Inventory
        </button>
      </div>
    </div>

    <!-- Analytics Dashboard -->
    <div class="analytics-grid grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-blue-600">{{ totalItems }}</div>
            <div class="text-sm text-gray-600">Total Items</div>
          </div>
          <div class="p-2 bg-blue-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-box" class="w-5 h-5 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-green-600">{{ activeItems }}</div>
            <div class="text-sm text-gray-600">Active Items</div>
          </div>
          <div class="p-2 bg-green-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-check-circle" class="w-5 h-5 text-green-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-yellow-600">{{ lowStockItems }}</div>
            <div class="text-sm text-gray-600">Low Stock Items</div>
          </div>
          <div class="p-2 bg-yellow-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-5 h-5 text-yellow-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-purple-600">{{ totalValue.toFixed(2) }}</div>
            <div class="text-sm text-gray-600">Total Value (₱)</div>
          </div>
          <div class="p-2 bg-purple-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-dollar-sign" class="w-5 h-5 text-purple-600" />
          </div>
        </div>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <!-- Search and Filters -->
        <div class="search-and-filters flex flex-col md:flex-row gap-4 items-start md:items-center flex-1">
          <div class="search-box relative">
            <input
              v-model="searchQuery"
              @input="debouncedSearch"
              type="text"
              placeholder="Search inventory items..."
              class="w-full md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          </div>

          <div class="filters-row flex flex-wrap items-center gap-4">
            <!-- Category Filter -->
            <select
              v-model="categoryFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Categories</option>
              <option value="equipment">Equipment</option>
              <option value="supplies">Supplies</option>
              <option value="consumables">Consumables</option>
              <option value="instruments">Instruments</option>
              <option value="materials">Materials</option>
            </select>

            <!-- Status Filter -->
            <select
              v-model="statusFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Statuses</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="discontinued">Discontinued</option>
            </select>

            <!-- Clear Filters -->
            <button
              v-if="hasActiveFilters"
              @click="clearFilters"
              class="btn btn-secondary text-sm"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-center">
        <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-5 h-5 text-red-600 mr-3" />
        <div>
          <h3 class="text-sm font-medium text-red-800">Error loading inventory</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Inventory Table -->
    <div v-else-if="inventoryItems.length > 0" class="inventory-table bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Item</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">SKU</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="item in inventoryItems"
              :key="item.id"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <div class="h-10 w-10 rounded-lg bg-gray-200 flex items-center justify-center">
                      <font-awesome-icon icon="fa-solid fa-box" class="w-5 h-5 text-gray-600" />
                    </div>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ item.name }}</div>
                    <div class="text-sm text-gray-500">{{ item.description?.substring(0, 50) }}...</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ item.sku }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-block text-xs px-2 py-1 rounded-full bg-gray-100 text-gray-800 capitalize">
                  {{ item.category }}
                </span>
              </td>
               <td class="px-6 py-4 whitespace-nowrap">
                 <div class="text-sm text-gray-900">
                   {{ item.current_stock }}
                   <span v-if="item.current_stock <= item.min_stock_level" class="text-red-600 font-medium">
                     (Low)
                   </span>
                 </div>
               </td>
               <td class="px-6 py-4 whitespace-nowrap">
                 <div class="text-sm font-medium text-gray-900">₱{{ item.selling_price?.toFixed(2) }}</div>
               </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  :class="getStatusBadgeClass(item.status)"
                  class="inline-block text-xs px-2 py-1 rounded-full font-medium"
                >
                  {{ getStatusText(item.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center space-x-2">
                   <router-link
                     :to="`/admin/shop/${item.id}/view`"
                     class="text-blue-600 hover:text-blue-900"
                   >
                     View
                   </router-link>
                  <router-link
                    :to="`/admin/shop/edit/${item.id}`"
                    class="text-green-600 hover:text-green-900"
                  >
                    Edit
                  </router-link>
                  <button
                    @click="toggleItemStatus(item)"
                    :class="item.status === 'active' ? 'text-yellow-600 hover:text-yellow-900' : 'text-green-600 hover:text-green-900'"
                  >
                    {{ item.status === 'active' ? 'Deactivate' : 'Activate' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <font-awesome-icon icon="fa-solid fa-box" class="w-16 h-16 text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900 mb-2">No inventory items found</h3>
      <p class="text-gray-600 mb-6">There are currently no items in the platform inventory.</p>
      <router-link
        to="/admin/shop/add"
        class="btn btn-primary"
      >
        Add First Item
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="inventoryItems.length > 0 && totalPages > 1" class="pagination mt-6 flex justify-center">
      <nav class="flex items-center space-x-2">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>

        <span
          v-for="page in visiblePages"
          :key="page"
          @click="goToPage(page)"
          :class="page === currentPage ? 'bg-blue-600 text-white' : 'text-gray-500 bg-white hover:bg-gray-50'"
          class="px-3 py-2 text-sm font-medium border border-gray-300 rounded-md cursor-pointer"
        >
          {{ page }}
        </span>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </nav>
    </div>



    <!-- Item Details Modal -->
    <div v-if="showDetailsModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeDetailsModal">
      <div class="relative top-20 mx-auto p-5 border w-11/12 max-w-2xl shadow-lg rounded-md bg-white max-h-screen overflow-y-auto" @click.stop>
        <div class="mt-3">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">Item Details</h3>
            <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600">
              <font-awesome-icon icon="fa-solid fa-times" class="w-5 h-5" />
            </button>
          </div>

          <div v-if="selectedItem" class="space-y-6">
            <!-- Item Header -->
            <div class="flex items-center space-x-4">
              <div class="h-16 w-16 rounded-lg bg-gray-200 flex items-center justify-center">
                <font-awesome-icon icon="fa-solid fa-box" class="w-8 h-8 text-gray-600" />
              </div>
              <div>
                <h4 class="text-xl font-bold text-gray-900">{{ selectedItem.name }}</h4>
                <p class="text-gray-600">{{ selectedItem.sku }}</p>
              </div>
            </div>

            <!-- Item Details -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700">Category</label>
                <div class="mt-1 text-sm text-gray-900 capitalize">{{ selectedItem.category }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Status</label>
                <span
                  :class="getStatusBadgeClass(selectedItem.status)"
                  class="inline-block text-sm px-3 py-1 rounded-full font-medium mt-1"
                >
                  {{ getStatusText(selectedItem.status) }}
                </span>
              </div>
               <div>
                 <label class="block text-sm font-medium text-gray-700">Price</label>
                 <div class="mt-1 text-sm font-medium text-gray-900">₱{{ selectedItem.selling_price?.toFixed(2) }}</div>
               </div>
               <div>
                 <label class="block text-sm font-medium text-gray-700">Stock Quantity</label>
                 <div class="mt-1 text-sm text-gray-900">
                   {{ selectedItem.current_stock }}
                   <span v-if="selectedItem.current_stock <= selectedItem.min_stock_level" class="text-red-600 font-medium">
                     (Low Stock)
                   </span>
                 </div>
               </div>
               <div>
                 <label class="block text-sm font-medium text-gray-700">Unit</label>
                 <div class="mt-1 text-sm text-gray-900">{{ selectedItem.unit_of_measure || 'N/A' }}</div>
               </div>
               <div>
                 <label class="block text-sm font-medium text-gray-700">Supplier</label>
                 <div class="mt-1 text-sm text-gray-900">{{ selectedItem.supplier_name || 'N/A' }}</div>
               </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <div class="text-sm text-gray-900 bg-gray-50 rounded-lg p-4">{{ selectedItem.description || 'No description available' }}</div>
            </div>

            <!-- Actions -->
            <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
              <button @click="closeDetailsModal" class="btn btn-secondary">Close</button>
              <button @click="editItem(selectedItem)" class="btn btn-primary">Edit Item</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'
import debounce from 'lodash/debounce'

export default {
  name: 'SuperAdminPlatformInventory',
  setup() {
    const authStore = useAuthStore()

    // Reactive data
    const inventoryItems = ref([])
    const loading = ref(false)
    const error = ref('')
    const currentPage = ref(1)
    const limit = ref(20)
    const totalPages = ref(0)

    // Filters
    const searchQuery = ref('')
    const categoryFilter = ref('')
    const statusFilter = ref('')

    // Selection
    const selectedItems = ref([])

    // Modals
    const showDetailsModal = ref(false)
    const selectedItem = ref(null)

    // Computed properties
    const totalItems = computed(() => inventoryItems.value.length)
    const activeItems = computed(() => {
      return inventoryItems.value.filter(item => item.status === 'active').length
    })
    const lowStockItems = computed(() => {
      return inventoryItems.value.filter(item => item.current_stock <= item.min_stock_level).length
    })
    const totalValue = computed(() => {
      return inventoryItems.value.reduce((sum, item) => sum + (item.selling_price * item.current_stock), 0)
    })

    const hasActiveFilters = computed(() => {
      return searchQuery.value || categoryFilter.value || statusFilter.value
    })

    const visiblePages = computed(() => {
      const pages = []
      const start = Math.max(1, currentPage.value - 2)
      const end = Math.min(totalPages.value, currentPage.value + 2)

      for (let i = start; i <= end; i++) {
        pages.push(i)
      }

      return pages
    })

    // Methods
    const debouncedSearch = debounce(() => {
      applyFilters()
    }, 300)

    const applyFilters = () => {
      currentPage.value = 1
      loadInventory()
    }

    const clearFilters = () => {
      searchQuery.value = ''
      categoryFilter.value = ''
      statusFilter.value = ''
      applyFilters()
    }

    const loadInventory = async () => {
      loading.value = true
      error.value = ''

      try {
        const params = {
          page: currentPage.value,
          limit: limit.value,
          search: searchQuery.value,
          category: categoryFilter.value,
          status: statusFilter.value
        }

        const response = await apiService.get('/api/inventory/shop/items', { params })
        if (response.success) {
          inventoryItems.value = response.data.items || []
          totalPages.value = Math.ceil((response.data.total || 0) / limit.value)
        } else {
          throw new Error(response.error || 'Failed to load inventory')
        }
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load inventory'
        console.error('Failed to load inventory:', err)
      } finally {
        loading.value = false
      }
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadInventory()
      }
    }

    const getStatusBadgeClass = (status) => {
      const classes = {
        active: 'bg-green-100 text-green-800',
        inactive: 'bg-gray-100 text-gray-800',
        discontinued: 'bg-red-100 text-red-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    }

    const getStatusText = (status) => {
      const texts = {
        active: 'Active',
        inactive: 'Inactive',
        discontinued: 'Discontinued'
      }
      return texts[status] || status
    }

    const viewItemDetails = (item) => {
      selectedItem.value = item
      showDetailsModal.value = true
    }



    const toggleItemStatus = async (item) => {
      const newStatus = item.status === 'active' ? 'inactive' : 'active'
      const action = newStatus === 'active' ? 'activate' : 'deactivate'

      if (!confirm(`Are you sure you want to ${action} this item?`)) {
        return
      }

      try {
        // Update the item status using the general update endpoint
        const result = await apiService.put(`/api/inventory/shop/items/${item.id}`, {
          ...item,
          status: newStatus
        })
        if (!result.success) {
          throw new Error(result.error || 'Failed to update item status')
        }

        alert(`Item ${action}d successfully!`)
        loadInventory()
      } catch (err) {
        alert('Failed to update item status. Please try again.')
        console.error('Failed to update item status:', err)
      }
    }



    const closeDetailsModal = () => {
      showDetailsModal.value = false
      selectedItem.value = null
    }

    const exportInventory = () => {
      // TODO: Implement export functionality
      alert('Export functionality will be implemented soon')
    }

    // Lifecycle
    onMounted(() => {
      loadInventory()
    })

    // Watchers
    watch([currentPage], () => {
      loadInventory()
    })

    return {
      // Data
      inventoryItems,
      loading,
      error,
      currentPage,
      totalPages,
      searchQuery,
      categoryFilter,
      statusFilter,
      selectedItems,
      showDetailsModal,
      selectedItem,

      // Computed
      totalItems,
      activeItems,
      lowStockItems,
      totalValue,
      hasActiveFilters,
      visiblePages,

      // Methods
      debouncedSearch,
      applyFilters,
      clearFilters,
      loadInventory,
      goToPage,
      getStatusBadgeClass,
      getStatusText,
      viewItemDetails,
      toggleItemStatus,
      closeDetailsModal,
      exportInventory
    }
  }
}
</script>

<style scoped>
@import "../styles/main.css";

.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors duration-200;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-700 hover:bg-gray-300;
}

.analytics-card {
  transition: all 0.2s ease-in-out;
}

.analytics-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}
</style>