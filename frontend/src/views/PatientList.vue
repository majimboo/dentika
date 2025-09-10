<template>
  <div class="patient-list-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
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
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          Add Patient
        </router-link>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <!-- Search and Filters -->
        <div class="search-and-filters flex flex-col md:flex-row gap-4 items-start md:items-center flex-1">
          <div class="search-box relative">
            <input
              v-model="searchQuery"
              @input="debouncedSearch"
              type="text"
              placeholder="Search patients by name, email, or phone..."
              class="w-full md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          </div>

          <div class="filters-row flex items-center space-x-4">
            <!-- Status Filter -->
            <select
              v-model="statusFilter"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Patients</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>

            <!-- Sort Options -->
            <select
              v-model="sortBy"
              @change="applyFilters"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="name">Sort by Name</option>
              <option value="created_at">Sort by Date Added</option>
              <option value="last_visit">Sort by Last Visit</option>
            </select>

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

        <!-- View Options -->
        <div class="view-options flex items-center space-x-2">
          <button
            @click="viewMode = 'grid'"
            :class="viewMode === 'grid' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
          >
            <font-awesome-icon icon="fa-solid fa-th" class="w-5 h-5" />
          </button>
          <button
            @click="viewMode = 'list'"
            :class="viewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
          >
            <font-awesome-icon icon="fa-solid fa-list" class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Patient Stats -->
    <div class="patient-stats grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-blue-600">{{ totalPatients }}</div>
        <div class="stat-label text-sm text-gray-600">Total Patients</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-green-600">{{ activePatients }}</div>
        <div class="stat-label text-sm text-gray-600">Active Patients</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-purple-600">{{ newPatientsThisMonth }}</div>
        <div class="stat-label text-sm text-gray-600">New This Month</div>
      </div>
      
      <div class="stat-card bg-white p-4 rounded-lg shadow-sm border border-gray-200">
        <div class="stat-value text-2xl font-bold text-orange-600">{{ upcomingAppointments }}</div>
        <div class="stat-label text-sm text-gray-600">Upcoming Appointments</div>
      </div>
    </div>

    <!-- Patient Grid/List View -->
    <!-- Loading State -->
    <div v-if="loading" class="loading-state flex items-center justify-center py-12 bg-white rounded-lg shadow-sm border border-gray-200">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Loading patients...</span>
    </div>

    <!-- Empty State -->
    <div v-else-if="patients.length === 0" class="empty-state text-center py-12 bg-white rounded-lg shadow-sm border border-gray-200">
      <div class="empty-icon text-gray-300 mb-4">
        <font-awesome-icon icon="fa-solid fa-users" class="w-16 h-16 mx-auto" />
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

    <!-- Patient Grid View -->
    <div v-else-if="viewMode === 'grid'" class="patients-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div
        v-for="patient in patients"
        :key="patient.id"
        class="patient-card group bg-white rounded-xl border border-gray-200 hover:border-blue-300 hover:shadow-lg transition-all duration-200 cursor-pointer overflow-hidden"
        @click="viewPatient(patient)"
      >
        <!-- Card Header with Avatar and Status -->
        <div class="relative bg-gradient-to-br from-blue-50 to-indigo-50 p-4 pb-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="relative">
                <div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-sm">
                  <span class="text-white font-semibold text-lg">
                    {{ getPatientInitials(patient) }}
                  </span>
                </div>
                <!-- Online status indicator -->
                <div class="absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-white" :class="patient.is_active ? 'bg-green-400' : 'bg-gray-400'"></div>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="font-semibold text-gray-900 truncate text-base leading-tight">
                  {{ patient.first_name }} {{ patient.last_name }}
                </h3>
                <p class="text-xs text-gray-500 font-medium">#{{ patient.patient_number }}</p>
              </div>
            </div>
            
            <!-- Quick actions dropdown -->
            <div class="patient-actions opacity-0 group-hover:opacity-100 transition-opacity">
              <div class="relative" @click.stop>
                <button
                  @click="togglePatientMenu(patient.id)"
                  class="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-white/50 transition-colors"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
                  </svg>
                </button>

                <!-- Dropdown Menu -->
                <div
                  v-if="activePatientMenu === patient.id"
                  class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg z-20 border border-gray-200 py-1"
                >
                  <button
                    @click="viewPatient(patient)"
                    class="flex items-center w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    <font-awesome-icon icon="fa-solid fa-eye" class="w-4 h-4 mr-3 text-gray-400" />
                    View Details
                  </button>
                  <button
                    v-if="canManagePatients"
                    @click="editPatient(patient)"
                    class="flex items-center w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    <font-awesome-icon icon="fa-solid fa-edit" class="w-4 h-4 mr-3 text-gray-400" />
                    Edit Patient
                  </button>
                  <button
                    @click="viewDentalChart(patient)"
                    class="flex items-center w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    <font-awesome-icon icon="fa-solid fa-tooth" class="w-4 h-4 mr-3 text-gray-400" />
                    Dental Chart
                  </button>
                  <button
                    @click="scheduleAppointment(patient)"
                    class="flex items-center w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    <font-awesome-icon icon="fa-solid fa-calendar-plus" class="w-4 h-4 mr-3 text-gray-400" />
                    Schedule
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Card Content -->
        <div class="p-4 pt-3">
          <!-- Key Info Grid -->
          <div class="grid grid-cols-2 gap-3 mb-3">
            <div class="bg-gray-50 rounded-lg p-2.5">
              <div class="flex items-center space-x-2">
                <div class="w-6 h-6 bg-blue-100 rounded-md flex items-center justify-center">
                  <font-awesome-icon icon="fa-solid fa-birthday-cake" class="w-3 h-3 text-blue-600" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 leading-none">Age</p>
                  <p class="text-sm font-semibold text-gray-900">{{ getPatientAge(patient) }}</p>
                </div>
              </div>
            </div>
            
            <div class="bg-gray-50 rounded-lg p-2.5">
              <div class="flex items-center space-x-2">
                <div class="w-6 h-6 bg-purple-100 rounded-md flex items-center justify-center">
                  <font-awesome-icon :icon="patient.gender === 'female' ? 'fa-solid fa-venus' : patient.gender === 'male' ? 'fa-solid fa-mars' : 'fa-solid fa-genderless'" class="w-3 h-3 text-purple-600" />
                </div>
                <div>
                  <p class="text-xs text-gray-500 leading-none">Gender</p>
                  <p class="text-sm font-semibold text-gray-900 capitalize">{{ patient.gender || 'N/A' }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Contact Info -->
          <div class="space-y-2">
            <div class="flex items-center text-sm text-gray-600">
              <font-awesome-icon icon="fa-solid fa-phone" class="w-4 h-4 mr-3 text-gray-400" />
              <span class="truncate">{{ patient.phone || 'No phone' }}</span>
            </div>
            <div class="flex items-center text-sm text-gray-600">
              <font-awesome-icon icon="fa-solid fa-envelope" class="w-4 h-4 mr-3 text-gray-400" />
              <span class="truncate">{{ patient.email || 'No email' }}</span>
            </div>
          </div>

          <!-- Quick Actions -->
          <div class="mt-4 pt-3 border-t border-gray-100">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <span class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium" :class="patient.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
                  <span class="w-1.5 h-1.5 rounded-full mr-1.5" :class="patient.is_active ? 'bg-green-400' : 'bg-red-400'"></span>
                  {{ patient.is_active ? 'Active' : 'Inactive' }}
                </span>
              </div>
              
              <div class="flex items-center space-x-1">
                <button
                  @click.stop="viewDentalChart(patient)"
                  class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                  title="Dental Chart"
                >
                  <font-awesome-icon icon="fa-solid fa-tooth" class="w-4 h-4" />
                </button>
                <button
                  @click.stop="scheduleAppointment(patient)"
                  class="p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-md transition-colors"
                  title="Schedule Appointment"
                >
                  <font-awesome-icon icon="fa-solid fa-calendar-plus" class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Patient List View -->
    <div v-else class="patients-list bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="patients-table">
        <table class="w-full">
          <thead class="bg-gray-50 border-b">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Patient</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Contact</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Age/Gender</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr
              v-for="patient in patients"
              :key="patient.id"
              class="hover:bg-gray-50 cursor-pointer"
              @click="viewPatient(patient)"
            >
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="patient-avatar w-8 h-8 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full flex items-center justify-center mr-3">
                    <span class="text-white font-medium">
                      {{ getPatientInitials(patient) }}
                    </span>
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ patient.first_name }} {{ patient.last_name }}</div>
                    <div class="text-sm text-gray-500">#{{ patient.patient_number }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ patient.phone || 'No phone' }}</div>
                <div class="text-sm text-gray-500">{{ patient.email || 'No email' }}</div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ getPatientAge(patient) }} years</div>
                <div class="text-sm text-gray-500">{{ patient.gender || 'Not specified' }}</div>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs px-2 py-1 rounded-full" :class="patient.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                  {{ patient.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center space-x-2">
                  <router-link
                    v-if="canManagePatients"
                    :to="`/patients/${patient.id}/edit`"
                    class="text-blue-600 hover:text-blue-700 text-sm font-medium"
                  >
                    Edit
                  </router-link>
                  <button
                    @click.stop="viewDentalChart(patient)"
                    class="text-green-600 hover:text-green-700 text-sm font-medium"
                  >
                    Dental Chart
                  </button>
                  <button
                    @click.stop="scheduleAppointment(patient)"
                    class="text-purple-600 hover:text-purple-700 text-sm font-medium"
                  >
                    Schedule
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
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
const viewMode = ref('grid') // 'grid' or 'list'

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

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
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

/* Patient card animations and effects */
.patient-card {
  transition: all 0.2s ease-in-out;
}

.patient-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px -12px rgba(0, 0, 0, 0.15), 0 8px 16px -8px rgba(0, 0, 0, 0.1);
}

.patient-card .bg-gradient-to-br {
  transition: all 0.3s ease-in-out;
}

.patient-card:hover .bg-gradient-to-br {
  background-size: 150% 150%;
}

/* Improved hover effects */
.patient-card:hover .w-12.h-12.bg-gradient-to-br {
  transform: scale(1.05) rotate(2deg);
  transition: all 0.2s ease-in-out;
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
  
  .patients-grid .grid {
    @apply gap-3;
  }
  
  .patient-card .p-4 {
    @apply p-3;
  }
  
  .patient-card .grid-cols-2 {
    @apply grid-cols-1;
  }
}
</style>