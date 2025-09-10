<template>
  <div class="tag-input-container">
    <div class="flex flex-wrap gap-2 mb-2">
      <span
        v-for="(tag, index) in tags"
        :key="index"
        class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-800 border border-primary-200"
      >
        {{ tag }}
        <button
          type="button"
          @click="removeTag(index)"
          class="ml-2 inline-flex items-center justify-center w-4 h-4 rounded-full hover:bg-primary-200 focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </span>
    </div>

    <div class="relative">
      <input
        ref="inputRef"
        v-model="inputValue"
        @keydown="handleKeydown"
        @blur="handleBlur"
        :placeholder="placeholder"
        class="block w-full px-4 py-3 border border-neutral-300 rounded-xl text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 bg-neutral-50 hover:bg-white focus:bg-white"
        :class="{ 'border-red-500': hasError }"
      />

      <div v-if="suggestions.length > 0" class="absolute z-10 w-full mt-1 bg-white border border-neutral-300 rounded-xl shadow-lg max-h-40 overflow-y-auto">
        <div
          v-for="(suggestion, index) in suggestions"
          :key="suggestion"
          @click="addSuggestion(suggestion)"
          :class="[
            'px-4 py-2 cursor-pointer hover:bg-neutral-50',
            { 'bg-primary-50': highlightedIndex === index }
          ]"
        >
          {{ suggestion }}
        </div>
      </div>
    </div>

    <p v-if="helpText" class="mt-2 text-sm text-neutral-500">{{ helpText }}</p>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  placeholder: {
    type: String,
    default: 'Type and press Enter to add...'
  },
  suggestions: {
    type: Array,
    default: () => []
  },
  helpText: {
    type: String,
    default: ''
  },
  hasError: {
    type: Boolean,
    default: false
  },
  separator: {
    type: String,
    default: ','
  }
})

const emit = defineEmits(['update:modelValue'])

const inputRef = ref(null)
const inputValue = ref('')
const highlightedIndex = ref(-1)

// Common allergy suggestions
const defaultSuggestions = [
  'Penicillin',
  'Aspirin',
  'Ibuprofen',
  'Codeine',
  'Morphine',
  'Latex',
  'Shellfish',
  'Peanuts',
  'Tree nuts',
  'Eggs',
  'Milk',
  'Soy',
  'Wheat',
  'Fish',
  'Dust mites',
  'Pollen',
  'Animal dander',
  'Mold',
  'Insect stings',
  'Contrast dye'
]

const tags = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const filteredSuggestions = computed(() => {
  if (!inputValue.value.trim()) return []

  const input = inputValue.value.toLowerCase()
  return defaultSuggestions
    .concat(props.suggestions)
    .filter(suggestion =>
      suggestion.toLowerCase().includes(input) &&
      !tags.value.some(tag => tag.toLowerCase() === suggestion.toLowerCase())
    )
    .slice(0, 5)
})

const suggestions = computed(() => {
  return inputValue.value.trim() ? filteredSuggestions.value : []
})

const addTag = (tag) => {
  const trimmedTag = tag.trim()
  if (trimmedTag && !tags.value.some(t => t.toLowerCase() === trimmedTag.toLowerCase())) {
    tags.value = [...tags.value, trimmedTag]
    inputValue.value = ''
    highlightedIndex.value = -1
  }
}

const removeTag = (index) => {
  tags.value = tags.value.filter((_, i) => i !== index)
}

const addSuggestion = (suggestion) => {
  addTag(suggestion)
  inputRef.value?.focus()
}

const handleKeydown = (event) => {
  if (event.key === 'Enter') {
    event.preventDefault()
    if (highlightedIndex.value >= 0 && suggestions.value[highlightedIndex.value]) {
      addSuggestion(suggestions.value[highlightedIndex.value])
    } else if (inputValue.value.trim()) {
      addTag(inputValue.value)
    }
  } else if (event.key === 'Backspace' && !inputValue.value && tags.value.length > 0) {
    removeTag(tags.value.length - 1)
  } else if (event.key === 'ArrowDown') {
    event.preventDefault()
    if (suggestions.value.length > 0) {
      highlightedIndex.value = Math.min(highlightedIndex.value + 1, suggestions.value.length - 1)
    }
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    highlightedIndex.value = Math.max(highlightedIndex.value - 1, -1)
  } else if (event.key === 'Escape') {
    highlightedIndex.value = -1
  } else if (event.key === ',') {
    event.preventDefault()
    if (inputValue.value.trim()) {
      addTag(inputValue.value)
    }
  }
}

const handleBlur = () => {
  // Add tag on blur if there's input
  if (inputValue.value.trim()) {
    addTag(inputValue.value)
  }
  highlightedIndex.value = -1
}

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  if (JSON.stringify(newValue) !== JSON.stringify(tags.value)) {
    tags.value = [...newValue]
  }
}, { deep: true })

// Auto-focus removed to prevent unwanted focus on page load

defineExpose({
  focus: () => inputRef.value?.focus()
})
</script>

<style scoped>
.tag-input-container {
  @apply relative;
}
</style>