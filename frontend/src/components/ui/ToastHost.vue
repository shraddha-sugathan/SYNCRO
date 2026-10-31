<script setup>
import { useToast } from '../../composables/useToast.js'
const { toasts, dismiss } = useToast()
</script>

<template>
  <div class="toast-host" aria-live="polite">
    <TransitionGroup name="toast">
      <div v-for="t in toasts" :key="t.id" class="toast" :class="`toast-${t.type}`" role="status" @click="dismiss(t.id)">
        {{ t.message }}
      </div>
    </TransitionGroup>
  </div>
</template>
