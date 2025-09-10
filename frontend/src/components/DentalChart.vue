<template>
  <div class="dental-chart">
    <div class="dental-chart-header mb-4">
      <h3 class="text-lg font-semibold text-gray-900">Dental Chart</h3>
      <div class="flex items-center space-x-4">
        <div class="flex items-center space-x-2">
          <label class="text-sm font-medium text-gray-700">Chart Type:</label>
          <select 
            v-model="chartType" 
            @change="switchChartType"
            class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="permanent">Adult Teeth (32)</option>
            <option value="primary">Primary Teeth (20)</option>
          </select>
        </div>
        
        <div v-if="isDoctor" class="flex items-center space-x-2">
          <button 
            @click="toggleEditMode"
            class="px-3 py-1 text-sm rounded-md transition-colors"
            :class="editMode ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-blue-100 text-blue-700 hover:bg-blue-200'"
          >
            {{ editMode ? 'Exit Edit' : 'Edit Chart' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <!-- Chart Container -->
    <div v-else class="dental-chart-container">
      <!-- Upper Jaw -->
      <div class="jaw-section upper-jaw mb-8">
        <div class="jaw-label text-center mb-2">
          <span class="text-sm font-medium text-gray-600">Upper Jaw</span>
        </div>
        
        <!-- Upper Right Quadrant -->
        <div class="quadrant upper-right">
          <div class="quadrant-label">Upper Right</div>
          <div class="teeth-row">
            <ToothComponent
              v-for="tooth in getQuadrantTeeth('upper-right')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :selected="selectedTooth === tooth.tooth_number"
              :editable="editMode && isDoctor"
              @click="selectTooth"
              @update="updateTooth"
            />
          </div>
        </div>

        <!-- Upper Left Quadrant -->
        <div class="quadrant upper-left">
          <div class="quadrant-label">Upper Left</div>
          <div class="teeth-row">
            <ToothComponent
              v-for="tooth in getQuadrantTeeth('upper-left')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :selected="selectedTooth === tooth.tooth_number"
              :editable="editMode && isDoctor"
              @click="selectTooth"
              @update="updateTooth"
            />
          </div>
        </div>
      </div>

      <!-- Lower Jaw -->
      <div class="jaw-section lower-jaw">
        <div class="jaw-label text-center mb-2">
          <span class="text-sm font-medium text-gray-600">Lower Jaw</span>
        </div>
        
        <!-- Lower Left Quadrant -->
        <div class="quadrant lower-left">
          <div class="quadrant-label">Lower Left</div>
          <div class="teeth-row">
            <ToothComponent
              v-for="tooth in getQuadrantTeeth('lower-left')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :selected="selectedTooth === tooth.tooth_number"
              :editable="editMode && isDoctor"
              @click="selectTooth"
              @update="updateTooth"
            />
          </div>
        </div>

        <!-- Lower Right Quadrant -->
        <div class="quadrant lower-right">
          <div class="quadrant-label">Lower Right</div>
          <div class="teeth-row">
            <ToothComponent
              v-for="tooth in getQuadrantTeeth('lower-right')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :selected="selectedTooth === tooth.tooth_number"
              :editable="editMode && isDoctor"
              @click="selectTooth"
              @update="updateTooth"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="dental-chart-legend mt-6">
      <div class="legend-title text-sm font-medium text-gray-700 mb-2">Legend:</div>
      <div class="legend-items grid grid-cols-3 md:grid-cols-5 gap-2 text-xs">
        <div v-for="(color, condition) in legendItems" :key="condition" class="legend-item flex items-center">
          <div 
            class="legend-color w-4 h-4 rounded border border-gray-300 mr-2" 
            :style="{ backgroundColor: color }"
          ></div>
          <span class="legend-text">{{ formatCondition(condition) }}</span>
        </div>
      </div>
    </div>

    <!-- Selected Tooth Details -->
    <ToothDetailsPanel
      v-if="selectedTooth && selectedToothData"
      :tooth="selectedToothData"
      :editable="editMode && isDoctor"
      @update="updateTooth"
      @close="deselectTooth"
    />

    <!-- Tooth Edit Modal -->
    <ToothEditModal
      v-if="editingTooth"
      :tooth="editingTooth"
      @save="saveToothEdit"
      @cancel="cancelToothEdit"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useDentalRecordStore } from '../stores/dentalRecord'
import ToothComponent from './ToothComponent.vue'
import ToothDetailsPanel from './ToothDetailsPanel.vue'
import ToothEditModal from './ToothEditModal.vue'

const props = defineProps({
  patientId: {
    type: Number,
    required: true
  },
  readonly: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['toothUpdated', 'chartUpdated'])

const authStore = useAuthStore()
const dentalRecordStore = useDentalRecordStore()

const chartType = ref('permanent')
const editMode = ref(false)
const selectedTooth = ref(null)
const editingTooth = ref(null)

const loading = computed(() => dentalRecordStore.isLoading)
const isDoctor = computed(() => authStore.isDoctor || authStore.isSuperAdmin)
const teethData = computed(() => dentalRecordStore.currentRecord?.teeth_data || [])

const selectedToothData = computed(() => {
  if (!selectedTooth.value) return null
  return teethData.value.find(tooth => tooth.tooth_number === selectedTooth.value)
})

const legendItems = {
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

// Permanent teeth numbering (1-32)
const permanentTeethQuadrants = {
  'upper-right': [1, 2, 3, 4, 5, 6, 7, 8],
  'upper-left': [9, 10, 11, 12, 13, 14, 15, 16],
  'lower-left': [17, 18, 19, 20, 21, 22, 23, 24],
  'lower-right': [25, 26, 27, 28, 29, 30, 31, 32]
}

// Primary teeth lettering (A-T)
const primaryTeethQuadrants = {
  'upper-right': ['A', 'B', 'C', 'D', 'E'],
  'upper-left': ['F', 'G', 'H', 'I', 'J'],
  'lower-left': ['K', 'L', 'M', 'N', 'O'],
  'lower-right': ['P', 'Q', 'R', 'S', 'T']
}

const getQuadrantTeeth = (quadrant) => {
  const toothNumbers = chartType.value === 'permanent' 
    ? permanentTeethQuadrants[quadrant]
    : primaryTeethQuadrants[quadrant]
    
  return toothNumbers
    .map(num => teethData.value.find(tooth => tooth.tooth_number === String(num)))
    .filter(Boolean)
}

const selectTooth = (toothNumber) => {
  selectedTooth.value = selectedTooth.value === toothNumber ? null : toothNumber
  dentalRecordStore.setSelectedTooth(selectedTooth.value)
}

const deselectTooth = () => {
  selectedTooth.value = null
  dentalRecordStore.clearSelectedTooth()
}

const toggleEditMode = () => {
  editMode.value = !editMode.value
  if (!editMode.value) {
    deselectTooth()
  }
}

const updateTooth = async (toothData) => {
  if (!dentalRecordStore.currentRecord) return
  
  const result = await dentalRecordStore.updateToothCondition(
    dentalRecordStore.currentRecord.id,
    toothData
  )
  
  if (result.success) {
    emit('toothUpdated', toothData)
    emit('chartUpdated')
  }
}

const switchChartType = async () => {
  const targetRecord = dentalRecordStore.dentalRecords.find(record => 
    record.record_type === chartType.value
  )
  
  if (targetRecord) {
    if (!targetRecord.is_active && isDoctor.value) {
      // Activate the record if it's not active
      await dentalRecordStore.activateDentalRecord(targetRecord.id)
    } else {
      dentalRecordStore.setCurrentRecord(targetRecord)
    }
  }
  
  deselectTooth()
}

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const saveToothEdit = async (toothData) => {
  await updateTooth(toothData)
  editingTooth.value = null
}

const cancelToothEdit = () => {
  editingTooth.value = null
}

// Initialize chart when component mounts
onMounted(async () => {
  if (props.patientId) {
    await dentalRecordStore.fetchPatientDentalRecords(props.patientId)
    
    // Set initial chart type based on available records
    const permanentRecord = dentalRecordStore.getPermanentTeethRecord
    const primaryRecord = dentalRecordStore.getPrimaryTeethRecord
    
    if (permanentRecord) {
      chartType.value = 'permanent'
      dentalRecordStore.setCurrentRecord(permanentRecord)
    } else if (primaryRecord) {
      chartType.value = 'primary'
      dentalRecordStore.setCurrentRecord(primaryRecord)
    }
  }
})

// Watch for readonly prop changes
watch(() => props.readonly, (newValue) => {
  if (newValue) {
    editMode.value = false
    deselectTooth()
  }
})
</script>

<style scoped>
@import "../styles/main.css";

.dental-chart {
  @apply bg-white rounded-lg shadow-sm border border-gray-200 p-6;
}

.dental-chart-container {
  @apply bg-gray-50 rounded-lg p-6;
}

.jaw-section {
  @apply flex flex-col items-center;
}

/* Wide screen layout - side by side quadrants */
@media (min-width: 1024px) {
  .upper-jaw {
    @apply flex-row justify-center items-start;
  }

  .lower-jaw {
    @apply flex-row justify-center items-start;
  }

  .upper-jaw .quadrant {
    @apply mr-8;
  }

  .lower-jaw .quadrant {
    @apply mr-8;
  }

  .lower-jaw .quadrant:last-child {
    @apply mr-0;
  }

  .upper-jaw .quadrant:last-child {
    @apply mr-0;
  }
}

.quadrant {
  @apply flex flex-col items-center m-4;
}

.quadrant-label {
  @apply text-xs font-medium text-gray-500 mb-2;
}

.teeth-row {
  @apply flex flex-wrap justify-center gap-1;
}

.upper-jaw {
  @apply border-b-2 border-gray-200 pb-4;
}

.upper-jaw .teeth-row {
  @apply flex-row;
}

.lower-jaw .teeth-row {
  @apply flex-row-reverse;
}

.dental-chart-legend {
  @apply bg-gray-50 rounded-lg p-4;
}

.legend-color {
  @apply border-gray-300;
}

.legend-text {
  @apply text-gray-700 capitalize;
}
</style>