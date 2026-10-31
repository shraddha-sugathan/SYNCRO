// Mirrors handlers/auth_handler.go validatePassword: len > 8, upper, lower, special.
export const PASSWORD_MIN_LENGTH = 9

export const passwordRules = [
  { id: 'length', label: 'More than 8 characters', test: (p) => p.length >= PASSWORD_MIN_LENGTH },
  { id: 'upper', label: 'One uppercase letter', test: (p) => /[A-Z]/.test(p) },
  { id: 'lower', label: 'One lowercase letter', test: (p) => /[a-z]/.test(p) },
  { id: 'special', label: 'One special character', test: (p) => /[^a-zA-Z0-9]/.test(p) },
]

export const checkPassword = (p = '') => passwordRules.map((r) => ({ id: r.id, label: r.label, ok: r.test(p) }))
export const isPasswordValid = (p = '') => passwordRules.every((r) => r.test(p))
export const isEmailValid = (e = '') => /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(e.trim())
export const PASSWORD_POLICY_MESSAGE =
  'Use 9+ characters with an uppercase letter, a lowercase letter and a special character'
