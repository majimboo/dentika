<template>
  <footer class="bg-white/80 backdrop-blur-md border-t border-gray-200/50 mt-auto">
    <div class="px-6 py-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
      <div class="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <div class="h-6 w-6 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
              <svg class="h-3 w-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
            </div>
            <span class="text-sm font-semibold text-gray-900">Dentika</span>
          </div>
          <span class="text-sm text-gray-500">•</span>
          <span class="text-sm text-gray-600">© {{ currentYear }} All rights reserved.</span>
        </div>
        
        <div class="flex items-center space-x-6">
          <a href="#" class="text-sm text-gray-500 hover:text-gray-700 transition-colors duration-200 flex items-center space-x-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span>Help</span>
          </a>
          <a href="#" class="text-sm text-gray-500 hover:text-gray-700 transition-colors duration-200 flex items-center space-x-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.031 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
            <span>Privacy</span>
          </a>
          <a href="#" class="text-sm text-gray-500 hover:text-gray-700 transition-colors duration-200 flex items-center space-x-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            <span>Terms</span>
          </a>
          
          <div class="flex items-center space-x-2 text-sm text-gray-500">
            <div class="flex items-center space-x-1">
              <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
              <span>System Online</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Build Info -->
      <div class="mt-4 pt-4 border-t border-gray-200/50">
        <div class="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
          <p class="text-xs text-gray-400">
            Built with Vue.js, Go, and TailwindCSS • Version 1.0.0
          </p>
          <div class="flex items-center space-x-4 text-xs text-gray-400">
            <span class="flex items-center space-x-1">
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span>Last updated: {{ lastUpdated }}</span>
            </span>
          </div>
        </div>
      </div>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'AppFooter',
  data() {
    return {
      currentYear: new Date().getFullYear(),
      lastUpdated: new Date().toLocaleDateString()
    }
  }
}
</script>