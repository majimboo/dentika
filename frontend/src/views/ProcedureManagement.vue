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
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          New Procedure
        </router-link>
        
        <router-link 
          to="/diagnoses/new"
          class="btn btn-primary flex items-center"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
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
            <svg class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
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
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path>
            </svg>
          </button>
          <button
            @click="procedureViewMode = 'list'"
            :class="procedureViewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-600 hover:text-gray-800'"
            class="p-2 rounded-lg transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Procedures Grid/List View -->
      <div v-if="procedureViewMode === 'grid'" class="procedures-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="procedure in filteredProcedures" 
          :key="procedure.id"
          class="procedure-card bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow"
        >
          <div class="p-6">
            <div class="flex items-start justify-between mb-4">
              <div class="procedure-icon w-12 h-12 rounded-lg flex items-center justify-center mr-4"
                   :class="getProcedureCategoryColor(procedure.category)">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getProcedureIcon(procedure.category)"></path>
                </svg>
              </div>
              <div class="procedure-actions">
                <router-link
                  :to="`/procedures/${procedure.id}/edit`"
                  class="text-gray-400 hover:text-gray-600 transition-colors"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                </router-link>
              </div>
            </div>
            
            <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ procedure.name }}</h3>
            <p class="text-sm text-gray-600 mb-4 line-clamp-2">{{ procedure.description }}</p>
            
            <div class="procedure-details space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Duration</span>
                <span class="text-sm text-gray-900">{{ procedure.duration }} min</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Cost</span>
                <span class="text-sm text-gray-900 font-semibold">{{ procedure.cost }}</span>
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

      <!-- Procedures List View -->
      <div v-else class="procedures-list bg-white rounded-lg shadow-sm border overflow-hidden">
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
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getProcedureIcon(procedure.category)"></path>
                      </svg>
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
                <td class="px-6 py-4 text-sm text-gray-900">{{ procedure.duration }} min</td>
                <td class="px-6 py-4 text-sm font-semibold text-gray-900">{{ procedure.cost }}</td>
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
            <svg class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
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
      <div class="diagnoses-list bg-white rounded-lg shadow-sm border overflow-hidden">
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
                <td class="px-6 py-4 text-sm font-mono text-gray-900">{{ diagnosis.icd10_code }}</td>
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
        <div class="treatment-stats grid grid-cols-1 md:grid-cols-4 gap-4">
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
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          New Treatment Plan
        </router-link>
      </div>

      <!-- Treatment Plans List -->
      <div class="treatments-list bg-white rounded-lg shadow-sm border overflow-hidden">
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

const activeTab = ref('procedures')
const procedureViewMode = ref('grid')

// Search and filters
const procedureSearchQuery = ref('')
const diagnosisSearchQuery = ref('')
const selectedProcedureCategory = ref('')
const selectedDiagnosisCategory = ref('')


// Data
const procedures = ref([
  {
    id: 1,
    name: 'Dental Cleaning',
    description: 'Professional teeth cleaning and plaque removal',
    category: 'preventive',
    duration: 45,
    cost: 120,
    created_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: 'Root Canal Treatment',
    description: 'Endodontic treatment to save infected tooth',
    category: 'restorative',
    duration: 90,
    cost: 800,
    created_at: '2024-01-16T14:30:00Z'
  },
  {
    id: 3,
    name: 'Tooth Extraction',
    description: 'Surgical removal of damaged or problematic tooth',
    category: 'surgical',
    duration: 30,
    cost: 200,
    created_at: '2024-01-17T09:15:00Z'
  }
])

const diagnoses = ref([
  {
    id: 1,
    name: 'Dental Caries',
    description: 'Tooth decay affecting enamel and dentin',
    icd10_code: 'K02.9',
    category: 'caries',
    severity: 'moderate',
    created_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: 'Gingivitis',
    description: 'Inflammation of gums due to bacterial infection',
    icd10_code: 'K05.0',
    category: 'periodontal',
    severity: 'mild',
    created_at: '2024-01-16T14:30:00Z'
  }
])

const treatmentPlans = ref([
  {
    id: 1,
    name: 'Full Mouth Restoration',
    patient: {
      id: 1,
      first_name: 'John',
      last_name: 'Smith',
      phone: '(555) 123-4567'
    },
    procedures: [1, 2, 3],
    status: 'active',
    progress: 65,
    total_cost: 2500,
    created_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: 'Orthodontic Treatment',
    patient: {
      id: 2,
      first_name: 'Jane',
      last_name: 'Doe',
      phone: '(555) 987-6543'
    },
    procedures: [4, 5],
    status: 'on_hold',
    progress: 30,
    total_cost: 4500,
    created_at: '2024-01-16T14:30:00Z'
  }
])

const treatmentStats = ref({
  active: 12,
  completed: 45,
  onHold: 3,
  revenue: 125000
})

// Computed properties
const filteredProcedures = computed(() => {
  let filtered = procedures.value

  if (procedureSearchQuery.value) {
    const query = procedureSearchQuery.value.toLowerCase()
    filtered = filtered.filter(proc => 
      proc.name.toLowerCase().includes(query) ||
      proc.description.toLowerCase().includes(query)
    )
  }

  if (selectedProcedureCategory.value) {
    filtered = filtered.filter(proc => proc.category === selectedProcedureCategory.value)
  }

  return filtered
})

const filteredDiagnoses = computed(() => {
  let filtered = diagnoses.value

  if (diagnosisSearchQuery.value) {
    const query = diagnosisSearchQuery.value.toLowerCase()
    filtered = filtered.filter(diag => 
      diag.name.toLowerCase().includes(query) ||
      diag.description.toLowerCase().includes(query) ||
      diag.icd10_code.toLowerCase().includes(query)
    )
  }

  if (selectedDiagnosisCategory.value) {
    filtered = filtered.filter(diag => diag.category === selectedDiagnosisCategory.value)
  }

  return filtered
})

// Methods
const searchProcedures = () => {
  // Search is reactive via computed property
}

const searchDiagnoses = () => {
  // Search is reactive via computed property
}

const filterProcedures = () => {
  // Filter is reactive via computed property
}

const filterDiagnoses = () => {
  // Filter is reactive via computed property
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

const getProcedureIcon = (category) => {
  const icons = {
    preventive: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z',
    restorative: 'M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547A8.014 8.014 0 004 20h16a8.014 8.014 0 00-.244-4.572z',
    surgical: 'M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4',
    orthodontic: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253',
    cosmetic: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z',
    emergency: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z'
  }
  return icons[category] || 'M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547A8.014 8.014 0 004 20h16a8.014 8.014 0 00-.244-4.572z'
}

onMounted(() => {
  // Load initial data
  // This would typically fetch from APIs
})
</script>

<style scoped>
@reference "tailwindcss";

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