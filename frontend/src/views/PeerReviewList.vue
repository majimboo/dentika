<template>
  <div class="peer-review-list">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Peer Review Cases</h1>
      <router-link
        to="/peer-review/create"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Create Case
      </router-link>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-lg shadow-sm border p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="filters.status"
            @change="fetchCases"
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All Status</option>
            <option value="open">Open</option>
            <option value="closed">Closed</option>
            <option value="resolved">Resolved</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Visibility</label>
          <select
            v-model="filters.visibility"
            @change="fetchCases"
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All Visibility</option>
            <option value="public">Public</option>
            <option value="in_clinic">In-Clinic</option>
            <option value="invite_only">Invite Only</option>
          </select>
        </div>

        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            v-model="filters.search"
            @input="debounceSearch"
            type="text"
            placeholder="Search cases..."
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
    </div>

    <!-- Cases List -->
    <div class="space-y-4">
      <div
        v-for="caseItem in cases"
        :key="caseItem.id"
        @click="viewCase(caseItem.id)"
        class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow cursor-pointer"
      >
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ caseItem.title }}</h3>
            <p class="text-sm text-gray-600">Case #{{ caseItem.case_number }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span
              :class="getStatusBadgeClass(caseItem.status)"
              class="px-2 py-1 text-xs font-medium rounded-full"
            >
              {{ caseItem.status }}
            </span>
            <span
              :class="getVisibilityBadgeClass(caseItem.visibility)"
              class="px-2 py-1 text-xs font-medium rounded-full"
            >
              {{ getVisibilityLabel(caseItem.visibility) }}
            </span>
          </div>
        </div>

        <p class="text-gray-700 mb-4 line-clamp-2">{{ caseItem.description }}</p>

        <div class="flex justify-between items-center text-sm text-gray-500">
          <div class="flex items-center gap-4">
            <span>By {{ caseItem.created_by.first_name }} {{ caseItem.created_by.last_name }}</span>
            <span>{{ caseItem.comments?.length || 0 }} comments</span>
          </div>
          <span>{{ formatDate(caseItem.created_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="pagination.total > pagination.limit" class="flex justify-center mt-8">
      <nav class="flex items-center gap-2">
        <button
          @click="changePage(pagination.page - 1)"
          :disabled="pagination.page <= 1"
          class="px-3 py-2 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>

        <span class="px-3 py-2 text-sm text-gray-700">
          Page {{ pagination.page }} of {{ Math.ceil(pagination.total / pagination.limit) }}
        </span>

        <button
          @click="changePage(pagination.page + 1)"
          :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)"
          class="px-3 py-2 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </nav>
    </div>


  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const cases = ref([])
const pagination = ref({
  page: 1,
  limit: 10,
  total: 0
})
const filters = ref({
  status: '',
  visibility: '',
  search: ''
})
let searchTimeout = null

const fetchCases = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page,
      limit: pagination.value.limit,
      ...filters.value
    })

    const response = await fetch(`/api/peer-review/cases?${params}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      cases.value = data.cases
      pagination.value = data.pagination
    }
  } catch (error) {
    console.error('Failed to fetch cases:', error)
  }
}

const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(fetchCases, 500)
}

const changePage = (page) => {
  pagination.value.page = page
  fetchCases()
}

const viewCase = (caseId) => {
  router.push(`/peer-review/${caseId}`)
}



const getStatusBadgeClass = (status) => {
  const classes = {
    open: 'bg-green-100 text-green-800',
    closed: 'bg-gray-100 text-gray-800',
    resolved: 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getVisibilityBadgeClass = (visibility) => {
  const classes = {
    public: 'bg-purple-100 text-purple-800',
    in_clinic: 'bg-orange-100 text-orange-800',
    invite_only: 'bg-red-100 text-red-800'
  }
  return classes[visibility] || 'bg-gray-100 text-gray-800'
}

const getVisibilityLabel = (visibility) => {
  const labels = {
    public: 'Public',
    in_clinic: 'In-Clinic',
    invite_only: 'Invite Only'
  }
  return labels[visibility] || visibility
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchCases()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>