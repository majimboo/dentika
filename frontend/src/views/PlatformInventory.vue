<template>
  <div class="platform-inventory-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Platform Inventory</h1>
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
          to="/orders"
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
              placeholder="Search platform inventory..."
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
          <h3 class="text-sm font-medium text-red-800">Error loading platform inventory</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Inventory Items Grid -->
    <div v-else-if="items.length > 0" class="inventory-items">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
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
            <!-- Platform Badge -->
            <div class="absolute top-2 right-2 bg-green-500 text-white text-xs px-2 py-1 rounded-full">
              Platform
            </div>
          </div>

          <!-- Item Details -->
          <div class="p-4">
            <h3 class="font-semibold text-gray-900 mb-1">{{ item.name }}</h3>
            <p class="text-sm text-gray-600 mb-2">{{ item.sku }}</p>

            <!-- Price Info -->
            <div class="price-info mb-3">
              <div class="flex justify-between items-center text-sm">
                <span class="text-gray-600">Price:</span>
                <span class="text-lg font-bold text-green-600">₱{{ item.selling_price.toFixed(2) }}</span>
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
              <button
                @click="viewItemDetails(item)"
                class="text-blue-600 hover:text-blue-800 text-sm font-medium"
              >
                View Details
              </button>
               <router-link
                 :to="`/platform-inventory/${item.id}/order`"
                 class="btn btn-primary text-sm flex items-center"
               >
                 <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-4 h-4 mr-1" />
                 Order Now
               </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <font-awesome-icon icon="fa-solid fa-box-open" class="w-16 h-16 text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900 mb-2">No platform inventory items found</h3>
      <p class="text-gray-600 mb-6">There are currently no items available in the platform inventory.</p>
      <router-link
        to="/orders"
        class="btn btn-primary inline-flex items-center"
      >
        <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-4 h-4 mr-2" />
        View My Orders
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
import { useRouter } from 'vue-router'
import { useInventoryStore } from '../stores/inventory'
import { useAuthStore } from '../stores/auth'
import debounce from 'lodash/debounce'

export default {
  name: 'PlatformInventory',
  setup() {
    const router = useRouter()
    const inventoryStore = useInventoryStore()
    const authStore = useAuthStore()

    // Reactive data
    const searchQuery = ref('')
    const categoryFilter = ref('')
    const currentPage = ref(1)
    const limit = ref(20)



    // Computed properties
    const items = computed(() => inventoryStore.platformItems)
    const totalPages = computed(() => Math.ceil(inventoryStore.platformPagination.total / limit.value))
    const loading = computed(() => inventoryStore.loading)
    const error = computed(() => inventoryStore.error)

    const hasActiveFilters = computed(() => {
      return searchQuery.value || categoryFilter.value
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
      applyFilters()
    }

    const loadItems = () => {
      const params = {
        page: currentPage.value,
        limit: limit.value,
        search: searchQuery.value,
        category: categoryFilter.value
      }

      inventoryStore.fetchPlatformItems(params)
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

    const viewItemDetails = (item) => {
      // Navigate to platform inventory item view page
      router.push(`/admin/platform-inventory/${item.id}/view`)
    }





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
      currentPage,

      // Computed
      items,
      totalPages,
      loading,
      error,
      hasActiveFilters,
      visiblePages,

      // Methods
      debouncedSearch,
      applyFilters,
      clearFilters,
      loadItems,
      goToPage,
      getImageUrl,
      viewItemDetails
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