<template>
  <div class="week-calendar">
    <div class="calendar-grid">
       <!-- Header with days -->
       <div class="calendar-header border-b px-2 py-3" :style="headerGridStyle">
        <!-- Empty cell for time column -->
        <div class="time-column-header"></div>

         <!-- Day headers -->
         <div
           v-for="(day, dayIndex) in weekDays"
           :key="day.date.getTime()"
           :class="[
             'day-header px-2 py-1 text-center',
             { 'bg-blue-50': isToday(day.date), 'text-blue-600': isToday(day.date) },
             dayIndex < weekDays.length - 1 ? 'border-r' : ''
           ]"
         >
          <div class="day-name text-xs font-medium text-gray-600 uppercase">
            {{ day.name }}
          </div>
          <div class="day-number text-lg font-semibold mt-1"
               :class="isToday(day.date) ? 'text-blue-600' : 'text-gray-900'">
            {{ day.number }}
          </div>
        </div>
      </div>

       <!-- Time Grid -->
       <div class="time-grid flex-1 overflow-y-auto relative" :style="combinedGridStyle">
          <!-- Time Labels Column -->
          <div class="time-labels border-r">
            <div
              v-for="hour in hours"
              :key="hour"
              class="time-label h-16 flex items-center justify-center text-xs text-gray-500 border-b px-1"
            >
             {{ formatHour(hour) }}
           </div>
         </div>

          <!-- Day Columns -->
          <div
            v-for="(day, dayIndex) in weekDays"
            :key="day.date.getTime()"
            :class="[
              'day-column relative',
              dayIndex < weekDays.length - 1 ? 'border-r' : ''
            ]"
          >
            <!-- Time Slots -->
            <div class="time-slots">
              <div
                v-for="hour in hours"
                :key="hour"
                :class="[
                  'time-slot h-16 border-b border-gray-100 transition-colors relative',
                  isTimeSlotAvailable(day.date, hour)
                    ? 'cursor-pointer hover:bg-blue-50'
                    : 'cursor-not-allowed bg-red-50 opacity-60'
                ]"
                @click="isTimeSlotAvailable(day.date, hour) ? selectTimeSlot(day.date, hour) : null"
              >
                <!-- 30-minute marker -->
                <div class="absolute top-1/2 left-0 right-0 h-px bg-gray-50"></div>

                <!-- Unavailable indicator -->
                <div
                  v-if="!isTimeSlotAvailable(day.date, hour)"
                  class="absolute inset-0 flex items-center justify-center"
                >
                  <div class="text-red-400 text-xs font-medium">
                    Busy
                  </div>
                </div>
              </div>
            </div>

           <!-- Appointments for this day -->
           <div class="appointments-layer absolute inset-0 pointer-events-none">
             <div
               v-for="appointment in getDayAppointments(day.date)"
               :key="appointment.id"
               class="appointment-block absolute left-0.5 right-0.5 text-white rounded-md shadow-sm cursor-pointer hover:shadow-md transition-shadow pointer-events-auto"
               :style="getAppointmentStyle(appointment)"
               @click="$emit('appointment-click', appointment)"
             >
               <div class="appointment-content p-1.5">
                 <div class="appointment-time text-xs font-medium">
                   {{ formatTime(appointment.start_time) }}
                 </div>
                 <div class="patient-name text-sm font-semibold truncate">
                   {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
                 </div>
                 <div class="appointment-type text-xs opacity-90 truncate">
                   {{ formatAppointmentType(appointment.type) }}
                 </div>
               </div>

               <!-- Status indicator -->
               <div
                 class="status-indicator absolute top-1 right-1 w-1.5 h-1.5 rounded-full"
                 :class="getStatusIndicatorColor(appointment.status)"
               ></div>

               <!-- Patient arrival indicator -->
               <div
                 v-if="appointment.patient_arrived"
                 class="arrival-indicator absolute bottom-1 right-1"
               >
                 <svg class="w-2.5 h-2.5 text-green-300" fill="currentColor" viewBox="0 0 20 20">
                   <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                 </svg>
               </div>
             </div>
           </div>
         </div>

         <!-- Current Time Indicator -->
         <!-- <div
           v-if="todayColumn !== -1"
           class="current-time-line absolute bg-red-500 z-10 pointer-events-none"
           :style="getCurrentTimeStyle()"
         >
           <div class="absolute -left-16 -top-2 bg-red-500 text-white text-xs px-1 py-0.5 rounded whitespace-nowrap">
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

const weekDays = computed(() => {
  const days = []
  const startOfWeek = new Date(props.currentDate)
  
  // Go to Sunday of current week
  startOfWeek.setDate(props.currentDate.getDate() - props.currentDate.getDay())
  
  for (let i = 0; i < 7; i++) {
    const day = new Date(startOfWeek)
    day.setDate(startOfWeek.getDate() + i)
    
    days.push({
      date: day,
      name: day.toLocaleDateString('en-US', { weekday: 'short' }),
      number: day.getDate()
    })
  }
  
  return days
})

const todayColumn = computed(() => {
  const today = new Date()
  const todayInTimezone = new Date(today.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  return weekDays.value.findIndex(day => {
    const dayInTimezone = new Date(day.date.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
    return dayInTimezone.toDateString() === todayInTimezone.toDateString()
  })
})

const currentTimeLabel = computed(() => {
  return currentTime.value.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
})

 const isToday = (date) => {
   const today = new Date()
   const todayInTimezone = new Date(today.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
   const dateInTimezone = new Date(date.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
   return dateInTimezone.toDateString() === todayInTimezone.toDateString()
 }

const timeLabelsWidth = computed(() => {
  // Match the width of time-labels column
  if (window.innerWidth <= 640) {
    return '40px'
  } else if (window.innerWidth <= 768) {
    return '48px'
  } else {
    return '64px' // w-16 = 4rem = 64px
  }
})

const headerGridStyle = computed(() => {
  const timeWidth = timeLabelsWidth.value
  return {
    display: 'grid',
    gridTemplateColumns: `${timeWidth} repeat(7, 1fr)`
  }
})

const timeGridStyle = computed(() => {
  // Responsive height based on screen size
  if (window.innerWidth <= 640) {
    return { height: 'calc(100vh - 400px)' }
  } else if (window.innerWidth <= 768) {
    return { height: 'calc(100vh - 350px)' }
  } else {
    return { height: 'calc(100vh - 280px)' }
  }
})

const combinedGridStyle = computed(() => {
  const timeWidth = timeLabelsWidth.value
  const height = timeGridStyle.value.height
  return {
    ...timeGridStyle.value,
    display: 'grid',
    gridTemplateColumns: `${timeWidth} repeat(7, 1fr)`
  }
})

const formatHour = (hour) => {
  if (hour === 0) return '12 AM'
  if (hour < 12) return `${hour} AM`
  if (hour === 12) return '12 PM'
  return `${hour - 12} PM`
}

const formatTime = (timeString) => {
  return new Intl.DateTimeFormat('en-US', {
    timeZone: 'Asia/Manila',
    hour: 'numeric',
    minute: '2-digit'
  }).format(new Date(timeString))
}

const formatAppointmentType = (type) => {
  return type?.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()) || ''
}

const getDayAppointments = (date) => {
  // Format date in Asia/Manila timezone for comparison
  const dateInTimezone = new Date(date.getTime() - (date.getTimezoneOffset() * 60000))
  const manilaTime = new Date(dateInTimezone.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  const dayStr = manilaTime.toISOString().split('T')[0]

  return props.appointments.filter(apt => {
    const aptDate = new Date(apt.start_time)
    const aptDateInTimezone = new Date(aptDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
    const aptDayStr = aptDateInTimezone.toISOString().split('T')[0]
    return aptDayStr === dayStr
  })
}

const isTimeSlotAvailable = (date, hour) => {
  const slotStart = new Date(date)
  slotStart.setHours(hour, 0, 0, 0)
  const slotEnd = new Date(slotStart)
  slotEnd.setHours(hour + 1, 0, 0, 0)

  const dayAppointments = getDayAppointments(date)

  // Check if this time slot conflicts with any existing appointments
  return !dayAppointments.some(appointment => {
    const aptStart = new Date(appointment.start_time)
    const aptEnd = new Date(appointment.end_time)

    // Check for overlap: slot starts before appointment ends AND slot ends after appointment starts
    return slotStart < aptEnd && slotEnd > aptStart
  })
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
    height: `${Math.max(heightPercent, 2.5)}%`, // Minimum height for visibility
    backgroundColor: getAppointmentColor(appointment.status),
    borderLeft: `3px solid ${getStatusBorderColor(appointment.status)}`
  }
}

const getAppointmentColor = (status) => {
  const colors = {
    scheduled: '#3B82F6',
    confirmed: '#10B981',
    in_progress: '#F59E0B',
    completed: '#6B7280',
    cancelled: '#EF4444',
    no_show: '#DC2626'
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

const getStatusIndicatorColor = (status) => {
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

const getCurrentTimeStyle = () => {
   if (todayColumn.value === -1) return { display: 'none' }

   const now = currentTime.value
   const hours = now.getHours()
   const minutes = now.getMinutes()

   const startHour = 6
   const endHour = 22
   const totalMinutes = (endHour - startHour) * 60
   const currentMinutes = (hours - startHour) * 60 + minutes

   if (currentMinutes < 0) {
     return { display: 'none' }
   }

   // If current time is after 10 PM, position at the bottom
   if (currentMinutes > totalMinutes) {
     return {
       top: '100%',
       left: `${(todayColumn.value + 1) * 12.5}%`, // Position in today's column (skip time column)
       width: '12.5%', // Each day column is 12.5% (100% / 8 columns, but skip time column)
       height: '2px',
       transform: 'translateY(-2px)'
     }
   }

   const topPercent = (currentMinutes / totalMinutes) * 100

   return {
     top: `${topPercent}%`,
     left: `${(todayColumn.value + 1) * 12.5}%`, // Position in today's column (skip time column)
     width: '12.5%', // Each day column is 12.5% (100% / 8 columns, but skip time column)
     height: '2px'
   }
 }

const selectTimeSlot = (date, hour) => {
  const slotDate = new Date(date)
  slotDate.setHours(hour, 0, 0, 0)
  
  emit('time-slot-click', slotDate, hour)
}

// Update current time every minute
let timeInterval = null

onMounted(() => {
  timeInterval = setInterval(() => {
    currentTime.value = new Date()
  }, 60000)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
@import "../../styles/main.css";

.week-calendar {
  @apply relative h-full;
}

.calendar-grid {
  @apply h-full flex flex-col;
}

.time-grid {
  @apply relative;
}

.day-column {
  @apply relative;
}

.appointment-block {
  @apply flex flex-col justify-between overflow-hidden;
  min-height: 20px;
}

.appointment-content {
  @apply flex-1 overflow-hidden;
}

 .current-time-line {
   @apply relative;
 }

/* Custom scrollbar */
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
  .appointment-type {
    @apply text-xs;
  }

  .patient-name {
    @apply text-xs font-medium;
  }

  .day-header {
    @apply p-1;
  }

  .day-number {
    @apply text-sm;
  }



  .time-label {
    @apply text-xs px-1 h-12;
  }

  .time-slot {
    @apply h-12;
  }

  .appointment-block {
    min-height: 30px;
  }
}

/* Grid layout adjustments for mobile */
@media (max-width: 640px) {
  .time-label {
    @apply text-xs px-1 h-10;
  }

  .time-slot {
    @apply h-10;
  }

  .appointment-block {
    min-height: 25px;
  }

  .appointment-content {
    @apply p-0.5;
  }

  .day-header {
    @apply p-0.5;
  }

  .day-number {
    @apply text-xs;
  }
}
</style>