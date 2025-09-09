<template>
  <div class="flex flex-col items-center space-y-4">
    <!-- Current Avatar Display -->
    <div class="relative">
      <UserAvatar 
        :user="user" 
        size="2xl"
        class="transition-all duration-200 hover:scale-105"
      />
      
      <!-- Upload overlay on hover -->
      <div class="absolute inset-0 bg-black bg-opacity-50 rounded-full flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity duration-200 cursor-pointer">
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
      </div>
    </div>

    <!-- Upload Controls -->
    <div class="flex flex-col sm:flex-row items-center gap-3">
      <!-- File Input -->
      <input 
        ref="fileInput"
        type="file" 
        accept="image/*"
        @change="handleFileSelect"
        class="hidden"
      />
      
      <!-- Upload Button -->
      <button
        @click="$refs.fileInput.click()"
        :disabled="uploading"
        class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
        </svg>
        {{ uploading ? 'Uploading...' : 'Choose Photo' }}
      </button>
      
      <!-- Remove Button -->
      <button
        v-if="user?.avatar"
        @click="removeAvatar"
        :disabled="uploading"
        class="inline-flex items-center px-4 py-2 border border-danger-300 rounded-lg text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
        Remove
      </button>
    </div>

    <!-- Upload Progress -->
    <BaseTransition name="fade">
      <div v-if="uploading" class="w-full max-w-xs">
        <div class="flex items-center space-x-3">
          <BaseLoading type="spinner" size="small" color="primary" :show-text="false" />
          <div class="flex-1">
            <div class="text-sm text-neutral-600">Uploading...</div>
            <div class="w-full bg-neutral-200 rounded-full h-2 mt-1">
              <div 
                class="bg-primary-600 h-2 rounded-full transition-all duration-300" 
                :style="{ width: uploadProgress + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </BaseTransition>

    <!-- Upload Guidelines -->
    <div class="text-center">
      <p class="text-sm text-neutral-500">
        Recommended: Square image, at least 200x200px
      </p>
      <p class="text-xs text-neutral-400 mt-1">
        JPG, PNG, GIF, or WebP. Max file size: 5MB
      </p>
    </div>

    <!-- Error Message -->
    <BaseTransition name="fade">
      <div v-if="error" class="bg-danger-50 border border-danger-200 rounded-lg p-3 max-w-sm">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-danger-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <p class="text-sm text-danger-700">{{ error }}</p>
        </div>
      </div>
    </BaseTransition>
  </div>
</template>

<script>
import UserAvatar from './UserAvatar.vue'
import BaseLoading from './BaseLoading.vue'
import BaseTransition from './BaseTransition.vue'
import apiService from '../services/api'

export default {
  name: 'AvatarUpload',
  components: {
    UserAvatar,
    BaseLoading,
    BaseTransition
  },
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      uploading: false,
      uploadProgress: 0,
      error: null
    }
  },
  methods: {
    async handleFileSelect(event) {
      const file = event.target.files[0]
      if (!file) return

      // Validate file
      const validationError = this.validateFile(file)
      if (validationError) {
        this.error = validationError
        return
      }

      this.error = null
      await this.uploadFile(file)
    },

    validateFile(file) {
      // Check file type
      const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
      if (!allowedTypes.includes(file.type)) {
        return 'Please select a valid image file (JPG, PNG, GIF, or WebP)'
      }

      // Check file size (5MB limit)
      const maxSize = 5 * 1024 * 1024
      if (file.size > maxSize) {
        return 'File size must be less than 5MB'
      }

      return null
    },

    async uploadFile(file) {
      this.uploading = true
      this.uploadProgress = 0

      try {
        const formData = new FormData()
        formData.append('avatar', file)

        // Simulate progress for better UX
        const progressInterval = setInterval(() => {
          if (this.uploadProgress < 90) {
            this.uploadProgress += Math.random() * 20
          }
        }, 200)

        const result = await apiService.uploadAvatar(formData)

        clearInterval(progressInterval)
        this.uploadProgress = 100

        if (result.success) {
          // Update user avatar
          this.$emit('avatar-updated', result.data.path)
          
          // Clear file input
          this.$refs.fileInput.value = ''
        } else {
          this.error = result.error || 'Upload failed'
        }
      } catch (error) {
        console.error('Upload error:', error)
        this.error = 'Upload failed. Please try again.'
      } finally {
        this.uploading = false
        
        // Reset progress after a delay
        setTimeout(() => {
          this.uploadProgress = 0
        }, 1000)
      }
    },

    async removeAvatar() {
      if (!confirm('Are you sure you want to remove your avatar?')) {
        return
      }

      try {
        const result = await apiService.deleteAvatar(this.user.avatar)
        
        if (result.success) {
          this.$emit('avatar-updated', null)
        } else {
          this.error = result.error || 'Failed to remove avatar'
        }
      } catch (error) {
        console.error('Remove avatar error:', error)
        this.error = 'Failed to remove avatar. Please try again.'
      }
    }
  }
}
</script>