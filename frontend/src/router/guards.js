import { useAuthStore } from '../stores/auth'

export function setupRouterGuards(router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // Initialize auth if not already done
    if (!authStore.initialized) {
      try {
        await authStore.initializeAuth()
      } catch (error) {
        console.error('Failed to initialize auth:', error)
        // If auth initialization fails, clear the token
        authStore.clearAuthState()
      }
    }

    const isAuthenticated = authStore.isAuthenticated
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
    const hideLayout = to.matched.some(record => record.meta.hideLayout)
    const requiresNonSuperAdmin = to.matched.some(record => record.meta.requiresNonSuperAdmin)
    const requiresSuperAdmin = to.matched.some(record => record.meta.requiresSuperAdmin)

    if (requiresAuth && !isAuthenticated) {
      next('/login')
    } else if ((to.name === 'Login' || to.name === 'Register') && isAuthenticated) {
      next('/')
    } else if (requiresNonSuperAdmin && isAuthenticated && authStore.user?.role === 'super_admin') {
      // Redirect super admins away from non-super-admin pages
      next('/')
    } else if (requiresSuperAdmin && isAuthenticated && authStore.user?.role !== 'super_admin') {
      // Redirect non-super-admins away from super-admin-only pages
      next('/')
    } else {
      next()
    }
  })
}