<template>
  <t-layout>
    <t-header class="app-header">
      <div class="header-content">
        <div class="header-left">
          <router-link to="/" class="logo">
            <h2>Vikunja</h2>
          </router-link>
        </div>
        
        <div class="header-center">
          <slot name="header" />
        </div>
        
        <div class="header-right">
          <t-dropdown>
            <t-button variant="text">
              <UserIcon />
              {{ authStore.user?.name || 'User' }}
              <ChevronDownIcon />
            </t-button>
            <template #dropdown>
              <t-dropdown-menu>
                <t-dropdown-item @click="authStore.logout(); $router.push('/login')">
                  <LogoutIcon /> Logout
                </t-dropdown-item>
              </t-dropdown-menu>
            </template>
          </t-dropdown>
        </div>
      </div>
    </t-header>
    
    <t-layout>
      <t-aside class="app-aside">
        <t-menu :value="currentRoute" theme="light">
          <t-menu-item value="/" @click="$router.push('/')">
            <template #icon><HomeIcon /></template>
            Dashboard
          </t-menu-item>
          <t-menu-item value="/projects" @click="$router.push('/projects')">
            <template #icon><FolderIcon /></template>
            Projects
          </t-menu-item>
        </t-menu>
      </t-aside>
      
      <t-content class="app-content">
        <slot />
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { UserIcon, ChevronDownIcon, LogoutIcon, HomeIcon, FolderIcon } from 'tdesign-icons-vue-next'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const currentRoute = computed(() => route.path)
</script>

<style scoped>
.app-header {
  background: white;
  border-bottom: 1px solid #e7e7e7;
  padding: 0 24px;
  height: 64px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.header-left {
  flex: 0 0 auto;
}

.header-center {
  flex: 1;
  padding: 0 24px;
}

.header-right {
  flex: 0 0 auto;
}

.logo {
  text-decoration: none;
  color: #0052d9;
}

.logo h2 {
  margin: 0;
  font-weight: 700;
}

.app-aside {
  width: 200px;
  background: white;
  border-right: 1px solid #e7e7e7;
}

.app-content {
  background: #f5f5f5;
  min-height: calc(100vh - 64px);
}
</style>