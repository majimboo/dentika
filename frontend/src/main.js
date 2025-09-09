import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import './styles/main.css'

import App from './App.vue'
import routes from './router/routes.js'
import { setupRouterGuards } from './router/guards.js'

const router = createRouter({
  history: createWebHistory(),
  routes
})

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)

setupRouterGuards(router)

app.mount('#app')