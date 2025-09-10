<template>
  <div class="avatar-upload-component">
    <!-- Avatar Display -->
    <div class="relative inline-block">
      <div 
        class="avatar-container w-24 h-24 rounded-full overflow-hidden border-4 border-gray-200 bg-gray-100 cursor-pointer transition-all duration-200 hover:border-blue-300 group"
        @click="triggerFileInput"
      >
        <!-- Current Avatar -->
        <img 
          v-if="currentAvatarUrl" 
          :src="currentAvatarUrl" 
          :alt="`${fullName} avatar`"
          class="w-full h-full object-cover"
          @error="handleImageError"
        />
        
        <!-- Fallback Initials -->
        <div 
          v-else 
          class="w-full h-full flex items-center justify-center bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-bold text-xl"
        >
          {{ initials }}
        </div>
        
        <!-- Upload Overlay -->
        <div class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <div class="text-white text-center">
            <font-awesome-icon icon="fa-solid fa-camera" class="w-6 h-6 mb-1" />
            <div class="text-xs font-medium">{{ currentAvatarUrl ? 'Change' : 'Upload' }}</div>
          </div>
        </div>
      </div>
      
      <!-- Remove Button -->
      <button
        v-if="currentAvatarUrl && !disabled"
        @click.stop="removeAvatar"
        class="absolute -top-2 -right-2 w-8 h-8 bg-red-500 hover:bg-red-600 text-white rounded-full flex items-center justify-center transition-colors duration-200"
        title="Remove avatar"
        :disabled="uploading"
      >
        <font-awesome-icon icon="fa-solid fa-times" class="w-4 h-4" />
      </button>
      
      <!-- Upload Progress -->
      <div 
        v-if="uploading" 
        class="absolute inset-0 rounded-full bg-black bg-opacity-70 flex items-center justify-center"
      >
        <div class="text-white text-center">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-white mb-1"></div>
          <div class="text-xs">Uploading...</div>
        </div>
      </div>
    </div>
    
    <!-- Hidden File Input -->
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
      @change="handleFileChange"
      class="hidden"
      :disabled="disabled || uploading"
    />
    
    <!-- Upload Instructions -->
    <div class="mt-2 text-center">
      <div class="text-sm text-gray-600">
        {{ currentAvatarUrl ? 'Click to change avatar' : 'Click to upload avatar' }}
      </div>
      <div class="text-xs text-gray-500 mt-1">
        JPEG, PNG, GIF, WebP (max 5MB)
      </div>
    </div>
    
    <!-- Error Message -->
    <div v-if="error" class="mt-2 text-center">
      <div class="text-sm text-red-600">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useNotificationStore } from '../stores/notification'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  user: {
    type: Object,
    default: null
  },
  fullName: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  size: {
    type: String,
    default: 'large', // 'small', 'medium', 'large'
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  }
})

const emit = defineEmits(['update:modelValue', 'upload-success', 'upload-error', 'avatar-updated'])

const notificationStore = useNotificationStore()

const fileInput = ref(null)
const uploading = ref(false)
const error = ref('')

const currentAvatarUrl = computed(() => {
  // Support both v-model pattern and user object pattern
  const avatarPath = props.user ? (props.user.avatar_path || props.user.avatar) : props.modelValue
  if (!avatarPath) return ''

  // If it's already a full URL, return as is
  if (avatarPath.startsWith('http')) {
    return avatarPath
  }

  // If it's a relative path, prepend the server URL
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:3000'
  return `${baseUrl}/uploads/${avatarPath}`
})

const initials = computed(() => {
  // Support both patterns: fullName prop or user object
  let fullName = props.fullName
  if (props.user && !fullName) {
    if (props.user.first_name && props.user.last_name) {
      fullName = `${props.user.first_name} ${props.user.last_name}`
    } else if (props.user.username) {
      return props.user.username.charAt(0).toUpperCase()
    }
  }
  
  if (!fullName) return '?'
  
  const nameParts = fullName.trim().split(' ')
  if (nameParts.length === 1) {
    return nameParts[0].charAt(0).toUpperCase()
  }
  
  return (nameParts[0].charAt(0) + nameParts[nameParts.length - 1].charAt(0)).toUpperCase()
})

const sizeClasses = computed(() => {
  const sizes = {
    small: 'w-12 h-12 text-sm',
    medium: 'w-16 h-16 text-base',
    large: 'w-24 h-24 text-xl'
  }
  return sizes[props.size] || sizes.large
})

const triggerFileInput = () => {
  if (!props.disabled && !uploading.value) {
    fileInput.value?.click()
  }
}

const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  error.value = ''
  
  // Validate file size (5MB)
  if (file.size > 5 * 1024 * 1024) {
    error.value = 'File size must be less than 5MB'
    return
  }
  
  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    error.value = 'Only JPEG, PNG, GIF, and WebP files are allowed'
    return
  }
  
  await uploadFile(file)
}

const uploadFile = async (file) => {
  uploading.value = true
  error.value = ''
  
  try {
    const formData = new FormData()
    formData.append('avatar', file)
    
    const response = await fetch('/api/upload/avatar', {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Upload failed')
    }
    
    const data = await response.json()
    
    // Update the model value with the relative path
    emit('update:modelValue', data.path)
    emit('upload-success', data)
    // Also emit avatar-updated for user object pattern
    emit('avatar-updated', data.path)
    
    notificationStore.showSuccess('Avatar uploaded successfully')
    
    // Clear file input
    if (fileInput.value) {
      fileInput.value.value = ''
    }
    
  } catch (err) {
    error.value = err.message || 'Failed to upload avatar'
    emit('upload-error', err)
    notificationStore.showError('Failed to upload avatar')
    
  } finally {
    uploading.value = false
  }
}

const removeAvatar = async () => {
  if (!props.modelValue || uploading.value) return
  
  uploading.value = true
  error.value = ''
  
  try {
    const response = await fetch(`/api/upload/avatar?path=${encodeURIComponent(props.modelValue)}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to remove avatar')
    }
    
    // Clear the model value
    emit('update:modelValue', '')
    emit('avatar-updated', '')
    notificationStore.showSuccess('Avatar removed successfully')
    
  } catch (err) {
    error.value = err.message || 'Failed to remove avatar'
    notificationStore.showError('Failed to remove avatar')
    
  } finally {
    uploading.value = false
  }
}

const handleImageError = () => {
  // If image fails to load, clear the avatar path
  emit('update:modelValue', '')
  emit('avatar-updated', '')
}

// Clear error when modelValue changes
watch(() => props.modelValue, () => {
  error.value = ''
})
</script>

<style scoped>
@import "../styles/main.css";

.avatar-upload-component {
  @apply flex flex-col items-center;
}

.avatar-container {
  @apply relative transition-all duration-200;
}

.avatar-container:hover {
  @apply shadow-lg;
}

/* Size variants */
.avatar-container.size-small {
  @apply w-12 h-12;
}

.avatar-container.size-medium {
  @apply w-16 h-16;
}

.avatar-container.size-large {
  @apply w-24 h-24;
}

/* Disabled state */
.avatar-upload-component[disabled] .avatar-container {
  @apply cursor-not-allowed opacity-60;
}

.avatar-upload-component[disabled] .avatar-container:hover {
  @apply border-gray-200 shadow-none;
}
</style>