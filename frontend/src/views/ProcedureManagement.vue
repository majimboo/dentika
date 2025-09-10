<template>
  <div class="procedure-management-page">
    <!-- Page Header -->
    <div class="page-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Procedure & Diagnosis Management</h1>
        <p class="text-gray-600">Manage dental procedures, diagnoses, and treatment plans</p>
      </div>
      
      <div class="header-actions flex items-center space-x-3">
        <router-link
          to="/procedures/new"
          class="btn btn-secondary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          New Procedure
        </router-link>
        
        <router-link
          to="/diagnoses/new"
          class="btn btn-primary flex items-center"
        >
          <font-awesome-icon icon="fa-solid fa-file-medical" class="w-4 h-4 mr-2" />
          New Diagnosis
        </router-link>
      </div>
    </div>

    <!-- Tab Navigation -->
    <div class="tabs-container mb-6">
      <div class="tab-navigation flex border-b border-gray-200">
        <button
          @click="activeTab = 'procedures'"
          :class="activeTab === 'procedures' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
          class="tab-button py-2 px-4 border-b-2 font-medium text-sm transition-colors"
        >
          Procedures
        </button>
        <button
          @click="activeTab = 'diagnoses'"
          :class="activeTab === 'diagnoses' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
          class="tab-button py-2 px-4 border-b-2 font-medium text-sm transition-colors"
        >
          Diagnoses
        </button>
        <button
          @click="activeTab = 'treatments'"
          :class="activeTab === 'treatments' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
          class="tab-button py-2 px-4 border-b-2 font-medium text-sm transition-colors"
        >
          Treatment Plans
        </button>
      </div>
    </div>

    <!-- Procedures Tab -->
    <div v-if="activeTab === 'procedures'" class="procedures-section">
      <div class="procedures-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <div class="search-and-filters flex items-center space-x-4">
           <div class="search-box relative">
             <input
               v-model="procedureSearchQuery"
               @input="searchProcedures"
               type="text"
               placeholder="Search procedures..."
               class="w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
             />
             <font-awesome-icon icon="fa-solid fa-search" class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" />
           </div>

          <select
            v-model="selectedProcedureCategory"
            @change="filterProcedures"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Categories</option>
            <option value="preventive">Preventive</option>
            <option value="restorative">Restorative</option>
            <option value="surgical">Surgical</option>
            <option value="orthodontic">Orthodontic</option>
            <option value="cosmetic">Cosmetic</option>
            <option value="emergency">Emergency</option>
          </select>
        </div>
        
        <div class="view-options flex items-center space-x-2">
           <button
             @click="procedureViewMode = 'grid'"
             :class="procedureViewMode === 'grid' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
             class="p-2 rounded-lg transition-colors"
           >
             <font-awesome-icon icon="fa-solid fa-th" class="w-5 h-5" />
           </button>
           <button
             @click="procedureViewMode = 'list'"
             :class="procedureViewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
             class="p-2 rounded-lg transition-colors"
           >
             <font-awesome-icon icon="fa-solid fa-list" class="w-5 h-5" />
           </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="proceduresLoading" class="loading-state flex items-center justify-center py-12 bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading procedures...</span>
      </div>

      <!-- Procedures Grid/List View -->
      <div v-else-if="procedureViewMode === 'grid'" class="procedures-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="procedure in filteredProcedures"
          :key="procedure.id"
          class="procedure-card bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
        >
          <div class="p-6">
            <div class="flex items-start justify-between mb-4">
               <div class="procedure-icon w-12 h-12 rounded-lg flex items-center justify-center mr-4"
                    :class="getProcedureCategoryColor(procedure.category)">
                 <font-awesome-icon icon="fa-solid fa-tooth" class="w-6 h-6" />
               </div>
              <div class="procedure-actions">
                 <router-link
                   :to="`/procedures/${procedure.id}/edit`"
                   class="text-gray-400 hover:text-gray-600 transition-colors"
                 >
                   <font-awesome-icon icon="fa-solid fa-edit" class="w-5 h-5" />
                 </router-link>
              </div>
            </div>
            
            <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ procedure.name }}</h3>
            <p class="text-sm text-gray-600 mb-4 line-clamp-2">{{ procedure.description }}</p>

            <div class="procedure-details space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Duration</span>
                <span class="text-sm text-gray-900">{{ procedure.estimated_duration }} min</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Cost</span>
                <span class="text-sm text-gray-900 font-semibold">{{ procedure.default_cost }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Category</span>
                <span class="text-xs px-2 py-1 rounded-full" :class="getProcedureCategoryBadge(procedure.category)">
                  {{ formatCategory(procedure.category) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State for Procedures -->
      <div v-else-if="filteredProcedures.length === 0" class="empty-state text-center py-12 bg-white rounded-lg shadow-sm border border-gray-200">
         <div class="empty-icon text-gray-300 mb-4">
           <font-awesome-icon icon="fa-solid fa-file-alt" class="w-16 h-16 mx-auto" />
         </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No procedures found</h3>
        <p class="text-gray-500 mb-4">
          {{ procedureSearchQuery ? 'No procedures match your search criteria.' : 'Get started by adding your first procedure.' }}
        </p>
      </div>

      <!-- Procedures List View -->
      <div v-else class="procedures-list bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="procedures-table">
          <table class="w-full">
            <thead class="bg-gray-50 border-b">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Procedure</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Duration</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cost</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="procedure in filteredProcedures" 
                :key="procedure.id"
                class="hover:bg-gray-50 cursor-pointer"
                @click="viewProcedure(procedure)"
              >
                <td class="px-6 py-4">
                  <div class="flex items-center">
                     <div class="procedure-icon w-8 h-8 rounded-lg flex items-center justify-center mr-3"
                          :class="getProcedureCategoryColor(procedure.category)">
                       <font-awesome-icon icon="fa-solid fa-tooth" class="w-4 h-4" />
                     </div>
                     <div>
                       <div class="text-sm font-medium text-gray-900">{{ procedure.name }}</div>
                       <div class="text-sm text-gray-500 line-clamp-1">{{ procedure.description }}</div>
                     </div>
                   </div>
                 </td>
                 <td class="px-6 py-4">
                   <span class="text-xs px-2 py-1 rounded-full" :class="getProcedureCategoryBadge(procedure.category)">
                     {{ formatCategory(procedure.category) }}
                   </span>
                 </td>
                 <td class="px-6 py-4 text-sm text-gray-900">{{ procedure.estimated_duration }} min</td>
                 <td class="px-6 py-4 text-sm font-semibold text-gray-900">{{ procedure.default_cost }}</td>
                <td class="px-6 py-4">
                  <div class="flex items-center space-x-2">
                    <router-link
                      :to="`/procedures/${procedure.id}/edit`"
                      class="text-blue-600 hover:text-blue-700 text-sm font-medium"
                    >
                      Edit
                    </router-link>
                    <button 
                      @click.stop="deleteProcedure(procedure.id)"
                      class="text-red-600 hover:text-red-700 text-sm font-medium"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Diagnoses Tab -->
    <div v-if="activeTab === 'diagnoses'" class="diagnoses-section">
      <div class="diagnoses-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <div class="search-and-filters flex items-center space-x-4">
           <div class="search-box relative">
             <input
               v-model="diagnosisSearchQuery"
               @input="searchDiagnoses"
               type="text"
               placeholder="Search diagnoses..."
               class="w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
             />
             <font-awesome-icon icon="fa-solid fa-search" class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" />
           </div>

          <select
            v-model="selectedDiagnosisCategory"
            @change="filterDiagnoses"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Categories</option>
            <option value="caries">Caries</option>
            <option value="periodontal">Periodontal</option>
            <option value="endodontic">Endodontic</option>
            <option value="orthodontic">Orthodontic</option>
            <option value="oral_surgery">Oral Surgery</option>
            <option value="other">Other</option>
          </select>
        </div>
      </div>

      <!-- Diagnoses List -->
      <div class="diagnoses-list bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="diagnoses-table">
          <table class="w-full">
            <thead class="bg-gray-50 border-b">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Diagnosis</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ICD-10 Code</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Severity</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="diagnosis in filteredDiagnoses" 
                :key="diagnosis.id"
                class="hover:bg-gray-50"
              >
                <td class="px-6 py-4">
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ diagnosis.name }}</div>
                    <div class="text-sm text-gray-500 line-clamp-1">{{ diagnosis.description }}</div>
                  </div>
                </td>
                 <td class="px-6 py-4 text-sm font-mono text-gray-900">{{ diagnosis.code }}</td>
                <td class="px-6 py-4">
                  <span class="text-xs px-2 py-1 rounded-full" :class="getDiagnosisCategoryBadge(diagnosis.category)">
                    {{ formatCategory(diagnosis.category) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <span class="text-xs px-2 py-1 rounded-full" :class="getSeverityBadge(diagnosis.severity)">
                    {{ formatSeverity(diagnosis.severity) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center space-x-2">
                    <router-link
                      :to="`/diagnoses/${diagnosis.id}/edit`"
                      class="text-blue-600 hover:text-blue-700 text-sm font-medium"
                    >
                      Edit
                    </router-link>
                    <button 
                      @click="deleteDiagnosis(diagnosis.id)"
                      class="text-red-600 hover:text-red-700 text-sm font-medium"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Treatment Plans Tab -->
    <div v-if="activeTab === 'treatments'" class="treatments-section">
      <div class="treatments-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <div class="treatment-stats grid grid-cols-1 md:grid-cols-4 gap-4 w-full">
          <div class="stat-card bg-blue-50 rounded-lg p-4">
            <div class="text-sm font-medium text-blue-600">Active Plans</div>
            <div class="text-2xl font-bold text-blue-900">{{ treatmentStats.active }}</div>
          </div>
          <div class="stat-card bg-green-50 rounded-lg p-4">
            <div class="text-sm font-medium text-green-600">Completed</div>
            <div class="text-2xl font-bold text-green-900">{{ treatmentStats.completed }}</div>
          </div>
          <div class="stat-card bg-yellow-50 rounded-lg p-4">
            <div class="text-sm font-medium text-yellow-600">On Hold</div>
            <div class="text-2xl font-bold text-yellow-900">{{ treatmentStats.onHold }}</div>
          </div>
          <div class="stat-card bg-gray-50 rounded-lg p-4">
            <div class="text-sm font-medium text-gray-600">Total Revenue</div>
            <div class="text-2xl font-bold text-gray-900">{{ formatCurrency(treatmentStats.revenue) }}</div>
          </div>
        </div>
        
        <router-link
          to="/treatments/new"
          class="btn btn-primary flex items-center whitespace-nowrap"
        >
           <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
          New Treatment Plan
        </router-link>
      </div>

      <!-- Treatment Plans List -->
      <div class="treatments-list bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="treatments-table">
          <table class="w-full">
            <thead class="bg-gray-50 border-b">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Patient</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Treatment</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Progress</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Cost</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="treatment in treatmentPlans" 
                :key="treatment.id"
                class="hover:bg-gray-50"
              >
                <td class="px-6 py-4">
                  <div class="flex items-center">
                    <div class="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center mr-3">
                      <span class="text-xs font-medium text-gray-600">
                        {{ getPatientInitials(treatment.patient) }}
                      </span>
                    </div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">
                        {{ treatment.patient?.first_name }} {{ treatment.patient?.last_name }}
                      </div>
                      <div class="text-sm text-gray-500">{{ treatment.patient?.phone }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm font-medium text-gray-900">{{ treatment.name }}</div>
                  <div class="text-sm text-gray-500">{{ treatment.procedures.length }} procedures</div>
                </td>
                <td class="px-6 py-4">
                  <span class="text-xs px-2 py-1 rounded-full" :class="getTreatmentStatusBadge(treatment.status)">
                    {{ formatTreatmentStatus(treatment.status) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center">
                    <div class="flex-1 bg-gray-200 rounded-full h-2 mr-2">
                      <div 
                        class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                        :style="{ width: `${treatment.progress}%` }"
                      ></div>
                    </div>
                    <span class="text-xs text-gray-600">{{ treatment.progress }}%</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-sm font-semibold text-gray-900">{{ formatCurrency(treatment.total_cost) }}</td>
                <td class="px-6 py-4">
                  <div class="flex items-center space-x-2">
                    <router-link
                      :to="`/treatments/${treatment.id}`"
                      class="text-blue-600 hover:text-blue-700 text-sm font-medium"
                    >
                      View
                    </router-link>
                    <router-link
                      :to="`/treatments/${treatment.id}/edit`"
                      class="text-green-600 hover:text-green-700 text-sm font-medium"
                    >
                      Edit
                    </router-link>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import apiService from '../services/api'

const activeTab = ref('procedures')
const procedureViewMode = ref('grid')

// Search and filters
const procedureSearchQuery = ref('')
const diagnosisSearchQuery = ref('')
const selectedProcedureCategory = ref('')
const selectedDiagnosisCategory = ref('')


// Data
const procedures = ref([])
const diagnoses = ref([])
const treatmentPlans = ref([])
const treatmentStats = ref({
  active: 0,
  completed: 0,
  onHold: 0,
  revenue: 0
})

// Loading states
const proceduresLoading = ref(false)
const diagnosesLoading = ref(false)
const treatmentsLoading = ref(false)

// Computed properties
const filteredProcedures = computed(() => {
  // Since we're now loading filtered data from API, just return the procedures
  return procedures.value
})

const filteredDiagnoses = computed(() => {
  // Since we're now loading filtered data from API, just return the diagnoses
  return diagnoses.value
})

// Methods
const loadProcedures = async () => {
  try {
    proceduresLoading.value = true
    const params = {}
    if (procedureSearchQuery.value) {
      params.search = procedureSearchQuery.value
    }
    if (selectedProcedureCategory.value) {
      params.category = selectedProcedureCategory.value
    }

    const result = await apiService.getProcedureTemplates(params)
    if (result.success) {
      procedures.value = result.data
    }
  } catch (error) {
    console.error('Error loading procedures:', error)
  } finally {
    proceduresLoading.value = false
  }
}

const loadDiagnoses = async () => {
  try {
    diagnosesLoading.value = true
    const params = {}
    if (diagnosisSearchQuery.value) {
      params.search = diagnosisSearchQuery.value
    }
    if (selectedDiagnosisCategory.value) {
      params.category = selectedDiagnosisCategory.value
    }

    const result = await apiService.getDiagnosisTemplates(params)
    if (result.success) {
      diagnoses.value = result.data
    }
  } catch (error) {
    console.error('Error loading diagnoses:', error)
  } finally {
    diagnosesLoading.value = false
  }
}

const loadTreatmentPlans = async () => {
  // For now, treatment plans are not implemented in the backend
  // This would load treatment plans from a future API endpoint
  treatmentsLoading.value = false
}

const searchProcedures = () => {
  loadProcedures()
}

const searchDiagnoses = () => {
  loadDiagnoses()
}

const filterProcedures = () => {
  loadProcedures()
}

const filterDiagnoses = () => {
  loadDiagnoses()
}

const deleteProcedure = (procedureId) => {
  if (confirm('Are you sure you want to delete this procedure?')) {
    procedures.value = procedures.value.filter(p => p.id !== procedureId)
  }
}


const deleteDiagnosis = (diagnosisId) => {
  if (confirm('Are you sure you want to delete this diagnosis?')) {
    diagnoses.value = diagnoses.value.filter(d => d.id !== diagnosisId)
  }
}



// Utility functions
const formatCategory = (category) => {
  return category.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatSeverity = (severity) => {
  return severity.charAt(0).toUpperCase() + severity.slice(1)
}

const formatTreatmentStatus = (status) => {
  return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount || 0)
}

const getPatientInitials = (patient) => {
  return `${patient.first_name?.[0] || ''}${patient.last_name?.[0] || ''}`.toUpperCase()
}

const getProcedureCategoryColor = (category) => {
  const colors = {
    preventive: 'bg-green-100 text-green-600',
    restorative: 'bg-blue-100 text-blue-600',
    surgical: 'bg-red-100 text-red-600',
    orthodontic: 'bg-purple-100 text-purple-600',
    cosmetic: 'bg-pink-100 text-pink-600',
    emergency: 'bg-yellow-100 text-yellow-600'
  }
  return colors[category] || 'bg-gray-100 text-gray-600'
}

const getProcedureCategoryBadge = (category) => {
  const badges = {
    preventive: 'bg-green-100 text-green-800',
    restorative: 'bg-blue-100 text-blue-800',
    surgical: 'bg-red-100 text-red-800',
    orthodontic: 'bg-purple-100 text-purple-800',
    cosmetic: 'bg-pink-100 text-pink-800',
    emergency: 'bg-yellow-100 text-yellow-800'
  }
  return badges[category] || 'bg-gray-100 text-gray-800'
}

const getDiagnosisCategoryBadge = (category) => {
  const badges = {
    caries: 'bg-red-100 text-red-800',
    periodontal: 'bg-orange-100 text-orange-800',
    endodontic: 'bg-purple-100 text-purple-800',
    orthodontic: 'bg-blue-100 text-blue-800',
    oral_surgery: 'bg-gray-100 text-gray-800',
    other: 'bg-gray-100 text-gray-800'
  }
  return badges[category] || 'bg-gray-100 text-gray-800'
}

const getSeverityBadge = (severity) => {
  const badges = {
    mild: 'bg-green-100 text-green-800',
    moderate: 'bg-yellow-100 text-yellow-800',
    severe: 'bg-red-100 text-red-800'
  }
  return badges[severity] || 'bg-gray-100 text-gray-800'
}

const getTreatmentStatusBadge = (status) => {
  const badges = {
    active: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    on_hold: 'bg-yellow-100 text-yellow-800',
    cancelled: 'bg-red-100 text-red-800'
  }
  return badges[status] || 'bg-gray-100 text-gray-800'
}



onMounted(() => {
  loadProcedures()
  loadDiagnoses()
  loadTreatmentPlans()
})
</script>

<style scoped>
@import "../styles/main.css";

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.tab-button:hover {
  @apply border-gray-300;
}

.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}

.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

/* Custom scrollbar */
.procedures-table::-webkit-scrollbar,
.diagnoses-table::-webkit-scrollbar,
.treatments-table::-webkit-scrollbar {
  width: 6px;
}

.procedures-table::-webkit-scrollbar-track,
.diagnoses-table::-webkit-scrollbar-track,
.treatments-table::-webkit-scrollbar-track {
  @apply bg-gray-100 rounded;
}

.procedures-table::-webkit-scrollbar-thumb,
.diagnoses-table::-webkit-scrollbar-thumb,
.treatments-table::-webkit-scrollbar-thumb {
  @apply bg-gray-400 rounded;
}

.procedures-table::-webkit-scrollbar-thumb:hover,
.diagnoses-table::-webkit-scrollbar-thumb:hover,
.treatments-table::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-500;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .procedures-grid {
    @apply grid-cols-1;
  }
  
  .treatment-stats {
    @apply grid-cols-2;
  }
  
  .search-and-filters {
    @apply flex-col items-start space-x-0 space-y-4;
  }
  
  .search-box input {
    @apply w-full;
  }
}

@media (max-width: 640px) {
  .treatment-stats {
    @apply grid-cols-1;
  }
  
  .procedures-table,
  .diagnoses-table,
  .treatments-table {
    @apply text-sm;
  }
}
</style>