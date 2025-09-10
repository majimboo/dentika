<template>
  <div class="appointment-countdown bg-white rounded-lg shadow-sm border p-4">
    <!-- Header -->
    <div class="countdown-header flex items-center justify-between mb-4">
      <div class="flex items-center">
        <div class="countdown-icon mr-2">
          <font-awesome-icon icon="fa-solid fa-clock" class="w-5 h-5 text-blue-600" />
        </div>
        <h3 class="text-sm font-semibold text-gray-900">Upcoming Appointments</h3>
      </div>
      
      <div class="countdown-controls flex items-center space-x-2">
        <button
          @click="refreshAppointments"
          :disabled="loading"
          class="text-gray-400 hover:text-gray-600 transition-colors"
          title="Refresh appointments"
        >
          <font-awesome-icon
            icon="fa-solid fa-sync"
            class="w-4 h-4"
            :class="{ 'animate-spin': loading }"
          />
        </button>
        
        <button
          @click="toggleMinimized"
          class="text-gray-400 hover:text-gray-600 transition-colors"
          :title="minimized ? 'Expand' : 'Minimize'"
        >
          <font-awesome-icon
            :icon="minimized ? 'fa-solid fa-chevron-down' : 'fa-solid fa-chevron-up'"
            class="w-4 h-4"
          />
        </button>
      </div>
    </div>

    <!-- Current Time -->
    <div class="current-time text-center mb-4">
      <div class="time-display text-lg font-mono font-semibold text-gray-900">
        {{ currentTime }}
      </div>
      <div class="date-display text-xs text-gray-500">
        {{ currentDate }}
      </div>
    </div>

    <!-- Appointments List -->
    <div v-if="!minimized" class="appointments-list">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state flex items-center justify-center py-4">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
        <span class="ml-2 text-sm text-gray-600">Loading appointments...</span>
      </div>

      <!-- No Appointments -->
      <div v-else-if="countdownAppointments.length === 0" class="no-appointments text-center py-6">
        <div class="no-appointments-icon text-gray-300 mb-2">
          <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
        </div>
        <p class="text-sm text-gray-500">No appointments in the next 30 minutes</p>
      </div>

      <!-- Appointment Cards -->
      <div v-else class="appointment-cards space-y-3">
        <div
          v-for="appointment in countdownAppointments"
          :key="appointment.appointment.id"
          class="appointment-card"
          :class="{
            'urgent-appointment': appointment.time_until_minutes <= 5,
            'warning-appointment': appointment.time_until_minutes <= 15 && appointment.time_until_minutes > 5,
            'normal-appointment': appointment.time_until_minutes > 15
          }"
        >
          <!-- Countdown Timer -->
          <div class="appointment-timer flex items-center justify-between mb-2">
            <div class="timer-display">
              <span class="countdown-time text-lg font-mono font-bold">
                {{ formatCountdown(appointment.time_until_minutes) }}
              </span>
              <span class="countdown-label text-xs text-gray-500 ml-1">until appointment</span>
            </div>
            
            <div class="appointment-status">
              <span 
                class="status-badge px-2 py-1 rounded-full text-xs font-medium"
                :class="getStatusClass(appointment.appointment.status)"
              >
                {{ formatStatus(appointment.appointment.status) }}
              </span>
            </div>
          </div>

          <!-- Appointment Details -->
          <div class="appointment-details">
            <div class="patient-info flex items-center mb-1">
              <div class="patient-avatar w-6 h-6 bg-gray-300 rounded-full flex items-center justify-center mr-2">
                <span class="text-xs font-medium text-gray-600">
                  {{ getPatientInitials(appointment.appointment.patient) }}
                </span>
              </div>
              <div class="patient-name text-sm font-medium text-gray-900 flex-1">
                {{ appointment.appointment.patient?.first_name }} {{ appointment.appointment.patient?.last_name }}
              </div>
            </div>
            
            <div class="appointment-meta text-xs text-gray-600 space-y-1">
              <div class="flex items-center">
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                Dr. {{ appointment.appointment.doctor?.first_name }} {{ appointment.appointment.doctor?.last_name }}
              </div>
              
              <div class="flex items-center">
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                {{ formatAppointmentTime(appointment.appointment.start_time) }}
              </div>
              
              <div v-if="appointment.appointment.type" class="flex items-center">
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                {{ formatType(appointment.appointment.type) }}
              </div>
            </div>
          </div>

          <!-- Quick Actions -->
          <div class="appointment-actions flex justify-between items-center mt-3 pt-2 border-t border-gray-100">
            <button
              v-if="!appointment.appointment.patient_arrived"
              @click="markPatientArrived(appointment.appointment)"
              class="action-btn arrived-btn text-xs px-2 py-1 bg-green-100 text-green-700 rounded hover:bg-green-200 transition-colors"
            >
              Mark Arrived
            </button>
            
            <span
              v-else
              class="arrived-status text-xs text-green-600 font-medium flex items-center"
            >
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
              Patient Arrived
              <span v-if="appointment.appointment.is_late" class="ml-1 text-orange-600">(Late)</span>
            </span>
            
            <button
              @click="viewAppointment(appointment.appointment)"
              class="action-btn view-btn text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
            >
              View Details
            </button>
          </div>
        </div>
      </div>

      <!-- Auto-refresh indicator -->
      <div v-if="!loading && countdownAppointments.length > 0" class="auto-refresh-info text-center mt-4 pt-3 border-t border-gray-100">
        <div class="flex items-center justify-center text-xs text-gray-500">
          <div class="refresh-indicator w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></div>
          <span>Auto-refreshing every {{ refreshInterval / 1000 }}s</span>
        </div>
      </div>
    </div>

    <!-- Minimized View -->
    <div v-else class="minimized-view">
      <div class="next-appointment-summary">
        <div v-if="nextAppointment" class="text-center">
          <div class="countdown-time text-xl font-mono font-bold text-blue-600">
            {{ formatCountdown(nextAppointment.time_until_minutes) }}
          </div>
          <div class="next-patient text-xs text-gray-600">
            Next: {{ nextAppointment.appointment.patient?.first_name }} {{ nextAppointment.appointment.patient?.last_name }}
          </div>
        </div>
        <div v-else class="text-center text-sm text-gray-500">
          No upcoming appointments
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAppointmentStore } from '../stores/appointment'

const appointmentStore = useAppointmentStore()

const currentTime = ref('')
const currentDate = ref('')
const loading = ref(false)
const minimized = ref(false)
const refreshInterval = ref(30000) // 30 seconds

let timeInterval = null
let refreshAppointmentsInterval = null

const countdownAppointments = computed(() => {
  return appointmentStore.getAppointmentsWithCountdown || []
})

const nextAppointment = computed(() => {
  return countdownAppointments.value[0] || null
})

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('en-US', {
    hour12: true,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
  currentDate.value = now.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const refreshAppointments = async () => {
  loading.value = true
  try {
    await appointmentStore.fetchUpcomingAppointments()
  } catch (error) {
    console.error('Error fetching upcoming appointments:', error)
  } finally {
    loading.value = false
  }
}

const toggleMinimized = () => {
  minimized.value = !minimized.value
}

const formatCountdown = (minutes) => {
  if (minutes < 0) return '00:00'
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (hours > 0) {
    return `${hours}h ${mins}m`
  } else {
    return `${mins}m`
  }
}

const formatStatus = (status) => {
  return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatType = (type) => {
  return type.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatAppointmentTime = (timeString) => {
  return new Date(timeString).toLocaleTimeString('en-US', {
    hour12: true,
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getPatientInitials = (patient) => {
  if (!patient) return '?'
  const first = patient.first_name?.charAt(0) || ''
  const last = patient.last_name?.charAt(0) || ''
  return (first + last).toUpperCase() || '?'
}

const getStatusClass = (status) => {
  const classes = {
    scheduled: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-red-100 text-red-800',
    no_show: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const markPatientArrived = async (appointment) => {
  try {
    const result = await appointmentStore.markPatientArrived(appointment.id)
    if (result.success) {
      // Refresh appointments to get updated data
      await refreshAppointments()
    }
  } catch (error) {
    console.error('Error marking patient as arrived:', error)
  }
}

const viewAppointment = (appointment) => {
  // Emit event to parent component or use router to navigate
  // This would typically open an appointment detail modal or navigate to appointment page
  console.log('View appointment:', appointment.id)
}

onMounted(() => {
  // Start time updates
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  
  // Initial fetch
  refreshAppointments()
  
  // Set up auto-refresh
  refreshAppointmentsInterval = setInterval(refreshAppointments, refreshInterval.value)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
  if (refreshAppointmentsInterval) {
    clearInterval(refreshAppointmentsInterval)
  }
})
</script>

<style scoped>
.appointment-countdown {
  min-width: 320px;
  max-width: 400px;
}

.appointment-card {
  @apply p-3 border rounded-lg transition-all duration-200;
}

.normal-appointment {
  @apply border-gray-200 bg-gray-50;
}

.warning-appointment {
  @apply border-yellow-300 bg-yellow-50;
}

.urgent-appointment {
  @apply border-red-300 bg-red-50 shadow-md;
}

.countdown-time {
  @apply tracking-wider;
}

.urgent-appointment .countdown-time {
  @apply text-red-600 animate-pulse;
}

.warning-appointment .countdown-time {
  @apply text-yellow-600;
}

.normal-appointment .countdown-time {
  @apply text-blue-600;
}

.patient-avatar {
  @apply text-white bg-gradient-to-r from-blue-400 to-blue-600;
}

.action-btn {
  @apply transition-all duration-200 font-medium;
}

.action-btn:hover {
  @apply transform scale-105;
}

.refresh-indicator {
  animation: pulse 2s infinite;
}

.time-display {
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

.minimized-view {
  @apply py-2;
}

.next-appointment-summary .countdown-time {
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .appointment-countdown {
    min-width: 280px;
  }
  
  .appointment-actions {
    @apply flex-col space-y-2;
  }
  
  .appointment-actions .action-btn {
    @apply w-full text-center;
  }
}

/* Animation for new appointments */
.appointment-card {
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>