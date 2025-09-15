<template>
  <div class="space-y-4">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Daily Agenda</h1>
        <p class="mt-1 text-sm text-gray-600">Manage today's appointments and patient flow</p>
      </div>

      <!-- Date Selector -->
      <div class="flex items-center gap-3">
        <div class="relative">
          <input v-model="selectedDate" type="date"
            class="px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all bg-white"
            @change="loadAppointments" />
        </div>

        <button @click="setToday"
          class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          Today
        </button>

        <button @click="refreshAgenda" :disabled="loading"
          class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 transition-all">
          <svg class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Kanban Board -->
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-3 lg:gap-0">
      <!-- Arriving Column -->
      <div class="bg-white shadow-sm border border-gray-200 overflow-hidden">
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 flex items-center justify-between">
            <div class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Arriving
            </div>
            <span class="bg-blue-100 text-blue-800 text-xs font-medium px-2 py-0.5 rounded-full">
              {{ getColumnCount('scheduled') }}
            </span>
          </h3>
        </div>
        <div class="p-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          <draggable
            v-model="scheduledAppointments"
            group="appointments"
            item-key="id"
            @change="handleDragChange"
            class="space-y-2 min-h-[250px]"
            :animation="200"
            ghost-class="dragging-ghost"
            chosen-class="dragging-chosen"
            drag-class="dragging-active"
            :touch-start-threshold="10"
            :force-fallback="false"
            :fallback-tolerance="3"
          >
            <template #item="{ element }">
              <AppointmentCard
                :appointment="element"
                @status-changed="handleStatusChange"
                class="mb-2 cursor-move"
              />
            </template>
            <template #footer>
              <div v-if="scheduledAppointments.length === 0"
                class="text-center py-8 text-gray-500 text-sm">
                No appointments
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <!-- Arrived Column -->
      <div class="bg-white shadow-sm border border-gray-200 overflow-hidden">
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 flex items-center justify-between">
            <div class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Arrived
            </div>
            <span class="bg-green-100 text-green-800 text-xs font-medium px-2 py-0.5 rounded-full">
              {{ getColumnCount('confirmed') }}
            </span>
          </h3>
        </div>
        <div class="p-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          <draggable
            v-model="confirmedAppointments"
            group="appointments"
            item-key="id"
            @change="handleDragChange"
            class="space-y-2 min-h-[250px]"
            :animation="200"
            ghost-class="dragging-ghost"
            chosen-class="dragging-chosen"
            drag-class="dragging-active"
            :touch-start-threshold="10"
            :force-fallback="false"
            :fallback-tolerance="3"
          >
            <template #item="{ element }">
              <AppointmentCard
                :appointment="element"
                @status-changed="handleStatusChange"
                class="mb-2 cursor-move"
              />
            </template>
            <template #footer>
              <div v-if="confirmedAppointments.length === 0"
                class="text-center py-8 text-gray-500 text-sm">
                No appointments
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <!-- Started Column -->
      <div class="bg-white shadow-sm border border-gray-200 overflow-hidden">
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 flex items-center justify-between">
            <div class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              Started
            </div>
            <span class="bg-yellow-100 text-yellow-800 text-xs font-medium px-2 py-0.5 rounded-full">
              {{ getColumnCount('in_progress') }}
            </span>
          </h3>
        </div>
        <div class="p-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          <draggable
            v-model="inProgressAppointments"
            group="appointments"
            item-key="id"
            @change="handleDragChange"
            class="space-y-2 min-h-[250px]"
            :animation="200"
            ghost-class="dragging-ghost"
            chosen-class="dragging-chosen"
            drag-class="dragging-active"
            :touch-start-threshold="10"
            :force-fallback="false"
            :fallback-tolerance="3"
          >
            <template #item="{ element }">
              <AppointmentCard
                :appointment="element"
                @status-changed="handleStatusChange"
                class="mb-2 cursor-move"
              />
            </template>
            <template #footer>
              <div v-if="inProgressAppointments.length === 0"
                class="text-center py-8 text-gray-500 text-sm">
                No appointments
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <!-- Finished Column -->
      <div class="bg-white shadow-sm border border-gray-200 overflow-hidden">
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 flex items-center justify-between">
            <div class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Finished
            </div>
            <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2 py-0.5 rounded-full">
              {{ getColumnCount('completed') }}
            </span>
          </h3>
        </div>
        <div class="p-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          <draggable
            v-model="completedAppointments"
            group="appointments"
            item-key="id"
            @change="handleDragChange"
            class="space-y-2 min-h-[250px]"
            :animation="200"
            ghost-class="dragging-ghost"
            chosen-class="dragging-chosen"
            drag-class="dragging-active"
            :touch-start-threshold="10"
            :force-fallback="false"
            :fallback-tolerance="3"
          >
            <template #item="{ element }">
              <AppointmentCard
                :appointment="element"
                @status-changed="handleStatusChange"
                class="mb-2 cursor-move"
              />
            </template>
            <template #footer>
              <div v-if="completedAppointments.length === 0"
                class="text-center py-8 text-gray-500 text-sm">
                No appointments
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <!-- No Show Column -->
      <div class="bg-white shadow-sm border border-gray-200 overflow-hidden">
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 flex items-center justify-between">
            <div class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              No Show
            </div>
            <span class="bg-red-100 text-red-800 text-xs font-medium px-2 py-0.5 rounded-full">
              {{ getColumnCount('cancelled') }}
            </span>
          </h3>
        </div>
        <div class="p-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          <draggable
            v-model="cancelledAppointments"
            group="appointments"
            item-key="id"
            @change="handleDragChange"
            class="space-y-2 min-h-[250px]"
            :animation="200"
            ghost-class="dragging-ghost"
            chosen-class="dragging-chosen"
            drag-class="dragging-active"
            :touch-start-threshold="10"
            :force-fallback="false"
            :fallback-tolerance="3"
          >
            <template #item="{ element }">
              <AppointmentCard
                :appointment="element"
                @status-changed="handleStatusChange"
                class="mb-2 cursor-move"
              />
            </template>
            <template #footer>
              <div v-if="cancelledAppointments.length === 0"
                class="text-center py-8 text-gray-500 text-sm">
                No appointments
              </div>
            </template>
          </draggable>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <BaseTransition name="fade">
      <div v-if="!loading && appointments.length === 0"
        class="bg-white rounded-lg p-8 shadow-sm border border-gray-200 text-center">
        <div class="inline-flex items-center justify-center w-12 h-12 bg-gray-100 rounded-lg mb-4">
          <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">No appointments for {{ formatDate(selectedDate, { relative: true, fallback: 'today' }) }}</h3>
        <p class="text-gray-600">There are no appointments scheduled for this date.</p>
      </div>
    </BaseTransition>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useAppointmentStore } from '../stores/appointment'
import { useTimezone } from '../composables/useTimezone'
import BaseTransition from '../components/BaseTransition.vue'
import AppointmentCard from '../components/AppointmentCard.vue'
import { formatDate } from '@/utils'
import draggable from 'vuedraggable'

export default {
  name: 'Agenda',
  components: {
    BaseTransition,
    AppointmentCard,
    draggable
  },
  setup() {
    const authStore = useAuthStore()
    const appointmentStore = useAppointmentStore()
    const { getCurrentTimeInTimezone } = useTimezone()

    const loading = computed(() => appointmentStore.isLoading)
    const appointments = computed(() => appointmentStore.appointments)
    const getLocalDateString = () => {
      const todayInTimezone = getCurrentTimeInTimezone()
      const year = todayInTimezone.getFullYear()
      const month = String(todayInTimezone.getMonth() + 1).padStart(2, '0')
      const day = String(todayInTimezone.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }

    const selectedDate = ref(getLocalDateString())

    const loadAppointments = async () => {
      await appointmentStore.fetchAppointments({
        date: selectedDate.value,
        limit: 100 // Get more appointments for the agenda
      })
    }

    const refreshAgenda = async () => {
      await loadAppointments()
    }

    const getColumnCount = (status) => {
      return appointments.value.filter(apt => apt.status === status).length
    }

    const getAppointmentsByStatus = (status) => {
      return appointments.value.filter(apt => apt.status === status)
    }

    // Computed properties for each column
    const scheduledAppointments = computed({
      get: () => appointments.value.filter(apt => apt.status === 'scheduled'),
      set: (value) => {
        updateAppointmentsStatus(value, 'scheduled')
      }
    })

    const confirmedAppointments = computed({
      get: () => appointments.value.filter(apt => apt.status === 'confirmed'),
      set: (value) => {
        updateAppointmentsStatus(value, 'confirmed')
      }
    })

    const inProgressAppointments = computed({
      get: () => appointments.value.filter(apt => apt.status === 'in_progress'),
      set: (value) => {
        updateAppointmentsStatus(value, 'in_progress')
      }
    })

    const completedAppointments = computed({
      get: () => appointments.value.filter(apt => apt.status === 'completed'),
      set: (value) => {
        updateAppointmentsStatus(value, 'completed')
      }
    })

    const cancelledAppointments = computed({
      get: () => appointments.value.filter(apt => apt.status === 'cancelled'),
      set: (value) => {
        updateAppointmentsStatus(value, 'cancelled')
      }
    })

    // Helper function to update appointment status when dragged
    const updateAppointmentsStatus = (appointments, newStatus) => {
      appointments.forEach(appointment => {
        if (appointment.status !== newStatus) {
          appointment.status = newStatus
        }
      })
    }

    const setToday = () => {
      selectedDate.value = getLocalDateString()
      loadAppointments()
    }

    const handleStatusChange = async (appointmentId, newStatus) => {
      try {
        const result = await appointmentStore.updateAppointmentStatus(appointmentId, {
          status: newStatus
        })

        if (result.success) {
          // After successful status update, refresh all appointments for the current date
          // to ensure patient information and other details are properly loaded
          await loadAppointments()
        } else {
          console.error('Failed to update appointment status:', result.error)
        }
      } catch (err) {
        console.error('Error updating appointment status:', err)
      }
    }

    // Handle drag and drop changes
    const handleDragChange = async (evt) => {
      if (evt.added) {
        // An appointment was dropped into this column
        const appointment = evt.added.element

        // The appointment status is already updated by the computed setter
        // We just need to sync it with the backend
        try {
          const result = await appointmentStore.updateAppointmentStatus(appointment.id, {
            status: appointment.status
          })

          if (!result.success) {
            console.error('Failed to update appointment status:', result.error)
            // Refresh to revert the change
            await loadAppointments()
          }
        } catch (err) {
          console.error('Error updating appointment status:', err)
          // Refresh to revert the change
          await loadAppointments()
        }
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
      loadAppointments,
      refreshAgenda,
      setToday,
      getColumnCount,
      getAppointmentsByStatus,
      handleStatusChange,
      handleDragChange,
      scheduledAppointments,
      confirmedAppointments,
      inProgressAppointments,
      completedAppointments,
      cancelledAppointments,
      cleanup,
      formatDate
    }
  },
  beforeUnmount() {
    this.cleanup()
  }
}
</script>

<style scoped>
@import "../styles/main.css";

.dragging-chosen {
  @apply shadow-lg transform rotate-2 z-10;
}

.dragging-ghost {
  @apply opacity-50;
}

.dragging-active {
  @apply opacity-75;
}
</style>