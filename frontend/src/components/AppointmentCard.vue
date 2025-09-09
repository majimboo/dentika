<template>
  <div class="bg-white border border-gray-200 rounded-lg p-3 hover:shadow-sm transition-all duration-200 cursor-pointer group">
    <!-- Patient Info -->
    <div class="flex items-start justify-between mb-2">
      <div class="flex-1">
        <h4 class="font-medium text-gray-900 text-sm mb-1">
          {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
        </h4>
        <p class="text-xs text-gray-500">
          ID: {{ appointment.patient?.patient_number || appointment.patient_id }}
        </p>
      </div>

      <!-- Status Badge -->
      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
            :class="getStatusBadgeClass(appointment.status)">
        {{ getStatusDisplayName(appointment.status) }}
      </span>
    </div>

    <!-- Appointment Details -->
    <div class="space-y-1 mb-2">
      <div class="flex items-center text-xs text-gray-600">
        <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        {{ formatTime(appointment.start_time) }} - {{ formatTime(appointment.end_time) }}
      </div>

      <div class="flex items-center text-xs text-gray-600">
        <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
        </svg>
        {{ appointment.type || 'General' }}
      </div>

      <div v-if="appointment.doctor" class="flex items-center text-xs text-gray-600">
        <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
        Dr. {{ appointment.doctor.first_name }} {{ appointment.doctor.last_name }}
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="flex items-center justify-between pt-2 border-t border-gray-100">
      <div class="flex space-x-1">
        <!-- Mark as Arrived (for scheduled appointments) -->
        <button
          v-if="appointment.status === 'scheduled'"
          @click="$emit('status-changed', appointment.id, 'confirmed')"
          class="inline-flex items-center px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md hover:bg-green-200 transition-colors"
          title="Mark as arrived"
        >
          <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          Arrived
        </button>

        <!-- Start Appointment (for confirmed appointments) -->
        <button
          v-if="appointment.status === 'confirmed'"
          @click="$emit('status-changed', appointment.id, 'in_progress')"
          class="inline-flex items-center px-2 py-1 text-xs font-medium text-yellow-700 bg-yellow-100 rounded-md hover:bg-yellow-200 transition-colors"
          title="Start appointment"
        >
          <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.586a1 1 0 01.707.293l.707.707A1 1 0 0012.414 11H13m-3 3h1.586a1 1 0 01.707.293l.707.707A1 1 0 0012.414 14H13"/>
          </svg>
          Start
        </button>

        <!-- Complete Appointment (for in-progress appointments) -->
        <button
          v-if="appointment.status === 'in_progress'"
          @click="$emit('status-changed', appointment.id, 'completed')"
          class="inline-flex items-center px-2 py-1 text-xs font-medium text-purple-700 bg-purple-100 rounded-md hover:bg-purple-200 transition-colors"
          title="Complete appointment"
        >
          <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          Complete
        </button>
      </div>

      <!-- View Details -->
      <router-link
        :to="`/appointments/${appointment.id}`"
        class="inline-flex items-center px-2 py-1 text-xs font-medium text-neutral-700 bg-neutral-100 rounded-md hover:bg-neutral-200 transition-colors"
        title="View appointment details"
      >
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
        </svg>
      </router-link>
    </div>

    <!-- Time until appointment (for scheduled appointments) -->
    <div v-if="appointment.status === 'scheduled'" class="mt-2 text-xs text-neutral-500">
      {{ getTimeUntilAppointment(appointment.start_time) }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppointmentCard',
  props: {
    appointment: {
      type: Object,
      required: true
    }
  },
  emits: ['status-changed'],
  methods: {
    getStatusBadgeClass(status) {
      const classes = {
        'scheduled': 'bg-blue-100 text-blue-800',
        'confirmed': 'bg-green-100 text-green-800',
        'in_progress': 'bg-yellow-100 text-yellow-800',
        'completed': 'bg-purple-100 text-purple-800',
        'cancelled': 'bg-red-100 text-red-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    },

    getStatusDisplayName(status) {
      const names = {
        'scheduled': 'Scheduled',
        'confirmed': 'Arrived',
        'in_progress': 'In Progress',
        'completed': 'Completed',
        'cancelled': 'No Show'
      }
      return names[status] || status
    },

    formatTime(timeString) {
      if (!timeString) return ''
      try {
        const date = new Date(timeString)
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      } catch {
        return timeString
      }
    },

    getTimeUntilAppointment(startTime) {
      if (!startTime) return ''

      const now = new Date()
      const appointmentTime = new Date(startTime)
      const diffMs = appointmentTime - now
      const diffMins = Math.floor(diffMs / (1000 * 60))

      if (diffMins < 0) {
        return `${Math.abs(diffMins)} minutes late`
      } else if (diffMins === 0) {
        return 'Starting now'
      } else if (diffMins < 60) {
        return `In ${diffMins} minutes`
      } else {
        const diffHours = Math.floor(diffMins / 60)
        const remainingMins = diffMins % 60
        return `In ${diffHours}h ${remainingMins}m`
      }
    }
  }
}
</script>