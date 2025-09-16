<template>
  <div class="space-y-2">
    <div v-if="!parsedSchedule" class="text-sm text-neutral-500">
      No schedule configured
    </div>
    <div v-else>
      <div v-for="day in daysOfWeek" :key="day.key" class="flex justify-between items-center text-sm">
        <span class="font-medium text-neutral-700 capitalize w-20">{{ day.label }}:</span>
        <span v-if="!parsedSchedule[day.key] || !parsedSchedule[day.key].enabled" class="text-neutral-500">
          Closed
        </span>
        <span v-else-if="!parsedSchedule[day.key].periods || parsedSchedule[day.key].periods.length === 0" class="text-neutral-500">
          No hours set
        </span>
        <span v-else class="text-neutral-900 font-mono text-xs">
          {{ formatPeriods(parsedSchedule[day.key].periods) }}
        </span>
      </div>
      <div class="flex justify-between items-center text-sm pt-2 border-t border-neutral-200">
        <span class="font-medium text-neutral-700">Holidays:</span>
        <span class="text-neutral-900">
          {{ parsedSchedule.holidays && parsedSchedule.holidays.enabled ? 'Open' : 'Closed' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'BranchSchedulePreview',
  props: {
    schedule: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    const daysOfWeek = [
      { key: 'monday', label: 'Mon' },
      { key: 'tuesday', label: 'Tue' },
      { key: 'wednesday', label: 'Wed' },
      { key: 'thursday', label: 'Thu' },
      { key: 'friday', label: 'Fri' },
      { key: 'saturday', label: 'Sat' },
      { key: 'sunday', label: 'Sun' }
    ]

    const parsedSchedule = computed(() => {
      if (!props.schedule) return null

      try {
        return JSON.parse(props.schedule)
      } catch (error) {
        console.error('Error parsing schedule:', error)
        return null
      }
    })

    const formatPeriods = (periods) => {
      if (!periods || periods.length === 0) return 'No hours'

      return periods.map(period => `${period.start}-${period.end}`).join(', ')
    }

    return {
      daysOfWeek,
      parsedSchedule,
      formatPeriods
    }
  }
}
</script>