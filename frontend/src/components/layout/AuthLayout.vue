<script setup>
import ApiModeBadge from '../ui/ApiModeBadge.vue'
import BrandMark from '../ui/BrandMark.vue'
defineProps({ title: String, subtitle: String })
</script>

<template>
  <main class="auth">
    <aside class="auth-brand">
      <BrandMark inverse />
      <div>
        <h2>Every renewal, accounted for.</h2>
        <p class="lede">Plans, subscriptions and invoices in one ledger, with renewals billed automatically on the day they fall due.</p>

        <div class="stmt" aria-hidden="true">
          <div class="stmt-head"><b>Invoice</b><span>INV-4F2A91</span></div>
          <dl>
            <div class="row"><dt>Studio plan, monthly</dt><dd>1,180.00</dd></div>
            <div class="row"><dt>Brand: Odoo</dt><dd>560.00</dd></div>
            <div class="row sub"><dt>Includes GST at 18%</dt><dd>265.42</dd></div>
          </dl>
          <div class="total"><span>Total, INR</span><b>₹1,740.00</b></div>
          <span class="stamp">PAID</span>
        </div>
      </div>
      <small class="foot">Built for Odoo Makeathon 9.0, Kochi</small>
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
