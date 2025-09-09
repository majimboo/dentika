<template>
  <div class="dashboard-page">
    <!-- Page Header -->
    <div class="page-header mb-8">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
          <p class="text-gray-600 mt-1">Welcome back, {{ userName }}. Here's what's happening at your clinic today.</p>
        </div>
        
        <div class="header-actions flex items-center space-x-3 mt-4 md:mt-0">
          <!-- Date Range Selector -->
          <div class="date-selector">
            <select
              v-model="selectedPeriod"
              @change="loadDashboardData"
              class="text-sm border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="today">Today</option>
              <option value="week">This Week</option>
              <option value="month">This Month</option>
              <option value="quarter">This Quarter</option>
              <option value="year">This Year</option>
            </select>
          </div>
          
          <!-- Refresh Button -->
          <button
            @click="refreshDashboard"
            :disabled="isLoading"
            class="btn btn-secondary flex items-center"
          >
            <svg 
              class="w-4 h-4 mr-2"
              :class="{ 'animate-spin': isLoading }"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Loading dashboard...</span>
    </div>

    <!-- Dashboard Content -->
    <div v-else class="dashboard-content space-y-8">
      <!-- Key Metrics Cards -->
      <div class="metrics-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Total Appointments -->
        <div class="metric-card bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="metric-icon w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value text-2xl font-bold text-gray-900">
                {{ dashboardData.totalAppointments || 0 }}
              </div>
              <div class="metric-label text-sm text-gray-600">Total Appointments</div>
              <div class="metric-change text-xs mt-1" :class="getChangeClass(dashboardData.appointmentsChange)">
                {{ formatChange(dashboardData.appointmentsChange) }} from last period
              </div>
            </div>
          </div>
        </div>

        <!-- Revenue -->
        <div class="metric-card bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="metric-icon w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value text-2xl font-bold text-gray-900">
                {{ formatCurrency(dashboardData.totalRevenue) }}
              </div>
              <div class="metric-label text-sm text-gray-600">Total Revenue</div>
              <div class="metric-change text-xs mt-1" :class="getChangeClass(dashboardData.revenueChange)">
                {{ formatChange(dashboardData.revenueChange) }} from last period
              </div>
            </div>
          </div>
        </div>

        <!-- New Patients -->
        <div class="metric-card bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="metric-icon w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value text-2xl font-bold text-gray-900">
                {{ dashboardData.newPatients || 0 }}
              </div>
              <div class="metric-label text-sm text-gray-600">New Patients</div>
              <div class="metric-change text-xs mt-1" :class="getChangeClass(dashboardData.patientsChange)">
                {{ formatChange(dashboardData.patientsChange) }} from last period
              </div>
            </div>
          </div>
        </div>

        <!-- Completion Rate -->
        <div class="metric-card bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="metric-icon w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value text-2xl font-bold text-gray-900">
                {{ formatPercentage(dashboardData.completionRate) }}%
              </div>
              <div class="metric-label text-sm text-gray-600">Completion Rate</div>
              <div class="metric-change text-xs mt-1" :class="getChangeClass(dashboardData.completionChange)">
                {{ formatChange(dashboardData.completionChange) }}% from last period
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="charts-grid grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Appointments Chart -->
        <div class="chart-card bg-white rounded-lg shadow-sm border p-6">
          <div class="chart-header flex items-center justify-between mb-6">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Appointments Overview</h3>
              <p class="text-sm text-gray-600">Daily appointment trends</p>
            </div>
            <div class="chart-legend flex items-center space-x-4 text-xs">
              <div class="legend-item flex items-center">
                <div class="w-3 h-3 bg-blue-500 rounded-full mr-1"></div>
                Scheduled
              </div>
              <div class="legend-item flex items-center">
                <div class="w-3 h-3 bg-green-500 rounded-full mr-1"></div>
                Completed
              </div>
            </div>
          </div>
          
          <!-- Simple Bar Chart Representation -->
          <div class="chart-content">
            <div class="chart-bars flex items-end justify-between h-40 space-x-2">
              <div 
                v-for="(day, index) in dashboardData.weeklyStats" 
                :key="index"
                class="chart-bar flex flex-col items-center flex-1"
              >
                <div class="bar-group flex space-x-1 items-end h-32">
                  <div 
                    class="bar bg-blue-500 w-3 rounded-t"
                    :style="{ height: `${(day.scheduled / maxDayAppointments) * 100}%` }"
                  ></div>
                  <div 
                    class="bar bg-green-500 w-3 rounded-t"
                    :style="{ height: `${(day.completed / maxDayAppointments) * 100}%` }"
                  ></div>
                </div>
                <div class="bar-label text-xs text-gray-600 mt-2">{{ day.label }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Revenue Chart -->
        <div class="chart-card bg-white rounded-lg shadow-sm border p-6">
          <div class="chart-header flex items-center justify-between mb-6">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Revenue Breakdown</h3>
              <p class="text-sm text-gray-600">Revenue by treatment type</p>
            </div>
          </div>
          
          <!-- Revenue Categories -->
          <div class="revenue-breakdown space-y-4">
            <div 
              v-for="category in dashboardData.revenueByCategory" 
              :key="category.name"
              class="category-item"
            >
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium text-gray-700">{{ category.name }}</span>
                <span class="text-sm text-gray-900 font-semibold">{{ formatCurrency(category.amount) }}</span>
              </div>
              <div class="progress-bar bg-gray-200 rounded-full h-2">
                <div 
                  class="progress-fill h-2 rounded-full transition-all duration-500"
                  :class="category.color"
                  :style="{ width: `${(category.amount / dashboardData.totalRevenue) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tables Row -->
      <div class="tables-grid grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Today's Appointments -->
        <div class="appointments-card bg-white rounded-lg shadow-sm border">
          <div class="card-header flex items-center justify-between p-6 border-b">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Today's Appointments</h3>
              <p class="text-sm text-gray-600">{{ todayAppointments.length }} appointments scheduled</p>
            </div>
            <router-link to="/appointments" class="text-blue-600 hover:text-blue-700 text-sm font-medium">
              View All →
            </router-link>
          </div>
          
          <div class="appointments-list max-h-80 overflow-y-auto">
            <div 
              v-for="appointment in todayAppointments.slice(0, 6)" 
              :key="appointment.id"
              class="appointment-item flex items-center p-4 border-b last:border-b-0 hover:bg-gray-50"
            >
              <div class="appointment-time flex-shrink-0">
                <div class="time-display text-sm font-medium text-gray-900">
                  {{ formatTime(appointment.start_time) }}
                </div>
                <div class="duration-display text-xs text-gray-500">
                  {{ calculateDuration(appointment.start_time, appointment.end_time) }}m
                </div>
              </div>
              
              <div class="appointment-details flex-1 ml-4">
                <div class="patient-name text-sm font-medium text-gray-900">
                  {{ appointment.patient?.first_name }} {{ appointment.patient?.last_name }}
                </div>
                <div class="appointment-type text-xs text-gray-600">
                  {{ formatAppointmentType(appointment.type) }}
                </div>
              </div>
              
              <div class="appointment-status">
                <span 
                  class="status-badge px-2 py-1 text-xs font-medium rounded-full"
                  :class="getStatusBadgeClass(appointment.status)"
                >
                  {{ formatStatus(appointment.status) }}
                </span>
              </div>
            </div>
            
            <div v-if="todayAppointments.length === 0" class="empty-state p-8 text-center">
              <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
              <p class="text-gray-500">No appointments scheduled for today</p>
            </div>
          </div>
        </div>

        <!-- Recent Patients -->
        <div class="patients-card bg-white rounded-lg shadow-sm border">
          <div class="card-header flex items-center justify-between p-6 border-b">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Recent Patients</h3>
              <p class="text-sm text-gray-600">Latest patient registrations</p>
            </div>
            <router-link to="/patients" class="text-blue-600 hover:text-blue-700 text-sm font-medium">
              View All →
            </router-link>
          </div>
          
          <div class="patients-list max-h-80 overflow-y-auto">
            <div 
              v-for="patient in recentPatients.slice(0, 6)" 
              :key="patient.id"
              class="patient-item flex items-center p-4 border-b last:border-b-0 hover:bg-gray-50"
            >
              <div class="patient-avatar flex-shrink-0">
                <div class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                  <span class="text-sm font-medium text-gray-600">
                    {{ getPatientInitials(patient) }}
                  </span>
                </div>
              </div>
              
              <div class="patient-details flex-1 ml-4">
                <div class="patient-name text-sm font-medium text-gray-900">
                  {{ patient.first_name }} {{ patient.last_name }}
                </div>
                <div class="patient-info text-xs text-gray-600">
                  {{ patient.phone }} • Registered {{ formatRelativeTime(patient.created_at) }}
                </div>
              </div>
              
              <div class="patient-actions">
                <button class="text-blue-600 hover:text-blue-700 text-xs font-medium">
                  View
                </button>
              </div>
            </div>
            
            <div v-if="recentPatients.length === 0" class="empty-state p-8 text-center">
              <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
              <p class="text-gray-500">No recent patients</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Alerts & Notifications -->
      <div class="alerts-section" v-if="alerts.length > 0">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Alerts & Notifications</h3>
        <div class="alerts-grid space-y-3">
          <div 
            v-for="alert in alerts" 
            :key="alert.id"
            class="alert-item flex items-start p-4 rounded-lg border"
            :class="getAlertClass(alert.type)"
          >
            <div class="alert-icon flex-shrink-0 mt-0.5">
              <svg 
                class="w-5 h-5"
                :class="getAlertIconClass(alert.type)"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getAlertIconPath(alert.type)"></path>
              </svg>
            </div>
            <div class="alert-content flex-1 ml-3">
              <div class="alert-title text-sm font-medium" :class="getAlertTextClass(alert.type)">
                {{ alert.title }}
              </div>
              <div class="alert-message text-sm mt-1" :class="getAlertMessageClass(alert.type)">
                {{ alert.message }}
              </div>
              <div class="alert-time text-xs mt-2 opacity-75">
                {{ formatRelativeTime(alert.created_at) }}
              </div>
            </div>
            <button 
              @click="dismissAlert(alert.id)"
              class="alert-dismiss flex-shrink-0 ml-3 p-1 hover:bg-black hover:bg-opacity-10 rounded"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useAppointmentStore } from '../stores/appointment'
import { usePatientStore } from '../stores/patient'
import apiService from '../services/api'

const authStore = useAuthStore()
const appointmentStore = useAppointmentStore()
const patientStore = usePatientStore()

const selectedPeriod = ref('today')
const isLoading = ref(false)

const dashboardData = ref({
  totalAppointments: 0,
  appointmentsChange: 0,
  totalRevenue: 0,
  revenueChange: 0,
  newPatients: 0,
  patientsChange: 0,
  completionRate: 0,
  completionChange: 0,
  weeklyStats: [
    { label: 'Mon', scheduled: 8, completed: 6 },
    { label: 'Tue', scheduled: 12, completed: 10 },
    { label: 'Wed', scheduled: 15, completed: 13 },
    { label: 'Thu', scheduled: 10, completed: 9 },
    { label: 'Fri', scheduled: 18, completed: 16 },
    { label: 'Sat', scheduled: 6, completed: 5 },
    { label: 'Sun', scheduled: 2, completed: 2 }
  ],
  revenueByCategory: [
    { name: 'Cleanings', amount: 2500, color: 'bg-blue-500' },
    { name: 'Fillings', amount: 1800, color: 'bg-green-500' },
    { name: 'Root Canals', amount: 3200, color: 'bg-yellow-500' },
    { name: 'Crowns', amount: 4500, color: 'bg-purple-500' },
    { name: 'Extractions', amount: 900, color: 'bg-red-500' }
  ]
})

const todayAppointments = ref([])
const recentPatients = ref([])
const alerts = ref([])

const userName = computed(() => {
  return authStore.user?.first_name || 'User'
})

const maxDayAppointments = computed(() => {
  if (dashboardData.value.weeklyStats.length === 0) return 10
  return Math.max(...dashboardData.value.weeklyStats.map(day => Math.max(day.scheduled || 0, day.completed || 0)))
})

const loadDashboardData = async () => {
  isLoading.value = true
  
  try {
    // Simulate API calls
    await Promise.all([
      loadMetrics(),
      loadTodayAppointments(),
      loadRecentPatients(),
      loadAlerts()
    ])
    
    // Update dashboard data based on selected period
    updateMetricsForPeriod()
  } catch (error) {
    console.error('Error loading dashboard data:', error)
  } finally {
    isLoading.value = false
  }
}

const loadMetrics = async () => {
  try {
    const response = await apiService.getDashboardMetrics(selectedPeriod.value)

    if (response.success) {
      const data = response.data

      // Ensure metrics object exists
      const metrics = data.metrics || {}

      // Update metrics from API response with safe defaults
      dashboardData.value.totalAppointments = metrics.total_appointments || 0
      dashboardData.value.totalRevenue = metrics.total_revenue || 0
      dashboardData.value.newPatients = metrics.new_patients || 0

      // Calculate completion rate from appointments
      dashboardData.value.completionRate = metrics.total_appointments > 0
        ? (metrics.completed_appointments / metrics.total_appointments) * 100
        : 0

      // Update weekly stats - handle both date-based and label-based formats
      dashboardData.value.weeklyStats = (data.weekly_stats || []).map(stat => ({
        label: stat.label || new Date(stat.date).toLocaleDateString('en-US', { weekday: 'short' }),
        scheduled: stat.scheduled || 0,
        completed: stat.completed || 0
      }))

      // Update revenue by category
      dashboardData.value.revenueByCategory = (data.revenue_by_type || []).map(item => ({
        name: item.type || 'Unknown',
        amount: item.amount || 0,
        color: item.color || 'bg-blue-500'
      }))

      // Calculate changes (simplified - in a real app you'd compare with previous period)
      dashboardData.value.appointmentsChange = 0
      dashboardData.value.revenueChange = 0
      dashboardData.value.patientsChange = 0
      dashboardData.value.completionChange = 0
    } else {
      console.error('Failed to load dashboard metrics:', response.error)
      // Set default values on error
      dashboardData.value.totalAppointments = 0
      dashboardData.value.totalRevenue = 0
      dashboardData.value.newPatients = 0
      dashboardData.value.completionRate = 0
    }
  } catch (error) {
    console.error('Error loading dashboard metrics:', error)
    // Set default values on error
    dashboardData.value.totalAppointments = 0
    dashboardData.value.totalRevenue = 0
    dashboardData.value.newPatients = 0
    dashboardData.value.completionRate = 0
  }
}

const loadTodayAppointments = async () => {
  // Load today's appointments
  await appointmentStore.fetchAppointments({
    date: new Date().toISOString().split('T')[0]
  })
  todayAppointments.value = appointmentStore.appointments
}

const loadRecentPatients = async () => {
  // Load recent patients
  await patientStore.fetchPatients({ limit: 10, sort: 'created_at' })
  recentPatients.value = patientStore.patients
}

const loadAlerts = async () => {
  // Simulate alerts
  alerts.value = [
    {
      id: 1,
      type: 'warning',
      title: 'Equipment Maintenance Due',
      message: 'X-ray machine requires maintenance check',
      created_at: new Date().toISOString()
    },
    {
      id: 2,
      type: 'info',
      title: 'New Patient Inquiry',
      message: '3 new patient inquiries received today',
      created_at: new Date().toISOString()
    }
  ]
}

const updateMetricsForPeriod = () => {
  // Metrics are now fetched from API based on selected period
  // No manual adjustment needed
}

const refreshDashboard = () => {
  loadDashboardData()
}

const dismissAlert = (alertId) => {
  alerts.value = alerts.value.filter(alert => alert.id !== alertId)
}

// Utility functions
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount || 0)
}

const formatPercentage = (value) => {
  return (value || 0).toFixed(1)
}

const formatChange = (change) => {
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change}%`
}

const getChangeClass = (change) => {
  return change >= 0 ? 'text-green-600' : 'text-red-600'
}

const formatTime = (timeString) => {
  return new Date(timeString).toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  })
}

const formatAppointmentType = (type) => {
  return type?.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()) || ''
}

const formatStatus = (status) => {
  const statusMap = {
    scheduled: 'Scheduled',
    confirmed: 'Confirmed',
    in_progress: 'In Progress',
    completed: 'Completed',
    cancelled: 'Cancelled',
    no_show: 'No Show'
  }
  return statusMap[status] || status
}

const getStatusBadgeClass = (status) => {
  const classes = {
    scheduled: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-red-100 text-red-800',
    no_show: 'bg-red-200 text-red-900'
  }
  return classes[status] || classes.scheduled
}

const calculateDuration = (startTime, endTime) => {
  const start = new Date(startTime)
  const end = new Date(endTime)
  return Math.round((end - start) / (1000 * 60))
}

const getPatientInitials = (patient) => {
  return `${patient.first_name?.[0] || ''}${patient.last_name?.[0] || ''}`.toUpperCase()
}

const formatRelativeTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now - date) / (1000 * 60 * 60))
  
  if (diffInHours < 1) return 'Just now'
  if (diffInHours < 24) return `${diffInHours}h ago`
  
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 7) return `${diffInDays}d ago`
  
  const diffInWeeks = Math.floor(diffInDays / 7)
  return `${diffInWeeks}w ago`
}

const getAlertClass = (type) => {
  const classes = {
    info: 'bg-blue-50 border-blue-200',
    warning: 'bg-yellow-50 border-yellow-200',
    error: 'bg-red-50 border-red-200',
    success: 'bg-green-50 border-green-200'
  }
  return classes[type] || classes.info
}

const getAlertIconClass = (type) => {
  const classes = {
    info: 'text-blue-500',
    warning: 'text-yellow-500',
    error: 'text-red-500',
    success: 'text-green-500'
  }
  return classes[type] || classes.info
}

const getAlertTextClass = (type) => {
  const classes = {
    info: 'text-blue-800',
    warning: 'text-yellow-800',
    error: 'text-red-800',
    success: 'text-green-800'
  }
  return classes[type] || classes.info
}

const getAlertMessageClass = (type) => {
  const classes = {
    info: 'text-blue-700',
    warning: 'text-yellow-700',
    error: 'text-red-700',
    success: 'text-green-700'
  }
  return classes[type] || classes.info
}

// Load dashboard data on component mount
onMounted(() => {
  loadDashboardData()
})

const getAlertIconPath = (type) => {
  const paths = {
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    warning: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z',
    error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
    success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
  }
  return paths[type] || paths.info
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
@import "../styles/main.css";

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.metric-card:hover {
  @apply shadow-md;
}

.chart-bar {
  @apply transition-all duration-300;
}

.progress-fill {
  @apply transition-all duration-500;
}

/* Custom scrollbars */
.appointments-list::-webkit-scrollbar,
.patients-list::-webkit-scrollbar {
  width: 4px;
}

.appointments-list::-webkit-scrollbar-track,
.patients-list::-webkit-scrollbar-track {
  @apply bg-gray-100;
}

.appointments-list::-webkit-scrollbar-thumb,
.patients-list::-webkit-scrollbar-thumb {
  @apply bg-gray-400 rounded;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .metrics-grid {
    @apply grid-cols-2;
  }
  
  .charts-grid,
  .tables-grid {
    @apply grid-cols-1;
  }
  
  .metric-card {
    @apply p-4;
  }
  
  .metric-icon {
    @apply w-10 h-10 mr-3;
  }
  
  .metric-value {
    @apply text-xl;
  }
}

@media (max-width: 640px) {
  .metrics-grid {
    @apply grid-cols-1;
  }
  
  .page-header {
    @apply flex-col items-start;
  }
  
  .header-actions {
    @apply w-full mt-4 flex-wrap;
  }
}
</style>