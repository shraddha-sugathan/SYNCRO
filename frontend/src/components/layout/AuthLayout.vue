<script setup>
import ApiModeBadge from '../ui/ApiModeBadge.vue'
defineProps({ title: String, subtitle: String })
</script>

<template>
  <main class="auth">
    <aside class="auth-brand">
      <div class="logo"><span class="logo-mark">S</span> Syncro</div>
      <div>
        <h2>Recurring revenue, on autopilot.</h2>
        <p>Plans, subscriptions, invoices and renewals in one place. No spreadsheets, no missed billing days.</p>
        <ul class="auth-points">
          <li>Automatic renewals &amp; invoice generation</li>
          <li>Role-based admin and customer portals</li>
          <li>Live MRR and overdue tracking</li>
        </ul>
      </div>
      <small>Odoo Makeathon 9.0 · Kochi</small>
    </aside>

    <section class="auth-panel">
      <div class="auth-card">
        <div class="auth-mode"><ApiModeBadge /></div>
        <h1>{{ title }}</h1>
        <p class="muted">{{ subtitle }}</p>
        <slot />
      </div>
    </section>
  </main>
</template>
