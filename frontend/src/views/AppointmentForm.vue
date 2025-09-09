<template>
  <div class="appointment-form-page">
    <!-- Page Header -->
    <div class="page-header flex items-center mb-6">
      <button
        @click="$router.go(-1)"
        class="mr-4 p-2 rounded-lg hover:bg-gray-100 transition-colors"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
        </svg>
      </button>
      
      <div>
        <h1 class="text-2xl font-bold text-gray-900">
          {{ isEditing ? 'Edit Appointment' : 'New Appointment' }}
        </h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update appointment details' : 'Schedule a new patient appointment' }}
        </p>
      </div>
    </div>

    <!-- Form -->
    <div class="form-container bg-white rounded-lg shadow-sm border">
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <!-- Patient Selection -->
        <div class="form-section">
          <h3 class="section-title text-lg font-semibold text-gray-900 mb-4">Patient Information</h3>
          
          <div class="form-group">
            <label class="form-label">Patient *</label>
            <div class="relative">
              <input
                v-model="patientSearchQuery"
                @input="searchPatients"
                @focus="showPatientDropdown = true"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.patient_id }"
                placeholder="Search for patient..."
                autocomplete="off"
                required
              />
              
              <!-- Patient Dropdown -->
              <div 
                v-if="showPatientDropdown && (filteredPatients.length > 0 || patientSearchQuery)"
                class="absolute z-10 w-full bg-white border border-gray-300 rounded-md shadow-lg max-h-48 overflow-y-auto mt-1"
              >
                <div 
                  v-for="patient in filteredPatients" 
                  :key="patient.id"
                  @click="selectPatient(patient)"
                  class="patient-option px-3 py-2 hover:bg-gray-100 cursor-pointer border-b last:border-b-0"
                >
                  <div class="font-medium text-gray-900">
                    {{ patient.first_name }} {{ patient.last_name }}
                  </div>
                  <div class="text-sm text-gray-600">
                    {{ patient.phone }} • {{ patient.email }}
                  </div>
                </div>
                
                <div 
                  v-if="filteredPatients.length === 0 && patientSearchQuery"
                  class="px-3 py-2 text-gray-500 text-sm"
                >
                  No patients found. 
                  <button type="button" @click="createNewPatient" class="text-blue-600 hover:text-blue-700">
                    Create new patient
                  </button>
                </div>
              </div>
            </div>
            <div v-if="errors.patient_id" class="error-message">{{ errors.patient_id }}</div>
          </div>
        </div>

        <!-- Date & Time -->
        <div class="form-section">
          <h3 class="section-title text-lg font-semibold text-gray-900 mb-4">Schedule</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-group">
              <label class="form-label">Date *</label>
              <input
                v-model="formData.date"
                type="date"
                class="form-input"
                :class="{ 'border-red-500': errors.date }"
                required
                :min="minDate"
              />
              <div v-if="errors.date" class="error-message">{{ errors.date }}</div>
            </div>
            
            <div class="form-group">
              <label class="form-label">Time *</label>
              <input
                v-model="formData.time"
                type="time"
                class="form-input"
                :class="{ 'border-red-500': errors.time }"
                required
              />
              <div v-if="errors.time" class="error-message">{{ errors.time }}</div>
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">Duration (minutes) *</label>
            <select
              v-model="formData.duration"
              class="form-input"
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
            <div v-if="errors.duration" class="error-message">{{ errors.duration }}</div>
          </div>
        </div>

        <!-- Appointment Details -->
        <div class="form-section">
          <h3 class="section-title text-lg font-semibold text-gray-900 mb-4">Appointment Details</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-group">
              <label class="form-label">Appointment Type *</label>
              <select
                v-model="formData.type"
                class="form-input"
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
              <div v-if="errors.type" class="error-message">{{ errors.type }}</div>
            </div>
            
            <div class="form-group">
              <label class="form-label">Doctor *</label>
              <select
                v-model="formData.doctor_id"
                class="form-input"
                :class="{ 'border-red-500': errors.doctor_id }"
                required
              >
                <option value="">Select doctor...</option>
                <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                  Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                </option>
              </select>
              <div v-if="errors.doctor_id" class="error-message">{{ errors.doctor_id }}</div>
            </div>
          </div>
          
          <div class="form-group" v-if="branches.length > 1">
            <label class="form-label">Branch</label>
            <select
              v-model="formData.branch_id"
              class="form-input"
            >
              <option value="">Select branch...</option>
              <option v-for="branch in branches" :key="branch.id" :value="branch.id">
                {{ branch.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="form-label">Notes</label>
            <textarea
              v-model="formData.notes"
              rows="3"
              class="form-input resize-none"
              placeholder="Additional notes or special instructions..."
            ></textarea>
          </div>
        </div>

        <!-- Contact & Reminders -->
        <div class="form-section">
          <h3 class="section-title text-lg font-semibold text-gray-900 mb-4">Notifications</h3>
          
          <div class="space-y-3">
            <div class="form-check">
              <input
                v-model="formData.send_reminder"
                type="checkbox"
                id="send_reminder"
                class="form-checkbox"
              />
              <label for="send_reminder" class="form-check-label">
                Send appointment reminder
              </label>
            </div>
            
            <div v-if="formData.send_reminder" class="ml-6 space-y-2">
              <div class="form-check">
                <input
                  v-model="formData.reminder_email"
                  type="checkbox"
                  id="reminder_email"
                  class="form-checkbox"
                />
                <label for="reminder_email" class="form-check-label">
                  Email reminder (24 hours before)
                </label>
              </div>
              
              <div class="form-check">
                <input
                  v-model="formData.reminder_sms"
                  type="checkbox"
                  id="reminder_sms"
                  class="form-checkbox"
                />
                <label for="reminder_sms" class="form-check-label">
                  SMS reminder (2 hours before)
                </label>
              </div>
            </div>
            
            <div class="form-check">
              <input
                v-model="formData.send_confirmation"
                type="checkbox"
                id="send_confirmation"
                class="form-checkbox"
              />
              <label for="send_confirmation" class="form-check-label">
                Send appointment confirmation
              </label>
            </div>
          </div>
        </div>

        <!-- Conflict Warning -->
        <div v-if="hasConflicts" class="conflict-warning bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex items-start">
            <svg class="w-5 h-5 text-red-500 mt-0.5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
            </svg>
            <div>
              <h4 class="text-sm font-medium text-red-800">Schedule Conflict Detected</h4>
              <p class="text-sm text-red-700 mt-1">This appointment time conflicts with an existing appointment. Please choose a different time.</p>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions flex flex-col sm:flex-row justify-end space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t">
          <button
            type="button"
            @click="$router.go(-1)"
            class="btn btn-secondary w-full sm:w-auto"
          >
            Cancel
          </button>
          
          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="btn btn-primary w-full sm:w-auto"
            :class="{ 'opacity-50 cursor-not-allowed': isSubmitting || !isFormValid }"
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
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useClinicStore } from '../stores/clinic'

const route = useRoute()
const router = useRouter()

const patientStore = usePatientStore()
const clinicStore = useClinicStore()

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
         formData.value.doctor_id
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
    
    // Set patient search query if patient is loaded
    if (mockAppointment.patient_id) {
      patientSearchQuery.value = 'John Smith' // This would come from the patient data
    }
  } catch (error) {
    console.error('Error loading appointment:', error)
  }
}

const searchPatients = async () => {
  if (patientSearchQuery.value.length >= 2) {
    await patientStore.searchPatients(patientSearchQuery.value)
    filteredPatients.value = patientStore.patients.slice(0, 10)
  } else {
    filteredPatients.value = []
  }
}

const selectPatient = (patient) => {
  formData.value.patient_id = patient.id
  patientSearchQuery.value = `${patient.first_name} ${patient.last_name}`
  showPatientDropdown.value = false
  filteredPatients.value = []
}

const createNewPatient = () => {
  // Navigate to new patient form with pre-filled name
  router.push(`/patients/new?name=${encodeURIComponent(patientSearchQuery.value)}`)
}

const validateForm = () => {
  errors.value = {}
  
  if (!formData.value.patient_id) {
    errors.value.patient_id = 'Please select a patient'
  }
  
  if (!formData.value.date) {
    errors.value.date = 'Please select a date'
  }
  
  if (!formData.value.time) {
    errors.value.time = 'Please select a time'
  }
  
  if (!formData.value.duration) {
    errors.value.duration = 'Please select duration'
  }
  
  if (!formData.value.type) {
    errors.value.type = 'Please select appointment type'
  }
  
  if (!formData.value.doctor_id) {
    errors.value.doctor_id = 'Please select a doctor'
  }
  
  return Object.keys(errors.value).length === 0
}

const checkForConflicts = async () => {
  // This would check for appointment conflicts with the API
  // For now, just simulate
  hasConflicts.value = false
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }
  
  isSubmitting.value = true
  
  try {
    const appointmentData = {
      patient_id: parseInt(formData.value.patient_id),
      doctor_id: parseInt(formData.value.doctor_id),
      branch_id: formData.value.branch_id ? parseInt(formData.value.branch_id) : null,
      type: formData.value.type,
      start_time: `${formData.value.date}T${formData.value.time}:00`,
      end_time: calculateEndTime(),
      notes: formData.value.notes,
      send_reminder: formData.value.send_reminder,
      reminder_email: formData.value.reminder_email,
      reminder_sms: formData.value.reminder_sms,
      send_confirmation: formData.value.send_confirmation
    }
    
    if (isEditing.value) {
      // Update existing appointment
      console.log('Updating appointment:', appointmentData)
      // API call would go here
    } else {
      // Create new appointment
      console.log('Creating appointment:', appointmentData)
      // API call would go here
    }
    
    // Navigate back to calendar
    router.push('/appointments')
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
</script>

<style scoped>
@import "../styles/main.css";

.form-label {
  @apply block text-sm font-medium text-gray-700 mb-1;
}

.form-input {
  @apply w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors;
}

.form-checkbox {
  @apply w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 focus:ring-2;
}

.form-check {
  @apply flex items-center;
}

.form-check-label {
  @apply ml-2 text-sm text-gray-700 cursor-pointer;
}

.error-message {
  @apply text-red-600 text-xs mt-1;
}

.btn {
  @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500;
}

.section-title {
  @apply border-b border-gray-200 pb-2;
}

.form-section {
  @apply border-b border-gray-200 pb-6 last:border-b-0 last:pb-0;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .form-container {
    @apply mx-4 rounded-lg;
  }
  
  .page-header {
    @apply mx-4;
  }
  
  .grid {
    @apply grid-cols-1;
  }
}
</style>