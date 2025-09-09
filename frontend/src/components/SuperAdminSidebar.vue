<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 w-64 bg-white/90 backdrop-blur-md shadow-xl border-r border-neutral-200/50 transform transition-transform duration-300 lg:translate-x-0"
    :class="{ '-translate-x-full': !isSidebarOpen }"
    role="navigation"
    aria-label="Super Admin navigation"
  >
    <div class="flex flex-col h-full">
      <!-- Sidebar Header -->
      <div class="p-6 border-b border-neutral-200/50">
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-gradient-to-r from-red-600 to-orange-600 rounded-lg flex items-center justify-center">
            <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-bold bg-gradient-to-r from-neutral-900 to-neutral-600 bg-clip-text text-transparent">Dentika</h2>
            <p class="text-xs text-neutral-500">Super Admin Panel</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 p-6 space-y-2">
        <div class="mb-6">
          <h3 class="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">System Management</h3>
          <ul class="space-y-2" role="list">
            <li role="listitem">
              <router-link
                to="/"
                @click="closeSidebarOnMobile"
                class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-red-500"
                :class="{ 'bg-red-100 text-red-700 shadow-sm': $route.name === 'Dashboard' }"
                :aria-current="$route.name === 'Dashboard' ? 'page' : null"
              >
                <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-red-200 transition-colors duration-200"
                     :class="{ 'bg-red-200': $route.name === 'Dashboard' }">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium">System Dashboard</div>
                  <div class="text-xs text-gray-400">Global system overview</div>
                </div>
              </router-link>
            </li>
            <li role="listitem">
              <router-link
                to="/clinics"
                @click="closeSidebarOnMobile"
                class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-red-500"
                :class="{ 'bg-red-100 text-red-700 shadow-sm': $route.path.startsWith('/clinics') }"
                :aria-current="$route.path.startsWith('/clinics') ? 'page' : null"
              >
                <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-red-200 transition-colors duration-200"
                     :class="{ 'bg-red-200': $route.path.startsWith('/clinics') }">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
                  </svg>
                </div>
                <div>
                  <div class="font-medium">Clinics Management</div>
                  <div class="text-xs text-gray-400">Manage all clinics & branches</div>
                </div>
              </router-link>
            </li>
            <li role="listitem">
              <router-link
                to="/users"
                @click="closeSidebarOnMobile"
                class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-red-500"
                :class="{ 'bg-red-100 text-red-700 shadow-sm': $route.path.startsWith('/users') }"
                :aria-current="$route.path.startsWith('/users') ? 'page' : null"
              >
                <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-red-200 transition-colors duration-200"
                     :class="{ 'bg-red-200': $route.path.startsWith('/users') }">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium">Global Users</div>
                  <div class="text-xs text-gray-400">Manage all system users</div>
                </div>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Analytics Section -->
        <div class="mb-6">
          <h3 class="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Analytics</h3>
          <ul class="space-y-2" role="list">
            <li role="listitem">
              <router-link
                to="/analytics"
                @click="closeSidebarOnMobile"
                class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-red-50 hover:text-red-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-red-500"
                :class="{ 'bg-red-100 text-red-700 shadow-sm': $route.path.startsWith('/analytics') }"
                :aria-current="$route.path.startsWith('/analytics') ? 'page' : null"
              >
                <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-red-200 transition-colors duration-200"
                     :class="{ 'bg-red-200': $route.path.startsWith('/analytics') }">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                </div>
                <div>
                  <div class="font-medium">System Analytics</div>
                  <div class="text-xs text-gray-400">Global performance metrics</div>
                </div>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Quick Actions -->
        <div class="mb-6">
          <h3 class="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Quick Actions</h3>
          <div class="space-y-2">
            <BaseTooltip text="Create a new clinic in the system" position="right">
              <button class="w-full flex items-center px-4 py-2 text-sm text-neutral-600 rounded-lg hover:bg-neutral-50 transition-colors duration-200">
                <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                </svg>
                Add New Clinic
              </button>
            </BaseTooltip>
            <BaseTooltip text="Export system-wide data and reports" position="right">
              <button class="w-full flex items-center px-4 py-2 text-sm text-neutral-600 rounded-lg hover:bg-neutral-50 transition-colors duration-200">
                <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                </svg>
                Export System Data
              </button>
            </BaseTooltip>
          </div>
        </div>
      </nav>

      <!-- Help Section -->
      <div class="p-6 border-t border-gray-200/50">
        <div class="bg-gradient-to-r from-red-50 to-orange-50 rounded-xl p-4">
          <div class="flex items-center space-x-3">
            <div class="flex items-center justify-center w-8 h-8 bg-red-100 rounded-lg">
              <svg class="w-4 h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div>
              <h4 class="text-sm font-semibold text-gray-900">Super Admin Support</h4>
              <p class="text-xs text-gray-600">System administration help</p>
            </div>
          </div>
          <button class="w-full mt-3 px-3 py-2 bg-white rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors duration-200">
            Get Support
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>

<script>
import BaseTooltip from './BaseTooltip.vue'
import { useSidebar } from '../composables/useSidebar'

export default {
  name: 'SuperAdminSidebar',
  components: {
    BaseTooltip
  },
  setup() {
    const { isSidebarOpen, closeSidebar } = useSidebar()

    const closeSidebarOnMobile = () => {
      // Only close on mobile screens
      if (window.innerWidth < 1024) {
        closeSidebar()
      }
    }

    return {
      isSidebarOpen,
      closeSidebarOnMobile
    }
  }
}
</script>