import { ref, computed } from 'vue'

// Default timezone - should match server configuration
const DEFAULT_TIMEZONE = 'Asia/Manila'

export function useTimezone() {
  // Reactive timezone setting
  const timezone = ref(DEFAULT_TIMEZONE)

  // Set timezone (can be called to update from user preferences)
  const setTimezone = (newTimezone) => {
    timezone.value = newTimezone
  }

  // Get current timezone
  const getTimezone = () => timezone.value

  // Format date/time for display using configured timezone
  const formatDateTime = (dateTime, options = {}) => {
    if (!dateTime) return ''

    const date = new Date(dateTime)
    if (isNaN(date.getTime())) return ''

    const defaultOptions = {
      timeZone: timezone.value,
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      ...options
    }

    return new Intl.DateTimeFormat('en-US', defaultOptions).format(date)
  }

  // Format date only
  const formatDate = (dateTime, options = {}) => {
    return formatDateTime(dateTime, {
      hour: undefined,
      minute: undefined,
      ...options
    })
  }

  // Format time only
  const formatTime = (dateTime, options = {}) => {
    return formatDateTime(dateTime, {
      year: undefined,
      month: undefined,
      day: undefined,
      ...options
    })
  }

  // Convert UTC date to local timezone
  const toLocalTime = (utcDateTime) => {
    if (!utcDateTime) return null

    const date = new Date(utcDateTime)
    if (isNaN(date.getTime())) return null

    // Convert to configured timezone
    return new Date(date.toLocaleString('en-US', { timeZone: timezone.value }))
  }

  // Get relative time (e.g., "2 hours ago")
  const getRelativeTime = (dateTime) => {
    if (!dateTime) return ''

    const date = new Date(dateTime)
    if (isNaN(date.getTime())) return ''

    const now = new Date()
    const diffInMs = now - date
    const diffInMinutes = Math.floor(diffInMs / (1000 * 60))
    const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60))
    const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))

    if (diffInMinutes < 1) return 'Just now'
    if (diffInMinutes < 60) return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`
    if (diffInHours < 24) return `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`
    if (diffInDays < 7) return `${diffInDays} day${diffInDays > 1 ? 's' : ''} ago`

    return formatDate(dateTime)
  }

  // Check if appointment is today
  const isToday = (dateTime) => {
    if (!dateTime) return false

    const date = new Date(dateTime)
    const today = new Date()

    return date.toDateString() === today.toDateString()
  }

  // Check if appointment is upcoming (within next 24 hours)
  const isUpcoming = (dateTime) => {
    if (!dateTime) return false

    const date = new Date(dateTime)
    const now = new Date()
    const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000)

    return date >= now && date <= tomorrow
  }

  // Get appointment urgency color based on time
  const getAppointmentUrgency = (dateTime) => {
    if (!dateTime) return 'gray'

    const date = new Date(dateTime)
    const now = new Date()
    const diffInMinutes = (date - now) / (1000 * 60)

    if (diffInMinutes < 0) return 'red' // Past due
    if (diffInMinutes <= 5) return 'red' // Within 5 minutes
    if (diffInMinutes <= 15) return 'yellow' // Within 15 minutes
    if (diffInMinutes <= 30) return 'blue' // Within 30 minutes

    return 'gray' // More than 30 minutes away
  }

  return {
    timezone: computed(() => timezone.value),
    setTimezone,
    getTimezone,
    formatDateTime,
    formatDate,
    formatTime,
    toLocalTime,
    getRelativeTime,
    isToday,
    isUpcoming,
    getAppointmentUrgency
  }
}