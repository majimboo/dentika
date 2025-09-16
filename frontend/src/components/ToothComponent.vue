<template>
  <div
    class="tooth-component"
    :class="[
      'tooth-' + tooth.tooth_number,
      {
        'tooth-selected': selected,
        'tooth-editable': editable,
        'tooth-has-notes': tooth.notes
      }
    ]"
    @click="handleClick"
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
          <!-- Tooth outline (simple circle like life guard donut) -->
          <circle
            cx="30"
            cy="30"
            r="18"
            fill="white"
            :stroke="getToothBorderColor()"
            stroke-width="3"
            class="tooth-shape"
          />

          <!-- Surface division lines (always visible for clarity) -->
          <g class="surface-dividers">
            <!-- Vertical center line (stops at center circle) -->
            <line x1="30" y1="12" x2="30" y2="22" stroke="#6B7280" stroke-width="2" />
            <line x1="30" y1="38" x2="30" y2="48" stroke="#6B7280" stroke-width="2" />
            <!-- Horizontal center line (stops at center circle) -->
            <line x1="12" y1="30" x2="22" y2="30" stroke="#6B7280" stroke-width="2" />
            <line x1="38" y1="30" x2="48" y2="30" stroke="#6B7280" stroke-width="2" />
          </g>

          <!-- Surface areas (clickable) - divided like a pie with center cutout -->
          <!-- Top surface (occlusal) -->
          <path
            d="M 30 12 A 18 18 0 0 1 48 30 L 38 30 A 8 8 0 0 0 30 22 Z"
            fill="transparent"
            class="surface-area surface-occlusal"
            :class="{ 'surface-selected': isSurfaceSelected('occlusal') }"
            @click="handleSurfaceClick('occlusal')"
            :title="'Occlusal surface - ' + (props.editable ? 'Click to edit' : 'View details')"
          />

          <!-- Right surface (mesial/distal depending on quadrant) -->
          <path
            d="M 48 30 A 18 18 0 0 1 30 48 L 30 38 A 8 8 0 0 0 38 30 Z"
            fill="transparent"
            class="surface-area"
            :class="[
              'surface-' + getRightSurface(),
              { 'surface-selected': isSurfaceSelected(getRightSurface()) }
            ]"
            @click="handleSurfaceClick(getRightSurface())"
            :title="getRightSurface() + ' surface - ' + (props.editable ? 'Click to edit' : 'View details')"
          />

          <!-- Bottom surface (lingual/buccal depending on quadrant) -->
          <path
            d="M 30 48 A 18 18 0 0 1 12 30 L 22 30 A 8 8 0 0 0 30 38 Z"
            fill="transparent"
            class="surface-area"
            :class="[
              'surface-' + getBottomSurface(),
              { 'surface-selected': isSurfaceSelected(getBottomSurface()) }
            ]"
            @click="handleSurfaceClick(getBottomSurface())"
            :title="getBottomSurface() + ' surface - ' + (props.editable ? 'Click to edit' : 'View details')"
          />

          <!-- Left surface (mesial/distal depending on quadrant) -->
          <path
            d="M 12 30 A 18 18 0 0 1 30 12 L 30 22 A 8 8 0 0 0 22 30 Z"
            fill="transparent"
            class="surface-area"
            :class="[
              'surface-' + getLeftSurface(),
              { 'surface-selected': isSurfaceSelected(getLeftSurface()) }
            ]"
            @click="handleSurfaceClick(getLeftSurface())"
            :title="getLeftSurface() + ' surface - ' + (props.editable ? 'Click to edit' : 'View details')"
          />

          <!-- Center area (whole tooth) -->
          <circle
            cx="30"
            cy="30"
            r="8"
            fill="transparent"
            stroke="#6B7280"
            stroke-width="2"
            class="surface-area surface-center"
            :class="{ 'surface-selected': isSurfaceSelected('center') }"
            @click="handleSurfaceClick('center')"
            :title="'Whole tooth - ' + (props.editable ? 'Click to edit' : 'View details')"
          />

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
        </g>
      </svg>

      <!-- Tooth number -->
      <div class="tooth-number">{{ tooth.tooth_number }}</div>

      <!-- Status indicators -->
      <div class="tooth-indicators">
        <div v-if="tooth.notes" class="indicator-note" title="Has notes">
          <font-awesome-icon icon="fa-solid fa-sticky-note" class="w-3 h-3" />
        </div>
        <div v-if="isRecentlyUpdated" class="indicator-recent" title="Recently updated">
          <font-awesome-icon icon="fa-solid fa-sync" class="w-3 h-3" />
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
  selected: {
    type: Boolean,
    default: false
  },
  editable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click', 'update'])

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

const conditionIndicators = {
  healthy: '✓',
  decay: '⚠',
  filled: 'F',
  crowned: 'C',
  root_canal: 'RC',
  extracted: '✗',
  implant: 'I',
  bridge: 'B',
  missing: '-'
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
  '31': 'Third Molar', '32': 'Second Molar', '33': 'First Molar', '34': 'Second Premolar',
  '35': 'First Premolar', '36': 'Canine', '37': 'Lateral Incisor', '38': 'Central Incisor',
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

const isRecentlyUpdated = computed(() => {
  if (!props.tooth.last_updated) return false
  const updatedTime = new Date(props.tooth.last_updated)
  const now = new Date()
  const diffHours = (now - updatedTime) / (1000 * 60 * 60)
  return diffHours < 24 // Updated within last 24 hours
})

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

const getConditionIndicator = () => {
  return conditionIndicators[props.tooth.condition] || '?'
}

const getToothName = () => {
  const toothNumber = props.tooth.tooth_number
  const isPrimary = isNaN(toothNumber)
  
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
  } else {
    tooltip += `\nNo specific surfaces selected`
  }

  if (props.editable) {
    tooltip += `\n\nClick on tooth areas to mark specific surface issues`
  }

  if (props.tooth.notes) {
    tooltip += `\nNotes: ${props.tooth.notes}`
  }

  if (props.tooth.last_updated) {
    const updated = new Date(props.tooth.last_updated).toLocaleDateString()
    tooltip += `\nLast updated: ${updated}`
  }

  return tooltip
}

// Surface mapping based on quadrant (FDI World Dental Federation notation)
const getRightSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'mesial' // primary teeth

  const quadrant = Math.floor(toothNum / 10) // First digit determines quadrant

  // Quadrant 1 (Upper right): distal on right
  if (quadrant === 1) {
    return 'distal'
  }
  // Quadrant 2 (Upper left): mesial on right
  if (quadrant === 2) {
    return 'mesial'
  }
  // Quadrant 3 (Lower left): mesial on right
  if (quadrant === 3) {
    return 'mesial'
  }
  // Quadrant 4 (Lower right): distal on right
  if (quadrant === 4) {
    return 'distal'
  }

  return 'mesial' // fallback
}

const getLeftSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'distal' // primary teeth

  const quadrant = Math.floor(toothNum / 10) // First digit determines quadrant

  // Quadrant 1 (Upper right): mesial on left
  if (quadrant === 1) {
    return 'mesial'
  }
  // Quadrant 2 (Upper left): distal on left
  if (quadrant === 2) {
    return 'distal'
  }
  // Quadrant 3 (Lower left): distal on left
  if (quadrant === 3) {
    return 'distal'
  }
  // Quadrant 4 (Lower right): mesial on left
  if (quadrant === 4) {
    return 'mesial'
  }

  return 'distal' // fallback
}

const getBottomSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'lingual' // primary teeth

  const quadrant = Math.floor(toothNum / 10) // First digit determines quadrant

  // Upper quadrants (1-2) use lingual on bottom
  if (quadrant === 1 || quadrant === 2) {
    return 'lingual'
  }
  // Lower quadrants (3-4) use buccal on bottom
  if (quadrant === 3 || quadrant === 4) {
    return 'buccal'
  }

  return 'lingual' // fallback
}

const isSurfaceSelected = (surface) => {
  return props.tooth.surfaces && props.tooth.surfaces.includes(surface)
}

const handleSurfaceClick = (surface) => {
  if (!props.editable) {
    // If not editable, just trigger the main tooth click
    handleClick()
    return
  }

  // In the new workflow, surface clicks should open the edit modal
  // instead of just selecting surfaces
  handleClick()
}

const getSurfaceIndicatorPosition = (surface) => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) {
    // Primary teeth - use standard positions
    const positions = {
      occlusal: { x: 37, y: 20 },  // Top pie slice center
      mesial: { x: 40, y: 37 },    // Right pie slice center
      distal: { x: 20, y: 37 },    // Left pie slice center
      buccal: { x: 30, y: 42 },    // Bottom pie slice center
      lingual: { x: 30, y: 42 },   // Bottom pie slice center
      center: { x: 30, y: 30 }     // Center of tooth
    }
    return positions[surface] || { x: 30, y: 30 }
  }

  // Permanent teeth - adjust positions based on quadrant
  let positions = {
    occlusal: { x: 37, y: 20 },  // Top pie slice center (always same)
    center: { x: 30, y: 30 }     // Center of tooth (always same)
  }

  // Determine which surface is on the right and left based on quadrant
  const rightSurface = getRightSurface()
  const leftSurface = getLeftSurface()
  const bottomSurface = getBottomSurface()

  // Assign positions based on actual surface names
  if (rightSurface === 'mesial') {
    positions.mesial = { x: 40, y: 37 }  // Right side
    positions.distal = { x: 20, y: 37 }  // Left side
  } else {
    positions.distal = { x: 40, y: 37 }  // Right side
    positions.mesial = { x: 20, y: 37 }  // Left side
  }

  // Bottom surface (buccal or lingual)
  positions.buccal = { x: 30, y: 42 }
  positions.lingual = { x: 30, y: 42 }

  return positions[surface] || { x: 30, y: 30 }
}

const getSurfaceIndicatorColor = (surface) => {
  const colors = {
    occlusal: '#FFD700', // gold
    mesial: '#FF6B6B',   // red
    distal: '#4ECDC4',   // teal
    buccal: '#45B7D1',   // blue
    lingual: '#96CEB4',  // green
    center: '#FFEAA7'    // yellow
  }
  return colors[surface] || '#666'
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

const getToothBorderColor = () => {
  if (props.selected) return '#2563EB' // blue-600
  if (props.editable) return '#6B7280' // gray-500
  return '#D1D5DB' // gray-300
}

const handleClick = () => {
  emit('click', props.tooth.tooth_number)
}
</script>

<style scoped>
@import "../styles/main.css";

.tooth-component {
  @apply relative cursor-pointer transition-all duration-200 ease-in-out;
}

.tooth-visual {
  @apply w-16 h-16 flex items-center justify-center relative;
  @apply shadow-sm hover:shadow-md transition-shadow;
}

.tooth-svg {
  @apply w-full h-full;
}

.tooth-shape {
  @apply transition-all duration-200;
}

.tooth-editable .tooth-shape {
  @apply hover:stroke-blue-400;
}

.tooth-selected .tooth-shape {
  @apply stroke-blue-600 shadow-lg;
}

.tooth-selected .tooth-visual {
  @apply ring-2 ring-blue-300 rounded-lg;
}

/* Surface area styling */
.surface-area {
  @apply cursor-pointer transition-all duration-150;
}

.surface-area:hover {
  @apply stroke-gray-500 stroke-2;
}

.tooth-editable .surface-area:hover {
  @apply fill-current opacity-35 stroke-gray-700 stroke-3;
}

.surface-selected {
  @apply fill-current opacity-55 stroke-gray-800 stroke-3;
}

/* Center circle specific styling */
.surface-center {
  @apply stroke-gray-500;
}

.surface-center:hover {
  @apply stroke-gray-700 stroke-2;
}

.tooth-editable .surface-center:hover {
  @apply stroke-gray-800 stroke-3;
}

.surface-center.surface-selected {
  @apply stroke-gray-900 stroke-3;
}

/* Surface colors for selection */
.surface-occlusal.surface-selected { fill: #FFD700 !important; }
.surface-mesial.surface-selected { fill: #FF6B6B !important; }
.surface-distal.surface-selected { fill: #4ECDC4 !important; }
.surface-buccal.surface-selected { fill: #45B7D1 !important; }
.surface-lingual.surface-selected { fill: #96CEB4 !important; }
.surface-center.surface-selected { fill: #FFEAA7 !important; }

/* Hover effects for surface areas */
.tooth-editable .surface-occlusal:hover { fill: rgba(255, 215, 0, 0.45) !important; }
.tooth-editable .surface-mesial:hover { fill: rgba(255, 107, 107, 0.45) !important; }
.tooth-editable .surface-distal:hover { fill: rgba(78, 205, 196, 0.45) !important; }
.tooth-editable .surface-buccal:hover { fill: rgba(69, 183, 209, 0.45) !important; }
.tooth-editable .surface-lingual:hover { fill: rgba(150, 206, 180, 0.45) !important; }
.tooth-editable .surface-center:hover { fill: rgba(255, 234, 167, 0.45) !important; }

/* Surface dividers styling - always visible for clarity */
.surface-dividers {
  @apply pointer-events-none;
}

.surface-dividers line {
  @apply opacity-70 stroke-gray-600;
}

.tooth-number {
  @apply absolute top-1 left-1 text-xs font-bold text-gray-700 z-10;
  @apply bg-white/50 rounded px-1;
}

.tooth-indicators {
  @apply absolute bottom-0 right-0 flex gap-0.5 z-10;
}

.indicator-note,
.indicator-recent {
  @apply text-xs bg-white/50 rounded p-0.5;
}

.tooth-label {
  @apply text-center mt-1;
}

.tooth-name {
  @apply text-xs text-gray-600 leading-tight;
  display: block;
  max-width: 64px;
  line-height: 1.2;
}

/* Special styling for different conditions */
.tooth-shape[fill="#FF6347"] { /* extracted */
  @apply opacity-75;
}

.tooth-shape[fill="#D3D3D3"] { /* missing */
  @apply opacity-60;
  stroke-dasharray: 4;
}

/* Hover effects */
.tooth-component:hover .tooth-visual {
  @apply transform scale-105;
}

.tooth-selected:hover .tooth-visual {
  @apply transform scale-110;
}

/* Surface indicator dots */
.surface-indicator {
  @apply transition-all duration-200;
}

.surface-indicator:hover {
  r: 3;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .tooth-visual {
    @apply w-14 h-14;
  }

  .tooth-number {
    @apply text-xs;
  }

  .tooth-name {
    @apply text-xs;
    max-width: 56px;
  }
}
</style>