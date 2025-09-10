<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full">
      <!-- Logo and Welcome Section -->
      <div class="text-center mb-8">
        <div class="mx-auto h-16 w-16 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-xl flex items-center justify-center mb-4">
          <font-awesome-icon icon="fa-solid fa-bolt" class="h-8 w-8 text-white" />
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
                <font-awesome-icon icon="fa-solid fa-user" class="w-4 h-4 mr-2 text-gray-500" />
                Username
                <BaseTooltip text="Enter the unique username you chose during registration. This is case-sensitive." position="top" :showMobileIcon="false">
                  <font-awesome-icon icon="fa-solid fa-info-circle" class="w-4 h-4 ml-2 text-neutral-400 hover:text-neutral-600 cursor-help" />
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
                <font-awesome-icon icon="fa-solid fa-info-circle" class="w-3 h-3 mr-1 text-neutral-400" />
                Use the username you created during registration
              </p>
            </div>

            <!-- Password Field -->
            <div class="space-y-2">
              <label for="password" class="block text-sm font-semibold text-gray-700 flex items-center">
                <font-awesome-icon icon="fa-solid fa-lock" class="w-4 h-4 mr-2 text-gray-500" />
                Password
                <BaseTooltip text="Click the eye icon to show/hide your password. We recommend using a strong, unique password." position="top" :showMobileIcon="false">
                  <font-awesome-icon icon="fa-solid fa-info-circle" class="w-4 h-4 ml-2 text-neutral-400 hover:text-neutral-600 cursor-help" />
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
                  <font-awesome-icon
                    :icon="!showPassword ? 'fa-solid fa-eye' : 'fa-solid fa-eye-slash'"
                    class="w-5 h-5"
                  />
                </button>
              </div>
              <p class="text-xs text-neutral-500 flex items-center">
                <font-awesome-icon icon="fa-solid fa-shield-alt" class="w-3 h-3 mr-1 text-neutral-400" />
                Your password is encrypted and secure
              </p>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="bg-danger-50 border border-danger-200 rounded-xl p-4">
              <div class="flex items-center">
                <font-awesome-icon icon="fa-solid fa-exclamation-circle" class="w-5 h-5 text-danger-400 mr-2" />
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
                <font-awesome-icon
                  v-if="loading"
                  icon="fa-solid fa-spinner"
                  class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                />
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
          <font-awesome-icon icon="fa-solid fa-shield-alt" class="w-4 h-4 mr-1" />
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