<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="flex items-center mb-6">
      <div>
         <h1 class="text-2xl font-bold text-neutral-900">
           {{ isReadOnly ? 'Appointment Details' : 'Edit Appointment' }}
         </h1>
         <p class="text-neutral-600">
           {{ isReadOnly ? 'View appointment information' : 'Update appointment details' }}
         </p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      <span class="ml-3 text-neutral-600">Loading appointment...</span>
    </div>

    <!-- Form Container -->
    <div v-else class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
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
                   :readonly="isReadOnly"
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
             <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Schedule & Procedures</h3>

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
                   :readonly="isReadOnly"
                   required
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
                   :readonly="isReadOnly"
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
                 :disabled="isReadOnly"
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
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                  </svg>
                  Doctor
                  <span class="text-danger-500 ml-1">*</span>
                </label>
                 <select
                   v-model="formData.doctor_id"
                   class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                   :class="{ 'border-red-500': errors.doctor_id }"
                   :disabled="isReadOnly"
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
                 :disabled="isReadOnly"
               >
                <option value="">Select branch...</option>
                <option v-for="branch in branches" :key="branch.id" :value="branch.id">
                  {{ branch.name }}
                </option>
              </select>
             </div>

             <!-- Procedure Selection -->
             <div class="space-y-4">
               <div class="flex items-center justify-between">
                 <label class="block text-sm font-semibold text-gray-700 flex items-center">
                   <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                     <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                   </svg>
                   Procedures
                 </label>
                  <button
                    v-if="!isReadOnly"
                    type="button"
                    @click="addProcedure"
                    class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md text-primary-700 bg-primary-100 hover:bg-primary-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all duration-200"
                  >
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                    </svg>
                    Add Procedure
                  </button>
               </div>

               <!-- Selected Procedures -->
               <div v-if="selectedProcedures.length > 0" class="space-y-3">
                 <div
                   v-for="(procedure, index) in selectedProcedures"
                   :key="index"
                   class="flex items-center justify-between p-3 bg-neutral-50 rounded-lg border border-neutral-200"
                 >
                   <div class="flex-1">
                     <div class="font-medium text-neutral-900">{{ procedure.name }}</div>
                     <div class="text-sm text-neutral-600">{{ procedure.description }}</div>
                     <div class="text-xs text-neutral-500 mt-1">
                       Duration: {{ procedure.estimated_duration }} min • Cost: {{ procedure.default_cost }}
                     </div>
                   </div>
                   <button
                     type="button"
                     @click="removeProcedure(index)"
                     class="ml-3 text-red-500 hover:text-red-700 transition-colors duration-200"
                   >
                     <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                     </svg>
                   </button>
                 </div>
               </div>

               <!-- No procedures message -->
               <div v-else class="text-center py-6 text-neutral-500">
                 <svg class="w-12 h-12 mx-auto mb-3 text-neutral-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                 </svg>
                 <p class="text-sm">No procedures selected</p>
                 <p class="text-xs mt-1">Click "Add Procedure" to select dental procedures for this appointment</p>
               </div>
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
                 :readonly="isReadOnly"
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
               class="inline-flex items-center justify-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200 min-h-[44px]"
             >
               {{ isReadOnly ? 'Back' : 'Cancel' }}
             </button>

             <button
               v-if="isReadOnly"
               type="button"
               @click="editAppointment"
               class="inline-flex items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-all duration-200"
             >
               <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
               </svg>
               Edit Appointment
             </button>

             <button
               v-else
               type="submit"
               :disabled="isSubmitting || !isFormValid"
               class="inline-flex items-center justify-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 min-h-[44px]"
             >
               <template v-if="isSubmitting">
                 <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                 Updating...
               </template>
               <template v-else>
                 Update Appointment
               </template>
             </button>
           </div>
        </form>
      </div>
    </div>

    <!-- Procedure Selection Modal -->
    <div v-if="showProcedureModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeProcedureModal"></div>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4" id="modal-title">
                  Select Procedure
                </h3>

                <!-- Procedure Search -->
                <div class="mb-4">
                  <input
                    v-model="procedureSearchQuery"
                    type="text"
                    class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                    placeholder="Search procedures..."
                  />
                </div>

                <!-- Procedure List -->
                <div class="max-h-64 overflow-y-auto">
                  <div
                    v-for="procedure in filteredProcedures"
                    :key="procedure.id"
                    @click="selectProcedure(procedure)"
                    class="cursor-pointer p-3 border border-gray-200 rounded-md mb-2 hover:bg-gray-50 transition-colors duration-200"
                  >
                    <div class="font-medium text-gray-900">{{ procedure.name }}</div>
                    <div class="text-sm text-gray-600">{{ procedure.description }}</div>
                    <div class="text-xs text-gray-500 mt-1">
                      Duration: {{ procedure.estimated_duration }} min • Cost: {{ procedure.default_cost }}
                    </div>
                  </div>

                  <!-- Loading state -->
                  <div v-if="loadingProcedures" class="text-center py-4">
                    <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-600 mx-auto"></div>
                    <p class="text-sm text-gray-500 mt-2">Loading procedures...</p>
                  </div>

                  <!-- No procedures found -->
                  <div v-if="!loadingProcedures && filteredProcedures.length === 0" class="text-center py-4 text-gray-500">
                    <p class="text-sm">No procedures found</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              @click="closeProcedureModal"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Cancel
            </button>
          </div>
        </div>
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
import { useAuthStore } from '../stores/auth'
import apiService from '../services/api'

const route = useRoute()
const router = useRouter()

const patientStore = usePatientStore()
const clinicStore = useClinicStore()
const appointmentStore = useAppointmentStore()
const authStore = useAuthStore()

const isSubmitting = ref(false)
const loading = ref(false)
const errors = ref({})
const hasConflicts = ref(false)

// Check if this is read-only view mode
const isReadOnly = computed(() => {
  return route.meta.readOnly === true
})

const patientSearchQuery = ref('')
const filteredPatients = ref([])
const showPatientDropdown = ref(false)

const selectedProcedures = ref([])
const availableProcedures = ref([])
const showProcedureModal = ref(false)
const procedureSearchQuery = ref('')
const loadingProcedures = ref(false)

const formData = ref({
  patient_id: '',
  date: '',
  time: '',
  duration: '30',
  doctor_id: '',
  branch_id: '',
  notes: '',
  send_reminder: true,
  reminder_email: true,
  reminder_sms: false,
  send_confirmation: true
})

// Procedure filtering computed property

const isFormValid = computed(() => {
  const branchId = getSelectedBranchId()
  return formData.value.patient_id &&
         formData.value.date &&
         formData.value.time &&
         formData.value.duration &&
         formData.value.doctor_id &&
         (branchId > 0 || branches.value.length > 0)
})

const doctors = computed(() => {
  return clinicStore.doctors || []
})

const branches = computed(() => {
  return clinicStore.branches || []
})

// Procedure filtering
const filteredProcedures = computed(() => {
  if (!procedureSearchQuery.value) {
    return availableProcedures.value
  }
  const query = procedureSearchQuery.value.toLowerCase()
  return availableProcedures.value.filter(procedure =>
    procedure.name.toLowerCase().includes(query) ||
    procedure.description.toLowerCase().includes(query) ||
    procedure.category.toLowerCase().includes(query)
  )
})

const loadAppointment = async (appointmentId) => {
  loading.value = true
  console.log('Loading appointment with ID:', appointmentId)
  try {
    const result = await appointmentStore.fetchAppointment(appointmentId)
    console.log('Fetch appointment result:', result)
    console.log('Current appointment from store:', appointmentStore.currentAppointment)

    if (result && appointmentStore.currentAppointment) {
      const appointment = appointmentStore.currentAppointment

      // Parse the start_time to get date and time
      const startTime = new Date(appointment.start_time)
      const date = startTime.toISOString().split('T')[0]
      const time = startTime.toTimeString().slice(0, 5) // HH:MM format

      // Load appointment data into form
      Object.assign(formData.value, {
        patient_id: appointment.patient_id?.toString() || '',
        date: date,
        time: time,
        duration: appointment.duration?.toString() || '30',
        doctor_id: appointment.doctor_id?.toString() || '',
        branch_id: appointment.branch_id?.toString() || '',
        notes: appointment.pre_appointment_notes || appointment.description || '',
        send_reminder: appointment.send_reminder ?? true,
        reminder_email: appointment.reminder_email ?? true,
        reminder_sms: appointment.reminder_sms ?? false,
        send_confirmation: appointment.send_confirmation ?? true
      })

      // Set patient search query
      if (appointment.patient) {
        patientSearchQuery.value = `${appointment.patient.first_name} ${appointment.patient.last_name}`
      }

      // Load existing procedures for this appointment
      await loadAppointmentProcedures(appointmentId)

    } else {
      console.error('Failed to load appointment')
      router.push('/appointments')
    }
  } catch (error) {
    console.error('Error loading appointment:', error)
    router.push('/appointments')
  } finally {
    loading.value = false
  }
}

const loadAppointmentProcedures = async (appointmentId) => {
  try {
    // Use apiService instead of raw fetch for proper authentication
    const result = await apiService.get(`/api/appointments/${appointmentId}/procedures`)

    if (result.success) {
      selectedProcedures.value = result.data || []
    } else {
      console.error('Failed to load appointment procedures:', result.error)
      selectedProcedures.value = []
    }
  } catch (error) {
    console.error('Error loading appointment procedures:', error)
    selectedProcedures.value = []
  }
}

const searchPatients = async () => {
  if (patientSearchQuery.value.length < 1) {
    filteredPatients.value = []
    return
  }

  try {
    const result = await patientStore.searchPatients({
      query: patientSearchQuery.value,
      limit: 10
    })

    if (result.success) {
      filteredPatients.value = result.data.patients || []
    } else {
      filteredPatients.value = []
    }
  } catch (error) {
    console.error('Error searching patients:', error)
    filteredPatients.value = []
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

// Procedure methods
const loadProcedures = async () => {
  if (loadingProcedures.value) return

  loadingProcedures.value = true
  try {
    // Use apiService instead of raw fetch for proper authentication
    const result = await apiService.get('/api/procedure-templates')

    if (result.success) {
      availableProcedures.value = result.data || []
      console.log('Loaded procedures:', availableProcedures.value.length)
    } else {
      console.error('Failed to load procedures:', result.error)
      availableProcedures.value = []
    }
  } catch (error) {
    console.error('Error loading procedures:', error)
    availableProcedures.value = []
  } finally {
    loadingProcedures.value = false
  }
}

const addProcedure = () => {
  showProcedureModal.value = true
  if (availableProcedures.value.length === 0 && !loadingProcedures.value) {
    loadProcedures()
  }
}

const selectProcedure = (procedure) => {
  // Check if procedure is already selected
  const exists = selectedProcedures.value.find(p => p.id === procedure.id)
  if (!exists) {
    selectedProcedures.value.push(procedure)
  }
  closeProcedureModal()
}

const removeProcedure = (index) => {
  selectedProcedures.value.splice(index, 1)
}

const closeProcedureModal = () => {
  showProcedureModal.value = false
  procedureSearchQuery.value = ''
}

const getSelectedBranchId = () => {
  // If branch is explicitly selected, use it
  if (formData.value.branch_id) {
    const branchId = parseInt(formData.value.branch_id)
    if (branchId > 0) return branchId
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

  // If no branches available yet, try to trigger loading
  if (branches.value.length === 0) {
    console.warn('No branches loaded yet, attempting to load...')
    if (authStore.userClinicId) {
      clinicStore.fetchBranches(authStore.userClinicId)
    } else if (authStore.isClinicOwner) {
      clinicStore.fetchBranches(1)
    }
    return 0 // Return 0 to indicate loading in progress
  }

  // If no branches available after loading attempt, this is an error condition
  console.error('No valid branch available for appointment')
  return 0 // This will cause validation to fail
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
    let branchId = getSelectedBranchId()

    // If no branch is selected but branches exist, try to get one
    if (branchId <= 0 && branches.value.length > 0) {
      const mainBranch = clinicStore.getMainBranch
      if (mainBranch) {
        branchId = mainBranch.id
        formData.value.branch_id = branchId.toString()
      } else if (branches.value.length === 1) {
        branchId = branches.value[0].id
        formData.value.branch_id = branchId.toString()
      }
    }

    const patientId = parseInt(formData.value.patient_id)
    const doctorId = parseInt(formData.value.doctor_id)
    const appointmentId = route.params.id

    if (isNaN(patientId) || patientId <= 0) {
      alert('Please select a valid patient')
      return
    }

    if (isNaN(doctorId) || doctorId <= 0) {
      alert('Please select a valid doctor')
      return
    }

    if (branchId <= 0) {
      alert('Please wait for branches to load or select a valid branch')
      return
    }

    if (!formData.value.date || !formData.value.time) {
      alert('Please select a valid date and time')
      return
    }

    // Generate title based on selected procedures or default
    const appointmentTitle = selectedProcedures.value.length > 0
      ? selectedProcedures.value[0].name
      : 'Dental Appointment'

    const appointmentData = {
      patient_id: Number(patientId),
      doctor_id: Number(doctorId),
      branch_id: Number(branchId),
      start_time: `${formData.value.date}T${formData.value.time}:00Z`,
      end_time: calculateEndTime(),
      duration: Number(formData.value.duration) || 30,
      pre_appointment_notes: formData.value.notes || '',
      title: appointmentTitle,
      description: formData.value.notes || '',
      status: 'scheduled'
    }

    console.log('Updating appointment:', appointmentData)

    const result = await appointmentStore.updateAppointment(appointmentId, appointmentData)

    if (result.success) {
      // Update procedures for this appointment
      await updateAppointmentProcedures(appointmentId)

      // Navigate back to calendar
      router.push('/appointments')
    } else {
      console.error('Failed to update appointment:', result.error)
      alert(`Failed to update appointment: ${result.error}`)
    }
  } catch (error) {
    console.error('Error updating appointment:', error)
  } finally {
    isSubmitting.value = false
  }
}

const updateAppointmentProcedures = async (appointmentId) => {
  try {
    // First, get existing procedures
    const existingResult = await apiService.get(`/api/appointments/${appointmentId}/procedures`)
    const existingProcedures = existingResult.success ? existingResult.data : []

    // Remove procedures that are no longer selected
    for (const existingProc of existingProcedures) {
      const stillSelected = selectedProcedures.value.find(p => p.id === existingProc.procedure_template_id)
      if (!stillSelected) {
        await apiService.delete(`/api/appointment-procedures/${existingProc.id}`)
      }
    }

    // Add new procedures
    for (const procedure of selectedProcedures.value) {
      const alreadyExists = existingProcedures.find(p => p.procedure_template_id === procedure.id)
      if (!alreadyExists) {
        const procedureResult = await apiService.post(`/api/appointments/${appointmentId}/procedures`, {
          procedure_template_id: procedure.id,
          cost: procedure.default_cost
        })

        if (!procedureResult.success) {
          console.error('Failed to add procedure:', procedure.name, procedureResult.error)
        }
      }
    }
  } catch (error) {
    console.error('Error updating appointment procedures:', error)
  }
}

const calculateEndTime = () => {
  const startDateTime = new Date(`${formData.value.date}T${formData.value.time}:00Z`)
  const endDateTime = new Date(startDateTime.getTime() + parseInt(formData.value.duration) * 60000)
  return endDateTime.toISOString()
}

const editAppointment = () => {
  router.push(`/appointments/${route.params.id}/edit`)
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showPatientDropdown.value = false
  }
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside)

  // Load initial data
  if (!doctors.value.length) {
    console.log('Fetching doctors...')
    clinicStore.fetchDoctors()
  }

   if (!branches.value.length) {
     console.log('Current user data:', authStore.user)
     console.log('User clinic ID:', authStore.userClinicId)
     console.log('Is super admin:', authStore.isSuperAdmin)
     console.log('Is clinic owner:', authStore.isClinicOwner)

     if (authStore.userClinicId) {
       console.log('Fetching branches for clinic:', authStore.userClinicId)
       clinicStore.fetchBranches(authStore.userClinicId)
     } else if (authStore.isSuperAdmin) {
       console.warn('Super admin detected, skipping branch fetch')
     } else if (authStore.isClinicOwner) {
       // Clinic owner without clinic_id - try to fetch clinic 1
       console.log('Clinic owner without clinic_id, trying clinic 1')
       clinicStore.fetchBranches(1)
     } else {
       console.warn('No clinic ID available for non-super-admin user')
     }
   }

  // Load the appointment data
  console.log('Route params:', route.params)
  console.log('Appointment ID from route:', route.params.id)
  if (route.params.id) {
    await loadAppointment(route.params.id)
  } else {
    console.error('No appointment ID found in route params')
    router.push('/appointments')
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