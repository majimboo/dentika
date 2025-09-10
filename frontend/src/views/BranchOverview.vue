<template>
  <div class="branch-overview-page">
    <!-- Page Header -->
    <div class="page-header mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Branch Overview</h1>
          <p class="text-gray-600 mt-1">Clinic-wide data across all branches</p>
        </div>

        <div class="header-actions flex items-center space-x-3">
          <!-- Date Range Selector -->
          <div class="date-selector">
            <select
              v-model="selectedPeriod"
              @change="loadOverviewData"
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
            @click="refreshOverview"
            :disabled="isLoading"
            class="btn btn-secondary flex items-center"
          >
            <font-awesome-icon
              icon="fa-solid fa-sync"
              class="w-4 h-4 mr-2"
              :class="{ 'animate-spin': isLoading }"
            />
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Loading overview data...</span>
    </div>

    <!-- Overview Content -->
    <div v-else class="overview-content space-y-8">
      <!-- Clinic Summary Cards -->
      <div class="summary-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Total Appointments -->
        <div class="summary-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="summary-icon w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
              <font-awesome-icon icon="fa-solid fa-calendar-alt" class="w-6 h-6 text-blue-600" />
            </div>
            <div class="summary-content">
              <div class="summary-value text-2xl font-bold text-gray-900">
                {{ overviewData.totalAppointments || 0 }}
              </div>
              <div class="summary-label text-sm text-gray-600">Total Appointments</div>
            </div>
          </div>
        </div>

        <!-- Total Revenue -->
        <div class="summary-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="summary-icon w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mr-4">
              <font-awesome-icon icon="fa-solid fa-peso-sign" class="w-6 h-6 text-green-600" />
            </div>
            <div class="summary-content">
              <div class="summary-value text-2xl font-bold text-gray-900">
                {{ formatCurrency(overviewData.totalRevenue) }}
              </div>
              <div class="summary-label text-sm text-gray-600">Total Revenue</div>
            </div>
          </div>
        </div>

        <!-- Active Branches -->
        <div class="summary-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="summary-icon w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
              <font-awesome-icon icon="fa-solid fa-building" class="w-6 h-6 text-purple-600" />
            </div>
            <div class="summary-content">
              <div class="summary-value text-2xl font-bold text-gray-900">
                {{ activeBranches.length }}
              </div>
              <div class="summary-label text-sm text-gray-600">Active Branches</div>
            </div>
          </div>
        </div>

        <!-- Average Completion Rate -->
        <div class="summary-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="summary-icon w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center mr-4">
              <font-awesome-icon icon="fa-solid fa-chart-line" class="w-6 h-6 text-yellow-600" />
            </div>
            <div class="summary-content">
              <div class="summary-value text-2xl font-bold text-gray-900">
                {{ formatPercentage(overviewData.averageCompletionRate) }}%
              </div>
              <div class="summary-label text-sm text-gray-600">Avg Completion Rate</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Branch Performance Table -->
      <div class="branch-performance bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="card-header flex items-center justify-between p-6 border-b border-gray-200">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Branch Performance</h3>
            <p class="text-sm text-gray-600">Performance metrics by branch</p>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Branch
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Appointments
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Revenue
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Completion Rate
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="branch in branchPerformance" :key="branch.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0">
                      <font-awesome-icon
                        :icon="branch.is_main_branch ? 'fa-solid fa-star' : 'fa-solid fa-building'"
                        class="w-5 h-5 text-gray-400"
                        :class="branch.is_main_branch ? 'text-yellow-500' : ''"
                      />
                    </div>
                    <div class="ml-3">
                      <div class="text-sm font-medium text-gray-900">
                        {{ branch.name }}
                      </div>
                      <div class="text-sm text-gray-500">
                        {{ branch.address }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ branch.appointments || 0 }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatCurrency(branch.revenue || 0) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="text-sm text-gray-900 mr-2">
                      {{ formatPercentage(branch.completionRate || 0) }}%
                    </div>
                    <div class="w-16 bg-gray-200 rounded-full h-2">
                      <div
                        class="bg-blue-600 h-2 rounded-full"
                        :style="{ width: `${branch.completionRate || 0}%` }"
                      ></div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
                    :class="branch.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  >
                    {{ branch.is_active ? 'Active' : 'Inactive' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Branch Comparison Charts -->
      <div class="charts-grid grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Appointments by Branch -->
        <div class="chart-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="chart-header mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Appointments by Branch</h3>
            <p class="text-sm text-gray-600">Distribution of appointments across branches</p>
          </div>

          <div class="chart-content">
            <div class="space-y-4">
              <div
                v-for="branch in branchPerformance.slice(0, 5)"
                :key="branch.id"
                class="branch-chart-item"
              >
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-700">{{ branch.name }}</span>
                  <span class="text-sm text-gray-900">{{ branch.appointments || 0 }}</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-3">
                  <div
                    class="bg-blue-600 h-3 rounded-full transition-all duration-500"
                    :style="{ width: `${getBranchPercentage(branch.appointments, 'appointments')}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Revenue by Branch -->
        <div class="chart-card bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div class="chart-header mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Revenue by Branch</h3>
            <p class="text-sm text-gray-600">Revenue distribution across branches</p>
          </div>

          <div class="chart-content">
            <div class="space-y-4">
              <div
                v-for="branch in branchPerformance.slice(0, 5)"
                :key="branch.id"
                class="branch-chart-item"
              >
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-700">{{ branch.name }}</span>
                  <span class="text-sm text-gray-900">{{ formatCurrency(branch.revenue || 0) }}</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-3">
                  <div
                    class="bg-green-600 h-3 rounded-full transition-all duration-500"
                    :style="{ width: `${getBranchPercentage(branch.revenue, 'revenue')}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useClinicStore } from '../stores/clinic'
import apiService from '../services/api'

const clinicStore = useClinicStore()

const selectedPeriod = ref('month')
const isLoading = ref(false)

const overviewData = ref({
  totalAppointments: 0,
  totalRevenue: 0,
  averageCompletionRate: 0
})

const branchPerformance = ref([])

const activeBranches = computed(() => clinicStore.getActiveBranches)

const loadOverviewData = async () => {
  isLoading.value = true

  try {
    // Load clinic-wide data across all branches
    const response = await apiService.getClinicOverview(selectedPeriod.value)

    if (response.success) {
      const data = response.data

      overviewData.value = {
        totalAppointments: data.total_appointments || 0,
        totalRevenue: data.total_revenue || 0,
        averageCompletionRate: data.average_completion_rate || 0
      }

      // Load branch-specific performance data
      branchPerformance.value = data.branch_performance || []
    } else {
      console.error('Failed to load overview data:', response.error)
    }
  } catch (error) {
    console.error('Error loading overview data:', error)
  } finally {
    isLoading.value = false
  }
}

const refreshOverview = () => {
  loadOverviewData()
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'PHP',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount || 0)
}

const formatPercentage = (value) => {
  return (value || 0).toFixed(1)
}

const getBranchPercentage = (value, type) => {
  if (!branchPerformance.value.length) return 0

  const total = branchPerformance.value.reduce((sum, branch) => {
    return sum + (branch[type] || 0)
  }, 0)

  if (total === 0) return 0

  return ((value || 0) / total) * 100
}

// Load overview data on component mount
onMounted(() => {
  loadOverviewData()
})

export default {
  name: 'BranchOverview',
  setup() {
    return {
      selectedPeriod,
      isLoading,
      overviewData,
      branchPerformance,
      activeBranches,
      loadOverviewData,
      refreshOverview,
      formatCurrency,
      formatPercentage,
      getBranchPercentage
    }
  }
}
</script>

<style scoped>
@import "../styles/main.css";

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.summary-card:hover {
  @apply shadow-md;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .summary-grid {
    @apply grid-cols-2;
  }

  .charts-grid {
    @apply grid-cols-1;
  }

  .summary-card {
    @apply p-4;
  }

  .summary-icon {
    @apply w-10 h-10 mr-3;
  }

  .summary-value {
    @apply text-xl;
  }
}

@media (max-width: 640px) {
  .summary-grid {
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