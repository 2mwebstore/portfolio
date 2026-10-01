const TOKEN_KEY = 'admin_token'
const USERNAME_KEY = 'admin_username'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function isLoggedIn() {
  return !!getToken()
}

export function getUsername() {
  return localStorage.getItem(USERNAME_KEY) || ''
}

export function setSession(token, username) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USERNAME_KEY, username)
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USERNAME_KEY)
}

// apiFetch: same signature as fetch(), but adds the Authorization
// header when logged in, and throws a readable Error on non-2xx.
export async function apiFetch(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(path, { ...options, headers })

  if (res.status === 401) {
    clearSession()
    throw new Error('Session expired — please log in again.')
  }

  if (!res.ok) {
    let message = `Request failed (${res.status})`
    try {
      const body = await res.json()
      if (body.error) message = body.error
    } catch (_) { /* ignore parse errors */ }
    throw new Error(message)
  }

  if (res.status === 204) return null
  return res.json()
}

// uploadImage: posts a File to /api/upload and returns { url }.
// Separate from apiFetch because multipart requests must NOT set
// Content-Type manually — the browser needs to add its own boundary.
export async function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)

  const token = getToken()
  const headers = token ? { Authorization: `Bearer ${token}` } : {}

  const res = await fetch('/api/upload', { method: 'POST', headers, body: formData })

  if (res.status === 401) {
    clearSession()
    throw new Error('Session expired — please log in again.')
  }
  if (!res.ok) {
    let message = `Upload failed (${res.status})`
    try {
      const body = await res.json()
      if (body.error) message = body.error
    } catch (_) { /* ignore parse errors */ }
    throw new Error(message)
  }
  return res.json()
}
