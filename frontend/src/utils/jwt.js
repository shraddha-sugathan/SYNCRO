const toB64Url = (obj) => {
  const bytes = new TextEncoder().encode(JSON.stringify(obj))
  let bin = ''
  bytes.forEach((b) => (bin += String.fromCharCode(b)))
  return btoa(bin).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

const fromB64Url = (s) => {
  const b64 = s.replace(/-/g, '+').replace(/_/g, '/').padEnd(Math.ceil(s.length / 4) * 4, '=')
  const bytes = Uint8Array.from(atob(b64), (c) => c.charCodeAt(0))
  return JSON.parse(new TextDecoder().decode(bytes))
}

const part = (token, i) => {
  try {
    return fromB64Url(String(token).split('.')[i])
  } catch {
    return null
  }
}

export const decodeJwt = (token) => (token ? part(token, 1) : null)
export const isMockToken = (token) => !!token && part(token, 0)?.mock === true

/** Unsigned JWT-shaped token used only when the Go API is unreachable (demo mode). */
export function createMockJwt({ userId, role, email }, ttlHours = 24) {
  const header = { alg: 'none', typ: 'JWT', mock: true }
  const payload = { user_id: userId, role, email, exp: Math.floor(Date.now() / 1000) + ttlHours * 3600 }
  return `${toB64Url(header)}.${toB64Url(payload)}.mock`
}
