<template>
  <div class="inventory-form-page">
    <!-- Page Header -->
    <div class="page-header mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">
            {{ isEditing ? 'Edit Inventory Item' : 'Add New Inventory Item' }}
          </h1>
          <p class="text-gray-600 mt-1">
            {{ isEditing ? 'Update item details and stock information' : 'Create a new inventory item for your clinic' }}
          </p>
        </div>
         <button
           @click="goBack"
           class="btn btn-secondary inline-flex items-center"
         >
           <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
           Back to Inventory
         </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error && !item" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-center">
        <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-5 h-5 text-red-600 mr-3" />
        <div>
          <h3 class="text-sm font-medium text-red-800">Error loading item</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Form -->
    <form v-else @submit.prevent="handleSubmit" class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="p-6">
        <!-- Basic Information Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Basic Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Item Name -->
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
                Item Name <span class="text-red-500">*</span>
              </label>
              <input
                id="name"
                v-model="form.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.name }"
                placeholder="Enter item name"
              >
              <span v-if="errors.name" class="text-red-500 text-sm mt-1">{{ errors.name }}</span>
            </div>

            <!-- SKU -->
            <div>
              <label for="sku" class="block text-sm font-medium text-gray-700 mb-2">
                SKU (Stock Keeping Unit)
              </label>
              <input
                id="sku"
                v-model="form.sku"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Auto-generated if empty"
              >
            </div>

            <!-- Category -->
            <div>
              <label for="category" class="block text-sm font-medium text-gray-700 mb-2">
                Category <span class="text-red-500">*</span>
              </label>
              <select
                id="category"
                v-model="form.category"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.category }"
              >
                <option value="">Select Category</option>
                <option value="consumables">Consumables</option>
                <option value="equipment">Equipment</option>
                <option value="medications">Medications</option>
                <option value="supplies">Supplies</option>
                <option value="disposables">Disposables</option>
              </select>
              <span v-if="errors.category" class="text-red-500 text-sm mt-1">{{ errors.category }}</span>
            </div>

            <!-- Unit of Measure -->
            <div>
              <label for="unit_of_measure" class="block text-sm font-medium text-gray-700 mb-2">
                Unit of Measure
              </label>
              <select
                id="unit_of_measure"
                v-model="form.unit_of_measure"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
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

          <!-- Description -->
          <div class="mt-6">
            <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              id="description"
              v-model="form.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter item description"
            ></textarea>
          </div>
        </div>

        <!-- Pricing Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Pricing Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Unit Cost -->
            <div>
              <label for="unit_cost" class="block text-sm font-medium text-gray-700 mb-2">
                Unit Cost (₱)
              </label>
              <input
                id="unit_cost"
                v-model.number="form.unit_cost"
                type="number"
                step="0.01"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="0.00"
              >
            </div>

            <!-- Selling Price -->
            <div>
              <label for="selling_price" class="block text-sm font-medium text-gray-700 mb-2">
                Selling Price (₱)
              </label>
              <input
                id="selling_price"
                v-model.number="form.selling_price"
                type="number"
                step="0.01"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="0.00"
              >
            </div>
          </div>
        </div>

        <!-- Stock Management Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Stock Management
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Current Stock -->
            <div>
              <label for="current_stock" class="block text-sm font-medium text-gray-700 mb-2">
                Current Stock
              </label>
              <input
                id="current_stock"
                v-model.number="form.current_stock"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="0"
              >
            </div>

            <!-- Minimum Stock Level -->
            <div>
              <label for="min_stock_level" class="block text-sm font-medium text-gray-700 mb-2">
                Minimum Stock Level
              </label>
              <input
                id="min_stock_level"
                v-model.number="form.min_stock_level"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="10"
              >
            </div>

            <!-- Reorder Point -->
            <div>
              <label for="reorder_point" class="block text-sm font-medium text-gray-700 mb-2">
                Reorder Point
              </label>
              <input
                id="reorder_point"
                v-model.number="form.reorder_point"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="5"
              >
            </div>
          </div>
        </div>

        <!-- Supplier Information Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Supplier Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Supplier Name -->
            <div>
              <label for="supplier_name" class="block text-sm font-medium text-gray-700 mb-2">
                Supplier Name
              </label>
              <input
                id="supplier_name"
                v-model="form.supplier_name"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter supplier name"
              >
            </div>

            <!-- Supplier SKU -->
            <div>
              <label for="supplier_sku" class="block text-sm font-medium text-gray-700 mb-2">
                Supplier SKU
              </label>
              <input
                id="supplier_sku"
                v-model="form.supplier_sku"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Supplier's SKU for this item"
              >
            </div>

            <!-- Supplier Email -->
            <div>
              <label for="supplier_email" class="block text-sm font-medium text-gray-700 mb-2">
                Supplier Email
              </label>
              <input
                id="supplier_email"
                v-model="form.supplier_email"
                type="email"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="supplier@example.com"
              >
            </div>

            <!-- Supplier Phone -->
            <div>
              <label for="supplier_phone" class="block text-sm font-medium text-gray-700 mb-2">
                Supplier Phone
              </label>
              <input
                id="supplier_phone"
                v-model="form.supplier_phone"
                type="tel"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="+1 (555) 123-4567"
              >
            </div>
          </div>
        </div>

        <!-- Expiration & Image Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Additional Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Expiration Settings -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-4">
                Expiration Tracking
              </label>
              <div class="space-y-3">
                <label class="flex items-center">
                  <input
                    v-model="form.has_expiration"
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  >
                  <span class="ml-2 text-sm text-gray-700">This item has an expiration date</span>
                </label>

                <div v-if="form.has_expiration">
                  <label for="expiry_date" class="block text-sm font-medium text-gray-700 mb-2">
                    Expiry Date
                  </label>
                  <input
                    id="expiry_date"
                    v-model="form.expiry_date"
                    type="date"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                </div>
              </div>
            </div>

            <!-- Status -->
            <div v-if="isEditing">
              <label for="status" class="block text-sm font-medium text-gray-700 mb-2">
                Status
              </label>
              <select
                id="status"
                v-model="form.status"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="discontinued">Discontinued</option>
              </select>
            </div>
          </div>

          <!-- Image Upload -->
          <div class="mt-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Item Image
            </label>
            <div class="flex items-center space-x-4">
              <div class="flex-shrink-0">
                <img
                  v-if="imagePreview"
                  :src="imagePreview"
                  alt="Item preview"
                  class="w-20 h-20 object-cover rounded-lg border border-gray-300"
                >
                <div v-else class="w-20 h-20 bg-gray-100 rounded-lg border border-gray-300 flex items-center justify-center">
                  <font-awesome-icon icon="fa-solid fa-box" class="w-8 h-8 text-gray-400" />
                </div>
              </div>
              <div class="flex-1">
                <input
                  ref="imageInput"
                  type="file"
                  accept="image/*"
                  @change="handleImageChange"
                  class="hidden"
                >
                <button
                  type="button"
                  @click="$refs.imageInput.click()"
                  class="btn btn-secondary text-sm"
                >
                  <font-awesome-icon icon="fa-solid fa-upload" class="w-4 h-4 mr-2" />
                  Choose Image
                </button>
                <p class="text-xs text-gray-500 mt-1">
                  PNG, JPG, GIF up to 5MB
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
           <button
             type="button"
             @click="goBack"
             class="btn btn-secondary"
           >
             Cancel
           </button>
          <button
            type="submit"
            :disabled="submitting"
            class="btn btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="submitting" class="inline-flex items-center">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              Saving...
            </span>
            <span v-else>
              {{ isEditing ? 'Update Item' : 'Create Item' }}
            </span>
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useInventoryStore } from '../stores/inventory'
import { useAuthStore } from '../stores/auth'
import { useNavigation } from '../composables/useNavigation'

export default {
  name: 'InventoryForm',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const { goBack } = useNavigation()
    const inventoryStore = useInventoryStore()
    const authStore = useAuthStore()

    // Reactive data
    const loading = ref(false)
    const submitting = ref(false)
    const error = ref('')
    const item = ref(null)
    const imagePreview = ref('')
    const imageFile = ref(null)

    const form = ref({
      name: '',
      description: '',
      sku: '',
      category: '',
      unit_of_measure: 'pieces',
      unit_cost: 0,
      selling_price: 0,
      min_stock_level: 10,
      reorder_point: 5,
      current_stock: 0,
      supplier_name: '',
      supplier_sku: '',
      supplier_email: '',
      supplier_phone: '',
      has_expiration: false,
      expiry_date: '',
      status: 'active'
    })

    const errors = ref({})

    // Computed properties
    const isEditing = computed(() => !!route.params.id)
    const itemId = computed(() => route.params.id)

    // Methods
    const loadItem = async () => {
      if (!isEditing.value) return

      loading.value = true
      error.value = ''

      try {
        const response = await inventoryStore.fetchItem(itemId.value)
        item.value = response

        // Populate form
        form.value = {
          name: response.name || '',
          description: response.description || '',
          sku: response.sku || '',
          category: response.category || '',
          unit_of_measure: response.unit_of_measure || 'pieces',
          unit_cost: response.unit_cost || 0,
          selling_price: response.selling_price || 0,
          min_stock_level: response.min_stock_level || 10,
          reorder_point: response.reorder_point || 5,
          current_stock: response.current_stock || 0,
          supplier_name: response.supplier_name || '',
          supplier_sku: response.supplier_sku || '',
          supplier_email: response.supplier_email || '',
          supplier_phone: response.supplier_phone || '',
          has_expiration: response.has_expiration || false,
          expiry_date: response.expiry_date ? new Date(response.expiry_date).toISOString().split('T')[0] : '',
          status: response.status || 'active'
        }

        // Set image preview if exists
        if (response.image_path) {
          imagePreview.value = getImageUrl(response.image_path)
        }
      } catch (err) {
        error.value = err.message || 'Failed to load item'
      } finally {
        loading.value = false
      }
    }

    const handleImageChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        imageFile.value = file

        // Create preview
        const reader = new FileReader()
        reader.onload = (e) => {
          imagePreview.value = e.target.result
        }
        reader.readAsDataURL(file)
      }
    }

    const uploadImage = async () => {
      if (!imageFile.value) return null

      try {
        const formData = new FormData()
        formData.append('image', imageFile.value)

        const response = await fetch('/api/upload/inventory-item-image', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${authStore.token}`
          },
          body: formData
        })

        if (!response.ok) {
          throw new Error('Failed to upload image')
        }

        const result = await response.json()
        return result.path
      } catch (err) {
        console.error('Image upload failed:', err)
        return null
      }
    }

    const validateForm = () => {
      errors.value = {}

      if (!form.value.name.trim()) {
        errors.value.name = 'Item name is required'
      }

      if (!form.value.category) {
        errors.value.category = 'Category is required'
      }

      return Object.keys(errors.value).length === 0
    }

    const handleSubmit = async () => {
      if (!validateForm()) return

      submitting.value = true

      try {
        // Upload image if selected
        let imagePath = item.value?.image_path || null
        if (imageFile.value) {
          imagePath = await uploadImage()
        }

        const itemData = {
          ...form.value,
          image_path: imagePath,
          branch_id: authStore.user?.branch_id || 1 // Default to first branch if not set
        }

        if (isEditing.value) {
          await inventoryStore.updateItem(itemId.value, itemData)
        } else {
          await inventoryStore.createItem(itemData)
        }

        // Redirect to inventory list
        router.push('/inventory')
      } catch (err) {
        error.value = err.message || 'Failed to save item'
      } finally {
        submitting.value = false
      }
    }

    const getImageUrl = (imagePath) => {
      if (!imagePath) return ''

      // Build the full backend URL for uploads
      const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:9483'

      // Remove leading slash if present to avoid double slashes
      const cleanPath = imagePath.startsWith('/') ? imagePath.substring(1) : imagePath
      return `${baseUrl}/uploads/${cleanPath}`
    }

    // Lifecycle
    onMounted(() => {
      loadItem()
    })

    return {
      // Data
      loading,
      submitting,
      error,
      item,
      form,
      errors,
      imagePreview,
      imageFile,

      // Computed
      isEditing,
      itemId,

      // Methods
      loadItem,
      handleImageChange,
      handleSubmit,
      getImageUrl
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