<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Staff Management</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Manage members of your team</p>
      </div>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
        <button 
          @click="refreshUsers"
          :disabled="loading"
          class="inline-flex items-center justify-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
           <span class="sm:ml-2">Refresh</span>
        </button>
        <button
          v-if="canManageUsers"
          @click="addNewUser"
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
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
             <p class="text-sm font-medium text-neutral-600">Total Staff</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ users.length }}</p>
            <p class="text-xs text-success-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
              </svg>
              Active accounts
            </p>
          </div>
          <div class="w-12 h-12 bg-primary-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Recent Activity</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ recentActivity }}</p>
            <p class="text-xs text-warning-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              Last 24 hours
            </p>
          </div>
          <div class="w-12 h-12 bg-warning-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-warning-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">System Status</p>
            <p class="text-3xl font-bold text-success-600 mt-2">Online</p>
            <p class="text-xs text-success-600 mt-1 flex items-center">
              <span class="w-2 h-2 bg-success-400 rounded-full mr-2 animate-pulse"></span>
              All services operational
            </p>
          </div>
          <div class="w-12 h-12 bg-success-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Search and Filter Section -->
    <div class="bg-white rounded-2xl p-4 md:p-6 shadow-lg border border-neutral-100">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex-1 max-w-md">
          <div class="relative">
            <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search staff by name or ID..."
              aria-label="Search staff"
              class="w-full pl-10 pr-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
            />
          </div>
        </div>
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-3">
          <select 
            v-model="sortBy"
            aria-label="Sort users by"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="username">Sort by Username</option>
            <option value="id">Sort by ID</option>
            <option value="created_at">Sort by Created Date</option>
          </select>
          <button 
            @click="toggleSortOrder"
            class="p-3 border border-neutral-300 rounded-xl text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
            :title="sortOrder === 'asc' ? 'Sort Ascending' : 'Sort Descending'"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" :class="{ 'rotate-180': sortOrder === 'desc' }">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"/>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Users List -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <!-- Loading State -->
      <BaseTransition name="fade">
        <div v-if="loading" class="p-12 text-center">
          <BaseLoading 
            type="spinner" 
            size="large" 
            color="primary" 
            text="Loading users..." 
          />
          <p class="text-neutral-600 mt-4">Please wait while we fetch the user data.</p>
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
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading users</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button 
          @click="refreshUsers"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
        >
          Try Again
        </button>
        </div>
      </BaseTransition>
      
      <!-- Users Grid -->
      <BaseTransition name="slide-up">
        <div v-if="!loading && !error && filteredUsers.length > 0" class="p-4 md:p-6">
          <TransitionGroup name="list" tag="div" class="grid grid-cols-1 xl:grid-cols-2 gap-4">
            <div v-for="user in filteredUsers" :key="user.id" class="group bg-neutral-50 hover:bg-white border border-neutral-200 hover:border-primary-300 rounded-xl p-6 transition-all duration-300 hover:shadow-lg hover:scale-[1.02]">
                <div class="flex items-center justify-between">
                   <div class="flex items-center space-x-4">
                     <div class="flex-shrink-0">
                       <UserAvatar :user="user" size="lg" />
                     </div>
                    <div class="flex-1">
                      <h3 class="text-lg font-semibold text-neutral-900 group-hover:text-primary-700 transition-colors">
                        {{ user.username }}
                      </h3>
                      <div class="flex items-center space-x-4 mt-1">
                        <span class="text-sm text-neutral-500">ID: {{ user.id }}</span>
                        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-success-100 text-success-800">
                          Active
                        </span>
                      </div>
                      <p class="text-xs text-neutral-400 mt-2 flex items-center">
                        <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        Created {{ formatDate(user.created_at) }}
                      </p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity duration-200">
                    <router-link 
                      :to="`/users/${user.id}/edit`"
                      class="inline-flex items-center px-2 sm:px-3 py-2 border border-neutral-300 rounded-lg text-xs sm:text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
                      title="Edit user"
                    >
                      <svg class="w-4 h-4 sm:mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                      </svg>
                      <span class="hidden sm:inline">Edit</span>
                    </router-link>
                    <button 
                      @click="confirmDelete(user)"
                      class="inline-flex items-center px-2 sm:px-3 py-2 border border-danger-300 rounded-lg text-xs sm:text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100 focus:outline-none focus:ring-2 focus:ring-danger-500 transition-all duration-200"
                      title="Delete user"
                    >
                      <svg class="w-4 h-4 sm:mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                      </svg>
                      <span class="hidden sm:inline">Delete</span>
                    </button>
                  </div>
                </div>
            </div>
          </TransitionGroup>
        </div>
      </BaseTransition>
      
      <!-- Empty State -->
      <BaseTransition name="fade">
        <div v-if="!loading && !error && filteredUsers.length === 0" class="p-12 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">
          {{ searchQuery ? 'No staff found' : 'No staff yet' }}
        </h3>
        <p class="text-neutral-600 mb-6">
          {{ searchQuery
            ? `No staff match your search for "${searchQuery}"`
            : 'Get started by adding your first staff member.'
          }}
        </p>
        <button 
          v-if="!searchQuery"
          @click="addNewUser"
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

    <!-- Bulk Actions (when users are selected) -->
    <div v-if="selectedUsers.length > 0" class="fixed bottom-4 sm:bottom-6 left-4 right-4 sm:left-1/2 sm:right-auto sm:transform sm:-translate-x-1/2 bg-white rounded-2xl shadow-2xl border border-neutral-200 p-3 sm:p-4">
      <div class="flex flex-col sm:flex-row items-center space-y-2 sm:space-y-0 sm:space-x-4">
        <span class="text-sm font-medium text-neutral-900">{{ selectedUsers.length }} users selected</span>
        <div class="flex items-center space-x-2 w-full sm:w-auto">
          <button class="flex-1 sm:flex-none px-3 sm:px-4 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 hover:bg-neutral-50 transition-colors">
            Export
          </button>
          <button class="flex-1 sm:flex-none px-3 sm:px-4 py-2 border border-danger-300 rounded-lg text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100 transition-colors">
            Delete Selected
          </button>
          <button @click="selectedUsers = []" class="p-2 text-neutral-400 hover:text-neutral-600 transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseTransition name="modal">
      <div v-if="userToDelete" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-2xl p-6 max-w-md w-full transform transition-all">
        <div class="flex items-center space-x-4 mb-4">
          <div class="w-12 h-12 bg-danger-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z"/>
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-neutral-900">Delete User</h3>
            <p class="text-sm text-neutral-600">This action cannot be undone</p>
          </div>
        </div>
        <p class="text-neutral-700 mb-6">
          Are you sure you want to delete <strong>{{ userToDelete.username }}</strong>? 
          This will permanently remove the user account and all associated data.
        </p>
        <div class="flex items-center justify-end space-x-3">
          <button 
            @click="userToDelete = null"
            class="px-4 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 hover:bg-neutral-50 transition-colors"
          >
            Cancel
          </button>
          <button 
            @click="deleteUser(userToDelete.id)"
            :disabled="deleting"
            class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 disabled:opacity-50 transition-colors"
          >
            <BaseLoading v-if="deleting" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
            {{ deleting ? 'Deleting...' : 'Delete User' }}
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
import UserAvatar from '../components/UserAvatar.vue'

export default {
  name: 'UserList',
  components: {
    BaseLoading,
    BaseTransition,
    TransitionGroup,
    UserAvatar
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const users = ref([])
    const loading = ref(true)
    const error = ref('')
    const searchQuery = ref('')
    const sortBy = ref('username')
    const sortOrder = ref('asc')
    const selectedUsers = ref([])
    const userToDelete = ref(null)
    const deleting = ref(false)
    const recentActivity = ref(0)
    
    const filteredUsers = computed(() => {
      let filtered = users.value
      
      // Apply search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(user => 
          user.username.toLowerCase().includes(query) ||
          user.id.toString().includes(query)
        )
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

      const canManageUsers = computed(() => {
        const role = authStore.user?.role
        return ['super_admin', 'admin'].includes(role)
      })

     const loadUsers = async () => {
      try {
        loading.value = true
        error.value = ''
        const result = await apiService.getUsers()
        
        if (result.success) {
          users.value = result.data
          // Simulate recent activity
          recentActivity.value = Math.floor(Math.random() * 25) + 5
        } else {
          error.value = result.error
        }
      } catch (err) {
        error.value = 'Failed to load users. Please check your connection and try again.'
        console.error('Error loading users:', err)
      } finally {
        loading.value = false
      }
    }
    
    const refreshUsers = async () => {
      await loadUsers()
    }
    
    const addNewUser = () => {
      const clinicId = authStore.userClinicId
      if (clinicId) {
        router.push(`/users/new?clinic_id=${clinicId}`)
      } else {
        router.push('/users/new')
      }
    }
    
    const confirmDelete = (user) => {
      userToDelete.value = user
    }
    
    const deleteUser = async (userId) => {
      try {
        deleting.value = true
        const result = await apiService.deleteUser(userId)
        
        if (result.success) {
          users.value = users.value.filter(user => user.id !== userId)
          userToDelete.value = null
        } else {
          error.value = result.error
        }
      } catch (err) {
        error.value = 'Failed to delete user. Please try again.'
        console.error('Error deleting user:', err)
      } finally {
        deleting.value = false
      }
    }
    
    const toggleSortOrder = () => {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    }
    
    const getInitials = (username) => {
      if (!username) return 'U'
      return username.charAt(0).toUpperCase()
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
      loadUsers()
    })
    
    return {
      users,
      loading,
      error,
      searchQuery,
      sortBy,
      sortOrder,
      selectedUsers,
      userToDelete,
      deleting,
      recentActivity,
      filteredUsers,
      refreshUsers,
      addNewUser,
      confirmDelete,
      deleteUser,
      toggleSortOrder,
      getInitials,
      formatDate
    }
  }
}
</script>