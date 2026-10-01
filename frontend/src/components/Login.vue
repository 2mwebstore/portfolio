<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch, setSession } from '../api'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const res = await apiFetch('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username: username.value, password: password.value }),
    })
    setSession(res.token, res.username)
    router.push({ name: 'admin' })
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <form class="login-card" @submit.prevent="submit">
      <h1>Admin Login</h1>
      <div class="field">
        <label for="username">Username</label>
        <input id="username" v-model="username" type="text" autocomplete="username" required />
      </div>
      <div class="field">
        <label for="password">Password</label>
        <input id="password" v-model="password" type="password" autocomplete="current-password" required />
      </div>
      <p v-if="error" class="error-note">{{ error }}</p>
      <button type="submit" class="btn-primary" :disabled="loading">
        {{ loading ? 'Signing in…' : 'Sign in' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login-page{
  min-height: 100vh;
  display: flex; align-items: center; justify-content: center;
  background: var(--bg, #f4f6f8);
}
.login-card{
  width: 100%; max-width: 340px;
  background: #fff; border: 1px solid #e3e7ec; border-radius: 12px;
  padding: 32px 28px;
  box-shadow: 0 20px 50px -30px rgba(22,41,79,0.4);
}
.login-card h1{ font-size: 20px; margin-bottom: 20px; color: #16294f; }
.field{ margin-bottom: 16px; }
.field label{ display:block; font-size: 13px; color:#66707d; margin-bottom: 6px; }
.field input{
  width: 100%; padding: 10px 12px; border: 1px solid #e3e7ec; border-radius: 6px;
  font-size: 14px;
}
.error-note{ color: #c0392b; font-size: 13px; margin-bottom: 14px; }
.btn-primary{
  width: 100%; padding: 11px; border: none; border-radius: 6px;
  background: #e8863a; color: #fff; font-weight: 600; font-size: 14px; cursor: pointer;
}
.btn-primary:disabled{ opacity: 0.6; cursor: default; }
</style>
