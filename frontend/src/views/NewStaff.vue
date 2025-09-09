<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Add New Staff Member</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Create a new staff account for your clinic</p>
      </div>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
        <button
          @click="$router.go(-1)"
          class="inline-flex items-center justify-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
           <span class="sm:ml-2">Back</span>
        </button>
      </div>
    </div>

    <!-- Create Form -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div class="p-6 sm:p-8">
        <form @submit.prevent="handleSubmit" class="space-y-8">

          <!-- Avatar Section -->
          <div class="pb-6 border-b border-neutral-200">
            <h3 class="text-lg font-semibold text-neutral-900 mb-4">Profile Picture</h3>
            <AvatarUpload
              :user="form"
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
                placeholder="staff@email.com"
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

            <!-- Role Selection -->
            <div class="space-y-2">
              <label for="role" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
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
                <option value="doctor">Doctor</option>
                <option value="secretary">Secretary</option>
                <option value="assistant">Assistant</option>
              </select>
              <p class="text-xs text-neutral-500 mt-1">
                <strong>Doctor:</strong> Full clinical access<br>
                <strong>Secretary:</strong> Front desk and scheduling<br>
                <strong>Assistant:</strong> Limited clinical support
              </p>
            </div>
          </div>

          <!-- Password Section -->
          <div class="space-y-6 pt-6 border-t border-neutral-200">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Set Password</h3>
            <p class="text-sm text-neutral-600">Create a secure password for this staff account</p>

            <!-- Password -->
            <div class="space-y-2">
              <label for="password" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
                Password
                <span class="text-danger-500 ml-1">*</span>
              </label>
              <div class="relative">
                <input
                  v-model="form.password"
                  id="password"
                  name="password"
                  type="password"
                  required
                  minlength="8"
                  class="block w-full px-4 py-3 pr-12 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
                  placeholder="Enter secure password"
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
              <p class="text-xs text-neutral-500">Minimum 8 characters required</p>
            </div>

            <!-- Confirm Password -->
            <div class="space-y-2">
              <label for="confirmPassword" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
                Confirm Password
                <span class="text-danger-500 ml-1">*</span>
              </label>
              <div class="relative">
                <input
                  v-model="form.confirmPassword"
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  required
                  class="block w-full px-4 py-3 pr-12 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
                  placeholder="Confirm password"
                />
                <button
                  type="button"
                  @click="showConfirmPassword = !showConfirmPassword"
                  class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-gray-600"
                >
                  <svg v-if="!showConfirmPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                  </svg>
                </button>
              </div>
              <p v-if="passwordMismatch" class="text-xs text-danger-600">Passwords do not match</p>
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
                <p class="text-sm text-success-700">Staff member added successfully!</p>
              </div>
            </div>
          </BaseTransition>

          <!-- Form Actions -->
          <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-6 border-t border-neutral-200">
            <button
              type="submit"
              :disabled="saving || !isFormValid"
              class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            >
              <BaseLoading v-if="saving" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
              {{ saving ? 'Adding Staff...' : 'Add Staff Member' }}
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
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'
import AvatarUpload from '../components/AvatarUpload.vue'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'

export default {
  name: 'NewStaff',
  components: {
    AvatarUpload,
    BaseLoading,
    BaseTransition
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()

    const saving = ref(false)
    const submitError = ref('')
    const submitSuccess = ref(false)
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)

    const form = ref({
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      gender: '',
      avatar: '',
      password: '',
      confirmPassword: '',
      role: ''
    })

    const passwordMismatch = computed(() => {
      return form.value.password && form.value.confirmPassword &&
             form.value.password !== form.value.confirmPassword
    })

    const isFormValid = computed(() => {
      return form.value.username &&
             form.value.password &&
             form.value.confirmPassword &&
             form.value.password === form.value.confirmPassword &&
             form.value.password.length >= 8 &&
             form.value.role &&
             !passwordMismatch.value
    })

    const resetForm = () => {
      form.value = {
        username: '',
        email: '',
        first_name: '',
        last_name: '',
        gender: '',
        avatar: '',
        password: '',
        confirmPassword: '',
        role: ''
      }
      submitError.value = ''
      submitSuccess.value = false
    }

    const handleAvatarUpdated = (avatarUrl) => {
      form.value.avatar = avatarUrl
    }

    const handleSubmit = async () => {
      if (!isFormValid.value) {
        submitError.value = 'Please fill in all required fields and ensure passwords match.'
        return
      }

      try {
        saving.value = true
        submitError.value = ''
        submitSuccess.value = false

        const submitData = {
          username: form.value.username,
          email: form.value.email,
          first_name: form.value.first_name,
          last_name: form.value.last_name,
          gender: form.value.gender,
          avatar: form.value.avatar,
          password: form.value.password,
          role: form.value.role
        }

        const result = await apiService.post('/api/users', submitData)

        if (result.success) {
          submitSuccess.value = true

          // Redirect to staff list after 2 seconds
          setTimeout(() => {
            router.push('/staff')
          }, 2000)
        } else {
          submitError.value = result.error || 'Failed to create staff member'
        }
      } catch (err) {
        console.error('Error creating staff member:', err)
        submitError.value = 'An unexpected error occurred. Please try again.'
      } finally {
        saving.value = false
      }
    }



    return {
      saving,
      submitError,
      submitSuccess,
      showPassword,
      showConfirmPassword,
      form,
      passwordMismatch,
      isFormValid,
      resetForm,
      handleAvatarUpdated,
      handleSubmit
    }
  }
}
</script>