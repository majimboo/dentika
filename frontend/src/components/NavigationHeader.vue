<template>
  <div class="navigation-header bg-white border-b border-gray-200 sticky top-0 z-40 lg:ml-64">
    <div class="flex items-center h-14 px-4 lg:px-6">
      <!-- Left Actions Slot (hamburger menu, etc.) -->
      <div v-if="$slots.leftActions" class="flex items-center space-x-2 flex-shrink-0 mr-3">
        <slot name="leftActions"></slot>
      </div>
      
      <!-- Back Button - Only shows if not at home level -->
      <button
        v-if="showBackButton"
        @click="goBack"
        class="mr-3 p-2 -ml-2 rounded-lg hover:bg-gray-100 transition-colors flex-shrink-0"
        :aria-label="backButtonLabel"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
        </svg>
      </button>

      <!-- Page Title -->
      <div class="flex-1 min-w-0">
        <h1 class="text-lg font-semibold text-gray-900 truncate">
          {{ currentPageTitle }}
        </h1>
        <div v-if="breadcrumb" class="text-xs text-gray-500 truncate">
          {{ breadcrumb }}
        </div>
      </div>

      <!-- Right Actions Slot -->
      <div v-if="$slots.actions" class="flex items-center space-x-2 flex-shrink-0 ml-3">
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// Computed properties for navigation state
const currentLevel = computed(() => {
  return route.meta?.level || 0
})

const showBackButton = computed(() => {
  return currentLevel.value > 0
})

const currentPageTitle = computed(() => {
  return route.meta?.title || route.name || 'Page'
})

const backButtonLabel = computed(() => {
  const parentRoute = getParentRouteName()
  return parentRoute ? `Back to ${parentRoute}` : 'Back'
})

const breadcrumb = computed(() => {
  if (currentLevel.value <= 1) return null
  
  const parentRoute = getParentRouteName()
  return parentRoute ? `${parentRoute}` : null
})

// Get the parent route name for breadcrumb
const getParentRouteName = () => {
  const parentName = route.meta?.parent
  if (!parentName) return null
  
  // Convert route names to display names
  const displayNames = {
    'Dashboard': 'Dashboard',
    'PatientList': 'Patients',
    'AppointmentCalendar': 'Appointments', 
    'ProcedureManagement': 'Procedures',
    'UserList': 'Users',
    'ClinicManagement': 'Clinics'
  }
  
  return displayNames[parentName] || parentName
}

// Navigation logic
const goBack = () => {
  const parentRouteName = route.meta?.parent
  
  if (parentRouteName) {
    // Navigate to the specific parent route
    const routeMap = {
      'Dashboard': '/',
      'PatientList': '/patients',
      'AppointmentCalendar': '/appointments',
      'ProcedureManagement': '/procedures',
      'UserList': '/users',
      'ClinicManagement': '/clinics'
    }
    
    const parentPath = routeMap[parentRouteName]
    if (parentPath) {
      router.push(parentPath)
    } else {
      // Fallback to browser back
      router.go(-1)
    }
  } else {
    // Fallback to browser back for unknown structure
    router.go(-1)
  }
}

// Optional: Handle hardware back button on mobile
const handleHardwareBack = (event) => {
  if (showBackButton.value) {
    event.preventDefault()
    goBack()
  }
}

// Listen for Android back button
if (typeof window !== 'undefined' && window.history) {
  window.addEventListener('popstate', handleHardwareBack)
}
</script>

<style scoped>
.navigation-header {
  /* Ensure proper layering above content */
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

/* iOS-style back button hover effect */
@media (hover: hover) {
  button:hover {
    background-color: rgba(0, 0, 0, 0.05);
  }
}

/* Active state for touch devices */
@media (hover: none) {
  button:active {
    background-color: rgba(0, 0, 0, 0.1);
  }
}
</style>