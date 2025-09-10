<template>
  <div class="relative">
    <button
      @click="toggleDropdown"
      class="flex items-center space-x-2 px-3 py-2 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors duration-200 min-w-[200px] justify-between"
      :class="{ 'ring-2 ring-primary-500 ring-opacity-50': showDropdown }"
    >
      <div class="flex items-center space-x-2 flex-1 min-w-0">
        <font-awesome-icon icon="fa-solid fa-building" class="w-4 h-4 text-neutral-500 flex-shrink-0" />
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-neutral-900 truncate">
            {{ currentBranchName }}
          </p>
          <p class="text-xs text-neutral-500 truncate">
            {{ currentClinicName }}
          </p>
        </div>
      </div>
      <font-awesome-icon
        :icon="showDropdown ? 'fa-solid fa-chevron-up' : 'fa-solid fa-chevron-down'"
        class="w-4 h-4 text-neutral-500 flex-shrink-0"
      />
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="showDropdown"
      class="absolute top-full mt-2 w-full bg-white rounded-lg shadow-lg border border-neutral-200 py-2 z-50 max-h-64 overflow-y-auto"
    >
      <!-- Current Clinic Header -->
      <div class="px-4 py-2 border-b border-neutral-100">
        <p class="text-xs font-semibold text-neutral-500 uppercase tracking-wide">
          Current Clinic
        </p>
        <p class="text-sm font-medium text-neutral-900 mt-1">
          {{ currentClinicName }}
        </p>
      </div>

      <!-- Branches List -->
      <div class="py-1">
        <button
          v-for="branch in activeBranches"
          :key="branch.id"
          @click="selectBranch(branch)"
          class="w-full px-4 py-2 text-left hover:bg-neutral-50 transition-colors duration-200 flex items-center space-x-3"
          :class="{ 'bg-primary-50 text-primary-700': branch.id === currentBranchId }"
        >
          <font-awesome-icon
            :icon="branch.is_main_branch ? 'fa-solid fa-star' : 'fa-solid fa-building'"
            class="w-4 h-4 flex-shrink-0"
            :class="branch.is_main_branch ? 'text-yellow-500' : 'text-neutral-400'"
          />
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-neutral-900 truncate">
              {{ branch.name }}
            </p>
            <p class="text-xs text-neutral-500 truncate">
              {{ branch.address }}
            </p>
          </div>
          <font-awesome-icon
            v-if="branch.id === currentBranchId"
            icon="fa-solid fa-check"
            class="w-4 h-4 text-primary-600 flex-shrink-0"
          />
        </button>
      </div>

      <!-- No branches message -->
      <div v-if="activeBranches.length === 0" class="px-4 py-3 text-center">
        <p class="text-sm text-neutral-500">No active branches found</p>
      </div>

      <!-- Clinic Overview Option -->
      <div class="border-t border-neutral-100 mt-1 pt-1">
        <button
          @click="viewAllBranches"
          class="w-full px-4 py-2 text-left hover:bg-neutral-50 transition-colors duration-200 flex items-center space-x-3"
        >
          <font-awesome-icon icon="fa-solid fa-eye" class="w-4 h-4 text-neutral-400" />
          <span class="text-sm text-neutral-700">View All Branches</span>
        </button>
      </div>
    </div>

    <!-- Overlay to close dropdown when clicking outside -->
    <div
      v-if="showDropdown"
      @click="closeDropdown"
      class="fixed inset-0 z-40"
    ></div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClinicStore } from '../stores/clinic'

export default {
  name: 'BranchSelector',
  setup() {
    const router = useRouter()
    const clinicStore = useClinicStore()
    const showDropdown = ref(false)

    const currentBranch = computed(() => clinicStore.getCurrentBranch)
    const currentBranchId = computed(() => clinicStore.getCurrentBranchId)
    const currentBranchName = computed(() => clinicStore.getCurrentBranchName)
    const currentClinicName = computed(() => clinicStore.currentClinic?.name || 'No Clinic')
    const activeBranches = computed(() => clinicStore.getActiveBranches)

    const toggleDropdown = () => {
      showDropdown.value = !showDropdown.value
    }

    const closeDropdown = () => {
      showDropdown.value = false
    }

    const selectBranch = async (branch) => {
      const result = await clinicStore.switchBranch(branch.id)
      if (result.success) {
        showDropdown.value = false
      } else {
        console.error('Failed to switch branch:', result.error)
        // Could show a toast notification here
      }
    }

    const viewAllBranches = () => {
      // Navigate to branch overview page
      router.push('/branch-overview')
      showDropdown.value = false
    }

    // Close dropdown when clicking escape
    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        showDropdown.value = false
      }
    }

    onMounted(() => {
      document.addEventListener('keydown', handleEscape)
    })

    onUnmounted(() => {
      document.removeEventListener('keydown', handleEscape)
    })

    return {
      showDropdown,
      currentBranch,
      currentBranchId,
      currentBranchName,
      currentClinicName,
      activeBranches,
      toggleDropdown,
      closeDropdown,
      selectBranch,
      viewAllBranches,
      router
    }
  }
}
</script>