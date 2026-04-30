const BASE = '/api/v1/admin'

function getToken(): string {
  return localStorage.getItem('admin_token') || ''
}

function headers(): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${getToken()}`,
  }
}

export interface Patient {
  id: number
  name: string
  id_card: string
  phone: string
  gender?: string
  age?: number
  address?: string
  created_at?: string
  updated_at?: string
}

export interface ListResp {
  total: number
  list: Patient[]
}

export interface ApiResp<T> {
  code: number
  data: T
  message: string
}

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + url, {
    ...init,
    headers: { ...headers(), ...(init?.headers || {}) },
  })
  const json: ApiResp<T> = await res.json()
  if (json.code !== 200) {
    throw new Error(json.message || '请求失败')
  }
  return json.data
}

export const patientApi = {
  list(keyword?: string, page = 1, pageSize = 10): Promise<ListResp> {
    const qs = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
    if (keyword) qs.set('keyword', keyword)
    return request(`/patients?${qs.toString()}`)
  },
  get(id: number): Promise<Patient> {
    return request(`/patients/${id}`)
  },
  update(id: number, data: Partial<Patient>): Promise<void> {
    return request(`/patients/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },
  remove(id: number): Promise<void> {
    return request(`/patients/${id}`, { method: 'DELETE' })
  },
}
