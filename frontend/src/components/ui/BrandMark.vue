<script setup>
defineProps({ inverse: Boolean })
</script>

<template>
  <span class="brand" :class="{ inverse }">
    <svg viewBox="0 0 32 32" aria-hidden="true" fill="none">
      <!-- a geometric S drawn as two linked arcs: one cycle flowing into the next -->
      <path d="M22.5 10.5A6.5 6.5 0 1 0 16 17A6.5 6.5 0 1 1 9.5 23.5" stroke="currentColor" stroke-width="3.2" stroke-linecap="butt" />
    </svg>
    <span class="brand-word">Syncro</span>
  </span>
</template>
