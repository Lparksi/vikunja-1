<template>
  <div class="login-container">
    <t-card class="login-card">
      <template #header>
        <div class="login-header">
          <h1>Vikunja</h1>
          <p>Welcome back! Please sign in to your account.</p>
        </div>
      </template>
      
      <t-form ref="formRef" :model="formData" :rules="rules" @submit="handleLogin">
        <t-form-item name="username">
          <t-input
            v-model="formData.username"
            placeholder="Username or Email"
            size="large"
            :prefixIcon="UserIcon"
          />
        </t-form-item>
        
        <t-form-item name="password">
          <t-input
            v-model="formData.password"
            type="password"
            placeholder="Password"
            size="large"
            :prefixIcon="LockOnIcon"
          />
        </t-form-item>
        
        <t-form-item>
          <t-checkbox v-model="formData.long_token">
            Keep me logged in for 30 days
          </t-checkbox>
        </t-form-item>
        
        <t-form-item>
          <t-button
            type="submit"
            theme="primary"
            size="large"
            block
            :loading="loading"
          >
            Sign In
          </t-button>
        </t-form-item>
      </t-form>
      
      <div class="login-footer">
        <p>
          Don't have an account?
          <router-link to="/register">Sign up here</router-link>
        </p>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { UserIcon, LockOnIcon } from 'tdesign-icons-vue-next'
import { useAuthStore } from '@/stores/auth'
import type { UserLogin } from '@/types'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const formRef = ref()

const formData = reactive<UserLogin & { long_token: boolean }>({
  username: '',
  password: '',
  long_token: false
})

const rules = {
  username: [
    { required: true, message: 'Username or email is required' }
  ],
  password: [
    { required: true, message: 'Password is required' }
  ]
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  loading.value = true
  try {
    await authStore.login({
      username: formData.username,
      password: formData.password,
      long_token: formData.long_token
    })
    
    MessagePlugin.success('Login successful!')
    router.push('/')
  } catch (error: any) {
    MessagePlugin.error(error.response?.data?.detail || 'Login failed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
}

.login-header h1 {
  color: #333;
  margin-bottom: 8px;
  font-size: 32px;
  font-weight: 700;
}

.login-header p {
  color: #666;
  font-size: 14px;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.login-footer a {
  color: #0052d9;
  text-decoration: none;
}

.login-footer a:hover {
  text-decoration: underline;
}
</style>