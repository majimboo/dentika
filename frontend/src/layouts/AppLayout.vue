<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Mobile-First Layout -->
    
    <!-- Navigation Header - Always visible on mobile -->
    <NavigationHeader>
      <template #leftActions>
        <!-- Menu button for mobile -->
        <button
          @click="toggleSidebar"
          class="lg:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
          </svg>
        </button>
      </template>
      
      <template #actions>
        <!-- Mobile & Desktop notification/profile actions -->
        <div class="flex lg:hidden items-center space-x-2">
          <!-- Mobile Notifications -->
          <NotificationPanel />
        </div>
        
        <div class="hidden lg:flex items-center space-x-2">
          <!-- Desktop Notifications -->
          <NotificationPanel />
          
          <!-- Profile -->
          <button class="p-1 rounded-lg hover:bg-gray-100 transition-colors">
            <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
              <span class="text-white text-sm font-medium">{{ userInitials }}</span>
            </div>
          </button>
        </div>
      </template>
    </NavigationHeader>

    <!-- Mobile overlay for sidebar -->
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm z-30 lg:hidden"
      @click="closeSidebar"
      @touchstart="closeSidebar"
      @touchend.prevent
      @touchmove.prevent
      @touchcancel.prevent
      @contextmenu.prevent
      @selectstart.prevent
      style="pointer-events: auto; touch-action: none; user-select: none; -webkit-user-select: none; -webkit-touch-callout: none; -webkit-tap-highlight-color: transparent; position: fixed; top: 0; left: 0; right: 0; bottom: 0;"
    ></div>
    
    <!-- Sidebar - Slides in from left on mobile, always visible on desktop -->
    <SuperAdminSidebar v-if="isSuperAdmin" />
    <AdminSidebar v-else-if="isAdmin" />
    <DoctorSidebar v-else-if="isDoctor" />
    <SecretarySidebar v-else-if="isSecretary" />
    <AssistantSidebar v-else-if="isAssistant" />
    <AppSidebar v-else />
    
    <div class="flex flex-col min-h-screen lg:ml-64">
      <!-- Main Content -->
      <main class="flex-1 bg-white lg:bg-gray-50 pb-16 lg:pb-0">
        <!-- Mobile: Full width, no padding top (NavigationHeader handles it), bottom padding for mobile nav -->
        <!-- Mobile: Small padding, Desktop: Regular layout with padding -->
        <div class="p-4 lg:p-6 lg:pt-4">
          <div class="lg:max-w-7xl lg:mx-auto">
            <router-view />
          </div>
        </div>
      </main>
      
      <!-- Footer - Hidden on mobile, visible on desktop -->
      <AppFooter class="hidden lg:block" />
    </div>

    <!-- Mobile Bottom Navigation (Alternative to sidebar) -->
    <MobileBottomNav class="lg:hidden" />
  </div>
</template>

<script>
import { computed, watch, onUnmounted } from 'vue'
import { mapState } from 'pinia'
import NavigationHeader from '../components/NavigationHeader.vue'
import AppSidebar from '../components/AppSidebar.vue'
import SuperAdminSidebar from '../components/SuperAdminSidebar.vue'
import AdminSidebar from '../components/AdminSidebar.vue'
import DoctorSidebar from '../components/DoctorSidebar.vue'
import SecretarySidebar from '../components/SecretarySidebar.vue'
import AssistantSidebar from '../components/AssistantSidebar.vue'
import AppFooter from '../components/AppFooter.vue'
import MobileBottomNav from '../components/MobileBottomNav.vue'
import NotificationPanel from '../components/NotificationPanel.vue'
import { useSidebar } from '../composables/useSidebar'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'AppLayout',
  components: {
    NavigationHeader,
    AppSidebar,
    SuperAdminSidebar,
    AdminSidebar,
    DoctorSidebar,
    SecretarySidebar,
    AssistantSidebar,
    AppFooter,
    MobileBottomNav,
    NotificationPanel
  },
  computed: {
    ...mapState(useAuthStore, ['isSuperAdmin', 'isAdmin', 'isDoctor', 'isSecretary', 'isAssistant'])
  },
  setup() {
    const { isSidebarOpen, closeSidebar, toggleSidebar } = useSidebar()
    const authStore = useAuthStore()

    const userInitials = computed(() => {
      const user = authStore.user
      if (!user) return 'U'
      return `${user.first_name?.[0] || ''}${user.last_name?.[0] || ''}`.toUpperCase() || 'U'
    })

    // Handle body scroll lock for mobile sidebar
    const updateBodyScroll = () => {
      if (typeof window !== 'undefined') {
        const isMobile = window.innerWidth < 1024
        if (isSidebarOpen.value && isMobile) {
          document.body.style.overflow = 'hidden'
        } else {
          document.body.style.overflow = ''
        }
      }
    }

    // Watch for sidebar state changes
    watch(isSidebarOpen, updateBodyScroll)

    // Watch for window resize
    const handleResize = () => {
      updateBodyScroll()
    }

    if (typeof window !== 'undefined') {
      window.addEventListener('resize', handleResize)
    }

    // Cleanup
    onUnmounted(() => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('resize', handleResize)
        document.body.style.overflow = ''
      }
    })

    return {
      isSidebarOpen,
      closeSidebar,
      toggleSidebar,
      userInitials
    }
  }
}
</script>