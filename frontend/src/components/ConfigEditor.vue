<script setup>
import { ref, onMounted } from 'vue'
import { apiFetch } from '../api'

const form = ref(null)
const loading = ref(true)
const saving = ref(false)
const note = ref('')
const noteType = ref('ok')

async function load() {
  loading.value = true
  form.value = await apiFetch('/api/config')
  loading.value = false
}

async function save() {
  saving.value = true
  note.value = ''
  try {
    await apiFetch('/api/config', { method: 'PUT', body: JSON.stringify(form.value) })
    note.value = 'Saved.'
    noteType.value = 'ok'
  } catch (e) {
    note.value = e.message
    noteType.value = 'error'
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-if="loading" class="loading-note">Loading…</div>
  <form v-else class="admin-card" @submit.prevent="save">
    <h2>Site Config</h2>

    <div class="admin-row">
      <div class="admin-field">
        <label>Site name</label>
        <input v-model="form.site_name" required />
      </div>
      <div class="admin-field">
        <label>Logo URL</label>
        <input v-model="form.logo_url" placeholder="/uploads/logo.png" />
      </div>
    </div>

    <div class="admin-row">
      <div class="admin-field">
        <label>Telegram URL</label>
        <input v-model="form.telegram_url" placeholder="https://t.me/yourchannel" />
      </div>
      <div class="admin-field">
        <label>Facebook URL</label>
        <input v-model="form.facebook_url" placeholder="https://facebook.com/yourpage" />
      </div>
    </div>

    <div class="admin-row">
      <div class="admin-field">
        <label>Telegram icon URL</label>
        <input v-model="form.telegram_icon_url" placeholder="/uploads/telegram.png" />
      </div>
      <div class="admin-field">
        <label>Telegram button color</label>
        <input v-model="form.telegram_color" type="color" style="height:38px; padding:2px;" />
      </div>
    </div>

    <div class="admin-row">
      <div class="admin-field">
        <label>Facebook icon URL</label>
        <input v-model="form.facebook_icon_url" placeholder="/uploads/facebook.png" />
      </div>
      <div class="admin-field">
        <label>Facebook button color</label>
        <input v-model="form.facebook_color" type="color" style="height:38px; padding:2px;" />
      </div>
    </div>

    <div class="admin-row">
      <div class="admin-field">
        <label>TikTok URL</label>
        <input v-model="form.tiktok_url" />
      </div>
      <div class="admin-field">
        <label>Instagram URL</label>
        <input v-model="form.instagram_url" />
      </div>
    </div>

    <div class="admin-row">
      <div class="admin-field">
        <label>Background color</label>
        <input v-model="form.background_color" type="color" style="height:38px; padding:2px;" />
      </div>
      <div class="admin-field">
        <label>Background image URL (optional, overlays the color)</label>
        <input v-model="form.background_image_url" placeholder="/uploads/background.jpg" />
      </div>
    </div>

    <div class="admin-field">
      <label>Footer note</label>
      <input v-model="form.footer_note" />
    </div>

    <p class="hint">Looking for the homepage banner? It moved to its own <strong>Banners</strong> tab, where you can add and reorder multiple banners for the swiper.</p>

    <button class="btn primary" type="submit" :disabled="saving">
      {{ saving ? 'Saving…' : 'Save changes' }}
    </button>
    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </form>
</template>

<style scoped>
.hint{ color: var(--text-dim); font-size: 13px; margin-top: -4px; }
</style>