<script setup>
import { computed, reactive, ref } from 'vue'
import { money } from '../../utils/format.js'
import { uid } from '../../services/mock/db.js'

defineProps({ busy: Boolean })
const emit = defineEmits(['submit'])

const TYPES = ['Software', 'Add-on', 'Service', 'Hardware']
const blank = () => ({ name: '', type: 'Software', price: '', cost: '', billing_interval: 'monthly', notes: '', variants: [] })

const form = reactive(blank())
const submitted = ref(false)

const price = computed(() => Number(form.price))
const cost = computed(() => Number(form.cost || 0))
const margin = computed(() => (price.value > 0 ? ((price.value - cost.value) / price.value) * 100 : null))

const errors = computed(() => ({
  name: form.name.trim().length < 2 ? 'Enter a product name' : '',
  price: !(price.value > 0) ? 'Sales price must be greater than 0' : '',
  cost: form.cost !== '' && !(cost.value >= 0) ? 'Cost cannot be negative' : cost.value > price.value && price.value > 0 ? 'Cost is higher than the sales price' : '',
  variants: form.variants.map((v) => ({
    value: v.value.trim() ? '' : 'Required',
    extra: v.extra !== '' && !(Number(v.extra) >= 0) ? 'Must be 0 or more' : '',
  })),
}))
const invalid = computed(
  () => errors.value.name || errors.value.price || (errors.value.cost && cost.value < 0) || errors.value.variants.some((v) => v.value || v.extra),
)

const addVariant = (preset) => form.variants.push({ attribute: 'Brand', value: preset?.value || '', extra: preset?.extra ?? '' })
const addOdoo = () => {
  if (!form.variants.some((v) => v.value.toLowerCase() === 'odoo')) addVariant({ value: 'Odoo', extra: 560 })
}
const removeVariant = (i) => form.variants.splice(i, 1)

function submit() {
  submitted.value = true
  if (invalid.value) return
  emit('submit', {
    name: form.name.trim(),
    type: form.type,
    price_cents: Math.round(price.value * 100),
    cost_cents: Math.round(cost.value * 100),
    billing_interval: form.billing_interval,
    notes: form.notes.trim(),
    variants: form.variants.map((v) => ({ id: uid(), attribute: v.attribute.trim() || 'Variant', value: v.value.trim(), extra_cents: Math.round(Number(v.extra || 0) * 100) })),
  })
}

function reset() {
  Object.assign(form, blank())
  submitted.value = false
}
defineExpose({ reset })
</script>

<template>
  <form class="card product-form" novalidate @submit.prevent="submit">
    <h3>New product</h3>
    <p class="muted small">Add a plan customers can subscribe to.</p>

    <div class="field" :class="{ invalid: submitted && errors.name }">
      <label for="p-name">Product name</label>
      <input id="p-name" v-model="form.name" placeholder="e.g. Syncro Business" maxlength="80" />
      <p v-if="submitted && errors.name" class="field-error">{{ errors.name }}</p>
    </div>

    <div class="field-row">
      <div class="field">
        <label for="p-type">Type</label>
        <select id="p-type" v-model="form.type"><option v-for="t in TYPES" :key="t">{{ t }}</option></select>
      </div>
      <div class="field">
        <label for="p-int">Billing</label>
        <select id="p-int" v-model="form.billing_interval">
          <option value="monthly">Monthly</option>
          <option value="yearly">Yearly</option>
        </select>
      </div>
    </div>

    <div class="field-row">
      <div class="field" :class="{ invalid: submitted && errors.price }">
        <label for="p-price">Sales price (₹)</label>
        <input id="p-price" v-model="form.price" type="number" min="0" step="0.01" inputmode="decimal" placeholder="1999" />
        <p v-if="submitted && errors.price" class="field-error">{{ errors.price }}</p>
      </div>
      <div class="field" :class="{ invalid: errors.cost }">
        <label for="p-cost">Cost price (₹)</label>
        <input id="p-cost" v-model="form.cost" type="number" min="0" step="0.01" inputmode="decimal" placeholder="540" />
        <p v-if="errors.cost" class="field-error">{{ errors.cost }}</p>
      </div>
    </div>
    <p v-if="margin !== null" class="margin-line" :class="{ neg: margin < 0 }">
      Margin {{ money(Math.round((price - cost) * 100)) }} · {{ margin.toFixed(0) }}%
    </p>

    <div class="field">
      <div class="label-row">
        <label>Variants <span class="muted">(optional)</span></label>
        <span class="variant-actions">
          <button type="button" class="link" @click="addOdoo">+ Odoo brand (+₹560)</button>
          <button type="button" class="link" @click="addVariant()">+ Custom</button>
        </span>
      </div>
      <p v-if="!form.variants.length" class="muted small">No variants. Price stays the same for everyone.</p>
      <div v-for="(v, i) in form.variants" :key="i" class="variant-row">
        <input v-model="v.attribute" aria-label="Attribute" placeholder="Attribute" />
        <div :class="{ invalid: submitted && errors.variants[i]?.value }">
          <input v-model="v.value" aria-label="Value" placeholder="Value" />
        </div>
        <div :class="{ invalid: submitted && errors.variants[i]?.extra }">
          <input v-model="v.extra" type="number" min="0" step="0.01" aria-label="Extra price" placeholder="+₹" />
        </div>
        <button type="button" class="icon-btn" aria-label="Remove variant" @click="removeVariant(i)">×</button>
      </div>
    </div>

    <div class="field">
      <label for="p-notes">Notes</label>
      <textarea id="p-notes" v-model="form.notes" rows="2" maxlength="240" placeholder="Seats, limits, what's included…" />
    </div>

    <button class="btn btn-primary btn-block" :disabled="busy">{{ busy ? 'Saving…' : 'Create product' }}</button>
  </form>
</template>
