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
             <option value="permanent">Adult Teeth (FDI)</option>
             <option value="primary">Primary Teeth (A-T)</option>
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
      <!-- Upper Jaw Label -->
      <div class="jaw-label-container upper-jaw-label">
        <span class="text-sm font-medium text-gray-600">Upper Jaw</span>
      </div>

      <!-- Upper Jaw -->
      <div class="jaw-section upper-jaw mb-8">
        <!-- Upper Left Quadrant (FDI Quadrant 2) -->
        <div class="quadrant upper-left">
          <div class="quadrant-label">Quadrant 2</div>
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

        <!-- Upper Right Quadrant (FDI Quadrant 1) -->
        <div class="quadrant upper-right">
          <div class="quadrant-label">Quadrant 1</div>
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
      </div>

      <!-- Lower Jaw -->
      <div class="jaw-section lower-jaw">
        <!-- Lower Left Quadrant (FDI Quadrant 3) -->
        <div class="quadrant lower-left">
          <div class="quadrant-label">Quadrant 3</div>
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

        <!-- Lower Right Quadrant (FDI Quadrant 4) -->
        <div class="quadrant lower-right">
          <div class="quadrant-label">Quadrant 4</div>
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

      <!-- Lower Jaw Label -->
      <div class="jaw-label-container lower-jaw-label">
        <span class="text-sm font-medium text-gray-600">Lower Jaw</span>
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
            :class="{ 'border-gray-400': condition === 'healthy' }"
          ></div>
          <span class="legend-text">{{ formatCondition(condition) }}</span>
        </div>
      </div>
       <div class="legend-note text-xs text-gray-500 mt-2">
         <em>White = Healthy, Colored = Requires attention | FDI Notation (ISO 3950)</em>
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
  healthy: '#FFFFFF',
  decay: '#FFB6C1',
  filled: '#87CEEB',
  crowned: '#DDA0DD',
  root_canal: '#F0E68C',
  extracted: '#FF6347',
  implant: '#20B2AA',
  bridge: '#DEB887',
  missing: '#D3D3D3'
}

// Permanent teeth numbering - FDI World Dental Federation notation (ISO 3950)
// Two-digit system: First digit = quadrant, Second digit = tooth position
// Quadrant 1: Upper right (11-18), Quadrant 2: Upper left (21-28)
// Quadrant 3: Lower left (31-38), Quadrant 4: Lower right (41-48)
const permanentTeethQuadrants = {
  'upper-right': [11, 12, 13, 14, 15, 16, 17, 18],
  'upper-left': [21, 22, 23, 24, 25, 26, 27, 28],
  'lower-left': [31, 32, 33, 34, 35, 36, 37, 38],
  'lower-right': [41, 42, 43, 44, 45, 46, 47, 48]
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

.jaw-label-container {
  @apply text-center mb-4;
}

.upper-jaw-label {
  @apply mb-6;
}

.lower-jaw-label {
  @apply mt-6;
}

/* Mobile responsive adjustments */
@media (max-width: 640px) {
  .upper-jaw-label {
    @apply mb-4;
  }

  .lower-jaw-label {
    @apply mt-4;
  }
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
  @apply flex-row;
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