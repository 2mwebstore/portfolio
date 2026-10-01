<script setup>
import { ref, onMounted } from 'vue'
import { apiFetch, uploadImage } from '../api'

const banners = ref([])
const loading = ref(true)
const note = ref('')
const noteType = ref('ok')
const editingId = ref(null)
const editForm = ref({})
const uploadingField = ref('') // e.g. 'new-image_url', 'edit-image_url'

const blankNewItem = () => ({ image_url: '', link: '', sort_order: 0, active: true })
const newItem = ref(blankNewItem())

async function load() {
  loading.value = true
  banners.value = await apiFetch('/api/banners/all')
  loading.value = false
}

function showNote(message, type = 'ok') {
  note.value = message
  noteType.value = type
  setTimeout(() => { note.value = '' }, 3000)
}

async function handleUpload(event, target, field, tag) {
  const file = event.target.files[0]
  if (!file) return
  uploadingField.value = tag
  try {
    const { url } = await uploadImage(file)
    target[field] = url
    showNote('Image uploaded.')
  } catch (e) {
    showNote(e.message, 'error')
  } finally {
    uploadingField.value = ''
    event.target.value = ''
  }
}

async function createBanner() {
  if (!newItem.value.image_url) {
    showNote('Add an image first.', 'error')
    return
  }
  try {
    await apiFetch('/api/banners', { method: 'POST', body: JSON.stringify(newItem.value) })
    newItem.value = blankNewItem()
    showNote('Banner added.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

function startEdit(b) {
  editingId.value = b.id
  editForm.value = { ...b }
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit() {
  try {
    await apiFetch(`/api/banners/${editForm.value.id}`, {
      method: 'PUT',
      body: JSON.stringify(editForm.value),
    })
    editingId.value = null
    showNote('Banner updated.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

async function toggleActive(b) {
  try {
    await apiFetch(`/api/banners/${b.id}`, {
      method: 'PUT',
      body: JSON.stringify({ ...b, active: !b.active }),
    })
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

async function removeBanner(id) {
  if (!confirm('Delete this banner?')) return
  try {
    await apiFetch(`/api/banners/${id}`, { method: 'DELETE' })
    showNote('Banner deleted.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-card">
    <h2>Homepage Banners</h2>
    <p class="hint">These rotate automatically in the hero swiper on the public site. Order is left to right; toggle Active to hide one without deleting it.</p>

    <div v-if="loading" class="loading-note">Loading…</div>

    <table v-else class="admin-table">
      <thead>
        <tr>
          <th></th><th>Link</th><th>Order</th><th>Active</th><th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="b in banners" :key="b.id">
          <template v-if="editingId === b.id">
            <td colspan="5">
              <div class="new-item-form">
                <input v-model="editForm.link" placeholder="Link (or leave blank / #)" />
                <input v-model.number="editForm.sort_order" type="number" placeholder="Order" style="width:70px" />
                <label class="active-toggle">
                  <input type="checkbox" v-model="editForm.active" /> Active
                </label>
              </div>
              <div class="new-item-form">
                <label class="upload-field">
                  Banner image:
                  <input v-model="editForm.image_url" placeholder="Image URL" />
                  <input type="file" accept="image/*" @change="e => handleUpload(e, editForm, 'image_url', 'edit-image_url')" />
                  <span v-if="uploadingField === 'edit-image_url'" class="upload-note">Uploading…</span>
                </label>
              </div>
              <div class="new-item-form">
                <button class="btn primary" @click="saveEdit">Save</button>
                <button class="btn" @click="cancelEdit">Cancel</button>
              </div>
            </td>
          </template>
          <template v-else>
            <td><img class="banner-thumb" :src="b.image_url" alt="" /></td>
            <td>{{ b.link || '—' }}</td>
            <td>{{ b.sort_order }}</td>
            <td>
              <button class="btn" :class="{ primary: b.active }" @click="toggleActive(b)">
                {{ b.active ? 'Active' : 'Hidden' }}
              </button>
            </td>
            <td class="table-actions">
              <button class="btn" @click="startEdit(b)">Edit</button>
              <button class="btn danger" @click="removeBanner(b.id)">Delete</button>
            </td>
          </template>
        </tr>
      </tbody>
    </table>

    <div class="new-item-form">
      <input v-model="newItem.link" placeholder="Link (or leave blank / #)" />
      <input v-model.number="newItem.sort_order" type="number" placeholder="Order" style="width:70px" />
    </div>
    <div class="new-item-form">
      <label class="upload-field">
        Banner image:
        <input v-model="newItem.image_url" placeholder="Image URL" />
        <input type="file" accept="image/*" @change="e => handleUpload(e, newItem, 'image_url', 'new-image_url')" />
        <span v-if="uploadingField === 'new-image_url'" class="upload-note">Uploading…</span>
      </label>
    </div>
    <div class="new-item-form">
      <button class="btn primary" @click="createBanner">Add banner</button>
    </div>

    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </div>
</template>

<style scoped>
.hint{ color: var(--text-dim); font-size: 13px; margin: -6px 0 14px; }
.banner-thumb{ width: 90px; height: 45px; object-fit: cover; border-radius: 4px; }
.active-toggle{ display:flex; align-items:center; gap:6px; font-size: 13px; color: var(--text-dim); }
.upload-field{
  display:flex; flex-direction:column; gap:5px;
  font-size: 12.5px; color: var(--text-dim);
  background: var(--bg); border: 1px solid var(--border); border-radius: 8px;
  padding: 10px 12px;
}
.upload-field input[type="text"], .upload-field input:not([type]){ font-size: 13px; }
.upload-note{ color: var(--accent); font-size: 12px; }
</style>
