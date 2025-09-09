<template>
  <div class="patient-list-page">
    <!-- Page Header -->
    <div class="page-header flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Patients</h1>
        <p class="text-gray-600">Manage patient records and information</p>
      </div>
      
      <div class="header-actions flex items-center space-x-3">
        <router-link
          v-if="canManagePatients"
          to="/patients/new"
          class="btn btn-primary flex items-center"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          Add Patient
        </router-link>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section bg-white rounded-lg shadow-sm border p-4 mb-6">
      <div class="flex flex-col md:flex-row gap-4 items-center">
        <!-- Search -->
        <div class="search-box flex-1">
          <div class="relative">
            <input
              v-model="searchQuery"
              @input="debouncedSearch"
              type="text"
              placeholder="Search patients by name, email, or phone..."
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <svg class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>

        <!-- Status Filter -->
        <div class="filter-group">
          <select
            v-model="statusFilter"
            @change="applyFilters"
            class="border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Patients</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>

        <!-- Sort Options -->
        <div class="filter-group">
          <select
            v-model="sortBy"
            @change="applyFilters"
            class="border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="name">Sort by Name</option>
            <option value="created_at">Sort by Date Added</option>
            <option value="last_visit">Sort by Last Visit</option>
          </select>
        </div>

        <!-- Clear Filters -->
        <button
          v-if="hasActiveFilters"
          @click="clearFilters"
          class="btn btn-secondary text-sm"
        >
          Clear Filters
        </button>
      </div>
    </div>

    <!-- Patient Stats -->
    <div class="patient-stats grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border">
        <div class="stat-value text-2xl font-bold text-blue-600">{{ totalPatients }}</div>
        <div class="stat-label text-sm text-gray-600">Total Patients</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border">
        <div class="stat-value text-2xl font-bold text-green-600">{{ activePatients }}</div>
        <div class="stat-label text-sm text-gray-600">Active Patients</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border">
        <div class="stat-value text-2xl font-bold text-purple-600">{{ newPatientsThisMonth }}</div>
        <div class="stat-label text-sm text-gray-600">New This Month</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border">
        <div class="stat-value text-2xl font-bold text-orange-600">{{ upcomingAppointments }}</div>
        <div class="stat-label text-sm text-gray-600">Upcoming Appointments</div>
      </div>
    </div>

    <!-- Patient List -->
    <div class="patient-list bg-white rounded-lg shadow-sm border">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading patients...</span>
      </div>

      <!-- Empty State -->
      <div v-else-if="patients.length === 0" class="empty-state text-center py-12">
        <div class="empty-icon text-gray-300 mb-4">
          <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No patients found</h3>
        <p class="text-gray-500 mb-4">
          {{ searchQuery ? 'No patients match your search criteria.' : 'Get started by adding your first patient.' }}
        </p>
        <router-link
          v-if="canManagePatients && !searchQuery"
          to="/patients/new"
          class="btn btn-primary"
        >
          Add First Patient
        </router-link>
      </div>

      <!-- Patient Table -->
      <div v-else class="patient-table">
        <div class="table-header grid grid-cols-12 gap-4 p-4 bg-gray-50 border-b text-sm font-medium text-gray-700">
          <div class="col-span-3">Patient</div>
          <div class="col-span-2">Contact</div>
          <div class="col-span-2">Age/Gender</div>
          <div class="col-span-2">Last Visit</div>
          <div class="col-span-2">Status</div>
          <div class="col-span-1">Actions</div>
        </div>

        <div class="table-body">
          <div
            v-for="patient in patients"
            :key="patient.id"
            class="patient-row grid grid-cols-12 gap-4 p-4 border-b hover:bg-gray-50 transition-colors cursor-pointer"
            @click="viewPatient(patient)"
          >
            <!-- Patient Info -->
            <div class="col-span-3 flex items-center">
              <div class="patient-avatar w-10 h-10 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full flex items-center justify-center mr-3">
                <span class="text-white font-medium">
                  {{ getPatientInitials(patient) }}
                </span>
              </div>
              <div>
                <div class="patient-name font-medium text-gray-900">
                  {{ patient.first_name }} {{ patient.last_name }}
                </div>
                <div class="patient-number text-xs text-gray-500">
                  #{{ patient.patient_number }}
                </div>
              </div>
            </div>

            <!-- Contact -->
            <div class="col-span-2">
              <div class="text-sm text-gray-900">{{ patient.phone || 'No phone' }}</div>
              <div class="text-xs text-gray-500">{{ patient.email || 'No email' }}</div>
            </div>

            <!-- Age/Gender -->
            <div class="col-span-2">
              <div class="text-sm text-gray-900">
                {{ getPatientAge(patient) }} years old
              </div>
              <div class="text-xs text-gray-500">{{ patient.gender || 'Not specified' }}</div>
            </div>

            <!-- Last Visit -->
            <div class="col-span-2">
              <div class="text-sm text-gray-900">
                {{ getLastVisit(patient) }}
              </div>
            </div>

            <!-- Status -->
            <div class="col-span-2">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="patient.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
              >
                {{ patient.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>

            <!-- Actions -->
            <div class="col-span-1 flex items-center">
              <div class="relative" @click.stop>
                <button
                  @click="togglePatientMenu(patient.id)"
                  class="text-gray-400 hover:text-gray-600 transition-colors"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
                  </svg>
                </button>

                <!-- Dropdown Menu -->
                <div
                  v-if="activePatientMenu === patient.id"
                  class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border"
                >
                  <div class="py-1">
                    <button
                      @click="viewPatient(patient)"
                      class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      View Details
                    </button>
                    <button
                      v-if="canManagePatients"
                      @click="editPatient(patient)"
                      class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Edit Patient
                    </button>
                    <button
                      @click="viewDentalChart(patient)"
                      class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Dental Chart
                    </button>
                    <button
                      @click="scheduleAppointment(patient)"
                      class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Schedule Appointment
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div class="pagination-section p-4 border-t bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="pagination-info text-sm text-gray-700">
              Showing {{ (currentPage - 1) * pageSize + 1 }} to {{ Math.min(currentPage * pageSize, totalPatients) }} of {{ totalPatients }} patients
            </div>
            
            <div class="pagination-controls flex items-center space-x-2">
              <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage <= 1"
                class="pagination-btn"
              >
                Previous
              </button>
              
              <span class="text-sm text-gray-700">
                Page {{ currentPage }} of {{ totalPages }}
              </span>
              
              <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage >= totalPages"
                class="pagination-btn"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { usePatientStore } from '../stores/patient'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const patientStore = usePatientStore()
const authStore = useAuthStore()
const router = useRouter()

// Reactive data
const searchQuery = ref('')
const statusFilter = ref('')
const sortBy = ref('name')
const activePatientMenu = ref(null)
const currentPage = ref(1)
const pageSize = ref(50)

// Computed properties
const loading = computed(() => patientStore.isLoading)
const patients = computed(() => patientStore.patients)
const totalPatients = computed(() => patientStore.getTotalPatients)
const canManagePatients = computed(() => authStore.canManagePatients)

const activePatients = computed(() => {
  return patients.value.filter(p => p.is_active).length
})

const newPatientsThisMonth = computed(() => {
  const thisMonth = new Date()
  thisMonth.setDate(1)
  return patients.value.filter(p => {
    return new Date(p.created_at) >= thisMonth
  }).length
})

const upcomingAppointments = computed(() => {
  // This would typically come from appointment store
  return 0 // Placeholder
})

const totalPages = computed(() => {
  return Math.ceil(totalPatients.value / pageSize.value)
})

const hasActiveFilters = computed(() => {
  return searchQuery.value || statusFilter.value || sortBy.value !== 'name'
})

// Methods
const loadPatients = async (params = {}) => {
  const searchParams = {
    page: currentPage.value,
    limit: pageSize.value,
    ...params
  }
  
  if (statusFilter.value) {
    searchParams.status = statusFilter.value
  }
  
  await patientStore.fetchPatients(searchParams)
}

const debouncedSearch = debounce(() => {
  if (searchQuery.value !== patientStore.searchQuery) {
    currentPage.value = 1
    patientStore.searchPatients(searchQuery.value)
  }
}, 300)

const applyFilters = () => {
  currentPage.value = 1
  loadPatients()
}

const clearFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  sortBy.value = 'name'
  currentPage.value = 1
  patientStore.clearSearch()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadPatients()
  }
}

const getPatientInitials = (patient) => {
  const first = patient.first_name?.charAt(0) || ''
  const last = patient.last_name?.charAt(0) || ''
  return (first + last).toUpperCase() || '?'
}

const getPatientAge = (patient) => {
  if (!patient.date_of_birth) return 'N/A'
  const birthDate = new Date(patient.date_of_birth)
  const today = new Date()
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--
  }
  return age
}

const getLastVisit = (patient) => {
  // This would come from appointments data
  return 'No visits'
}

const togglePatientMenu = (patientId) => {
  activePatientMenu.value = activePatientMenu.value === patientId ? null : patientId
}

const viewPatient = (patient) => {
  router.push(`/patients/${patient.id}`)
}

const editPatient = (patient) => {
  router.push(`/patients/${patient.id}/edit`)
  activePatientMenu.value = null
}

const viewDentalChart = (patient) => {
  router.push(`/patients/${patient.id}/dental-chart`)
}

const scheduleAppointment = (patient) => {
  router.push(`/appointments/new?patient=${patient.id}`)
}


// Utility function
function debounce(func, wait) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

// Lifecycle
onMounted(() => {
  loadPatients()
})

// Close dropdown when clicking outside
onMounted(() => {
  document.addEventListener('click', () => {
    activePatientMenu.value = null
  })
})
</script>

<style scoped>
@reference "tailwindcss";

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.pagination-btn {
  @apply px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed;
}

.patient-row {
  @apply cursor-pointer transition-colors;
}

.patient-row:hover {
  @apply bg-blue-50;
}

@media (max-width: 768px) {
  .table-header,
  .patient-row {
    @apply grid-cols-1 gap-2;
  }
  
  .table-header > div,
  .patient-row > div {
    @apply col-span-1;
  }
}
</style>