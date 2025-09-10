<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Staff Management</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Manage members of your clinic team</p>
      </div>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
        <button
          @click="refreshStaff"
          :disabled="loading"
          class="inline-flex items-center justify-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
           <span class="sm:ml-2">Refresh</span>
        </button>
        <button
          @click="addNewStaff"
          class="inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
            <span class="sm:ml-2">Add Staff</span>
        </button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
             <p class="text-sm font-medium text-neutral-600">Total Staff</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ staff.length }}</p>
            <p class="text-xs text-success-600 mt-1 flex items-center">
              <font-awesome-icon icon="fa-solid fa-arrow-up" class="w-3 h-3 mr-1" />
              Active members
            </p>
          </div>
           <div class="w-12 h-12 bg-primary-100 rounded-xl flex items-center justify-center">
             <font-awesome-icon icon="fa-solid fa-users" class="w-6 h-6 text-primary-600" />
           </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Doctors</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ doctorsCount }}</p>
            <p class="text-xs text-blue-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
              </svg>
              Medical staff
            </p>
          </div>
           <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
             <font-awesome-icon icon="fa-solid fa-stethoscope" class="w-6 h-6 text-blue-600" />
           </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Secretaries</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ secretariesCount }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
              </svg>
              Administrative staff
            </p>
          </div>
           <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
             <font-awesome-icon icon="fa-solid fa-clipboard" class="w-6 h-6 text-green-600" />
           </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Assistants</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ assistantsCount }}</p>
            <p class="text-xs text-purple-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              Support staff
            </p>
          </div>
           <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
             <font-awesome-icon icon="fa-solid fa-bolt" class="w-6 h-6 text-purple-600" />
           </div>
        </div>
      </div>
    </div>

    <!-- Search and Filter Section -->
    <div class="bg-white rounded-2xl p-4 md:p-6 shadow-lg border border-neutral-100">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex-1 max-w-md">
          <div class="relative">
            <font-awesome-icon icon="fa-solid fa-search" class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-neutral-400" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search staff by name, role, or ID..."
              aria-label="Search staff"
              class="w-full pl-10 pr-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
            />
          </div>
        </div>
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-3">
          <select
            v-model="roleFilter"
            aria-label="Filter by role"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="">All Roles</option>
            <option value="doctor">Doctors</option>
            <option value="secretary">Secretaries</option>
            <option value="assistant">Assistants</option>
          </select>
          <select
            v-model="sortBy"
            aria-label="Sort staff by"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="username">Sort by Username</option>
            <option value="role">Sort by Role</option>
            <option value="id">Sort by ID</option>
            <option value="created_at">Sort by Created Date</option>
          </select>
          <button
            @click="toggleSortOrder"
            class="p-3 border border-neutral-300 rounded-xl text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
            :title="sortOrder === 'asc' ? 'Sort Ascending' : 'Sort Descending'"
          >
            <font-awesome-icon icon="fa-solid fa-sort" class="w-5 h-5" :class="{ 'rotate-180': sortOrder === 'desc' }" />
          </button>
        </div>
      </div>
    </div>

    <!-- Staff List -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <!-- Loading State -->
      <BaseTransition name="fade">
        <div v-if="loading" class="p-12 text-center">
          <BaseLoading
            type="spinner"
            size="large"
            color="primary"
            text="Loading staff..."
          />
          <p class="text-neutral-600 mt-4">Please wait while we fetch the staff data.</p>
        </div>
      </BaseTransition>

      <!-- Error State -->
      <BaseTransition name="fade">
        <div v-if="!loading && error" class="p-12 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading staff</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button
          @click="refreshStaff"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
        >
          Try Again
        </button>
        </div>
      </BaseTransition>

      <!-- Staff Grid -->
      <BaseTransition name="slide-up">
        <div v-if="!loading && !error && filteredStaff.length > 0" class="p-4 md:p-6">
          <TransitionGroup name="list" tag="div" class="grid grid-cols-1 xl:grid-cols-2 gap-4">
            <div v-for="member in filteredStaff" :key="member.id" class="group bg-neutral-50 hover:bg-white border border-neutral-200 hover:border-primary-300 rounded-xl p-6 transition-all duration-300 hover:shadow-lg hover:scale-[1.02]">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-4">
                    <div class="flex-shrink-0">
                      <div class="h-12 w-12 rounded-xl bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center text-white font-semibold text-lg">
                        {{ getInitials(member.username) }}
                      </div>
                    </div>
                    <div class="flex-1">
                      <h3 class="text-lg font-semibold text-neutral-900 group-hover:text-primary-700 transition-colors">
                        {{ member.username }}
                      </h3>
                      <div class="flex items-center space-x-4 mt-1">
                        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                              :class="getRoleBadgeClass(member.role)">
                          {{ getRoleDisplayName(member.role) }}
                        </span>
                        <span class="text-sm text-neutral-500">ID: {{ member.id }}</span>
                      </div>
                      <p class="text-xs text-neutral-400 mt-2 flex items-center">
                        <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        Joined {{ formatDate(member.created_at) }}
                      </p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity duration-200">
                    <router-link
                      :to="`/staff/${member.id}/edit`"
                      class="inline-flex items-center px-2 sm:px-3 py-2 border border-neutral-300 rounded-lg text-xs sm:text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
                      title="Edit staff member"
                    >
                      <svg class="w-4 h-4 sm:mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                      </svg>
                      <span class="hidden sm:inline">Edit</span>
                    </router-link>
                    <button
                      @click="confirmDelete(member)"
                      class="inline-flex items-center px-2 sm:px-3 py-2 border border-danger-300 rounded-lg text-xs sm:text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100 focus:outline-none focus:ring-2 focus:ring-danger-500 transition-all duration-200"
                      title="Remove staff member"
                    >
                      <svg class="w-4 h-4 sm:mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                      </svg>
                      <span class="hidden sm:inline">Remove</span>
                    </button>
                  </div>
                </div>
            </div>
          </TransitionGroup>
        </div>
      </BaseTransition>

      <!-- Empty State -->
      <BaseTransition name="fade">
        <div v-if="!loading && !error && filteredStaff.length === 0" class="p-12 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">
          {{ searchQuery || roleFilter ? 'No staff found' : 'No staff yet' }}
        </h3>
        <p class="text-neutral-600 mb-6">
          {{ searchQuery || roleFilter
            ? `No staff match your search criteria`
            : 'Get started by adding your first staff member.'
          }}
        </p>
        <button
          v-if="!searchQuery && !roleFilter"
          @click="addNewStaff"
          class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
           Add First Staff Member
        </button>
        </div>
      </BaseTransition>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseTransition name="modal">
      <div v-if="staffToDelete" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-2xl p-6 max-w-md w-full transform transition-all">
        <div class="flex items-center space-x-4 mb-4">
          <div class="w-12 h-12 bg-danger-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z"/>
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-neutral-900">Remove Staff Member</h3>
            <p class="text-sm text-neutral-600">This action cannot be undone</p>
          </div>
        </div>
        <p class="text-neutral-700 mb-6">
          Are you sure you want to remove <strong>{{ staffToDelete.username }}</strong> from your clinic staff?
          This will revoke their access to clinic systems and data.
        </p>
        <div class="flex items-center justify-end space-x-3">
          <button
            @click="staffToDelete = null"
            class="px-4 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 hover:bg-neutral-50 transition-colors"
          >
            Cancel
          </button>
          <button
            @click="removeStaff(staffToDelete.id)"
            :disabled="removing"
            class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 disabled:opacity-50 transition-colors"
          >
            <BaseLoading v-if="removing" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
            {{ removing ? 'Removing...' : 'Remove Staff' }}
          </button>
        </div>
        </div>
      </div>
    </BaseTransition>
  </div>
</template>

<script>
import { ref, computed, onMounted, TransitionGroup } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'

export default {
  name: 'StaffManagement',
  components: {
    BaseLoading,
    BaseTransition,
    TransitionGroup
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const staff = ref([])
    const loading = ref(true)
    const error = ref('')
    const searchQuery = ref('')
    const roleFilter = ref('')
    const sortBy = ref('username')
    const sortOrder = ref('asc')
    const staffToDelete = ref(null)
    const removing = ref(false)

    const filteredStaff = computed(() => {
      let filtered = staff.value

      // Apply search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(member =>
          member.username.toLowerCase().includes(query) ||
          member.role.toLowerCase().includes(query) ||
          member.id.toString().includes(query)
        )
      }

      // Apply role filter
      if (roleFilter.value) {
        filtered = filtered.filter(member => member.role === roleFilter.value)
      }

      // Apply sorting
      filtered.sort((a, b) => {
        let aVal = a[sortBy.value]
        let bVal = b[sortBy.value]

        if (typeof aVal === 'string') {
          aVal = aVal.toLowerCase()
          bVal = bVal.toLowerCase()
        }

        if (sortOrder.value === 'desc') {
          return aVal < bVal ? 1 : -1
        }
        return aVal > bVal ? 1 : -1
      })

        return filtered
      })

    const doctorsCount = computed(() => staff.value.filter(member => member.role === 'doctor').length)
    const secretariesCount = computed(() => staff.value.filter(member => member.role === 'secretary').length)
    const assistantsCount = computed(() => staff.value.filter(member => member.role === 'assistant').length)

    const loadStaff = async () => {
      try {
        loading.value = true
        error.value = ''
        const result = await apiService.getUsers()

        if (result.success) {
          // Filter to only show staff from current clinic
          const clinicId = authStore.userClinicId
          staff.value = result.data.filter(user => user.clinic_id === clinicId)
        } else {
          error.value = result.error
        }
      } catch (err) {
        error.value = 'Failed to load staff. Please check your connection and try again.'
        console.error('Error loading staff:', err)
      } finally {
        loading.value = false
      }
    }

    const refreshStaff = async () => {
      await loadStaff()
    }

    const addNewStaff = () => {
      router.push('/staff/new')
    }

    const confirmDelete = (member) => {
      staffToDelete.value = member
    }

    const removeStaff = async (staffId) => {
      try {
        removing.value = true
        const result = await apiService.deleteUser(staffId)

        if (result.success) {
          staff.value = staff.value.filter(member => member.id !== staffId)
          staffToDelete.value = null
        } else {
          error.value = result.error
        }
      } catch (err) {
        error.value = 'Failed to remove staff member. Please try again.'
        console.error('Error removing staff:', err)
      } finally {
        removing.value = false
      }
    }

    const toggleSortOrder = () => {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    }

    const getInitials = (username) => {
      if (!username) return 'U'
      return username.charAt(0).toUpperCase()
    }

    const getRoleDisplayName = (role) => {
      const roleNames = {
        'doctor': 'Doctor',
        'secretary': 'Secretary',
        'assistant': 'Assistant'
      }
      return roleNames[role] || role
    }

    const getRoleBadgeClass = (role) => {
      const classes = {
        'doctor': 'bg-blue-100 text-blue-800',
        'secretary': 'bg-green-100 text-green-800',
        'assistant': 'bg-purple-100 text-purple-800'
      }
      return classes[role] || 'bg-gray-100 text-gray-800'
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'Unknown'
      try {
        return new Date(dateString).toLocaleDateString()
      } catch {
        return 'Unknown'
      }
    }

    onMounted(() => {
      loadStaff()
    })

    return {
      staff,
      loading,
      error,
      searchQuery,
      roleFilter,
      sortBy,
      sortOrder,
      staffToDelete,
      removing,
      filteredStaff,
      doctorsCount,
      secretariesCount,
      assistantsCount,
      refreshStaff,
      addNewStaff,
      confirmDelete,
      removeStaff,
      toggleSortOrder,
      getInitials,
      getRoleDisplayName,
      getRoleBadgeClass,
      formatDate
    }
  }
}
</script>