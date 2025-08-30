import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User, UserLogin, UserCreate } from '@/types'
import apiService from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('vikunja_token'))
  const isLoggedIn = ref<boolean>(!!token.value)

  const login = async (credentials: UserLogin) => {
    try {
      const tokenData = await apiService.login(credentials)
      token.value = tokenData.token
      localStorage.setItem('vikunja_token', tokenData.token)
      isLoggedIn.value = true
      
      // Test the token to get user info
      const testResponse = await apiService.testToken()
      // For now, create a minimal user object - in real app you'd fetch full user info
      user.value = {
        id: 0,
        username: testResponse.user,
        email: '',
        name: testResponse.user,
        timezone: null,
        week_start: 0,
        language: null,
        is_active: true,
        is_admin: false,
        is_email_confirmed: false,
        totp_enabled: false,
        avatar_provider: 'initials',
        avatar_file_id: null,
        created_at: '',
        updated_at: ''
      }
      
      return tokenData
    } catch (error) {
      console.error('Login error:', error)
      throw error
    }
  }

  const register = async (userData: UserCreate) => {
    try {
      const newUser = await apiService.register(userData)
      // Auto-login after registration
      await login({ username: userData.username, password: userData.password })
      return newUser
    } catch (error) {
      console.error('Registration error:', error)
      throw error
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('vikunja_token')
    isLoggedIn.value = false
  }

  const checkAuth = async () => {
    if (!token.value) {
      isLoggedIn.value = false
      return false
    }

    try {
      const testResponse = await apiService.testToken()
      isLoggedIn.value = true
      if (!user.value) {
        user.value = {
          id: 0,
          username: testResponse.user,
          email: '',
          name: testResponse.user,
          timezone: null,
          week_start: 0,
          language: null,
          is_active: true,
          is_admin: false,
          is_email_confirmed: false,
          totp_enabled: false,
          avatar_provider: 'initials',
          avatar_file_id: null,
          created_at: '',
          updated_at: ''
        }
      }
      return true
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    user,
    token,
    isLoggedIn,
    login,
    register,
    logout,
    checkAuth
  }
})