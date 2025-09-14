<template>
  <div class="shop-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Dentika Shop</h1>
        <p class="text-gray-600">Browse and order items from Dentika's inventory</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/inventory"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to My Inventory
        </router-link>
        <router-link
          to="/shop/orders"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-4 h-4 mr-2" />
          View My Orders
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
              placeholder="Search Dentika shop..."
              class="w-full md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          </div>

          <div class="filters-row flex items-center space-x-4">
            <!-- Category Filter -->
            <select
              v-model="categoryFilter"
              @change="applyFilters"
              class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Categories</option>
              <option value="consumables">Consumables</option>
              <option value="equipment">Equipment</option>
              <option value="medications">Medications</option>
              <option value="supplies">Supplies</option>
              <option value="disposables">Disposables</option>
            </select>
          </div>
        </div>

        <!-- View Toggle -->
        <div class="view-toggle flex items-center border border-gray-200 rounded-lg p-1">
          <button
            @click="viewMode = 'grid'"
            :class="[
              'px-3 py-1 rounded-md text-sm transition-colors',
              viewMode === 'grid' ? 'bg-blue-500 text-white' : 'text-gray-600 hover:text-gray-900'
            ]"
          >
            <font-awesome-icon icon="fa-solid fa-th" class="w-4 h-4" />
          </button>
          <button
            @click="viewMode = 'list'"
            :class="[
              'px-3 py-1 rounded-md text-sm transition-colors',
              viewMode === 'list' ? 'bg-blue-500 text-white' : 'text-gray-600 hover:text-gray-900'
            ]"
          >
            <font-awesome-icon icon="fa-solid fa-list" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      <p class="mt-2 text-gray-600">Loading shop items...</p>
    </div>

    <!-- Items Display -->
    <div v-else-if="items.length > 0">
      <!-- Grid View -->
      <div v-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
        <div v-for="item in items" :key="item.id" class="shop-item-card bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow">
          <div class="item-image-container relative">
            <img
              v-if="item.image_path"
              :src="item.image_path"
              :alt="item.name"
              class="w-full h-48 object-cover rounded-t-lg"
            >
            <div v-else class="w-full h-48 bg-gray-100 rounded-t-lg flex items-center justify-center">
              <font-awesome-icon icon="fa-solid fa-box" class="w-12 h-12 text-gray-400" />
            </div>
            <div class="absolute top-2 right-2">
              <span :class="[
                'px-2 py-1 text-xs font-medium rounded-full',
                item.current_stock > item.min_stock_level ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              ]">
                {{ item.current_stock }} in stock
              </span>
            </div>
          </div>

          <div class="p-4">
            <div class="mb-2">
              <h3 class="font-semibold text-gray-900 mb-1">{{ item.name }}</h3>
              <p class="text-sm text-gray-600 line-clamp-2">{{ item.description }}</p>
            </div>

            <div class="flex justify-between items-center mb-3">
              <span class="text-lg font-bold text-blue-600">₱{{ formatPrice(item.selling_price) }}</span>
              <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">{{ formatCategory(item.category) }}</span>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.unit_of_measure }}</span>
              <router-link
                :to="{ name: 'ShopOrder', params: { id: item.id } }"
                :class="[
                  'btn btn-sm',
                  item.current_stock > 0 ? 'btn-primary' : 'btn-disabled'
                ]"
                :disabled="item.current_stock === 0"
              >
                <font-awesome-icon icon="fa-solid fa-cart-plus" class="w-4 h-4 mr-1" />
                {{ item.current_stock > 0 ? 'Order' : 'Out of Stock' }}
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 mb-8">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Item</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-12 w-12">
                      <img v-if="item.image_path" :src="item.image_path" :alt="item.name" class="h-12 w-12 rounded-lg object-cover">
                      <div v-else class="h-12 w-12 bg-gray-100 rounded-lg flex items-center justify-center">
                        <font-awesome-icon icon="fa-solid fa-box" class="w-6 h-6 text-gray-400" />
                      </div>
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ item.name }}</div>
                      <div class="text-sm text-gray-500">{{ item.description }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs font-medium bg-gray-100 text-gray-800 rounded-full">
                    {{ formatCategory(item.category) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                  ₱{{ formatPrice(item.selling_price) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="[
                    'px-2 py-1 text-xs font-medium rounded-full',
                    item.current_stock > item.min_stock_level ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                  ]">
                    {{ item.current_stock }} {{ item.unit_of_measure }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <router-link
                    :to="{ name: 'ShopOrder', params: { id: item.id } }"
                    :class="[
                      'btn btn-sm',
                      item.current_stock > 0 ? 'btn-primary' : 'btn-disabled'
                    ]"
                    :disabled="item.current_stock === 0"
                  >
                    <font-awesome-icon icon="fa-solid fa-cart-plus" class="w-4 h-4 mr-1" />
                    {{ item.current_stock > 0 ? 'Order' : 'Out of Stock' }}
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination-container flex justify-center items-center space-x-2 mt-6">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="btn btn-secondary btn-sm"
        >
          <font-awesome-icon icon="fa-solid fa-chevron-left" class="w-4 h-4" />
        </button>

        <span class="px-4 py-2 text-sm text-gray-700">
          Page {{ currentPage }} of {{ totalPages }}
        </span>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="btn btn-secondary btn-sm"
        >
          <font-awesome-icon icon="fa-solid fa-chevron-right" class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state text-center py-12">
      <div class="inline-block p-4 bg-gray-100 rounded-full mb-4">
        <font-awesome-icon icon="fa-solid fa-store" class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No items found</h3>
      <p class="text-gray-600 mb-4">
        {{ searchQuery ? 'Try adjusting your search or filters.' : 'No items are currently available in the shop.' }}
      </p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import apiService from '../services/api'

export default {
  name: 'Shop',
  setup() {
    const items = ref([])
    const loading = ref(false)
    const error = ref(null)

    // Filters and search
    const searchQuery = ref('')
    const categoryFilter = ref('')
    const viewMode = ref('grid')

    // Pagination
    const currentPage = ref(1)
    const totalItems = ref(0)
    const totalPages = ref(0)
    const itemsPerPage = 20

    // Search debounce
    let searchTimeout = null

    const debouncedSearch = () => {
      clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        currentPage.value = 1
        fetchItems()
      }, 300)
    }

    const applyFilters = () => {
      currentPage.value = 1
      fetchItems()
    }

    const fetchItems = async () => {
      loading.value = true
      error.value = null

      try {
        const params = {
          page: currentPage.value,
          limit: itemsPerPage,
          search: searchQuery.value,
          category: categoryFilter.value
        }

        const response = await apiService.get('/api/shop', { params })

        if (response.success) {
          items.value = response.data.items || []
          totalItems.value = response.data.total || 0
          totalPages.value = response.data.total_pages || 0
          currentPage.value = response.data.page || 1
        } else {
          error.value = response.error || 'Failed to load shop items'
        }
      } catch (err) {
        console.error('Error fetching shop items:', err)
        error.value = 'Failed to load shop items'
      } finally {
        loading.value = false
      }
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        fetchItems()
      }
    }

    const formatPrice = (price) => {
      return new Intl.NumberFormat('en-PH', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(price)
    }

    const formatCategory = (category) => {
      return category.charAt(0).toUpperCase() + category.slice(1).replace('_', ' ')
    }

    onMounted(() => {
      fetchItems()
    })

    return {
      items,
      loading,
      error,
      searchQuery,
      categoryFilter,
      viewMode,
      currentPage,
      totalItems,
      totalPages,
      debouncedSearch,
      applyFilters,
      goToPage,
      formatPrice,
      formatCategory
    }
  }
}
</script>

<style scoped>
@import "../styles/main.css";

.shop-item-card {
  transition: all 0.2s ease-in-out;
}

.shop-item-card:hover {
  transform: translateY(-2px);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.btn {
  @apply inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-100 text-gray-700 hover:bg-gray-200 focus:ring-gray-500 border-gray-300;
}

.btn-disabled {
  @apply bg-gray-100 text-gray-400 cursor-not-allowed;
}

.btn-sm {
  @apply px-3 py-1.5 text-xs;
}
</style>