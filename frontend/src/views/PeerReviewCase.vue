<template>
  <div class="peer-review-case" v-if="caseData">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
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
          <span>Created by {{ caseData.created_by.first_name }} {{ caseData.created_by.last_name }}</span>
          <span>{{ formatDate(caseData.created_at) }}</span>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Patient Information -->
          <div class="bg-white rounded-lg shadow-sm border p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Patient Information</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <span class="text-sm font-medium text-gray-500">Age</span>
                <p class="text-gray-900">{{ caseData.patient_age || 'Not specified' }}</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Gender</span>
                <p class="text-gray-900">{{ caseData.patient_gender || 'Not specified' }}</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Blood Type</span>
                <p class="text-gray-900">{{ caseData.patient_blood_type || 'Not specified' }}</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Allergies</span>
                <p class="text-gray-900">{{ caseData.allergies || 'None specified' }}</p>
              </div>
            </div>

            <div class="mt-4 space-y-3">
              <div>
                <span class="text-sm font-medium text-gray-500">Chief Complaint</span>
                <p class="text-gray-900 mt-1">{{ caseData.chief_complaint || 'Not specified' }}</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Medical History</span>
                <p class="text-gray-900 mt-1">{{ caseData.medical_history || 'Not specified' }}</p>
              </div>
              <div>
                <span class="text-sm font-medium text-gray-500">Current Medications</span>
                <p class="text-gray-900 mt-1">{{ caseData.current_medications || 'None specified' }}</p>
              </div>
            </div>
          </div>

          <!-- Dental Chart -->
          <div class="bg-white rounded-lg shadow-sm border p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Dental Chart</h2>
            <div v-if="caseData.dental_chart_data" class="bg-gray-50 rounded-lg p-4">
              <pre class="text-sm text-gray-800 whitespace-pre-wrap">{{ formatDentalChart(caseData.dental_chart_data) }}</pre>
            </div>
            <p v-else class="text-gray-500">No dental chart data available</p>
          </div>

          <!-- Comments Section -->
          <div class="bg-white rounded-lg shadow-sm border p-6">
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
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              ></textarea>
              <div class="flex justify-between items-center mt-2">
                <select
                  v-model="commentType"
                  class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
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
            <div class="space-y-4">
              <div
                v-for="comment in caseData.comments"
                :key="comment.id"
                class="border-l-4 border-blue-200 pl-4"
              >
                <div class="flex justify-between items-start mb-2">
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-gray-900">
                      {{ comment.author.first_name }} {{ comment.author.last_name }}
                    </span>
                    <span
                      :class="getCommentTypeBadgeClass(comment.comment_type)"
                      class="px-2 py-1 text-xs font-medium rounded-full"
                    >
                      {{ comment.comment_type }}
                    </span>
                  </div>
                  <span class="text-sm text-gray-500">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="text-gray-700 whitespace-pre-wrap">{{ comment.content }}</p>

                <!-- Replies -->
                <div v-if="comment.replies && comment.replies.length > 0" class="mt-3 ml-6 space-y-2">
                  <div
                    v-for="reply in comment.replies"
                    :key="reply.id"
                    class="bg-gray-50 rounded-lg p-3"
                  >
                    <div class="flex justify-between items-start mb-1">
                      <span class="font-medium text-gray-900 text-sm">
                        {{ reply.author.first_name }} {{ reply.author.last_name }}
                      </span>
                      <span class="text-xs text-gray-500">{{ formatDate(reply.created_at) }}</span>
                    </div>
                    <p class="text-gray-700 text-sm whitespace-pre-wrap">{{ reply.content }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Case Actions -->
          <div class="bg-white rounded-lg shadow-sm border p-6">
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
          <div class="bg-white rounded-lg shadow-sm border p-6">
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

const route = useRoute()

const caseData = ref(null)
const newComment = ref('')
const commentType = ref('comment')
const loading = ref(false)

const fetchCase = async () => {
  try {
    const response = await fetch(`/api/peer-review/cases/${route.params.id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      caseData.value = await response.json()
    } else if (response.status === 404) {
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
    const response = await fetch(`/api/peer-review/cases/${route.params.id}/comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        content: newComment.value,
        comment_type: commentType.value
      })
    })

    if (response.ok) {
      newComment.value = ''
      commentType.value = 'comment'
      await fetchCase() // Refresh comments
    } else {
      const error = await response.json()
      alert(error.error || 'Failed to add comment')
    }
  } catch (error) {
    console.error('Failed to add comment:', error)
    alert('Failed to add comment')
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
    const response = await fetch(`/api/peer-review/cases/${route.params.id}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        status: newStatus
      })
    })

    if (response.ok) {
      await fetchCase() // Refresh case data
    } else {
      const error = await response.json()
      alert(error.error || 'Failed to update status')
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