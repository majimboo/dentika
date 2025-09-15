<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 relative overflow-hidden">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-5">
      <div class="absolute top-20 left-10 w-32 h-32 bg-blue-400 rounded-full blur-3xl"></div>
      <div class="absolute top-40 right-20 w-40 h-40 bg-indigo-400 rounded-full blur-3xl"></div>
      <div class="absolute bottom-20 left-1/3 w-36 h-36 bg-purple-400 rounded-full blur-3xl"></div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center min-h-screen relative z-10">
      <div class="text-center bg-white/80 backdrop-blur-sm rounded-3xl p-12 shadow-2xl border border-white/20">
        <div class="relative mb-8">
          <div class="w-20 h-20 border-4 border-blue-200 rounded-full animate-spin mx-auto"></div>
          <div class="absolute inset-0 w-20 h-20 border-4 border-transparent border-t-blue-600 rounded-full animate-spin mx-auto"></div>
        </div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">Loading Your Experience</h3>
        <p class="text-gray-600">Preparing your personalized scheduling page...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex items-center justify-center min-h-screen relative z-10">
      <div class="text-center bg-white/90 backdrop-blur-sm rounded-3xl p-12 shadow-2xl border border-white/20 max-w-md mx-auto">
        <div class="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg class="w-10 h-10 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.081 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-gray-900 mb-3">Oops! Something went wrong</h2>
        <p class="text-gray-600 mb-8 leading-relaxed">{{ error }}</p>
        <button @click="loadClinicInfo"
          class="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-8 py-3 rounded-xl font-semibold hover:from-blue-700 hover:to-indigo-700 transform hover:scale-105 transition-all duration-200 shadow-lg hover:shadow-xl">
          Try Again
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else class="relative z-10">
      <!-- Full-Width Hero Header -->
      <div class="relative min-h-[85vh] md:min-h-screen flex items-center justify-center overflow-hidden">
        <!-- Hero Background -->
        <div class="absolute inset-0 bg-gradient-to-br from-blue-600 via-indigo-700 to-purple-800"></div>

        <!-- Animated Background Elements -->
        <div class="absolute inset-0">
          <div class="hidden md:block absolute top-20 left-10 w-96 h-96 bg-white/5 rounded-full blur-3xl animate-float"></div>
          <div class="hidden md:block absolute top-40 right-20 w-80 h-80 bg-purple-300/10 rounded-full blur-3xl animate-float" style="animation-delay: 2s;"></div>
          <div class="hidden md:block absolute bottom-20 left-1/3 w-64 h-64 bg-indigo-300/10 rounded-full blur-3xl animate-float" style="animation-delay: 4s;"></div>
          <div class="absolute top-1/4 left-1/4 w-2 h-2 bg-white/20 rounded-full animate-pulse-slow"></div>
          <div class="absolute top-1/3 right-1/3 w-1 h-1 bg-white/30 rounded-full animate-pulse-slow" style="animation-delay: 1s;"></div>
        </div>

        <!-- Hero Content -->
        <div class="relative z-10 text-center px-6 py-12 md:px-4 max-w-6xl mx-auto">
          <!-- Main Logo/Icon -->
          <div class="flex items-center justify-center mb-8 md:mb-12">
            <div class="relative">
              <div class="w-24 h-24 md:w-32 md:h-32 bg-white/10 backdrop-blur-sm rounded-3xl flex items-center justify-center shadow-2xl border border-white/20">
                <svg class="w-12 h-12 md:w-16 md:h-16 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <div class="absolute inset-0 w-24 h-24 md:w-32 md:h-32 bg-white/20 rounded-3xl animate-pulse-slow"></div>
              <div class="absolute -top-2 -right-2 md:-top-3 md:-right-3 w-8 h-8 md:w-10 md:h-10 bg-green-400 rounded-full border-4 border-white flex items-center justify-center shadow-lg">
                <svg class="w-4 h-4 md:w-5 md:h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Clinic Name -->
          <h1 class="text-4xl md:text-6xl lg:text-8xl font-bold text-white mb-4 md:mb-6 tracking-tight leading-tight">
            {{ clinic.name }}
          </h1>

          <!-- Tagline -->
          <p class="text-lg md:text-xl lg:text-2xl text-blue-100 mb-8 md:mb-12 max-w-3xl mx-auto leading-relaxed px-4">
            Your trusted partner in dental excellence
          </p>

          <!-- Clinic Details -->
          <div class="bg-white/10 backdrop-blur-md rounded-2xl p-6 md:p-8 max-w-4xl mx-auto border border-white/20 shadow-2xl">
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 md:gap-8">
              <!-- Address -->
              <div class="text-center">
                <div class="w-10 h-10 md:w-12 md:h-12 bg-white/20 rounded-xl flex items-center justify-center mx-auto mb-3 md:mb-4">
                  <svg class="w-5 h-5 md:w-6 md:h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
                <h3 class="text-white font-semibold mb-2 text-sm md:text-base">Location</h3>
                <p class="text-blue-100 text-xs md:text-sm leading-relaxed">{{ clinic.address }}</p>
              </div>

              <!-- Phone -->
              <div class="text-center">
                <div class="w-10 h-10 md:w-12 md:h-12 bg-white/20 rounded-xl flex items-center justify-center mx-auto mb-3 md:mb-4">
                  <svg class="w-5 h-5 md:w-6 md:h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                  </svg>
                </div>
                <h3 class="text-white font-semibold mb-2 text-sm md:text-base">Phone</h3>
                <p class="text-blue-100 font-medium text-sm md:text-base">{{ clinic.phone }}</p>
              </div>

              <!-- Email -->
              <div class="text-center">
                <div class="w-10 h-10 md:w-12 md:h-12 bg-white/20 rounded-xl flex items-center justify-center mx-auto mb-3 md:mb-4">
                  <svg class="w-5 h-5 md:w-6 md:h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <h3 class="text-white font-semibold mb-2 text-sm md:text-base">Email</h3>
                <p class="text-blue-100 font-medium text-sm md:text-base">{{ clinic.email }}</p>
              </div>
            </div>
          </div>

          <!-- Call to Action -->
          <div class="mt-8 md:mt-12">
            <div class="inline-flex items-center space-x-2 md:space-x-3 bg-white/20 backdrop-blur-sm rounded-full px-6 py-3 md:px-8 md:py-4 border border-white/30">
              <svg class="w-4 h-4 md:w-5 md:h-5 text-white animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
              </svg>
              <span class="text-white font-medium text-sm md:text-base">Book your appointment</span>
            </div>
          </div>
        </div>

        <!-- Scroll Indicator -->
        <div class="hidden md:block absolute bottom-8 left-1/2 transform -translate-x-1/2">
          <div class="w-6 h-10 border-2 border-white/30 rounded-full flex justify-center">
            <div class="w-1 h-3 bg-white/50 rounded-full mt-2 animate-pulse"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Form Container -->
    <div class="max-w-4xl mx-auto my-10 px-4 md:px-6">
      <div class="bg-white/90 backdrop-blur-sm rounded-2xl md:rounded-3xl shadow-2xl overflow-hidden border border-white/20">
        <!-- Form Header -->
        <div class="bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 text-white p-6 md:p-8 relative overflow-hidden">
          <div class="absolute inset-0 bg-black/10"></div>
          <div class="relative z-10">
            <div class="flex items-center space-x-3 md:space-x-4 mb-3 md:mb-4">
              <div class="w-10 h-10 md:w-12 md:h-12 bg-white/20 rounded-xl flex items-center justify-center flex-shrink-0">
                <svg class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <div class="min-w-0 flex-1">
                <h2 class="text-xl md:text-3xl font-bold leading-tight">Request an Appointment</h2>
                <p class="text-blue-100 mt-1 text-sm md:text-base">Fill out the form below and we'll get back to you soon</p>
              </div>
            </div>
          </div>
          <div class="hidden md:block absolute top-4 right-4 w-20 h-20 bg-white/10 rounded-full blur-xl"></div>
          <div class="hidden md:block absolute bottom-4 left-4 w-16 h-16 bg-white/10 rounded-full blur-xl"></div>
        </div>

        <!-- Step Navigation -->
        <div class="px-4 md:px-8 pt-6 md:pt-8">
          <!-- Progress Bar -->
          <div class="mb-8">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg md:text-xl font-semibold text-gray-900">Step {{ currentStep }} of {{ totalSteps }}</h2>
              <span class="text-sm text-gray-500">{{ Math.round((currentStep / totalSteps) * 100) }}% Complete</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div class="bg-gradient-to-r from-blue-600 to-indigo-600 h-2 rounded-full transition-all duration-300 ease-out"
                   :style="{ width: `${(currentStep / totalSteps) * 100}%` }"></div>
            </div>
          </div>

          <!-- Step Indicators -->
          <div class="flex items-center justify-between mb-8 overflow-x-auto">
            <div v-for="step in steps" :key="step.id" class="flex items-center flex-shrink-0"
                 :class="{ 'flex-1': step.id < totalSteps }">
              <div class="flex flex-col items-center">
                <!-- Step Circle -->
                <button @click="goToStep(step.id)"
                        :disabled="!completedSteps.includes(step.id) && step.id > currentStep"
                        class="relative w-10 h-10 md:w-12 md:h-12 rounded-full flex items-center justify-center transition-all duration-200 mb-2"
                        :class="[
                          step.id < currentStep || completedSteps.includes(step.id)
                            ? 'bg-green-500 text-white shadow-lg'
                            : step.id === currentStep
                            ? 'bg-blue-600 text-white shadow-lg ring-4 ring-blue-100'
                            : 'bg-gray-200 text-gray-400 cursor-not-allowed'
                        ]">
                  <!-- Step Icon -->
                  <svg v-if="step.icon === 'phone'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                  </svg>
                  <svg v-else-if="step.icon === 'user'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <svg v-else-if="step.icon === 'calendar'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <svg v-else-if="step.icon === 'clock'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else-if="step.icon === 'clipboard'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <svg v-else-if="step.icon === 'check-circle'" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>

                  <!-- Check mark for completed steps -->
                  <svg v-if="step.id < currentStep || completedSteps.includes(step.id)" class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>

                <!-- Step Title -->
                <div class="text-center">
                  <div class="text-xs md:text-sm font-medium text-gray-900 mb-1">{{ step.title }}</div>
                  <div class="text-xs text-gray-500 hidden md:block">{{ step.description }}</div>
                </div>
              </div>

              <!-- Connector Line -->
              <div v-if="step.id < totalSteps" class="flex-1 h-px bg-gray-300 mx-4 hidden md:block"></div>
            </div>
          </div>
        </div>

        <!-- Form Content -->
        <form @submit.prevent="handleSubmit" class="p-4 md:p-8 space-y-6 md:space-y-8">
          <!-- Step 1: Phone Number Check -->
          <div v-if="currentStep === 1" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-green-500 to-emerald-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Phone Verification</h3>
            </div>

            <div class="space-y-4">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Phone Number <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <input v-model="form.phone" type="tel" required
                         class="w-full px-4 py-4 pl-12 border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-green-100 focus:border-green-500 transition-all duration-200 text-gray-900 placeholder-gray-400"
                         placeholder="+63 9XX XXX XXXX" @input="resetPhoneCheck" />
                  <div class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                    </svg>
                  </div>
                </div>
              </div>

              <div class="flex items-center space-x-4">
                <button type="button" @click="checkPatientByPhone"
                        :disabled="!form.phone || checkingPhone || patientFound"
                        class="flex items-center space-x-2 bg-gradient-to-r from-green-600 to-emerald-600 text-white px-6 py-3 rounded-xl font-semibold hover:from-green-700 hover:to-emerald-700 focus:ring-4 focus:ring-green-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200">
                  <svg v-if="checkingPhone" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>{{ checkingPhone ? 'Checking...' : patientFound ? 'Phone Verified' : 'Verify Phone' }}</span>
                </button>

                <div v-if="phoneChecked" class="flex items-center space-x-2">
                  <div v-if="patientFound" class="flex items-center space-x-2 text-green-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span class="text-sm font-medium">Patient found - information auto-filled</span>
                  </div>
                  <div v-else class="flex items-center space-x-2 text-blue-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    <span class="text-sm font-medium">New patient - please fill in your information</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 2: Personal Information -->
          <div v-if="currentStep === 2" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Personal Information</h3>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  First Name <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <input v-model="form.first_name" type="text" required :readonly="patientFound"
                         :class="[
                           'w-full px-4 py-4 pl-12 border-2 rounded-xl transition-all duration-200 text-gray-900 placeholder-gray-400',
                           patientFound
                             ? 'border-green-200 bg-green-50 cursor-not-allowed'
                             : 'border-gray-200 focus:ring-4 focus:ring-blue-100 focus:border-blue-500'
                         ]"
                         placeholder="Enter your first name" />
                  <div class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <div v-if="patientFound" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-green-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Last Name <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <input v-model="form.last_name" type="text" required :readonly="patientFound"
                         :class="[
                           'w-full px-4 py-4 pl-12 border-2 rounded-xl transition-all duration-200 text-gray-900 placeholder-gray-400',
                           patientFound
                             ? 'border-green-200 bg-green-50 cursor-not-allowed'
                             : 'border-gray-200 focus:ring-4 focus:ring-blue-100 focus:border-blue-500'
                         ]"
                         placeholder="Enter your last name" />
                  <div class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <div v-if="patientFound" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-green-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
              </div>

              <div class="space-y-2 md:col-span-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Email Address <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <input v-model="form.email" type="email" required :readonly="patientFound"
                         :class="[
                           'w-full px-4 py-4 pl-12 border-2 rounded-xl transition-all duration-200 text-gray-900 placeholder-gray-400',
                           patientFound
                             ? 'border-green-200 bg-green-50 cursor-not-allowed'
                             : 'border-gray-200 focus:ring-4 focus:ring-blue-100 focus:border-blue-500'
                         ]"
                         placeholder="your.email@example.com" />
                  <div class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div v-if="patientFound" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-green-600">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 3: Date Selection -->
          <div v-if="currentStep === 3" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Select Date</h3>
            </div>

            <div class="space-y-6">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Preferred Date <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <input v-model="form.preferred_date" type="date" required :min="minDate"
                         class="w-full px-4 py-4 pl-12 border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-purple-100 focus:border-purple-500 transition-all duration-200 text-gray-900 placeholder-gray-400"
                         @change="fetchAvailableTimeSlots" />
                  <div class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 4: Time Selection -->
          <div v-if="currentStep === 4" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Select Time</h3>
            </div>

            <div class="space-y-6">
              <div class="space-y-4">
                <div class="flex items-center justify-between">
                  <label class="block text-sm font-semibold text-gray-700">
                    Preferred Time <span class="text-red-500">*</span>
                  </label>
                  <div v-if="loadingTimeSlots" class="flex items-center text-sm text-gray-500">
                    <div class="w-4 h-4 border-2 border-purple-200 border-t-purple-500 rounded-full animate-spin mr-2"></div>
                    Loading availability...
                  </div>
                </div>

                <p class="text-xs text-gray-500 mb-3">
                  Select your preferred appointment time. We'll confirm the closest available slot.
                </p>

                <!-- Time Slots Grid -->
                <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-2 md:gap-3">
                  <label v-for="slot in allTimeSlots" :key="slot.time"
                         class="relative flex items-center justify-center p-2 md:p-3 border-2 border-gray-200 rounded-lg md:rounded-xl cursor-pointer hover:border-purple-300 hover:bg-purple-50 transition-all duration-200 group active:scale-95"
                         :class="{
                           'border-purple-500 bg-purple-50 ring-2 ring-purple-100': form.preferred_time === slot.time,
                           'border-green-300 bg-green-50': form.preferred_time === slot.time,
                           'border-gray-100 bg-gray-50 cursor-not-allowed opacity-50': !slot.available
                         }">
                    <input v-model="form.preferred_time" :value="slot.time" type="radio" name="timeSlot" class="sr-only"
                           :disabled="!slot.available" />
                    <div class="text-center">
                      <div class="font-semibold text-gray-900 group-hover:text-purple-900 transition-colors text-xs md:text-sm">
                        {{ slot.display }}
                      </div>
                      <div class="text-xs text-gray-500 mt-1 hidden md:block">
                        {{ slot.time }}
                      </div>
                      <!-- Availability indicator -->
                      <div v-if="!slot.available"
                           class="absolute inset-0 flex items-center justify-center bg-gray-900/10 rounded-lg md:rounded-xl">
                        <div class="w-1 h-1 bg-gray-400 rounded-full"></div>
                      </div>
                    </div>
                    <div v-if="form.preferred_time === slot.time"
                         class="absolute -top-1 -right-1 md:-top-2 md:-right-2 w-4 h-4 md:w-5 md:h-5 bg-purple-500 rounded-full flex items-center justify-center">
                      <svg class="w-2 h-2 md:w-3 md:h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd"
                              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                              clip-rule="evenodd" />
                      </svg>
                    </div>
                  </label>
                </div>

                <!-- Legend -->
                <div class="flex items-center justify-center space-x-6 text-xs text-gray-500 mt-4">
                  <div class="flex items-center space-x-2">
                    <div class="w-3 h-3 bg-purple-100 border border-purple-300 rounded"></div>
                    <span>Selected</span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <div class="w-3 h-3 bg-gray-50 border border-gray-200 rounded"></div>
                    <span>Available</span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <div class="w-3 h-3 bg-gray-100 border border-gray-300 rounded relative">
                      <div class="absolute inset-0 flex items-center justify-center">
                        <div class="w-1 h-1 bg-gray-400 rounded-full"></div>
                      </div>
                    </div>
                    <span>Unavailable</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 5: Additional Information -->
          <div v-if="currentStep === 5" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Additional Details</h3>
            </div>

            <div class="space-y-6">
              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  What brings you to the dentist? <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <textarea v-model="form.symptoms" required rows="4"
                            class="w-full px-4 py-4 pl-12 border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-purple-100 focus:border-purple-500 transition-all duration-200 text-gray-900 placeholder-gray-400 resize-none"
                            placeholder="Please describe your dental concern or the reason for your visit..."></textarea>
                  <div class="absolute left-4 top-4 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Additional Notes
                </label>
                <div class="relative">
                  <textarea v-model="form.additional_notes" rows="3"
                            class="w-full px-4 py-4 pl-12 border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-purple-100 focus:border-purple-500 transition-all duration-200 text-gray-900 placeholder-gray-400 resize-none"
                            placeholder="Any additional information you'd like us to know..."></textarea>
                  <div class="absolute left-4 top-4 text-gray-400">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 6: Review & Submit -->
          <div v-if="currentStep === 6" class="space-y-6">
            <div class="flex items-center space-x-3 mb-6">
              <div class="w-10 h-10 bg-gradient-to-br from-green-500 to-emerald-500 rounded-xl flex items-center justify-center">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">Review & Submit</h3>
            </div>

            <div class="space-y-6">
              <!-- Review Summary -->
              <div class="bg-gray-50 rounded-xl p-6 space-y-4">
                <h4 class="text-lg font-semibold text-gray-900 mb-4">Appointment Summary</h4>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="space-y-3">
                    <div>
                      <span class="text-sm font-medium text-gray-500">Patient:</span>
                      <p class="text-gray-900">{{ form.first_name }} {{ form.last_name }}</p>
                    </div>
                    <div>
                      <span class="text-sm font-medium text-gray-500">Phone:</span>
                      <p class="text-gray-900">{{ form.phone }}</p>
                    </div>
                    <div>
                      <span class="text-sm font-medium text-gray-500">Email:</span>
                      <p class="text-gray-900">{{ form.email }}</p>
                    </div>
                  </div>

                  <div class="space-y-3">
                    <div>
                      <span class="text-sm font-medium text-gray-500">Date:</span>
                      <p class="text-gray-900">{{ new Date(form.preferred_date).toLocaleDateString('en-US', {
                        weekday: 'long',
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      }) }}</p>
                    </div>
                    <div>
                      <span class="text-sm font-medium text-gray-500">Time:</span>
                      <p class="text-gray-900">{{ formatTimeDisplay(parseInt(form.preferred_time.split(':')[0]),
                        parseInt(form.preferred_time.split(':')[1])) }}</p>
                    </div>
                    <div v-if="form.branch_id">
                      <span class="text-sm font-medium text-gray-500">Location:</span>
                      <p class="text-gray-900">{{ clinic.branches?.find(b => b.id == form.branch_id)?.name }}</p>
                    </div>
                  </div>
                </div>

                <div class="border-t border-gray-200 pt-4 mt-4">
                  <div>
                    <span class="text-sm font-medium text-gray-500">Reason for visit:</span>
                    <p class="text-gray-900 mt-1">{{ form.symptoms }}</p>
                  </div>
                  <div v-if="form.additional_notes" class="mt-3">
                    <span class="text-sm font-medium text-gray-500">Additional notes:</span>
                    <p class="text-gray-900 mt-1">{{ form.additional_notes }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Navigation Controls -->
          <div class="pt-6 md:pt-8 border-t border-gray-100">
            <div class="flex flex-col gap-3 md:gap-4">
              <!-- Step Navigation Buttons -->
              <div class="flex justify-between items-center">
                <button type="button" @click="previousStep" :disabled="isFirstStep"
                        class="flex items-center space-x-2 px-6 py-3 bg-gray-100 text-gray-700 rounded-xl font-semibold hover:bg-gray-200 focus:ring-4 focus:ring-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                  <span>Previous</span>
                </button>

                <div class="text-sm text-gray-500">
                  Step {{ currentStep }} of {{ totalSteps }}
                </div>

                <button v-if="!isLastStep" type="button" @click="nextStep" :disabled="!canGoNext"
                        class="flex items-center space-x-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-xl font-semibold hover:from-blue-700 hover:to-indigo-700 focus:ring-4 focus:ring-blue-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200">
                  <span>Next</span>
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>

                <!-- Submit Button for Last Step -->
                <button v-if="isLastStep" type="submit" :disabled="isSubmitting"
                        class="flex items-center space-x-2 px-6 py-3 bg-gradient-to-r from-green-600 to-emerald-600 text-white rounded-xl font-semibold hover:from-green-700 hover:to-emerald-700 focus:ring-4 focus:ring-green-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200">
                  <svg v-if="isSubmitting" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                  <span>{{ isSubmitting ? 'Submitting...' : 'Submit Request' }}</span>
                </button>
              </div>
            </div>

            <!-- Trust indicators -->
            <div class="mt-6 text-center">
              <div class="flex items-center justify-center space-x-6 text-sm text-gray-500">
                <div class="flex items-center space-x-2">
                  <svg class="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clip-rule="evenodd" />
                  </svg>
                  <span>Secure & Private</span>
                </div>
                <div class="flex items-center space-x-2">
                  <svg class="w-4 h-4 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                      clip-rule="evenodd" />
                  </svg>
                  <span>24h Response</span>
                </div>
                <div class="flex items-center space-x-2">
                  <svg class="w-4 h-4 text-purple-500" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>Expert Care</span>
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="showSuccess" class="mt-8 bg-gradient-to-r from-green-50 to-emerald-50 border-2 border-green-200 rounded-3xl p-8 shadow-xl">
      <div class="text-center">
        <div class="w-20 h-20 bg-gradient-to-br from-green-400 to-emerald-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 class="text-2xl font-bold text-green-800 mb-3">Request Submitted Successfully!</h3>
        <p class="text-green-700 text-lg leading-relaxed mb-6">
          Thank you for choosing {{ clinic.name }}. We'll review your information and contact you within 24 hours to schedule your appointment.
        </p>
        <div class="bg-white/50 rounded-xl p-4 inline-block">
          <p class="text-sm text-green-600 font-medium">
            📧 Check your email for confirmation details
          </p>
        </div>
      </div>
    </div>

    <!-- Contact Information -->
    <div class="mt-12 bg-gradient-to-r from-gray-50 to-blue-50 rounded-3xl p-8 border border-gray-100">
      <div class="text-center mb-8">
        <h3 class="text-2xl font-bold text-gray-900 mb-3">Need to speak with us directly?</h3>
        <p class="text-gray-600">Our friendly team is here to help you with any questions</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white rounded-2xl p-6 shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
          <div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-500 rounded-xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-900 mb-2 text-center">Call Us</h4>
          <p class="text-gray-600 text-center font-medium">{{ clinic.phone }}</p>
          <p class="text-sm text-gray-500 text-center mt-2">Mon-Fri: 8AM-6PM</p>
        </div>

        <div class="bg-white rounded-2xl p-6 shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
          <div class="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-500 rounded-xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-900 mb-2 text-center">Email Us</h4>
          <p class="text-gray-600 text-center font-medium">{{ clinic.email }}</p>
          <p class="text-sm text-gray-500 text-center mt-2">We'll respond within 2 hours</p>
        </div>

        <div class="bg-white rounded-2xl p-6 shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
          <div class="w-12 h-12 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-900 mb-2 text-center">Visit Us</h4>
          <p class="text-gray-600 text-center font-medium text-sm">{{ clinic.address }}</p>
          <p class="text-sm text-gray-500 text-center mt-2">Open daily</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import apiService from '../services/api'

const route = useRoute()

// Reactive data
const loading = ref(true)
const error = ref('')
const clinic = ref({})
const showSuccess = ref(false)
const isSubmitting = ref(false)
const loadingTimeSlots = ref(false)
const availableTimeSlots = ref([])

// Phone check state
const checkingPhone = ref(false)
const patientFound = ref(false)
const existingPatient = ref(null)
const phoneChecked = ref(false)

// Step management
const currentStep = ref(1)
const totalSteps = 6
const completedSteps = ref([])

// Step definitions
const steps = [
  { id: 1, title: 'Phone Verification', icon: 'phone', description: 'Verify your phone number' },
  { id: 2, title: 'Personal Information', icon: 'user', description: 'Your basic details' },
  { id: 3, title: 'Select Date', icon: 'calendar', description: 'Choose your preferred date' },
  { id: 4, title: 'Select Time', icon: 'clock', description: 'Pick available time slot' },
  { id: 5, title: 'Additional Details', icon: 'clipboard', description: 'Tell us about your visit' },
  { id: 6, title: 'Review & Submit', icon: 'check-circle', description: 'Confirm your appointment' }
]

const form = ref({
  phone: '',
  first_name: '',
  last_name: '',
  email: '',
  preferred_date: '',
  preferred_time: '',
  symptoms: '',
  additional_notes: '',
  branch_id: null
})

// Computed properties
const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

const allTimeSlots = computed(() => {
  const slots = []
  // Generate time slots from 9:00 AM to 5:00 PM in 30-minute intervals
  for (let hour = 9; hour < 17; hour++) {
    for (let minute = 0; minute < 60; minute += 30) {
      const timeString = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
      const displayTime = formatTimeDisplay(hour, minute)

      // Check if this slot is available
      const isAvailable = availableTimeSlots.value.some(slot => slot.time === timeString)

      slots.push({
        time: timeString,
        display: displayTime,
        available: isAvailable
      })
    }
  }
  return slots
})

// Step navigation computed properties
const currentStepData = computed(() => steps.find(step => step.id === currentStep.value))
const isFirstStep = computed(() => currentStep.value === 1)
const isLastStep = computed(() => currentStep.value === totalSteps)
const canGoNext = computed(() => {
  switch (currentStep.value) {
    case 1: return phoneChecked.value
    case 2: return form.value.first_name && form.value.last_name && form.value.email
    case 3: return form.value.preferred_date
    case 4: return form.value.preferred_time
    case 5: return form.value.symptoms
    case 6: return true
    default: return false
  }
})
const canGoPrevious = computed(() => currentStep.value > 1)

const formatTimeDisplay = (hour, minute) => {
  const period = hour >= 12 ? 'PM' : 'AM'
  const displayHour = hour > 12 ? hour - 12 : hour === 0 ? 12 : hour
  return `${displayHour}:${minute.toString().padStart(2, '0')} ${period}`
}

// Check if patient exists by phone number
const checkPatientByPhone = async () => {
  if (!form.value.phone) {
    error.value = 'Please enter a phone number'
    return
  }

  try {
    checkingPhone.value = true
    error.value = ''

    const clinicIdentifier = route.params.clinicIdentifier
    const result = await apiService.request('get', `/api/public/patient/${clinicIdentifier}`, null, {
      params: { phone: form.value.phone }
    })

    if (result.success) {
      phoneChecked.value = true
      if (result.data.found) {
        patientFound.value = true
        existingPatient.value = result.data.patient
        // Auto-fill the form with existing patient data
        form.value.first_name = result.data.patient.first_name
        form.value.last_name = result.data.patient.last_name
        form.value.email = result.data.patient.email
      } else {
        patientFound.value = false
        existingPatient.value = null
      }
    } else {
      error.value = result.error || 'Failed to check phone number'
    }
  } catch (err) {
    console.error('Error checking phone:', err)
    error.value = 'Failed to check phone number'
  } finally {
    checkingPhone.value = false
  }
}

// Reset phone check when phone changes
const resetPhoneCheck = () => {
  if (phoneChecked.value) {
    phoneChecked.value = false
    patientFound.value = false
    existingPatient.value = null
    // Clear form fields except phone
    form.value.first_name = ''
    form.value.last_name = ''
    form.value.email = ''
  }
}

// Step navigation methods
const nextStep = () => {
  if (canGoNext.value && !isLastStep.value) {
    completedSteps.value.push(currentStep.value)
    currentStep.value++
  }
}

const previousStep = () => {
  if (canGoPrevious.value) {
    currentStep.value--
  }
}

const goToStep = (stepId) => {
  if (stepId <= currentStep.value || completedSteps.value.includes(stepId - 1)) {
    currentStep.value = stepId
  }
}

// Methods
const loadClinicInfo = async () => {
  try {
    loading.value = true
    error.value = ''

    const clinicIdentifier = route.params.clinicIdentifier
    const result = await apiService.getClinicInfo(clinicIdentifier)

    if (result.success) {
      clinic.value = result.data.clinic

      // Set default branch if only one exists
      if (clinic.value.branches && clinic.value.branches.length === 1) {
        form.value.branch_id = clinic.value.branches[0].id
      }
    } else {
      error.value = result.error || 'Failed to load clinic information'
    }
  } catch (err) {
    console.error('Error loading clinic info:', err)
    error.value = 'Failed to load clinic information'
  } finally {
    loading.value = false
  }
}

const fetchAvailableTimeSlots = async () => {
  if (!form.value.preferred_date) return

  try {
    loadingTimeSlots.value = true
    const clinicIdentifier = route.params.clinicIdentifier
    const result = await apiService.request('get', `/api/public/timeslots/${clinicIdentifier}`, null, {
      params: {
        date: form.value.preferred_date,
        branch_id: form.value.branch_id || ''
      }
    })

    if (result.success) {
      availableTimeSlots.value = result.data.available_slots || []
    } else {
      availableTimeSlots.value = []
    }
  } catch (err) {
    console.error('Error fetching time slots:', err)
    availableTimeSlots.value = []
  } finally {
    loadingTimeSlots.value = false
  }
}

const handleSubmit = async () => {
  // Validate that phone has been checked
  if (!phoneChecked.value) {
    error.value = 'Please verify your phone number first'
    return
  }

  try {
    isSubmitting.value = true
    error.value = ''

    const clinicIdentifier = route.params.clinicIdentifier
    const submitData = { ...form.value }

    // Convert branch_id to number if it exists
    if (submitData.branch_id) {
      submitData.branch_id = parseInt(submitData.branch_id)
    }

    const result = await apiService.createPatientSelfSchedule(clinicIdentifier, submitData)

    if (result.success) {
      showSuccess.value = true

      // Scroll to success message
      setTimeout(() => {
        const successElement = document.querySelector('.bg-green-50')
        if (successElement) {
          successElement.scrollIntoView({ behavior: 'smooth' })
        }
      }, 100)
    } else {
      error.value = result.error || 'Failed to submit appointment request'
    }

  } catch (err) {
    console.error('Error submitting request:', err)
    error.value = 'Failed to submit appointment request'
  } finally {
    isSubmitting.value = false
  }
}

// Watchers
watch(() => form.value.branch_id, () => {
  if (form.value.preferred_date) {
    fetchAvailableTimeSlots()
  }
})

watch(() => form.value.phone, () => {
  resetPhoneCheck()
})

// Lifecycle
onMounted(() => {
  loadClinicInfo()
})
</script>

<style scoped>
/* Enhanced form styling */
input:focus,
select:focus,
textarea:focus {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

/* Custom scrollbar for textareas */
textarea::-webkit-scrollbar {
  width: 6px;
}

textarea::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

textarea::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #cbd5e1, #94a3b8);
  border-radius: 3px;
}

textarea::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #94a3b8, #64748b);
}

/* Enhanced animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.7;
  }
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0px);
  }

  50% {
    transform: translateY(-10px);
  }
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

.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out;
}

.animate-pulse-slow {
  animation: pulse 3s ease-in-out infinite;
}

.animate-float {
  animation: float 6s ease-in-out infinite;
}

/* Gradient text effect */
.gradient-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Glass morphism effect */
.glass {
  backdrop-filter: blur(16px) saturate(180%);
  background-color: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(209, 213, 219, 0.3);
}

/* Enhanced hover effects */
.hover-lift {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.hover-lift:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

/* Custom radio button styling */
input[type="radio"] {
  transform: scale(1.2);
  accent-color: #3b82f6;
}

/* Focus ring improvements */
.focus-ring:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .container {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .gradient-text {
    font-size: 2.5rem;
  }

  /* Mobile-specific form improvements */
  .mobile-form-spacing {
    margin-bottom: 1rem;
  }

  /* Better touch targets on mobile */
  button,
  .cursor-pointer {
    min-height: 44px;
  }

  /* Improve readability on small screens */
  .mobile-text {
    font-size: 0.875rem;
    line-height: 1.25rem;
  }

  /* Better spacing for mobile grids */
  .mobile-grid {
    gap: 0.5rem;
  }
}

/* Tablet adjustments */
@media (min-width: 768px) and (max-width: 1024px) {
  .tablet-spacing {
    padding: 1.5rem;
  }

  .tablet-text {
    font-size: 1rem;
  }
}

/* Desktop enhancements */
@media (min-width: 1024px) {
  .desktop-shadow {
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }

  .desktop-hover {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .desktop-hover:hover {
    transform: translateY(-2px);
  }
}

/* Smooth transitions for all interactive elements */
button,
input,
select,
textarea,
label {
  transition: all 0.2s ease-in-out;
}

/* Custom placeholder styling */
::placeholder {
  color: #9ca3af;
  opacity: 0.8;
}

/* Enhanced shadow effects */
.shadow-soft {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.shadow-glow {
  box-shadow: 0 0 30px rgba(59, 130, 246, 0.2);
}

/* Background pattern animation */
@keyframes backgroundMove {
  0% {
    transform: translateX(0) translateY(0);
  }

  33% {
    transform: translateX(30px) translateY(-30px);
  }

  66% {
    transform: translateX(-20px) translateY(20px);
  }

  100% {
    transform: translateX(0) translateY(0);
  }
}

.animate-background {
  animation: backgroundMove 20s ease-in-out infinite;
}
</style>