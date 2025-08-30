import axios, { AxiosInstance, AxiosResponse } from 'axios'
import type { 
  User, UserLogin, UserCreate, Token, 
  Project, ProjectCreate, 
  Task, TaskCreate, 
  Label, LabelCreate, 
  Team, TeamCreate, 
  ApiInfo 
} from '@/types'

class ApiService {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: '/api/v1',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Add auth token to requests if available
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('vikunja_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Handle auth errors
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('vikunja_token')
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  // Authentication
  async login(credentials: UserLogin): Promise<Token> {
    const response: AxiosResponse<Token> = await this.client.post('/login', credentials)
    return response.data
  }

  async register(userData: UserCreate): Promise<User> {
    const response: AxiosResponse<User> = await this.client.post('/register', userData)
    return response.data
  }

  async testToken(): Promise<{ message: string; user: string }> {
    const response = await this.client.get('/token/test')
    return response.data
  }

  // System Info
  async getInfo(): Promise<ApiInfo> {
    const response: AxiosResponse<ApiInfo> = await this.client.get('/info')
    return response.data
  }

  // Projects
  async getProjects(): Promise<Project[]> {
    const response: AxiosResponse<Project[]> = await this.client.get('/projects')
    return response.data
  }

  async getProject(id: number): Promise<Project> {
    const response: AxiosResponse<Project> = await this.client.get(`/projects/${id}`)
    return response.data
  }

  async createProject(project: ProjectCreate): Promise<Project> {
    const response: AxiosResponse<Project> = await this.client.post('/projects', project)
    return response.data
  }

  async updateProject(id: number, project: Partial<ProjectCreate>): Promise<Project> {
    const response: AxiosResponse<Project> = await this.client.put(`/projects/${id}`, project)
    return response.data
  }

  async deleteProject(id: number): Promise<void> {
    await this.client.delete(`/projects/${id}`)
  }

  // Tasks
  async getProjectTasks(projectId: number): Promise<Task[]> {
    const response: AxiosResponse<Task[]> = await this.client.get(`/projects/${projectId}/tasks`)
    return response.data
  }

  async getTask(taskId: number): Promise<Task> {
    const response: AxiosResponse<Task> = await this.client.get(`/tasks/${taskId}`)
    return response.data
  }

  async createTask(task: TaskCreate): Promise<Task> {
    const response: AxiosResponse<Task> = await this.client.post(`/projects/${task.project_id}/tasks`, task)
    return response.data
  }

  async updateTask(id: number, task: Partial<TaskCreate>): Promise<Task> {
    const response: AxiosResponse<Task> = await this.client.put(`/tasks/${id}`, task)
    return response.data
  }

  async deleteTask(id: number): Promise<void> {
    await this.client.delete(`/tasks/${id}`)
  }

  // Labels
  async getLabels(): Promise<Label[]> {
    const response: AxiosResponse<Label[]> = await this.client.get('/labels')
    return response.data
  }

  async createLabel(label: LabelCreate): Promise<Label> {
    const response: AxiosResponse<Label> = await this.client.post('/labels', label)
    return response.data
  }

  async updateLabel(id: number, label: Partial<LabelCreate>): Promise<Label> {
    const response: AxiosResponse<Label> = await this.client.put(`/labels/${id}`, label)
    return response.data
  }

  async deleteLabel(id: number): Promise<void> {
    await this.client.delete(`/labels/${id}`)
  }

  // Teams
  async getTeams(): Promise<Team[]> {
    const response: AxiosResponse<Team[]> = await this.client.get('/teams')
    return response.data
  }

  async createTeam(team: TeamCreate): Promise<Team> {
    const response: AxiosResponse<Team> = await this.client.post('/teams', team)
    return response.data
  }

  async updateTeam(id: number, team: Partial<TeamCreate>): Promise<Team> {
    const response: AxiosResponse<Team> = await this.client.put(`/teams/${id}`, team)
    return response.data
  }

  async deleteTeam(id: number): Promise<void> {
    await this.client.delete(`/teams/${id}`)
  }
}

export const apiService = new ApiService()
export default apiService