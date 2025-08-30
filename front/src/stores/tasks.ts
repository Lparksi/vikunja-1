import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task, TaskCreate } from '@/types'
import apiService from '@/services/api'

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchProjectTasks = async (projectId: number) => {
    loading.value = true
    error.value = null
    try {
      const projectTasks = await apiService.getProjectTasks(projectId)
      // Replace tasks for this project
      tasks.value = tasks.value.filter(t => t.project_id !== projectId)
      tasks.value.push(...projectTasks)
      return projectTasks
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to fetch tasks'
      console.error('Error fetching tasks:', err)
    } finally {
      loading.value = false
    }
  }

  const createTask = async (taskData: TaskCreate) => {
    try {
      const newTask = await apiService.createTask(taskData)
      tasks.value.push(newTask)
      return newTask
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to create task'
      console.error('Error creating task:', err)
      throw err
    }
  }

  const updateTask = async (id: number, taskData: Partial<TaskCreate>) => {
    try {
      const updatedTask = await apiService.updateTask(id, taskData)
      const index = tasks.value.findIndex(t => t.id === id)
      if (index !== -1) {
        tasks.value[index] = updatedTask
      }
      return updatedTask
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to update task'
      console.error('Error updating task:', err)
      throw err
    }
  }

  const deleteTask = async (id: number) => {
    try {
      await apiService.deleteTask(id)
      tasks.value = tasks.value.filter(t => t.id !== id)
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to delete task'
      console.error('Error deleting task:', err)
      throw err
    }
  }

  const getTasksByProject = (projectId: number) => {
    return tasks.value.filter(t => t.project_id === projectId)
  }

  const getTaskById = (id: number) => {
    return tasks.value.find(t => t.id === id)
  }

  const toggleTaskDone = async (id: number) => {
    const task = getTaskById(id)
    if (task) {
      await updateTask(id, { done: !task.done } as any)
    }
  }

  return {
    tasks,
    loading,
    error,
    fetchProjectTasks,
    createTask,
    updateTask,
    deleteTask,
    getTasksByProject,
    getTaskById,
    toggleTaskDone
  }
})