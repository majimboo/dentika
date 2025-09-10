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
          <font-awesome-icon icon="fa-solid fa-home" class="w-6 h-6 mb-1" />
          <span class="text-xs font-medium">Home</span>
        </router-link>
         <router-link
           to="/patients"
           class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
           :class="isActive('/patients') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
         >
           <font-awesome-icon icon="fa-solid fa-users" class="w-6 h-6 mb-1" />
          <span class="text-xs font-medium">Patients</span>
        </router-link>
         <router-link
           to="/appointments"
           class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
           :class="isActive('/appointments') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
         >
           <font-awesome-icon icon="fa-solid fa-calendar-alt" class="w-6 h-6 mb-1" />
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
           <font-awesome-icon icon="fa-solid fa-building" class="w-6 h-6 mb-1" />
          <span class="text-xs font-medium">Clinics</span>
        </router-link>
         <router-link
           to="/users"
           class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors"
           :class="isActive('/users') ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
         >
           <font-awesome-icon icon="fa-solid fa-users" class="w-6 h-6 mb-1" />
          <span class="text-xs font-medium">Users</span>
        </router-link>
      </template>

      <!-- More/Profile -->
      <div
        @click="toggleMoreMenu"
        class="nav-item flex flex-col items-center py-2 px-3 rounded-lg transition-colors cursor-pointer"
        :class="showMoreMenu ? 'text-blue-600 bg-blue-50' : 'text-gray-600'"
      >
        <font-awesome-icon icon="fa-solid fa-ellipsis-h" class="w-6 h-6 mb-1" />
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
          <font-awesome-icon icon="fa-solid fa-users" class="w-6 h-6 mb-1 text-gray-600" />
          <span class="text-xs font-medium text-gray-700">Users</span>
        </router-link>

        <!-- Reports -->
        <div class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer opacity-50">
          <font-awesome-icon icon="fa-solid fa-chart-bar" class="w-6 h-6 mb-1 text-gray-600" />
          <span class="text-xs font-medium text-gray-700">Reports</span>
        </div>

        <!-- Settings -->
        <div class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer opacity-50">
          <font-awesome-icon icon="fa-solid fa-cog" class="w-6 h-6 mb-1 text-gray-600" />
          <span class="text-xs font-medium text-gray-700">Settings</span>
        </div>

        <!-- Logout -->
        <div
          @click="logout"
          class="flex flex-col items-center py-3 px-2 rounded-lg hover:bg-red-50 transition-colors cursor-pointer"
        >
          <font-awesome-icon icon="fa-solid fa-sign-out-alt" class="w-6 h-6 mb-1 text-red-600" />
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