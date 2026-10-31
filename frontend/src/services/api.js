/**
 * Single entry point for all data access.
 *
 * Every call tries the real Go API first. If it is unreachable, not implemented yet (404/405),
 * failing (5xx) or returns non-JSON (e.g. dev-server HTML), the same call is served from the
 * demo database so the UI is always demoable. Real 4xx answers (400/401/409...) are surfaced as errors.
 * Sessions started with a demo token never touch the network.
 */
import { ref } from 'vue'
import { ApiError, mock, isDemoEmail } from './mock/handlers.js'
import { decodeJwt, isMockToken } from '../utils/jwt.js'

export { ApiError }
export const TOKEN_KEY = 'syncro_token'
const META_KEY = 'syncro_product_meta'

/** 'live' | 'mock' | 'idle' — drives the status badge in the UI. */
export const apiMode = ref('idle')

const store = {
  get: (k) => { try { return localStorage.getItem(k) } catch { return null } },
  set: (k, v) => { try { localStorage.setItem(k, v) } catch { /* ignore */ } },
}
const getToken = () => store.get(TOKEN_KEY)
const claims = () => decodeJwt(getToken()) || {}

async function live(path, { method = 'GET', body, auth = true } = {}) {
  const headers = { 'Content-Type': 'application/json' }
  const token = getToken()
  if (auth && token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(`/api${path}`, { method, headers, body: body ? JSON.stringify(body) : undefined })
  if (!res.ok) throw new ApiError((await res.text()).trim() || res.statusText, res.status)
  return res.json() // throws on HTML/empty bodies -> triggers fallback
}

const shouldFallback = (e) => !(e instanceof ApiError) || [0, 404, 405].includes(e.status) || e.status >= 500

async function call(path, opts, fromMock, adapt = (d) => d) {
  if (isMockToken(getToken())) {
    apiMode.value = 'mock'
    return fromMock()
  }
  try {
    const data = await adapt(await live(path, opts))
    apiMode.value = 'live'
    return data
  } catch (e) {
    if (!shouldFallback(e)) throw e
    apiMode.value = 'mock'
    return fromMock()
  }
}

/* ---------- adapters: Go payloads -> UI models ---------- */

const readMeta = () => { try { return JSON.parse(store.get(META_KEY) || '{}') } catch { return {} } }
const writeMeta = (id, m) => store.set(META_KEY, JSON.stringify({ ...readMeta(), [id]: m }))

const adaptProduct = (p) => ({
  type: 'Software', cost_cents: 0, notes: '', variants: [], ...readMeta()[p.id],
  id: p.id, name: p.name, price_cents: p.price_cents, billing_interval: p.billing_interval, created_at: p.created_at,
})

const TAX_RATE = 0.18
const adaptInvoice = (i, subs, products) => {
  const sub = subs.find((s) => s.id === i.subscription_id)
  const name = products.find((p) => p.id === sub?.product_id)?.name || 'Subscription'
  const total = i.amount_due_cents
  const subtotal = Math.round(total / (1 + TAX_RATE)) // Go API stores one tax-inclusive amount
  const status = i.status === 'paid' ? 'paid' : i.status === 'failed' || new Date(i.due_date) < new Date() ? 'overdue' : 'confirmed'
  return {
    id: i.id, number: `INV-${i.id.slice(0, 6).toUpperCase()}`, user_id: i.user_id, subscription_id: i.subscription_id,
    product_name: name, status, issue_date: i.created_at, due_date: i.due_date, paid_date: status === 'paid' ? i.created_at : null,
    lines: [{ description: name, quantity: 1, unit_cents: subtotal, amount_cents: subtotal }],
    subtotal_cents: subtotal, discount_cents: 0, tax_rate: TAX_RATE, tax_cents: total - subtotal, total_cents: total,
  }
}

const adaptSub = (s, products) => {
  const p = products.find((x) => x.id === s.product_id)
  return {
    id: s.id, number: `SUB-${s.id.slice(0, 6).toUpperCase()}`, user_id: s.user_id, product_id: s.product_id,
    product_name: p?.name || 'Unknown plan', billing_interval: p?.billing_interval || 'monthly', price_cents: p?.price_cents || 0,
    variant_label: '', status: s.status === 'canceled' ? 'closed' : s.status, start_date: s.current_period_start,
    renewal_date: s.current_period_end, created_at: s.created_at,
  }
}

/* ---------- public API ---------- */

export const api = {
  async login(email, password) {
    if (isDemoEmail(email)) { apiMode.value = 'mock'; return mock.login(email, password) }
    try {
      const data = await live('/login', { method: 'POST', body: { email, password }, auth: false })
      apiMode.value = 'live'
      return data
    } catch (e) {
      if (!shouldFallback(e)) throw e
      apiMode.value = 'mock'
      return mock.login(email, password)
    }
  },

  async register({ email, password, role }) {
    const fromMock = () => mock.register({ email, password, role })
    if (isDemoEmail(email)) return fromMock()
    try {
      // Go API ignores role and always creates customers.
      const data = await live('/register', { method: 'POST', body: { email, password }, auth: false })
      apiMode.value = 'live'
      return data
    } catch (e) {
      if (!shouldFallback(e)) throw e
      apiMode.value = 'mock'
      return fromMock()
    }
  },

  getProducts: () => call('/products', {}, mock.getProducts, (list) => list.map(adaptProduct)),

  async createProduct(p) {
    return call(
      '/products',
      { method: 'POST', body: { name: p.name, price_cents: p.price_cents, billing_interval: p.billing_interval } },
      () => mock.createProduct(p),
      (res) => {
        const extra = { type: p.type, cost_cents: p.cost_cents, notes: p.notes, variants: p.variants }
        writeMeta(res.id, extra) // fields the Go model doesn't store yet are kept client-side
        return { ...p, id: res.id, created_at: new Date().toISOString() }
      },
    )
  },

  subscribe: (productId, variantId) =>
    call('/subscribe', { method: 'POST', body: { product_id: productId } }, () => mock.subscribe(claims().user_id, productId, variantId)),

  async getMySubscriptions() {
    return call('/my-subscriptions', {}, () => mock.getMySubscriptions(claims().user_id), async (subs) => {
      const products = await this.getProducts()
      return subs.map((s) => adaptSub(s, products))
    })
  },

  async getMyInvoices() {
    return call('/invoices', {}, () => mock.getMyInvoices(claims().user_id), (list) => list).then(async (list) => {
      if (apiMode.value !== 'live') return list
      const [subs, products] = await Promise.all([this.getMySubscriptions(), this.getProducts()])
      return list.map((i) => adaptInvoice(i, subs, products))
    })
  },

  /** Admin KPIs. The Go API only exposes active_mrr_cents so far; the rest is merged from demo data until it lands. */
  async getMetrics() {
    const demo = await mock.getMetrics()
    return call('/revenue', {}, async () => ({ ...demo, partial: false }), (r) => ({ ...demo, mrr_cents: r.active_mrr_cents, partial: true }))
  },
}
