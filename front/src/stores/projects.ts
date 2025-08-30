import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Project, ProjectCreate } from '@/types'
import apiService from '@/services/api'

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchProjects = async () => {
    loading.value = true
    error.value = null
    try {
      projects.value = await apiService.getProjects()
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to fetch projects'
      console.error('Error fetching projects:', err)
    } finally {
      loading.value = false
    }
  }

  const createProject = async (projectData: ProjectCreate) => {
    try {
      const newProject = await apiService.createProject(projectData)
      projects.value.push(newProject)
      return newProject
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to create project'
      console.error('Error creating project:', err)
      throw err
    }
  }

  const updateProject = async (id: number, projectData: Partial<ProjectCreate>) => {
    try {
      const updatedProject = await apiService.updateProject(id, projectData)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = updatedProject
      }
      return updatedProject
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to update project'
      console.error('Error updating project:', err)
      throw err
    }
  }

  const deleteProject = async (id: number) => {
    try {
      await apiService.deleteProject(id)
      projects.value = projects.value.filter(p => p.id !== id)
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to delete project'
      console.error('Error deleting project:', err)
      throw err
    }
  }

  const getProjectById = (id: number) => {
    return projects.value.find(p => p.id === id)
  }

  return {
    projects,
    loading,
    error,
    fetchProjects,
    createProject,
    updateProject,
    deleteProject,
    getProjectById
  }
})