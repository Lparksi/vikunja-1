<template>
  <div class="project-detail">
    <AppLayout>
      <template #header>
        <div class="project-header" v-if="project">
          <div class="project-info">
            <t-tag v-if="project.color" :color="project.color" variant="light" size="large">
              {{ project.title }}
            </t-tag>
            <h1 v-else>{{ project.title }}</h1>
            <p v-if="project.description">{{ project.description }}</p>
          </div>
          <t-button theme="primary" @click="showCreateTaskDialog = true">
            <template #icon><AddIcon /></template>
            New Task
          </t-button>
        </div>
      </template>
      
      <div class="project-content">
        <div v-if="loading" class="loading">
          <t-skeleton />
        </div>
        
        <div v-else-if="!project" class="error">
          <h3>Project not found</h3>
          <t-button @click="$router.push('/projects')">Back to Projects</t-button>
        </div>
        
        <div v-else class="tasks-section">
          <div v-if="tasks.length === 0" class="empty-tasks">
            <h3>No tasks yet</h3>
            <p>Create your first task to get started.</p>
            <t-button theme="primary" @click="showCreateTaskDialog = true">
              Create Task
            </t-button>
          </div>
          
          <div v-else class="tasks-list">
            <div v-for="task in tasks" :key="task.id" class="task-item">
              <t-card>
                <div class="task-content">
                  <div class="task-main">
                    <t-checkbox
                      :checked="task.done"
                      @change="toggleTask(task.id)"
                    />
                    <div class="task-info">
                      <h4 :class="{ completed: task.done }">{{ task.title }}</h4>
                      <p v-if="task.description" class="task-description">{{ task.description }}</p>
                      <div v-if="task.due_date" class="task-due-date">
                        <t-tag variant="outline" size="small">
                          Due: {{ formatDate(task.due_date) }}
                        </t-tag>
                      </div>
                    </div>
                  </div>
                  <div class="task-actions">
                    <t-button variant="text" @click="editTask(task)">
                      <EditIcon />
                    </t-button>
                    <t-button variant="text" @click="deleteTask(task.id)">
                      <DeleteIcon />
                    </t-button>
                  </div>
                </div>
              </t-card>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>

    <!-- Create/Edit Task Dialog -->
    <t-dialog
      v-model:visible="showCreateTaskDialog"
      :header="editingTask ? 'Edit Task' : 'Create Task'"
      @confirm="handleSaveTask"
      @cancel="resetTaskForm"
      width="600px"
    >
      <t-form ref="taskFormRef" :model="taskFormData" :rules="taskRules">
        <t-form-item label="Title" name="title">
          <t-input v-model="taskFormData.title" placeholder="Enter task title" />
        </t-form-item>
        
        <t-form-item label="Description" name="description">
          <t-textarea v-model="taskFormData.description" placeholder="Enter task description (optional)" />
        </t-form-item>
        
        <t-form-item label="Due Date" name="due_date">
          <t-date-picker
            v-model="taskFormData.due_date"
            clearable
            placeholder="Select due date (optional)"
          />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { AddIcon, EditIcon, DeleteIcon } from 'tdesign-icons-vue-next'
import { useProjectStore } from '@/stores/projects'
import { useTaskStore } from '@/stores/tasks'
import AppLayout from '@/components/AppLayout.vue'
import type { Task, TaskCreate } from '@/types'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const loading = ref(true)
const showCreateTaskDialog = ref(false)
const editingTask = ref<Task | null>(null)
const taskFormRef = ref()

const projectId = computed(() => parseInt(route.params.id as string))

const project = computed(() => {
  return projectStore.getProjectById(projectId.value)
})

const tasks = computed(() => {
  return taskStore.getTasksByProject(projectId.value)
})

const taskFormData = reactive<TaskCreate & { due_date: Date | null }>({
  title: '',
  description: '',
  project_id: projectId.value,
  due_date: null
})

const taskRules = {
  title: [
    { required: true, message: 'Task title is required' }
  ]
}

const resetTaskForm = () => {
  taskFormData.title = ''
  taskFormData.description = ''
  taskFormData.due_date = null
  editingTask.value = null
  showCreateTaskDialog.value = false
}

const editTask = (task: Task) => {
  editingTask.value = task
  taskFormData.title = task.title
  taskFormData.description = task.description || ''
  taskFormData.due_date = task.due_date ? new Date(task.due_date) : null
  showCreateTaskDialog.value = true
}

const handleSaveTask = async () => {
  const valid = await taskFormRef.value?.validate()
  if (!valid) return

  try {
    const taskData = {
      title: taskFormData.title,
      description: taskFormData.description || null,
      project_id: projectId.value,
      due_date: taskFormData.due_date ? taskFormData.due_date.toISOString() : null
    }

    if (editingTask.value) {
      await taskStore.updateTask(editingTask.value.id, taskData as any)
      MessagePlugin.success('Task updated successfully!')
    } else {
      await taskStore.createTask(taskData as TaskCreate)
      MessagePlugin.success('Task created successfully!')
    }
    resetTaskForm()
  } catch (error) {
    MessagePlugin.error('Failed to save task')
  }
}

const toggleTask = async (taskId: number) => {
  try {
    await taskStore.toggleTaskDone(taskId)
  } catch (error) {
    MessagePlugin.error('Failed to update task')
  }
}

const deleteTask = async (taskId: number) => {
  try {
    await taskStore.deleteTask(taskId)
    MessagePlugin.success('Task deleted successfully!')
  } catch (error) {
    MessagePlugin.error('Failed to delete task')
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

onMounted(async () => {
  loading.value = true
  try {
    // Ensure projects are loaded
    if (projectStore.projects.length === 0) {
      await projectStore.fetchProjects()
    }
    
    // Load tasks for this project
    await taskStore.fetchProjectTasks(projectId.value)
    
    // If project doesn't exist, redirect
    if (!project.value) {
      MessagePlugin.error('Project not found')
      router.push('/projects')
    }
  } catch (error) {
    MessagePlugin.error('Failed to load project data')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.project-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
}

.project-info h1 {
  margin: 0 0 8px 0;
}

.project-info p {
  margin: 0;
  color: #666;
}

.project-content {
  padding: 20px;
}

.loading, .error {
  padding: 40px;
  text-align: center;
}

.empty-tasks {
  text-align: center;
  padding: 60px 20px;
}

.empty-tasks h3 {
  margin-bottom: 8px;
  color: #333;
}

.empty-tasks p {
  margin-bottom: 20px;
  color: #666;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.task-item {
  transition: all 0.2s ease;
}

.task-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.task-main {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  flex: 1;
}

.task-info {
  flex: 1;
}

.task-info h4 {
  margin: 0 0 8px 0;
  transition: all 0.2s ease;
}

.task-info h4.completed {
  text-decoration: line-through;
  color: #999;
}

.task-description {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
}

.task-due-date {
  margin-top: 8px;
}

.task-actions {
  display: flex;
  gap: 8px;
}
</style>