<template>
  <div class="day-calendar">
    <!-- Time slots -->
    <div class="calendar-grid">
       <!-- Header -->
       <div class="calendar-header px-4 py-3 border-b">
        <div class="day-header text-center">
          <div class="day-name text-lg font-semibold text-gray-900">
            {{ dayName }}
          </div>
          <div class="day-number text-2xl font-bold text-blue-600">
            {{ dayNumber }}
          </div>
        </div>
      </div>

       <!-- Time Grid -->
       <div class="time-grid flex relative">
         <!-- Time Labels -->
         <div class="time-labels w-16 border-r">
           <div
             v-for="hour in hours"
             :key="hour"
             class="time-label h-16 flex items-center justify-center text-xs text-gray-500 border-b"
           >
             {{ formatHour(hour) }}
           </div>
         </div>

         <!-- Appointment Slots -->
         <div class="appointment-slots flex-1 relative">
           <!-- Time Slot Grid -->
           <div class="time-slot-grid">
             <div
               v-for="hour in hours"
               :key="hour"
               class="hour-slot h-16 border-b border-gray-100 relative cursor-pointer hover:bg-blue-50 transition-colors"
               @click="selectTimeSlot(hour)"
             >
               <!-- 30-minute markers -->
               <div class="absolute top-1/2 left-0 right-0 h-px bg-gray-100"></div>
             </div>
           </div>

           <!-- Appointments -->
           <div class="appointments-layer absolute inset-0">
             <div
               v-for="appointment in dayAppointments"
               :key="appointment.id"
               class="appointment-block absolute left-1 right-1 bg-blue-500 text-white rounded-md shadow-sm cursor-pointer hover:bg-blue-600 transition-colors"
               :style="getAppointmentStyle(appointment)"
               @click="$emit('appointment-click', appointment)"
             >
               <div class="appointment-content p-2">
                 <div class="appointment-time text-xs font-medium">
                   {{ formatAppointmentTime(appointment.start_time) }} - {{ formatAppointmentTime(appointment.end_time) }}
                 </div>
                 <div class="patient-name text-sm font-semibold truncate">
                   {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
                 </div>
                 <div class="appointment-type text-xs opacity-90 truncate">
                   {{ formatAppointmentType(appointment.type) }}
                 </div>
                 <div class="doctor-name text-xs opacity-75 truncate">
                   Dr. {{ appointment.doctor?.first_name }} {{ appointment.doctor?.last_name }}
                 </div>
               </div>

               <!-- Status indicator -->
               <div
                 class="status-indicator absolute top-1 right-1 w-2 h-2 rounded-full"
                 :class="getStatusColor(appointment.status)"
               ></div>

               <!-- Patient arrival indicator -->
               <div
                 v-if="appointment.patient_arrived"
                 class="arrival-indicator absolute bottom-1 right-1"
               >
                 <svg class="w-3 h-3 text-green-300" fill="currentColor" viewBox="0 0 20 20">
                   <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                 </svg>
               </div>
             </div>
           </div>
         </div>

         <!-- Current Time Indicator -->
         <!-- <div
           v-if="isToday"
           class="current-time-line absolute left-16 right-0 h-0.5 bg-red-500 z-10 pointer-events-none"
           :style="currentTimeStyle"
         >
           <div class="absolute -left-16 -top-2 bg-red-500 text-white text-xs px-1 py-0.5 rounded">
             {{ currentTimeLabel }}
           </div>
         </div> -->
       </div>
     </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  currentDate: {
    type: Date,
    required: true
  },
  appointments: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['appointment-click', 'time-slot-click'])

const currentTime = ref(new Date())

// Generate hours from 6 AM to 10 PM
const hours = Array.from({ length: 16 }, (_, i) => i + 6)

const dayName = computed(() => {
  return props.currentDate.toLocaleDateString('en-US', { weekday: 'long' })
})

const dayNumber = computed(() => {
  return props.currentDate.getDate()
})

const isToday = computed(() => {
  const today = new Date()
  const todayInTimezone = new Date(today.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  const currentDateInTimezone = new Date(props.currentDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  return currentDateInTimezone.toDateString() === todayInTimezone.toDateString()
})

const dayAppointments = computed(() => {
  // Format date in Asia/Manila timezone for comparison
  const dateInTimezone = new Date(props.currentDate.getTime() - (props.currentDate.getTimezoneOffset() * 60000))
  const manilaTime = new Date(dateInTimezone.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  const dayStr = manilaTime.toISOString().split('T')[0]

  return props.appointments.filter(apt => {
    const aptDate = new Date(apt.start_time)
    const aptDateInTimezone = new Date(aptDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
    const aptDayStr = aptDateInTimezone.toISOString().split('T')[0]
    return aptDayStr === dayStr
  })
})

  const currentTimeStyle = computed(() => {
   if (!isToday.value) return {}

   const now = currentTime.value
   const hours = now.getHours()
   const minutes = now.getMinutes()

   // Calculate position (6 AM = 0%, 10 PM = 100%)
   const startHour = 6
   const endHour = 22
   const totalMinutes = (endHour - startHour) * 60
   const currentMinutes = (hours - startHour) * 60 + minutes

   if (currentMinutes < 0) {
     return { display: 'none' }
   }

   // If current time is after 10 PM, position at the bottom
   if (currentMinutes > totalMinutes) {
     return { top: '100%', transform: 'translateY(-2px)' }
   }

   const percentage = (currentMinutes / totalMinutes) * 100
   return { top: `${percentage}%` }
 })

const currentTimeLabel = computed(() => {
  return currentTime.value.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
})

const formatHour = (hour) => {
  if (hour === 0) return '12 AM'
  if (hour < 12) return `${hour} AM`
  if (hour === 12) return '12 PM'
  return `${hour - 12} PM`
}

const formatAppointmentTime = (timeString) => {
  return new Intl.DateTimeFormat('en-US', {
    timeZone: 'Asia/Manila',
    hour: 'numeric',
    minute: '2-digit'
  }).format(new Date(timeString))
}

const formatAppointmentType = (type) => {
  return type?.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()) || ''
}

const getAppointmentStyle = (appointment) => {
  const startTime = new Date(appointment.start_time)
  const endTime = new Date(appointment.end_time)
  
  const startHour = 6
  const endHour = 22
  const totalMinutes = (endHour - startHour) * 60
  
  // Calculate start position
  const startMinutes = (startTime.getHours() - startHour) * 60 + startTime.getMinutes()
  const endMinutes = (endTime.getHours() - startHour) * 60 + endTime.getMinutes()
  
  const topPercent = (startMinutes / totalMinutes) * 100
  const heightPercent = ((endMinutes - startMinutes) / totalMinutes) * 100
  
  return {
    top: `${topPercent}%`,
    height: `${heightPercent}%`,
    backgroundColor: getAppointmentColor(appointment.status),
    borderLeft: `4px solid ${getStatusBorderColor(appointment.status)}`
  }
}

const getAppointmentColor = (status) => {
  const colors = {
    scheduled: '#3B82F6', // blue
    confirmed: '#10B981', // green
    in_progress: '#F59E0B', // yellow
    completed: '#6B7280', // gray
    cancelled: '#EF4444', // red
    no_show: '#DC2626' // dark red
  }
  return colors[status] || colors.scheduled
}

const getStatusBorderColor = (status) => {
  const colors = {
    scheduled: '#1D4ED8',
    confirmed: '#059669',
    in_progress: '#D97706',
    completed: '#374151',
    cancelled: '#DC2626',
    no_show: '#991B1B'
  }
  return colors[status] || colors.scheduled
}

const getStatusColor = (status) => {
  const colors = {
    scheduled: 'bg-blue-300',
    confirmed: 'bg-green-300',
    in_progress: 'bg-yellow-300',
    completed: 'bg-gray-300',
    cancelled: 'bg-red-300',
    no_show: 'bg-red-400'
  }
  return colors[status] || colors.scheduled
}

const selectTimeSlot = (hour) => {
  const slotDate = new Date(props.currentDate)
  slotDate.setHours(hour, 0, 0, 0)
  
  emit('time-slot-click', slotDate, hour)
}

// Update current time every minute
let timeInterval = null

onMounted(() => {
  timeInterval = setInterval(() => {
    currentTime.value = new Date()
  }, 60000) // Update every minute
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
@import "../../styles/main.css";

.day-calendar {
  @apply relative h-full;
}

.calendar-grid {
  @apply h-full flex flex-col;
}

.time-grid {
  @apply flex-1 relative overflow-y-auto;
  height: calc(100vh - 200px);
}

/* Mobile responsive height */
@media (max-width: 768px) {
  .time-grid {
    height: calc(100vh - 300px);
  }
}

@media (max-width: 640px) {
  .time-grid {
    height: calc(100vh - 350px);
  }
}

.appointment-block {
  @apply flex flex-col justify-between;
  min-height: 60px;
}

.appointment-content {
  @apply flex-1 overflow-hidden;
}

.current-time-line {
  @apply relative;
}

.status-indicator {
  @apply shadow-sm;
}

.arrival-indicator {
  @apply opacity-90;
}

/* Scrollbar styling */
.time-grid::-webkit-scrollbar {
  width: 6px;
}

.time-grid::-webkit-scrollbar-track {
  @apply bg-gray-100 rounded;
}

.time-grid::-webkit-scrollbar-thumb {
  @apply bg-gray-400 rounded;
}

.time-grid::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-500;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .appointment-content {
    @apply p-1;
  }

  .appointment-time,
  .patient-name,
  .appointment-type,
  .doctor-name {
    @apply text-xs;
  }

  .patient-name {
    @apply text-sm;
  }

  .time-labels {
    @apply w-12;
  }

  .time-label {
    @apply text-xs px-1;
  }

  .hour-slot {
    @apply h-12;
  }

  .appointment-block {
    min-height: 40px;
  }
}

@media (max-width: 640px) {
  .time-labels {
    @apply w-10;
  }

  .time-label {
    @apply text-xs px-0.5;
  }

  .hour-slot {
    @apply h-10;
  }

  .appointment-block {
    min-height: 35px;
  }

  .appointment-content {
    @apply p-0.5;
  }
}
</style>