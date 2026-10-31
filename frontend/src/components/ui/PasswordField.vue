<script setup>
import { ref } from 'vue'

defineProps({
  modelValue: { type: String, default: '' },
  id: { type: String, required: true },
  label: { type: String, default: 'Password' },
  error: { type: String, default: '' },
  autocomplete: { type: String, default: 'current-password' },
})
defineEmits(['update:modelValue', 'blur'])
const visible = ref(false)
</script>

<template>
  <div class="field" :class="{ invalid: error }">
    <label :for="id">{{ label }}</label>
    <div class="input-wrap">
      <input
        :id="id" :type="visible ? 'text' : 'password'" :value="modelValue" :autocomplete="autocomplete"
        :aria-invalid="!!error" :aria-describedby="error ? `${id}-err` : undefined" placeholder="••••••••••"
        @input="$emit('update:modelValue', $event.target.value)" @blur="$emit('blur')"
      />
      <button type="button" class="reveal" :aria-label="visible ? 'Hide password' : 'Show password'" @click="visible = !visible">
        {{ visible ? 'Hide' : 'Show' }}
      </button>
    </div>
    <p v-if="error" :id="`${id}-err`" class="field-error">{{ error }}</p>
  </div>
</template>
