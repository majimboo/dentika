export const appointmentStatusClasses = {
  scheduled: 'bg-blue-100 text-blue-800',
  confirmed: 'bg-green-100 text-green-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-gray-100 text-gray-800',
  cancelled: 'bg-red-100 text-red-800',
  no_show: 'bg-red-100 text-red-800',
  pending: 'bg-orange-100 text-orange-800'
}

export const diagnosisStatusClasses = {
  draft: 'bg-gray-100 text-gray-800',
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800'
}

export const treatmentStatusClasses = {
  planned: 'bg-blue-100 text-blue-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  completed: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
  on_hold: 'bg-orange-100 text-orange-800'
}

export const inventoryStatusClasses = {
  in_stock: 'bg-green-100 text-green-800',
  low_stock: 'bg-yellow-100 text-yellow-800',
  out_of_stock: 'bg-red-100 text-red-800',
  discontinued: 'bg-gray-100 text-gray-800'
}

export const getStatusClass = (status, type = 'appointment') => {
  const statusMaps = {
    appointment: appointmentStatusClasses,
    diagnosis: diagnosisStatusClasses,
    treatment: treatmentStatusClasses,
    inventory: inventoryStatusClasses
  }

  const statusMap = statusMaps[type] || appointmentStatusClasses
  return statusMap[status] || 'bg-gray-100 text-gray-800'
}

export const getAppointmentStatusClass = (status) => {
  return getStatusClass(status, 'appointment')
}

export const getDiagnosisStatusClass = (status) => {
  return getStatusClass(status, 'diagnosis')
}

export const getTreatmentStatusClass = (status) => {
  return getStatusClass(status, 'treatment')
}

export const getInventoryStatusClass = (status) => {
  return getStatusClass(status, 'inventory')
}

export const formatStatus = (status) => {
  if (!status) return ''
  return status.split('_').map(word =>
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}