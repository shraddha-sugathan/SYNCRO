<script setup>
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '../components/layout/AuthLayout.vue'
import PasswordField from '../components/ui/PasswordField.vue'
import RoleToggle from '../components/ui/RoleToggle.vue'
import { useAuthStore } from '../stores/auth.js'
import { isEmailValid, isPasswordValid, PASSWORD_POLICY_MESSAGE } from '../utils/validators.js'
import { DEMO_ACCOUNTS } from '../services/mock/db.js'

const auth = useAuthStore()
const router = useRouter()

const form = reactive({ email: '', password: '', role: 'customer' })
const touched = reactive({ email: false, password: false })
const submitted = ref(false)
const busy = ref(false)
const serverError = ref('')

const errors = computed(() => ({
  email: !form.email ? 'Email is required' : !isEmailValid(form.email) ? 'Enter a valid email address' : '',
  password: !form.password ? 'Password is required' : !isPasswordValid(form.password) ? PASSWORD_POLICY_MESSAGE : '',
}))
const shown = (k) => ((touched[k] || submitted.value) && errors.value[k]) || ''

function fillDemo(role) {
  Object.assign(form, { role, ...DEMO_ACCOUNTS[role] })
  serverError.value = ''
}

async function submit() {
  submitted.value = true
  serverError.value = ''
  if (errors.value.email || errors.value.password) return
  busy.value = true
  try {
    await auth.login({ ...form })
    router.replace('/')
  } catch (e) {
    serverError.value = e.message || 'Unable to sign in'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <AuthLayout title="Welcome back" subtitle="Sign in to manage your subscriptions.">
    <form novalidate @submit.prevent="submit">
      <div class="field">
        <label>Sign in as</label>
        <RoleToggle v-model="form.role" />
      </div>

      <div class="field" :class="{ invalid: shown('email') }">
        <label for="email">Email</label>
        <input id="email" v-model.trim="form.email" type="email" autocomplete="email" placeholder="you@company.com" @blur="touched.email = true" />
        <p v-if="shown('email')" class="field-error">{{ shown('email') }}</p>
      </div>

      <PasswordField id="password" v-model="form.password" :error="shown('password')" @blur="touched.password = true" />

      <p v-if="serverError" class="form-error" role="alert">{{ serverError }}</p>

      <button class="btn btn-primary btn-block" :disabled="busy">{{ busy ? 'Signing in…' : 'Sign in' }}</button>
    </form>

    <div class="demo-box">
      <span>Demo accounts</span>
      <button type="button" class="link" @click="fillDemo('admin')">Admin</button>
      <button type="button" class="link" @click="fillDemo('customer')">Customer</button>
    </div>

    <p class="auth-switch">New to Syncro? <RouterLink to="/register">Create an account</RouterLink></p>
  </AuthLayout>
</template>
