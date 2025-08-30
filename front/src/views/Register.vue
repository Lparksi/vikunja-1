<template>
  <div class="register-container">
    <t-card class="register-card">
      <template #header>
        <div class="register-header">
          <h1>Vikunja</h1>
          <p>Create your account to get started.</p>
        </div>
      </template>
      
      <t-form ref="formRef" :model="formData" :rules="rules" @submit="handleRegister">
        <t-form-item name="username">
          <t-input
            v-model="formData.username"
            placeholder="Username"
            size="large"
            :prefixIcon="UserIcon"
          />
        </t-form-item>
        
        <t-form-item name="email">
          <t-input
            v-model="formData.email"
            placeholder="Email"
            size="large"
            :prefixIcon="MailIcon"
          />
        </t-form-item>
        
        <t-form-item name="name">
          <t-input
            v-model="formData.name"
            placeholder="Full Name"
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
        
        <t-form-item name="confirmPassword">
          <t-input
            v-model="formData.confirmPassword"
            type="password"
            placeholder="Confirm Password"
            size="large"
            :prefixIcon="LockOnIcon"
          />
        </t-form-item>
        
        <t-form-item>
          <t-button
            type="submit"
            theme="primary"
            size="large"
            block
            :loading="loading"
          >
            Create Account
          </t-button>
        </t-form-item>
      </t-form>
      
      <div class="register-footer">
        <p>
          Already have an account?
          <router-link to="/login">Sign in here</router-link>
        </p>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { UserIcon, LockOnIcon, MailIcon } from 'tdesign-icons-vue-next'
import { useAuthStore } from '@/stores/auth'
import type { UserCreate } from '@/types'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const formRef = ref()

const formData = reactive<UserCreate & { confirmPassword: string }>({
  username: '',
  email: '',
  name: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  username: [
    { required: true, message: 'Username is required' },
    { min: 3, message: 'Username must be at least 3 characters' }
  ],
  email: [
    { required: true, message: 'Email is required' },
    { email: true, message: 'Please enter a valid email' }
  ],
  name: [
    { required: true, message: 'Name is required' }
  ],
  password: [
    { required: true, message: 'Password is required' },
    { min: 6, message: 'Password must be at least 6 characters' }
  ],
  confirmPassword: [
    { required: true, message: 'Please confirm your password' },
    {
      validator: (val: string) => val === formData.password,
      message: 'Passwords do not match'
    }
  ]
}

const handleRegister = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  loading.value = true
  try {
    const newUser = await authStore.register({
      username: formData.username,
      email: formData.email,
      name: formData.name,
      password: formData.password
    })
    
    MessagePlugin.success('Account created successfully!')
    router.push('/')
  } catch (error: any) {
    MessagePlugin.error(error.response?.data?.detail || 'Registration failed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 450px;
}

.register-header {
  text-align: center;
  margin-bottom: 20px;
}

.register-header h1 {
  color: #333;
  margin-bottom: 8px;
  font-size: 32px;
  font-weight: 700;
}

.register-header p {
  color: #666;
  font-size: 14px;
}

.register-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.register-footer a {
  color: #0052d9;
  text-decoration: none;
}

.register-footer a:hover {
  text-decoration: underline;
}
</style>