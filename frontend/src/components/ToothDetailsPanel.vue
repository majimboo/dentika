<template>
  <div class="tooth-details-panel bg-white border rounded-lg shadow-lg p-4 mt-4">
    <div class="panel-header flex justify-between items-start mb-4">
      <div>
        <h4 class="text-lg font-semibold text-gray-900">
          Tooth {{ tooth.tooth_number }} Details
        </h4>
        <p class="text-sm text-gray-600">{{ getToothName() }}</p>
      </div>
      <button 
        @click="$emit('close')"
        class="text-gray-400 hover:text-gray-600 transition-colors"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
    </div>

    <div class="tooth-info grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Current Condition -->
      <div class="condition-section">
        <div class="section-header mb-2">
          <h5 class="font-medium text-gray-900">Current Condition</h5>
        </div>
        
        <div class="condition-display flex items-center mb-3">
          <div 
            class="condition-color w-6 h-6 rounded border mr-3"
            :style="{ backgroundColor: getToothColor() }"
          ></div>
          <span class="condition-text text-lg font-medium">
            {{ formatCondition(tooth.condition) }}
          </span>
        </div>

        <!-- Edit Condition (if editable) -->
        <div v-if="editable" class="edit-condition mb-3">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Update Condition:
          </label>
          <select 
            v-model="editData.condition"
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option v-for="(label, value) in conditionOptions" :key="value" :value="value">
              {{ label }}
            </option>
          </select>
        </div>

        <!-- Affected Surfaces -->
        <div class="surfaces-section">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Affected Surfaces:
          </label>
          <div v-if="!editable" class="surfaces-display">
            <div v-if="tooth.surfaces && tooth.surfaces.length > 0" class="flex flex-wrap gap-2">
              <span 
                v-for="surface in tooth.surfaces" 
                :key="surface"
                class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                {{ formatSurface(surface) }}
              </span>
            </div>
            <span v-else class="text-gray-500 text-sm">No specific surfaces affected</span>
          </div>
          
          <div v-else class="surfaces-edit">
            <div class="flex flex-wrap gap-2">
              <label 
                v-for="surface in availableSurfaces" 
                :key="surface.value"
                class="inline-flex items-center"
              >
                <input
                  type="checkbox"
                  :value="surface.value"
                  v-model="editData.surfaces"
                  class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
                >
                <span class="ml-1 text-sm text-gray-700">{{ surface.label }}</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      <!-- Notes and History -->
      <div class="notes-section">
        <div class="section-header mb-2">
          <h5 class="font-medium text-gray-900">Notes & Information</h5>
        </div>

        <!-- Notes -->
        <div class="notes-display mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Notes:
          </label>
          <div v-if="!editable" class="notes-text">
            <p v-if="tooth.notes" class="text-sm text-gray-900 bg-gray-50 p-2 rounded">
              {{ tooth.notes }}
            </p>
            <p v-else class="text-sm text-gray-500 italic">No notes recorded</p>
          </div>
          
          <textarea
            v-else
            v-model="editData.notes"
            rows="3"
            placeholder="Enter notes about this tooth..."
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
          ></textarea>
        </div>

        <!-- Last Updated -->
        <div class="last-updated mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Last Updated:
          </label>
          <p class="text-sm text-gray-600">
            {{ tooth.last_updated ? formatDate(tooth.last_updated) : 'Never' }}
          </p>
        </div>

        <!-- Edit Reason (if editing) -->
        <div v-if="editable" class="edit-reason">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Reason for Change:
          </label>
          <input
            type="text"
            v-model="editData.reason"
            placeholder="e.g., Regular checkup, treatment completed..."
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
          >
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div v-if="editable" class="action-buttons flex justify-end space-x-3 mt-4 pt-4 border-t">
      <button
        @click="cancelEdit"
        class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
      >
        Cancel
      </button>
      <button
        @click="saveChanges"
        :disabled="!hasChanges"
        class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Save Changes
      </button>
    </div>

    <!-- View History Button -->
    <div class="view-history mt-4 pt-4 border-t">
      <button
        @click="viewHistory"
        class="text-sm text-blue-600 hover:text-blue-800 font-medium"
      >
        View Change History
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, watch } from 'vue'

const props = defineProps({
  tooth: {
    type: Object,
    required: true
  },
  editable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update', 'close', 'viewHistory'])

const editData = reactive({
  condition: props.tooth.condition,
  surfaces: [...(props.tooth.surfaces || [])],
  notes: props.tooth.notes || '',
  reason: ''
})

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
         editData.notes !== (props.tooth.notes || '')
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

const getToothColor = () => {
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

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatSurface = (surface) => {
  const surfaceMap = {
    mesial: 'Mesial',
    distal: 'Distal', 
    occlusal: 'Occlusal',
    buccal: 'Buccal',
    lingual: 'Lingual',
    incisal: 'Incisal'
  }
  return surfaceMap[surface] || surface
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const saveChanges = () => {
  emit('update', {
    tooth_number: props.tooth.tooth_number,
    condition: editData.condition,
    surfaces: editData.surfaces,
    notes: editData.notes,
    reason: editData.reason
  })
}

const cancelEdit = () => {
  // Reset edit data
  editData.condition = props.tooth.condition
  editData.surfaces = [...(props.tooth.surfaces || [])]
  editData.notes = props.tooth.notes || ''
  editData.reason = ''
  
  emit('close')
}

const viewHistory = () => {
  emit('viewHistory', props.tooth)
}

// Reset edit data when tooth changes
watch(() => props.tooth, (newTooth) => {
  editData.condition = newTooth.condition
  editData.surfaces = [...(newTooth.surfaces || [])]
  editData.notes = newTooth.notes || ''
  editData.reason = ''
}, { immediate: true })
</script>

<style scoped>
.tooth-details-panel {
  max-width: 800px;
}

.section-header {
  @apply border-b border-gray-200 pb-1;
}

.condition-color {
  @apply border-gray-300 shadow-sm;
}

.surfaces-display .inline-flex {
  @apply transition-colors;
}

.surfaces-edit label {
  @apply py-1 px-2 rounded hover:bg-gray-50 transition-colors cursor-pointer;
}

.notes-text .bg-gray-50 {
  @apply border border-gray-200;
}

.action-buttons button:disabled {
  @apply bg-gray-400 cursor-not-allowed;
}

@media (max-width: 768px) {
  .tooth-info {
    @apply grid-cols-1;
  }
  
  .action-buttons {
    @apply flex-col space-y-2 space-x-0;
  }
  
  .action-buttons button {
    @apply w-full;
  }
}
</style>