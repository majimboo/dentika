<template>
  <div v-if="shouldShowOverlay" class="connection-overlay">
    <BaseLoading
      type="spinner"
      size="large"
      color="primary"
      overlay
      text="Connecting to server..."
      class="connection-loading"
    />

    <div class="connection-message">
      <div class="message-content">
        <svg class="warning-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>

        <h3 class="message-title">Connection Lost</h3>

        <p class="message-text">
          {{ connectionMessage }}
        </p>

        <div class="retry-info">
          <p class="retry-text">
            Retrying automatically...
            <span v-if="lastChecked" class="last-checked">
              Last checked: {{ formatLastChecked }}
            </span>
          </p>
        </div>

        <button
          @click="manualRetry"
          class="retry-button"
          :disabled="isChecking"
        >
          {{ isChecking ? 'Checking...' : 'Retry Now' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useConnectionStore } from '../stores/connection'
import BaseLoading from './BaseLoading.vue'

const connectionStore = useConnectionStore()

const shouldShowOverlay = computed(() => connectionStore.shouldShowOverlay)
const isChecking = computed(() => connectionStore.isChecking)
const lastChecked = computed(() => connectionStore.lastChecked)

const connectionMessage = computed(() => {
  if (!connectionStore.isOnline) {
    return "You're currently offline. Please check your internet connection."
  }
  return "Unable to connect to the server. The server might be down or experiencing issues."
})

const formatLastChecked = computed(() => {
  if (!lastChecked.value) return ''

  const now = new Date()
  const diff = now - lastChecked.value
  const seconds = Math.floor(diff / 1000)

  if (seconds < 60) return `${seconds} seconds ago`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes} minutes ago`
  const hours = Math.floor(minutes / 60)
  return `${hours} hours ago`
})

const manualRetry = async () => {
  await connectionStore.checkServerHealth()
}
</script>

<style scoped>
@import "../styles/main.css";

.connection-overlay {
  @apply fixed inset-0 z-50 pointer-events-none;
}

.connection-loading {
  @apply pointer-events-auto;
}

.connection-message {
  @apply fixed inset-0 flex items-center justify-center pointer-events-auto z-10;
}

.message-content {
  @apply bg-white rounded-lg shadow-xl border border-gray-200 p-6 max-w-md mx-4 text-center;
}

.warning-icon {
  @apply w-12 h-12 text-amber-500 mx-auto mb-4;
}

.message-title {
  @apply text-lg font-semibold text-gray-900 mb-2;
}

.message-text {
  @apply text-gray-600 mb-4 text-sm leading-relaxed;
}

.retry-info {
  @apply mb-4;
}

.retry-text {
  @apply text-xs text-gray-500;
}

.last-checked {
  @apply block mt-1 text-gray-400;
}

.retry-button {
  @apply bg-primary-600 hover:bg-primary-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition-colors duration-200 text-sm;
}

.retry-button:disabled {
  @apply cursor-not-allowed;
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
  .message-content {
    @apply bg-gray-800 border-gray-700;
  }

  .message-title {
    @apply text-gray-100;
  }

  .message-text {
    @apply text-gray-300;
  }

  .retry-text {
    @apply text-gray-400;
  }

  .last-checked {
    @apply text-gray-500;
  }
}
</style>