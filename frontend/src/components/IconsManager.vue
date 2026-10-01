<script setup>
import { ref, computed, onMounted } from 'vue'
import { apiFetch } from '../api'

const icons = ref([])
const loading = ref(true)
const note = ref('')
const noteType = ref('ok')
const editingId = ref(null)
const editForm = ref({})

const groups = [
  { key: 'bank', title: 'Bank icons' },
  { key: 'social', title: 'Social icons' },
]

const blankNewItem = () => ({ name: '', icon_url: '', link: '', category: 'bank', sort_order: 0 })
const newItem = ref(blankNewItem())

const iconsByCategory = computed(() => {
  const out = { bank: [], social: [] }
  for (const icon of icons.value) (out[icon.category] ||= []).push(icon)
  return out
})

async function load() {
  loading.value = true
  icons.value = await apiFetch('/api/payment-icons')
  loading.value = false
}

function showNote(message, type = 'ok') {
  note.value = message
  noteType.value = type
  setTimeout(() => { note.value = '' }, 3000)
}

async function createIcon() {
  try {
    await apiFetch('/api/payment-icons', { method: 'POST', body: JSON.stringify(newItem.value) })
    newItem.value = blankNewItem()
    showNote('Icon added.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

function startEdit(icon) {
  editingId.value = icon.id
  editForm.value = { ...icon }
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit() {
  try {
    await apiFetch(`/api/payment-icons/${editForm.value.id}`, {
      method: 'PUT',
      body: JSON.stringify(editForm.value),
    })
    editingId.value = null
    showNote('Icon updated.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

async function removeIcon(id) {
  if (!confirm('Delete this icon?')) return
  try {
    await apiFetch(`/api/payment-icons/${id}`, { method: 'DELETE' })
    showNote('Icon deleted.')
    await load()
  } catch (e) {
    showNote(e.message, 'error')
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-card">
    <h2>Bank / Social Icons</h2>
    <p class="hint">These icons appear in the site footer. Add a link to make an icon clickable.</p>

    <div v-if="loading" class="loading-note">Loading…</div>

    <template v-else>
      <section v-for="group in groups" :key="group.key" class="icon-section">
        <h3 class="section-heading">
          {{ group.title }}
          <span class="count">{{ iconsByCategory[group.key].length }}</span>
        </h3>

        <p v-if="!iconsByCategory[group.key].length" class="empty">No icons yet.</p>

        <ul v-else class="icon-list">
          <li v-for="icon in iconsByCategory[group.key]" :key="icon.id" class="icon-item">
            <!-- Edit mode -->
            <div v-if="editingId === icon.id" class="icon-form">
              <div class="form-preview">
                <img v-if="editForm.icon_url" :src="editForm.icon_url" alt="" />
                <span v-else>?</span>
              </div>
              <div class="form-fields">
                <div class="admin-field">
                  <label>Name</label>
                  <input v-model="editForm.name" />
                </div>
                <div class="admin-field">
                  <label>Icon URL</label>
                  <input v-model="editForm.icon_url" />
                </div>
                <div class="admin-field wide">
                  <label>Link</label>
                  <input v-model="editForm.link" placeholder="https://… (optional)" />
                </div>
                <div class="admin-field">
                  <label>Category</label>
                  <select v-model="editForm.category">
                    <option value="bank">Bank</option>
                    <option value="social">Social</option>
                  </select>
                </div>
                <div class="admin-field">
                  <label>Order</label>
                  <input v-model.number="editForm.sort_order" type="number" />
                </div>
                <div class="form-actions wide">
                  <button class="btn" @click="cancelEdit">Cancel</button>
                  <button class="btn primary" @click="saveEdit">Save</button>
                </div>
              </div>
            </div>

            <!-- View mode -->
            <div v-else class="icon-row-view">
              <img class="icon-thumb" :src="icon.icon_url" alt="" />
              <div class="icon-info">
                <div class="icon-name">
                  {{ icon.name }}
                  <span class="order">#{{ icon.sort_order }}</span>
                </div>
                <a v-if="icon.link" class="icon-link" :href="icon.link" target="_blank" rel="noopener">{{ icon.link }}</a>
                <span v-else class="icon-link muted">No link</span>
              </div>
              <div class="table-actions">
                <button class="btn" @click="startEdit(icon)">Edit</button>
                <button class="btn danger" @click="removeIcon(icon.id)">Delete</button>
              </div>
            </div>
          </li>
        </ul>
      </section>
    </template>

    <section class="add-panel">
      <h3 class="section-heading">Add new icon</h3>
      <div class="icon-form">
        <div class="form-preview">
          <img v-if="newItem.icon_url" :src="newItem.icon_url" alt="" />
          <span v-else>+</span>
        </div>
        <form class="form-fields" @submit.prevent="createIcon">
          <div class="admin-field">
            <label>Name</label>
            <input v-model="newItem.name" placeholder="ABA" required />
          </div>
          <div class="admin-field">
            <label>Icon URL</label>
            <input v-model="newItem.icon_url" placeholder="https://…/icon.png" required />
          </div>
          <div class="admin-field wide">
            <label>Link</label>
            <input v-model="newItem.link" placeholder="https://… (optional)" />
          </div>
          <div class="admin-field">
            <label>Category</label>
            <select v-model="newItem.category">
              <option value="bank">Bank</option>
              <option value="social">Social</option>
            </select>
          </div>
          <div class="admin-field">
            <label>Order</label>
            <input v-model.number="newItem.sort_order" type="number" />
          </div>
          <div class="form-actions wide">
            <button class="btn primary" type="submit">Add icon</button>
          </div>
        </form>
      </div>
    </section>

    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </div>
</template>

<style scoped>
.hint{ color: var(--text-dim); font-size: 13px; margin: -8px 0 18px; }

.icon-section{ margin-bottom: 22px; }
.section-heading{
  display:flex; align-items:center; gap: 8px;
  font-size: 14px; color: var(--text); margin-bottom: 10px;
}
.count{
  font-size: 11.5px; font-weight: 600; color: var(--text-dim);
  background: #eef1f4; border-radius: 999px; padding: 1px 8px;
}
.empty{ color: var(--text-dim); font-size: 13px; padding: 10px 0; }

.icon-list{ list-style: none; border: 1px solid var(--border); border-radius: 8px; overflow: hidden; }
.icon-item + .icon-item{ border-top: 1px solid var(--border); }

.icon-row-view{ display:flex; align-items:center; gap: 12px; padding: 10px 14px; }
.icon-row-view:hover{ background: #fafbfc; }
.icon-thumb{
  width: 40px; height: 40px; border-radius: 50%; object-fit: cover; flex-shrink: 0;
  border: 1px solid var(--border);
}
.icon-info{ flex: 1; min-width: 0; }
.icon-name{ font-weight: 600; font-size: 14px; }
.order{ font-weight: 400; font-size: 12px; color: var(--text-dim); margin-left: 4px; }
.icon-link{
  display:block; font-size: 12.5px; color: #2a7fc1;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.icon-link:hover{ text-decoration: underline; }
.icon-link.muted{ color: var(--text-dim); font-style: italic; }

.icon-form{ display:flex; gap: 16px; padding: 14px; background: #fafbfc; }
.form-preview{
  width: 56px; height: 56px; border-radius: 50%; flex-shrink: 0;
  display:flex; align-items:center; justify-content:center;
  background: #fff; border: 1px dashed var(--border); color: var(--text-dim); font-size: 20px;
  overflow: hidden;
}
.form-preview img{ width: 100%; height: 100%; object-fit: cover; }
.form-fields{ flex: 1; display:grid; grid-template-columns: 1fr 1fr; gap: 0 12px; }
.form-fields .wide{ grid-column: 1 / -1; }
.form-fields select{
  width: 100%; padding: 9px 11px; border: 1px solid var(--border); border-radius: 6px;
  font-size: 14px; font-family: var(--sans); background: #fff;
}
.form-actions{ display:flex; justify-content:flex-end; gap: 8px; }

.add-panel{ border-top: 1px dashed var(--border); padding-top: 18px; }
.add-panel .icon-form{ border: 1px solid var(--border); border-radius: 8px; }

@media (max-width: 560px){
  .icon-form{ flex-direction: column; }
  .form-fields{ grid-template-columns: 1fr; }
  .icon-row-view{ flex-wrap: wrap; }
}
</style>
