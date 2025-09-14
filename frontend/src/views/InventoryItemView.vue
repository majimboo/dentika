<template>
  <div class="inventory-item-view-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div class="flex items-center space-x-4">
        <router-link
          :to="backRoute"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back
        </router-link>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">{{ item?.name }}</h1>
          <p class="text-gray-600">{{ item?.sku }}</p>
        </div>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          :to="`/inventory/${item?.id}/restock`"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          Restock Item
        </router-link>
        <router-link
          :to="`/inventory/${item?.id}/edit`"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-edit" class="w-4 h-4 mr-2" />
          Edit Item
        </router-link>
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
          <h3 class="text-sm font-medium text-red-800">Error loading item</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Item Details -->
    <div v-else-if="item" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Item Overview Card -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="flex items-start space-x-6">
            <!-- Item Image -->
            <div class="flex-shrink-0">
              <div class="h-24 w-24 rounded-lg bg-gray-200 flex items-center justify-center">
                <img
                  v-if="item.image_path"
                  :src="getImageUrl(item.image_path)"
                  :alt="item.name"
                  class="h-24 w-24 rounded-lg object-cover"
                >
                <font-awesome-icon
                  v-else
                  icon="fa-solid fa-box"
                  class="w-12 h-12 text-gray-600"
                />
              </div>
            </div>

            <!-- Item Details -->
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <span class="inline-block bg-blue-100 text-blue-800 text-sm px-3 py-1 rounded-full capitalize">
                  {{ item.category }}
                </span>
                <span
                  :class="getStatusBadgeClass(item.status)"
                  class="inline-block text-sm px-3 py-1 rounded-full font-medium"
                >
                  {{ getStatusText(item.status) }}
                </span>
              </div>

              <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700">Current Stock</label>
                  <div class="text-lg font-semibold text-gray-900">
                    {{ item.current_stock }} {{ item.unit_of_measure }}
                  </div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700">Min Level</label>
                  <div class="text-lg font-semibold text-gray-900">{{ item.min_stock_level }}</div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700">Unit Price</label>
                  <div class="text-lg font-semibold text-gray-900">₱{{ item.selling_price?.toFixed(2) }}</div>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700">Stock Status</label>
                  <div
                    :class="item.current_stock <= item.min_stock_level ? 'text-red-600' : 'text-green-600'"
                    class="text-lg font-semibold"
                  >
                    {{ item.current_stock <= item.min_stock_level ? 'Low Stock' : 'In Stock' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Description</h3>
          <div class="text-gray-700">
            {{ item.description || 'No description available' }}
          </div>
        </div>

        <!-- Stock History -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Recent Stock Transactions</h3>

          <div v-if="stockHistory.length === 0" class="text-center py-8 text-gray-500">
            No stock transactions found
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="transaction in stockHistory"
              :key="transaction.id"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div class="flex items-center space-x-3">
                <div
                  :class="transaction.quantity > 0 ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'"
                  class="p-2 rounded-full"
                >
                  <font-awesome-icon
                    :icon="transaction.quantity > 0 ? 'fa-solid fa-plus' : 'fa-solid fa-minus'"
                    class="w-4 h-4"
                  />
                </div>
                <div>
                  <div class="font-medium text-gray-900 capitalize">
                    {{ transaction.transaction_type }}
                  </div>
                  <div class="text-sm text-gray-500">
                    {{ formatDate(transaction.created_at) }}
                  </div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-medium text-gray-900">
                  {{ transaction.quantity > 0 ? '+' : '' }}{{ transaction.quantity }}
                </div>
                <div class="text-sm text-gray-500">
                  {{ transaction.reason || 'No reason provided' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Quick Actions -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
          <div class="space-y-3">
            <router-link
              :to="`/inventory/${item.id}/restock`"
              class="w-full btn btn-primary flex items-center justify-center"
            >
              <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
              Restock Item
            </router-link>
            <router-link
              :to="`/inventory/${item.id}/edit`"
              class="w-full btn btn-secondary flex items-center justify-center"
            >
              <font-awesome-icon icon="fa-solid fa-edit" class="w-4 h-4 mr-2" />
              Edit Item
            </router-link>
          </div>
        </div>

        <!-- Item Information -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Item Information</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-700">SKU</label>
              <div class="text-sm text-gray-900">{{ item.sku }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Category</label>
              <div class="text-sm text-gray-900 capitalize">{{ item.category }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Unit of Measure</label>
              <div class="text-sm text-gray-900">{{ item.unit_of_measure }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Type</label>
              <div class="text-sm text-gray-900 capitalize">{{ item.type }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Created</label>
              <div class="text-sm text-gray-900">{{ formatDate(item.created_at) }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Last Updated</label>
              <div class="text-sm text-gray-900">{{ formatDate(item.updated_at) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'

export default {
  name: 'InventoryItemView',
  setup() {
    const route = useRoute()
    const authStore = useAuthStore()

    // Reactive data
    const item = ref(null)
    const stockHistory = ref([])
    const loading = ref(false)
    const error = ref('')

    // Computed properties
    const backRoute = computed(() => {
      // Check if we're coming from platform inventory
      if (route.path.includes('/admin/platform-inventory/')) {
        return '/admin/platform-inventory'
      }
      return '/inventory'
    })

    // Methods
    const loadItem = async () => {
      loading.value = true
      error.value = ''

      try {
        const itemId = route.params.id
        let endpoint = `/api/inventory/items/${itemId}`

        // Use platform inventory endpoint if we're in admin section
        if (route.path.includes('/admin/platform-inventory/')) {
          endpoint = `/api/inventory/platform/${itemId}`
        }

        const response = await apiService.get(endpoint)

        if (response.success) {
          item.value = response.data.item || response.data
          stockHistory.value = response.data.stock_history || []
        } else {
          throw new Error(response.error || 'Failed to load item')
        }
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load item'
        console.error('Failed to load item:', err)
      } finally {
        loading.value = false
      }
    }

    const getImageUrl = (imagePath) => {
      if (!imagePath) return ''
      return `/uploads/${imagePath}`
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

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // Lifecycle
    onMounted(() => {
      loadItem()
    })

    return {
      // Data
      item,
      stockHistory,
      loading,
      error,

      // Computed
      backRoute,

      // Methods
      loadItem,
      getImageUrl,
      getStatusBadgeClass,
      getStatusText,
      formatDate
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
</style>