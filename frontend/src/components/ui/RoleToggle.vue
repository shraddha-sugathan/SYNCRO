<script setup>
defineProps({ modelValue: { type: String, required: true } })
defineEmits(['update:modelValue'])
const roles = [
  { value: 'customer', label: 'Customer', hint: 'Manage my plans & invoices' },
  { value: 'admin', label: 'Admin', hint: 'Products & revenue' },
]
</script>

<template>
  <div class="role-toggle" role="radiogroup" aria-label="Account role">
    <button
      v-for="r in roles" :key="r.value" type="button" role="radio"
      :aria-checked="modelValue === r.value" :class="{ active: modelValue === r.value }"
      @click="$emit('update:modelValue', r.value)"
    >
      <strong>{{ r.label }}</strong>
      <small>{{ r.hint }}</small>
    </button>
  </div>
</template>
