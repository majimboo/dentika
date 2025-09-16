<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h4 class="text-lg font-semibold text-neutral-900">Operating Schedule</h4>
      <div class="flex space-x-2">
        <button
          type="button"
          @click="applyTemplate('weekdays')"
          class="px-3 py-1 text-xs font-medium text-primary-700 bg-primary-100 rounded-lg hover:bg-primary-200 transition-colors"
        >
          Weekdays 9-5
        </button>
        <button
          type="button"
          @click="applyTemplate('24/7')"
          class="px-3 py-1 text-xs font-medium text-primary-700 bg-primary-100 rounded-lg hover:bg-primary-200 transition-colors"
        >
          24/7
        </button>
        <button
          type="button"
          @click="applyTemplate('clear')"
          class="px-3 py-1 text-xs font-medium text-neutral-700 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors"
        >
          Clear All
        </button>
      </div>
    </div>

    <div class="space-y-4">
      <div
        v-for="day in daysOfWeek"
        :key="day.key"
        class="border border-neutral-200 rounded-xl p-4 bg-neutral-50"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center space-x-3">
            <input
              :id="`${day.key}-enabled`"
              type="checkbox"
              v-model="schedule[day.key].enabled"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
            >
            <label
              :for="`${day.key}-enabled`"
              class="text-sm font-medium text-neutral-900 cursor-pointer"
            >
              {{ day.label }}
            </label>
          </div>
          <button
            v-if="schedule[day.key].enabled"
            type="button"
            @click="addPeriod(day.key)"
            class="text-xs font-medium text-primary-700 hover:text-primary-800 transition-colors"
          >
            + Add Period
          </button>
        </div>

        <div v-if="schedule[day.key].enabled" class="space-y-2">
          <div
            v-for="(period, periodIndex) in schedule[day.key].periods"
            :key="periodIndex"
            class="flex items-center space-x-3 bg-white p-3 rounded-lg border border-neutral-200"
          >
            <div class="flex items-center space-x-2">
              <label class="text-xs font-medium text-neutral-600">From:</label>
              <input
                type="time"
                v-model="period.start"
                class="px-2 py-1 text-sm border border-neutral-300 rounded focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
              >
            </div>
            <div class="flex items-center space-x-2">
              <label class="text-xs font-medium text-neutral-600">To:</label>
              <input
                type="time"
                v-model="period.end"
                class="px-2 py-1 text-sm border border-neutral-300 rounded focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
              >
            </div>
            <button
              v-if="schedule[day.key].periods.length > 1"
              type="button"
              @click="removePeriod(day.key, periodIndex)"
              class="text-danger-600 hover:text-danger-700 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="border border-neutral-200 rounded-xl p-4 bg-neutral-50">
      <div class="flex items-center space-x-3">
        <input
          id="holidays-enabled"
          type="checkbox"
          v-model="schedule.holidays.enabled"
          class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
        >
        <label for="holidays-enabled" class="text-sm font-medium text-neutral-900 cursor-pointer">
          Open on holidays
        </label>
      </div>
    </div>

    <div class="bg-primary-50 border border-primary-200 rounded-xl p-4">
      <h5 class="text-sm font-semibold text-primary-900 mb-2">Schedule Preview</h5>
      <div class="text-xs text-primary-800 space-y-1">
        <div v-for="day in daysOfWeek" :key="day.key">
          <span class="font-medium">{{ day.label }}:</span>
          <span v-if="!schedule[day.key].enabled" class="text-neutral-600 ml-2">Closed</span>
          <span v-else-if="schedule[day.key].periods.length === 0" class="text-neutral-600 ml-2">No hours set</span>
          <span v-else class="ml-2">
            {{ schedule[day.key].periods.map(p => `${p.start}-${p.end}`).join(', ') }}
          </span>
        </div>
        <div class="pt-2 border-t border-primary-200">
          <span class="font-medium">Holidays:</span>
          <span class="ml-2">{{ schedule.holidays.enabled ? 'Open' : 'Closed' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue'

export default {
  name: 'ScheduleEditor',
  props: {
    modelValue: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const daysOfWeek = [
      { key: 'monday', label: 'Monday' },
      { key: 'tuesday', label: 'Tuesday' },
      { key: 'wednesday', label: 'Wednesday' },
      { key: 'thursday', label: 'Thursday' },
      { key: 'friday', label: 'Friday' },
      { key: 'saturday', label: 'Saturday' },
      { key: 'sunday', label: 'Sunday' }
    ]

    const defaultSchedule = () => ({
      monday: { enabled: false, periods: [] },
      tuesday: { enabled: false, periods: [] },
      wednesday: { enabled: false, periods: [] },
      thursday: { enabled: false, periods: [] },
      friday: { enabled: false, periods: [] },
      saturday: { enabled: false, periods: [] },
      sunday: { enabled: false, periods: [] },
      holidays: { enabled: false },
      timezone: 'UTC'
    })

    const schedule = ref(defaultSchedule())

    const parseSchedule = (scheduleString) => {
      if (!scheduleString) return defaultSchedule()

      try {
        const parsed = JSON.parse(scheduleString)
        const result = defaultSchedule()

        daysOfWeek.forEach(day => {
          if (parsed[day.key]) {
            result[day.key] = {
              enabled: parsed[day.key].enabled || false,
              periods: parsed[day.key].periods || []
            }
          }
        })

        if (parsed.holidays) {
          result.holidays = { enabled: parsed.holidays.enabled || false }
        }

        result.timezone = parsed.timezone || 'UTC'

        return result
      } catch (error) {
        console.error('Error parsing schedule:', error)
        return defaultSchedule()
      }
    }

    const addPeriod = (dayKey) => {
      const lastPeriod = schedule.value[dayKey].periods[schedule.value[dayKey].periods.length - 1]
      const defaultStart = lastPeriod ? lastPeriod.end : '09:00'
      const defaultEnd = lastPeriod ? '18:00' : '17:00'

      schedule.value[dayKey].periods.push({
        start: defaultStart,
        end: defaultEnd
      })
    }

    const removePeriod = (dayKey, periodIndex) => {
      schedule.value[dayKey].periods.splice(periodIndex, 1)
    }

    const applyTemplate = (template) => {
      switch (template) {
        case 'weekdays':
          daysOfWeek.forEach(day => {
            const isWeekday = !['saturday', 'sunday'].includes(day.key)
            schedule.value[day.key] = {
              enabled: isWeekday,
              periods: isWeekday ? [{ start: '09:00', end: '17:00' }] : []
            }
          })
          schedule.value.holidays.enabled = false
          break
        case '24/7':
          daysOfWeek.forEach(day => {
            schedule.value[day.key] = {
              enabled: true,
              periods: [{ start: '00:00', end: '23:59' }]
            }
          })
          schedule.value.holidays.enabled = true
          break
        case 'clear':
          schedule.value = defaultSchedule()
          break
      }
    }

    watch(schedule, () => {
      const scheduleData = {
        ...schedule.value,
        timezone: schedule.value.timezone || 'UTC'
      }
      emit('update:modelValue', JSON.stringify(scheduleData))
    }, { deep: true })

    onMounted(() => {
      schedule.value = parseSchedule(props.modelValue)

      daysOfWeek.forEach(day => {
        if (schedule.value[day.key].enabled && schedule.value[day.key].periods.length === 0) {
          addPeriod(day.key)
        }
      })
    })

    watch(() => props.modelValue, (newValue) => {
      if (newValue !== JSON.stringify(schedule.value)) {
        schedule.value = parseSchedule(newValue)
      }
    })

    return {
      daysOfWeek,
      schedule,
      addPeriod,
      removePeriod,
      applyTemplate
    }
  }
}
</script>