<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
      <!-- Background overlay -->
      <div 
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        @click="$emit('cancel')"
      ></div>

      <!-- Modal panel -->
      <div class="relative transform overflow-hidden rounded-lg bg-white w-full max-w-lg mx-auto shadow-xl transition-all my-4 sm:my-8">
        <!-- Mobile-optimized content wrapper -->
        <div class="px-4 py-5 sm:p-6 max-h-[85vh] overflow-y-auto">
        <div class="modal-header flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold leading-6 text-gray-900">
              Edit Tooth {{ tooth.tooth_number }}
            </h3>
            <p class="text-sm text-gray-600">{{ getToothName() }}</p>
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

        <form @submit.prevent="handleSave" class="space-y-4">
          <!-- Current condition display -->
          <div class="current-condition p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center">
              <div 
                class="w-4 h-4 rounded mr-2 border border-gray-300"
                :style="{ backgroundColor: getCurrentColor() }"
              ></div>
              <span class="text-sm font-medium">
                Current: {{ formatCondition(tooth.condition) }}
              </span>
            </div>
          </div>

          <!-- New condition selection -->
          <div class="form-group">
            <label for="condition" class="block text-sm font-medium text-gray-700 mb-1">
              New Condition *
            </label>
            <select 
              id="condition"
              v-model="editData.condition"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select condition...</option>
              <option v-for="(label, value) in conditionOptions" :key="value" :value="value">
                {{ label }}
              </option>
            </select>
          </div>

          <!-- Surface-specific conditions -->
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Surface-Specific Conditions
            </label>
            <div class="surface-conditions-container">
              <div
                v-for="surface in availableSurfaces"
                :key="surface.value"
                class="surface-condition-item p-3 border border-gray-200 rounded-lg mb-2"
              >
                <div class="flex items-center justify-between mb-2">
                  <label class="text-sm font-medium text-gray-700">
                    {{ surface.label }}
                  </label>
                  <button
                    type="button"
                    @click="toggleSurfaceCondition(surface.value)"
                    class="text-xs px-2 py-1 rounded transition-colors"
                    :class="hasSurfaceCondition(surface.value)
                      ? 'bg-blue-100 text-blue-700 hover:bg-blue-200'
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
                  >
                    {{ hasSurfaceCondition(surface.value) ? 'Remove' : 'Add Condition' }}
                  </button>
                </div>

                <div v-if="hasSurfaceCondition(surface.value)" class="surface-condition-details">
                  <select
                    :value="getSurfaceCondition(surface.value)"
                    @change="updateSurfaceCondition(surface.value, $event.target.value)"
                    class="w-full border border-gray-300 rounded-md px-2 py-1 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">Select condition...</option>
                    <option v-for="(label, value) in conditionOptions" :key="value" :value="value">
                      {{ label }}
                    </option>
                  </select>
                </div>

                <div v-else class="text-xs text-gray-500">
                  This surface is healthy
                </div>
              </div>

              <!-- Overall tooth condition note -->
              <div class="overall-condition-note p-2 bg-yellow-50 border border-yellow-200 rounded text-sm">
                <div class="flex items-start">
                  <font-awesome-icon icon="fa-solid fa-info-circle" class="w-3 h-3 text-yellow-600 mt-0.5 mr-2" />
                  <div class="text-yellow-800">
                    <strong>Overall tooth condition:</strong> {{ formatCondition(editData.condition) }}
                    <div class="text-xs mt-1">
                      The overall condition represents the tooth's general state, while surface conditions show specific issues on each area.
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Notes -->
          <div class="form-group">
            <label for="notes" class="block text-sm font-medium text-gray-700 mb-1">
              Notes
            </label>
            <textarea
              id="notes"
              v-model="editData.notes"
              rows="3"
              placeholder="Add any relevant notes about this tooth condition..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
            ></textarea>
          </div>

          <!-- Link to appointment/visit -->
          <div class="form-group">
            <label for="appointment" class="block text-sm font-medium text-gray-700 mb-1">
              Link to Appointment/Visit
            </label>
            <select
              id="appointment"
              v-model="editData.appointmentId"
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
              Linking to an appointment helps track when this change was made during a specific visit.
            </p>
          </div>

          <!-- Reason for change -->
          <div class="form-group">
            <label for="reason" class="block text-sm font-medium text-gray-700 mb-1">
              Reason for Change *
            </label>
            <input
              id="reason"
              type="text"
              v-model="editData.reason"
              required
              placeholder="e.g., Treatment completed, new findings, regular checkup..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
            >
          </div>

          <!-- Preview of changes -->
          <div v-if="hasChanges" class="changes-preview p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
            <h4 class="text-sm font-medium text-yellow-800 mb-2">Preview Changes:</h4>
            <div class="text-sm text-yellow-700">
              <div v-if="editData.condition !== tooth.condition" class="flex items-center mb-1">
                <span class="font-medium">Condition:</span>
                <span class="ml-2">{{ formatCondition(tooth.condition) }}</span>
                <svg class="w-4 h-4 mx-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                </svg>
                <span>{{ formatCondition(editData.condition) }}</span>
              </div>
              
              <div v-if="surfacesChanged" class="flex items-start mb-1">
                <span class="font-medium">Surface Conditions:</span>
                <div class="ml-2">
                  <div>{{ formatSurfaceConditions(tooth.surface_conditions || []) }}</div>
                  <div class="text-xs">↓</div>
                  <div>{{ formatSurfaceConditions(editData.surfaceConditions) }}</div>
                </div>
              </div>
              
              <div v-if="editData.notes !== (tooth.notes || '')" class="mb-1">
                <span class="font-medium">Notes:</span>
                <span class="ml-2">{{ editData.notes ? 'Updated' : 'Cleared' }}</span>
              </div>
            </div>
          </div>

          <!-- Action buttons -->
          <div class="flex flex-col sm:flex-row justify-end space-y-2 sm:space-y-0 sm:space-x-3 pt-4">
            <button
              type="button"
              @click="$emit('cancel')"
              class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="!isFormValid"
              class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="creating">Saving...</span>
              <span v-else>Save Changes</span>
            </button>
          </div>
        </form>
        </div> <!-- End mobile-optimized content wrapper -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useAppointmentStore } from '../stores/appointment'

const props = defineProps({
  tooth: {
    type: Object,
    required: true
  },
  patientId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['save', 'cancel'])

const appointmentStore = useAppointmentStore()

const editData = reactive({
  condition: props.tooth.condition,
  surfaces: [...(props.tooth.surfaces || [])],
  surfaceConditions: [...(props.tooth.surface_conditions || [])],
  notes: props.tooth.notes || '',
  reason: '',
  appointmentId: props.tooth.appointment_id || null
})

const availableAppointments = ref([])
const creating = ref(false)

const conditionOptions = {
  healthy: 'Healthy',
  decay: 'Decay/Caries',
  filled: 'Filled',
  crowned: 'Crowned',
  root_canal: 'Root Canal',
  extracted: 'Extracted',
  implant: 'Implant',
  bridge: 'Bridge',
  missing: 'Missing'
}

const availableSurfaces = [
  { value: 'mesial', label: 'Mesial (M)' },
  { value: 'distal', label: 'Distal (D)' },
  { value: 'occlusal', label: 'Occlusal (O)' },
  { value: 'buccal', label: 'Buccal (B)' },
  { value: 'lingual', label: 'Lingual (L)' },
  { value: 'incisal', label: 'Incisal (I)' }
]

const permanentToothNames = {
  '1': 'Upper Right Central Incisor', '2': 'Upper Right Lateral Incisor', 
  '3': 'Upper Right Canine', '4': 'Upper Right First Premolar',
  '5': 'Upper Right Second Premolar', '6': 'Upper Right First Molar', 
  '7': 'Upper Right Second Molar', '8': 'Upper Right Third Molar',
  '9': 'Upper Left Central Incisor', '10': 'Upper Left Lateral Incisor',
  '11': 'Upper Left Canine', '12': 'Upper Left First Premolar',
  '13': 'Upper Left Second Premolar', '14': 'Upper Left First Molar',
  '15': 'Upper Left Second Molar', '16': 'Upper Left Third Molar',
  '17': 'Lower Left Third Molar', '18': 'Lower Left Second Molar',
  '19': 'Lower Left First Molar', '20': 'Lower Left Second Premolar',
  '21': 'Lower Left First Premolar', '22': 'Lower Left Canine',
  '23': 'Lower Left Lateral Incisor', '24': 'Lower Left Central Incisor',
  '25': 'Lower Right Central Incisor', '26': 'Lower Right Lateral Incisor',
  '27': 'Lower Right Canine', '28': 'Lower Right First Premolar',
  '29': 'Lower Right Second Premolar', '30': 'Lower Right First Molar',
  '31': 'Lower Right Second Molar', '32': 'Lower Right Third Molar'
}

const primaryToothNames = {
  'A': 'Upper Right Central Incisor', 'B': 'Upper Right Lateral Incisor',
  'C': 'Upper Right Canine', 'D': 'Upper Right First Molar', 'E': 'Upper Right Second Molar',
  'F': 'Upper Left Second Molar', 'G': 'Upper Left First Molar',
  'H': 'Upper Left Canine', 'I': 'Upper Left Lateral Incisor', 'J': 'Upper Left Central Incisor',
  'K': 'Lower Left Central Incisor', 'L': 'Lower Left Lateral Incisor',
  'M': 'Lower Left Canine', 'N': 'Lower Left First Molar', 'O': 'Lower Left Second Molar',
  'P': 'Lower Right Second Molar', 'Q': 'Lower Right First Molar',
  'R': 'Lower Right Canine', 'S': 'Lower Right Lateral Incisor', 'T': 'Lower Right Central Incisor'
}

const hasChanges = computed(() => {
  return editData.condition !== props.tooth.condition ||
         JSON.stringify(editData.surfaces.sort()) !== JSON.stringify((props.tooth.surfaces || []).sort()) ||
         JSON.stringify(editData.surfaceConditions) !== JSON.stringify(props.tooth.surface_conditions || []) ||
         editData.notes !== (props.tooth.notes || '')
})

const surfacesChanged = computed(() => {
  return JSON.stringify(editData.surfaces.sort()) !== JSON.stringify((props.tooth.surfaces || []).sort()) ||
         JSON.stringify(editData.surfaceConditions) !== JSON.stringify(props.tooth.surface_conditions || [])
})

const isFormValid = computed(() => {
  return editData.condition && editData.reason.trim() && hasChanges.value
})

const getToothName = () => {
  const toothNumber = props.tooth.tooth_number
  const isPrimary = isNaN(toothNumber)
  
  if (isPrimary) {
    return primaryToothNames[toothNumber] || 'Unknown Primary Tooth'
  } else {
    return permanentToothNames[toothNumber] || 'Unknown Permanent Tooth'
  }
}

const getCurrentColor = () => {
  const colors = {
    healthy: '#90EE90',
    decay: '#FFB6C1',
    filled: '#87CEEB',
    crowned: '#DDA0DD',
    root_canal: '#F0E68C',
    extracted: '#FF6347',
    implant: '#20B2AA',
    bridge: '#DEB887',
    missing: '#D3D3D3'
  }
  return colors[props.tooth.condition] || '#ffffff'
}

// Surface condition methods
const hasSurfaceCondition = (surface) => {
  return editData.surfaceConditions.some(sc => sc.surface === surface)
}

const getSurfaceCondition = (surface) => {
  const surfaceCondition = editData.surfaceConditions.find(sc => sc.surface === surface)
  return surfaceCondition ? surfaceCondition.condition : ''
}

const updateSurfaceCondition = (surface, condition) => {
  // Remove existing condition for this surface
  editData.surfaceConditions = editData.surfaceConditions.filter(sc => sc.surface !== surface)

  // Add new condition if not empty
  if (condition) {
    editData.surfaceConditions.push({
      surface: surface,
      condition: condition
    })
  }

  // Update surfaces array for backward compatibility
  updateSurfacesFromConditions()
}

const toggleSurfaceCondition = (surface) => {
  if (hasSurfaceCondition(surface)) {
    // Remove the condition
    editData.surfaceConditions = editData.surfaceConditions.filter(sc => sc.surface !== surface)
  } else {
    // Add a default condition (decay)
    editData.surfaceConditions.push({
      surface: surface,
      condition: 'decay'
    })
  }

  // Update surfaces array for backward compatibility
  updateSurfacesFromConditions()
}

const updateSurfacesFromConditions = () => {
  // Update the surfaces array based on surface conditions
  editData.surfaces = editData.surfaceConditions.map(sc => sc.surface)
}

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatSurfaceConditions = (surfaceConditions) => {
  if (!surfaceConditions || surfaceConditions.length === 0) {
    return 'No specific surface conditions'
  }
  return surfaceConditions.map(sc =>
    `${sc.surface}: ${formatCondition(sc.condition)}`
  ).join(', ')
}

const formatAppointmentOption = (appointment) => {
  const date = new Date(appointment.start_time).toLocaleDateString()
  const time = new Date(appointment.start_time).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit'
  })
  return `${date} ${time} - ${appointment.title || 'Appointment'}`
}

const loadAppointments = async () => {
  try {
    // Load recent and upcoming appointments for this patient
    await appointmentStore.fetchPatientAppointments(props.patientId)

    // Filter to get relevant appointments (completed or in progress for linking)
    const now = new Date()
    const thirtyDaysAgo = new Date(now.getTime() - (30 * 24 * 60 * 60 * 1000))

    availableAppointments.value = appointmentStore.appointments.filter(appointment => {
      const appointmentDate = new Date(appointment.start_time)
      const isRecent = appointmentDate >= thirtyDaysAgo
      const isCompletedOrInProgress = ['completed', 'in_progress'].includes(appointment.status)

      return isRecent && isCompletedOrInProgress
    }).sort((a, b) => new Date(b.start_time) - new Date(a.start_time)) // Most recent first
  } catch (error) {
    console.error('Error loading appointments:', error)
  }
}

const handleSave = async () => {
  if (!isFormValid.value) return

  creating.value = true

  try {
    emit('save', {
      tooth_number: props.tooth.tooth_number,
      condition: editData.condition,
      surfaces: editData.surfaces,
      surface_conditions: editData.surfaceConditions,
      notes: editData.notes,
      reason: editData.reason,
      appointment_id: editData.appointmentId
    })
  } finally {
    // Reset after a short delay to show the saving state
    setTimeout(() => {
      creating.value = false
    }, 500)
  }
}

onMounted(() => {
  loadAppointments()
})
</script>

<style scoped>
@import "../styles/main.css";

.surfaces-grid label {
  @apply transition-colors;
}

.changes-preview {
  @apply animate-pulse;
}

.modal-header h3 {
  @apply text-gray-900;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .surfaces-grid {
    @apply grid-cols-1;
  }

  .modal-header {
    @apply flex-col items-start space-y-2;
  }

  .modal-header button {
    @apply self-end;
  }

  .form-group label {
    @apply text-sm;
  }

  .surface-condition-item {
    @apply p-2;
  }

  .chart-type-buttons {
    @apply flex-col space-x-0 space-y-2;
  }

  .overall-condition-note {
    @apply text-xs;
  }
}
</style>