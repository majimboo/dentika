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
            <font-awesome-icon icon="fa-solid fa-bars" class="w-5 h-5 sm:w-6 sm:h-6 text-gray-600" />
          </button>
          <div class="flex items-center space-x-2 sm:space-x-3">
            <div class="h-7 w-7 sm:h-8 sm:w-8 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-lg flex items-center justify-center">
            <font-awesome-icon icon="fa-solid fa-bolt" class="h-4 w-4 sm:h-5 sm:w-5 text-white" />
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
              <font-awesome-icon icon="fa-solid fa-bell" class="w-5 h-5 sm:w-6 sm:h-6 text-neutral-600" />
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
               <UserAvatar
                 :user="user"
                 size="sm"
                 :show-status="false"
                 v-if="user"
               />
            </div>
            
            <div class="relative">
              <button 
                @click="showUserMenu = !showUserMenu"
                class="p-1.5 sm:p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200"
              >
                <font-awesome-icon icon="fa-solid fa-chevron-down" class="w-4 h-4 sm:w-5 sm:h-5 text-neutral-600" />
              </button>
              
              <!-- User Menu Dropdown -->
              <div v-if="showUserMenu" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-10" role="menu" aria-label="User account menu">
                <a href="#" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors duration-200">
                  <font-awesome-icon icon="fa-solid fa-user" class="w-4 h-4 mr-3" />
                  Profile
                </a>
                <a href="#" class="flex items-center px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors duration-200">
                  <font-awesome-icon icon="fa-solid fa-sign-out-alt" class="w-4 h-4 mr-3" />
                  Settings
                </a>
                <div class="border-t border-neutral-100 my-1"></div>
                <button 
                  @click="logout" 
                  class="flex items-center w-full px-4 py-2 text-sm text-danger-700 hover:bg-danger-50 transition-colors duration-200"
                >
                  <font-awesome-icon icon="fa-solid fa-cog" class="w-4 h-4 mr-3" />
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
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSidebar } from '../composables/useSidebar'
import NotificationPanel from './NotificationPanel.vue'
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'AppHeader',
  components: {
    NotificationPanel,
    UserAvatar
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const showUserMenu = ref(false)
    const showNotifications = ref(false)
    const unreadNotifications = ref(2) // Mock unread count
    
    const user = computed(() => authStore.user)
    
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
      toggleSidebar,
      logout,
      toggleNotifications
    }
  }
}
</script>