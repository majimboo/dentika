<template>
  <div class="inventory-view-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Inventory Item Details</h1>
        <p v-if="item" class="text-gray-600">
          {{ item.name }} - {{ item.sku }}
        </p>
        <p v-else class="text-gray-600">Loading item information...</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/inventory"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to Inventory
        </router-link>
        <router-link
          v-if="item"
          :to="`/inventory/${item.id}/edit`"
          class="btn btn-primary flex items-center"
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
    <div v-else-if="item" class="space-y-6">
      <!-- Item Overview Card -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <div class="flex flex-col lg:flex-row gap-6">
            <!-- Item Image -->
            <div class="flex-shrink-0">
              <div class="w-32 h-32 bg-gray-100 rounded-lg flex items-center justify-center">
                <img
                  v-if="item.image_path"
                  :src="getImageUrl(item.image_path)"
                  :alt="item.name"
                  class="w-full h-full object-cover rounded-lg"
                >
                <font-awesome-icon
                  v-else
                  icon="fa-solid fa-box"
                  class="w-16 h-16 text-gray-400"
                />
              </div>
            </div>

            <!-- Item Basic Info -->
            <div class="flex-1">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4">
                <div>
                  <h2 class="text-2xl font-bold text-gray-900">{{ item.name }}</h2>
                  <p class="text-gray-600">{{ item.sku }}</p>
                </div>
                <div class="mt-2 sm:mt-0">
                  <span
                    :class="getStatusBadgeClass(item)"
                    class="inline-block text-sm px-3 py-1 rounded-full font-medium"
                  >
                    {{ getStatusText(item) }}
                  </span>
                </div>
              </div>

              <p v-if="item.description" class="text-gray-700 mb-4">{{ item.description }}</p>

              <!-- Key Stats -->
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="text-center">
                  <div class="text-2xl font-bold text-blue-600">{{ item.current_stock }}</div>
                  <div class="text-sm text-gray-600">Current Stock</div>
                </div>
                <div class="text-center">
                  <div class="text-2xl font-bold text-gray-600">{{ item.min_stock_level }}</div>
                  <div class="text-sm text-gray-600">Min Level</div>
                </div>
                 <div class="text-center">
                  <div class="text-2xl font-bold text-green-600">₱{{ (item.average_unit_cost || 0).toFixed(2) }}</div>
                  <div class="text-sm text-gray-600">Avg Unit Cost</div>
                </div>
                <div class="text-center">
                  <div class="text-2xl font-bold text-purple-600">₱{{ item.selling_price.toFixed(2) }}</div>
                  <div class="text-sm text-gray-600">Selling Price</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Detailed Information -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Product Details -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200">
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Product Details</h3>
             <dl class="space-y-3">
               <div class="flex justify-between">
                 <dt class="text-sm font-medium text-gray-600">Type:</dt>
                 <dd class="text-sm text-gray-900">
                   <span class="inline-block bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
                     {{ item.type === 'platform' ? 'Platform Inventory' : 'Clinic Inventory' }}
                   </span>
                 </dd>
               </div>
               <div class="flex justify-between">
                 <dt class="text-sm font-medium text-gray-600">Category:</dt>
                 <dd class="text-sm text-gray-900">
                   <span class="inline-block bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs">
                     {{ item.category }}
                   </span>
                 </dd>
               </div>
              <div class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Unit of Measure:</dt>
                <dd class="text-sm text-gray-900">{{ item.unit_of_measure }}</dd>
              </div>
              <div v-if="item.has_expiration" class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Expiry Date:</dt>
                <dd class="text-sm text-gray-900">
                  <span :class="isExpired(item) ? 'text-red-600' : isExpiringSoon(item) ? 'text-orange-600' : 'text-gray-900'">
                    {{ formatDate(item.expiry_date) }}
                  </span>
                </dd>
              </div>
              <div class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Reorder Point:</dt>
                <dd class="text-sm text-gray-900">{{ item.reorder_point }}</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Supplier Information -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200">
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Supplier Information</h3>
            <dl class="space-y-3">
              <div v-if="item.supplier_name" class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Supplier:</dt>
                <dd class="text-sm text-gray-900">{{ item.supplier_name }}</dd>
              </div>
              <div v-if="item.supplier_sku" class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Supplier SKU:</dt>
                <dd class="text-sm text-gray-900">{{ item.supplier_sku }}</dd>
              </div>
              <div v-if="item.supplier_email" class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Email:</dt>
                <dd class="text-sm text-gray-900">{{ item.supplier_email }}</dd>
              </div>
              <div v-if="item.supplier_phone" class="flex justify-between">
                <dt class="text-sm font-medium text-gray-600">Phone:</dt>
                <dd class="text-sm text-gray-900">{{ item.supplier_phone }}</dd>
              </div>
            </dl>
            <div v-if="!item.supplier_name && !item.supplier_sku && !item.supplier_email && !item.supplier_phone" class="text-sm text-gray-500">
              No supplier information available
            </div>
          </div>
        </div>
      </div>

      <!-- Stock Transaction History -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Stock Transaction History</h3>
            <button
              @click="loadStockTransactions"
              class="btn btn-secondary text-sm"
              :disabled="loadingTransactions"
            >
              <font-awesome-icon
                v-if="loadingTransactions"
                icon="fa-solid fa-spinner"
                class="w-4 h-4 mr-2 animate-spin"
              />
              Refresh
            </button>
          </div>

          <div v-if="loadingTransactions" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>

          <div v-else-if="stockTransactions.length === 0" class="text-center py-8 text-gray-500">
            No stock transactions found
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Quantity</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Reason</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Performed By</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="transaction in stockTransactions" :key="transaction.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ formatDate(transaction.created_at) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span
                      :class="getTransactionTypeClass(transaction.transaction_type)"
                      class="inline-block text-xs px-2 py-1 rounded-full font-medium"
                    >
                      {{ transaction.transaction_type }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span :class="transaction.quantity > 0 ? 'text-green-600' : 'text-red-600'">
                      {{ transaction.quantity > 0 ? '+' : '' }}{{ transaction.quantity }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ transaction.reason || 'N/A' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ transaction.user?.first_name }} {{ transaction.user?.last_name }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useInventoryStore } from '../stores/inventory'
import { formatDate, getInventoryStatusClass } from '@/utils'

export default {
  name: 'InventoryView',
  setup() {
    const route = useRoute()
    const inventoryStore = useInventoryStore()

    // Reactive data
    const loading = ref(false)
    const loadingTransactions = ref(false)
    const error = ref('')
    const item = ref(null)
    const stockTransactions = ref([])

    // Computed properties
    const itemId = route.params.id

    // Methods
    const loadItem = async () => {
      loading.value = true
      error.value = ''

      try {
        const response = await inventoryStore.fetchItem(itemId)
        if (response.success) {
          item.value = response.data
          await loadStockTransactions()
        } else {
          error.value = response.error
        }
      } catch (err) {
        error.value = err.message || 'Failed to load item'
      } finally {
        loading.value = false
      }
    }

    const loadStockTransactions = async () => {
      loadingTransactions.value = true

      try {
        const response = await inventoryStore.fetchStockTransactions(itemId)
        if (response.success) {
          stockTransactions.value = response.data.transactions
        }
      } catch (err) {
        console.error('Failed to load stock transactions:', err)
      } finally {
        loadingTransactions.value = false
      }
    }

    const getImageUrl = (imagePath) => {
      if (!imagePath) return ''
      const cleanPath = imagePath.startsWith('/') ? imagePath.substring(1) : imagePath
      return `/uploads/${cleanPath}`
    }


    const getStatusBadgeClass = (item) => {
      const status = item.current_stock <= item.min_stock_level ? 'low_stock' : 'in_stock'
      return getInventoryStatusClass(status)
    }

    const getStatusText = (item) => {
      if (item.current_stock <= item.min_stock_level) {
        return 'Low Stock'
      }
      return 'In Stock'
    }

    const isExpired = (item) => {
      if (!item.has_expiration || !item.expiry_date) return false
      return new Date(item.expiry_date) < new Date()
    }

    const isExpiringSoon = (item) => {
      if (!item.has_expiration || !item.expiry_date) return false
      const expiryDate = new Date(item.expiry_date)
      const now = new Date()
      const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24))
      return daysUntilExpiry <= 30 && daysUntilExpiry > 0
    }

    const getTransactionTypeClass = (type) => {
      switch (type) {
        case 'restock':
          return 'bg-green-100 text-green-800'
        case 'usage':
          return 'bg-red-100 text-red-800'
        case 'adjustment':
          return 'bg-yellow-100 text-yellow-800'
        case 'transfer':
          return 'bg-blue-100 text-blue-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    // Lifecycle
    onMounted(() => {
      loadItem()
    })

    return {
      // Data
      loading,
      loadingTransactions,
      error,
      item,
      stockTransactions,

      // Computed
      itemId,

      // Methods
      loadItem,
      loadStockTransactions,
      getImageUrl,
      getStatusBadgeClass,
      getStatusText,
      isExpired,
      isExpiringSoon,
      getTransactionTypeClass,
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