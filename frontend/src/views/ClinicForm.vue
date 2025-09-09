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
          <ul class="space-y-2">
            <li v-for="branch in clinic.branches" :key="branch.id" class="p-4 border rounded-md flex justify-between items-center">
              <div>
                <p class="font-semibold">{{ branch.name }} <span v-if="branch.is_main_branch" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">Main</span></p>
                <p class="text-sm text-neutral-600">{{ branch.address }}</p>
              </div>
              <div class="space-x-2">
                <button class="text-sm text-blue-600">Edit</button>
                <button @click="deleteBranch(branch.id)" class="text-sm text-red-600">Delete</button>
              </div>
            </li>
          </ul>
          <button class="mt-4 inline-flex justify-center items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700">Add Branch</button>
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
            <li v-for="user in clinic.staff" :key="user.id" class="p-4 border rounded-md flex justify-between items-center">
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

    onMounted(() => {
      loadClinicData()
    })

    return {
      clinic,
      isEditMode,
      saving,
      loading,
      error,
      loadClinicData,
      saveClinic,
      deleteBranch,
      cancel,
      addUserToClinic,
      editUser,
    }
  }
}
</script>