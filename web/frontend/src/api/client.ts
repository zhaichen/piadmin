const API_BASE = '/api'

function getToken(): string | null {
  return localStorage.getItem('piadmin_token')
}

export function setToken(token: string) {
  localStorage.setItem('piadmin_token', token)
}

export function clearToken() {
  localStorage.removeItem('piadmin_token')
}

export function hasToken(): boolean {
  return !!getToken()
}

export function getAuthToken(): string {
  return getToken() || ''
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

  if (res.status === 401) {
    clearToken()
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }

  return res.json()
}

async function requestRaw(path: string, options: RequestInit = {}): Promise<Response> {
  const token = getToken()
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return fetch(`${API_BASE}${path}`, { ...options, headers })
}

// Auth
export async function login(password: string): Promise<string> {
  const data = await request<{ token: string }>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ password }),
  })
  setToken(data.token)
  return data.token
}

// System
export async function getSnapshot() {
  return request<import('@/types/system').SystemSnapshot>('/system/snapshot')
}

// Processes
export async function getProcesses() {
  return request<import('@/types/process').ProcessInfo[]>('/processes')
}

export async function killProcess(pid: number, force: boolean = false) {
  return request<{ status: string }>(`/processes?pid=${pid}&force=${force}`, { method: 'DELETE' })
}

// Services
export async function getServices() {
  return request<import('@/types/service').ServiceInfo[]>('/services')
}

export async function getServiceStatus(name: string) {
  return request<{ status: string }>(`/services/${name}`)
}

export async function serviceAction(name: string, action: string) {
  return request<{ status: string }>(`/services/${name}`, {
    method: 'POST',
    body: JSON.stringify({ action }),
  })
}

// Network
export async function getNetworkInterfaces() {
  return request<import('@/types/network').InterfaceInfo[]>('/network/interfaces')
}

// Files
export async function listFiles(path: string) {
  return request<import('@/types/file').FileEntry[]>(`/files?path=${encodeURIComponent(path)}`)
}

export async function downloadFile(path: string) {
  const res = await requestRaw(`/files/download?path=${encodeURIComponent(path)}`)
  return res.blob()
}

export async function uploadFile(dir: string, file: File) {
  const formData = new FormData()
  formData.append('path', dir)
  formData.append('file', file)
  const token = getToken()
  return fetch(`${API_BASE}/files/upload`, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: formData,
  })
}

export async function deleteFile(path: string) {
  return request<{ status: string }>('/files', {
    method: 'DELETE',
    body: JSON.stringify({ path }),
  })
}

export async function mkdirFile(path: string) {
  return request<{ status: string }>('/files/mkdir', {
    method: 'POST',
    body: JSON.stringify({ path }),
  })
}

export async function renameFile(oldPath: string, newPath: string) {
  return request<{ status: string }>('/files/rename', {
    method: 'PUT',
    body: JSON.stringify({ old_path: oldPath, new_path: newPath }),
  })
}

// GPIO
export async function gpioAvailable() {
  return request<{ available: boolean }>('/gpio/available')
}

export async function getGpioPins() {
  return request<import('@/types/gpio').PinInfo[]>('/gpio/pins')
}

export async function setGpioPin(pin: number, direction: string, value: number) {
  return request<{ status: string }>('/gpio/pins', {
    method: 'POST',
    body: JSON.stringify({ pin, direction, value }),
  })
}

export async function exportGpioPin(pin: number) {
  return request<{ status: string }>('/gpio/export', {
    method: 'POST',
    body: JSON.stringify({ pin }),
  })
}

export async function unexportGpioPin(pin: number) {
  return request<{ status: string }>('/gpio/unexport', {
    method: 'POST',
    body: JSON.stringify({ pin }),
  })
}
