import { useAuthStore } from '../stores/auth'

export function setupRouterGuards(router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    
    const isAuthenticated = authStore.isAuthenticated
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
    const hideLayout = to.matched.some(record => record.meta.hideLayout)
    
    if (requiresAuth && !isAuthenticated) {
      next('/login')
    } else if ((to.name === 'Login' || to.name === 'Register') && isAuthenticated) {
      next('/')
    } else {
      next()
    }
  })
}