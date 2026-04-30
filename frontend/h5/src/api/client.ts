const BASE = '/api/v1'

function getPhone(): string {
  return localStorage.getItem('visitor_phone') || ''
}

function headers(): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'X-Visitor-Phone': getPhone(),
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
  list(): Promise<ListResp> {
    return request('/patients')
  },
  get(id: number): Promise<Patient> {
    return request(`/patients/${id}`)
  },
  create(data: Omit<Patient, 'id'>): Promise<Patient> {
    return request('/patients', { method: 'POST', body: JSON.stringify(data) })
  },
  update(id: number, data: Partial<Patient>): Promise<void> {
    return request(`/patients/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },
  remove(id: number): Promise<void> {
    return request(`/patients/${id}`, { method: 'DELETE' })
  },
}
