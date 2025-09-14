<template>
  <div class="super-admin-platform-inventory-add-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Add Platform Inventory Item</h1>
        <p class="text-gray-600">Add a new item to Dentika's platform inventory catalog</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/admin/platform-inventory"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to Inventory
        </router-link>
      </div>
    </div>

    <!-- Form Container -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <form @submit.prevent="saveItem" class="space-y-8">
        <!-- Basic Information -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Item Name *</label>
              <input
                v-model="itemForm.name"
                type="text"
                required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter item name"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">SKU *</label>
              <input
                v-model="itemForm.sku"
                type="text"
                required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter SKU"
              >
            </div>
          </div>

          <div class="mt-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
            <textarea
              v-model="itemForm.description"
              rows="4"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter item description"
            ></textarea>
          </div>
        </div>

        <!-- Category and Status -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Category & Status</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Category *</label>
              <select
                v-model="itemForm.category"
                required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Select Category</option>
                <option value="equipment">Equipment</option>
                <option value="supplies">Supplies</option>
                <option value="consumables">Consumables</option>
                <option value="instruments">Instruments</option>
                <option value="materials">Materials</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
              <select
                v-model="itemForm.status"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="discontinued">Discontinued</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Pricing and Stock -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Pricing & Stock</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Price (₱) *</label>
              <input
                v-model.number="itemForm.price"
                type="number"
                step="0.01"
                min="0"
                required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="0.00"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Low Stock Threshold</label>
              <input
                v-model.number="itemForm.low_stock_threshold"
                type="number"
                min="0"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="10"
              >
            </div>
          </div>
        </div>

        <!-- Unit Information -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Unit Information</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
             <div>
               <label class="block text-sm font-medium text-gray-700 mb-2">Unit of Measure</label>
               <select
                 v-model="itemForm.unit"
                 class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
               >
                 <option value="">Select Unit</option>
                 <option value="pieces">Pieces</option>
                 <option value="boxes">Boxes</option>
                 <option value="bottles">Bottles</option>
                 <option value="packs">Packs</option>
                 <option value="tubes">Tubes</option>
                 <option value="vials">Vials</option>
                 <option value="kits">Kits</option>
               </select>
             </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
          <router-link
            to="/admin/platform-inventory"
            class="btn btn-secondary"
          >
            Cancel
          </router-link>
          <button
            type="submit"
            :disabled="saving"
            class="btn btn-primary"
          >
            <span v-if="saving">Creating Item...</span>
            <span v-else>Create Item</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api'

export default {
  name: 'SuperAdminPlatformInventoryAdd',
  setup() {
    const router = useRouter()
    const saving = ref(false)

    // Form data
    const itemForm = ref({
      name: '',
      sku: '',
      description: '',
      category: '',
      status: 'active',
      price: 0,
      low_stock_threshold: 10,
      unit: ''
    })

    const saveItem = async () => {
      saving.value = true

      try {
        const result = await apiService.post('/api/inventory/platform', itemForm.value)
        if (result.success) {
          alert('Item created successfully!')
          router.push('/admin/platform-inventory')
        } else {
          alert(`Failed to create item: ${result.error}`)
        }
      } catch (err) {
        alert('Failed to create item. Please try again.')
        console.error('Failed to create item:', err)
      } finally {
        saving.value = false
      }
    }

    return {
      itemForm,
      saving,
      saveItem
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