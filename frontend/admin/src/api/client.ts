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

export interface Schedule {
  id: number
  date: string
  department: string
  doctor_name: string
  start_time: string
  end_time: string
  total_quota: number
  remaining: number
  status: string
  created_at?: string
  updated_at?: string
}

export interface ListResp {
  total: number
  list: any[]
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
  list(keyword?: string, page = 1, pageSize = 10): Promise<{ total: number; list: Patient[] }> {
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

export interface Order {
  id: number
  order_no: string
  patient_name: string
  department: string
  doctor_name: string
  date: string
  start_time: string
  end_time: string
  status: string
  created_at: string
}

export interface OrderDetail {
  id: number
  order_no: string
  status: string
  schedule: {
    id: number
    department: string
    doctor_name: string
    date: string
    start_time: string
    end_time: string
  }
  patient: {
    id: number
    name: string
    gender: string
    age: number
  }
  visitor_phone: string
  created_at: string
  updated_at: string
}

export const scheduleApi = {
  list(params: { date: string; department?: string; doctor_name?: string; page?: number; pageSize?: number }): Promise<{ total: number; list: Schedule[] }> {
    const qs = new URLSearchParams({
      page: String(params.page || 1),
      page_size: String(params.pageSize || 10),
      date: params.date,
    })
    if (params.department) qs.set('department', params.department)
    if (params.doctor_name) qs.set('doctor_name', params.doctor_name)
    return request(`/schedules?${qs.toString()}`)
  },
  get(id: number): Promise<Schedule> {
    return request(`/schedules/${id}`)
  },
  create(data: Partial<Schedule>): Promise<Schedule> {
    return request(`/schedules`, { method: 'POST', body: JSON.stringify(data) })
  },
  update(id: number, data: Partial<Schedule>): Promise<void> {
    return request(`/schedules/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },
  remove(id: number): Promise<void> {
    return request(`/schedules/${id}`, { method: 'DELETE' })
  },
}

export const orderApi = {
  list(params?: { keyword?: string; status?: string; date?: string; department?: string; page?: number; pageSize?: number }): Promise<{ total: number; list: Order[] }> {
    const qs = new URLSearchParams()
    if (params?.keyword) qs.set('keyword', params.keyword)
    if (params?.status) qs.set('status', params.status)
    if (params?.date) qs.set('date', params.date)
    if (params?.department) qs.set('department', params.department)
    qs.set('page', String(params?.page || 1))
    qs.set('page_size', String(params?.pageSize || 10))
    return request(`/orders?${qs.toString()}`)
  },
  get(id: number): Promise<OrderDetail> {
    return request(`/orders/${id}`)
  },
  cancel(id: number, reason?: string): Promise<{ id: number; order_no: string; status: string }> {
    return request(`/orders/${id}/cancel`, { method: 'PUT', body: JSON.stringify({ reason }) })
  },
  change(id: number, new_schedule_id: number): Promise<OrderDetail> {
    return request(`/orders/${id}/change`, { method: 'PUT', body: JSON.stringify({ new_schedule_id }) })
  },
  complete(id: number): Promise<{ id: number; order_no: string; status: string }> {
    return request(`/orders/${id}/complete`, { method: 'PUT' })
  },
}
