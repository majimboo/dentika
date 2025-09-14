export const getInitials = (user, fallback = 'U') => {
  if (!user) return fallback

  if (typeof user === 'string') {
    return user.charAt(0).toUpperCase()
  }

  const firstName = user.first_name || user.firstName || ''
  const lastName = user.last_name || user.lastName || ''

  if (firstName && lastName) {
    return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase()
  } else if (firstName) {
    return firstName.charAt(0).toUpperCase()
  } else if (user.username) {
    return user.username.charAt(0).toUpperCase()
  } else if (user.name) {
    const nameParts = user.name.trim().split(' ')
    if (nameParts.length >= 2) {
      return (nameParts[0].charAt(0) + nameParts[nameParts.length - 1].charAt(0)).toUpperCase()
    }
    return user.name.charAt(0).toUpperCase()
  }

  return fallback
}

export const getFullName = (user, fallback = 'Unknown User') => {
  if (!user) return fallback

  if (typeof user === 'string') return user

  const firstName = user.first_name || user.firstName || ''
  const lastName = user.last_name || user.lastName || ''

  if (firstName && lastName) {
    return `${firstName} ${lastName}`
  } else if (firstName) {
    return firstName
  } else if (user.name) {
    return user.name
  } else if (user.username) {
    return user.username
  }

  return fallback
}

export const getUserDisplayName = (user, fallback = 'Unknown') => {
  const fullName = getFullName(user, null)
  return fullName || user?.username || user?.email?.split('@')[0] || fallback
}