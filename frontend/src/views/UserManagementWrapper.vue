<template>
  <div>
    <!-- Super Admin View -->
    <UserManagement v-if="isSuperAdmin" />

    <!-- Clinic Owner/Staff View -->
    <StaffManagement v-else-if="isClinicOwner" />

    <!-- Fallback for other roles (shouldn't happen but just in case) -->
    <div v-else class="p-12 text-center">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
        <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Access Denied</h3>
      <p class="text-neutral-600">You don't have permission to access user management.</p>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import UserManagement from './UserManagement.vue'
import StaffManagement from './StaffManagement.vue'

export default {
  name: 'UserManagementWrapper',
  components: {
    UserManagement,
    StaffManagement
  },
  setup() {
    const authStore = useAuthStore()

    const isSuperAdmin = computed(() => authStore.user?.role === 'super_admin')
    const isClinicOwner = computed(() => authStore.user?.role === 'clinic_owner')

    return {
      isSuperAdmin,
      isClinicOwner
    }
  }
}
</script>