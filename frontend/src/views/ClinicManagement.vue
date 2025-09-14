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
           v-if="isSuperAdmin"
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

     <!-- Clinics Table -->
     <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
       <!-- Loading State -->
       <div v-if="loading" class="p-12 text-center">
         <BaseLoading
           type="spinner"
           size="large"
           color="primary"
           text="Loading clinics..."
         />
         <p class="text-neutral-600 mt-4">Please wait while we fetch the clinic data.</p>
       </div>

       <!-- Error State -->
       <div v-else-if="error" class="p-12 text-center">
         <div class="inline-flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mb-4">
           <svg class="w-8 h-8 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
           </svg>
         </div>
         <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading clinics</h3>
         <p class="text-danger-600 mb-4">{{ error }}</p>
         <button
           @click="loadClinics"
           class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
         >
           Try Again
         </button>
       </div>

       <!-- Empty State -->
       <div v-else-if="clinics.length === 0" class="p-12 text-center">
         <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
           <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
           </svg>
         </div>
         <h3 class="text-lg font-semibold text-neutral-900 mb-2">No clinics yet</h3>
         <p class="text-neutral-600 mb-6">Get started by creating your first clinic to manage your dental practice.</p>
         <button
           @click="addNewClinic"
           class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
         >
           <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
           </svg>
           Create First Clinic
         </button>
       </div>

       <!-- Clinics Grid -->
       <div v-else class="p-4 md:p-6">
         <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
           <div
             v-for="clinic in clinics"
             :key="clinic.id"
             class="group bg-neutral-50 hover:bg-white border border-neutral-200 hover:border-primary-300 rounded-xl p-6 transition-all duration-300 hover:shadow-lg hover:scale-[1.02]"
           >
             <div class="flex items-start justify-between">
               <div class="flex-1 min-w-0">
                 <!-- Clinic Header -->
                 <div class="flex items-center space-x-3 mb-3">
                   <div class="flex-shrink-0">
                     <div class="h-10 w-10 rounded-xl bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center text-white font-semibold text-sm">
                       {{ getInitials(clinic.name) }}
                     </div>
                   </div>
                   <div class="flex-1 min-w-0">
                     <h3 class="text-lg font-semibold text-neutral-900 group-hover:text-primary-700 transition-colors truncate">
                       {{ clinic.name }}
                     </h3>
                     <div class="flex items-center space-x-2 mt-1">
                       <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-success-100 text-success-800">
                         Active
                       </span>
                       <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                         {{ clinic.branches?.length || 0 }} Branches
                       </span>
                     </div>
                   </div>
                 </div>

                 <!-- Clinic Details -->
                 <div class="space-y-2 mb-4">
                   <div v-if="clinic.address" class="flex items-start space-x-2">
                     <svg class="w-4 h-4 text-neutral-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                     </svg>
                     <p class="text-sm text-neutral-600">{{ clinic.address }}</p>
                   </div>

                   <div class="flex items-center space-x-4">
                     <div v-if="clinic.phone" class="flex items-center space-x-2">
                       <svg class="w-4 h-4 text-neutral-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
                       </svg>
                       <p class="text-sm text-neutral-600">{{ clinic.phone }}</p>
                     </div>

                     <div v-if="clinic.email" class="flex items-center space-x-2">
                       <svg class="w-4 h-4 text-neutral-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                       </svg>
                       <p class="text-sm text-neutral-600">{{ clinic.email }}</p>
                     </div>
                   </div>

                   <div v-if="clinic.website" class="flex items-center space-x-2">
                     <svg class="w-4 h-4 text-neutral-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
                     </svg>
                     <a :href="clinic.website" target="_blank" class="text-sm text-primary-600 hover:text-primary-700 transition-colors">
                       {{ clinic.website }}
                     </a>
                   </div>
                 </div>

                 <!-- Metadata -->
                 <div class="flex items-center justify-between text-xs text-neutral-400">
                   <span>ID: {{ clinic.id }}</span>
                   <span v-if="clinic.created_at">Created {{ formatDate(clinic.created_at) }}</span>
                 </div>
               </div>

               <!-- Actions -->
               <div class="flex items-center space-x-2 ml-4 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity duration-200">
                 <button
                   @click="editClinic(clinic.id)"
                   class="inline-flex items-center px-3 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
                   title="Edit clinic"
                 >
                   <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                   </svg>
                   Edit
                 </button>
                 <button
                   @click="confirmDelete(clinic)"
                   class="inline-flex items-center px-3 py-2 border border-danger-300 rounded-lg text-sm font-medium text-danger-700 bg-danger-50 hover:bg-danger-100 focus:outline-none focus:ring-2 focus:ring-danger-500 transition-all duration-200"
                   title="Delete clinic"
                 >
                   <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                   </svg>
                   Delete
                 </button>
               </div>
             </div>
           </div>
         </div>
       </div>
     </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="clinicToDelete" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-lg font-semibold">Delete Clinic</h3>
        <p class="mt-2 text-neutral-600">Are you sure you want to delete <strong>{{ clinicToDelete.name }}</strong>? This will also deactivate all associated users and branches. This action cannot be undone.</p>
        <div class="mt-6 flex justify-end space-x-3">
          <button @click="clinicToDelete = null" class="px-4 py-2 border border-gray-300 rounded-lg">Cancel</button>
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
import { useAuthStore } from '../stores/auth'
import { storeToRefs } from 'pinia'
import BaseLoading from '../components/BaseLoading.vue'

export default {
  name: 'ClinicManagement',
  components: {
    BaseLoading
  },
  setup() {
    const router = useRouter()
    const clinicStore = useClinicStore()
    const authStore = useAuthStore()

    const { clinics, loading, error } = storeToRefs(clinicStore)
    const { isSuperAdmin } = storeToRefs(authStore)
    
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

    const getInitials = (name) => {
      if (!name) return 'C'
      return name.charAt(0).toUpperCase()
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'Unknown'
      try {
        return new Date(dateString).toLocaleDateString()
      } catch {
        return 'Unknown'
      }
    }

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
      getInitials,
      formatDate,
      isSuperAdmin,
    }
  }
}
</script>