<template>
  <div class="dental-chart-history">
    <!-- Header -->
    <div class="history-header flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
      <div>
        <h3 class="text-lg font-semibold text-gray-900">Dental Chart History</h3>
        <p class="text-gray-600 text-sm">View dental chart states from different visits</p>
      </div>

      <div class="header-actions flex items-center space-x-3">
        <!-- Chart type filter -->
        <div class="chart-type-filter">
          <label for="chart-type-filter" class="block text-sm font-medium text-gray-700 mb-1">
            Chart Type:
          </label>
          <select
            id="chart-type-filter"
            v-model="selectedChartType"
            @change="loadSnapshots"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Types</option>
            <option value="permanent">Adult (Permanent)</option>
            <option value="primary">Child (Primary)</option>
          </select>
        </div>

        <!-- Create snapshot button -->
        <div v-if="isDoctor">
          <button
            @click="showCreateSnapshotModal = true"
            class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors text-sm"
          >
            Create Snapshot
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <!-- Snapshots Timeline -->
    <div v-else-if="snapshots.length > 0" class="snapshots-timeline">
      <!-- Timeline -->
      <div class="timeline-container">
        <div class="timeline-header mb-4">
          <h4 class="text-md font-medium text-gray-800">Visit Timeline</h4>
          <p class="text-sm text-gray-600">Click on any visit to view the dental chart at that time</p>
        </div>

        <div class="timeline relative">
          <!-- Timeline line -->
          <div class="timeline-line absolute left-4 top-0 bottom-0 w-0.5 bg-gray-300"></div>

          <!-- Timeline items -->
          <div
            v-for="(snapshot, index) in snapshots"
            :key="snapshot.id"
            class="timeline-item relative pl-10 pb-8 cursor-pointer"
            :class="{ 'active': selectedSnapshot?.id === snapshot.id }"
            @click="selectSnapshot(snapshot)"
          >
            <!-- Timeline dot -->
            <div
              class="timeline-dot absolute left-2 w-4 h-4 rounded-full border-2 bg-white transition-all"
              :class="
                selectedSnapshot?.id === snapshot.id
                  ? 'border-blue-500 bg-blue-500'
                  : 'border-gray-300 hover:border-blue-400'
              "
            ></div>

            <!-- Timeline content -->
            <div
              class="timeline-content bg-white rounded-lg border border-gray-200 p-4 shadow-sm hover:shadow-md transition-shadow"
              :class="
                selectedSnapshot?.id === snapshot.id
                  ? 'border-blue-300 bg-blue-50'
                  : 'hover:border-gray-300'
              "
            >
              <div class="timeline-date text-sm font-medium text-gray-900">
                {{ formatSnapshotDate(snapshot.created_at) }}
              </div>

              <div v-if="snapshot.appointment" class="timeline-appointment text-sm text-gray-600 mt-1">
                <font-awesome-icon icon="fa-solid fa-calendar" class="w-3 h-3 mr-1" />
                {{ snapshot.appointment.title || 'Appointment' }}
              </div>

              <div class="timeline-chart-type text-xs text-gray-500 mt-1">
                {{ formatChartType(snapshot.chart_type) }} Chart
              </div>

              <div v-if="snapshot.visit_notes" class="timeline-notes text-sm text-gray-700 mt-2">
                <font-awesome-icon icon="fa-solid fa-sticky-note" class="w-3 h-3 mr-1" />
                {{ snapshot.visit_notes.length > 50 ? snapshot.visit_notes.substring(0, 50) + '...' : snapshot.visit_notes }}
              </div>

              <div class="timeline-doctor text-xs text-gray-500 mt-2">
                By: {{ snapshot.created_by?.first_name }} {{ snapshot.created_by?.last_name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- No snapshots state -->
    <div v-else class="no-snapshots-state bg-gray-50 border border-gray-200 rounded-lg p-12 text-center">
      <div class="text-gray-400 mb-4">
        <font-awesome-icon icon="fa-solid fa-history" class="w-16 h-16" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No Chart History</h3>
      <p class="text-gray-600 mb-4">No dental chart snapshots have been created for this patient yet.</p>
      <button
        v-if="isDoctor"
        @click="showCreateSnapshotModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
      >
        Create First Snapshot
      </button>
    </div>

    <!-- Selected Snapshot Viewer -->
    <div v-if="selectedSnapshot" class="snapshot-viewer mt-8">
      <div class="viewer-header flex justify-between items-center mb-4">
        <div>
          <h4 class="text-lg font-semibold text-gray-900">
            Chart State - {{ formatSnapshotDate(selectedSnapshot.created_at) }}
          </h4>
          <p class="text-sm text-gray-600">
            {{ formatChartType(selectedSnapshot.chart_type) }} dental chart as of this visit
          </p>
        </div>
        <button
          @click="selectedSnapshot = null"
          class="text-gray-400 hover:text-gray-600"
        >
          <font-awesome-icon icon="fa-solid fa-times" class="w-5 h-5" />
        </button>
      </div>

      <!-- Snapshot metadata -->
      <div class="snapshot-metadata bg-gray-50 rounded-lg p-4 mb-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div>
            <span class="font-medium text-gray-700">Visit Date:</span>
            <div>{{ formatSnapshotDate(selectedSnapshot.created_at) }}</div>
          </div>
          <div v-if="selectedSnapshot.appointment">
            <span class="font-medium text-gray-700">Appointment:</span>
            <div>{{ selectedSnapshot.appointment.title || 'General Appointment' }}</div>
          </div>
          <div>
            <span class="font-medium text-gray-700">Doctor:</span>
            <div>{{ selectedSnapshot.created_by?.first_name }} {{ selectedSnapshot.created_by?.last_name }}</div>
          </div>
        </div>
        <div v-if="selectedSnapshot.visit_notes" class="mt-4">
          <span class="font-medium text-gray-700">Visit Notes:</span>
          <div class="mt-1 text-gray-600">{{ selectedSnapshot.visit_notes }}</div>
        </div>
      </div>

      <!-- Snapshot chart display -->
      <div class="snapshot-chart bg-white rounded-lg border border-gray-200 p-6">
        <DentalChartSnapshot
          :snapshotData="selectedSnapshot.teeth_data"
          :chartType="selectedSnapshot.chart_type"
          readonly
        />
      </div>
    </div>

    <!-- Create Snapshot Modal -->
    <CreateSnapshotModal
      v-if="showCreateSnapshotModal"
      :patientId="patientId"
      @created="handleSnapshotCreated"
      @cancel="showCreateSnapshotModal = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useDentalRecordStore } from '../stores/dentalRecord'
import { useNotificationStore } from '../stores/notification'
import DentalChartSnapshot from './DentalChartSnapshot.vue'
import CreateSnapshotModal from './CreateSnapshotModal.vue'

const props = defineProps({
  patientId: {
    type: Number,
    required: true
  }
})

const authStore = useAuthStore()
const dentalRecordStore = useDentalRecordStore()
const notificationStore = useNotificationStore()

const snapshots = ref([])
const selectedSnapshot = ref(null)
const selectedChartType = ref('')
const loading = ref(false)
const showCreateSnapshotModal = ref(false)

const isDoctor = computed(() => authStore.isDoctor || authStore.isSuperAdmin)

const formatSnapshotDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatChartType = (chartType) => {
  return chartType === 'permanent' ? 'Adult' : 'Child'
}

const loadSnapshots = async () => {
  loading.value = true
  try {
    const result = await dentalRecordStore.fetchPatientDentalSnapshots(
      props.patientId,
      selectedChartType.value || null
    )
    if (result.success) {
      snapshots.value = result.data
    }
  } catch (error) {
    console.error('Error loading snapshots:', error)
    notificationStore.showError('Failed to load dental chart history')
  } finally {
    loading.value = false
  }
}

const selectSnapshot = (snapshot) => {
  selectedSnapshot.value = snapshot
}

const handleSnapshotCreated = (newSnapshot) => {
  showCreateSnapshotModal.value = false
  notificationStore.showSuccess('Dental chart snapshot created successfully')
  loadSnapshots() // Reload snapshots to include the new one
}

onMounted(() => {
  loadSnapshots()
})
</script>

<style scoped>
@import "../styles/main.css";

.dental-chart-history {
  @apply bg-white rounded-lg shadow-sm border border-gray-200 p-6;
}

.timeline-item.active .timeline-content {
  @apply border-blue-300 bg-blue-50;
}

.timeline-item.active .timeline-dot {
  @apply border-blue-500 bg-blue-500;
}

.timeline-item:hover .timeline-dot {
  @apply border-blue-400;
}

.timeline-content {
  @apply transition-all duration-200;
}

.timeline-content:hover {
  @apply shadow-md;
}

.snapshot-viewer {
  @apply border-t border-gray-200 pt-6;
}

@media (max-width: 768px) {
  .history-header {
    @apply flex-col items-start;
  }

  .header-actions {
    @apply w-full flex-col items-stretch space-x-0 space-y-3;
  }

  .chart-type-filter {
    @apply w-full;
  }

  .timeline-item {
    @apply pl-8;
  }

  .timeline-dot {
    @apply left-1;
  }

  .timeline-line {
    @apply left-2;
  }
}
</style>