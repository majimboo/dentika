import { useAuthStore } from '../stores/auth'

export function setupRouterGuards(router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    const isAuthenticated = authStore.isAuthenticated
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
    const hideLayout = to.matched.some(record => record.meta.hideLayout)
    const requiresNonSuperAdmin = to.matched.some(record => record.meta.requiresNonSuperAdmin)

    if (requiresAuth && !isAuthenticated) {
      next('/login')
    } else if ((to.name === 'Login' || to.name === 'Register') && isAuthenticated) {
      next('/')
    } else if (requiresNonSuperAdmin && isAuthenticated && authStore.user?.role === 'super_admin') {
      // Redirect super admins away from non-super-admin pages
      next('/')
    } else {
      next()
    }
  })
}