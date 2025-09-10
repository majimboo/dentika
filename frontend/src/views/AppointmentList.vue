<template>
  <div class="appointment-list-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Appointment List</h1>
        <p class="text-gray-600">View and manage all appointments</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <!-- View Options -->
        <div class="view-options flex items-center space-x-2 mr-4">
          <button
            @click="viewMode = 'grid'"
            :class="viewMode === 'grid' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
            title="Card View"
          >
            <font-awesome-icon icon="fa-solid fa-th" class="w-5 h-5" />
          </button>
          <button
            @click="viewMode = 'list'"
            :class="viewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
            title="Table View"
          >
            <font-awesome-icon icon="fa-solid fa-list" class="w-5 h-5" />
          </button>
        </div>

        <button
          @click="refreshList"
          :disabled="loading"
          class="inline-flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 transition-all"
        >
          <font-awesome-icon
            icon="fa-solid fa-sync"
            class="w-4 h-4 mr-2"
            :class="{ 'animate-spin': loading }"
          />
          Refresh
        </button>

        <router-link
          to="/appointments"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-calendar-alt" class="w-4 h-4 mr-2" />
          Calendar View
        </router-link>

        <router-link
          to="/appointments/new"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          New Appointment
        </router-link>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section bg-white rounded-lg shadow-sm border p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Search -->
        <div class="search-group">
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              class="block w-full pl-9 pr-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Search by patient name..."
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <font-awesome-icon icon="fa-solid fa-search" class="w-4 h-4 text-gray-400" />
            </div>
          </div>
        </div>

        <!-- Date Filter -->
        <div class="date-group">
          <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <select
            v-model="dateFilter"
            @change="applyFilters"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Dates</option>
            <option value="today">Today</option>
            <option value="tomorrow">Tomorrow</option>
            <option value="week">This Week</option>
            <option value="month">This Month</option>
          </select>
        </div>

        <!-- Status Filter -->
        <div class="status-group">
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="statusFilter"
            @change="applyFilters"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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

        <!-- Doctor Filter -->
        <div class="doctor-group">
          <label class="block text-sm font-medium text-gray-700 mb-1">Doctor</label>
          <select
            v-model="doctorFilter"
            @change="applyFilters"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Doctors</option>
            <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
              Dr. {{ doctor.first_name }} {{ doctor.last_name }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Appointments Grid View -->
    <div v-if="viewMode === 'grid'" class="appointments-grid">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state flex items-center justify-center py-12 bg-white rounded-lg shadow-sm border">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading appointments...</span>
      </div>

      <!-- Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="appointment in paginatedAppointments"
          :key="appointment.id"
          class="appointment-card bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow cursor-pointer"
          @click="$router.push(`/appointments/${appointment.id}`)"
        >
          <div class="p-6">
            <div class="flex items-start justify-between mb-4">
              <div class="patient-avatar w-12 h-12 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full flex items-center justify-center mr-4">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
              </div>
              <div class="appointment-actions">
                <div class="relative" @click.stop>
                  <button
                    @click="toggleStatusMenu(appointment.id)"
                    class="text-gray-400 hover:text-gray-600 transition-colors p-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
                    </svg>
                  </button>

                  <!-- Status Menu -->
                  <div
                    v-if="activeStatusMenu === appointment.id"
                    class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border"
                  >
                    <div class="py-1">
                      <button
                        v-if="appointment.status === 'scheduled'"
                        @click="updateAppointmentStatus(appointment.id, 'confirmed')"
                        class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >
                        Mark as Confirmed
                      </button>
                      <button
                        v-if="appointment.status === 'confirmed'"
                        @click="updateAppointmentStatus(appointment.id, 'in_progress')"
                        class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >
                        Start Appointment
                      </button>
                      <button
                        v-if="appointment.status === 'in_progress'"
                        @click="updateAppointmentStatus(appointment.id, 'completed')"
                        class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >
                        Complete Appointment
                      </button>
                      <button
                        v-if="appointment.status !== 'cancelled'"
                        @click="updateAppointmentStatus(appointment.id, 'cancelled')"
                        class="block w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-100"
                      >
                        Cancel Appointment
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <h3 class="text-lg font-semibold text-gray-900 mb-2">
              {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
            </h3>
            <p class="text-sm text-gray-600 mb-4">ID: {{ appointment.patient?.patient_number || appointment.patient_id }}</p>

            <div class="appointment-details space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Date</span>
                <span class="text-sm text-gray-900">{{ formatDate(appointment.start_time) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Time</span>
                <span class="text-sm text-gray-900">{{ formatTime(appointment.start_time) }} - {{ formatTime(appointment.end_time) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Doctor</span>
                <span class="text-sm text-gray-900">Dr. {{ appointment.doctor?.first_name }} {{ appointment.doctor?.last_name }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Status</span>
                <span class="text-xs px-2 py-1 rounded-full" :class="getStatusBadgeClass(appointment.status)">
                  {{ formatStatus(appointment.status) }}
                </span>
              </div>
            </div>

            <div class="mt-4 flex items-center justify-between">
              <router-link
                :to="`/appointments/${appointment.id}/edit`"
                class="text-indigo-600 hover:text-indigo-900 text-sm font-medium transition-colors"
                @click.stop
              >
                Edit
              </router-link>
              <router-link
                :to="`/appointments/${appointment.id}`"
                class="text-blue-600 hover:text-blue-900 text-sm font-medium transition-colors"
                @click.stop
              >
                View Details
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State for Grid -->
      <div v-if="!loading && filteredAppointments.length === 0" class="empty-state text-center py-12 bg-white rounded-lg shadow-sm border">
        <div class="empty-icon text-gray-300 mb-4">
          <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No appointments found</h3>
        <p class="text-gray-500 mb-4">
          {{ searchQuery || dateFilter || statusFilter || doctorFilter ? 'Try adjusting your filters.' : 'Get started by creating a new appointment.' }}
        </p>
        <router-link
          v-if="!searchQuery && !dateFilter && !statusFilter && !doctorFilter"
          to="/appointments/new"
          class="btn btn-primary"
        >
          Create First Appointment
        </router-link>
      </div>
    </div>

    <!-- Appointments Table View -->
    <div v-else class="table-container bg-white rounded-lg shadow-sm border overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading appointments...</span>
      </div>

      <!-- Table -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Patient
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date & Time
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Doctor
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="appointment in paginatedAppointments"
              :key="appointment.id"
              class="hover:bg-gray-50 transition-colors"
            >
              <!-- Patient Info -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <div class="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                      <svg class="h-5 w-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                      </svg>
                    </div>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">
                      {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
                    </div>
                    <div class="text-sm text-gray-500">
                      ID: {{ appointment.patient?.patient_number || appointment.patient_id }}
                    </div>
                  </div>
                </div>
              </td>

              <!-- Date & Time -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ formatDate(appointment.start_time) }}
                </div>
                <div class="text-sm text-gray-500">
                  {{ formatTime(appointment.start_time) }} - {{ formatTime(appointment.end_time) }}
                </div>
              </td>

              <!-- Doctor -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  Dr. {{ appointment.doctor?.first_name }} {{ appointment.doctor?.last_name }}
                </div>
              </td>

              <!-- Status -->
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="getStatusBadgeClass(appointment.status)"
                >
                  {{ formatStatus(appointment.status) }}
                </span>
              </td>

              <!-- Actions -->
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center space-x-2">
                  <router-link
                    :to="`/appointments/${appointment.id}`"
                    class="text-blue-600 hover:text-blue-900 transition-colors"
                    title="View Details"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </router-link>

                  <router-link
                    :to="`/appointments/${appointment.id}/edit`"
                    class="text-indigo-600 hover:text-indigo-900 transition-colors"
                    title="Edit Appointment"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                  </router-link>

                  <!-- Quick Status Actions -->
                  <div class="relative">
                    <button
                      @click="toggleStatusMenu(appointment.id)"
                      class="text-gray-400 hover:text-gray-600 transition-colors"
                      title="Change Status"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/>
                      </svg>
                    </button>

                    <!-- Status Menu -->
                    <div
                      v-if="activeStatusMenu === appointment.id"
                      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border"
                    >
                      <div class="py-1">
                        <button
                          v-if="appointment.status === 'scheduled'"
                          @click="updateAppointmentStatus(appointment.id, 'confirmed')"
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          Mark as Confirmed
                        </button>
                        <button
                          v-if="appointment.status === 'confirmed'"
                          @click="updateAppointmentStatus(appointment.id, 'in_progress')"
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          Start Appointment
                        </button>
                        <button
                          v-if="appointment.status === 'in_progress'"
                          @click="updateAppointmentStatus(appointment.id, 'completed')"
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          Complete Appointment
                        </button>
                        <button
                          v-if="appointment.status !== 'cancelled'"
                          @click="updateAppointmentStatus(appointment.id, 'cancelled')"
                          class="block w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-100"
                        >
                          Cancel Appointment
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State for Table -->
      <div v-if="!loading && filteredAppointments.length === 0" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No appointments found</h3>
        <p class="mt-1 text-sm text-gray-500">
          {{ searchQuery || dateFilter || statusFilter || doctorFilter ? 'Try adjusting your filters.' : 'Get started by creating a new appointment.' }}
        </p>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between bg-white px-4 py-3 border-t border-gray-200 sm:px-6 rounded-b-lg mt-6">
      <div class="flex-1 flex justify-between sm:hidden">
        <button
          @click="currentPage = Math.max(1, currentPage - 1)"
          :disabled="currentPage === 1"
          class="relative inline-flex items-center px-4 py-2 text-sm font-medium rounded-md text-gray-700 bg-white border border-gray-300 hover:bg-gray-50 disabled:opacity-50"
        >
          Previous
        </button>
        <button
          @click="currentPage = Math.min(totalPages, currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="ml-3 relative inline-flex items-center px-4 py-2 text-sm font-medium rounded-md text-gray-700 bg-white border border-gray-300 hover:bg-gray-50 disabled:opacity-50"
        >
          Next
        </button>
      </div>
      <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            Showing
            <span class="font-medium">{{ (currentPage - 1) * itemsPerPage + 1 }}</span>
            to
            <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, filteredAppointments.length) }}</span>
            of
            <span class="font-medium">{{ filteredAppointments.length }}</span>
            results
          </p>
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
            <button
              @click="currentPage = Math.max(1, currentPage - 1)"
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="currentPage = page"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                page === currentPage
                  ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                  : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
            <button
              @click="currentPage = Math.min(totalPages, currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAppointmentStore } from '../stores/appointment'
import { useClinicStore } from '../stores/clinic'
import { useAuthStore } from '../stores/auth'

const appointmentStore = useAppointmentStore()
const clinicStore = useClinicStore()
const authStore = useAuthStore()

// Reactive data
const loading = computed(() => appointmentStore.isLoading)
const appointments = computed(() => appointmentStore.appointments)

// View mode
const viewMode = ref('grid') // 'grid' or 'list'

// Filters
const searchQuery = ref('')
const dateFilter = ref('')
const statusFilter = ref('')
const doctorFilter = ref('')
const activeStatusMenu = ref(null)

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Computed properties
const filteredAppointments = computed(() => {
  let filtered = appointments.value

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(apt =>
      apt.patient?.first_name?.toLowerCase().includes(query) ||
      apt.patient?.last_name?.toLowerCase().includes(query) ||
      apt.patient?.patient_number?.toLowerCase().includes(query)
    )
  }

  // Date filter
  if (dateFilter.value) {
    const today = new Date()
    const tomorrow = new Date(today)
    tomorrow.setDate(today.getDate() + 1)

    filtered = filtered.filter(apt => {
      const aptDate = new Date(apt.start_time)
      switch (dateFilter.value) {
        case 'today':
          return aptDate.toDateString() === today.toDateString()
        case 'tomorrow':
          return aptDate.toDateString() === tomorrow.toDateString()
        case 'week':
          const weekStart = new Date(today)
          weekStart.setDate(today.getDate() - today.getDay())
          const weekEnd = new Date(weekStart)
          weekEnd.setDate(weekStart.getDate() + 6)
          return aptDate >= weekStart && aptDate <= weekEnd
        case 'month':
          return aptDate.getMonth() === today.getMonth() && aptDate.getFullYear() === today.getFullYear()
        default:
          return true
      }
    })
  }

  // Status filter
  if (statusFilter.value) {
    filtered = filtered.filter(apt => apt.status === statusFilter.value)
  }

  // Doctor filter
  if (doctorFilter.value) {
    filtered = filtered.filter(apt => apt.doctor_id === parseInt(doctorFilter.value))
  }

  return filtered
})

const totalPages = computed(() => Math.ceil(filteredAppointments.value.length / itemsPerPage.value))

const paginatedAppointments = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredAppointments.value.slice(start, end)
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages.filter(page => page !== '...')
})

const doctors = computed(() => {
  return clinicStore.doctors || []
})

// Methods
const loadAppointments = async () => {
  await appointmentStore.fetchAppointments()
}

const refreshList = async () => {
  await loadAppointments()
}

const applyFilters = () => {
  currentPage.value = 1 // Reset to first page when filters change
}

const toggleStatusMenu = (appointmentId) => {
  activeStatusMenu.value = activeStatusMenu.value === appointmentId ? null : appointmentId
}

const updateAppointmentStatus = async (appointmentId, newStatus) => {
  try {
    await appointmentStore.updateAppointmentStatus(appointmentId, { status: newStatus })
    activeStatusMenu.value = null
    await loadAppointments() // Refresh the list
  } catch (error) {
    console.error('Failed to update appointment status:', error)
  }
}

// Utility methods
const getStatusBadgeClass = (status) => {
  const classes = {
    scheduled: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-purple-100 text-purple-800',
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

const formatDate = (dateTime) => {
  const date = new Date(dateTime)
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatTime = (dateTime) => {
  return new Date(dateTime).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Close status menu when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    activeStatusMenu.value = null
  }
}

// Lifecycle
onMounted(async () => {
  document.addEventListener('click', handleClickOutside)
  await loadAppointments()

  // Load clinic data if needed
  if (authStore.userClinicId && !clinicStore.currentClinic) {
    await clinicStore.fetchClinic(authStore.userClinicId)
  }
})

// Watch for filter changes to reset pagination
watch([searchQuery, dateFilter, statusFilter, doctorFilter], () => {
  currentPage.value = 1
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

/* Custom scrollbar for table */
.table-container::-webkit-scrollbar {
  height: 8px;
}

.table-container::-webkit-scrollbar-track {
  @apply bg-gray-100;
}

.table-container::-webkit-scrollbar-thumb {
  @apply bg-gray-400 rounded;
}

.table-container::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-500;
}

/* Status menu animation */
.relative .absolute {
  animation: fadeIn 0.1s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Card animations and effects */
.appointment-card {
  transition: all 0.2s ease-in-out;
}

.appointment-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px -8px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.patient-avatar {
  transition: transform 0.2s ease-in-out;
}

.appointment-card:hover .patient-avatar {
  transform: scale(1.05);
}

/* View toggle buttons */
.view-options button {
  transition: all 0.2s ease-in-out;
}

.view-options button:hover {
  transform: scale(1.05);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .page-header {
    @apply flex-col items-start;
  }

  .header-actions {
    @apply flex-wrap gap-2;
  }

  .view-options {
    @apply mr-2;
  }

  .filters-section .grid {
    @apply grid-cols-1 gap-3;
  }

  .table-container {
    @apply overflow-x-scroll;
  }

  .appointments-grid .grid {
    @apply grid-cols-1 gap-4;
  }

  .appointment-card .p-6 {
    @apply p-4;
  }
}

@media (max-width: 640px) {
  .appointment-card .patient-avatar {
    @apply w-10 h-10;
  }

  .appointment-card h3 {
    @apply text-base;
  }

  .appointment-details .space-y-3 > * + * {
    @apply mt-2;
  }
}
</style>