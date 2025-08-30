export interface User {
  id: number
  username: string
  email: string
  name: string
  timezone?: string | null
  week_start: number
  language?: string | null
  is_active: boolean
  is_admin: boolean
  is_email_confirmed: boolean
  totp_enabled: boolean
  avatar_provider: string
  avatar_file_id?: number | null
  created_at: string
  updated_at: string
}

export interface UserLogin {
  username: string
  password: string
  long_token?: boolean
}

export interface UserCreate {
  username: string
  email: string
  name: string
  password: string
  timezone?: string | null
  week_start?: number
  language?: string | null
}

export interface Token {
  token: string
  type?: string
}

export interface Project {
  id: number
  title: string
  description?: string | null
  owner_id: number
  color?: string | null
  is_favorite: boolean
  position: number
  created_at: string
  updated_at: string
}

export interface ProjectCreate {
  title: string
  description?: string | null
  color?: string | null
}

export interface Task {
  id: number
  title: string
  description?: string | null
  project_id: number
  position: number
  done: boolean
  due_date?: string | null
  created_at: string
  updated_at: string
}

export interface TaskCreate {
  title: string
  description?: string | null
  project_id: number
  due_date?: string | null
}

export interface Label {
  id: number
  title: string
  description?: string | null
  hex_color: string
  created_by: number
  created_at: string
  updated_at: string
}

export interface LabelCreate {
  title: string
  description?: string | null
  hex_color: string
}

export interface Team {
  id: number
  name: string
  description?: string | null
  created_by: number
  created_at: string
  updated_at: string
}

export interface TeamCreate {
  name: string
  description?: string | null
}

export interface ApiInfo {
  version: string
  auth: {
    local: {
      enabled: boolean
      registration_enabled: boolean
    }
  }
  enabled_background_providers: string[]
  public_teams_enabled: boolean
}