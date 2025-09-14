export const formatDate = (dateString, options = {}) => {
  if (!dateString) return options.fallback || 'N/A'

  try {
    const date = new Date(dateString)

    if (options.relative) {
      const today = new Date()
      const tomorrow = new Date(today)
      tomorrow.setDate(tomorrow.getDate() + 1)
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)

      if (date.toDateString() === today.toDateString()) {
        return 'today'
      } else if (date.toDateString() === tomorrow.toDateString()) {
        return 'tomorrow'
      } else if (date.toDateString() === yesterday.toDateString()) {
        return 'yesterday'
      }
    }

    const formatOptions = options.format || {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }

    return date.toLocaleDateString('en-US', formatOptions)
  } catch (error) {
    console.error('Error formatting date:', error)
    return options.fallback || 'Invalid Date'
  }
}

export const formatTime = (dateTime, options = {}) => {
  if (!dateTime) return options.fallback || 'N/A'

  try {
    const date = new Date(dateTime)

    if (options.relative) {
      const now = new Date()
      const diffMs = now - date
      const diffMins = Math.floor(diffMs / (1000 * 60))
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

      if (diffMins < 1) return 'just now'
      if (diffMins < 60) return `${diffMins}m ago`
      if (diffHours < 24) return `${diffHours}h ago`
      if (diffDays < 7) return `${diffDays}d ago`
    }

    const formatOptions = options.format || {
      hour: 'numeric',
      minute: '2-digit'
    }

    return date.toLocaleTimeString('en-US', formatOptions)
  } catch (error) {
    console.error('Error formatting time:', error)
    return options.fallback || 'Invalid Time'
  }
}

export const formatDateTime = (dateTime, options = {}) => {
  if (!dateTime) return options.fallback || 'N/A'

  const dateStr = formatDate(dateTime, { format: options.dateFormat })
  const timeStr = formatTime(dateTime, { format: options.timeFormat })

  return `${dateStr} ${timeStr}`
}

export const calculateAge = (birthDate) => {
  if (!birthDate) return null

  const today = new Date()
  const birth = new Date(birthDate)
  let age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()

  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--
  }

  return age
}

export const isToday = (date) => {
  const today = new Date()
  const checkDate = new Date(date)
  return checkDate.toDateString() === today.toDateString()
}

export const isTomorrow = (date) => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const checkDate = new Date(date)
  return checkDate.toDateString() === tomorrow.toDateString()
}

export const isYesterday = (date) => {
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)
  const checkDate = new Date(date)
  return checkDate.toDateString() === yesterday.toDateString()
}