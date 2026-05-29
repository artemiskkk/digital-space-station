import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export const login = (password: string) =>
  api.post<{ token: string }>('/api/auth/login', { password })

// Posts
export interface Post {
  id?: number
  title: string
  slug: string
  excerpt: string
  content: string
  cover_url: string
  tags: string[]
  status: 'draft' | 'published'
  read_time: string
  created_at?: string
  updated_at?: string
}

export const getPosts = (page = 1, size = 50) =>
  api.get('/api/admin/posts', { params: { page, size } })

export const createPost = (data: Post) =>
  api.post<Post>('/api/admin/posts', data)

export const updatePost = (id: number, data: Post) =>
  api.put<Post>(`/api/admin/posts/${id}`, data)

export const deletePost = (id: number) =>
  api.delete(`/api/admin/posts/${id}`)

// Moments
export interface Moment {
  id?: number
  text: string
  images: string[]
  mood: string
  created_at?: string
}

export const getMoments = (page = 1, size = 20) =>
  api.get('/api/moments', { params: { page, size } })

export const createMoment = (data: Moment) =>
  api.post<Moment>('/api/admin/moments', data)

export const deleteMoment = (id: number) =>
  api.delete(`/api/admin/moments/${id}`)

export default api
