<template>
  <div class="fixed bottom-0 left-0 right-0 mobile-bottom-nav bg-white border-t border-gray-200 px-2 py-1 z-30">
    <div class="flex justify-around items-center">
      <!-- Clinic User Links -->
      <template v-if="!isSuperAdmin">
        <router-link
          to="/"
          class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
          :class="isActive('/') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
        >
          <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
          </svg>
          <span class="text-xs font-medium">Home</span>
        </router-link>
        <router-link
          to="/patients"
          class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
          :class="isActive('/patients') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
        >
          <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
          </svg>
          <span class="text-xs font-medium">Patients</span>
        </router-link>
        <router-link
          to="/appointments"
          class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
          :class="isActive('/appointments') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
        >
          <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <span class="text-xs font-medium">Schedule</span>
        </router-link>
      </template>

      <!-- Super Admin Links -->
      <template v-if="isSuperAdmin">
        <router-link
          to="/clinics"
          class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
          :class="isActive('/clinics') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
        >
          <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
          </svg>
          <span class="text-xs font-medium">Clinics</span>
        </router-link>
        <router-link
          to="/users"
          class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
          :class="isActive('/users') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
        >
          <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
          </svg>
          <span class="text-xs font-medium">Users</span>
        </router-link>
      </template>

      <!-- More/Profile -->
      <div
        @click="toggleMoreMenu"
        class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors cursor-pointer"
        :class="showMoreMenu ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
      >
        <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z"></path>
        </svg>
        <span class="text-xs font-medium">More</span>
      </div>
    </div>

    <!-- More Menu Overlay -->
    <div
      v-if="showMoreMenu"
      class="fixed bottom-16 left-0 right-0 bg-white border-t border-gray-200 shadow-lg z-50"
    >
      <div class="grid grid-cols-3 gap-2 p-4">
        <!-- Users (Admin only) -->
        <router-link
          v-if="canManageUsers"
          to="/users"
          @click="closeMoreMenu"
          class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <svg class="w-6 h-6 mb-1 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
          </svg>
          <span class="text-xs font-medium text-gray-700">Users</span>
        </router-link>

        <!-- Reports -->
        <div class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer opacity-50">
          <svg class="w-6 h-6 mb-1 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
          </svg>
          <span class="text-xs font-medium text-gray-700">Reports</span>
        </div>

        <!-- Settings -->
        <div class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer opacity-50">
          <svg class="w-6 h-6 mb-1 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
          </svg>
          <span class="text-xs font-medium text-gray-700">Settings</span>
        </div>

        <!-- Logout -->
        <div
          @click="logout"
          class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-red-50 transition-colors cursor-pointer"
        >
          <svg class="w-6 h-6 mb-1 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
          </svg>
          <span class="text-xs font-medium text-red-600">Logout</span>
        </div>
      </div>
    </div>

    <!-- Overlay to close more menu -->
    <div
      v-if="showMoreMenu"
      class="fixed inset-0 bg-black/75 z-40"
      @click="closeMoreMenu"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const showMoreMenu = ref(false)

// Computed properties
const userClinicId = computed(() => authStore.userClinicId)
const isSuperAdmin = computed(() => authStore.isSuperAdmin)

const canManageUsers = computed(() => {
  const role = authStore.user?.role
  return role === 'super_admin' || role === 'clinic_owner'
})

// Methods
const isActive = (path) => {
  const currentPath = route.path
  const currentName = route.name

  // Home is only active when exactly on the dashboard route
  if (path === '/') {
    return currentName === 'Dashboard' || currentPath === '/'
  }

  // For other paths, check if current path starts with the nav path
  if (path === '/patients') {
    return currentPath.startsWith('/patients')
  }

  if (path === '/appointments') {
    return currentPath.startsWith('/appointments')
  }

  if (path === '/clinics') {
    return currentPath.startsWith('/clinics')
  }

  if (path === '/users') {
    return currentPath.startsWith('/users')
  }

  return false
}

const toggleMoreMenu = () => {
  showMoreMenu.value = !showMoreMenu.value
}

const closeMoreMenu = () => {
  showMoreMenu.value = false
}

const logout = async () => {
  closeMoreMenu()
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
@import "../styles/main.css";

.mobile-bottom-nav {
  /* Fixed height to prevent layout shifts */
  height: 60px;
  /* Safe area for iOS devices */
  padding-bottom: env(safe-area-inset-bottom);
}

.nav-item {
  min-width: 60px;
  /* Touch target size */
  min-height: 44px;
}

/* Active state animations */
.nav-item.router-link-active {
  @apply text-blue-600 bg-blue-50;
}

/* Smooth transitions */
.nav-item svg,
.nav-item span {
  transition: all 0.2s ease-in-out;
}

/* Touch feedback */
@media (hover: none) {
  .nav-item:active {
    @apply scale-95;
    transition: transform 0.1s ease-in-out;
  }
}

/* More menu animation */
.fixed.bottom-16 {
  animation: slideUp 0.2s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>