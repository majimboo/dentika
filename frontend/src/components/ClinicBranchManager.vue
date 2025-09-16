<template>
  <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
    <div class="p-6 sm:p-8">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-xl font-semibold text-neutral-900">Branch Management</h2>
          <p class="text-neutral-600 mt-1">Manage schedules and settings for your clinic branches</p>
        </div>
        <div class="flex items-center space-x-3">
          <div class="h-10 w-10 bg-gradient-to-r from-primary-100 to-secondary-100 rounded-lg flex items-center justify-center">
            <font-awesome-icon icon="fa-solid fa-building" class="h-5 w-5 text-primary-600" />
          </div>
        </div>
      </div>

      <!-- No Branches State -->
      <div v-if="!clinic.branches || clinic.branches.length === 0" class="text-center py-12">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-neutral-100 rounded-full mb-4">
          <font-awesome-icon icon="fa-solid fa-building" class="w-8 h-8 text-neutral-400" />
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">No Branches Found</h3>
        <p class="text-neutral-600 mb-4">Your clinic doesn't have any branches set up yet.</p>
        <p class="text-sm text-neutral-500">Contact your system administrator to add branches.</p>
      </div>

      <!-- Branches List -->
      <div v-else class="space-y-6">
        <div
          v-for="branch in clinic.branches"
          :key="branch.id"
          class="border border-neutral-200 rounded-xl p-6 hover:shadow-md transition-shadow duration-200"
          :class="{ 'ring-2 ring-primary-500 border-primary-200': branch.is_main_branch }"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3 class="text-lg font-semibold text-neutral-900">{{ branch.name }}</h3>
                <span
                  v-if="branch.is_main_branch"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800"
                >
                  Main Branch
                </span>
                <span
                  v-if="!branch.is_active"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-danger-100 text-danger-800"
                >
                  Inactive
                </span>
                <span
                  v-if="branch.is_closed_today"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-warning-100 text-warning-800"
                >
                  Closed Today
                </span>
              </div>
              <div class="text-sm text-neutral-600 space-y-1">
                <p v-if="branch.address">
                  <font-awesome-icon icon="fa-solid fa-map-marker-alt" class="w-3 h-3 mr-2" />
                  {{ branch.address }}
                </p>
                <p v-if="branch.phone">
                  <font-awesome-icon icon="fa-solid fa-phone" class="w-3 h-3 mr-2" />
                  {{ branch.phone }}
                </p>
              </div>
            </div>
            <div class="flex items-center space-x-2 ml-4">
              <button
                @click="toggleBranchEdit(branch.id)"
                class="inline-flex items-center px-3 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
              >
                <font-awesome-icon :icon="editingBranch === branch.id ? 'fa-solid fa-times' : 'fa-solid fa-edit'" class="w-4 h-4 mr-2" />
                {{ editingBranch === branch.id ? 'Cancel' : 'Edit' }}
              </button>
            </div>
          </div>

          <!-- Branch Edit Form -->
          <BaseTransition name="slide-down">
            <div v-if="editingBranch === branch.id" class="mt-6 space-y-6">
              <form @submit.prevent="saveBranch(branch)" class="space-y-6">

                <!-- Basic Information -->
                <div class="space-y-4">
                  <h4 class="text-sm font-semibold text-neutral-700 border-b border-neutral-200 pb-2">Basic Information</h4>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label :for="`branch-name-${branch.id}`" class="block text-sm font-medium text-neutral-700 mb-2">Branch Name <span class="text-danger-500">*</span></label>
                      <input
                        :id="`branch-name-${branch.id}`"
                        type="text"
                        v-model="editedBranches[branch.id].name"
                        required
                        class="block w-full px-3 py-2 border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                      >
                    </div>
                    <div>
                      <label :for="`branch-phone-${branch.id}`" class="block text-sm font-medium text-neutral-700 mb-2">Phone</label>
                      <input
                        :id="`branch-phone-${branch.id}`"
                        type="tel"
                        v-model="editedBranches[branch.id].phone"
                        class="block w-full px-3 py-2 border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                      >
                    </div>
                    <div class="md:col-span-2">
                      <label :for="`branch-address-${branch.id}`" class="block text-sm font-medium text-neutral-700 mb-2">Address</label>
                      <input
                        :id="`branch-address-${branch.id}`"
                        type="text"
                        v-model="editedBranches[branch.id].address"
                        class="block w-full px-3 py-2 border border-neutral-300 rounded-lg text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                      >
                    </div>
                  </div>
                </div>

                <!-- Operating Schedule -->
                <div class="space-y-4">
                  <h4 class="text-sm font-semibold text-neutral-700 border-b border-neutral-200 pb-2">Operating Schedule</h4>
                  <ScheduleEditor v-model="editedBranches[branch.id].schedule" />
                </div>

                <!-- Status Settings -->
                <div class="space-y-4">
                  <h4 class="text-sm font-semibold text-neutral-700 border-b border-neutral-200 pb-2">Status Settings</h4>
                  <div class="flex items-center">
                    <input
                      :id="`branch-closed-today-${branch.id}`"
                      type="checkbox"
                      v-model="editedBranches[branch.id].is_closed_today"
                      class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
                    >
                    <label :for="`branch-closed-today-${branch.id}`" class="ml-2 text-sm font-medium text-neutral-700">
                      Closed today (emergency override)
                    </label>
                  </div>
                </div>

                <!-- Form Actions -->
                <div class="flex justify-end space-x-3 pt-4 border-t border-neutral-200">
                  <button
                    type="button"
                    @click="cancelBranchEdit(branch.id)"
                    :disabled="savingBranches[branch.id]"
                    class="inline-flex items-center px-4 py-2 border border-neutral-300 rounded-lg text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    :disabled="savingBranches[branch.id]"
                    class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-secondary-600 hover:from-primary-700 hover:to-secondary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                  >
                    <BaseLoading v-if="savingBranches[branch.id]" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                    {{ savingBranches[branch.id] ? 'Saving...' : 'Save Changes' }}
                  </button>
                </div>
              </form>
            </div>
          </BaseTransition>

          <!-- Current Schedule Preview -->
          <div v-if="editingBranch !== branch.id" class="mt-4">
            <div class="bg-neutral-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-neutral-700 mb-3">Current Schedule</h4>
              <BranchSchedulePreview :schedule="branch.schedule" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useClinicStore } from '../stores/clinic'
import BaseLoading from './BaseLoading.vue'
import BaseTransition from './BaseTransition.vue'
import ScheduleEditor from './ScheduleEditor.vue'
import BranchSchedulePreview from './BranchSchedulePreview.vue'

export default {
  name: 'ClinicBranchManager',
  components: {
    BaseLoading,
    BaseTransition,
    ScheduleEditor,
    BranchSchedulePreview
  },
  props: {
    clinic: {
      type: Object,
      required: true
    }
  },
  emits: ['branch-updated'],
  setup(props, { emit }) {
    const clinicStore = useClinicStore()

    const editingBranch = ref(null)
    const editedBranches = reactive({})
    const savingBranches = reactive({})

    const toggleBranchEdit = (branchId) => {
      if (editingBranch.value === branchId) {
        cancelBranchEdit(branchId)
      } else {
        startBranchEdit(branchId)
      }
    }

    const startBranchEdit = (branchId) => {
      const branch = props.clinic.branches.find(b => b.id === branchId)
      if (branch) {
        editedBranches[branchId] = {
          name: branch.name || '',
          address: branch.address || '',
          phone: branch.phone || '',
          schedule: branch.schedule || '',
          is_closed_today: branch.is_closed_today || false
        }
        editingBranch.value = branchId
      }
    }

    const cancelBranchEdit = (branchId) => {
      editingBranch.value = null
      delete editedBranches[branchId]
    }

    const saveBranch = async (branch) => {
      savingBranches[branch.id] = true

      try {
        const result = await clinicStore.updateBranch(
          props.clinic.id,
          branch.id,
          editedBranches[branch.id]
        )

        if (result.success) {
          editingBranch.value = null
          delete editedBranches[branch.id]
          emit('branch-updated')
        } else {
          alert('Failed to update branch: ' + (result.error || 'Unknown error'))
        }
      } catch (err) {
        console.error('Error updating branch:', err)
        alert('Failed to update branch')
      } finally {
        savingBranches[branch.id] = false
      }
    }

    return {
      editingBranch,
      editedBranches,
      savingBranches,
      toggleBranchEdit,
      startBranchEdit,
      cancelBranchEdit,
      saveBranch
    }
  }
}
</script>