import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import './styles/main.css'

// FontAwesome setup
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
  faBars,
  faBolt,
  faBell,
  faChevronDown,
  faUser,
  faCog,
  faSignOutAlt,
  faHome,
  faCheckSquare,
  faUsers,
  faCalendarAlt,
  faStethoscope,
  faUserMd,
  faPlus,
  faChartBar,
  faQuestionCircle,
  faTimes,
  faCheckCircle,
  faExclamationTriangle,
  faInfoCircle,
  faClock,
  faCalendar,
  faUser as faUserSolid,
  faStickyNote,
  faSync,
  faSpinner,
  faBuilding,
  faEllipsisH,
  faDollarSign,
  faUserPlus,
  faSearch,
  faTh,
  faList,
  faLock,
  faEye,
  faEyeSlash,
  faShieldAlt,
  faExclamationCircle,
  faEnvelope,
  faCheck,
  faBoxes
} from '@fortawesome/free-solid-svg-icons'

library.add(
  faBars,
  faBolt,
  faBell,
  faChevronDown,
  faUser,
  faCog,
  faSignOutAlt,
  faHome,
  faCheckSquare,
  faUsers,
  faCalendarAlt,
  faStethoscope,
  faUserMd,
  faPlus,
  faChartBar,
  faQuestionCircle,
  faTimes,
  faCheckCircle,
  faExclamationTriangle,
  faInfoCircle,
  faClock,
  faCalendar,
  faUserSolid,
  faStickyNote,
  faSync,
  faSpinner,
  faBuilding,
  faEllipsisH,
  faDollarSign,
  faUserPlus,
  faSearch,
  faTh,
  faList,
  faLock,
  faEye,
  faEyeSlash,
  faShieldAlt,
  faExclamationCircle,
  faEnvelope,
  faCheck,
  faBoxes
)

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