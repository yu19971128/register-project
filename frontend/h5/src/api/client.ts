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
}

export interface RegistrationResult {
  order_no: string
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
  status: string
  created_at: string
  ticket_url: string
}

export interface TicketResult {
  order_no: string
  qrcode_data: string
  department: string
  doctor_name: string
  date: string
  start_time: string
  end_time: string
  patient_name: string
  patient_gender: string
  patient_age: number
  location: string
  status: string
  notice: string[]
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

export const scheduleApi = {
  list(date: string): Promise<{ total: number; list: Schedule[] }> {
    return request(`/schedules?date=${date}`)
  },
  get(id: number): Promise<Schedule> {
    return request(`/schedules/${id}`)
  },
}

export const registrationApi = {
  submit(schedule_id: number, patient_id: number, visitor_phone: string): Promise<RegistrationResult> {
    return request('/registrations', {
      method: 'POST',
      body: JSON.stringify({ schedule_id, patient_id, visitor_phone }),
    })
  },
  getTicket(order_no: string): Promise<TicketResult> {
    return request(`/registrations/ticket/${order_no}`)
  },
}
