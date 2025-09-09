<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="flex items-center mb-6">
      <button
        @click="$router.go(-1)"
        class="mr-4 p-2 rounded-xl hover:bg-neutral-100 transition-all duration-200"
      >
        <svg class="w-5 h-5 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
        </svg>
      </button>

      <div>
        <h1 class="text-2xl font-bold text-neutral-900">
          {{ isEditing ? 'Edit Appointment' : 'New Appointment' }}
        </h1>
        <p class="text-neutral-600">
          {{ isEditing ? 'Update appointment details' : 'Schedule a new patient appointment' }}
        </p>
      </div>
    </div>

    <!-- Form Container -->
    <div class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div class="p-6 sm:p-8">
        <form @submit.prevent="handleSubmit" class="space-y-8">
          <!-- Patient Selection -->
          <div class="space-y-6">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Patient Information</h3>

            <div class="space-y-2">
              <label class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                Patient
                <span class="text-danger-500 ml-1">*</span>
              </label>
              <div class="relative">
                <input
                  v-model="patientSearchQuery"
                  @input="searchPatients"
                  @focus="showPatientDropdown = true"
                  type="text"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  :class="{ 'border-red-500': errors.patient_id }"
                  placeholder="Search for patient..."
                  autocomplete="off"
                  required
                />

                <!-- Patient Dropdown -->
                <div
                  v-if="showPatientDropdown && (filteredPatients.length > 0 || patientSearchQuery)"
                  class="absolute z-10 w-full bg-white border border-neutral-300 rounded-xl shadow-lg max-h-48 overflow-y-auto mt-1"
                >
                  <div
                    v-for="patient in filteredPatients"
                    :key="patient.id"
                    @click="selectPatient(patient)"
                    class="patient-option px-4 py-3 hover:bg-neutral-50 cursor-pointer border-b last:border-b-0 transition-all duration-200"
                  >
                    <div class="font-medium text-neutral-900">
                      {{ patient.first_name }} {{ patient.last_name }}
                    </div>
                    <div class="text-sm text-neutral-600">
                      {{ patient.phone }} • {{ patient.email }}
                    </div>
                  </div>

                  <div
                    v-if="filteredPatients.length === 0 && patientSearchQuery"
                    class="px-4 py-3 text-neutral-500 text-sm"
                  >
                    No patients found.
                    <button type="button" @click="createNewPatient" class="text-primary-600 hover:text-primary-700 ml-1">
                      Create new patient
                    </button>
                  </div>
                </div>
              </div>
              <span v-if="errors.patient_id" class="text-red-500 text-sm mt-1">{{ errors.patient_id }}</span>
            </div>
          </div>

          <!-- Schedule -->
          <div class="space-y-6">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Schedule</h3>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                  </svg>
                  Date
                  <span class="text-danger-500 ml-1">*</span>
                </label>
                <input
                  v-model="formData.date"
                  type="date"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  :class="{ 'border-red-500': errors.date }"
                  required
                  :min="minDate"
                />
                <span v-if="errors.date" class="text-red-500 text-sm mt-1">{{ errors.date }}</span>
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  Time
                  <span class="text-danger-500 ml-1">*</span>
                </label>
                <input
                  v-model="formData.time"
                  type="time"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  :class="{ 'border-red-500': errors.time }"
                  required
                />
                <span v-if="errors.time" class="text-red-500 text-sm mt-1">{{ errors.time }}</span>
              </div>
            </div>

            <div class="space-y-2">
              <label class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Duration (minutes)
                <span class="text-danger-500 ml-1">*</span>
              </label>
              <select
                v-model="formData.duration"
                class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                :class="{ 'border-red-500': errors.duration }"
                required
              >
                <option value="15">15 minutes</option>
                <option value="30">30 minutes</option>
                <option value="45">45 minutes</option>
                <option value="60">1 hour</option>
                <option value="90">1.5 hours</option>
                <option value="120">2 hours</option>
              </select>
              <span v-if="errors.duration" class="text-red-500 text-sm mt-1">{{ errors.duration }}</span>
            </div>
          </div>

          <!-- Appointment Details -->
          <div class="space-y-6">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Appointment Details</h3>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                  </svg>
                  Appointment Type
                  <span class="text-danger-500 ml-1">*</span>
                </label>
                <select
                  v-model="formData.type"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  :class="{ 'border-red-500': errors.type }"
                  required
                >
                  <option value="">Select type...</option>
                  <option value="consultation">Consultation</option>
                  <option value="cleaning">Cleaning</option>
                  <option value="filling">Filling</option>
                  <option value="extraction">Extraction</option>
                  <option value="root_canal">Root Canal</option>
                  <option value="crown">Crown</option>
                  <option value="bridge">Bridge</option>
                  <option value="implant">Implant</option>
                  <option value="orthodontics">Orthodontics</option>
                  <option value="emergency">Emergency</option>
                  <option value="follow_up">Follow Up</option>
                  <option value="other">Other</option>
                </select>
                <span v-if="errors.type" class="text-red-500 text-sm mt-1">{{ errors.type }}</span>
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                  </svg>
                  Doctor
                  <span class="text-danger-500 ml-1">*</span>
                </label>
                <select
                  v-model="formData.doctor_id"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                  :class="{ 'border-red-500': errors.doctor_id }"
                  required
                >
                  <option value="">Select doctor...</option>
                  <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                    Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                  </option>
                </select>
                <span v-if="errors.doctor_id" class="text-red-500 text-sm mt-1">{{ errors.doctor_id }}</span>
              </div>
            </div>

            <div v-if="branches.length > 1" class="space-y-2">
              <label class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
                </svg>
                Branch
              </label>
              <select
                v-model="formData.branch_id"
                class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
              >
                <option value="">Select branch...</option>
                <option v-for="branch in branches" :key="branch.id" :value="branch.id">
                  {{ branch.name }}
                </option>
              </select>
            </div>

            <div class="space-y-2">
              <label class="block text-sm font-semibold text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                </svg>
                Notes
              </label>
              <textarea
                v-model="formData.notes"
                rows="3"
                class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                placeholder="Additional notes or special instructions..."
              ></textarea>
            </div>
          </div>

          <!-- Notifications -->
          <div class="space-y-6">
            <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Notifications</h3>

            <div class="space-y-4">
              <div class="flex items-center">
                <input
                  v-model="formData.send_reminder"
                  type="checkbox"
                  id="send_reminder"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded transition-all duration-200"
                />
                <label for="send_reminder" class="ml-3 block text-sm font-medium text-neutral-700">
                  Send appointment reminder
                </label>
              </div>

              <div v-if="formData.send_reminder" class="ml-7 space-y-3">
                <div class="flex items-center">
                  <input
                    v-model="formData.reminder_email"
                    type="checkbox"
                    id="reminder_email"
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded transition-all duration-200"
                  />
                  <label for="reminder_email" class="ml-3 block text-sm text-neutral-600">
                    Email reminder (24 hours before)
                  </label>
                </div>

                <div class="flex items-center">
                  <input
                    v-model="formData.reminder_sms"
                    type="checkbox"
                    id="reminder_sms"
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded transition-all duration-200"
                  />
                  <label for="reminder_sms" class="ml-3 block text-sm text-neutral-600">
                    SMS reminder (2 hours before)
                  </label>
                </div>
              </div>

              <div class="flex items-center">
                <input
                  v-model="formData.send_confirmation"
                  type="checkbox"
                  id="send_confirmation"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded transition-all duration-200"
                />
                <label for="send_confirmation" class="ml-3 block text-sm font-medium text-neutral-700">
                  Send appointment confirmation
                </label>
              </div>
            </div>
          </div>

          <!-- Conflict Warning -->
          <div v-if="hasConflicts" class="bg-danger-50 border border-danger-200 rounded-xl p-4">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-danger-500 mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
              </svg>
              <div>
                <h4 class="text-sm font-medium text-danger-800">Schedule Conflict Detected</h4>
                <p class="text-sm text-danger-700 mt-1">This appointment time conflicts with an existing appointment. Please choose a different time.</p>
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="flex flex-col sm:flex-row justify-end space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t border-neutral-200">
            <button
              type="button"
              @click="$router.go(-1)"
              class="inline-flex items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
            >
              Cancel
            </button>

            <button
              type="submit"
              :disabled="isSubmitting || !isFormValid"
              class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            >
              <div v-if="isSubmitting" class="flex items-center justify-center">
                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                {{ isEditing ? 'Updating...' : 'Creating...' }}
              </div>
              <span v-else>
                {{ isEditing ? 'Update Appointment' : 'Create Appointment' }}
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useClinicStore } from '../stores/clinic'
import { useAppointmentStore } from '../stores/appointment'

const route = useRoute()
const router = useRouter()

const patientStore = usePatientStore()
const clinicStore = useClinicStore()
const appointmentStore = useAppointmentStore()

const isEditing = ref(false)
const isSubmitting = ref(false)
const errors = ref({})
const hasConflicts = ref(false)

const patientSearchQuery = ref('')
const filteredPatients = ref([])
const showPatientDropdown = ref(false)

const formData = ref({
  patient_id: '',
  date: '',
  time: '',
  duration: '30',
  type: '',
  doctor_id: '',
  branch_id: '',
  notes: '',
  send_reminder: true,
  reminder_email: true,
  reminder_sms: false,
  send_confirmation: true
})

const minDate = computed(() => {
  return new Date().toISOString().split('T')[0]
})

const isFormValid = computed(() => {
  return formData.value.patient_id &&
         formData.value.date &&
         formData.value.time &&
         formData.value.duration &&
         formData.value.type &&
         formData.value.doctor_id &&
         getSelectedBranchId()
})

const doctors = computed(() => {
  return clinicStore.doctors || []
})

const branches = computed(() => {
  return clinicStore.branches || []
})

const initializeForm = () => {
  // Check for URL parameters (from calendar time slot clicks)
  if (route.query.date) {
    formData.value.date = route.query.date
  }

  if (route.query.time) {
    const hour = parseInt(route.query.time)
    formData.value.time = `${hour.toString().padStart(2, '0')}:00`
  }

  // Auto-select main branch if available
  const mainBranch = clinicStore.getMainBranch
  if (mainBranch && !formData.value.branch_id) {
    formData.value.branch_id = mainBranch.id.toString()
  }

  // If editing, load appointment data
  if (route.params.id && route.params.id !== 'new') {
    isEditing.value = true
    loadAppointment(route.params.id)
  }
}

const loadAppointment = async (appointmentId) => {
  // In a real app, this would fetch from API
  // For now, simulate loading
  try {
    // Simulate API call - replace with actual API call
    const mockAppointment = {
      id: appointmentId,
      patient_id: 1,
      date: '2024-01-20',
      time: '10:00',
      duration: '60',
      type: 'cleaning',
      doctor_id: 1,
      branch_id: 1,
      notes: 'Regular cleaning appointment',
      send_reminder: true,
      reminder_email: true,
      reminder_sms: false,
      send_confirmation: true
    }

    Object.assign(formData.value, mockAppointment)

  } catch (error) {
    console.error('Error loading appointment:', error)
  }
}

const searchPatients = async () => {
  if (patientSearchQuery.value.length < 2) {
    filteredPatients.value = []
    return
  }

  try {
    const result = await patientStore.searchPatients({
      query: patientSearchQuery.value,
      limit: 10
    })

    if (result.success) {
      filteredPatients.value = result.data.patients || result.data
    }
  } catch (error) {
    console.error('Error searching patients:', error)
  }
}

const selectPatient = (patient) => {
  formData.value.patient_id = patient.id.toString()
  patientSearchQuery.value = `${patient.first_name} ${patient.last_name}`
  showPatientDropdown.value = false
}

const createNewPatient = () => {
  // Navigate to patient creation
  router.push('/patients/new')
}

const getSelectedBranchId = () => {
  // If branch is explicitly selected, use it
  if (formData.value.branch_id) {
    return parseInt(formData.value.branch_id)
  }

  // If only one branch, use it
  if (branches.value.length === 1) {
    return branches.value[0].id
  }

  // If main branch exists, use it
  const mainBranch = clinicStore.getMainBranch
  if (mainBranch) {
    return mainBranch.id
  }

  return null
}

const getEstimatedCost = () => {
  const costMap = {
    consultation: 100.00,
    cleaning: 120.00,
    filling: 150.00,
    extraction: 200.00,
    root_canal: 800.00,
    crown: 1200.00,
    bridge: 1500.00,
    implant: 3000.00,
    orthodontics: 5000.00,
    emergency: 150.00,
    follow_up: 80.00,
    other: 100.00
  }

  return costMap[formData.value.type] || 100.00
}

const checkForConflicts = async () => {
  if (!formData.value.doctor_id || !formData.value.date || !formData.value.time) {
    return
  }

  try {
    // This would check for scheduling conflicts
    // For now, just simulate no conflicts
    hasConflicts.value = false
  } catch (error) {
    console.error('Error checking conflicts:', error)
  }
}

const handleSubmit = async () => {
  isSubmitting.value = true

  try {
    const appointmentData = {
      patient_id: parseInt(formData.value.patient_id),
      doctor_id: parseInt(formData.value.doctor_id),
      branch_id: getSelectedBranchId(),
      type: formData.value.type,
      start_time: `${formData.value.date}T${formData.value.time}:00`,
      end_time: calculateEndTime(),
      duration: parseInt(formData.value.duration),
      pre_appointment_notes: formData.value.notes,
      title: generateAppointmentTitle(),
      description: formData.value.notes || '',
      status: 'scheduled',
      estimated_cost: getEstimatedCost()
    }

    let result
    if (isEditing.value) {
      // Update existing appointment
      console.log('Updating appointment:', appointmentData)
      // For now, we'll implement create only - update would need a separate API endpoint
      result = { success: false, error: 'Update not implemented yet' }
    } else {
      // Create new appointment
      console.log('Creating appointment:', appointmentData)
      result = await appointmentStore.createAppointment(appointmentData)
    }

    if (result.success) {
      // Navigate back to calendar
      router.push('/appointments')
    } else {
      console.error('Failed to save appointment:', result.error)
      // Could show error message to user here
    }
  } catch (error) {
    console.error('Error saving appointment:', error)
    // Handle error - could show toast notification
  } finally {
    isSubmitting.value = false
  }
}

const calculateEndTime = () => {
  const startDateTime = new Date(`${formData.value.date}T${formData.value.time}:00`)
  const endDateTime = new Date(startDateTime.getTime() + parseInt(formData.value.duration) * 60000)
  return endDateTime.toISOString().slice(0, 19)
}

const generateAppointmentTitle = () => {
  const typeLabels = {
    consultation: 'Consultation',
    cleaning: 'Cleaning',
    filling: 'Filling',
    extraction: 'Extraction',
    root_canal: 'Root Canal',
    crown: 'Crown',
    bridge: 'Bridge',
    implant: 'Implant',
    orthodontics: 'Orthodontics',
    emergency: 'Emergency',
    follow_up: 'Follow Up',
    other: 'Appointment'
  }
  return typeLabels[formData.value.type] || 'Appointment'
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showPatientDropdown.value = false
  }
}

onMounted(() => {
  initializeForm()
  document.addEventListener('click', handleClickOutside)

  // Load initial data
  if (!doctors.value.length) {
    clinicStore.fetchDoctors()
  }

  if (!branches.value.length) {
    clinicStore.fetchBranches()
  }
})

watch([() => formData.value.date, () => formData.value.time, () => formData.value.doctor_id], () => {
  if (formData.value.date && formData.value.time && formData.value.doctor_id) {
    checkForConflicts()
  }
})

// Watch for branches to be loaded and auto-select main branch
watch(() => clinicStore.branches, (newBranches) => {
  if (newBranches.length > 0 && !formData.value.branch_id) {
    const mainBranch = clinicStore.getMainBranch
    if (mainBranch) {
      formData.value.branch_id = mainBranch.id.toString()
    } else if (newBranches.length === 1) {
      formData.value.branch_id = newBranches[0].id.toString()
    }
  }
}, { immediate: true })
</script>

<style scoped>
/* Custom styles for patient dropdown */
.patient-option:hover {
  background-color: rgb(249 250 251);
}

/* Ensure proper spacing and alignment */
.space-y-8 > * + * {
  margin-top: 2rem;
}

.space-y-6 > * + * {
  margin-top: 1.5rem;
}

.space-y-4 > * + * {
  margin-top: 1rem;
}

.space-y-3 > * + * {
  margin-top: 0.75rem;
}

.space-y-2 > * + * {
  margin-top: 0.5rem;
}

/* Form transitions */
form input,
form select,
form textarea {
  transition: all 0.2s ease-in-out;
}

/* Focus states */
form input:focus,
form select:focus,
form textarea:focus {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Button hover effects */
button:hover:not(:disabled) {
  transform: translateY(-1px);
}

/* Loading spinner animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Error message styling */
.text-red-500 {
  color: rgb(239 68 68);
}

.text-danger-500 {
  color: rgb(239 68 68);
}

.text-danger-600 {
  color: rgb(220 38 38);
}

.text-danger-700 {
  color: rgb(185 28 28);
}

.text-danger-800 {
  color: rgb(153 27 27);
}

/* Primary color styling */
.text-primary-600 {
  color: rgb(37 99 235);
}

.text-primary-700 {
  color: rgb(29 78 216);
}

.bg-primary-600 {
  background-color: rgb(37 99 235);
}

.bg-primary-700 {
  background-color: rgb(29 78 216);
}

.focus\:ring-primary-500:focus {
  --tw-ring-color: rgb(59 130 246 / var(--tw-ring-opacity));
}

/* Neutral color styling */
.text-neutral-600 {
  color: rgb(75 85 99);
}

.text-neutral-700 {
  color: rgb(55 65 81);
}

.text-neutral-900 {
  color: rgb(17 24 39);
}

.bg-neutral-50 {
  background-color: rgb(249 250 251);
}

.bg-neutral-100 {
  background-color: rgb(243 244 246);
}

.border-neutral-100 {
  border-color: rgb(243 244 246);
}

.border-neutral-200 {
  border-color: rgb(229 231 235);
}

.border-neutral-300 {
  border-color: rgb(209 213 219);
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .space-y-8 > * + * {
    margin-top: 1.5rem;
  }

  .space-y-6 > * + * {
    margin-top: 1.25rem;
  }
}
</style>