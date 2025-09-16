<template>
  <div class="space-y-8">
    <!-- Loading State -->
    <BaseTransition name="fade">
      <div v-if="loading" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
        <BaseLoading type="spinner" size="large" color="primary" text="Loading branch data..." />
        <p class="text-neutral-600 mt-4">Please wait while we fetch the branch information.</p>
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
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading branch</h3>
        <p class="text-danger-600 mb-4">{{ error }}</p>
        <button
          @click="loadBranchData"
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
          <form @submit.prevent="saveBranch" class="space-y-8">

            <!-- Branch Details -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Branch Details</h3>

              <div class="space-y-2">
                <label for="name" class="block text-sm font-semibold text-gray-700">Branch Name <span class="text-danger-500 ml-1">*</span></label>
                <input type="text" id="name" v-model="branch.name" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white" required>
              </div>

              <div class="space-y-2">
                <label for="address" class="block text-sm font-semibold text-gray-700">Address</label>
                <input type="text" id="address" v-model="branch.address" class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white">
              </div>

              <div class="space-y-4">
                <ScheduleEditor v-model="branch.schedule" />
              </div>

              <div class="flex items-center">
                <input
                  type="checkbox"
                  id="is_closed_today"
                  v-model="branch.is_closed_today"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
                >
                <label for="is_closed_today" class="ml-2 text-sm font-medium text-neutral-600">Closed today (emergency override)</label>
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex justify-between items-center pt-6 border-t border-neutral-200">
              <button
                type="button"
                @click="cancel"
                :disabled="saving"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="saving"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <BaseLoading v-if="saving" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                {{ saving ? 'Saving...' : 'Update Branch' }}
              </button>
            </div>
          </form>
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
import ScheduleEditor from '../components/ScheduleEditor.vue'

export default {
  name: 'BranchForm',
  components: {
    BaseLoading,
    BaseTransition,
    ScheduleEditor
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const clinicStore = useClinicStore()

    const branch = ref({
      name: '',
      address: '',
      schedule: '',
      is_closed_today: false
    })
    const saving = ref(false)
    const loading = ref(true)
    const error = ref('')

    const clinicId = route.params.clinicId
    const branchId = route.params.branchId
    const isEditMode = computed(() => !!branchId)

    console.log('BranchForm route params:', route.params)
    console.log('clinicId:', clinicId, 'branchId:', branchId)

    const loadBranchData = async () => {
      if (isEditMode.value) {
        loading.value = true
        error.value = ''

        try {
          console.log('Loading branch data for clinicId:', clinicId, 'branchId:', branchId)

          // First check if we already have the clinic data cached
          let clinicData = clinicStore.currentClinic

          if (!clinicData || clinicData.id !== clinicId) {
            console.log('Fetching clinic data from API...')
            // First fetch the clinic to get branch data
            const result = await clinicStore.fetchClinic(clinicId)
            console.log('Clinic fetch result:', result)

            if (!result.success) {
              error.value = result.error || 'Failed to fetch clinic data'
              return
            }
            clinicData = result.data
          } else {
            console.log('Using cached clinic data')
          }

          console.log('Clinic branches:', clinicData.branches)

          if (!clinicData.branches || clinicData.branches.length === 0) {
            error.value = 'No branches found for this clinic'
            return
          }

          // Handle type conversion - route params are strings, but IDs might be numbers
          const branchIdNum = parseInt(branchId)
          const foundBranch = clinicData.branches.find(b => b.id === branchId || b.id === branchIdNum)

          console.log('Looking for branch with ID:', branchId, 'or', branchIdNum)
          console.log('Found branch:', foundBranch)

          if (foundBranch) {
            branch.value = { ...foundBranch }
          } else {
            error.value = `Branch not found. Available branches: ${clinicData.branches.map(b => b.id).join(', ')}`
          }
        } catch (err) {
          error.value = 'Failed to load branch data.'
          console.error('Error loading branch data:', err)
        } finally {
          loading.value = false
        }
      } else {
        loading.value = false
      }
    }

    const saveBranch = async () => {
      saving.value = true

      const result = await clinicStore.updateBranch(clinicId, branchId, branch.value)

      if (result.success) {
        router.push(`/clinics/${clinicId}/edit`)
      } else {
        // Handle error - could show a toast notification
        console.error(result.error)
        alert('Failed to update branch: ' + result.error)
      }
      saving.value = false
    }

    const cancel = () => {
      router.push(`/clinics/${clinicId}/edit`)
    }

    onMounted(() => {
      loadBranchData()
    })

    return {
      branch,
      isEditMode,
      saving,
      loading,
      error,
      loadBranchData,
      saveBranch,
      cancel
    }
  }
}
</script>