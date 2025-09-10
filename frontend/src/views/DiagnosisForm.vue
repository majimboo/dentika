<template>
  <div class="diagnosis-form-page">
    <!-- Page Header -->
    <div class="page-header flex items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">
          {{ isEditing ? 'Edit Diagnosis' : 'New Diagnosis' }}
        </h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update diagnosis details' : 'Add a new diagnosis to your catalog' }}
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
              <label class="form-label">Diagnosis Name *</label>
              <input
                v-model="formData.name"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.name }"
                placeholder="e.g. Dental Caries"
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
                placeholder="Describe the diagnosis in detail..."
              ></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">ICD-10 Code</label>
                <input
                  v-model="formData.icd10_code"
                  type="text"
                  class="form-input"
                  placeholder="e.g. K02.9"
                />
                <p class="text-xs text-gray-500 mt-1">International Classification of Diseases code</p>
              </div>

              <div class="form-group">
                <label class="form-label">Category *</label>
                <select
                  v-model="formData.category"
                  class="form-input"
                  :class="{ 'border-red-500': errors.category }"
                  required
                >
                  <option value="">Select category...</option>
                  <option value="caries">Caries</option>
                  <option value="periodontal">Periodontal</option>
                  <option value="endodontic">Endodontic</option>
                  <option value="orthodontic">Orthodontic</option>
                  <option value="oral_surgery">Oral Surgery</option>
                  <option value="other">Other</option>
                </select>
                <div v-if="errors.category" class="error-message">{{ errors.category }}</div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Severity *</label>
                <select
                  v-model="formData.severity"
                  class="form-input"
                  :class="{ 'border-red-500': errors.severity }"
                  required
                >
                  <option value="">Select severity...</option>
                  <option value="mild">Mild</option>
                  <option value="moderate">Moderate</option>
                  <option value="severe">Severe</option>
                </select>
                <div v-if="errors.severity" class="error-message">{{ errors.severity }}</div>
              </div>

              <div class="form-group">
                <label class="form-label">Affected Tooth/Teeth</label>
                <input
                  v-model="formData.affected_teeth"
                  type="text"
                  class="form-input"
                  placeholder="e.g. #14, #15 or Upper Left Molar"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Clinical Details -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Clinical Details</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Symptoms</label>
              <textarea
                v-model="formData.symptoms"
                rows="3"
                class="form-input resize-none"
                placeholder="List patient symptoms..."
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Clinical Findings</label>
              <textarea
                v-model="formData.clinical_findings"
                rows="3"
                class="form-input resize-none"
                placeholder="Describe clinical examination findings..."
              ></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Radiographic Findings</label>
                <textarea
                  v-model="formData.radiographic_findings"
                  rows="2"
                  class="form-input resize-none"
                  placeholder="X-ray findings..."
                ></textarea>
              </div>

              <div class="form-group">
                <label class="form-label">Differential Diagnosis</label>
                <textarea
                  v-model="formData.differential_diagnosis"
                  rows="2"
                  class="form-input resize-none"
                  placeholder="Alternative diagnoses considered..."
                ></textarea>
              </div>
            </div>
          </div>
        </div>

        <!-- Treatment & Prognosis -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Treatment & Prognosis</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Recommended Treatment</label>
              <textarea
                v-model="formData.recommended_treatment"
                rows="3"
                class="form-input resize-none"
                placeholder="Recommended treatment plan..."
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Prognosis</label>
              <select
                v-model="formData.prognosis"
                class="form-input"
              >
                <option value="">Select prognosis...</option>
                <option value="excellent">Excellent</option>
                <option value="good">Good</option>
                <option value="fair">Fair</option>
                <option value="poor">Poor</option>
                <option value="guarded">Guarded</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">Complications & Risks</label>
              <textarea
                v-model="formData.complications"
                rows="2"
                class="form-input resize-none"
                placeholder="Potential complications or risks..."
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Additional Settings -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Additional Settings</h3>
          
          <div class="space-y-4">
            <div class="form-check">
              <input
                v-model="formData.is_chronic"
                type="checkbox"
                id="is_chronic"
                class="form-checkbox"
              />
              <label for="is_chronic" class="form-check-label">
                Chronic condition
              </label>
            </div>

            <div class="form-check">
              <input
                v-model="formData.requires_monitoring"
                type="checkbox"
                id="requires_monitoring"
                class="form-checkbox"
              />
              <label for="requires_monitoring" class="form-check-label">
                Requires ongoing monitoring
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
                Active (available for use)
              </label>
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
            @click="goBack"
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
              {{ isEditing ? 'Update Diagnosis' : 'Create Diagnosis' }}
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
import { useNavigation } from '../composables/useNavigation'

const route = useRoute()
const router = useRouter()
const { goBack } = useNavigation()

const isEditing = ref(false)
const isSubmitting = ref(false)
const errors = ref({})

const formData = ref({
  name: '',
  description: '',
  icd10_code: '',
  category: '',
  severity: '',
  affected_teeth: '',
  symptoms: '',
  clinical_findings: '',
  radiographic_findings: '',
  differential_diagnosis: '',
  recommended_treatment: '',
  prognosis: '',
  complications: '',
  is_chronic: false,
  requires_monitoring: false,
  is_active: true,
  internal_notes: ''
})

onMounted(() => {
  // Check if editing existing diagnosis
  if (route.params.id && route.params.id !== 'new') {
    isEditing.value = true
    loadDiagnosis(route.params.id)
  }
})

const loadDiagnosis = async (diagnosisId) => {
  // In a real app, this would fetch from API
  // For now, simulate loading
  try {
    // Simulate API call
    const mockDiagnosis = {
      id: diagnosisId,
      name: 'Dental Caries',
      description: 'Tooth decay affecting enamel and dentin',
      icd10_code: 'K02.9',
      category: 'caries',
      severity: 'moderate',
      affected_teeth: '#14, #15',
      symptoms: 'Pain when chewing, sensitivity to hot/cold',
      clinical_findings: 'Cavitation visible on occlusal surface',
      radiographic_findings: 'Radiolucency extending to DEJ',
      differential_diagnosis: 'Arrested caries, enamel hypoplasia',
      recommended_treatment: 'Composite restoration, fluoride treatment',
      prognosis: 'good',
      complications: 'Possible pulpal involvement if left untreated',
      is_chronic: false,
      requires_monitoring: true,
      is_active: true,
      internal_notes: 'Patient has high caries risk'
    }
    
    Object.assign(formData.value, mockDiagnosis)
  } catch (error) {
    console.error('Error loading diagnosis:', error)
    // Handle error - could show toast notification
  }
}

const validateForm = () => {
  errors.value = {}
  
  if (!formData.value.name.trim()) {
    errors.value.name = 'Diagnosis name is required'
  }
  
  if (!formData.value.category) {
    errors.value.category = 'Category is required'
  }
  
  if (!formData.value.severity) {
    errors.value.severity = 'Severity is required'
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
      // Update existing diagnosis
      console.log('Updating diagnosis:', formData.value)
      // API call would go here
    } else {
      // Create new diagnosis
      console.log('Creating diagnosis:', formData.value)
      // API call would go here
    }
    
    // Navigate back to procedures list (diagnoses tab)
    router.push('/procedures?tab=diagnoses')
  } catch (error) {
    console.error('Error saving diagnosis:', error)
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