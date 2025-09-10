<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="p-6">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-bold text-gray-900">Create Peer Review Case</h2>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitForm" class="space-y-6">
          <!-- Basic Information -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Case Number *</label>
              <input
                v-model="form.case_number"
                type="text"
                required
                placeholder="e.g., PR-2024-001"
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Visibility *</label>
              <select
                v-model="form.visibility"
                required
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="invite_only">Invite Only</option>
                <option value="in_clinic">In-Clinic</option>
                <option value="public">Public</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
            <input
              v-model="form.title"
              type="text"
              required
              placeholder="Brief description of the case"
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Detailed description of the case..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            ></textarea>
          </div>

          <!-- Patient Information -->
          <div class="border-t pt-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Patient Information</h3>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Age</label>
                <input
                  v-model.number="form.patient_age"
                  type="number"
                  min="0"
                  max="150"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Gender</label>
                <select
                  v-model="form.patient_gender"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select Gender</option>
                  <option value="male">Male</option>
                  <option value="female">Female</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Blood Type</label>
                <select
                  v-model="form.patient_blood_type"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select Blood Type</option>
                  <option value="A+">A+</option>
                  <option value="A-">A-</option>
                  <option value="B+">B+</option>
                  <option value="B-">B-</option>
                  <option value="AB+">AB+</option>
                  <option value="AB-">AB-</option>
                  <option value="O+">O+</option>
                  <option value="O-">O-</option>
                </select>
              </div>
            </div>

            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Chief Complaint</label>
                <textarea
                  v-model="form.chief_complaint"
                  rows="2"
                  placeholder="Main reason for visit..."
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                ></textarea>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Medical History</label>
                <textarea
                  v-model="form.medical_history"
                  rows="2"
                  placeholder="Relevant medical history..."
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                ></textarea>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Current Medications</label>
                <textarea
                  v-model="form.current_medications"
                  rows="2"
                  placeholder="Current medications..."
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                ></textarea>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Allergies</label>
                <input
                  v-model="form.allergies"
                  type="text"
                  placeholder="Known allergies..."
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          <!-- Dental Chart -->
          <div class="border-t pt-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Dental Chart</h3>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Dental Chart Data (JSON)</label>
              <textarea
                v-model="form.dental_chart_data"
                rows="4"
                placeholder='{"teeth": [...]}'
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              ></textarea>
              <p class="text-xs text-gray-500 mt-1">Paste dental chart data in JSON format</p>
            </div>
          </div>

          <!-- Share With (only for invite_only) -->
          <div v-if="form.visibility === 'invite_only'" class="border-t pt-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Share With</h3>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Invite Doctors</label>
              <input
                v-model="inviteSearch"
                type="text"
                placeholder="Search doctors by name..."
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                @input="searchDoctors"
              />
              <div v-if="searchResults.length > 0" class="mt-2 max-h-40 overflow-y-auto border border-gray-200 rounded-md">
                <div
                  v-for="doctor in searchResults"
                  :key="doctor.id"
                  @click="addToInviteList(doctor)"
                  class="p-2 hover:bg-gray-50 cursor-pointer"
                >
                  {{ doctor.first_name }} {{ doctor.last_name }} ({{ doctor.email }})
                </div>
              </div>
              <div v-if="form.share_with.length > 0" class="mt-2 flex flex-wrap gap-2">
                <span
                  v-for="doctor in form.share_with"
                  :key="doctor.id"
                  class="bg-blue-100 text-blue-800 px-2 py-1 rounded-md text-sm flex items-center gap-1"
                >
                  {{ doctor.first_name }} {{ doctor.last_name }}
                  <button @click="removeFromInviteList(doctor.id)" class="text-blue-600 hover:text-blue-800">×</button>
                </span>
              </div>
            </div>
          </div>

          <!-- Submit Buttons -->
          <div class="flex justify-end gap-3 pt-6 border-t">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ loading ? 'Creating...' : 'Create Case' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  // No props needed
})

const emit = defineEmits(['close', 'created'])

const form = reactive({
  case_number: '',
  title: '',
  description: '',
  visibility: 'invite_only',
  patient_age: null,
  patient_gender: '',
  patient_blood_type: '',
  chief_complaint: '',
  medical_history: '',
  current_medications: '',
  allergies: '',
  dental_chart_data: '',
  share_with: []
})

const loading = ref(false)
const inviteSearch = ref('')
const searchResults = ref([])
let searchTimeout = null

const searchDoctors = async () => {
  if (inviteSearch.value.length < 2) {
    searchResults.value = []
    return
  }

  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    try {
      const response = await fetch(`/api/users?role=doctor&search=${encodeURIComponent(inviteSearch.value)}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        searchResults.value = data.users.filter(user =>
          !form.share_with.some(invited => invited.id === user.id)
        )
      }
    } catch (error) {
      console.error('Failed to search doctors:', error)
    }
  }, 300)
}

const addToInviteList = (doctor) => {
  if (!form.share_with.some(d => d.id === doctor.id)) {
    form.share_with.push(doctor)
  }
  inviteSearch.value = ''
  searchResults.value = []
}

const removeFromInviteList = (doctorId) => {
  form.share_with = form.share_with.filter(d => d.id !== doctorId)
}

const submitForm = async () => {
  loading.value = true

  try {
    const payload = {
      ...form,
      share_with: form.share_with.map(d => d.id)
    }

    const response = await fetch('/api/peer-review/cases', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(payload)
    })

    if (response.ok) {
      emit('created')
    } else {
      const error = await response.json()
      alert(error.error || 'Failed to create case')
    }
  } catch (error) {
    console.error('Failed to create case:', error)
    alert('Failed to create case')
  } finally {
    loading.value = false
  }
}
</script>