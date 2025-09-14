<template>
  <div class="space-y-8">
    <!-- Loading State -->
    <BaseTransition name="fade">
      <div v-if="loading" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <BaseLoading type="spinner" size="large" color="primary" text="Loading staff data..." />
        <p class="text-neutral-600 mt-4">Please wait while we fetch the staff information.</p>
      </div>
    </BaseTransition>

    <!-- Error State -->
    <BaseTransition name="fade">
      <div v-if="!loading && error" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
         <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading staff member</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button
          @click="loadUser"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
        >
          Try Again
        </button>
      </div>
    </BaseTransition>

    <!-- Edit Form -->
    <BaseTransition name="slide-up">
      <div v-if="!loading && !error && originalUser" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="p-6 sm:p-8">
          <form @submit.prevent="handleSubmit" class="space-y-8">

            <!-- Avatar Section -->
            <div class="pb-6 border-b border-neutral-200">
              <h3 class="text-lg font-semibold text-neutral-900 mb-4">Profile Picture</h3>
              <AvatarUpload
                :user="form"
                entity-type="user"
                :entity-id="form.id || 0"
                @avatar-updated="handleAvatarUpdated"
              />
            </div>

            <!-- Basic Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Basic Information</h3>

               <!-- Username -->
               <div class="space-y-2">
                 <label for="username" class="block text-sm font-semibold text-gray-700 flex items-center">
                   <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                   </svg>
                   Username
                   <span class="text-danger-500 ml-1">*</span>
                 </label>
                 <input
                   v-model="form.username"
                   id="username"
                   name="username"
                   type="text"
                   required
                   class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                   placeholder="Enter username"
                 />
               </div>


               <!-- Role (for admins in edit mode) -->
                <div v-if="currentUserRole === 'admin' || currentUserRole === 'super_admin'" class="space-y-2">
                 <label for="role" class="block text-sm font-semibold text-gray-700 flex items-center">
                   <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                   </svg>
                   Role
                   <span class="text-danger-500 ml-1">*</span>
                 </label>
                 <select
                   v-model="form.role"
                   id="role"
                   name="role"
                   required
                   class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                 >
                   <option value="">Select role</option>
                   <!-- Super Admin can see all roles -->
                   <option v-if="currentUserRole === 'super_admin'" value="super_admin">Super Admin</option>
                   <option v-if="currentUserRole === 'super_admin'" value="admin">Admin</option>
                   <!-- Admin can only assign staff roles -->
                   <option v-if="currentUserRole === 'super_admin' || currentUserRole === 'admin'" value="doctor">Doctor</option>
                   <option v-if="currentUserRole === 'super_admin' || currentUserRole === 'admin'" value="secretary">Secretary</option>
                   <option v-if="currentUserRole === 'super_admin' || currentUserRole === 'admin'" value="assistant">Assistant</option>
                 </select>
                  <p v-if="currentUserRole === 'admin'" class="text-xs text-neutral-500 mt-1">
                    As an Admin, you can update staff roles within your clinic.
                  </p>
               </div>

               <!-- Current Role Display (for non-admins in edit mode) -->
               <div v-if="currentUserRole !== 'admin' && currentUserRole !== 'super_admin'" class="space-y-2">
                 <label class="block text-sm font-semibold text-gray-700 flex items-center">
                   <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                   </svg>
                   Current Role
                 </label>
                 <div class="bg-gray-50 border border-gray-200 rounded-xl p-4">
                   <p class="text-sm text-gray-800 capitalize">{{ form.role.replace('_', ' ') }}</p>
                   <p class="text-xs text-gray-500 mt-1">Contact your clinic administrator to change your role.</p>
                 </div>
                </div>

               <!-- Email -->
              <div class="space-y-2">
                <label for="email" class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"/>
                  </svg>
                  Email Address
                </label>
                <input
                  v-model="form.email"
                  id="email"
                  name="email"
                  type="email"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  placeholder="your@email.com"
                />
              </div>

              <!-- Name Fields -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="firstName" class="block text-sm font-semibold text-gray-700">
                    First Name
                  </label>
                  <input
                    v-model="form.first_name"
                    id="firstName"
                    name="firstName"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="First name"
                  />
                </div>
                <div class="space-y-2">
                  <label for="lastName" class="block text-sm font-semibold text-gray-700">
                    Last Name
                  </label>
                  <input
                    v-model="form.last_name"
                    id="lastName"
                    name="lastName"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="Last name"
                  />
                </div>
              </div>

              <!-- Gender -->
              <div class="space-y-2">
                <label for="gender" class="block text-sm font-semibold text-gray-700">
                  Gender
                </label>
                <select
                  v-model="form.gender"
                  id="gender"
                  name="gender"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                >
                  <option value="">Select gender</option>
                  <option value="male">Male</option>
                  <option value="female">Female</option>
                  <option value="other">Other</option>
                  <option value="prefer_not_to_say">Prefer not to say</option>
                </select>
              </div>
            </div>

            <!-- Password Section -->
            <div class="space-y-6 pt-6 border-t border-neutral-200">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Change Password</h3>
            <p class="text-sm text-neutral-600">Leave blank to keep current password</p>

               <!-- New Password -->
               <div class="space-y-2">
                 <label for="password" class="block text-sm font-semibold text-gray-700 flex items-center">
                   <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                   </svg>
                  New Password
                 </label>
                 <div class="relative">
                   <input
                     v-model="form.password"
                     id="password"
                     name="password"
                     :type="showPassword ? 'text' : 'password'"
                  class="block w-full px-4 py-3 pr-12 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
                  placeholder="Enter new password (optional)"
                   />
                   <button
                     type="button"
                     @click="showPassword = !showPassword"
                     class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-gray-600"
                   >
                     <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                     </svg>
                     <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                     </svg>
                   </button>
                 </div>
               </div>
            </div>

            <!-- Submit Error -->
            <BaseTransition name="fade">
              <div v-if="submitError" class="bg-danger-50 border border-danger-200 rounded-xl p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-danger-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  <p class="text-sm text-danger-700">{{ submitError }}</p>
                </div>
              </div>
            </BaseTransition>

            <!-- Success Message -->
            <BaseTransition name="fade">
              <div v-if="submitSuccess" class="bg-success-50 border border-success-200 rounded-xl p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-success-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                     <p class="text-sm text-success-700">Staff profile updated successfully!</p>
                </div>
              </div>
            </BaseTransition>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-6 border-t border-neutral-200">
               <button
                 type="submit"
                 :disabled="saving || !hasChanges"
                 class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
               >
                  <BaseLoading v-if="saving" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                    {{ saving ? 'Saving...' : 'Save Changes' }}
               </button>

              <button
                type="button"
                @click="resetForm"
                :disabled="saving"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Reset
              </button>
            </div>
          </form>
        </div>
      </div>
    </BaseTransition>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'
import UserAvatar from '../components/UserAvatar.vue'
import AvatarUpload from '../components/AvatarUpload.vue'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'

export default {
  name: 'StaffEdit',
  components: {
    UserAvatar,
    AvatarUpload,
    BaseLoading,
    BaseTransition
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const authStore = useAuthStore()

    const loading = ref(true)
    const saving = ref(false)
    const error = ref('')
    const submitError = ref('')
    const submitSuccess = ref(false)
    const showPassword = ref(false)

    const originalUser = ref(null)
    const form = ref({
      id: 0,
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      gender: '',
      avatar_path: '',
      password: '',
      role: ''
    })

    const userId = computed(() => route.params.id)
    const currentUserRole = computed(() => authStore.user?.role || '')

    const hasChanges = computed(() => {
      if (!originalUser.value) return false

      // Check basic field changes
      let hasBasicChanges = (
        form.value.username !== originalUser.value.username ||
        form.value.email !== (originalUser.value.email || '') ||
        form.value.first_name !== (originalUser.value.first_name || '') ||
        form.value.last_name !== (originalUser.value.last_name || '') ||
        form.value.gender !== (originalUser.value.gender || '') ||
        form.value.avatar_path !== (originalUser.value.avatar_path || originalUser.value.avatar || '') ||
        form.value.password !== ''
      )

      // Check role changes (only for admins and super admins)
      let hasRoleChanges = false
      if (currentUserRole.value === 'admin' || currentUserRole.value === 'super_admin') {
        hasRoleChanges = form.value.role !== (originalUser.value.role || '')
      }

      return hasBasicChanges || hasRoleChanges
    })

    const loadUser = async () => {
      try {
        loading.value = true
        error.value = ''

        const result = await apiService.getUser(userId.value)

        if (result.success) {
          originalUser.value = result.data
          resetForm()
        } else {
          error.value = result.error || 'Failed to load staff member. Please try again.'
        }
      } catch (err) {
        error.value = 'An unexpected error occurred while loading the staff member.'
        console.error(err)
      } finally {
        loading.value = false
      }
    }

    const resetForm = () => {
      if (originalUser.value) {
        // Reset to original user data
        form.value = {
          id: originalUser.value.id || 0,
          username: originalUser.value.username,
          email: originalUser.value.email || '',
          first_name: originalUser.value.first_name || '',
          last_name: originalUser.value.last_name || '',
          gender: originalUser.value.gender || '',
          avatar_path: originalUser.value.avatar_path || originalUser.value.avatar || '',
          password: '',
          role: originalUser.value.role || ''
        }
      }
    }

    const handleAvatarUpdated = (avatarPath) => {
      form.value.avatar_path = avatarPath || ''
    }

    const handleSubmit = async () => {
      try {
        saving.value = true
        submitError.value = ''
        submitSuccess.value = false

        // Edit mode - prepare update data (only include changed fields)
        const updateData = {}

        if (form.value.username !== originalUser.value.username) {
          updateData.username = form.value.username
        }
        if (form.value.email !== (originalUser.value.email || '')) {
          updateData.email = form.value.email
        }
        if (form.value.first_name !== (originalUser.value.first_name || '')) {
          updateData.first_name = form.value.first_name
        }
        if (form.value.last_name !== (originalUser.value.last_name || '')) {
          updateData.last_name = form.value.last_name
        }
        if (form.value.gender !== (originalUser.value.gender || '')) {
          updateData.gender = form.value.gender
        }
        if (form.value.avatar_path !== (originalUser.value.avatar_path || originalUser.value.avatar || '')) {
          updateData.avatar_path = form.value.avatar_path
        }
        if (form.value.password) {
          updateData.password = form.value.password
        }

        // Include role changes for admins and super admins
        if ((currentUserRole.value === 'admin' || currentUserRole.value === 'super_admin') &&
            form.value.role !== (originalUser.value.role || '')) {
          updateData.role = form.value.role
        }

        const result = await apiService.updateUser(userId.value, updateData)

        if (result.success) {
          originalUser.value = { ...originalUser.value, ...updateData }
          form.value.password = '' // Clear password field
          submitSuccess.value = true

          // Hide success message after 3 seconds
          setTimeout(() => {
            submitSuccess.value = false
          }, 3000)
        } else {
          submitError.value = result.error || 'Failed to update staff member'
        }
      } catch (err) {
        console.error('Update staff error:', err)
        submitError.value = 'Failed to update staff member. Please try again.'
      } finally {
        saving.value = false
      }
    }

    // Clear messages when form changes
    watch(() => form.value, () => {
      submitError.value = ''
      submitSuccess.value = false
    }, { deep: true })

    onMounted(async () => {
      await loadUser()
    })

    return {
      loading,
      saving,
      error,
      submitError,
      submitSuccess,
      showPassword,
      originalUser,
      form,
      currentUserRole,
      hasChanges,
      loadUser,
      resetForm,
      handleAvatarUpdated,
      handleSubmit
    }
  }
}
</script>