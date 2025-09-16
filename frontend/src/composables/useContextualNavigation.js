import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

// Global navigation context store
const navigationContext = ref(null)

export function useContextualNavigation() {
  const router = useRouter()
  const route = useRoute()

  // Set navigation context when navigating to a page
  const setNavigationContext = (context) => {
    navigationContext.value = context
  }

  // Navigate with context
  const navigateWithContext = (to, fromContext = null) => {
    const context = fromContext || {
      routeName: route.name,
      routePath: route.path,
      params: { ...route.params },
      query: { ...route.query },
      title: route.meta?.title || route.name
    }

    setNavigationContext(context)
    router.push(to)
  }

  // Get the appropriate back destination
  const getBackDestination = () => {
    // If we have navigation context, use it
    if (navigationContext.value) {
      return {
        path: navigationContext.value.routePath,
        title: navigationContext.value.title
      }
    }

    // Fallback to parent route logic
    const parentRouteName = route.meta?.parent
    if (parentRouteName) {
      const routeMap = {
        'Dashboard': { path: '/', title: 'Dashboard' },
        'PatientList': { path: '/patients', title: 'Patients' },
        'PatientView': { path: '/patients', title: 'Patients' },
        'AppointmentCalendar': { path: '/appointments', title: 'Appointments' },
        'AppointmentList': { path: '/appointments/list', title: 'Appointment List' },
        'ProcedureManagement': { path: '/procedures', title: 'Procedures' },
        'UserManagement': { path: '/users', title: 'Users' },
        'StaffManagement': { path: '/staff', title: 'Staff' },
        'ClinicManagement': { path: '/clinics', title: 'Clinics' },
        'ClinicEdit': { path: '/clinics', title: 'Clinics' },
        'PeerReviewList': { path: '/peer-review', title: 'Peer Review' },
        'InventoryList': { path: '/inventory', title: 'Inventory' },
        'InventoryView': { path: '/inventory', title: 'Inventory' },
        'Shop': { path: '/shop', title: 'Shop' },
        'SuperAdminShop': { path: '/admin/shop', title: 'Dentika Shop' },
        'SuperAdminShopView': { path: '/admin/shop', title: 'Dentika Shop' },
        'Notifications': { path: '/notifications', title: 'Notifications' }
      }

      return routeMap[parentRouteName] || null
    }

    return null
  }

  // Contextual back navigation
  const goBack = () => {
    const destination = getBackDestination()

    if (destination) {
      // Clear context after using it
      navigationContext.value = null
      router.push(destination.path)
    } else {
      // Fallback to browser back
      router.go(-1)
    }
  }

  // Get back button label
  const getBackButtonLabel = () => {
    const destination = getBackDestination()
    return destination ? `Back to ${destination.title}` : 'Back'
  }

  // Helper to determine if current page has contextual navigation
  const hasContextualNavigation = () => {
    return navigationContext.value !== null
  }

  // Clear navigation context (useful for direct navigation)
  const clearNavigationContext = () => {
    navigationContext.value = null
  }

  return {
    navigateWithContext,
    goBack,
    getBackButtonLabel,
    getBackDestination,
    hasContextualNavigation,
    clearNavigationContext,
    setNavigationContext
  }
}