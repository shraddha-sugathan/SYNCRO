<script setup>
import { computed } from 'vue'
import { apiMode } from '../../services/api.js'

const meta = computed(() =>
  apiMode.value === 'live'
    ? { label: 'Live API', cls: 'is-live', tip: 'Connected to the Go backend' }
    : apiMode.value === 'mock'
      ? { label: 'Demo data', cls: 'is-mock', tip: 'Backend unreachable or endpoint missing: serving built-in demo data' }
      : { label: 'Ready', cls: '', tip: 'No request made yet' },
)
</script>

<template>
  <span class="mode-badge" :class="meta.cls" :title="meta.tip"><i /> {{ meta.label }}</span>
</template>
