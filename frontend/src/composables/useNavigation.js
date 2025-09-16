import { useRoute, useRouter } from 'vue-router'
import { useContextualNavigation } from './useContextualNavigation'

export function useNavigation() {
  const route = useRoute()
  const router = useRouter()
  const { goBack: contextualGoBack, navigateWithContext } = useContextualNavigation()

  const goBack = () => {
    // Use contextual navigation which handles both context-aware and fallback navigation
    contextualGoBack()
  }

  const getParentRouteName = () => {
    const parentName = route.meta?.parent
    if (!parentName) return null

    // Convert route names to display names
    const displayNames = {
      'Dashboard': 'Dashboard',
      'PatientList': 'Patients',
      'PatientView': 'Patients',
      'AppointmentCalendar': 'Appointments',
      'AppointmentList': 'Appointment List',
      'ProcedureManagement': 'Procedures',
      'UserManagement': 'Users',
      'StaffManagement': 'Staff',
      'ClinicManagement': 'Clinics',
      'ClinicEdit': 'Clinics',
      'PeerReviewList': 'Peer Review',
      'InventoryList': 'Inventory',
      'InventoryView': 'Inventory',
      'Shop': 'Shop',
      'SuperAdminShop': 'Dentika Shop',
      'SuperAdminShopView': 'Dentika Shop'
    }

    return displayNames[parentName] || parentName
  }

  return {
    goBack,
    getParentRouteName,
    navigateWithContext
  }
}