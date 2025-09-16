<template>
  <div
    class="tooth-snapshot-component"
    :class="[
      'tooth-' + tooth.tooth_number,
      {
        'tooth-has-notes': tooth.notes
      }
    ]"
    :title="getTooltipText()"
  >
    <div class="tooth-visual">
      <!-- Tooth Shape SVG -->
      <svg
        class="tooth-svg"
        viewBox="0 0 60 60"
        xmlns="http://www.w3.org/2000/svg"
      >
        <!-- Rotated tooth group -->
        <g transform="rotate(45 30 30)">
          <!-- Tooth outline -->
          <circle
            cx="30"
            cy="30"
            r="18"
            :fill="getToothColor()"
            stroke="#6B7280"
            stroke-width="2"
            class="tooth-shape"
          />

          <!-- Surface division lines (for reference) -->
          <g class="surface-dividers">
            <!-- Vertical center line -->
            <line x1="30" y1="12" x2="30" y2="22" stroke="#6B7280" stroke-width="1" opacity="0.5" />
            <line x1="30" y1="38" x2="30" y2="48" stroke="#6B7280" stroke-width="1" opacity="0.5" />
            <!-- Horizontal center line -->
            <line x1="12" y1="30" x2="22" y2="30" stroke="#6B7280" stroke-width="1" opacity="0.5" />
            <line x1="38" y1="30" x2="48" y2="30" stroke="#6B7280" stroke-width="1" opacity="0.5" />
          </g>

          <!-- Surface indicators (dots) -->
          <g v-if="tooth.surface_conditions && tooth.surface_conditions.length > 0">
            <circle
              v-for="surfaceCondition in tooth.surface_conditions"
              :key="surfaceCondition.surface"
              :cx="getSurfaceIndicatorPosition(surfaceCondition.surface).x"
              :cy="getSurfaceIndicatorPosition(surfaceCondition.surface).y"
              r="2"
              :fill="getSurfaceConditionColor(surfaceCondition.condition)"
              class="surface-indicator"
            />
          </g>

          <!-- Center circle -->
          <circle
            cx="30"
            cy="30"
            r="8"
            fill="transparent"
            stroke="#6B7280"
            stroke-width="1"
            opacity="0.7"
          />
        </g>
      </svg>

      <!-- Tooth number -->
      <div class="tooth-number">{{ tooth.tooth_number }}</div>

      <!-- Status indicators -->
      <div class="tooth-indicators">
        <div v-if="tooth.notes" class="indicator-note" title="Has notes">
          <font-awesome-icon icon="fa-solid fa-sticky-note" class="w-2 h-2" />
        </div>
        <div v-if="tooth.appointment_id" class="indicator-appointment" title="Linked to appointment">
          <font-awesome-icon icon="fa-solid fa-calendar" class="w-2 h-2" />
        </div>
      </div>
    </div>

    <!-- Tooth label -->
    <div class="tooth-label">
      <span class="tooth-name">{{ getToothName() }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  tooth: {
    type: Object,
    required: true
  },
  chartType: {
    type: String,
    required: true
  }
})

const conditionColors = {
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

// Standard tooth names for permanent teeth
const permanentToothNames = {
  // Upper right (Quadrant 1)
  '11': 'Central Incisor', '12': 'Lateral Incisor', '13': 'Canine', '14': 'First Premolar',
  '15': 'Second Premolar', '16': 'First Molar', '17': 'Second Molar', '18': 'Third Molar',
  // Upper left (Quadrant 2)
  '21': 'Central Incisor', '22': 'Lateral Incisor', '23': 'Canine', '24': 'First Premolar',
  '25': 'Second Premolar', '26': 'First Molar', '27': 'Second Molar', '28': 'Third Molar',
  // Lower left (Quadrant 3)
  '31': 'Central Incisor', '32': 'Lateral Incisor', '33': 'Canine', '34': 'First Premolar',
  '35': 'Second Premolar', '36': 'First Molar', '37': 'Second Molar', '38': 'Third Molar',
  // Lower right (Quadrant 4)
  '41': 'Central Incisor', '42': 'Lateral Incisor', '43': 'Canine', '44': 'First Premolar',
  '45': 'Second Premolar', '46': 'First Molar', '47': 'Second Molar', '48': 'Third Molar'
}

// Primary tooth names
const primaryToothNames = {
  'A': 'Central Incisor', 'B': 'Lateral Incisor', 'C': 'Canine', 'D': 'First Molar', 'E': 'Second Molar',
  'F': 'Second Molar', 'G': 'First Molar', 'H': 'Canine', 'I': 'Lateral Incisor', 'J': 'Central Incisor',
  'K': 'Central Incisor', 'L': 'Lateral Incisor', 'M': 'Canine', 'N': 'First Molar', 'O': 'Second Molar',
  'P': 'Second Molar', 'Q': 'First Molar', 'R': 'Canine', 'S': 'Lateral Incisor', 'T': 'Central Incisor'
}

const getToothColor = () => {
  // If there are surface conditions, use the most severe one
  if (props.tooth.surface_conditions && props.tooth.surface_conditions.length > 0) {
    const severityOrder = ['healthy', 'filled', 'decay', 'crowned', 'root_canal', 'extracted', 'implant', 'bridge', 'missing']
    let mostSevere = 'healthy'

    for (const surfaceCondition of props.tooth.surface_conditions) {
      const currentIndex = severityOrder.indexOf(surfaceCondition.condition)
      const severeIndex = severityOrder.indexOf(mostSevere)
      if (currentIndex > severeIndex) {
        mostSevere = surfaceCondition.condition
      }
    }

    return conditionColors[mostSevere] || '#ffffff'
  }

  // Otherwise use the overall tooth condition
  return conditionColors[props.tooth.condition] || '#ffffff'
}

const getToothName = () => {
  const toothNumber = props.tooth.tooth_number
  const isPrimary = props.chartType === 'primary'

  if (isPrimary) {
    return primaryToothNames[toothNumber] || 'Unknown'
  } else {
    return permanentToothNames[toothNumber] || 'Unknown'
  }
}

const getTooltipText = () => {
  const condition = props.tooth.condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
  let tooltip = `Tooth ${props.tooth.tooth_number}: ${getToothName()}\nCondition: ${condition}`

  if (props.tooth.surfaces && props.tooth.surfaces.length > 0) {
    const surfaceNames = props.tooth.surfaces.map(surface => {
      switch(surface) {
        case 'occlusal': return 'Top Surface (Occlusal)'
        case 'mesial': return 'Mesial Surface'
        case 'distal': return 'Distal Surface'
        case 'buccal': return 'Buccal Surface (Cheek Side)'
        case 'lingual': return 'Lingual Surface (Tongue Side)'
        case 'center': return 'Whole Tooth'
        default: return surface
      }
    })
    tooltip += `\nAffected surfaces: ${surfaceNames.join(', ')}`
  }

  if (props.tooth.surface_conditions && props.tooth.surface_conditions.length > 0) {
    tooltip += '\nSurface conditions:'
    props.tooth.surface_conditions.forEach(sc => {
      tooltip += `\n  ${sc.surface}: ${sc.condition.replace(/_/g, ' ')}`
    })
  }

  if (props.tooth.notes) {
    tooltip += `\nNotes: ${props.tooth.notes}`
  }

  if (props.tooth.last_updated) {
    const updated = new Date(props.tooth.last_updated).toLocaleDateString()
    tooltip += `\nLast updated: ${updated}`
  }

  tooltip += '\n\nThis is a historical snapshot - read-only view'

  return tooltip
}

const getSurfaceIndicatorPosition = (surface) => {
  const positions = {
    occlusal: { x: 37, y: 20 },  // Top pie slice center
    mesial: { x: 40, y: 37 },    // Right pie slice center (varies by quadrant)
    distal: { x: 20, y: 37 },    // Left pie slice center (varies by quadrant)
    buccal: { x: 30, y: 42 },    // Bottom pie slice center
    lingual: { x: 30, y: 42 },   // Bottom pie slice center
    center: { x: 30, y: 30 }     // Center of tooth
  }
  return positions[surface] || { x: 30, y: 30 }
}

const getSurfaceConditionColor = (condition) => {
  const colors = {
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
  return colors[condition] || '#666'
}
</script>

<style scoped>
@import "../styles/main.css";

.tooth-snapshot-component {
  @apply relative transition-all duration-200 ease-in-out opacity-90;
}

.tooth-visual {
  @apply w-12 h-12 flex items-center justify-center relative;
  @apply shadow-sm;
}

.tooth-svg {
  @apply w-full h-full;
}

.tooth-shape {
  @apply transition-all duration-200;
}

.tooth-number {
  @apply absolute top-0 left-0 text-xs font-bold text-gray-700 z-10;
  @apply bg-white/70 rounded px-1;
  font-size: 10px;
}

.tooth-indicators {
  @apply absolute bottom-0 right-0 flex gap-0.5 z-10;
}

.indicator-note,
.indicator-appointment {
  @apply text-xs bg-white/70 rounded p-0.5;
}

.tooth-label {
  @apply text-center mt-1;
}

.tooth-name {
  @apply text-xs text-gray-600 leading-tight;
  display: block;
  max-width: 48px;
  line-height: 1.1;
  font-size: 10px;
}

/* Surface dividers styling */
.surface-dividers {
  @apply pointer-events-none;
}

/* Special styling for different conditions */
.tooth-shape[fill="#FF6347"] { /* extracted */
  @apply opacity-75;
}

.tooth-shape[fill="#D3D3D3"] { /* missing */
  @apply opacity-60;
  stroke-dasharray: 2;
}

/* Surface indicator dots */
.surface-indicator {
  @apply transition-all duration-200;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .tooth-visual {
    @apply w-10 h-10;
  }

  .tooth-number {
    font-size: 9px;
  }

  .tooth-name {
    font-size: 9px;
    max-width: 40px;
  }
}
</style>