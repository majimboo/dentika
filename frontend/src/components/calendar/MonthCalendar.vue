<template>
  <div class="month-calendar">
    <div class="calendar-grid">
      <!-- Days of week header -->
      <div class="calendar-header grid grid-cols-7 border-b">
        <div 
          v-for="day in daysOfWeek" 
          :key="day"
          class="day-header p-2 text-center text-sm font-medium text-gray-600 uppercase"
        >
          {{ day }}
        </div>
      </div>

      <!-- Calendar days grid -->
      <div class="calendar-body grid grid-cols-7">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          class="calendar-day border-r border-b min-h-24 relative cursor-pointer transition-colors"
          :class="{
            'bg-gray-50 text-gray-400': !day.isCurrentMonth,
            'bg-blue-50 border-blue-200': day.isToday,
            'hover:bg-blue-25': day.isCurrentMonth && !day.isToday
          }"
          @click="selectDate(day.date)"
        >
          <!-- Day number -->
          <div class="day-number p-1">
            <span 
              class="inline-block w-6 h-6 text-sm leading-6 text-center rounded-full"
              :class="{
                'bg-blue-600 text-white font-semibold': day.isToday,
                'text-gray-900 font-medium': day.isCurrentMonth && !day.isToday,
                'text-gray-400': !day.isCurrentMonth
              }"
            >
              {{ day.date.getDate() }}
            </span>
          </div>

          <!-- Appointment indicators -->
          <div class="appointments-container absolute inset-x-1 top-7 bottom-1 overflow-hidden">
            <!-- Show first few appointments as colored bars -->
            <div 
              v-for="(appointment, idx) in getDayAppointments(day.date).slice(0, 4)"
              :key="appointment.id"
              class="appointment-indicator h-1.5 rounded-sm mb-0.5 text-xs"
              :style="{ backgroundColor: getAppointmentColor(appointment.status) }"
            >
              <div 
                v-if="idx === 0"
                class="appointment-text text-xs text-white truncate px-1 leading-none pt-px"
              >
                {{ formatTime(appointment.start_time) }} {{ appointment.patient?.first_name }}
              </div>
            </div>

            <!-- More indicator if there are additional appointments -->
            <div 
              v-if="getDayAppointments(day.date).length > 4"
              class="more-appointments text-xs text-gray-500 px-1 mt-0.5"
            >
              +{{ getDayAppointments(day.date).length - 4 }} more
            </div>
          </div>

          <!-- Day status indicators -->
          <div class="day-status absolute bottom-1 right-1 flex space-x-1">
            <!-- Today indicator -->
            <div 
              v-if="day.isToday"
              class="w-1 h-1 bg-blue-600 rounded-full"
            ></div>
            
            <!-- Has appointments indicator -->
            <div 
              v-if="getDayAppointments(day.date).length > 0"
              class="w-1 h-1 bg-green-500 rounded-full"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Appointment count summary -->
    <div class="appointment-summary mt-4 p-3 bg-gray-50 rounded-lg">
      <div class="summary-stats grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
        <div class="stat-item">
          <div class="stat-value text-lg font-semibold text-gray-900">
            {{ totalAppointments }}
          </div>
          <div class="stat-label text-gray-600">Total Appointments</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value text-lg font-semibold text-blue-600">
            {{ confirmedAppointments }}
          </div>
          <div class="stat-label text-gray-600">Confirmed</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value text-lg font-semibold text-yellow-600">
            {{ pendingAppointments }}
          </div>
          <div class="stat-label text-gray-600">Pending</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value text-lg font-semibold text-green-600">
            {{ completedAppointments }}
          </div>
          <div class="stat-label text-gray-600">Completed</div>
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="calendar-legend mt-3 flex flex-wrap items-center gap-4 text-xs text-gray-600">
      <div class="legend-item flex items-center">
        <div class="w-3 h-1.5 bg-blue-500 rounded-sm mr-1"></div>
        Scheduled
      </div>
      <div class="legend-item flex items-center">
        <div class="w-3 h-1.5 bg-green-500 rounded-sm mr-1"></div>
        Confirmed
      </div>
      <div class="legend-item flex items-center">
        <div class="w-3 h-1.5 bg-yellow-500 rounded-sm mr-1"></div>
        In Progress
      </div>
      <div class="legend-item flex items-center">
        <div class="w-3 h-1.5 bg-gray-500 rounded-sm mr-1"></div>
        Completed
      </div>
      <div class="legend-item flex items-center">
        <div class="w-3 h-1.5 bg-red-500 rounded-sm mr-1"></div>
        Cancelled
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

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

const emit = defineEmits(['appointment-click', 'date-click'])

const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

const calendarDays = computed(() => {
  const days = []
  const year = props.currentDate.getFullYear()
  const month = props.currentDate.getMonth()
  
  // Get first day of the month and find the starting date (Sunday of the first week)
  const firstDayOfMonth = new Date(year, month, 1)
  const startDate = new Date(firstDayOfMonth)
  startDate.setDate(startDate.getDate() - firstDayOfMonth.getDay())
  
  // Generate 42 days (6 weeks) for the calendar grid
  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    
    const isCurrentMonth = date.getMonth() === month
    const isToday = date.toDateString() === new Date().toDateString()
    
    days.push({
      date,
      isCurrentMonth,
      isToday
    })
  }
  
  return days
})

const totalAppointments = computed(() => {
  return props.appointments.length
})

const confirmedAppointments = computed(() => {
  return props.appointments.filter(apt => apt.status === 'confirmed').length
})

const pendingAppointments = computed(() => {
  return props.appointments.filter(apt => apt.status === 'scheduled').length
})

const completedAppointments = computed(() => {
  return props.appointments.filter(apt => apt.status === 'completed').length
})

const getDayAppointments = (date) => {
  const dayStr = date.toISOString().split('T')[0]
  return props.appointments.filter(apt => 
    apt.start_time.startsWith(dayStr)
  ).sort((a, b) => new Date(a.start_time) - new Date(b.start_time))
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

const formatTime = (timeString) => {
  return new Date(timeString).toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  })
}

const selectDate = (date) => {
  emit('date-click', date)
}
</script>

<style scoped>
@import "../../styles/main.css";

.month-calendar {
  @apply h-full;
}

.calendar-grid {
  @apply border border-gray-200 rounded-lg overflow-hidden;
}

.calendar-day {
  @apply relative;
  min-height: 96px;
}

.calendar-day:hover .appointment-indicator {
  @apply shadow-sm;
}

.appointment-indicator {
  @apply transition-shadow;
}

.hover\:bg-blue-25:hover {
  background-color: rgba(59, 130, 246, 0.05);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .calendar-day {
    min-height: 80px;
  }
  
  .appointment-text {
    @apply text-xs;
  }
  
  .more-appointments {
    @apply text-xs;
  }
  
  .summary-stats {
    @apply grid-cols-2;
  }
}

@media (max-width: 640px) {
  .calendar-day {
    min-height: 70px;
  }
  
  .day-number {
    @apply p-0.5;
  }
  
  .appointment-indicator {
    @apply h-1 mb-0.5;
  }
  
  .appointment-text {
    @apply hidden;
  }
  
  .appointments-container {
    @apply top-6;
  }
}

/* Calendar day animations */
.calendar-day {
  @apply transition-all duration-200;
}

.appointment-indicator {
  @apply transition-all duration-150;
}

/* Today highlight animation */
.calendar-day.bg-blue-50 {
  @apply shadow-sm;
  animation: todayPulse 2s ease-in-out infinite;
}

@keyframes todayPulse {
  0%, 100% { 
    box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.1);
  }
  50% { 
    box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
  }
}
</style>