const inr = new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 2 })
const date = new Intl.DateTimeFormat('en-IN', { day: '2-digit', month: 'short', year: 'numeric' })

export const money = (cents = 0) => inr.format((cents || 0) / 100)
export const formatDate = (v) => (v ? date.format(new Date(v)) : '—')
export const cap = (s = '') => s.charAt(0).toUpperCase() + s.slice(1)
export const intervalLabel = (i) => ({ daily: 'day', weekly: 'week', monthly: 'month', yearly: 'year' }[i] || i)

export function daysUntil(v) {
  return Math.ceil((new Date(v).getTime() - Date.now()) / 86400000)
}
