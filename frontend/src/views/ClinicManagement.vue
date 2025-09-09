<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-neutral-900">Clinic Management</h1>
        <p class="mt-2 text-base sm:text-lg text-neutral-600">Oversee all clinics and their branches.</p>
      </div>
      <div class="flex items-center space-x-4">
        <button 
          @click="loadClinics"
          :disabled="loading"
          class="inline-flex items-center justify-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          <span class="hidden sm:inline">Refresh</span>
        </button>
        <button 
          @click="addNewClinic"
          class="inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          <span class="hidden sm:inline">Add New Clinic</span>
        </button>
      </div>
    </div>

    <!-- Clinics List -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div v-if="loading" class="p-12 text-center">
        <p>Loading clinics...</p>
      </div>
      <div v-else-if="error" class="p-12 text-center text-danger-600">
        <p>Error loading clinics: {{ error }}</p>
        <button @click="loadClinics" class="mt-4 px-4 py-2 border rounded-lg">Try Again</button>
      </div>
      <div v-else-if="clinics.length === 0" class="p-12 text-center text-neutral-500">
        <p>No clinics found. Get started by adding a new one.</p>
      </div>
      <div v-else>
        <ul class="divide-y divide-neutral-200">
          <li v-for="clinic in clinics" :key="clinic.id" class="p-6 hover:bg-neutral-50 transition-colors">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <h3 class="text-lg font-semibold text-primary-700">{{ clinic.name }}</h3>
                <p class="text-sm text-neutral-600 mt-1">{{ clinic.address }}</p>
                <p class="text-sm text-neutral-500 mt-1">{{ clinic.phone }} | {{ clinic.email }}</p>
                <div class="mt-2">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {{ clinic.branches.length }} Branches
                  </span>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <button @click="editClinic(clinic.id)" class="px-3 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 hover:bg-neutral-100">Edit</button>
                <button @click="confirmDelete(clinic)" class="px-3 py-2 border border-danger-300 rounded-lg text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100">Delete</button>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="clinicToDelete" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-lg font-semibold">Delete Clinic</h3>
        <p class="mt-2 text-neutral-600">Are you sure you want to delete <strong>{{ clinicToDelete.name }}</strong>? This will also deactivate all associated users and branches. This action cannot be undone.</p>
        <div class="mt-6 flex justify-end space-x-3">
          <button @click="clinicToDelete = null" class="px-4 py-2 border rounded-lg">Cancel</button>
          <button @click="deleteClinic(clinicToDelete.id)" :disabled="deleting" class="px-4 py-2 border border-transparent rounded-lg text-white bg-danger-600 hover:bg-danger-700 disabled:opacity-50">
            {{ deleting ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClinicStore } from '../stores/clinic'
import { storeToRefs } from 'pinia'

export default {
  name: 'ClinicManagement',
  setup() {
    const router = useRouter()
    const clinicStore = useClinicStore()
    
    const { clinics, loading, error } = storeToRefs(clinicStore)
    
    const clinicToDelete = ref(null)
    const deleting = ref(false)

    const loadClinics = () => {
      clinicStore.fetchClinics()
    }

    const addNewClinic = () => {
      router.push('/clinics/new')
    }

    const editClinic = (clinicId) => {
      router.push(`/clinics/${clinicId}/edit`)
    }

    const confirmDelete = (clinic) => {
      clinicToDelete.value = clinic
    }

    const deleteClinic = async (clinicId) => {
      deleting.value = true
      const result = await clinicStore.deleteClinic(clinicId)
      if (result.success) {
        clinicToDelete.value = null
      } else {
        // Handle error, maybe show a toast notification
        console.error(result.error)
      }
      deleting.value = false
    }

    onMounted(() => {
      loadClinics()
    })

    return {
      clinics,
      loading,
      error,
      clinicToDelete,
      deleting,
      loadClinics,
      addNewClinic,
      editClinic,
      confirmDelete,
      deleteClinic,
    }
  }
}
</script>