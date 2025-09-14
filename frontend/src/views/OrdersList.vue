<template>
  <div class="orders-list-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">My Orders</h1>
        <p class="text-gray-600">Track and manage your orders from Dentika</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/platform-inventory"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          Order More Items
        </router-link>
      </div>
    </div>

    <!-- Order Stats -->
    <div class="order-stats grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-blue-600">{{ pendingOrders.length }}</div>
        <div class="stat-label text-sm text-gray-600">Pending Orders</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-yellow-600">{{ confirmedOrders.length }}</div>
        <div class="stat-label text-sm text-gray-600">Confirmed Orders</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-orange-600">{{ shippedOrders.length }}</div>
        <div class="stat-label text-sm text-gray-600">Shipped Orders</div>
      </div>

      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-green-600">{{ deliveredOrders.length }}</div>
        <div class="stat-label text-sm text-gray-600">Delivered Orders</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="filters-section bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div class="filters-row flex items-center space-x-4">
          <!-- Status Filter -->
          <select
            v-model="statusFilter"
            @change="applyFilters"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Orders</option>
            <option value="pending">Pending</option>
            <option value="confirmed">Confirmed</option>
            <option value="shipped">Shipped</option>
            <option value="delivered">Delivered</option>
            <option value="cancelled">Cancelled</option>
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

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-center">
        <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-5 h-5 text-red-600 mr-3" />
        <div>
          <h3 class="text-sm font-medium text-red-800">Error loading orders</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Orders List -->
    <div v-else-if="orders.length > 0" class="orders-list space-y-4">
      <div
        v-for="order in orders"
        :key="order.id"
        class="order-card bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="p-6">
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Order #{{ order.order_number }}</h3>
              <p class="text-sm text-gray-600">{{ formatDate(order.order_date) }}</p>
            </div>
            <div class="mt-2 lg:mt-0 flex items-center space-x-4">
              <span
                :class="getStatusBadgeClass(order.status)"
                class="inline-block text-sm px-3 py-1 rounded-full font-medium"
              >
                {{ getStatusText(order.status) }}
              </span>
              <div class="text-right">
                <div class="text-lg font-bold text-gray-900">₱{{ order.total_amount.toFixed(2) }}</div>
                <div class="text-sm text-gray-600">{{ order.order_items.length }} items</div>
              </div>
            </div>
          </div>

          <!-- Order Items -->
          <div class="order-items mb-4">
            <div class="space-y-2">
              <div
                v-for="item in order.items"
                :key="item.id"
                class="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0"
              >
                <div class="flex items-center space-x-3">
                  <div class="flex-shrink-0 w-8 h-8 bg-gray-100 rounded-lg flex items-center justify-center">
                    <font-awesome-icon icon="fa-solid fa-box" class="w-4 h-4 text-gray-400" />
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ item.item.name }}</div>
                    <div class="text-xs text-gray-600">{{ item.item.sku }}</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-sm text-gray-900">{{ item.quantity }} × ₱{{ item.unit_price.toFixed(2) }}</div>
                  <div class="text-sm font-medium text-gray-900">₱{{ item.line_total.toFixed(2) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Order Actions -->
          <div class="flex justify-between items-center">
            <div class="text-sm text-gray-600">
              <span v-if="order.tracking_number">Tracking: {{ order.tracking_number }}</span>
            </div>
            <div class="flex space-x-2">
              <button
                @click="viewOrderDetails(order)"
                class="btn btn-secondary text-sm"
              >
                View Details
              </button>
              <button
                v-if="order.status === 'pending'"
                @click="cancelOrder(order)"
                class="btn btn-danger text-sm"
              >
                Cancel Order
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-16 h-16 text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900 mb-2">No orders found</h3>
      <p class="text-gray-600 mb-6">You haven't placed any orders yet.</p>
      <router-link
        to="/platform-inventory"
        class="btn btn-primary inline-flex items-center"
      >
        <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
        Browse Platform Inventory
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="orders.length > 0 && totalPages > 1" class="pagination mt-6 flex justify-center">
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

export default {
  name: 'OrdersList',
  setup() {
    const inventoryStore = useInventoryStore()
    const authStore = useAuthStore()

    // Reactive data
    const statusFilter = ref('')
    const currentPage = ref(1)
    const limit = ref(20)

    // Computed properties
    const orders = computed(() => inventoryStore.orders)
    const totalPages = computed(() => Math.ceil(inventoryStore.ordersPagination.total / limit.value))
    const loading = computed(() => inventoryStore.loading)
    const error = computed(() => inventoryStore.error)

    const pendingOrders = computed(() => inventoryStore.getPendingOrders)
    const confirmedOrders = computed(() => inventoryStore.orders.filter(order => order.status === 'confirmed'))
    const shippedOrders = computed(() => inventoryStore.getShippedOrders)
    const deliveredOrders = computed(() => inventoryStore.getDeliveredOrders)

    const hasActiveFilters = computed(() => {
      return statusFilter.value
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
    const applyFilters = () => {
      currentPage.value = 1
      loadOrders()
    }

    const clearFilters = () => {
      statusFilter.value = ''
      applyFilters()
    }

    const loadOrders = () => {
      const params = {
        page: currentPage.value,
        limit: limit.value,
        status: statusFilter.value
      }

      inventoryStore.fetchOrders(params)
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadOrders()
      }
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString()
    }

    const getStatusBadgeClass = (status) => {
      const classes = {
        pending: 'bg-yellow-100 text-yellow-800',
        confirmed: 'bg-blue-100 text-blue-800',
        shipped: 'bg-orange-100 text-orange-800',
        delivered: 'bg-green-100 text-green-800',
        cancelled: 'bg-red-100 text-red-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    }

    const getStatusText = (status) => {
      const texts = {
        pending: 'Pending',
        confirmed: 'Confirmed',
        shipped: 'Shipped',
        delivered: 'Delivered',
        cancelled: 'Cancelled'
      }
      return texts[status] || status
    }

    const viewOrderDetails = (order) => {
      // For now, just show an alert. In a real app, you'd navigate to a detail page
      alert(`Viewing details for order: ${order.order_number}`)
    }

    const cancelOrder = async (order) => {
      if (!confirm('Are you sure you want to cancel this order?')) return

      try {
        // In a real implementation, you'd call an API to cancel the order
        alert('Order cancellation not implemented yet')
      } catch (error) {
        console.error('Failed to cancel order:', error)
        alert('Failed to cancel order. Please try again.')
      }
    }

    // Lifecycle
    onMounted(() => {
      loadOrders()
    })

    // Watchers
    watch([currentPage], () => {
      loadOrders()
    })

    return {
      // Data
      statusFilter,
      currentPage,

      // Computed
      orders,
      totalPages,
      loading,
      error,
      pendingOrders,
      confirmedOrders,
      shippedOrders,
      deliveredOrders,
      hasActiveFilters,
      visiblePages,

      // Methods
      applyFilters,
      clearFilters,
      loadOrders,
      goToPage,
      formatDate,
      getStatusBadgeClass,
      getStatusText,
      viewOrderDetails,
      cancelOrder
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

.btn-danger {
  @apply bg-red-600 text-white hover:bg-red-700;
}

.order-card {
  transition: all 0.2s ease-in-out;
}

.order-card:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style>