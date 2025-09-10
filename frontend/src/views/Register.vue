<template>
  <div class="min-h-screen bg-gradient-to-br from-success-50 via-white to-secondary-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full">
      <!-- Logo and Welcome Section -->
      <div class="text-center mb-8">
        <div class="mx-auto h-16 w-16 bg-gradient-to-r from-success-600 to-secondary-600 rounded-xl flex items-center justify-center mb-4">
          <font-awesome-icon icon="fa-solid fa-user-plus" class="h-8 w-8 text-white" />
        </div>
        <h2 class="text-3xl font-bold text-gray-900 mb-2">Join Our Platform</h2>
        <p class="text-gray-600">Create your account to get started</p>
      </div>

      <!-- Register Card -->
      <div class="bg-white rounded-2xl shadow-xl border-0 overflow-hidden">
        <div class="px-8 py-8">
          <form @submit.prevent="handleRegister" class="space-y-6">
            <!-- Username Field -->
            <div class="space-y-2">
              <label for="username" class="block text-sm font-semibold text-gray-700 flex items-center">
                <font-awesome-icon icon="fa-solid fa-user" class="w-4 h-4 mr-2 text-gray-500" />
                Choose a Username
                <BaseTooltip text="Your username will be unique and used for login. Choose something memorable but secure." position="top" :showMobileIcon="false">
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
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-success-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  placeholder="Enter a unique username"
                  :class="{ 'border-danger-300 focus:ring-danger-500': usernameError }"
                  @input="validateUsername"
                />
                <div v-if="usernameValid && form.username" class="absolute inset-y-0 right-0 pr-4 flex items-center">
                  <font-awesome-icon icon="fa-solid fa-check-circle" class="w-5 h-5 text-green-500" />
                </div>
              </div>
              <div class="space-y-1">
                <p class="text-xs text-gray-500 flex items-center">
                  <font-awesome-icon icon="fa-solid fa-info-circle" class="w-3 h-3 mr-1" />
                  Must be 3-20 characters, letters and numbers only
                </p>
                <p v-if="usernameError" class="text-xs text-danger-600 flex items-center">
                  <font-awesome-icon icon="fa-solid fa-exclamation-circle" class="w-3 h-3 mr-1" />
                  {{ usernameError }}
                </p>
              </div>
            </div>

            <!-- Email Field -->
            <div class="space-y-2">
              <label for="email" class="block text-sm font-semibold text-gray-700 flex items-center">
                <font-awesome-icon icon="fa-solid fa-envelope" class="w-4 h-4 mr-2 text-gray-500" />
                Email Address (Optional)
              </label>
              <input
                v-model="form.email"
                id="email"
                name="email"
                type="email"
                class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-success-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                placeholder="your@email.com"
              />
            </div>

            <!-- Password Field -->
            <div class="space-y-2">
              <label for="password" class="block text-sm font-semibold text-gray-700 flex items-center">
                <font-awesome-icon icon="fa-solid fa-lock" class="w-4 h-4 mr-2 text-gray-500" />
                Create Password
                <BaseTooltip text="Use a combination of uppercase, lowercase, numbers, and special characters for a strong password." position="top" :showMobileIcon="false">
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
                  class="block w-full px-4 py-3 pr-12 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
                  placeholder="Create a strong password"
                  @input="validatePassword"
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
              <div class="space-y-1">
                <div class="flex items-center space-x-2 text-xs">
                  <div class="flex items-center" :class="passwordStrength.length ? 'text-success-600' : 'text-neutral-400'">
                    <font-awesome-icon icon="fa-solid fa-check" class="w-3 h-3 mr-1" />
                    8+ characters
                  </div>
                  <div class="flex items-center" :class="passwordStrength.uppercase ? 'text-success-600' : 'text-neutral-400'">
                    <font-awesome-icon icon="fa-solid fa-check" class="w-3 h-3 mr-1" />
                    Uppercase
                  </div>
                  <div class="flex items-center" :class="passwordStrength.number ? 'text-success-600' : 'text-neutral-400'">
                    <font-awesome-icon icon="fa-solid fa-check" class="w-3 h-3 mr-1" />
                    Number
                  </div>
                </div>
              </div>
            </div>

            <!-- Confirm Password Field -->
            <div class="space-y-2">
              <label for="confirmPassword" class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.031 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                </svg>
                Confirm Password
              </label>
              <div class="relative">
                <input
                  v-model="form.confirmPassword"
                  id="confirmPassword"
                  name="confirmPassword"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  required
                  class="block w-full px-4 py-3 pr-12 border border-gray-300 rounded-xl text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
                  placeholder="Confirm your password"
                  :class="{ 'border-success-300 focus:ring-success-500': passwordsMatch && form.confirmPassword, 'border-danger-300 focus:ring-danger-500': !passwordsMatch && form.confirmPassword }"
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
              <div v-if="form.confirmPassword">
                <p v-if="passwordsMatch" class="text-xs text-success-600 flex items-center">
                  <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4"/>
                  </svg>
                  Passwords match
                </p>
                <p v-else class="text-xs text-danger-600 flex items-center">
                  <font-awesome-icon icon="fa-solid fa-times" class="w-3 h-3 mr-1" />
                  Passwords don't match
                </p>
              </div>
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
                :disabled="loading || !isFormValid"
                class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-xl text-white bg-gradient-to-r from-success-600 to-secondary-600 hover:from-success-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-success-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed font-semibold text-sm transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98]"
              >
                <font-awesome-icon
                  v-if="loading"
                  icon="fa-solid fa-spinner"
                  class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                />
                {{ loading ? 'Creating Account...' : 'Create Account' }}
              </button>
              
              <!-- Login Link -->
              <div class="text-center">
                <p class="text-sm text-gray-600">
                  Already have an account? 
                  <router-link to="/login" class="font-semibold text-success-600 hover:text-success-700 transition-colors duration-200">
                    Sign in here
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
          Your information is protected with enterprise-grade security
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import BaseTooltip from '../components/BaseTooltip.vue'

export default {
  name: 'Register',
  components: {
    BaseTooltip
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })
    
    const error = ref('')
    const loading = ref(false)
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const usernameError = ref('')
    const usernameValid = ref(false)
    
    const passwordStrength = computed(() => ({
      length: form.value.password.length >= 8,
      uppercase: /[A-Z]/.test(form.value.password),
      number: /[0-9]/.test(form.value.password)
    }))
    
    const passwordsMatch = computed(() => {
      return form.value.password === form.value.confirmPassword
    })
    
    const isFormValid = computed(() => {
      return usernameValid.value && 
             passwordStrength.value.length && 
             passwordStrength.value.uppercase && 
             passwordStrength.value.number && 
             passwordsMatch.value
    })
    
    const validateUsername = () => {
      const username = form.value.username
      usernameError.value = ''
      usernameValid.value = false
      
      if (!username) {
        return
      }
      
      if (username.length < 3) {
        usernameError.value = 'Username must be at least 3 characters'
        return
      }
      
      if (username.length > 20) {
        usernameError.value = 'Username must be less than 20 characters'
        return
      }
      
      if (!/^[a-zA-Z0-9]+$/.test(username)) {
        usernameError.value = 'Username can only contain letters and numbers'
        return
      }
      
      usernameValid.value = true
    }
    
    const validatePassword = () => {
      // Password validation is handled by computed property
    }
    
    const handleRegister = async () => {
      if (!isFormValid.value) {
        error.value = 'Please fix the errors above'
        return
      }
      
      loading.value = true
      error.value = ''
      
      const result = await authStore.register({
        username: form.value.username,
        password: form.value.password
      })
      
      if (result.success) {
        router.push('/')
      } else {
        error.value = result.error
      }
      
      loading.value = false
    }
    
    return {
      form,
      error,
      loading,
      showPassword,
      showConfirmPassword,
      usernameError,
      usernameValid,
      passwordStrength,
      passwordsMatch,
      isFormValid,
      validateUsername,
      validatePassword,
      handleRegister
    }
  }
}
</script>