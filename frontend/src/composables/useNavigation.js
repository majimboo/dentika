import { useRoute, useRouter } from 'vue-router'

export function useNavigation() {
  const route = useRoute()
  const router = useRouter()

  const goBack = () => {
    const parentRouteName = route.meta?.parent

    if (parentRouteName) {
      // Navigate to the specific parent route
      const routeMap = {
        'Dashboard': '/',
        'PatientList': '/patients',
        'AppointmentCalendar': '/appointments',
        'ProcedureManagement': '/procedures',
        'UserList': '/users',
        'ClinicManagement': '/clinics',
        'StaffManagement': '/staff',
        'PeerReviewList': '/peer-review',
        'InventoryList': '/inventory'
      }

      const parentPath = routeMap[parentRouteName]
      if (parentPath) {
        router.push(parentPath)
        return
      }
    }

    // Fallback to browser back
    router.go(-1)
  }

  const getParentRouteName = () => {
    const parentName = route.meta?.parent
    if (!parentName) return null

    // Convert route names to display names
    const displayNames = {
      'Dashboard': 'Dashboard',
      'PatientList': 'Patients',
      'AppointmentCalendar': 'Appointments',
      'ProcedureManagement': 'Procedures',
      'UserList': 'Users',
      'ClinicManagement': 'Clinics',
      'StaffManagement': 'Staff',
      'PeerReviewList': 'Peer Review',
      'InventoryList': 'Inventory'
    }

    return displayNames[parentName] || parentName
  }

  return {
    goBack,
    getParentRouteName
  }
}