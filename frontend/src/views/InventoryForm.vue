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

            <!-- Type -->
            <div>
              <label for="type" class="block text-sm font-medium text-gray-700 mb-2">
                Inventory Type <span class="text-red-500">*</span>
              </label>
              <select
                id="type"
                v-model="form.type"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.type }"
              >
                <option value="clinic">Clinic Inventory</option>
                <option value="platform">Platform Inventory</option>
              </select>
              <span v-if="errors.type" class="text-red-500 text-sm mt-1">{{ errors.type }}</span>
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

          <div class="grid grid-cols-1 md:grid-cols-1 gap-6">
            <!-- Selling Price -->
            <div>
              <label for="selling_price" class="block text-sm font-medium text-gray-700 mb-2">
                Selling Price (₱) <span class="text-red-500">*</span>
              </label>
              <input
                id="selling_price"
                v-model.number="form.selling_price"
                type="number"
                step="0.01"
                min="0"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.selling_price }"
                placeholder="0.00"
              >
              <span v-if="errors.selling_price" class="text-red-500 text-sm mt-1">{{ errors.selling_price }}</span>
              <p class="text-xs text-gray-500 mt-1">
                For clinic inventory: price you charge patients. For platform inventory: price clinics pay when ordering.
              </p>
            </div>
          </div>
        </div>

        <!-- Stock Management Section -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2 mb-6">
            Stock Management
          </h3>

           <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
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
                placeholder="+63 917 234 5678"
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
import { useClinicStore } from '../stores/clinic'
import { useNavigation } from '../composables/useNavigation'

export default {
  name: 'InventoryForm',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const { goBack } = useNavigation()
    const inventoryStore = useInventoryStore()
    const authStore = useAuthStore()
    const clinicStore = useClinicStore()

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
      type: 'clinic',
      unit_of_measure: 'pieces',
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
        if (response.success) {
          item.value = response.data

          // Populate form
          form.value = {
            name: response.data.name || '',
            description: response.data.description || '',
            sku: response.data.sku || '',
            category: response.data.category || '',
            type: response.data.type || 'clinic',
            unit_of_measure: response.data.unit_of_measure || 'pieces',
            selling_price: response.data.selling_price || 0,
            min_stock_level: response.data.min_stock_level || 10,
            reorder_point: response.data.reorder_point || 5,
            current_stock: response.data.current_stock || 0,
            supplier_name: response.data.supplier_name || '',
            supplier_sku: response.data.supplier_sku || '',
            supplier_email: response.data.supplier_email || '',
            supplier_phone: response.data.supplier_phone || '',
            has_expiration: response.data.has_expiration || false,
            expiry_date: response.data.expiry_date ? new Date(response.data.expiry_date).toISOString().split('T')[0] : '',
            status: response.data.status || 'active'
          }

          // Set image preview if exists
          if (response.data.image_path) {
            imagePreview.value = getImageUrl(response.data.image_path)
          }
        } else {
          error.value = response.error || 'Failed to load item'
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

      if (!form.value.type) {
        errors.value.type = 'Inventory type is required'
      }

      if (form.value.selling_price <= 0) {
        errors.value.selling_price = 'Selling price must be greater than 0'
      }

      return Object.keys(errors.value).length === 0
    }

    const getBranchId = async () => {
      // If user is super admin, they can choose any branch
      if (authStore.isSuperAdmin) {
        // For now, default to branch 1 for super admin
        return 1
      }

      // For regular users, get branches for their clinic
      const clinicId = authStore.userClinicId
      if (!clinicId) {
        throw new Error('User is not assigned to any clinic')
      }

      // Fetch branches if not already loaded
      if (clinicStore.branches.length === 0) {
        await clinicStore.fetchBranches(clinicId)
      }

      // Use the first active branch
      const activeBranches = clinicStore.getActiveBranches
      if (activeBranches.length === 0) {
        throw new Error('No active branches found for your clinic')
      }

      return activeBranches[0].id
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

        // Get the branch ID
        const branchId = await getBranchId()

        // Only send fields that the server expects for create request
        const itemData = {
          name: form.value.name,
          description: form.value.description,
          sku: form.value.sku,
          category: form.value.category,
          type: form.value.type,
          unit_of_measure: form.value.unit_of_measure,
          selling_price: form.value.selling_price,
          min_stock_level: form.value.min_stock_level,
          reorder_point: form.value.reorder_point,
          supplier_name: form.value.supplier_name,
          supplier_sku: form.value.supplier_sku,
          supplier_email: form.value.supplier_email,
          supplier_phone: form.value.supplier_phone,
          has_expiration: form.value.has_expiration,
          expiry_date: form.value.expiry_date ? new Date(form.value.expiry_date) : null,
          branch_id: form.value.type === 'clinic' ? branchId : undefined
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

      // Remove leading slash if present to avoid double slashes
      const cleanPath = imagePath.startsWith('/') ? imagePath.substring(1) : imagePath
      return `/uploads/${cleanPath}`
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
      getImageUrl,
      goBack
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