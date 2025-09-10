<template>
  <div class="patient-dental-chart-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Patient Dental Chart</h1>
        <p v-if="currentPatient" class="text-gray-600">
          {{ currentPatient.first_name }} {{ currentPatient.last_name }} - 
          <span class="text-sm">ID: {{ currentPatient.id }}</span>
        </p>
        <p v-else class="text-gray-600">Loading patient information...</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <!-- Patient Selector -->
        <div class="patient-selector">
          <label for="patient-select" class="block text-sm font-medium text-gray-700 mb-1">
            Select Patient:
          </label>
          <select
            id="patient-select"
            v-model="selectedPatientId"
            @change="loadPatientChart"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Choose a patient...</option>
            <option 
              v-for="patient in patientOptions" 
              :key="patient.id" 
              :value="patient.id"
            >
              {{ patient.first_name }} {{ patient.last_name }}
            </option>
          </select>
        </div>

        <!-- Quick Actions -->
        <div v-if="currentPatient" class="quick-actions flex items-center space-x-2">
          <router-link
            :to="`/patients/${currentPatient.id}`"
            class="btn btn-secondary text-sm"
          >
            <font-awesome-icon icon="fa-solid fa-user" class="w-4 h-4 mr-1" />
            View Patient
          </router-link>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state flex items-center justify-center h-64">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p class="text-gray-600">Loading dental chart...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state bg-red-50 border border-red-200 rounded-lg p-6 text-center">
      <div class="text-red-600 mb-4">
        <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-8 h-8" />
      </div>
      <h3 class="text-lg font-medium text-red-900 mb-2">Error Loading Dental Chart</h3>
      <p class="text-red-700 mb-4">{{ error }}</p>
      <button
        @click="retryLoad"
        class="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
      >
        Try Again
      </button>
    </div>

    <!-- No Patient Selected -->
    <div v-else-if="!currentPatient" class="no-patient-state bg-gray-50 border border-gray-200 rounded-lg p-12 text-center">
      <div class="text-gray-400 mb-4">
        <font-awesome-icon icon="fa-solid fa-tooth" class="w-16 h-16" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Select a Patient</h3>
      <p class="text-gray-600">Choose a patient from the dropdown above to view their dental chart.</p>
    </div>

    <!-- Dental Chart -->
    <div v-else-if="currentPatient" class="dental-chart-container">
      <!-- Patient Info Summary -->
      <div class="patient-summary bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
        <div class="flex items-center justify-between">
          <div class="patient-details">
            <h3 class="text-lg font-semibold text-gray-900">
              {{ currentPatient.first_name }} {{ currentPatient.last_name }}
            </h3>
            <div class="text-sm text-gray-600 space-x-4">
              <span v-if="currentPatient.date_of_birth">
                Age: {{ calculateAge(currentPatient.date_of_birth) }}
              </span>
              <span v-if="currentPatient.phone">
                Phone: {{ currentPatient.phone }}
              </span>
              <span v-if="currentPatient.email">
                Email: {{ currentPatient.email }}
              </span>
            </div>
          </div>
          
          <!-- Chart Stats -->
          <div class="chart-stats flex items-center space-x-6 text-sm">
            <div class="stat-item text-center">
              <div class="text-lg font-bold text-blue-600">{{ healthyTeethCount }}</div>
              <div class="text-gray-600">Healthy</div>
            </div>
            <div class="stat-item text-center">
              <div class="text-lg font-bold text-yellow-600">{{ treatedTeethCount }}</div>
              <div class="text-gray-600">Treated</div>
            </div>
            <div class="stat-item text-center">
              <div class="text-lg font-bold text-red-600">{{ problemTeethCount }}</div>
              <div class="text-gray-600">Issues</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Dental Chart Component -->
      <DentalChart
        :patientId="currentPatient.id"
        :readonly="readonly"
        @toothUpdated="handleToothUpdated"
        @chartUpdated="handleChartUpdated"
      />

      <!-- Additional Information -->
      <div class="additional-info mt-6 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Recent Updates -->
        <div class="recent-updates bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <h4 class="text-lg font-semibold text-gray-900 mb-3">Recent Updates</h4>
          <div v-if="recentUpdates.length > 0" class="space-y-3">
            <div 
              v-for="update in recentUpdates" 
              :key="update.id"
              class="update-item flex items-start space-x-3 p-3 bg-gray-50 rounded-lg"
            >
              <div class="update-icon">
                <font-awesome-icon icon="fa-solid fa-clock" class="w-4 h-4 text-gray-500 mt-1" />
              </div>
              <div class="update-content flex-1">
                <div class="update-header flex items-center justify-between">
                  <span class="font-medium text-gray-900">Tooth {{ update.tooth_number }}</span>
                  <span class="text-xs text-gray-500">{{ formatRelativeTime(update.updated_at) }}</span>
                </div>
                <div class="text-sm text-gray-600">
                  Updated to: {{ formatCondition(update.condition) }}
                </div>
                <div v-if="update.notes" class="text-xs text-gray-500 mt-1">
                  {{ update.notes }}
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-gray-500 text-sm">
            No recent updates to display.
          </div>
        </div>

        <!-- Treatment Recommendations -->
        <div class="treatment-recommendations bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <h4 class="text-lg font-semibold text-gray-900 mb-3">Treatment Recommendations</h4>
          <div v-if="treatmentRecommendations.length > 0" class="space-y-3">
            <div 
              v-for="recommendation in treatmentRecommendations" 
              :key="recommendation.id"
              class="recommendation-item p-3 border-l-4 border-yellow-400 bg-yellow-50"
            >
              <div class="font-medium text-gray-900">Tooth {{ recommendation.tooth_number }}</div>
              <div class="text-sm text-gray-700 mt-1">{{ recommendation.recommendation }}</div>
              <div class="text-xs text-gray-500 mt-2">
                Priority: {{ recommendation.priority }}
              </div>
            </div>
          </div>
          <div v-else class="text-gray-500 text-sm">
            No treatment recommendations at this time.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { usePatientStore } from '../stores/patient'
import { useDentalRecordStore } from '../stores/dentalRecord'
import { useNotificationStore } from '../stores/notification'
import DentalChart from '../components/DentalChart.vue'

const route = useRoute()
const authStore = useAuthStore()
const patientStore = usePatientStore()
const dentalRecordStore = useDentalRecordStore()
const notificationStore = useNotificationStore()

const selectedPatientId = ref('')
const loading = ref(false)
const error = ref('')

const readonly = computed(() => {
  return !authStore.isDoctor && !authStore.isSuperAdmin
})

const currentPatient = computed(() => {
  return patientStore.patients.find(p => p.id === parseInt(selectedPatientId.value))
})

const patientOptions = computed(() => {
  return patientStore.patients || []
})

const teethData = computed(() => {
  return dentalRecordStore.currentRecord?.teeth_data || []
})

const healthyTeethCount = computed(() => {
  return teethData.value.filter(tooth => tooth.condition === 'healthy').length
})

const treatedTeethCount = computed(() => {
  return teethData.value.filter(tooth => 
    ['filled', 'crowned', 'root_canal', 'implant', 'bridge'].includes(tooth.condition)
  ).length
})

const problemTeethCount = computed(() => {
  return teethData.value.filter(tooth => 
    ['decay', 'extracted', 'missing'].includes(tooth.condition)
  ).length
})

const recentUpdates = computed(() => {
  return teethData.value
    .filter(tooth => tooth.last_updated)
    .sort((a, b) => new Date(b.last_updated) - new Date(a.last_updated))
    .slice(0, 5)
})

const treatmentRecommendations = computed(() => {
  const recommendations = []
  
  teethData.value.forEach(tooth => {
    if (tooth.condition === 'decay') {
      recommendations.push({
        id: `${tooth.tooth_number}-decay`,
        tooth_number: tooth.tooth_number,
        recommendation: 'Requires filling or restoration',
        priority: 'High'
      })
    } else if (tooth.condition === 'missing') {
      recommendations.push({
        id: `${tooth.tooth_number}-missing`,
        tooth_number: tooth.tooth_number,
        recommendation: 'Consider implant or bridge replacement',
        priority: 'Medium'
      })
    }
  })
  
  return recommendations
})

const calculateAge = (dateOfBirth) => {
  const today = new Date()
  const birthDate = new Date(dateOfBirth)
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--
  }
  
  return age
}

const formatCondition = (condition) => {
  return condition.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatRelativeTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 1) return 'Yesterday'
  if (diffDays < 7) return `${diffDays} days ago`
  if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`
  return `${Math.ceil(diffDays / 30)} months ago`
}

const loadPatientChart = async () => {
  if (!selectedPatientId.value) {
    error.value = ''
    return
  }

  loading.value = true
  error.value = ''

  try {
    // Load patient dental records
    await dentalRecordStore.fetchPatientDentalRecords(parseInt(selectedPatientId.value))
    
    // Set initial chart type based on available records
    const permanentRecord = dentalRecordStore.getPermanentTeethRecord
    const primaryRecord = dentalRecordStore.getPrimaryTeethRecord
    
    if (permanentRecord) {
      dentalRecordStore.setCurrentRecord(permanentRecord)
    } else if (primaryRecord) {
      dentalRecordStore.setCurrentRecord(primaryRecord)
    }
    
  } catch (err) {
    console.error('Error loading patient dental chart:', err)
    error.value = err.message || 'Failed to load dental chart'
  } finally {
    loading.value = false
  }
}

const retryLoad = () => {
  loadPatientChart()
}

const handleToothUpdated = (toothData) => {
  notificationStore.showSuccess(`Tooth ${toothData.tooth_number} updated successfully`)
}

const handleChartUpdated = () => {
  // Refresh any computed values or additional data if needed
  console.log('Dental chart updated')
}

// Initialize component
onMounted(async () => {
  // Load patients if not already loaded
  if (patientStore.patients.length === 0) {
    await patientStore.fetchPatients()
  }
  
  // Check for patient ID in route params
  if (route.params.patientId) {
    selectedPatientId.value = route.params.patientId
    await loadPatientChart()
  }
})

// Watch for route changes
watch(() => route.params.patientId, (newPatientId) => {
  if (newPatientId) {
    selectedPatientId.value = newPatientId
    loadPatientChart()
  }
})
</script>

<style scoped>
@import "../styles/main.css";

.patient-dental-chart-page {
  @apply p-6;
}

.patient-selector {
  @apply min-w-[200px];
}

.loading-state {
  @apply bg-white rounded-lg shadow-sm border border-gray-200;
}

.error-state,
.no-patient-state {
  @apply mx-auto max-w-md;
}

.patient-summary {
  @apply transition-shadow hover:shadow-md;
}

.stat-item {
  @apply px-3;
}

.update-item {
  @apply transition-colors hover:bg-gray-100;
}

.recommendation-item {
  @apply transition-colors hover:bg-yellow-100;
}

@media (max-width: 768px) {
  .page-header {
    @apply flex-col items-start;
  }
  
  .header-actions {
    @apply w-full flex-col items-stretch space-x-0 space-y-3;
  }
  
  .patient-selector {
    @apply w-full;
  }
  
  .patient-summary .flex {
    @apply flex-col space-y-4;
  }
  
  .chart-stats {
    @apply justify-center;
  }
  
  .additional-info {
    @apply grid-cols-1;
  }
}
</style>