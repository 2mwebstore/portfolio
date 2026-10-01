<script setup>
import { ref, onMounted } from 'vue'
import { apiFetch, uploadImage } from '../api'

const partners = ref([])
const loading = ref(true)
const note = ref('')
const noteType = ref('ok')
const editingId = ref(null)
const editForm = ref({})
const uploadingField = ref('') // e.g. 'new-logo_url', 'edit-media_image_url'

const blankNewItem = () => ({ name: '', logo_url: '', media_image_url: '', link: '', detail: '', sort_order: 0 })
const newItem = ref(blankNewItem())

async function load() {
  loading.value = true
  partners.value = await apiFetch('/api/partners')
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

async function createPartner() {
  try {
    await apiFetch('/api/partners', { method: 'POST', body: JSON.stringify(newItem.value) })
    newItem.value = blankNewItem()
    showNote('Partner added.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

function startEdit(p) {
  editingId.value = p.id
  editForm.value = { ...p }
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit() {
  try {
    await apiFetch(`/api/partners/${editForm.value.id}`, {
      method: 'PUT',
      body: JSON.stringify(editForm.value),
    })
    editingId.value = null
    showNote('Partner updated.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

async function removePartner(id) {
  if (!confirm('Delete this partner?')) return
  try {
    await apiFetch(`/api/partners/${id}`, { method: 'DELETE' })
    showNote('Partner deleted.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-card">
    <h2>Partners</h2>

    <div v-if="loading" class="loading-note">Loading…</div>

    <table v-else class="admin-table">
      <thead>
        <tr>
          <th></th><th>Name</th><th>Detail</th><th>Link</th><th>Order</th><th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="p in partners" :key="p.id">
          <template v-if="editingId === p.id">
            <td colspan="6">
              <div class="new-item-form">
                <input v-model="editForm.name" placeholder="Name" />
                <input v-model="editForm.link" placeholder="Link" />
                <input v-model="editForm.detail" placeholder="Detail" />
                <input v-model.number="editForm.sort_order" type="number" style="width:70px" />
              </div>
              <div class="new-item-form">
                <label class="upload-field">
                  Logo:
                  <input v-model="editForm.logo_url" placeholder="Logo URL" />
                  <input type="file" accept="image/*" @change="e => handleUpload(e, editForm, 'logo_url', 'edit-logo_url')" />
                  <span v-if="uploadingField === 'edit-logo_url'" class="upload-note">Uploading…</span>
                </label>
                <label class="upload-field">
                  Media image:
                  <input v-model="editForm.media_image_url" placeholder="Media image URL" />
                  <input type="file" accept="image/*" @change="e => handleUpload(e, editForm, 'media_image_url', 'edit-media_image_url')" />
                  <span v-if="uploadingField === 'edit-media_image_url'" class="upload-note">Uploading…</span>
                </label>
              </div>
              <div class="new-item-form">
                <button class="btn primary" @click="saveEdit">Save</button>
                <button class="btn" @click="cancelEdit">Cancel</button>
              </div>
            </td>
          </template>
          <template v-else>
            <td><img :src="p.logo_url" alt="" /></td>
            <td>{{ p.name }}</td>
            <td>{{ p.detail }}</td>
            <td>{{ p.link }}</td>
            <td>{{ p.sort_order }}</td>
            <td class="table-actions">
              <button class="btn" @click="startEdit(p)">Edit</button>
              <button class="btn danger" @click="removePartner(p.id)">Delete</button>
            </td>
          </template>
        </tr>
      </tbody>
    </table>

    <div class="new-item-form">
      <input v-model="newItem.name" placeholder="Name" />
      <input v-model="newItem.link" placeholder="Link" />
      <input v-model="newItem.detail" placeholder="Detail" />
      <input v-model.number="newItem.sort_order" type="number" placeholder="Order" style="width:70px" />
    </div>
    <div class="new-item-form">
      <label class="upload-field">
        Logo:
        <input v-model="newItem.logo_url" placeholder="Logo URL" />
        <input type="file" accept="image/*" @change="e => handleUpload(e, newItem, 'logo_url', 'new-logo_url')" />
        <span v-if="uploadingField === 'new-logo_url'" class="upload-note">Uploading…</span>
      </label>
      <label class="upload-field">
        Media image:
        <input v-model="newItem.media_image_url" placeholder="Media image URL" />
        <input type="file" accept="image/*" @change="e => handleUpload(e, newItem, 'media_image_url', 'new-media_image_url')" />
        <span v-if="uploadingField === 'new-media_image_url'" class="upload-note">Uploading…</span>
      </label>
    </div>
    <div class="new-item-form">
      <button class="btn primary" @click="createPartner">Add partner</button>
    </div>

    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </div>
</template>

<style scoped>
.upload-field{
  display:flex; flex-direction:column; gap:5px;
  font-size: 12.5px; color: var(--text-dim);
  background: var(--bg); border: 1px solid var(--border); border-radius: 8px;
  padding: 10px 12px;
}
.upload-field input[type="text"], .upload-field input:not([type]){ font-size: 13px; }
.upload-note{ color: var(--accent); font-size: 12px; }
</style>
