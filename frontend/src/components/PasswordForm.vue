<script setup>
import { ref } from 'vue'
import { apiFetch } from '../api'

const currentPassword = ref('')
const newPassword = ref('')
const note = ref('')
const noteType = ref('ok')
const saving = ref(false)

async function submit() {
  saving.value = true
  note.value = ''
  try {
    await apiFetch('/api/auth/password', {
      method: 'PUT',
      body: JSON.stringify({
        current_password: currentPassword.value,
        new_password: newPassword.value,
      }),
    })
    note.value = 'Password updated.'
    noteType.value = 'ok'
    currentPassword.value = ''
    newPassword.value = ''
  } catch (e) {
    note.value = e.message
    noteType.value = 'error'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <form class="admin-card" @submit.prevent="submit">
    <h2>Change Password</h2>
    <div class="admin-field">
      <label>Current password</label>
      <input v-model="currentPassword" type="password" required />
    </div>
    <div class="admin-field">
      <label>New password (min 8 characters)</label>
      <input v-model="newPassword" type="password" minlength="8" required />
    </div>
    <button class="btn primary" type="submit" :disabled="saving">
      {{ saving ? 'Updating…' : 'Update password' }}
    </button>
    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </form>
</template>
