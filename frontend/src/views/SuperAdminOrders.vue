<template>
  <div class="super-admin-orders-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Order Management</h1>
        <p class="text-gray-600">Manage all orders across clinics</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <button
          @click="exportOrders"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-download" class="w-4 h-4 mr-2" />
          Export Orders
        </button>
      </div>
    </div>

    <!-- Analytics Dashboard -->
    <div class="analytics-grid grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-blue-600">{{ totalOrders }}</div>
            <div class="text-sm text-gray-600">Total Orders</div>
          </div>
          <div class="p-2 bg-blue-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-5 h-5 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-green-600">₱{{ totalRevenue.toFixed(2) }}</div>
            <div class="text-sm text-gray-600">Total Revenue</div>
          </div>
          <div class="p-2 bg-green-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-dollar-sign" class="w-5 h-5 text-green-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-yellow-600">{{ pendingOrders }}</div>
            <div class="text-sm text-gray-600">Pending Orders</div>
          </div>
          <div class="p-2 bg-yellow-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-clock" class="w-5 h-5 text-yellow-600" />
          </div>
        </div>
      </div>

      <div class="analytics-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-purple-600">{{ deliveredOrders }}</div>
            <div class="text-sm text-gray-600">Delivered Orders</div>
          </div>
          <div class="p-2 bg-purple-100 rounded-lg">
            <font-awesome-icon icon="fa-solid fa-check-circle" class="w-5 h-5 text-purple-600" />
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
              placeholder="Search orders..."
              class="w-full md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          </div>

          <div class="filters-row flex flex-wrap items-center gap-4">
            <!-- Status Filter -->
            <select
              v-model="statusFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="confirmed">Confirmed</option>
              <option value="shipped">Shipped</option>
              <option value="delivered">Delivered</option>
              <option value="cancelled">Cancelled</option>
            </select>

            <!-- Clinic Filter -->
            <select
              v-model="clinicFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Clinics</option>
              <option v-for="clinic in clinics" :key="clinic.id" :value="clinic.id">
                {{ clinic.name }}
              </option>
            </select>

            <!-- Date Range -->
            <div class="flex items-center space-x-2">
              <input
                v-model="dateFrom"
                @change="applyFilters"
                type="date"
                class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
              <span class="text-gray-500">to</span>
              <input
                v-model="dateTo"
                @change="applyFilters"
                type="date"
                class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
            </div>

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

    <!-- Bulk Actions -->
    <div v-if="selectedOrders.length > 0" class="bulk-actions bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-sm font-medium text-blue-900">
            {{ selectedOrders.length }} order{{ selectedOrders.length > 1 ? 's' : '' }} selected
          </span>
          <button
            @click="selectAllOrders"
            class="text-sm text-blue-600 hover:text-blue-800"
          >
            Select all {{ orders.length }}
          </button>
        </div>
        <div class="flex items-center space-x-2">
          <select
            v-model="bulkStatus"
            class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Update Status</option>
            <option value="confirmed">Mark as Confirmed</option>
            <option value="shipped">Mark as Shipped</option>
            <option value="delivered">Mark as Delivered</option>
            <option value="cancelled">Cancel Orders</option>
          </select>
          <button
            @click="applyBulkStatusUpdate"
            :disabled="!bulkStatus"
            class="btn btn-primary text-sm"
          >
            Apply
          </button>
          <button
            @click="clearSelection"
            class="btn btn-secondary text-sm"
          >
            Clear
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

    <!-- Orders Table -->
    <div v-else-if="orders.length > 0" class="orders-table bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <input
                  type="checkbox"
                  @change="toggleSelectAll"
                  :checked="selectedOrders.length === orders.length && orders.length > 0"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                >
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order #</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Clinic</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Items</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="order in orders"
              :key="order.id"
              :class="{ 'bg-blue-50': selectedOrders.includes(order.id) }"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <input
                  type="checkbox"
                  @change="toggleOrderSelection(order.id)"
                  :checked="selectedOrders.includes(order.id)"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                >
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ order.order_number }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ order.clinic?.name || 'N/A' }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ formatDate(order.order_date) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  :class="getStatusBadgeClass(order.status)"
                  class="inline-block text-xs px-2 py-1 rounded-full font-medium"
                >
                  {{ getStatusText(order.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ order.order_items?.length || 0 }} items</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">₱{{ order.total_amount?.toFixed(2) || '0.00' }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center space-x-2">
                  <button
                    @click="viewOrderDetails(order)"
                    class="text-blue-600 hover:text-blue-900"
                  >
                    View
                  </button>
                  <button
                    @click="updateOrderStatus(order, 'confirmed')"
                    v-if="order.status === 'pending'"
                    class="text-green-600 hover:text-green-900"
                  >
                    Confirm
                  </button>
                  <button
                    @click="updateOrderStatus(order, 'shipped')"
                    v-if="order.status === 'confirmed'"
                    class="text-orange-600 hover:text-orange-900"
                  >
                    Ship
                  </button>
                  <button
                    @click="updateOrderStatus(order, 'delivered')"
                    v-if="order.status === 'shipped'"
                    class="text-purple-600 hover:text-purple-900"
                  >
                    Deliver
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
      <font-awesome-icon icon="fa-solid fa-shopping-cart" class="w-16 h-16 text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900 mb-2">No orders found</h3>
      <p class="text-gray-600 mb-6">There are currently no orders matching your criteria.</p>
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

    <!-- Order Details Modal -->
    <div v-if="showOrderModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeOrderModal">
      <div class="relative top-20 mx-auto p-5 border w-11/12 max-w-4xl shadow-lg rounded-md bg-white max-h-screen overflow-y-auto" @click.stop>
        <div class="mt-3">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">Order Details - {{ selectedOrder?.order_number }}</h3>
            <button @click="closeOrderModal" class="text-gray-400 hover:text-gray-600">
              <font-awesome-icon icon="fa-solid fa-times" class="w-5 h-5" />
            </button>
          </div>

          <div v-if="selectedOrder" class="space-y-6">
            <!-- Order Header -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Order Number</label>
                <div class="mt-1 text-sm text-gray-900">{{ selectedOrder.order_number }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Order Date</label>
                <div class="mt-1 text-sm text-gray-900">{{ formatDate(selectedOrder.order_date) }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Status</label>
                <span
                  :class="getStatusBadgeClass(selectedOrder.status)"
                  class="inline-block text-sm px-3 py-1 rounded-full font-medium mt-1"
                >
                  {{ getStatusText(selectedOrder.status) }}
                </span>
              </div>
            </div>

            <!-- Clinic Info -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Clinic Information</label>
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ selectedOrder.clinic?.name }}</div>
                    <div class="text-sm text-gray-600">{{ selectedOrder.clinic?.address }}</div>
                    <div class="text-sm text-gray-600">{{ selectedOrder.clinic?.phone }}</div>
                  </div>
                  <div>
                    <div class="text-sm text-gray-600">{{ selectedOrder.clinic?.email }}</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Order Items -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Order Items</label>
              <div class="bg-white border border-gray-200 rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                    <tr>
                      <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Item</th>
                      <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">SKU</th>
                      <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Qty</th>
                      <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Unit Price</th>
                      <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-200">
                    <tr v-for="item in selectedOrder.order_items" :key="item.id">
                      <td class="px-4 py-3 text-sm text-gray-900">{{ item.item?.name }}</td>
                      <td class="px-4 py-3 text-sm text-gray-600">{{ item.item?.sku }}</td>
                      <td class="px-4 py-3 text-sm text-gray-900">{{ item.quantity }}</td>
                      <td class="px-4 py-3 text-sm text-gray-900">₱{{ item.unit_price?.toFixed(2) }}</td>
                      <td class="px-4 py-3 text-sm font-medium text-gray-900">₱{{ item.line_total?.toFixed(2) }}</td>
                    </tr>
                  </tbody>
                  <tfoot class="bg-gray-50">
                    <tr>
                      <td colspan="4" class="px-4 py-3 text-sm font-medium text-gray-900 text-right">Total:</td>
                      <td class="px-4 py-3 text-sm font-bold text-gray-900">₱{{ selectedOrder.total_amount?.toFixed(2) }}</td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            </div>

            <!-- Shipping & Notes -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Shipping Address</label>
                <div class="text-sm text-gray-900">{{ selectedOrder.shipping_address || 'N/A' }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Shipping Notes</label>
                <div class="text-sm text-gray-900">{{ selectedOrder.shipping_notes || 'N/A' }}</div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Order Notes</label>
              <div class="text-sm text-gray-900">{{ selectedOrder.notes || 'N/A' }}</div>
            </div>

            <!-- Tracking Info -->
            <div v-if="selectedOrder.tracking_number">
              <label class="block text-sm font-medium text-gray-700 mb-2">Tracking Information</label>
              <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <div class="flex items-center space-x-2">
                  <font-awesome-icon icon="fa-solid fa-truck" class="w-4 h-4 text-blue-600" />
                  <span class="text-sm font-medium text-blue-900">Tracking Number: {{ selectedOrder.tracking_number }}</span>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
              <button @click="closeOrderModal" class="btn btn-secondary">Close</button>
              <button
                v-if="selectedOrder.status === 'confirmed'"
                @click="showTrackingModal = true"
                class="btn btn-primary"
              >
                Add Tracking
              </button>
              <button
                v-if="selectedOrder.status === 'pending'"
                @click="updateOrderStatus(selectedOrder, 'confirmed')"
                class="btn btn-success"
              >
                Confirm Order
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tracking Modal -->
    <div v-if="showTrackingModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeTrackingModal">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Add Tracking Information</h3>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Tracking Number</label>
            <input
              v-model="trackingNumber"
              type="text"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter tracking number"
            >
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Notes</label>
            <textarea
              v-model="trackingNotes"
              rows="3"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Optional notes about the shipment"
            ></textarea>
          </div>

          <div class="flex justify-end space-x-3">
            <button @click="closeTrackingModal" class="btn btn-secondary">Cancel</button>
            <button @click="submitTrackingInfo" class="btn btn-primary" :disabled="!trackingNumber">
              Add Tracking & Ship
            </button>
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
  name: 'SuperAdminOrders',
  setup() {
    const authStore = useAuthStore()

    // Reactive data
    const orders = ref([])
    const clinics = ref([])
    const loading = ref(false)
    const error = ref('')
    const currentPage = ref(1)
    const limit = ref(20)
    const totalPages = ref(0)

    // Filters
    const searchQuery = ref('')
    const statusFilter = ref('')
    const clinicFilter = ref('')
    const dateFrom = ref('')
    const dateTo = ref('')

    // Selection
    const selectedOrders = ref([])
    const bulkStatus = ref('')

    // Modals
    const showOrderModal = ref(false)
    const showTrackingModal = ref(false)
    const selectedOrder = ref(null)
    const trackingNumber = ref('')
    const trackingNotes = ref('')

    // Computed properties
    const totalOrders = computed(() => orders.value.length)
    const totalRevenue = computed(() => {
      return orders.value.reduce((sum, order) => sum + (order.total_amount || 0), 0)
    })
    const pendingOrders = computed(() => {
      return orders.value.filter(order => order.status === 'pending').length
    })
    const deliveredOrders = computed(() => {
      return orders.value.filter(order => order.status === 'delivered').length
    })

    const hasActiveFilters = computed(() => {
      return searchQuery.value || statusFilter.value || clinicFilter.value || dateFrom.value || dateTo.value
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
      loadOrders()
    }

    const clearFilters = () => {
      searchQuery.value = ''
      statusFilter.value = ''
      clinicFilter.value = ''
      dateFrom.value = ''
      dateTo.value = ''
      applyFilters()
    }

    const loadOrders = async () => {
      loading.value = true
      error.value = ''

      try {
        const params = {
          page: currentPage.value,
          limit: limit.value,
          search: searchQuery.value,
          status: statusFilter.value,
          clinic_id: clinicFilter.value,
          date_from: dateFrom.value,
          date_to: dateTo.value
        }

        const response = await apiService.get('/api/inventory/orders', { params })
        if (response.success) {
          orders.value = response.data.orders || []
          totalPages.value = Math.ceil((response.data.total || 0) / limit.value)
        } else {
          error.value = response.error || 'Failed to load orders'
        }
      } catch (err) {
        error.value = 'Failed to load orders'
        console.error('Failed to load orders:', err)
      } finally {
        loading.value = false
      }
    }

    const loadClinics = async () => {
      try {
        const response = await apiService.getClinics()
        if (response.success) {
          clinics.value = response.data || []
        } else {
          console.error('Failed to load clinics:', response.error)
        }
      } catch (err) {
        console.error('Failed to load clinics:', err)
      }
    }

    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadOrders()
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
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

    const toggleOrderSelection = (orderId) => {
      const index = selectedOrders.value.indexOf(orderId)
      if (index > -1) {
        selectedOrders.value.splice(index, 1)
      } else {
        selectedOrders.value.push(orderId)
      }
    }

    const toggleSelectAll = () => {
      if (selectedOrders.value.length === orders.value.length) {
        selectedOrders.value = []
      } else {
        selectedOrders.value = orders.value.map(order => order.id)
      }
    }

    const selectAllOrders = () => {
      selectedOrders.value = orders.value.map(order => order.id)
    }

    const clearSelection = () => {
      selectedOrders.value = []
      bulkStatus.value = ''
    }

    const applyBulkStatusUpdate = async () => {
      if (!bulkStatus.value || selectedOrders.value.length === 0) return

      if (!confirm(`Are you sure you want to update ${selectedOrders.value.length} order(s) to ${getStatusText(bulkStatus.value)}?`)) {
        return
      }

      loading.value = true
      let successCount = 0
      let errorCount = 0

      for (const orderId of selectedOrders.value) {
        try {
          const response = await apiService.put(`/api/inventory/orders/${orderId}/status`, {
            status: bulkStatus.value,
            notes: `Bulk status update to ${getStatusText(bulkStatus.value)}`
          })
          if (response.success) {
            successCount++
          } else {
            errorCount++
            console.error(`Failed to update order ${orderId}:`, response.error)
          }
        } catch (err) {
          errorCount++
          console.error(`Failed to update order ${orderId}:`, err)
        }
      }

      loading.value = false
      clearSelection()

      if (successCount > 0) {
        alert(`Successfully updated ${successCount} order(s). ${errorCount > 0 ? `${errorCount} failed.` : ''}`)
        loadOrders()
      } else {
        alert('Failed to update any orders.')
      }
    }

    const updateOrderStatus = async (order, newStatus) => {
      const statusText = getStatusText(newStatus)
      if (!confirm(`Are you sure you want to update order ${order.order_number} to ${statusText}?`)) {
        return
      }

      try {
        if (newStatus === 'shipped') {
          // For shipped status, we need tracking info
          selectedOrder.value = order
          showTrackingModal.value = true
          return
        }

        const response = await apiService.put(`/api/inventory/orders/${order.id}/status`, {
          status: newStatus,
          notes: `Status updated to ${statusText}`
        })
        if (!response.success) {
          throw new Error(response.error || 'Failed to update order status')
        }

        alert(`Order ${order.order_number} has been updated to ${statusText}`)
        loadOrders()
      } catch (err) {
        alert('Failed to update order status. Please try again.')
        console.error('Failed to update order status:', err)
      }
    }

    const submitTrackingInfo = async () => {
      if (!trackingNumber.value || !selectedOrder.value) return

      try {
        const response = await apiService.put(`/api/inventory/orders/${selectedOrder.value.id}/status`, {
          status: 'shipped',
          tracking_number: trackingNumber.value,
          notes: trackingNotes.value || `Order shipped with tracking number: ${trackingNumber.value}`
        })
        if (!response.success) {
          throw new Error(response.error || 'Failed to update tracking information')
        }

        alert(`Order ${selectedOrder.value.order_number} has been marked as shipped with tracking number: ${trackingNumber.value}`)
        closeTrackingModal()
        loadOrders()
      } catch (err) {
        alert('Failed to update tracking information. Please try again.')
        console.error('Failed to update tracking info:', err)
      }
    }

    const viewOrderDetails = (order) => {
      selectedOrder.value = order
      showOrderModal.value = true
    }

    const closeOrderModal = () => {
      showOrderModal.value = false
      selectedOrder.value = null
    }

    const closeTrackingModal = () => {
      showTrackingModal.value = false
      trackingNumber.value = ''
      trackingNotes.value = ''
      selectedOrder.value = null
    }

    const exportOrders = () => {
      // TODO: Implement export functionality
      alert('Export functionality will be implemented soon')
    }

    // Lifecycle
    onMounted(async () => {
      await loadClinics()
      loadOrders()
    })

    // Watchers
    watch([currentPage], () => {
      loadOrders()
    })

    return {
      // Data
      orders,
      clinics,
      loading,
      error,
      currentPage,
      totalPages,
      searchQuery,
      statusFilter,
      clinicFilter,
      dateFrom,
      dateTo,
      selectedOrders,
      bulkStatus,
      showOrderModal,
      showTrackingModal,
      selectedOrder,
      trackingNumber,
      trackingNotes,

      // Computed
      totalOrders,
      totalRevenue,
      pendingOrders,
      deliveredOrders,
      hasActiveFilters,
      visiblePages,

      // Methods
      debouncedSearch,
      applyFilters,
      clearFilters,
      loadOrders,
      goToPage,
      formatDate,
      getStatusBadgeClass,
      getStatusText,
      toggleOrderSelection,
      toggleSelectAll,
      selectAllOrders,
      clearSelection,
      applyBulkStatusUpdate,
      updateOrderStatus,
      submitTrackingInfo,
      viewOrderDetails,
      closeOrderModal,
      closeTrackingModal,
      exportOrders
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

.btn-success {
  @apply bg-green-600 text-white hover:bg-green-700;
}

.analytics-card {
  transition: all 0.2s ease-in-out;
}

.analytics-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.order-row {
  transition: all 0.2s ease-in-out;
}

.order-row:hover {
  background-color: #f9fafb;
}
</style>