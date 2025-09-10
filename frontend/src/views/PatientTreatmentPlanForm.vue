<template>
  <div class="patient-treatment-plan-form-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">{{ isEditing ? 'Edit Treatment Plan' : 'Add Treatment Plan' }}</h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update treatment plan details' : 'Create a comprehensive treatment plan for the patient' }}
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
    <div v-else class="space-y-6">
      <!-- Basic Information -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <form @submit.prevent="handleSubmit" class="p-6">
          <div class="space-y-6">
            <!-- Title -->
            <div>
              <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
                Treatment Plan Title <span class="text-red-500">*</span>
              </label>
              <input
                id="title"
                v-model="form.title"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="{ 'border-red-500': errors.title }"
                placeholder="e.g., Root Canal Treatment for Tooth 14"
              >
              <span v-if="errors.title" class="text-red-500 text-sm mt-1">{{ errors.title }}</span>
            </div>

            <!-- Description -->
            <div>
              <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <textarea
                id="description"
                v-model="form.description"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Detailed description of the treatment plan..."
              ></textarea>
            </div>

            <!-- Priority and Status -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label for="priority" class="block text-sm font-medium text-gray-700 mb-2">
                  Priority
                </label>
                <select
                  id="priority"
                  v-model="form.priority"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="urgent">Urgent</option>
                </select>
              </div>

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
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                  <option value="on_hold">On Hold</option>
                </select>
              </div>
            </div>

            <!-- Related Diagnosis -->
            <div>
              <label for="diagnosis_id" class="block text-sm font-medium text-gray-700 mb-2">
                Related Diagnosis
              </label>
              <select
                id="diagnosis_id"
                v-model="form.diagnosis_id"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">No related diagnosis</option>
                <option
                  v-for="diagnosis in patientDiagnoses"
                  :key="diagnosis.id"
                  :value="diagnosis.id"
                >
                  {{ diagnosis.diagnosis_template?.name }} - {{ diagnosis.tooth_number || 'General' }}
                </option>
              </select>
            </div>

            <!-- Cost and Duration Estimates -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label for="estimated_cost" class="block text-sm font-medium text-gray-700 mb-2">
                  Estimated Cost (₱)
                </label>
                <input
                  id="estimated_cost"
                  v-model.number="form.estimated_cost"
                  type="number"
                  step="0.01"
                  min="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="0.00"
                >
              </div>

              <div>
                <label for="estimated_visits" class="block text-sm font-medium text-gray-700 mb-2">
                  Estimated Visits
                </label>
                <input
                  id="estimated_visits"
                  v-model.number="form.estimated_visits"
                  type="number"
                  min="1"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="1"
                >
              </div>

              <div>
                <label for="estimated_duration" class="block text-sm font-medium text-gray-700 mb-2">
                  Duration (Days)
                </label>
                <input
                  id="estimated_duration"
                  v-model.number="form.estimated_duration"
                  type="number"
                  min="1"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="30"
                >
              </div>
            </div>

            <!-- Date Range -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label for="start_date" class="block text-sm font-medium text-gray-700 mb-2">
                  Start Date
                </label>
                <input
                  id="start_date"
                  v-model="form.start_date"
                  type="date"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
              </div>

              <div>
                <label for="target_completion" class="block text-sm font-medium text-gray-700 mb-2">
                  Target Completion
                </label>
                <input
                  id="target_completion"
                  v-model="form.target_completion"
                  type="date"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
              </div>
            </div>

            <!-- Progress Tracking (Edit Mode Only) -->
            <div v-if="isEditing" class="border-t pt-6">
              <h3 class="text-lg font-medium text-gray-900 mb-4">Progress Tracking</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="completed_visits" class="block text-sm font-medium text-gray-700 mb-2">
                    Completed Visits
                  </label>
                  <input
                    id="completed_visits"
                    v-model.number="form.completed_visits"
                    type="number"
                    min="0"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="0"
                  >
                </div>

                <div>
                  <label for="actual_cost" class="block text-sm font-medium text-gray-700 mb-2">
                    Actual Cost (₱)
                  </label>
                  <input
                    id="actual_cost"
                    v-model.number="form.actual_cost"
                    type="number"
                    step="0.01"
                    min="0"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="0.00"
                  >
                </div>
              </div>
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
                {{ isEditing ? 'Update Treatment Plan' : 'Create Treatment Plan' }}
              </span>
            </button>
          </div>
        </form>
      </div>

      <!-- Procedures Section (Future Enhancement) -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900">Associated Procedures</h3>
          <button
            @click="showProcedureModal = true"
            class="btn btn-secondary text-sm"
          >
            <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
            Add Procedure
          </button>
        </div>

        <div v-if="treatmentPlanProcedures.length === 0" class="text-center py-8 text-gray-500">
          No procedures added yet. Click "Add Procedure" to get started.
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="procedure in treatmentPlanProcedures"
            :key="procedure.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
          >
            <div>
              <h4 class="font-medium text-gray-900">{{ procedure.procedure_template?.name }}</h4>
              <p class="text-sm text-gray-600">{{ procedure.notes }}</p>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm text-gray-500">Seq: {{ procedure.sequence }}</span>
              <button
                @click="removeProcedure(procedure)"
                class="text-red-600 hover:text-red-800"
              >
                <font-awesome-icon icon="fa-solid fa-trash" class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export default {
  name: 'PatientTreatmentPlanForm',
  setup() {
    const route = useRoute()
    const router = useRouter()

    // Reactive data
    const loading = ref(false)
    const submitting = ref(false)
    const error = ref('')
    const patientDiagnoses = ref([])
    const treatmentPlanProcedures = ref([])
    const showProcedureModal = ref(false)

    const form = reactive({
      title: '',
      description: '',
      priority: 'medium',
      status: 'active',
      diagnosis_id: null,
      estimated_cost: 0,
      estimated_visits: 1,
      estimated_duration: 30,
      start_date: '',
      target_completion: '',
      completed_visits: 0,
      actual_cost: 0
    })

    const errors = ref({})

    // Computed properties
    const patientId = computed(() => route.params.patientId)
    const treatmentPlanId = computed(() => route.params.treatmentPlanId)
    const isEditing = computed(() => !!treatmentPlanId.value)

    // Methods
    const loadPatientDiagnoses = async () => {
      try {
        const response = await fetch(`/api/patients/${patientId.value}/diagnoses`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          patientDiagnoses.value = data.diagnoses || []
        }
      } catch (error) {
        console.error('Error loading patient diagnoses:', error)
      }
    }

    const loadTreatmentPlan = async () => {
      if (!isEditing.value) return

      try {
        loading.value = true
        const response = await fetch(`/api/patients/${patientId.value}/treatment-plans/${treatmentPlanId.value}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const treatmentPlan = await response.json()
          // Populate form with treatment plan data
          form.title = treatmentPlan.title || ''
          form.description = treatmentPlan.description || ''
          form.priority = treatmentPlan.priority || 'medium'
          form.status = treatmentPlan.status || 'active'
          form.diagnosis_id = treatmentPlan.diagnosis_id
          form.estimated_cost = treatmentPlan.estimated_cost || 0
          form.estimated_visits = treatmentPlan.estimated_visits || 1
          form.estimated_duration = treatmentPlan.estimated_duration || 30
          form.start_date = treatmentPlan.start_date ? new Date(treatmentPlan.start_date).toISOString().split('T')[0] : ''
          form.target_completion = treatmentPlan.target_completion ? new Date(treatmentPlan.target_completion).toISOString().split('T')[0] : ''
          form.completed_visits = treatmentPlan.completed_visits || 0
          form.actual_cost = treatmentPlan.actual_cost || 0
        } else {
          error.value = 'Failed to load treatment plan'
        }
      } catch (err) {
        error.value = 'Failed to load treatment plan'
        console.error('Error loading treatment plan:', err)
      } finally {
        loading.value = false
      }
    }

    const validateForm = () => {
      errors.value = {}

      if (!form.title.trim()) {
        errors.value.title = 'Treatment plan title is required'
      }

      return Object.keys(errors.value).length === 0
    }

    const handleSubmit = async () => {
      if (!validateForm()) return

      try {
        submitting.value = true
        const url = isEditing.value
          ? `/api/patients/${patientId.value}/treatment-plans/${treatmentPlanId.value}`
          : `/api/patients/${patientId.value}/treatment-plans`

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
          error.value = 'Failed to save treatment plan'
        }
      } catch (err) {
        error.value = 'Failed to save treatment plan'
        console.error('Error saving treatment plan:', err)
      } finally {
        submitting.value = false
      }
    }

    const removeProcedure = (procedure) => {
      // TODO: Implement procedure removal
      console.log('Remove procedure:', procedure)
    }

    // Lifecycle
    onMounted(async () => {
      await Promise.all([
        loadPatientDiagnoses(),
        loadTreatmentPlan()
      ])
    })

    return {
      // Data
      loading,
      submitting,
      error,
      patientDiagnoses,
      treatmentPlanProcedures,
      showProcedureModal,
      form,
      errors,

      // Computed
      patientId,
      treatmentPlanId,
      isEditing,

      // Methods
      handleSubmit,
      removeProcedure
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