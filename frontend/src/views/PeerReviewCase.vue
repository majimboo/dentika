<template>
  <div class="peer-review-case" v-if="caseData">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6 mb-6">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">{{ caseData.title }}</h1>
            <p class="text-gray-600">Case #{{ caseData.case_number }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span
              :class="getStatusBadgeClass(caseData.status)"
              class="px-3 py-1 text-sm font-medium rounded-full"
            >
              {{ caseData.status }}
            </span>
            <span
              :class="getVisibilityBadgeClass(caseData.visibility)"
              class="px-3 py-1 text-sm font-medium rounded-full"
            >
              {{ getVisibilityLabel(caseData.visibility) }}
            </span>
          </div>
        </div>

        <p class="text-gray-700 mb-4">{{ caseData.description }}</p>

        <div class="flex justify-between items-center text-sm text-gray-500">
          <span>Created by {{ caseData.created_by.first_name }} {{ caseData.created_by.last_name }}
            <span v-if="isCaseMaker" class="ml-2 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">You</span>
          </span>
          <span>{{ formatDate(caseData.created_at) }}</span>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Patient & Appointment Information -->
          <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Patient & Appointment Information</h2>
            <div v-if="patientData">
              <div class="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <span class="text-sm font-medium text-gray-500">Age</span>
                  <p class="text-gray-900">{{ patientData.patient_age || 'Not specified' }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Gender</span>
                  <p class="text-gray-900">{{ patientData.patient_gender || 'Not specified' }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Blood Type</span>
                  <p class="text-gray-900">{{ patientData.patient_blood_type || 'Not specified' }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Allergies</span>
                  <p class="text-gray-900">{{ patientData.allergies || 'None specified' }}</p>
                </div>
              </div>

              <div class="space-y-3">
                <div v-if="patientData.appointment_date">
                  <span class="text-sm font-medium text-gray-500">Appointment Date</span>
                  <p class="text-gray-900 mt-1">{{ formatDate(patientData.appointment_date) }}</p>
                </div>
                <div v-if="patientData.title">
                  <span class="text-sm font-medium text-gray-500">Appointment Title</span>
                  <p class="text-gray-900 mt-1">{{ patientData.title }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Chief Complaint</span>
                  <p class="text-gray-900 mt-1">{{ patientData.chief_complaint || 'Not specified' }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Medical History</span>
                  <p class="text-gray-900 mt-1">{{ patientData.medical_history || 'Not specified' }}</p>
                </div>
                <div>
                  <span class="text-sm font-medium text-gray-500">Current Medications</span>
                  <p class="text-gray-900 mt-1">{{ patientData.current_medications || 'None specified' }}</p>
                </div>
                <div v-if="patientData.pre_notes">
                  <span class="text-sm font-medium text-gray-500">Pre-Appointment Notes</span>
                  <p class="text-gray-900 mt-1">{{ patientData.pre_notes }}</p>
                </div>
                <div v-if="patientData.post_notes">
                  <span class="text-sm font-medium text-gray-500">Post-Appointment Notes</span>
                  <p class="text-gray-900 mt-1">{{ patientData.post_notes }}</p>
                </div>
              </div>
            </div>
            <div v-else class="text-gray-500">
              Loading patient information...
            </div>
          </div>


          <!-- Comments Section -->
          <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-semibold text-gray-900">Discussion</h2>
              <span class="text-sm text-gray-500">{{ caseData.comments?.length || 0 }} comments</span>
            </div>

            <!-- Add Comment -->
            <div class="mb-6">
              <textarea
                v-model="newComment"
                rows="3"
                placeholder="Add your thoughts, recommendations, or questions..."
                class="w-full border border-neutral-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              ></textarea>
              <div class="flex justify-between items-center mt-2">
                <select
                  v-model="commentType"
                  class="border border-neutral-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="comment">Comment</option>
                  <option value="recommendation">Recommendation</option>
                  <option value="question">Question</option>
                  <option value="diagnosis">Diagnosis</option>
                </select>
                <button
                  @click="addComment"
                  :disabled="!newComment.trim() || loading"
                  class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ loading ? 'Posting...' : 'Post Comment' }}
                </button>
              </div>
            </div>

            <!-- Comments List -->
            <div class="space-y-6">
              <!-- Empty state -->
              <div v-if="!caseData.comments || caseData.comments.length === 0" class="text-center py-8">
                <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
                  <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
                  </svg>
                </div>
                <p class="text-gray-500 font-medium">No comments yet</p>
                <p class="text-gray-400 text-sm">Be the first to share your thoughts on this case</p>
              </div>

              <div
                v-for="comment in caseData.comments"
                :key="comment.id"
                class="flex gap-3"
              >
                <!-- Avatar -->
                <div class="flex-shrink-0">
                  <div class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                    <img
                      v-if="comment.author.avatar_path"
                      :src="`/uploads/${comment.author.avatar_path}`"
                      :alt="`${comment.author.first_name} ${comment.author.last_name}`"
                      class="w-full h-full object-cover"
                    />
                    <span
                      v-else
                      class="text-sm font-medium text-gray-600"
                    >
                      {{ getInitials(comment.author.first_name, comment.author.last_name) }}
                    </span>
                  </div>
                </div>

                <!-- Comment Content -->
                <div class="flex-1 min-w-0">
                  <!-- Header with author info -->
                  <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-2">
                      <span class="font-semibold text-gray-900">
                        {{ comment.author.first_name }} {{ comment.author.last_name }}
                      </span>
                      <span
                        :class="getCommentTypeBadgeClass(comment.comment_type)"
                        class="px-2 py-1 text-xs font-medium rounded-full"
                      >
                        {{ comment.comment_type }}
                      </span>
                    </div>
                    <span class="text-sm text-gray-500 flex-shrink-0">
                      {{ formatDateWithAMPM(comment.created_at) }}
                    </span>
                  </div>

                  <!-- Comment text -->
                  <div class="bg-gray-50 rounded-lg px-4 py-3 mb-3">
                    <p class="text-gray-800 whitespace-pre-wrap leading-relaxed">{{ comment.content }}</p>
                  </div>

                  <!-- Reply button -->
                  <div class="flex items-center gap-2 mb-3">
                    <button
                      @click="toggleReplyForm(comment.id)"
                      class="text-sm text-blue-600 hover:text-blue-800 font-medium"
                    >
                      Reply
                    </button>
                    <span v-if="comment.replies && comment.replies.length > 0" class="text-sm text-gray-500">
                      {{ comment.replies.length }} {{ comment.replies.length === 1 ? 'reply' : 'replies' }}
                    </span>
                  </div>

                  <!-- Reply form -->
                  <div v-if="showReplyForm === comment.id" class="mb-4 ml-4">
                    <textarea
                      v-model="replyText"
                      rows="2"
                      placeholder="Write a reply..."
                      class="w-full border border-neutral-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    ></textarea>
                    <div class="flex justify-end gap-2 mt-2">
                      <button
                        @click="cancelReply"
                        class="px-3 py-1 text-sm text-gray-600 hover:text-gray-800"
                      >
                        Cancel
                      </button>
                      <button
                        @click="submitReply(comment.id)"
                        :disabled="!replyText.trim() || loading"
                        class="px-3 py-1 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {{ loading ? 'Posting...' : 'Reply' }}
                      </button>
                    </div>
                  </div>

                  <!-- Replies -->
                  <div v-if="comment.replies && comment.replies.length > 0" class="ml-4 space-y-3">
                    <div
                      v-for="reply in comment.replies"
                      :key="reply.id"
                      class="flex gap-3"
                    >
                      <!-- Reply Avatar (smaller) -->
                      <div class="flex-shrink-0">
                        <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                          <img
                            v-if="reply.author.avatar_path"
                            :src="`/uploads/${reply.author.avatar_path}`"
                            :alt="`${reply.author.first_name} ${reply.author.last_name}`"
                            class="w-full h-full object-cover"
                          />
                          <span
                            v-else
                            class="text-xs font-medium text-gray-600"
                          >
                            {{ getInitials(reply.author.first_name, reply.author.last_name) }}
                          </span>
                        </div>
                      </div>

                      <!-- Reply Content -->
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center justify-between mb-1">
                          <span class="font-medium text-gray-900 text-sm">
                            {{ reply.author.first_name }} {{ reply.author.last_name }}
                          </span>
                          <span class="text-xs text-gray-500 flex-shrink-0">
                            {{ formatDateWithAMPM(reply.created_at) }}
                          </span>
                        </div>
                        <div class="bg-white border border-neutral-200 rounded-lg px-3 py-2">
                          <p class="text-gray-700 text-sm whitespace-pre-wrap leading-relaxed">{{ reply.content }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Case Actions -->
          <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Actions</h3>
            <div class="space-y-2">
              <button
                v-if="canEditCase"
                @click="updateStatus"
                class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm"
              >
                Update Status
              </button>
              <button
                @click="shareCase"
                class="w-full bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-md text-sm"
              >
                Share Case
              </button>
            </div>
          </div>

          <!-- Participants -->
          <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Participants</h3>
            <div class="space-y-2">
              <div
                v-for="participant in caseData.participants"
                :key="participant.id"
                class="flex items-center justify-between"
              >
                <span class="text-sm text-gray-700">
                  {{ participant.user.first_name }} {{ participant.user.last_name }}
                </span>
                <span
                  :class="getPermissionBadgeClass(participant.permission)"
                  class="px-2 py-1 text-xs font-medium rounded-full"
                >
                  {{ participant.permission }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  <div v-else class="flex justify-center items-center min-h-screen">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Loading case...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import apiService from '../services/api.js'

const route = useRoute()

const caseData = ref(null)
const isCaseMaker = ref(false)
const patientData = ref(null)
const newComment = ref('')
const commentType = ref('comment')
const loading = ref(false)
const showReplyForm = ref(null)
const replyText = ref('')

const fetchCase = async () => {
  try {
    const result = await apiService.get(`/api/peer-review/cases/${route.params.id}`)

    if (result.success) {
      const responseData = result.data
      caseData.value = responseData.case || responseData // Handle both response formats
      isCaseMaker.value = responseData.is_case_maker || false
      patientData.value = responseData.patient_data || null
    } else if (result.status === 404) {
      // Handle case not found
      console.error('Case not found')
    } else {
      console.error('Failed to fetch case')
    }
  } catch (error) {
    console.error('Failed to fetch case:', error)
  }
}

const addComment = async () => {
  if (!newComment.value.trim()) return

  loading.value = true

  try {
    const result = await apiService.post(`/api/peer-review/cases/${route.params.id}/comments`, {
      content: newComment.value,
      comment_type: commentType.value
    })

    if (result.success) {
      newComment.value = ''
      commentType.value = 'comment'
      await fetchCase() // Refresh comments
    } else {
      alert(result.error || 'Failed to add comment')
    }
  } catch (error) {
    console.error('Failed to add comment:', error)
    alert('Failed to add comment')
  } finally {
    loading.value = false
  }
}

const toggleReplyForm = (commentId) => {
  if (showReplyForm.value === commentId) {
    showReplyForm.value = null
    replyText.value = ''
  } else {
    showReplyForm.value = commentId
    replyText.value = ''
  }
}

const cancelReply = () => {
  showReplyForm.value = null
  replyText.value = ''
}

const submitReply = async (parentId) => {
  if (!replyText.value.trim()) return

  loading.value = true

  try {
    const result = await apiService.post(`/api/peer-review/cases/${route.params.id}/comments`, {
      content: replyText.value,
      comment_type: 'comment',
      parent_id: parentId
    })

    if (result.success) {
      replyText.value = ''
      showReplyForm.value = null
      await fetchCase() // Refresh comments
    } else {
      alert(result.error || 'Failed to add reply')
    }
  } catch (error) {
    console.error('Failed to add reply:', error)
    alert('Failed to add reply')
  } finally {
    loading.value = false
  }
}

const updateStatus = async () => {
  const newStatus = prompt('Enter new status (open, closed, resolved):', caseData.value.status)
  if (!newStatus || !['open', 'closed', 'resolved'].includes(newStatus)) {
    return
  }

  try {
    const result = await apiService.put(`/api/peer-review/cases/${route.params.id}/status`, {
      status: newStatus
    })

    if (result.success) {
      await fetchCase() // Refresh case data
    } else {
      alert(result.error || 'Failed to update status')
    }
  } catch (error) {
    console.error('Failed to update status:', error)
    alert('Failed to update status')
  }
}

const shareCase = () => {
  const shareUrl = `${window.location.origin}/peer-review/${caseData.value.id}`
  navigator.clipboard.writeText(shareUrl).then(() => {
    alert('Case URL copied to clipboard!')
  })
}

const canEditCase = () => {
  return caseData.value?.participants?.some(p =>
    p.user_id === caseData.value.created_by_id && p.permission === 'edit'
  )
}

const getStatusBadgeClass = (status) => {
  const classes = {
    open: 'bg-green-100 text-green-800',
    closed: 'bg-gray-100 text-gray-800',
    resolved: 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getVisibilityBadgeClass = (visibility) => {
  const classes = {
    public: 'bg-purple-100 text-purple-800',
    in_clinic: 'bg-orange-100 text-orange-800',
    invite_only: 'bg-red-100 text-red-800'
  }
  return classes[visibility] || 'bg-gray-100 text-gray-800'
}

const getVisibilityLabel = (visibility) => {
  const labels = {
    public: 'Public',
    in_clinic: 'In-Clinic',
    invite_only: 'Invite Only'
  }
  return labels[visibility] || visibility
}

const getCommentTypeBadgeClass = (type) => {
  const classes = {
    comment: 'bg-gray-100 text-gray-800',
    recommendation: 'bg-blue-100 text-blue-800',
    question: 'bg-yellow-100 text-yellow-800',
    diagnosis: 'bg-green-100 text-green-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getPermissionBadgeClass = (permission) => {
  const classes = {
    view: 'bg-gray-100 text-gray-800',
    comment: 'bg-blue-100 text-blue-800',
    edit: 'bg-green-100 text-green-800'
  }
  return classes[permission] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const formatDateWithAMPM = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  })
}

const getInitials = (firstName, lastName) => {
  const first = firstName ? firstName.charAt(0).toUpperCase() : ''
  const last = lastName ? lastName.charAt(0).toUpperCase() : ''
  return first + last
}

const formatDentalChart = (chartData) => {
  try {
    const parsed = JSON.parse(chartData)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return chartData
  }
}

onMounted(() => {
  fetchCase()
})
</script>