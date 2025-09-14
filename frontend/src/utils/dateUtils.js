export const formatDate = (dateString, options = {}) => {
  if (!dateString) return options.fallback || 'N/A'

  try {
    const date = new Date(dateString)

    if (options.relative) {
      const today = new Date()
      const todayInTimezone = new Date(today.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
      const tomorrow = new Date(todayInTimezone)
      tomorrow.setDate(tomorrow.getDate() + 1)
      const yesterday = new Date(todayInTimezone)
      yesterday.setDate(yesterday.getDate() - 1)

      const dateInTimezone = new Date(date.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))

      if (dateInTimezone.toDateString() === todayInTimezone.toDateString()) {
        return 'today'
      } else if (dateInTimezone.toDateString() === tomorrow.toDateString()) {
        return 'tomorrow'
      } else if (dateInTimezone.toDateString() === yesterday.toDateString()) {
        return 'yesterday'
      }
    }

    const formatOptions = {
      timeZone: 'Asia/Manila',
      ...(options.format || {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }

    return new Intl.DateTimeFormat('en-US', formatOptions).format(date)
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
      const nowInTimezone = new Date(now.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
      const dateInTimezone = new Date(date.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
      const diffMs = nowInTimezone - dateInTimezone
      const diffMins = Math.floor(diffMs / (1000 * 60))
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

      if (diffMins < 1) return 'just now'
      if (diffMins < 60) return `${diffMins}m ago`
      if (diffHours < 24) return `${diffHours}h ago`
      if (diffDays < 7) return `${diffDays}d ago`
    }

    const formatOptions = {
      timeZone: 'Asia/Manila',
      ...(options.format || {
        hour: 'numeric',
        minute: '2-digit'
      })
    }

    return new Intl.DateTimeFormat('en-US', formatOptions).format(date)
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
  const todayInTimezone = new Date(today.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  const checkDate = new Date(date)
  const checkDateInTimezone = new Date(checkDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  return checkDateInTimezone.toDateString() === todayInTimezone.toDateString()
}

export const isTomorrow = (date) => {
  const tomorrow = new Date()
  const tomorrowInTimezone = new Date(tomorrow.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  tomorrowInTimezone.setDate(tomorrowInTimezone.getDate() + 1)
  const checkDate = new Date(date)
  const checkDateInTimezone = new Date(checkDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  return checkDateInTimezone.toDateString() === tomorrowInTimezone.toDateString()
}

export const isYesterday = (date) => {
  const yesterday = new Date()
  const yesterdayInTimezone = new Date(yesterday.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  yesterdayInTimezone.setDate(yesterdayInTimezone.getDate() - 1)
  const checkDate = new Date(date)
  const checkDateInTimezone = new Date(checkDate.toLocaleString('en-US', { timeZone: 'Asia/Manila' }))
  return checkDateInTimezone.toDateString() === yesterdayInTimezone.toDateString()
}