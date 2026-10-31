<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth.js'
import ApiModeBadge from '../ui/ApiModeBadge.vue'
import { cap } from '../../utils/format.js'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.replace('/login')
}
</script>

<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar-inner">
        <div class="logo"><span class="logo-mark">S</span> Syncro</div>
        <div class="topbar-right">
          <ApiModeBadge />
          <span class="chip chip-brand">{{ cap(auth.role || '') }}</span>
          <span class="who" :title="auth.email">{{ auth.email }}</span>
          <button class="btn btn-ghost btn-sm" @click="logout">Sign out</button>
        </div>
      </div>
    </header>

    <div v-if="auth.notice" class="banner" role="status">
      {{ auth.notice }}
      <button class="banner-x" aria-label="Dismiss" @click="auth.notice = ''">×</button>
    </div>

    <main class="page"><RouterView /></main>
  </div>
</template>
