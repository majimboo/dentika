<template>
  <div class="patient-diagnosis-form-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">{{ isEditing ? 'Edit Diagnosis' : 'Add Diagnosis' }}</h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update diagnosis details' : 'Record a new diagnosis for the patient' }}
        </p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <router-link
          :to="`/patients/${patientId}`"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-arrow-left" class="w-4 h-4 mr-2" />
          Back to Patient
        </router-link>
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
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Form -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <form @submit.prevent="handleSubmit" class="p-6">
        <div class="space-y-6">
          <!-- Diagnosis Template Selection -->
          <div>
            <label for="diagnosis_template_id" class="block text-sm font-medium text-gray-700 mb-2">
              Diagnosis Template <span class="text-red-500">*</span>
            </label>
            <select
              id="diagnosis_template_id"
              v-model="form.diagnosis_template_id"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.diagnosis_template_id }"
            >
              <option value="">Select a diagnosis template</option>
              <option
                v-for="template in diagnosisTemplates"
                :key="template.id"
                :value="template.id"
              >
                {{ template.code }} - {{ template.name }}
              </option>
            </select>
            <span v-if="errors.diagnosis_template_id" class="text-red-500 text-sm mt-1">{{ errors.diagnosis_template_id }}</span>
          </div>

          <!-- Selected Template Details -->
          <div v-if="selectedTemplate" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h4 class="font-medium text-blue-900 mb-2">{{ selectedTemplate.name }}</h4>
            <p class="text-sm text-blue-700 mb-2">{{ selectedTemplate.description }}</p>
            <div class="flex items-center space-x-4 text-xs text-blue-600">
              <span>Category: {{ selectedTemplate.category }}</span>
              <span>Severity: {{ selectedTemplate.severity }}</span>
            </div>
          </div>

          <!-- Tooth Number -->
          <div>
            <label for="tooth_number" class="block text-sm font-medium text-gray-700 mb-2">
              Tooth Number
            </label>
            <input
              id="tooth_number"
              v-model="form.tooth_number"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., 11, 12, A, B"
            >
          </div>

          <!-- Surface -->
          <div>
            <label for="surface" class="block text-sm font-medium text-gray-700 mb-2">
              Surface
            </label>
            <select
              id="surface"
              v-model="form.surface"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select surface</option>
              <option value="mesial">Mesial</option>
              <option value="distal">Distal</option>
              <option value="occlusal">Occlusal</option>
              <option value="buccal">Buccal</option>
              <option value="lingual">Lingual</option>
              <option value="incisal">Incisal</option>
            </select>
          </div>

          <!-- Severity Override -->
          <div>
            <label for="severity" class="block text-sm font-medium text-gray-700 mb-2">
              Severity
            </label>
            <select
              id="severity"
              v-model="form.severity"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Use template default</option>
              <option value="mild">Mild</option>
              <option value="moderate">Moderate</option>
              <option value="severe">Severe</option>
            </select>
          </div>

          <!-- Notes -->
          <div>
            <label for="notes" class="block text-sm font-medium text-gray-700 mb-2">
              Additional Notes
            </label>
            <textarea
              id="notes"
              v-model="form.notes"
              rows="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Additional notes about this diagnosis..."
            ></textarea>
          </div>

          <!-- Associated Appointment -->
          <div>
            <label for="appointment_id" class="block text-sm font-medium text-gray-700 mb-2">
              Associated Appointment
            </label>
            <select
              id="appointment_id"
              v-model="form.appointment_id"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">No associated appointment</option>
              <option
                v-for="appointment in patientAppointments"
                :key="appointment.id"
                :value="appointment.id"
              >
                {{ formatDate(appointment.appointment_date) }} - {{ appointment.appointment_type }}
              </option>
            </select>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200 mt-6">
          <router-link
            :to="`/patients/${patientId}`"
            class="btn btn-secondary"
          >
            Cancel
          </router-link>
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
              {{ isEditing ? 'Update Diagnosis' : 'Save Diagnosis' }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export default {
  name: 'PatientDiagnosisForm',
  setup() {
    const route = useRoute()
    const router = useRouter()

    // Reactive data
    const loading = ref(false)
    const submitting = ref(false)
    const error = ref('')
    const diagnosisTemplates = ref([])
    const patientAppointments = ref([])

    const form = reactive({
      diagnosis_template_id: '',
      tooth_number: '',
      surface: '',
      severity: '',
      notes: '',
      appointment_id: null
    })

    const errors = ref({})

    // Computed properties
    const patientId = computed(() => route.params.patientId)
    const diagnosisId = computed(() => route.params.diagnosisId)
    const isEditing = computed(() => !!diagnosisId.value)
    const selectedTemplate = computed(() => {
      return diagnosisTemplates.value.find(template => template.id == form.diagnosis_template_id)
    })

    // Methods
    const loadDiagnosisTemplates = async () => {
      try {
        const response = await fetch('/api/diagnosis-templates', {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          diagnosisTemplates.value = data.templates || []
        }
      } catch (error) {
        console.error('Error loading diagnosis templates:', error)
      }
    }

    const loadPatientAppointments = async () => {
      try {
        const response = await fetch(`/api/appointments?patient_id=${patientId.value}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          patientAppointments.value = data.appointments || []
        }
      } catch (error) {
        console.error('Error loading patient appointments:', error)
      }
    }

    const loadDiagnosis = async () => {
      if (!isEditing.value) return

      try {
        loading.value = true
        const response = await fetch(`/api/patients/${patientId.value}/diagnoses/${diagnosisId.value}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const diagnosis = await response.json()
          // Populate form with diagnosis data
          form.diagnosis_template_id = diagnosis.diagnosis_template_id
          form.tooth_number = diagnosis.tooth_number || ''
          form.surface = diagnosis.surface || ''
          form.severity = diagnosis.severity || ''
          form.notes = diagnosis.notes || ''
          form.appointment_id = diagnosis.appointment_id
        } else {
          error.value = 'Failed to load diagnosis'
        }
      } catch (err) {
        error.value = 'Failed to load diagnosis'
        console.error('Error loading diagnosis:', err)
      } finally {
        loading.value = false
      }
    }

    const validateForm = () => {
      errors.value = {}

      if (!form.diagnosis_template_id) {
        errors.value.diagnosis_template_id = 'Diagnosis template is required'
      }

      return Object.keys(errors.value).length === 0
    }

    const handleSubmit = async () => {
      if (!validateForm()) return

      try {
        submitting.value = true
        const url = isEditing.value
          ? `/api/patients/${patientId.value}/diagnoses/${diagnosisId.value}`
          : `/api/patients/${patientId.value}/diagnoses`

        const method = isEditing.value ? 'PUT' : 'POST'

        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: JSON.stringify(form)
        })

        if (response.ok) {
          router.push(`/patients/${patientId.value}`)
        } else {
          error.value = 'Failed to save diagnosis'
        }
      } catch (err) {
        error.value = 'Failed to save diagnosis'
        console.error('Error saving diagnosis:', err)
      } finally {
        submitting.value = false
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString()
    }

    // Lifecycle
    onMounted(async () => {
      await Promise.all([
        loadDiagnosisTemplates(),
        loadPatientAppointments(),
        loadDiagnosis()
      ])
    })

    return {
      // Data
      loading,
      submitting,
      error,
      diagnosisTemplates,
      patientAppointments,
      form,
      errors,

      // Computed
      patientId,
      diagnosisId,
      isEditing,
      selectedTemplate,

      // Methods
      handleSubmit,
      formatDate
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