import { db, save, uid, buildInvoice, addInterval, iso, DEMO_ACCOUNTS } from './db.js'
import { createMockJwt } from '../../utils/jwt.js'
import { isPasswordValid, PASSWORD_POLICY_MESSAGE } from '../../utils/validators.js'

export class ApiError extends Error {
  constructor(message, status = 0) {
    super(message)
    this.status = status
  }
}

const sleep = (ms = 280) => new Promise((r) => setTimeout(r, ms))
const clone = (v) => JSON.parse(JSON.stringify(v))
export const isDemoEmail = (email) =>
  Object.values(DEMO_ACCOUNTS).some((a) => a.email === String(email).trim().toLowerCase())

// Tokens issued by the live backend carry user ids unknown to the demo DB; show the demo customer's data then.
const resolveUser = (userId) => {
  const { users } = db()
  return users.some((u) => u.id === userId) ? userId : users.find((u) => u.role === 'customer').id
}

const withStatus = (inv) =>
  inv.status === 'confirmed' && new Date(inv.due_date) < new Date() ? { ...inv, status: 'overdue' } : inv

const monthly = (s) => (s.billing_interval === 'yearly' ? s.price_cents / 12 : s.billing_interval === 'weekly' ? s.price_cents * 4.345 : s.billing_interval === 'daily' ? s.price_cents * 30 : s.price_cents)

export const mock = {
  async login(email, password) {
    await sleep()
    const u = db().users.find((x) => x.email === String(email).trim().toLowerCase())
    if (!u || u.password !== password) throw new ApiError('Invalid credentials', 401)
    return { token: createMockJwt({ userId: u.id, role: u.role, email: u.email }) }
  },

  async register({ email, password, role }) {
    await sleep()
    const data = db()
    const e = String(email).trim().toLowerCase()
    if (!isPasswordValid(password)) throw new ApiError(PASSWORD_POLICY_MESSAGE, 400)
    if (data.users.some((u) => u.email === e)) throw new ApiError('Email already exists', 409)
    data.users.push({ id: uid(), email: e, password, role: role === 'admin' ? 'admin' : 'customer' })
    save()
    return { message: 'Account created successfully' }
  },

  async getProducts() {
    await sleep(180)
    return clone(db().products)
  },

  async createProduct(p) {
    await sleep()
    const product = { id: uid(), created_at: iso(0), variants: [], notes: '', cost_cents: 0, type: 'Software', ...p }
    db().products.unshift(product)
    save()
    return clone(product)
  },

  async subscribe(userId, productId, variantId) {
    await sleep(450)
    const data = db()
    const user = resolveUser(userId)
    const p = data.products.find((x) => x.id === productId)
    if (!p) throw new ApiError('Product not found', 404)
    const v = p.variants?.find((x) => x.id === variantId)
    const price = p.price_cents + (v?.extra_cents || 0)
    const start = iso(0)
    const sub = {
      id: uid(), number: `SUB-${String(++data.counters.sub).padStart(4, '0')}`, user_id: user, product_id: p.id,
      product_name: p.name, billing_interval: p.billing_interval, price_cents: price,
      variant_label: v ? `${v.attribute}: ${v.value}` : '', status: 'active',
      start_date: start, renewal_date: addInterval(start, p.billing_interval), created_at: start,
    }
    data.subscriptions.unshift(sub)
    data.invoices.unshift(buildInvoice({
      number: `INV-2026-${String(++data.counters.inv).padStart(4, '0')}`, user_id: user, subscription_id: sub.id,
      product_name: p.name, unit_cents: price, status: 'confirmed', issue: 0, due: 7,
    }))
    save()
    return clone(sub)
  },

  async getMySubscriptions(userId) {
    await sleep(200)
    const user = resolveUser(userId)
    return clone(db().subscriptions.filter((s) => s.user_id === user))
  },

  async getMyInvoices(userId) {
    await sleep(200)
    const user = resolveUser(userId)
    return clone(db().invoices.filter((i) => i.user_id === user).map(withStatus))
  },

  async getMetrics() {
    await sleep(200)
    const { subscriptions, invoices } = db()
    const active = subscriptions.filter((s) => s.status === 'active')
    const all = invoices.map(withStatus)
    const overdue = all.filter((i) => i.status === 'overdue')
    return {
      mrr_cents: Math.round(active.reduce((sum, s) => sum + monthly(s), 0)),
      active_subscriptions: active.length,
      overdue_invoices: overdue.length,
      overdue_amount_cents: overdue.reduce((sum, i) => sum + i.total_cents, 0),
      collected_cents: all.filter((i) => i.status === 'paid').reduce((sum, i) => sum + i.total_cents, 0),
    }
  },
}
