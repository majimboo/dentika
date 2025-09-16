<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
      <!-- Background overlay -->
      <div
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        @click="$emit('cancel')"
      ></div>

      <!-- Modal panel -->
      <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
        <div class="modal-header flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold leading-6 text-gray-900">
              Create Dental Chart Snapshot
            </h3>
            <p class="text-sm text-gray-600">Save current chart state for this visit</p>
          </div>
          <button
            @click="$emit('cancel')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleCreate" class="space-y-4">
          <!-- Chart type selection -->
          <div class="form-group">
            <label for="chart-type" class="block text-sm font-medium text-gray-700 mb-1">
              Chart Type *
            </label>
            <select
              id="chart-type"
              v-model="formData.chartType"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select chart type...</option>
              <option value="permanent">Adult Teeth (Permanent)</option>
              <option value="primary">Child Teeth (Primary)</option>
            </select>
          </div>

          <!-- Link to appointment -->
          <div class="form-group">
            <label for="appointment" class="block text-sm font-medium text-gray-700 mb-1">
              Link to Appointment/Visit
            </label>
            <select
              id="appointment"
              v-model="formData.appointmentId"
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">No specific appointment</option>
              <option
                v-for="appointment in availableAppointments"
                :key="appointment.id"
                :value="appointment.id"
              >
                {{ formatAppointmentOption(appointment) }}
              </option>
            </select>
            <p class="text-xs text-gray-500 mt-1">
              Link this snapshot to a specific appointment or visit for better tracking.
            </p>
          </div>

          <!-- Visit notes -->
          <div class="form-group">
            <label for="visit-notes" class="block text-sm font-medium text-gray-700 mb-1">
              Visit Notes
            </label>
            <textarea
              id="visit-notes"
              v-model="formData.visitNotes"
              rows="3"
              placeholder="Add any relevant notes about this visit or chart updates..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
            ></textarea>
          </div>

          <!-- Current chart preview -->
          <div v-if="formData.chartType" class="chart-preview p-3 bg-gray-50 border border-gray-200 rounded-lg">
            <h4 class="text-sm font-medium text-gray-800 mb-2">
              Current {{ formatChartType(formData.chartType) }} Chart Summary
            </h4>
            <div v-if="currentRecord" class="chart-summary text-sm text-gray-600">
              <div class="summary-stats flex space-x-4">
                <div>Total teeth: {{ getTotalTeethCount() }}</div>
                <div>Healthy: {{ getHealthyTeethCount() }}</div>
                <div>Issues: {{ getProblemsCount() }}</div>
              </div>
              <div class="mt-2 text-xs text-gray-500">
                This snapshot will preserve the current state of all teeth in this chart.
              </div>
            </div>
            <div v-else class="text-sm text-red-600">
              No active {{ formatChartType(formData.chartType) }} chart found for this patient.
            </div>
          </div>

          <!-- Action buttons -->
          <div class="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              @click="$emit('cancel')"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="!isFormValid || creating"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="creating">Creating...</span>
              <span v-else>Create Snapshot</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useAppointmentStore } from '../stores/appointment'
import { useDentalRecordStore } from '../stores/dentalRecord'
import { useNotificationStore } from '../stores/notification'

const props = defineProps({
  patientId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['created', 'cancel'])

const appointmentStore = useAppointmentStore()
const dentalRecordStore = useDentalRecordStore()
const notificationStore = useNotificationStore()

const formData = reactive({
  chartType: '',
  appointmentId: null,
  visitNotes: ''
})

const availableAppointments = ref([])
const creating = ref(false)

const isFormValid = computed(() => {
  return formData.chartType && currentRecord.value
})

const currentRecord = computed(() => {
  if (!formData.chartType) return null
  return dentalRecordStore.dentalRecords.find(record =>
    record.record_type === formData.chartType && record.is_active
  )
})

const formatChartType = (chartType) => {
  return chartType === 'permanent' ? 'Adult' : 'Child'
}

const formatAppointmentOption = (appointment) => {
  const date = new Date(appointment.start_time).toLocaleDateString()
  const time = new Date(appointment.start_time).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit'
  })
  const status = appointment.status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
  return `${date} ${time} - ${appointment.title || 'Appointment'} (${status})`
}

const getTotalTeethCount = () => {
  if (!currentRecord.value?.teeth_data) return 0
  return currentRecord.value.teeth_data.length
}

const getHealthyTeethCount = () => {
  if (!currentRecord.value?.teeth_data) return 0
  return currentRecord.value.teeth_data.filter(tooth => tooth.condition === 'healthy').length
}

const getProblemsCount = () => {
  if (!currentRecord.value?.teeth_data) return 0
  return currentRecord.value.teeth_data.filter(tooth =>
    ['decay', 'extracted', 'missing'].includes(tooth.condition)
  ).length
}

const loadAppointments = async () => {
  try {
    await appointmentStore.fetchPatientAppointments(props.patientId)

    // Filter to get relevant appointments (recent and completed/in progress)
    const now = new Date()
    const thirtyDaysAgo = new Date(now.getTime() - (30 * 24 * 60 * 60 * 1000))
    const sevenDaysAhead = new Date(now.getTime() + (7 * 24 * 60 * 60 * 1000))

    availableAppointments.value = appointmentStore.appointments.filter(appointment => {
      const appointmentDate = new Date(appointment.start_time)
      const isInRange = appointmentDate >= thirtyDaysAgo && appointmentDate <= sevenDaysAhead
      const isRelevant = ['completed', 'in_progress', 'scheduled'].includes(appointment.status)

      return isInRange && isRelevant
    }).sort((a, b) => new Date(b.start_time) - new Date(a.start_time))
  } catch (error) {
    console.error('Error loading appointments:', error)
  }
}

const handleCreate = async () => {
  if (!isFormValid.value) return

  creating.value = true

  try {
    const snapshotData = {
      chart_type: formData.chartType,
      appointment_id: formData.appointmentId || null,
      visit_notes: formData.visitNotes
    }

    const result = await dentalRecordStore.createDentalChartSnapshot(props.patientId, snapshotData)

    if (result.success) {
      emit('created', result.data)
    } else {
      notificationStore.showError(result.error || 'Failed to create snapshot')
    }
  } catch (error) {
    console.error('Error creating snapshot:', error)
    notificationStore.showError('Failed to create snapshot')
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  loadAppointments()
  // Auto-select the most common chart type if only one is active
  const activeRecords = dentalRecordStore.dentalRecords.filter(r => r.is_active)
  if (activeRecords.length === 1) {
    formData.chartType = activeRecords[0].record_type
  }
})
</script>

<style scoped>
@import "../styles/main.css";

.chart-preview {
  @apply transition-all duration-200;
}

.summary-stats {
  @apply text-sm;
}

@media (max-width: 640px) {
  .summary-stats {
    @apply flex-col space-x-0 space-y-2;
  }
}
</style>