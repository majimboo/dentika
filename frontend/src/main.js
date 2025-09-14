import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import './styles/main.css'

// FontAwesome setup - all icons loaded at once
import { FontAwesomeIcon } from './plugins/fontawesome.js'

import App from './App.vue'
import routes from './router/routes.js'
import { setupRouterGuards } from './router/guards.js'

const router = createRouter({
  history: createWebHistory(),
  routes
})

const pinia = createPinia()
const app = createApp(App)

app.component('font-awesome-icon', FontAwesomeIcon)
app.use(pinia)
app.use(router)

setupRouterGuards(router)

app.mount('#app')

// Hide loading screen when app is ready
const loadingElement = document.getElementById('loading')
if (loadingElement) {
  // Add a small delay for smooth transition
  setTimeout(() => {
    loadingElement.style.opacity = '0'
    loadingElement.style.transition = 'opacity 0.5s ease-out'
    setTimeout(() => {
      loadingElement.style.display = 'none'
    }, 500)
  }, 100)
}