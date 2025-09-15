export const validateRequired = (value, fieldName = 'This field') => {
  if (!value || (typeof value === 'string' && !value.trim())) {
    return `${fieldName} is required`
  }
  return null
}

export const validateEmail = (email) => {
  if (!email) return 'Email is required'
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email)) {
    return 'Please enter a valid email address'
  }
  return null
}

export const validatePhone = (phone) => {
  if (!phone) return null

  // Remove spaces for validation
  const cleanPhone = phone.replace(/\s+/g, '')

  // Philippine phone number patterns:
  // +639xxxxxxxxx (13 digits with +)
  // 639xxxxxxxxx (12 digits)
  // 09xxxxxxxxx (11 digits, local format)
  // 9xxxxxxxxx (10 digits, mobile only)
  const phoneRegex = /^(\+63[9]\d{9}|63[9]\d{9}|09[9]\d{8}|[9]\d{9})$/

  if (!phoneRegex.test(cleanPhone)) {
    return 'Please enter a valid Philippine phone number'
  }
  return null
}

export const validateDate = (date, fieldName = 'Date') => {
  if (!date) return `${fieldName} is required`
  const parsedDate = new Date(date)
  if (isNaN(parsedDate.getTime())) {
    return `Please enter a valid ${fieldName.toLowerCase()}`
  }
  return null
}

export const validateFutureDate = (date, fieldName = 'Date') => {
  const dateError = validateDate(date, fieldName)
  if (dateError) return dateError

  const selectedDate = new Date(date)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  if (selectedDate < today) {
    return `${fieldName} must be in the future`
  }
  return null
}

export const validatePastDate = (date, fieldName = 'Date') => {
  const dateError = validateDate(date, fieldName)
  if (dateError) return dateError

  const selectedDate = new Date(date)
  const today = new Date()
  today.setHours(23, 59, 59, 999)

  if (selectedDate > today) {
    return `${fieldName} cannot be in the future`
  }
  return null
}

export const validateMinLength = (value, minLength, fieldName = 'This field') => {
  if (!value) return null
  if (value.length < minLength) {
    return `${fieldName} must be at least ${minLength} characters long`
  }
  return null
}

export const validateMaxLength = (value, maxLength, fieldName = 'This field') => {
  if (!value) return null
  if (value.length > maxLength) {
    return `${fieldName} must be no more than ${maxLength} characters long`
  }
  return null
}

export const validateNumeric = (value, fieldName = 'This field') => {
  if (!value && value !== 0) return null
  if (isNaN(value) || isNaN(parseFloat(value))) {
    return `${fieldName} must be a valid number`
  }
  return null
}

export const validatePositive = (value, fieldName = 'This field') => {
  const numError = validateNumeric(value, fieldName)
  if (numError) return numError
  if (parseFloat(value) <= 0) {
    return `${fieldName} must be a positive number`
  }
  return null
}

export const validatePassword = (password) => {
  if (!password) return 'Password is required'
  if (password.length < 8) {
    return 'Password must be at least 8 characters long'
  }
  if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password)) {
    return 'Password must contain at least one uppercase letter, one lowercase letter, and one number'
  }
  return null
}

export const validatePasswordConfirm = (password, confirmPassword) => {
  if (!confirmPassword) return 'Please confirm your password'
  if (password !== confirmPassword) {
    return 'Passwords do not match'
  }
  return null
}

export const validateForm = (formData, validationRules) => {
  const errors = {}

  for (const [field, rules] of Object.entries(validationRules)) {
    const value = formData[field]

    for (const rule of rules) {
      const error = rule(value)
      if (error) {
        errors[field] = error
        break
      }
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  }
}