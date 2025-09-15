export const debounce = (func, wait) => {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func.apply(this, args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

export const throttle = (func, limit) => {
  let inThrottle
  return function executedFunction(...args) {
    if (!inThrottle) {
      func.apply(this, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

export const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
}

export const capitalizeWords = (str) => {
  if (!str) return ''
  return str.split(' ').map(word => capitalize(word)).join(' ')
}

export const truncateText = (text, maxLength = 50, suffix = '...') => {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength - suffix.length) + suffix
}

export const generateRandomId = (prefix = 'id') => {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

export const deepClone = (obj) => {
  if (obj === null || typeof obj !== 'object') return obj
  if (obj instanceof Date) return new Date(obj.getTime())
  if (obj instanceof Array) return obj.map(item => deepClone(item))
  if (typeof obj === 'object') {
    const clonedObj = {}
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key])
      }
    }
    return clonedObj
  }
}

export const isEmpty = (value) => {
  if (value == null) return true
  if (typeof value === 'string') return value.trim().length === 0
  if (Array.isArray(value)) return value.length === 0
  if (typeof value === 'object') return Object.keys(value).length === 0
  return false
}

export const formatCurrency = (amount, currency = 'PHP', locale = 'en-US') => {
  if (amount == null || isNaN(amount)) return 'N/A'
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: currency
  }).format(amount)
}

export const formatPercentage = (value, decimals = 1) => {
  if (value == null || isNaN(value)) return 'N/A'
  return `${(value * 100).toFixed(decimals)}%`
}

export const sortByKey = (array, key, direction = 'asc') => {
  return [...array].sort((a, b) => {
    const aVal = a[key]
    const bVal = b[key]

    if (aVal < bVal) return direction === 'asc' ? -1 : 1
    if (aVal > bVal) return direction === 'asc' ? 1 : -1
    return 0
  })
}

export const groupBy = (array, key) => {
  return array.reduce((groups, item) => {
    const group = item[key]
    groups[group] = groups[group] || []
    groups[group].push(item)
    return groups
  }, {})
}

export const removeFromArray = (array, item) => {
  const index = array.indexOf(item)
  if (index > -1) {
    array.splice(index, 1)
  }
  return array
}

export const formatPhoneNumber = (phoneNumber) => {
  if (!phoneNumber) return ''

  // Remove any existing spaces and normalize
  const cleaned = phoneNumber.replace(/\s+/g, '')

  // Handle Philippine phone numbers (+63xxxxxxxxx format)
  if (cleaned.startsWith('+63') && cleaned.length === 13) {
    // Format as: +63 917 123 4567
    const numberPart = cleaned.substring(3) // Remove +63
    return `+63 ${numberPart.substring(0, 3)} ${numberPart.substring(3, 6)} ${numberPart.substring(6)}`
  }

  // Handle other international formats (keep +XX format with spaces)
  if (cleaned.startsWith('+')) {
    const countryCode = cleaned.substring(0, cleaned.length - 10) // Extract country code
    const numberPart = cleaned.substring(cleaned.length - 10) // Last 10 digits
    return `${countryCode} ${numberPart.substring(0, 3)} ${numberPart.substring(3, 6)} ${numberPart.substring(6)}`
  }

  // Handle local formats without country code
  if (cleaned.length === 10 && cleaned.startsWith('9')) {
    // Format as: 917 123 4567
    return `${cleaned.substring(0, 3)} ${cleaned.substring(3, 6)} ${cleaned.substring(6)}`
  }

  if (cleaned.length === 11 && cleaned.startsWith('09')) {
    // Format as: 0917 123 4567
    return `${cleaned.substring(0, 4)} ${cleaned.substring(4, 7)} ${cleaned.substring(7)}`
  }

  // Return as-is for other formats
  return phoneNumber
}

export const handleAsyncError = async (operation, errorCallback = null) => {
  try {
    return await operation()
  } catch (error) {
    console.error('Async operation failed:', error)
    if (errorCallback) {
      errorCallback(error)
    }
    throw error
  }
}