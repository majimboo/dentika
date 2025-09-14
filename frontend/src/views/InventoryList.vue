<template>
  <div class="inventory-list-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Inventory Management</h1>
        <p class="text-gray-600">Track and manage clinic consumables and supplies</p>
      </div>

       <div class="header-actions flex items-center space-x-3">
         <router-link
           to="/inventory/new"
           class="btn btn-secondary flex items-center"
         >
           <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
           Add Item
         </router-link>
         <router-link
           to="/platform-inventory"
           class="btn btn-primary flex items-center"
         >
           <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-4 h-4 mr-2" />
           Order Supplies
         </router-link>
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
              placeholder="Search inventory by name, SKU, or supplier..."
              class="w-full md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          </div>

          <div class="filters-row flex items-center space-x-4">
            <!-- Category Filter -->
            <select
              v-model="categoryFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Categories</option>
              <option value="consumables">Consumables</option>
              <option value="equipment">Equipment</option>
              <option value="medications">Medications</option>
              <option value="supplies">Supplies</option>
              <option value="disposables">Disposables</option>
            </select>

             <!-- Status Filter -->
             <select
               v-model="statusFilter"
               @change="applyFilters"
               class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
             >
               <option value="">All Status</option>
               <option value="active">Active</option>
               <option value="inactive">Inactive</option>
               <option value="discontinued">Discontinued</option>
             </select>

             <!-- Type Filter -->
             <select
               v-model="typeFilter"
               @change="applyFilters"
               class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
             >
               <option value="">All Types</option>
               <option value="clinic">Clinic Inventory</option>
               <option value="platform">Platform Inventory</option>
             </select>

            <!-- Low Stock Filter -->
            <label class="flex items-center space-x-2">
              <input
                v-model="lowStockOnly"
                @change="applyFilters"
                type="checkbox"
                class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              >
              <span class="text-sm text-gray-700">Low Stock Only</span>
            </label>

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

        <!-- View Options -->
        <div class="view-options flex items-center space-x-2">
          <button
            @click="viewMode = 'grid'"
            :class="viewMode === 'grid' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
          >
            <font-awesome-icon icon="fa-solid fa-th" class="w-5 h-5" />
          </button>
          <button
            @click="viewMode = 'list'"
            :class="viewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
          >
            <font-awesome-icon icon="fa-solid fa-list" class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Inventory Stats -->
    <div class="inventory-stats grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-blue-600">{{ totalItems }}</div>
        <div class="stat-label text-sm text-gray-600">Total Items</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow cursor-pointer" @click="showLowStockOnly = true">
        <div class="stat-value text-2xl font-bold text-red-600">{{ lowStockItems }}</div>
        <div class="stat-label text-sm text-gray-600">Low Stock Items</div>
        <div class="text-xs text-red-500 mt-1">Click to view</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-orange-600">{{ expiringItems }}</div>
        <div class="stat-label text-sm text-gray-600">Expiring Soon</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-green-600">₱{{ totalValue.toFixed(2) }}</div>
        <div class="stat-label text-sm text-gray-600">Total Value</div>
      </div>
    </div>

    <!-- Low Stock Alert -->
    <div v-if="lowStockItems > 0 && !showLowStockOnly" class="low-stock-alert bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="flex items-center justify-center w-10 h-10 bg-red-100 rounded-full">
            <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-5 h-5 text-red-600" />
          </div>
          <div>
            <div class="text-sm font-medium text-red-900">{{ lowStockItems }} items are running low on stock</div>
            <div class="text-sm text-red-700">Consider reordering these items to avoid shortages</div>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="showLowStockOnly = true"
            class="btn btn-secondary text-sm"
          >
            View Low Stock
          </button>
          <router-link
            to="/platform-inventory"
            class="btn btn-primary text-sm"
          >
            Order Supplies
          </router-link>
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

    <!-- Inventory Items Grid/List View -->
    <div v-else-if="items.length > 0" class="inventory-items">
      <!-- Grid View -->
      <div v-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div
          v-for="item in items"
          :key="item.id"
          class="inventory-item-card bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
        >
          <!-- Item Image -->
          <div class="item-image h-48 bg-gray-100 flex items-center justify-center relative">
            <img
              v-if="item.image_path"
              :src="getImageUrl(item.image_path)"
              :alt="item.name"
              class="w-full h-full object-cover"
            >
            <div v-else class="text-gray-400">
              <font-awesome-icon icon="fa-solid fa-box" class="w-16 h-16" />
            </div>
            <!-- Stock Status Badge -->
            <div
              v-if="item.current_stock <= item.min_stock_level"
              class="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded-full"
            >
              Low Stock
            </div>
          </div>

          <!-- Item Details -->
          <div class="p-4">
            <h3 class="font-semibold text-gray-900 mb-1">{{ item.name }}</h3>
            <p class="text-sm text-gray-600 mb-2">{{ item.sku }}</p>

            <!-- Stock Info -->
            <div class="stock-info mb-3">
              <div class="flex justify-between items-center text-sm">
                <span class="text-gray-600">Stock:</span>
                <span :class="item.current_stock <= item.min_stock_level ? 'text-red-600 font-semibold' : 'text-gray-900'">
                  {{ item.current_stock }} {{ item.unit_of_measure }}
                </span>
              </div>
              <div class="flex justify-between items-center text-sm">
                <span class="text-gray-600">Min Level:</span>
                <span class="text-gray-900">{{ item.min_stock_level }}</span>
              </div>
            </div>

            <!-- Category Badge -->
            <div class="category-badge mb-3">
              <span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                {{ item.category }}
              </span>
            </div>

             <!-- Actions -->
             <div class="actions flex justify-between items-center">
                <router-link
                  :to="`/inventory/${item.id}/view`"
                  class="text-blue-600 hover:text-blue-800 text-sm font-medium"
                >
                  View Details
                </router-link>
               <router-link
                 :to="`/inventory/${item.id}/restock`"
                 class="text-green-600 hover:text-green-800 text-sm font-medium"
               >
                 Restock Item
               </router-link>
             </div>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Item</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">SKU</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10">
                      <img
                        v-if="item.image_path"
                        :src="getImageUrl(item.image_path)"
                        :alt="item.name"
                        class="h-10 w-10 rounded-lg object-cover"
                      >
                      <div v-else class="h-10 w-10 rounded-lg bg-gray-200 flex items-center justify-center">
                        <font-awesome-icon icon="fa-solid fa-box" class="w-5 h-5 text-gray-400" />
                      </div>
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ item.name }}</div>
                      <div class="text-sm text-gray-500">{{ item.supplier_name }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ item.sku }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                    {{ item.category }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-900">
                    {{ item.current_stock }} {{ item.unit_of_measure }}
                  </div>
                  <div class="text-xs text-gray-500">
                    Min: {{ item.min_stock_level }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="item.current_stock <= item.min_stock_level ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'"
                    class="inline-block text-xs px-2 py-1 rounded-full"
                  >
                    {{ item.current_stock <= item.min_stock_level ? 'Low Stock' : 'In Stock' }}
                  </span>
                </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                   <div class="flex items-center space-x-3">
                      <router-link
                        :to="`/inventory/${item.id}/view`"
                        class="text-blue-600 hover:text-blue-900"
                      >
                        View
                      </router-link>
                     <router-link
                       :to="`/inventory/${item.id}/restock`"
                       class="text-green-600 hover:text-green-900"
                     >
                       Restock
                     </router-link>
                   </div>
                 </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <font-awesome-icon icon="fa-solid fa-box-open" class="w-16 h-16 text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900 mb-2">No inventory items found</h3>
      <p class="text-gray-600 mb-6">Get started by adding your first inventory item.</p>
      <router-link
        to="/inventory/new"
        class="btn btn-primary inline-flex items-center"
      >
        <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
        Add First Item
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="items.length > 0 && totalPages > 1" class="pagination mt-6 flex justify-center">
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


  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useInventoryStore } from '../stores/inventory'
import { useAuthStore } from '../stores/auth'
import debounce from 'lodash/debounce'

export default {
  name: 'InventoryList',
  setup() {
    const inventoryStore = useInventoryStore()
    const authStore = useAuthStore()

    // Reactive data
    const searchQuery = ref('')
    const categoryFilter = ref('')
    const statusFilter = ref('')
    const typeFilter = ref('')
    const lowStockOnly = ref(false)
    const showLowStockOnly = ref(false)
    const viewMode = ref('grid')
    const currentPage = ref(1)
    const limit = ref(20)



    // Computed properties
    const items = computed(() => inventoryStore.items)
    const totalItems = computed(() => inventoryStore.getTotalItems)
    const totalPages = computed(() => Math.ceil(inventoryStore.getTotalItems / limit.value))
    const loading = computed(() => inventoryStore.loading)
    const error = computed(() => inventoryStore.error)

    // Use server-provided stats instead of client-side calculations
    const lowStockItems = computed(() => inventoryStore.getLowStockItems)
    const expiringItems = computed(() => inventoryStore.getExpiringItems)
    const totalValue = computed(() => inventoryStore.getTotalValue)

    const hasActiveFilters = computed(() => {
      return searchQuery.value || categoryFilter.value || statusFilter.value || typeFilter.value || lowStockOnly.value
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
      loadItems()
    }

    const clearFilters = () => {
      searchQuery.value = ''
      categoryFilter.value = ''
      statusFilter.value = ''
      typeFilter.value = ''
      lowStockOnly.value = false
      applyFilters()
    }

    const loadItems = () => {
      const params = {
        page: currentPage.value,
        limit: limit.value,
        search: searchQuery.value,
        category: categoryFilter.value,
        status: statusFilter.value,
        type: typeFilter.value,
        low_stock: lowStockOnly.value ? 'true' : ''
      }

      inventoryStore.fetchItems(params)
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadItems()
      }
    }

    const getImageUrl = (imagePath) => {
      if (!imagePath) return ''

      return `/uploads/${imagePath}`
    }



    // Watchers
    watch(showLowStockOnly, (newValue) => {
      if (newValue) {
        lowStockOnly.value = true
        applyFilters()
      }
    })

    // Lifecycle
    onMounted(() => {
      loadItems()
    })

    // Watchers
    watch([currentPage], () => {
      loadItems()
    })

    return {
      // Data
      searchQuery,
      categoryFilter,
      statusFilter,
      typeFilter,
      lowStockOnly,
      showLowStockOnly,
      viewMode,
      currentPage,


      // Computed
      items,
      totalItems,
      totalPages,
      loading,
      error,
      lowStockItems,
      expiringItems,
      totalValue,
      hasActiveFilters,
      visiblePages,


      // Methods
      debouncedSearch,
      applyFilters,
      clearFilters,
      loadItems,
      goToPage,
      getImageUrl,

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

.inventory-item-card {
  transition: all 0.2s ease-in-out;
}

.inventory-item-card:hover {
  transform: translateY(-2px);
}
</style>