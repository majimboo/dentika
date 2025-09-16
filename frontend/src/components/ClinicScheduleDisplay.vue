<template>
  <div class="text-sm">
    <div v-if="!parsedSchedule" class="text-gray-500">
      Schedule not available
    </div>
    <div v-else class="space-y-1">
      <div v-for="day in daysOfWeek" :key="day.key" class="flex justify-between items-center">
        <span class="font-medium text-gray-700 capitalize min-w-[4rem]">{{ day.label }}:</span>
        <span v-if="!parsedSchedule[day.key] || !parsedSchedule[day.key].enabled" class="text-gray-500">
          Closed
        </span>
        <span v-else-if="!parsedSchedule[day.key].periods || parsedSchedule[day.key].periods.length === 0" class="text-gray-500">
          No hours set
        </span>
        <span v-else class="text-gray-900 font-mono text-xs">
          {{ formatPeriods(parsedSchedule[day.key].periods) }}
        </span>
      </div>
      <div v-if="showHolidays" class="flex justify-between items-center pt-2 border-t border-gray-200">
        <span class="font-medium text-gray-700">Holidays:</span>
        <span class="text-gray-900">
          {{ parsedSchedule.holidays && parsedSchedule.holidays.enabled ? 'Open' : 'Closed' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'ClinicScheduleDisplay',
  props: {
    schedule: {
      type: String,
      default: ''
    },
    showHolidays: {
      type: Boolean,
      default: true
    },
    compact: {
      type: Boolean,
      default: false
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

      return periods.map(period => {
        const startTime = formatTime(period.start)
        const endTime = formatTime(period.end)
        return `${startTime}-${endTime}`
      }).join(', ')
    }

    const formatTime = (timeString) => {
      if (!timeString) return ''

      try {
        const [hours, minutes] = timeString.split(':')
        const hour = parseInt(hours)
        const minute = parseInt(minutes)

        const period = hour >= 12 ? 'PM' : 'AM'
        const displayHour = hour > 12 ? hour - 12 : hour === 0 ? 12 : hour

        if (minute === 0) {
          return `${displayHour}${period}`
        } else {
          return `${displayHour}:${minute.toString().padStart(2, '0')}${period}`
        }
      } catch (error) {
        return timeString
      }
    }

    return {
      daysOfWeek,
      parsedSchedule,
      formatPeriods
    }
  }
}
</script>