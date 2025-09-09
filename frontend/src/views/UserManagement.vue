<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Global User Management</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Manage all users across the entire system</p>
      </div>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
        <button
          @click="refreshUsers"
          :disabled="loading"
          class="inline-flex items-center justify-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
           <span class="sm:ml-2">Refresh</span>
        </button>
        <button
          @click="addNewUser"
          class="inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-700 hover:to-orange-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
            <span class="sm:ml-2">Add User</span>
        </button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-6">
      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
             <p class="text-sm font-medium text-neutral-600">Total Users</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ users.length }}</p>
            <p class="text-xs text-success-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
              </svg>
              System wide
            </p>
          </div>
          <div class="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Super Admins</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ superAdminsCount }}</p>
            <p class="text-xs text-red-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              System admins
            </p>
          </div>
          <div class="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Clinic Owners</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ clinicOwnersCount }}</p>
            <p class="text-xs text-orange-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              Clinic managers
            </p>
          </div>
          <div class="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Clinic Staff</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ clinicStaffCount }}</p>
            <p class="text-xs text-blue-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
              </svg>
              Doctors & staff
            </p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 shadow-lg border border-neutral-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-neutral-600">Clinics</p>
            <p class="text-3xl font-bold text-neutral-900 mt-2">{{ clinicsCount }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              Active clinics
            </p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
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
              placeholder="Search users by name, role, clinic, or ID..."
              aria-label="Search users"
              class="w-full pl-10 pr-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
            />
          </div>
        </div>
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-3">
          <select
            v-model="roleFilter"
            aria-label="Filter by role"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="">All Roles</option>
            <option value="super_admin">Super Admins</option>
            <option value="clinic_owner">Clinic Owners</option>
            <option value="doctor">Doctors</option>
            <option value="secretary">Secretaries</option>
            <option value="assistant">Assistants</option>
          </select>
          <select
            v-model="clinicFilter"
            aria-label="Filter by clinic"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="">All Clinics</option>
            <option v-for="clinic in clinics" :key="clinic.id" :value="clinic.id">
              {{ clinic.name }}
            </option>
          </select>
          <select
            v-model="sortBy"
            aria-label="Sort users by"
            class="px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
          >
            <option value="username">Sort by Username</option>
            <option value="role">Sort by Role</option>
            <option value="clinic_name">Sort by Clinic</option>
            <option value="id">Sort by ID</option>
            <option value="created_at">Sort by Created Date</option>
          </select>
          <button
            @click="toggleSortOrder"
            class="p-3 border border-neutral-300 rounded-xl text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-red-500 transition-all duration-200"
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
            color="red"
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
            <div v-for="user in filteredUsers" :key="user.id" class="group bg-neutral-50 hover:bg-white border border-neutral-200 hover:border-red-300 rounded-xl p-6 transition-all duration-300 hover:shadow-lg hover:scale-[1.02]">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-4">
                    <div class="flex-shrink-0">
                      <div class="h-12 w-12 rounded-xl bg-gradient-to-r from-red-500 to-orange-500 flex items-center justify-center text-white font-semibold text-lg">
                        {{ getInitials(user.username) }}
                      </div>
                    </div>
                    <div class="flex-1">
                      <h3 class="text-lg font-semibold text-neutral-900 group-hover:text-red-700 transition-colors">
                        {{ user.username }}
                      </h3>
                      <div class="flex items-center space-x-4 mt-1">
                        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                              :class="getRoleBadgeClass(user.role)">
                          {{ getRoleDisplayName(user.role) }}
                        </span>
                        <span class="text-sm text-neutral-500">ID: {{ user.id }}</span>
                      </div>
                      <div class="flex items-center space-x-2 mt-2">
                        <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                          <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
                          </svg>
                          {{ user.clinic_name || 'No Clinic' }}
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
                      class="inline-flex items-center px-2 sm:px-3 py-2 border border-neutral-300 rounded-lg text-xs sm:text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-red-500 transition-all duration-200"
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
          {{ searchQuery || roleFilter || clinicFilter ? 'No users found' : 'No users yet' }}
        </h3>
        <p class="text-neutral-600 mb-6">
          {{ searchQuery || roleFilter || clinicFilter
            ? `No users match your search criteria`
            : 'Get started by adding your first user to the system.'
          }}
        </p>
        <button
          v-if="!searchQuery && !roleFilter && !clinicFilter"
          @click="addNewUser"
          class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-700 hover:to-orange-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
           Add First User
        </button>
        </div>
      </BaseTransition>
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
          This will permanently remove the user account and all associated data from the system.
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

export default {
  name: 'UserManagement',
  components: {
    BaseLoading,
    BaseTransition,
    TransitionGroup
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const users = ref([])
    const clinics = ref([])
    const loading = ref(true)
    const error = ref('')
    const searchQuery = ref('')
    const roleFilter = ref('')
    const clinicFilter = ref('')
    const sortBy = ref('username')
    const sortOrder = ref('asc')
    const userToDelete = ref(null)
    const deleting = ref(false)

    const filteredUsers = computed(() => {
      let filtered = users.value.map(user => ({
        ...user,
        clinic_name: clinics.value.find(c => c.id === user.clinic_id)?.name || 'No Clinic'
      }))

      // Apply search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(user =>
          user.username.toLowerCase().includes(query) ||
          user.role.toLowerCase().includes(query) ||
          user.clinic_name.toLowerCase().includes(query) ||
          user.id.toString().includes(query)
        )
      }

      // Apply role filter
      if (roleFilter.value) {
        filtered = filtered.filter(user => user.role === roleFilter.value)
      }

      // Apply clinic filter
      if (clinicFilter.value) {
        filtered = filtered.filter(user => user.clinic_id === parseInt(clinicFilter.value))
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

    const superAdminsCount = computed(() => users.value.filter(user => user.role === 'super_admin').length)
    const clinicOwnersCount = computed(() => users.value.filter(user => user.role === 'clinic_owner').length)
    const clinicStaffCount = computed(() => users.value.filter(user => ['doctor', 'secretary', 'assistant'].includes(user.role)).length)
    const clinicsCount = computed(() => clinics.value.length)

    const loadUsers = async () => {
      try {
        loading.value = true
        error.value = ''
        const result = await apiService.getUsers()

        if (result.success) {
          users.value = result.data
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

    const loadClinics = async () => {
      try {
        const result = await apiService.getClinics()
        if (result.success) {
          clinics.value = result.data
        }
      } catch (err) {
        console.error('Error loading clinics:', err)
      }
    }

    const refreshUsers = async () => {
      await loadUsers()
    }

    const addNewUser = () => {
      router.push('/users/new')
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

    const getRoleDisplayName = (role) => {
      const roleNames = {
        'super_admin': 'Super Admin',
        'clinic_owner': 'Clinic Owner',
        'doctor': 'Doctor',
        'secretary': 'Secretary',
        'assistant': 'Assistant'
      }
      return roleNames[role] || role
    }

    const getRoleBadgeClass = (role) => {
      const classes = {
        'super_admin': 'bg-red-100 text-red-800',
        'clinic_owner': 'bg-orange-100 text-orange-800',
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

    onMounted(async () => {
      await Promise.all([loadUsers(), loadClinics()])
    })

    return {
      users,
      clinics,
      loading,
      error,
      searchQuery,
      roleFilter,
      clinicFilter,
      sortBy,
      sortOrder,
      userToDelete,
      deleting,
      filteredUsers,
      superAdminsCount,
      clinicOwnersCount,
      clinicStaffCount,
      clinicsCount,
      refreshUsers,
      addNewUser,
      confirmDelete,
      deleteUser,
      toggleSortOrder,
      getInitials,
      getRoleDisplayName,
      getRoleBadgeClass,
      formatDate
    }
  }
}
</script>