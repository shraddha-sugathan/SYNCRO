import { defineStore } from 'pinia'
import { api, TOKEN_KEY } from '../services/api.js'
import { decodeJwt } from '../utils/jwt.js'
import { cap } from '../utils/format.js'

const EMAIL_KEY = 'syncro_email'
const read = (k) => { try { return localStorage.getItem(k) || '' } catch { return '' } }
const write = (k, v) => { try { v ? localStorage.setItem(k, v) : localStorage.removeItem(k) } catch { /* ignore */ } }

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: read(TOKEN_KEY), email: read(EMAIL_KEY), notice: '' }),

  getters: {
    claims: (s) => (s.token ? decodeJwt(s.token) : null),
    role() { return this.claims?.role ?? null },
    isAuthenticated() {
      const c = this.claims
      return !!c && (!c.exp || c.exp * 1000 > Date.now())
    },
  },

  actions: {
    _start(token, email) {
      this.token = token
      this.email = email.trim().toLowerCase()
      write(TOKEN_KEY, this.token)
      write(EMAIL_KEY, this.email)
    },

    /** `role` is the role picked on the form: the account's real role (from the JWT) must match it. */
    async login({ email, password, role }) {
      const { token } = await api.login(email, password)
      const actual = decodeJwt(token)?.role
      if (role && actual !== role) {
        throw new Error(`This is a ${actual} account. Switch the role selector to ${cap(actual)} to sign in.`)
      }
      this.notice = ''
      this._start(token, email)
    },

    async register({ email, password, role }) {
      await api.register({ email, password, role })
      const { token } = await api.login(email, password)
      this._start(token, email)
      const actual = decodeJwt(token)?.role
      this.notice = actual !== role ? `Self-registered accounts are created as ${cap(actual)}. An admin can grant other roles.` : ''
    },

    logout() {
      this.token = ''
      this.email = ''
      this.notice = ''
      write(TOKEN_KEY, '')
      write(EMAIL_KEY, '')
    },
  },
})
