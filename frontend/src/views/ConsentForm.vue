<template>
  <div class="consent-form-page min-h-screen bg-gray-50 lg:bg-white">
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
            <h1 class="text-xl font-semibold text-gray-900">Treatment Consent Form</h1>
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
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div>
                    <span class="text-blue-700 font-medium">Name:</span>
                    {{ selectedPatient.first_name }} {{ selectedPatient.last_name }}
                  </div>
                  <div>
                    <span class="text-blue-700 font-medium">DOB:</span>
                    {{ formatDate(selectedPatient.date_of_birth) }}
                  </div>
                  <div>
                    <span class="text-blue-700 font-medium">Phone:</span>
                    {{ selectedPatient.phone }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Treatment Information -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Treatment Information</h2>
              <div class="space-y-4">
                <div>
                  <label for="treatment_type" class="block text-sm font-medium text-gray-700 mb-1">
                    Treatment Type *
                  </label>
                  <select
                    id="treatment_type"
                    v-model="form.treatment_type"
                    @change="loadConsentTemplate"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :class="{ 'border-red-500': errors.treatment_type }"
                  >
                    <option value="">Select treatment type</option>
                    <option value="dental_cleaning">Dental Cleaning</option>
                    <option value="tooth_extraction">Tooth Extraction</option>
                    <option value="root_canal">Root Canal Treatment</option>
                    <option value="dental_filling">Dental Filling</option>
                    <option value="crown_placement">Crown Placement</option>
                    <option value="dental_implant">Dental Implant</option>
                    <option value="orthodontic">Orthodontic Treatment</option>
                    <option value="oral_surgery">Oral Surgery</option>
                    <option value="periodontal">Periodontal Treatment</option>
                    <option value="cosmetic">Cosmetic Dentistry</option>
                    <option value="other">Other</option>
                  </select>
                  <span v-if="errors.treatment_type" class="text-red-500 text-sm mt-1">{{ errors.treatment_type }}</span>
                </div>

                <div v-if="form.treatment_type === 'other'">
                  <label for="treatment_description" class="block text-sm font-medium text-gray-700 mb-1">
                    Treatment Description *
                  </label>
                  <textarea
                    id="treatment_description"
                    v-model="form.treatment_description"
                    rows="3"
                    required
                    placeholder="Please describe the treatment..."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  ></textarea>
                </div>

                <div>
                  <label for="doctor_id" class="block text-sm font-medium text-gray-700 mb-1">
                    Treating Doctor *
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
                  <label for="estimated_duration" class="block text-sm font-medium text-gray-700 mb-1">
                    Estimated Duration (minutes)
                  </label>
                  <input
                    id="estimated_duration"
                    v-model.number="form.estimated_duration"
                    type="number"
                    min="15"
                    max="480"
                    step="15"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="estimated_cost" class="block text-sm font-medium text-gray-700 mb-1">
                    Estimated Cost
                  </label>
                  <input
                    id="estimated_cost"
                    v-model.number="form.estimated_cost"
                    type="number"
                    min="0"
                    step="0.01"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Consent Content -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Consent Information</h2>
              
              <!-- Pre-filled consent template -->
              <div class="prose prose-sm max-w-none bg-gray-50 p-4 rounded-lg border">
                <div v-html="consentTemplate"></div>
              </div>

              <!-- Additional notes -->
              <div class="mt-4">
                <label for="additional_notes" class="block text-sm font-medium text-gray-700 mb-1">
                  Additional Notes or Modifications
                </label>
                <textarea
                  id="additional_notes"
                  v-model="form.additional_notes"
                  rows="4"
                  placeholder="Any additional information or modifications to the standard consent..."
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                ></textarea>
              </div>
            </div>

            <!-- Risks and Complications -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Risks and Complications</h2>
              <div class="space-y-3">
                <div v-for="risk in treatmentRisks" :key="risk.id" class="flex items-start">
                  <input
                    :id="`risk-${risk.id}`"
                    v-model="form.acknowledged_risks"
                    :value="risk.id"
                    type="checkbox"
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label :for="`risk-${risk.id}`" class="ml-2 text-sm text-gray-700">
                    <span class="font-medium">{{ risk.title }}:</span> {{ risk.description }}
                  </label>
                </div>
              </div>
            </div>

            <!-- Alternative Treatments -->
            <div class="form-section" v-if="alternativeTreatments.length > 0">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Alternative Treatments</h2>
              <div class="space-y-2">
                <div v-for="alternative in alternativeTreatments" :key="alternative" class="text-sm text-gray-700">
                  • {{ alternative }}
                </div>
              </div>
              <div class="mt-3">
                <label class="flex items-center">
                  <input
                    v-model="form.alternatives_discussed"
                    type="checkbox"
                    class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    Alternative treatments have been discussed with the patient
                  </span>
                </label>
              </div>
            </div>

            <!-- Patient Consent -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Patient Consent</h2>
              <div class="space-y-3">
                <label class="flex items-start">
                  <input
                    v-model="form.understands_treatment"
                    type="checkbox"
                    required
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    I understand the nature of the proposed treatment and the procedures to be performed
                  </span>
                </label>

                <label class="flex items-start">
                  <input
                    v-model="form.understands_risks"
                    type="checkbox"
                    required
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    I understand the risks, complications, and potential outcomes of the treatment
                  </span>
                </label>

                <label class="flex items-start">
                  <input
                    v-model="form.consents_to_treatment"
                    type="checkbox"
                    required
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    I consent to the proposed treatment and authorize the doctor to proceed
                  </span>
                </label>

                <label class="flex items-start">
                  <input
                    v-model="form.had_opportunity_to_ask"
                    type="checkbox"
                    required
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">
                    I have had the opportunity to ask questions, and all my questions have been answered to my satisfaction
                  </span>
                </label>
              </div>
            </div>

            <!-- Signature Section -->
            <div class="form-section">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Digital Signature</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="patient_signature" class="block text-sm font-medium text-gray-700 mb-1">
                    Patient Signature *
                  </label>
                  <input
                    id="patient_signature"
                    v-model="form.patient_signature"
                    type="text"
                    required
                    placeholder="Type full name as digital signature"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="signature_date" class="block text-sm font-medium text-gray-700 mb-1">
                    Date *
                  </label>
                  <input
                    id="signature_date"
                    v-model="form.signature_date"
                    type="date"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="witness_name" class="block text-sm font-medium text-gray-700 mb-1">
                    Witness Name
                  </label>
                  <input
                    id="witness_name"
                    v-model="form.witness_name"
                    type="text"
                    placeholder="Staff member witnessing consent"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <div>
                  <label for="witness_signature" class="block text-sm font-medium text-gray-700 mb-1">
                    Witness Signature
                  </label>
                  <input
                    id="witness_signature"
                    v-model="form.witness_signature"
                    type="text"
                    placeholder="Type full name as digital signature"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row sm:justify-between items-center space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t border-gray-200">
              <div class="flex space-x-3">
                <button
                  type="button"
                  @click="previewConsent"
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
              
              <button
                type="submit"
                :disabled="loading"
                class="w-full sm:w-auto px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                <span v-if="loading" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Saving...
                </span>
                <span v-else>Save Consent Form</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
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
  treatment_type: '',
  treatment_description: '',
  doctor_id: '',
  estimated_duration: 60,
  estimated_cost: 0,
  additional_notes: '',
  acknowledged_risks: [],
  alternatives_discussed: false,
  understands_treatment: false,
  understands_risks: false,
  consents_to_treatment: false,
  had_opportunity_to_ask: false,
  patient_signature: '',
  signature_date: new Date().toISOString().split('T')[0],
  witness_name: '',
  witness_signature: ''
})

// Mock data - in real app, this would come from stores
const patients = ref([])
const doctors = ref([
  { id: 1, first_name: 'John', last_name: 'Smith' },
  { id: 2, first_name: 'Sarah', last_name: 'Johnson' },
  { id: 3, first_name: 'Michael', last_name: 'Brown' }
])

const patientAppointments = ref([])
const consentTemplate = ref('')
const treatmentRisks = ref([])
const alternativeTreatments = ref([])

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

const loadConsentTemplate = () => {
  const templates = {
    dental_cleaning: {
      content: `
        <h3>Dental Cleaning Consent</h3>
        <p>I understand that dental cleaning involves the removal of plaque, tartar, and stains from my teeth. This procedure may include scaling, root planing, and polishing.</p>
        <p>The cleaning may cause temporary sensitivity and minor bleeding of the gums, which is normal and should subside within a few days.</p>
      `,
      risks: [
        { id: 'sensitivity', title: 'Temporary Sensitivity', description: 'Mild sensitivity to hot/cold temperatures may occur for 24-48 hours' },
        { id: 'bleeding', title: 'Gum Bleeding', description: 'Minor bleeding during and after cleaning is normal' }
      ],
      alternatives: ['No treatment (risks of untreated periodontal disease)']
    },
    tooth_extraction: {
      content: `
        <h3>Tooth Extraction Consent</h3>
        <p>I understand that tooth extraction involves the complete removal of one or more teeth from the socket in the bone.</p>
        <p>Post-operative care instructions will be provided and must be followed to prevent complications such as dry socket.</p>
      `,
      risks: [
        { id: 'pain', title: 'Post-operative Pain', description: 'Pain and discomfort for several days following extraction' },
        { id: 'infection', title: 'Infection', description: 'Risk of infection at the extraction site' },
        { id: 'dry_socket', title: 'Dry Socket', description: 'Delayed healing if blood clot is disturbed' },
        { id: 'nerve_damage', title: 'Nerve Damage', description: 'Rare risk of temporary or permanent nerve damage' }
      ],
      alternatives: ['Root canal treatment', 'No treatment (risks of infected tooth)']
    },
    root_canal: {
      content: `
        <h3>Root Canal Treatment Consent</h3>
        <p>I understand that root canal treatment involves removing infected or damaged nerve tissue from inside the tooth, cleaning and shaping the root canals, and sealing them.</p>
        <p>Multiple appointments may be required to complete the treatment.</p>
      `,
      risks: [
        { id: 'incomplete_healing', title: 'Incomplete Healing', description: 'Treatment may not be successful and extraction may be necessary' },
        { id: 'instrument_breakage', title: 'Instrument Breakage', description: 'Rare risk of instruments breaking inside the tooth' },
        { id: 'perforation', title: 'Root Perforation', description: 'Small risk of creating a hole in the tooth root' }
      ],
      alternatives: ['Tooth extraction and replacement', 'No treatment (risks of abscess and spreading infection)']
    },
    // Add more templates as needed
    other: {
      content: `
        <h3>Treatment Consent</h3>
        <p>I understand the nature of the proposed dental treatment as explained by my doctor.</p>
        <p>All procedures, risks, benefits, and alternatives have been discussed with me.</p>
      `,
      risks: [
        { id: 'general_risks', title: 'General Dental Risks', description: 'All dental procedures carry some inherent risks' }
      ],
      alternatives: ['No treatment']
    }
  }

  const template = templates[form.treatment_type] || templates.other
  consentTemplate.value = template.content
  treatmentRisks.value = template.risks
  alternativeTreatments.value = template.alternatives
}

const validateForm = () => {
  errors.value = {}

  if (!form.patient_id) {
    errors.value.patient_id = 'Patient selection is required'
  }

  if (!form.treatment_type) {
    errors.value.treatment_type = 'Treatment type is required'
  }

  if (!form.doctor_id) {
    errors.value.doctor_id = 'Doctor selection is required'
  }

  if (!form.patient_signature.trim()) {
    errors.value.patient_signature = 'Patient signature is required'
  }

  if (!form.understands_treatment || !form.understands_risks || 
      !form.consents_to_treatment || !form.had_opportunity_to_ask) {
    errors.value.consent = 'All consent checkboxes must be checked'
  }

  return Object.keys(errors.value).length === 0
}

const previewConsent = () => {
  // Open preview in new window or modal
  const previewContent = generateConsentDocument()
  const previewWindow = window.open('', '_blank', 'width=800,height=600')
  previewWindow.document.write(previewContent)
  previewWindow.document.close()
}

const generateConsentDocument = () => {
  const patient = selectedPatient.value
  return `
    <!DOCTYPE html>
    <html>
    <head>
      <title>Consent Form - ${patient?.first_name} ${patient?.last_name}</title>
      <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        h1 { color: #333; border-bottom: 2px solid #333; }
        h2 { color: #555; margin-top: 30px; }
        .patient-info { background: #f5f5f5; padding: 15px; margin: 20px 0; }
        .signature-section { margin-top: 40px; border-top: 1px solid #ccc; padding-top: 20px; }
        @media print { body { margin: 20px; } }
      </style>
    </head>
    <body>
      <h1>Treatment Consent Form</h1>
      <div class="patient-info">
        <h2>Patient Information</h2>
        <p><strong>Name:</strong> ${patient?.first_name} ${patient?.last_name}</p>
        <p><strong>Date of Birth:</strong> ${formatDate(patient?.date_of_birth)}</p>
        <p><strong>Phone:</strong> ${patient?.phone}</p>
      </div>
      
      <h2>Treatment Information</h2>
      <p><strong>Treatment:</strong> ${form.treatment_type.replace('_', ' ').toUpperCase()}</p>
      ${form.treatment_description ? `<p><strong>Description:</strong> ${form.treatment_description}</p>` : ''}
      
      <h2>Consent</h2>
      ${consentTemplate.value}
      
      ${form.additional_notes ? `
        <h2>Additional Notes</h2>
        <p>${form.additional_notes}</p>
      ` : ''}
      
      <div class="signature-section">
        <p><strong>Patient Signature:</strong> ${form.patient_signature}</p>
        <p><strong>Date:</strong> ${formatDate(form.signature_date)}</p>
        ${form.witness_name ? `<p><strong>Witness:</strong> ${form.witness_name}</p>` : ''}
      </div>
    </body>
    </html>
  `
}

const handleSubmit = async () => {
  if (!validateForm()) {
    notificationStore.showError('Please complete all required fields')
    return
  }

  try {
    loading.value = true

    // Save consent form (mock API call)
    const consentData = {
      ...form,
      created_at: new Date().toISOString(),
      clinic_id: 1 // From auth store in real app
    }

    // Mock save
    console.log('Saving consent form:', consentData)

    notificationStore.showSuccess('Consent form saved successfully')
    router.push('/procedures')
  } catch (error) {
    console.error('Error saving consent form:', error)
    notificationStore.showError('Failed to save consent form')
  } finally {
    loading.value = false
  }
}

// Load patients on mount
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

// Watch treatment type changes
watch(() => form.treatment_type, () => {
  if (form.treatment_type) {
    loadConsentTemplate()
  }
})
</script>

<style scoped>
@import "../styles/main.css";

.form-section {
  @apply bg-gray-50 lg:bg-white p-4 lg:p-0 rounded-lg lg:rounded-none mb-6;
}

.form-section h2 {
  @apply border-b border-gray-200 pb-2 mb-4 lg:border-b-0 lg:pb-0;
}

.prose {
  color: #374151;
}

.prose h3 {
  @apply text-lg font-semibold text-gray-900 mt-0 mb-3;
}

.prose p {
  @apply mb-3;
}

/* Mobile optimizations */
@media (max-width: 1024px) {
  .consent-form-page {
    @apply p-0;
  }
  
  .form-section {
    @apply mx-0;
  }
}
</style>