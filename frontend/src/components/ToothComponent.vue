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
        <!-- Tooth outline (simple circle like life guard donut) -->
        <circle
          cx="30"
          cy="30"
          r="18"
          :fill="getToothColor()"
          :stroke="getToothBorderColor()"
          stroke-width="2"
          class="tooth-shape"
        />

        <!-- Surface division lines (always visible for clarity) -->
        <g class="surface-dividers">
          <!-- Vertical center line (stops at center circle) -->
          <line x1="30" y1="12" x2="30" y2="22" stroke="#6B7280" stroke-width="1" />
          <line x1="30" y1="38" x2="30" y2="48" stroke="#6B7280" stroke-width="1" />
          <!-- Horizontal center line (stops at center circle) -->
          <line x1="12" y1="30" x2="22" y2="30" stroke="#6B7280" stroke-width="1" />
          <line x1="38" y1="30" x2="48" y2="30" stroke="#6B7280" stroke-width="1" />
        </g>

        <!-- Surface areas (clickable) - divided like a pie with center cutout -->
        <!-- Top surface (occlusal) -->
        <path
          d="M 30 12 A 18 18 0 0 1 48 30 L 38 30 A 8 8 0 0 0 30 22 Z"
          fill="transparent"
          class="surface-area surface-occlusal"
          :class="{ 'surface-selected': isSurfaceSelected('occlusal') }"
          @click.stop="handleSurfaceClick('occlusal')"
          :title="'Occlusal surface - ' + (isSurfaceSelected('occlusal') ? 'Selected' : 'Click to select')"
        />

        <!-- Right surface (mesial/distal depending on quadrant) -->
        <path
          d="M 48 30 A 18 18 0 0 1 30 48 L 30 38 A 8 8 0 0 0 38 30 Z"
          fill="transparent"
          class="surface-area surface-mesial"
          :class="{ 'surface-selected': isSurfaceSelected(getRightSurface()) }"
          @click.stop="handleSurfaceClick(getRightSurface())"
          :title="getRightSurface() + ' surface - ' + (isSurfaceSelected(getRightSurface()) ? 'Selected' : 'Click to select')"
        />

        <!-- Bottom surface (lingual/buccal depending on quadrant) -->
        <path
          d="M 30 48 A 18 18 0 0 1 12 30 L 22 30 A 8 8 0 0 0 30 38 Z"
          fill="transparent"
          class="surface-area surface-lingual"
          :class="{ 'surface-selected': isSurfaceSelected(getBottomSurface()) }"
          @click.stop="handleSurfaceClick(getBottomSurface())"
          :title="getBottomSurface() + ' surface - ' + (isSurfaceSelected(getBottomSurface()) ? 'Selected' : 'Click to select')"
        />

        <!-- Left surface (mesial/distal depending on quadrant) -->
        <path
          d="M 12 30 A 18 18 0 0 1 30 12 L 30 22 A 8 8 0 0 0 22 30 Z"
          fill="transparent"
          class="surface-area surface-distal"
          :class="{ 'surface-selected': isSurfaceSelected(getLeftSurface()) }"
          @click.stop="handleSurfaceClick(getLeftSurface())"
          :title="getLeftSurface() + ' surface - ' + (isSurfaceSelected(getLeftSurface()) ? 'Selected' : 'Click to select')"
        />

        <!-- Center area (whole tooth) -->
        <circle
          cx="30"
          cy="30"
          r="8"
          fill="transparent"
          stroke="#6B7280"
          stroke-width="1"
          class="surface-area surface-center"
          :class="{ 'surface-selected': isSurfaceSelected('center') }"
          @click.stop="handleSurfaceClick('center')"
          :title="'Whole tooth - ' + (isSurfaceSelected('center') ? 'Selected' : 'Click to select')"
        />

        <!-- Surface indicators (dots) -->
        <g v-if="tooth.surfaces && tooth.surfaces.length > 0">
          <circle
            v-for="surface in tooth.surfaces"
            :key="surface"
            :cx="getSurfaceIndicatorPosition(surface).x"
            :cy="getSurfaceIndicatorPosition(surface).y"
            r="2"
            :fill="getSurfaceIndicatorColor(surface)"
            class="surface-indicator"
          />
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
  // Upper right
  '1': 'Central Incisor', '2': 'Lateral Incisor', '3': 'Canine', '4': 'First Premolar',
  '5': 'Second Premolar', '6': 'First Molar', '7': 'Second Molar', '8': 'Third Molar',
  // Upper left  
  '9': 'Central Incisor', '10': 'Lateral Incisor', '11': 'Canine', '12': 'First Premolar',
  '13': 'Second Premolar', '14': 'First Molar', '15': 'Second Molar', '16': 'Third Molar',
  // Lower left
  '17': 'Third Molar', '18': 'Second Molar', '19': 'First Molar', '20': 'Second Premolar',
  '21': 'First Premolar', '22': 'Canine', '23': 'Lateral Incisor', '24': 'Central Incisor',
  // Lower right
  '25': 'Central Incisor', '26': 'Lateral Incisor', '27': 'Canine', '28': 'First Premolar',
  '29': 'Second Premolar', '30': 'First Molar', '31': 'Second Molar', '32': 'Third Molar'
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

// Surface mapping based on quadrant
const getRightSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'mesial' // primary teeth

  // Upper right and lower right quadrants use distal on the right
  if (toothNum <= 16 || (toothNum >= 25 && toothNum <= 32)) {
    return 'distal'
  }
  // Upper left and lower left quadrants use mesial on the right
  return 'mesial'
}

const getLeftSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'distal' // primary teeth

  // Upper right and lower right quadrants use mesial on the left
  if (toothNum <= 16 || (toothNum >= 25 && toothNum <= 32)) {
    return 'mesial'
  }
  // Upper left and lower left quadrants use distal on the left
  return 'distal'
}

const getBottomSurface = () => {
  const toothNum = parseInt(props.tooth.tooth_number)
  if (isNaN(toothNum)) return 'lingual' // primary teeth

  // Upper quadrants use lingual on bottom
  if (toothNum <= 16) {
    return 'lingual'
  }
  // Lower quadrants use buccal on bottom
  return 'buccal'
}

const isSurfaceSelected = (surface) => {
  return props.tooth.surfaces && props.tooth.surfaces.includes(surface)
}

const handleSurfaceClick = (surface) => {
  if (!props.editable) return

  const currentSurfaces = props.tooth.surfaces || []
  let newSurfaces

  if (surface === 'center') {
    // Selecting center toggles whole tooth selection
    newSurfaces = currentSurfaces.includes('center') ? [] : ['center']
  } else {
    // For specific surfaces, toggle that surface
    if (currentSurfaces.includes(surface)) {
      newSurfaces = currentSurfaces.filter(s => s !== surface)
    } else {
      // Remove center if selecting specific surface
      newSurfaces = [...currentSurfaces.filter(s => s !== 'center'), surface]
    }
  }

  const updatedTooth = {
    ...props.tooth,
    surfaces: newSurfaces
  }

  emit('update', updatedTooth)
}

const getSurfaceIndicatorPosition = (surface) => {
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