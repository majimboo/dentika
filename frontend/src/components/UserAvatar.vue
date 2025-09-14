<template>
  <div class="relative inline-block">
    <div 
      :class="sizeClasses"
      class="rounded-full overflow-hidden bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center text-white font-semibold border-2 border-white shadow-lg"
    >
      <img 
        v-if="user?.avatar" 
        :src="getAvatarUrl(user.avatar)" 
        :alt="`${user.first_name || user.username}'s avatar`"
        class="w-full h-full object-cover"
        @error="handleImageError"
      />
      <span 
        v-else 
        :class="textSizeClasses"
      >
        {{ getInitials(user) }}
      </span>
    </div>
    
    <!-- Online indicator -->
    <div 
      v-if="showStatus && isOnline" 
      :class="statusSizeClasses"
      class="absolute bottom-0 right-0 bg-success-400 rounded-full border-2 border-white"
    ></div>
  </div>
</template>

<script>
import { getInitials } from '@/utils'

export default {
  name: 'UserAvatar',
  props: {
    user: {
      type: Object,
      required: true
    },
    size: {
      type: String,
      default: 'md',
      validator: value => ['xs', 'sm', 'md', 'lg', 'xl', '2xl'].includes(value)
    },
    showStatus: {
      type: Boolean,
      default: false
    },
    isOnline: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    sizeClasses() {
      const sizes = {
        xs: 'w-6 h-6',
        sm: 'w-8 h-8',
        md: 'w-10 h-10',
        lg: 'w-12 h-12',
        xl: 'w-16 h-16',
        '2xl': 'w-20 h-20'
      }
      return sizes[this.size]
    },
    textSizeClasses() {
      const sizes = {
        xs: 'text-xs',
        sm: 'text-sm',
        md: 'text-base',
        lg: 'text-lg',
        xl: 'text-xl',
        '2xl': 'text-2xl'
      }
      return sizes[this.size]
    },
    statusSizeClasses() {
      const sizes = {
        xs: 'w-2 h-2',
        sm: 'w-2 h-2',
        md: 'w-3 h-3',
        lg: 'w-3 h-3',
        xl: 'w-4 h-4',
        '2xl': 'w-5 h-5'
      }
      return sizes[this.size]
    }
  },
  methods: {
    getInitials(user) {
      return getInitials(user, 'U')
    },
    getAvatarUrl(avatarPath) {
      if (!avatarPath) return ''

      // If it's already a full URL, return as is
      if (avatarPath.startsWith('http')) {
        return avatarPath
      }

      // Simply use current host
      const baseUrl = window.location.host

      return `/uploads/${avatarPath}`
    },
    handleImageError(event) {
      // Hide the image and show initials when image fails to load
      event.target.style.display = 'none'
    }
  }
}
</script>