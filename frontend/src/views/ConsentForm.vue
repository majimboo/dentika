<template>
  <div class="consent-form-page">
    <!-- Page Header -->
    <div class="page-header flex items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">
          {{ isEditing ? 'Edit Consent Form' : 'New Consent Form' }}
        </h1>
        <p class="text-gray-600">
          {{ isEditing ? 'Update consent form details' : 'Create a new consent form for patient treatment' }}
        </p>
      </div>
    </div>

    <!-- Form -->
    <div class="form-container bg-white rounded-lg shadow-sm border border-gray-200">
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <!-- Basic Information -->
        <div class="form-section">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Basic Information</h3>

          <div class="space-y-4">
            <div class="form-group">
              <label class="form-label">Consent Form Title *</label>
              <input
                v-model="formData.title"
                type="text"
                class="form-input"
                :class="{ 'border-red-500': errors.title }"
                placeholder="e.g. Dental Treatment Consent"
                required
              />
              <div v-if="errors.title" class="error-message">{{ errors.title }}</div>
            </div>

            <div class="form-group">
              <label class="form-label">Description</label>
              <textarea
                v-model="formData.description"
                rows="3"
                class="form-input resize-none"
                placeholder="Describe the consent form purpose..."
              ></textarea>
            </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-group">
                 <label class="form-label">Patient</label>
                 <div class="form-input bg-gray-50 text-gray-700 cursor-not-allowed">
                   {{ selectedPatient?.first_name }} {{ selectedPatient?.last_name }}
                 </div>
                 <p class="text-sm text-gray-500 mt-1">Patient is automatically selected from the current context</p>
               </div>

               <div class="form-group">
                 <label class="form-label">Consent Template *</label>
                 <select
                   v-model="formData.consent_template_id"
                   class="form-input"
                   :class="{ 'border-red-500': errors.consent_template_id }"
                   required
                 >
                   <option value="">Select consent template...</option>
                   <option
                     v-for="template in consentStore.getActiveConsentTemplates"
                     :key="template.id"
                     :value="template.id"
                   >
                     {{ template.code ? template.code + ' - ' + template.name : template.name }}
                   </option>
                 </select>
                 <div v-if="errors.consent_template_id" class="error-message">{{ errors.consent_template_id }}</div>
               </div>
             </div>
          </div>
        </div>

          <!-- Consent Template Content -->
          <div class="form-section" v-if="selectedConsentTemplate">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Consent Form Content</h3>

            <div class="bg-gray-50 p-4 rounded-lg border">
              <div v-html="populatedConsentContent" class="prose prose-sm max-w-none"></div>
            </div>
          </div>

         <!-- Treatment Details -->
         <div class="form-section">
           <h3 class="text-lg font-semibold text-gray-900 mb-4">Treatment Details</h3>

           <div class="space-y-4">
             <div class="form-group">
               <label class="form-label">Procedure Description</label>
               <textarea
                 v-model="formData.procedure_description"
                 rows="3"
                 class="form-input resize-none"
                 placeholder="Detailed description of the proposed procedure..."
               ></textarea>
             </div>

             <div class="form-group">
               <label class="form-label">Risks and Complications</label>
               <textarea
                 v-model="formData.risks"
                 rows="3"
                 class="form-input resize-none"
                 placeholder="List potential risks and complications..."
               ></textarea>
             </div>

             <div class="form-group">
               <label class="form-label">Alternatives</label>
               <textarea
                 v-model="formData.alternatives"
                 rows="2"
                 class="form-input resize-none"
                 placeholder="Alternative treatments available..."
               ></textarea>
             </div>

             <div class="form-group">
               <label class="form-label">Post-Care Instructions</label>
               <textarea
                 v-model="formData.post_care_instructions"
                 rows="3"
                 class="form-input resize-none"
                 placeholder="Instructions for after the procedure..."
               ></textarea>
             </div>
           </div>
         </div>

         <!-- Patient Agreement -->
         <div class="form-section">
           <h3 class="text-lg font-semibold text-gray-900 mb-4">Patient Agreement</h3>

           <div class="space-y-4">
             <div class="form-check">
               <input
                 v-model="formData.understands_treatment"
                 type="checkbox"
                 id="understands_treatment"
                 class="form-checkbox"
               />
               <label for="understands_treatment" class="form-check-label">
                 I understand the nature of the proposed treatment as explained by my doctor
               </label>
             </div>

             <div class="form-check">
               <input
                 v-model="formData.understands_risks"
                 type="checkbox"
                 id="understands_risks"
                 class="form-checkbox"
               />
               <label for="understands_risks" class="form-check-label">
                 I understand the risks and complications that have been explained to me
               </label>
             </div>

             <div class="form-check">
               <input
                 v-model="formData.consents_to_treatment"
                 type="checkbox"
                 id="consents_to_treatment"
                 class="form-checkbox"
               />
               <label for="consents_to_treatment" class="form-check-label">
                 I voluntarily consent to the proposed treatment
               </label>
             </div>

             <div class="form-check">
               <input
                 v-model="formData.had_opportunity_to_ask"
                 type="checkbox"
                 id="had_opportunity_to_ask"
                 class="form-checkbox"
               />
               <label for="had_opportunity_to_ask" class="form-check-label">
                 I have had the opportunity to ask questions and all my questions have been answered
               </label>
             </div>
           </div>
         </div>

         <!-- Signatures -->
         <div class="form-section">
           <h3 class="text-lg font-semibold text-gray-900 mb-4">Signatures</h3>

           <div class="space-y-4">
             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-group">
                 <label class="form-label">Patient Signature *</label>
                 <input
                   v-model="formData.patient_signature"
                   type="text"
                   class="form-input"
                   :class="{ 'border-red-500': errors.patient_signature }"
                   placeholder="Patient's full name as signature"
                   required
                 />
                 <div v-if="errors.patient_signature" class="error-message">{{ errors.patient_signature }}</div>
               </div>

               <div class="form-group">
                 <label class="form-label">Date</label>
                 <input
                   v-model="formData.signature_date"
                   type="date"
                   class="form-input"
                   readonly
                 />
               </div>
             </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-group">
                 <label class="form-label">Witness Name (Optional)</label>
                 <input
                   v-model="formData.witness_name"
                   type="text"
                   class="form-input"
                   placeholder="Witness's full name"
                 />
               </div>

               <div class="form-group">
                 <label class="form-label">Witness Signature (Optional)</label>
                 <input
                   v-model="formData.witness_signature"
                   type="text"
                   class="form-input"
                   placeholder="Witness's signature"
                 />
               </div>
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
            :disabled="isSubmitting"
            class="btn btn-primary w-full sm:w-auto"
            :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
          >
            <div v-if="isSubmitting" class="flex items-center justify-center">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              {{ isEditing ? 'Updating...' : 'Creating...' }}
            </div>
            <span v-else>
              {{ isEditing ? 'Update Consent Form' : 'Create Consent Form' }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useAppointmentStore } from '../stores/appointment'
import { useNotificationStore } from '../stores/notification'
import { useConsentStore } from '../stores/consent'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const patientStore = usePatientStore()
const appointmentStore = useAppointmentStore()
const notificationStore = useNotificationStore()
const consentStore = useConsentStore()
const userStore = useUserStore()
const authStore = useAuthStore()

// Component state
const loading = ref(false)
const isSubmitting = ref(false)
const isEditing = ref(false)
const errors = ref({})

// Form data
const formData = reactive({
  consent_template_id: '',
  procedure_description: '',
  risks: '',
  alternatives: '',
  post_care_instructions: '',
  understands_treatment: false,
  understands_risks: false,
  consents_to_treatment: false,
  had_opportunity_to_ask: false,
  patient_signature: '',
  signature_date: new Date().toISOString().split('T')[0],
  witness_name: '',
  witness_signature: ''
})

// Auto-population state
const autoPopulateFields = ref(true)

// Data from stores
const patients = ref([])
const patientAppointments = ref([])
const consentTemplate = ref('')
const treatmentRisks = ref([])
const alternativeTreatments = ref([])

// Computed properties
const selectedPatient = computed(() => {
  const patientId = route.params.patientId
  return patients.value.find(p => p.id === parseInt(patientId))
})

const selectedConsentTemplate = computed(() => {
  if (!formData.consent_template_id) return null
  return consentStore.getConsentTemplateById(parseInt(formData.consent_template_id))
})

const populatedConsentContent = computed(() => {
  if (!selectedConsentTemplate.value || !selectedPatient.value) return ''

  const patient = selectedPatient.value
  const template = selectedConsentTemplate.value
  const currentDate = new Date().toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  // Replace common placeholders
  let content = template.content
    .replace(/\[PATIENT_NAME\]/g, `${patient.first_name} ${patient.last_name}`)
    .replace(/\[PATIENT_FIRST_NAME\]/g, patient.first_name)
    .replace(/\[PATIENT_LAST_NAME\]/g, patient.last_name)
    .replace(/\[CURRENT_DATE\]/g, currentDate)
    .replace(/\[TODAY\]/g, currentDate)
    .replace(/\[DATE\]/g, currentDate)
    .replace(/\[DOCTOR_NAME\]/g, authStore.user?.first_name ? `${authStore.user.first_name} ${authStore.user.last_name}` : 'Doctor')
    .replace(/\[CLINIC_NAME\]/g, authStore.userClinic?.name || 'Clinic')

  // Replace patient-specific placeholders if available
  if (patient.date_of_birth) {
    const dob = new Date(patient.date_of_birth).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
    content = content.replace(/\[PATIENT_DOB\]/g, dob)
  }

  if (patient.phone) {
    content = content.replace(/\[PATIENT_PHONE\]/g, patient.phone)
  }

  if (patient.email) {
    content = content.replace(/\[PATIENT_EMAIL\]/g, patient.email)
  }

  return content
})

// Methods
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const loadPatientInfo = async () => {
  const patientId = route.params.patientId
  if (!patientId) {
    patientAppointments.value = []
    return
  }

  try {
    // Load patient appointments using the existing fetchAppointments method
    await appointmentStore.fetchAppointments({ patient_id: patientId })
    patientAppointments.value = appointmentStore.appointments || []
  } catch (error) {
    console.error('Error loading patient appointments:', error)
  }
}

const loadConsentTemplate = () => {
  const selectedTemplate = selectedConsentTemplate.value

  if (selectedTemplate) {
    consentTemplate.value = selectedTemplate.content
  } else {
    // Default template when no template is selected
    consentTemplate.value = `
      <h3>Treatment Consent</h3>
      <p>I understand the nature of the proposed dental treatment as explained by my doctor.</p>
      <p>All procedures, risks, benefits, and alternatives have been discussed with me.</p>
    `
  }
}

const handleAutoPopulateToggle = () => {
  // Auto-populate logic based on selected procedure template
  if (autoPopulateFields.value && formData.procedure_template_id) {
    const selectedTemplate = selectedProcedureTemplate.value
    if (selectedTemplate) {
      formData.estimated_duration = selectedTemplate.estimated_duration || 60
      formData.estimated_cost = selectedTemplate.default_cost || 0
    }
  }
}

const validateForm = () => {
  errors.value = {}

  if (!formData.consent_template_id) {
    errors.value.consent_template_id = 'Consent template selection is required'
  }

  if (!formData.patient_signature.trim()) {
    errors.value.patient_signature = 'Patient signature is required'
  }

  if (!formData.understands_treatment || !formData.understands_risks ||
      !formData.consents_to_treatment || !formData.had_opportunity_to_ask) {
    errors.value.agreement = 'All patient agreement checkboxes must be checked'
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
  const selectedTemplate = selectedConsentTemplate.value

  if (!patient || !selectedTemplate) return ''

  const patientName = `${patient.first_name} ${patient.last_name}`
  const currentDate = new Date().toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  // Replace placeholders in template content
  let populatedContent = selectedTemplate.content
    .replace(/\[PATIENT_NAME\]/g, patientName)
    .replace(/\[PATIENT_FIRST_NAME\]/g, patient.first_name)
    .replace(/\[PATIENT_LAST_NAME\]/g, patient.last_name)
    .replace(/\[CURRENT_DATE\]/g, currentDate)
    .replace(/\[TODAY\]/g, currentDate)
    .replace(/\[DATE\]/g, currentDate)
    .replace(/\[DOCTOR_NAME\]/g, authStore.user?.first_name ? `${authStore.user.first_name} ${authStore.user.last_name}` : 'Doctor')
    .replace(/\[CLINIC_NAME\]/g, authStore.userClinic?.name || 'Clinic')

  // Replace patient-specific placeholders if available
  if (patient.date_of_birth) {
    const dob = new Date(patient.date_of_birth).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
    populatedContent = populatedContent.replace(/\[PATIENT_DOB\]/g, dob)
  }

  if (patient.phone) {
    populatedContent = populatedContent.replace(/\[PATIENT_PHONE\]/g, patient.phone)
  }

  if (patient.email) {
    populatedContent = populatedContent.replace(/\[PATIENT_EMAIL\]/g, patient.email)
  }

  const procedureName = selectedTemplate ?
    (selectedTemplate.code ? `${selectedTemplate.code} - ${selectedTemplate.name}` : selectedTemplate.name) :
    'Other Procedure'

  return `
    <!DOCTYPE html>
    <html>
    <head>
      <title>${selectedTemplate.name} - ${patientName}</title>
      <style>
        body { font-family: 'Times New Roman', serif; font-size: 12pt; line-height: 1.5; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { margin: 0; font-size: 18pt; color: #333; text-align: center; }
        h2 { color: #555; margin-top: 30px; }
        h3 { color: #333; border-bottom: 1px solid #ccc; padding-bottom: 5px; }
        .patient-info { background: #f9f9f9; padding: 15px; margin: 20px 0; }
        table { width: 100%; border-collapse: collapse; }
        td { padding: 8px; border: 1px solid #ccc; }
        .signature-section { margin-top: 50px; }
        .signature-line { border-bottom: 1px solid #333; height: 40px; margin-top: 10px; }
        @media print { body { margin: 20px; } }
      </style>
    </head>
    <body>
      ${populatedContent}

      <!-- Additional form data if provided -->
      ${formData.procedure_description ? `
        <div style="margin-top: 30px;">
          <h3>Additional Procedure Details</h3>
          <p style="text-align: justify;">${formData.procedure_description}</p>
        </div>
      ` : ''}

      ${formData.risks ? `
        <div style="margin-top: 30px;">
          <h3>Additional Risks Noted</h3>
          <p style="text-align: justify;">${formData.risks}</p>
        </div>
      ` : ''}

      ${formData.alternatives ? `
        <div style="margin-top: 30px;">
          <h3>Additional Alternatives Discussed</h3>
          <p style="text-align: justify;">${formData.alternatives}</p>
        </div>
      ` : ''}

      ${formData.post_care_instructions ? `
        <div style="margin-top: 30px;">
          <h3>Additional Post-Care Instructions</h3>
          <p style="text-align: justify;">${formData.post_care_instructions}</p>
        </div>
      ` : ''}
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
    isSubmitting.value = true

    // Save consent form (mock API call)
    const consentData = {
      ...formData,
      patient_id: parseInt(route.params.patientId),
      created_at: new Date().toISOString(),
      clinic_id: authStore.userClinicId
    }

    // Mock save
    console.log('Saving consent form:', consentData)

    notificationStore.showSuccess('Consent form saved successfully')
    router.push(`/patients/${route.params.patientId}`)
  } catch (error) {
    console.error('Error saving consent form:', error)
    notificationStore.showError('Failed to save consent form')
  } finally {
    isSubmitting.value = false
  }
}

// Load data on mount
onMounted(async () => {
  try {
    // Load patients
    await patientStore.fetchPatients()
    patients.value = patientStore.patients

    // Load consent templates
    await consentStore.fetchConsentTemplates()

    // Load doctors
    await userStore.fetchDoctors()

    // Patient is automatically selected from route params
    const patientId = route.params.patientId
    if (patientId) {
      await loadPatientInfo()
    }
  } catch (error) {
    console.error('Error loading data:', error)
    // Don't show error notification for now to avoid blocking the component
    // notificationStore.showError('Failed to load form data')
  }
})

// Watch consent template changes
watch(() => formData.consent_template_id, () => {
  if (formData.consent_template_id) {
    loadConsentTemplate()
  }
})

// Watch route params changes
watch(() => route.params.patientId, () => {
  loadPatientInfo()
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

.form-section {
  @apply border-b border-gray-200 pb-6 last:border-b-0 last:pb-0;
}

.form-actions {
  @apply flex flex-col sm:flex-row justify-end space-y-3 sm:space-y-0 sm:space-x-3 pt-6 border-t;
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
