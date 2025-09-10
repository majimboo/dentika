<template>
  <div :class="containerClasses">
    <div v-if="type === 'spinner'" :class="spinnerClasses">
      <font-awesome-icon icon="fa-solid fa-spinner" :class="iconClasses" class="animate-spin" />
    </div>

    <div v-else-if="type === 'dots'" :class="dotsClasses">
      <div class="dot bg-current animate-pulse" :style="{ animationDelay: '0ms' }"></div>
      <div class="dot bg-current animate-pulse" :style="{ animationDelay: '150ms' }"></div>
      <div class="dot bg-current animate-pulse" :style="{ animationDelay: '300ms' }"></div>
    </div>

    <div v-else-if="type === 'skeleton'" :class="skeletonClasses">
      <div class="animate-pulse">
        <div class="h-4 bg-neutral-200 rounded mb-2"></div>
        <div class="h-3 bg-neutral-200 rounded w-3/4"></div>
      </div>
    </div>

    <div v-else-if="type === 'bar'" :class="barClasses">
      <div class="bg-primary-200 rounded-full h-2 overflow-hidden">
        <div class="bg-primary-600 h-full rounded-full loading-bar"></div>
      </div>
    </div>

    <div v-if="text && showText" :class="textClasses">
      {{ text }}
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'BaseLoading',
  props: {
    type: {
      type: String,
      default: 'spinner',
      validator: (value) => ['spinner', 'dots', 'skeleton', 'bar'].includes(value)
    },
    size: {
      type: String,
      default: 'medium',
      validator: (value) => ['small', 'medium', 'large'].includes(value)
    },
    color: {
      type: String,
      default: 'primary',
      validator: (value) => ['primary', 'secondary', 'neutral', 'white'].includes(value)
    },
    text: {
      type: String,
      default: ''
    },
    showText: {
      type: Boolean,
      default: true
    },
    overlay: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const containerClasses = computed(() => {
      let classes = 'flex flex-col items-center justify-center'
      
      if (props.overlay) {
        classes += ' fixed inset-0 bg-white/80 backdrop-blur-sm z-50'
      }
      
      return classes
    })

    const spinnerClasses = computed(() => {
      const colorClasses = {
        primary: 'text-primary-600',
        secondary: 'text-secondary-600',
        neutral: 'text-neutral-600',
        white: 'text-white'
      }
      
      return `inline-block ${colorClasses[props.color]}`
    })

    const iconClasses = computed(() => {
      const sizeClasses = {
        small: 'w-5 h-5',
        medium: 'w-8 h-8',
        large: 'w-12 h-12'
      }
      
      return sizeClasses[props.size]
    })

    const dotsClasses = computed(() => {
      const colorClasses = {
        primary: 'text-primary-600',
        secondary: 'text-secondary-600',
        neutral: 'text-neutral-600',
        white: 'text-white'
      }
      
      const sizeClasses = {
        small: 'space-x-1',
        medium: 'space-x-2',
        large: 'space-x-3'
      }
      
      return `flex items-center ${colorClasses[props.color]} ${sizeClasses[props.size]}`
    })

    const skeletonClasses = computed(() => {
      const sizeClasses = {
        small: 'w-32',
        medium: 'w-48',
        large: 'w-64'
      }
      
      return `animate-pulse ${sizeClasses[props.size]}`
    })

    const barClasses = computed(() => {
      const sizeClasses = {
        small: 'w-32',
        medium: 'w-48',
        large: 'w-64'
      }
      
      return sizeClasses[props.size]
    })

    const textClasses = computed(() => {
      const colorClasses = {
        primary: 'text-primary-700',
        secondary: 'text-secondary-700',
        neutral: 'text-neutral-700',
        white: 'text-white'
      }
      
      const sizeClasses = {
        small: 'text-sm mt-2',
        medium: 'text-base mt-3',
        large: 'text-lg mt-4'
      }
      
      return `font-medium ${colorClasses[props.color]} ${sizeClasses[props.size]}`
    })

    return {
      containerClasses,
      spinnerClasses,
      iconClasses,
      dotsClasses,
      skeletonClasses,
      barClasses,
      textClasses
    }
  }
}
</script>

<style scoped>
.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.loading-bar {
  animation: loading-progress 2s ease-in-out infinite;
}

@keyframes loading-progress {
  0% {
    width: 0%;
  }
  50% {
    width: 60%;
  }
  100% {
    width: 0%;
  }
}
</style>