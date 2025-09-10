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
    <div 
      class="tooth-visual"
      :style="{ backgroundColor: getToothColor() }"
    >
      <div class="tooth-number">{{ tooth.tooth_number }}</div>
      <div class="tooth-condition-indicator">{{ getConditionIndicator() }}</div>
      
      <!-- Surface indicators -->
      <div v-if="tooth.surfaces && tooth.surfaces.length > 0" class="tooth-surfaces">
        <div 
          v-for="surface in tooth.surfaces" 
          :key="surface"
          class="surface-dot"
          :class="'surface-' + surface"
        ></div>
      </div>
      
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
    tooltip += `\nSurfaces: ${props.tooth.surfaces.join(', ')}`
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
  @apply w-12 h-12 border-2 border-gray-400 rounded-lg flex flex-col items-center justify-center relative;
  @apply shadow-sm hover:shadow-md transition-shadow;
}

.tooth-editable .tooth-visual {
  @apply hover:border-blue-400 hover:shadow-blue-100;
}

.tooth-selected .tooth-visual {
  @apply border-blue-600 shadow-lg ring-2 ring-blue-300;
}

.tooth-number {
  @apply text-xs font-bold text-gray-700;
}

.tooth-condition-indicator {
  @apply text-xs font-medium mt-1;
}

.tooth-surfaces {
  @apply absolute top-1 right-1 flex flex-wrap gap-0.5;
}

.surface-dot {
  @apply w-1.5 h-1.5 bg-red-500 rounded-full;
}

.surface-mesial { @apply bg-blue-500; }
.surface-distal { @apply bg-green-500; }
.surface-occlusal { @apply bg-yellow-500; }
.surface-buccal { @apply bg-purple-500; }
.surface-lingual { @apply bg-pink-500; }

.tooth-indicators {
  @apply absolute bottom-0 right-0 flex gap-0.5;
}

.indicator-note,
.indicator-recent {
  @apply text-xs;
}

.tooth-label {
  @apply text-center mt-1;
}

.tooth-name {
  @apply text-xs text-gray-600 leading-tight;
  display: block;
  max-width: 48px;
  line-height: 1.2;
}

/* Special styling for different conditions */
.tooth-visual[style*="#FF6347"] { /* extracted */
  @apply opacity-75;
}

.tooth-visual[style*="#D3D3D3"] { /* missing */
  @apply border-dashed opacity-60;
}

/* Hover effects */
.tooth-component:hover .tooth-visual {
  @apply transform scale-105;
}

.tooth-selected:hover .tooth-visual {
  @apply transform scale-110;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .tooth-visual {
    @apply w-10 h-10;
  }
  
  .tooth-number,
  .tooth-condition-indicator {
    @apply text-xs;
  }
  
  .tooth-name {
    @apply text-xs;
    max-width: 40px;
  }
}
</style>