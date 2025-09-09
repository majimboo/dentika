<template>
  <div class="relative">
    <!-- Desktop: Show content with hover -->
    <div class="hidden sm:block" @mouseenter="showTooltip" @mouseleave="hideTooltip">
      <slot></slot>
    </div>
    
    <!-- Mobile: Show content with optional (?) in top right -->
    <div class="sm:hidden relative" @click="!showMobileIcon ? toggleTooltip : null">
      <slot></slot>
      <button 
        v-if="showMobileIcon"
        @click="toggleTooltip"
        class="absolute top-2 right-2 inline-flex items-center justify-center w-4 h-4 text-xs font-bold text-neutral-500 bg-neutral-100 border border-neutral-300 rounded-full hover:bg-neutral-200 transition-colors z-10"
        aria-label="Help"
      >
        ?
      </button>
    </div>
    
    <!-- Tooltip content -->
    <div 
      v-if="isVisible" 
      :class="tooltipClasses"
      class="absolute z-50 px-4 py-3 text-sm font-medium text-white bg-neutral-900 rounded-lg shadow-xl transition-opacity duration-200 max-w-xs w-max"
      :style="tooltipStyle"
      ref="tooltip"
    >
      {{ text }}
      <div class="tooltip-arrow absolute w-2 h-2 bg-neutral-900 transform rotate-45" :class="arrowClasses"></div>
    </div>
    
    <!-- Mobile overlay backdrop -->
    <div 
      v-if="isVisible && isMobile" 
      class="fixed inset-0 bg-black/25 z-40"
      @click="hideTooltip"
    ></div>
  </div>
</template>

<script>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'

export default {
  name: 'BaseTooltip',
  props: {
    text: {
      type: String,
      required: true
    },
    position: {
      type: String,
      default: 'top',
      validator: (value) => ['top', 'bottom', 'left', 'right'].includes(value)
    },
    delay: {
      type: Number,
      default: 500
    },
    showMobileIcon: {
      type: Boolean,
      default: true
    }
  },
  setup(props) {
    const isVisible = ref(false)
    let showTimeout = null
    let hideTimeout = null

    const isMobile = ref(false)
    
    const updateIsMobile = () => {
      isMobile.value = window.innerWidth < 640
    }

    const tooltipClasses = computed(() => {
      // On mobile, center the tooltip in viewport
      if (isMobile.value) {
        return 'fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2'
      }
      
      // Desktop positioning with better spacing
      const positions = {
        top: 'bottom-full left-1/2 transform -translate-x-1/2 mb-3',
        bottom: 'top-full left-1/2 transform -translate-x-1/2 mt-3', 
        left: 'right-full top-1/2 transform -translate-y-1/2 mr-3',
        right: 'left-full top-1/2 transform -translate-y-1/2 ml-3'
      }
      return positions[props.position]
    })

    const arrowClasses = computed(() => {
      // Hide arrow on mobile since tooltip is centered
      if (isMobile.value) {
        return 'hidden'
      }
      
      const arrows = {
        top: 'top-full left-1/2 transform -translate-x-1/2 -mt-1',
        bottom: 'bottom-full left-1/2 transform -translate-x-1/2 -mb-1',
        left: 'left-full top-1/2 transform -translate-y-1/2 -ml-1',
        right: 'right-full top-1/2 transform -translate-y-1/2 -mr-1'
      }
      return arrows[props.position]
    })

    const tooltipStyle = computed(() => {
      return isVisible.value ? { opacity: '1' } : { opacity: '0' }
    })

    const showTooltip = () => {
      if (hideTimeout) {
        clearTimeout(hideTimeout)
        hideTimeout = null
      }
      
      showTimeout = setTimeout(() => {
        isVisible.value = true
        nextTick(() => {
          const tooltip = document.querySelector('.tooltip-arrow').parentElement
          if (tooltip) {
            tooltip.style.opacity = '1'
          }
        })
      }, props.delay)
    }

    const hideTooltip = () => {
      if (showTimeout) {
        clearTimeout(showTimeout)
        showTimeout = null
      }
      
      const delay = isMobile.value ? 0 : 100
      hideTimeout = setTimeout(() => {
        isVisible.value = false
      }, delay)
    }

    const toggleTooltip = () => {
      if (isVisible.value) {
        hideTooltip()
      } else {
        showTooltip()
      }
    }

    // Initialize and handle resize
    onMounted(() => {
      updateIsMobile()
      window.addEventListener('resize', updateIsMobile)
    })

    onUnmounted(() => {
      window.removeEventListener('resize', updateIsMobile)
    })

    return {
      isVisible,
      isMobile,
      tooltipClasses,
      arrowClasses,
      tooltipStyle,
      showTooltip,
      hideTooltip,
      toggleTooltip
    }
  }
}
</script>

<style scoped>
.tooltip-arrow {
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}
</style>