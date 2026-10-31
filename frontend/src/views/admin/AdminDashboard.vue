<script setup>
import { onMounted, ref } from 'vue'
import MetricCard from '../../components/admin/MetricCard.vue'
import ProductForm from '../../components/admin/ProductForm.vue'
import ProductTable from '../../components/admin/ProductTable.vue'
import { api } from '../../services/api.js'
import { useToast } from '../../composables/useToast.js'
import { money } from '../../utils/format.js'

const toast = useToast()
const metrics = ref(null)
const products = ref([])
const loadingMetrics = ref(true)
const loadingProducts = ref(true)
const saving = ref(false)
const formRef = ref(null)

async function loadMetrics() {
  loadingMetrics.value = true
  try {
    metrics.value = await api.getMetrics()
  } catch (e) {
    toast.error(e.message || 'Could not load metrics')
  } finally {
    loadingMetrics.value = false
  }
}

async function loadProducts() {
  loadingProducts.value = true
  try {
    products.value = await api.getProducts()
  } catch (e) {
    toast.error(e.message || 'Could not load products')
  } finally {
    loadingProducts.value = false
  }
}

const refresh = () => Promise.all([loadMetrics(), loadProducts()])

async function createProduct(payload) {
  saving.value = true
  try {
    const created = await api.createProduct(payload)
    products.value = [created, ...products.value.filter((p) => p.id !== created.id)]
    formRef.value?.reset()
    toast.success(`“${created.name}” created`)
    loadMetrics()
  } catch (e) {
    toast.error(e.message || 'Could not create product')
  } finally {
    saving.value = false
  }
}

onMounted(refresh)
</script>

<template>
  <div class="dashboard">
    <header class="page-head">
      <div>
        <h1>Admin dashboard</h1>
        <p class="muted">Revenue health and your product catalogue at a glance.</p>
      </div>
      <button class="btn btn-ghost btn-sm" :disabled="loadingMetrics || loadingProducts" @click="refresh">↻ Refresh</button>
    </header>

    <section class="grid grid-3" aria-label="Key metrics">
      <MetricCard
        label="Monthly recurring revenue" tone="brand" :loading="loadingMetrics"
        :value="metrics ? money(metrics.mrr_cents) : ''" hint="Active subscriptions, yearly plans ÷ 12"
      />
      <MetricCard
        label="Active subscriptions" tone="ok" :loading="loadingMetrics" :demo="!!metrics?.partial"
        :value="metrics?.active_subscriptions ?? ''" hint="Currently billing customers"
      />
      <MetricCard
        label="Overdue invoices" tone="bad" :loading="loadingMetrics"
        :value="metrics?.overdue_invoices ?? ''"
        :hint="metrics ? `${money(metrics.overdue_amount_cents)} to collect` : ''"
      />
    </section>

    <div class="split">
      <ProductForm ref="formRef" :busy="saving" @submit="createProduct" />
      <ProductTable :products="products" :loading="loadingProducts" />
    </div>
  </div>
</template>
