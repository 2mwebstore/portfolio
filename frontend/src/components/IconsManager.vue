<script setup>
import { ref, onMounted } from 'vue'
import { apiFetch } from '../api'

const icons = ref([])
const loading = ref(true)
const note = ref('')
const noteType = ref('ok')
const editingId = ref(null)
const editForm = ref({})

const blankNewItem = () => ({ name: '', icon_url: '', category: 'bank', sort_order: 0 })
const newItem = ref(blankNewItem())

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

    <div v-if="loading" class="loading-note">Loading…</div>

    <table v-else class="admin-table">
      <thead>
        <tr><th></th><th>Name</th><th>Category</th><th>Order</th><th></th></tr>
      </thead>
      <tbody>
        <tr v-for="icon in icons" :key="icon.id">
          <template v-if="editingId === icon.id">
            <td colspan="5">
              <div class="new-item-form">
                <input v-model="editForm.name" placeholder="Name" />
                <input v-model="editForm.icon_url" placeholder="Icon URL" />
                <select v-model="editForm.category">
                  <option value="bank">bank</option>
                  <option value="social">social</option>
                </select>
                <input v-model.number="editForm.sort_order" type="number" style="width:70px" />
                <button class="btn primary" @click="saveEdit">Save</button>
                <button class="btn" @click="cancelEdit">Cancel</button>
              </div>
            </td>
          </template>
          <template v-else>
            <td><img :src="icon.icon_url" alt="" /></td>
            <td>{{ icon.name }}</td>
            <td>{{ icon.category }}</td>
            <td>{{ icon.sort_order }}</td>
            <td class="table-actions">
              <button class="btn" @click="startEdit(icon)">Edit</button>
              <button class="btn danger" @click="removeIcon(icon.id)">Delete</button>
            </td>
          </template>
        </tr>
      </tbody>
    </table>

    <div class="new-item-form">
      <input v-model="newItem.name" placeholder="Name" />
      <input v-model="newItem.icon_url" placeholder="Icon URL" />
      <select v-model="newItem.category">
        <option value="bank">bank</option>
        <option value="social">social</option>
      </select>
      <input v-model.number="newItem.sort_order" type="number" placeholder="Order" style="width:70px" />
      <button class="btn primary" @click="createIcon">Add icon</button>
    </div>

    <p v-if="note" class="status-note" :class="noteType">{{ note }}</p>
  </div>
</template>
