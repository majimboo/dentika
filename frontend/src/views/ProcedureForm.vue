<template>
  <div class="procedure-form-page">
    <!-- Page Header -->
    <div class="page-header flex items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">
          {{ isEditing ? 'Edit Procedure' : 'New Procedure' }}
        </h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update procedure details' : 'Add a new dental procedure to your catalog' }}
        </p>
      </div>
    </div>

    <!-- Form -->
    <div class="form-container bg-white rounded-lg shadow-sm border border-gray-200">
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <!-- Basic Information -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Basic Information</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Procedure Name *</label>
              <input
                v-model="formData.name"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.name }"
                placeholder="e.g. Dental Cleaning"
                required
              />
              <div v-if="errors.name" class="error-message">{{ errors.name }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">Description</label>
              <textarea
                v-model="formData.description"
                rows="3"
                class="form-input resize-none"
                placeholder="Describe the procedure in detail..."
              ></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Category *</label>
                <select
                  v-model="formData.category"
                  class="form-input"
                  :class="{ 'border-red-500': errors.category }"
                  required
                >
                  <option value="">Select category...</option>
                  <option value="preventive">Preventive</option>
                  <option value="restorative">Restorative</option>
                  <option value="surgical">Surgical</option>
                  <option value="orthodontic">Orthodontic</option>
                  <option value="cosmetic">Cosmetic</option>
                  <option value="emergency">Emergency</option>
                </select>
                <div v-if="errors.category" class="error-message">{{ errors.category }}</div>
              </div>

              <div class="form-group">
                <label class="form-label">Duration (minutes) *</label>
                <input
                  v-model.number="formData.duration"
                  type="number"
                  min="5"
                  max="480"
                  class="form-input"
                  :class="{ 'border-red-500': errors.duration }"
                  placeholder="30"
                  required
                />
                <div v-if="errors.duration" class="error-message">{{ errors.duration }}</div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Cost (PHP) *</label>
                <div class="relative">
                  <input
                    v-model.number="formData.cost"
                    type="number"
                    min="0"
                    step="0.01"
                    class="form-input pl-8"
                    :class="{ 'border-red-500': errors.cost }"
                    placeholder="0.00"
                    required
                  />
                </div>
                <div v-if="errors.cost" class="error-message">{{ errors.cost }}</div>
              </div>

              <div class="form-group">
                <label class="form-label">Insurance Code</label>
                <input
                  v-model="formData.insurance_code"
                  type="text"
                  class="form-input"
                  placeholder="D1110"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Additional Settings -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Additional Settings</h3>
          
          <div class="space-y-4">
            <div class="form-check">
              <input
                v-model="formData.requires_anesthesia"
                type="checkbox"
                id="requires_anesthesia"
                class="form-checkbox"
              />
              <label for="requires_anesthesia" class="form-check-label">
                Requires anesthesia
              </label>
            </div>

            <div class="form-check">
              <input
                v-model="formData.requires_followup"
                type="checkbox"
                id="requires_followup"
                class="form-checkbox"
              />
              <label for="requires_followup" class="form-check-label">
                Requires follow-up appointment
              </label>
            </div>

            <div class="form-check">
              <input
                v-model="formData.is_active"
                type="checkbox"
                id="is_active"
                class="form-checkbox"
              />
              <label for="is_active" class="form-check-label">
                Active (available for scheduling)
              </label>
            </div>
          </div>
        </div>

        <!-- Prerequisites -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Prerequisites & Notes</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Prerequisites</label>
              <textarea
                v-model="formData.prerequisites"
                rows="2"
                class="form-input resize-none"
                placeholder="List any requirements before this procedure..."
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Post-procedure Instructions</label>
              <textarea
                v-model="formData.post_instructions"
                rows="3"
                class="form-input resize-none"
                placeholder="Instructions for patient after procedure..."
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Internal Notes</label>
              <textarea
                v-model="formData.internal_notes"
                rows="2"
                class="form-input resize-none"
                placeholder="Internal notes for staff..."
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions flex flex-col sm:flex-row justify-end space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t">
          <button
            type="button"
            @click="$router.go(-1)"
            class="btn btn-secondary w-full sm:w-auto"
          >
            Cancel
          </button>
          
          <button
            type="submit"
            :disabled="isSubmitting"
            class="btn btn-primary w-full sm:w-auto"
            :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
          >
            <div v-if="isSubmitting" class="flex items-center justify-center">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              {{ isEditing ? 'Updating...' : 'Creating...' }}
            </div>
            <span v-else>
              {{ isEditing ? 'Update Procedure' : 'Create Procedure' }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isEditing = ref(false)
const isSubmitting = ref(false)
const errors = ref({})

const formData = ref({
  name: '',
  description: '',
  category: '',
  duration: 30,
  cost: 0,
  insurance_code: '',
  requires_anesthesia: false,
  requires_followup: false,
  is_active: true,
  prerequisites: '',
  post_instructions: '',
  internal_notes: ''
})

onMounted(() => {
  // Check if editing existing procedure
  if (route.params.id && route.params.id !== 'new') {
    isEditing.value = true
    loadProcedure(route.params.id)
  }
})

const loadProcedure = async (procedureId) => {
  // In a real app, this would fetch from API
  // For now, simulate loading
  try {
    // Simulate API call
    const mockProcedure = {
      id: procedureId,
      name: 'Dental Cleaning',
      description: 'Professional teeth cleaning and plaque removal',
      category: 'preventive',
      duration: 45,
      cost: 120,
      insurance_code: 'D1110',
      requires_anesthesia: false,
      requires_followup: false,
      is_active: true,
      prerequisites: 'No eating 2 hours before appointment',
      post_instructions: 'Avoid eating for 30 minutes after cleaning',
      internal_notes: 'Use ultrasonic scaler for heavy tartar'
    }
    
    Object.assign(formData.value, mockProcedure)
  } catch (error) {
    console.error('Error loading procedure:', error)
    // Handle error - could show toast notification
  }
}

const validateForm = () => {
  errors.value = {}
  
  if (!formData.value.name.trim()) {
    errors.value.name = 'Procedure name is required'
  }
  
  if (!formData.value.category) {
    errors.value.category = 'Category is required'
  }
  
  if (!formData.value.duration || formData.value.duration < 5) {
    errors.value.duration = 'Duration must be at least 5 minutes'
  }
  
  if (!formData.value.cost || formData.value.cost < 0) {
    errors.value.cost = 'Cost must be greater than 0'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }
  
  isSubmitting.value = true
  
  try {
    if (isEditing.value) {
      // Update existing procedure
      console.log('Updating procedure:', formData.value)
      // API call would go here
    } else {
      // Create new procedure
      console.log('Creating procedure:', formData.value)
      // API call would go here
    }
    
    // Navigate back to procedures list
    router.push('/procedures')
  } catch (error) {
    console.error('Error saving procedure:', error)
    // Handle error - could show toast notification
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
@import "../styles/main.css";

.form-label {
  @apply block text-sm font-medium text-gray-700 mb-1;
}

.form-input {
  @apply w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors;
}

.form-checkbox {
  @apply w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 focus:ring-2;
}

.form-check {
  @apply flex items-center;
}

.form-check-label {
  @apply ml-2 text-sm text-gray-700 cursor-pointer;
}

.error-message {
  @apply text-red-600 text-xs mt-1;
}

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.form-section {
  @apply border-b border-gray-200 pb-6 last:border-b-0 last:pb-0;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .form-container {
    @apply mx-4 rounded-lg;
  }
  
  .page-header {
    @apply mx-4;
  }
  
  .grid {
    @apply grid-cols-1;
  }
}
</style>