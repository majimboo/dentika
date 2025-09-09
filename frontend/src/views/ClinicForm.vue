<template>
  <div class="space-y-8">
    <h1 class="text-2xl font-bold">{{ isEditMode ? 'Edit Clinic' : 'Create New Clinic' }}</h1>

    <!-- Clinic Details Form -->
    <div class="bg-white rounded-2xl shadow-lg border p-6">
      <h2 class="text-xl font-semibold mb-4">Clinic Details</h2>
      <form @submit.prevent="saveClinic" class="space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-neutral-700">Clinic Name</label>
          <input type="text" id="name" v-model="clinic.name" class="mt-1 block w-full border-neutral-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500" required>
        </div>
        <div>
          <label for="address" class="block text-sm font-medium text-neutral-700">Address</label>
          <input type="text" id="address" v-model="clinic.address" class="mt-1 block w-full border-neutral-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500">
        </div>
        <div>
          <label for="phone" class="block text-sm font-medium text-neutral-700">Phone</label>
          <input type="text" id="phone" v-model="clinic.phone" class="mt-1 block w-full border-neutral-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500">
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-neutral-700">Email</label>
          <input type="email" id="email" v-model="clinic.email" class="mt-1 block w-full border-neutral-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500">
        </div>
        <div>
          <label for="website" class="block text-sm font-medium text-neutral-700">Website</label>
          <input type="url" id="website" v-model="clinic.website" class="mt-1 block w-full border-neutral-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500">
        </div>
        <div class="flex justify-end space-x-4">
          <button type="button" @click="cancel" class="px-4 py-2 border rounded-lg">Cancel</button>
          <button type="submit" :disabled="saving" class="px-4 py-2 border border-transparent rounded-lg text-white bg-primary-600 hover:bg-primary-700 disabled:opacity-50">
            {{ saving ? 'Saving...' : 'Save Clinic' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Branch Management -->
    <div v-if="isEditMode" class="bg-white rounded-2xl shadow-lg border p-6">
      <h2 class="text-xl font-semibold mb-4">Branches</h2>
      <ul class="space-y-2">
        <li v-for="branch in clinic.branches" :key="branch.id" class="p-4 border rounded-md flex justify-between items-center">
          <div>
            <p class="font-semibold">{{ branch.name }} <span v-if="branch.is_main_branch" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">Main</span></p>
            <p class="text-sm text-neutral-600">{{ branch.address }}</p>
          </div>
          <div class="space-x-2">
            <button class="text-sm text-blue-600">Edit</button>
            <button @click="deleteBranch(branch.id)" class="text-sm text-red-600">Delete</button>
          </div>
        </li>
      </ul>
      <button class="mt-4 px-4 py-2 border rounded-lg">Add Branch</button>
    </div>
  <!-- Staff Management -->
    <div v-if="isEditMode" class="bg-white rounded-2xl shadow-lg border p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold">Staff</h2>
        <button @click="addUserToClinic" class="px-4 py-2 border rounded-lg text-sm bg-primary-600 text-white hover:bg-primary-700">Add User</button>
      </div>
      <ul class="space-y-2">
        <li v-for="user in clinic.staff" :key="user.id" class="p-4 border rounded-md flex justify-between items-center">
          <div>
            <p class="font-semibold">{{ user.username }} <span class="text-xs bg-gray-200 text-gray-800 px-2 py-1 rounded-full">{{ user.role }}</span></p>
            <p class="text-sm text-neutral-600">{{ user.email }}</p>
          </div>
          <div class="space-x-2">
            <button @click="editUser(user.id)" class="text-sm text-blue-600">Edit</button>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClinicStore } from '../stores/clinic'

export default {
  name: 'ClinicForm',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const clinicStore = useClinicStore()

    const clinic = ref({ name: '', address: '', phone: '', email: '', website: '', branches: [], staff: [] })
    const saving = ref(false)

    const clinicId = route.params.id
    const isEditMode = computed(() => !!clinicId)

    const loadClinicData = async () => {
      if (isEditMode.value) {
        await clinicStore.fetchClinic(clinicId)
        clinic.value = { ...clinicStore.currentClinic }
      }
    }

    const saveClinic = async () => {
      saving.value = true
      let result
      if (isEditMode.value) {
        result = await clinicStore.updateClinic(clinicId, clinic.value)
      } else {
        result = await clinicStore.createClinic(clinic.value)
      }

      if (result.success) {
        router.push('/clinics')
      } else {
        // Handle error
        console.error(result.error)
      }
      saving.value = false
    }

    const deleteBranch = async (branchId) => {
        if(confirm('Are you sure you want to delete this branch?')){
            const result = await clinicStore.deleteBranch(clinicId, branchId);
            if(result.success){
                loadClinicData(); // Refresh data
            }
        }
    }

    const cancel = () => {
      router.push('/clinics')
    }

    const addUserToClinic = () => {
      router.push(`/users/new?clinic_id=${clinicId}`)
    }

    const editUser = (userId) => {
      router.push(`/users/${userId}/edit`)
    }

    onMounted(() => {
      loadClinicData()
    })

    return {
      clinic,
      isEditMode,
      saving,
      saveClinic,
      deleteBranch,
      cancel,
      addUserToClinic,
      editUser,
    }
  }
}
</script>
