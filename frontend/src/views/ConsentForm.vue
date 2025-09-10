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
             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-group">
                 <label class="form-label">Patient</label>
                 <div class="form-input bg-gray-50 text-gray-700 cursor-not-allowed">
                   {{ selectedPatient?.first_name }} {{ selectedPatient?.last_name }}
                 </div>
                 <p class="text-sm text-gray-500 mt-1">Patient is automatically selected from the current context</p>
               </div>

               <div class="form-group">
                 <label class="form-label">Date</label>
                 <input
                   v-model="formData.signature_date"
                   type="date"
                   class="form-input bg-gray-50 text-gray-700 cursor-not-allowed"
                   readonly
                 />
                 <p class="text-sm text-gray-500 mt-1">Date is set to today automatically</p>
               </div>
             </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-group">
                 <label class="form-label">Treating Doctor *</label>
                 <div class="relative">
                   <input
                     v-model="doctorSearchQuery"
                     @input="searchDoctors"
                     @focus="showDoctorSuggestions = true"
                     @blur="hideDoctorSuggestions"
                     type="text"
                     class="form-input"
                     :class="{ 'border-red-500': errors.doctor_id }"
                     placeholder="Type doctor name..."
                     autocomplete="off"
                   />
                   <div v-if="showDoctorSuggestions && filteredDoctors.length > 0" 
                        class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto">
                     <div
                       v-for="doctor in filteredDoctors"
                       :key="doctor.id"
                       @mousedown="selectDoctor(doctor)"
                       class="px-4 py-2 cursor-pointer hover:bg-gray-100 transition-colors"
                     >
                       <div class="font-medium">Dr. {{ doctor.first_name }} {{ doctor.last_name }}</div>
                       <div v-if="doctor.specialization" class="text-sm text-gray-600">{{ doctor.specialization }}</div>
                     </div>
                   </div>
                 </div>
                 <div v-if="errors.doctor_id" class="error-message">{{ errors.doctor_id }}</div>
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
                <div class="max-w-4xl mx-auto">
                  <div class="consent-content prose prose-sm max-w-none">
                     <div v-html="consentContent" class="whitespace-pre-wrap font-sans text-gray-800 leading-relaxed"></div>

                     <!-- Signature Section -->
                     <div class="signature-section mt-5 p-5 bg-slate-50 rounded-lg border border-slate-300">
                       <h3 class="mb-5 text-gray-800 text-lg font-semibold text-center border-b-2 border-blue-500 pb-2.5">Signatures Required</h3>

                       <div class="flex flex-col md:flex-row gap-8 mb-4">
                         <div class="flex-1 bg-white p-4 rounded border border-gray-300 shadow-sm">
                           <div class="mb-3 font-semibold text-gray-700 text-sm">
                             Patient Signature:
                           </div>
                           <div class="signature-pad-container mb-2">
                             <canvas
                               ref="patientCanvasRef"
                               width="400"
                               height="150"
                               class="w-full h-36 border-2 border-gray-700 rounded bg-gray-50 cursor-crosshair"
                             ></canvas>
                             <div class="flex gap-2 mt-2">
                               <button type="button" @click="clearPatientSignature" class="text-xs px-2 py-1 bg-red-500 text-white border-0 rounded cursor-pointer hover:bg-red-600">Clear</button>
                               <span class="text-xs text-gray-500">Draw your signature above</span>
                             </div>
                             <div v-if="errors.patient_signature_data" class="mt-2 p-2 bg-red-50 border border-red-200 rounded-md">
                               <div class="flex items-center">
                                 <font-awesome-icon icon="fa-solid fa-exclamation-triangle" class="w-4 h-4 text-red-500 mr-2" />
                                 <p class="text-sm text-red-700 mb-0">{{ errors.patient_signature_data }}</p>
                               </div>
                             </div>
                           </div>
                           <div class="text-xs text-gray-500 text-center">
                             Date: {{ formData.signature_date || new Date().toISOString().split('T')[0] }}
                           </div>
                         </div>

                         <div class="flex-1 bg-white p-4 rounded border border-gray-300 shadow-sm">
                           <div class="mb-3 font-semibold text-gray-700 text-sm">
                             Witness Signature:
                           </div>
                           <div class="signature-pad-container mb-2">
                             <canvas
                               ref="witnessCanvasRef"
                               width="400"
                               height="150"
                               class="w-full h-36 border-2 border-gray-700 rounded bg-gray-50 cursor-crosshair"
                             ></canvas>
                             <div class="flex gap-2 mt-2">
                               <button type="button" @click="clearWitnessSignature" class="text-xs px-2 py-1 bg-red-500 text-white border-0 rounded cursor-pointer hover:bg-red-600">Clear</button>
                               <span class="text-xs text-gray-500">Draw witness signature above</span>
                             </div>
                           </div>
                           <div class="text-xs text-gray-500 text-center">
                             Date: {{ formData.signature_date || new Date().toISOString().split('T')[0] }}
                           </div>
                           <div v-if="formData.witness_name" class="mt-2 text-xs text-gray-500 text-center italic">
                             Witness: {{ formData.witness_name }}
                           </div>
                         </div>
                       </div>
                     </div>

                     <!-- Footer -->
                     <div class="mt-5 text-center text-xs text-gray-500 pt-2.5 font-normal">
                       <p class="m-0 mb-1">This consent form has been electronically generated and is legally binding.</p>
                       <p class="m-0">Please retain a copy for your records.</p>
                     </div>
                  </div>
                </div>
              </div>
          </div>

         <!-- Treatment Details (Collapsible) -->
         <div class="form-section">
           <div class="flex items-center justify-between mb-4">
             <h3 class="text-lg font-semibold text-gray-900">Treatment Details</h3>
             <button
               type="button"
               @click="showTreatmentDetails = !showTreatmentDetails"
               class="flex items-center space-x-1 text-sm text-blue-600 hover:text-blue-800 transition-colors"
             >
               <span>{{ showTreatmentDetails ? 'Hide' : 'Show' }} Details</span>
               <font-awesome-icon 
                 :icon="['fas', 'chevron-down']" 
                 :class="{ 'rotate-180': showTreatmentDetails }"
                 class="transition-transform duration-200"
               />
             </button>
           </div>

           <div 
             v-show="showTreatmentDetails" 
             class="space-y-4 transition-all duration-300"
           >
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
              {{ isEditing ? 'Updating...' : 'Saving...' }}
            </div>
            <span v-else>
              {{ isEditing ? 'Update Consent Form' : 'Save Consent Form' }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
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
const showTreatmentDetails = ref(false)

// Form data
const formData = reactive({
  consent_template_id: '',
  doctor_id: '',
  procedure_description: '',
  risks: '',
  alternatives: '',
  post_care_instructions: '',
  understands_treatment: false,
  understands_risks: false,
  consents_to_treatment: false,
  had_opportunity_to_ask: false,
  patient_signature_data: '',
  signature_date: new Date().toISOString().split('T')[0],
  witness_name: '',
  witness_signature: '',
  witness_signature_data: ''
})

// Signature pad variables
let patientCanvas = null
let patientCtx = null
let witnessCanvas = null
let witnessCtx = null
let isDrawing = false
let lastPoint = null
let lastTime = null

// Vue refs
const patientCanvasRef = ref(null)
const witnessCanvasRef = ref(null)

// Auto-population state
const autoPopulateFields = ref(true)

// Data from stores
const patients = ref([])
const patientAppointments = ref([])
const doctors = ref([])
const consentTemplate = ref('')
const treatmentRisks = ref([])
const alternativeTreatments = ref([])

// Doctor search functionality
const doctorSearchQuery = ref('')
const showDoctorSuggestions = ref(false)
const filteredDoctors = ref([])
const selectedDoctorInfo = ref(null)

// Computed properties
const selectedPatient = computed(() => {
  const patientId = route.params.patientId
  return patients.value.find(p => p.id === parseInt(patientId))
})

const selectedConsentTemplate = computed(() => {
  if (!formData.consent_template_id) return null
  return consentStore.getConsentTemplateById(parseInt(formData.consent_template_id))
})

const selectedDoctor = computed(() => {
  if (!formData.doctor_id) return null
  return doctors.value.find(d => d.id === parseInt(formData.doctor_id)) || selectedDoctorInfo.value
})

const consentContent = computed(() => {
  if (!selectedConsentTemplate.value || !selectedPatient.value) return ''

  const patient = selectedPatient.value
  const template = selectedConsentTemplate.value
  const currentDate = new Date().toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  // Replace common placeholders
  const doctorName = selectedDoctor.value ? `${selectedDoctor.value.first_name} ${selectedDoctor.value.last_name}` : 'Doctor'
  
  let content = template.content
    .replace(/\[PATIENT_NAME\]/g, `${patient.first_name} ${patient.last_name}`)
    .replace(/\[PATIENT_FIRST_NAME\]/g, patient.first_name)
    .replace(/\[PATIENT_LAST_NAME\]/g, patient.last_name)
    .replace(/\[CURRENT_DATE\]/g, currentDate)
    .replace(/\[TODAY\]/g, currentDate)
    .replace(/\[DATE\]/g, currentDate)
    .replace(/\[DOCTOR_NAME\]/g, doctorName)
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

  // Apply formatting to the template content
  content = formatConsentContent(content)

  return content
})

// Methods
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const formatConsentContent = (content) => {
  if (!content) return ''

  // Split content into lines and format each section
  const lines = content.split('\n')
  const formattedLines = lines.map(line => {
    // Format main headings (all caps)
    if (line.match(/^[A-Z\s]+$/)) {
      return `<h2 class="text-xl font-bold text-gray-900 mt-6 mb-3 border-b border-gray-300 pb-2">${line}</h2>`
    }

    // Format section headings (UPPERCASE)
    if (line.match(/^[A-Z\s]+:$/)) {
      return `<h3 class="text-lg font-semibold text-gray-800 mt-4 mb-2">${line.replace(':', '')}</h3>`
    }

    // Format bullet points
    if (line.trim().startsWith('- ')) {
      return `<li class="ml-4 mb-1">${line.trim().substring(2)}</li>`
    }

    // Format patient info lines (contain colons)
    if (line.includes(':') && !line.includes('http') && line.split(':').length === 2) {
      const [label, value] = line.split(':')
      return `<div class="mb-2"><strong class="text-gray-700">${label.trim()}:</strong> ${value.trim()}</div>`
    }

    // Regular paragraphs
    if (line.trim()) {
      return `<p class="mb-3 text-gray-700 leading-relaxed">${line.trim()}</p>`
    }

    // Empty lines
    return '<br>'
  })

  return formattedLines.join('')
}

// Signature pad methods
const initSignaturePads = () => {
  // Initialize patient signature canvas
  patientCanvas = patientCanvasRef.value
  if (patientCanvas) {
    setupCanvas(patientCanvas, 'patient')
  }

  // Initialize witness signature canvas
  witnessCanvas = witnessCanvasRef.value
  if (witnessCanvas) {
    setupCanvas(witnessCanvas, 'witness')
  }
}

const setupCanvas = (canvas, type) => {
  if (!canvas) {
    console.error(`Canvas element for ${type} not found`)
    return
  }

  const rect = canvas.getBoundingClientRect()
  
  if (rect.width === 0 || rect.height === 0) {
    console.warn(`Canvas ${type} has zero dimensions: ${rect.width}x${rect.height}`)
    return
  }
  
  // Set canvas dimensions to match display size
  canvas.width = rect.width
  canvas.height = rect.height

  const ctx = canvas.getContext('2d')
  
  // Set drawing styles
  ctx.strokeStyle = '#374151'
  ctx.lineWidth = 1.0
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.globalCompositeOperation = 'source-over'

  // Store context reference
  if (type === 'patient') {
    patientCtx = ctx
  } else {
    witnessCtx = ctx
  }


  // Remove existing event listeners to prevent duplicates
  if (canvas._mousedownHandler) {
    canvas.removeEventListener('mousedown', canvas._mousedownHandler)
    canvas.removeEventListener('mousemove', canvas._mousemoveHandler)
    canvas.removeEventListener('mouseup', canvas._mouseupHandler)
    canvas.removeEventListener('mouseout', canvas._mouseoutHandler)
    canvas.removeEventListener('touchstart', canvas._touchstartHandler)
    canvas.removeEventListener('touchmove', canvas._touchmoveHandler)
    canvas.removeEventListener('touchend', canvas._touchendHandler)
  }

  // Create bound handlers and store references for removal
  canvas._mousedownHandler = (e) => startDrawing(type, e)
  canvas._mousemoveHandler = (e) => draw(type, e)
  canvas._mouseupHandler = stopDrawing
  canvas._mouseoutHandler = stopDrawing
  canvas._touchstartHandler = (e) => handleTouchStart(type, e)
  canvas._touchmoveHandler = (e) => handleTouchMove(type, e)
  canvas._touchendHandler = stopDrawing

  // Mouse events
  canvas.addEventListener('mousedown', canvas._mousedownHandler)
  canvas.addEventListener('mousemove', canvas._mousemoveHandler)
  canvas.addEventListener('mouseup', canvas._mouseupHandler)
  canvas.addEventListener('mouseout', canvas._mouseoutHandler)

  // Touch events
  canvas.addEventListener('touchstart', canvas._touchstartHandler, { passive: false })
  canvas.addEventListener('touchmove', canvas._touchmoveHandler, { passive: false })
  canvas.addEventListener('touchend', canvas._touchendHandler)
}

const startDrawing = (type, e) => {
  isDrawing = true
  const canvas = type === 'patient' ? patientCanvas : witnessCanvas
  const ctx = type === 'patient' ? patientCtx : witnessCtx

  if (!canvas || !ctx) return

  const rect = canvas.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  lastPoint = { x, y }
  lastTime = Date.now()
  
  ctx.beginPath()
  ctx.moveTo(x, y)
}

const draw = (type, e) => {
  if (!isDrawing || !lastPoint) return

  const canvas = type === 'patient' ? patientCanvas : witnessCanvas
  const ctx = type === 'patient' ? patientCtx : witnessCtx

  if (!canvas || !ctx) return

  // Clear signature validation errors when user starts drawing
  if (type === 'patient' && errors.value.patient_signature_data) {
    delete errors.value.patient_signature_data
  }

  const rect = canvas.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  // Calculate distance and time to determine stroke width
  const currentTime = Date.now()
  const distance = Math.sqrt(Math.pow(x - lastPoint.x, 2) + Math.pow(y - lastPoint.y, 2))
  const timeDelta = currentTime - lastTime
  
  // Calculate velocity (distance per time)
  const velocity = timeDelta > 0 ? distance / timeDelta : 0
  
  // Variable stroke width based on velocity (slower = thicker, faster = thinner)
  // Clamp between 0.5 and 3.0 for signature-like appearance with smaller dots
  const baseWidth = 1.5
  const maxWidth = 3.0
  const minWidth = 0.5
  const velocityFactor = Math.max(0.1, Math.min(1, 1 / (velocity * 50 + 1)))
  const strokeWidth = minWidth + (maxWidth - minWidth) * velocityFactor

  // Only draw if there's meaningful movement to reduce unwanted dots
  if (distance > 0.5) {
    ctx.lineWidth = strokeWidth
    ctx.lineTo(x, y)
    ctx.stroke()
  }

  // Update last point and time
  lastPoint = { x, y }
  lastTime = currentTime
}

const stopDrawing = () => {
  if (!isDrawing) return
  
  isDrawing = false
  lastPoint = null
  lastTime = null

  // Save signature data
  if (patientCanvas) {
    formData.patient_signature_data = patientCanvas.toDataURL()
  }
  if (witnessCanvas) {
    formData.witness_signature_data = witnessCanvas.toDataURL()
  }
}

const handleTouchStart = (type, e) => {
  e.preventDefault()
  const touch = e.touches[0]
  const mouseEvent = new MouseEvent('mousedown', {
    clientX: touch.clientX,
    clientY: touch.clientY
  })
  startDrawing(type, mouseEvent)
}

const handleTouchMove = (type, e) => {
  e.preventDefault()
  if (!isDrawing) return

  const touch = e.touches[0]
  const mouseEvent = new MouseEvent('mousemove', {
    clientX: touch.clientX,
    clientY: touch.clientY
  })
  draw(type, mouseEvent)
}

const clearPatientSignature = () => {
  if (patientCtx && patientCanvas) {
    patientCtx.clearRect(0, 0, patientCanvas.width, patientCanvas.height)
    formData.patient_signature_data = ''
    lastPoint = null
    lastTime = null
    // Clear signature validation errors
    if (errors.value.patient_signature_data) {
      delete errors.value.patient_signature_data
    }
  }
}

const clearWitnessSignature = () => {
  if (witnessCtx && witnessCanvas) {
    witnessCtx.clearRect(0, 0, witnessCanvas.width, witnessCanvas.height)
    formData.witness_signature_data = ''
    lastPoint = null
    lastTime = null
  }
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

  if (!formData.doctor_id) {
    errors.value.doctor_id = 'Please select a treating doctor from the suggestions'
  }

  if (!formData.understands_treatment || !formData.understands_risks ||
      !formData.consents_to_treatment || !formData.had_opportunity_to_ask) {
    errors.value.agreement = 'All patient agreement checkboxes must be checked'
  }

  // Validate that signature pads have been signed (not empty canvas)
  if (!formData.patient_signature_data || formData.patient_signature_data === '') {
    errors.value.patient_signature_data = 'Patient signature is required - please draw your signature in the signature area'
  } else {
    // Check if it's just an empty/transparent canvas
    const isEmptySignature = isCanvasEmptySync(formData.patient_signature_data)
    if (isEmptySignature) {
      errors.value.patient_signature_data = 'Please draw your signature in the patient signature area - the canvas appears to be empty'
    }
  }

  return Object.keys(errors.value).length === 0
}

// Helper function to check if canvas signature is empty
const isCanvasEmpty = (signatureData) => {
  if (!signatureData || signatureData === '') return true
  
  // Check if it's a data URL
  if (!signatureData.startsWith('data:image/png;base64,')) return true
  
  // Extract base64 part
  const base64Data = signatureData.split(',')[1]
  if (!base64Data || base64Data.length < 100) return true // Very short means likely empty
  
  // More thorough check: create a temporary canvas to analyze pixel data
  try {
    // Create a temporary image and canvas to check pixel data
    const img = new Image()
    const tempCanvas = document.createElement('canvas')
    const tempCtx = tempCanvas.getContext('2d')
    
    return new Promise((resolve) => {
      img.onload = () => {
        tempCanvas.width = img.width
        tempCanvas.height = img.height
        tempCtx.drawImage(img, 0, 0)
        
        // Get image data
        const imageData = tempCtx.getImageData(0, 0, tempCanvas.width, tempCanvas.height)
        const data = imageData.data
        
        // Check if any pixel has opacity > 0 (non-transparent)
        for (let i = 3; i < data.length; i += 4) { // Check alpha channel (every 4th value)
          if (data[i] > 0) {
            resolve(false) // Found non-transparent pixel, signature exists
            return
          }
        }
        resolve(true) // All pixels are transparent, signature is empty
      }
      
      img.onerror = () => resolve(true) // If image fails to load, consider it empty
      img.src = signatureData
    })
  } catch (error) {
    console.error('Error checking canvas signature:', error)
    return false // If error occurs, assume signature exists to be safe
  }
}

// Synchronous fallback for immediate validation
const isCanvasEmptySync = (signatureData) => {
  if (!signatureData || signatureData === '') return true
  if (!signatureData.startsWith('data:image/png;base64,')) return true
  
  const base64Data = signatureData.split(',')[1]
  if (!base64Data || base64Data.length < 100) return true
  
  // Quick check: compare with known empty canvas signature
  // This is the base64 of a typical empty transparent canvas
  const emptyCanvasPatterns = [
    'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==', // 1x1 transparent
    'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=' // Another common empty pattern
  ]
  
  // If it matches known empty patterns, it's empty
  for (const pattern of emptyCanvasPatterns) {
    if (base64Data === pattern) return true
  }
  
  return false // Assume it has content if it doesn't match empty patterns
}

// Doctor search methods
const searchDoctors = () => {
  if (!doctorSearchQuery.value.trim()) {
    filteredDoctors.value = []
    return
  }
  
  const query = doctorSearchQuery.value.toLowerCase()
  filteredDoctors.value = doctors.value.filter(doctor => {
    const fullName = `${doctor.first_name} ${doctor.last_name}`.toLowerCase()
    const firstNameMatch = doctor.first_name.toLowerCase().includes(query)
    const lastNameMatch = doctor.last_name.toLowerCase().includes(query)
    const fullNameMatch = fullName.includes(query)
    const specializationMatch = doctor.specialization?.toLowerCase().includes(query) || false
    
    return firstNameMatch || lastNameMatch || fullNameMatch || specializationMatch
  })
}

const selectDoctor = (doctor) => {
  formData.doctor_id = doctor.id
  selectedDoctorInfo.value = doctor
  doctorSearchQuery.value = `Dr. ${doctor.first_name} ${doctor.last_name}`
  showDoctorSuggestions.value = false
  // Clear any doctor validation errors
  if (errors.value.doctor_id) {
    delete errors.value.doctor_id
  }
}

const hideDoctorSuggestions = () => {
  setTimeout(() => {
    showDoctorSuggestions.value = false
  }, 200) // Delay to allow click events on suggestions
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

  // Add signature areas for PDF generation
  const signatureAreas = `
    <div class="signature-section" style="margin-top: 40px; padding-top: 30px; page-break-inside: avoid;">
      <h3 style="margin-bottom: 25px; color: #374151; font-size: 16px; font-weight: 600; border-bottom: 1px solid #d1d5db; padding-bottom: 8px;">Signatures</h3>

      <div style="display: flex; justify-content: space-between; gap: 40px; margin-bottom: 20px;">
        <div style="flex: 1;">
          <div style="margin-bottom: 8px; font-weight: 500; color: #374151;">
            Patient Signature:
          </div>
          <div class="signature-line" style="border-bottom: 1px solid #6b7280; height: 50px; margin-top: 8px; display: flex; align-items: flex-end; padding-bottom: 4px; background: #f9fafb;">
            <span style="font-size: 11px; color: #6b7280; font-family: 'Courier New', monospace;">${formData.patient_signature || ''}</span>
          </div>
          <div style="margin-top: 4px; font-size: 11px; color: #6b7280;">
            Date: ${formData.signature_date || currentDate}
          </div>
        </div>

        <div style="flex: 1;">
          <div style="margin-bottom: 8px; font-weight: 500; color: #374151;">
            Witness Signature:
          </div>
          <div class="signature-line" style="border-bottom: 1px solid #6b7280; height: 50px; margin-top: 8px; display: flex; align-items: flex-end; padding-bottom: 4px; background: #f9fafb;">
            <span style="font-size: 11px; color: #6b7280; font-family: 'Courier New', monospace;">${formData.witness_signature || ''}</span>
          </div>
          <div style="margin-top: 4px; font-size: 11px; color: #6b7280;">
            Date: ${formData.signature_date || currentDate}
          </div>
          ${formData.witness_name ? `<div style="margin-top: 4px; font-size: 11px; color: #6b7280;">Witness: ${formData.witness_name}</div>` : ''}
        </div>
      </div>
    </div>
  `

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

      <!-- Signature areas -->
      ${signatureAreas}
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

    // Generate title and description based on template and patient
    const template = selectedConsentTemplate.value
    const patient = selectedPatient.value
    const title = template ? `${template.name} - ${patient.first_name} ${patient.last_name}` : 'Consent Form'
    const description = `Consent form for ${patient.first_name} ${patient.last_name} created on ${new Date().toLocaleDateString()}`

    // Clean up signature data - set to null if empty/transparent
    const patientSignatureClean = (formData.patient_signature_data && !isCanvasEmptySync(formData.patient_signature_data)) 
      ? formData.patient_signature_data 
      : null
      
    const witnessSignatureClean = (formData.witness_signature_data && !isCanvasEmptySync(formData.witness_signature_data)) 
      ? formData.witness_signature_data 
      : null

    // Prepare consent form data - map frontend field names to backend field names
    const consentData = {
      ...formData,
      title,
      description,
      patient_id: parseInt(route.params.patientId),
      created_at: new Date().toISOString(),
      clinic_id: authStore.userClinicId,
      // Map signature data to backend field names - use cleaned versions
      patient_signature: patientSignatureClean,
      witness_signature: witnessSignatureClean,
      // Set signature timestamps only if signatures exist and are not empty
      patient_signed_at: patientSignatureClean ? new Date().toISOString() : null,
      witness_signed_at: witnessSignatureClean ? new Date().toISOString() : null,
      // Set created by user
      created_by_id: authStore.user?.id
    }

    // Save consent form using the store
    await consentStore.createConsentForm(consentData)

    notificationStore.showSuccess('Consent form saved successfully')
    router.push(`/patients/${route.params.patientId}`)
  } catch (error) {
    console.error('Error saving consent form:', error)
    notificationStore.showError('Failed to save consent form')
  } finally {
    isSubmitting.value = false
  }
}

// Handle window resize to reinitialize canvases
const handleResize = () => {
  // Debounce resize events
  clearTimeout(window.canvasResizeTimeout)
  window.canvasResizeTimeout = setTimeout(() => {
    if (patientCanvasRef.value || witnessCanvasRef.value) {
      // Store current signatures before reinitializing
      const patientData = formData.patient_signature_data
      const witnessData = formData.witness_signature_data
      
      // Reinitialize canvases
      initSignaturePads()
      
      // Restore signatures if they existed
      if (patientData && patientCanvas && patientCtx) {
        const img = new Image()
        img.onload = () => {
          patientCtx.drawImage(img, 0, 0)
        }
        img.src = patientData
      }
      
      if (witnessData && witnessCanvas && witnessCtx) {
        const img = new Image()
        img.onload = () => {
          witnessCtx.drawImage(img, 0, 0)
        }
        img.src = witnessData
      }
    }
  }, 250)
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
    doctors.value = userStore.doctors || []

    // Patient is automatically selected from route params
    const patientId = route.params.patientId
    if (patientId) {
      await loadPatientInfo()
    }

    // Initialize signature pads after DOM is ready
    await nextTick()
    setTimeout(() => {
      initSignaturePads()
    }, 500)
    
    // Add resize listener
    window.addEventListener('resize', handleResize)
  } catch (error) {
    console.error('Error loading data:', error)
    // Don't show error notification for now to avoid blocking the component
    // notificationStore.showError('Failed to load form data')
  }
})

// Watch consent template changes
watch(() => formData.consent_template_id, async () => {
  if (formData.consent_template_id) {
    loadConsentTemplate()
    // Initialize signature pads after template loads and DOM updates
    await nextTick()
    setTimeout(() => {
      initSignaturePads()
    }, 200)
  }
})

// Watch route params changes
watch(() => route.params.patientId, () => {
  loadPatientInfo()
})

// Cleanup on unmount
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  clearTimeout(window.canvasResizeTimeout)
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

/* Chevron rotation animation */
.rotate-180 {
  transform: rotate(180deg);
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
