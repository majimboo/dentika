<template>
  <div class="peer-review-create">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Create Peer Review Case</h1>
            <p class="text-gray-600 mt-1">Share a case anonymously for peer collaboration</p>
          </div>
          <router-link
            to="/peer-review"
            class="text-gray-400 hover:text-gray-600 p-2"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </router-link>
        </div>
      </div>

      <!-- Form -->
      <form @submit.prevent="submitForm" class="space-y-6">
        <!-- Basic Information -->
        <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Case Information</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Case Number *</label>
              <input
                v-model="form.case_number"
                type="text"
                required
                placeholder="e.g., PR-2024-001"
                class="w-full border border-neutral-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Visibility *</label>
              <select
                v-model="form.visibility"
                required
                class="w-full border border-neutral-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="invite_only">Invite Only</option>
                <option value="in_clinic">In-Clinic</option>
                <option value="public">Public</option>
              </select>
            </div>
          </div>

          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
            <input
              v-model="form.title"
              type="text"
              required
              placeholder="Brief description of the case"
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Detailed description of the case..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            ></textarea>
          </div>
        </div>

        <!-- Patient Selection -->
        <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Select Patient</h2>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Patient *</label>
            <div class="relative">
              <input
                v-model="patientSearch"
                @input="searchPatients"
                @focus="showPatientDropdown = true"
                type="text"
                placeholder="Search for patient by name..."
                class="w-full border border-neutral-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <div
                v-if="showPatientDropdown && (filteredPatients.length > 0 || patientSearch)"
                class="absolute z-10 w-full bg-white border border-neutral-300 rounded-md shadow-lg max-h-48 overflow-y-auto mt-1"
              >
                <div
                  v-for="patient in filteredPatients"
                  :key="patient.id"
                  @click="selectPatient(patient)"
                  class="px-3 py-2 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                >
                  <div class="font-medium">{{ patient.first_name }} {{ patient.last_name }}</div>
                  <div class="text-sm text-gray-500">{{ patient.email || patient.phone || 'No contact info' }}</div>
                </div>
                <div v-if="filteredPatients.length === 0 && patientSearch" class="px-3 py-2 text-gray-500 text-sm">
                  No patients found
                </div>
              </div>
            </div>
            <div v-if="selectedPatient" class="mt-2 p-3 bg-blue-50 border border-blue-200 rounded-md">
              <div class="flex justify-between items-center">
                <div>
                  <div class="font-medium text-blue-900">{{ selectedPatient.first_name }} {{ selectedPatient.last_name }}</div>
                  <div class="text-sm text-blue-700">Age: {{ selectedPatient.age || 'Unknown' }} | Gender: {{ selectedPatient.gender || 'Unknown' }}</div>
                </div>
                <button @click="clearPatientSelection" class="text-blue-600 hover:text-blue-800">×</button>
              </div>
            </div>
          </div>
        </div>

        <!-- Appointment Selection -->
        <div v-if="selectedPatient" class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Select Appointment</h2>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Appointment *</label>
            <div class="relative">
              <select
                v-model="selectedAppointment"
                @change="selectAppointmentData"
                class="w-full border border-neutral-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Select an appointment...</option>
                <option
                  v-for="appointment in patientAppointments"
                  :key="appointment.id"
                  :value="appointment"
                >
                  {{ formatAppointmentDate(appointment.start_time) }} - {{ appointment.description || appointment.title || 'No description' }}
                </option>
              </select>
            </div>
            <div v-if="loadingAppointments" class="mt-2 text-sm text-gray-500">
              Loading appointments...
            </div>
            <div v-if="selectedAppointment" class="mt-3 p-3 bg-green-50 border border-green-200 rounded-md">
              <h4 class="font-medium text-green-900 mb-2">Selected Appointment Details:</h4>
              <div class="text-sm text-green-800">
                <div><strong>Date:</strong> {{ formatAppointmentDate(selectedAppointment.start_time) }}</div>
                <div v-if="selectedAppointment.title"><strong>Title:</strong> {{ selectedAppointment.title }}</div>
                <div v-if="selectedAppointment.description"><strong>Description:</strong> {{ selectedAppointment.description }}</div>
                <div v-if="selectedAppointment.status"><strong>Status:</strong> {{ selectedAppointment.status }}</div>
                <div v-if="selectedAppointment.pre_appointment_notes"><strong>Pre-Notes:</strong> {{ selectedAppointment.pre_appointment_notes }}</div>
                <div v-if="selectedAppointment.post_appointment_notes"><strong>Post-Notes:</strong> {{ selectedAppointment.post_appointment_notes }}</div>
              </div>
            </div>
          </div>
        </div>


        <!-- Share With (only for invite_only) -->
        <div v-if="form.visibility === 'invite_only'" class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Share With</h2>
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

        <!-- Submit Actions -->
        <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
          <div class="flex justify-end gap-3">
            <router-link
              to="/peer-review"
              class="px-4 py-2 text-gray-700 border border-neutral-300 rounded-md hover:bg-gray-50"
            >
              Cancel
            </router-link>
            <button
              type="submit"
              :disabled="loading"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ loading ? 'Creating...' : 'Create Case' }}
            </button>
          </div>
        </div>
      </form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()

const form = reactive({
  case_number: '',
  title: '',
  description: '',
  visibility: 'invite_only',
  patient_id: null,
  appointment_id: null,
  share_with: []
})

const loading = ref(false)
const inviteSearch = ref('')
const searchResults = ref([])
let searchTimeout = null

// Patient selection
const selectedPatient = ref(null)
const patientSearch = ref('')
const filteredPatients = ref([])
const showPatientDropdown = ref(false)
let patientSearchTimeout = null

// Appointment selection
const selectedAppointment = ref(null)
const patientAppointments = ref([])
const loadingAppointments = ref(false)

const searchDoctors = async () => {
  if (inviteSearch.value.length < 2) {
    searchResults.value = []
    return
  }

  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    try {
      const result = await apiService.get(`/api/users?role=doctor&search=${encodeURIComponent(inviteSearch.value)}`)

      if (result.success) {
        searchResults.value = result.data.users.filter(user =>
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

// Patient search and selection
const searchPatients = async () => {
  if (patientSearch.value.length < 2) {
    filteredPatients.value = []
    return
  }

  clearTimeout(patientSearchTimeout)
  patientSearchTimeout = setTimeout(async () => {
    try {
      const result = await apiService.get(`/api/patients?search=${encodeURIComponent(patientSearch.value)}&limit=10`)

      if (result.success) {
        filteredPatients.value = result.data.patients || []
      }
    } catch (error) {
      console.error('Failed to search patients:', error)
    }
  }, 300)
}

const selectPatient = async (patient) => {
  selectedPatient.value = patient
  form.patient_id = patient.id
  patientSearch.value = `${patient.first_name} ${patient.last_name}`
  showPatientDropdown.value = false
  filteredPatients.value = []

  // Load patient appointments
  await loadPatientAppointments(patient.id)
}

const clearPatientSelection = () => {
  selectedPatient.value = null
  form.patient_id = null
  patientSearch.value = ''
  selectedAppointment.value = null
  form.appointment_id = null
  patientAppointments.value = []
}

const loadPatientAppointments = async (patientId) => {
  loadingAppointments.value = true
  try {
    const result = await apiService.get(`/api/appointments?patient_id=${patientId}&status=completed&limit=50`)

    if (result.success) {
      patientAppointments.value = result.data.appointments || []
    }
  } catch (error) {
    console.error('Failed to load appointments:', error)
  } finally {
    loadingAppointments.value = false
  }
}

const selectAppointmentData = () => {
  if (selectedAppointment.value) {
    form.appointment_id = selectedAppointment.value.id
  } else {
    form.appointment_id = null
  }
}

const formatAppointmentDate = (dateTime) => {
  if (!dateTime) return 'Unknown Date'
  const date = new Date(dateTime)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// Close dropdown when clicking outside
document.addEventListener('click', (event) => {
  if (!event.target.closest('.relative')) {
    showPatientDropdown.value = false
  }
})

const submitForm = async () => {
  loading.value = true

  try {
    const payload = {
      ...form,
      share_with: form.share_with.map(d => d.id)
    }

    const result = await apiService.post('/api/peer-review/cases', payload)

    if (result.success) {
      router.push('/peer-review')
    } else {
      alert(result.error || 'Failed to create case')
    }
  } catch (error) {
    console.error('Failed to create case:', error)
    alert('Failed to create case')
  } finally {
    loading.value = false
  }
}
</script>