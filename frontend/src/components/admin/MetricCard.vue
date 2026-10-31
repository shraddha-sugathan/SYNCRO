<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], default: '' },
  hint: { type: String, default: '' },
  tone: { type: String, default: 'brand' }, // brand | ok | bad
  loading: Boolean,
  demo: Boolean,
})
</script>

<template>
  <article class="card metric" :class="`tone-${tone}`">
    <header>
      <span class="metric-label">{{ label }}</span>
      <span v-if="demo" class="chip chip-muted" title="Backend has no endpoint for this yet: showing demo data">demo</span>
    </header>
    <div v-if="loading" class="skeleton metric-skeleton" aria-busy="true" />
    <strong v-else class="metric-value">{{ value }}</strong>
    <small class="metric-hint">{{ hint }}</small>
  </article>
</template>
