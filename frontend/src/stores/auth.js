import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)

  async function register(username, email, password) {
    try {
      const data = await api.post('/auth/register', {
        username,
        email,
        password,
      })
      token.value = data.token
      user.value = data.user
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Registration failed' }
    }
  }

  async function login(username, password) {
    try {
      const data = await api.post('/auth/login', {
        username,
        password,
      })
      token.value = data.token
      user.value = data.user
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Login failed' }
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchProfile() {
    try {
      const data = await api.get('/user/profile')
      user.value = data.user
      localStorage.setItem('user', JSON.stringify(data.user))
    } catch (error) {
      console.error('Failed to fetch profile:', error)
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    register,
    login,
    logout,
    fetchProfile,
  }
})
