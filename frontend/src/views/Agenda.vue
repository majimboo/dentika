<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Daily Agenda</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Manage today's appointments and patient flow</p>
      </div>

      <!-- Date Selector -->
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div class="relative">
          <input
            v-model="selectedDate"
            type="date"
            :min="minDate"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
            @change="loadAppointments"
          />
        </div>

        <button
          @click="setToday"
          class="inline-flex items-center justify-center px-4 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
          Today
        </button>

        <button
          @click="refreshAgenda"
          :disabled="loading"
          class="inline-flex items-center justify-center px-4 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
      <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-blue-800">Arriving</p>
            <p class="text-2xl font-bold text-blue-900">{{ getColumnCount('scheduled') }}</p>
          </div>
          <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-green-50 border border-green-200 rounded-xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-green-800">Arrived</p>
            <p class="text-2xl font-bold text-green-900">{{ getColumnCount('confirmed') }}</p>
          </div>
          <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-yellow-50 border border-yellow-200 rounded-xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-yellow-800">Started</p>
            <p class="text-2xl font-bold text-yellow-900">{{ getColumnCount('in_progress') }}</p>
          </div>
          <div class="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-purple-50 border border-purple-200 rounded-xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-purple-800">Finished</p>
            <p class="text-2xl font-bold text-purple-900">{{ getColumnCount('completed') }}</p>
          </div>
          <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-red-50 border border-red-200 rounded-xl p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-red-800">No Show</p>
            <p class="text-2xl font-bold text-red-900">{{ getColumnCount('cancelled') }}</p>
          </div>
          <div class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Kanban Board -->
    <div class="grid grid-cols-1 xl:grid-cols-5 gap-6">
      <!-- Arriving Column -->
      <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="bg-blue-50 px-6 py-4 border-b border-blue-200">
          <h3 class="text-lg font-semibold text-blue-900 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Arriving
            <span class="ml-2 bg-blue-100 text-blue-800 text-xs font-medium px-2 py-1 rounded-full">
              {{ getColumnCount('scheduled') }}
            </span>
          </h3>
        </div>
        <div class="p-4 space-y-3 min-h-[400px]">
          <AppointmentCard
            v-for="appointment in getAppointmentsByStatus('scheduled')"
            :key="appointment.id"
            :appointment="appointment"
            @status-changed="handleStatusChange"
          />
        </div>
      </div>

      <!-- Arrived Column -->
      <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="bg-green-50 px-6 py-4 border-b border-green-200">
          <h3 class="text-lg font-semibold text-green-900 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Arrived
            <span class="ml-2 bg-green-100 text-green-800 text-xs font-medium px-2 py-1 rounded-full">
              {{ getColumnCount('confirmed') }}
            </span>
          </h3>
        </div>
        <div class="p-4 space-y-3 min-h-[400px]">
          <AppointmentCard
            v-for="appointment in getAppointmentsByStatus('confirmed')"
            :key="appointment.id"
            :appointment="appointment"
            @status-changed="handleStatusChange"
          />
        </div>
      </div>

      <!-- Started Column -->
      <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="bg-yellow-50 px-6 py-4 border-b border-yellow-200">
          <h3 class="text-lg font-semibold text-yellow-900 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            Started
            <span class="ml-2 bg-yellow-100 text-yellow-800 text-xs font-medium px-2 py-1 rounded-full">
              {{ getColumnCount('in_progress') }}
            </span>
          </h3>
        </div>
        <div class="p-4 space-y-3 min-h-[400px]">
          <AppointmentCard
            v-for="appointment in getAppointmentsByStatus('in_progress')"
            :key="appointment.id"
            :appointment="appointment"
            @status-changed="handleStatusChange"
          />
        </div>
      </div>

      <!-- Finished Column -->
      <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="bg-purple-50 px-6 py-4 border-b border-purple-200">
          <h3 class="text-lg font-semibold text-purple-900 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Finished
            <span class="ml-2 bg-purple-100 text-purple-800 text-xs font-medium px-2 py-1 rounded-full">
              {{ getColumnCount('completed') }}
            </span>
          </h3>
        </div>
        <div class="p-4 space-y-3 min-h-[400px]">
          <AppointmentCard
            v-for="appointment in getAppointmentsByStatus('completed')"
            :key="appointment.id"
            :appointment="appointment"
            @status-changed="handleStatusChange"
          />
        </div>
      </div>

      <!-- No Show Column -->
      <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="bg-red-50 px-6 py-4 border-b border-red-200">
          <h3 class="text-lg font-semibold text-red-900 flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            No Show
            <span class="ml-2 bg-red-100 text-red-800 text-xs font-medium px-2 py-1 rounded-full">
              {{ getColumnCount('cancelled') }}
            </span>
          </h3>
        </div>
        <div class="p-4 space-y-3 min-h-[400px]">
          <AppointmentCard
            v-for="appointment in getAppointmentsByStatus('cancelled')"
            :key="appointment.id"
            :appointment="appointment"
            @status-changed="handleStatusChange"
          />
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <BaseTransition name="fade">
      <div v-if="!loading && appointments.length === 0" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">No appointments for {{ formatDate(selectedDate) }}</h3>
        <p class="text-neutral-600">There are no appointments scheduled for this date.</p>
      </div>
    </BaseTransition>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'
import BaseTransition from '../components/BaseTransition.vue'
import AppointmentCard from '../components/AppointmentCard.vue'

export default {
  name: 'Agenda',
  components: {
    BaseTransition,
    AppointmentCard
  },
  setup() {
    const authStore = useAuthStore()

    const loading = ref(true)
    const appointments = ref([])
    const selectedDate = ref(new Date().toISOString().split('T')[0])
    const minDate = ref(new Date().toISOString().split('T')[0])

    const loadAppointments = async () => {
      try {
        loading.value = true

        // Get appointments for the selected date
        const result = await apiService.getAppointments({
          date: selectedDate.value,
          limit: 100 // Get more appointments for the agenda
        })

        if (result.success) {
          appointments.value = result.data
        } else {
          console.error('Failed to load appointments:', result.error)
          appointments.value = []
        }
      } catch (err) {
        console.error('Error loading appointments:', err)
        appointments.value = []
      } finally {
        loading.value = false
      }
    }

    const refreshAgenda = async () => {
      await loadAppointments()
    }

    const setToday = () => {
      selectedDate.value = new Date().toISOString().split('T')[0]
      loadAppointments()
    }

    const getColumnCount = (status) => {
      return appointments.value.filter(apt => apt.status === status).length
    }

    const getAppointmentsByStatus = (status) => {
      return appointments.value.filter(apt => apt.status === status)
    }

    const handleStatusChange = async (appointmentId, newStatus) => {
      try {
        const result = await apiService.updateAppointmentStatus(appointmentId, {
          status: newStatus
        })

        if (result.success) {
          // Update the appointment status locally
          const appointment = appointments.value.find(apt => apt.id === appointmentId)
          if (appointment) {
            appointment.status = newStatus
          }
        } else {
          console.error('Failed to update appointment status:', result.error)
        }
      } catch (err) {
        console.error('Error updating appointment status:', err)
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'today'
      try {
        const date = new Date(dateString)
        const today = new Date()
        const tomorrow = new Date(today)
        tomorrow.setDate(today.getDate() + 1)
        const yesterday = new Date(today)
        yesterday.setDate(today.getDate() - 1)

        if (date.toDateString() === today.toDateString()) {
          return 'today'
        } else if (date.toDateString() === tomorrow.toDateString()) {
          return 'tomorrow'
        } else if (date.toDateString() === yesterday.toDateString()) {
          return 'yesterday'
        } else {
          return date.toLocaleDateString()
        }
      } catch {
        return dateString
      }
    }

    // Auto-refresh every 30 seconds
    const autoRefreshInterval = setInterval(() => {
      if (!loading.value) {
        loadAppointments()
      }
    }, 30000)

    onMounted(async () => {
      await loadAppointments()
    })

    // Cleanup interval on unmount
    const cleanup = () => {
      if (autoRefreshInterval) {
        clearInterval(autoRefreshInterval)
      }
    }

    return {
      loading,
      appointments,
      selectedDate,
      minDate,
      loadAppointments,
      refreshAgenda,
      setToday,
      getColumnCount,
      getAppointmentsByStatus,
      handleStatusChange,
      formatDate,
      cleanup
    }
  },
  beforeUnmount() {
    this.cleanup()
  }
}
</script>