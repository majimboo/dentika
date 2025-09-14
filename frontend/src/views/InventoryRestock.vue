<template>
  <div class="inventory-restock-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div class="flex items-center space-x-4">
        <router-link
          :to="backRoute"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to Item
        </router-link>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Restock Item</h1>
          <p class="text-gray-600">{{ item?.name }} ({{ item?.sku }})</p>
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
          <h3 class="text-sm font-medium text-red-800">Error loading item</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

     <!-- Restock Form -->
     <div v-else-if="item">
       <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-6">Restock Information</h2>

        <form @submit.prevent="submitRestock" class="space-y-6">
          <!-- Supplier Information -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Supplier Information</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Supplier Name *</label>
                <input
                  v-model="restockForm.supplier_name"
                  type="text"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter supplier name"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Supplier Email</label>
                <input
                  v-model="restockForm.supplier_email"
                  type="email"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="supplier@example.com"
                >
              </div>
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">Supplier Phone</label>
                <input
                  v-model="restockForm.supplier_phone"
                  type="tel"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="+63 912 345 6789"
                >
              </div>
            </div>
          </div>

          <!-- Restock Details -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Restock Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Quantity *</label>
                <input
                  v-model.number="restockForm.quantity_ordered"
                  type="number"
                  min="1"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter quantity"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Unit Cost (₱) *</label>
                <input
                  v-model.number="restockForm.unit_cost"
                  type="number"
                  step="0.01"
                  min="0"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="0.00"
                >
              </div>
            </div>

            <!-- Total Cost Display -->
            <div class="mt-4 p-3 bg-blue-50 rounded-lg">
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium text-gray-700">Total Cost:</span>
                <span class="text-lg font-bold text-blue-600">
                  ₱{{ ((restockForm.unit_cost || 0) * (restockForm.quantity_ordered || 0)).toFixed(2) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Order Details -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Order Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Order Date *</label>
                <input
                  v-model="restockForm.order_date"
                  type="date"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Expected Date</label>
                <input
                  v-model="restockForm.expected_date"
                  type="date"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
              </div>
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-2">Supplier Invoice Number</label>
                <input
                  v-model="restockForm.supplier_invoice"
                  type="text"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Invoice number from supplier"
                >
              </div>
            </div>
          </div>

          <!-- Notes -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Notes</label>
            <textarea
              v-model="restockForm.notes"
              rows="3"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Additional notes about this restock order"
            ></textarea>
          </div>

          <!-- Form Actions -->
          <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
            <router-link
              :to="backRoute"
              class="btn btn-secondary"
            >
              Cancel
            </router-link>
            <button
              type="submit"
              :disabled="!canSubmit"
              class="btn btn-primary"
            >
              <font-awesome-icon
                v-if="submitting"
                icon="fa-solid fa-spinner"
                class="w-4 h-4 mr-2 animate-spin"
              />
              Create Restock Order
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'

export default {
  name: 'InventoryRestock',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const authStore = useAuthStore()

    // Reactive data
    const itemId = route.params.id
    const item = ref(null)
    const loading = ref(false)
    const error = ref('')
    const submitting = ref(false)

    const restockForm = ref({
      supplier_name: '',
      supplier_email: '',
      supplier_phone: '',
      quantity_ordered: 0,
      unit_cost: 0,
      order_date: new Date().toISOString().split('T')[0], // Today's date
      expected_date: '',
      supplier_invoice: '',
      notes: ''
    })

    // Computed properties
    const backRoute = computed(() => {
      // Check if we're coming from platform inventory
      if (route.path.includes('/admin/platform-inventory/')) {
        return `/admin/platform-inventory/${itemId}/view`
      }
      return `/inventory/${itemId}/view`
    })

    const canSubmit = computed(() => {
      return restockForm.value.supplier_name &&
             restockForm.value.quantity_ordered > 0 &&
             restockForm.value.unit_cost >= 0 &&
             restockForm.value.order_date &&
             !submitting.value
    })

    // Methods
    const loadItem = async () => {
      loading.value = true
      error.value = ''

      try {
        let endpoint = `/api/inventory/items/${itemId}`

        // Use platform inventory endpoint if we're in admin section
        if (route.path.includes('/admin/platform-inventory/')) {
          endpoint = `/api/inventory/platform/${itemId}`
        }

        const response = await apiService.get(endpoint)

        if (response.success) {
          item.value = response.data.item || response.data
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

    const submitRestock = async () => {
      if (!canSubmit.value) return

      submitting.value = true

      try {
        const payload = {
          item_id: parseInt(itemId),
          ...restockForm.value,
          total_cost: (restockForm.value.unit_cost || 0) * (restockForm.value.quantity_ordered || 0)
        }

        const endpoint = `/api/inventory/${itemId}/restock`

        const response = await apiService.post(endpoint, payload)

        if (response.success) {
          // Redirect back to item view
          router.push(backRoute.value)
        } else {
          throw new Error(response.error || 'Failed to create restock order')
        }
      } catch (err) {
        alert('Failed to create restock order. Please try again.')
        console.error('Failed to create restock order:', err)
      } finally {
        submitting.value = false
      }
    }

    // Lifecycle
    onMounted(() => {
      loadItem()
    })

    return {
      // Data
      itemId,
      item,
      loading,
      error,
      submitting,
      restockForm,

      // Computed
      backRoute,
      canSubmit,

      // Methods
      loadItem,
      submitRestock
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
  @apply bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-700 hover:bg-gray-300;
}
</style>