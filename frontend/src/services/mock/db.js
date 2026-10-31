// Demo database persisted in localStorage. Shapes match what services/api.js returns to the UI.
const KEY = 'syncro_mock_db_v1'
const DAY = 86400000
const TAX_RATE = 0.18 // GST

export const DEMO_ACCOUNTS = {
  admin: { email: 'admin@syncro.demo', password: 'Admin@12345' },
  customer: { email: 'customer@syncro.demo', password: 'Customer@123' },
}

export const uid = () =>
  globalThis.crypto?.randomUUID?.() ??
  'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    return (c === 'x' ? r : (r & 3) | 8).toString(16)
  })

const iso = (offsetDays) => new Date(Date.now() + offsetDays * DAY).toISOString()
const addInterval = (from, interval) => {
  const d = new Date(from)
  if (interval === 'yearly') d.setFullYear(d.getFullYear() + 1)
  else if (interval === 'weekly') d.setDate(d.getDate() + 7)
  else if (interval === 'daily') d.setDate(d.getDate() + 1)
  else d.setMonth(d.getMonth() + 1)
  return d.toISOString()
}

// First period end after "now", stepping from the start date (mirrors the billing worker's cycle).
function nextRenewal(start, interval) {
  let d = start
  while (new Date(d).getTime() <= Date.now()) d = addInterval(d, interval)
  return d
}

export function buildInvoice({ number, user_id, subscription_id, product_name, unit_cents, status, issue, due, paid = null, discount_cents = 0 }) {
  const subtotal = unit_cents
  const tax = Math.round((subtotal - discount_cents) * TAX_RATE)
  return {
    id: uid(),
    number,
    user_id,
    subscription_id,
    product_name,
    status,
    issue_date: iso(issue),
    due_date: iso(due),
    paid_date: paid === null ? null : iso(paid),
    lines: [{ description: product_name, quantity: 1, unit_cents, amount_cents: unit_cents }],
    subtotal_cents: subtotal,
    discount_cents,
    tax_rate: TAX_RATE,
    tax_cents: tax,
    total_cents: subtotal - discount_cents + tax,
  }
}

function seed() {
  const adminId = uid()
  const custId = uid()
  const odoo = { id: uid(), attribute: 'Brand', value: 'Odoo', extra_cents: 56000 }
  const std = { id: uid(), attribute: 'Brand', value: 'Standard', extra_cents: 0 }
  const mk = (name, type, price, cost, interval, notes, variants = []) => ({
    id: uid(), name, type, price_cents: price * 100, cost_cents: cost * 100, billing_interval: interval, notes, variants, created_at: iso(-40),
  })
  const products = [
    mk('Syncro Starter', 'Software', 499, 120, 'monthly', 'For small shops getting started.'),
    mk('Syncro Business', 'Software', 1999, 540, 'monthly', 'Teams up to 25 seats with priority email support.', [std, odoo]),
    mk('Syncro Enterprise', 'Software', 19999, 6200, 'yearly', 'Unlimited seats, SSO and a dedicated success manager.', [std, odoo]),
    mk('Analytics Add-on', 'Add-on', 799, 150, 'monthly', 'Cohort and churn dashboards.'),
    mk('Priority Support', 'Service', 4999, 1800, 'yearly', '24x7 phone and chat support.'),
  ]
  const [starter, business, enterprise, analytics] = products

  const mkSub = (n, user_id, p, status, startOffset, variant_label = '', price = p.price_cents) => ({
    id: uid(), number: `SUB-${String(n).padStart(4, '0')}`, user_id, product_id: p.id, product_name: p.name,
    billing_interval: p.billing_interval, price_cents: price, variant_label, status,
    start_date: iso(startOffset), renewal_date: nextRenewal(iso(startOffset), p.billing_interval),
    created_at: iso(startOffset),
  })
  const subs = [
    mkSub(1, custId, business, 'active', -64, 'Brand: Odoo', business.price_cents + 56000),
    mkSub(2, custId, analytics, 'active', -33),
    mkSub(3, custId, starter, 'closed', -150),
  ]
  subs[2].renewal_date = iso(-120) // closed plans don't renew
  // Subscriptions of other customers so admin metrics look realistic.
  const others = [uid(), uid(), uid(), uid()]
  subs.push(
    mkSub(4, others[0], enterprise, 'active', -200),
    mkSub(5, others[1], business, 'active', -90),
    mkSub(6, others[2], starter, 'past_due', -45),
    mkSub(7, others[3], analytics, 'active', -12),
  )

  const inv = (n, sub, p, status, issue, due, paid, unit = sub.price_cents, discount = 0) =>
    buildInvoice({ number: `INV-2026-${String(n).padStart(4, '0')}`, user_id: sub.user_id, subscription_id: sub.id, product_name: p.name, unit_cents: unit, status, issue, due, paid, discount_cents: discount })
  const invoices = [
    inv(1, subs[0], business, 'paid', -64, -57, -60),
    inv(2, subs[0], business, 'paid', -34, -27, -30),
    inv(3, subs[1], analytics, 'paid', -33, -26, -31),
    inv(4, subs[0], business, 'confirmed', -4, 3, null),
    inv(5, subs[1], analytics, 'confirmed', -9, -2, null),
    inv(6, subs[0], business, 'draft', 0, 7, null),
    inv(7, subs[5], starter, 'confirmed', -20, -13, null),
    inv(8, subs[4], business, 'paid', -30, -23, -22),
    inv(9, subs[3], enterprise, 'paid', -200, -193, -195),
    inv(10, subs[6], analytics, 'confirmed', -5, 2, null),
  ]

  return {
    users: [
      { id: adminId, email: DEMO_ACCOUNTS.admin.email, password: DEMO_ACCOUNTS.admin.password, role: 'admin' },
      { id: custId, email: DEMO_ACCOUNTS.customer.email, password: DEMO_ACCOUNTS.customer.password, role: 'customer' },
    ],
    products, subscriptions: subs, invoices, counters: { sub: 7, inv: 10 },
  }
}

let cache = null
export function db() {
  if (cache) return cache
  try {
    const raw = localStorage.getItem(KEY)
    if (raw) return (cache = JSON.parse(raw))
  } catch { /* storage unavailable: fall through to in-memory seed */ }
  cache = seed()
  save()
  return cache
}

export function save() {
  try {
    localStorage.setItem(KEY, JSON.stringify(cache))
  } catch { /* ignore */ }
}

export function resetMockDb() {
  cache = null
  try { localStorage.removeItem(KEY) } catch { /* ignore */ }
  return db()
}

export { addInterval, iso, TAX_RATE }
