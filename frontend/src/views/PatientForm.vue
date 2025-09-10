<template>
  <!-- Loading State -->
  <BaseTransition name="fade">
    <div v-if="loading" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
      <BaseLoading type="spinner" size="large" color="primary" text="Loading patient data..." />
      <p class="text-neutral-600 mt-4">Please wait while we fetch the patient information.</p>
    </div>
  </BaseTransition>

  <!-- Error State -->
  <BaseTransition name="fade">
    <div v-if="!loading && error" class="bg-white rounded-2xl p-12 shadow-lg border border-neutral-100 text-center">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mb-4">
        <svg class="w-8 h-8 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Error loading patient</h3>
      <p class="text-danger-600 mb-4">{{ error }}</p>
      <button
        @click="loadPatient"
        class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-danger-600 hover:bg-danger-700 focus:outline-none focus:ring-2 focus:ring-danger-500 focus:ring-offset-2 transition-all duration-200"
      >
        Try Again
      </button>
    </div>
  </BaseTransition>

  <!-- Patient Form -->
  <BaseTransition name="slide-up">
    <div v-if="!loading && !error" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div class="p-6 sm:p-8">
        <!-- Header with Actions -->
        <div v-if="isViewMode" class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 pb-4 border-b border-neutral-200">
          <div>
            <h2 class="text-xl font-semibold text-neutral-900">Patient Details</h2>
            <p class="text-sm text-neutral-600 mt-1">View patient information and manage records</p>
          </div>
          <div class="flex items-center space-x-3 mt-4 sm:mt-0">
            <router-link
              :to="`/patients/${route.params.id}/edit`"
              class="inline-flex items-center px-4 py-2 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all duration-200"
            >
              <font-awesome-icon icon="fa-solid fa-edit" class="w-4 h-4 mr-2" />
              Edit Patient
            </router-link>
            <router-link
              :to="`/patients/${route.params.id}/consent/new`"
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 transition-all duration-200"
            >
              <font-awesome-icon icon="fa-solid fa-file-signature" class="w-4 h-4 mr-2" />
              New Consent Form
            </router-link>
          </div>
        </div>

        <form v-if="!isViewMode" @submit.prevent="handleSubmit" class="space-y-8">
            <!-- Basic Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Basic Information</h3>
              
              <!-- Avatar Upload -->
              <div v-if="!isViewMode" class="flex justify-center">
                <AvatarUpload
                  :user="form"
                  entity-type="patient"
                  :entity-id="form.id || 0"
                  @avatar-updated="handleAvatarUpdated"
                />
              </div>

              <!-- Name Fields -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="first_name" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                    </svg>
                    First Name
                    <span class="text-danger-500 ml-1">*</span>
                  </label>
                  <input
                    id="first_name"
                    v-model="form.first_name"
                    type="text"
                    required
                    :readonly="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.first_name, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                    placeholder="Enter first name"
                  />
                  <span v-if="errors.first_name" class="text-red-500 text-sm mt-1">{{ errors.first_name }}</span>
                </div>

                <div class="space-y-2">
                  <label for="last_name" class="block text-sm font-semibold text-gray-700">
                    Last Name
                    <span class="text-danger-500 ml-1">*</span>
                  </label>
                  <input
                    id="last_name"
                    v-model="form.last_name"
                    type="text"
                    required
                    :readonly="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.last_name, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                    placeholder="Enter last name"
                  />
                  <span v-if="errors.last_name" class="text-red-500 text-sm mt-1">{{ errors.last_name }}</span>
                </div>
              </div>

              <!-- Date of Birth and Gender -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="date_of_birth" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    Date of Birth
                    <span class="text-danger-500 ml-1">*</span>
                  </label>
                  <input
                    id="date_of_birth"
                    v-model="form.date_of_birth"
                    type="date"
                    required
                    :readonly="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.date_of_birth, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                  />
                  <span v-if="errors.date_of_birth" class="text-red-500 text-sm mt-1">{{ errors.date_of_birth }}</span>
                </div>

                <div class="space-y-2">
                  <label for="gender" class="block text-sm font-semibold text-gray-700">
                    Gender
                    <span class="text-danger-500 ml-1">*</span>
                  </label>
                  <select
                    id="gender"
                    v-model="form.gender"
                    required
                    :disabled="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.gender, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                  >
                    <option value="">Select gender</option>
                    <option value="male">Male</option>
                    <option value="female">Female</option>
                    <option value="other">Other</option>
                    <option value="prefer_not_to_say">Prefer not to say</option>
                  </select>
                  <span v-if="errors.gender" class="text-red-500 text-sm mt-1">{{ errors.gender }}</span>
                </div>
              </div>
            </div>

            <!-- Contact Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Contact Information</h3>

              <!-- Phone and Email -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="phone" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
                    </svg>
                    Phone Number
                    <span class="text-danger-500 ml-1">*</span>
                  </label>
                  <input
                    id="phone"
                    v-model="form.phone"
                    type="tel"
                    required
                    :readonly="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.phone, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                    placeholder="Enter phone number"
                  />
                  <span v-if="errors.phone" class="text-red-500 text-sm mt-1">{{ errors.phone }}</span>
                </div>

                <div class="space-y-2">
                  <label for="email" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"/>
                    </svg>
                    Email Address
                  </label>
                  <input
                    id="email"
                    v-model="form.email"
                    type="email"
                    :readonly="isViewMode"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.email, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                    placeholder="your@email.com"
                  />
                  <span v-if="errors.email" class="text-red-500 text-sm mt-1">{{ errors.email }}</span>
                </div>
              </div>

              <!-- Address -->
              <div class="space-y-2">
                <label for="address" class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                  </svg>
                  Address
                </label>
                <textarea
                  id="address"
                  v-model="form.address"
                  rows="3"
                  :readonly="isViewMode"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                  :class="{ 'border-red-500': errors.address, 'bg-gray-100 cursor-not-allowed': isViewMode }"
                  placeholder="Enter full address"
                ></textarea>
                <span v-if="errors.address" class="text-red-500 text-sm mt-1">{{ errors.address }}</span>
              </div>
            </div>

            <!-- Emergency Contact -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Emergency Contact</h3>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="emergency_contact_name" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                    </svg>
                    Contact Name
                  </label>
                  <input
                    id="emergency_contact_name"
                    v-model="form.emergency_contact_name"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="Enter contact name"
                  />
                </div>

                <div class="space-y-2">
                  <label for="emergency_contact_phone" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
                    </svg>
                    Contact Phone
                  </label>
                  <input
                    id="emergency_contact_phone"
                    v-model="form.emergency_contact_phone"
                    type="tel"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="Enter contact phone"
                  />
                </div>

                <div class="space-y-2">
                  <label for="emergency_contact_relation" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                    </svg>
                    Relationship
                  </label>
                  <input
                    id="emergency_contact_relation"
                    v-model="form.emergency_contact_relation"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="e.g., Spouse, Parent, Sibling"
                  />
                </div>
              </div>
            </div>

            <!-- Medical Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Medical Information</h3>

              <div class="space-y-6">
                <div class="space-y-2">
                  <label for="medical_conditions" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                    </svg>
                    Medical Conditions
                  </label>
                  <textarea
                    id="medical_conditions"
                    v-model="form.medical_conditions"
                    rows="4"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                    placeholder="Previous surgeries, chronic conditions, etc."
                  ></textarea>
                </div>

                <div class="space-y-2">
                  <label for="allergies" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z"/>
                    </svg>
                    Allergies
                  </label>
                  <TagInput
                    id="allergies"
                    v-model="form.allergies"
                    placeholder="Type allergy and press Enter..."
                    help-text="Common allergies will be suggested as you type"
                  />
                </div>

                <div class="space-y-2">
                  <label for="current_medications" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
                    </svg>
                    Current Medications
                  </label>
                  <textarea
                    id="current_medications"
                    v-model="form.current_medications"
                    rows="4"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                    placeholder="Current medications, dosages, etc."
                  ></textarea>
                </div>
              </div>
            </div>

            <!-- Insurance Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Insurance Information</h3>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label for="insurance_provider" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                    </svg>
                    Insurance Provider
                  </label>
                  <input
                    id="insurance_provider"
                    v-model="form.insurance_provider"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="Enter insurance provider"
                  />
                </div>

                <div class="space-y-2">
                  <label for="insurance_number" class="block text-sm font-semibold text-gray-700 flex items-center">
                    <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
                    </svg>
                    Policy Number
                  </label>
                  <input
                    id="insurance_number"
                    v-model="form.insurance_number"
                    type="text"
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    placeholder="Enter policy number"
                  />
                </div>
              </div>
            </div>

            <!-- Additional Notes -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Additional Notes</h3>

              <div class="space-y-2">
                <label for="notes" class="block text-sm font-semibold text-gray-700 flex items-center">
                  <svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                  </svg>
                  Notes
                </label>
                <textarea
                  id="notes"
                  v-model="form.notes"
                  rows="4"
                  :readonly="isViewMode"
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                  :class="{ 'bg-gray-100 cursor-not-allowed': isViewMode }"
                  placeholder="Any additional information about the patient..."
                ></textarea>
              </div>
            </div>

            <!-- Submit Error -->
            <BaseTransition name="fade">
              <div v-if="submitError" class="bg-danger-50 border border-danger-200 rounded-xl p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-danger-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  <p class="text-sm text-danger-700">{{ submitError }}</p>
                </div>
              </div>
            </BaseTransition>

            <!-- Success Message -->
            <BaseTransition name="fade">
              <div v-if="submitSuccess" class="bg-success-50 border border-success-200 rounded-xl p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-success-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  <p class="text-sm text-success-700">{{ isEditing ? 'Patient updated successfully!' : 'Patient created successfully!' }}</p>
                </div>
              </div>
            </BaseTransition>

            <!-- Form Actions -->
            <div v-if="!isViewMode" class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-6 border-t border-neutral-200">
              <button
                type="submit"
                :disabled="loading || !hasChanges"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-xl text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <BaseLoading v-if="loading" type="spinner" size="small" color="white" :show-text="false" class="mr-2" />
                {{ loading ? 'Saving...' : (isEditing ? 'Update Patient' : 'Create Patient') }}
              </button>

              <button
                type="button"
                @click="resetForm"
                :disabled="loading"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Reset
              </button>

              <button
                type="button"
                @click="goBack"
                :disabled="loading"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Cancel
              </button>
            </div>
        </form>

        <!-- Consent Forms Section (View Mode Only) -->
        <div v-if="isViewMode" class="">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-neutral-900">Consent Forms</h3>

          </div>

          <!-- Consent Forms List -->
          <div v-if="consentForms && consentForms.length > 0" class="space-y-3">
            <div
              v-for="consentForm in consentForms"
              :key="consentForm.id"
              class="bg-gray-50 rounded-lg border border-gray-200 p-4 hover:bg-gray-100 transition-colors"
            >
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div class="flex-1 mb-3 sm:mb-0">
                  <h4 class="font-medium text-gray-900">{{ consentForm.title }}</h4>
                  <p class="text-sm text-gray-600 mt-1">{{ consentForm.description }}</p>
                  <div class="flex flex-col sm:flex-row sm:items-center mt-2 text-xs text-gray-500">
                    <div class="flex items-center">
                      <font-awesome-icon icon="fa-solid fa-calendar" class="w-3 h-3 mr-1" />
                      Created: {{ formatDate(consentForm.created_at) }}
                    </div>
                    <span class="hidden sm:inline mx-2">•</span>
                    <div class="flex items-center mt-1 sm:mt-0">
                      <font-awesome-icon icon="fa-solid fa-user" class="w-3 h-3 mr-1" />
                      {{ consentForm.patient_signed_at ? `Signed: ${formatDate(consentForm.patient_signed_at)}` : 'Not signed' }}
                    </div>
                  </div>
                </div>
                <div class="flex items-center space-x-2 sm:ml-4">
                  <button
                    @click="viewConsentForm(consentForm)"
                    class="flex-1 sm:flex-none inline-flex items-center justify-center px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
                  >
                    <font-awesome-icon icon="fa-solid fa-eye" class="w-4 h-4 sm:mr-2" />
                    <span class="ml-2 sm:ml-0">View</span>
                  </button>
                  <button
                    @click="downloadConsentFormPDF(consentForm)"
                    class="flex-1 sm:flex-none inline-flex items-center justify-center px-3 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 transition-all duration-200"
                  >
                    <font-awesome-icon icon="fa-solid fa-download" class="w-4 h-4 sm:mr-2" />
                    <span class="ml-2 sm:ml-0">PDF</span>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-8 bg-gray-50 rounded-lg border border-gray-200">
            <font-awesome-icon icon="fa-solid fa-file-signature" class="w-12 h-12 text-gray-400 mb-4" />
            <h4 class="text-lg font-medium text-gray-900 mb-2">No Consent Forms</h4>
            <p class="text-gray-600 mb-4">This patient doesn't have any consent forms yet.</p>
            <router-link
              :to="`/patients/${route.params.id}/consent/new`"
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
            >
              <font-awesome-icon icon="fa-solid fa-plus" class="w-4 h-4 mr-2" />
              Create First Consent Form
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </BaseTransition>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useNotificationStore } from '../stores/notification'
import { useConsentStore } from '../stores/consent'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'
import TagInput from '../components/TagInput.vue'
import AvatarUpload from '../components/AvatarUpload.vue'
import { useNavigation } from '../composables/useNavigation'

const route = useRoute()
const router = useRouter()
const { goBack } = useNavigation()
const patientStore = usePatientStore()
const notificationStore = useNotificationStore()
const consentStore = useConsentStore()

// Component state
const loading = ref(false)
const error = ref('')
const submitError = ref('')
const submitSuccess = ref(false)
const errors = ref({})
const originalPatient = ref(null)
const consentForms = ref([])

// Form data
const form = reactive({
  first_name: '',
  last_name: '',
  date_of_birth: '',
  gender: '',
  phone: '',
  email: '',
  address: '',
  emergency_contact_name: '',
  emergency_contact_phone: '',
  emergency_contact_relation: '',
  medical_conditions: '',
  allergies: [],
  current_medications: '',
  insurance_provider: '',
  insurance_number: '',
  notes: '',
  avatar_path: ''
})

// Computed properties
const isEditing = computed(() => !!route.params.id)
const isViewMode = computed(() => isEditing.value && route.path.includes('/patients/') && !route.path.includes('/edit'))
const isEditMode = computed(() => isEditing.value && route.path.includes('/edit'))
const isCreateMode = computed(() => !isEditing.value)

const hasChanges = computed(() => {
  if (!isEditing.value) {
    // In create mode, check if required fields are filled
    return !!(form.first_name && form.last_name && form.date_of_birth && form.gender && form.phone)
  }

  if (!originalPatient.value) return false

  // Check basic field changes
  const allergiesChanged = JSON.stringify(form.allergies.sort()) !== JSON.stringify((originalPatient.value.allergies || '').split(',').map(a => a.trim()).filter(a => a).sort())

  return (
    form.first_name !== (originalPatient.value.first_name || '') ||
    form.last_name !== (originalPatient.value.last_name || '') ||
    form.date_of_birth !== (originalPatient.value.date_of_birth ? new Date(originalPatient.value.date_of_birth).toISOString().split('T')[0] : '') ||
    form.gender !== (originalPatient.value.gender || '') ||
    form.phone !== (originalPatient.value.phone || '') ||
    form.email !== (originalPatient.value.email || '') ||
    form.address !== (originalPatient.value.address || '') ||
    form.emergency_contact_name !== (originalPatient.value.emergency_contact_name || '') ||
    form.emergency_contact_phone !== (originalPatient.value.emergency_contact_phone || '') ||
    form.emergency_contact_relation !== (originalPatient.value.emergency_contact_relation || '') ||
    form.medical_conditions !== (originalPatient.value.medical_conditions || '') ||
    allergiesChanged ||
    form.current_medications !== (originalPatient.value.current_medications || '') ||
    form.insurance_provider !== (originalPatient.value.insurance_provider || '') ||
    form.insurance_number !== (originalPatient.value.insurance_number || '') ||
    form.notes !== (originalPatient.value.notes || '') ||
    form.avatar_path !== (originalPatient.value.avatar_path || '')
  )
})

// Methods
const resetForm = () => {
  if (isEditing.value && originalPatient.value) {
    // Reset to original patient data for edit mode
    const patient = originalPatient.value
    form.first_name = patient.first_name || ''
    form.last_name = patient.last_name || ''
    form.email = patient.email || ''
    form.phone = patient.phone || ''
    form.gender = patient.gender || ''
    form.address = patient.address || ''
    form.emergency_contact_name = patient.emergency_contact_name || ''
    form.emergency_contact_phone = patient.emergency_contact_phone || ''
    form.emergency_contact_relation = patient.emergency_contact_relation || ''
    // Convert allergies string to array
    form.allergies = patient.allergies ? patient.allergies.split(',').map(a => a.trim()).filter(a => a) : []
    form.medical_conditions = patient.medical_conditions || ''
    form.current_medications = patient.current_medications || ''
    form.insurance_provider = patient.insurance_provider || ''
    form.insurance_number = patient.insurance_number || ''
    form.notes = patient.notes || ''
    form.avatar_path = patient.avatar_path || ''

    // Format date of birth for date input (YYYY-MM-DD)
    if (patient.date_of_birth) {
      const date = new Date(patient.date_of_birth)
      if (!isNaN(date.getTime())) {
        form.date_of_birth = date.toISOString().split('T')[0]
      }
    } else {
      form.date_of_birth = ''
    }
  } else {
    // Reset to empty form for create mode
    form.first_name = ''
    form.last_name = ''
    form.date_of_birth = ''
    form.gender = ''
    form.phone = ''
    form.email = ''
    form.address = ''
    form.emergency_contact_name = ''
    form.emergency_contact_phone = ''
    form.emergency_contact_relation = ''
    form.medical_conditions = ''
    form.allergies = []
    form.current_medications = ''
    form.insurance_provider = ''
    form.insurance_number = ''
    form.notes = ''
    form.avatar_path = ''
  }
  errors.value = {}
  submitError.value = ''
  submitSuccess.value = false
}

const loadPatient = async (patientId) => {
  try {
    loading.value = true
    error.value = ''
    const result = await patientStore.fetchPatient(patientId)

    if (result && result.success && patientStore.currentPatient) {
      const patient = patientStore.currentPatient
      originalPatient.value = { ...patient }

      // Populate form with patient data
      form.id = patient.id
      form.first_name = patient.first_name || ''
      form.last_name = patient.last_name || ''
      form.email = patient.email || ''
      form.phone = patient.phone || ''
      form.gender = patient.gender || ''
      form.address = patient.address || ''
      form.emergency_contact_name = patient.emergency_contact_name || ''
      form.emergency_contact_phone = patient.emergency_contact_phone || ''
      form.emergency_contact_relation = patient.emergency_contact_relation || ''
      // Convert allergies string to array
      form.allergies = patient.allergies ? patient.allergies.split(',').map(a => a.trim()).filter(a => a) : []
      form.medical_conditions = patient.medical_conditions || ''
      form.current_medications = patient.current_medications || ''
      form.insurance_provider = patient.insurance_provider || ''
      form.insurance_number = patient.insurance_number || ''
      form.notes = patient.notes || ''
      form.avatar_path = patient.avatar_path || ''

      // Format date of birth for date input (YYYY-MM-DD)
      if (patient.date_of_birth) {
        const date = new Date(patient.date_of_birth)
        if (!isNaN(date.getTime())) {
          form.date_of_birth = date.toISOString().split('T')[0]
        }
      } else {
        form.date_of_birth = ''
      }

      // Load consent forms if in view mode
      if (isViewMode.value) {
        await loadConsentForms(patientId)
      }
    } else {
      const errorMessage = result?.error || 'Patient not found'
      error.value = errorMessage
      console.error('Failed to load patient:', errorMessage)
      router.push('/patients')
    }
  } catch (err) {
    console.error('Error loading patient:', err)
    error.value = 'Failed to load patient data'
    router.push('/patients')
  } finally {
    loading.value = false
  }
}

const handleAvatarUpdated = (avatarPath) => {
  form.avatar_path = avatarPath || ''
}

const validateForm = () => {
  errors.value = {}
  
  // Required fields
  if (!form.first_name.trim()) {
    errors.value.first_name = 'First name is required'
  }
  
  if (!form.last_name.trim()) {
    errors.value.last_name = 'Last name is required'
  }
  
  if (!form.date_of_birth) {
    errors.value.date_of_birth = 'Date of birth is required'
  }
  
  if (!form.gender) {
    errors.value.gender = 'Gender is required'
  }
  
  if (!form.phone.trim()) {
    errors.value.phone = 'Phone number is required'
  }
  
  // Email validation
  if (form.email && !/\S+@\S+\.\S+/.test(form.email)) {
    errors.value.email = 'Please enter a valid email address'
  }
  
  // Phone validation (basic)
  if (form.phone && !/^[\d\s\-\+\(\)]+$/.test(form.phone)) {
    errors.value.phone = 'Please enter a valid phone number'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) {
    submitError.value = 'Please fix the form errors before submitting'
    return
  }

  try {
    loading.value = true
    submitError.value = ''
    submitSuccess.value = false

    // Convert allergies array to string for backend
    const submitData = {
      ...form,
      allergies: form.allergies.join(', ')
    }

    if (isEditing.value) {
      await patientStore.updatePatient(route.params.id, submitData)
      // Update original patient with converted data
      originalPatient.value = {
        ...form,
        allergies: submitData.allergies
      }
      submitSuccess.value = true
      notificationStore.showSuccess('Patient updated successfully')

      // Hide success message after 3 seconds
      setTimeout(() => {
        submitSuccess.value = false
      }, 3000)
    } else {
      await patientStore.createPatient(submitData)
      notificationStore.showSuccess('Patient created successfully')
      router.push('/patients')
    }
  } catch (err) {
    console.error('Error saving patient:', err)
    submitError.value = isEditing.value ? 'Failed to update patient' : 'Failed to create patient'
    notificationStore.showError(submitError.value)
  } finally {
    loading.value = false
  }
}

// Consent form methods
const loadConsentForms = async (patientId) => {
  try {
    // Load consent forms for this patient
    await consentStore.fetchConsentForms({ patient_id: patientId })
    consentForms.value = consentStore.getPatientConsentForms(patientId)
    
    // Also make sure consent templates are loaded for preview generation
    if (!consentStore.consentTemplates || consentStore.consentTemplates.length === 0) {
      await consentStore.fetchConsentTemplates()
    }
  } catch (error) {
    console.error('Error loading consent forms:', error)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const viewConsentForm = (consentForm) => {
  // Open consent form in new window for viewing
  const previewContent = generateConsentPreview(consentForm)
  const previewWindow = window.open('', '_blank', 'width=800,height=600,scrollbars=yes')
  previewWindow.document.write(previewContent)
  previewWindow.document.close()
}

const generateConsentPreview = (consentForm) => {
  // Get the patient data
  const patient = form
  const currentDate = new Date().toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  // Get the template content (we need to load it or store it when saving)
  const template = consentStore.getConsentTemplateById(parseInt(consentForm.consent_template_id))
  let consentContent = ''
  
  if (template) {
    // Get doctor information if available
    const doctorName = consentForm.doctor_id ? 
      `${consentForm.doctor?.first_name || 'Doctor'} ${consentForm.doctor?.last_name || ''}`.trim() : 
      'Doctor'
    
    // Apply the same placeholder replacement logic as in the create form
    consentContent = template.content
      .replace(/\[PATIENT_NAME\]/g, `${patient.first_name} ${patient.last_name}`)
      .replace(/\[PATIENT_FIRST_NAME\]/g, patient.first_name)
      .replace(/\[PATIENT_LAST_NAME\]/g, patient.last_name)
      .replace(/\[CURRENT_DATE\]/g, currentDate)
      .replace(/\[TODAY\]/g, currentDate)
      .replace(/\[DATE\]/g, currentDate)
      .replace(/\[DOCTOR_NAME\]/g, doctorName)
      .replace(/\[CLINIC_NAME\]/g, 'Clinic')
  }

  // Format the consent content with the same styling as the create form
  const formattedContent = formatConsentContent(consentContent)

  return `
    <!DOCTYPE html>
    <html>
    <head>
      <title>${consentForm.title}</title>
      <meta charset="utf-8">
      <style>
        body { 
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          line-height: 1.6; 
          max-width: 900px; 
          margin: 0 auto; 
          padding: 20px;
          background: white;
          color: #374151;
        }
        .print-btn {
          margin-bottom: 20px; 
          padding: 10px 20px; 
          background: #3b82f6; 
          color: white; 
          border: none; 
          border-radius: 8px; 
          cursor: pointer;
          font-size: 14px;
        }
        .print-btn:hover { background: #2563eb; }
        
        .consent-content {
          background: white;
          padding: 24px;
          margin-bottom: 20px;
        }
        
        .consent-content h2 {
          color: #1f2937;
          font-size: 20px;
          font-weight: 700;
          margin-top: 24px;
          margin-bottom: 12px;
          border-bottom: 2px solid #3b82f6;
          padding-bottom: 8px;
        }
        
        .consent-content h3 {
          color: #374151;
          font-size: 18px;
          font-weight: 600;
          margin-top: 16px;
          margin-bottom: 8px;
        }
        
        .consent-content p {
          color: #4b5563;
          margin-bottom: 12px;
          line-height: 1.7;
        }
        
        .consent-content li {
          margin-left: 16px;
          margin-bottom: 4px;
          color: #4b5563;
        }
        
        .consent-content div {
          margin-bottom: 8px;
        }
        
        .consent-content strong {
          color: #374151;
        }

        /* Signature Section Styling - Traditional document style */
        .signature-section {
          margin-top: 30px;
          padding: 20px;
          page-break-inside: avoid;
        }
        
        .signature-section h3 {
          margin-bottom: 25px;
          color: #374151;
          font-size: 16px;
          font-weight: 600;
          border-bottom: 1px solid #d1d5db;
          padding-bottom: 8px;
        }
        
        .footer-text {
          margin-top: 20px;
          text-align: center;
          font-size: 12px;
          color: #6b7280;
        }
        
        @media print { 
          body { margin: 0; padding: 15px; }
          .print-btn { display: none; }
          .signature-container { display: flex !important; }
        }
      </style>
    </head>
    <body>
      <button class="print-btn" onclick="window.print()">🖨️ Print / Save as PDF</button>
      
      <div class="consent-content">
        <div>${formattedContent}</div>

        <!-- Signature Section - Traditional document style -->
        <div class="signature-section">
          <h3>Signatures</h3>
          
          <div style="display: flex; justify-content: space-between; gap: 40px; margin-bottom: 20px;">
            <div style="flex: 1;">
              <div style="margin-bottom: 8px; font-weight: 500; color: #374151;">
                Patient Signature:
              </div>
              <div style="border-bottom: 2px solid #374151; height: 120px; margin-top: 8px; display: flex; align-items: center; justify-content: center; position: relative; background: white;">
                ${consentForm.patient_signature 
                  ? `<img src="${consentForm.patient_signature}" alt="Patient Signature" style="max-width: 100%; max-height: 110px; object-fit: contain;" />`
                  : '<span style="font-size: 12px; color: #9ca3af;">Not signed</span>'
                }
              </div>
              <div style="margin-top: 8px; font-size: 12px; color: #6b7280;">
                Date: ${formatDate(consentForm.patient_signed_at || consentForm.created_at)}
              </div>
            </div>

            <div style="flex: 1;">
              <div style="margin-bottom: 8px; font-weight: 500; color: #374151;">
                Witness Signature:
              </div>
              <div style="border-bottom: 2px solid #374151; height: 120px; margin-top: 8px; display: flex; align-items: center; justify-content: center; position: relative; background: white;">
                ${consentForm.witness_signature
                  ? `<img src="${consentForm.witness_signature}" alt="Witness Signature" style="max-width: 100%; max-height: 110px; object-fit: contain;" />`
                  : '<span style="font-size: 12px; color: #9ca3af;">No witness signature</span>'
                }
              </div>
              <div style="margin-top: 8px; font-size: 12px; color: #6b7280;">
                Date: ${formatDate(consentForm.witness_signed_at || consentForm.created_at)}
                ${consentForm.witness_name ? `<br>Witness: ${consentForm.witness_name}` : ''}
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="footer-text">
          <p>This consent form has been electronically generated and is legally binding.</p>
          <p>Please retain a copy for your records.</p>
        </div>
      </div>
    </body>
    </html>
  `
}

// Helper function to format consent content - same as in create form
const formatConsentContent = (content) => {
  if (!content) return ''

  // Split content into lines and format each section
  const lines = content.split('\n')
  const formattedLines = lines.map(line => {
    // Format main headings (all caps)
    if (line.match(/^[A-Z\s]+$/)) {
      return `<h2>${line}</h2>`
    }

    // Format section headings (UPPERCASE)
    if (line.match(/^[A-Z\s]+:$/)) {
      return `<h3>${line.replace(':', '')}</h3>`
    }

    // Format bullet points
    if (line.trim().startsWith('- ')) {
      return `<li>${line.trim().substring(2)}</li>`
    }

    // Format patient info lines (contain colons)
    if (line.includes(':') && !line.includes('http') && line.split(':').length === 2) {
      const [label, value] = line.split(':')
      return `<div><strong>${label.trim()}:</strong> ${value.trim()}</div>`
    }

    // Regular paragraphs
    if (line.trim()) {
      return `<p>${line.trim()}</p>`
    }

    // Empty lines
    return '<br>'
  })

  return formattedLines.join('')
}

const downloadConsentFormPDF = async (consentForm) => {
  try {
    // For now, open the preview window - in a real app you'd generate a PDF
    // You could use libraries like jsPDF or html2pdf, or call a backend endpoint
    viewConsentForm(consentForm)
    notificationStore.showInfo('PDF generation would be implemented here. For now, use the print function in the preview window.')
  } catch (error) {
    console.error('Error downloading PDF:', error)
    notificationStore.showError('Failed to generate PDF')
  }
}

// Clear messages when form changes
watch(() => form, () => {
  submitError.value = ''
  submitSuccess.value = false
}, { deep: true })

// Lifecycle
onMounted(() => {
  if (isEditing.value) {
    const patientId = route.params.id

    if (!patientId || patientId === 'new') {
      router.push('/patients')
      return
    }

    // Ensure patientId is a number
    const numericId = parseInt(patientId, 10)
    if (isNaN(numericId)) {
      router.push('/patients')
      return
    }

    loadPatient(numericId)
  } else {
    // For create mode, reset form to empty
    resetForm()
  }
})


</script>

<style scoped>
@import "../styles/main.css";
</style>