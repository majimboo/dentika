<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div class="p-6 sm:p-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-neutral-900">Clinic Settings</h1>
            <p class="text-neutral-600 mt-1">Manage your clinic configuration and branches</p>
          </div>
          <div class="flex items-center space-x-3">
            <div class="h-12 w-12 bg-gradient-to-r from-primary-600 to-secondary-600 rounded-xl flex items-center justify-center">
              <font-awesome-icon icon="fa-solid fa-cog" class="h-6 w-6 text-white" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <BaseTransition name="fade">
      <div v-if="loading" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <BaseLoading type="spinner" size="large" color="primary" text="Loading clinic settings..." />
        <p class="text-neutral-600 mt-4">Please wait while we fetch your clinic information.</p>
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
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading settings</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button
          @click="loadClinicData"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
        >
          Try Again
        </button>
      </div>
    </BaseTransition>

    <!-- Settings Content -->
    <BaseTransition name="slide-up">
      <div v-if="!loading && !error && clinic" class="space-y-8">

        <!-- Clinic Information -->
        <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
          <div class="p-6 sm:p-8">
            <div class="flex items-center justify-between mb-6">
              <h2 class="text-xl font-semibold text-neutral-900">Clinic Information</h2>
              <button
                @click="editMode.clinic = !editMode.clinic"
                class="inline-flex items-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
              >
                <font-awesome-icon :icon="editMode.clinic ? 'fa-solid fa-times' : 'fa-solid fa-edit'" class="w-4 h-4 mr-2" />
                {{ editMode.clinic ? 'Cancel' : 'Edit' }}
              </button>
            </div>

            <div v-if="!editMode.clinic" class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Clinic Name</label>
                <p class="text-neutral-900 font-medium">{{ clinic.name }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Address</label>
                <p class="text-neutral-900">{{ clinic.address || 'Not provided' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Phone</label>
                <p class="text-neutral-900">{{ clinic.phone || 'Not provided' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Email</label>
                <p class="text-neutral-900">{{ clinic.email || 'Not provided' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Website</label>
                <p class="text-neutral-900">{{ clinic.website || 'Not provided' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Tagline</label>
                <p class="text-neutral-900">{{ clinic.tagline || 'Not provided' }}</p>
              </div>
            </div>

            <form v-else @submit.prevent="saveClinicInfo" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label for="clinic-name" class="block text-sm font-medium text-neutral-700 mb-2">Clinic Name <span class="text-danger-500">*</span></label>
                  <input
                    id="clinic-name"
                    type="text"
                    v-model="editedClinic.name"
                    required
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
                <div>
                  <label for="clinic-address" class="block text-sm font-medium text-neutral-700 mb-2">Address</label>
                  <input
                    id="clinic-address"
                    type="text"
                    v-model="editedClinic.address"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
                <div>
                  <label for="clinic-phone" class="block text-sm font-medium text-neutral-700 mb-2">Phone</label>
                  <input
                    id="clinic-phone"
                    type="tel"
                    v-model="editedClinic.phone"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
                <div>
                  <label for="clinic-email" class="block text-sm font-medium text-neutral-700 mb-2">Email</label>
                  <input
                    id="clinic-email"
                    type="email"
                    v-model="editedClinic.email"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
                <div>
                  <label for="clinic-website" class="block text-sm font-medium text-neutral-700 mb-2">Website</label>
                  <input
                    id="clinic-website"
                    type="url"
                    v-model="editedClinic.website"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
                <div>
                  <label for="clinic-tagline" class="block text-sm font-medium text-neutral-700 mb-2">Tagline</label>
                  <input
                    id="clinic-tagline"
                    type="text"
                    v-model="editedClinic.tagline"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  >
                </div>
              </div>

              <div class="flex justify-end space-x-3">
                <button
                  type="button"
                  @click="cancelClinicEdit"
                  :disabled="savingClinic"
                  class="inline-flex items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="savingClinic"
                  class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                >
                  <BaseLoading v-if="savingClinic" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                  {{ savingClinic ? 'Saving...' : 'Save Changes' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Branch Management -->
        <ClinicBranchManager
          :clinic="clinic"
          @branch-updated="handleBranchUpdated"
        />

      </div>
    </BaseTransition>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useClinicStore } from '../stores/clinic'
import { useAuthStore } from '../stores/auth'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'
import ClinicBranchManager from '../components/ClinicBranchManager.vue'

export default {
  name: 'ClinicSettings',
  components: {
    BaseLoading,
    BaseTransition,
    ClinicBranchManager
  },
  setup() {
    const clinicStore = useClinicStore()
    const authStore = useAuthStore()

    const clinic = ref(null)
    const loading = ref(true)
    const error = ref('')
    const savingClinic = ref(false)

    const editMode = ref({
      clinic: false
    })

    const editedClinic = ref({
      name: '',
      address: '',
      phone: '',
      email: '',
      website: '',
      tagline: ''
    })

    const loadClinicData = async () => {
      loading.value = true
      error.value = ''

      try {
        const result = await clinicStore.fetchClinic(authStore.user.clinic_id)

        if (result.success) {
          clinic.value = result.data
        } else {
          error.value = result.error || 'Failed to load clinic data'
        }
      } catch (err) {
        error.value = 'Failed to load clinic data'
        console.error('Error loading clinic data:', err)
      } finally {
        loading.value = false
      }
    }

    const startClinicEdit = () => {
      editedClinic.value = {
        name: clinic.value.name || '',
        address: clinic.value.address || '',
        phone: clinic.value.phone || '',
        email: clinic.value.email || '',
        website: clinic.value.website || '',
        tagline: clinic.value.tagline || ''
      }
      editMode.value.clinic = true
    }

    const cancelClinicEdit = () => {
      editMode.value.clinic = false
      editedClinic.value = {
        name: '',
        address: '',
        phone: '',
        email: '',
        website: '',
        tagline: ''
      }
    }

    const saveClinicInfo = async () => {
      savingClinic.value = true

      try {
        const result = await clinicStore.updateClinic(clinic.value.id, editedClinic.value)

        if (result.success) {
          clinic.value = { ...clinic.value, ...editedClinic.value }
          editMode.value.clinic = false
        } else {
          alert('Failed to update clinic: ' + (result.error || 'Unknown error'))
        }
      } catch (err) {
        console.error('Error updating clinic:', err)
        alert('Failed to update clinic')
      } finally {
        savingClinic.value = false
      }
    }

    const handleBranchUpdated = () => {
      loadClinicData()
    }

    onMounted(() => {
      loadClinicData()
    })

    return {
      clinic,
      loading,
      error,
      editMode,
      editedClinic,
      savingClinic,
      loadClinicData,
      startClinicEdit,
      cancelClinicEdit,
      saveClinicInfo,
      handleBranchUpdated
    }
  }
}
</script>