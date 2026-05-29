import { ref } from 'vue'
import { defineStore } from 'pinia'
import { login as apiLogin } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const ready = ref(false)

  async function ensureAuth() {
    if (token.value) { ready.value = true; return }
    try {
      const password = import.meta.env.VITE_ADMIN_PASSWORD
      const res = await apiLogin(password)
      token.value = res.data.token
      localStorage.setItem('token', token.value)
    } catch {
      console.error('auto-auth failed')
    }
    ready.value = true
  }

  return { token, ready, ensureAuth }
})
