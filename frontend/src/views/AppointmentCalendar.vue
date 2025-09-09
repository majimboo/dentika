<template>
  <div class="appointment-calendar-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Appointment Calendar</h1>
        <p class="text-gray-600">Schedule and manage patient appointments</p>
      </div>
      
      <div class="header-actions flex items-center space-x-3">
        <!-- View Toggle -->
        <div class="view-toggle flex bg-gray-100 rounded-lg p-1">
          <button
            @click="currentView = 'day'"
            :class="currentView === 'day' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600'"
            class="px-3 py-1 text-sm font-medium rounded-md transition-colors"
          >
            Day
          </button>
          <button
            @click="currentView = 'week'"
            :class="currentView === 'week' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600'"
            class="px-3 py-1 text-sm font-medium rounded-md transition-colors"
          >
            Week
          </button>
          <button
            @click="currentView = 'month'"
            :class="currentView === 'month' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600'"
            class="px-3 py-1 text-sm font-medium rounded-md transition-colors"
          >
            Month
          </button>
        </div>

        <button
          @click="goToToday"
          class="btn btn-secondary"
        >
          Today
        </button>
        
        <router-link
          to="/appointments/list"
          class="btn btn-secondary flex items-center"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
          </svg>
          List View
        </router-link>

        <router-link
          to="/appointments/new"
          class="btn btn-primary flex items-center"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          New Appointment
        </router-link>
      </div>
    </div>

    <!-- Calendar Navigation -->
    <div class="calendar-navigation flex items-center justify-between mb-6 bg-white rounded-lg shadow-sm border p-4">
      <div class="nav-controls flex items-center space-x-4">
        <button
          @click="navigateCalendar(-1)"
          class="nav-btn p-2 rounded-md hover:bg-gray-100 transition-colors"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
        </button>
        
        <div class="current-period">
          <h2 class="text-lg font-semibold text-gray-900">{{ currentPeriodLabel }}</h2>
        </div>
        
        <button
          @click="navigateCalendar(1)"
          class="nav-btn p-2 rounded-md hover:bg-gray-100 transition-colors"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
          </svg>
        </button>
      </div>

      <!-- Filters -->
      <div class="calendar-filters flex items-center space-x-3">
        <!-- Doctor Filter -->
        <div class="filter-group">
          <select
            v-model="selectedDoctor"
            @change="applyFilters"
            class="text-sm border border-gray-300 rounded-md px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Doctors</option>
            <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
              Dr. {{ doctor.first_name }} {{ doctor.last_name }}
            </option>
          </select>
        </div>

        <!-- Branch Filter -->
        <div class="filter-group" v-if="branches.length > 1">
          <select
            v-model="selectedBranch"
            @change="applyFilters"
            class="text-sm border border-gray-300 rounded-md px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Branches</option>
            <option v-for="branch in branches" :key="branch.id" :value="branch.id">
              {{ branch.name }}
            </option>
          </select>
        </div>

        <!-- Status Filter -->
        <div class="filter-group">
          <select
            v-model="selectedStatus"
            @change="applyFilters"
            class="text-sm border border-gray-300 rounded-md px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Status</option>
            <option value="scheduled">Scheduled</option>
            <option value="confirmed">Confirmed</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
            <option value="no_show">No Show</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Calendar Views -->
    <div class="calendar-container bg-white rounded-lg shadow-sm border">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading appointments...</span>
      </div>

      <!-- Day View -->
      <div v-else-if="currentView === 'day'" class="day-view">
         <div v-if="filteredAppointments.length === 0" class="text-center py-12 px-4 text-gray-500">
           <p>No appointments for this day</p>
           <p class="text-sm mt-2">Click on a time slot to create a new appointment</p>
         </div>
        <DayCalendar
          v-else
          :current-date="currentDate"
          :appointments="filteredAppointments"
          @appointment-click="viewAppointment"
          @time-slot-click="createAppointmentAtTime"
        />
      </div>

      <!-- Week View -->
      <div v-else-if="currentView === 'week'" class="week-view">
         <div v-if="filteredAppointments.length === 0" class="text-center py-12 px-4 text-gray-500">
           <p>No appointments for this week</p>
           <p class="text-sm mt-2">Click on a time slot to create a new appointment</p>
         </div>
        <WeekCalendar
          v-else
          :current-date="currentDate"
          :appointments="filteredAppointments"
          @appointment-click="viewAppointment"
          @time-slot-click="createAppointmentAtTime"
        />
      </div>

      <!-- Month View -->
      <div v-else-if="currentView === 'month'" class="month-view">
         <div v-if="filteredAppointments.length === 0" class="text-center py-12 px-4 text-gray-500">
           <p>No appointments for this month</p>
           <p class="text-sm mt-2">Click on a date to create a new appointment</p>
         </div>
        <MonthCalendar
          v-else
          :current-date="currentDate"
          :appointments="filteredAppointments"
          @appointment-click="viewAppointment"
          @date-click="selectDate"
        />
      </div>
    </div>

    <!-- Appointment Details Sidebar -->
    <div
      v-if="selectedAppointment"
      class="fixed right-0 top-0 h-full w-96 bg-white shadow-xl z-40 transform transition-transform"
      :class="selectedAppointment ? 'translate-x-0' : 'translate-x-full'"
    >
     <AppointmentDetails
         :appointment="selectedAppointment"
         @close="closeAppointmentDetails"
         @update="handleAppointmentUpdate"
         @edit="handleAppointmentEdit"
       />
    </div>


    <!-- Overlay when sidebar is open -->
    <div
      v-if="selectedAppointment"
      class="fixed inset-0 bg-black/25 z-30"
      @click="closeAppointmentDetails"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAppointmentStore } from '../stores/appointment'
import { useClinicStore } from '../stores/clinic'
import { useAuthStore } from '../stores/auth'

import DayCalendar from '../components/calendar/DayCalendar.vue'
import WeekCalendar from '../components/calendar/WeekCalendar.vue'
import MonthCalendar from '../components/calendar/MonthCalendar.vue'
import AppointmentDetails from '../components/AppointmentDetails.vue'

const router = useRouter()

const appointmentStore = useAppointmentStore()
const clinicStore = useClinicStore()
const authStore = useAuthStore()

// Reactive data
const currentView = ref('week')
const currentDate = ref(new Date())
const selectedDoctor = ref('')
const selectedBranch = ref('')
const selectedStatus = ref('')
const selectedAppointment = ref(null)

// Computed properties
const loading = computed(() => appointmentStore.isLoading)
const appointments = computed(() => appointmentStore.appointments)

const filteredAppointments = computed(() => {
  let filtered = appointments.value

  if (selectedDoctor.value) {
    filtered = filtered.filter(apt => apt.doctor_id === parseInt(selectedDoctor.value))
  }

  if (selectedBranch.value) {
    filtered = filtered.filter(apt => apt.branch_id === parseInt(selectedBranch.value))
  }

  if (selectedStatus.value) {
    filtered = filtered.filter(apt => apt.status === selectedStatus.value)
  }

  return filtered
})

const currentPeriodLabel = computed(() => {
  const date = currentDate.value
  
  switch (currentView.value) {
    case 'day':
      return date.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    case 'week':
      const startOfWeek = new Date(date)
      startOfWeek.setDate(date.getDate() - date.getDay())
      const endOfWeek = new Date(startOfWeek)
      endOfWeek.setDate(startOfWeek.getDate() + 6)
      
      if (startOfWeek.getMonth() === endOfWeek.getMonth()) {
        return `${startOfWeek.toLocaleDateString('en-US', { month: 'long', day: 'numeric' })} - ${endOfWeek.getDate()}, ${endOfWeek.getFullYear()}`
      } else {
        return `${startOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${endOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}, ${endOfWeek.getFullYear()}`
      }
    case 'month':
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long'
      })
    default:
      return ''
  }
})

const doctors = computed(() => {
  // This would come from a doctors store or clinic store
  return []
})

const branches = computed(() => {
  return clinicStore.branches || []
})

// Methods
const navigateCalendar = (direction) => {
  const newDate = new Date(currentDate.value)
  
  switch (currentView.value) {
    case 'day':
      newDate.setDate(newDate.getDate() + direction)
      break
    case 'week':
      newDate.setDate(newDate.getDate() + (direction * 7))
      break
    case 'month':
      newDate.setMonth(newDate.getMonth() + direction)
      break
  }
  
  currentDate.value = newDate
  loadAppointments()
}

const goToToday = () => {
  currentDate.value = new Date()
  loadAppointments()
}

const applyFilters = () => {
  // Filters are applied via computed property
  // This method can be used for additional filter logic if needed
}

const loadAppointments = async () => {
  const params = {}
  
  switch (currentView.value) {
    case 'day':
      params.date = currentDate.value.toISOString().split('T')[0]
      break
    case 'week':
      const startOfWeek = new Date(currentDate.value)
      startOfWeek.setDate(currentDate.value.getDate() - currentDate.value.getDay())
      const endOfWeek = new Date(startOfWeek)
      endOfWeek.setDate(startOfWeek.getDate() + 6)
      
      params.start_date = startOfWeek.toISOString().split('T')[0]
      params.end_date = endOfWeek.toISOString().split('T')[0]
      break
    case 'month':
      const startOfMonth = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth(), 1)
      const endOfMonth = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 0)
      
      params.start_date = startOfMonth.toISOString().split('T')[0]
      params.end_date = endOfMonth.toISOString().split('T')[0]
      break
  }
  
  await appointmentStore.fetchAppointments(params)
}

const viewAppointment = (appointment) => {
  selectedAppointment.value = appointment
}

const closeAppointmentDetails = () => {
  selectedAppointment.value = null
}

const createAppointmentAtTime = (date, hour, minute = 0) => {
  // Navigate to appointment form with date/time parameters
  const dateStr = date.toISOString().split('T')[0]
  const timeStr = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
  router.push(`/appointments/new?date=${dateStr}&time=${timeStr}`)
}

const selectDate = (date) => {
  currentDate.value = date
  currentView.value = 'day'
  loadAppointments()
}


const handleAppointmentUpdate = async (appointmentId, updateData) => {
  await appointmentStore.updateAppointment(appointmentId, updateData)
  await loadAppointments()

  // Update selected appointment if it's the one being updated
  if (selectedAppointment.value && selectedAppointment.value.id === appointmentId) {
    const updated = appointments.value.find(apt => apt.id === appointmentId)
    if (updated) {
      selectedAppointment.value = updated
    }
  }
}

const handleAppointmentEdit = (appointment) => {
  console.log('Editing appointment:', appointment)
  console.log('Appointment ID:', appointment.id)
  if (appointment && appointment.id) {
    router.push({ name: 'AppointmentEdit', params: { id: appointment.id } })
    closeAppointmentDetails()
  } else {
    console.error('Appointment or appointment ID is missing:', appointment)
  }
}


// Lifecycle
onMounted(async () => {
  await loadAppointments()
  
  // Load clinic data if needed
  if (authStore.userClinicId && !clinicStore.currentClinic) {
    await clinicStore.fetchClinic(authStore.userClinicId)
  }
})

// Watch for view or date changes
watch([currentView, currentDate], () => {
  loadAppointments()
})
</script>

<style scoped>
@import "../styles/main.css";

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.nav-btn:hover {
  @apply bg-gray-100;
}

.view-toggle button {
  @apply transition-all duration-200;
}

.calendar-container {
  min-height: 500px;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .page-header {
    @apply flex-col items-start gap-3;
  }

  .header-actions {
    @apply w-full flex-wrap gap-2;
  }

  .calendar-navigation {
    @apply flex-col gap-3 p-3;
  }

  .calendar-filters {
    @apply flex-wrap gap-2;
  }

  .calendar-container {
    @apply mx-2;
  }

  /* Sidebar on mobile takes full screen */
  .fixed.right-0.w-96 {
    @apply w-full;
  }
}

@media (max-width: 640px) {
  .page-header {
    @apply px-3 py-2;
  }

  .calendar-navigation {
    @apply p-2 gap-2;
  }

  .calendar-container {
    @apply mx-1;
  }

  .btn {
    @apply px-3 py-2 text-sm;
  }
}

/* Animation for sidebar */
.transform.transition-transform {
  transition: transform 0.3s ease-in-out;
}
</style>