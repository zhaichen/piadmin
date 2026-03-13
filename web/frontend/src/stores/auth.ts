import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { hasToken, clearToken } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(hasToken())

  const isAuthenticated = computed(() => authenticated.value)

  function setAuthenticated(val: boolean) {
    authenticated.value = val
  }

  function logout() {
    clearToken()
    authenticated.value = false
  }

  return { isAuthenticated, setAuthenticated, logout }
})
