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

  <!-- Edit Form -->
  <BaseTransition name="slide-up">
    <div v-if="!loading && !error" class="bg-white rounded-2xl shadow-lg border border-neutral-100 overflow-hidden">
      <div class="p-6 sm:p-8">
        <form @submit.prevent="handleSubmit" class="space-y-8">
            <!-- Basic Information -->
            <div class="space-y-6">
              <h3 class="text-lg font-semibold text-neutral-900 border-b border-neutral-200 pb-2">Basic Information</h3>
              
              <!-- Avatar Upload -->
              <div class="flex justify-center">
                <AvatarUpload
                  :user="form"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.first_name }"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.last_name }"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.date_of_birth }"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.gender }"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.phone }"
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
                    class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
                    :class="{ 'border-red-500': errors.email }"
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
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
                  :class="{ 'border-red-500': errors.address }"
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
                  class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white resize-none"
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
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-6 border-t border-neutral-200">
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
                @click="$router.go(-1)"
                :disabled="loading"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-6 py-3 border border-neutral-300 rounded-xl text-sm font-medium text-neutral-700 bg-white hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                Cancel
              </button>
            </div>
        </form>
      </div>
    </div>
  </BaseTransition>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePatientStore } from '../stores/patient'
import { useNotificationStore } from '../stores/notification'
import BaseLoading from '../components/BaseLoading.vue'
import BaseTransition from '../components/BaseTransition.vue'
import TagInput from '../components/TagInput.vue'
import AvatarUpload from '../components/AvatarUpload.vue'

const route = useRoute()
const router = useRouter()
const patientStore = usePatientStore()
const notificationStore = useNotificationStore()

// Component state
const loading = ref(false)
const error = ref('')
const submitError = ref('')
const submitSuccess = ref(false)
const errors = ref({})
const originalPatient = ref(null)

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