<template>
  <svg class="tick" viewBox="0 0 14 14" fill="none" aria-hidden="true">
    <rect x=".75" y=".75" width="12.5" height="12.5" rx="1.5" stroke="currentColor" stroke-width="1.5" />
    <path v-if="$attrs.on !== undefined && $attrs.on !== false" d="M3.8 7.2l2.1 2.1 4.3-4.6" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
  </svg>
</template>
<script>
export default { inheritAttrs: false }
</script>
