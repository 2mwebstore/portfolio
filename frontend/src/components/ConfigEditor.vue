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

    <h3 class="section-heading">Social buttons</h3>
    <div class="social-panels">
      <div class="social-panel">
        <div class="panel-head">
          <img v-if="form.telegram_icon_url" :src="form.telegram_icon_url" alt="" />
          <span>Telegram</span>
        </div>
        <div class="admin-field">
          <label>Title</label>
          <input v-model="form.telegram_label" placeholder="តេលេក្រាមផ្លូវការ" />
        </div>
        <div class="admin-field">
          <label>Link</label>
          <input v-model="form.telegram_url" placeholder="https://t.me/yourchannel" />
        </div>
        <div class="admin-field">
          <label>Icon URL</label>
          <input v-model="form.telegram_icon_url" placeholder="/uploads/telegram.png" />
        </div>
        <div class="admin-field">
          <label>Button color</label>
          <input v-model="form.telegram_color" type="color" class="color-input" />
        </div>
      </div>

      <div class="social-panel">
        <div class="panel-head">
          <img v-if="form.facebook_icon_url" :src="form.facebook_icon_url" alt="" />
          <span>Facebook</span>
        </div>
        <div class="admin-field">
          <label>Title</label>
          <input v-model="form.facebook_label" placeholder="ហ្វេសប៊ុកផ្លូវការ" />
        </div>
        <div class="admin-field">
          <label>Link</label>
          <input v-model="form.facebook_url" placeholder="https://facebook.com/yourpage" />
        </div>
        <div class="admin-field">
          <label>Icon URL</label>
          <input v-model="form.facebook_icon_url" placeholder="/uploads/facebook.png" />
        </div>
        <div class="admin-field">
          <label>Button color</label>
          <input v-model="form.facebook_color" type="color" class="color-input" />
        </div>
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
.section-heading{ font-size: 14px; color: var(--text); margin: 6px 0 10px; }
.social-panels{ display:grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 16px; }
@media (max-width: 560px){ .social-panels{ grid-template-columns: 1fr; } }
.social-panel{ border: 1px solid var(--border); border-radius: 8px; padding: 14px; background: #fafbfc; }
.social-panel .admin-field:last-child{ margin-bottom: 0; }
.panel-head{ display:flex; align-items:center; gap: 8px; font-weight: 600; font-size: 14px; margin-bottom: 12px; }
.panel-head img{ width: 24px; height: 24px; border-radius: 50%; object-fit: cover; }
.color-input{ height: 38px; padding: 2px !important; cursor: pointer; }
</style>