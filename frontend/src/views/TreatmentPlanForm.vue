<template>
  <div class="treatment-plan-form-page">
    <!-- Page Header -->
    <div class="page-header flex items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">
          {{ isEditing ? 'Edit Treatment Plan' : 'New Treatment Plan' }}
        </h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update treatment plan details' : 'Create a comprehensive treatment plan for a patient' }}
        </p>
      </div>
    </div>

    <!-- Form -->
    <div class="form-container bg-white rounded-lg shadow-sm border border-gray-200">
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <!-- Patient Selection -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Patient Information</h3>
          
          <div class="form-group">
            <label class="form-label">Patient *</label>
            <div class="relative">
              <input
                v-model="patientSearchQuery"
                @input="searchPatients"
                @focus="showPatientDropdown = true"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.patient_id }"
                placeholder="Search for patient..."
                autocomplete="off"
                required
              />
              
              <!-- Patient Dropdown -->
              <div 
                v-if="showPatientDropdown && (filteredPatients.length > 0 || patientSearchQuery)"
                class="absolute z-10 w-full bg-white border border-gray-300 rounded-md shadow-lg max-h-48 overflow-y-auto mt-1"
              >
                <div 
                  v-for="patient in filteredPatients" 
                  :key="patient.id"
                  @click="selectPatient(patient)"
                  class="patient-option px-3 py-2 hover:bg-gray-100 cursor-pointer border-b last:border-b-0"
                >
                  <div class="font-medium text-gray-900">
                    {{ patient.first_name }} {{ patient.last_name }}
                  </div>
                  <div class="text-sm text-gray-600">
                    {{ patient.phone }} • {{ patient.email }}
                  </div>
                </div>
                
                <div 
                  v-if="filteredPatients.length === 0 && patientSearchQuery"
                  class="px-3 py-2 text-gray-500 text-sm"
                >
                  No patients found
                </div>
              </div>
            </div>
            <div v-if="errors.patient_id" class="error-message">{{ errors.patient_id }}</div>
          </div>
        </div>

        <!-- Basic Information -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Treatment Plan Details</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Treatment Plan Name *</label>
              <input
                v-model="formData.name"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.name }"
                placeholder="e.g. Full Mouth Restoration"
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
                placeholder="Describe the treatment plan objectives..."
              ></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Status</label>
                <select
                  v-model="formData.status"
                  class="form-input"
                >
                  <option value="draft">Draft</option>
                  <option value="active">Active</option>
                  <option value="on_hold">On Hold</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>

              <div class="form-group">
                <label class="form-label">Priority</label>
                <select
                  v-model="formData.priority"
                  class="form-input"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="urgent">Urgent</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="form-group">
                <label class="form-label">Start Date</label>
                <input
                  v-model="formData.start_date"
                  type="date"
                  class="form-input"
                />
              </div>

              <div class="form-group">
                <label class="form-label">Target Completion Date</label>
                <input
                  v-model="formData.target_completion_date"
                  type="date"
                  class="form-input"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Diagnoses -->
        <div class="form-section">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Primary Diagnoses</h3>
             <button
               type="button"
               @click="addDiagnosis"
               class="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center"
             >
               <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-1" />
               Add Diagnosis
             </button>
          </div>
          
          <div class="space-y-3">
            <div 
              v-for="(diagnosis, index) in formData.diagnoses"
              :key="index"
              class="diagnosis-item flex items-center space-x-3 p-3 bg-gray-50 rounded-lg"
            >
              <select
                v-model="diagnosis.id"
                class="flex-1 form-input"
                required
              >
                <option value="">Select diagnosis...</option>
                <option v-for="diag in availableDiagnoses" :key="diag.id" :value="diag.id">
                  {{ diag.name }} ({{ diag.icd10_code }})
                </option>
              </select>
              
              <textarea
                v-model="diagnosis.notes"
                rows="1"
                class="flex-1 form-input resize-none"
                placeholder="Notes for this diagnosis..."
              ></textarea>
              
               <button
                 type="button"
                 @click="removeDiagnosis(index)"
                 class="text-red-600 hover:text-red-700 p-1"
               >
                 <font-awesome-icon icon="fa-solid fa-times" class="w-4 h-4" />
               </button>
            </div>
            
            <div v-if="formData.diagnoses.length === 0" class="text-center py-6 text-gray-500">
              <p>No diagnoses added yet</p>
              <button
                type="button"
                @click="addDiagnosis"
                class="mt-2 text-blue-600 hover:text-blue-700 text-sm font-medium"
              >
                Add your first diagnosis
              </button>
            </div>
          </div>
        </div>

        <!-- Procedures -->
        <div class="form-section">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Treatment Procedures</h3>
             <button
               type="button"
               @click="addProcedure"
               class="text-blue-600 hover:text-blue-700 text-sm font-medium flex items-center"
             >
               <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-1" />
               Add Procedure
             </button>
          </div>
          
          <div class="space-y-3">
            <div 
              v-for="(procedure, index) in formData.procedures"
              :key="index"
              class="procedure-item p-4 bg-gray-50 rounded-lg space-y-3"
            >
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <select
                  v-model="procedure.id"
                  @change="updateProcedureCost(index)"
                  class="form-input"
                  required
                >
                  <option value="">Select procedure...</option>
                  <option v-for="proc in availableProcedures" :key="proc.id" :value="proc.id">
                    {{ proc.name }} - ${{ proc.cost }}
                  </option>
                </select>
                
                <input
                  v-model="procedure.tooth_number"
                  type="text"
                  class="form-input"
                  placeholder="Tooth/Teeth (e.g. #14, #15)"
                />
              </div>
              
              <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                <input
                  v-model.number="procedure.quantity"
                  type="number"
                  min="1"
                  class="form-input"
                  placeholder="Quantity"
                />
                
                <input
                  v-model.number="procedure.cost"
                  type="number"
                  step="0.01"
                  class="form-input"
                  placeholder="Cost"
                />
                
                <select
                  v-model="procedure.status"
                  class="form-input"
                >
                  <option value="planned">Planned</option>
                  <option value="in_progress">In Progress</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
              
              <div class="flex items-start space-x-3">
                <textarea
                  v-model="procedure.notes"
                  rows="2"
                  class="flex-1 form-input resize-none"
                  placeholder="Notes for this procedure..."
                ></textarea>
                
                 <button
                   type="button"
                   @click="removeProcedure(index)"
                   class="text-red-600 hover:text-red-700 p-1 mt-1"
                 >
                   <font-awesome-icon icon="fa-solid fa-times" class="w-4 h-4" />
                 </button>
              </div>
            </div>
            
            <div v-if="formData.procedures.length === 0" class="text-center py-6 text-gray-500">
              <p>No procedures added yet</p>
              <button
                type="button"
                @click="addProcedure"
                class="mt-2 text-blue-600 hover:text-blue-700 text-sm font-medium"
              >
                Add your first procedure
              </button>
            </div>
          </div>
        </div>

        <!-- Financial Information -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Financial Information</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="form-group">
              <label class="form-label">Total Cost</label>
              <div class="relative">
                <input
                  :value="totalCost"
                  type="text"
                  class="form-input pl-8 bg-gray-50"
                  readonly
                />
              </div>
              <p class="text-xs text-gray-500 mt-1">Automatically calculated</p>
            </div>

            <div class="form-group">
              <label class="form-label">Insurance Coverage (%)</label>
              <input
                v-model.number="formData.insurance_coverage"
                type="number"
                min="0"
                max="100"
                class="form-input"
                placeholder="0"
              />
            </div>

            <div class="form-group">
              <label class="form-label">Patient Portion</label>
              <div class="relative">
                <input
                  :value="patientPortion"
                  type="text"
                  class="form-input pl-8 bg-gray-50"
                  readonly
                />
              </div>
              <p class="text-xs text-gray-500 mt-1">After insurance</p>
            </div>
          </div>
        </div>

        <!-- Notes -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Additional Information</h3>
          
          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Treatment Goals</label>
              <textarea
                v-model="formData.treatment_goals"
                rows="3"
                class="form-input resize-none"
                placeholder="What are the goals of this treatment plan?"
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">Special Instructions</label>
              <textarea
                v-model="formData.special_instructions"
                rows="2"
                class="form-input resize-none"
                placeholder="Any special instructions for the treatment team..."
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
              {{ isEditing ? 'Update Treatment Plan' : 'Create Treatment Plan' }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNavigation } from '../composables/useNavigation'

const route = useRoute()
const router = useRouter()
const { goBack } = useNavigation()

const isEditing = ref(false)
const isSubmitting = ref(false)
const errors = ref({})

const patientSearchQuery = ref('')
const showPatientDropdown = ref(false)
const filteredPatients = ref([])

const formData = ref({
  patient_id: '',
  name: '',
  description: '',
  status: 'draft',
  priority: 'medium',
  start_date: '',
  target_completion_date: '',
  diagnoses: [],
  procedures: [],
  insurance_coverage: 0,
  treatment_goals: '',
  special_instructions: '',
  internal_notes: ''
})

// Mock data
const availableDiagnoses = ref([
  { id: 1, name: 'Dental Caries', icd10_code: 'K02.9' },
  { id: 2, name: 'Gingivitis', icd10_code: 'K05.0' },
  { id: 3, name: 'Periodontitis', icd10_code: 'K05.3' }
])

const availableProcedures = ref([
  { id: 1, name: 'Dental Cleaning', cost: 120 },
  { id: 2, name: 'Root Canal Treatment', cost: 800 },
  { id: 3, name: 'Crown Placement', cost: 1200 }
])

const mockPatients = ref([
  { id: 1, first_name: 'John', last_name: 'Smith', phone: '(555) 123-4567', email: 'john@example.com' },
  { id: 2, first_name: 'Jane', last_name: 'Doe', phone: '(555) 987-6543', email: 'jane@example.com' }
])

// Computed properties
const totalCost = computed(() => {
  return formData.value.procedures.reduce((total, proc) => {
    return total + (proc.cost * (proc.quantity || 1))
  }, 0).toFixed(2)
})

const patientPortion = computed(() => {
  const total = parseFloat(totalCost.value)
  const coverage = formData.value.insurance_coverage || 0
  return (total * (100 - coverage) / 100).toFixed(2)
})

// Methods
const searchPatients = () => {
  if (patientSearchQuery.value.length >= 2) {
    const query = patientSearchQuery.value.toLowerCase()
    filteredPatients.value = mockPatients.value.filter(patient => 
      patient.first_name.toLowerCase().includes(query) ||
      patient.last_name.toLowerCase().includes(query) ||
      patient.email.toLowerCase().includes(query)
    )
  } else {
    filteredPatients.value = []
  }
}

const selectPatient = (patient) => {
  formData.value.patient_id = patient.id
  patientSearchQuery.value = `${patient.first_name} ${patient.last_name}`
  showPatientDropdown.value = false
  filteredPatients.value = []
}

const addDiagnosis = () => {
  formData.value.diagnoses.push({
    id: '',
    notes: ''
  })
}

const removeDiagnosis = (index) => {
  formData.value.diagnoses.splice(index, 1)
}

const addProcedure = () => {
  formData.value.procedures.push({
    id: '',
    tooth_number: '',
    quantity: 1,
    cost: 0,
    status: 'planned',
    notes: ''
  })
}

const removeProcedure = (index) => {
  formData.value.procedures.splice(index, 1)
}

const updateProcedureCost = (index) => {
  const procedure = formData.value.procedures[index]
  const selectedProc = availableProcedures.value.find(p => p.id === procedure.id)
  if (selectedProc) {
    procedure.cost = selectedProc.cost
  }
}

const validateForm = () => {
  errors.value = {}
  
  if (!formData.value.patient_id) {
    errors.value.patient_id = 'Please select a patient'
  }
  
  if (!formData.value.name.trim()) {
    errors.value.name = 'Treatment plan name is required'
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
      // Update existing treatment plan
      console.log('Updating treatment plan:', formData.value)
      // API call would go here
    } else {
      // Create new treatment plan
      console.log('Creating treatment plan:', formData.value)
      // API call would go here
    }
    
    // Navigate back to procedures list (treatments tab)
    router.push('/procedures?tab=treatments')
  } catch (error) {
    console.error('Error saving treatment plan:', error)
    // Handle error - could show toast notification
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  // Check if editing existing treatment plan
  if (route.params.id && route.params.id !== 'new') {
    isEditing.value = true
    // Load treatment plan data
  }
})

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showPatientDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
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
  
  .procedure-item .grid {
    @apply grid-cols-1;
  }
}
</style>