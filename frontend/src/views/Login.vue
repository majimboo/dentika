<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full">
      <!-- Logo and Welcome Section -->
      <div class="text-center mb-8">
        <div class="mx-auto h-16 w-16 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-xl flex items-center justify-center mb-4">
          <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
        </div>
        <h2 class="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h2>
        <p class="text-gray-600">Please sign in to your account to continue</p>
      </div>

      <!-- Login Card -->
      <div class="bg-white rounded-2xl shadow-xl border-0 overflow-hidden">
        <div class="px-8 py-8">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <!-- Username Field -->
            <div class="space-y-2">
              <label for="username" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                Username
                <BaseTooltip text="Enter the unique username you chose during registration. This is case-sensitive." position="top" :showMobileIcon="false">
                  <svg class="w-4 h-4 ml-2 text-neutral-400 hover:text-neutral-600 cursor-help" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </BaseTooltip>
              </label>
              <div class="relative">
                <input
                  v-model="form.username"
                  id="username"
                  name="username"
                  type="text"
                  required
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  placeholder="Enter your username"
                />
              </div>
              <p class="text-xs text-neutral-500 flex items-center">
                <svg class="w-3 h-3 mr-1 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Use the username you created during registration
              </p>
            </div>

            <!-- Password Field -->
            <div class="space-y-2">
              <label for="password" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                </svg>
                Password
                <BaseTooltip text="Click the eye icon to show/hide your password. We recommend using a strong, unique password." position="top" :showMobileIcon="false">
                  <svg class="w-4 h-4 ml-2 text-neutral-400 hover:text-neutral-600 cursor-help" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </BaseTooltip>
              </label>
              <div class="relative">
                <input
                  v-model="form.password"
                  id="password"
                  name="password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  class="block w-full px-4 py-3 pr-12 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  placeholder="Enter your password"
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
              <p class="text-xs text-neutral-500 flex items-center">
                <svg class="w-3 h-3 mr-1 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.031 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                </svg>
                Your password is encrypted and secure
              </p>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="bg-danger-50 border border-danger-200 rounded-xl p-4">
              <div class="flex items-center">
                <svg class="w-5 h-5 text-danger-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <p class="text-sm text-danger-700">{{ error }}</p>
              </div>
            </div>

            <!-- Submit Button -->
            <div class="space-y-4">
              <button
                type="submit"
                :disabled="loading"
                class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-xl text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed font-semibold text-sm transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98] disabled:animate-pulse"
              >
                <svg v-if="loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ loading ? 'Signing In...' : 'Sign In' }}
              </button>
              
              <!-- Register Link -->
              <div class="text-center">
                <p class="text-sm text-gray-600">
                  Don't have an account? 
                  <router-link to="/register" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">
                    Create one here
                  </router-link>
                </p>
              </div>
            </div>
          </form>
        </div>
      </div>

      <!-- Help Text -->
      <div class="mt-6 text-center">
        <p class="text-xs text-neutral-500 flex items-center justify-center">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.031 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
          Secure login protected by industry-standard encryption
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import BaseTooltip from '../components/BaseTooltip.vue'

export default {
  name: 'Login',
  components: {
    BaseTooltip
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      username: '',
      password: ''
    })
    
    const error = ref('')
    const loading = ref(false)
    const showPassword = ref(false)
    
    const handleLogin = async () => {
      loading.value = true
      error.value = ''
      
      console.log('Login form submitted with:', form.value)
      
      try {
        const result = await authStore.login(form.value)
        console.log('Login result received:', result)
        
        if (result.success) {
          console.log('Login successful, navigating to dashboard...')
          await router.push('/')
          console.log('Navigation completed')
        } else {
          error.value = result.error
        }
      } catch (err) {
        console.error('Unexpected error during login:', err)
        error.value = 'An unexpected error occurred'
      }
      
      loading.value = false
    }
    
    return {
      form,
      error,
      loading,
      showPassword,
      handleLogin
    }
  }
}
</script>