<template>
  <div class="space-y-8">
    <!-- Loading State -->
    <BaseTransition name="fade">
      <div v-if="loading" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <BaseLoading type="spinner" size="large" color="primary" text="Loading clinic data..." />
        <p class="text-neutral-600 mt-4">Please wait while we fetch the clinic information.</p>
      </div>
    </BaseTransition>

    <!-- Error State -->
    <BaseTransition name="fade">
      <div v-if="!loading && error" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mb-4">
          <svg class="w-8 h-8 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading clinic</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button 
          @click="loadClinicData"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
        >
          Try Again
        </button>
      </div>
    </BaseTransition>

    <!-- Edit Form -->
    <BaseTransition name="slide-up">
      <div v-if="!loading && !error" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="p-6 sm:p-8">
          <form @submit.prevent="saveClinic" class="space-y-8">
            
            <!-- Clinic Details -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Clinic Details</h3>
              
              <div class="space-y-2">
                <label for="name" class="block text-sm font-semibold text-gray-700">Clinic Name <span class="text-danger-500 ml-1">*</span></label>
                <input type="text" id="name" v-model="clinic.name" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white" required>
              </div>
              <div class="space-y-2">
                <label for="address" class="block text-sm font-semibold text-gray-700">Address</label>
                <input type="text" id="address" v-model="clinic.address" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white">
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="phone" class="block text-sm font-semibold text-gray-700">Phone</label>
                  <input type="text" id="phone" v-model="clinic.phone" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white">
                </div>
                <div class="space-y-2">
                  <label for="email" class="block text-sm font-semibold text-gray-700">Email</label>
                  <input type="email" id="email" v-model="clinic.email" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white">
                </div>
              </div>
              <div class="space-y-2">
                <label for="website" class="block text-sm font-semibold text-gray-700">Website</label>
                <input type="url" id="website" v-model="clinic.website" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white">
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-6 border-t border-neutral-200">
              <button
                type="submit"
                :disabled="saving"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <BaseLoading v-if="saving" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                {{ saving ? 'Saving...' : 'Save Clinic' }}
              </button>
              
              <button
                type="button"
                @click="cancel"
                :disabled="saving"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </BaseTransition>

    <!-- Branch Management -->
    <BaseTransition name="slide-up">
      <div v-if="isEditMode && !loading && !error" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="p-6 sm:p-8">
          <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2 mb-4">Branches</h3>
           <ul class="space-y-3">
             <li v-for="branch in clinic.branches" :key="branch.id" class="border border-neutral-200 rounded-xl overflow-hidden shadow-sm bg-white">
               <!-- Branch Header -->
               <div class="p-4 flex justify-between items-start bg-white">
                 <div class="flex-1 min-w-0">
                   <div class="flex items-center space-x-2 mb-1">
                     <p class="font-semibold text-neutral-900 truncate">{{ branch.name }}</p>
                     <span v-if="branch.is_main_branch" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full font-medium">Main</span>
                   </div>
                   <p class="text-sm text-neutral-600 leading-relaxed">{{ branch.address || 'No address specified' }}</p>
                 </div>
                 <div class="flex items-center space-x-1 ml-3">
                   <button
                     @click="editBranch(branch)"
                     class="p-2 text-neutral-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200 touch-manipulation"
                     :class="{ 'text-blue-600 bg-blue-50': editingBranch && editingBranch.id === branch.id }"
                   >
                     <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                     </svg>
                   </button>
                   <button
                     @click="deleteBranch(branch.id)"
                     class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200 touch-manipulation"
                   >
                     <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                     </svg>
                   </button>
                 </div>
               </div>

               <!-- Inline Edit Form -->
               <BaseTransition name="slide-down">
                 <div v-if="editingBranch && editingBranch.id === branch.id" class="border-t border-neutral-200 bg-neutral-50 p-4">
                   <form @submit.prevent="saveBranch" class="space-y-3">
                     <div>
                       <label class="block text-xs font-medium text-neutral-600 mb-1">Branch Name <span class="text-danger-500">*</span></label>
                       <input
                         type="text"
                         v-model="newBranch.name"
                         class="block w-full px-3 py-2 text-sm border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                         placeholder="Enter branch name"
                         required
                       >
                     </div>
                     <div>
                       <label class="block text-xs font-medium text-neutral-600 mb-1">Address</label>
                       <input
                         type="text"
                         v-model="newBranch.address"
                         class="block w-full px-3 py-2 text-sm border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                         placeholder="Enter branch address"
                       >
                     </div>
                     <div class="flex justify-end space-x-2 pt-2">
                       <button type="button" @click="cancelBranchForm" class="px-3 py-2 text-xs font-medium text-neutral-600 bg-white border border-neutral-300 rounded-lg hover:bg-neutral-50 transition-colors">
                         Cancel
                       </button>
                       <button type="submit" :disabled="savingBranch" class="px-3 py-2 text-xs font-medium text-white bg-primary-600 hover:bg-primary-700 disabled:opacity-50 rounded-lg transition-colors flex items-center">
                         <BaseLoading v-if="savingBranch" type="spinner" size="small" color="white" :show-text="false" class="mr-1" />
                         {{ savingBranch ? 'Saving...' : 'Update Branch' }}
                       </button>
                     </div>
                   </form>
                 </div>
               </BaseTransition>
             </li>
           </ul>
           <button
             @click="addBranch"
             class="mt-4 w-full inline-flex justify-center items-center px-4 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200 shadow-sm"
           >
             <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
               <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
             </svg>
             Add New Branch
           </button>
        </div>
      </div>
    </BaseTransition>

    <!-- Staff Management -->
    <BaseTransition name="slide-up">
      <div v-if="isEditMode && !loading && !error" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
        <div class="p-6 sm:p-8">
          <div class="flex justify-between items-center border-b border-neutral-200 pb-2 mb-4">
            <h3 class="text-lg font-semibold text-neutral-900">Staff</h3>
            <button @click="addUserToClinic" class="inline-flex justify-center items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700">Add Staff</button>
          </div>
          <ul class="space-y-2">
            <li v-for="user in clinic.staff" :key="user.id" class="p-4 border border-gray-200 rounded-md flex justify-between items-center">
              <div>
                <p class="font-semibold">{{ user.username }} <span class="text-xs bg-gray-200 text-gray-800 px-2 py-1 rounded-full">{{ user.role }}</span></p>
                <p class="text-sm text-neutral-600">{{ user.email }}</p>
              </div>
              <div class="space-x-2">
                <button @click="editUser(user.id)" class="text-sm text-blue-600">Edit</button>
              </div>
            </li>
          </ul>
        </div>
      </div>
     </BaseTransition>

     <!-- Inline Branch Form -->
     <BaseTransition name="slide-down">
       <div v-if="showBranchForm" class="bg-neutral-50 rounded-xl p-4 mb-4 border border-neutral-200">
         <form @submit.prevent="saveBranch" class="space-y-4">
           <div class="flex items-center justify-between mb-3">
             <h4 class="text-sm font-semibold text-neutral-900">{{ editingBranch ? 'Edit Branch' : 'Add New Branch' }}</h4>
             <button type="button" @click="cancelBranchForm" class="text-neutral-400 hover:text-neutral-600 transition-colors">
               <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
               </svg>
             </button>
           </div>

           <div class="space-y-3">
             <div>
               <label for="branchName" class="block text-xs font-medium text-neutral-600 mb-1">Branch Name <span class="text-danger-500">*</span></label>
               <input
                 type="text"
                 id="branchName"
                 v-model="newBranch.name"
                 class="block w-full px-3 py-2 text-sm border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                 placeholder="Enter branch name"
                 required
               >
             </div>
             <div>
               <label for="branchAddress" class="block text-xs font-medium text-neutral-600 mb-1">Address</label>
               <input
                 type="text"
                 id="branchAddress"
                 v-model="newBranch.address"
                 class="block w-full px-3 py-2 text-sm border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                 placeholder="Enter branch address"
               >
             </div>
           </div>

           <div class="flex justify-end space-x-2 pt-2">
             <button type="button" @click="cancelBranchForm" class="px-3 py-2 text-xs font-medium text-neutral-600 bg-white border border-neutral-300 rounded-lg hover:bg-neutral-50 transition-colors">
               Cancel
             </button>
             <button type="submit" :disabled="savingBranch" class="px-3 py-2 text-xs font-medium text-white bg-primary-600 hover:bg-primary-700 disabled:opacity-50 rounded-lg transition-colors flex items-center">
               <BaseLoading v-if="savingBranch" type="spinner" size="small" color="white" :show-text="false" class="mr-1" />
               {{ savingBranch ? 'Saving...' : (editingBranch ? 'Update' : 'Add Branch') }}
             </button>
           </div>
         </form>
       </div>
     </BaseTransition>

   </div>
 </template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClinicStore } from '../stores/clinic'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'

export default {
  name: 'ClinicForm',
  components: {
    BaseLoading,
    BaseTransition
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const clinicStore = useClinicStore()

    const clinic = ref({ name: '', address: '', phone: '', email: '', website: '', branches: [], staff: [] })
    const saving = ref(false)
    const loading = ref(true)
    const error = ref('')
    const showBranchForm = ref(false)
    const savingBranch = ref(false)
    const editingBranch = ref(null)
    const newBranch = ref({ name: '', address: '' })

    const clinicId = route.params.id
    const isEditMode = computed(() => !!clinicId)

    const loadClinicData = async () => {
      if (isEditMode.value) {
        loading.value = true
        error.value = ''
        try {
          await clinicStore.fetchClinic(clinicId)
          clinic.value = { ...clinicStore.currentClinic }
        } catch (err) {
          error.value = 'Failed to load clinic data.'
          console.error(err)
        } finally {
          loading.value = false
        }
      } else {
        loading.value = false
      }
    }

    const saveClinic = async () => {
      saving.value = true
      let result
      if (isEditMode.value) {
        result = await clinicStore.updateClinic(clinicId, clinic.value)
      } else {
        result = await clinicStore.createClinic(clinic.value)
      }

      if (result.success) {
        router.push('/clinics')
      } else {
        // Handle error
        console.error(result.error)
      }
      saving.value = false
    }

    const deleteBranch = async (branchId) => {
        if(confirm('Are you sure you want to delete this branch?')){
            const result = await clinicStore.deleteBranch(clinicId, branchId);
            if(result.success){
                loadClinicData(); // Refresh data
            }
        }
    }

    const cancel = () => {
      router.push('/clinics')
    }

    const addUserToClinic = () => {
      router.push(`/staff/new?clinic_id=${clinicId}`)
    }

    const editUser = (userId) => {
      router.push(`/users/${userId}/edit`)
    }

    const addBranch = () => {
      editingBranch.value = null
      newBranch.value = { name: '', address: '' }
      showBranchForm.value = true
    }

    const editBranch = (branch) => {
      editingBranch.value = branch
      newBranch.value = { name: branch.name, address: branch.address }
      showBranchForm.value = false // Hide add form if open
    }

    const cancelBranchForm = () => {
      showBranchForm.value = false
      editingBranch.value = null
      newBranch.value = { name: '', address: '' }
    }

    const saveBranch = async () => {
      savingBranch.value = true
      let result

      if (editingBranch.value) {
        // Update existing branch
        result = await clinicStore.updateBranch(clinicId, editingBranch.value.id, newBranch.value)
      } else {
        // Create new branch
        result = await clinicStore.createBranch(clinicId, newBranch.value)
      }

      if (result.success) {
        closeBranchModal()
        loadClinicData() // Refresh data
      } else {
        // Handle error - could show a toast notification
        console.error(result.error)
      }
      savingBranch.value = false
    }

    onMounted(() => {
      loadClinicData()
    })

    return {
      clinic,
      isEditMode,
      saving,
      loading,
      error,
      showBranchForm,
      savingBranch,
      newBranch,
      loadClinicData,
      saveClinic,
      deleteBranch,
      cancel,
      addUserToClinic,
      editUser,
      addBranch,
      editBranch,
      cancelBranchForm,
      saveBranch,
    }
  }
}
</script>