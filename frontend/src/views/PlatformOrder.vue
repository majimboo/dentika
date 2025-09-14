<template>
  <div class="platform-order-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div class="flex items-center space-x-4">
        <router-link
          to="/platform-inventory"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to Platform Inventory
        </router-link>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Order Supplies</h1>
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

    <!-- Order Form -->
    <div v-else-if="item" class="max-w-2xl">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-6">Order Details</h2>

        <form @submit.prevent="submitOrder" class="space-y-6">
          <!-- Item Information -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Item Information</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Item Name</label>
                <div class="mt-1 text-sm text-gray-900 font-medium">{{ item.name }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">SKU</label>
                <div class="mt-1 text-sm text-gray-900">{{ item.sku }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Category</label>
                <div class="mt-1 text-sm text-gray-900 capitalize">{{ item.category }}</div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Unit Price</label>
                <div class="mt-1 text-sm font-medium text-green-600">₱{{ item.selling_price?.toFixed(2) }}</div>
              </div>
            </div>
          </div>

          <!-- Order Details -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Order Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Quantity *</label>
                <input
                  v-model.number="orderForm.quantity"
                  type="number"
                  min="1"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter quantity"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Order Date *</label>
                <input
                  v-model="orderForm.order_date"
                  type="date"
                  required
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
              </div>
            </div>

            <!-- Total Cost Display -->
            <div class="mt-4 p-3 bg-blue-50 rounded-lg">
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium text-gray-700">Total Cost:</span>
                <span class="text-lg font-bold text-blue-600">
                  ₱{{ ((item.selling_price || 0) * (orderForm.quantity || 0)).toFixed(2) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Shipping Information -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Shipping Information</h3>
            <div class="grid grid-cols-1 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Shipping Address</label>
                <textarea
                  v-model="orderForm.shipping_address"
                  rows="3"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter shipping address"
                ></textarea>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Shipping Notes</label>
                <textarea
                  v-model="orderForm.shipping_notes"
                  rows="2"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Any special delivery instructions"
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Additional Notes -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Additional Notes</label>
            <textarea
              v-model="orderForm.notes"
              rows="3"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Additional notes about this order"
            ></textarea>
          </div>

          <!-- Form Actions -->
          <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
            <router-link
              to="/platform-inventory"
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
              <font-awesome-icon
                v-else
                icon="fa-solid fa-shopping-cart"
                class="w-4 h-4 mr-2"
              />
              Place Order
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
  name: 'PlatformOrder',
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

    const orderForm = ref({
      quantity: 1,
      order_date: new Date().toISOString().split('T')[0], // Today's date
      shipping_address: '',
      shipping_notes: '',
      notes: ''
    })

    // Computed properties
    const canSubmit = computed(() => {
      return orderForm.value.quantity > 0 &&
             orderForm.value.order_date &&
             !submitting.value
    })

    // Methods
    const loadItem = async () => {
      loading.value = true
      error.value = ''

      try {
        const response = await apiService.get(`/api/inventory/items/${itemId}`)

        if (response.success) {
          item.value = response.data
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

    const submitOrder = async () => {
      if (!canSubmit.value) return

      submitting.value = true

      try {
        const orderData = {
          items: [{
            item_id: parseInt(itemId),
            quantity: orderForm.value.quantity
          }],
          shipping_address: orderForm.value.shipping_address,
          shipping_notes: orderForm.value.shipping_notes,
          notes: orderForm.value.notes
        }

        const response = await apiService.post('/api/inventory/orders', orderData)

        if (response.success) {
          // Redirect back to platform inventory
          router.push('/platform-inventory')
        } else {
          throw new Error(response.error || 'Failed to create order')
        }
      } catch (err) {
        alert('Failed to create order. Please try again.')
        console.error('Failed to create order:', err)
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
      orderForm,

      // Computed
      canSubmit,

      // Methods
      loadItem,
      submitOrder
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