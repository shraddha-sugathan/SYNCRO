<script setup>
import { computed, ref } from 'vue'
import { money, formatDate, intervalLabel } from '../../utils/format.js'

const props = defineProps({ products: { type: Array, default: () => [] }, loading: Boolean })

const query = ref('')
const type = ref('all')
const sortKey = ref('created_at')
const sortDir = ref('desc')

const columns = [
  { key: 'name', label: 'Product' },
  { key: 'type', label: 'Type' },
  { key: 'price_cents', label: 'Sales price', num: true },
  { key: 'cost_cents', label: 'Cost', num: true },
  { key: 'margin', label: 'Margin', num: true },
  { key: 'billing_interval', label: 'Billing' },
  { key: 'created_at', label: 'Added' },
]

const types = computed(() => [...new Set(props.products.map((p) => p.type).filter(Boolean))])
const marginPct = (p) => (p.price_cents > 0 ? ((p.price_cents - (p.cost_cents || 0)) / p.price_cents) * 100 : 0)
const valueOf = (p, k) => (k === 'margin' ? marginPct(p) : p[k])

const rows = computed(() => {
  const q = query.value.trim().toLowerCase()
  const list = props.products.filter(
    (p) => (type.value === 'all' || p.type === type.value) && (!q || `${p.name} ${p.notes || ''}`.toLowerCase().includes(q)),
  )
  const dir = sortDir.value === 'asc' ? 1 : -1
  return [...list].sort((a, b) => {
    const x = valueOf(a, sortKey.value)
    const y = valueOf(b, sortKey.value)
    return (typeof x === 'string' ? x.localeCompare(y) : (x ?? 0) - (y ?? 0)) * dir
  })
})

function sortBy(k) {
  if (sortKey.value === k) sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  else { sortKey.value = k; sortDir.value = k === 'name' || k === 'type' ? 'asc' : 'desc' }
}
const arrow = (k) => (sortKey.value === k ? (sortDir.value === 'asc' ? '▲' : '▼') : '')
</script>

<template>
  <section class="card">
    <header class="table-head">
      <div>
        <h3>Product catalogue</h3>
        <p class="muted small">{{ rows.length }} of {{ products.length }} products</p>
      </div>
      <div class="table-tools">
        <input v-model="query" type="search" placeholder="Search products…" aria-label="Search products" />
        <select v-model="type" aria-label="Filter by type">
          <option value="all">All types</option>
          <option v-for="t in types" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>
    </header>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th v-for="c in columns" :key="c.key" :class="{ num: c.num }" :aria-sort="sortKey === c.key ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'">
              <button class="th-btn" @click="sortBy(c.key)">{{ c.label }} <i>{{ arrow(c.key) }}</i></button>
            </th>
          </tr>
        </thead>
        <tbody v-if="loading">
          <tr v-for="n in 4" :key="n"><td v-for="c in columns" :key="c.key"><div class="skeleton skeleton-line" /></td></tr>
        </tbody>
        <tbody v-else-if="rows.length">
          <tr v-for="p in rows" :key="p.id">
            <td>
              <strong>{{ p.name }}</strong>
              <div v-if="p.variants?.length" class="chips">
                <span v-for="v in p.variants" :key="v.id || v.value" class="chip chip-info">
                  {{ v.attribute }}: {{ v.value }}<template v-if="v.extra_cents"> (+{{ money(v.extra_cents) }})</template>
                </span>
              </div>
            </td>
            <td><span class="chip chip-brand">{{ p.type || 'Software' }}</span></td>
            <td class="num">{{ money(p.price_cents) }}</td>
            <td class="num">{{ p.cost_cents ? money(p.cost_cents) : '—' }}</td>
            <td class="num" :class="marginPct(p) < 0 ? 'neg-text' : ''">{{ p.cost_cents ? `${marginPct(p).toFixed(0)}%` : '—' }}</td>
            <td>per {{ intervalLabel(p.billing_interval) }}</td>
            <td>{{ formatDate(p.created_at) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="!loading && !rows.length" class="empty-note">No products match. Create one using the form.</p>
    </div>
  </section>
</template>
