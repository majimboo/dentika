<template>
  <header class="bg-white/80 backdrop-blur-md shadow-lg border-b border-neutral-200/50 sticky top-0 z-30">
    <div class="px-3 sm:px-6 py-3 sm:py-4">
      <div class="max-w-7xl mx-auto">
      <div class="flex items-center justify-between gap-2 sm:gap-4">
        <div class="flex items-center space-x-2 sm:space-x-4">
          <button 
            @click="toggleSidebar"
            class="lg:hidden p-1.5 sm:p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
            aria-label="Toggle navigation menu"
            aria-expanded="false"
          >
            <svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
            </svg>
          </button>
          <div class="flex items-center space-x-2 sm:space-x-3">
            <div class="h-7 w-7 sm:h-8 sm:w-8 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-lg flex items-center justify-center">
              <svg class="h-4 w-4 sm:h-5 sm:w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
            </div>
            <div class="hidden sm:block">
              <h1 class="text-lg sm:text-xl font-bold bg-gradient-to-r from-neutral-900 to-neutral-600 bg-clip-text text-transparent">Dentika</h1>
            </div>
          </div>
        </div>
        
        <div class="flex items-center space-x-1 sm:space-x-4">
          <!-- Notifications -->
          <div class="relative">
            <button 
              @click="toggleNotifications"
              class="p-1.5 sm:p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200 relative"
              aria-label="Notifications"
            >
              <svg class="w-5 h-5 sm:w-6 sm:h-6 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
              </svg>
              <span v-if="unreadNotifications > 0" class="absolute top-0.5 right-0.5 sm:top-1 sm:right-1 block h-2 w-2 rounded-full bg-danger-400"></span>
            </button>
          </div>
          
          <!-- User Menu -->
          <div class="flex items-center space-x-2 sm:space-x-3 pl-2 sm:pl-4 border-l border-neutral-200">
            <div class="flex items-center space-x-2 sm:space-x-3">
              <div class="text-right hidden md:block">
                <p class="text-sm font-semibold text-neutral-900">
                  {{ authStore.isUserLoading ? 'Loading...' : (user?.username || 'User') }}
                </p>
                <p class="text-xs text-neutral-500">Administrator</p>
              </div>
              <div class="h-8 w-8 sm:h-10 sm:w-10 rounded-full bg-gradient-to-r from-primary-500 to-secondary-600 flex items-center justify-center text-white font-semibold text-xs sm:text-sm">
                <div v-if="authStore.isUserLoading" class="animate-pulse">
                  <div class="w-2 h-2 bg-white rounded-full"></div>
                </div>
                <span v-else>{{ getInitials(user?.username || 'U') }}</span>
              </div>
            </div>
            
            <div class="relative">
              <button 
                @click="showUserMenu = !showUserMenu"
                class="p-1.5 sm:p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200"
              >
                <svg class="w-4 h-4 sm:w-5 sm:h-5 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </button>
              
              <!-- User Menu Dropdown -->
              <div v-if="showUserMenu" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-10" role="menu" aria-label="User account menu">
                <a href="#" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                  </svg>
                  Profile
                </a>
                <a href="#" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  </svg>
                  Settings
                </a>
                <div class="border-t border-neutral-100 my-1"></div>
                <button 
                  @click="logout" 
                  class="flex items-center w-full px-4 py-2 text-sm text-danger-700 hover:bg-danger-50 transition-colors duration-200"
                >
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
                  </svg>
                  Sign Out
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      </div>
    </div>
    
    <!-- Notification Panel -->
    <NotificationPanel 
      :is-open="showNotifications" 
      @close="showNotifications = false" 
    />
  </header>
</template>

<script>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSidebar } from '../composables/useSidebar'
import NotificationPanel from './NotificationPanel.vue'

export default {
  name: 'AppHeader',
  components: {
    NotificationPanel
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const showUserMenu = ref(false)
    const showNotifications = ref(false)
    const unreadNotifications = ref(2) // Mock unread count
    
    const user = computed(() => authStore.user)
    
    const getInitials = (username) => {
      if (!username) return 'U'
      return username.charAt(0).toUpperCase()
    }
    
    const { toggleSidebar } = useSidebar()
    
    const logout = () => {
      authStore.logout()
      router.push('/login')
    }
    
    const toggleNotifications = () => {
      showNotifications.value = !showNotifications.value
    }
    
    return {
      user,
      authStore,
      showUserMenu,
      showNotifications,
      unreadNotifications,
      getInitials,
      toggleSidebar,
      logout,
      toggleNotifications
    }
  }
}
</script>