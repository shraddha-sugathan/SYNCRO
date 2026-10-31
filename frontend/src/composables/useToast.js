import { reactive } from 'vue'

export const toasts = reactive([])
let seq = 0

export function useToast() {
  const push = (message, type = 'info', ms = 3800) => {
    const id = ++seq
    toasts.push({ id, message, type })
    setTimeout(() => dismiss(id), ms)
  }
  const dismiss = (id) => {
    const i = toasts.findIndex((t) => t.id === id)
    if (i > -1) toasts.splice(i, 1)
  }
  return {
    toasts, dismiss,
    success: (m) => push(m, 'success'),
    error: (m) => push(m, 'error', 5500),
    info: (m) => push(m, 'info'),
  }
}
