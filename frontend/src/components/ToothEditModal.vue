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

          <!-- Affected surfaces -->
          <div class="form-group">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Affected Surfaces
            </label>
            <div class="surfaces-grid grid grid-cols-2 gap-2">
              <label 
                v-for="surface in availableSurfaces" 
                :key="surface.value"
                class="flex items-center p-2 border border-gray-300 rounded cursor-pointer hover:bg-gray-50"
                :class="{ 'bg-blue-50 border-blue-300': editData.surfaces.includes(surface.value) }"
              >
                <input
                  type="checkbox"
                  :value="surface.value"
                  v-model="editData.surfaces"
                  class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50 mr-2"
                >
                <span class="text-sm">{{ surface.label }}</span>
              </label>
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
                <span class="font-medium">Surfaces:</span>
                <div class="ml-2">
                  <div>{{ (tooth.surfaces || []).join(', ') || 'None' }}</div>
                  <div class="text-xs">↓</div>
                  <div>{{ editData.surfaces.join(', ') || 'None' }}</div>
                </div>
              </div>
              
              <div v-if="editData.notes !== (tooth.notes || '')" class="mb-1">
                <span class="font-medium">Notes:</span>
                <span class="ml-2">{{ editData.notes ? 'Updated' : 'Cleared' }}</span>
              </div>
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
              :disabled="!isFormValid"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Save Changes
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'

const props = defineProps({
  tooth: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['save', 'cancel'])

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

const surfacesChanged = computed(() => {
  return JSON.stringify(editData.surfaces.sort()) !== JSON.stringify((props.tooth.surfaces || []).sort())
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

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const handleSave = () => {
  if (!isFormValid.value) return
  
  emit('save', {
    tooth_number: props.tooth.tooth_number,
    condition: editData.condition,
    surfaces: editData.surfaces,
    notes: editData.notes,
    reason: editData.reason
  })
}
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

@media (max-width: 640px) {
  .surfaces-grid {
    @apply grid-cols-1;
  }
}
</style>