<template>
  <div class="home">
    <AppLayout>
      <template #header>
        <h1>Dashboard</h1>
      </template>
      
      <div class="dashboard-content">
        <t-row :gutter="16">
          <t-col :span="8">
            <t-card title="Welcome">
              <p>Welcome to Vikunja, {{ authStore.user?.name || 'User' }}!</p>
              <p>This is your personal todo and project management dashboard.</p>
              <div class="quick-actions">
                <t-button theme="primary" @click="$router.push('/projects')">
                  View Projects
                </t-button>
              </div>
            </t-card>
          </t-col>
          
          <t-col :span="8">
            <t-card title="Recent Projects">
              <div v-if="projectStore.loading">
                <t-skeleton />
              </div>
              <div v-else-if="recentProjects.length === 0">
                <p>No projects yet. <router-link to="/projects">Create your first project</router-link></p>
              </div>
              <div v-else>
                <div v-for="project in recentProjects" :key="project.id" class="project-item">
                  <router-link :to="`/projects/${project.id}`" class="project-link">
                    <t-tag v-if="project.color" :color="project.color" variant="light">
                      {{ project.title }}
                    </t-tag>
                    <span v-else>{{ project.title }}</span>
                  </router-link>
                </div>
              </div>
            </t-card>
          </t-col>
          
          <t-col :span="8">
            <t-card title="Quick Stats">
              <div class="stats">
                <div class="stat-item">
                  <span class="stat-number">{{ projectStore.projects.length }}</span>
                  <span class="stat-label">Projects</span>
                </div>
                <div class="stat-item">
                  <span class="stat-number">{{ totalTasks }}</span>
                  <span class="stat-label">Total Tasks</span>
                </div>
              </div>
            </t-card>
          </t-col>
        </t-row>
      </div>
    </AppLayout>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'
import { useTaskStore } from '@/stores/tasks'
import AppLayout from '@/components/AppLayout.vue'

const authStore = useAuthStore()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const recentProjects = computed(() => {
  return projectStore.projects.slice(0, 5)
})

const totalTasks = computed(() => {
  return taskStore.tasks.length
})

onMounted(async () => {
  await projectStore.fetchProjects()
})
</script>

<style scoped>
.dashboard-content {
  padding: 20px;
}

.quick-actions {
  margin-top: 16px;
}

.project-item {
  margin-bottom: 8px;
}

.project-link {
  text-decoration: none;
  color: inherit;
}

.project-link:hover {
  color: #0052d9;
}

.stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #0052d9;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
}
</style>