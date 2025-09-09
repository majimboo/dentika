<template>
  <div class="patient-form-page min-h-screen bg-gray-50 lg:bg-white">
    <!-- Mobile: Full screen form -->
    <div class="lg:max-w-4xl lg:mx-auto lg:py-8">
      <div class="bg-white lg:rounded-lg lg:shadow-sm lg:border min-h-screen lg:min-h-0">
        <!-- Form Header - Only show on desktop as mobile has NavigationHeader -->
        <div class="hidden lg:block border-b border-gray-200 px-6 py-4">
          <div class="flex items-center">
            <button 
              @click="$router.go(-1)" 
              class="mr-4 p-2 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <h1 class="text-xl font-semibold text-gray-900">
              {{ isEditing ? 'Edit Patient' : 'New Patient' }}
            </h1>
          </div>
        </div>

        <!-- Form Content -->
        <div class="p-4 lg:p-6">
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <!-- Basic Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="first_name" class="block text-sm font-medium text-gray-700 mb-1">
                    First Name *
                  </label>
                  <input
                    id="first_name"
                    v-model="form.first_name"
                    type="text"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.first_name }"
                  />
                  <span v-if="errors.first_name" class="text-red-500 text-sm mt-1">{{ errors.first_name }}</span>
                </div>

                <div>
                  <label for="last_name" class="block text-sm font-medium text-gray-700 mb-1">
                    Last Name *
                  </label>
                  <input
                    id="last_name"
                    v-model="form.last_name"
                    type="text"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.last_name }"
                  />
                  <span v-if="errors.last_name" class="text-red-500 text-sm mt-1">{{ errors.last_name }}</span>
                </div>

                <div>
                  <label for="date_of_birth" class="block text-sm font-medium text-gray-700 mb-1">
                    Date of Birth *
                  </label>
                  <input
                    id="date_of_birth"
                    v-model="form.date_of_birth"
                    type="date"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.date_of_birth }"
                  />
                  <span v-if="errors.date_of_birth" class="text-red-500 text-sm mt-1">{{ errors.date_of_birth }}</span>
                </div>

                <div>
                  <label for="gender" class="block text-sm font-medium text-gray-700 mb-1">
                    Gender *
                  </label>
                  <select
                    id="gender"
                    v-model="form.gender"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.gender }"
                  >
                    <option value="">Select Gender</option>
                    <option value="male">Male</option>
                    <option value="female">Female</option>
                    <option value="other">Other</option>
                  </select>
                  <span v-if="errors.gender" class="text-red-500 text-sm mt-1">{{ errors.gender }}</span>
                </div>
              </div>
            </div>

            <!-- Contact Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Contact Information</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="phone" class="block text-sm font-medium text-gray-700 mb-1">
                    Phone Number *
                  </label>
                  <input
                    id="phone"
                    v-model="form.phone"
                    type="tel"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.phone }"
                  />
                  <span v-if="errors.phone" class="text-red-500 text-sm mt-1">{{ errors.phone }}</span>
                </div>

                <div>
                  <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
                    Email
                  </label>
                  <input
                    id="email"
                    v-model="form.email"
                    type="email"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.email }"
                  />
                  <span v-if="errors.email" class="text-red-500 text-sm mt-1">{{ errors.email }}</span>
                </div>
              </div>

              <div class="mt-4">
                <label for="address" class="block text-sm font-medium text-gray-700 mb-1">
                  Address
                </label>
                <textarea
                  id="address"
                  v-model="form.address"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  :class="{ 'border-red-500': errors.address }"
                ></textarea>
                <span v-if="errors.address" class="text-red-500 text-sm mt-1">{{ errors.address }}</span>
              </div>
            </div>

            <!-- Emergency Contact -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Emergency Contact</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="emergency_contact_name" class="block text-sm font-medium text-gray-700 mb-1">
                    Contact Name
                  </label>
                  <input
                    id="emergency_contact_name"
                    v-model="form.emergency_contact_name"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="emergency_contact_phone" class="block text-sm font-medium text-gray-700 mb-1">
                    Contact Phone
                  </label>
                  <input
                    id="emergency_contact_phone"
                    v-model="form.emergency_contact_phone"
                    type="tel"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="emergency_contact_relationship" class="block text-sm font-medium text-gray-700 mb-1">
                    Relationship
                  </label>
                  <input
                    id="emergency_contact_relationship"
                    v-model="form.emergency_contact_relationship"
                    type="text"
                    placeholder="e.g., Spouse, Parent, Sibling"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Medical Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Medical Information</h2>
              
              <div class="space-y-4">
                <div>
                  <label for="medical_history" class="block text-sm font-medium text-gray-700 mb-1">
                    Medical History
                  </label>
                  <textarea
                    id="medical_history"
                    v-model="form.medical_history"
                    rows="4"
                    placeholder="Previous surgeries, chronic conditions, medications, etc."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  ></textarea>
                </div>

                <div>
                  <label for="allergies" class="block text-sm font-medium text-gray-700 mb-1">
                    Allergies
                  </label>
                  <textarea
                    id="allergies"
                    v-model="form.allergies"
                    rows="2"
                    placeholder="Food allergies, medication allergies, etc."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  ></textarea>
                </div>

                <div>
                  <label for="dental_history" class="block text-sm font-medium text-gray-700 mb-1">
                    Dental History
                  </label>
                  <textarea
                    id="dental_history"
                    v-model="form.dental_history"
                    rows="4"
                    placeholder="Previous dental treatments, oral health issues, etc."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  ></textarea>
                </div>
              </div>
            </div>

            <!-- Insurance Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Insurance Information</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="insurance_provider" class="block text-sm font-medium text-gray-700 mb-1">
                    Insurance Provider
                  </label>
                  <input
                    id="insurance_provider"
                    v-model="form.insurance_provider"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="insurance_number" class="block text-sm font-medium text-gray-700 mb-1">
                    Policy Number
                  </label>
                  <input
                    id="insurance_number"
                    v-model="form.insurance_number"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Additional Notes -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Additional Notes</h2>
              <textarea
                v-model="form.notes"
                rows="4"
                placeholder="Any additional information about the patient..."
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              ></textarea>
            </div>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row sm:justify-end space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t border-gray-200">
              <button
                type="button"
                @click="$router.go(-1)"
                class="w-full sm:w-auto px-4 py-2 text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="w-full sm:w-auto px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                <span v-if="loading" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Saving...
                </span>
                <span v-else>{{ isEditing ? 'Update Patient' : 'Create Patient' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useNotificationStore } from '../stores/notification'

const route = useRoute()
const router = useRouter()
const patientStore = usePatientStore()
const notificationStore = useNotificationStore()

// Component state
const loading = ref(false)
const errors = ref({})

// Form data
const form = reactive({
  first_name: '',
  last_name: '',
  date_of_birth: '',
  gender: '',
  phone: '',
  email: '',
  address: '',
  emergency_contact_name: '',
  emergency_contact_phone: '',
  emergency_contact_relationship: '',
  medical_history: '',
  allergies: '',
  dental_history: '',
  insurance_provider: '',
  insurance_number: '',
  notes: ''
})

// Computed properties
const isEditing = computed(() => !!route.params.id)

// Methods
const resetForm = () => {
  Object.keys(form).forEach(key => {
    form[key] = ''
  })
  errors.value = {}
}

const loadPatient = async (patientId) => {
  try {
    loading.value = true
    const patient = await patientStore.fetchPatient(patientId)
    
    if (patient) {
      // Populate form with patient data
      Object.keys(form).forEach(key => {
        if (patient[key] !== undefined) {
          form[key] = patient[key] || ''
        }
      })
    }
  } catch (error) {
    console.error('Error loading patient:', error)
    notificationStore.showError('Failed to load patient data')
    router.push('/patients')
  } finally {
    loading.value = false
  }
}

const validateForm = () => {
  errors.value = {}
  
  // Required fields
  if (!form.first_name.trim()) {
    errors.value.first_name = 'First name is required'
  }
  
  if (!form.last_name.trim()) {
    errors.value.last_name = 'Last name is required'
  }
  
  if (!form.date_of_birth) {
    errors.value.date_of_birth = 'Date of birth is required'
  }
  
  if (!form.gender) {
    errors.value.gender = 'Gender is required'
  }
  
  if (!form.phone.trim()) {
    errors.value.phone = 'Phone number is required'
  }
  
  // Email validation
  if (form.email && !/\S+@\S+\.\S+/.test(form.email)) {
    errors.value.email = 'Please enter a valid email address'
  }
  
  // Phone validation (basic)
  if (form.phone && !/^[\d\s\-\+\(\)]+$/.test(form.phone)) {
    errors.value.phone = 'Please enter a valid phone number'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) {
    notificationStore.showError('Please fix the form errors before submitting')
    return
  }
  
  try {
    loading.value = true
    
    if (isEditing.value) {
      await patientStore.updatePatient(route.params.id, form)
      notificationStore.showSuccess('Patient updated successfully')
    } else {
      await patientStore.createPatient(form)
      notificationStore.showSuccess('Patient created successfully')
    }
    
    router.push('/patients')
  } catch (error) {
    console.error('Error saving patient:', error)
    notificationStore.showError(
      isEditing.value ? 'Failed to update patient' : 'Failed to create patient'
    )
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  if (isEditing.value) {
    loadPatient(route.params.id)
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.form-section {
  @apply bg-gray-50 lg:bg-white p-4 lg:p-0 rounded-lg lg:rounded-none;
}

.form-section h2 {
  @apply border-b border-gray-200 pb-2 mb-4 lg:border-b-0 lg:pb-0;
}

/* Mobile optimizations */
@media (max-width: 1024px) {
  .patient-form-page {
    @apply p-0;
  }
  
  .form-section {
    @apply mx-0 mb-4;
  }
  
  .form-section:last-child {
    @apply mb-0;
  }
}
</style>