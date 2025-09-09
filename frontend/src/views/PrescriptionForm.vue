<template>
  <div class="prescription-form-page min-h-screen bg-gray-50 lg:bg-white">
    <!-- Mobile: Full screen form -->
    <div class="lg:max-w-4xl lg:mx-auto lg:py-8">
      <div class="bg-white lg:rounded-lg lg:shadow-sm lg:border min-h-screen lg:min-h-0">
        <!-- Form Header - Only show on desktop as mobile has NavigationHeader -->
        <div class="hidden lg:block border-b border-gray-200 px-6 py-4">
          <div class="flex items-center">
            <button 
              @click="$router.go(-1)" 
              class="mr-4 p-2 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <h1 class="text-xl font-semibold text-gray-900">Create Prescription</h1>
          </div>
        </div>

        <!-- Form Content -->
        <div class="p-4 lg:p-6">
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <!-- Patient Selection -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Patient Information</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="patient_id" class="block text-sm font-medium text-gray-700 mb-1">
                    Patient *
                  </label>
                  <select
                    id="patient_id"
                    v-model="form.patient_id"
                    @change="loadPatientInfo"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.patient_id }"
                  >
                    <option value="">Select a patient</option>
                    <option v-for="patient in patients" :key="patient.id" :value="patient.id">
                      {{ patient.first_name }} {{ patient.last_name }} - {{ formatDate(patient.date_of_birth) }}
                    </option>
                  </select>
                  <span v-if="errors.patient_id" class="text-red-500 text-sm mt-1">{{ errors.patient_id }}</span>
                </div>

                <div>
                  <label for="appointment_id" class="block text-sm font-medium text-gray-700 mb-1">
                    Related Appointment (Optional)
                  </label>
                  <select
                    id="appointment_id"
                    v-model="form.appointment_id"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">No specific appointment</option>
                    <option v-for="appointment in patientAppointments" :key="appointment.id" :value="appointment.id">
                      {{ formatDateTime(appointment.appointment_date) }} - {{ appointment.procedure_name }}
                    </option>
                  </select>
                </div>
              </div>

              <!-- Patient Details Preview -->
              <div v-if="selectedPatient" class="mt-4 p-4 bg-blue-50 rounded-lg">
                <h3 class="font-medium text-blue-900 mb-2">Patient Details</h3>
                <div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-sm">
                  <div>
                    <span class="text-blue-700 font-medium">Name:</span>
                    {{ selectedPatient.first_name }} {{ selectedPatient.last_name }}
                  </div>
                  <div>
                    <span class="text-blue-700 font-medium">DOB:</span>
                    {{ formatDate(selectedPatient.date_of_birth) }}
                  </div>
                  <div>
                    <span class="text-blue-700 font-medium">Age:</span>
                    {{ calculateAge(selectedPatient.date_of_birth) }}
                  </div>
                  <div>
                    <span class="text-blue-700 font-medium">Weight:</span>
                    {{ selectedPatient.weight || 'Not specified' }}
                  </div>
                  <div v-if="selectedPatient.allergies" class="md:col-span-2">
                    <span class="text-red-700 font-medium">⚠️ Allergies:</span>
                    {{ selectedPatient.allergies }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Doctor Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Prescribing Doctor</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="doctor_id" class="block text-sm font-medium text-gray-700 mb-1">
                    Doctor *
                  </label>
                  <select
                    id="doctor_id"
                    v-model="form.doctor_id"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.doctor_id }"
                  >
                    <option value="">Select doctor</option>
                    <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                      Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                    </option>
                  </select>
                  <span v-if="errors.doctor_id" class="text-red-500 text-sm mt-1">{{ errors.doctor_id }}</span>
                </div>

                <div>
                  <label for="prescription_date" class="block text-sm font-medium text-gray-700 mb-1">
                    Prescription Date *
                  </label>
                  <input
                    id="prescription_date"
                    v-model="form.prescription_date"
                    type="date"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Medications -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Medications</h2>
              
              <div v-for="(medication, index) in form.medications" :key="index" class="medication-item border border-gray-200 rounded-lg p-4 mb-4">
                <div class="flex justify-between items-center mb-3">
                  <h3 class="text-sm font-medium text-gray-900">Medication {{ index + 1 }}</h3>
                  <button
                    v-if="form.medications.length > 1"
                    @click="removeMedication(index)"
                    type="button"
                    class="text-red-600 hover:text-red-800 text-sm"
                  >
                    Remove
                  </button>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Medication Name *
                    </label>
                    <input
                      v-model="medication.name"
                      type="text"
                      required
                      placeholder="e.g., Amoxicillin"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Strength *
                    </label>
                    <input
                      v-model="medication.strength"
                      type="text"
                      required
                      placeholder="e.g., 500mg"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Form
                    </label>
                    <select
                      v-model="medication.form"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="tablet">Tablet</option>
                      <option value="capsule">Capsule</option>
                      <option value="liquid">Liquid</option>
                      <option value="injection">Injection</option>
                      <option value="cream">Cream/Ointment</option>
                      <option value="drops">Drops</option>
                      <option value="other">Other</option>
                    </select>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Dosage *
                    </label>
                    <input
                      v-model="medication.dosage"
                      type="text"
                      required
                      placeholder="e.g., 1 tablet"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Frequency *
                    </label>
                    <select
                      v-model="medication.frequency"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="once_daily">Once daily</option>
                      <option value="twice_daily">Twice daily</option>
                      <option value="three_times_daily">Three times daily</option>
                      <option value="four_times_daily">Four times daily</option>
                      <option value="every_4_hours">Every 4 hours</option>
                      <option value="every_6_hours">Every 6 hours</option>
                      <option value="every_8_hours">Every 8 hours</option>
                      <option value="every_12_hours">Every 12 hours</option>
                      <option value="as_needed">As needed</option>
                      <option value="custom">Custom</option>
                    </select>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Duration
                    </label>
                    <input
                      v-model="medication.duration"
                      type="text"
                      placeholder="e.g., 7 days"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Quantity
                    </label>
                    <input
                      v-model.number="medication.quantity"
                      type="number"
                      min="1"
                      placeholder="e.g., 21"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      Refills
                    </label>
                    <select
                      v-model.number="medication.refills"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option :value="0">No refills</option>
                      <option :value="1">1 refill</option>
                      <option :value="2">2 refills</option>
                      <option :value="3">3 refills</option>
                      <option :value="4">4 refills</option>
                      <option :value="5">5 refills</option>
                    </select>
                  </div>
                </div>

                <div class="mt-4">
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    Instructions for Use
                  </label>
                  <textarea
                    v-model="medication.instructions"
                    rows="2"
                    placeholder="e.g., Take with food. Complete full course even if symptoms improve."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  ></textarea>
                </div>
              </div>

              <button
                @click="addMedication"
                type="button"
                class="w-full border-2 border-dashed border-gray-300 rounded-lg p-4 text-center text-gray-500 hover:border-blue-400 hover:text-blue-600 transition-colors"
              >
                <svg class="w-5 h-5 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                </svg>
                Add Another Medication
              </button>
            </div>

            <!-- Additional Instructions -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Additional Instructions</h2>
              <textarea
                v-model="form.additional_instructions"
                rows="4"
                placeholder="Any additional instructions or precautions for the patient..."
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              ></textarea>
            </div>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row sm:justify-between items-center space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t border-gray-200">
              <div class="flex space-x-3">
                <button
                  type="button"
                  @click="previewPrescription"
                  class="px-4 py-2 text-blue-700 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors"
                >
                  Preview
                </button>
                <button
                  type="button"
                  @click="$router.go(-1)"
                  class="px-4 py-2 text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-colors"
                >
                  Cancel
                </button>
              </div>
              
              <div class="flex space-x-3">
                <button
                  type="button"
                  @click="saveDraft"
                  :disabled="loading"
                  class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-colors"
                >
                  Save Draft
                </button>
                <button
                  type="submit"
                  :disabled="loading"
                  class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <span v-if="loading" class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Saving...
                  </span>
                  <span v-else>Generate Prescription</span>
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useAppointmentStore } from '../stores/appointment'
import { useNotificationStore } from '../stores/notification'

const route = useRoute()
const router = useRouter()
const patientStore = usePatientStore()
const appointmentStore = useAppointmentStore()
const notificationStore = useNotificationStore()

// Component state
const loading = ref(false)
const errors = ref({})

// Form data
const form = reactive({
  patient_id: '',
  appointment_id: '',
  doctor_id: '',
  prescription_date: new Date().toISOString().split('T')[0],
  medications: [
    {
      name: '',
      strength: '',
      form: 'tablet',
      dosage: '',
      frequency: 'twice_daily',
      duration: '',
      quantity: null,
      refills: 0,
      instructions: ''
    }
  ],
  additional_instructions: ''
})

// Mock data - in real app, this would come from stores
const patients = ref([])
const doctors = ref([
  { id: 1, first_name: 'John', last_name: 'Smith' },
  { id: 2, first_name: 'Sarah', last_name: 'Johnson' },
  { id: 3, first_name: 'Michael', last_name: 'Brown' }
])

const patientAppointments = ref([])

// Computed properties
const selectedPatient = computed(() => {
  return patients.value.find(p => p.id === parseInt(form.patient_id))
})

// Methods
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const calculateAge = (birthDate) => {
  const today = new Date()
  const birth = new Date(birthDate)
  let age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--
  }
  
  return age + ' years'
}

const loadPatientInfo = async () => {
  if (!form.patient_id) {
    patientAppointments.value = []
    return
  }

  try {
    // Load patient appointments
    const appointments = await appointmentStore.fetchPatientAppointments(form.patient_id)
    patientAppointments.value = appointments || []
  } catch (error) {
    console.error('Error loading patient appointments:', error)
  }
}

const addMedication = () => {
  form.medications.push({
    name: '',
    strength: '',
    form: 'tablet',
    dosage: '',
    frequency: 'twice_daily',
    duration: '',
    quantity: null,
    refills: 0,
    instructions: ''
  })
}

const removeMedication = (index) => {
  if (form.medications.length > 1) {
    form.medications.splice(index, 1)
  }
}

const validateForm = () => {
  errors.value = {}

  if (!form.patient_id) {
    errors.value.patient_id = 'Patient selection is required'
  }

  if (!form.doctor_id) {
    errors.value.doctor_id = 'Doctor selection is required'
  }

  // Validate medications
  let hasValidMedication = false
  form.medications.forEach((med, index) => {
    if (med.name && med.strength && med.dosage && med.frequency) {
      hasValidMedication = true
    }
  })

  if (!hasValidMedication) {
    errors.value.medications = 'At least one complete medication is required'
  }

  return Object.keys(errors.value).length === 0
}

const generatePrescriptionNumber = () => {
  const timestamp = Date.now().toString(36).toUpperCase()
  const random = Math.random().toString(36).substr(2, 4).toUpperCase()
  return `RX${timestamp}${random}`
}

const previewPrescription = () => {
  if (!validateForm()) {
    notificationStore.showError('Please complete all required fields')
    return
  }

  const prescriptionDocument = generatePrescriptionDocument()
  const previewWindow = window.open('', '_blank', 'width=800,height=600')
  previewWindow.document.write(prescriptionDocument)
  previewWindow.document.close()
}

const generatePrescriptionDocument = () => {
  const patient = selectedPatient.value
  const doctor = doctors.value.find(d => d.id === parseInt(form.doctor_id))
  const prescriptionNumber = generatePrescriptionNumber()

  return `
    <!DOCTYPE html>
    <html>
    <head>
      <title>Prescription - ${patient?.first_name} ${patient?.last_name}</title>
      <style>
        body { 
          font-family: Arial, sans-serif; 
          margin: 40px; 
          line-height: 1.4; 
          color: #333;
        }
        .header { 
          border-bottom: 2px solid #333; 
          padding-bottom: 20px; 
          margin-bottom: 30px;
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        .logo { font-size: 24px; font-weight: bold; }
        .rx-number { font-size: 14px; color: #666; }
        .patient-info, .doctor-info { 
          background: #f8f9fa; 
          padding: 15px; 
          margin: 20px 0; 
          border-radius: 5px;
        }
        .medications { margin: 30px 0; }
        .medication-item { 
          border: 1px solid #ddd; 
          margin: 15px 0; 
          padding: 15px; 
          border-radius: 5px;
        }
        .medication-name { font-size: 18px; font-weight: bold; margin-bottom: 10px; }
        .medication-details { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
        .instructions { 
          margin-top: 15px; 
          padding: 10px; 
          background: #fffbf0; 
          border-left: 4px solid #ffc107;
        }
        .signature { 
          margin-top: 40px; 
          border-top: 1px solid #ccc; 
          padding-top: 20px;
          display: flex;
          justify-content: space-between;
        }
        @media print { 
          body { margin: 20px; }
          .no-print { display: none; }
        }
      </style>
    </head>
    <body>
      <div class="header">
        <div class="logo">Dental Clinic</div>
        <div class="rx-number">Rx #${prescriptionNumber}</div>
      </div>

      <div class="doctor-info">
        <h2>Prescribing Doctor</h2>
        <p><strong>Dr. ${doctor?.first_name} ${doctor?.last_name}</strong></p>
        <p>License #: DC123456</p>
        <p>Date: ${formatDate(form.prescription_date)}</p>
      </div>
      
      <div class="patient-info">
        <h2>Patient Information</h2>
        <p><strong>Name:</strong> ${patient?.first_name} ${patient?.last_name}</p>
        <p><strong>Date of Birth:</strong> ${formatDate(patient?.date_of_birth)} (Age: ${calculateAge(patient?.date_of_birth)})</p>
        <p><strong>Phone:</strong> ${patient?.phone}</p>
        ${patient?.allergies ? `<p style="color: #dc3545;"><strong>⚠️ Allergies:</strong> ${patient.allergies}</p>` : ''}
      </div>
      
      <div class="medications">
        <h2>Medications Prescribed</h2>
        ${form.medications.filter(med => med.name).map((med, index) => `
          <div class="medication-item">
            <div class="medication-name">${index + 1}. ${med.name} ${med.strength}</div>
            <div class="medication-details">
              <div><strong>Form:</strong> ${med.form}</div>
              <div><strong>Quantity:</strong> ${med.quantity || 'As directed'}</div>
              <div><strong>Dosage:</strong> ${med.dosage}</div>
              <div><strong>Frequency:</strong> ${med.frequency.replace('_', ' ')}</div>
              <div><strong>Duration:</strong> ${med.duration || 'As directed'}</div>
              <div><strong>Refills:</strong> ${med.refills}</div>
            </div>
            ${med.instructions ? `
              <div class="instructions">
                <strong>Instructions:</strong> ${med.instructions}
              </div>
            ` : ''}
          </div>
        `).join('')}
      </div>
      
      ${form.additional_instructions ? `
        <div class="instructions">
          <h3>Additional Instructions</h3>
          <p>${form.additional_instructions}</p>
        </div>
      ` : ''}
      
      <div class="signature">
        <div>
          <p>Doctor's Signature: _________________________</p>
          <p>Dr. ${doctor?.first_name} ${doctor?.last_name}</p>
        </div>
        <div>
          <p>Date: ${formatDate(form.prescription_date)}</p>
        </div>
      </div>

      <div style="margin-top: 40px; font-size: 12px; color: #666; text-align: center;">
        <p>This prescription was generated electronically. Please present this to your pharmacy.</p>
      </div>
    </body>
    </html>
  `
}

const saveDraft = async () => {
  try {
    loading.value = true

    // Save as draft (mock API call)
    const draftData = {
      ...form,
      status: 'draft',
      created_at: new Date().toISOString(),
      clinic_id: 1 // From auth store in real app
    }

    console.log('Saving prescription draft:', draftData)
    notificationStore.showSuccess('Prescription draft saved successfully')
  } catch (error) {
    console.error('Error saving draft:', error)
    notificationStore.showError('Failed to save draft')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!validateForm()) {
    notificationStore.showError('Please complete all required fields')
    return
  }

  try {
    loading.value = true

    // Generate prescription (mock API call)
    const prescriptionData = {
      ...form,
      prescription_number: generatePrescriptionNumber(),
      status: 'active',
      created_at: new Date().toISOString(),
      clinic_id: 1 // From auth store in real app
    }

    console.log('Generating prescription:', prescriptionData)
    
    // Auto-open preview
    setTimeout(() => {
      previewPrescription()
    }, 500)

    notificationStore.showSuccess('Prescription generated successfully')
    router.push('/procedures')
  } catch (error) {
    console.error('Error generating prescription:', error)
    notificationStore.showError('Failed to generate prescription')
  } finally {
    loading.value = false
  }
}

// Load data on mount
onMounted(async () => {
  try {
    await patientStore.fetchPatients()
    patients.value = patientStore.patients

    // Pre-select patient if passed via route params
    if (route.query.patient_id) {
      form.patient_id = route.query.patient_id
      await loadPatientInfo()
    }
  } catch (error) {
    console.error('Error loading patients:', error)
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.form-section {
  @apply bg-gray-50 lg:bg-white p-4 lg:p-0 rounded-lg lg:rounded-none mb-6;
}

.form-section h2 {
  @apply border-b border-gray-200 pb-2 mb-4 lg:border-b-0 lg:pb-0;
}

.medication-item {
  @apply bg-white border-gray-200;
}

/* Mobile optimizations */
@media (max-width: 1024px) {
  .prescription-form-page {
    @apply p-0;
  }
  
  .form-section {
    @apply mx-0;
  }
}
</style>