<template>
  <div class="dental-chart">
    <!-- Read-only banner for non-doctors -->
    <div v-if="!isDoctor" class="read-only-banner bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <font-awesome-icon icon="fa-solid fa-eye" class="w-5 h-5 text-blue-600" />
          </div>
          <div class="ml-3">
            <h4 class="text-sm font-medium text-blue-800">
              Read-Only Access
            </h4>
            <p class="text-sm text-blue-600">
              You are viewing this dental chart in read-only mode. Only doctors can modify chart data and settings.
            </p>
          </div>
        </div>
        <div class="flex-shrink-0">
          <font-awesome-icon icon="fa-solid fa-lock" class="w-4 h-4 text-blue-500" />
        </div>
      </div>
    </div>

    <div class="dental-chart-header mb-6">
      <div class="header-top flex flex-col md:flex-row justify-between items-start md:items-center mb-4 gap-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Dental Chart</h3>
          <p class="text-sm text-gray-600">Interactive dental chart with visit tracking</p>
        </div>

        <div class="header-actions flex items-center space-x-3">
          <!-- Chart History Button -->
          <button
            @click="showChartHistory = !showChartHistory"
            class="px-3 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
          >
            <font-awesome-icon icon="fa-solid fa-history" class="w-4 h-4 mr-1" />
            {{ showChartHistory ? 'Hide History' : 'Chart History' }}
          </button>

          <!-- Edit mode toggle -->
          <div v-if="isDoctor">
            <button
              @click="toggleEditMode"
              class="px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="editMode ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-blue-100 text-blue-700 hover:bg-blue-200'"
            >
              {{ editMode ? 'Exit Edit' : 'Edit Chart' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Chart Type Selection - More Prominent -->
      <div class="chart-type-selection bg-gray-50 rounded-lg p-4 border border-gray-200">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <label class="text-sm font-medium text-gray-700 mb-2 block">Chart Type:</label>
            <div class="chart-type-buttons flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-2">
              <button
                @click="setChartType('permanent')"
                :disabled="!isDoctor"
                class="chart-type-btn"
                :class="chartType === 'permanent' ? 'chart-type-btn-active' : 'chart-type-btn-inactive'"
                title="Switch to adult teeth chart (automatically saved)"
              >
                <font-awesome-icon icon="fa-solid fa-user" class="w-4 h-4 mr-2" />
                Adult Teeth
                <span class="chart-notation">(FDI Notation)</span>
              </button>
              <button
                @click="setChartType('primary')"
                :disabled="!isDoctor"
                class="chart-type-btn"
                :class="chartType === 'primary' ? 'chart-type-btn-active' : 'chart-type-btn-inactive'"
                title="Switch to child teeth chart (automatically saved)"
              >
                <font-awesome-icon icon="fa-solid fa-child" class="w-4 h-4 mr-2" />
                Child Teeth
                <span class="chart-notation">(A-T Notation)</span>
              </button>
            </div>

            <!-- Chart type status -->
            <div class="chart-type-status mt-2">
              <div class="text-xs text-gray-600">
                <font-awesome-icon icon="fa-solid fa-check-circle" class="w-3 h-3 text-green-500 mr-1" />
                Chart type selection is automatically saved
              </div>
            </div>
          </div>

          <!-- Chart Info -->
          <div class="chart-info text-sm text-gray-600">
            <div class="info-stats flex space-x-4">
              <div class="stat-item">
                <span class="font-medium">Total:</span> {{ totalTeethCount }}
              </div>
              <div class="stat-item">
                <span class="font-medium">Healthy:</span> {{ healthyTeethCount }}
              </div>
              <div class="stat-item">
                <span class="font-medium">Issues:</span> {{ problemTeethCount }}
              </div>
            </div>
          </div>
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

    <!-- Chart History Section -->
    <div v-if="showChartHistory" class="chart-history-section mt-8">
      <DentalChartHistory :patientId="patientId" />
    </div>

    <!-- Tooth Edit Modal -->
    <ToothEditModal
      v-if="editingTooth"
      :tooth="editingTooth"
      :patientId="patientId"
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
import DentalChartHistory from './DentalChartHistory.vue'

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
const showChartHistory = ref(false)

const loading = computed(() => dentalRecordStore.isLoading)
const isDoctor = computed(() => authStore.isDoctor || authStore.isSuperAdmin)
const teethData = computed(() => dentalRecordStore.currentRecord?.teeth_data || [])

const selectedToothData = computed(() => {
  if (!selectedTooth.value) return null
  return teethData.value.find(tooth => tooth.tooth_number === selectedTooth.value)
})

const totalTeethCount = computed(() => teethData.value.length)

const healthyTeethCount = computed(() => {
  return teethData.value.filter(tooth => tooth.condition === 'healthy').length
})

const problemTeethCount = computed(() => {
  return teethData.value.filter(tooth =>
    ['decay', 'extracted', 'missing'].includes(tooth.condition)
  ).length
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
  if (editMode.value && isDoctor.value) {
    // In edit mode, directly open the edit modal
    const tooth = teethData.value.find(t => t.tooth_number === toothNumber)
    if (tooth) {
      editingTooth.value = tooth
    }
  } else {
    // In view mode, just select the tooth for details
    selectedTooth.value = selectedTooth.value === toothNumber ? null : toothNumber
    dentalRecordStore.setSelectedTooth(selectedTooth.value)
  }
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

  let result
  if (toothData.appointment_id) {
    // Use the appointment-linked update method
    result = await dentalRecordStore.updateToothConditionWithAppointment(
      dentalRecordStore.currentRecord.id,
      toothData,
      toothData.appointment_id
    )
  } else {
    // Use the regular update method
    result = await dentalRecordStore.updateToothCondition(
      dentalRecordStore.currentRecord.id,
      toothData
    )
  }

  if (result.success) {
    emit('toothUpdated', toothData)
    emit('chartUpdated')
  }
}

const setChartType = async (type) => {
  chartType.value = type
  await switchChartType()
}

const switchChartType = async () => {
  const targetRecord = dentalRecordStore.dentalRecords.find(record =>
    record.record_type === chartType.value
  )

  if (targetRecord) {
    if (isDoctor.value) {
      // For doctors, always activate the selected record to persist the choice
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

    // Get the currently active record to determine chart type
    const activeRecord = dentalRecordStore.getActiveDentalRecord

    if (activeRecord) {
      // Set chart type based on the active record
      chartType.value = activeRecord.record_type
      dentalRecordStore.setCurrentRecord(activeRecord)
    } else {
      // If no active record, prioritize permanent teeth but check what's available
      const permanentRecord = dentalRecordStore.dentalRecords.find(r => r.record_type === 'permanent')
      const primaryRecord = dentalRecordStore.dentalRecords.find(r => r.record_type === 'primary')

      if (permanentRecord) {
        chartType.value = 'permanent'
        if (isDoctor.value) {
          // Activate permanent record for doctors
          await dentalRecordStore.activateDentalRecord(permanentRecord.id)
        } else {
          dentalRecordStore.setCurrentRecord(permanentRecord)
        }
      } else if (primaryRecord) {
        chartType.value = 'primary'
        if (isDoctor.value) {
          // Activate primary record for doctors
          await dentalRecordStore.activateDentalRecord(primaryRecord.id)
        } else {
          dentalRecordStore.setCurrentRecord(primaryRecord)
        }
      }
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

/* Chart Type Selection Styles */
.chart-type-selection {
  @apply transition-all duration-200;
}

.chart-type-btn {
  @apply flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200;
  @apply border-2;
}

.chart-type-btn-active {
  @apply bg-blue-100 border-blue-300 text-blue-800;
}

.chart-type-btn-inactive {
  @apply bg-white border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400;
}

.chart-type-btn:disabled {
  @apply opacity-50 cursor-not-allowed hover:bg-white hover:border-gray-300;
}

.chart-notation {
  @apply ml-2 text-xs opacity-75;
}

.info-stats {
  @apply text-sm;
}

.stat-item {
  @apply whitespace-nowrap;
}

/* Chart History Section */
.chart-history-section {
  @apply border-t border-gray-200 pt-8;
}

/* Read-only banner styles */
.read-only-banner {
  @apply transition-all duration-200;
}

.read-only-banner:hover {
  @apply bg-blue-100 border-blue-300;
}

/* Enhanced disabled styles for chart type buttons */
.chart-type-btn:disabled .chart-notation {
  @apply opacity-50;
}

/* Mobile responsive read-only banner */
@media (max-width: 640px) {
  .read-only-banner {
    @apply px-3 py-3;
  }

  .read-only-banner .flex {
    @apply flex-col items-start space-y-2;
  }

  .read-only-banner .justify-between {
    @apply justify-start;
  }
}
</style>