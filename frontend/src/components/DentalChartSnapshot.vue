<template>
  <div class="dental-chart-snapshot">
    <div class="snapshot-header mb-4">
      <h4 class="text-md font-semibold text-gray-900">
        {{ formatChartType(chartType) }} Dental Chart
      </h4>
      <p class="text-sm text-gray-600">Historical snapshot - read-only view</p>
    </div>

    <!-- Chart Container -->
    <div class="dental-chart-container">
      <!-- Upper Jaw Label -->
      <div class="jaw-label-container upper-jaw-label">
        <span class="text-sm font-medium text-gray-600">Upper Jaw</span>
      </div>

      <!-- Upper Jaw -->
      <div class="jaw-section upper-jaw mb-8">
        <!-- Upper Left Quadrant -->
        <div class="quadrant upper-left">
          <div class="quadrant-label">{{ chartType === 'permanent' ? 'Quadrant 2' : 'Upper Left' }}</div>
          <div class="teeth-row">
            <ToothSnapshotComponent
              v-for="tooth in getQuadrantTeeth('upper-left')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :chartType="chartType"
            />
          </div>
        </div>

        <!-- Upper Right Quadrant -->
        <div class="quadrant upper-right">
          <div class="quadrant-label">{{ chartType === 'permanent' ? 'Quadrant 1' : 'Upper Right' }}</div>
          <div class="teeth-row">
            <ToothSnapshotComponent
              v-for="tooth in getQuadrantTeeth('upper-right')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :chartType="chartType"
            />
          </div>
        </div>
      </div>

      <!-- Lower Jaw -->
      <div class="jaw-section lower-jaw">
        <!-- Lower Left Quadrant -->
        <div class="quadrant lower-left">
          <div class="quadrant-label">{{ chartType === 'permanent' ? 'Quadrant 3' : 'Lower Left' }}</div>
          <div class="teeth-row">
            <ToothSnapshotComponent
              v-for="tooth in getQuadrantTeeth('lower-left')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :chartType="chartType"
            />
          </div>
        </div>

        <!-- Lower Right Quadrant -->
        <div class="quadrant lower-right">
          <div class="quadrant-label">{{ chartType === 'permanent' ? 'Quadrant 4' : 'Lower Right' }}</div>
          <div class="teeth-row">
            <ToothSnapshotComponent
              v-for="tooth in getQuadrantTeeth('lower-right')"
              :key="tooth.tooth_number"
              :tooth="tooth"
              :chartType="chartType"
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
        <em>Historical snapshot - {{ formatChartType(chartType) }} chart notation</em>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ToothSnapshotComponent from './ToothSnapshotComponent.vue'

const props = defineProps({
  snapshotData: {
    type: Array,
    required: true
  },
  chartType: {
    type: String,
    required: true
  }
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

// Permanent teeth numbering - FDI World Dental Federation notation
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
  const toothNumbers = props.chartType === 'permanent'
    ? permanentTeethQuadrants[quadrant]
    : primaryTeethQuadrants[quadrant]

  return toothNumbers
    .map(num => props.snapshotData.find(tooth => tooth.tooth_number === String(num)))
    .filter(Boolean)
}

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatChartType = (chartType) => {
  return chartType === 'permanent' ? 'Adult (FDI)' : 'Child (A-T)'
}
</script>

<style scoped>
@import "../styles/main.css";

.dental-chart-snapshot {
  @apply bg-gray-50 rounded-lg p-6;
}

.dental-chart-container {
  @apply bg-white rounded-lg p-6 border border-gray-200;
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

.dental-chart-legend {
  @apply bg-gray-100 rounded-lg p-4;
}

.legend-color {
  @apply border-gray-300;
}

.legend-text {
  @apply text-gray-700 capitalize;
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
</style>