<script setup>
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '../components/layout/AuthLayout.vue'
import PasswordField from '../components/ui/PasswordField.vue'
import PasswordRules from '../components/ui/PasswordRules.vue'
import RoleToggle from '../components/ui/RoleToggle.vue'
import { useAuthStore } from '../stores/auth.js'
import { isEmailValid, isPasswordValid, PASSWORD_POLICY_MESSAGE } from '../utils/validators.js'

const auth = useAuthStore()
const router = useRouter()

const form = reactive({ email: '', password: '', confirm: '', role: 'customer' })
const touched = reactive({ email: false, password: false, confirm: false })
const submitted = ref(false)
const busy = ref(false)
const serverError = ref('')

const errors = computed(() => ({
  email: !form.email ? 'Email is required' : !isEmailValid(form.email) ? 'Enter a valid email address' : '',
  password: !form.password ? 'Password is required' : !isPasswordValid(form.password) ? PASSWORD_POLICY_MESSAGE : '',
  confirm: form.confirm !== form.password ? 'Passwords do not match' : '',
}))
const shown = (k) => ((touched[k] || submitted.value) && errors.value[k]) || ''

async function submit() {
  submitted.value = true
  serverError.value = ''
  if (Object.values(errors.value).some(Boolean)) return
  busy.value = true
  try {
    await auth.register({ email: form.email, password: form.password, role: form.role })
    router.replace('/')
  } catch (e) {
    serverError.value = e.message || 'Unable to create account'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <AuthLayout title="Open your account" subtitle="Start managing plans, invoices and renewals.">
    <form novalidate @submit.prevent="submit">
      <div class="field">
        <label>Account type</label>
        <RoleToggle v-model="form.role" />
      </div>

      <div class="field" :class="{ invalid: shown('email') }">
        <label for="email">Email</label>
        <input id="email" v-model.trim="form.email" type="email" autocomplete="email" placeholder="you@company.com" @blur="touched.email = true" />
        <p v-if="shown('email')" class="field-error">{{ shown('email') }}</p>
      </div>

      <PasswordField id="password" v-model="form.password" autocomplete="new-password" :error="shown('password')" @blur="touched.password = true" />
      <PasswordRules :password="form.password" />

      <PasswordField id="confirm" v-model="form.confirm" label="Confirm password" autocomplete="new-password" :error="shown('confirm')" @blur="touched.confirm = true" />

      <p v-if="serverError" class="form-error" role="alert">{{ serverError }}</p>

      <button class="btn btn-primary btn-block" :disabled="busy">{{ busy ? 'Creating account…' : 'Create account' }}</button>
    </form>

    <p class="auth-switch">Already have an account? <RouterLink to="/login">Sign in</RouterLink></p>
  </AuthLayout>
</template>
