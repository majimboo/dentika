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