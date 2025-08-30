<template>
  <div class="projects">
    <AppLayout>
      <template #header>
        <div class="projects-header">
          <h1>Projects</h1>
          <t-button theme="primary" @click="showCreateDialog = true">
            <template #icon><AddIcon /></template>
            New Project
          </t-button>
        </div>
      </template>
      
      <div class="projects-content">
        <div v-if="projectStore.loading" class="loading">
          <t-skeleton />
        </div>
        
        <div v-else-if="projectStore.projects.length === 0" class="empty-state">
          <div class="empty-content">
            <h3>No projects yet</h3>
            <p>Create your first project to get started with organizing your tasks.</p>
            <t-button theme="primary" @click="showCreateDialog = true">
              Create Project
            </t-button>
          </div>
        </div>
        
        <div v-else class="projects-grid">
          <div v-for="project in projectStore.projects" :key="project.id" class="project-card">
            <t-card hover>
              <template #header>
                <div class="project-header">
                  <div class="project-title">
                    <t-tag v-if="project.color" :color="project.color" variant="light">
                      {{ project.title }}
                    </t-tag>
                    <span v-else>{{ project.title }}</span>
                  </div>
                  <t-dropdown>
                    <t-button variant="text" shape="square">
                      <MoreIcon />
                    </t-button>
                    <template #dropdown>
                      <t-dropdown-menu>
                        <t-dropdown-item @click="editProject(project)">
                          <EditIcon /> Edit
                        </t-dropdown-item>
                        <t-dropdown-item @click="deleteProject(project.id)">
                          <DeleteIcon /> Delete
                        </t-dropdown-item>
                      </t-dropdown-menu>
                    </template>
                  </t-dropdown>
                </div>
              </template>
              
              <p v-if="project.description">{{ project.description }}</p>
              <p v-else class="no-description">No description</p>
              
              <template #actions>
                <t-button variant="text" @click="$router.push(`/projects/${project.id}`)">
                  View Tasks
                </t-button>
              </template>
            </t-card>
          </div>
        </div>
      </div>
    </AppLayout>

    <!-- Create/Edit Project Dialog -->
    <t-dialog
      v-model:visible="showCreateDialog"
      :header="editingProject ? 'Edit Project' : 'Create Project'"
      @confirm="handleSaveProject"
      @cancel="resetForm"
    >
      <t-form ref="formRef" :model="formData" :rules="rules">
        <t-form-item label="Title" name="title">
          <t-input v-model="formData.title" placeholder="Enter project title" />
        </t-form-item>
        
        <t-form-item label="Description" name="description">
          <t-textarea v-model="formData.description" placeholder="Enter project description (optional)" />
        </t-form-item>
        
        <t-form-item label="Color" name="color">
          <t-color-picker v-model="formData.color" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { AddIcon, MoreIcon, EditIcon, DeleteIcon } from 'tdesign-icons-vue-next'
import { useProjectStore } from '@/stores/projects'
import AppLayout from '@/components/AppLayout.vue'
import type { Project, ProjectCreate } from '@/types'

const projectStore = useProjectStore()

const showCreateDialog = ref(false)
const editingProject = ref<Project | null>(null)
const formRef = ref()

const formData = reactive<ProjectCreate & { color: string }>({
  title: '',
  description: '',
  color: '#1976d2'
})

const rules = {
  title: [
    { required: true, message: 'Project title is required' }
  ]
}

const resetForm = () => {
  formData.title = ''
  formData.description = ''
  formData.color = '#1976d2'
  editingProject.value = null
  showCreateDialog.value = false
}

const editProject = (project: Project) => {
  editingProject.value = project
  formData.title = project.title
  formData.description = project.description || ''
  formData.color = project.color || '#1976d2'
  showCreateDialog.value = true
}

const handleSaveProject = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  try {
    if (editingProject.value) {
      await projectStore.updateProject(editingProject.value.id, {
        title: formData.title,
        description: formData.description || null,
        color: formData.color
      })
      MessagePlugin.success('Project updated successfully!')
    } else {
      await projectStore.createProject({
        title: formData.title,
        description: formData.description || null,
        color: formData.color
      })
      MessagePlugin.success('Project created successfully!')
    }
    resetForm()
  } catch (error) {
    MessagePlugin.error('Failed to save project')
  }
}

const deleteProject = async (id: number) => {
  try {
    await projectStore.deleteProject(id)
    MessagePlugin.success('Project deleted successfully!')
  } catch (error) {
    MessagePlugin.error('Failed to delete project')
  }
}

onMounted(async () => {
  await projectStore.fetchProjects()
})
</script>

<style scoped>
.projects-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.projects-content {
  padding: 20px;
}

.loading {
  padding: 40px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-content {
  text-align: center;
}

.empty-content h3 {
  margin-bottom: 8px;
  color: #333;
}

.empty-content p {
  margin-bottom: 20px;
  color: #666;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.project-card {
  height: 200px;
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.project-title {
  font-weight: 500;
}

.no-description {
  color: #999;
  font-style: italic;
}
</style>