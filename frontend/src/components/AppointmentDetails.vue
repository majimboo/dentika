<template>
  <div class="appointment-details h-full flex flex-col bg-white">
    <!-- Header -->
    <div class="details-header flex items-center justify-between p-4 border-b bg-gray-50">
      <div class="header-content">
        <h3 class="text-lg font-semibold text-gray-900">Appointment Details</h3>
        <div class="appointment-status mt-1">
          <span 
            class="status-badge px-2 py-1 text-xs font-medium rounded-full"
            :class="getStatusBadgeClass(appointment.status)"
          >
            {{ formatStatus(appointment.status) }}
          </span>
        </div>
      </div>
      
      <button
        @click="$emit('close')"
        class="close-btn p-2 hover:bg-gray-200 rounded-lg transition-colors"
      >
        <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
    </div>

    <!-- Content -->
    <div class="details-content flex-1 overflow-y-auto">
      <!-- Patient Information -->
      <div class="patient-section p-4 border-b">
        <h4 class="section-title text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">
          Patient Information
        </h4>
        
        <div class="patient-info space-y-3">
          <div class="info-row flex items-center">
            <div class="info-icon w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3">
              <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
            </div>
            <div class="info-content">
              <div class="patient-name text-base font-semibold text-gray-900">
                {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
              </div>
              <div class="patient-details text-sm text-gray-600">
                {{ formatAge(appointment.patient?.date_of_birth) }} years old
                <span v-if="appointment.patient?.gender"> • {{ appointment.patient.gender }}</span>
              </div>
            </div>
          </div>
          
          <div class="contact-info">
            <div v-if="appointment.patient?.phone" class="info-item flex items-center text-sm text-gray-600 mb-1">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
              </svg>
              {{ appointment.patient.phone }}
            </div>
            
            <div v-if="appointment.patient?.email" class="info-item flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
              </svg>
              {{ appointment.patient.email }}
            </div>
          </div>
        </div>
      </div>

      <!-- Appointment Information -->
      <div class="appointment-section p-4 border-b">
        <h4 class="section-title text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">
          Appointment Details
        </h4>
        
        <div class="appointment-info space-y-3">
          <div class="info-row">
            <label class="info-label text-xs font-medium text-gray-500 uppercase tracking-wide">Date & Time</label>
            <div class="info-value text-sm text-gray-900 mt-1">
              <div class="flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
                {{ formatDateTime(appointment.start_time) }}
              </div>
              <div class="duration text-xs text-gray-500 mt-1 ml-6">
                Duration: {{ calculateDuration(appointment.start_time, appointment.end_time) }} minutes
              </div>
            </div>
          </div>
          
          <div class="info-row">
            <label class="info-label text-xs font-medium text-gray-500 uppercase tracking-wide">Type</label>
            <div class="info-value text-sm text-gray-900 mt-1">
              {{ formatAppointmentType(appointment.type) }}
            </div>
          </div>
          
          <div class="info-row" v-if="appointment.doctor">
            <label class="info-label text-xs font-medium text-gray-500 uppercase tracking-wide">Doctor</label>
            <div class="info-value text-sm text-gray-900 mt-1">
              Dr. {{ appointment.doctor.first_name }} {{ appointment.doctor.last_name }}
            </div>
          </div>
          
          <div class="info-row" v-if="appointment.notes">
            <label class="info-label text-xs font-medium text-gray-500 uppercase tracking-wide">Notes</label>
            <div class="info-value text-sm text-gray-700 mt-1 p-2 bg-gray-50 rounded">
              {{ appointment.notes }}
            </div>
          </div>
        </div>
      </div>

      <!-- Patient Arrival -->
      <div class="arrival-section p-4 border-b">
        <h4 class="section-title text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">
          Patient Status
        </h4>
        
        <div class="arrival-status">
          <div class="flex items-center justify-between">
            <div class="status-info">
              <div class="flex items-center">
                <div 
                  class="status-indicator w-3 h-3 rounded-full mr-2"
                  :class="appointment.patient_arrived ? 'bg-green-500' : 'bg-gray-300'"
                ></div>
                <span class="text-sm font-medium" :class="appointment.patient_arrived ? 'text-green-700' : 'text-gray-600'">
                  {{ appointment.patient_arrived ? 'Patient Arrived' : 'Waiting for Patient' }}
                </span>
              </div>
              <div v-if="appointment.patient_arrived && appointment.arrived_at" class="text-xs text-gray-500 mt-1 ml-5">
                Arrived at {{ formatTime(appointment.arrived_at) }}
              </div>
            </div>
            
            <button
              v-if="!appointment.patient_arrived && canMarkArrival"
              @click="markPatientArrived"
              class="btn-small btn-primary"
            >
              Mark Arrived
            </button>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="actions-section p-4">
        <h4 class="section-title text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">
          Quick Actions
        </h4>
        
        <div class="actions-grid space-y-2">
          <button
            v-if="appointment.status === 'scheduled'"
            @click="updateStatus('confirmed')"
            class="action-btn w-full flex items-center justify-center px-3 py-2 bg-green-50 text-green-700 hover:bg-green-100 rounded-md transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
            Confirm Appointment
          </button>
          
          <button
            v-if="appointment.status === 'confirmed'"
            @click="updateStatus('in_progress')"
            class="action-btn w-full flex items-center justify-center px-3 py-2 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded-md transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            Start Appointment
          </button>
          
          <button
            v-if="appointment.status === 'in_progress'"
            @click="updateStatus('completed')"
            class="action-btn w-full flex items-center justify-center px-3 py-2 bg-green-50 text-green-700 hover:bg-green-100 rounded-md transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            Complete Appointment
          </button>
          
          <button
            @click="editAppointment"
            class="action-btn w-full flex items-center justify-center px-3 py-2 bg-gray-50 text-gray-700 hover:bg-gray-100 rounded-md transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
            </svg>
            Edit Appointment
          </button>
          
          <button
            v-if="appointment.status !== 'cancelled'"
            @click="showCancelConfirm = true"
            class="action-btn w-full flex items-center justify-center px-3 py-2 bg-red-50 text-red-700 hover:bg-red-100 rounded-md transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
            Cancel Appointment
          </button>
        </div>
      </div>
    </div>

    <!-- Cancel Confirmation Modal -->
    <div v-if="showCancelConfirm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
        <div class="flex items-center mb-4">
          <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center mr-4">
            <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Cancel Appointment</h3>
            <p class="text-sm text-gray-600 mt-1">This action cannot be undone.</p>
          </div>
        </div>
        
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Cancellation Reason (Optional)
          </label>
          <textarea
            v-model="cancelReason"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-red-500 focus:border-red-500 resize-none"
            placeholder="Enter reason for cancellation..."
          ></textarea>
        </div>
        
        <div class="flex justify-end space-x-3">
          <button
            @click="showCancelConfirm = false"
            class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
          >
            Keep Appointment
          </button>
          <button
            @click="cancelAppointment"
            class="px-4 py-2 text-white bg-red-600 hover:bg-red-700 rounded-md transition-colors"
          >
            Cancel Appointment
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  appointment: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'update'])

const authStore = useAuthStore()
const showCancelConfirm = ref(false)
const cancelReason = ref('')

const canMarkArrival = computed(() => {
  return authStore.hasPermission('appointments.update')
})

const getStatusBadgeClass = (status) => {
  const classes = {
    scheduled: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-red-100 text-red-800',
    no_show: 'bg-red-200 text-red-900'
  }
  return classes[status] || classes.scheduled
}

const formatStatus = (status) => {
  const statusMap = {
    scheduled: 'Scheduled',
    confirmed: 'Confirmed',
    in_progress: 'In Progress',
    completed: 'Completed',
    cancelled: 'Cancelled',
    no_show: 'No Show'
  }
  return statusMap[status] || status
}

const formatAge = (dateOfBirth) => {
  if (!dateOfBirth) return 'N/A'
  const today = new Date()
  const birth = new Date(dateOfBirth)
  const age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    return age - 1
  }
  return age
}

const formatDateTime = (dateTime) => {
  const date = new Date(dateTime)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatTime = (dateTime) => {
  return new Date(dateTime).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatAppointmentType = (type) => {
  return type?.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()) || ''
}

const calculateDuration = (startTime, endTime) => {
  const start = new Date(startTime)
  const end = new Date(endTime)
  return Math.round((end - start) / (1000 * 60))
}

const markPatientArrived = async () => {
  await emit('update', props.appointment.id, {
    patient_arrived: true,
    arrived_at: new Date().toISOString()
  })
}

const updateStatus = async (newStatus) => {
  await emit('update', props.appointment.id, {
    status: newStatus
  })
}

const editAppointment = () => {
  // This would trigger edit mode
  emit('edit', props.appointment)
}

const cancelAppointment = async () => {
  await emit('update', props.appointment.id, {
    status: 'cancelled',
    cancellation_reason: cancelReason.value,
    cancelled_at: new Date().toISOString()
  })
  showCancelConfirm.value = false
  cancelReason.value = ''
}
</script>

<style scoped>
@reference "tailwindcss";

.btn-small {
  @apply px-2 py-1 text-xs font-medium rounded;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 transition-colors;
}

.section-title {
  @apply border-b border-gray-100 pb-2;
}

.info-label {
  @apply block;
}

.info-value {
  @apply font-medium;
}

.action-btn:hover {
  @apply shadow-sm;
}

/* Custom scrollbar for content */
.details-content::-webkit-scrollbar {
  width: 4px;
}

.details-content::-webkit-scrollbar-track {
  @apply bg-gray-100;
}

.details-content::-webkit-scrollbar-thumb {
  @apply bg-gray-400 rounded;
}

.details-content::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-500;
}
</style>