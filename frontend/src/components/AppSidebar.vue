<template>
  <aside 
    class="fixed inset-y-0 left-0 z-40 w-64 bg-white/90 backdrop-blur-md shadow-xl border-r border-neutral-200/50 transform transition-transform duration-300 lg:translate-x-0"
    :class="{ '-translate-x-full': !isSidebarOpen }"
    role="navigation"
    aria-label="Main navigation"
  >
    <div class="flex flex-col h-full">
      <!-- Sidebar Header -->
      <div class="p-6 border-b border-neutral-200/50">
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-lg flex items-center justify-center">
            <font-awesome-icon icon="fa-solid fa-bolt" class="h-5 w-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-bold bg-gradient-to-r from-neutral-900 to-neutral-600 bg-clip-text text-transparent">Dentika</h2>
            <p class="text-xs text-neutral-500">Management Panel</p>
          </div>
        </div>
      </div>

       <!-- Navigation -->
       <nav class="flex-1 p-6 space-y-2 overflow-y-auto">
        <div class="mb-6">
          <h3 class="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Main Menu</h3>
          <ul class="space-y-2" role="list">
             <li role="listitem">
               <router-link
                 to="/"
                 @click="closeSidebarOnMobile"
                 class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                 :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.name === 'Dashboard' }"
                 :aria-current="$route.name === 'Dashboard' ? 'page' : null"
               >
                 <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                      :class="{ 'bg-blue-200': $route.name === 'Dashboard' }">
                    <font-awesome-icon icon="fa-solid fa-home" class="w-5 h-5" />
                 </div>
                  <div>
                    <div class="font-medium">Dashboard</div>
                    <div class="text-xs text-gray-400">Overview & stats</div>
                  </div>
                </router-link>
              </li>
              <li v-if="!isSuperAdmin" role="listitem">
                <router-link
                  to="/agenda"
                  @click="closeSidebarOnMobile"
                  class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.name === 'Agenda' }"
                  :aria-current="$route.name === 'Agenda' ? 'page' : null"
                >
                  <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                       :class="{ 'bg-blue-200': $route.name === 'Agenda' }">
                     <font-awesome-icon icon="fa-solid fa-check-square" class="w-5 h-5" />
                  </div>
                  <div>
                    <div class="font-medium">Agenda</div>
                    <div class="text-xs text-gray-400">Daily workflow</div>
                  </div>
                </router-link>
              </li>
              <li role="listitem">
                <router-link
                  to="/patients"
                 @click="closeSidebarOnMobile"
                 class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                 :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.path.startsWith('/patients') }"
                 :aria-current="$route.path.startsWith('/patients') ? 'page' : null"
               >
                 <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                      :class="{ 'bg-blue-200': $route.path.startsWith('/patients') }">
                    <font-awesome-icon icon="fa-solid fa-users" class="w-5 h-5" />
                 </div>
                 <div>
                   <div class="font-medium">Patients</div>
                   <div class="text-xs text-gray-400">Manage patient records</div>
                 </div>
               </router-link>
             </li>
              <li role="listitem">
                <router-link
                  to="/appointments/list"
                  @click="closeSidebarOnMobile"
                  class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.path.startsWith('/appointments') }"
                  :aria-current="$route.path.startsWith('/appointments') ? 'page' : null"
                >
                 <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                      :class="{ 'bg-blue-200': $route.path.startsWith('/appointments') }">
                    <font-awesome-icon icon="fa-solid fa-calendar-alt" class="w-5 h-5" />
                 </div>
                 <div>
                   <div class="font-medium">Appointments</div>
                   <div class="text-xs text-gray-400">Schedule & calendar</div>
                 </div>
               </router-link>
             </li>
             <li role="listitem">
               <router-link
                 to="/procedures"
                 @click="closeSidebarOnMobile"
                 class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                 :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.path.startsWith('/procedures') || $route.path.startsWith('/diagnoses') || $route.path.startsWith('/treatments') || $route.path.startsWith('/consent') || $route.path.startsWith('/prescriptions') }"
                 :aria-current="$route.path.startsWith('/procedures') ? 'page' : null"
               >
                 <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                      :class="{ 'bg-blue-200': $route.path.startsWith('/procedures') || $route.path.startsWith('/diagnoses') || $route.path.startsWith('/treatments') || $route.path.startsWith('/consent') || $route.path.startsWith('/prescriptions') }">
                    <font-awesome-icon icon="fa-solid fa-stethoscope" class="w-5 h-5" />
                 </div>
                 <div>
                   <div class="font-medium">Procedures</div>
                   <div class="text-xs text-gray-400">Treatments & diagnoses</div>
                 </div>
               </router-link>
             </li>
             <li role="listitem" v-if="isClinicOwner">
               <router-link
                 to="/staff"
                 @click="closeSidebarOnMobile"
                 class="group flex items-center px-4 py-3 text-sm font-medium text-neutral-600 rounded-xl hover:bg-primary-50 hover:text-primary-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
                 :class="{ 'bg-primary-100 text-primary-700 shadow-sm': $route.path.startsWith('/staff') }"
                 :aria-current="$route.path.startsWith('/staff') ? 'page' : null"
               >
                 <div class="flex items-center justify-center w-10 h-10 mr-3 rounded-lg bg-gray-100 group-hover:bg-blue-200 transition-colors duration-200"
                      :class="{ 'bg-blue-200': $route.path.startsWith('/staff') }">
                    <font-awesome-icon icon="fa-solid fa-user-md" class="w-5 h-5" />
                 </div>
                   <div>
                     <div class="font-medium">Staff</div>
                     <div class="text-xs text-gray-400">Manage your team</div>
                   </div>
               </router-link>
             </li>
          </ul>
        </div>

        <!-- Quick Actions -->
        <div class="mb-6">
          <h3 class="text-xs font-semibold text-neutral-400 uppercase tracking-wider mb-3">Quick Actions</h3>
          <div class="space-y-2">
             <BaseTooltip text="Quickly add a new patient to the system" position="right">
               <button class="w-full flex items-center px-4 py-2 text-sm text-neutral-600 rounded-lg hover:bg-neutral-50 transition-colors duration-200">
                  <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-3" />
                 Add Patient
               </button>
             </BaseTooltip>
            <BaseTooltip text="Export patient and appointment data" position="right">
              <button class="w-full flex items-center px-4 py-2 text-sm text-neutral-600 rounded-lg hover:bg-neutral-50 transition-colors duration-200">
                 <font-awesome-icon icon="fa-solid fa-chart-bar" class="w-4 h-4 mr-3" />
                Export Data
              </button>
            </BaseTooltip>
          </div>
        </div>
      </nav>

      <!-- Help Section -->
      <div class="p-6 border-t border-gray-200/50">
        <div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-4">
          <div class="flex items-center space-x-3">
            <div class="flex items-center justify-center w-8 h-8 bg-blue-100 rounded-lg">
               <font-awesome-icon icon="fa-solid fa-question-circle" class="w-4 h-4 text-blue-600" />
            </div>
            <div>
              <h4 class="text-sm font-semibold text-gray-900">Need Help?</h4>
              <p class="text-xs text-gray-600">Access documentation</p>
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
import { useAuthStore } from '../stores/auth'
import { mapState } from 'pinia'

export default {
  name: 'AppSidebar',
  components: {
    BaseTooltip
  },
  computed: {
    ...mapState(useAuthStore, ['userClinicId', 'isSuperAdmin', 'isClinicOwner'])
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